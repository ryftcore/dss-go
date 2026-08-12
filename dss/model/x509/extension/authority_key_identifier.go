// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/extension/AuthorityKeyIdentifier.java (DSS 6.5.RC1).
package extension

import "github.com/utain/esig/dss/enumerations"

// AuthorityKeyIdentifier is RFC 5280 4.2.1.1. Authority Key Identifier.
//
// The authority key identifier extension provides a means of identifying the public key
// corresponding to the private key used to sign a certificate. This extension is used
// where an issuer has multiple signing keys (either due to multiple concurrent key pairs
// or due to changeover). The identification MAY be based on either the key identifier
// (the subject key identifier in the issuer's certificate) or the issuer name and serial
// number.
type AuthorityKeyIdentifier struct {
	CertificateExtension

	// keyIdentifier is the key identifier property.
	keyIdentifier []byte

	// authorityCertIssuerSerial is the DER-encoded representation of the combination of
	// authorityCertIssuer and authorityCertSerialNumber parameters.
	authorityCertIssuerSerial []byte
}

// NewAuthorityKeyIdentifier builds an AuthorityKeyIdentifier extension.
func NewAuthorityKeyIdentifier() *AuthorityKeyIdentifier {
	return &AuthorityKeyIdentifier{
		CertificateExtension: NewCertificateExtensionFromEnum(enumerations.CertificateExtensionEnum_AUTHORITY_KEY_IDENTIFIER),
	}
}

// KeyIdentifier returns the key identifier.
func (a *AuthorityKeyIdentifier) KeyIdentifier() []byte {
	return a.keyIdentifier
}

// SetKeyIdentifier sets the key identifier.
func (a *AuthorityKeyIdentifier) SetKeyIdentifier(keyIdentifier []byte) {
	a.keyIdentifier = keyIdentifier
}

// AuthorityCertIssuerSerial gets the DER-encoded IssuerSerial built on authorityCertIssuer
// and authorityCertSerialNumber parameters.
func (a *AuthorityKeyIdentifier) AuthorityCertIssuerSerial() []byte {
	return a.authorityCertIssuerSerial
}

// SetAuthorityCertIssuerSerial sets the DER-encoded IssuerSerial built on
// authorityCertIssuer and authorityCertSerialNumber parameters.
func (a *AuthorityKeyIdentifier) SetAuthorityCertIssuerSerial(authorityCertIssuerSerial []byte) {
	a.authorityCertIssuerSerial = authorityCertIssuerSerial
}
