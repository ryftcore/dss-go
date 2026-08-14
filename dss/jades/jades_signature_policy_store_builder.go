// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESSignaturePolicyStoreBuilder.java (DSS 6.5.RC1).
//
// Java overloads addSignaturePolicyStore on (document, store, base64UrlInstance) and
// (document, store, base64UrlInstance, signatureId); Go cannot overload, so the second becomes
// AddSignaturePolicyStoreForSignature.
//
// The 'sigPSt' component members are inserted in Java's statement order (sigPolDoc,
// sigPolLocalURI, spDSpec) into a jose.Object, the ordered counterpart of the LinkedHashMap
// upstream hands to addComponent unwrapped: the component may be base64url-encoded into the
// 'etsiU' array, so its byte order is observable.
//
// Objects.requireNonNull becomes a panic carrying the Java message; every
// IllegalArgumentException / IllegalInputException becomes a returned error; slf4j logging is
// dropped.
package jades

import (
	"errors"
	"fmt"

	"github.com/utain/esig/dss/internal/jose"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/exception"
	"github.com/utain/esig/dss/utils"
)

// JAdESSignaturePolicyStoreBuilder is the builder used to incorporate a SignaturePolicyStore into
// a JAdESSignature document.
type JAdESSignaturePolicyStoreBuilder struct {
	JAdESExtensionBuilder
}

// NewJAdESSignaturePolicyStoreBuilder is the default constructor.
func NewJAdESSignaturePolicyStoreBuilder() *JAdESSignaturePolicyStoreBuilder {
	return &JAdESSignaturePolicyStoreBuilder{}
}

// AddSignaturePolicyStore adds signaturePolicyStore to all signatures inside the document
// matching the given SignaturePolicyStore. base64UrlInstance is TRUE when the signature policy
// store shall be incorporated as a base64url-encoded component of the 'etsiU' header, FALSE when
// it shall be incorporated in its clear JSON representation.
// Port of #addSignaturePolicyStore(DSSDocument, SignaturePolicyStore, boolean).
func (b *JAdESSignaturePolicyStoreBuilder) AddSignaturePolicyStore(doc model.DSSDocument,
	signaturePolicyStore *model.SignaturePolicyStore, base64UrlInstance bool) (model.DSSDocument, error) {
	if doc == nil {
		panic("Signature document must be provided!")
	}
	if err := jadesSignaturePolicyStoreBuilderAssertConfigurationValid(signaturePolicyStore); err != nil {
		return nil, err
	}

	documentAnalyzer := jadesSignaturePolicyStoreBuilderDocumentAnalyzer(doc)
	jwsJsonSerializationObject := documentAnalyzer.JwsJsonSerializationObject()
	if err := b.AssertJSONSerializationObjectMayBeExtended(jwsJsonSerializationObject); err != nil {
		return nil, err
	}

	signatures := documentAnalyzer.Signatures()

	signaturePolicyStoreAdded := false
	for _, signature := range signatures {
		jadesSignature, ok := signature.(*JAdESSignature)
		if !ok {
			return nil, fmt.Errorf("unexpected signature type %T", signature)
		}
		added, err := b.AddSignaturePolicyStoreIfDigestMatch(jadesSignature, signaturePolicyStore,
			base64UrlInstance, documentAnalyzer)
		if err != nil {
			return nil, err
		}
		signaturePolicyStoreAdded = signaturePolicyStoreAdded || added
	}
	if !signaturePolicyStoreAdded {
		return nil, exception.NewIllegalInputException(
			"The process did not find a signature to add SignaturePolicyStore!")
	}

	generator := NewJWSJsonSerializationGenerator(jwsJsonSerializationObject,
		jwsJsonSerializationObject.JWSSerializationType())
	return generator.Generate()
}

// AddSignaturePolicyStoreForSignature adds signaturePolicyStore to the signature inside the
// document with the given signatureId.
// Port of #addSignaturePolicyStore(DSSDocument, SignaturePolicyStore, boolean, String).
func (b *JAdESSignaturePolicyStoreBuilder) AddSignaturePolicyStoreForSignature(doc model.DSSDocument,
	signaturePolicyStore *model.SignaturePolicyStore, base64UrlInstance bool,
	signatureId string) (model.DSSDocument, error) {
	if doc == nil {
		panic("Signature document must be provided!")
	}
	if err := jadesSignaturePolicyStoreBuilderAssertConfigurationValid(signaturePolicyStore); err != nil {
		return nil, err
	}

	documentAnalyzer := jadesSignaturePolicyStoreBuilderDocumentAnalyzer(doc)
	jwsJsonSerializationObject := documentAnalyzer.JwsJsonSerializationObject()
	if err := b.AssertJSONSerializationObjectMayBeExtended(jwsJsonSerializationObject); err != nil {
		return nil, err
	}

	signature := documentAnalyzer.SignatureByID(signatureId)
	if signature == nil {
		return nil, exception.NewIllegalInputException(fmt.Sprintf(
			"Unable to find a signature with Id : %s!", signatureId))
	}
	jadesSignature, ok := signature.(*JAdESSignature)
	if !ok {
		return nil, fmt.Errorf("unexpected signature type %T", signature)
	}

	added, err := b.AddSignaturePolicyStoreIfDigestMatch(jadesSignature, signaturePolicyStore,
		base64UrlInstance, documentAnalyzer)
	if err != nil {
		return nil, err
	}
	if !added {
		return nil, exception.NewIllegalInputException(fmt.Sprintf(
			"The process was not able to add SignaturePolicyStore to a signature with Id : %s!", signatureId))
	}

	generator := NewJWSJsonSerializationGenerator(jwsJsonSerializationObject,
		jwsJsonSerializationObject.JWSSerializationType())
	return generator.Generate()
}

// AddSignaturePolicyStoreIfDigestMatch adds the SignaturePolicyStore to jadesSignature if
// required, and reports whether it has been added.
// Port of the protected #addSignaturePolicyStoreIfDigestMatch.
func (b *JAdESSignaturePolicyStoreBuilder) AddSignaturePolicyStoreIfDigestMatch(
	jadesSignature *JAdESSignature, signaturePolicyStore *model.SignaturePolicyStore,
	base64UrlInstance bool, documentAnalyzer *AbstractJWSDocumentAnalyzer) (bool, error) {
	if _, err := b.AssertEtsiUComponentsConsistentWithEncoding(jadesSignature.Jws(),
		&base64UrlInstance); err != nil {
		return false, err
	}

	digestMatch, err := b.CheckDigest(jadesSignature, signaturePolicyStore, documentAnalyzer)
	if err != nil {
		return false, err
	}
	if !digestMatch {
		return false, nil
	}

	sigPolicyStoreParams := jose.NewObject()

	signaturePolicyContent := signaturePolicyStore.SignaturePolicyContent()
	if signaturePolicyContent != nil {
		binaries, err := spi.DSSUtilsToByteArrayOfDocument(signaturePolicyContent)
		if err != nil {
			return false, err
		}
		sigPolicyStoreParams.Put(JAdESHeaderParameterNamesSigPolDoc, utils.ToBase64(binaries))
	}

	sigPolDocLocalURI := signaturePolicyStore.SigPolDocLocalURI()
	if utils.IsStringNotEmpty(sigPolDocLocalURI) {
		sigPolicyStoreParams.Put(JAdESHeaderParameterNamesSigPolLocalURI, sigPolDocLocalURI)
	}

	spDocSpecification := signaturePolicyStore.SpDocSpecification()
	oidObject := DSSJsonUtilsOidObjectFromURI(spDocSpecification.Id(),
		spDocSpecification.Description(), spDocSpecification.DocumentationReferences())
	sigPolicyStoreParams.Put(JAdESHeaderParameterNamesSpDspec, oidObject)

	etsiUHeader := jadesSignature.EtsiUHeader()
	if err := etsiUHeader.AddComponent(JAdESHeaderParameterNamesSigPSt, sigPolicyStoreParams,
		base64UrlInstance); err != nil {
		return false, err
	}

	return true, nil
}

// CheckDigest verifies that the digests computed in the provided SignaturePolicyStore match the
// digest defined in the incorporated signature policy identifier.
// Port of the protected #checkDigest.
func (b *JAdESSignaturePolicyStoreBuilder) CheckDigest(jadesSignature *JAdESSignature,
	signaturePolicyStore *model.SignaturePolicyStore,
	documentAnalyzer *AbstractJWSDocumentAnalyzer) (bool, error) {
	signaturePolicy := jadesSignature.SignaturePolicy()
	if signaturePolicy == nil {
		// Upstream warns "No defined SignaturePolicyIdentifier for signature with Id : {}".
		return false, nil
	}
	expectedDigest := signaturePolicy.Digest()
	if expectedDigest.Algorithm() == "" && expectedDigest.Value() == nil {
		// Upstream warns "No defined digest for signature with Id : {}"; Java's null Digest is
		// the zero value here, since SignaturePolicy#getDigest answers a value type in the port.
		return false, nil
	}

	signaturePolicyContent := signaturePolicyStore.SignaturePolicyContent()
	if signaturePolicyContent == nil {
		// Upstream logs "No policy document has been provided. Digests are not checked!"
		return true, nil
	}
	signaturePolicy.SetPolicyContent(signaturePolicyContent)

	validator := documentAnalyzer.SignaturePolicyValidatorLoader().LoadValidator(signaturePolicy)
	computedDigest := validator.GetComputedDigest(signaturePolicyContent, expectedDigest.Algorithm())

	digestMatch := expectedDigest.Equals(computedDigest)
	// Upstream warns "Signature policy's digest {} doesn't match the digest extracted from
	// document {} for signature with Id : {}" when they differ.
	return digestMatch, nil
}

// jadesSignaturePolicyStoreBuilderDocumentAnalyzer ports the private getDocumentAnalyzer.
func jadesSignaturePolicyStoreBuilderDocumentAnalyzer(
	doc model.DSSDocument) *AbstractJWSDocumentAnalyzer {
	documentAnalyzerFactory := NewJWSDocumentAnalyzerFactory()
	return jwsDocumentAnalyzerBase(documentAnalyzerFactory.Create(doc))
}

// jadesSignaturePolicyStoreBuilderAssertConfigurationValid ports the private
// assertConfigurationValid.
func jadesSignaturePolicyStoreBuilderAssertConfigurationValid(
	signaturePolicyStore *model.SignaturePolicyStore) error {
	if signaturePolicyStore == nil {
		panic("SignaturePolicyStore must be provided")
	}
	if signaturePolicyStore.SpDocSpecification() == nil {
		panic("SpDocSpecification must be provided")
	}
	if signaturePolicyStore.SpDocSpecification().Id() == "" {
		panic("ID (OID or URI) for SpDocSpecification must be provided")
	}

	signaturePolicyContentPresent := signaturePolicyStore.SignaturePolicyContent() != nil
	sigPolDocLocalURIPresent := signaturePolicyStore.SigPolDocLocalURI() != ""
	if signaturePolicyContentPresent == sigPolDocLocalURIPresent {
		return errors.New("SignaturePolicyStore shall contain either " +
			"SignaturePolicyContent document or sigPolDocLocalURI!")
	}
	return nil
}
