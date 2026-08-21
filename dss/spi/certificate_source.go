// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CertificateSource.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: this file references CertificateRef, a spi.x509 type flattened into
// this package and ported in a sibling chunk of phase 2a; see certificate_ref_identifier.go
// for its assumed shape.
//
// Java's Set<CertificateToken> return values are ported as map[string]*model.CertificateToken
// keyed by the token's DSSIDAsString(), matching the convention equivalent_certificates_entity.go
// (chunk X509-B) already established for the same Java Set<CertificateToken>-over-hashCode
// problem: CertificateToken has no comparable Go representation usable directly as a map key
// (its identity digest is a []byte), so callers needing equals()-based deduplication key on the
// identifier string instead of on the pointer.
package spi

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// CertificateSource provides an abstraction for accessing a certificate, regardless of the
// source: the validation of a certificate requires access to some other certificates from
// multiple sources (Trusted List, Trust Store, the signature itself).
type CertificateSource interface {
	// AddCertificate manually adds any certificate to the source. The type of the source is
	// automatically set by each specific implementation. Port of addCertificate(CertificateToken).
	AddCertificate(certificate *model.CertificateToken) *model.CertificateToken

	// CertificateSourceType returns the certificate source type associated with the
	// implementation. Port of getCertificateSourceType().
	CertificateSourceType() enumerations.CertificateSourceType

	// Certificates retrieves the unmodifiable list of all certificate tokens from this source.
	// Port of getCertificates().
	Certificates() []*model.CertificateToken

	// IsTrusted checks if a given certificate is trusted. Port of isTrusted(CertificateToken).
	IsTrusted(certificateToken *model.CertificateToken) bool

	// IsTrustedAtTime checks if a given certificate is trusted at controlTime.
	// Port of isTrustedAtTime(CertificateToken, Date).
	IsTrustedAtTime(certificateToken *model.CertificateToken, controlTime time.Time) bool

	// IsKnown checks if a given certificate is known in the current source.
	// Port of isKnown(CertificateToken).
	IsKnown(certificateToken *model.CertificateToken) bool

	// BySubject returns the certificates with the same subjectDN, keyed by DSSIDAsString();
	// empty when no match is found. Port of getBySubject(X500PrincipalHelper).
	BySubject(subject *model.X500PrincipalHelper) map[string]*model.CertificateToken

	// BySignerIdentifier returns the certificates with the given SignerIdentifier, keyed by
	// DSSIDAsString(); empty when no match is found. Port of getBySignerIdentifier(SignerIdentifier).
	BySignerIdentifier(signerIdentifier *SignerIdentifier) map[string]*model.CertificateToken

	// ByCertificateDigest returns the certificates with the given Digest, keyed by
	// DSSIDAsString(). Port of getByCertificateDigest(Digest).
	ByCertificateDigest(digest model.Digest) map[string]*model.CertificateToken

	// ByPublicKey returns the certificates with the given PublicKey, keyed by DSSIDAsString().
	// Port of getByPublicKey(PublicKey).
	ByPublicKey(publicKey *model.PublicKey) map[string]*model.CertificateToken

	// ByEntityKey returns the certificates with the given EntityIdentifier, keyed by
	// DSSIDAsString(). Port of getByEntityKey(EntityIdentifier).
	ByEntityKey(entityKey *model.EntityIdentifier) map[string]*model.CertificateToken

	// BySki returns the certificates with the given SKI (SubjectKeyIdentifier, SHA-1 of the
	// PublicKey), keyed by DSSIDAsString(). Port of getBySki(byte[]).
	BySki(ski []byte) map[string]*model.CertificateToken

	// FindTokensFromCertRef returns the certificate tokens for the provided CertificateRef,
	// keyed by DSSIDAsString(). Port of findTokensFromCertRef(CertificateRef).
	FindTokensFromCertRef(certificateRef *CertificateRef) map[string]*model.CertificateToken

	// Entities returns the certificates grouped by their public keys. Port of getEntities().
	Entities() []CertificateSourceEntity

	// IsAllSelfSigned checks if all certificates are self-signed. Port of isAllSelfSigned().
	IsAllSelfSigned() bool

	// IsCertificateSourceEqual checks if the current and the given CertificateSources contain
	// the same certificate tokens. Port of isCertificateSourceEqual(CertificateSource).
	IsCertificateSourceEqual(certificateSource CertificateSource) bool

	// IsCertificateSourceEquivalent checks if the current and the given CertificateSources
	// contain the same entity keys. Port of isCertificateSourceEquivalent(CertificateSource).
	IsCertificateSourceEquivalent(certificateSource CertificateSource) bool
}
