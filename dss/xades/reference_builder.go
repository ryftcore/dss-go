// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/ReferenceBuilder.java (DSS 6.5.RC1).
//
// # Errors
//
// Java's two constructors validate with Objects.requireNonNull and therefore keep bare returns
// here, panicking with the Java message (PORTING.md). build() and the private per-packaging
// helpers gain an error return: DomUtils.buildDOM, which Java performs behind a
// DSSException-free API, can fail in Go, and the IllegalInputException / DSSException upstream
// throws become returned errors. slf4j logging is dropped.
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// referenceBuilderObjectIdPrefix is the id-prefix for the ds:Object element. Port of the
// private OBJECT_ID_PREFIX constant.
const referenceBuilderObjectIdPrefix = "o-"

// ReferenceBuilder creates references and processes defined transformations.
type ReferenceBuilder struct {
	// documents is the list of documents to create references for.
	documents []model.DSSDocument

	// digestAlgorithm is the DigestAlgorithm to use.
	digestAlgorithm enumerations.DigestAlgorithm

	// referenceIdProvider creates an identifier for a signature reference.
	referenceIdProvider *ReferenceIdProvider

	// signatureParameters are the used SignatureParameters. Nil for the detached
	// references constructor.
	signatureParameters *SignatureParameters
}

// NewReferenceBuilder is the default constructor for a signature references creation. Ports
// ReferenceBuilder(List, SignatureParameters, ReferenceIdProvider).
func NewReferenceBuilder(documents []model.DSSDocument, xadesSignatureParameters *SignatureParameters,
	referenceIdProvider *ReferenceIdProvider) *ReferenceBuilder {
	if documents == nil {
		panic("List of documents shall be provided!")
	}
	if xadesSignatureParameters == nil {
		panic("Signature parameters shall be provided!")
	}
	return &ReferenceBuilder{
		documents:           documents,
		signatureParameters: xadesSignatureParameters,
		digestAlgorithm:     DSSXMLUtilsGetReferenceDigestAlgorithmOrDefault(xadesSignatureParameters),
		referenceIdProvider: referenceIdProvider,
	}
}

// NewReferenceBuilderWithDigestAlgorithm is the default constructor for a detached references
// creation. Ports ReferenceBuilder(List, DigestAlgorithm, ReferenceIdProvider).
func NewReferenceBuilderWithDigestAlgorithm(documents []model.DSSDocument,
	digestAlgorithm enumerations.DigestAlgorithm, referenceIdProvider *ReferenceIdProvider) *ReferenceBuilder {
	if documents == nil {
		panic("List of documents shall be provided!")
	}
	if digestAlgorithm == "" {
		panic("Digest Algorithm shall be provided!")
	}
	return &ReferenceBuilder{
		documents:           documents,
		digestAlgorithm:     digestAlgorithm,
		referenceIdProvider: referenceIdProvider,
	}
}

// Build builds a list of references based on the configuration. Ports build().
func (b *ReferenceBuilder) Build() ([]*DSSReference, error) {
	references := make([]*DSSReference, 0, len(b.documents))
	for _, dssDocument := range b.documents {
		reference, err := b.createDSSReferenceForDocument(dssDocument)
		if err != nil {
			return nil, err
		}
		references = append(references, reference)
	}
	return references, nil
}

// createDSSReferenceForDocument ports the private createDSSReferenceForDocument(DSSDocument).
func (b *ReferenceBuilder) createDSSReferenceForDocument(document model.DSSDocument) (*DSSReference, error) {
	if b.signatureParameters != nil {
		if b.signatureParameters.SignaturePackaging() == "" {
			panic("SignaturePackaging must be defined!")
		}
		switch b.signatureParameters.SignaturePackaging() {
		case enumerations.SignaturePackagingEnveloped:
			return b.envelopedDSSReference(document)
		case enumerations.SignaturePackagingEnveloping:
			return b.envelopingDSSReference(document)
		case enumerations.SignaturePackagingDetached:
			return b.detachedDSSReference(document), nil
		case enumerations.SignaturePackagingInternallyDetached:
			return b.internallyDetachedDSSReference(document)
		default:
			return nil, model.NewDSSError(fmt.Sprintf("The given signature packaging %s is not supported!",
				b.signatureParameters.SignaturePackaging()))
		}
	}
	// detached reference creation
	return b.detachedDSSReference(document), nil
}

// envelopedDSSReference ports the private envelopedDSSReference(DSSDocument).
func (b *ReferenceBuilder) envelopedDSSReference(document model.DSSDocument) (*DSSReference, error) {
	if err := b.assertEnvelopedSignaturePossible(document); err != nil {
		return nil, err
	}

	dssReference := NewDSSReference()
	dssReference.SetId(b.referenceIdProvider.ReferenceId())
	// XMLDSIG : 4.4.3.2
	// URI=""
	// Identifies the node-set (minus any comment nodes) of the XML resource
	// containing the signature
	dssReference.SetUri("")
	dssReference.SetContents(document)
	dssReference.SetDigestMethodAlgorithm(b.digestAlgorithm)

	dssTransformList := make([]DSSTransform, 0, 2)

	// For parallel signatures
	xPathTransform := NewXPath2FilterEnvelopedSignatureTransformWithNamespace(
		b.signatureParameters.XmldsigNamespace())
	dssTransformList = append(dssTransformList, xPathTransform)

	// Canonicalization is the last operation, it is better to operate the canonicalization on
	// the smaller document
	canonicalizationTransform := NewCanonicalizationTransformWithNamespace(
		b.signatureParameters.XmldsigNamespace(), xmlutils.XMLCanonicalizerDefaultDSSC14NMethod)
	dssTransformList = append(dssTransformList, canonicalizationTransform)

	dssReference.SetTransforms(dssTransformList)

	return dssReference, nil
}

// assertEnvelopedSignaturePossible ports the private
// assertEnvelopedSignaturePossible(DSSDocument).
func (b *ReferenceBuilder) assertEnvelopedSignaturePossible(document model.DSSDocument) error {
	if !xmlutils.DomUtilsIsDOM(document) {
		return exception.NewIllegalInputException(
			"Enveloped signature cannot be created. Reason : the provided document is not XML!")
	}
	dom, err := xmlutils.DomUtilsBuildDOMFromDocument(document)
	if err != nil {
		return err
	}
	documentElement := dom.DocumentElement()
	if common.XMLDSigElementSignature.IsSameTagName(documentElement.Name.Local) {
		return exception.NewIllegalInputException(
			"Unable to create an enveloped signature for another XML signature document!")
	}
	return nil
}

// assertEnvelopingSignatureWithEmbeddedXMLPossible ports the private
// assertEnvelopingSignatureWithEmbeddedXMLPossible(DSSDocument).
func (b *ReferenceBuilder) assertEnvelopingSignatureWithEmbeddedXMLPossible(document model.DSSDocument) error {
	if !xmlutils.DomUtilsIsDOM(document) {
		return exception.NewIllegalInputException("Enveloping signature with embedded XML cannot be created. " +
			"Reason : the provided document is not XML!")
	}
	return nil
}

// envelopingDSSReference ports the private envelopingDSSReference(DSSDocument).
func (b *ReferenceBuilder) envelopingDSSReference(document model.DSSDocument) (*DSSReference, error) {
	// <ds:Reference Id="signed-data-ref" Type="http://www.w3.org/2000/09/xmldsig#Object"
	// URI="#signed-data-idfc5ff27ee49763d9ba88ba5bbc49f732">
	refId := b.referenceIdProvider.ReferenceId()
	reference := NewDSSReference()
	reference.SetId(refId)
	reference.SetContents(document)
	reference.SetDigestMethodAlgorithm(b.digestAlgorithm)

	switch {
	case b.signatureParameters.IsManifestSignature():
		if err := b.assertEnvelopingSignatureWithEmbeddedXMLPossible(document); err != nil {
			return nil, err
		}

		manifestDoc, err := xmlutils.DomUtilsBuildDOMFromDocument(document)
		if err != nil {
			return nil, err
		}
		manifestElement := manifestDoc.DocumentElement()
		if err := b.assertXmlManifestSignaturePossible(manifestElement); err != nil {
			return nil, err
		}

		reference.SetType(common.XMLDSigPathManifestType)
		reference.SetUri(xmlutils.DomUtilsToElementReference(
			manifestElement.AttrValue("", common.XMLDSigAttributeID.AttributeName())))
		var xmlTransform DSSTransform = NewCanonicalizationTransformWithNamespace(
			b.signatureParameters.XmldsigNamespace(), xmlutils.XMLCanonicalizerDefaultDSSC14NMethod)
		reference.SetTransforms([]DSSTransform{xmlTransform})

	case b.signatureParameters.IsEmbedXML():
		if err := b.assertEnvelopingSignatureWithEmbeddedXMLPossible(document); err != nil {
			return nil, err
		}

		reference.SetType(common.XMLDSigPathObjectType)
		reference.SetUri(xmlutils.DomUtilsToElementReference(referenceBuilderObjectIdPrefix + refId))

		var xmlTransform DSSTransform = NewCanonicalizationTransformWithNamespace(
			b.signatureParameters.XmldsigNamespace(), xmlutils.XMLCanonicalizerDefaultDSSC14NMethod)
		reference.SetTransforms([]DSSTransform{xmlTransform})

	default:
		reference.SetType(common.XMLDSigPathObjectType)
		reference.SetUri(xmlutils.DomUtilsToElementReference(referenceBuilderObjectIdPrefix + refId))

		var base64Transform DSSTransform = NewBase64TransformWithNamespace(
			b.signatureParameters.XmldsigNamespace())
		reference.SetTransforms([]DSSTransform{base64Transform})
	}
	return reference, nil
}

// assertXmlManifestSignaturePossible ports the private
// assertXmlManifestSignaturePossible(Element).
func (b *ReferenceBuilder) assertXmlManifestSignaturePossible(manifestElement *xmldom.Node) error {
	idAttr := manifestElement.AttrValue("", common.XMLDSigAttributeID.AttributeName())
	if utils.IsStringBlank(idAttr) {
		return exception.NewIllegalInputException(
			"Manifest signature is not possible for an XML file without Id attribute in the root element!")
	}
	return nil
}

// detachedDSSReference ports the private detachedDSSReference(DSSDocument).
func (b *ReferenceBuilder) detachedDSSReference(document model.DSSDocument) *DSSReference {
	reference := NewDSSReference()
	reference.SetId(b.referenceIdProvider.ReferenceId())
	if utils.IsStringNotEmpty(document.Name()) {
		reference.SetUri(document.Name())
	}
	reference.SetContents(document)
	reference.SetDigestMethodAlgorithm(b.digestAlgorithm)
	return reference
}

// internallyDetachedDSSReference ports the private
// internallyDetachedDSSReference(DSSDocument).
func (b *ReferenceBuilder) internallyDetachedDSSReference(document model.DSSDocument) (*DSSReference, error) {
	reference := NewDSSReference()
	reference.SetId(b.referenceIdProvider.ReferenceId())

	dom, err := xmlutils.DomUtilsBuildDOMFromDocument(document)
	if err != nil {
		return nil, err
	}
	identifier := DSSXMLUtilsGetIDIdentifier(dom.DocumentElement())
	if identifier == "" {
		panic("ID not defined on the root xml element")
	}
	reference.SetUri(xmlutils.DomUtilsToElementReference(identifier))

	reference.SetContents(document)
	reference.SetDigestMethodAlgorithm(b.digestAlgorithm)

	dssTransformList := make([]DSSTransform, 0, 1)
	canonicalization := NewCanonicalizationTransformWithNamespace(
		b.signatureParameters.XmldsigNamespace(), xmlutils.XMLCanonicalizerDefaultDSSC14NMethod)
	dssTransformList = append(dssTransformList, canonicalization)
	reference.SetTransforms(dssTransformList)
	return reference, nil
}
