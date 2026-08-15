// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// Performs cryptographic validation of a token's signature value algorithm
// and algorithms used within the signed properties, when applicable. T is a
// diagnostic.TokenProxy implementation.
package aov

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation performs
// cryptographic validation of a token's signature value algorithm and
// algorithms used within the signed properties, when applicable.
type SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation[T diagnostic.TokenProxy] struct {
	TokenAlgorithmObsolescenceValidation[T]
}

// NewSignatureValueAndSignedAttributesAlgorithmObsolescenceValidation is the
// default constructor.
func NewSignatureValueAndSignedAttributesAlgorithmObsolescenceValidation[T diagnostic.TokenProxy](i18nProvider *i18n.I18nProvider,
	token T, context enumerations.Context, validationDate time.Time,
	validationPolicy policy.ValidationPolicy) *SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation[T] {
	c := &SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation[T]{}
	c.InitAlgorithmObsolescenceValidation(i18nProvider, token, context, validationDate, validationPolicy, c)
	c.InitChainBase(c)
	return c
}

// BuildChain builds a chain of checks to be executed during the process. Port
// of the overridden protected ChainItem<XmlAOV> buildChain().
func (c *SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation[T]) BuildChain() process.ChainItem[*jaxb.XmlAOV] {
	item := c.TokenAlgorithmObsolescenceValidation.BuildChain()
	item = c.buildSignedAttributesValidationChain(item)
	return item
}

// buildSignedAttributesValidationChain builds a chain of crypto checks to be
// executed on a signature's signed attributes. Port of
// buildSignedAttributesValidationChain(ChainItem).
func (c *SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation[T]) buildSignedAttributesValidationChain(
	item process.ChainItem[*jaxb.XmlAOV]) process.ChainItem[*jaxb.XmlAOV] {
	if !c.token.IsSigningCertificateReferencePresent() {
		return item
	}

	signingCertificateReferences := c.token.SigningCertificateReferences()
	signingCertificateReference := c.token.SigningCertificateReference()

	var cryptographicValidation *jaxb.XmlCryptographicValidation

	// This code ensures that at least one good digest algorithm is found for
	// every defined signing certificate reference. Java iterates a HashMap
	// keyed by certificate id, whose order is unspecified; this port groups
	// by first-seen certificate id instead, for deterministic output (see
	// PORTING.md).
	var certificateIds []string
	signCertRefsMap := make(map[string][]*diagnostic.CertificateRefWrapper)
	for _, r := range signingCertificateReferences {
		id := r.CertificateId()
		if _, ok := signCertRefsMap[id]; !ok {
			certificateIds = append(certificateIds, id)
		}
		signCertRefsMap[id] = append(signCertRefsMap[id], r)
	}

	for _, certificateId := range certificateIds {
		certificateRefWrappers := signCertRefsMap[certificateId]

		var subContext enumerations.SubContext
		if signingCertificateReference != nil && signingCertificateReference.CertificateId() == certificateId {
			subContext = enumerations.SubContext_SIGNING_CERT
		} else {
			subContext = enumerations.SubContext_CA_CERTIFICATE
		}

		signCertCheck := c.signingCertificateRefDigestAlgoCheckResult(certificateRefWrappers, certificateId, subContext)

		if item == nil {
			item = signCertCheck
			c.FirstItem = item
		} else {
			item = item.SetNextItem(signCertCheck)
		}

		cryptoValidationResult := signCertCheck.CryptographicValidationResult()
		if cryptographicValidation == nil || (c.isValid(cryptographicValidation) &&
			!c.IsValid(&cryptoValidationResult.XmlConstraintsConclusionContent)) {
			cryptographicValidation = cryptoValidationResult.CryptographicValidation
			tokenId := certificateId
			cryptographicValidation.TokenId = &tokenId
		}
	}

	c.signedAttributesCryptographicValidation = cryptographicValidation
	if c.signedAttributesCryptographicValidation != nil {
		tokenId := c.token.Id()
		c.signedAttributesCryptographicValidation.TokenId = &tokenId
	}

	return item
}

// signingCertificateRefDigestAlgoCheckResult ports the private
// signingCertificateRefDigestAlgoCheckResult(List, String, SubContext).
func (c *SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation[T]) signingCertificateRefDigestAlgoCheckResult(
	signCertRefs []*diagnostic.CertificateRefWrapper, certificateId string,
	subContext enumerations.SubContext) *SigningCertificateRefDigestAlgorithmCheck[*jaxb.XmlAOV] {
	constraint := c.validationPolicy.SigningCertificateDigestAlgorithmConstraint(c.context)
	return NewSigningCertificateRefDigestAlgorithmCheck(c.I18nProvider, c.Result, c.validationDate,
		signCertRefs, certificateId, c.context, subContext, c.validationPolicy, constraint)
}
