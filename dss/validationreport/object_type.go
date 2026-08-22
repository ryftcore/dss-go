// Ported from specs-validation-report/src/main/java/eu/europa/esig/validationreport/enums/ObjectType.java
// (DSS 6.5.RC1). The concrete type lives in jaxb/jaxb_enums.go and is
// re-exported here as an alias - see doc.go's "Enum and parser re-exports".
package validationreport

import "github.com/ryftcore/dss-go/dss/validationreport/jaxb"

// ObjectType defines object types. Implements enumerations.UriBasedEnum.
type ObjectType = jaxb.ObjectType

// The ObjectType values, re-exported from jaxb.
const (
	ObjectTypeCertificate    = jaxb.ObjectTypeCertificate
	ObjectTypeCRL            = jaxb.ObjectTypeCRL
	ObjectTypeOCSPResponse   = jaxb.ObjectTypeOCSPResponse
	ObjectTypeTimestamp      = jaxb.ObjectTypeTimestamp
	ObjectTypeEvidenceRecord = jaxb.ObjectTypeEvidenceRecord
	ObjectTypePublicKey      = jaxb.ObjectTypePublicKey
	ObjectTypeSignedData     = jaxb.ObjectTypeSignedData
	ObjectTypeOther          = jaxb.ObjectTypeOther
)

// ObjectTypeValues returns all ObjectType constants in declaration order.
func ObjectTypeValues() []ObjectType { return jaxb.ObjectTypeValues() }

// ObjectTypeValueOf returns the ObjectType matching the given Java enum name.
func ObjectTypeValueOf(name string) (ObjectType, error) { return jaxb.ObjectTypeValueOf(name) }
