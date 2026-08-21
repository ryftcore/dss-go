package fc

import (
	"bufio"
	"encoding/json"
	"math/big"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// The fc direct-check KAT: testdata/oracle/fc_direct.jsonl holds the XmlFC
// upstream produces when each fc check the top-level chains do not exercise
// under the default policy is run alone at Level.FAIL - over every signature of
// the marshal-parity corpus, plus a few synthetic literals. See
// ../testdata/gen/FcSavDirectOracle.java. Each row carries a "check" label; this
// test rebuilds the same invocation in Go and compares.

type fcDirectRow struct {
	fcOracleRow
	Check string `json:"check"`
}

// singleFCChain is a chain of exactly one item, mirroring the oracle's
// SingleFCChain.
type singleFCChain struct {
	*process.ChainBase[*jaxb.XmlFC]
	factory func(result *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC]
}

func newSingleFCChain(i18nProvider *i18n.I18nProvider,
	factory func(result *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC],
) *singleFCChain {
	xmlFC := &jaxb.XmlFC{}
	c := &singleFCChain{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlFC,
			&xmlFC.XmlConstraintsConclusionContent, &xmlFC.XmlConstraintsConclusionAttrs)),
		factory: factory,
	}
	c.InitChainBase(c)
	return c
}

func (c *singleFCChain) InitChain() {
	c.FirstItem = c.factory(c.Result, process.GetLevelRule(enumerations.Level_FAIL))
}

// multiValuesRule is the oracle's ANY / NONE MultiValuesRule at Level.FAIL.
type fcMultiValuesRule struct{ values []string }

func (r fcMultiValuesRule) Level() enumerations.Level { return enumerations.Level_FAIL }
func (r fcMultiValuesRule) Values() []string          { return r.values }

var (
	fcAnyRule  = fcMultiValuesRule{values: []string{"*"}}
	fcNoneRule = fcMultiValuesRule{values: []string{"NO-SUCH-VALUE"}}
)

func loadFCDirectRowsWithCheck(t *testing.T, path string) []*fcDirectRow {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer file.Close()

	var rows []*fcDirectRow
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		row := &fcDirectRow{}
		if err := json.Unmarshal(scanner.Bytes(), row); err != nil {
			t.Fatalf("parse oracle row: %v", err)
		}
		rows = append(rows, row)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("read oracle: %v", err)
	}
	return rows
}

// checkBaseName strips the "-any"/"-none" rule variant suffix from a check label.
func checkBaseName(check string) string {
	if i := strings.LastIndex(check, "-"); i > 0 {
		switch check[i+1:] {
		case "any", "none", "match":
			return check[:i]
		}
	}
	return check
}

// corpusFileNames returns the corpus dumps the direct oracle recorded rows for,
// in first-seen order.
func corpusFileNames(rows []*fcDirectRow) []string {
	var files []string
	seen := map[string]bool{}
	for _, row := range rows {
		if row.File == "synthetic" || seen[row.File] {
			continue
		}
		seen[row.File] = true
		files = append(files, row.File)
	}
	return files
}

// fcKnownDeviation lists the (file, token, check) rows the port cannot reproduce, with the
// reason. The set is asserted to be exactly this one, so a new failure cannot hide in it.
// It is empty: the four model-*.xml schema-coverage fixtures, the only inputs on which the
// two runtimes disagreed, are excluded from this corpus altogether (see the package README -
// their IDREF graph is dangling and their attribute values carry raw control characters, so
// upstream and the port are not reading the same graph). The mechanism stays so that a future
// deviation has to be recorded rather than silently tolerated.
var fcKnownDeviation = map[string]string{}

func TestFCDirectChecksAgainstJavaOracle(t *testing.T) {
	rows := loadFCDirectRowsWithCheck(t, corpustest.Path(t, "oracle/fc_direct.jsonl"))
	if len(rows) == 0 {
		t.Fatal("empty direct oracle")
	}
	i18nProvider := i18n.NewI18nProvider()

	// index by (file, token, check)
	type key struct{ file, token, check string }
	index := map[key]*fcDirectRow{}
	for _, row := range rows {
		index[key{row.File, row.Token, row.Check}] = row
	}

	statuses := map[string]map[string]int{}
	seenFCDeviation := map[string]bool{}
	matched := 0

	run := func(file, token, check string,
		factory func(result *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC]) {
		want, ok := index[key{file, token, check}]
		if !ok {
			return // Java threw here and recorded no row
		}
		if want.Title == nil {
			// A single-item chain defines no title MessageTag: Java leaves the
			// attribute null, which the generated non-pointer Go member spells "".
			empty := ""
			want.Title = &empty
		}
		var result *jaxb.XmlFC
		if safeExecute(func() { result = newSingleFCChain(i18nProvider, factory).Execute() }) {
			t.Errorf("%s/%s/%s: Go panicked where Java produced a row", file, token, check)
			return
		}
		got := toFCRow(file, token, enumerations.Context(want.Context), "FC",
			&result.XmlConstraintsConclusionContent, result.Title)
		matched++
		base := checkBaseName(check)
		if statuses[base] == nil {
			statuses[base] = map[string]int{}
		}
		for _, constraint := range want.Constraints {
			if constraint.Status != nil {
				statuses[base][*constraint.Status]++
			}
		}
		key := file + "/" + token + "/" + check
		if reason, known := fcKnownDeviation[key]; known {
			seenFCDeviation[key] = true
			if reflect.DeepEqual(&want.fcOracleRow, got) {
				t.Errorf("%s: listed in fcKnownDeviation but now matches - drop the entry", key)
			} else {
				t.Logf("%s: known deviation, not compared (%s)", key, reason)
			}
			return
		}
		if !reflect.DeepEqual(&want.fcOracleRow, got) {
			t.Errorf("%s/%s/%s: mismatch\nexpected: %s\nactual:   %s",
				file, token, check, mustFCJSON(t, &want.fcOracleRow), mustFCJSON(t, got))
		}
	}

	for _, name := range corpusFileNames(rows) {
		if strings.HasPrefix(name, "model-") {
			continue // schema-coverage fixture, excluded from the direct corpus (see README)
		}
		diagnosticData := loadFCDiagnosticData(t, name)
		for _, signature := range diagnosticData.Signatures() {
			sig := signature
			id := sig.Id()
			run(name, id, "EllipticCurveKeySizeCheck",
				func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
					return NewEllipticCurveKeySizeCheck(i18nProvider, r, sig, rule)
				})
			run(name, id, "FullScopeCheck",
				func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
					return NewFullScopeCheck(i18nProvider, r, sig.SignatureScopes(), rule)
				})
			run(name, id, "ByteRangeAllDocumentCheck",
				func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
					return NewByteRangeAllDocumentCheck(i18nProvider, r, diagnosticData, rule)
				})
			run(name, id, "ReferencesNotAmbiguousCheck",
				func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
					return NewReferencesNotAmbiguousCheck(i18nProvider, r, sig, rule)
				})
			run(name, id, "SignatureNotAmbiguousCheck",
				func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
					return NewSignatureNotAmbiguousCheck(i18nProvider, r, sig, rule)
				})
			run(name, id, "SignerInformationStoreCheck",
				func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
					return NewSignerInformationStoreCheck(i18nProvider, r, sig, rule)
				})
			run(name, id, "FormatCheck-any",
				func(r *process.Result[*jaxb.XmlFC], _ policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
					return NewFormatCheck(i18nProvider, r, sig, fcAnyRule)
				})
			run(name, id, "FormatCheck-none",
				func(r *process.Result[*jaxb.XmlFC], _ policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
					return NewFormatCheck(i18nProvider, r, sig, fcNoneRule)
				})
			run(name, id, "ByteRangeCollisionCheck",
				func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
					return NewByteRangeCollisionCheck(i18nProvider, r, sig, diagnosticData, rule)
				})
			if containerInfo := diagnosticData.ContainerInfo(); containerInfo != nil {
				run(name, id, "SignedFilesPresentCheck",
					func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewSignedFilesPresentCheck(i18nProvider, r, containerInfo, rule)
					})
				run(name, id, "AllFilesSignedCheck",
					func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewAllFilesSignedCheck(i18nProvider, r, sig, containerInfo, rule)
					})
				run(name, id, "ManifestFilePresentCheck",
					func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewManifestFilePresentCheck(i18nProvider, r, containerInfo, rule)
					})
				containerType := diagnosticData.ContainerType()
				run(name, id, "ContainerTypeCheck-any",
					func(r *process.Result[*jaxb.XmlFC], _ policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewContainerTypeCheck(i18nProvider, r, containerType, fcAnyRule)
					})
				run(name, id, "ContainerTypeCheck-none",
					func(r *process.Result[*jaxb.XmlFC], _ policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewContainerTypeCheck(i18nProvider, r, containerType, fcNoneRule)
					})
				var mimeTypeContent string
				if containerInfo.MimeTypeContent != nil {
					mimeTypeContent = *containerInfo.MimeTypeContent
				}
				run(name, id, "AcceptableMimetypeFileContentCheck-any",
					func(r *process.Result[*jaxb.XmlFC], _ policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewAcceptableMimetypeFileContentCheck(i18nProvider, r, mimeTypeContent, fcAnyRule)
					})
				run(name, id, "AcceptableMimetypeFileContentCheck-none",
					func(r *process.Result[*jaxb.XmlFC], _ policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewAcceptableMimetypeFileContentCheck(i18nProvider, r, mimeTypeContent, fcNoneRule)
					})
			}
			if rev := sig.PDFRevision(); rev != nil {
				run(name, id, "PdfSignatureDictionaryCheck",
					func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewPdfSignatureDictionaryCheck(i18nProvider, r, rev, rule)
					})
				run(name, id, "PdfAnnotationOverlapCheck",
					func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewPdfAnnotationOverlapCheck(i18nProvider, r, rev, rule)
					})
				run(name, id, "UndefinedChangesCheck",
					func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewUndefinedChangesCheck(i18nProvider, r, rev, rule)
					})
				run(name, id, "SigFieldLockCheck",
					func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewSigFieldLockCheck(i18nProvider, r, rev, rule)
					})
				run(name, id, "FieldMDPCheck",
					func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewFieldMDPCheck(i18nProvider, r, rev, rule)
					})
				run(name, id, "PdfVisualDifferenceCheck",
					func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewPdfVisualDifferenceCheck(i18nProvider, r, rev, rule)
					})
				run(name, id, "AnnotationChangesCheck",
					func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewAnnotationChangesCheck(i18nProvider, r, rev, rule)
					})
				run(name, id, "DocMDPCheck",
					func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewDocMDPCheck(i18nProvider, r, rev, rule)
					})
				run(name, id, "FormFillChangesCheck",
					func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
						return NewFormFillChangesCheck(i18nProvider, r, rev, rule)
					})
			}
		}
		for _, timestamp := range diagnosticData.TimestampList() {
			tst := timestamp
			run(name, tst.Id(), "CAdESV3HashIndexCheck",
				func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
					return NewCAdESV3HashIndexCheck(i18nProvider, r, tst, rule)
				})
		}
	}

	// --- synthetic literals
	for _, zip := range []struct{ label, value string }{
		{"null", ""}, {"empty", ""}, {"asice", "mimetype=application/vnd.etsi.asic-e+zip"},
	} {
		comment := zip.value
		run("synthetic", "zip-comment-"+zip.label, "ZipCommentPresentCheck",
			func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
				return NewZipCommentPresentCheck(i18nProvider, r, comment, rule)
			})
		run("synthetic", "zip-comment-"+zip.label+"-any", "AcceptableZipCommentCheck",
			func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
				return NewAcceptableZipCommentCheck(i18nProvider, r, comment, fcAnyRule)
			})
		run("synthetic", "zip-comment-"+zip.label+"-none", "AcceptableZipCommentCheck",
			func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
				return NewAcceptableZipCommentCheck(i18nProvider, r, comment, fcNoneRule)
			})
	}
	for _, present := range []bool{true, false} {
		value := present
		label := "mimetype-present-false"
		if present {
			label = "mimetype-present-true"
		}
		run("synthetic", label, "MimeTypeFilePresentCheck",
			func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
				return NewMimeTypeFilePresentCheck(i18nProvider, r, value, rule)
			})
	}
	for _, compliant := range []bool{true, false} {
		value := compliant
		label := "pdfa-compliant-false"
		if compliant {
			label = "pdfa-compliant-true"
		}
		run("synthetic", label, "PDFAComplianceCheck",
			func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
				return NewPDFAComplianceCheck(i18nProvider, r, value, rule)
			})
	}
	for _, profile := range []struct{ label, value string }{{"null", ""}, {"PDF/A-2A", "PDF/A-2A"}} {
		value := profile.value
		run("synthetic", "pdfa-profile-"+profile.label+"-any", "PDFAProfileCheck",
			func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
				return NewPDFAProfileCheck(i18nProvider, r, value, fcAnyRule)
			})
		run("synthetic", "pdfa-profile-"+profile.label+"-none", "PDFAProfileCheck",
			func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
				return NewPDFAProfileCheck(i18nProvider, r, value, fcNoneRule)
			})
	}
	emptyContainer := &diagjaxb.XmlContainerInfo{}
	containerType := diagjaxb.ASiCContainerTypeValue(enumerations.ASiCContainerType_ASiC_E)
	emptyContainer.ContainerType = &containerType
	run("synthetic", "no-content-files", "SignedFilesPresentCheck",
		func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
			return NewSignedFilesPresentCheck(i18nProvider, r, emptyContainer, rule)
		})
	collidingData := collidingByteRangeDiagnosticData()
	run("synthetic", "byte-range-collision", "ByteRangeCollisionCheck",
		func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
			return NewByteRangeCollisionCheck(i18nProvider, r, collidingData.Signatures()[0], collidingData, rule)
		})
	badByteRangeData := failingByteRangeDiagnosticData()
	run("synthetic", "byte-range-invalid", "ByteRangeAllDocumentCheck",
		func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
			return NewByteRangeAllDocumentCheck(i18nProvider, r, badByteRangeData, rule)
		})
	run("synthetic", "duplicated-reference", "ReferencesNotAmbiguousCheck",
		func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
			return NewReferencesNotAmbiguousCheck(i18nProvider, r, signatureWithDuplicatedReference(), rule)
		})
	run("synthetic", "duplicated-signature", "SignatureNotAmbiguousCheck",
		func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
			return NewSignatureNotAmbiguousCheck(i18nProvider, r, signatureWithDuplicatedSignature(), rule)
		})
	failingRevision := failingPdfRevision()
	run("synthetic", "failing-revision", "DocMDPCheck",
		func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
			return NewDocMDPCheck(i18nProvider, r, failingRevision, rule)
		})
	run("synthetic", "failing-revision", "SigFieldLockCheck",
		func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
			return NewSigFieldLockCheck(i18nProvider, r, failingRevision, rule)
		})
	run("synthetic", "failing-revision", "FieldMDPCheck",
		func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
			return NewFieldMDPCheck(i18nProvider, r, failingRevision, rule)
		})
	run("synthetic", "failing-revision", "PdfSignatureDictionaryCheck",
		func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
			return NewPdfSignatureDictionaryCheck(i18nProvider, r, failingRevision, rule)
		})
	run("synthetic", "failing-revision", "PdfAnnotationOverlapCheck",
		func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
			return NewPdfAnnotationOverlapCheck(i18nProvider, r, failingRevision, rule)
		})
	run("synthetic", "ecdsa-mismatch", "EllipticCurveKeySizeCheck",
		func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
			return NewEllipticCurveKeySizeCheck(i18nProvider, r, signatureWithEcdsaKeySizeMismatch(), rule)
		})

	if matched != len(rows) {
		t.Errorf("matched %d of %d direct oracle rows", matched, len(rows))
	}
	for key := range fcKnownDeviation {
		if !seenFCDeviation[key] {
			t.Errorf("fcKnownDeviation lists %q, but that row never deviated", key)
		}
	}
	for check, seen := range statuses {
		if seen["OK"] == 0 || seen["NOT OK"] == 0 {
			t.Errorf("%s: the direct corpus only produced %v - both a happy and a failure row are required",
				check, seen)
		}
	}
	t.Logf("fc direct checks exercised: %d rows over %d check classes", matched, len(statuses))
}

// signatureWithEcdsaKeySizeMismatch mirrors the oracle's synthetic ECDSA
// signature whose key size does not match its digest algorithm.
func signatureWithEcdsaKeySizeMismatch() *diagnostic.SignatureWrapper {
	xml := &diagjaxb.XmlSignature{}
	id := diagjaxb.CollapsedString("S-SYNTHETIC")
	xml.Id = &id
	keyLength := "256"
	encryption := diagjaxb.EncryptionAlgorithmValue(enumerations.EncryptionAlgorithm_ECDSA)
	digest := diagjaxb.DigestAlgorithmValue(enumerations.DigestAlgorithm_SHA512)
	xml.BasicSignature = &diagjaxb.XmlBasicSignature{
		EncryptionAlgoUsedToSignThisToken: &encryption,
		DigestAlgoUsedToSignThisToken:     &digest,
		KeyLengthUsedToSignThisToken:      &keyLength,
	}
	return diagnostic.NewSignatureWrapper(xml)
}

// collidingByteRangeDiagnosticData mirrors the oracle's synthetic dump: two PAdES
// signatures whose /ByteRange intervals overlap.
func collidingByteRangeDiagnosticData() *diagnostic.DiagnosticData {
	jaxbData := &diagjaxb.XmlDiagnosticData{}
	for i := 0; i < 2; i++ {
		xml := &diagjaxb.XmlSignature{}
		id := diagjaxb.CollapsedString("S-COLLIDE-" + string(rune('0'+i)))
		xml.Id = &id
		byteRange := &diagjaxb.XmlByteRange{Value: diagjaxb.BigIntegerList{
			big.NewInt(0), big.NewInt(100), big.NewInt(200), big.NewInt(300),
		}}
		xml.PDFRevision = &diagjaxb.XmlPDFRevision{
			PDFSignatureDictionary: &diagjaxb.XmlPDFSignatureDictionary{SignatureByteRange: byteRange},
		}
		jaxbData.Signatures = &diagjaxb.SignaturesWrapper{Items: append(jaxbData.Signatures.All(), xml)}
	}
	return diagnostic.NewDiagnosticData(jaxbData)
}

// failingByteRangeDiagnosticData mirrors the oracle's dump whose single PAdES signature
// carries an invalid /ByteRange.
func failingByteRangeDiagnosticData() *diagnostic.DiagnosticData {
	jaxbData := &diagjaxb.XmlDiagnosticData{}
	xml := &diagjaxb.XmlSignature{}
	id := diagjaxb.CollapsedString("S-BAD-BYTERANGE")
	xml.Id = &id
	byteRange := &diagjaxb.XmlByteRange{
		Value: diagjaxb.BigIntegerList{big.NewInt(0), big.NewInt(100), big.NewInt(200), big.NewInt(300)},
		Valid: false,
	}
	xml.PDFRevision = &diagjaxb.XmlPDFRevision{
		PDFSignatureDictionary: &diagjaxb.XmlPDFSignatureDictionary{SignatureByteRange: byteRange},
	}
	jaxbData.Signatures = &diagjaxb.SignaturesWrapper{Items: []*diagjaxb.XmlSignature{xml}}
	return diagnostic.NewDiagnosticData(jaxbData)
}

// signatureWithDuplicatedReference carries a duplicated digest-matcher reference.
func signatureWithDuplicatedReference() *diagnostic.SignatureWrapper {
	xml := &diagjaxb.XmlSignature{}
	id := diagjaxb.CollapsedString("S-SYNTHETIC")
	xml.Id = &id
	uri := "#r-id"
	duplicated := true
	matcherType := diagjaxb.DigestMatcherTypeValue(enumerations.DigestMatcherType_REFERENCE)
	xml.DigestMatchers = &diagjaxb.DigestMatchersWrapper{Items: []*diagjaxb.XmlDigestMatcher{
		{Type: &matcherType, Uri: &uri, Duplicated: &duplicated},
	}}
	return diagnostic.NewSignatureWrapper(xml)
}

// signatureWithDuplicatedSignature is flagged as a duplicated signature.
func signatureWithDuplicatedSignature() *diagnostic.SignatureWrapper {
	xml := &diagjaxb.XmlSignature{}
	id := diagjaxb.CollapsedString("S-SYNTHETIC")
	xml.Id = &id
	duplicated := true
	xml.Duplicated = &duplicated
	return diagnostic.NewSignatureWrapper(xml)
}

// failingPdfRevision mirrors the oracle's revision that fails every lock / consistency /
// overlap check.
func failingPdfRevision() *diagnostic.PDFRevisionWrapper {
	permissions := diagjaxb.CertificationPermissionValue(enumerations.CertificationPermission_NO_CHANGE_PERMITTED)
	action := diagjaxb.PdfLockActionValue(enumerations.PdfLockAction_ALL)
	lock := &diagjaxb.XmlPDFLockDictionary{Action: &action, Permissions: &permissions}

	fieldName := "field-1"
	revision := &diagjaxb.XmlPDFRevision{
		SignatureField: []*diagjaxb.XmlPDFSignatureField{{SigFieldLock: lock}},
		PDFSignatureDictionary: &diagjaxb.XmlPDFSignatureDictionary{
			DocMDP:     &diagjaxb.XmlDocMDP{Permissions: &permissions},
			FieldMDP:   lock,
			Consistent: false,
		},
		ModificationDetection: &diagjaxb.XmlModificationDetection{
			AnnotationOverlap: []*diagjaxb.XmlModification{{Page: diagjaxb.NewBigIntegerFromInt64(1)}},
			ObjectModifications: &diagjaxb.XmlObjectModifications{
				Undefined: []*diagjaxb.XmlObjectModification{{FieldName: &fieldName}},
			},
		},
	}
	return diagnostic.NewPDFRevisionWrapper(revision)
}
