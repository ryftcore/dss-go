// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CertificateRef.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi (see PORTING_PLAN.md), so the
// type keeps its Java name unqualified.
package spi

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// CertificateRef represents a Certificate Reference entry extracted from a signature.
type CertificateRef struct {
	// certDigest is the digest of the certificate.
	certDigest model.Digest

	// signerIdentifier is the ASN.1 SignerId (signature or timestamp).
	signerIdentifier *SignerIdentifier

	// responderId is the ResponderId in case of an OCSP response.
	responderId *ResponderId

	// kid is the key identifier.
	kid string

	// x509Url identifies the location URI of the X.509 public key certificate.
	x509Url string

	// publicKey is the public key of the signer's certificate.
	publicKey *model.PublicKey

	// identifier is the unique identifier of the reference, computed lazily.
	identifier model.Identifier
}

// NewCertificateRef instantiates the object with null values. Port of the default constructor.
func NewCertificateRef() *CertificateRef {
	return &CertificateRef{}
}

// CertDigest gets the certificate digest. Port of getCertDigest().
func (r *CertificateRef) CertDigest() model.Digest {
	return r.certDigest
}

// SetCertDigest sets the certificate digest. Port of setCertDigest(Digest).
func (r *CertificateRef) SetCertDigest(certDigest model.Digest) {
	r.certDigest = certDigest
}

// CertificateIdentifier gets the SignerIdentifier (for a reference extracted from a signature
// or timestamp, when present). Port of getCertificateIdentifier().
func (r *CertificateRef) CertificateIdentifier() *SignerIdentifier {
	return r.signerIdentifier
}

// SetCertificateIdentifier sets the SignerIdentifier. Port of setCertificateIdentifier(SignerIdentifier).
func (r *CertificateRef) SetCertificateIdentifier(signerIdentifier *SignerIdentifier) {
	r.signerIdentifier = signerIdentifier
}

// ResponderId gets the ResponderId (for a reference extracted from an OCSP response).
// Port of getResponderId().
func (r *CertificateRef) ResponderId() *ResponderId {
	return r.responderId
}

// SetResponderId sets the ResponderId. Port of setResponderId(ResponderId).
func (r *CertificateRef) SetResponderId(responderId *ResponderId) {
	r.responderId = responderId
}

// Kid gets a key identifier (KID). Port of getKid().
func (r *CertificateRef) Kid() string {
	return r.kid
}

// SetKid sets a key identifier (KID). Port of setKid(String).
func (r *CertificateRef) SetKid(kid string) {
	r.kid = kid
}

// X509Url gets the X.509 Public Key Certificate location URL. Port of getX509Url().
func (r *CertificateRef) X509Url() string {
	return r.x509Url
}

// SetX509Url sets the X.509 Public Key Certificate location URL. Port of setX509Url(String).
func (r *CertificateRef) SetX509Url(x509Url string) {
	r.x509Url = x509Url
}

// PublicKey gets the public key of the signer certificate. Port of getPublicKey().
func (r *CertificateRef) PublicKey() *model.PublicKey {
	return r.publicKey
}

// SetPublicKey sets the public key of the signer's certificate. Port of setPublicKey(PublicKey).
func (r *CertificateRef) SetPublicKey(publicKey *model.PublicKey) {
	r.publicKey = publicKey
}

// DSSID returns the certificate reference identifier, building it lazily on first use.
// Port of getDSSId().
//
// NewCertificateRefIdentifier's data-dependent DSSException ("One of [certDigest,
// publicKeyDigest, issuerInfo, kid, x509Uri, publicKey] must be defined for a
// CertificateRef!") is unchecked in Java and propagates straight out of getDSSId(); since
// this method's signature (like Java's) carries no error, it is reproduced as a panic here.
func (r *CertificateRef) DSSID() model.Identifier {
	if r.identifier == nil {
		identifier, err := NewCertificateRefIdentifier(r)
		if err != nil {
			panic(err.Error())
		}
		r.identifier = identifier
	}
	return r.identifier
}

// DSSIDAsString returns the certificate reference String id. Port of getDSSIdAsString().
func (r *CertificateRef) DSSIDAsString() string {
	return r.DSSID().AsXmlID()
}

// String returns a string representation of the certificate reference. Port of toString().
func (r *CertificateRef) String() string {
	return "CertificateRef [" +
		"certDigest=" + r.certDigest.String() +
		", signerIdentifier=" + certificateRefSignerIdentifierString(r.signerIdentifier) +
		", responderId=" + certificateRefResponderIdString(r.responderId) +
		", kid='" + r.kid + "'" +
		", x509Url='" + r.x509Url + "'" +
		", publicKey=" + certificateRefPublicKeyString(r.publicKey) +
		", identifier=" + certificateRefIdentifierString(r.identifier) +
		"]"
}

// certificateRefSignerIdentifierString renders a possibly nil SignerIdentifier the way Java's
// string concatenation does.
func certificateRefSignerIdentifierString(signerIdentifier *SignerIdentifier) string {
	if signerIdentifier == nil {
		return "null"
	}
	return signerIdentifier.String()
}

// certificateRefResponderIdString renders a possibly nil ResponderId.
//
// DEVIATION: Java's ResponderId has no toString() override, so string concatenation there
// calls Object's default (implementation-dependent, identity-hash-based) representation -
// not meaningful to reproduce. As with certificateRefPublicKeyString, a stable field-based
// rendering is used instead; a cosmetic-only deviation with no effect on equality/DSS-Id.
func certificateRefResponderIdString(responderId *ResponderId) string {
	if responderId == nil {
		return "null"
	}
	principal := "null"
	if p := responderId.X500Principal(); p != nil {
		principal = p.String()
	}
	return "ResponderId [subjectX500Principal=" + principal + ", ski=" + utils.ToHex(responderId.Ski()) + "]"
}

// certificateRefPublicKeyString renders a possibly nil public key. Java interpolates the
// java.security.PublicKey's default Object#toString() (implementation-dependent, algorithm +
// identity hash); model.PublicKey has no such method, so a stable base64-of-the-encoding
// rendering is used instead - a documented cosmetic deviation with no effect on equality.
func certificateRefPublicKeyString(publicKey *model.PublicKey) string {
	if publicKey == nil {
		return "null"
	}
	return publicKey.Algorithm() + " PublicKey [" + utils.ToBase64(publicKey.Encoded()) + "]"
}

// certificateRefIdentifierString renders a possibly nil identifier the way Java's string
// concatenation does.
func certificateRefIdentifierString(identifier model.Identifier) string {
	if identifier == nil {
		return "null"
	}
	return identifier.String()
}

// Equals reports whether both references carry the same digest, signer identifier, responder
// Id, kid, X.509 URL and public key. Port of equals(Object).
//
// NOTE: hashCode() has no Go counterpart; upstream needs it only to key the JDK hash
// collections, which the port replaces with slices/maps keyed on Equals or DSSIDAsString().
func (r *CertificateRef) Equals(other *CertificateRef) bool {
	if r == other {
		return true
	}
	if other == nil {
		return false
	}
	if !r.certDigest.Equals(other.certDigest) {
		return false
	}
	if !certificateRefSignerIdentifiersEqual(r.signerIdentifier, other.signerIdentifier) {
		return false
	}
	if !certificateRefResponderIdsEqual(r.responderId, other.responderId) {
		return false
	}
	if r.kid != other.kid {
		return false
	}
	if r.x509Url != other.x509Url {
		return false
	}
	return certificateRefPublicKeysEqual(r.publicKey, other.publicKey)
}

// certificateRefSignerIdentifiersEqual is Objects.equals for two possibly-nil SignerIdentifiers.
func certificateRefSignerIdentifiersEqual(first, second *SignerIdentifier) bool {
	if first == nil || second == nil {
		return first == second
	}
	return first.Equals(second)
}

// certificateRefResponderIdsEqual is Objects.equals for two possibly-nil ResponderIds.
func certificateRefResponderIdsEqual(first, second *ResponderId) bool {
	if first == nil || second == nil {
		return first == second
	}
	return first.Equals(second)
}

// certificateRefPublicKeysEqual is Objects.equals for two possibly-nil public keys.
func certificateRefPublicKeysEqual(first, second *model.PublicKey) bool {
	if first == nil || second == nil {
		return first == second
	}
	return first.Equals(second)
}

// compile-time assertion: a CertificateRef is an IdentifierBasedObject.
var _ model.IdentifierBasedObject = (*CertificateRef)(nil)
