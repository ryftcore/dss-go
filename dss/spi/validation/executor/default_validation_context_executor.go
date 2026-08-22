// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/executor/DefaultValidationContextExecutor.java (DSS 6.5.RC1).
package executor

import "github.com/ryftcore/dss-go/dss/spi/validation"

// DefaultValidationContextExecutor performs basic validation of a Context, including
// certificate chain building and revocation data extraction, without executing different
// validity checks.
type DefaultValidationContextExecutor struct{}

// DefaultValidationContextExecutorInstance is the singleton instance.
// Port of the public static final INSTANCE field.
var DefaultValidationContextExecutorInstance = &DefaultValidationContextExecutor{}

// Validate requires validationContext to be non-nil (Objects.requireNonNull panics with Java's
// message) then delegates to validationContext.Validate().
func (e *DefaultValidationContextExecutor) Validate(validationContext validation.Context) {
	if validationContext == nil {
		panic("ValidationContext cannot be null!")
	}
	validationContext.Validate()
}

// compile-time interface assertion.
var _ ValidationContextExecutor = (*DefaultValidationContextExecutor)(nil)
