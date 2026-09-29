//go:build !eaa

package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	dsspolicy "github.com/ryftcore/dss-go/dss/policy"
	validationpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
)

// TestEAAPresentationProcessExecutor_WithoutEAABuildTag covers T32-SEC-001. Diagnostic data
// parsing and the EAA presentation executor are untagged, but a build without the `eaa` build
// tag has no EAA format-checking block, so the basic building blocks of an EAA carry no FC. The
// presentation process must then refuse with an explicit "requires the 'eaa' build tag" panic
// (turned into an error at the facade) rather than crash on a nil pointer dereference. The
// tagged build validates the same document (TestFullCorpusReportsByteParity covers its reports).
func TestEAAPresentationProcessExecutor_WithoutEAABuildTag(t *testing.T) {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())
	raw, err := os.ReadFile(corpustest.Path(t, filepath.Join("oracle", "full-corpus", "eaa-validation", "diag_data_eaa.xml")))
	if err != nil {
		t.Fatal(err)
	}
	diagnosticData, err := diagnosticjaxb.Unmarshal(raw)
	if err != nil {
		t.Fatal(err)
	}

	executor := NewEAAPresentationProcessExecutor()
	executor.SetDiagnosticData(diagnosticData)
	executor.SetValidationPolicy(validationpolicy.FromDefaultValidationPolicy().Create())
	executor.SetCurrentTime(fullCorpusValidationTime)
	executor.SetLocale("en")

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected the EAA presentation process to refuse to run without the eaa build tag")
		}
		if _, isRuntimeError := recovered.(error); isRuntimeError {
			t.Fatalf("got a bare runtime error instead of an explicit message: %v", recovered)
		}
		if message := fmt.Sprint(recovered); !strings.Contains(message, "'eaa' build tag") {
			t.Fatalf("panic message %q does not say that the eaa build tag is required", message)
		}
	}()
	executor.Execute()
}
