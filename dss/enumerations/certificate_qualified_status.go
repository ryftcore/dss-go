// Ported from dss-enumerations/.../CertificateQualifiedStatus.java (DSS 6.5.RC1).
package enumerations

// CertificateQualifiedStatus defines the qualification status of a certificate.
type CertificateQualifiedStatus string

const (
	// CertificateQualifiedStatusQC is Qualified.
	CertificateQualifiedStatusQC CertificateQualifiedStatus = "QC"
	// CertificateQualifiedStatusNotQC is Not qualified.
	CertificateQualifiedStatusNotQC CertificateQualifiedStatus = "NOT_QC"
)

// certificateQualifiedStatusLabels holds the user-friendly label for each constant.
var certificateQualifiedStatusLabels = map[CertificateQualifiedStatus]string{
	CertificateQualifiedStatusQC:    "Qualified",
	CertificateQualifiedStatusNotQC: "Not qualified",
}

// CertificateQualifiedStatusValues returns all constants in declaration order.
func CertificateQualifiedStatusValues() []CertificateQualifiedStatus {
	return []CertificateQualifiedStatus{
		CertificateQualifiedStatusQC,
		CertificateQualifiedStatusNotQC,
	}
}

// Label returns a user-friendly certificate qualification label.
func (c CertificateQualifiedStatus) Label() string {
	return certificateQualifiedStatusLabels[c]
}

// CertificateQualifiedStatusIsQC verifies if the given CertificateQualifiedStatus
// is related to a qualified certificate.
func CertificateQualifiedStatusIsQC(status CertificateQualifiedStatus) bool {
	return CertificateQualifiedStatusQC == status
}
