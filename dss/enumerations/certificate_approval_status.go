// Ported from dss-enumerations/.../CertificateApprovalStatus.java (DSS 6.5.RC1).
package enumerations

// CertificateApprovalStatus represents a certificate approval status, e.g.
// in the context of EUDI Wallet.
type CertificateApprovalStatus interface {
	// ListType gets a list type URI of the List of Trusted Entities
	// providing the certificates of the given usage.
	ListType() ListType

	// ServiceTypeIdentifier gets the ServiceTypeIdentifier related to the
	// certificate approval status.
	ServiceTypeIdentifier() LoTEServiceTypeIdentifier

	// ServiceStatus gets the ServiceStatus corresponding to the
	// certificate approval status.
	ServiceStatus() LoTEServiceStatus

	// Label gets user-friendly description of the certificate approval
	// status.
	Label() string
}

// CertificateApprovalStatusFromLabel returns a CertificateApprovalStatus for
// the given URI, or nil if none of the registered LoTELoaders resolve it.
func CertificateApprovalStatusFromLabel(label string) CertificateApprovalStatus {
	for _, loader := range loTELoaders() {
		if status := loader.CertificateApprovalStatusFromLabel(label); status != nil {
			return status
		}
	}
	return nil
}

// CertificateApprovalStatusFromDefinition returns a CertificateApprovalStatus
// for the given definition, or nil if none of the registered LoTELoaders
// resolve it.
func CertificateApprovalStatusFromDefinition(listType ListType, sti LoTEServiceTypeIdentifier, status LoTEServiceStatus) CertificateApprovalStatus {
	for _, loader := range loTELoaders() {
		if certApprovalStatus := loader.CertificateApprovalStatusFromDefinition(listType, sti, status); certApprovalStatus != nil {
			return certApprovalStatus
		}
	}
	return nil
}

// certificateApprovalStatus is a plain CertificateApprovalStatus
// implementation backing NewCertificateApprovalStatus (Java's
// CertificateApprovalStatus.create anonymous class).
type certificateApprovalStatus struct {
	label    string
	listType ListType
	sti      LoTEServiceTypeIdentifier
	status   LoTEServiceStatus
}

func (c *certificateApprovalStatus) ListType() ListType                               { return c.listType }
func (c *certificateApprovalStatus) ServiceTypeIdentifier() LoTEServiceTypeIdentifier { return c.sti }
func (c *certificateApprovalStatus) ServiceStatus() LoTEServiceStatus                 { return c.status }
func (c *certificateApprovalStatus) Label() string                                    { return c.label }

// NewCertificateApprovalStatus creates a new CertificateApprovalStatus from
// the provided data.
func NewCertificateApprovalStatus(label string, listType ListType, sti LoTEServiceTypeIdentifier, status LoTEServiceStatus) CertificateApprovalStatus {
	return &certificateApprovalStatus{
		label:    label,
		listType: listType,
		sti:      sti,
		status:   status,
	}
}
