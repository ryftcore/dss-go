// Oracle test for the phase-8f executors and report builders.
//
// The rows in testdata/oracle/reports.jsonl are a pure Java dump, produced by
// testdata/oracle/gen/ReportsOracle.java, which drives the upstream classes
// eu.europa.esig.dss.validation.executor.signature.DefaultSignatureProcessExecutor
// and ...executor.certificate.DefaultCertificateProcessExecutor - the two
// executors this package ports - and records the SHA-256 of every report they
// marshal. The Go side runs the same executors over the same inputs and must
// reproduce those digests, which is byte-parity: equal digests mean equal
// bytes.
//
// The inputs are NOT a private fixture: they are the marshal-parity
// diagnostic-data corpus already shipped at dss/diagnostic/jaxb/testdata/oracle,
// which the phase-8c BasicBuildingBlocks corpus reads too. The four model-*.xml
// schema-coverage fixtures are excluded, for the reason given in
// dss/validation/process/bbb/fc/testdata/README.md.
//
// testdata/oracle/xml holds the FULL Java-marshalled reports for five of the
// rows (one per format family), so a digest mismatch can be diffed rather than
// merely reported.

package executor

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	dsspolicy "github.com/ryftcore/dss-go/dss/policy"
	simplecertjaxb "github.com/ryftcore/dss-go/dss/simplecertificatereport/jaxb"
	simplereportjaxb "github.com/ryftcore/dss-go/dss/simplereport/jaxb"
	validationpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
	validationreportjaxb "github.com/ryftcore/dss-go/dss/validationreport/jaxb"
)

// oracleCorpusRoot is the module-root-relative path (inside the external
// corpus/ tree, see internal/corpustest) of the shipped diagnostic-data
// corpus used as input.
const oracleCorpusRoot = "diagnostic/jaxb/testdata/oracle"

// oracleValidationTime is the validation time the Java dump was produced with
// (1700000000000 ms since the epoch = 2023-11-14T22:13:20Z).
var oracleValidationTime = time.UnixMilli(1700000000000).UTC()

// reportsOracleRow is one row of testdata/oracle/reports.jsonl.
type reportsOracleRow struct {
	File                      string `json:"file"`
	SimpleReport              string `json:"simpleReport"`
	DetailedReport            string `json:"detailedReport"`
	EtsiValidationReport      string `json:"etsiValidationReport"`
	CertificateID             string `json:"certificateId"`
	SimpleCertificateReport   string `json:"simpleCertificateReport"`
	CertificateDetailedReport string `json:"certificateDetailedReport"`
}

func TestReportBuildersOracle(t *testing.T) {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())

	certificatePolicyData, err := os.ReadFile("../resources/certificate-constraint.xml")
	if err != nil {
		t.Fatalf("reading the certificate validation policy: %v", err)
	}

	rows := readReportsOracle(t)
	if len(rows) == 0 {
		t.Fatal("empty oracle corpus")
	}

	for _, row := range rows {
		t.Run(row.File, func(t *testing.T) {
			data, err := os.ReadFile(corpustest.RootPath(t, filepath.Join(oracleCorpusRoot, row.File+".xml")))
			if err != nil {
				t.Fatalf("reading the diagnostic data: %v", err)
			}
			diagnosticData, err := diagnosticjaxb.Unmarshal(data)
			if err != nil {
				t.Fatalf("unmarshalling the diagnostic data: %v", err)
			}

			signatureExecutor := NewDefaultSignatureProcessExecutor()
			signatureExecutor.SetDiagnosticData(diagnosticData)
			signatureExecutor.SetValidationPolicy(validationpolicy.FromDefaultValidationPolicy().Create())
			signatureExecutor.SetCurrentTime(oracleValidationTime)
			signatureExecutor.SetValidationLevel(enumerations.ValidationLevelArchivalData)
			signatureExecutor.SetEnableEtsiValidationReport(true)
			signatureExecutor.SetLocale("en")
			reports := signatureExecutor.Execute()

			simpleReport, err := simplereportjaxb.Marshal(reports.GetSimpleReportJaxb())
			if err != nil {
				t.Fatalf("marshalling the simple report: %v", err)
			}
			assertOracleDigest(t, row.File, ".sr.xml", "SimpleReport", simpleReport, row.SimpleReport)

			detailedReport, err := detailedreportjaxb.Marshal(reports.GetDetailedReportJaxb())
			if err != nil {
				t.Fatalf("marshalling the detailed report: %v", err)
			}
			assertOracleDigest(t, row.File, ".dr.xml", "DetailedReport", detailedReport, row.DetailedReport)

			validationReport, err := validationreportjaxb.Marshal(reports.GetEtsiValidationReportJaxb())
			if err != nil {
				t.Fatalf("marshalling the ETSI validation report: %v", err)
			}
			assertOracleDigest(t, row.File, ".vr.xml", "ETSI ValidationReport", validationReport, row.EtsiValidationReport)

			if row.CertificateID == "" {
				return
			}

			// The certificate executor is run over the very same diagnostic
			// data, against the certificate-constraint policy, exactly as the
			// Java dump does.
			certificateData, err := diagnosticjaxb.Unmarshal(data)
			if err != nil {
				t.Fatalf("unmarshalling the diagnostic data: %v", err)
			}
			certificateExecutor := NewDefaultCertificateProcessExecutor()
			certificateExecutor.SetDiagnosticData(certificateData)
			certificateExecutor.SetValidationPolicy(validationpolicy.FromValidationPolicyDocument(
				model.NewInMemoryDocument(certificatePolicyData)).Create())
			certificateExecutor.SetCurrentTime(oracleValidationTime)
			certificateExecutor.SetCertificateId(row.CertificateID)
			certificateExecutor.SetLocale("en")
			certificateReports := certificateExecutor.Execute()

			simpleCertificateReport, err := simplecertjaxb.Marshal(certificateReports.GetSimpleReportJaxb())
			if err != nil {
				t.Fatalf("marshalling the simple certificate report: %v", err)
			}
			assertOracleDigest(t, row.File, ".scr.xml", "SimpleCertificateReport",
				simpleCertificateReport, row.SimpleCertificateReport)

			certificateDetailedReport, err := detailedreportjaxb.Marshal(certificateReports.GetDetailedReportJaxb())
			if err != nil {
				t.Fatalf("marshalling the certificate detailed report: %v", err)
			}
			assertOracleDigest(t, row.File, ".cdr.xml", "certificate DetailedReport",
				certificateDetailedReport, row.CertificateDetailedReport)
		})
	}
}

// knownOrderingDeviations lists the (row, report) pairs whose byte parity is
// blocked by an ACCEPTED, documented iteration-order deviation in a package
// outside this one - not by anything the executors or the report builders do.
// Each entry names the deviation; a row that starts matching is reported as an
// error so the entry gets removed rather than silently rotting.
//
// It is EMPTY: the phase-8f audit closed both entries it used to hold, by
// reproducing the real java.util.HashMap iteration order (see
// utils.JavaHashMapStringKeyOrder / JavaHashMapComputeIfAbsentKeyOrder) at the
// two sites concerned - dss/validation/process/vpfswatsp/evidencerecord's
// <Timestamp> children of <EvidenceRecord>, and dss/validation/process/bbb/aov's
// signCertRefsMap behind the AOV_XCV <Constraint> elements. Every one of this
// corpus' reports is now byte-identical to the Java oracle's.
var knownOrderingDeviations = map[string]string{}

// assertOracleDigest compares the marshalled report against the Java digest,
// falling back to a full diff hint when the row is one of those whose complete
// Java output is shipped under testdata/oracle/xml.
func assertOracleDigest(t *testing.T, file, suffix, label string, got []byte, want string) {
	t.Helper()
	sum := sha256.Sum256(got)
	matched := hex.EncodeToString(sum[:]) == want
	if reason, known := knownOrderingDeviations[file+"/"+label]; known {
		if matched {
			t.Errorf("%s now matches the Java oracle: remove the knownOrderingDeviations entry (%s)", label, reason)
		} else {
			t.Logf("%s differs from the Java oracle, as expected: %s", label, reason)
		}
		return
	}
	if matched {
		return
	}
	reference, err := os.ReadFile(corpustest.Path(t, filepath.Join("oracle", "xml", file+suffix)))
	if err != nil {
		t.Errorf("%s differs from the Java oracle (sha256 %s, want %s)",
			label, hex.EncodeToString(sum[:]), want)
		return
	}
	t.Errorf("%s differs from the Java oracle:\n%s", label, firstDifference(string(reference), string(got)))
}

// firstDifference renders the first line at which two reports diverge.
func firstDifference(want, got string) string {
	wantLines := strings.Split(want, "\n")
	gotLines := strings.Split(got, "\n")
	for i := 0; i < len(wantLines) || i < len(gotLines); i++ {
		var w, g string
		if i < len(wantLines) {
			w = wantLines[i]
		}
		if i < len(gotLines) {
			g = gotLines[i]
		}
		if w != g {
			return "line " + itoa(i+1) + "\n  java: " + w + "\n  go  : " + g
		}
	}
	return "(no line differs; trailing bytes only)"
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	var digits []byte
	for v > 0 {
		digits = append([]byte{byte('0' + v%10)}, digits...)
		v /= 10
	}
	return string(digits)
}

// readReportsOracle loads testdata/oracle/reports.jsonl, sorted by file name.
func readReportsOracle(t *testing.T) []reportsOracleRow {
	t.Helper()
	data, err := os.ReadFile(corpustest.Path(t, filepath.Join("oracle", "reports.jsonl")))
	if err != nil {
		t.Fatalf("reading the oracle dump: %v", err)
	}
	var rows []reportsOracleRow
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var row reportsOracleRow
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			t.Fatalf("parsing the oracle dump: %v", err)
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].File < rows[j].File })
	return rows
}
