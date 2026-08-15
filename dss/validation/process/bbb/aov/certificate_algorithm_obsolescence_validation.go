// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/CertificateAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Performs cryptographic validation of the certificate.
package aov

import (
	"time"

	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// CertificateAlgorithmObsolescenceValidation performs cryptographic
// validation of the certificate.
type CertificateAlgorithmObsolescenceValidation struct {
	TokenAlgorithmObsolescenceValidation[*diagnostic.CertificateWrapper]

	// subContext is the SubContext of the validating certificate.
	subContext enumerations.SubContext
}

// NewCertificateAlgorithmObsolescenceValidation is the default constructor.
func NewCertificateAlgorithmObsolescenceValidation(i18nProvider *i18n.I18nProvider, token *diagnostic.CertificateWrapper,
	context enumerations.Context, subContext enumerations.SubContext, validationDate time.Time,
	validationPolicy policy.ValidationPolicy) *CertificateAlgorithmObsolescenceValidation {
	c := &CertificateAlgorithmObsolescenceValidation{subContext: subContext}
	c.InitAlgorithmObsolescenceValidation(i18nProvider, token, context, validationDate, validationPolicy, c)
	c.InitChainBase(c)
	return c
}

// Position gets position of the currently verifying token, based on the
// context and subContext. Port of the overridden protected MessageTag
// getPosition().
func (c *CertificateAlgorithmObsolescenceValidation) Position() i18n.MessageTag {
	position, err := process.GetSubContextPosition(c.context, c.subContext)
	if err != nil {
		panic(err)
	}
	return position
}

// CryptographicSuite gets the cryptographic suite based on the currently
// validating context and subContext. Port of the overridden protected
// CryptographicSuite getCryptographicSuite().
func (c *CertificateAlgorithmObsolescenceValidation) CryptographicSuite() policy.CryptographicSuite {
	return c.validationPolicy.CertificateCryptographicConstraint(c.context, c.subContext)
}
