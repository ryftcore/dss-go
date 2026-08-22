// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/TokenAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Performs cryptographic validation for a given TokenProxy, including the
// signature value, signed properties and certificate chain validation, when
// applicable. T is a diagnostic.TokenProxy implementation.
package aov

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TokenAlgorithmObsolescenceValidation is the Go form of the abstract Java
// class TokenAlgorithmObsolescenceValidation<T extends TokenProxy>. Java's
// constructor is a pure pass-through to the parent, so no additional wiring
// method is needed here.
type TokenAlgorithmObsolescenceValidation[T diagnostic.TokenProxy] struct {
	DigestAlgorithmObsolescenceValidation[T]
}

// BuildChain is the concrete override of the abstract buildChain() this class
// provides: it is inherited, unmodified, by subclasses that do not further
// override it (e.g. CertificateAlgorithmObsolescenceValidation). Port of the
// overridden protected ChainItem<XmlAOV> buildChain().
func (c *TokenAlgorithmObsolescenceValidation[T]) BuildChain() process.ChainItem[*jaxb.XmlAOV] {
	return c.buildSignatureValidationChain(c.FirstItem)
}

// buildSignatureValidationChain builds a chain of crypto checks to be
// executed on a signature. Port of buildSignatureValidationChain(ChainItem).
func (c *TokenAlgorithmObsolescenceValidation[T]) buildSignatureValidationChain(item process.ChainItem[*jaxb.XmlAOV]) process.ChainItem[*jaxb.XmlAOV] {
	cc := NewSignatureAlgorithmCryptographicChecker(c.I18nProvider, c.token.SignatureAlgorithm(),
		c.token.KeyLengthUsedToSignThisToken(), c.validationDate, c.position, c.cryptographicSuite)
	ccResult := cc.Execute()

	if item == nil {
		item = c.signatureAlgorithmCryptographicCheckResult(ccResult, c.position, c.cryptographicSuite)
		c.FirstItem = item
	} else {
		item = item.SetNextItem(c.signatureAlgorithmCryptographicCheckResult(ccResult, c.position, c.cryptographicSuite))
	}

	c.signatureCryptographicValidation = ccResult.CryptographicValidation
	if c.signatureCryptographicValidation != nil {
		tokenId := c.token.Id()
		c.signatureCryptographicValidation.TokenId = &tokenId
	}

	return item
}

// buildCertificateChainValidationChain builds a chain of crypto checks to be
// executed on a signature's certificate chain. Port of
// buildCertificateChainValidationChain(ChainItem, CertificateWrapper, List).
func (c *TokenAlgorithmObsolescenceValidation[T]) buildCertificateChainValidationChain(item process.ChainItem[*jaxb.XmlAOV],
	signingCertificate *diagnostic.CertificateWrapper, certificateChain []*diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlAOV] {
	if signingCertificate == nil {
		return item
	}

	c.certificateChainCryptographicValidation = nil

	for _, certificate := range certificateChain {
		var subContext enumerations.SubContext
		if signingCertificate.Equals(certificate) {
			subContext = enumerations.SubContextSigningCert
		} else {
			subContext = enumerations.SubContextCACertificate
		}

		if c.isTrustAnchor(certificate, subContext) {
			break
		}

		certificateCryptographicSuite := c.validationPolicy.CertificateCryptographicConstraint(c.context, subContext)
		certificatePosition, err := process.GetSubContextPosition(c.context, subContext)
		if err != nil {
			panic(err)
		}

		cc := NewSignatureAlgorithmCryptographicChecker(c.I18nProvider, certificate.SignatureAlgorithm(),
			certificate.KeyLengthUsedToSignThisToken(), c.validationDate, certificatePosition, certificateCryptographicSuite)
		ccResult := cc.Execute()

		if item == nil {
			item = c.signatureAlgorithmCryptographicCheckResultWithTokenId(ccResult, certificatePosition, certificateCryptographicSuite, certificate.Id())
			c.FirstItem = item
		} else {
			item = item.SetNextItem(c.signatureAlgorithmCryptographicCheckResultWithTokenId(ccResult, certificatePosition, certificateCryptographicSuite, certificate.Id()))
		}

		certificateCryptographicValidation := ccResult.CryptographicValidation
		tokenId := certificate.Id()
		certificateCryptographicValidation.TokenId = &tokenId

		c.certificateChainCryptographicValidation = append(c.certificateChainCryptographicValidation, certificateCryptographicValidation)
	}

	return item
}

// signatureAlgorithmCryptographicCheckResult ports the private
// signatureAlgorithmCryptographicCheckResult(XmlCC, MessageTag, CryptographicSuite).
func (c *TokenAlgorithmObsolescenceValidation[T]) signatureAlgorithmCryptographicCheckResult(ccResult *jaxb.XmlCC,
	position i18n.MessageTag, cryptographicSuite policy.CryptographicSuite) process.ChainItem[*jaxb.XmlAOV] {
	return NewSignatureAlgorithmCryptographicCheckerResultCheck(c.I18nProvider, c.Result, c.validationDate, position, ccResult, cryptographicSuite)
}

// signatureAlgorithmCryptographicCheckResultWithTokenId ports the private
// signatureAlgorithmCryptographicCheckResult(XmlCC, MessageTag, CryptographicSuite, String).
func (c *TokenAlgorithmObsolescenceValidation[T]) signatureAlgorithmCryptographicCheckResultWithTokenId(ccResult *jaxb.XmlCC,
	position i18n.MessageTag, cryptographicSuite policy.CryptographicSuite, tokenId string) process.ChainItem[*jaxb.XmlAOV] {
	return NewSignatureAlgorithmCryptographicCheckerResultCheckWithContext(c.I18nProvider, c.Result, c.validationDate,
		enumerations.ContextCertificate, position, ccResult, cryptographicSuite, tokenId)
}

// isTrustAnchor ports the private isTrustAnchor(CertificateWrapper, SubContext).
func (c *TokenAlgorithmObsolescenceValidation[T]) isTrustAnchor(certificateWrapper *diagnostic.CertificateWrapper,
	subContext enumerations.SubContext) bool {
	constraint := c.validationPolicy.CertificateSunsetDateConstraint(c.context, subContext)
	return process.IsTrustAnchor(certificateWrapper, c.validationDate, constraint)
}
