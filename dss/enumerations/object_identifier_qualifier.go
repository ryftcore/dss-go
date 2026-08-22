// Ported from dss-enumerations/.../ObjectIdentifierQualifier.java (DSS 6.5.RC1).
package enumerations

// ObjectIdentifierQualifier declares the type of the defined identifier.
// Used in XAdES:
//
//	<xsd:simpleType name="QualifierType">
//		<xsd:restriction base="xsd:string">
//			<xsd:enumeration value="OIDAsURI"/>
//			<xsd:enumeration value="OIDAsURN"/>
//		</xsd:restriction>
//	</xsd:simpleType>
type ObjectIdentifierQualifier string

const (
	// ObjectIdentifierQualifierOIDAsURI identifies object Identifier
	// encoded as URI (e.g. 'http://test/public').
	ObjectIdentifierQualifierOIDAsURI ObjectIdentifierQualifier = "OID_AS_URI"
	// ObjectIdentifierQualifierOIDAsURN identifies object Identifier
	// encoded as URN (e.g. 'urn:oid:1.2.840.113549.1.9.16.6.3').
	ObjectIdentifierQualifierOIDAsURN ObjectIdentifierQualifier = "OID_AS_URN"
)

var objectIdentifierQualifierValueTable = map[ObjectIdentifierQualifier]string{
	ObjectIdentifierQualifierOIDAsURI: "OIDAsURI",
	ObjectIdentifierQualifierOIDAsURN: "OIDAsURN",
}

// ObjectIdentifierQualifierValues returns all constants in declaration
// order.
func ObjectIdentifierQualifierValues() []ObjectIdentifierQualifier {
	return []ObjectIdentifierQualifier{
		ObjectIdentifierQualifierOIDAsURI,
		ObjectIdentifierQualifierOIDAsURN,
	}
}

// Value returns the XML value of the qualifier.
func (o ObjectIdentifierQualifier) Value() string {
	return objectIdentifierQualifierValueTable[o]
}

// ObjectIdentifierQualifierFromValue returns an ObjectIdentifierQualifier
// instance from the given value. Returns "" (zero value) if unknown,
// mirroring Java's null return.
func ObjectIdentifierQualifierFromValue(v string) ObjectIdentifierQualifier {
	for _, c := range ObjectIdentifierQualifierValues() {
		if objectIdentifierQualifierValueTable[c] == v {
			return c
		}
	}
	return ""
}

// ObjectIdentifierQualifierValueOf returns the constant matching the given
// Java enum name.
func ObjectIdentifierQualifierValueOf(name string) (ObjectIdentifierQualifier, error) {
	for _, v := range ObjectIdentifierQualifierValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &objectIdentifierQualifierInvalidValueError{name}
}

type objectIdentifierQualifierInvalidValueError struct {
	name string
}

func (e *objectIdentifierQualifierInvalidValueError) Error() string {
	return "no enum constant ObjectIdentifierQualifier." + e.name
}
