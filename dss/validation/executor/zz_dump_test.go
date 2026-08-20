package executor

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	detailedreportjaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	dsspolicy "github.com/utain/esig/dss/policy"
	validationpolicy "github.com/utain/esig/dss/validation/policy"
)

func TestZZDump(t *testing.T) {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())
	out := os.Getenv("DUMP_OUT")
	for _, base := range []string{"er-asn1-incorrect-hash.asice", "pades-5-signatures-and-1-document-timestamp.pdf"} {
		data, err := os.ReadFile(filepath.Join("..", "..", "diagnostic", "jaxb", "testdata", "oracle", base+".xml"))
		if err != nil {
			t.Fatal(err)
		}
		dd, err := diagnosticjaxb.Unmarshal(data)
		if err != nil {
			t.Fatal(err)
		}
		e := NewDefaultSignatureProcessExecutor()
		e.SetDiagnosticData(dd)
		e.SetValidationPolicy(validationpolicy.FromDefaultValidationPolicy().Create())
		e.SetCurrentTime(time.UnixMilli(1700000000000).UTC())
		e.SetValidationLevel(enumerations.ValidationLevel_ARCHIVAL_DATA)
		e.SetEnableEtsiValidationReport(true)
		e.SetLocale("en")
		r := e.Execute()
		dr, err := detailedreportjaxb.Marshal(r.GetDetailedReportJaxb())
		if err != nil {
			t.Fatal(err)
		}
		_ = os.WriteFile(filepath.Join(out, base+".go.dr.xml"), dr, 0644)
	}
}
