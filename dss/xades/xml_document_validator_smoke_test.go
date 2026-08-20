// Smoke test for the un-gated XAdES validator pair (xml_document_validator.go,
// xml_document_validator_factory.go), now that the phase 8 validation engine (dss/validation,
// dss/validation/executor, dss/validation/policy, dss/simplereport, dss/policy) has landed and
// their `phase8` build tags were removed.
//
// Exercises the full pipeline end to end - SignedDocumentValidator.fromDocument dispatch, the
// XMLDocumentAnalyzer, the base SignedDocumentDiagnosticDataBuilder (XAdES ships no
// format-specific override - see xml_document_validator.go's file header), the default validation
// policy, and the executor/report-builder tree - against real signed XAdES fixtures already
// committed under testdata/upstream.
//
// Test-local ServiceLoader wiring: same as cades/cms_document_validator_smoke_test.go and
// validation/policy/validation_policy_loader_test.go - registers the frozen dss-policy-jaxb
// equivalent factories that PORTING.md forbids editing.
package xades

import (
	"testing"

	"github.com/utain/esig/dss/alert"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	dsspolicy "github.com/utain/esig/dss/policy"
	cryptoxml "github.com/utain/esig/dss/policy/crypto/xml"
	"github.com/utain/esig/dss/spi/validation"
	dssvalidation "github.com/utain/esig/dss/validation"
	validationpolicy "github.com/utain/esig/dss/validation/policy"
)

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

func TestXMLDocumentValidator_Smoke(t *testing.T) {
	tests := []struct {
		name           string
		file           string
		wantSignatures int
	}{
		{"baseline-b-with-cert-values", "testdata/upstream/BaselineBWithCertificateValues.xml", 1},
		{"signature-x-at-1", "testdata/upstream/Signature-X-AT-1.xml", 1},
		{"xades-lta", "testdata/upstream/XAdESLTA.xml", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := model.NewFileDocument(tt.file)
			if err != nil {
				t.Fatalf("NewFileDocument(%s): %v", tt.file, err)
			}

			if !NewXMLDocumentValidatorFactory().IsSupported(doc) {
				t.Fatalf("XMLDocumentValidatorFactory.IsSupported() = false for %s, want true", tt.file)
			}

			// Exercise the format-dispatch registry (SignedDocumentValidator.fromDocument's
			// port), not just the direct constructor, since that is how a real caller with an
			// unknown document reaches this validator.
			documentValidator, err := dssvalidation.SignedDocumentValidatorFromDocument(doc)
			if err != nil {
				t.Fatalf("SignedDocumentValidatorFromDocument: %v", err)
			}
			validator, ok := documentValidator.(*XMLDocumentValidator)
			if !ok {
				t.Fatalf("SignedDocumentValidatorFromDocument() = %T, want *XMLDocumentValidator", documentValidator)
			}

			validator.SetCertificateVerifier(permissiveCertificateVerifier())
			validator.SetValidationLevel(enumerations.ValidationLevel_BASIC_SIGNATURES)
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

func TestXMLDocumentValidator_RootElement(t *testing.T) {
	doc, err := model.NewFileDocument("testdata/upstream/Signature-X-AT-1.xml")
	if err != nil {
		t.Fatalf("NewFileDocument: %v", err)
	}
	validator, err := NewXMLDocumentValidator(doc)
	if err != nil {
		t.Fatalf("NewXMLDocumentValidator: %v", err)
	}
	if validator.RootElement() == nil {
		t.Error("RootElement() = nil, want the parsed document root")
	}
}

func TestXMLDocumentValidatorFactory_RegistersItself(t *testing.T) {
	doc, err := model.NewFileDocument("testdata/upstream/Signature-X-AT-1.xml")
	if err != nil {
		t.Fatalf("NewFileDocument: %v", err)
	}

	found := false
	for _, factory := range dssvalidation.DocumentValidatorFactories() {
		if _, ok := factory.(*XMLDocumentValidatorFactory); ok {
			found = true
			if !factory.IsSupported(doc) {
				t.Error("registered XMLDocumentValidatorFactory.IsSupported() = false, want true")
			}
		}
	}
	if !found {
		t.Fatal("no *XMLDocumentValidatorFactory found in the registered DocumentValidatorFactories - init() registration missing")
	}
}
