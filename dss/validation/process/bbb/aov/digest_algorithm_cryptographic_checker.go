// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/cc/DigestAlgorithmCryptographicChecker.java (DSS 6.5.RC1).
//
// Package placement deviation: see signature_algorithm_cryptographic_checker.go.
package aov

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	vpolicy "github.com/utain/esig/dss/validation/policy"
	"github.com/utain/esig/dss/validation/process"
)

// DigestAlgorithmCryptographicChecker checks the digest algorithm.
type DigestAlgorithmCryptographicChecker struct {
	AbstractAlgorithmCryptographicChecker

	// digestAlgorithm is the Digest algorithm.
	digestAlgorithm enumerations.DigestAlgorithm
}

// NewDigestAlgorithmCryptographicChecker is the default constructor.
func NewDigestAlgorithmCryptographicChecker(i18nProvider *i18n.I18nProvider, digestAlgorithm enumerations.DigestAlgorithm,
	validationDate time.Time, position i18n.MessageTag, constraint policy.CryptographicSuite) *DigestAlgorithmCryptographicChecker {
	c := &DigestAlgorithmCryptographicChecker{digestAlgorithm: digestAlgorithm}
	c.InitAbstractAlgorithmCryptographicChecker(i18nProvider, validationDate, position, constraint, c)
	c.InitChainBase(c)
	return c
}

// InitChain initializes the chain. Port of the overridden initChain().
func (c *DigestAlgorithmCryptographicChecker) InitChain() {
	item := c.digestAlgorithmReliable()
	c.FirstItem = item

	item.SetNextItem(c.digestAlgorithmOnValidationTime())
}

// digestAlgorithmReliable checks if the digestAlgorithm is acceptable. Port
// of digestAlgorithmReliable().
func (c *DigestAlgorithmCryptographicChecker) digestAlgorithmReliable() process.ChainItem[*jaxb.XmlCC] {
	return NewDigestAlgorithmReliableCheck(c.I18nProvider, c.digestAlgorithm, c.Result, c.position, c.cryptographicSuite)
}

// digestAlgorithmOnValidationTime checks if the digestAlgorithm is not
// expired in validation time. Port of digestAlgorithmOnValidationTime().
func (c *DigestAlgorithmCryptographicChecker) digestAlgorithmOnValidationTime() process.ChainItem[*jaxb.XmlCC] {
	return NewDigestAlgorithmAtValidationTimeCheck(c.I18nProvider, c.digestAlgorithm, c.validationDate, c.Result, c.position, c.cryptographicSuite)
}

// Algorithm builds and returns the validated algorithm. Port of the
// overridden getAlgorithm().
func (c *DigestAlgorithmCryptographicChecker) Algorithm() *jaxb.XmlCryptographicAlgorithm {
	if c.cryptographicAlgorithm == nil {
		c.cryptographicAlgorithm = &jaxb.XmlCryptographicAlgorithm{}
		if c.digestAlgorithm != "" {
			// if DigestAlgorithm is defined
			c.cryptographicAlgorithm.Name = c.digestAlgorithm.Name()
			c.cryptographicAlgorithm.Uri = c.getDigestAlgorithmUri(c.digestAlgorithm)

		} else {
			// if DigestAlgorithm is not found (unable to build either SignatureAlgorithm nor DigestAlgorithm)
			c.cryptographicAlgorithm.Name = algorithmUnidentified
			c.cryptographicAlgorithm.Uri = algorithmUnidentifiedURN
		}
	}
	return c.cryptographicAlgorithm
}

// getDigestAlgorithmUri ports the private getDigestAlgorithmUri(DigestAlgorithm).
func (c *DigestAlgorithmCryptographicChecker) getDigestAlgorithmUri(digestAlgorithm enumerations.DigestAlgorithm) string {
	if digestAlgorithm != "" {
		if digestAlgorithm.URI() != "" {
			return digestAlgorithm.URI()
		}
		if digestAlgorithm.OID() != "" {
			return digestAlgorithm.OID()
		}
	}
	return algorithmUnidentifiedURN
}

// NotAfter returns time after which the used cryptographic algorithm(s) is no
// longer considered secure. Port of the overridden getNotAfter().
func (c *DigestAlgorithmCryptographicChecker) NotAfter() *time.Time {
	if vpolicy.IsDigestAlgorithmReliable(c.cryptographicSuite, c.digestAlgorithm) {
		return vpolicy.GetExpirationDateForDigestAlgorithm(c.cryptographicSuite, c.digestAlgorithm)
	}
	return nil
}
