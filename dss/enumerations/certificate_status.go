// Ported from dss-enumerations/.../CertificateStatus.java (DSS 6.5.RC1).
package enumerations

// CertificateStatus defines the certificate revocation status.
type CertificateStatus string

const (
	// CertificateStatus_GOOD means the certificate is not revoked.
	CertificateStatus_GOOD CertificateStatus = "GOOD"
	// CertificateStatus_REVOKED means the certificate is revoked.
	CertificateStatus_REVOKED CertificateStatus = "REVOKED"
	// CertificateStatus_UNKNOWN means the certificate status is not known.
	CertificateStatus_UNKNOWN CertificateStatus = "UNKNOWN"
)

// CertificateStatusValues returns all constants in declaration order.
func CertificateStatusValues() []CertificateStatus {
	return []CertificateStatus{
		CertificateStatus_GOOD,
		CertificateStatus_REVOKED,
		CertificateStatus_UNKNOWN,
	}
}

// IsGood checks if the certificate status is valid.
func (c CertificateStatus) IsGood() bool {
	return CertificateStatus_GOOD == c
}

// IsRevoked checks if the certificate is revoked.
func (c CertificateStatus) IsRevoked() bool {
	return CertificateStatus_REVOKED == c
}

// IsKnown checks if the certificate status is known.
func (c CertificateStatus) IsKnown() bool {
	return CertificateStatus_UNKNOWN != c
}
