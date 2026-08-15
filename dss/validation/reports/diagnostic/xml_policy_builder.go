// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/reports/diagnostic/XmlPolicyBuilder.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/signature"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
)

// XmlPolicyBuilder is used to validate a SignaturePolicy and build a XmlPolicy.
type XmlPolicyBuilder struct {
	// signaturePolicy is the SignaturePolicy to incorporate into the DiagnosticData.
	signaturePolicy *signature.SignaturePolicy

	// signaturePolicyStore is the found SignaturePolicyStore from a signature.
	signaturePolicyStore *model.SignaturePolicyStore
}

// NewXmlPolicyBuilder is the port of the default constructor: panics (Java requireNonNull) if
// signaturePolicy is nil.
func NewXmlPolicyBuilder(signaturePolicy *signature.SignaturePolicy) *XmlPolicyBuilder {
	if signaturePolicy == nil {
		panic("SignaturePolicy cannot be null!")
	}
	return &XmlPolicyBuilder{signaturePolicy: signaturePolicy}
}

// SetSignaturePolicyStore sets the SignaturePolicyStore extracted from a signature when
// applicable. Port of setSignaturePolicyStore(SignaturePolicyStore).
func (b *XmlPolicyBuilder) SetSignaturePolicyStore(signaturePolicyStore *model.SignaturePolicyStore) *XmlPolicyBuilder {
	b.signaturePolicyStore = signaturePolicyStore
	return b
}

// Build validates a SignaturePolicy and builds an XmlPolicy. Port of build().
func (b *XmlPolicyBuilder) Build() *jaxb.XmlPolicy {
	xmlPolicy := &jaxb.XmlPolicy{}

	id := b.signaturePolicy.Identifier()
	xmlPolicy.Id = &id
	description := b.signaturePolicy.Description()
	xmlPolicy.Description = &description
	if utils.IsCollectionNotEmpty(b.signaturePolicy.DocumentationReferences()) {
		xmlPolicy.DocumentationReferences = &jaxb.DocumentationReferencesWrapper{Items: b.signaturePolicy.DocumentationReferences()}
	}

	url := spi.DSSUtilsRemoveControlCharacters(b.signaturePolicy.URI())
	xmlPolicy.Url = &url
	userNotice := b.signaturePolicy.UserNotice()
	if userNotice != nil {
		xmlUserNotice := &jaxb.XmlUserNotice{}
		organization := userNotice.Organization()
		xmlUserNotice.Organization = &organization
		if len(userNotice.NoticeNumbers()) > 0 {
			xmlUserNotice.NoticeNumbers = &jaxb.BigIntegerList{}
			*xmlUserNotice.NoticeNumbers = append(*xmlUserNotice.NoticeNumbers, spi.DSSUtilsToBigIntegerList(userNotice.NoticeNumbers())...)
		}
		explicitText := userNotice.ExplicitText()
		xmlUserNotice.ExplicitText = &explicitText
		xmlPolicy.UserNotice = xmlUserNotice
	}
	spDocSpecification := b.signaturePolicy.DocSpecification()
	if spDocSpecification != nil {
		xmlSPDocSpecification := &jaxb.XmlSPDocSpecification{}
		id := spDocSpecification.Id()
		xmlSPDocSpecification.Id = &id
		description := spDocSpecification.Description()
		xmlSPDocSpecification.Description = &description
		documentationReferences := spDocSpecification.DocumentationReferences()
		if utils.IsArrayNotEmpty(documentationReferences) {
			xmlSPDocSpecification.DocumentationReferences = &jaxb.DocumentationReferencesWrapper{Items: documentationReferences}
		}
		xmlPolicy.DocSpecification = xmlSPDocSpecification
	}

	transformsDescription := b.signaturePolicy.TransformsDescription()
	if utils.IsCollectionNotEmpty(transformsDescription) {
		xmlPolicy.Transformations = &jaxb.TransformationsWrapper{Items: transformsDescription}
	}

	validationResult := b.signaturePolicy.ValidationResult()

	xmlPolicyDigestAlgoAndValue := &jaxb.XmlPolicyDigestAlgoAndValue{}
	if b.signaturePolicy.IsZeroHash() {
		zeroHash := b.signaturePolicy.IsZeroHash()
		xmlPolicyDigestAlgoAndValue.ZeroHash = &zeroHash
	} else {
		digestAlgorithmsEqual := validationResult.IsDigestAlgorithmsEqual()
		xmlPolicyDigestAlgoAndValue.DigestAlgorithmsEqual = &digestAlgorithmsEqual
	}
	digest := b.signaturePolicy.Digest()
	if !digest.IsEmpty() {
		xmlDigestAlgoAndValue := b.getXmlDigestAlgoAndValue(digest)
		xmlPolicyDigestAlgoAndValue.DigestMethod = xmlDigestAlgoAndValue.DigestMethod
		xmlPolicyDigestAlgoAndValue.DigestValue = xmlDigestAlgoAndValue.DigestValue
	}
	match := validationResult.IsDigestValid()
	xmlPolicyDigestAlgoAndValue.Match = &match
	xmlPolicy.DigestAlgoAndValue = xmlPolicyDigestAlgoAndValue

	asn1Processable := validationResult.IsAsn1Processable()
	xmlPolicy.Asn1Processable = &asn1Processable
	identified := validationResult.IsIdentified()
	xmlPolicy.Identified = &identified
	if utils.IsStringNotBlank(validationResult.ProcessingErrors()) {
		processingErrors := validationResult.ProcessingErrors()
		xmlPolicy.ProcessingError = &processingErrors
	}

	return xmlPolicy
}

// BuildSignaturePolicyStore builds an XmlSignaturePolicyStore. Port of buildSignaturePolicyStore().
func (b *XmlPolicyBuilder) BuildSignaturePolicyStore() *jaxb.XmlSignaturePolicyStore {
	if b.signaturePolicyStore == nil {
		return nil
	}

	xmlSignaturePolicyStore := &jaxb.XmlSignaturePolicyStore{}
	spDocSpecification := b.signaturePolicyStore.SpDocSpecification()
	if spDocSpecification != nil {
		id := spDocSpecification.Id()
		xmlSignaturePolicyStore.Id = &id
		description := spDocSpecification.Description()
		xmlSignaturePolicyStore.Description = &description
		documentationReferences := spDocSpecification.DocumentationReferences()
		if utils.IsArrayNotEmpty(documentationReferences) {
			xmlSignaturePolicyStore.DocumentationReferences = &jaxb.DocumentationReferencesWrapper{Items: documentationReferences}
		}
	}
	signaturePolicyContent := b.signaturePolicyStore.SignaturePolicyContent()
	if signaturePolicyContent != nil {
		validationResult := b.signaturePolicy.ValidationResult()
		recalculatedDigest := validationResult.Digest()
		if !recalculatedDigest.IsEmpty() {
			xmlSignaturePolicyStore.DigestAlgoAndValue = b.getXmlDigestAlgoAndValue(recalculatedDigest)
		}
	}
	sigPolDocLocalURI := b.signaturePolicyStore.SigPolDocLocalURI()
	xmlSignaturePolicyStore.SigPolDocLocalURI = &sigPolDocLocalURI

	return xmlSignaturePolicyStore
}

func (b *XmlPolicyBuilder) getXmlDigestAlgoAndValue(digest model.Digest) *jaxb.XmlDigestAlgoAndValue {
	xmlDigestAlgAndValue := &jaxb.XmlDigestAlgoAndValue{}
	algo := jaxb.DigestAlgorithmValue(digest.Algorithm())
	xmlDigestAlgAndValue.DigestMethod = &algo
	value := digest.Value()
	if value == nil {
		value = spi.DSSUtilsEmptyByteArray
	}
	bin := jaxb.Base64Binary(value)
	xmlDigestAlgAndValue.DigestValue = &bin
	return xmlDigestAlgAndValue
}
