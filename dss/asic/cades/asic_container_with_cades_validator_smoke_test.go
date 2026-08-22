// Smoke test for the ASiC-CAdES container validator tree
// (asic_container_with_cades_analyzer.go, asic_container_with_cades_analyzer_factory.go,
// asic_container_with_cades_validator.go, asic_container_with_cades_validator_factory.go,
// asic_with_cades_diagnostic_data_builder.go, asic_with_cades_timestamp_analyzer.go,
// asic_with_cades_timestamp_validator.go).
//
// Exercises the full pipeline end to end - SignedDocumentValidator.fromDocument dispatch, the
// ASiC container extraction, the nested CAdES signature analyzers,
// ASiCContainerDiagnosticDataBuilder/ASiCWithCAdESDiagnosticDataBuilder, the default validation
// policy, and the executor/report-builder tree - against real signed ASiC-CAdES fixtures already
// committed under testdata/upstream.
//
// Test-local ServiceLoader wiring: same as cades/cms_document_validator_smoke_test.go and
// xades/xml_document_validator_smoke_test.go - registers the frozen dss-policy-jaxb equivalent
// factories that PORTING.md forbids editing.
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

// asicCadesFixturePath resolves relFromPkg — a path relative to this
// package's own directory (e.g. "../testdata/upstream/x") pointing at a
// fixture vendored under the parent asic package's own testdata/upstream
// tree. Most of those fixtures still ship in-package; the heavier ones moved
// into the external corpus/ tree under that parent package's mirror, so a
// local miss is retried there: "../testdata/x" becomes the module-root-
// relative "asic/testdata/x" (this package's own module path, "asic/cades",
// joined with the "../" climb).
func asicCadesFixturePath(t *testing.T, relFromPkg string) string {
	t.Helper()
	if _, err := os.Stat(relFromPkg); err == nil {
		return relFromPkg
	}
	moduleRel := filepath.Clean(filepath.Join("asic/cades", relFromPkg))
	return corpustest.RootPath(t, moduleRel)
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

func TestASiCContainerWithCAdESValidator_Smoke(t *testing.T) {
	tests := []struct {
		name           string
		file           string
		wantSignatures int
	}{
		{"asics-onefile", "../testdata/upstream/dss-asic-cades/src/test/resources/validation/onefile-ok.asics", 1},
		{"asics-counter-signature", "../testdata/upstream/dss-asic-cades/src/test/resources/validation/containerWithCounterSig.asics", 2},
		// Exercises the ASiCWithCAdESTimestampAnalyzer container-timestamp path
		// (asic_with_cades_timestamp_analyzer.go) via a real ASiC-E LTA fixture with an
		// embedded archive (container) timestamp.
		{"asice-lta", "../testdata/upstream/dss-asic-cades/src/test/resources/validation/ASICE-CAdES-BpLTA-2-BpLTA-DSS5.4.asice", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := model.NewFileDocument(asicCadesFixturePath(t, tt.file))
			if err != nil {
				t.Fatalf("NewFileDocument(%s): %v", tt.file, err)
			}

			if !NewASiCContainerWithCAdESValidatorFactory().IsSupported(doc) {
				t.Fatalf("ASiCContainerWithCAdESValidatorFactory.IsSupported() = false for %s, want true", tt.file)
			}

			// Exercise the format-dispatch registry (SignedDocumentValidator.fromDocument's
			// port), not just the direct constructor, since that is how a real caller with an
			// unknown document reaches this validator.
			documentValidator, err := dssvalidation.SignedDocumentValidatorFromDocument(doc)
			if err != nil {
				t.Fatalf("SignedDocumentValidatorFromDocument: %v", err)
			}
			validator, ok := documentValidator.(*ASiCContainerWithCAdESValidator)
			if !ok {
				t.Fatalf("SignedDocumentValidatorFromDocument() = %T, want *ASiCContainerWithCAdESValidator", documentValidator)
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
			if diagnosticData.ContainerInfo == nil {
				t.Error("diagnostic data ContainerInfo = nil, want a built XmlContainerInfo")
			}
		})
	}
}

func TestASiCContainerWithCAdESValidatorFactory_RegistersItself(t *testing.T) {
	doc, err := model.NewFileDocument(asicCadesFixturePath(t, "../testdata/upstream/dss-asic-cades/src/test/resources/validation/onefile-ok.asics"))
	if err != nil {
		t.Fatalf("NewFileDocument: %v", err)
	}

	found := false
	for _, factory := range dssvalidation.DocumentValidatorFactories() {
		if _, ok := factory.(*ASiCContainerWithCAdESValidatorFactory); ok {
			found = true
			if !factory.IsSupported(doc) {
				t.Error("registered ASiCContainerWithCAdESValidatorFactory.IsSupported() = false, want true")
			}
		}
	}
	if !found {
		t.Fatal("no *ASiCContainerWithCAdESValidatorFactory found in the registered DocumentValidatorFactories - init() registration missing")
	}
}
