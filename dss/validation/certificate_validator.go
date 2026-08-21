// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/CertificateValidator.java
// (DSS 6.5.RC1).

package validation

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/model"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/validation/executor"
	"github.com/ryftcore/dss-go/dss/validation/reports"
)

// certificateValidationPolicyLocation is the path for the default certificate
// validation policy. Port of the private static final
// CERTIFICATE_VALIDATION_POLICY_LOCATION.
const certificateValidationPolicyLocation = "/policy/certificate-constraint.xml"

// CertificateValidator validates a CertificateToken. Port of
// CertificateValidator (extends AbstractCertificateValidator<CertificateReports,
// CertificateProcessExecutor>).
type CertificateValidator struct {
	AbstractCertificateValidator[*reports.CertificateReports, executor.CertificateProcessExecutor]

	// token is the certificateToken to be validated. Port of the private
	// final token field.
	token *model.CertificateToken
}

// compile-time interface assertion: the concrete validator supplies every
// member AbstractCertificateValidator dispatches through its overrides.
var _ AbstractCertificateValidatorOverrides[executor.CertificateProcessExecutor] = (*CertificateValidator)(nil)

// newCertificateValidator is the port of the private
// CertificateValidator(CertificateToken) constructor.
func newCertificateValidator(token *model.CertificateToken) *CertificateValidator {
	v := &CertificateValidator{
		AbstractCertificateValidator: NewAbstractCertificateValidator[*reports.CertificateReports, executor.CertificateProcessExecutor](),
		token:                        token,
	}
	v.InitAbstractCertificateValidator(v)
	return v
}

// CertificateValidatorFromCertificate creates a CertificateValidator from a
// certificateToken. Port of the static fromCertificate(CertificateToken).
//
// Panics when token is nil (Java's Objects.requireNonNull("The certificate is
// missing")).
func CertificateValidatorFromCertificate(token *model.CertificateToken) *CertificateValidator {
	if token == nil {
		panic("The certificate is missing")
	}
	return newCertificateValidator(token)
}

// DefaultValidationPolicyPath is the port of the overridden
// getDefaultValidationPolicyPath().
func (v *CertificateValidator) DefaultValidationPolicyPath() string {
	return certificateValidationPolicyLocation
}

// PrepareValidationContext is the port of the overridden
// prepareValidationContext(CertificateVerifier).
func (v *CertificateValidator) PrepareValidationContext(
	certificateVerifier spivalidation.CertificateVerifier) spivalidation.ValidationContext {
	svc := v.AbstractCertificateValidator.PrepareValidationContext(certificateVerifier)
	svc.AddCertificateTokenForVerification(v.token)
	return svc
}

// ProvideProcessExecutorInstance gets the CertificateProcessExecutor. Port of
// the protected provideProcessExecutorInstance().
func (v *CertificateValidator) ProvideProcessExecutorInstance() executor.CertificateProcessExecutor {
	if v.ProcessExecutor == nil {
		v.ProcessExecutor = v.DefaultProcessExecutor()
	}
	v.ProcessExecutor.SetCertificateId(v.IdentifierProvider.IDAsString(v.token))
	return v.ProcessExecutor
}

// DefaultProcessExecutor is the port of the overridden
// getDefaultProcessExecutor().
func (v *CertificateValidator) DefaultProcessExecutor() executor.CertificateProcessExecutor {
	return executor.NewDefaultCertificateProcessExecutor()
}

// AssertConfigurationValid is the port of the overridden
// assertConfigurationValid().
func (v *CertificateValidator) AssertConfigurationValid() error {
	if err := v.AbstractCertificateValidator.AssertConfigurationValid(); err != nil {
		return err
	}
	if v.token == nil {
		return fmt.Errorf("Certificate token is not provided to the validator")
	}
	return nil
}
