// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/Token.java (DSS 6.5.RC1).
package model

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// Token is the contract of the different token types (certificate, OCSP, CRL, timestamp)
// used in the process of signature validation. It is the polymorphic half of the Java
// abstract class Token; the state and the concrete method bodies live in TokenBase.
type Token interface {
	IdentifierBasedObject

	// IsSelfSigned reports whether the token is self-signed. Only a CertificateToken can
	// answer true.
	IsSelfSigned() bool
	// DSSIDAsString returns the unique string of the token, i.e. its identifier's XML Id.
	DSSIDAsString() string
	// IsSignedByToken reports whether the token is signed by the given certificate token.
	// Port of the isSignedBy(CertificateToken) overload.
	IsSignedByToken(token *CertificateToken) bool
	// IsSignedBy reports whether the token is signed by the given public key. Port of the
	// isSignedBy(PublicKey) overload.
	IsSignedBy(publicKey *PublicKey) bool
	// IssuerX500Principal returns the X500Principal of the certificate that signed this token.
	IssuerX500Principal() *X500Principal
	// CreationDate returns the creation date of the token (notBefore for a certificate,
	// productionDate for revocation data, ...).
	CreationDate() time.Time
	// Abbreviation returns the DSS abbreviation of the token, used for debugging.
	Abbreviation() string
	// IssuerEntityKey returns the entity key identifier of the token's issuer, or nil when
	// the signer has not been established with a successful CheckIsSignedBy call.
	IssuerEntityKey() *EntityIdentifier
	// SignatureAlgorithm returns the algorithm that was used to sign the token.
	SignatureAlgorithm() enumerations.SignatureAlgorithm
	// IsSignatureIntact reports whether the token's signature is intact.
	IsSignatureIntact() bool
	// IsValid reports whether the conditions corresponding to the token validity are met.
	IsValid() bool
	// SignatureValidity returns the three-state status of the token's signature validity.
	SignatureValidity() enumerations.SignatureValidity
	// InvalidityReason returns the token invalidity reason, empty when the token is valid.
	InvalidityReason() string
	// PublicKeyOfTheSigner returns the public key which signed this token.
	PublicKeyOfTheSigner() *PublicKey
	// ToString returns a string representation of the token using the given indentation.
	ToString(indentStr string) string
	// String returns ToString("").
	String() string
	// Encoded returns the encoded form of the wrapped token.
	Encoded() []byte
	// Digest returns the digest value of the wrapped token for the requested algorithm.
	Digest(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error)
}

// TokenOverrides declares the operations Java's abstract Token class declares abstract, or
// expects a subclass to override, and that the base implementation itself calls back into.
// It stands in for the virtual dispatch a Java abstract class gets for free; a concrete
// token registers itself with TokenBase.InitToken so that the base can reach them.
type TokenOverrides interface {
	// BuildTokenIdentifier builds the token's unique identifier. Port of the abstract
	// protected buildTokenIdentifier().
	BuildTokenIdentifier() *TokenIdentifier
	// CheckIsSignedBy verifies whether the token has been signed by the given public key,
	// updating the token's validity state and invalidity reason. Port of the abstract
	// protected checkIsSignedBy(PublicKey).
	CheckIsSignedBy(publicKey *PublicKey) enumerations.SignatureValidity
	// IssuerX500Principal returns the issuer's X500Principal. Port of the abstract
	// getIssuerX500Principal().
	IssuerX500Principal() *X500Principal
	// IsSelfSigned reports whether the token is self-signed; TokenBase supplies the
	// always-false default and CertificateToken overrides it.
	IsSelfSigned() bool
}

// TokenBase carries the state and the concrete behaviour of the Java abstract class Token.
// Concrete tokens embed it and register themselves with InitToken.
type TokenBase struct {
	// overrides points back at the concrete token; see InitToken.
	overrides TokenOverrides

	// tokenIdentifier caches the identifier so the digest is computed only once.
	tokenIdentifier *TokenIdentifier

	// publicKeyOfTheSigner is the public key of the signing certificate(s).
	publicKeyOfTheSigner *PublicKey
	// signatureValidity is the status of the token's signature; IsSignedBy must be called
	// to establish it. Default: NOT_EVALUATED.
	signatureValidity enumerations.SignatureValidity
	// signatureInvalidityReason indicates why the token's signature is invalid.
	signatureInvalidityReason string
	// signatureAlgorithm is the algorithm used to sign the token.
	signatureAlgorithm enumerations.SignatureAlgorithm
}

// NewTokenBase instantiates the base state of a token with the Java default values.
// Port of the protected Token() constructor. The concrete token must still call InitToken.
func NewTokenBase() TokenBase {
	return TokenBase{
		signatureValidity:         enumerations.SignatureValidityNotEvaluated,
		signatureInvalidityReason: "",
	}
}

// InitToken registers the concrete token with its base so that the base can dispatch to the
// operations Java would reach through virtual dispatch. It must be called exactly once, by
// the concrete token's constructor, before any other method.
func (t *TokenBase) InitToken(overrides TokenOverrides) {
	t.overrides = overrides
}

// tokenBaseOverrides returns the registered overrides, panicking when the concrete token
// forgot to call InitToken.
func (t *TokenBase) tokenBaseOverrides() TokenOverrides {
	if t.overrides == nil {
		panic("Token was not initialised: the concrete token must call InitToken in its constructor")
	}
	return t.overrides
}

// IsSelfSigned reports whether the token is self-signed. For all tokens other than a
// CertificateToken this always returns false; the method exists so that the different
// tokens can be managed uniformly.
func (t *TokenBase) IsSelfSigned() bool {
	return false
}

// DSSID returns the DSS unique token identifier, building it on first use.
// Port of getDSSId(); Java narrows the return type to TokenIdentifier, which Go cannot
// express, so callers needing the token identifier assert on the result.
func (t *TokenBase) DSSID() Identifier {
	return t.dssTokenID()
}

// dssTokenID is the internal, narrowly typed counterpart of DSSID.
func (t *TokenBase) dssTokenID() *TokenIdentifier {
	if t.tokenIdentifier == nil {
		t.tokenIdentifier = t.tokenBaseOverrides().BuildTokenIdentifier()
	}
	return t.tokenIdentifier
}

// DSSIDAsString returns a string representation of the unique DSS token identifier.
// Port of getDSSIdAsString().
func (t *TokenBase) DSSIDAsString() string {
	return t.dssTokenID().AsXmlID()
}

// IsSignedByToken reports whether the token is signed by the given certificate token.
// Port of the isSignedBy(CertificateToken) overload; Go has no overloading, so the two
// isSignedBy methods carry different names.
func (t *TokenBase) IsSignedByToken(token *CertificateToken) bool {
	return t.IsSignedBy(token.PublicKey())
}

// IsSignedBy reports whether the token is signed by the given public key, remembering the
// signer on success. Port of the isSignedBy(PublicKey) overload.
//
// Upstream declares both isSignedBy overloads synchronized; the Go port is not
// goroutine-safe, matching the rest of the value objects in this package.
func (t *TokenBase) IsSignedBy(publicKey *PublicKey) bool {
	if t.publicKeyOfTheSigner != nil {
		return t.publicKeyOfTheSigner.Equals(publicKey)
	} else if enumerations.SignatureValidityValid == t.tokenBaseOverrides().CheckIsSignedBy(publicKey) {
		if !t.tokenBaseOverrides().IsSelfSigned() {
			t.publicKeyOfTheSigner = publicKey
		}
		return true
	}
	return false
}

// Abbreviation returns the DSS abbreviation of the token, used for debugging.
// Port of getAbbreviation().
func (t *TokenBase) Abbreviation() string {
	return "?"
}

// IssuerEntityKey returns the identifier of the entity key of the issuer of the current
// token, or nil when the signer has not been established by a successful CheckIsSignedBy.
// Port of getIssuerEntityKey().
func (t *TokenBase) IssuerEntityKey() *EntityIdentifier {
	if t.publicKeyOfTheSigner != nil && t.tokenBaseOverrides().IssuerX500Principal() != nil {
		return NewEntityIdentifierBuilder(t.publicKeyOfTheSigner, t.tokenBaseOverrides().IssuerX500Principal()).Build()
	}
	return nil
}

// SignatureAlgorithm returns the algorithm that was used to sign the token (e.g.
// RSA_SHA256). Port of getSignatureAlgorithm().
func (t *TokenBase) SignatureAlgorithm() enumerations.SignatureAlgorithm {
	return t.signatureAlgorithm
}

// SetSignatureAlgorithm sets the algorithm that was used to sign the token. Java writes the
// protected field directly from its subclasses; Go subclasses outside this package need this
// setter instead.
func (t *TokenBase) SetSignatureAlgorithm(signatureAlgorithm enumerations.SignatureAlgorithm) {
	t.signatureAlgorithm = signatureAlgorithm
}

// IsSignatureIntact reports whether the token's signature is intact. IsSignedBy must have
// been called first; the method returns false both when the check has not run and when the
// signer's public key does not match. Port of isSignatureIntact().
func (t *TokenBase) IsSignatureIntact() bool {
	return enumerations.SignatureValidityValid == t.signatureValidity
}

// IsValid reports whether the conditions corresponding to the token validity are met.
// Port of isValid().
func (t *TokenBase) IsValid() bool {
	return t.IsSignatureIntact()
}

// SignatureValidity returns the three-state status of the token's signature validity.
// Port of getSignatureValidity().
func (t *TokenBase) SignatureValidity() enumerations.SignatureValidity {
	return t.signatureValidity
}

// SetSignatureValidity sets the status of the token's signature validity; the Go counterpart
// of writing Java's protected signatureValidity field from a subclass.
func (t *TokenBase) SetSignatureValidity(signatureValidity enumerations.SignatureValidity) {
	t.signatureValidity = signatureValidity
}

// InvalidityReason returns the token invalidity reason when applicable, empty when the token
// is valid. Port of getInvalidityReason().
func (t *TokenBase) InvalidityReason() string {
	return t.signatureInvalidityReason
}

// SetInvalidityReason sets the token invalidity reason; the Go counterpart of writing Java's
// protected signatureInvalidityReason field from a subclass.
func (t *TokenBase) SetInvalidityReason(reason string) {
	t.signatureInvalidityReason = reason
}

// PublicKeyOfTheSigner returns the public key which signed this token.
// Port of getPublicKeyOfTheSigner().
func (t *TokenBase) PublicKeyOfTheSigner() *PublicKey {
	return t.publicKeyOfTheSigner
}

// SetPublicKeyOfTheSigner sets the public key which signed this token; the Go counterpart of
// writing Java's protected publicKeyOfTheSigner field from a subclass.
func (t *TokenBase) SetPublicKeyOfTheSigner(publicKey *PublicKey) {
	t.publicKeyOfTheSigner = publicKey
}

// Digest returns the digest value of the wrapped token for the requested algorithm.
// Port of getDigest(DigestAlgorithm); Java's DSSException for an unavailable algorithm
// becomes the returned error.
func (t *TokenBase) Digest(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	return t.dssTokenID().DigestValue(digestAlgorithm)
}
