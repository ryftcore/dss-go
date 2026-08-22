// Ported from dss-enumerations/.../CertificateOrigin.java (DSS 6.5.RC1).
package enumerations

import "fmt"

// CertificateOrigin represents possible origin types for a certificate.
type CertificateOrigin string

const (
	// CertificateOriginKeyInfo: certificates extracted from KeyInfo
	// element, XAdES specific.
	CertificateOriginKeyInfo CertificateOrigin = "KEY_INFO"
	// CertificateOriginSignedData: certificates extracted from a signed
	// attribute (CAdES).
	CertificateOriginSignedData CertificateOrigin = "SIGNED_DATA"
	// CertificateOriginCertificateValues: certificates extracted from
	// CertificateValues element.
	CertificateOriginCertificateValues CertificateOrigin = "CERTIFICATE_VALUES"
	// CertificateOriginAttrAuthoritiesCertValues: certificates extracted
	// from AttrAuthoritiesCertValues element, XAdES specific.
	CertificateOriginAttrAuthoritiesCertValues CertificateOrigin = "ATTR_AUTHORITIES_CERT_VALUES"
	// CertificateOriginTimestampValidationData: certificates extracted
	// from TimeStampValidationData element.
	CertificateOriginTimestampValidationData CertificateOrigin = "TIMESTAMP_VALIDATION_DATA"
	// CertificateOriginAnyValidationData: certificates extracted from
	// AnyValidationData element.
	CertificateOriginAnyValidationData CertificateOrigin = "ANY_VALIDATION_DATA"
	// CertificateOriginDSSDictionary: certificates extracted from DSS
	// dictionary, PAdES specific.
	CertificateOriginDSSDictionary CertificateOrigin = "DSS_DICTIONARY"
	// CertificateOriginVRIDictionary: certificates extracted from VRI
	// dictionary, PAdES specific.
	CertificateOriginVRIDictionary CertificateOrigin = "VRI_DICTIONARY"
	// CertificateOriginBasicOCSPResp: certificates extracted from an
	// OCSP Response.
	CertificateOriginBasicOCSPResp CertificateOrigin = "BASIC_OCSP_RESP"
	// CertificateOriginEvidenceRecord: certificates extracted from an
	// Evidence Record.
	CertificateOriginEvidenceRecord CertificateOrigin = "EVIDENCE_RECORD"
	// CertificateOriginUnprotectedHeader: certificates present within an
	// unprotected header parameter (JWS or COSE).
	CertificateOriginUnprotectedHeader CertificateOrigin = "UNPROTECTED_HEADER"
	// CertificateOriginEAA: certificates present within a Token Status
	// List claim of an EAA.
	CertificateOriginEAA CertificateOrigin = "EAA"
)

// CertificateOriginValues returns all constants in declaration order.
func CertificateOriginValues() []CertificateOrigin {
	return []CertificateOrigin{
		CertificateOriginKeyInfo,
		CertificateOriginSignedData,
		CertificateOriginCertificateValues,
		CertificateOriginAttrAuthoritiesCertValues,
		CertificateOriginTimestampValidationData,
		CertificateOriginAnyValidationData,
		CertificateOriginDSSDictionary,
		CertificateOriginVRIDictionary,
		CertificateOriginBasicOCSPResp,
		CertificateOriginEvidenceRecord,
		CertificateOriginUnprotectedHeader,
		CertificateOriginEAA,
	}
}

// CertificateOriginValueOf returns the CertificateOrigin matching the given Java enum name.
func CertificateOriginValueOf(name string) (CertificateOrigin, error) {
	for _, v := range CertificateOriginValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant CertificateOrigin.%s", name)
}
