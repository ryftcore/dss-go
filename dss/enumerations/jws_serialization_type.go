// Ported from dss-enumerations/.../JWSSerializationType.java (DSS 6.5.RC1).
package enumerations

// JWSSerializationType represents JWS types defined in RFC 7515, 3. JSON
// Web Signature (JWS) Overview.
type JWSSerializationType string

const (
	// JWSSerializationTypeCompactSerialization is the JWS Compact
	// Serialization (RFC 7515, 3.1).
	JWSSerializationTypeCompactSerialization JWSSerializationType = "COMPACT_SERIALIZATION"
	// JWSSerializationTypeJSONSerialization is the general JWS JSON
	// Serialization Syntax (RFC 7515, 7.2.1).
	JWSSerializationTypeJSONSerialization JWSSerializationType = "JSON_SERIALIZATION"
	// JWSSerializationTypeFlattenedJSONSerialization is the Flattened
	// JWS JSON Serialization Syntax (RFC 7515, 7.2.2).
	JWSSerializationTypeFlattenedJSONSerialization JWSSerializationType = "FLATTENED_JSON_SERIALIZATION"
)

// JWSSerializationTypeValues returns all constants in declaration order.
func JWSSerializationTypeValues() []JWSSerializationType {
	return []JWSSerializationType{
		JWSSerializationTypeCompactSerialization,
		JWSSerializationTypeJSONSerialization,
		JWSSerializationTypeFlattenedJSONSerialization,
	}
}
