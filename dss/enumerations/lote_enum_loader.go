// Ported from dss-enumerations/.../loader/LoTEEnumLoader.java (DSS 6.5.RC1).
//
// Files under loader/ are flattened into package enumerations (a Go
// subpackage would create an import cycle with the enumerations that depend
// on it).
package enumerations

import "strings"

// LoTEEnumLoader is a LoTELoader implementation backed by the built-in
// LoTETypeEnum, LoTEServiceTypeIdentifierEnum, LoTEServiceStatusEnum and
// CertificateApprovalStatusEnum enumerations.
//
// NOTE: LoTETypeEnum, LoTEServiceTypeIdentifierEnum, LoTEServiceStatusEnum
// and CertificateApprovalStatusEnum (with their *Values() accessors and the
// CertificateApprovalStatusEnum_CERT_FOR_UNKNOWN constant) are defined
// outside this file's manifest and are assumed to exist per the porting
// brief.
type LoTEEnumLoader struct{}

// NewLoTEEnumLoader creates a new LoTEEnumLoader.
func NewLoTEEnumLoader() *LoTEEnumLoader {
	return &LoTEEnumLoader{}
}

func (l *LoTEEnumLoader) ListTypeFromURI(uri string) ListType {
	for _, t := range LoTETypeEnumValues() {
		if strings.EqualFold(uri, t.URI()) {
			return t
		}
	}
	return nil
}

func (l *LoTEEnumLoader) ServiceTypeIdentifierFromURI(uri string) LoTEServiceTypeIdentifier {
	for _, sti := range LoTEServiceTypeIdentifierEnumValues() {
		if strings.EqualFold(uri, sti.URI()) {
			return sti
		}
	}
	return nil
}

func (l *LoTEEnumLoader) ServiceStatusFromURI(uri string) LoTEServiceStatus {
	for _, status := range LoTEServiceStatusEnumValues() {
		if strings.EqualFold(uri, status.URI()) {
			return status
		}
	}
	return nil
}

func (l *LoTEEnumLoader) CertificateApprovalStatusFromLabel(label string) CertificateApprovalStatus {
	for _, certApprovalStatus := range CertificateApprovalStatusEnumValues() {
		if strings.EqualFold(label, certApprovalStatus.Label()) {
			return certApprovalStatus
		}
	}
	return nil
}

// CertificateApprovalStatusFromDefinition ported verbatim from Java,
// including its operator-precedence quirk: the source condition mixes ||
// and && without full parenthesization, so && binds tighter than ||
// (identically in Java and Go), producing the same non-obvious matching
// behavior in both languages. This is flagged upstream behavior, not
// something this port "corrects" — see PORTER_BRIEF notes.
func (l *LoTEEnumLoader) CertificateApprovalStatusFromDefinition(listType ListType, sti LoTEServiceTypeIdentifier, status LoTEServiceStatus) CertificateApprovalStatus {
	for _, certApprovalStatus := range CertificateApprovalStatusEnumValues() {
		if (listType == nil && certApprovalStatus.ListType() == nil) ||
			listType != nil && listType == certApprovalStatus.ListType() &&
				(sti == nil && certApprovalStatus.ServiceTypeIdentifier() == nil) ||
			sti != nil && sti == certApprovalStatus.ServiceTypeIdentifier() &&
				(status == nil && certApprovalStatus.ServiceStatus() == nil) ||
			status != nil && status == certApprovalStatus.ServiceStatus() {
			return certApprovalStatus
		}
	}
	return CertificateApprovalStatusEnum_CERT_FOR_UNKNOWN
}

// compile-time interface assertion.
var _ LoTELoader = (*LoTEEnumLoader)(nil)
