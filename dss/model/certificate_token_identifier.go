// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/CertificateTokenIdentifier.java (DSS 6.5.RC1).
package model

// CertificateTokenIdentifier is the unique id of a CertificateToken.
type CertificateTokenIdentifier struct {
	TokenIdentifier
}

// NewCertificateTokenIdentifier builds the identifier of the given certificate token: the
// SHA-256 digest of its DER encoding, rendered with the "C-" prefix.
func NewCertificateTokenIdentifier(certificateToken *CertificateToken) *CertificateTokenIdentifier {
	return &CertificateTokenIdentifier{
		NewTokenIdentifierFromToken("CertificateTokenIdentifier", "C-", certificateToken),
	}
}

// compile-time interface assertion.
var _ Identifier = (*CertificateTokenIdentifier)(nil)
