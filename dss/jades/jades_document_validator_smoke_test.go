// Smoke test for the JAdES validator tree (abstract_jws_document_validator.go,
// jades_diagnostic_data_builder.go, jades_document_validator_factory.go,
// jws_compact_document_validator.go, jws_serialization_document_validator.go).
//
// Exercises the full pipeline end to end - SignedDocumentValidator.fromDocument dispatch, both
// JWS analyzers (compact and serialization/flattened), JAdESDiagnosticDataBuilder, the default
// validation policy, and the executor/report-builder tree - against real signed JAdES fixtures
// already committed under testdata/upstream.
//
// Test-local ServiceLoader wiring: same as cades/cms_document_validator_smoke_test.go and
// xades/xml_document_validator_smoke_test.go - registers the frozen dss-policy-jaxb equivalent
// factories that PORTING.md forbids editing.
package jades

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ryftcore/dss-go/dss/alert"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
	dsspolicy "github.com/ryftcore/dss-go/dss/policy"
	cryptoxml "github.com/ryftcore/dss-go/dss/policy/crypto/xml"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	dssvalidation "github.com/ryftcore/dss-go/dss/validation"
	validationpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
)

// jadesFixturePath resolves rel (relative to this package's testdata/) to a
// real file: most fixtures ship in-package, a few larger ones live in the
// external corpus/ instead, so a local miss falls through to corpustest.
func jadesFixturePath(t *testing.T, rel string) string {
	t.Helper()
	local := filepath.Join("testdata", rel)
	if _, err := os.Stat(local); err == nil {
		return local
	}
	return corpustest.Path(t, rel)
}

func init() {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())
	validationpolicy.RegisterCryptographicSuiteFactory(cryptoxml.NewCryptographicSuiteXmlFactory())
}

// permissiveCertificateVerifier builds a CertificateVerifier that silently tolerates the missing
// revocation/AIA data and stale certificates in offline test fixtures, rather than raising the
// exception-by-default alerts CommonCertificateVerifier ships with.
func permissiveCertificateVerifier() *validation.CommonCertificateVerifier {
	v := validation.NewCommonCertificateVerifierSimple(true)
	v.SetAlertOnMissingRevocationData(alert.NewSilentOnStatusAlert())
	v.SetAlertOnRevokedCertificate(alert.NewSilentOnStatusAlert())
	v.SetAlertOnInvalidSignature(alert.NewSilentOnStatusAlert())
	v.SetAlertOnInvalidTimestamp(alert.NewSilentOnStatusAlert())
	v.SetAlertOnUncoveredPOE(alert.NewSilentOnStatusAlert())
	v.SetAlertOnExpiredCertificate(alert.NewSilentOnStatusAlert())
	v.SetAlertOnNotYetValidCertificate(alert.NewSilentOnStatusAlert())
	return v
}

func TestJAdESDocumentValidator_Smoke(t *testing.T) {
	tests := []struct {
		name           string
		file           string
		wantSignatures int
	}{
		{"compact-with-certified", "upstream/validation/jades-with-certified.json", 1},
		{"serialization-lta", "upstream/validation/jades-lta.json", 1},
		{"serialization-with-xvals", "upstream/validation/jades-b-with-xvals.json", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := model.NewFileDocument(jadesFixturePath(t, tt.file))
			if err != nil {
				t.Fatalf("NewFileDocument(%s): %v", tt.file, err)
			}

			if !NewJAdESDocumentValidatorFactory().IsSupported(doc) {
				t.Fatalf("JAdESDocumentValidatorFactory.IsSupported() = false for %s, want true", tt.file)
			}

			// Exercise the format-dispatch registry (SignedDocumentValidator.fromDocument's
			// port), not just the direct constructor, since that is how a real caller with an
			// unknown document reaches this validator.
			documentValidator, err := dssvalidation.SignedDocumentValidatorFromDocument(doc)
			if err != nil {
				t.Fatalf("SignedDocumentValidatorFromDocument: %v", err)
			}

			validator, ok := documentValidator.(interface {
				dssvalidation.SignedDocumentValidator
				JwsJsonSerializationObject() *JWSJsonSerializationObject
			})
			if !ok {
				t.Fatalf("SignedDocumentValidatorFromDocument() = %T, want a JAdES validator", documentValidator)
			}

			validator.SetCertificateVerifier(permissiveCertificateVerifier())
			validator.SetValidationLevel(enumerations.ValidationLevelBasicSignatures)
			validator.SetLocale("en")

			reports, err := validator.ValidateDocument()
			if err != nil {
				t.Fatalf("ValidateDocument: %v", err)
			}
			if reports == nil {
				t.Fatal("ValidateDocument returned a nil Reports")
			}

			simpleReport := reports.GetSimpleReport()
			if got := simpleReport.GetSignaturesCount(); got != tt.wantSignatures {
				t.Errorf("GetSignaturesCount() = %d, want %d", got, tt.wantSignatures)
			}
			for _, id := range simpleReport.GetSignatureIdList() {
				if id == "" {
					t.Error("GetSignatureIdList() contains an empty signature id")
				}
			}

			diagnosticData := reports.GetDiagnosticDataJaxb()
			if diagnosticData == nil || len(diagnosticData.Signatures.All()) == 0 {
				t.Fatal("expected the diagnostic data to contain at least one XmlSignature")
			}
		})
	}
}

func TestJAdESDocumentValidatorFactory_RegistersItself(t *testing.T) {
	doc, err := model.NewFileDocument(jadesFixturePath(t, "upstream/validation/jades-lta.json"))
	if err != nil {
		t.Fatalf("NewFileDocument: %v", err)
	}

	found := false
	for _, factory := range dssvalidation.DocumentValidatorFactories() {
		if _, ok := factory.(*JAdESDocumentValidatorFactory); ok {
			found = true
			if !factory.IsSupported(doc) {
				t.Error("registered JAdESDocumentValidatorFactory.IsSupported() = false, want true")
			}
		}
	}
	if !found {
		t.Fatal("no *JAdESDocumentValidatorFactory found in the registered DocumentValidatorFactories - init() registration missing")
	}
}

func TestJWSCompactDocumentValidator_DirectConstructor(t *testing.T) {
	doc, err := model.NewFileDocument(jadesFixturePath(t, "upstream/validation/jades-with-certified.json"))
	if err != nil {
		t.Fatalf("NewFileDocument: %v", err)
	}
	validator := NewJWSCompactDocumentValidatorFromDocument(doc)
	validator.SetCertificateVerifier(permissiveCertificateVerifier())
	if validator.JwsJsonSerializationObject() == nil {
		t.Fatal("JwsJsonSerializationObject() = nil")
	}
}

func TestJWSSerializationDocumentValidator_DirectConstructor(t *testing.T) {
	doc, err := model.NewFileDocument(jadesFixturePath(t, "upstream/validation/jades-lta.json"))
	if err != nil {
		t.Fatalf("NewFileDocument: %v", err)
	}
	validator := NewJWSSerializationDocumentValidatorFromDocument(doc)
	validator.SetCertificateVerifier(permissiveCertificateVerifier())
	if validator.JwsJsonSerializationObject() == nil {
		t.Fatal("JwsJsonSerializationObject() = nil")
	}
}
