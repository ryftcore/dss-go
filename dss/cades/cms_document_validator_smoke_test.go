// Smoke test for the CAdES validation trio (cms_document_validator.go,
// cms_document_validator_factory.go, cades_diagnostic_data_builder.go).
//
// Exercises the full pipeline end to end - SignedDocumentValidatorFromDocument dispatch, the
// CMSDocumentAnalyzer, CAdESDiagnosticDataBuilder's virtual-dispatch override of
// BuildDetachedXmlSignature/BuildDetachedXmlTimestamp, the default validation policy, and the
// executor/report-builder tree - against real signed CAdES fixtures already committed under
// testdata/upstream/validation.
//
// Test-local ServiceLoader wiring: dss-policy-jaxb's default-policy factory and dss-policy-jaxb's
// XML cryptographic-suite factory are separate modules Java loads via the classpath/ServiceLoader,
// mirrored here exactly as validation/policy/validation_policy_loader_test.go's own init() does
// (those factory packages are not edited here).
package cades

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

func init() {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())
	validationpolicy.RegisterCryptographicSuiteFactory(cryptoxml.NewCryptographicSuiteXmlFactory())
}

// cadesFixturePath resolves rel (relative to this package's testdata/) to a
// real file: most fixtures ship in-package, a few larger ones live in the
// external corpus/ instead, so a local miss falls through to corpustest.
func cadesFixturePath(t *testing.T, rel string) string {
	t.Helper()
	local := filepath.Join("testdata", rel)
	if _, err := os.Stat(local); err == nil {
		return local
	}
	return corpustest.Path(t, rel)
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

func TestCMSDocumentValidator_Smoke(t *testing.T) {
	tests := []struct {
		name           string
		file           string
		wantSignatures int
	}{
		{"enveloping-single", "upstream/validation/cades-bes-signeddata-enveloping.p7m", 1},
		{"detached-single", "upstream/validation/cades-bes-signeddata-detached.p7s", 1},
		{"baseline-b", "upstream/validation/Signature-C-B-B-8.p7m", 1},
		{"counter-signature", "upstream/validation/counterSig.p7m", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := model.NewFileDocument(cadesFixturePath(t, tt.file))
			if err != nil {
				t.Fatalf("NewFileDocument(%s): %v", tt.file, err)
			}

			if !NewCMSDocumentValidatorFactory().IsSupported(doc) {
				t.Fatalf("CMSDocumentValidatorFactory.IsSupported() = false for %s, want true", tt.file)
			}

			// Exercise the format-dispatch registry (SignedDocumentValidator.fromDocument's
			// port), not just the direct constructor, since that is how a real caller with an
			// unknown document reaches this validator.
			documentValidator, err := dssvalidation.SignedDocumentValidatorFromDocument(doc)
			if err != nil {
				t.Fatalf("SignedDocumentValidatorFromDocument: %v", err)
			}
			validator, ok := documentValidator.(*CMSDocumentValidator)
			if !ok {
				t.Fatalf("SignedDocumentValidatorFromDocument() = %T, want *CMSDocumentValidator", documentValidator)
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

			// The diagnostic data must have gone through CAdESDiagnosticDataBuilder's override,
			// not the base SignedDocumentDiagnosticDataBuilder default - the fixtures carry
			// SignerInformationStore data built only by the CAdES-specific override.
			diagnosticData := reports.GetDiagnosticDataJaxb()
			if diagnosticData == nil || len(diagnosticData.Signatures.All()) == 0 {
				t.Fatal("expected the diagnostic data to contain at least one XmlSignature")
			}
		})
	}
}

func TestCMSDocumentValidatorFactory_RegistersItself(t *testing.T) {
	doc, err := model.NewFileDocument(cadesFixturePath(t, "upstream/validation/cades-bes-signeddata-enveloping.p7m"))
	if err != nil {
		t.Fatalf("NewFileDocument: %v", err)
	}

	found := false
	for _, factory := range dssvalidation.DocumentValidatorFactories() {
		if _, ok := factory.(*CMSDocumentValidatorFactory); ok {
			found = true
			if !factory.IsSupported(doc) {
				t.Error("registered CMSDocumentValidatorFactory.IsSupported() = false, want true")
			}
		}
	}
	if !found {
		t.Fatal("no *CMSDocumentValidatorFactory found in the registered DocumentValidatorFactories - init() registration missing")
	}
}
