package validation

import (
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	"github.com/ryftcore/dss-go/dss/validation/executor"
)

// stubAnalyzer is the smallest analyzer.DocumentAnalyzer that
// SignedDocumentValidatorBase can be built on: only IsSupported is reached by
// these tests.
type stubAnalyzer struct {
	analyzer.DocumentAnalyzer
	supported bool
}

func (s *stubAnalyzer) IsSupported(model.DSSDocument) bool { return s.supported }

// stubValidator is a SignedDocumentValidator a test factory can hand back.
type stubValidator struct {
	SignedDocumentValidatorBase
	name string
}

// stubFactory registers under a marker prefix of the document's name, so the
// registry walk can be observed.
type stubFactory struct {
	prefix string
	name   string
}

func (f *stubFactory) IsSupported(document model.DSSDocument) bool {
	return strings.HasPrefix(document.Name(), f.prefix)
}

func (f *stubFactory) Create(document model.DSSDocument) SignedDocumentValidator {
	return &stubValidator{
		SignedDocumentValidatorBase: NewSignedDocumentValidatorBase(&stubAnalyzer{supported: true}),
		name:                        f.name,
	}
}

// TestSignedDocumentValidatorFromDocument covers the ServiceLoader replacement:
// the registry is walked front to back, in registration order, and the first
// supporting factory wins - Java's fromDocument() iterating the provider list.
func TestSignedDocumentValidatorFromDocument(t *testing.T) {
	before := len(documentValidatorFactoryRegistry)
	t.Cleanup(func() { documentValidatorFactoryRegistry = documentValidatorFactoryRegistry[:before] })

	RegisterDocumentValidatorFactory(&stubFactory{prefix: "a", name: "first"})
	RegisterDocumentValidatorFactory(&stubFactory{prefix: "a", name: "second"})
	RegisterDocumentValidatorFactory(&stubFactory{prefix: "b", name: "third"})

	if got := len(DocumentValidatorFactories()); got != before+3 {
		t.Fatalf("registered %d factories, want %d", got, before+3)
	}

	validator, err := SignedDocumentValidatorFromDocument(model.NewInMemoryDocumentWithName([]byte("x"), "abc"))
	if err != nil {
		t.Fatalf("fromDocument: %v", err)
	}
	if name := validator.(*stubValidator).name; name != "first" {
		t.Errorf("first supporting factory = %q, want %q (registration order decides)", name, "first")
	}

	validator, err = SignedDocumentValidatorFromDocument(model.NewInMemoryDocumentWithName([]byte("x"), "bcd"))
	if err != nil {
		t.Fatalf("fromDocument: %v", err)
	}
	if name := validator.(*stubValidator).name; name != "third" {
		t.Errorf("selected factory = %q, want %q", name, "third")
	}

	// Java throws UnsupportedOperationException("Document format not
	// recognized/handled") when no provider supports the document.
	if _, err := SignedDocumentValidatorFromDocument(model.NewInMemoryDocumentWithName([]byte("x"), "zzz")); err == nil {
		t.Error("expected an error for an unsupported document")
	} else if !strings.Contains(err.Error(), "Document format not recognized/handled") {
		t.Errorf("error = %q, want Java's message", err.Error())
	}
}

// TestSignedDocumentValidatorFromDocumentNil covers Java's
// Objects.requireNonNull("DSSDocument is null").
func TestSignedDocumentValidatorFromDocumentNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic for a nil document")
		} else if r != "DSSDocument is null" {
			t.Errorf("panic = %v, want Java's message", r)
		}
	}()
	_, _ = SignedDocumentValidatorFromDocument(nil)
}

// TestNewSignedDocumentValidatorBaseDefaults pins the field initializers Java's
// SignedDocumentValidator declares, which decide what processValidationPolicy
// hands the executor.
func TestNewSignedDocumentValidatorBaseDefaults(t *testing.T) {
	base := NewSignedDocumentValidatorBase(&stubAnalyzer{})
	if base.defaultDigestAlgorithm != "SHA256" {
		t.Errorf("defaultDigestAlgorithm = %q, want SHA256", base.defaultDigestAlgorithm)
	}
	if base.tokenExtractionStrategy != "NONE" {
		t.Errorf("tokenExtractionStrategy = %q, want NONE", base.tokenExtractionStrategy)
	}
	if base.validationLevel != "ARCHIVAL_DATA" {
		t.Errorf("validationLevel = %q, want ARCHIVAL_DATA", base.validationLevel)
	}
	if base.includeSemantics {
		t.Error("includeSemantics = true, want false")
	}
	if !base.enableEtsiValidationReport {
		t.Error("enableEtsiValidationReport = false, want true")
	}
	// getDefaultProcessExecutor() returns a DefaultSignatureProcessExecutor.
	if _, ok := base.DefaultProcessExecutor().(*executor.DefaultSignatureProcessExecutor); !ok {
		t.Errorf("DefaultProcessExecutor() = %T, want *executor.DefaultSignatureProcessExecutor",
			base.DefaultProcessExecutor())
	}
}

// TestNewSignedDocumentValidatorBaseNilAnalyzer covers Java's
// Objects.requireNonNull("DocumentAnalyzer cannot be null!").
func TestNewSignedDocumentValidatorBaseNilAnalyzer(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic for a nil analyzer")
		} else if r != "DocumentAnalyzer cannot be null!" {
			t.Errorf("panic = %v, want Java's message", r)
		}
	}()
	NewSignedDocumentValidatorBase(nil)
}

// compile-time assertion that a validator built on the base satisfies the
// public interface the factory hands back.
var _ SignedDocumentValidator = (*stubValidator)(nil)

// compile-time assertion that the base exposes the analyzer-backed members
// DocumentValidator declares.
var _ interface {
	Signatures() []spivalidation.AdvancedSignature
} = (*SignedDocumentValidatorBase)(nil)
