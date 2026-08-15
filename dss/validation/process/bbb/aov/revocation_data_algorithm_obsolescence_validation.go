// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/RevocationDataAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Performs cryptographic validation of the revocation data and its
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

// RevocationDataAlgorithmObsolescenceValidation performs cryptographic
// validation of the revocation data and its certificate chain.
type RevocationDataAlgorithmObsolescenceValidation struct {
	TokenAlgorithmObsolescenceValidation[*diagnostic.RevocationWrapper]
}

// NewRevocationDataAlgorithmObsolescenceValidation is the default
// constructor.
func NewRevocationDataAlgorithmObsolescenceValidation(i18nProvider *i18n.I18nProvider, token *diagnostic.RevocationWrapper,
	validationDate time.Time, validationPolicy policy.ValidationPolicy) *RevocationDataAlgorithmObsolescenceValidation {
	c := &RevocationDataAlgorithmObsolescenceValidation{}
	c.InitAlgorithmObsolescenceValidation(i18nProvider, token, enumerations.Context_REVOCATION, validationDate, validationPolicy, c)
	c.InitChainBase(c)
	return c
}

// BuildChain builds a chain of checks to be executed during the process. Port
// of the overridden protected ChainItem<XmlAOV> buildChain().
//
// TODO : add processing of CertID algorithm (carried over from the Java
// source's own TODO comment).
func (c *RevocationDataAlgorithmObsolescenceValidation) BuildChain() process.ChainItem[*jaxb.XmlAOV] {
	item := c.TokenAlgorithmObsolescenceValidation.BuildChain()
	item = c.buildCertificateChainValidationChain(item, c.token.SigningCertificate(), c.token.CertificateChain())
	return item
}
