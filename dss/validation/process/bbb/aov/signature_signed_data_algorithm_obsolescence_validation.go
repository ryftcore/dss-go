// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/SignatureSignedDataAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Performs cryptographic validation of algorithms used on digest computation
// for the original signed data objects.
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

// SignatureSignedDataAlgorithmObsolescenceValidation performs cryptographic
// validation of algorithms used on digest computation for the original
// signed data objects.
type SignatureSignedDataAlgorithmObsolescenceValidation struct {
	DigestAlgorithmObsolescenceValidation[*diagnostic.SignatureWrapper]
}

// NewSignatureSignedDataAlgorithmObsolescenceValidation is the default
// constructor.
func NewSignatureSignedDataAlgorithmObsolescenceValidation(i18nProvider *i18n.I18nProvider, token *diagnostic.SignatureWrapper,
	context enumerations.Context, validationDate time.Time,
	validationPolicy policy.ValidationPolicy) *SignatureSignedDataAlgorithmObsolescenceValidation {
	c := &SignatureSignedDataAlgorithmObsolescenceValidation{}
	c.InitAlgorithmObsolescenceValidation(i18nProvider, token, context, validationDate, validationPolicy, c)
	c.InitChainBase(c)
	return c
}

// BuildChain builds a chain of checks to be executed during the process. Port
// of the overridden protected ChainItem<XmlAOV> buildChain().
func (c *SignatureSignedDataAlgorithmObsolescenceValidation) BuildChain() process.ChainItem[*jaxb.XmlAOV] {
	return c.buildDigestMatchersValidationChain(c.FirstItem, c.token.DigestMatchers(), c.token.Id())
}
