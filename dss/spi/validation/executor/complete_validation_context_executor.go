// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/executor/CompleteValidationContextExecutor.java (DSS 6.5.RC1).
package executor

import "github.com/ryftcore/dss-go/dss/spi/validation"

// CompleteValidationContextExecutor executes complete validation of the Context,
// including running of all checks with the alerts processing specified in CertificateVerifier.
type CompleteValidationContextExecutor struct{}

// CompleteValidationContextExecutorInstance is the singleton instance.
// Port of the public static final INSTANCE field.
var CompleteValidationContextExecutorInstance = &CompleteValidationContextExecutor{}

// Validate requires validationContext to be non-nil and a *validation.SignatureValidationContext
// (Objects.requireNonNull and the instanceof check both panic with Java's messages, matching
// the UnsupportedOperationException/NullPointerException convention used throughout this port),
// then runs the full validation + alerting sequence.
func (e *CompleteValidationContextExecutor) Validate(validationContext validation.Context) {
	sigValidationContext := assertValidationContextSupported(validationContext)

	validationContext.Validate()
	assertSignaturesValid(sigValidationContext)
}

// assertValidationContextSupported ports the private static assertValidationContextSupported;
// it additionally returns the asserted *SignatureValidationContext, which the Java version's
// caller obtains from a separate cast at the call site.
func assertValidationContextSupported(validationContext validation.Context) *validation.SignatureValidationContext {
	if validationContext == nil {
		panic("ValidationContext cannot be null!")
	}
	sigValidationContext, ok := validationContext.(*validation.SignatureValidationContext)
	if !ok {
		panic("CompleteValidationContextExecutor supports only SignatureValidationContext class type!")
	}
	return sigValidationContext
}

// assertSignaturesValid ports the private assertSignaturesValid.
func assertSignaturesValid(sigValidationContext *validation.SignatureValidationContext) {
	validationAlerter := validation.NewSignatureValidationAlerter(sigValidationContext)
	validationAlerter.AssertAllTimestampsValid()
	validationAlerter.AssertAllRequiredRevocationDataPresent()
	validationAlerter.AssertAllPOECoveredByRevocationData()

	validationAlerter.AssertAllSignaturesAreYetValid()
	validationAlerter.AssertAllSignaturesNotExpired()
	validationAlerter.AssertAllSignatureCertificatesNotRevoked()
	validationAlerter.AssertAllSignatureCertificateHaveFreshRevocationData()
}

// compile-time interface assertion.
var _ ValidationContextExecutor = (*CompleteValidationContextExecutor)(nil)
