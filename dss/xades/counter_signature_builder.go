// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/CounterSignatureBuilder.java (DSS 6.5.RC1).
//
// Java extends ExtensionBuilder and overrides nothing; it registers itself with
// InitExtensionBuilderWithVerifier so the three XAdESBuilder hooks resolve through this concrete
// type, the TokenBase.InitToken(self) convention of PORTING.md.
//
// # javax.xml.crypto.dsig.XMLSignature.XMLNS and getElementsByTagNameNS
//
// incorporateCounterSignature looks the counter signature up with
// getElementsByTagNameNS(XMLNS, "Signature") - XMLNS being the XMLDSig namespace URI, which
// common.XMLDSigElementSignature.URI() already carries, so the JSR-105 constant does not need a
// separate home. xmldom has no getElementsByTagNameNS, so the document-order preorder walk that
// method performs is spelled out in a file-local helper; the only property the caller depends on
// is that walk order and the count.
//
// # Deterministic Id priming (DEVIATION, flagged for the integrator)
//
// Java's XAdESCounterSignatureParameters overrides getDeterministicId(), and every call reaching
// it through an XAdESSignatureParameters reference dispatches virtually to the counter-signature
// value. Go has no such dispatch: a *XAdESSignatureParameters holding the embedded base of an
// XAdESCounterSignatureParameters answers the base implementation. Both implementations cache
// their result in the shared XAdESProfileParameters, so each public method here that receives the
// concrete *XAdESCounterSignatureParameters calls its GetDeterministicId() once up front; every
// later base-typed read - inside ReferenceIdProvider, XAdESSignatureBuilder and this file - then
// returns that same cached counter-signature Id, exactly as Java's virtual call would. The value
// is identical to Java's and computed at the same point in the flow, since Java's first virtual
// call also happens inside these methods.
//
// Errors: IllegalInputException and IllegalArgumentException become returned errors, and
// Objects.requireNonNull keeps panicking with the Java message (PORTING.md). slf4j logging is
// dropped.
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// counterSignatureBuilderCounterSignaturePrefix is the default prefix for a counter signature Id.
// Port of the private COUNTER_SIGNATURE_PREFIX.
const counterSignatureBuilderCounterSignaturePrefix = "CS-"

// CounterSignatureBuilder creates a counter signature signing the ds:SignatureValue of a parent
// signature.
type CounterSignatureBuilder struct {
	ExtensionBuilder
}

// NewCounterSignatureBuilder is the default constructor.
// Port of the protected CounterSignatureBuilder(CertificateVerifier) constructor; the type is
// exported (see detached_signature_builder.go for the package-private deviation note), and so is
// its constructor, because XAdESService builds one.
func NewCounterSignatureBuilder(certificateVerifier validation.CertificateVerifier) *CounterSignatureBuilder {
	builder := &CounterSignatureBuilder{}
	builder.InitExtensionBuilderWithVerifier(builder, certificateVerifier)
	return builder
}

// GetCanonicalizedSignatureValue extracts a canonicalized ds:SignatureValue element from the
// provided XAdES signature. Port of #getCanonicalizedSignatureValue.
func (b *CounterSignatureBuilder) GetCanonicalizedSignatureValue(signatureDocument model.DSSDocument,
	parameters *XAdESCounterSignatureParameters) (model.DSSDocument, error) {
	// See the file header: prime the counter-signature deterministic Id in the shared context.
	parameters.GetDeterministicId()
	b.Params = &parameters.XAdESSignatureParameters

	documentAnalyzer, err := NewXMLDocumentAnalyzer(signatureDocument)
	if err != nil {
		return nil, err
	}
	b.DocumentAnalyzer = documentAnalyzer
	b.DocumentDom = b.DocumentAnalyzer.RootElement()

	xadesSignature, err := b.extractSignatureById(parameters)
	if err != nil {
		return nil, err
	}
	b.XadesSignature = xadesSignature
	if err := b.AssertUnsignedPropertiesExtensionPossible(b.XadesSignature); err != nil {
		return nil, err
	}

	signatureValueElement, err := counterSignatureBuilderSignatureValueElement(b.XadesSignature)
	if err != nil {
		return nil, err
	}
	canonicalizer, err := xmlutils.XMLCanonicalizerCreateInstanceWithMethod(
		parameters.CounterSignatureCanonicalizationMethod())
	if err != nil {
		return nil, err
	}
	canonicalizedSignatureValue, err := canonicalizer.CanonicalizeNode(signatureValueElement)
	if err != nil {
		return nil, err
	}

	// Upstream traces the canonicalized SignatureValue here.

	return model.NewInMemoryDocument(canonicalizedSignatureValue), nil
}

// BuildCounterSignatureDSSReference builds a DSSReference for the ds:SignatureValue to
// counter sign. Port of #buildCounterSignatureDSSReference.
func (b *CounterSignatureBuilder) BuildCounterSignatureDSSReference(signatureDocument model.DSSDocument,
	parameters *XAdESCounterSignatureParameters) (*DSSReference, error) {
	// See the file header: prime the counter-signature deterministic Id in the shared context.
	parameters.GetDeterministicId()

	documentAnalyzer, err := NewXMLDocumentAnalyzer(signatureDocument)
	if err != nil {
		return nil, err
	}
	b.DocumentAnalyzer = documentAnalyzer
	b.DocumentDom = b.DocumentAnalyzer.RootElement()

	xadesSignature, err := b.extractSignatureById(parameters)
	if err != nil {
		return nil, err
	}
	b.XadesSignature = xadesSignature
	if _, err := b.InitializeSignatureBuilder(b.XadesSignature); err != nil {
		return nil, err
	}

	reference := NewDSSReference()

	referenceIdProvider := NewReferenceIdProvider()
	referenceIdProvider.SetSignatureParameters(&parameters.XAdESSignatureParameters)
	reference.SetId(referenceIdProvider.ReferenceId())

	signatureElementBinaries, err := xmlutils.DomUtilsSerializeNode(b.XadesSignature.SignatureElement())
	if err != nil {
		return nil, err
	}
	reference.SetContents(model.NewInMemoryDocument(signatureElementBinaries))
	reference.SetDigestMethodAlgorithm(
		DSSXMLUtilsGetReferenceDigestAlgorithmOrDefault(&parameters.XAdESSignatureParameters))
	reference.SetType(b.XadesPath.CounterSignatureUri())

	signatureValueId := b.XadesSignature.SignatureValueId()
	if utils.IsStringNotEmpty(signatureValueId) {
		reference.SetUri(xmlutils.DomUtilsToElementReference(signatureValueId))
		var transform DSSTransform = NewCanonicalizationTransform(
			parameters.CounterSignatureCanonicalizationMethod())
		reference.SetTransforms([]DSSTransform{transform})

	} else {
		return nil, exception.NewIllegalInputException(fmt.Sprintf(
			"The signature with Id '%s' does not have an Id for a SignatureValue element! "+
				"Unable to counter sign.", parameters.SignatureIdToCounterSign()))
	}
	return reference, nil
}

// BuildEmbeddedCounterSignature embeds the counter signature into the UnsignedProperties of the
// original XAdES signature and returns the resulting document.
// Port of #buildEmbeddedCounterSignature.
func (b *CounterSignatureBuilder) BuildEmbeddedCounterSignature(signatureDocument model.DSSDocument,
	counterSignature model.DSSDocument,
	parameters *XAdESCounterSignatureParameters) (model.DSSDocument, error) {
	// See the file header: prime the counter-signature deterministic Id in the shared context.
	parameters.GetDeterministicId()
	b.Params = &parameters.XAdESSignatureParameters

	documentAnalyzer, err := NewXMLDocumentAnalyzer(signatureDocument)
	if err != nil {
		return nil, err
	}
	b.DocumentAnalyzer = documentAnalyzer
	b.DocumentDom = b.DocumentAnalyzer.RootElement()

	xadesSignature, err := b.extractSignatureById(parameters)
	if err != nil {
		return nil, err
	}
	b.XadesSignature = xadesSignature
	if err := b.AssertUnsignedPropertiesExtensionPossible(b.XadesSignature); err != nil {
		return nil, err
	}

	if _, err := b.InitializeSignatureBuilder(b.XadesSignature); err != nil {
		return nil, err
	}

	levelBUnsignedProperties := b.UnsignedSignaturePropertiesDom.Clone(true)

	if err := b.incorporateCounterSignature(counterSignature); err != nil {
		return nil, err
	}

	indented, err := b.IndentIfPrettyPrint(b.UnsignedSignaturePropertiesDom, levelBUnsignedProperties)
	if err != nil {
		return nil, err
	}
	b.UnsignedSignaturePropertiesDom = indented

	return b.CreateXmlDocument()
}

// incorporateCounterSignature ports the private incorporateCounterSignature.
func (b *CounterSignatureBuilder) incorporateCounterSignature(counterSignature model.DSSDocument) error {
	counterSignatureDom, err := xmlutils.DomUtilsBuildDOMFromDocument(counterSignature)
	if err != nil {
		return err
	}

	signatureNodeList := counterSignatureBuilderElementsByTagNameNS(counterSignatureDom,
		common.XMLDSigElementSignature.URI(), common.XMLDSigElementSignature.TagName())
	if len(signatureNodeList) != 1 {
		return exception.NewIllegalInputException(fmt.Sprintf(
			"The counterSignature document shall have one counter signature, when %d signatures found!",
			len(signatureNodeList)))
	}

	signatureElement := signatureNodeList[0]
	adopted := b.DocumentDom.Import(signatureElement, true)

	currentElements, err := b.CurrentXAdESElements()
	if err != nil {
		return err
	}
	counterSignatureElement := xmlutils.DomUtilsAddElement(b.DocumentDom, b.UnsignedSignaturePropertiesDom,
		b.overrides.XadesNamespace(), currentElements.ElementCounterSignature())
	counterSignatureElement.SetAttr(
		xmldom.Name{Local: common.XMLDSigAttributeID.AttributeName()},
		counterSignatureBuilderCounterSignaturePrefix+b.Params.GetDeterministicId())
	counterSignatureElement.AppendChild(adopted)
	return nil
}

// counterSignatureBuilderElementsByTagNameNS reproduces org.w3c.dom's
// Document#getElementsByTagNameNS: every element of the subtree whose namespace URI and local
// name match, in document order (preorder).
func counterSignatureBuilderElementsByTagNameNS(root *xmldom.Node, uri, localName string) []*xmldom.Node {
	var matches []*xmldom.Node
	root.Walk(func(node *xmldom.Node) bool {
		if node.Kind == xmldom.Element && node.Name.Space == uri && node.Name.Local == localName {
			matches = append(matches, node)
		}
		return true
	})
	return matches
}

// extractSignatureById ports the private extractSignatureById.
func (b *CounterSignatureBuilder) extractSignatureById(
	parameters *XAdESCounterSignatureParameters) (*XAdESSignature, error) {
	if parameters.SignatureIdToCounterSign() == "" {
		panic("The Id of a signature to be counter signed shall be defined! " +
			"Please use SerializableCounterSignatureParameters.setSignatureIdToCounterSign(signatureId) method.")
	}

	signatures := b.DocumentAnalyzer.Signatures()
	for _, signature := range signatures {
		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			return nil, fmt.Errorf("unexpected signature type %T", signature)
		}
		signatureById, err := counterSignatureBuilderSignatureOrItsCounterSignatureById(xadesSignature,
			parameters.SignatureIdToCounterSign())
		if err != nil {
			return nil, err
		}
		if signatureById != nil {
			return signatureById, nil
		}
	}

	return nil, fmt.Errorf("A signature with Id '%s' has not been found in the file! Unable to counter sign.",
		parameters.SignatureIdToCounterSign())
}

// counterSignatureBuilderSignatureOrItsCounterSignatureById ports the private
// getSignatureOrItsCounterSignatureById.
func counterSignatureBuilderSignatureOrItsCounterSignatureById(signature *XAdESSignature,
	signatureId string) (*XAdESSignature, error) {
	if signatureId == signature.DAIdentifier() || signatureId == signature.ID() {
		return signature, nil
	}

	for _, counterSignature := range signature.CounterSignatures() {
		xadesCounterSignature, ok := counterSignature.(*XAdESSignature)
		if !ok {
			return nil, fmt.Errorf("unexpected signature type %T", counterSignature)
		}
		counterSignatureById, err := counterSignatureBuilderSignatureOrItsCounterSignatureById(
			xadesCounterSignature, signatureId)
		if err != nil {
			return nil, err
		}
		if counterSignatureById != nil {
			// check if not timestamped
			if signature.TimestampSource().IsTimestamped(signatureId,
				enumerations.TimestampedObjectTypeSignature) {
				return nil, exception.NewIllegalInputException(fmt.Sprintf(
					"Unable to counter sign a signature with Id '%s'. "+
						"The signature is timestamped by a master signature!", signatureId))
			}
			return counterSignatureById, nil
		}
	}

	return nil, nil
}

// counterSignatureBuilderSignatureValueElement ports the private getSignatureValueElement.
func counterSignatureBuilderSignatureValueElement(xadesSignature *XAdESSignature) (*xmldom.Node, error) {
	signatureElement := xadesSignature.SignatureElement()

	signatureValueElement, err := xmlutils.XPathUtilsGetElement(signatureElement,
		common.XMLDSigPathSignatureValuePath)
	if err != nil {
		return nil, err
	}
	if signatureValueElement != nil {
		return signatureValueElement, nil
	}

	return nil, exception.NewIllegalInputException(fmt.Sprintf(
		"Unable to counter sign a signature with Id '%s'. The SignatureValue element is not found!",
		xadesSignature.DAIdentifier()))
}
