// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/cc/checks/SignatureAlgorithmReliableCheck.java (DSS 6.5.RC1).
//
// Package placement deviation: see abstract_cryptographic_check.go.
package aov

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	vpolicy "github.com/utain/esig/dss/validation/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SignatureAlgorithmReliableCheck checks if SignatureAlgorithm is acceptable.
type SignatureAlgorithmReliableCheck struct {
	*AbstractCryptographicCheck

	// signatureAlgorithm is the algorithm to check.
	signatureAlgorithm enumerations.SignatureAlgorithm

	// cryptographicSuite is the cryptographic rules.
	cryptographicSuite policy.CryptographicSuite
}

// NewSignatureAlgorithmReliableCheck is the default constructor.
func NewSignatureAlgorithmReliableCheck(i18nProvider *i18n.I18nProvider, signatureAlgorithm enumerations.SignatureAlgorithm,
	result *process.Result[*jaxb.XmlCC], position i18n.MessageTag,
	cryptographicSuite policy.CryptographicSuite) *SignatureAlgorithmReliableCheck {
	c := &SignatureAlgorithmReliableCheck{
		signatureAlgorithm: signatureAlgorithm,
		cryptographicSuite: cryptographicSuite,
	}
	c.AbstractCryptographicCheck = NewAbstractCryptographicCheck(i18nProvider, result, position,
		process.GetLevelRule(cryptographicSuite.AcceptableSignatureAlgorithmsLevel()))
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of the overridden process().
func (c *SignatureAlgorithmReliableCheck) Process() bool {
	return vpolicy.IsSignatureAlgorithmReliable(c.cryptographicSuite, c.signatureAlgorithm)
}

// BuildConstraintMessage builds a constraint message. Port of the overridden
// buildConstraintMessage().
func (c *SignatureAlgorithmReliableCheck) BuildConstraintMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTag_ASCCM_CAA, c.SignatureAlgorithmName(c.signatureAlgorithm))
}

// BuildErrorMessage builds an error message. Port of the overridden
// buildErrorMessage().
func (c *SignatureAlgorithmReliableCheck) BuildErrorMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTag_ASCCM_CAA_ANS, c.SignatureAlgorithmName(c.signatureAlgorithm), c.position)
}
