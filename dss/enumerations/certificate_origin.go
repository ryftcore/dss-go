// Ported from dss-enumerations/.../CertificateOrigin.java (DSS 6.5.RC1).
package enumerations

import "fmt"

// CertificateOrigin represents possible origin types for a certificate.
type CertificateOrigin string

const (
	// CertificateOrigin_KEY_INFO: certificates extracted from KeyInfo
	// element, XAdES specific.
	CertificateOrigin_KEY_INFO CertificateOrigin = "KEY_INFO"
	// CertificateOrigin_SIGNED_DATA: certificates extracted from a signed
	// attribute (CAdES).
	CertificateOrigin_SIGNED_DATA CertificateOrigin = "SIGNED_DATA"
	// CertificateOrigin_CERTIFICATE_VALUES: certificates extracted from
	// CertificateValues element.
	CertificateOrigin_CERTIFICATE_VALUES CertificateOrigin = "CERTIFICATE_VALUES"
	// CertificateOrigin_ATTR_AUTHORITIES_CERT_VALUES: certificates extracted
	// from AttrAuthoritiesCertValues element, XAdES specific.
	CertificateOrigin_ATTR_AUTHORITIES_CERT_VALUES CertificateOrigin = "ATTR_AUTHORITIES_CERT_VALUES"
	// CertificateOrigin_TIMESTAMP_VALIDATION_DATA: certificates extracted
	// from TimeStampValidationData element.
	CertificateOrigin_TIMESTAMP_VALIDATION_DATA CertificateOrigin = "TIMESTAMP_VALIDATION_DATA"
	// CertificateOrigin_ANY_VALIDATION_DATA: certificates extracted from
	// AnyValidationData element.
	CertificateOrigin_ANY_VALIDATION_DATA CertificateOrigin = "ANY_VALIDATION_DATA"
	// CertificateOrigin_DSS_DICTIONARY: certificates extracted from DSS
	// dictionary, PAdES specific.
	CertificateOrigin_DSS_DICTIONARY CertificateOrigin = "DSS_DICTIONARY"
	// CertificateOrigin_VRI_DICTIONARY: certificates extracted from VRI
	// dictionary, PAdES specific.
	CertificateOrigin_VRI_DICTIONARY CertificateOrigin = "VRI_DICTIONARY"
	// CertificateOrigin_BASIC_OCSP_RESP: certificates extracted from an
	// OCSP Response.
	CertificateOrigin_BASIC_OCSP_RESP CertificateOrigin = "BASIC_OCSP_RESP"
	// CertificateOrigin_EVIDENCE_RECORD: certificates extracted from an
	// Evidence Record.
	CertificateOrigin_EVIDENCE_RECORD CertificateOrigin = "EVIDENCE_RECORD"
	// CertificateOrigin_UNPROTECTED_HEADER: certificates present within an
	// unprotected header parameter (JWS or COSE).
	CertificateOrigin_UNPROTECTED_HEADER CertificateOrigin = "UNPROTECTED_HEADER"
	// CertificateOrigin_EAA: certificates present within a Token Status
	// List claim of an EAA.
	CertificateOrigin_EAA CertificateOrigin = "EAA"
)

// CertificateOriginValues returns all constants in declaration order.
func CertificateOriginValues() []CertificateOrigin {
	return []CertificateOrigin{
		CertificateOrigin_KEY_INFO,
		CertificateOrigin_SIGNED_DATA,
		CertificateOrigin_CERTIFICATE_VALUES,
		CertificateOrigin_ATTR_AUTHORITIES_CERT_VALUES,
		CertificateOrigin_TIMESTAMP_VALIDATION_DATA,
		CertificateOrigin_ANY_VALIDATION_DATA,
		CertificateOrigin_DSS_DICTIONARY,
		CertificateOrigin_VRI_DICTIONARY,
		CertificateOrigin_BASIC_OCSP_RESP,
		CertificateOrigin_EVIDENCE_RECORD,
		CertificateOrigin_UNPROTECTED_HEADER,
		CertificateOrigin_EAA,
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
