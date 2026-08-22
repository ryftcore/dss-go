// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/cc/checks/DigestAlgorithmAtValidationTimeCheck.java (DSS 6.5.RC1).
//
// Package placement deviation: see abstract_cryptographic_check.go.
package aov

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	vpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// DigestAlgorithmAtValidationTimeCheck checks DigestAlgorithm at validation
// time.
type DigestAlgorithmAtValidationTimeCheck struct {
	*AbstractCryptographicCheck

	// digestAlgo is the algorithm to check.
	digestAlgo enumerations.DigestAlgorithm

	// validationDate is the validation date.
	validationDate time.Time

	// cryptographicSuite is the cryptographic rules.
	cryptographicSuite policy.CryptographicSuite
}

// NewDigestAlgorithmAtValidationTimeCheck is the default constructor.
func NewDigestAlgorithmAtValidationTimeCheck(i18nProvider *i18n.I18nProvider, digestAlgo enumerations.DigestAlgorithm,
	validationDate time.Time, result *process.Result[*jaxb.XmlCC], position i18n.MessageTag,
	cryptographicSuite policy.CryptographicSuite) *DigestAlgorithmAtValidationTimeCheck {
	c := &DigestAlgorithmAtValidationTimeCheck{
		digestAlgo:         digestAlgo,
		validationDate:     validationDate,
		cryptographicSuite: cryptographicSuite,
	}
	c.AbstractCryptographicCheck = NewAbstractCryptographicCheck(i18nProvider, result, position,
		process.GetLevelRule(cryptographicSuite.AlgorithmsExpirationDateLevel()))
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of the overridden process().
func (c *DigestAlgorithmAtValidationTimeCheck) Process() bool {
	return vpolicy.IsDigestAlgorithmReliableAtTime(c.cryptographicSuite, c.digestAlgo, c.validationDate)
}

// Level returns an execution Level of the chain item. Port of the overridden
// getLevel().
func (c *DigestAlgorithmAtValidationTimeCheck) Level() enumerations.Level {
	algoExpirationDate := vpolicy.GetExpirationDateForDigestAlgorithm(c.cryptographicSuite, c.digestAlgo)
	cryptographicSuiteUpdateDate := c.cryptographicSuite.CryptographicSuiteUpdateDate()
	if algoExpirationDate != nil && cryptographicSuiteUpdateDate != nil && cryptographicSuiteUpdateDate.Before(*algoExpirationDate) {
		return c.cryptographicSuite.AlgorithmsExpirationDateAfterUpdateLevel()
	}
	return c.AbstractCryptographicCheck.Level()
}

// BuildConstraintMessage builds a constraint message. Port of the overridden
// buildConstraintMessage().
func (c *DigestAlgorithmAtValidationTimeCheck) BuildConstraintMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTag_ASCCM_AR, c.Name(c.digestAlgo))
}

// BuildErrorMessage builds an error message. Port of the overridden
// buildErrorMessage().
func (c *DigestAlgorithmAtValidationTimeCheck) BuildErrorMessage() *jaxb.XmlMessage {
	messageTag := i18n.MessageTag_ASCCM_AR_ANS_ANR_2 // other cases
	algoExpirationDate := vpolicy.GetExpirationDateForDigestAlgorithm(c.cryptographicSuite, c.digestAlgo)
	if algoExpirationDate != nil && algoExpirationDate.Before(c.validationDate) {
		messageTag = i18n.MessageTag_ASCCM_AR_ANS_ANR // expired case
	}
	return c.BuildXmlMessage(messageTag, c.Name(c.digestAlgo), c.position)
}
