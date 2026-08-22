// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/TimestampTokenVerifier.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY (flagged per S2B_BRIEF.md): TrustAnchorVerifier (Java
// spi.validation.TrustAnchorVerifier) is assigned to sibling chunk VAL-C (s2b_VAL-C.txt),
// which lands it in this same package. It is referenced here by name only, using the shape
// already documented by certificate_verifier.go's and revocation_data_verifier.go's header
// comments: NewDefaultTrustAnchorVerifier() *TrustAnchorVerifier and
// IsTrustedCertificateChain(certChain []*model.CertificateToken, controlTime time.Time,
// context enumerations.Context) bool.
//
// slf4j logging is dropped per the phase 2a handoff fact ("slf4j dropped unless
// load-bearing").
package validation

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// TimestampTokenVerifier is used to verify applicability of a timestamp token within the
// signature validation process.
type TimestampTokenVerifier struct {
	// trustAnchorVerifier verifies whether a given certificate token is a trust anchor at the
	// control time.
	trustAnchorVerifier *TrustAnchorVerifier

	// revocationDataVerifier verifies validity of the certificate's revocation data for
	// timestamps's certificate chain.
	revocationDataVerifier *RevocationDataVerifier
}

// NewEmptyTimestampTokenVerifier creates an empty instance of TimestampTokenVerifier.
// All constraints should be configured manually. Port of
// createEmptyTimestampTokenVerifier().
func NewEmptyTimestampTokenVerifier() *TimestampTokenVerifier {
	return &TimestampTokenVerifier{}
}

// NewDefaultTimestampTokenVerifier creates a default instance of TimestampTokenVerifier, with
// pre-configured constraints. Port of createDefaultTimestampTokenVerifier().
func NewDefaultTimestampTokenVerifier() *TimestampTokenVerifier {
	// no configuration available
	return &TimestampTokenVerifier{}
}

// TrustAnchorVerifier gets a trust anchor verifier. This method is used internally within
// eu.europa.esig.dss.validation.SignatureValidationContext to identify whether the
// configuration is already present and a trustAnchorVerifier should be set. Port of
// getTrustAnchorVerifier().
func (v *TimestampTokenVerifier) TrustAnchorVerifier() *TrustAnchorVerifier {
	return v.trustAnchorVerifier
}

// SetTrustAnchorVerifier sets whether a certificate token can be considered as a trust anchor
// at the given control time.
//
// NOTE: This method is used internally during a
// eu.europa.esig.dss.validation.SignatureValidationContext initialization, when not defined
// explicitly, in order to provide the same configuration as the one used within a
// eu.europa.esig.dss.validation.CertificateVerifier.
//
// Port of setTrustAnchorVerifier(...).
func (v *TimestampTokenVerifier) SetTrustAnchorVerifier(trustAnchorVerifier *TrustAnchorVerifier) {
	v.trustAnchorVerifier = trustAnchorVerifier
}

// RevocationDataVerifier gets a revocation data verifier. This method is used internally
// within eu.europa.esig.dss.validation.SignatureValidationContext to identify whether the
// configuration is already present and a trustAnchorVerifier should be set. Port of
// getRevocationDataVerifier().
func (v *TimestampTokenVerifier) RevocationDataVerifier() *RevocationDataVerifier {
	if v.revocationDataVerifier != nil && v.revocationDataVerifier.TrustAnchorVerifier() == nil {
		v.revocationDataVerifier.SetTrustAnchorVerifier(v.TrustAnchorVerifier())
	}
	return v.revocationDataVerifier
}

// SetRevocationDataVerifier sets a revocation data verifier for validation of timestamp's
// certificate chain revocation data validity.
//
// NOTE: This method is used internally during a
// eu.europa.esig.dss.validation.SignatureValidationContext initialization, when not defined
// explicitly, in order to provide the same configuration as the one used within a
// eu.europa.esig.dss.validation.CertificateVerifier.
//
// Port of setRevocationDataVerifier(...).
func (v *TimestampTokenVerifier) SetRevocationDataVerifier(revocationDataVerifier *RevocationDataVerifier) {
	v.revocationDataVerifier = revocationDataVerifier
}

// IsAcceptable verifies whether the given timestampToken is valid and acceptable at the
// current time, and its POE can be extracted to the validation process.
//
// NOTE: The method does not accept certificate chain, thus validity of the timestamp's
// certificate chain is not verified. To successfully execute this method, the
// acceptOnlyTrustedCertificateChains constraint shall be set to false. For validation with a
// certificate chain, use IsAcceptableWithChain.
//
// Port of isAcceptable(TimestampToken).
func (v *TimestampTokenVerifier) IsAcceptable(timestampToken *TimestampToken) bool {
	return v.IsAcceptableAt(timestampToken, time.Now())
}

// IsAcceptableAt verifies whether the given timestampToken is valid and acceptable at the
// given control time, and its POE can be extracted to the validation process.
//
// NOTE: The method does not accept certificate chain, thus validity of the timestamp's
// certificate chain is not verified. To successfully execute this method, the
// acceptOnlyTrustedCertificateChains constraint shall be set to false. For validation with a
// certificate chain, use IsAcceptableWithChain.
//
// Port of isAcceptable(TimestampToken, Date).
func (v *TimestampTokenVerifier) IsAcceptableAt(timestampToken *TimestampToken, controlTime time.Time) bool {
	return v.IsAcceptableWithChainAt(timestampToken, nil, controlTime)
}

// IsAcceptableWithChain verifies whether the given timestampToken is valid and acceptable at
// the current time, and its POE can be extracted to the validation process. Port of
// isAcceptable(TimestampToken, List).
func (v *TimestampTokenVerifier) IsAcceptableWithChain(timestampToken *TimestampToken, certificateChain []*model.CertificateToken) bool {
	return v.IsAcceptableWithChainAt(timestampToken, certificateChain, time.Now())
}

// IsAcceptableWithChainAt verifies whether the given timestampToken is valid and acceptable at
// the given control time, and its POE can be extracted to the validation process. Port of
// isAcceptable(TimestampToken, List, Date).
func (v *TimestampTokenVerifier) IsAcceptableWithChainAt(timestampToken *TimestampToken, certificateChain []*model.CertificateToken, controlTime time.Time) bool {
	return v.isTrustedTimestampToken(timestampToken, certificateChain, controlTime) &&
		v.isCryptographicallyValid(timestampToken) &&
		v.isCertificateChainValid(certificateChain, controlTime)
}

// isTrustedTimestampToken verifies whether the timestampToken is trusted to continue the
// process at the control time. The method expects the certificate chain of the timestamp to
// reach a trustedCertificateSource or to have acceptOnlyTrustedCertificateChains constraint to
// accept untrusted certificate chains as well. Port of isTrustedTimestampToken(...); protected
// in Java.
func (v *TimestampTokenVerifier) isTrustedTimestampToken(timestampToken *TimestampToken, certificateChain []*model.CertificateToken, controlTime time.Time) bool {
	return v.containsTrustAnchor(certificateChain, controlTime)
}

// containsTrustAnchor verifies whether the certificate chain is trusted at the given time.
// Port of containsTrustAnchor(...); protected in Java.
func (v *TimestampTokenVerifier) containsTrustAnchor(certChain []*model.CertificateToken, controlTime time.Time) bool {
	currentTrustAnchorVerifier := v.TrustAnchorVerifier()
	if currentTrustAnchorVerifier == nil {
		return false
	}
	return currentTrustAnchorVerifier.IsTrustedCertificateChain(certChain, controlTime, enumerations.ContextTimestamp)
}

// isCryptographicallyValid verifies whether the timestampToken is cryptographically valid
// (signature and message imprint match). Port of isCryptographicallyValid(...); protected in
// Java.
func (v *TimestampTokenVerifier) isCryptographicallyValid(timestampToken *TimestampToken) bool {
	if !timestampToken.IsMessageImprintDataIntact() {
		return false
	}
	if !timestampToken.IsSignatureIntact() {
		return false
	}
	return true
}

// isCertificateChainValid verifies certificate chain and presence of a valid revocation data
// for certificates. Port of isCertificateChainValid(...); protected in Java.
func (v *TimestampTokenVerifier) isCertificateChainValid(certificateChain []*model.CertificateToken, controlTime time.Time) bool {
	currentRevocationDataVerifier := v.RevocationDataVerifier()
	if v.revocationDataVerifier == nil {
		return true
	}
	return currentRevocationDataVerifier.IsCertificateChainValid(certificateChain, controlTime, enumerations.ContextTimestamp)
}
