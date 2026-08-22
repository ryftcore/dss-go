// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/cc/checks/PublicKeySizeAcceptableCheck.java (DSS 6.5.RC1).
//
// Package placement deviation: see abstract_cryptographic_check.go.
package aov

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	vpolicy "github.com/ryftcore/dss-go/dss/validation/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// PublicKeySizeAcceptableCheck checks if public key size is acceptable.
type PublicKeySizeAcceptableCheck struct {
	*AbstractCryptographicCheck

	// signatureAlgorithm is the algorithm to check.
	signatureAlgorithm enumerations.SignatureAlgorithm

	// keyLength is the used public key length.
	keyLength string

	// cryptographicSuite is the cryptographic rules.
	cryptographicSuite policy.CryptographicSuite
}

// NewPublicKeySizeAcceptableCheck is the default constructor.
func NewPublicKeySizeAcceptableCheck(i18nProvider *i18n.I18nProvider, signatureAlgorithm enumerations.SignatureAlgorithm,
	keyLength string, result *process.Result[*jaxb.XmlCC], position i18n.MessageTag,
	cryptographicSuite policy.CryptographicSuite) *PublicKeySizeAcceptableCheck {
	c := &PublicKeySizeAcceptableCheck{
		signatureAlgorithm: signatureAlgorithm,
		keyLength:          keyLength,
		cryptographicSuite: cryptographicSuite,
	}
	c.AbstractCryptographicCheck = NewAbstractCryptographicCheck(i18nProvider, result, position,
		process.GetLevelRule(cryptographicSuite.AcceptableSignatureAlgorithmsMiniKeySizeLevel()))
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of the overridden process().
func (c *PublicKeySizeAcceptableCheck) Process() bool {
	return vpolicy.IsSignatureAlgorithmWithKeySizeReliable(c.cryptographicSuite, c.signatureAlgorithm, c.keyLength)
}

// BuildConstraintMessage builds a constraint message. Port of the overridden
// buildConstraintMessage().
func (c *PublicKeySizeAcceptableCheck) BuildConstraintMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTag_ASCCM_APKSA, c.SignatureAlgorithmName(c.signatureAlgorithm), c.keyLength)
}

// BuildErrorMessage builds an error message. Port of the overridden
// buildErrorMessage().
func (c *PublicKeySizeAcceptableCheck) BuildErrorMessage() *jaxb.XmlMessage {
	messageTag := i18n.MessageTag_ASCCM_APKSA_ANS
	if vpolicy.IsSignatureAlgorithmKeyLengthBigEnough(c.cryptographicSuite, c.signatureAlgorithm, c.keyLength) {
		messageTag = i18n.MessageTag_ASCCM_APKSA_ANS_2
	}
	return c.BuildXmlMessage(messageTag, c.SignatureAlgorithmName(c.signatureAlgorithm), c.keyLength, c.position)
}
