// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/RevocationWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"math/big"
	"time"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// RevocationWrapper contains common revocation information.
type RevocationWrapper struct {
	AbstractTokenProxyBase

	// revocation is the wrapped XmlRevocation.
	revocation *jaxb.XmlRevocation
}

// NewRevocationWrapper is the default constructor. Port of RevocationWrapper(XmlRevocation);
// panics per Objects.requireNonNull(revocation, "XMLRevocation cannot be null!").
func NewRevocationWrapper(revocation *jaxb.XmlRevocation) *RevocationWrapper {
	if revocation == nil {
		panic("XMLRevocation cannot be null!")
	}
	w := &RevocationWrapper{revocation: revocation}
	w.InitTokenProxy(w)
	return w
}

// Id is the AbstractTokenProxy override. Port of getId().
func (w *RevocationWrapper) Id() string {
	if w.revocation.Id != nil {
		return string(*w.revocation.Id)
	}
	return ""
}

// CurrentBasicSignature is the AbstractTokenProxy override. Port of getCurrentBasicSignature().
func (w *RevocationWrapper) CurrentBasicSignature() *jaxb.XmlBasicSignature {
	return w.revocation.BasicSignature
}

// CurrentCertificateChain is the AbstractTokenProxy override. Port of
// getCurrentCertificateChain().
func (w *RevocationWrapper) CurrentCertificateChain() []*jaxb.XmlChainItem {
	return w.revocation.CertificateChain.All()
}

// CurrentSigningCertificate is the AbstractTokenProxy override. Port of
// getCurrentSigningCertificate().
func (w *RevocationWrapper) CurrentSigningCertificate() *jaxb.XmlSigningCertificate {
	return w.revocation.SigningCertificate
}

// FoundCertificates is the AbstractTokenProxy override. Port of foundCertificates().
func (w *RevocationWrapper) FoundCertificates() *FoundCertificatesProxy {
	return NewFoundCertificatesProxy(w.revocation.FoundCertificates)
}

// FoundRevocations is the AbstractTokenProxy default (not overridden in Java). Port of
// foundRevocations().
func (w *RevocationWrapper) FoundRevocations() *FoundRevocationsProxy {
	return DefaultFoundRevocations()
}

// DigestMatchers is the AbstractTokenProxy default (not overridden in Java). Port of
// getDigestMatchers().
func (w *RevocationWrapper) DigestMatchers() []*jaxb.XmlDigestMatcher {
	return nil
}

// ProductionDate returns the revocation data production time. Port of getProductionDate().
func (w *RevocationWrapper) ProductionDate() *time.Time {
	if w.revocation.ProductionDate != nil {
		t := w.revocation.ProductionDate.Time()
		return &t
	}
	return nil
}

// ThisUpdate returns the revocation data ThisUpdate time. Port of getThisUpdate().
func (w *RevocationWrapper) ThisUpdate() *time.Time {
	if w.revocation.ThisUpdate != nil {
		t := w.revocation.ThisUpdate.Time()
		return &t
	}
	return nil
}

// NextUpdate returns the revocation data NextUpdate time. Port of getNextUpdate().
func (w *RevocationWrapper) NextUpdate() *time.Time {
	if w.revocation.NextUpdate != nil {
		t := w.revocation.NextUpdate.Time()
		return &t
	}
	return nil
}

// CRLNumber returns the value of CRLNumber extension, when present. NOTE: CRL only. Port of
// getCRLNumber().
func (w *RevocationWrapper) CRLNumber() *big.Int {
	return w.revocation.CRLNumber.BigInt()
}

// ExpiredCertsOnCRL returns the expired-certs-on-crl attribute time, when present. Port of
// getExpiredCertsOnCRL().
func (w *RevocationWrapper) ExpiredCertsOnCRL() *time.Time {
	if w.revocation.ExpiredCertsOnCRL != nil {
		t := w.revocation.ExpiredCertsOnCRL.Time()
		return &t
	}
	return nil
}

// ArchiveCutOff returns the archive-cut-off attribute time, when present. Port of
// getArchiveCutOff().
func (w *RevocationWrapper) ArchiveCutOff() *time.Time {
	if w.revocation.ArchiveCutOff != nil {
		t := w.revocation.ArchiveCutOff.Time()
		return &t
	}
	return nil
}

// IsCertHashExtensionPresent gets if a certHash extension if present. Port of
// isCertHashExtensionPresent().
func (w *RevocationWrapper) IsCertHashExtensionPresent() bool {
	return w.revocation.CertHashExtensionPresent != nil && *w.revocation.CertHashExtensionPresent
}

// IsCertHashExtensionMatch gets if a certHash extension matches to the hash of the concerned
// certificate. Port of isCertHashExtensionMatch().
func (w *RevocationWrapper) IsCertHashExtensionMatch() bool {
	return w.revocation.CertHashExtensionMatch != nil && *w.revocation.CertHashExtensionMatch
}

// Origin returns the origin of the revocation token. Port of getOrigin().
func (w *RevocationWrapper) Origin() enumerations.RevocationOrigin {
	if w.revocation.Origin != nil {
		return enumerations.RevocationOrigin(*w.revocation.Origin)
	}
	return ""
}

// RevocationType returns the revocation data type. Port of getRevocationType().
func (w *RevocationWrapper) RevocationType() enumerations.RevocationType {
	if w.revocation.Type != nil {
		return enumerations.RevocationType(*w.revocation.Type)
	}
	return ""
}

// Binaries is the AbstractTokenProxy override. Port of getBinaries().
func (w *RevocationWrapper) Binaries() []byte {
	if w.revocation.Base64Encoded != nil {
		return []byte(*w.revocation.Base64Encoded)
	}
	return nil
}

// DigestAlgoAndValue returns the digest of the revocation token. Port of
// getDigestAlgoAndValue().
func (w *RevocationWrapper) DigestAlgoAndValue() *jaxb.XmlDigestAlgoAndValue {
	return w.revocation.DigestAlgoAndValue
}

// IsInternalRevocationOrigin returns true if the Revocation data was obtained from a signature
// container. Port of isInternalRevocationOrigin().
func (w *RevocationWrapper) IsInternalRevocationOrigin() bool {
	originType := w.Origin()
	if originType != "" {
		return w.Origin().IsInternalOrigin()
	}
	return false
}

// SourceAddress returns the remote URI used to obtain the revocation data. Port of
// getSourceAddress().
func (w *RevocationWrapper) SourceAddress() string {
	if w.revocation.SourceAddress != nil {
		return *w.revocation.SourceAddress
	}
	return ""
}

// Equals reports whether other wraps a revocation token with the same Id. Port of
// equals(Object): Java's `!(obj instanceof RevocationWrapper)` check accepts any
// AbstractTokenProxy, comparing only Id - unlike AbstractTokenProxy.equals() it does NOT
// additionally compare getClass(), so this compares only by Id, faithfully. hashCode() has no
// Go equivalent (nothing here keys a hash-based collection on a RevocationWrapper) and is
// dropped, matching the omission pattern documented elsewhere in this port (see
// model/certificate_token.go).
func (w *RevocationWrapper) Equals(other TokenProxy) bool {
	if other == nil {
		return false
	}
	return w.Id() == other.Id()
}
