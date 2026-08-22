// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CertificateValidity.java (DSS 6.5.RC1).
package spi

import "github.com/ryftcore/dss-go/dss/model"

// CertificateValidity stores the information about the validity of the signing certificate.
type CertificateValidity struct {
	// publicKey is used when only the public key is available (non AdES signature).
	publicKey *model.PublicKey

	// certificateToken is the certificate token, when available.
	certificateToken *model.CertificateToken

	// signerIdentifier is the signer identifier (used in CAdES).
	signerIdentifier *SignerIdentifier

	// signerIdMatch reports whether the CMS Signer id matches.
	signerIdMatch bool
	// digestPresent reports whether a digest is present.
	digestPresent bool
	// digestEqual reports whether the digest is equal.
	digestEqual bool
	// issuerSerialPresent reports whether the Issuer Serial is present.
	issuerSerialPresent bool
	// serialNumberEqual reports whether the Issuer Serial Number is equal.
	serialNumberEqual bool
	// distinguishedNameEqual reports whether the Issuer Serial Distinguished Name is equal.
	distinguishedNameEqual bool
	// responderIdPresent reports whether the OCSP Responder Id is present.
	responderIdPresent bool
	// responderIdMatch reports whether the OCSP Responder Id matches.
	responderIdMatch bool
}

// NewCertificateValidity creates an object containing all information concerning the
// validity of a candidate for the signing certificate.
// Port of the CertificateValidity(CertificateToken) constructor.
//
// Panics with the Java message when certificateToken is missing (Objects.requireNonNull).
func NewCertificateValidity(certificateToken *model.CertificateToken) *CertificateValidity {
	if certificateToken == nil {
		panic("CertificateToken cannot be null!")
	}
	return &CertificateValidity{certificateToken: certificateToken}
}

// NewCertificateValidityFromPublicKey creates an object containing all information
// concerning the validity of a candidate for the signing certificate which is based only on
// the PublicKey. To be used in case of a non AdES signature.
// Port of the CertificateValidity(PublicKey) constructor.
//
// Panics with the Java message when publicKey is missing (Objects.requireNonNull).
func NewCertificateValidityFromPublicKey(publicKey *model.PublicKey) *CertificateValidity {
	if publicKey == nil {
		panic("PublicKey cannot be null!")
	}
	return &CertificateValidity{publicKey: publicKey}
}

// NewCertificateValidityFromSignerIdentifier creates an object containing all information
// concerning the validity of a candidate for the signing certificate which is based only on
// the SignerIdentifier. To be used in case of a non AdES signature.
// Port of the CertificateValidity(SignerIdentifier) constructor.
//
// Panics with the Java message when signerIdentifier is missing (Objects.requireNonNull).
func NewCertificateValidityFromSignerIdentifier(signerIdentifier *SignerIdentifier) *CertificateValidity {
	if signerIdentifier == nil {
		panic("CertificateIdentifier cannot be null!")
	}
	return &CertificateValidity{signerIdentifier: signerIdentifier}
}

// PublicKey returns the PublicKey associated with this instance: if the certificateToken is
// not nil the associated PublicKey is returned, otherwise the provided publicKey is
// returned. NOTE: can return nil. Port of getPublicKey().
func (v *CertificateValidity) PublicKey() *model.PublicKey {
	if v.certificateToken == nil {
		return v.publicKey
	}
	return v.certificateToken.PublicKey()
}

// SignerInfo returns the associated SignerIdentifier. NOTE: can return nil.
// Port of getSignerInfo().
func (v *CertificateValidity) SignerInfo() *SignerIdentifier {
	if v.certificateToken == nil {
		return v.signerIdentifier
	}
	signerIdentifierFromCert := NewSignerIdentifier()
	signerIdentifierFromCert.SetIssuerName(v.certificateToken.IssuerX500Principal())
	signerIdentifierFromCert.SetSerialNumber(v.certificateToken.SerialNumber())
	return signerIdentifierFromCert
}

// CertificateToken gets the CertificateToken. Port of getCertificateToken().
func (v *CertificateValidity) CertificateToken() *model.CertificateToken {
	return v.certificateToken
}

// IsSignerIdMatch gets whether the CMS Signer Id matches. Port of isSignerIdMatch().
func (v *CertificateValidity) IsSignerIdMatch() bool {
	return v.signerIdMatch
}

// SetSignerIdMatch sets whether the CMS Signer Id matches. Port of setSignerIdMatch(boolean).
func (v *CertificateValidity) SetSignerIdMatch(signerIdMatch bool) {
	v.signerIdMatch = signerIdMatch
}

// IsDigestPresent gets whether the digest is present. Port of isDigestPresent().
func (v *CertificateValidity) IsDigestPresent() bool {
	return v.digestPresent
}

// SetDigestPresent sets whether the digest is present. Port of setDigestPresent(boolean).
func (v *CertificateValidity) SetDigestPresent(digestPresent bool) {
	v.digestPresent = digestPresent
}

// IsDigestEqual gets whether the digest is equal. Port of isDigestEqual().
func (v *CertificateValidity) IsDigestEqual() bool {
	return v.digestEqual
}

// SetDigestEqual sets whether the digest is equal. Port of setDigestEqual(boolean).
func (v *CertificateValidity) SetDigestEqual(digestEqual bool) {
	v.digestEqual = digestEqual
}

// IsIssuerSerialPresent indicates whether the IssuerSerial (issuerAndSerialNumber) is present
// in the signature. Port of isIssuerSerialPresent().
func (v *CertificateValidity) IsIssuerSerialPresent() bool {
	return v.issuerSerialPresent
}

// SetIssuerSerialPresent sets whether the IssuerSerial is present.
// Port of setIssuerSerialPresent(boolean).
func (v *CertificateValidity) SetIssuerSerialPresent(issuerSerialPresent bool) {
	v.issuerSerialPresent = issuerSerialPresent
}

// IsSerialNumberEqual indicates whether the SerialNumber equals. Port of isSerialNumberEqual().
func (v *CertificateValidity) IsSerialNumberEqual() bool {
	return v.serialNumberEqual
}

// SetSerialNumberEqual sets whether the serial number matches.
// Port of setSerialNumberEqual(boolean).
func (v *CertificateValidity) SetSerialNumberEqual(serialNumberEqual bool) {
	v.serialNumberEqual = serialNumberEqual
}

// IsDistinguishedNameEqual gets whether the distinguished name equals.
// Port of isDistinguishedNameEqual().
func (v *CertificateValidity) IsDistinguishedNameEqual() bool {
	return v.distinguishedNameEqual
}

// SetDistinguishedNameEqual sets whether the distinguished name equals.
// Port of setDistinguishedNameEqual(boolean).
func (v *CertificateValidity) SetDistinguishedNameEqual(distinguishedNameEqual bool) {
	v.distinguishedNameEqual = distinguishedNameEqual
}

// IsResponderIdPresent gets whether the ResponderId is present. Port of isResponderIdPresent().
func (v *CertificateValidity) IsResponderIdPresent() bool {
	return v.responderIdPresent
}

// SetResponderIdPresent sets whether the ResponderId is present.
// Port of setResponderIdPresent(boolean).
func (v *CertificateValidity) SetResponderIdPresent(responderIdPresent bool) {
	v.responderIdPresent = responderIdPresent
}

// IsResponderIdMatch gets whether the Responder Id matches. Port of isResponderIdMatch().
func (v *CertificateValidity) IsResponderIdMatch() bool {
	return v.responderIdMatch
}

// SetResponderIdMatch sets whether the ResponderId matches. Port of setResponderIdMatch(boolean).
func (v *CertificateValidity) SetResponderIdMatch(responderIdMatch bool) {
	v.responderIdMatch = responderIdMatch
}

// IsValid returns true if the certificate digest or IssuerSerial/issuerAndSerialNumber match
// or the certificate is signed. Port of isValid().
func (v *CertificateValidity) IsValid() bool {
	return v.IsDigestEqual() || (v.IsDistinguishedNameEqual() && v.IsSerialNumberEqual()) || v.IsResponderIdMatch()
}
