// Ported from dss-enumerations/.../LoTEServiceTypeIdentifier.java (DSS 6.5.RC1).
package enumerations

// LoTEServiceTypeIdentifier represents a LoTE service type identifier.
type LoTEServiceTypeIdentifier interface {
	UriBasedEnum

	// Label gets user-friendly identifier.
	Label() string
}

// LoTEServiceTypeIdentifierFromURI returns a LoTEServiceTypeIdentifier for
// the given URI, or nil if none of the registered LoTELoaders resolve it.
func LoTEServiceTypeIdentifierFromURI(uri string) LoTEServiceTypeIdentifier {
	for _, loader := range loTELoaders() {
		if identifier := loader.ServiceTypeIdentifierFromURI(uri); identifier != nil {
			return identifier
		}
	}
	return nil
}
