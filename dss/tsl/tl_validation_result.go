// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/validation/TLValidationResult.java (DSS 6.5.RC1).
package tsl

import (
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
)

// TLValidationResult is the validation task result.
//
// It implements eu.europa.esig.dss.validation.job.validation.ValidationResult (Go package
// dss/validation/job). Go satisfies interfaces structurally, so that package is deliberately NOT
// imported here: this file only has to expose the five accessors the interface declares.
type TLValidationResult struct {
	// certificateSource is the used certificate source.
	certificateSource spi.CertificateSource

	// indication is the validation Indication.
	indication enumerations.Indication

	// subIndication is the validation SubIndication.
	subIndication enumerations.SubIndication

	// signingTime is the claimed signing time. The zero time.Time stands in for Java's null.
	signingTime time.Time

	// signingCertificate is the signing certificate.
	signingCertificate *model.CertificateToken
}

// NewTLValidationResult is the default constructor. Port of TLValidationResult(Indication,
// SubIndication, Date, CertificateToken, CertificateSource).
func NewTLValidationResult(indication enumerations.Indication, subIndication enumerations.SubIndication,
	signingTime time.Time, signingCertificate *model.CertificateToken,
	certificateSource spi.CertificateSource) *TLValidationResult {
	return &TLValidationResult{
		indication:         indication,
		subIndication:      subIndication,
		signingTime:        signingTime,
		signingCertificate: signingCertificate,
		certificateSource:  certificateSource,
	}
}

// Indication gets the validation Indication. Port of getIndication().
func (r *TLValidationResult) Indication() enumerations.Indication {
	return r.indication
}

// SubIndication gets the validation SubIndication. Port of getSubIndication().
func (r *TLValidationResult) SubIndication() enumerations.SubIndication {
	return r.subIndication
}

// SigningTime gets the (claimed) signing time. Port of getSigningTime().
func (r *TLValidationResult) SigningTime() time.Time {
	return r.signingTime
}

// SigningCertificate gets the signing certificate. Port of getSigningCertificate().
func (r *TLValidationResult) SigningCertificate() *model.CertificateToken {
	return r.signingCertificate
}

// PotentialSigners gets the list of signing candidates. Port of getPotentialSigners(), which
// copies the certificate source's certificates into a new ArrayList; the Go port copies the
// slice for the same reason (the caller must not be able to mutate the source's backing array).
func (r *TLValidationResult) PotentialSigners() []*model.CertificateToken {
	certificates := r.certificateSource.Certificates()
	potentialSigners := make([]*model.CertificateToken, len(certificates))
	copy(potentialSigners, certificates)
	return potentialSigners
}
