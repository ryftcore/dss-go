// Ported from dss-enumerations/.../ListType.java (DSS 6.5.RC1).
package enumerations

// ListType defines a List type.
type ListType interface {
	UriBasedEnum

	// Label gets label.
	Label() string
}

// ListTypeFromURI returns a ListType for the given URI, or nil if none of
// the registered LoTELoaders resolve it.
func ListTypeFromURI(uri string) ListType {
	for _, loader := range loTELoaders() {
		if listType := loader.ListTypeFromURI(uri); listType != nil {
			return listType
		}
	}
	return nil
}
