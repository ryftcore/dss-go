// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/cc/checks/SignatureAlgorithmAtValidationTimeCheck.java (DSS 6.5.RC1).
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

// SignatureAlgorithmAtValidationTimeCheck checks SignatureAlgorithm at
// validation time.
type SignatureAlgorithmAtValidationTimeCheck struct {
	*AbstractCryptographicCheck

	// signatureAlgorithm is the algorithm to check.
	signatureAlgorithm enumerations.SignatureAlgorithm

	// keyLength is the used key size.
	keyLength string

	// validationDate is the validation time.
	validationDate time.Time

	// cryptographicSuite is the cryptographic rules.
	cryptographicSuite policy.CryptographicSuite
}

// NewSignatureAlgorithmAtValidationTimeCheck is the default constructor.
func NewSignatureAlgorithmAtValidationTimeCheck(i18nProvider *i18n.I18nProvider, signatureAlgorithm enumerations.SignatureAlgorithm,
	keyLength string, validationDate time.Time, result *process.Result[*jaxb.XmlCC], position i18n.MessageTag,
	cryptographicSuite policy.CryptographicSuite) *SignatureAlgorithmAtValidationTimeCheck {
	c := &SignatureAlgorithmAtValidationTimeCheck{
		signatureAlgorithm: signatureAlgorithm,
		keyLength:          keyLength,
		validationDate:     validationDate,
		cryptographicSuite: cryptographicSuite,
	}
	c.AbstractCryptographicCheck = NewAbstractCryptographicCheck(i18nProvider, result, position,
		process.GetLevelRule(cryptographicSuite.AlgorithmsExpirationDateLevel()))
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of the overridden process().
func (c *SignatureAlgorithmAtValidationTimeCheck) Process() bool {
	return vpolicy.IsSignatureAlgorithmReliableAtTime(c.cryptographicSuite, c.signatureAlgorithm, c.keyLength, c.validationDate)
}

// Level returns an execution Level of the chain item. Port of the overridden
// getLevel().
func (c *SignatureAlgorithmAtValidationTimeCheck) Level() enumerations.Level {
	algoExpirationDate := vpolicy.GetExpirationDateForSignatureAlgorithm(c.cryptographicSuite, c.signatureAlgorithm, c.keyLength)
	cryptographicSuiteUpdateDate := c.cryptographicSuite.CryptographicSuiteUpdateDate()
	if algoExpirationDate != nil && cryptographicSuiteUpdateDate != nil && cryptographicSuiteUpdateDate.Before(*algoExpirationDate) {
		return c.cryptographicSuite.AlgorithmsExpirationDateAfterUpdateLevel()
	}
	return c.AbstractCryptographicCheck.Level()
}

// BuildConstraintMessage builds a constraint message. Port of the overridden
// buildConstraintMessage().
func (c *SignatureAlgorithmAtValidationTimeCheck) BuildConstraintMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTagASCCMAR, c.SignatureAlgorithmName(c.signatureAlgorithm))
}

// BuildErrorMessage builds an error message. Port of the overridden
// buildErrorMessage().
func (c *SignatureAlgorithmAtValidationTimeCheck) BuildErrorMessage() *jaxb.XmlMessage {
	messageTag := i18n.MessageTagASCCMARANSAKSNR2 // other cases
	algoExpirationDate := vpolicy.GetExpirationDateForSignatureAlgorithm(c.cryptographicSuite, c.signatureAlgorithm, c.keyLength)
	if algoExpirationDate != nil && algoExpirationDate.Before(c.validationDate) {
		messageTag = i18n.MessageTagASCCMARANSAKSNR // expired case
	}
	return c.BuildXmlMessage(messageTag, c.SignatureAlgorithmName(c.signatureAlgorithm), c.keyLength, c.position)
}
