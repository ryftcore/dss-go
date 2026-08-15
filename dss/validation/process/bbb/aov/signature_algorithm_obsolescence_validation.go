// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/SignatureAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Performs validation of the signature token (or timestamp), including
// validation of the signature value, signed attributes and certificate
// chain. T is a diagnostic.TokenProxy implementation.
package aov

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SignatureAlgorithmObsolescenceValidation performs validation of the
// signature token (or timestamp), including validation of the signature
// value, signed attributes and certificate chain.
type SignatureAlgorithmObsolescenceValidation[T diagnostic.TokenProxy] struct {
	SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation[T]
}

// NewSignatureAlgorithmObsolescenceValidation is the default constructor.
func NewSignatureAlgorithmObsolescenceValidation[T diagnostic.TokenProxy](i18nProvider *i18n.I18nProvider, token T,
	context enumerations.Context, validationDate time.Time,
	validationPolicy policy.ValidationPolicy) *SignatureAlgorithmObsolescenceValidation[T] {
	c := &SignatureAlgorithmObsolescenceValidation[T]{}
	c.InitAlgorithmObsolescenceValidation(i18nProvider, token, context, validationDate, validationPolicy, c)
	c.InitChainBase(c)
	return c
}

// BuildChain builds a chain of checks to be executed during the process. Port
// of the overridden protected ChainItem<XmlAOV> buildChain().
func (c *SignatureAlgorithmObsolescenceValidation[T]) BuildChain() process.ChainItem[*jaxb.XmlAOV] {
	item := c.SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation.BuildChain()
	item = c.buildDigestMatchersValidationChain(item, c.token.DigestMatchers(), c.token.Id())
	item = c.buildCertificateChainValidationChain(item, c.token.SigningCertificate(), c.token.CertificateChain())
	return item
}
