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
	return w.eaaStatusToken.CertificateChain
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
	return w.eaaStatusToken.Origin
}

// Type gets the claimed type of the EAA revocation token. Port of getType().
func (w *EAARevocationTokenWrapper) Type() string {
	return w.eaaStatusToken.Type
}

// SourceAddress gets the location URI used to access the original EAA source token. Port of
// getSourceAddress().
func (w *EAARevocationTokenWrapper) SourceAddress() string {
	return w.eaaStatusToken.SourceAddress
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
	return w.eaaStatusToken.IssuedAt
}

// ExpirationTime gets time of the expiration of the EAA revocation token. Port of
// getExpirationTime().
func (w *EAARevocationTokenWrapper) ExpirationTime() *time.Time {
	return w.eaaStatusToken.ExpirationTime
}

// TimeToLive gets number of seconds after which a new EAA Status token should be requested.
// Port of getTimeToLive().
func (w *EAARevocationTokenWrapper) TimeToLive() *big.Int {
	return w.eaaStatusToken.TimeToLive
}

// Binaries is the AbstractTokenProxy override. Port of getBinaries().
func (w *EAARevocationTokenWrapper) Binaries() []byte {
	return w.eaaStatusToken.Base64Encoded
}

// Id is the AbstractTokenProxy override. Port of getId().
func (w *EAARevocationTokenWrapper) Id() string {
	return w.eaaStatusToken.Id
}
