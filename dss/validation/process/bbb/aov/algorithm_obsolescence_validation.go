// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/AlgorithmObsolescenceValidation.java (DSS 6.5.RC1).
//
// This performs validation of cryptographic algorithms against the provided
// cryptographic suite constraints; the result of this block is used within
// the XCV and SAV building blocks. T is the validation token wrapper type.
package aov

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AlgorithmObsolescenceValidationOverrides captures the members Java's
// AlgorithmObsolescenceValidation<T> treats virtually: the abstract
// buildChain(), and the overridable getPosition()/getCryptographicSuite()
// (CertificateAlgorithmObsolescenceValidation, EAAAlgorithmObsolescenceValidation
// and EvidenceRecordAlgorithmObsolescenceValidation override one or both of the
// latter). initChain() (defined once, below) dispatches to these through the
// registered overrides, the way process.Chain dispatches through ChainOverrides.
type AlgorithmObsolescenceValidationOverrides interface {
	// Position gets position of the currently verifying token, based on the
	// context. Port of getPosition().
	Position() i18n.MessageTag
	// CryptographicSuite gets the cryptographic suite based on the currently
	// validating context. Port of getCryptographicSuite().
	CryptographicSuite() policy.CryptographicSuite
	// BuildChain builds a chain of checks to be executed during the process.
	// Port of the abstract buildChain().
	BuildChain() process.ChainItem[*jaxb.XmlAOV]
}

// AlgorithmObsolescenceValidation is the Go form of the abstract Java class
// AlgorithmObsolescenceValidation<T>.
type AlgorithmObsolescenceValidation[T any] struct {
	*process.ChainBase[*jaxb.XmlAOV]

	// token is the token to be validated.
	token T

	// context is the validation context.
	context enumerations.Context

	// validationDate is the validation time.
	validationDate time.Time

	// validationPolicy contains the corresponding cryptographic constraints.
	validationPolicy policy.ValidationPolicy

	// position is the position of the token, resolved once in InitChain.
	position i18n.MessageTag

	// cryptographicSuite is the cryptographic suite for the token's
	// validation, resolved once in InitChain.
	cryptographicSuite policy.CryptographicSuite

	// signatureCryptographicValidation is the cryptographic information for
	// the SignatureValue algorithms validation.
	signatureCryptographicValidation *jaxb.XmlCryptographicValidation

	// signedAttributesCryptographicValidation is the cryptographic
	// information for the signed attributes algorithms validation.
	signedAttributesCryptographicValidation *jaxb.XmlCryptographicValidation

	// digestMatchersCryptographicValidation is the cryptographic information
	// for the digest matchers (signed references) algorithms validation.
	digestMatchersCryptographicValidation *jaxb.XmlCryptographicValidation

	// certificateChainCryptographicValidation is the cryptographic
	// information for the certificate chain algorithms validation.
	certificateChainCryptographicValidation []*jaxb.XmlCryptographicValidation

	// overrides points back at the concrete AOV type; see
	// InitAlgorithmObsolescenceValidation.
	overrides AlgorithmObsolescenceValidationOverrides
}

// InitAlgorithmObsolescenceValidation wires the shared state; called by the
// concrete constructor before InitChainBase. Port of the common constructor
// AlgorithmObsolescenceValidation(I18nProvider, T, Context, Date, ValidationPolicy).
func (c *AlgorithmObsolescenceValidation[T]) InitAlgorithmObsolescenceValidation(i18nProvider *i18n.I18nProvider,
	token T, context enumerations.Context, validationDate time.Time, validationPolicy policy.ValidationPolicy,
	overrides AlgorithmObsolescenceValidationOverrides) {
	xmlAOV := &jaxb.XmlAOV{}
	result := process.NewResult(xmlAOV, &xmlAOV.XmlConstraintsConclusionContent, &xmlAOV.XmlConstraintsConclusionAttrs)
	c.ChainBase = process.NewChainBase(i18nProvider, result)
	c.token = token
	c.context = context
	c.validationDate = validationDate
	c.validationPolicy = validationPolicy
	c.overrides = overrides
}

// Title returns the title of the chain (i.e. the BasicBuildingBlock title).
// Port of the overridden protected MessageTag getTitle().
func (c *AlgorithmObsolescenceValidation[T]) Title() i18n.MessageTag {
	return i18n.MessageTagAOV
}

// InitChain initializes the chain. Port of the overridden protected void
// initChain().
func (c *AlgorithmObsolescenceValidation[T]) InitChain() {
	c.position = c.overrides.Position()
	c.cryptographicSuite = c.overrides.CryptographicSuite()

	// Java discards buildChain()'s return value: the builders it calls set
	// firstItem themselves, on the first item they append, and return the LAST
	// item of the chain so that a caller can keep appending to it. Assigning
	// the return value to FirstItem here would leave only the final check
	// reachable.
	c.overrides.BuildChain()
}

// Position is the default implementation of getPosition(). Port of protected
// MessageTag getPosition().
func (c *AlgorithmObsolescenceValidation[T]) Position() i18n.MessageTag {
	position, err := process.GetCryptoPosition(c.context)
	if err != nil {
		panic(err)
	}
	return position
}

// CryptographicSuite is the default implementation of getCryptographicSuite().
// Port of protected CryptographicSuite getCryptographicSuite().
func (c *AlgorithmObsolescenceValidation[T]) CryptographicSuite() policy.CryptographicSuite {
	return c.validationPolicy.SignatureCryptographicConstraint(c.context)
}

// isValid checks whether CryptographicValidation returned a successful
// validation result. Port of protected boolean isValid(XmlCryptographicValidation).
func (c *AlgorithmObsolescenceValidation[T]) isValid(cryptographicValidation *jaxb.XmlCryptographicValidation) bool {
	return cryptographicValidation != nil && cryptographicValidation.Conclusion != nil &&
		enumerations.IndicationPassed == cryptographicValidation.Conclusion.Indication.Indication()
}

// AddAdditionalInfo adds additional info to the chain. Port of the overridden
// protected void addAdditionalInfo(). Java's super.addAdditionalInfo() call is
// the Chain default, which is empty, and is therefore not ported.
func (c *AlgorithmObsolescenceValidation[T]) AddAdditionalInfo() {
	result := c.Result.Value
	result.ValidationTime = jaxb.NewXSDateTime(c.validationDate)
	result.SignatureCryptographicValidation = c.signatureCryptographicValidation
	result.SignedAttributesValidation = c.signedAttributesCryptographicValidation
	result.DigestMatchersValidation = c.digestMatchersCryptographicValidation
	if utils.IsCollectionNotEmpty(c.certificateChainCryptographicValidation) {
		xmlCertificateChainCryptographicValidation := &jaxb.XmlCertificateChainCryptographicValidation{}
		xmlCertificateChainCryptographicValidation.CertificateCryptographicValidation = append(
			xmlCertificateChainCryptographicValidation.CertificateCryptographicValidation,
			c.certificateChainCryptographicValidation...)
		result.CertificateChainCryptographicValidation = xmlCertificateChainCryptographicValidation
	}
}
