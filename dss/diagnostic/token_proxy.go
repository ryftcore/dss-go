// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/TokenProxy.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// TokenProxy provides a user-friendly API for dealing with JAXB objects from a DiagnosticData.
// Every concrete wrapper (CertificateWrapper, SignatureWrapper, TimestampWrapper,
// RevocationWrapper) satisfies this interface by embedding AbstractTokenProxyBase and
// implementing its own Id/FoundCertificates/FoundRevocations/DigestMatchers overrides; see
// abstract_token_proxy.go.
type TokenProxy interface {
	// Id returns the unique identifier of the object. Port of getId().
	Id() string
	// IsSignatureIntact reports whether the signatureValue of the token is valid against the
	// identifier signing certificate's public key. Port of isSignatureIntact().
	IsSignatureIntact() bool
	// IsSignatureValid reports whether the signature and all signed data is cryptographically
	// correct. Port of isSignatureValid().
	IsSignatureValid() bool
	// SignatureAlgorithm returns the SignatureAlgorithm used to create the signatureValue. Port
	// of getSignatureAlgorithm().
	SignatureAlgorithm() enumerations.SignatureAlgorithm
	// DigestAlgorithm returns the DigestAlgorithm used to create the signatureValue. Port of
	// getDigestAlgorithm().
	DigestAlgorithm() enumerations.DigestAlgorithm
	// EncryptionAlgorithm returns the EncryptionAlgorithm used to create the signature. Port of
	// getEncryptionAlgorithm().
	EncryptionAlgorithm() enumerations.EncryptionAlgorithm
	// KeyLengthUsedToSignThisToken returns the length of the private key used to create the
	// signatureValue of the token. Port of getKeyLengthUsedToSignThisToken().
	KeyLengthUsedToSignThisToken() string
	// SigningCertificate returns the signing certificate of the token if identified. Port of
	// getSigningCertificate().
	SigningCertificate() *CertificateWrapper
	// IsSigningCertificateReferencePresent reports whether a reference to the SigningCertificate
	// is present within the token (used for signatures and timestamps). Port of
	// isSigningCertificateReferencePresent().
	IsSigningCertificateReferencePresent() bool
	// IsSigningCertificateReferenceUnique reports whether the reference to the signing
	// certificate is unique and present only once. Port of
	// isSigningCertificateReferenceUnique().
	IsSigningCertificateReferenceUnique() bool
	// SigningCertificateReference returns the reference to the signing certificate present
	// within the token (for signature or timestamp). Port of getSigningCertificateReference().
	SigningCertificateReference() *CertificateRefWrapper
	// SigningCertificateReferences returns a list of all references to the signing certificate
	// present within the token (for signature or timestamp). Port of
	// getSigningCertificateReferences().
	SigningCertificateReferences() []*CertificateRefWrapper
	// SigningCertificatePublicKey returns the public key binaries linked to a private key used
	// to create the signature, when a signing-certificate is not available. Port of
	// getSigningCertificatePublicKey().
	SigningCertificatePublicKey() []byte
	// CertificateChain returns the certificate chain. Port of getCertificateChain().
	CertificateChain() []*CertificateWrapper
	// IsTrustedChain reports whether the certificate chain is trusted. Port of isTrustedChain().
	IsTrustedChain() bool
	// DigestMatchers returns a list of DigestMatchers used in the validation process for a
	// signature or timestamp. Port of getDigestMatchers().
	DigestMatchers() []*jaxb.XmlDigestMatcher
	// FoundCertificates returns a FoundCertificatesProxy to access embedded certificates. Port
	// of foundCertificates().
	FoundCertificates() *FoundCertificatesProxy
	// FoundRevocations returns a FoundRevocationsProxy to access embedded revocation data. Port
	// of foundRevocations().
	FoundRevocations() *FoundRevocationsProxy
}
