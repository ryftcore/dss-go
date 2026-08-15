// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/EAAAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Performs algorithm obsolescence validation for cryptographic algorithms
// used within an EAA.
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

// EAAAlgorithmObsolescenceValidation performs algorithm obsolescence
// validation for cryptographic algorithms used within an EAA.
type EAAAlgorithmObsolescenceValidation struct {
	DigestAlgorithmObsolescenceValidation[*diagnostic.EAAWrapper]
}

// NewEAAAlgorithmObsolescenceValidation is the default constructor.
func NewEAAAlgorithmObsolescenceValidation(i18nProvider *i18n.I18nProvider, token *diagnostic.EAAWrapper,
	validationDate time.Time, validationPolicy policy.ValidationPolicy) *EAAAlgorithmObsolescenceValidation {
	c := &EAAAlgorithmObsolescenceValidation{}
	c.InitAlgorithmObsolescenceValidation(i18nProvider, token, enumerations.Context_EAA, validationDate, validationPolicy, c)
	c.InitChainBase(c)
	return c
}

// BuildChain builds a chain of checks to be executed during the process. Port
// of the overridden protected ChainItem<XmlAOV> buildChain().
func (c *EAAAlgorithmObsolescenceValidation) BuildChain() process.ChainItem[*jaxb.XmlAOV] {
	return c.buildDigestMatchersValidationChain(c.FirstItem, c.token.DigestMatchers(), c.token.Id())
}

// CryptographicSuite gets the cryptographic suite based on the currently
// validating context. Port of the overridden protected CryptographicSuite
// getCryptographicSuite().
func (c *EAAAlgorithmObsolescenceValidation) CryptographicSuite() policy.CryptographicSuite {
	return c.validationPolicy.EAACryptographicConstraint()
}
