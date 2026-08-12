// Ported from dss-enumerations/.../CertificateRefOrigin.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// CertificateRefOrigin lists the supported origin types for a certificate reference.
type CertificateRefOrigin string

const (
	// CertificateRefOrigin_ATTRIBUTE_CERTIFICATE_REFS: the certificate reference was
	// embedded in the signature 'attribute-certificate-references' attribute.
	CertificateRefOrigin_ATTRIBUTE_CERTIFICATE_REFS CertificateRefOrigin = "ATTRIBUTE_CERTIFICATE_REFS"
	// CertificateRefOrigin_COMPLETE_CERTIFICATE_REFS: the certificate reference was
	// embedded in the signature 'complete-certificate-references' attribute.
	CertificateRefOrigin_COMPLETE_CERTIFICATE_REFS CertificateRefOrigin = "COMPLETE_CERTIFICATE_REFS"
	// CertificateRefOrigin_SIGNING_CERTIFICATE: the certificate reference was
	// embedded in the signature 'signing-certificate' attribute.
	CertificateRefOrigin_SIGNING_CERTIFICATE CertificateRefOrigin = "SIGNING_CERTIFICATE"
	// CertificateRefOrigin_KEY_IDENTIFIER is used as a hint to identify the signing
	// certificate (used in JAdES).
	CertificateRefOrigin_KEY_IDENTIFIER CertificateRefOrigin = "KEY_IDENTIFIER"
	// CertificateRefOrigin_X509_URL is used as a hint to identify the resource
	// containing the signing certificate or certificate chain (used in JAdES).
	CertificateRefOrigin_X509_URL CertificateRefOrigin = "X509_URL"
	// CertificateRefOrigin_PUBLIC_KEY contains a public key of the signing certificate.
	CertificateRefOrigin_PUBLIC_KEY CertificateRefOrigin = "PUBLIC_KEY"
	// CertificateRefOrigin_UNPROTECTED_HEADER_REFS: certificate reference present
	// within an unprotected header parameter (JWS or COSE).
	CertificateRefOrigin_UNPROTECTED_HEADER_REFS CertificateRefOrigin = "UNPROTECTED_HEADER_REFS"
)

// CertificateRefOriginValues returns all CertificateRefOrigin constants in declaration order.
func CertificateRefOriginValues() []CertificateRefOrigin {
	return []CertificateRefOrigin{
		CertificateRefOrigin_ATTRIBUTE_CERTIFICATE_REFS,
		CertificateRefOrigin_COMPLETE_CERTIFICATE_REFS,
		CertificateRefOrigin_SIGNING_CERTIFICATE,
		CertificateRefOrigin_KEY_IDENTIFIER,
		CertificateRefOrigin_X509_URL,
		CertificateRefOrigin_PUBLIC_KEY,
		CertificateRefOrigin_UNPROTECTED_HEADER_REFS,
	}
}

// CertificateRefOriginValueOf returns the CertificateRefOrigin matching the given Java enum name.
func CertificateRefOriginValueOf(name string) (CertificateRefOrigin, error) {
	for _, v := range CertificateRefOriginValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant CertificateRefOrigin.%s", name)
}
