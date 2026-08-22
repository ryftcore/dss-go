// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/ReferenceProcessor.java (DSS 6.5.RC1).
//
// # Errors
//
// getReferenceOutput and incorporateReferences gain an error return: DomUtils.buildDOM and
// DomUtils.excludeComments, which Java performs behind DSSException-free APIs, can fail in Go,
// DSSXMLUtils.applyTransforms throws IllegalInputException/DSSException, and the digest
// computation can fail. slf4j trace logging of the reference output is dropped.
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// ReferenceProcessor contains methods for processing a DSSReference.
type ReferenceProcessor struct {
	// signatureParameters are the signature parameters used on signature creation. Nil for
	// non-signature references, e.g. for a Manifest.
	signatureParameters *SignatureParameters
}

// NewReferenceProcessorEmpty ports the empty constructor, to be used for non-signature
// references (e.g. for a Manifest).
func NewReferenceProcessorEmpty() *ReferenceProcessor {
	return &ReferenceProcessor{}
}

// NewReferenceProcessor is the constructor to be used for reference processing on signature
// creation. Ports ReferenceProcessor(XAdESSignatureParameters).
func NewReferenceProcessor(signatureParameters *SignatureParameters) *ReferenceProcessor {
	return &ReferenceProcessor{signatureParameters: signatureParameters}
}

// ReferenceOutput returns an output content after processing the given DSSReference.
// Ports getReferenceOutput(DSSReference).
func (p *ReferenceProcessor) ReferenceOutput(reference *DSSReference) (model.DSSDocument, error) {
	if _, ok := reference.Contents().(*model.DigestDocument); ok {
		return reference.Contents(), nil
	}

	nodeToTransform, err := p.dereferenceNode(reference)
	if err != nil {
		return nil, err
	}
	if nodeToTransform == nil {
		return reference.Contents(), nil
	}
	transforms := reference.Transforms()
	if p.isUniqueBase64Transform(transforms) {
		return reference.Contents(), nil
	}

	referenceOutputResult, err := DSSXMLUtilsApplyTransforms(nodeToTransform, reference.Transforms())
	if err != nil {
		return nil, err
	}
	// NodeSet canonicalization is performed by internal/xmldsig within ApplyTransforms

	return model.NewInMemoryDocument(referenceOutputResult), nil
}

// dereferenceNode ports the private dereferenceNode(DSSReference).
func (p *ReferenceProcessor) dereferenceNode(reference *DSSReference) (*xmldom.Node, error) {
	document, err := p.documentToTransform(reference)
	if err != nil {
		return nil, err
	}
	if document == nil {
		// not XML, return NULL node
		return nil, nil
	}

	/*
	 * 4.4.3.3 Same-Document URI-References
	 *
	 * The application must behave as if the result of XPointer processing [XPTR-FRAMEWORK] were
	 * a node-set derived from the resultant subresource as follows:
	 * 1. include XPath nodes having full or partial content within the subresource
	 * 2. replace the root node with its children (if it is in the node-set)
	 * 3. replace any element node E with E plus all descendants of E (text, comment, PI, element)
	 *    and all namespace and attribute nodes of E and its descendant elements.
	 * 4. if the URI has no fragment identifier or the fragment identifier is a shortname XPointer,
	 *    then delete all comment nodes
	 */
	// reference.HasUri() carries Java's "the URI is not null": isSameDocumentReference answers
	// true for the empty URI and false for a null one, which a bare "" could not tell apart.
	if reference.HasUri() && DSSXMLUtilsIsSameDocumentReference(reference.Uri()) &&
		!xmlutils.DomUtilsIsXPointerQuery(reference.Uri()) {
		document, err = xmlutils.DomUtilsExcludeComments(document)
		if err != nil {
			return nil, err
		}
	}
	// extract the target Node after performing comments removal (otherwise the relation to
	// parent nodes can be lost)
	return p.nodeToTransform(document, reference), nil
}

// documentToTransform ports the private getDocumentToTransform(DSSReference).
func (p *ReferenceProcessor) documentToTransform(reference *DSSReference) (*xmldom.Node, error) {
	contents := reference.Contents()
	if !xmlutils.DomUtilsIsDOM(contents) {
		// cannot be transformed
		return nil, nil
	}
	return xmlutils.DomUtilsBuildDOMFromDocument(contents)
}

// nodeToTransform ports the private getNodeToTransform(Document, DSSReference).
func (p *ReferenceProcessor) nodeToTransform(doc *xmldom.Node, reference *DSSReference) *xmldom.Node {
	uri := reference.Uri()
	if p.signatureParameters != nil && p.signatureParameters.IsEmbedXML() {
		doc2 := xmlutils.DomUtilsBuildDOMEmpty()
		dom := xmlutils.DomUtilsCreateElementNS(doc2, p.signatureParameters.XmldsigNamespace(),
			common.XMLDSigElementObject)
		dom2 := xmlutils.DomUtilsCreateElementNS(doc2, p.signatureParameters.XmldsigNamespace(),
			common.XMLDSigElementObject)
		doc2.AppendChild(dom2)
		dom2.AppendChild(dom)
		dom.SetAttr(xmldom.Name{Local: common.XMLDSigAttributeID.AttributeName()}, xmlutils.DomUtilsGetId(uri))

		xmlutils.DomUtilsAdoptChildren(dom, doc)
		return dom

	} else if xmlutils.DomUtilsIsElementReference(uri) {
		targetId := xmlutils.DomUtilsGetId(uri)
		elementById := xmlutils.XPathUtilsGetElementById(doc, targetId)
		if elementById != nil {
			return elementById
		}
		// continue for on-fly document creation

	}
	// TODO : add support of xPointer
	if utils.IsCollectionNotEmpty(reference.Transforms()) {
		return doc
	}
	return nil
}

// isUniqueBase64Transform ports the private isUniqueBase64Transform(List).
func (p *ReferenceProcessor) isUniqueBase64Transform(transforms []DSSTransform) bool {
	if len(transforms) != 1 {
		return false
	}
	_, ok := transforms[0].(*Base64Transform)
	return ok
}

// IncorporateReferences incorporates a list of references within the provided
// referenceContainer element. Ports incorporateReferences(Element, List, DSSNamespace).
func (p *ReferenceProcessor) IncorporateReferences(referenceContainer *xmldom.Node,
	references []*DSSReference, namespace *common.DSSNamespace) error {
	if utils.IsCollectionNotEmpty(references) {
		documentDom := referenceContainer.OwnerDocument()
		for _, dssReference := range references {
			referenceDom := xmlutils.DomUtilsCreateElementNS(documentDom, namespace,
				common.XMLDSigElementReference)
			referenceContainer.AppendChild(referenceDom)

			if dssReference.Id() != "" {
				referenceDom.SetAttr(
					xmldom.Name{Local: common.XMLDSigAttributeID.AttributeName()}, dssReference.Id())
			}
			// Java's "if (uri != null)": the enveloped reference's URI is the EMPTY string and
			// upstream does write URI="" for it, so the test is HasUri, not uri != "".
			uri := dssReference.Uri()
			if dssReference.HasUri() {
				if !DSSXMLUtilsIsSameDocumentReference(uri) {
					uri = spi.DSSUtilsEncodeURI(uri)
				}
				referenceDom.SetAttr(
					xmldom.Name{Local: common.XMLDSigAttributeURI.AttributeName()}, uri)
			}
			referenceType := dssReference.Type()
			if referenceType != "" {
				referenceDom.SetAttr(
					xmldom.Name{Local: common.XMLDSigAttributeType.AttributeName()}, referenceType)
			}

			DSSXMLUtilsIncorporateTransforms(referenceDom, dssReference.Transforms(), namespace)
			DSSXMLUtilsIncorporateDigestMethod(referenceDom, dssReference.DigestMethodAlgorithm(), namespace)

			documentAfterTransforms, err := p.ReferenceOutput(dssReference)
			if err != nil {
				return err
			}
			digestBytes, err := documentAfterTransforms.DigestValue(dssReference.DigestMethodAlgorithm())
			if err != nil {
				return err
			}
			base64EncodedDigestBytes := utils.ToBase64(digestBytes)
			DSSXMLUtilsIncorporateDigestValue(referenceDom, base64EncodedDigestBytes, namespace)
		}
	}
	return nil
}
