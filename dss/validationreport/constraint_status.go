// Ported from specs-validation-report/src/main/java/eu/europa/esig/validationreport/enums/ConstraintStatus.java
// (DSS 6.5.RC1). The concrete type lives in jaxb/jaxb_enums.go and is
// re-exported here as an alias - see doc.go's "Enum and parser re-exports".
package validationreport

import "github.com/ryftcore/dss-go/dss/validationreport/jaxb"

// ConstraintStatus defines the ConstraintStatus type. Implements
// enumerations.UriBasedEnum.
type ConstraintStatus = jaxb.ConstraintStatus

// The ConstraintStatus values, re-exported from jaxb.
const (
	ConstraintStatus_APPLIED    = jaxb.ConstraintStatus_APPLIED
	ConstraintStatus_DISABLED   = jaxb.ConstraintStatus_DISABLED
	ConstraintStatus_OVERRIDDEN = jaxb.ConstraintStatus_OVERRIDDEN
)

// ConstraintStatusValues returns all ConstraintStatus constants in declaration order.
func ConstraintStatusValues() []ConstraintStatus { return jaxb.ConstraintStatusValues() }

// ConstraintStatusValueOf returns the ConstraintStatus matching the given Java enum name.
func ConstraintStatusValueOf(name string) (ConstraintStatus, error) {
	return jaxb.ConstraintStatusValueOf(name)
}
