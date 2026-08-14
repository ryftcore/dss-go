// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/executor/SkipValidationContextExecutor.java (DSS 6.5.RC1).
package executor

import "github.com/utain/esig/dss/spi/validation"

// SkipValidationContextExecutor skips validation of the ValidationContext.
type SkipValidationContextExecutor struct{}

// SkipValidationContextExecutorInstance is the singleton instance.
// Port of the public static final INSTANCE field.
var SkipValidationContextExecutorInstance = &SkipValidationContextExecutor{}

// Validate skips validation entirely.
func (e *SkipValidationContextExecutor) Validate(validationContext validation.ValidationContext) {
	// skip
}

// compile-time interface assertion.
var _ ValidationContextExecutor = (*SkipValidationContextExecutor)(nil)
