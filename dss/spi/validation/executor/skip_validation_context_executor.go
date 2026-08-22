// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/executor/SkipValidationContextExecutor.java (DSS 6.5.RC1).
package executor

import "github.com/ryftcore/dss-go/dss/spi/validation"

// SkipValidationContextExecutor skips validation of the Context.
type SkipValidationContextExecutor struct{}

// SkipValidationContextExecutorInstance is the singleton instance.
// Port of the public static final INSTANCE field.
var SkipValidationContextExecutorInstance = &SkipValidationContextExecutor{}

// Validate skips validation entirely.
func (e *SkipValidationContextExecutor) Validate(validationContext validation.Context) {
	// skip
}

// compile-time interface assertion.
var _ ValidationContextExecutor = (*SkipValidationContextExecutor)(nil)
