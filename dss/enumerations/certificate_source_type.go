// Ported from dss-enumerations/.../CertificateSourceType.java (DSS 6.5.RC1).
package enumerations

import "fmt"

// CertificateSourceType represents standard sources for a certificate.
// Indicates where the certificate comes from.
type CertificateSourceType string

const (
	// CertificateSourceTypeTrustedStore defines a pre-defined trusted source.
	CertificateSourceTypeTrustedStore CertificateSourceType = "TRUSTED_STORE"
	// CertificateSourceTypeTrustedList defines a certificate source
	// populated by a TLValidationJob.
	CertificateSourceTypeTrustedList CertificateSourceType = "TRUSTED_LIST"
	// CertificateSourceTypeTrustedEntities defines a certificate source
	// populated by processing List(s) of Trusted Entities.
	CertificateSourceTypeTrustedEntities CertificateSourceType = "TRUSTED_ENTITIES"
	// CertificateSourceTypeSignature: certificate source extracted from a
	// signature.
	CertificateSourceTypeSignature CertificateSourceType = "SIGNATURE"
	// CertificateSourceTypeOCSPResponse: certificate source extracted
	// from an OCSP response.
	CertificateSourceTypeOCSPResponse CertificateSourceType = "OCSP_RESPONSE"
	// CertificateSourceTypeOther: other types of certificate sources.
	CertificateSourceTypeOther CertificateSourceType = "OTHER"
	// CertificateSourceTypeAIA: the certificate source has been obtained
	// by AIA.
	CertificateSourceTypeAIA CertificateSourceType = "AIA"
	// CertificateSourceTypeTimestamp: certificate source extracted from a
	// timestamp.
	CertificateSourceTypeTimestamp CertificateSourceType = "TIMESTAMP"
	// CertificateSourceTypeEvidenceRecord: certificate source extracted
	// from an Evidence record.
	CertificateSourceTypeEvidenceRecord CertificateSourceType = "EVIDENCE_RECORD"
	// CertificateSourceTypeEAA: certificate source extracted from an EAA
	// token's claims.
	CertificateSourceTypeEAA CertificateSourceType = "EAA"
	// CertificateSourceTypeUnknown: the unknown origin of a certificate source.
	CertificateSourceTypeUnknown CertificateSourceType = "UNKNOWN"
)

// CertificateSourceTypeValues returns all constants in declaration order.
func CertificateSourceTypeValues() []CertificateSourceType {
	return []CertificateSourceType{
		CertificateSourceTypeTrustedStore,
		CertificateSourceTypeTrustedList,
		CertificateSourceTypeTrustedEntities,
		CertificateSourceTypeSignature,
		CertificateSourceTypeOCSPResponse,
		CertificateSourceTypeOther,
		CertificateSourceTypeAIA,
		CertificateSourceTypeTimestamp,
		CertificateSourceTypeEvidenceRecord,
		CertificateSourceTypeEAA,
		CertificateSourceTypeUnknown,
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
	return c == CertificateSourceTypeTrustedStore || c == CertificateSourceTypeTrustedList || c == CertificateSourceTypeTrustedEntities
}
