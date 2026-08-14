// Ported from dss-enumerations/.../loader/LoTELoader.java (DSS 6.5.RC1).
//
// Files under loader/ are flattened into package enumerations (a Go
// subpackage would create an import cycle with the enumerations that depend
// on it).
package enumerations

// LoTELoader is used to load TS 119 602 LoTE related properties.
//
// Java used java.util.ServiceLoader to discover implementations at runtime.
// Go has no equivalent runtime discovery mechanism, so implementations
// register themselves explicitly via RegisterLoTELoader.
type LoTELoader interface {
	// ListTypeFromURI gets a ListType from the given URI.
	ListTypeFromURI(uri string) ListType

	// ServiceTypeIdentifierFromURI gets a LoTEServiceTypeIdentifier from the
	// given URI.
	ServiceTypeIdentifierFromURI(uri string) LoTEServiceTypeIdentifier

	// ServiceStatusFromURI gets a LoTEServiceStatus from the given URI.
	ServiceStatusFromURI(uri string) LoTEServiceStatus

	// CertificateApprovalStatusFromLabel gets a CertificateApprovalStatus
	// from the given label string.
	CertificateApprovalStatusFromLabel(label string) CertificateApprovalStatus

	// CertificateApprovalStatusFromDefinition gets a CertificateApprovalStatus
	// from the given listType, service type identifier and service status.
	CertificateApprovalStatusFromDefinition(listType ListType, sti LoTEServiceTypeIdentifier, status LoTEServiceStatus) CertificateApprovalStatus
}

// loTELoaderRegistry holds the LoTELoader implementations registered via
// RegisterLoTELoader, consulted in registration order — the Go equivalent
// of Java's ServiceLoader.load(LoTELoader.class) iteration.
var loTELoaderRegistry []LoTELoader

// RegisterLoTELoader registers a LoTELoader to be consulted by ListType,
// LoTEServiceStatus, LoTEServiceTypeIdentifier and CertificateApprovalStatus
// lookup functions, in the order loaders are registered.
func RegisterLoTELoader(l LoTELoader) {
	loTELoaderRegistry = append(loTELoaderRegistry, l)
}

// loTELoaders returns the currently registered LoTELoader instances.
func loTELoaders() []LoTELoader {
	return loTELoaderRegistry
}

// Upstream ships these two loaders in
// META-INF/services/eu.europa.esig.dss.enumerations.loader.LoTELoader, listed in this
// order, so ServiceLoader hands them to ListType.fromUri and friends with LoTEEnumLoader
// first and LoTEEmptyLoader second. Registering them here reproduces that default: a URI
// that matches a built-in enumeration resolves to the enum constant, and anything else
// falls through to LoTEEmptyLoader, which echoes the URI back in a label-less value rather
// than yielding nothing.
//
// The order is load-bearing and the registration deliberately lives in one init() in this
// file rather than one per loader file: Go runs a package's init() functions in filename
// order, and "lote_empty_loader.go" sorts before "lote_enum_loader.go", so per-file inits
// would register LoTEEmptyLoader first. Because LoTEEmptyLoader answers every lookup with a
// non-nil value, that inversion would shadow LoTEEnumLoader completely and no URI would
// ever resolve to its enum constant.
func init() {
	RegisterLoTELoader(NewLoTEEnumLoader())
	RegisterLoTELoader(NewLoTEEmptyLoader())
}
