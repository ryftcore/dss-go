// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/EvidenceRecordAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Performs cryptographic validation of the algorithms used on evidence
// record creation.
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

// EvidenceRecordAlgorithmObsolescenceValidation performs cryptographic
// validation of the algorithms used on evidence record creation.
type EvidenceRecordAlgorithmObsolescenceValidation struct {
	DigestAlgorithmObsolescenceValidation[*diagnostic.EvidenceRecordWrapper]
}

// NewEvidenceRecordAlgorithmObsolescenceValidation is the default
// constructor.
func NewEvidenceRecordAlgorithmObsolescenceValidation(i18nProvider *i18n.I18nProvider, token *diagnostic.EvidenceRecordWrapper,
	validationDate time.Time, validationPolicy policy.ValidationPolicy) *EvidenceRecordAlgorithmObsolescenceValidation {
	c := &EvidenceRecordAlgorithmObsolescenceValidation{}
	c.InitAlgorithmObsolescenceValidation(i18nProvider, token, enumerations.ContextEvidenceRecord, validationDate, validationPolicy, c)
	c.InitChainBase(c)
	return c
}

// BuildChain builds a chain of checks to be executed during the process. Port
// of the overridden protected ChainItem<XmlAOV> buildChain().
func (c *EvidenceRecordAlgorithmObsolescenceValidation) BuildChain() process.ChainItem[*jaxb.XmlAOV] {
	return c.buildDigestMatchersValidationChain(c.FirstItem, c.token.DigestMatchers(), c.token.Id())
}

// CryptographicSuite gets the cryptographic suite based on the currently
// validating context. Port of the overridden protected CryptographicSuite
// getCryptographicSuite().
func (c *EvidenceRecordAlgorithmObsolescenceValidation) CryptographicSuite() policy.CryptographicSuite {
	return c.validationPolicy.EvidenceRecordCryptographicConstraint()
}
