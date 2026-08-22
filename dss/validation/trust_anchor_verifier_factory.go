// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/TrustAnchorVerifierFactory.java
// (DSS 6.5.RC1).

package validation

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	modelpolicy "github.com/ryftcore/dss-go/dss/model/policy"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
)

// TrustAnchorVerifierFactory loads a TrustAnchorVerifier from a provided
// eu.europa.esig.dss.model.policy.ValidationPolicy. Port of
// TrustAnchorVerifierFactory.
type TrustAnchorVerifierFactory struct {
	// validationPolicy is the validation policy to load the
	// TrustAnchorVerifier from. Port of the private final validationPolicy
	// field.
	validationPolicy modelpolicy.ValidationPolicy
}

// NewTrustAnchorVerifierFactory is the default constructor. Port of
// TrustAnchorVerifierFactory(ValidationPolicy).
func NewTrustAnchorVerifierFactory(validationPolicy modelpolicy.ValidationPolicy) *TrustAnchorVerifierFactory {
	return &TrustAnchorVerifierFactory{validationPolicy: validationPolicy}
}

// Create creates the TrustAnchorVerifier. Port of create().
func (f *TrustAnchorVerifierFactory) Create() *spivalidation.TrustAnchorVerifier {
	trustAnchorVerifier := spivalidation.NewEmptyTrustAnchorVerifier()
	f.instantiateAcceptUntrustedCertificateChains(trustAnchorVerifier, f.validationPolicy)
	f.instantiateUseSunsetDate(trustAnchorVerifier, f.validationPolicy)
	return trustAnchorVerifier
}

// instantiateAcceptUntrustedCertificateChains is the port of the private
// instantiateAcceptUntrustedCertificateChains(TrustAnchorVerifier, ValidationPolicy).
func (f *TrustAnchorVerifierFactory) instantiateAcceptUntrustedCertificateChains(
	trustAnchorVerifier *spivalidation.TrustAnchorVerifier, validationPolicy modelpolicy.ValidationPolicy) {
	acceptUntrustedCertificateChains := f.acceptUntrustedCertificateChains(validationPolicy, enumerations.ContextRevocation)
	trustAnchorVerifier.SetAcceptRevocationUntrustedCertificateChains(acceptUntrustedCertificateChains)

	acceptUntrustedCertificateChains = f.acceptUntrustedCertificateChains(validationPolicy, enumerations.ContextTimestamp)
	trustAnchorVerifier.SetAcceptTimestampUntrustedCertificateChains(acceptUntrustedCertificateChains)
}

// acceptUntrustedCertificateChains is the port of the private
// getAcceptUntrustedCertificateChains(ValidationPolicy, Context).
func (f *TrustAnchorVerifierFactory) acceptUntrustedCertificateChains(validationPolicy modelpolicy.ValidationPolicy,
	context enumerations.Context) bool {
	constraint := validationPolicy.ProspectiveCertificateChainConstraint(context)
	return constraint == nil || enumerations.LevelFail != constraint.Level()
}

// instantiateUseSunsetDate is the port of the private
// instantiateUseSunsetDate(TrustAnchorVerifier, ValidationPolicy).
func (f *TrustAnchorVerifierFactory) instantiateUseSunsetDate(
	trustAnchorVerifier *spivalidation.TrustAnchorVerifier, validationPolicy modelpolicy.ValidationPolicy) {
	useSunsetDate := f.useSunsetDate(validationPolicy, enumerations.ContextSignature)
	useSunsetDate = useSunsetDate || f.useSunsetDate(validationPolicy, enumerations.ContextCounterSignature)
	useSunsetDate = useSunsetDate || f.useSunsetDate(validationPolicy, enumerations.ContextKeyBindingSignature)
	useSunsetDate = useSunsetDate || f.useSunsetDate(validationPolicy, enumerations.ContextTimestamp)
	useSunsetDate = useSunsetDate || f.useSunsetDate(validationPolicy, enumerations.ContextRevocation)
	trustAnchorVerifier.SetUseSunsetDate(useSunsetDate)
}

// useSunsetDate is the port of the private getUseSunsetDate(ValidationPolicy,
// Context).
func (f *TrustAnchorVerifierFactory) useSunsetDate(validationPolicy modelpolicy.ValidationPolicy,
	context enumerations.Context) bool {
	constraint := validationPolicy.CertificateSunsetDateConstraint(context, enumerations.SubContextSigningCert)
	if constraint != nil && enumerations.LevelFail == constraint.Level() {
		return true
	}
	constraint = validationPolicy.CertificateSunsetDateConstraint(context, enumerations.SubContextCACertificate)
	if constraint != nil && enumerations.LevelFail == constraint.Level() {
		return true
	}
	return false
}
