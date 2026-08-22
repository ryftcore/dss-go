// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/cc/checks/DigestAlgorithmReliableCheck.java (DSS 6.5.RC1).
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

// DigestAlgorithmReliableCheck checks if DigestAlgorithm is acceptable.
type DigestAlgorithmReliableCheck struct {
	*AbstractCryptographicCheck

	// digestAlgo is the algorithm to check.
	digestAlgo enumerations.DigestAlgorithm

	// cryptographicSuite is the cryptographic rules.
	cryptographicSuite policy.CryptographicSuite
}

// NewDigestAlgorithmReliableCheck is the default constructor.
func NewDigestAlgorithmReliableCheck(i18nProvider *i18n.I18nProvider, digestAlgo enumerations.DigestAlgorithm,
	result *process.Result[*jaxb.XmlCC], position i18n.MessageTag,
	cryptographicSuite policy.CryptographicSuite) *DigestAlgorithmReliableCheck {
	c := &DigestAlgorithmReliableCheck{
		digestAlgo:         digestAlgo,
		cryptographicSuite: cryptographicSuite,
	}
	c.AbstractCryptographicCheck = NewAbstractCryptographicCheck(i18nProvider, result, position,
		process.GetLevelRule(cryptographicSuite.AcceptableDigestAlgorithmsLevel()))
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of the overridden process().
func (c *DigestAlgorithmReliableCheck) Process() bool {
	return vpolicy.IsDigestAlgorithmReliable(c.cryptographicSuite, c.digestAlgo)
}

// BuildConstraintMessage builds a constraint message. Port of the overridden
// buildConstraintMessage().
func (c *DigestAlgorithmReliableCheck) BuildConstraintMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(i18n.MessageTagASCCMDAA, c.Name(c.digestAlgo))
}

// BuildErrorMessage builds an error message. Port of the overridden
// buildErrorMessage().
func (c *DigestAlgorithmReliableCheck) BuildErrorMessage() *jaxb.XmlMessage {
	if c.digestAlgo == "" {
		return c.BuildXmlMessage(i18n.MessageTagASCCMDAAANS2, c.position)
	}
	return c.BuildXmlMessage(i18n.MessageTagASCCMDAAANS, c.Name(c.digestAlgo), c.position)
}
