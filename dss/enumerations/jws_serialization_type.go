// Ported from dss-enumerations/.../JWSSerializationType.java (DSS 6.5.RC1).
package enumerations

// JWSSerializationType represents JWS types defined in RFC 7515, 3. JSON
// Web Signature (JWS) Overview.
type JWSSerializationType string

const (
	// JWSSerializationType_COMPACT_SERIALIZATION is the JWS Compact
	// Serialization (RFC 7515, 3.1).
	JWSSerializationType_COMPACT_SERIALIZATION JWSSerializationType = "COMPACT_SERIALIZATION"
	// JWSSerializationType_JSON_SERIALIZATION is the general JWS JSON
	// Serialization Syntax (RFC 7515, 7.2.1).
	JWSSerializationType_JSON_SERIALIZATION JWSSerializationType = "JSON_SERIALIZATION"
	// JWSSerializationType_FLATTENED_JSON_SERIALIZATION is the Flattened
	// JWS JSON Serialization Syntax (RFC 7515, 7.2.2).
	JWSSerializationType_FLATTENED_JSON_SERIALIZATION JWSSerializationType = "FLATTENED_JSON_SERIALIZATION"
)

// JWSSerializationTypeValues returns all constants in declaration order.
func JWSSerializationTypeValues() []JWSSerializationType {
	return []JWSSerializationType{
		JWSSerializationType_COMPACT_SERIALIZATION,
		JWSSerializationType_JSON_SERIALIZATION,
		JWSSerializationType_FLATTENED_JSON_SERIALIZATION,
	}
}
