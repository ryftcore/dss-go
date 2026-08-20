package executor

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	detailedreportjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	dsspolicy "github.com/utain/esig/dss/policy"
	simplecertjaxb "github.com/utain/esig/dss/simplecertificatereport/jaxb"
	simplereportjaxb "github.com/utain/esig/dss/simplereport/jaxb"
	validationpolicy "github.com/utain/esig/dss/validation/policy"
	validationreportjaxb "github.com/utain/esig/dss/validationreport/jaxb"
)

type fcrRow struct {
	File                      string `json:"file"`
	SimpleReport              string `json:"simpleReport"`
	DetailedReport            string `json:"detailedReport"`
	EtsiValidationReport      string `json:"etsiValidationReport"`
	SignatureError            string `json:"signatureError"`
	CertificateID             string `json:"certificateId"`
	SimpleCertificateReport   string `json:"simpleCertificateReport"`
	CertificateDetailedReport string `json:"certificateDetailedReport"`
	CertificateError          string `json:"certificateError"`
}

func TestZZFullCorpusReports(t *testing.T) {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())
	certPolicy, err := os.ReadFile("../resources/certificate-constraint.xml")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(os.Getenv("FCR_ORACLE"))
	if err != nil {
		t.Fatal(err)
	}
	now := time.UnixMilli(1700000000000).UTC()
	var total, matched int
	mismatch := map[string]int{}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		var row fcrRow
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(filepath.Join("testdata", "oracle", "full-corpus", row.File))
		if err != nil {
			t.Fatal(err)
		}
		dd, err := diagnosticjaxb.Unmarshal(raw)
		if err != nil {
			t.Logf("SKIP-GO-UNMARSHAL %s: %v", row.File, err)
			continue
		}
		func() {
			defer func() {
				if p := recover(); p != nil {
					t.Errorf("PANIC %s: %v", row.File, p)
				}
			}()
			e := NewDefaultSignatureProcessExecutor()
			e.SetDiagnosticData(dd)
			e.SetValidationPolicy(validationpolicy.FromDefaultValidationPolicy().Create())
			e.SetCurrentTime(now)
			e.SetValidationLevel(enumerations.ValidationLevel_ARCHIVAL_DATA)
			e.SetEnableEtsiValidationReport(true)
			e.SetLocale("en")
			r := e.Execute()
			check := func(label string, got []byte, err error, want string) {
				if err != nil {
					t.Errorf("%s %s marshal: %v", row.File, label, err)
					return
				}
				total++
				sum := sha256.Sum256(got)
				if hex.EncodeToString(sum[:]) == want {
					matched++
				} else {
					mismatch[label]++
					fmt.Printf("MISMATCH %s %s\n", row.File, label)
				}
			}
			b, e1 := simplereportjaxb.Marshal(r.GetSimpleReportJaxb())
			check("SimpleReport", b, e1, row.SimpleReport)
			b, e1 = detailedreportjaxb.Marshal(r.GetDetailedReportJaxb())
			check("DetailedReport", b, e1, row.DetailedReport)
			b, e1 = validationreportjaxb.Marshal(r.GetEtsiValidationReportJaxb())
			check("ETSIValidationReport", b, e1, row.EtsiValidationReport)
		}()
		if row.CertificateID == "" {
			continue
		}
		dd2, err := diagnosticjaxb.Unmarshal(raw)
		if err != nil {
			continue
		}
		func() {
			defer func() {
				if p := recover(); p != nil {
					t.Errorf("PANIC-CERT %s: %v", row.File, p)
				}
			}()
			c := NewDefaultCertificateProcessExecutor()
			c.SetDiagnosticData(dd2)
			c.SetValidationPolicy(validationpolicy.FromValidationPolicyDocument(model.NewInMemoryDocument(certPolicy)).Create())
			c.SetCurrentTime(now)
			c.SetCertificateId(row.CertificateID)
			c.SetLocale("en")
			cr := c.Execute()
			check := func(label string, got []byte, err error, want string) {
				if err != nil {
					t.Errorf("%s %s marshal: %v", row.File, label, err)
					return
				}
				total++
				sum := sha256.Sum256(got)
				if hex.EncodeToString(sum[:]) == want {
					matched++
				} else {
					mismatch[label]++
					fmt.Printf("MISMATCH %s %s\n", row.File, label)
				}
			}
			b, e1 := simplecertjaxb.Marshal(cr.GetSimpleReportJaxb())
			check("SimpleCertificateReport", b, e1, row.SimpleCertificateReport)
			b, e1 = detailedreportjaxb.Marshal(cr.GetDetailedReportJaxb())
			check("CertificateDetailedReport", b, e1, row.CertificateDetailedReport)
		}()
	}
	fmt.Printf("FCR-SUMMARY total=%d matched=%d mismatchByLabel=%v\n", total, matched, mismatch)
}
