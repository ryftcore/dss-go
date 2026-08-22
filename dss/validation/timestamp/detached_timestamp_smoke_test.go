// Smoke test for the detached timestamp analyzer/validator pair and their factories, using the
// raw CMS TimeStampToken fixture already ported (byte-identical) for
// spi/validation/timestamp_token_kat_test.go: DER-encoded PKCS#7 SignedData with an
// id-ct-TSTInfo encapsulated content, SHA-256 message imprint, SHA-512 signature, TSA and CA
// certificates embedded, no timestamped (message-imprint) content attached. Covers format
// detection (IsSupported on both the analyzer and the validator, plus their ServiceLoader-registry
// stand-in factories) and basic TimestampToken construction.
package timestamp

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
)

func loadTestTimestampDocument(t *testing.T) model.DSSDocument {
	t.Helper()
	doc, err := model.NewFileDocument("testdata/timestamp-token.tst")
	if err != nil {
		t.Fatalf("NewFileDocument: %v", err)
	}
	return doc
}

func TestDetachedTimestampAnalyzer_IsSupported(t *testing.T) {
	doc := loadTestTimestampDocument(t)

	if !NewDetachedTimestampAnalyzer(doc).IsSupported(doc) {
		t.Error("DetachedTimestampAnalyzer.IsSupported() = false, want true for a raw CMS TimeStampToken")
	}
	if !NewDetachedTimestampAnalyzerFactory().IsSupported(doc) {
		t.Error("DetachedTimestampAnalyzerFactory.IsSupported() = false, want true")
	}

	notATimestamp := model.NewInMemoryDocument([]byte("not a timestamp token"))
	if NewDetachedTimestampAnalyzer(notATimestamp).IsSupported(notATimestamp) {
		t.Error("DetachedTimestampAnalyzer.IsSupported() = true for non-ASN.1 data, want false")
	}
}

func TestDetachedTimestampValidator_IsSupported(t *testing.T) {
	doc := loadTestTimestampDocument(t)

	if !NewDetachedTimestampValidator(doc).IsSupported(doc) {
		t.Error("DetachedTimestampValidator.IsSupported() = false, want true for a raw CMS TimeStampToken")
	}
	if !NewDetachedTimestampValidatorFactory().IsSupported(doc) {
		t.Error("DetachedTimestampValidatorFactory.IsSupported() = false, want true")
	}
}

func TestDetachedTimestampAnalyzerFactory_RegistersItself(t *testing.T) {
	// init() in detached_timestamp_analyzer_factory.go must have self-registered a
	// *DetachedTimestampAnalyzerFactory with the shared analyzer registry (mirroring
	// upstream's META-INF/services/....DocumentAnalyzerFactory entry), so the generic
	// registry-walking lookup resolves our fixture to a *DetachedTimestampAnalyzer.
	doc := loadTestTimestampDocument(t)
	documentAnalyzer, err := analyzer.DocumentAnalyzerFromDocument(doc)
	if err != nil {
		t.Fatalf("DocumentAnalyzerFromDocument: %v", err)
	}
	if _, ok := documentAnalyzer.(*DetachedTimestampAnalyzer); !ok {
		t.Errorf("DocumentAnalyzerFromDocument() = %T, want *DetachedTimestampAnalyzer", documentAnalyzer)
	}
}

func TestDetachedTimestampValidator_Timestamp(t *testing.T) {
	doc := loadTestTimestampDocument(t)

	timestampValidator := NewDetachedTimestampValidator(doc)
	timestampValidator.SetCertificateVerifier(validation.NewCommonCertificateVerifier())

	timestampToken := timestampValidator.Timestamp()
	if timestampToken == nil {
		t.Fatal("Timestamp() returned nil")
	}
	if timestampToken.GenerationTime().IsZero() {
		t.Error("built TimestampToken's GenerationTime() is zero, want the TSTInfo genTime")
	}

	// Timestamp() caches: a second call must return the identical instance (Java's
	// `if (timestampToken == null) { ... }` guard).
	if second := timestampValidator.Timestamp(); second != timestampToken {
		t.Error("Timestamp() did not return the cached instance on a second call")
	}

	// TimestampedData is unset: matching against it must not error, only fail to match.
	if data := timestampValidator.TimestampedData(); data != nil {
		t.Errorf("TimestampedData() = %v, want nil (none was set)", data)
	}
}

func TestDetachedTimestampValidator_OriginalDocumentsUnsupported(t *testing.T) {
	doc := loadTestTimestampDocument(t)
	timestampValidator := NewDetachedTimestampValidator(doc)

	assertPanics(t, func() { timestampValidator.OriginalDocuments("sig-1") })
	assertPanics(t, func() { timestampValidator.OriginalDocumentsForSignature(nil) })
}

func assertPanics(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("expected a panic, got none")
		}
	}()
	fn()
}
