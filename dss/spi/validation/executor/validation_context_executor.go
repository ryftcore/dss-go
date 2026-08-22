// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/executor/ValidationContextExecutor.java (DSS 6.5.RC1).
package executor

import "github.com/ryftcore/dss-go/dss/spi/validation"

// ValidationContextExecutor defines a strategy for execution of a Context's
// validation.
type ValidationContextExecutor interface {
	// Validate performs validation of validationContext. Port of validate(ValidationContext).
	Validate(validationContext validation.Context)
}
