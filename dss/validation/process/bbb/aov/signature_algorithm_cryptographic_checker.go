// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/cc/SignatureAlgorithmCryptographicChecker.java (DSS 6.5.RC1).
//
// Package placement deviation: Java's eu.europa.esig.dss.validation.process.bbb.aov.cc
// package (the CryptographicChecker family) is flattened into this pkg aov,
// per the phase 8d porter brief ("cc + both checks subpackages flattened").
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

// SignatureAlgorithmCryptographicChecker runs the cryptographic validation
// for a SignatureAlgorithm with a given key length.
type SignatureAlgorithmCryptographicChecker struct {
	AbstractAlgorithmCryptographicChecker

	// signatureAlgorithm is the Signature algorithm.
	signatureAlgorithm enumerations.SignatureAlgorithm

	// keyLengthUsedToSignThisToken is the used Key length.
	keyLengthUsedToSignThisToken string
}

// NewSignatureAlgorithmCryptographicChecker is the default constructor.
func NewSignatureAlgorithmCryptographicChecker(i18nProvider *i18n.I18nProvider, signatureAlgorithm enumerations.SignatureAlgorithm,
	keyLengthUsedToSignThisToken string, validationDate time.Time, position i18n.MessageTag,
	cryptographicSuite policy.CryptographicSuite) *SignatureAlgorithmCryptographicChecker {
	c := &SignatureAlgorithmCryptographicChecker{
		signatureAlgorithm:           signatureAlgorithm,
		keyLengthUsedToSignThisToken: keyLengthUsedToSignThisToken,
	}
	c.InitAbstractAlgorithmCryptographicChecker(i18nProvider, validationDate, position, cryptographicSuite, c)
	c.InitChainBase(c)
	return c
}

// InitChain initializes the chain. Port of the overridden initChain().
func (c *SignatureAlgorithmCryptographicChecker) InitChain() {
	item := c.signatureAlgorithmReliable()
	c.FirstItem = item

	item = item.SetNextItem(c.publicKeySizeKnown())

	item = item.SetNextItem(c.publicKeySizeAcceptable())

	item.SetNextItem(c.signatureAlgorithmOnValidationTime())
}

// signatureAlgorithmReliable checks if the signatureAlgorithm is acceptable.
// Port of signatureAlgorithmReliable().
func (c *SignatureAlgorithmCryptographicChecker) signatureAlgorithmReliable() process.ChainItem[*jaxb.XmlCC] {
	return NewSignatureAlgorithmReliableCheck(c.I18nProvider, c.signatureAlgorithm, c.Result, c.position, c.cryptographicSuite)
}

// publicKeySizeKnown checks if the keyLengthUsedToSignThisToken is known.
// Port of publicKeySizeKnown().
func (c *SignatureAlgorithmCryptographicChecker) publicKeySizeKnown() process.ChainItem[*jaxb.XmlCC] {
	return NewPublicKeySizeKnownCheck(c.I18nProvider, c.keyLengthUsedToSignThisToken, c.Result, c.position, c.cryptographicSuite)
}

// publicKeySizeAcceptable checks if the keyLengthUsedToSignThisToken is
// acceptable. Port of publicKeySizeAcceptable().
func (c *SignatureAlgorithmCryptographicChecker) publicKeySizeAcceptable() process.ChainItem[*jaxb.XmlCC] {
	return NewPublicKeySizeAcceptableCheck(c.I18nProvider, c.signatureAlgorithm, c.keyLengthUsedToSignThisToken,
		c.Result, c.position, c.cryptographicSuite)
}

// signatureAlgorithmOnValidationTime checks if the signatureAlgorithm is not
// expired in validation time. Port of signatureAlgorithmOnValidationTime().
func (c *SignatureAlgorithmCryptographicChecker) signatureAlgorithmOnValidationTime() process.ChainItem[*jaxb.XmlCC] {
	return NewSignatureAlgorithmAtValidationTimeCheck(c.I18nProvider, c.signatureAlgorithm, c.keyLengthUsedToSignThisToken,
		c.validationDate, c.Result, c.position, c.cryptographicSuite)
}

// Algorithm builds and returns the validated algorithm. Port of the
// overridden getAlgorithm().
func (c *SignatureAlgorithmCryptographicChecker) Algorithm() *jaxb.XmlCryptographicAlgorithm {
	if c.cryptographicAlgorithm == nil {
		c.cryptographicAlgorithm = &jaxb.XmlCryptographicAlgorithm{}
		if c.signatureAlgorithm != "" {
			// if SignatureAlgorithm is defined
			c.cryptographicAlgorithm.Name = c.signatureAlgorithm.Name()
			c.cryptographicAlgorithm.Uri = c.getSignatureAlgorithmUri(c.signatureAlgorithm)
			// Java sets the nullable String straight through, so a null key
			// length leaves the KeyLength element out. TokenProxy's Go form
			// collapses that null to "" (see diagnostic.AbstractTokenProxyBase),
			// so "" is what stands in for it here and must be left out too.
			if c.keyLengthUsedToSignThisToken != "" {
				keyLength := c.keyLengthUsedToSignThisToken
				c.cryptographicAlgorithm.KeyLength = &keyLength
			}

		} else {
			// if SignatureAlgorithm is not found
			c.cryptographicAlgorithm.Name = algorithmUnidentified
			c.cryptographicAlgorithm.Uri = algorithmUnidentifiedURN
		}
	}
	return c.cryptographicAlgorithm
}

// getSignatureAlgorithmUri ports the private getSignatureAlgorithmUri(SignatureAlgorithm).
func (c *SignatureAlgorithmCryptographicChecker) getSignatureAlgorithmUri(signatureAlgorithm enumerations.SignatureAlgorithm) string {
	if signatureAlgorithm != "" {
		if signatureAlgorithm.URI() != "" {
			return signatureAlgorithm.URI()
		}
		if signatureAlgorithm.OID() != "" {
			return signatureAlgorithm.URIBasedOnOID()
		}
	}
	return algorithmUnidentifiedURN
}

// NotAfter returns time after which the used cryptographic algorithm(s) is no
// longer considered secure. Port of the overridden getNotAfter().
func (c *SignatureAlgorithmCryptographicChecker) NotAfter() *time.Time {
	if vpolicy.IsSignatureAlgorithmReliable(c.cryptographicSuite, c.signatureAlgorithm) &&
		vpolicy.IsSignatureAlgorithmWithKeySizeReliable(c.cryptographicSuite, c.signatureAlgorithm, c.keyLengthUsedToSignThisToken) {
		return vpolicy.GetExpirationDateForSignatureAlgorithm(c.cryptographicSuite, c.signatureAlgorithm, c.keyLengthUsedToSignThisToken)
	}
	return nil
}
