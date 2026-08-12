// Ported from dss-enumerations/.../LoTEServiceStatus.java (DSS 6.5.RC1).
package enumerations

// LoTEServiceStatus represents a LoTE service status.
type LoTEServiceStatus interface {
	UriBasedEnum

	// Label gets user-friendly label.
	Label() string
}

// LoTEServiceStatusFromURI returns a LoTEServiceStatus for the given URI, or
// nil if none of the registered LoTELoaders resolve it.
func LoTEServiceStatusFromURI(uri string) LoTEServiceStatus {
	for _, loader := range loTELoaders() {
		if status := loader.ServiceStatusFromURI(uri); status != nil {
			return status
		}
	}
	return nil
}
