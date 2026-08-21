// Ported from dss-enumerations/.../loader/LoTEEmptyLoader.java (DSS 6.5.RC1).
//
// Files under loader/ are flattened into package enumerations (a Go
// subpackage would create an import cycle with the enumerations that depend
// on it).
package enumerations

// LoTEEmptyLoader is a LoTELoader implementation that echoes back the given
// URI/label wrapped in a value with an empty Label, rather than resolving
// against any known enumeration. It mirrors Java's LoTEEmptyLoader, whose
// anonymous getLabel() implementations returned null; the Go equivalent is
// the zero value "".
type LoTEEmptyLoader struct{}

// NewLoTEEmptyLoader creates a new LoTEEmptyLoader.
func NewLoTEEmptyLoader() *LoTEEmptyLoader {
	return &LoTEEmptyLoader{}
}

type loteEmptyListType struct {
	uri string
}

func (l *loteEmptyListType) Label() string { return "" }
func (l *loteEmptyListType) URI() string   { return l.uri }

// ListTypeFromURI implements LoTELoader by wrapping uri in a ListType with
// an empty label, without resolving it against any known enumeration.
func (l *LoTEEmptyLoader) ListTypeFromURI(uri string) ListType {
	return &loteEmptyListType{uri: uri}
}

type loteEmptyServiceTypeIdentifier struct {
	uri string
}

func (l *loteEmptyServiceTypeIdentifier) Label() string { return "" }
func (l *loteEmptyServiceTypeIdentifier) URI() string   { return l.uri }

// ServiceTypeIdentifierFromURI implements LoTELoader by wrapping uri in a
// LoTEServiceTypeIdentifier with an empty label, without resolving it
// against any known enumeration.
func (l *LoTEEmptyLoader) ServiceTypeIdentifierFromURI(uri string) LoTEServiceTypeIdentifier {
	return &loteEmptyServiceTypeIdentifier{uri: uri}
}

type loteEmptyServiceStatus struct {
	uri string
}

func (l *loteEmptyServiceStatus) Label() string { return "" }
func (l *loteEmptyServiceStatus) URI() string   { return l.uri }

// ServiceStatusFromURI implements LoTELoader by wrapping uri in a
// LoTEServiceStatus with an empty label, without resolving it against any
// known enumeration.
func (l *LoTEEmptyLoader) ServiceStatusFromURI(uri string) LoTEServiceStatus {
	return &loteEmptyServiceStatus{uri: uri}
}

type loteEmptyCertificateApprovalStatus struct {
	label string
}

func (l *loteEmptyCertificateApprovalStatus) ListType() ListType { return nil }
func (l *loteEmptyCertificateApprovalStatus) ServiceTypeIdentifier() LoTEServiceTypeIdentifier {
	return nil
}
func (l *loteEmptyCertificateApprovalStatus) ServiceStatus() LoTEServiceStatus { return nil }
func (l *loteEmptyCertificateApprovalStatus) Label() string                    { return l.label }

// CertificateApprovalStatusFromLabel implements LoTELoader by wrapping label
// in a CertificateApprovalStatus, without resolving it against any known
// enumeration.
func (l *LoTEEmptyLoader) CertificateApprovalStatusFromLabel(label string) CertificateApprovalStatus {
	return &loteEmptyCertificateApprovalStatus{label: label}
}

// CertificateApprovalStatusFromDefinition is not supported by LoTEEmptyLoader.
func (l *LoTEEmptyLoader) CertificateApprovalStatusFromDefinition(listType ListType, sti LoTEServiceTypeIdentifier, status LoTEServiceStatus) CertificateApprovalStatus {
	return nil
}

// compile-time interface assertion.
var _ LoTELoader = (*LoTEEmptyLoader)(nil)
