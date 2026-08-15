// Ported from specs-validation-report/src/main/java/eu/europa/esig/validationreport/enums/ObjectType.java
// (DSS 6.5.RC1). The concrete type lives in jaxb/jaxb_enums.go and is
// re-exported here as an alias - see doc.go's "Enum and parser re-exports".
package validationreport

import "github.com/utain/esig/dss/validationreport/jaxb"

// ObjectType defines object types. Implements enumerations.UriBasedEnum.
type ObjectType = jaxb.ObjectType

// The ObjectType values, re-exported from jaxb.
const (
	ObjectType_CERTIFICATE     = jaxb.ObjectType_CERTIFICATE
	ObjectType_CRL             = jaxb.ObjectType_CRL
	ObjectType_OCSP_RESPONSE   = jaxb.ObjectType_OCSP_RESPONSE
	ObjectType_TIMESTAMP       = jaxb.ObjectType_TIMESTAMP
	ObjectType_EVIDENCE_RECORD = jaxb.ObjectType_EVIDENCE_RECORD
	ObjectType_PUBLIC_KEY      = jaxb.ObjectType_PUBLIC_KEY
	ObjectType_SIGNED_DATA     = jaxb.ObjectType_SIGNED_DATA
	ObjectType_OTHER           = jaxb.ObjectType_OTHER
)

// ObjectTypeValues returns all ObjectType constants in declaration order.
func ObjectTypeValues() []ObjectType { return jaxb.ObjectTypeValues() }

// ObjectTypeValueOf returns the ObjectType matching the given Java enum name.
func ObjectTypeValueOf(name string) (ObjectType, error) { return jaxb.ObjectTypeValueOf(name) }
