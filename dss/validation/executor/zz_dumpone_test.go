package executor

import (
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
	validationpolicy "github.com/utain/esig/dss/validation/policy"
)

func TestZZDumpOne(t *testing.T) {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())
	certPolicy, _ := os.ReadFile("../resources/certificate-constraint.xml")
	out := os.Getenv("DUMP_OUT")
	now := time.UnixMilli(1700000000000).UTC()
	for _, p := range strings.Split(os.Getenv("DUMP_FILES"), ",") {
		raw, err := os.ReadFile(filepath.Join("testdata", "oracle", "full-corpus", p))
		if err != nil {
			t.Fatal(err)
		}
		base := strings.ReplaceAll(p, "/", "_")
		dd, err := diagnosticjaxb.Unmarshal(raw)
		if err != nil {
			t.Fatal(err)
		}
		e := NewDefaultSignatureProcessExecutor()
		e.SetDiagnosticData(dd)
		e.SetValidationPolicy(validationpolicy.FromDefaultValidationPolicy().Create())
		e.SetCurrentTime(now)
		e.SetValidationLevel(enumerations.ValidationLevel_ARCHIVAL_DATA)
		e.SetEnableEtsiValidationReport(true)
		e.SetLocale("en")
		r := e.Execute()
		b, _ := detailedreportjaxb.Marshal(r.GetDetailedReportJaxb())
		_ = os.WriteFile(filepath.Join(out, base+".dr.xml"), b, 0644)

		dd2, _ := diagnosticjaxb.Unmarshal(raw)
		certs := dd2.UsedCertificates.All()
		if len(certs) == 0 || certs[0].Id == nil {
			continue
		}
		c := NewDefaultCertificateProcessExecutor()
		c.SetDiagnosticData(dd2)
		c.SetValidationPolicy(validationpolicy.FromValidationPolicyDocument(model.NewInMemoryDocument(certPolicy)).Create())
		c.SetCurrentTime(now)
		c.SetCertificateId(string(*certs[0].Id))
		c.SetLocale("en")
		cr := c.Execute()
		b, _ = detailedreportjaxb.Marshal(cr.GetDetailedReportJaxb())
		_ = os.WriteFile(filepath.Join(out, base+".cdr.xml"), b, 0644)
		b, _ = simplecertjaxb.Marshal(cr.GetSimpleReportJaxb())
		_ = os.WriteFile(filepath.Join(out, base+".scr.xml"), b, 0644)
	}
}
