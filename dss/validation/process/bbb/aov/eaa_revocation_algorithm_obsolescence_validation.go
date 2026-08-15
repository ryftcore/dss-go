// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/EAARevocationAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Performs cryptographic validation of the EAA revocation token and its
// certificate chain.
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

// EAARevocationAlgorithmObsolescenceValidation performs cryptographic
// validation of the EAA revocation token and its certificate chain.
type EAARevocationAlgorithmObsolescenceValidation struct {
	TokenAlgorithmObsolescenceValidation[*diagnostic.EAARevocationTokenWrapper]
}

// NewEAARevocationAlgorithmObsolescenceValidation is the common constructor.
func NewEAARevocationAlgorithmObsolescenceValidation(i18nProvider *i18n.I18nProvider, token *diagnostic.EAARevocationTokenWrapper,
	validationDate time.Time, validationPolicy policy.ValidationPolicy) *EAARevocationAlgorithmObsolescenceValidation {
	c := &EAARevocationAlgorithmObsolescenceValidation{}
	c.InitAlgorithmObsolescenceValidation(i18nProvider, token, enumerations.Context_EAA_REVOCATION, validationDate, validationPolicy, c)
	c.InitChainBase(c)
	return c
}

// BuildChain builds a chain of checks to be executed during the process. Port
// of the overridden protected ChainItem<XmlAOV> buildChain().
func (c *EAARevocationAlgorithmObsolescenceValidation) BuildChain() process.ChainItem[*jaxb.XmlAOV] {
	item := c.TokenAlgorithmObsolescenceValidation.BuildChain()
	item = c.buildCertificateChainValidationChain(item, c.token.SigningCertificate(), c.token.CertificateChain())
	return item
}
