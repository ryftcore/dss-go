// Ported from dss-enumerations/.../CertificateQualifiedStatus.java (DSS 6.5.RC1).
package enumerations

// CertificateQualifiedStatus defines the qualification status of a certificate.
type CertificateQualifiedStatus string

const (
	// CertificateQualifiedStatus_QC is Qualified.
	CertificateQualifiedStatus_QC CertificateQualifiedStatus = "QC"
	// CertificateQualifiedStatus_NOT_QC is Not qualified.
	CertificateQualifiedStatus_NOT_QC CertificateQualifiedStatus = "NOT_QC"
)

// certificateQualifiedStatusLabels holds the user-friendly label for each constant.
var certificateQualifiedStatusLabels = map[CertificateQualifiedStatus]string{
	CertificateQualifiedStatus_QC:     "Qualified",
	CertificateQualifiedStatus_NOT_QC: "Not qualified",
}

// CertificateQualifiedStatusValues returns all constants in declaration order.
func CertificateQualifiedStatusValues() []CertificateQualifiedStatus {
	return []CertificateQualifiedStatus{
		CertificateQualifiedStatus_QC,
		CertificateQualifiedStatus_NOT_QC,
	}
}

// Label returns a user-friendly certificate qualification label.
func (c CertificateQualifiedStatus) Label() string {
	return certificateQualifiedStatusLabels[c]
}

// CertificateQualifiedStatusIsQC verifies if the given CertificateQualifiedStatus
// is related to a qualified certificate.
func CertificateQualifiedStatusIsQC(status CertificateQualifiedStatus) bool {
	return CertificateQualifiedStatus_QC == status
}
