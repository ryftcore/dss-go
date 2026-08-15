package sav

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
	"github.com/utain/esig/dss/validation/process/bbb/aov"
	"github.com/utain/esig/dss/validation/process/vpfltvd"
)

// The sav direct-check KAT: testdata/oracle/sav_direct.jsonl holds the XmlSAV
// upstream produces when each sav check the top-level chains do not exercise
// under the default policy is run alone at Level.FAIL - over every signature and
// time-stamp of the marshal-parity corpus, plus synthetic signatures for the
// branches no dump reaches. See ../testdata/gen/FcSavDirectOracle.java.

type savDirectRow struct {
	savOracleRow
	Check string `json:"check"`
}

// singleSAVChain is a chain of exactly one item, mirroring the oracle's
// SingleSAVChain.
type singleSAVChain struct {
	*process.ChainBase[*jaxb.XmlSAV]
	factory savCheckFactory
}

type savCheckFactory func(result *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV]

func newSingleSAVChain(i18nProvider *i18n.I18nProvider, factory savCheckFactory) *singleSAVChain {
	xmlSAV := &jaxb.XmlSAV{}
	c := &singleSAVChain{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlSAV,
			&xmlSAV.XmlConstraintsConclusionContent, &xmlSAV.XmlConstraintsConclusionAttrs)),
		factory: factory,
	}
	c.InitChainBase(c)
	return c
}

func (c *singleSAVChain) InitChain() {
	c.FirstItem = c.factory(c.Result, process.GetLevelRule(enumerations.Level_FAIL))
}

// savMultiValuesRule is the oracle's ANY / NONE / MATCH_COMMITMENT rule at Level.FAIL.
type savMultiValuesRule struct{ values []string }

func (r savMultiValuesRule) Level() enumerations.Level { return enumerations.Level_FAIL }
func (r savMultiValuesRule) Values() []string          { return r.values }

var (
	savAnyRule      = savMultiValuesRule{values: []string{"*"}}
	savNoneRule     = savMultiValuesRule{values: []string{"NO-SUCH-VALUE"}}
	savCommitmentOK = savMultiValuesRule{values: []string{"1.2.840.113549.1.9.16.6.1"}}
)

func loadSAVDirectRows(t *testing.T, path string) []*savDirectRow {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open oracle: %v", err)
	}
	defer file.Close()

	var rows []*savDirectRow
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<24)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		row := &savDirectRow{}
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

func savCheckBaseName(check string) string {
	if i := strings.LastIndex(check, "-"); i > 0 {
		switch check[i+1:] {
		case "any", "none", "match":
			return check[:i]
		}
	}
	return check
}

func TestSAVDirectChecksAgainstJavaOracle(t *testing.T) {
	rows := loadSAVDirectRows(t, "testdata/oracle/sav_direct.jsonl")
	if len(rows) == 0 {
		t.Fatal("empty direct oracle")
	}
	i18nProvider := i18n.NewI18nProvider()

	type key struct{ file, token, check string }
	index := map[key]*savDirectRow{}
	for _, row := range rows {
		index[key{row.File, row.Token, row.Check}] = row
	}

	statuses := map[string]map[string]int{}
	matched := 0

	run := func(file, token, check string, factory savCheckFactory) {
		want, ok := index[key{file, token, check}]
		if !ok {
			return // Java threw here and recorded no row
		}
		if want.Title == nil {
			empty := ""
			want.Title = &empty
		}
		var result *jaxb.XmlSAV
		if savSafeExecute(func() { result = newSingleSAVChain(i18nProvider, factory).Execute() }) {
			t.Errorf("%s/%s/%s: Go panicked where Java produced a row", file, token, check)
			return
		}
		got := toSAVRow(file, token, enumerations.Context(want.Context),
			&result.XmlConstraintsConclusionContent, result.Title)
		matched++
		base := savCheckBaseName(check)
		if statuses[base] == nil {
			statuses[base] = map[string]int{}
		}
		for _, constraint := range want.Constraints {
			if constraint.Status != nil {
				statuses[base][*constraint.Status]++
			}
		}
		if !reflect.DeepEqual(&want.savOracleRow, got) {
			t.Errorf("%s/%s/%s: mismatch\nexpected: %s\nactual:   %s",
				file, token, check, mustSAVJSON(t, &want.savOracleRow), mustSAVJSON(t, got))
		}
	}

	var files []string
	seen := map[string]bool{}
	for _, row := range rows {
		if row.File == "synthetic" || seen[row.File] {
			continue
		}
		seen[row.File] = true
		files = append(files, row.File)
	}

	for _, name := range files {
		if strings.HasPrefix(name, "model-") {
			continue // schema-coverage fixture, excluded from the direct corpus (see README)
		}
		diagnosticData := loadSAVDiagnosticData(t, name)
		for _, signature := range diagnosticData.Signatures() {
			sig := signature
			id := sig.Id()

			run(name, id, "KeyIdentifierMatchCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewKeyIdentifierMatchCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "MessageDigestOrSignedPropertiesCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewMessageDigestOrSignedPropertiesCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "SigningCertificateAttributePresentCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewSigningCertificateAttributePresentCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "SigningCertificateReferencesValidityCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewSigningCertificateReferencesValidityCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "SigningTimeInCertificateValidityRangeCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewSigningTimeInCertificateValidityRangeCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "StructuralValidationCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewStructuralValidationCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "UnicitySigningCertificateAttributeCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewUnicitySigningCertificateAttributeCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "AllCertificatesInPathReferencedCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewAllCertificatesInPathReferencedCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "ArchiveTimeStampCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewArchiveTimeStampCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "ContentTimeStampCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewContentTimeStampCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "CounterSignatureCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewCounterSignatureCheck(i18nProvider, r, diagnosticData, sig, rule)
			})
			run(name, id, "DocumentTimeStampCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewDocumentTimeStampCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "KeyIdentifierPresentCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewKeyIdentifierPresentCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "SignatureTimeStampCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewSignatureTimeStampCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "SignerLocationCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewSignerLocationCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "ValidationDataRefsOnlyTimeStampCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewValidationDataRefsOnlyTimeStampCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "ValidationDataTimeStampCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewValidationDataTimeStampCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "X509UrlPresentCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewX509UrlPresentCheck(i18nProvider, r, sig, rule)
			})
			run(name, id, "X509UrlMatchCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewX509UrlMatchCheck(i18nProvider, r, sig, rule)
			})

			for _, indication := range []enumerations.Indication{
				enumerations.Indication_PASSED, enumerations.Indication_FAILED,
			} {
				reportTimestamps := savReportTimestamps(diagnosticData, indication)
				token := id + "/" + string(indication)
				run(name, token, "TLevelTimeStampCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
					return NewTLevelTimeStampCheck(i18nProvider, r, sig,
						map[string]*jaxb.XmlBasicBuildingBlocks{}, reportTimestamps, rule)
				})
				run(name, token, "LTALevelTimeStampCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
					return NewLTALevelTimeStampCheck(i18nProvider, r, sig,
						map[string]*jaxb.XmlBasicBuildingBlocks{}, reportTimestamps, rule)
				})
			}

			for _, variant := range []string{"any", "none"} {
				values := savAnyRule
				if variant == "none" {
					values = savNoneRule
				}
				rule := values
				run(name, id, "CertifiedRolesCheck-"+variant, func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
					return NewCertifiedRolesCheck(i18nProvider, r, sig, rule)
				})
				run(name, id, "ClaimedRolesCheck-"+variant, func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
					return NewClaimedRolesCheck(i18nProvider, r, sig, rule)
				})
				run(name, id, "CommitmentTypeIndicationsCheck-"+variant, func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
					return NewCommitmentTypeIndicationsCheck(i18nProvider, r, sig, rule)
				})
				run(name, id, "ContentHintsCheck-"+variant, func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
					return NewContentHintsCheck(i18nProvider, r, sig, rule)
				})
				run(name, id, "ContentIdentifierCheck-"+variant, func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
					return NewContentIdentifierCheck(i18nProvider, r, sig, rule)
				})
				run(name, id, "ContentTypeCheck-"+variant, func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
					return NewContentTypeCheck(i18nProvider, r, sig, rule)
				})
				run(name, id, "SignatureTypeCheck-"+variant, func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
					return NewSignatureTypeCheck(i18nProvider, r, sig, rule)
				})
			}

			var contentTimestamps []*diagnostic.TimestampWrapper
			if savSafeExecute(func() { contentTimestamps = sig.ContentTimestamps() }) {
				contentTimestamps = nil
			}
			for _, contentTimestamp := range contentTimestamps {
				ct := contentTimestamp
				for _, indication := range []enumerations.Indication{
					enumerations.Indication_PASSED, enumerations.Indication_FAILED,
				} {
					conclusion := &jaxb.XmlConclusion{Indication: jaxb.IndicationValue(indication)}
					token := id + "/" + ct.Id() + "/" + string(indication)
					run(name, token, "ContentTimestampBasicValidationCheck",
						func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
							return NewContentTimestampBasicValidationCheck(i18nProvider, r, ct, conclusion, rule)
						})
				}
			}
		}

		var allTimestamps []*diagnostic.TimestampWrapper
		if savSafeExecute(func() { allTimestamps = diagnosticData.TimestampList() }) {
			allTimestamps = nil
		}
		for _, timestamp := range allTimestamps {
			tst := timestamp
			run(name, tst.Id(), "TSAGeneralNameValueMatchCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewTSAGeneralNameValueMatchCheck(i18nProvider, r, tst, rule)
			})
			run(name, tst.Id(), "TimestampMessageImprintWithIdCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return vpfltvd.NewTimestampMessageImprintWithIdCheck(i18nProvider, r, tst, rule)
			})
			run(name, tst.Id(), "TSAGeneralNameFieldPresentCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewTSAGeneralNameFieldPresentCheck(i18nProvider, r, tst, rule)
			})
			run(name, tst.Id(), "TSAGeneralNameOrderMatchCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return NewTSAGeneralNameOrderMatchCheck(i18nProvider, r, tst, rule)
			})
		}
	}

	// --- aov: AlgorithmObsolescenceValidationCheck at each conclusion shape
	for _, shape := range []string{"passed", "passed-with-algo", "passed-with-algo-nokeysize",
		"error", "warning", "info"} {
		aovResult := aovOfShape(shape)
		run("synthetic", "aov-"+shape, "AlgorithmObsolescenceValidationCheck",
			func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
				return aov.NewAlgorithmObsolescenceValidationCheck(i18nProvider, r, aovResult,
					savCurrentTime, i18n.MessageTag_ACCM_POS_SIG_SIG, "T-AOV")
			})
	}

	// --- synthetic signatures for the branches no corpus dump reaches
	roleSignature := savSignatureWithRoles()
	run("synthetic", "certified-roles", "CertifiedRolesCheck-any", func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
		return NewCertifiedRolesCheck(i18nProvider, r, roleSignature, savAnyRule)
	})
	run("synthetic", "claimed-roles", "ClaimedRolesCheck-any", func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
		return NewClaimedRolesCheck(i18nProvider, r, roleSignature, savAnyRule)
	})

	commitmentSignature := savSignatureWithCommitment()
	run("synthetic", "commitment", "CommitmentTypeIndicationsCheck-any", func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
		return NewCommitmentTypeIndicationsCheck(i18nProvider, r, commitmentSignature, savAnyRule)
	})
	run("synthetic", "commitment", "CommitmentTypeIndicationsCheck-match", func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
		return NewCommitmentTypeIndicationsCheck(i18nProvider, r, commitmentSignature, savCommitmentOK)
	})

	x509Signature := savSignatureWithX509Url()
	run("synthetic", "x509url", "X509UrlPresentCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
		return NewX509UrlPresentCheck(i18nProvider, r, x509Signature, rule)
	})
	run("synthetic", "x509url", "X509UrlMatchCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
		return NewX509UrlMatchCheck(i18nProvider, r, x509Signature, rule)
	})

	contentAttrs := savSignatureWithContentAttributes()
	run("synthetic", "content-attrs", "ContentHintsCheck-any", func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
		return NewContentHintsCheck(i18nProvider, r, contentAttrs, savAnyRule)
	})
	run("synthetic", "content-attrs", "ContentIdentifierCheck-any", func(r *process.Result[*jaxb.XmlSAV], _ policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
		return NewContentIdentifierCheck(i18nProvider, r, contentAttrs, savAnyRule)
	})

	kidSignature := savSignatureWithMismatchedKeyIdentifier()
	run("synthetic", "kid-mismatch", "KeyIdentifierMatchCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
		return NewKeyIdentifierMatchCheck(i18nProvider, r, kidSignature, rule)
	})

	vdRefsSignature := savSignatureWithValidationDataRefsOnlyTimestamp()
	run("synthetic", "vd-refs-only", "ValidationDataRefsOnlyTimeStampCheck", func(r *process.Result[*jaxb.XmlSAV], rule policy.LevelRule) process.ChainItem[*jaxb.XmlSAV] {
		return NewValidationDataRefsOnlyTimeStampCheck(i18nProvider, r, vdRefsSignature, rule)
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
	t.Logf("sav direct checks exercised: %d rows over %d check classes", matched, len(statuses))
}

// savReportTimestamps mirrors the oracle's reportTimestamps: one detailed-report
// XmlTimestamp per used time-stamp, each with a basic-validation conclusion
// carrying the given indication.
func savReportTimestamps(diagnosticData *diagnostic.DiagnosticData,
	indication enumerations.Indication) []*jaxb.XmlTimestamp {
	var wrappers []*diagnostic.TimestampWrapper
	if savSafeExecute(func() { wrappers = diagnosticData.TimestampList() }) {
		return nil
	}
	timestamps := make([]*jaxb.XmlTimestamp, 0, len(wrappers))
	for _, wrapper := range wrappers {
		basic := &jaxb.XmlValidationProcessBasicTimestamp{}
		basic.Conclusion = &jaxb.XmlConclusion{Indication: jaxb.IndicationValue(indication)}
		xmlTimestamp := &jaxb.XmlTimestamp{}
		id := wrapper.Id()
		xmlTimestamp.Id = &id
		xmlTimestamp.ValidationProcessBasicTimestamp = basic
		timestamps = append(timestamps, xmlTimestamp)
	}
	return timestamps
}

func savBaseSignature() *diagjaxb.XmlSignature {
	xml := &diagjaxb.XmlSignature{}
	id := diagjaxb.CollapsedString("S-SYNTHETIC")
	xml.Id = &id
	return xml
}

func savSignatureWithRoles() *diagnostic.SignatureWrapper {
	xml := savBaseSignature()
	claimedCategory := diagjaxb.EndorsementTypeValue(enumerations.EndorsementType_CLAIMED)
	certifiedCategory := diagjaxb.EndorsementTypeValue(enumerations.EndorsementType_CERTIFIED)
	xml.SignerRole = []*diagjaxb.XmlSignerRole{
		{Role: strPtr("claimed-role"), Category: &claimedCategory},
		{Role: strPtr("certified-role"), Category: &certifiedCategory},
	}
	return diagnostic.NewSignatureWrapper(xml)
}

func savSignatureWithCommitment() *diagnostic.SignatureWrapper {
	xml := savBaseSignature()
	xml.CommitmentTypeIndications = &diagjaxb.CommitmentTypeIndicationsWrapper{
		Items: []*diagjaxb.XmlCommitmentTypeIndication{{Identifier: strPtr("1.2.840.113549.1.9.16.6.1")}},
	}
	return diagnostic.NewSignatureWrapper(xml)
}

func savSignatureWithX509Url() *diagnostic.SignatureWrapper {
	xml := savBaseSignature()
	certificateId := diagjaxb.CollapsedString("C-SYNTHETIC")
	certificate := &diagjaxb.XmlCertificate{}
	certificate.Id = &certificateId

	xml.SigningCertificate = &diagjaxb.XmlSigningCertificate{Certificate: certificate}

	origin := diagjaxb.CertificateRefOriginValue(enumerations.CertificateRefOrigin_X509_URL)
	related := &diagjaxb.XmlRelatedCertificate{Certificate: certificate}
	related.CertificateRef = []*diagjaxb.XmlCertificateRef{{Origin: &origin}}
	xml.FoundCertificates = &diagjaxb.XmlFoundCertificates{
		RelatedCertificate: []*diagjaxb.XmlRelatedCertificate{related},
	}
	return diagnostic.NewSignatureWrapper(xml)
}

func savSignatureWithValidationDataRefsOnlyTimestamp() *diagnostic.SignatureWrapper {
	xml := savBaseSignature()
	timestampId := diagjaxb.CollapsedString("T-SYNTHETIC")
	timestampType := diagjaxb.TimestampTypeValue(enumerations.TimestampType_VALIDATION_DATA_REFSONLY_TIMESTAMP)
	timestamp := &diagjaxb.XmlTimestamp{Type: &timestampType}
	timestamp.Id = &timestampId
	xml.FoundTimestamps = &diagjaxb.FoundTimestampsWrapper{
		Items: []*diagjaxb.XmlFoundTimestamp{{Timestamp: timestamp}},
	}
	return diagnostic.NewSignatureWrapper(xml)
}

func strPtr(value string) *string { return &value }

// aovOfShape mirrors the oracle's aovOfShape: an XmlAOV whose conclusion carries the
// requested message shape, plus - for the "with-algo" shapes - the cryptographic
// validation the check renders into its additional info.
func aovOfShape(shape string) *jaxb.XmlAOV {
	result := &jaxb.XmlAOV{}
	conclusion := &jaxb.XmlConclusion{}
	switch shape {
	case "error":
		conclusion.Indication = jaxb.IndicationValue(enumerations.Indication_INDETERMINATE)
		sub := jaxb.SubIndicationValue(enumerations.SubIndication_CRYPTO_CONSTRAINTS_FAILURE)
		conclusion.SubIndication = &sub
		conclusion.Errors = append(conclusion.Errors,
			aovMessage("ASCCM_AR_ANS_ANR", "The algorithm is no longer reliable!"))
	case "warning":
		conclusion.Indication = jaxb.IndicationValue(enumerations.Indication_PASSED)
		conclusion.Warnings = append(conclusion.Warnings,
			aovMessage("ASCCM_AR_ANS_AKSNR", "The key size is no longer reliable!"))
	case "info":
		conclusion.Indication = jaxb.IndicationValue(enumerations.Indication_PASSED)
		conclusion.Infos = append(conclusion.Infos,
			aovMessage("ASCCM_AR_ANS_ANR", "The algorithm expires soon."))
	default:
		conclusion.Indication = jaxb.IndicationValue(enumerations.Indication_PASSED)
	}
	result.Conclusion = conclusion
	if shape == "passed-with-algo" || shape == "passed-with-algo-nokeysize" {
		algorithm := &jaxb.XmlCryptographicAlgorithm{Name: "RSA with SHA256"}
		if shape == "passed-with-algo" {
			keyLength := "2048"
			algorithm.KeyLength = &keyLength
		}
		result.SignatureCryptographicValidation = &jaxb.XmlCryptographicValidation{Algorithm: algorithm}
	}
	return result
}

func aovMessage(key, value string) *jaxb.XmlMessage {
	k := key
	return &jaxb.XmlMessage{Key: &k, Value: value}
}

// savSignatureWithMismatchedKeyIdentifier mirrors the oracle's synthetic signature whose
// KEY_IDENTIFIER reference does not match the issuer serial.
func savSignatureWithMismatchedKeyIdentifier() *diagnostic.SignatureWrapper {
	xml := savBaseSignature()
	certificateId := diagjaxb.CollapsedString("C-SYNTHETIC")
	certificate := &diagjaxb.XmlCertificate{}
	certificate.Id = &certificateId

	match := false
	origin := diagjaxb.CertificateRefOriginValue(enumerations.CertificateRefOrigin_KEY_IDENTIFIER)
	ref := &diagjaxb.XmlCertificateRef{
		Origin:       &origin,
		IssuerSerial: &diagjaxb.XmlIssuerSerial{Value: diagjaxb.Base64Binary([]byte{1, 2, 3}), Match: &match},
	}
	related := &diagjaxb.XmlRelatedCertificate{Certificate: certificate}
	related.CertificateRef = []*diagjaxb.XmlCertificateRef{ref}
	xml.FoundCertificates = &diagjaxb.XmlFoundCertificates{
		RelatedCertificate: []*diagjaxb.XmlRelatedCertificate{related},
	}
	return diagnostic.NewSignatureWrapper(xml)
}

// savSignatureWithContentAttributes carries content-type / content-hints /
// content-identifier attributes.
func savSignatureWithContentAttributes() *diagnostic.SignatureWrapper {
	xml := savBaseSignature()
	xml.ContentType = strPtr("1.2.840.113549.1.7.1")
	xml.ContentHints = strPtr("content-hints")
	xml.ContentIdentifier = strPtr("content-identifier")
	return diagnostic.NewSignatureWrapper(xml)
}
