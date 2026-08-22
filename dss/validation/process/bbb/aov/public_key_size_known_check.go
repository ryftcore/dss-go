// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/cc/checks/PublicKeySizeKnownCheck.java (DSS 6.5.RC1).
//
// Package placement deviation: see abstract_cryptographic_check.go.
package aov

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// PublicKeySizeKnownCheck checks if EncryptionAlgorithm's public key size is
// known.
type PublicKeySizeKnownCheck struct {
	*AbstractCryptographicCheck

	// keySize is the used key size.
	keySize string
}

// NewPublicKeySizeKnownCheck is the default constructor.
func NewPublicKeySizeKnownCheck(i18nProvider *i18n.I18nProvider, keySize string, result *process.Result[*jaxb.XmlCC],
	position i18n.MessageTag, cryptographicSuite policy.CryptographicSuite) *PublicKeySizeKnownCheck {
	c := &PublicKeySizeKnownCheck{keySize: keySize}
	c.AbstractCryptographicCheck = NewAbstractCryptographicCheck(i18nProvider, result, position,
		process.GetLevelRule(cryptographicSuite.AcceptableSignatureAlgorithmsMiniKeySizeLevel()))
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of the overridden process().
func (c *PublicKeySizeKnownCheck) Process() bool {
	return utils.IsStringDigits(c.keySize)
}

// BuildConstraintMessage builds a constraint message. Port of the overridden
// buildConstraintMessage().
func (c *PublicKeySizeKnownCheck) BuildConstraintMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTagASCCMPKSK)
}

// BuildErrorMessage builds an error message. Port of the overridden
// buildErrorMessage().
func (c *PublicKeySizeKnownCheck) BuildErrorMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTagASCCMPKSKANS, c.position)
}
