// Ported from specs-validation-report/src/main/java/eu/europa/esig/validationreport/enums/SignatureValidationProcessID.java
// (DSS 6.5.RC1). The concrete type lives in jaxb/jaxb_enums.go and is
// re-exported here as an alias - see doc.go's "Enum and parser re-exports".
package validationreport

import "github.com/ryftcore/dss-go/dss/validationreport/jaxb"

// SignatureValidationProcessID defines SignatureValidationProcessID.
// Implements enumerations.UriBasedEnum.
type SignatureValidationProcessID = jaxb.SignatureValidationProcessID

// The SignatureValidationProcessID values, re-exported from jaxb.
const (
	SignatureValidationProcessIDBasic = jaxb.SignatureValidationProcessIDBasic
	SignatureValidationProcessIDLTVM  = jaxb.SignatureValidationProcessIDLTVM
	SignatureValidationProcessIDLTA   = jaxb.SignatureValidationProcessIDLTA
)

// SignatureValidationProcessIDValues returns all SignatureValidationProcessID
// constants in declaration order.
func SignatureValidationProcessIDValues() []SignatureValidationProcessID {
	return jaxb.SignatureValidationProcessIDValues()
}

// SignatureValidationProcessIDValueOf returns the SignatureValidationProcessID
// matching the given Java enum name.
func SignatureValidationProcessIDValueOf(name string) (SignatureValidationProcessID, error) {
	return jaxb.SignatureValidationProcessIDValueOf(name)
}
