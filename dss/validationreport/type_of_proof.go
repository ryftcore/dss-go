// Ported from specs-validation-report/src/main/java/eu/europa/esig/validationreport/enums/TypeOfProof.java
// (DSS 6.5.RC1). The concrete type lives in jaxb/jaxb_enums.go and is
// re-exported here as an alias - see doc.go's "Enum and parser re-exports".
package validationreport

import "github.com/ryftcore/dss-go/dss/validationreport/jaxb"

// TypeOfProof defines a TypeOfProof. Implements enumerations.UriBasedEnum.
type TypeOfProof = jaxb.TypeOfProof

// The TypeOfProof values, re-exported from jaxb.
const (
	TypeOfProof_VALIDATION = jaxb.TypeOfProof_VALIDATION
	TypeOfProof_PROVIDED   = jaxb.TypeOfProof_PROVIDED
	TypeOfProof_POLICY     = jaxb.TypeOfProof_POLICY
)

// TypeOfProofValues returns all TypeOfProof constants in declaration order.
func TypeOfProofValues() []TypeOfProof { return jaxb.TypeOfProofValues() }

// TypeOfProofValueOf returns the TypeOfProof matching the given Java enum name.
func TypeOfProofValueOf(name string) (TypeOfProof, error) { return jaxb.TypeOfProofValueOf(name) }
