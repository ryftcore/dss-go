// THE PHASE 8 EXIT CRITERION - report BYTE parity beyond the harness subset.
//
// TestReportBuildersOracle (report_builders_oracle_test.go) compares marshalled
// report bytes on the 50-document marshal-parity corpus. This test does the
// same over EVERY diagnostic-data document of the full upstream diag-data
// corpus already vendored under testdata/oracle/full-corpus/ for
// TestFullCorpusExecutorOracle - all 252 of them, 1255 marshalled reports:
//
//   - SimpleReport, DetailedReport and the ETSI ValidationReport of
//     DefaultSignatureProcessExecutor (default validation policy,
//     ValidationLevel.ARCHIVAL_DATA, validation time 1700000000000, locale en,
//     ETSI report ENABLED - unlike the verdict-only full-corpus oracle);
//   - SimpleCertificateReport and DetailedReport of
//     DefaultCertificateProcessExecutor over the first <UsedCertificates> entry
//     with /policy/certificate-constraint.xml.
//
// testdata/oracle/full_corpus_reports.jsonl is a pure Java dump produced by
// testdata/oracle/gen/FullCorpusReportsOracle.java; it holds the SHA-256 of
// each report, so an equal digest means the Go bytes are the Java bytes.
//
// The 21 files the Java oracle could not unmarshal - policy and crypto-suite
// constraint fixtures misfiled into the diag-data tree - carry no row at all:
// this file only checks byte parity, and TestFullCorpusExecutorOracle already
// asserts the unmarshal outcome of all 273 files on both sides.
//
// It admits NO tolerances. Five separate parity defects were found by exactly
// this comparison and fixed rather than recorded:
// AOV_XCV constraint order and the digest algorithm reported for the signed
// attributes (dss/validation/process/bbb/aov), <Timestamp> order under
// <EvidenceRecord> (dss/validation/process/vpfswatsp/evidencerecord), trusted-
// list constraint order (dss/validation/process/qualification), a stray empty
// <Indication> in the simple certificate report (dss/simplecertificatereport/
// jaxb), and the whole certificate-approval-status ("certificate usage",
// ETSI TS 119 602) process concluding FAILED because <TrustedEntity LoTE="..."/>
// IDREFs never resolved (dss/diagnostic/jaxb).
package executor

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	detailedreportjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
	dsspolicy "github.com/ryftcore/dss-go/dss/policy"
	simplecertjaxb "github.com/ryftcore/dss-go/dss/simplecertificatereport/jaxb"
	simplereportjaxb "github.com/ryftcore/dss-go/dss/simplereport/jaxb"
	validationpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
	validationreportjaxb "github.com/ryftcore/dss-go/dss/validationreport/jaxb"
)

type fullCorpusReportsRow struct {
	File                      string `json:"file"`
	SimpleReport              string `json:"simpleReport"`
	DetailedReport            string `json:"detailedReport"`
	EtsiValidationReport      string `json:"etsiValidationReport"`
	CertificateID             string `json:"certificateId"`
	SimpleCertificateReport   string `json:"simpleCertificateReport"`
	CertificateDetailedReport string `json:"certificateDetailedReport"`
}

// eaaOnlyDetailedReports lists the documents whose signature DetailedReport
// carries the EAA validation blocks the `eaa` build tag gates in. Under the
// default build the Go report legitimately lacks them, so their digest is
// checked only when the tag is set - and the entry is asserted to be NEEDED in
// the default build, so it cannot rot if the gating ever goes away.
var eaaOnlyDetailedReports = map[string]bool{
	"eaa-validation/diag_data_eaa.xml":                true,
	"eaa-validation/diag_data_eaa_no_disclosures.xml": true,
	"eaa-validation/diag_data_pid.xml":                true,
}

func TestFullCorpusReportsByteParity(t *testing.T) {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())

	certificatePolicyData, err := os.ReadFile("../resources/certificate-constraint.xml")
	if err != nil {
		t.Fatalf("reading the certificate validation policy: %v", err)
	}

	data, err := os.ReadFile(corpustest.Path(t, filepath.Join("oracle", "full_corpus_reports.jsonl")))
	if err != nil {
		t.Fatalf("reading the full-corpus report oracle dump: %v", err)
	}

	var rows []fullCorpusReportsRow
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var row fullCorpusReportsRow
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			t.Fatalf("parsing the full-corpus report oracle dump: %v", err)
		}
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		t.Fatal("empty full-corpus report oracle")
	}

	compared := 0
	for _, row := range rows {
		row := row
		t.Run(row.File, func(t *testing.T) {
			raw, err := os.ReadFile(corpustest.Path(t, filepath.Join("oracle", "full-corpus", row.File)))
			if err != nil {
				t.Fatalf("reading %s: %v", row.File, err)
			}
			diagnosticData, err := diagnosticjaxb.Unmarshal(raw)
			if err != nil {
				t.Fatalf("Go failed to unmarshal %s (%v) but the Java oracle produced reports for it - DEFECT", row.File, err)
			}

			signatureExecutor := NewDefaultSignatureProcessExecutor()
			signatureExecutor.SetDiagnosticData(diagnosticData)
			signatureExecutor.SetValidationPolicy(validationpolicy.FromDefaultValidationPolicy().Create())
			signatureExecutor.SetCurrentTime(fullCorpusValidationTime)
			signatureExecutor.SetValidationLevel(enumerations.ValidationLevelArchivalData)
			signatureExecutor.SetEnableEtsiValidationReport(true)
			signatureExecutor.SetLocale("en")
			reports := signatureExecutor.Execute()

			simpleReport, err := simplereportjaxb.Marshal(reports.GetSimpleReportJaxb())
			if err != nil {
				t.Fatalf("marshalling the simple report: %v", err)
			}
			compared += assertFullCorpusDigest(t, "SimpleReport", simpleReport, row.SimpleReport, true)

			detailedReport, err := detailedreportjaxb.Marshal(reports.GetDetailedReportJaxb())
			if err != nil {
				t.Fatalf("marshalling the detailed report: %v", err)
			}
			compared += assertFullCorpusDigest(t, "DetailedReport", detailedReport, row.DetailedReport,
				eaaBuildTagEnabled || !eaaOnlyDetailedReports[row.File])

			validationReport, err := validationreportjaxb.Marshal(reports.GetEtsiValidationReportJaxb())
			if err != nil {
				t.Fatalf("marshalling the ETSI validation report: %v", err)
			}
			compared += assertFullCorpusDigest(t, "ETSIValidationReport", validationReport, row.EtsiValidationReport, true)

			if row.CertificateID == "" {
				return
			}
			certificateData, err := diagnosticjaxb.Unmarshal(raw)
			if err != nil {
				t.Fatalf("re-reading %s for the certificate executor: %v", row.File, err)
			}
			certificateExecutor := NewDefaultCertificateProcessExecutor()
			certificateExecutor.SetDiagnosticData(certificateData)
			certificateExecutor.SetValidationPolicy(validationpolicy.FromValidationPolicyDocument(
				model.NewInMemoryDocument(certificatePolicyData)).Create())
			certificateExecutor.SetCurrentTime(fullCorpusValidationTime)
			certificateExecutor.SetCertificateId(row.CertificateID)
			certificateExecutor.SetLocale("en")
			certificateReports := certificateExecutor.Execute()

			simpleCertificateReport, err := simplecertjaxb.Marshal(certificateReports.GetSimpleReportJaxb())
			if err != nil {
				t.Fatalf("marshalling the simple certificate report: %v", err)
			}
			compared += assertFullCorpusDigest(t, "SimpleCertificateReport", simpleCertificateReport, row.SimpleCertificateReport, true)

			certificateDetailedReport, err := detailedreportjaxb.Marshal(certificateReports.GetDetailedReportJaxb())
			if err != nil {
				t.Fatalf("marshalling the certificate detailed report: %v", err)
			}
			compared += assertFullCorpusDigest(t, "certificate DetailedReport", certificateDetailedReport, row.CertificateDetailedReport, true)
		})
	}
	t.Logf("full-corpus report byte parity: %d rows, %d marshalled reports compared", len(rows), compared)
}

// assertFullCorpusDigest compares one marshalled report against the Java
// digest. When mustMatch is false (an EAA-gated DetailedReport in the default
// build) it asserts the opposite: the report must still DIFFER, so the
// eaaOnlyDetailedReports entry cannot outlive the gating it describes.
func assertFullCorpusDigest(t *testing.T, label string, got []byte, want string, mustMatch bool) int {
	t.Helper()
	sum := sha256.Sum256(got)
	matched := hex.EncodeToString(sum[:]) == want
	switch {
	case mustMatch && !matched:
		t.Errorf("%s differs from the Java oracle (sha256 %s, want %s)", label, hex.EncodeToString(sum[:]), want)
	case !mustMatch && matched:
		t.Errorf("%s now matches the Java oracle without the eaa build tag: remove the eaaOnlyDetailedReports entry", label)
	}
	if mustMatch {
		return 1
	}
	return 0
}
