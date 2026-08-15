// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/EAARevocationTokenWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"math/big"
	"time"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// EAARevocationTokenWrapper wraps validation information of the EAA revocation token.
type EAARevocationTokenWrapper struct {
	AbstractTokenProxyBase

	// eaaStatusToken is the wrapped XmlEAARevocationToken.
	eaaStatusToken *jaxb.XmlEAARevocationToken
}

// NewEAARevocationTokenWrapper is the default constructor. Port of
// EAARevocationTokenWrapper(XmlEAARevocationToken); panics per
// Objects.requireNonNull(eaaRevocationToken, "XmlEAARevocationToken cannot be null!").
func NewEAARevocationTokenWrapper(eaaRevocationToken *jaxb.XmlEAARevocationToken) *EAARevocationTokenWrapper {
	if eaaRevocationToken == nil {
		panic("XmlEAARevocationToken cannot be null!")
	}
	w := &EAARevocationTokenWrapper{eaaStatusToken: eaaRevocationToken}
	w.InitTokenProxy(w)
	return w
}

// CurrentBasicSignature is the AbstractTokenProxy override. Port of
// getCurrentBasicSignature().
func (w *EAARevocationTokenWrapper) CurrentBasicSignature() *jaxb.XmlBasicSignature {
	return w.eaaStatusToken.BasicSignature
}

// CurrentCertificateChain is the AbstractTokenProxy override. Port of
// getCurrentCertificateChain().
func (w *EAARevocationTokenWrapper) CurrentCertificateChain() []*jaxb.XmlChainItem {
	return w.eaaStatusToken.CertificateChain.All()
}

// CurrentSigningCertificate is the AbstractTokenProxy override. Port of
// getCurrentSigningCertificate().
func (w *EAARevocationTokenWrapper) CurrentSigningCertificate() *jaxb.XmlSigningCertificate {
	return w.eaaStatusToken.SigningCertificate
}

// FoundCertificates returns a proxy to access embedded certificates. Port of
// foundCertificates() (overridden).
func (w *EAARevocationTokenWrapper) FoundCertificates() *FoundCertificatesProxy {
	return NewFoundCertificatesProxy(w.eaaStatusToken.FoundCertificates)
}

// FoundRevocations is the AbstractTokenProxy default (not overridden in Java). Port of
// foundRevocations().
func (w *EAARevocationTokenWrapper) FoundRevocations() *FoundRevocationsProxy {
	return DefaultFoundRevocations()
}

// DigestMatchers is the AbstractTokenProxy default (not overridden in Java). Port of
// getDigestMatchers().
func (w *EAARevocationTokenWrapper) DigestMatchers() []*jaxb.XmlDigestMatcher {
	return nil
}

// Origin gets origin of the EAA revocation token (e.g. EXTERNAL or CACHED). Port of
// getOrigin().
func (w *EAARevocationTokenWrapper) Origin() enumerations.EAARevocationOrigin {
	if w.eaaStatusToken.Origin != nil {
		return enumerations.EAARevocationOrigin(*w.eaaStatusToken.Origin)
	}
	return ""
}

// Type gets the claimed type of the EAA revocation token. Port of getType().
func (w *EAARevocationTokenWrapper) Type() string {
	if w.eaaStatusToken.Type != nil {
		return *w.eaaStatusToken.Type
	}
	return ""
}

// SourceAddress gets the location URI used to access the original EAA source token. Port of
// getSourceAddress().
func (w *EAARevocationTokenWrapper) SourceAddress() string {
	if w.eaaStatusToken.SourceAddress != nil {
		return *w.eaaStatusToken.SourceAddress
	}
	return ""
}

// Subject gets the subject of the EAA revocation token. Port of getSubject().
func (w *EAARevocationTokenWrapper) Subject() string {
	if w.eaaStatusToken.Subject != nil {
		return w.eaaStatusToken.Subject.Value
	}
	return ""
}

// SubjectMatch gets whether the subject of the EAA revocation token matches the subject of
// the related EAA. Port of getSubjectMatch().
func (w *EAARevocationTokenWrapper) SubjectMatch() bool {
	return w.eaaStatusToken.Subject != nil && w.eaaStatusToken.Subject.Match != nil && *w.eaaStatusToken.Subject.Match
}

// IssuedAt gets time of the issuance of the EAA revocation token. Port of getIssuedAt().
func (w *EAARevocationTokenWrapper) IssuedAt() *time.Time {
	if w.eaaStatusToken.IssuedAt == nil {
		return nil
	}
	t := w.eaaStatusToken.IssuedAt.Time()
	return &t
}

// ExpirationTime gets time of the expiration of the EAA revocation token. Port of
// getExpirationTime().
func (w *EAARevocationTokenWrapper) ExpirationTime() *time.Time {
	if w.eaaStatusToken.ExpirationTime == nil {
		return nil
	}
	t := w.eaaStatusToken.ExpirationTime.Time()
	return &t
}

// TimeToLive gets number of seconds after which a new EAA Status token should be requested.
// Port of getTimeToLive().
func (w *EAARevocationTokenWrapper) TimeToLive() *big.Int {
	return w.eaaStatusToken.TimeToLive
}

// Binaries is the AbstractTokenProxy override. Port of getBinaries().
func (w *EAARevocationTokenWrapper) Binaries() []byte {
	if w.eaaStatusToken.Base64Encoded == nil {
		return nil
	}
	return []byte(*w.eaaStatusToken.Base64Encoded)
}

// Id is the AbstractTokenProxy override. Port of getId().
func (w *EAARevocationTokenWrapper) Id() string {
	if w.eaaStatusToken.Id != nil {
		return string(*w.eaaStatusToken.Id)
	}
	return ""
}
