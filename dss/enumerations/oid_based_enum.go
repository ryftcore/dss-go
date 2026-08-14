// Ported from dss-enumerations/.../OidBasedEnum.java (DSS 6.5.RC1).
package enumerations

// OidBasedEnum represents an OID-based property.
type OidBasedEnum interface {
	// OID returns the OID value.
	OID() string
}
