// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/TokenCertificateChainAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Performs cryptographic validation of the token's certificate chain only. T
// is a diagnostic.TokenProxy implementation.
package aov

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TokenCertificateChainAlgorithmObsolescenceValidation performs cryptographic
// validation of the token's certificate chain only.
type TokenCertificateChainAlgorithmObsolescenceValidation[T diagnostic.TokenProxy] struct {
	TokenAlgorithmObsolescenceValidation[T]
}

// NewTokenCertificateChainAlgorithmObsolescenceValidation is the default
// constructor.
func NewTokenCertificateChainAlgorithmObsolescenceValidation[T diagnostic.TokenProxy](i18nProvider *i18n.I18nProvider,
	token T, context enumerations.Context, validationDate time.Time,
	validationPolicy policy.ValidationPolicy) *TokenCertificateChainAlgorithmObsolescenceValidation[T] {
	c := &TokenCertificateChainAlgorithmObsolescenceValidation[T]{}
	c.InitAlgorithmObsolescenceValidation(i18nProvider, token, context, validationDate, validationPolicy, c)
	c.InitChainBase(c)
	return c
}

// BuildChain builds a chain of checks to be executed during the process. Port
// of the overridden protected ChainItem<XmlAOV> buildChain().
func (c *TokenCertificateChainAlgorithmObsolescenceValidation[T]) BuildChain() process.ChainItem[*jaxb.XmlAOV] {
	return c.buildCertificateChainValidationChain(c.FirstItem, c.token.SigningCertificate(), c.token.CertificateChain())
}
