// Ported from dss-enumerations/.../ObjectIdentifier.java (DSS 6.5.RC1).
package enumerations

// ObjectIdentifier represents an identifier type with the following properties:
// identifier (URI for XAdES and/or OID for CAdES), identifier qualifier
// (URI or URN encoding), description, and document references.
type ObjectIdentifier interface {
	OidAndUriBasedEnum
	OidDescription

	// DocumentationReferences returns a slice of URI-based references.
	// NOTE: used in XAdES.
	DocumentationReferences() []string

	// Qualifier returns the Object Identifier Qualifier.
	Qualifier() ObjectIdentifierQualifier
}
