package fc

import (
	"bufio"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	diagjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
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

func TestFCDirectChecksAgainstJavaOracle(t *testing.T) {
	rows := loadFCDirectRowsWithCheck(t, "testdata/oracle/fc_direct.jsonl")
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
		if !reflect.DeepEqual(&want.fcOracleRow, got) {
			t.Errorf("%s/%s/%s: mismatch\nexpected: %s\nactual:   %s",
				file, token, check, mustFCJSON(t, &want.fcOracleRow), mustFCJSON(t, got))
		}
	}

	for _, name := range corpusFileNames(rows) {
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
			if rev := sig.PDFRevision(); rev != nil {
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
	run("synthetic", "ecdsa-mismatch", "EllipticCurveKeySizeCheck",
		func(r *process.Result[*jaxb.XmlFC], rule policy.LevelRule) process.ChainItem[*jaxb.XmlFC] {
			return NewEllipticCurveKeySizeCheck(i18nProvider, r, signatureWithEcdsaKeySizeMismatch(), rule)
		})

	if matched != len(rows) {
		t.Errorf("matched %d of %d direct oracle rows", matched, len(rows))
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
