// Ported from dss-enumerations/.../OidDescription.java (DSS 6.5.RC1).
package enumerations

// OidDescription represents an OID-based property with a description.
type OidDescription interface {
	OidBasedEnum

	// Description returns the literal description of the OID.
	Description() string
}
