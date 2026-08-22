// Ported from dss-enumerations/.../CertificateRefOrigin.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// CertificateRefOrigin lists the supported origin types for a certificate reference.
type CertificateRefOrigin string

const (
	// CertificateRefOriginAttributeCertificateRefs: the certificate reference was
	// embedded in the signature 'attribute-certificate-references' attribute.
	CertificateRefOriginAttributeCertificateRefs CertificateRefOrigin = "ATTRIBUTE_CERTIFICATE_REFS"
	// CertificateRefOriginCompleteCertificateRefs: the certificate reference was
	// embedded in the signature 'complete-certificate-references' attribute.
	CertificateRefOriginCompleteCertificateRefs CertificateRefOrigin = "COMPLETE_CERTIFICATE_REFS"
	// CertificateRefOriginSigningCertificate: the certificate reference was
	// embedded in the signature 'signing-certificate' attribute.
	CertificateRefOriginSigningCertificate CertificateRefOrigin = "SIGNING_CERTIFICATE"
	// CertificateRefOriginKeyIdentifier is used as a hint to identify the signing
	// certificate (used in JAdES).
	CertificateRefOriginKeyIdentifier CertificateRefOrigin = "KEY_IDENTIFIER"
	// CertificateRefOriginX509URL is used as a hint to identify the resource
	// containing the signing certificate or certificate chain (used in JAdES).
	CertificateRefOriginX509URL CertificateRefOrigin = "X509_URL"
	// CertificateRefOriginPublicKey contains a public key of the signing certificate.
	CertificateRefOriginPublicKey CertificateRefOrigin = "PUBLIC_KEY"
	// CertificateRefOriginUnprotectedHeaderRefs: certificate reference present
	// within an unprotected header parameter (JWS or COSE).
	CertificateRefOriginUnprotectedHeaderRefs CertificateRefOrigin = "UNPROTECTED_HEADER_REFS"
)

// CertificateRefOriginValues returns all CertificateRefOrigin constants in declaration order.
func CertificateRefOriginValues() []CertificateRefOrigin {
	return []CertificateRefOrigin{
		CertificateRefOriginAttributeCertificateRefs,
		CertificateRefOriginCompleteCertificateRefs,
		CertificateRefOriginSigningCertificate,
		CertificateRefOriginKeyIdentifier,
		CertificateRefOriginX509URL,
		CertificateRefOriginPublicKey,
		CertificateRefOriginUnprotectedHeaderRefs,
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
