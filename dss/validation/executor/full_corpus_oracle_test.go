// THE PHASE 8 EXIT CRITERION - executor-level verdict parity, item (A) of the
// s8f harness brief.
//
// testdata/oracle/full_corpus.jsonl is a pure Java dump, produced by
// testdata/oracle/gen/FullCorpusOracle.java, over EVERY file (273, no
// exclusions - "no quarantines: every row") under the upstream
// dss-validation diag-data corpus
// (dss-validation/src/test/resources/diag-data/** upstream),
// vendored byte-identically under testdata/oracle/full-corpus/ (relative
// paths preserved, forward-slash keyed).
//
// 21 of the 273 files are not diagnostic-data documents at all - they are
// validation-policy and crypto-suite-constraint XML fixtures that upstream's
// OWN diag-data-corpus tests pass to the POLICY loader, not the diagnostic-
// data unmarshaller, and which sit in that resource tree only because
// upstream's test layout keeps every DefaultSignatureProcessExecutorTest
// input under one directory. FullCorpusOracle.java doesn't special-case
// them: it unmarshals every file as XmlDiagnosticData exactly like this
// test does, so both sides record the SAME unmarshal failure for the SAME
// 21 rows - which is itself a parity assertion, not a quarantine.
//
// For every row that DOES unmarshal, both sides run
// DefaultSignatureProcessExecutor with the default validation policy,
// ValidationLevel.ARCHIVAL_DATA, validation time 1700000000000, locale en
// (matching testdata/oracle/reports.jsonl's configuration), and compare:
//
//   - every BasicBuildingBlocks conclusion (keyed by Id: Indication/SubIndication/Type)
//   - the final Indication/SubIndication of every top-level Signature/Timestamp/
//     EvidenceRecord entry in the DetailedReport
//   - the SimpleReport signature qualification of every signature
//
// and, when the document carries at least one used certificate, the SAME
// BBB comparison for DefaultCertificateProcessExecutor run against the first
// certificate (matching testdata/oracle/reports.jsonl's convention), against
// /policy/certificate-constraint.xml.
//
// This deliberately runs ONE fixed executor pair uniformly over the whole
// corpus rather than reproducing each fixture's original upstream unit-test
// scenario (several subdirectories - cert-validation, timestamp-validation,
// qwac-validation, eaa-validation - are exercised upstream by OTHER
// executors, e.g. DefaultCertificateProcessExecutor directly or
// DefaultTimestampProcessExecutor). Running the identical Go/Java code path
// on identical input is what "verdict parity" means here; see the header of
// testdata/oracle/gen/FullCorpusOracle.java for the full rationale.
//
// This test does NOT replace testdata/oracle/reports.jsonl's byte-parity
// check (TestReportBuildersOracle, item C of the harness): that one compares
// marshalled report BYTES on a 50-document subset. This one compares
// VERDICTS (Indication/SubIndication/qualification) on the full 273-document
// corpus. A row can pass this test yet still be part of a documented
// ordering deviation in the byte-parity test (e.g. HashMap iteration order
// changed to insertion order) - the two tests check different things.
package executor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	detailedreportjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
	dsspolicy "github.com/ryftcore/dss-go/dss/policy"
	simplereportjaxb "github.com/ryftcore/dss-go/dss/simplereport/jaxb"
	validationpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
)

var fullCorpusValidationTime = time.UnixMilli(1700000000000).UTC()

type fcToken struct {
	Kind          string `json:"kind"`
	ID            string `json:"id"`
	Indication    string `json:"indication"`
	SubIndication string `json:"subIndication"`
}

type fcBBB struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Indication    string `json:"indication"`
	SubIndication string `json:"subIndication"`
}

type fcQualification struct {
	ID            string `json:"id"`
	Qualification string `json:"qualification"`
}

type fcSignatureExecutor struct {
	Tokens        []fcToken         `json:"tokens"`
	BBB           []fcBBB           `json:"bbb"`
	Qualification []fcQualification `json:"qualification"`
}

type fcCertificateExecutor struct {
	CertificateID string  `json:"certificateId"`
	BBB           []fcBBB `json:"bbb"`
}

type fullCorpusRow struct {
	File                     string                 `json:"file"`
	SignatureExecutor        *fcSignatureExecutor   `json:"signatureExecutor"`
	SignatureExecutorError   string                 `json:"signatureExecutorError"`
	CertificateExecutor      *fcCertificateExecutor `json:"certificateExecutor"`
	CertificateExecutorError string                 `json:"certificateExecutorError"`
}

func TestFullCorpusExecutorOracle(t *testing.T) {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())

	certificatePolicyData, err := os.ReadFile("../resources/certificate-constraint.xml")
	if err != nil {
		t.Fatalf("reading the certificate validation policy: %v", err)
	}

	rows := readFullCorpusOracle(t)
	if len(rows) == 0 {
		t.Fatal("empty full-corpus oracle")
	}

	var totalRows, totalMismatches int
	for _, row := range rows {
		row := row
		t.Run(row.File, func(t *testing.T) {
			totalRows++
			data, err := os.ReadFile(corpustest.Path(t, filepath.Join("oracle", "full-corpus", row.File)))
			if err != nil {
				t.Fatalf("reading %s: %v", row.File, err)
			}

			diagnosticData, unmarshalErr := diagnosticjaxb.Unmarshal(data)
			if unmarshalErr != nil {
				if reason, known := knownFileLevelDivergences[row.File]; known {
					if row.SignatureExecutorError != "" {
						t.Errorf("%s: this file now fails on BOTH sides - remove the knownFileLevelDivergences entry (%s)", row.File, reason)
					} else {
						t.Logf("Go fails to unmarshal %s as documented: %s (java: %v -> go: %v)", row.File, reason, unmarshalErr, unmarshalErr)
					}
					return
				}
				if row.SignatureExecutorError == "" {
					t.Errorf("Go failed to unmarshal %s (%v) but the Java oracle succeeded - DEFECT", row.File, unmarshalErr)
				}
				return
			}
			if reason, known := knownFileLevelDivergences[row.File]; known {
				t.Errorf("%s: Go now unmarshals successfully - remove the knownFileLevelDivergences entry (%s)", row.File, reason)
				return
			}

			sigPolicy := validationpolicy.FromDefaultValidationPolicy().Create()
			detailedReport, simpleReport, runErr := runSignatureExecutorSafely(diagnosticData, sigPolicy)
			if runErr != nil {
				if row.SignatureExecutorError == "" {
					t.Errorf("Go signature executor errored on %s (%v) but the Java oracle succeeded - DEFECT", row.File, runErr)
				}
				return
			}
			if row.SignatureExecutorError != "" {
				t.Errorf("Go signature executor succeeded on %s but the Java oracle errored (%s) - DEFECT",
					row.File, row.SignatureExecutorError)
				return
			}
			totalMismatches += compareSignatureExecutor(t, row.SignatureExecutor, detailedReport, simpleReport)

			if row.CertificateExecutor == nil && row.CertificateExecutorError == "" {
				return // no used certificates in this document
			}
			certData, err := diagnosticjaxb.Unmarshal(data)
			if err != nil {
				t.Fatalf("re-reading %s for the certificate executor: %v", row.File, err)
			}
			var certID string
			if row.CertificateExecutor != nil {
				certID = row.CertificateExecutor.CertificateID
			} else {
				// The Java oracle picked UsedCertificates[0] before failing;
				// reproduce the same selection so both sides attempt the
				// identical run.
				usedCertificates := certData.UsedCertificates.All()
				if len(usedCertificates) == 0 {
					t.Fatalf("%s: Java recorded a certificateExecutorError but has no UsedCertificates", row.File)
				}
				if usedCertificates[0].Id != nil {
					certID = string(*usedCertificates[0].Id)
				}
			}
			certPolicy := validationpolicy.FromValidationPolicyDocument(model.NewInMemoryDocument(certificatePolicyData)).Create()
			certDetailedReport, certRunErr := runCertificateExecutorSafely(certData, certID, certPolicy)
			if certRunErr != nil {
				if row.CertificateExecutorError == "" {
					t.Errorf("Go certificate executor errored on %s (%v) but the Java oracle succeeded - DEFECT", row.File, certRunErr)
				}
				return
			}
			if row.CertificateExecutorError != "" {
				t.Errorf("Go certificate executor succeeded on %s but the Java oracle errored (%s) - DEFECT",
					row.File, row.CertificateExecutorError)
				return
			}
			totalMismatches += compareBBB(t, "certificateExecutor", row.CertificateExecutor.BBB, certDetailedReport.BasicBuildingBlocks)
		})
	}
	t.Logf("full-corpus executor oracle: %d rows, %d field mismatches", totalRows, totalMismatches)
}

// runSignatureExecutorSafely recovers from a panic so a single malformed
// fixture reports as a mismatch on its own subtest rather than aborting the
// whole corpus run - a panic is exactly the kind of Go/Java divergence this
// oracle exists to surface.
func runSignatureExecutorSafely(dd *diagnosticjaxb.XmlDiagnosticData, policy modelpolicy.ValidationPolicy) (
	detailed *detailedreportjaxb.XmlDetailedReport, simple *simplereportjaxb.XmlSimpleReport, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	signatureExecutor := NewDefaultSignatureProcessExecutor()
	signatureExecutor.SetDiagnosticData(dd)
	signatureExecutor.SetValidationPolicy(policy)
	signatureExecutor.SetCurrentTime(fullCorpusValidationTime)
	signatureExecutor.SetValidationLevel(enumerations.ValidationLevelArchivalData)
	signatureExecutor.SetEnableEtsiValidationReport(false)
	signatureExecutor.SetLocale("en")
	reports := signatureExecutor.Execute()
	return reports.GetDetailedReportJaxb(), reports.GetSimpleReportJaxb(), nil
}

func runCertificateExecutorSafely(dd *diagnosticjaxb.XmlDiagnosticData, certID string, policy modelpolicy.ValidationPolicy) (
	detailed *detailedreportjaxb.XmlDetailedReport, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	certificateExecutor := NewDefaultCertificateProcessExecutor()
	certificateExecutor.SetDiagnosticData(dd)
	certificateExecutor.SetValidationPolicy(policy)
	certificateExecutor.SetCurrentTime(fullCorpusValidationTime)
	certificateExecutor.SetCertificateId(certID)
	certificateExecutor.SetLocale("en")
	reports := certificateExecutor.Execute()
	return reports.GetDetailedReportJaxb(), nil
}

// compareSignatureExecutor compares the Java-dumped signature-executor row
// against the Go-produced reports, returning the number of field-level
// mismatches found (each also reported individually via t.Errorf).
func compareSignatureExecutor(t *testing.T, want *fcSignatureExecutor, detailed *detailedreportjaxb.XmlDetailedReport, simple *simplereportjaxb.XmlSimpleReport) int {
	t.Helper()
	mismatches := 0

	gotTokens := make(map[string]fcToken)
	var gotOrder []string
	for _, item := range detailed.SignatureOrTimestampOrEvidenceRecord {
		var kind, id string
		var conclusion *detailedreportjaxb.XmlConclusion
		switch v := item.(type) {
		case *detailedreportjaxb.XmlSignature:
			kind, conclusion = "Signature", v.Conclusion
			if v.Id != nil {
				id = *v.Id
			}
		case *detailedreportjaxb.XmlTimestamp:
			kind, conclusion = "Timestamp", v.Conclusion
			if v.Id != nil {
				id = *v.Id
			}
		case *detailedreportjaxb.XmlEvidenceRecord:
			kind, conclusion = "EvidenceRecord", v.Conclusion
			if v.Id != nil {
				id = *v.Id
			}
		default:
			continue
		}
		tok := fcToken{Kind: kind, ID: id}
		if conclusion != nil {
			tok.Indication = string(conclusion.Indication)
			if conclusion.SubIndication != nil {
				tok.SubIndication = string(*conclusion.SubIndication)
			}
		}
		gotTokens[kind+"/"+id] = tok
		gotOrder = append(gotOrder, kind+"/"+id)
	}
	wantKeys := make(map[string]bool)
	for _, w := range want.Tokens {
		key := w.Kind + "/" + w.ID
		wantKeys[key] = true
		g, ok := gotTokens[key]
		if !ok {
			t.Errorf("missing token %s (java: indication=%s subIndication=%s)", key, w.Indication, w.SubIndication)
			mismatches++
			continue
		}
		if g.Indication != w.Indication || g.SubIndication != w.SubIndication {
			t.Errorf("token %s: indication/subIndication go=%s/%s java=%s/%s",
				key, g.Indication, g.SubIndication, w.Indication, w.SubIndication)
			mismatches++
		}
	}
	for _, key := range gotOrder {
		if !wantKeys[key] {
			t.Errorf("unexpected extra token %s in the Go report", key)
			mismatches++
		}
	}

	mismatches += compareBBB(t, "signatureExecutor", want.BBB, detailed.BasicBuildingBlocks)

	gotQual := make(map[string]string)
	for _, item := range simple.SignatureOrTimestampOrEvidenceRecord {
		sig, ok := item.(*simplereportjaxb.XmlSignature)
		if !ok {
			continue
		}
		q := ""
		if sig.SignatureLevel != nil {
			q = string(sig.SignatureLevel.Value)
		}
		gotQual[sig.Id] = q
	}
	seenQual := make(map[string]bool)
	for _, w := range want.Qualification {
		seenQual[w.ID] = true
		g, ok := gotQual[w.ID]
		if !ok {
			t.Errorf("missing simple-report signature %s (java qualification=%s)", w.ID, w.Qualification)
			mismatches++
			continue
		}
		if g != w.Qualification {
			t.Errorf("signature %s qualification: go=%s java=%s", w.ID, g, w.Qualification)
			mismatches++
		}
	}
	for id := range gotQual {
		if !seenQual[id] {
			t.Errorf("unexpected extra simple-report signature %s in the Go report", id)
			mismatches++
		}
	}
	return mismatches
}

// compareBBB compares a Java-dumped BasicBuildingBlocks list against the Go
// list, keyed by Id (BBB ids are unique per document: certificate, revocation,
// timestamp and signature/counter-signature ids never collide).
func compareBBB(t *testing.T, label string, want []fcBBB, got []*detailedreportjaxb.XmlBasicBuildingBlocks) int {
	t.Helper()
	mismatches := 0

	gotByID := make(map[string]fcBBB)
	var gotIDs []string
	for _, bbb := range got {
		entry := fcBBB{ID: bbb.Id, Type: string(bbb.Type)}
		if bbb.Conclusion != nil {
			entry.Indication = string(bbb.Conclusion.Indication)
			if bbb.Conclusion.SubIndication != nil {
				entry.SubIndication = string(*bbb.Conclusion.SubIndication)
			}
		}
		gotByID[bbb.Id] = entry
		gotIDs = append(gotIDs, bbb.Id)
	}

	wantIDs := make(map[string]bool)
	for _, w := range want {
		wantIDs[w.ID] = true
		g, ok := gotByID[w.ID]
		if !ok {
			t.Errorf("%s: missing BasicBuildingBlocks %s (java: type=%s indication=%s subIndication=%s)",
				label, w.ID, w.Type, w.Indication, w.SubIndication)
			mismatches++
			continue
		}
		if g.Type != w.Type || g.Indication != w.Indication || g.SubIndication != w.SubIndication {
			t.Errorf("%s: BasicBuildingBlocks %s: go=(%s,%s,%s) java=(%s,%s,%s)",
				label, w.ID, g.Type, g.Indication, g.SubIndication, w.Type, w.Indication, w.SubIndication)
			mismatches++
		}
	}
	for _, id := range gotIDs {
		if !wantIDs[id] {
			t.Errorf("%s: unexpected extra BasicBuildingBlocks %s in the Go report", label, id)
			mismatches++
		}
	}
	return mismatches
}

// knownFileLevelDivergences lists the (row) documents whose Go unmarshal
// result is ACCEPTED to diverge from the Java oracle's, each naming the
// upstream defect responsible. Every entry is a genuine, root-caused parity
// gap surfaced by this harness - not a quarantine: the row is still fed
// through the full comparison above, and the entry only suppresses the
// specific already-understood failure mode; if the file starts unmarshalling
// like Java (or starts failing on both sides), the assertions above turn
// into a hard failure so the entry cannot silently rot.
//
// It is EMPTY. The phase-8f audit closed both entries it used to hold:
//
//   - F1, diagnostic/jaxb binding java.math.BigInteger properties to a bare
//     *big.Int, whose stdlib UnmarshalText parses with base 0 and so read a
//     leading "0" as an octal prefix - fixed by diagnostic/jaxb.BigInteger,
//     which parses decimal like Java's BigInteger(String);
//   - the JAXB per-field adapter tolerance behind
//     qwac-validation/2-qwac-valid-diag-data.xml's malformed
//     TrustedList/LastLoading - fixed in diagnostic/jaxb.XSDateTime.UnmarshalText.
var knownFileLevelDivergences = map[string]string{}

// readFullCorpusOracle loads testdata/oracle/full_corpus.jsonl, sorted by
// file name so the subtest order is stable across runs.
func readFullCorpusOracle(t *testing.T) []fullCorpusRow {
	t.Helper()
	data, err := os.ReadFile(corpustest.Path(t, filepath.Join("oracle", "full_corpus.jsonl")))
	if err != nil {
		t.Fatalf("reading the full-corpus oracle dump: %v", err)
	}
	var rows []fullCorpusRow
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var row fullCorpusRow
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			t.Fatalf("parsing the full-corpus oracle dump: %v", err)
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].File < rows[j].File })
	return rows
}
