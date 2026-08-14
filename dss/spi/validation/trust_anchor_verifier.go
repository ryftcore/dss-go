// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/TrustAnchorVerifier.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY: two already-landed sibling files (spi/validation/revocation_data_verifier.go
// and spi/validation/signature_validation_context.go) forward-declare their expectation of this
// type's shape as:
//
//	NewDefaultTrustAnchorVerifier() *TrustAnchorVerifier
//	TrustedCertificateSource() / SetTrustedCertificateSource(*spi.ListCertificateSource)
//	IsTrustedAtTime(cert *model.CertificateToken, controlTime time.Time, context enumerations.Context) bool
//
// This file matches that exact 3-argument IsTrustedAtTime signature as the canonical method
// (signature_validation_context.go:1801 and revocation_data_verifier.go:501 already call it
// this way). Java additionally declares a 2-argument overload isTrustedAtTime(certificateToken,
// controlTime) that simply forwards with a null Context; since Go has no overloading and the
// 3-argument name is already load-bearing, that convenience overload is ported under the
// distinct name IsTrustedAtTimeAnyContext (same treatment applied to the analogous
// isTrustedCertificateChain(certChain, controlTime) 2-argument overload, as
// IsTrustedCertificateChainAnyContext).
//
// Java's setTrustedCertificateSource(CertificateSource) takes the interface type; the two call
// sites above pass a *spi.ListCertificateSource (which satisfies spi.CertificateSource
// structurally, per the "compile-time assertion" precedent noted throughout this package), so
// the parameter here is kept as the interface spi.CertificateSource for Java fidelity - this
// remains source-compatible with both existing call sites.
package validation

import (
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
)

// TrustAnchorVerifier is used to verify whether a given certificate token is trusted at the
// control time.
type TrustAnchorVerifier struct {
	// trustedCertificateSource provides a source to trust anchors.
	trustedCertificateSource spi.CertificateSource

	// acceptTimestampUntrustedCertificateChains indicates whether timestamp's untrusted
	// certificate chains shall be accepted.
	acceptTimestampUntrustedCertificateChains bool

	// acceptRevocationUntrustedCertificateChains indicates whether revocation data's untrusted
	// certificate chains shall be accepted.
	acceptRevocationUntrustedCertificateChains bool

	// useSunsetDate indicates whether the sunset date should be used for trust anchor
	// determinations.
	useSunsetDate bool
}

// newTrustAnchorVerifier is the port of the protected default constructor.
func newTrustAnchorVerifier() *TrustAnchorVerifier {
	return &TrustAnchorVerifier{}
}

// NewEmptyTrustAnchorVerifier creates an empty instance of TrustAnchorVerifier. All constraints
// should be configured manually. Port of createEmptyTrustAnchorVerifier().
func NewEmptyTrustAnchorVerifier() *TrustAnchorVerifier {
	return newTrustAnchorVerifier()
}

// NewDefaultTrustAnchorVerifier creates a default instance of TrustAnchorVerifier, with
// pre-configured constraints. Port of createDefaultTrustAnchorVerifier().
func NewDefaultTrustAnchorVerifier() *TrustAnchorVerifier {
	t := newTrustAnchorVerifier()
	t.SetUseSunsetDate(true)
	return t
}

// IsAcceptTimestampUntrustedCertificateChains gets whether untrusted certificate chains of
// timestamps should be accepted. Port of isAcceptTimestampUntrustedCertificateChains().
func (t *TrustAnchorVerifier) IsAcceptTimestampUntrustedCertificateChains() bool {
	return t.acceptTimestampUntrustedCertificateChains
}

// SetAcceptTimestampUntrustedCertificateChains sets whether untrusted certificate chains of
// timestamps should be accepted. Default: FALSE (only timestamps created with trusted CAs are
// considered as valid, untrusted timestamps are ignored). Port of
// setAcceptTimestampUntrustedCertificateChains(boolean).
func (t *TrustAnchorVerifier) SetAcceptTimestampUntrustedCertificateChains(acceptTimestampUntrustedCertificateChains bool) {
	t.acceptTimestampUntrustedCertificateChains = acceptTimestampUntrustedCertificateChains
}

// IsAcceptRevocationUntrustedCertificateChains gets whether untrusted certificate chains of
// revocation data should be accepted. Port of isAcceptRevocationUntrustedCertificateChains().
func (t *TrustAnchorVerifier) IsAcceptRevocationUntrustedCertificateChains() bool {
	return t.acceptRevocationUntrustedCertificateChains
}

// SetAcceptRevocationUntrustedCertificateChains sets whether untrusted certificate chains of
// revocation data should be accepted. Default: FALSE (only revocation data created with
// trusted CAs are considered as valid, untrusted revocation data is ignored). Port of
// setAcceptRevocationUntrustedCertificateChains(boolean).
func (t *TrustAnchorVerifier) SetAcceptRevocationUntrustedCertificateChains(acceptRevocationUntrustedCertificateChains bool) {
	t.acceptRevocationUntrustedCertificateChains = acceptRevocationUntrustedCertificateChains
}

// TrustedCertificateSource gets trusted certificate source, when present. Port of
// getTrustedCertificateSource().
func (t *TrustAnchorVerifier) TrustedCertificateSource() spi.CertificateSource {
	return t.trustedCertificateSource
}

// SetTrustedCertificateSource sets a trusted certificate source in order to provide information
// about the available trust anchors.
//
// NOTE: This method is used internally during a SignatureValidationContext initialization, in
// order to provide the same trusted source as the one used within a CertificateVerifier. Port
// of setTrustedCertificateSource(CertificateSource).
func (t *TrustAnchorVerifier) SetTrustedCertificateSource(trustedCertificateSource spi.CertificateSource) {
	t.trustedCertificateSource = trustedCertificateSource
}

// IsUseSunsetDate defines whether sunset date shall be considered during trust anchor
// validation. Port of isUseSunsetDate().
func (t *TrustAnchorVerifier) IsUseSunsetDate() bool {
	return t.useSunsetDate
}

// SetUseSunsetDate sets whether a trust anchor's sunset date shall be taken into account when
// checking a trust anchor. Default : TRUE (sunset date is used for a trust anchor
// determination, when applicable). Port of setUseSunsetDate(boolean).
func (t *TrustAnchorVerifier) SetUseSunsetDate(useSunsetDate bool) {
	t.useSunsetDate = useSunsetDate
}

// IsTrustedAtTimeAnyContext verifies whether the certificateToken is trusted at controlTime,
// without considering a specific validation Context. Port of the 2-argument overload
// isTrustedAtTime(CertificateToken, Date).
func (t *TrustAnchorVerifier) IsTrustedAtTimeAnyContext(certificateToken *model.CertificateToken, controlTime time.Time) bool {
	return t.IsTrustedAtTime(certificateToken, controlTime, "")
}

// IsTrustedAtTime verifies whether the certificateToken is trusted at controlTime. Port of the
// 3-argument overload isTrustedAtTime(CertificateToken, Date, Context).
func (t *TrustAnchorVerifier) IsTrustedAtTime(certificateToken *model.CertificateToken, controlTime time.Time, context enumerations.Context) bool {
	if t.isAcceptUntrustedCertificateChains(context) {
		return true
	} else if t.trustedCertificateSource == nil {
		return false
	} else if t.useSunsetDate && !controlTime.IsZero() {
		return t.trustedCertificateSource.IsTrustedAtTime(certificateToken, controlTime)
	}
	return t.trustedCertificateSource.IsTrusted(certificateToken)
}

// IsTrustedCertificateChainAnyContext verifies whether the certificate chain contains a trust
// anchor, without considering a specific validation Context. Port of the 2-argument overload
// isTrustedCertificateChain(List, Date).
func (t *TrustAnchorVerifier) IsTrustedCertificateChainAnyContext(certChain []*model.CertificateToken, controlTime time.Time) bool {
	return t.IsTrustedCertificateChain(certChain, controlTime, "")
}

// IsTrustedCertificateChain verifies whether the certificate chain contains a trust anchor.
// Port of the 3-argument overload isTrustedCertificateChain(List, Date, Context).
func (t *TrustAnchorVerifier) IsTrustedCertificateChain(certChain []*model.CertificateToken, controlTime time.Time, context enumerations.Context) bool {
	if t.isAcceptUntrustedCertificateChains(context) {
		return true
	}
	if utils.IsCollectionNotEmpty(certChain) {
		for _, token := range certChain {
			if t.IsTrustedAtTime(token, controlTime, context) {
				return true
			}
		}
	}
	return false
}

// isAcceptUntrustedCertificateChains is the port of the private isAcceptUntrustedCertificateChains(Context).
func (t *TrustAnchorVerifier) isAcceptUntrustedCertificateChains(context enumerations.Context) bool {
	switch context {
	case enumerations.Context_TIMESTAMP:
		return t.acceptTimestampUntrustedCertificateChains
	case enumerations.Context_REVOCATION:
		return t.acceptRevocationUntrustedCertificateChains
	default:
		return false // continue in other cases
	}
}
