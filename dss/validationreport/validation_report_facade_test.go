package validationreport

import (
	"os"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/corpustest"
)

// TestFacadeRoundTrip exercises Facade.Marshal/Unmarshal
// (the hand facade path, distinct from jaxb.Marshal/Unmarshal - see
// validation_report_facade.go's header) over one of the marshal-parity
// oracles, checking it parses without error and reproduces the same
// content on a second parse (facade marshalling is not itself pinned
// byte-for-byte against the Java oracle; jaxb.Marshal is - see that
// package's xml_kat_test.go).
func TestFacadeRoundTrip(t *testing.T) {
	data, err := os.ReadFile(corpustest.RootPath(t, "validationreport/jaxb/testdata/oracle/dss1770.xml.xml"))
	if err != nil {
		t.Fatal(err)
	}
	f := NewFacade()
	vr, err := f.UnmarshalString(string(data))
	if err != nil {
		t.Fatalf("UnmarshalString: %v", err)
	}
	if len(vr.SignatureValidationReport) == 0 {
		t.Fatal("expected at least one SignatureValidationReport")
	}
	out, err := f.Marshal(vr)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !strings.Contains(out, "<ValidationReport") {
		t.Errorf("marshalled output missing document element: %s", out)
	}

	vr2, err := f.UnmarshalString(out)
	if err != nil {
		t.Fatalf("second UnmarshalString: %v", err)
	}
	if len(vr2.SignatureValidationReport) != len(vr.SignatureValidationReport) {
		t.Errorf("SignatureValidationReport count changed across re-marshal: %d vs %d",
			len(vr2.SignatureValidationReport), len(vr.SignatureValidationReport))
	}

	if _, err := f.Marshal(nil); err == nil {
		t.Error("Marshal(nil): expected error")
	}
	if _, err := NewFacade().Unmarshal(nil); err == nil {
		t.Error("Unmarshal(nil reader): expected error")
	}
}

// TestSchemaDeferred confirms Schema() surfaces the documented deferral
// rather than silently succeeding.
func TestSchemaDeferred(t *testing.T) {
	if _, err := Schema(); err != ErrXSDSchemaNotSupported {
		t.Errorf("Schema() error = %v, want ErrXSDSchemaNotSupported", err)
	}
}
