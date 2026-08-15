// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/AbstractTokenProxy.java (DSS 6.5.RC1).
package diagnostic

import (
	"reflect"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// AbstractTokenProxyOverrides declares the operations Java's abstract AbstractTokenProxy class
// leaves abstract, or provides a default body for that a subclass may override, and that the
// base implementation itself calls back into (getSigningCertificateReference() calling
// foundCertificates(), for instance). It stands in for the virtual dispatch a Java abstract
// class gets for free; a concrete wrapper registers itself with
// AbstractTokenProxyBase.InitTokenProxy so the base can reach it. Every concrete wrapper embeds
// AbstractTokenProxyBase and must implement the full interface itself, even for the members it
// does not override in Java: those implementations simply reproduce the AbstractTokenProxy
// default body (see CertificateWrapper for an example that does not override any of them).
type AbstractTokenProxyOverrides interface {
	// CurrentBasicSignature returns a basic signature validation. Port of the abstract
	// protected getCurrentBasicSignature().
	CurrentBasicSignature() *jaxb.XmlBasicSignature
	// CurrentCertificateChain returns the token's certificate chain. Port of the abstract
	// protected getCurrentCertificateChain().
	CurrentCertificateChain() []*jaxb.XmlChainItem
	// CurrentSigningCertificate returns the signing certificate of the token. Port of the
	// abstract protected getCurrentSigningCertificate().
	CurrentSigningCertificate() *jaxb.XmlSigningCertificate
	// FoundCertificates returns a collection of certificates embedded within the token. Port
	// of the interface method foundCertificates() (TokenProxy), overridable by subclasses.
	FoundCertificates() *FoundCertificatesProxy
	// FoundRevocations returns a collection of revocation data embedded within the token.
	// Port of the interface method foundRevocations() (TokenProxy), overridable by subclasses.
	FoundRevocations() *FoundRevocationsProxy
	// DigestMatchers returns the digest matchers of the token. Port of the interface method
	// getDigestMatchers() (TokenProxy), overridable by subclasses.
	DigestMatchers() []*jaxb.XmlDigestMatcher
	// Binaries returns binaries of the token, when present. Port of the abstract
	// public getBinaries().
	Binaries() []byte
	// Id returns the unique identifier of the token. Port of the interface method getId()
	// (TokenProxy), left abstract by AbstractTokenProxy.
	Id() string
}

// AbstractTokenProxyBase carries the state and the concrete behaviour of the Java abstract
// class AbstractTokenProxy. Concrete wrappers embed it and register themselves with
// InitTokenProxy.
type AbstractTokenProxyBase struct {
	// overrides points back at the concrete wrapper; see InitTokenProxy.
	overrides AbstractTokenProxyOverrides
}

// InitTokenProxy registers the concrete wrapper with its base so that the base can dispatch to
// the operations Java would reach through virtual dispatch. It must be called exactly once, by
// the concrete wrapper's constructor, before any other method.
func (a *AbstractTokenProxyBase) InitTokenProxy(overrides AbstractTokenProxyOverrides) {
	a.overrides = overrides
}

// tokenProxyOverrides returns the registered overrides, panicking when the concrete wrapper
// forgot to call InitTokenProxy.
func (a *AbstractTokenProxyBase) tokenProxyOverrides() AbstractTokenProxyOverrides {
	if a.overrides == nil {
		panic("AbstractTokenProxy was not initialised: the concrete wrapper must call InitTokenProxy in its constructor")
	}
	return a.overrides
}

// DefaultFoundCertificates is the AbstractTokenProxy default body for foundCertificates():
// an empty proxy over a nil XmlFoundCertificates. Concrete wrappers that do not override
// foundCertificates() in Java call this from their own FoundCertificates() implementation.
func DefaultFoundCertificates() *FoundCertificatesProxy {
	return NewFoundCertificatesProxy(nil)
}

// DefaultFoundRevocations is the AbstractTokenProxy default body for foundRevocations(): an
// empty proxy over a nil XmlFoundRevocations. Concrete wrappers that do not override
// foundRevocations() in Java call this from their own FoundRevocations() implementation.
func DefaultFoundRevocations() *FoundRevocationsProxy {
	return NewFoundRevocationsProxy(nil)
}

// CertificateChain returns the token's certificate chain as a list of CertificateWrapper.
// Port of getCertificateChain().
func (a *AbstractTokenProxyBase) CertificateChain() []*CertificateWrapper {
	var result []*CertificateWrapper
	certificateChain := a.tokenProxyOverrides().CurrentCertificateChain()
	if certificateChain != nil {
		for _, xmlChainCertificate := range certificateChain {
			if xmlChainCertificate.Certificate != nil {
				result = append(result, NewCertificateWrapper(xmlChainCertificate.Certificate))
			}
		}
	}
	return result
}

// IsSignatureIntact reports whether the token's basic signature is intact. Port of
// isSignatureIntact().
func (a *AbstractTokenProxyBase) IsSignatureIntact() bool {
	basicSignature := a.tokenProxyOverrides().CurrentBasicSignature()
	if basicSignature != nil {
		return basicSignature.SignatureIntact != nil && *basicSignature.SignatureIntact
	}
	return false
}

// IsSignatureValid reports whether the token's basic signature is valid. Port of
// isSignatureValid().
func (a *AbstractTokenProxyBase) IsSignatureValid() bool {
	basicSignature := a.tokenProxyOverrides().CurrentBasicSignature()
	if basicSignature != nil {
		return basicSignature.SignatureValid != nil && *basicSignature.SignatureValid
	}
	return false
}

// SignatureAlgorithm returns the signature algorithm combining the encryption and digest
// algorithms used to sign the token. Port of getSignatureAlgorithm().
func (a *AbstractTokenProxyBase) SignatureAlgorithm() enumerations.SignatureAlgorithm {
	encryptionAlgorithm := a.EncryptionAlgorithm()
	digestAlgorithm := a.DigestAlgorithm()
	if encryptionAlgorithm != "" && digestAlgorithm != "" {
		return enumerations.SignatureAlgorithmGetAlgorithm(encryptionAlgorithm, digestAlgorithm)
	}
	return ""
}

// EncryptionAlgorithm returns the encryption algorithm used to sign the token. Port of
// getEncryptionAlgorithm().
func (a *AbstractTokenProxyBase) EncryptionAlgorithm() enumerations.EncryptionAlgorithm {
	basicSignature := a.tokenProxyOverrides().CurrentBasicSignature()
	if basicSignature != nil && basicSignature.EncryptionAlgoUsedToSignThisToken != nil {
		return enumerations.EncryptionAlgorithm(*basicSignature.EncryptionAlgoUsedToSignThisToken)
	}
	return ""
}

// DigestAlgorithm returns the digest algorithm used to sign the token. Port of
// getDigestAlgorithm().
func (a *AbstractTokenProxyBase) DigestAlgorithm() enumerations.DigestAlgorithm {
	basicSignature := a.tokenProxyOverrides().CurrentBasicSignature()
	if basicSignature != nil && basicSignature.DigestAlgoUsedToSignThisToken != nil {
		return enumerations.DigestAlgorithm(*basicSignature.DigestAlgoUsedToSignThisToken)
	}
	return ""
}

// KeyLengthUsedToSignThisToken returns the key length used to sign the token. Port of
// getKeyLengthUsedToSignThisToken().
func (a *AbstractTokenProxyBase) KeyLengthUsedToSignThisToken() string {
	basicSignature := a.tokenProxyOverrides().CurrentBasicSignature()
	if basicSignature != nil && basicSignature.KeyLengthUsedToSignThisToken != nil {
		return *basicSignature.KeyLengthUsedToSignThisToken
	}
	return ""
}

// SigningCertificate returns the signing certificate of the token. Port of
// getSigningCertificate().
func (a *AbstractTokenProxyBase) SigningCertificate() *CertificateWrapper {
	currentSigningCertificate := a.tokenProxyOverrides().CurrentSigningCertificate()
	if currentSigningCertificate != nil && currentSigningCertificate.Certificate != nil {
		return NewCertificateWrapper(currentSigningCertificate.Certificate)
	}
	return nil
}

// SigningCertificatePublicKey returns the public key of the signing certificate. Port of
// getSigningCertificatePublicKey().
func (a *AbstractTokenProxyBase) SigningCertificatePublicKey() []byte {
	currentSigningCertificate := a.tokenProxyOverrides().CurrentSigningCertificate()
	if currentSigningCertificate != nil && currentSigningCertificate.PublicKey != nil {
		return []byte(*currentSigningCertificate.PublicKey)
	}
	return nil
}

// IsSigningCertificateReferencePresent reports whether a signing certificate reference is
// present. Port of isSigningCertificateReferencePresent().
func (a *AbstractTokenProxyBase) IsSigningCertificateReferencePresent() bool {
	return len(a.SigningCertificateReferences()) != 0
}

// IsSigningCertificateReferenceUnique reports whether exactly one signing certificate
// reference is present. Port of isSigningCertificateReferenceUnique().
func (a *AbstractTokenProxyBase) IsSigningCertificateReferenceUnique() bool {
	return len(a.SigningCertificateReferences()) == 1
}

// SigningCertificateReference returns a reference matching the signing certificate, or the
// first orphan signing certificate reference when the signing certificate itself was not
// found. Port of getSigningCertificateReference().
func (a *AbstractTokenProxyBase) SigningCertificateReference() *CertificateRefWrapper {
	signingCertificateReferences := a.tokenProxyOverrides().FoundCertificates().
		RelatedCertificateRefsByRefOrigin(enumerations.CertificateRefOrigin_SIGNING_CERTIFICATE)
	if len(signingCertificateReferences) != 0 {
		// return a reference matching a signing certificate
		signingCertificate := a.SigningCertificate()
		if signingCertificate != nil {
			return a.getCertificateReferenceOfReferenceOriginType(signingCertificate, enumerations.CertificateRefOrigin_SIGNING_CERTIFICATE)
		}
	} else {
		orphanSigningCertificateReferences := a.tokenProxyOverrides().FoundCertificates().
			OrphanCertificateRefsByRefOrigin(enumerations.CertificateRefOrigin_SIGNING_CERTIFICATE)
		if len(orphanSigningCertificateReferences) != 0 {
			return orphanSigningCertificateReferences[0]
		}
	}
	return nil
}

func (a *AbstractTokenProxyBase) getCertificateReferenceOfReferenceOriginType(certificate *CertificateWrapper, refOrigin enumerations.CertificateRefOrigin) *CertificateRefWrapper {
	for _, relatedCertificate := range a.tokenProxyOverrides().FoundCertificates().RelatedCertificates() {
		signCertRefs := relatedCertificate.References()
		if certificate.Id() == relatedCertificate.Id() && len(signCertRefs) != 0 {
			for _, signCertRef := range signCertRefs {
				if refOrigin == signCertRef.Origin() {
					return signCertRef
				}
			}
		}
	}
	return nil
}

// SigningCertificateReferences returns all signing certificate references, related and
// orphan. Port of getSigningCertificateReferences().
func (a *AbstractTokenProxyBase) SigningCertificateReferences() []*CertificateRefWrapper {
	var certificateRefs []*CertificateRefWrapper
	certificateRefs = append(certificateRefs, a.tokenProxyOverrides().FoundCertificates().
		RelatedCertificateRefsByRefOrigin(enumerations.CertificateRefOrigin_SIGNING_CERTIFICATE)...)
	certificateRefs = append(certificateRefs, a.tokenProxyOverrides().FoundCertificates().
		OrphanCertificateRefsByRefOrigin(enumerations.CertificateRefOrigin_SIGNING_CERTIFICATE)...)
	return certificateRefs
}

// IsTrustedChain reports whether any certificate of the token's certificate chain is trusted.
// Port of isTrustedChain().
func (a *AbstractTokenProxyBase) IsTrustedChain() bool {
	for _, certificate := range a.CertificateChain() {
		if certificate.IsTrusted() {
			return true
		}
	}
	return false
}

// IsCertificateChainFromTrustedStore reports whether the certificate chain is trusted from a
// Trusted Store. NOTE: not from a Trusted List!
// Port of isCertificateChainFromTrustedStore().
func (a *AbstractTokenProxyBase) IsCertificateChainFromTrustedStore() bool {
	for _, certificate := range a.CertificateChain() {
		for _, source := range certificate.Sources() {
			if source == enumerations.CertificateSourceType_TRUSTED_STORE {
				return true
			}
		}
	}
	return false
}

// String returns a string representation of the token. Port of toString().
func (a *AbstractTokenProxyBase) String() string {
	return "Token Id='" + a.tokenProxyOverrides().Id() + "'"
}

// Equals reports whether other wraps a token of the same concrete wrapper type and the same
// Id. Port of equals(Object obj): Java compares getClass() != obj.getClass(); since Go has no
// class hierarchy, the concrete wrapper's dynamic type is compared via reflection instead.
func (a *AbstractTokenProxyBase) Equals(other AbstractTokenProxyOverrides) bool {
	if other == nil {
		return false
	}
	overrides := a.tokenProxyOverrides()
	if reflect.TypeOf(overrides) != reflect.TypeOf(other) {
		return false
	}
	return overrides.Id() == other.Id()
}
