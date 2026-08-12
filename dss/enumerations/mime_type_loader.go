// Ported from dss-enumerations/.../MimeTypeLoader.java (DSS 6.5.RC1).
package enumerations

// MimeTypeLoader is used to load an enumeration(s) of the MimeType interface.
//
// Java used java.util.ServiceLoader to discover implementations at runtime.
// Go has no equivalent runtime discovery mechanism, so implementations
// register themselves explicitly via RegisterMimeTypeLoader.
type MimeTypeLoader interface {
	// FromMimeTypeString returns a MimeType matching the MimeType string,
	// or nil if no associated MimeType is found.
	FromMimeTypeString(mimeTypeString string) MimeType

	// FromFileExtension returns a MimeType matching the provided file
	// extension, or nil if no associated MimeType is found.
	FromFileExtension(fileExtension string) MimeType
}

// mimeTypeLoaderRegistry holds the MimeTypeLoader implementations registered
// via RegisterMimeTypeLoader, consulted in registration order — the Go
// equivalent of Java's ServiceLoader.load(MimeTypeLoader.class) iteration.
var mimeTypeLoaderRegistry []MimeTypeLoader

// RegisterMimeTypeLoader registers a MimeTypeLoader to be consulted by
// MimeTypeFromMimeTypeString, MimeTypeFromFileExtension, and
// MimeTypeFromFileName, in the order loaders are registered.
func RegisterMimeTypeLoader(l MimeTypeLoader) {
	mimeTypeLoaderRegistry = append(mimeTypeLoaderRegistry, l)
}

// mimeTypeLoaders returns the currently registered MimeTypeLoader instances.
func mimeTypeLoaders() []MimeTypeLoader {
	return mimeTypeLoaderRegistry
}
