// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/EAAAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Performs algorithm obsolescence validation for cryptographic algorithms
// used within an EAA.
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

// EAAAlgorithmObsolescenceValidation performs algorithm obsolescence
// validation for cryptographic algorithms used within an EAA.
type EAAAlgorithmObsolescenceValidation struct {
	DigestAlgorithmObsolescenceValidation[*diagnostic.EAAWrapper]
}

// NewEAAAlgorithmObsolescenceValidation is the default constructor.
func NewEAAAlgorithmObsolescenceValidation(i18nProvider *i18n.Provider, token *diagnostic.EAAWrapper,
	validationDate time.Time, validationPolicy policy.ValidationPolicy) *EAAAlgorithmObsolescenceValidation {
	c := &EAAAlgorithmObsolescenceValidation{}
	c.InitAlgorithmObsolescenceValidation(i18nProvider, token, enumerations.ContextEAA, validationDate, validationPolicy, c)
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
