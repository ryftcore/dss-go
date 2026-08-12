// Ported from dss-enumerations/.../CertificateSourceType.java (DSS 6.5.RC1).
package enumerations

import "fmt"

// CertificateSourceType represents standard sources for a certificate.
// Indicates where the certificate comes from.
type CertificateSourceType string

const (
	// CertificateSourceType_TRUSTED_STORE defines a pre-defined trusted source.
	CertificateSourceType_TRUSTED_STORE CertificateSourceType = "TRUSTED_STORE"
	// CertificateSourceType_TRUSTED_LIST defines a certificate source
	// populated by a TLValidationJob.
	CertificateSourceType_TRUSTED_LIST CertificateSourceType = "TRUSTED_LIST"
	// CertificateSourceType_TRUSTED_ENTITIES defines a certificate source
	// populated by processing List(s) of Trusted Entities.
	CertificateSourceType_TRUSTED_ENTITIES CertificateSourceType = "TRUSTED_ENTITIES"
	// CertificateSourceType_SIGNATURE: certificate source extracted from a
	// signature.
	CertificateSourceType_SIGNATURE CertificateSourceType = "SIGNATURE"
	// CertificateSourceType_OCSP_RESPONSE: certificate source extracted
	// from an OCSP response.
	CertificateSourceType_OCSP_RESPONSE CertificateSourceType = "OCSP_RESPONSE"
	// CertificateSourceType_OTHER: other types of certificate sources.
	CertificateSourceType_OTHER CertificateSourceType = "OTHER"
	// CertificateSourceType_AIA: the certificate source has been obtained
	// by AIA.
	CertificateSourceType_AIA CertificateSourceType = "AIA"
	// CertificateSourceType_TIMESTAMP: certificate source extracted from a
	// timestamp.
	CertificateSourceType_TIMESTAMP CertificateSourceType = "TIMESTAMP"
	// CertificateSourceType_EVIDENCE_RECORD: certificate source extracted
	// from an Evidence record.
	CertificateSourceType_EVIDENCE_RECORD CertificateSourceType = "EVIDENCE_RECORD"
	// CertificateSourceType_EAA: certificate source extracted from an EAA
	// token's claims.
	CertificateSourceType_EAA CertificateSourceType = "EAA"
	// CertificateSourceType_UNKNOWN: the unknown origin of a certificate source.
	CertificateSourceType_UNKNOWN CertificateSourceType = "UNKNOWN"
)

// CertificateSourceTypeValues returns all constants in declaration order.
func CertificateSourceTypeValues() []CertificateSourceType {
	return []CertificateSourceType{
		CertificateSourceType_TRUSTED_STORE,
		CertificateSourceType_TRUSTED_LIST,
		CertificateSourceType_TRUSTED_ENTITIES,
		CertificateSourceType_SIGNATURE,
		CertificateSourceType_OCSP_RESPONSE,
		CertificateSourceType_OTHER,
		CertificateSourceType_AIA,
		CertificateSourceType_TIMESTAMP,
		CertificateSourceType_EVIDENCE_RECORD,
		CertificateSourceType_EAA,
		CertificateSourceType_UNKNOWN,
	}
}

// CertificateSourceTypeValueOf returns the CertificateSourceType matching the given Java enum name.
func CertificateSourceTypeValueOf(name string) (CertificateSourceType, error) {
	for _, v := range CertificateSourceTypeValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant CertificateSourceType.%s", name)
}

// IsTrusted gets whether the certificate source is trusted.
func (c CertificateSourceType) IsTrusted() bool {
	return c == CertificateSourceType_TRUSTED_STORE || c == CertificateSourceType_TRUSTED_LIST || c == CertificateSourceType_TRUSTED_ENTITIES
}
