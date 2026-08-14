// Ported from dss-enumerations/.../UriBasedEnum.java (DSS 6.5.RC1).
package enumerations

// UriBasedEnum defines an enumeration containing a URI.
type UriBasedEnum interface {
	// URI returns a URI.
	URI() string
}
