// Ported from dss-enumerations/.../CertificateStatus.java (DSS 6.5.RC1).
package enumerations

// CertificateStatus defines the certificate revocation status.
type CertificateStatus string

const (
	// CertificateStatusGood means the certificate is not revoked.
	CertificateStatusGood CertificateStatus = "GOOD"
	// CertificateStatusRevoked means the certificate is revoked.
	CertificateStatusRevoked CertificateStatus = "REVOKED"
	// CertificateStatusUnknown means the certificate status is not known.
	CertificateStatusUnknown CertificateStatus = "UNKNOWN"
)

// CertificateStatusValues returns all constants in declaration order.
func CertificateStatusValues() []CertificateStatus {
	return []CertificateStatus{
		CertificateStatusGood,
		CertificateStatusRevoked,
		CertificateStatusUnknown,
	}
}

// IsGood checks if the certificate status is valid.
func (c CertificateStatus) IsGood() bool {
	return CertificateStatusGood == c
}

// IsRevoked checks if the certificate is revoked.
func (c CertificateStatus) IsRevoked() bool {
	return CertificateStatusRevoked == c
}

// IsKnown checks if the certificate status is known.
func (c CertificateStatus) IsKnown() bool {
	return CertificateStatusUnknown != c
}
