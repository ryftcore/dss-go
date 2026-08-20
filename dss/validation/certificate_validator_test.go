package validation

import (
	"crypto/x509"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/utain/esig/dss/model"
	dsspolicy "github.com/utain/esig/dss/policy"
	spivalidation "github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/validation/executor"
	validationpolicy "github.com/utain/esig/dss/validation/policy"
)

// loadTestCertificate reads one certificate of the corpus the frozen spi
// package's own KAT test uses, so this test ships no private fixture.
func loadTestCertificate(t *testing.T) *model.CertificateToken {
	t.Helper()
	der, err := os.ReadFile(filepath.Join("..", "spi", "testdata", "certificate_extensions", "cert_00.der"))
	if err != nil {
		t.Fatalf("reading the certificate: %v", err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parsing the certificate: %v", err)
	}
	token, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("building the certificate token: %v", err)
	}
	return token
}

// TestCertificateValidatorWiring pins what CertificateValidator overrides:
// the default policy path, the process executor it defaults to, and that the
// executor is handed the identifier of the certificate under validation.
func TestCertificateValidatorWiring(t *testing.T) {
	token := loadTestCertificate(t)
	validator := CertificateValidatorFromCertificate(token)

	if got := validator.DefaultValidationPolicyPath(); got != "/policy/certificate-constraint.xml" {
		t.Errorf("DefaultValidationPolicyPath() = %q, want the upstream constant", got)
	}
	if _, ok := validator.DefaultProcessExecutor().(*executor.DefaultCertificateProcessExecutor); !ok {
		t.Errorf("DefaultProcessExecutor() = %T, want *executor.DefaultCertificateProcessExecutor",
			validator.DefaultProcessExecutor())
	}

	provided := validator.ProvideProcessExecutorInstance()
	if provided == nil {
		t.Fatal("ProvideProcessExecutorInstance() returned nil")
	}
	// Java's provideProcessExecutorInstance() caches the executor and sets the
	// certificate id on every call.
	if again := validator.ProvideProcessExecutorInstance(); again != provided {
		t.Error("ProvideProcessExecutorInstance() did not cache the executor")
	}
	concrete := provided.(*executor.DefaultCertificateProcessExecutor)
	if concrete.CertificateId != token.DSSIDAsString() {
		t.Errorf("certificate id = %q, want %q", concrete.CertificateId, token.DSSIDAsString())
	}
}

// TestCertificateValidatorFromCertificateNil covers Java's
// Objects.requireNonNull("The certificate is missing").
func TestCertificateValidatorFromCertificateNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic for a nil certificate")
		} else if r != "The certificate is missing" {
			t.Errorf("panic = %v, want Java's message", r)
		}
	}()
	CertificateValidatorFromCertificate(nil)
}

// TestCertificateValidatorAssertConfigurationValid covers both requireNonNull
// guards Java's assertConfigurationValid() chain performs.
func TestCertificateValidatorAssertConfigurationValid(t *testing.T) {
	validator := CertificateValidatorFromCertificate(loadTestCertificate(t))
	err := validator.AssertConfigurationValid()
	if err == nil || !strings.Contains(err.Error(), "CertificateVerifier is not defined") {
		t.Errorf("error = %v, want the missing-CertificateVerifier message", err)
	}

	validator.SetCertificateVerifier(spivalidation.NewCommonCertificateVerifier())
	if err := validator.AssertConfigurationValid(); err != nil {
		t.Errorf("AssertConfigurationValid() = %v, want nil once the verifier is set", err)
	}
}

// TestCertificateValidatorValidate runs the whole certificate validation, the
// way Java's CertificateValidator.validate() does: prepare the validation
// context, build the diagnostic data, and let DefaultCertificateProcessExecutor
// produce the reports.
func TestCertificateValidatorValidate(t *testing.T) {
	validationpolicy.RegisterValidationPolicyFactory(dsspolicy.NewEtsiValidationPolicyFactory())

	token := loadTestCertificate(t)
	validator := CertificateValidatorFromCertificate(token)
	validator.SetCertificateVerifier(spivalidation.NewCommonCertificateVerifier())
	validator.SetValidationTime(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC))

	reports, err := validator.Validate()
	if err != nil {
		t.Fatalf("Validate(): %v", err)
	}
	if reports == nil {
		t.Fatal("Validate() returned no reports")
	}
	simpleReport := reports.GetSimpleReportJaxb()
	if simpleReport == nil || simpleReport.Certificate == nil {
		t.Fatal("the simple certificate report carries no certificate")
	}
	if simpleReport.Certificate.Id != token.DSSIDAsString() {
		t.Errorf("report certificate id = %q, want %q", simpleReport.Certificate.Id, token.DSSIDAsString())
	}
	if simpleReport.ValidationPolicy == nil || simpleReport.ValidationPolicy.PolicyName == nil {
		t.Error("the simple certificate report carries no validation policy name")
	}
	if detailedReport := reports.GetDetailedReportJaxb(); detailedReport == nil ||
		len(detailedReport.BasicBuildingBlocks) == 0 {
		t.Error("the detailed report carries no basic building blocks")
	}
	// The default policy path resolves against the embedded resources; see
	// abstract_certificate_validator.go's header.
	if _, err := validator.FromDefaultCertificateValidationPolicyLoader().Create(), error(nil); err != nil {
		t.Errorf("loading the embedded default certificate policy: %v", err)
	}
}
