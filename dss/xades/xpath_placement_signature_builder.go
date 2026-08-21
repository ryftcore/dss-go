// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XPathPlacementSignatureBuilder.java (DSS 6.5.RC1).
//
// Java's class is abstract and sits between XAdESSignatureBuilder and the two packagings that
// place the signature by XPath (Enveloped, InternallyDetached). It overrides two of the hooks
// collected in XAdESSignatureBuilderOverrides - getParentNodeOfSignature and
// incorporateSignatureDom(Node) - and adds assertOriginalXmlDocumentValid for its subclasses. Go
// has no method overriding across embedding, so the concrete subclass embeds this type and
// registers itself with InitXAdESSignatureBuilder; the two overridden methods are promoted to it
// and reached through the base's overrides field, the TokenBase.InitToken(self) convention of
// PORTING.md.
//
// Santuario replacement (see internal/xmldsig/doc.go): the only org.apache.xml.security member
// this file touches is the constant Transforms.TRANSFORM_ENVELOPED_SIGNATURE, which is the
// enveloped-signature transform URI - taken here from the transform registry's own identifier so
// the value is stated once for the whole port.
//
// Java's LinkedHashSet<Node> of ancestor nodes becomes an insertion-ordered slice plus a
// membership test: the set is only ever probed with contains(), and a slice keeps the iteration
// deterministic (PORTING.md; the determinism sweep).
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// XPathPlacementSignatureBuilder creates signatures that are enveloped into the parent document
// based on the defined (if any) XPath location.
type XPathPlacementSignatureBuilder struct {
	XAdESSignatureBuilder
}

// InitXPathPlacementSignatureBuilder registers the concrete builder with this intermediate base
// and with XAdESSignatureBuilder. Port of the two protected XPathPlacementSignatureBuilder
// constructors, which do nothing beyond delegating to super.
func (b *XPathPlacementSignatureBuilder) InitXPathPlacementSignatureBuilder(
	self XAdESSignatureBuilderOverrides, params *XAdESSignatureParameters,
	documents []model.DSSDocument, certificateVerifier validation.CertificateVerifier) {
	b.InitXAdESSignatureBuilder(self, params, documents, certificateVerifier)
}

// AssertOriginalXmlDocumentValid verifies the conformance of the original document for an
// enveloped signature creation. Port of the protected #assertOriginalXmlDocumentValid.
func (b *XPathPlacementSignatureBuilder) AssertOriginalXmlDocumentValid() error {
	if utils.CollectionSize(b.Documents) > 1 {
		return fmt.Errorf("Only one original document is allowed for '%s' signature packaging!",
			b.Params.SignaturePackaging())
	}
	if !xmlutils.DomUtilsIsDOM(b.Documents[0]) {
		return exception.NewIllegalInputException(
			"Enveloped signature cannot be created. Reason : the provided document is not XML!")
	}

	b.initRootDocumentDom()

	signatureNodeList, err := DSSXMLUtilsGetAllSignaturesExceptCounterSignatures(b.DocumentDom)
	if err != nil {
		return err
	}
	if len(signatureNodeList) == 0 {
		return nil
	}

	parentSignatureNode := b.overrides.ParentNodeOfSignature()
	parentNodes := xpathPlacementSignatureBuilderParentNodesChain(parentSignatureNode)

	for _, signatureNode := range signatureNodeList {
		referenceNodeList, err := DSSXMLUtilsGetReferenceNodeList(signatureNode)
		if err != nil {
			return err
		}
		if len(referenceNodeList) == 0 {
			continue
		}

		for _, referenceNode := range referenceNodeList {
			affected, err := b.isSignatureCoveredNodeAffected(referenceNode, parentNodes)
			if err != nil {
				return err
			}
			if affected {
				if err := xpathPlacementSignatureBuilderAssertDoesNotContainEnvelopedTransform(referenceNode); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// xpathPlacementSignatureBuilderParentNodesChain ports the private getParentNodesChain; Java's
// LinkedHashSet becomes an insertion-ordered slice (the set is only probed for membership).
func xpathPlacementSignatureBuilderParentNodesChain(node *xmldom.Node) []*xmldom.Node {
	nodesChain := []*xmldom.Node{node}
	for parentNode := node.Parent; parentNode != nil; parentNode = parentNode.Parent {
		nodesChain = append(nodesChain, parentNode)
	}
	return nodesChain
}

// isSignatureCoveredNodeAffected ports the private isSignatureCoveredNodeAffected.
func (b *XPathPlacementSignatureBuilder) isSignatureCoveredNodeAffected(referenceNode *xmldom.Node,
	affectedNodes []*xmldom.Node) (bool, error) {
	// Java reads DSSXMLUtils.getAttribute(node, name), i.e. getAttributes().getNamedItem(name),
	// and distinguishes its null (attribute absent -> false) from "" (URI="" covers the whole
	// document -> true). A Go string cannot carry that distinction, so the attribute node itself
	// is read here; the value is the same one DSSXMLUtilsGetAttribute would return.
	uriAttr := referenceNode.Attr("", common.XMLDSigAttribute_URI.AttributeName())
	if uriAttr == nil {
		return false, nil
	}
	id := uriAttr.Value
	if utils.IsStringEmpty(id) {
		// covers the whole file
		return true, nil
	}
	referencedNode := xmlutils.XPathUtilsGetElementById(b.DocumentDom, id)
	for _, affectedNode := range affectedNodes {
		if affectedNode == referencedNode {
			return true, nil
		}
	}
	return false, nil
}

// xpathPlacementSignatureBuilderAssertDoesNotContainEnvelopedTransform ports the private
// assertDoesNotContainEnvelopedTransform.
func xpathPlacementSignatureBuilderAssertDoesNotContainEnvelopedTransform(referenceNode *xmldom.Node) error {
	transformList, err := xmlutils.XPathUtilsGetNodeList(referenceNode,
		common.XMLDSigPath_TRANSFORMS_TRANSFORM_PATH)
	if err != nil {
		return err
	}
	for _, transformElement := range transformList {
		transformAlgorithm := transformElement.AttrValue("", common.XMLDSigAttribute_ALGORITHM.AttributeName())
		if string(xmldsig.TransformEnvelopedSignature) == transformAlgorithm {
			return exception.NewIllegalInputException(fmt.Sprintf(
				"The parallel signature is not possible! The provided file contains a signature with an '%s' transform.",
				xmldsig.TransformEnvelopedSignature))
		}
	}
	return nil
}

// ParentNodeOfSignature returns the node the ds:Signature element is placed under: the single
// element the configured XPath location selects, or the document element when no location is
// configured. Port of the overridden protected #getParentNodeOfSignature.
func (b *XPathPlacementSignatureBuilder) ParentNodeOfSignature() *xmldom.Node {
	xPathLocationString := b.Params.XPathLocationString()
	if utils.IsStringNotEmpty(xPathLocationString) {
		nodeList, err := xmlutils.NewXPathQueryExecutorLoader().GetXPathStringExecutor().
			GetNodeListByString(b.DocumentDom, xPathLocationString)
		if err == nil && len(nodeList) == 1 {
			return nodeList[0]
		}
		// Java throws IllegalArgumentException here; this hook has no error channel (its Java
		// signature has none either, and the DetachedSignatureBuilder override matches), so the
		// unchecked exception becomes a panic carrying the same message - PORTING.md's rule for
		// an unchecked exception that has nowhere to be returned.
		panic(fmt.Sprintf("Unable to find an element corresponding to XPath location '%s'",
			xPathLocationString))
	}
	return b.DocumentDom.DocumentElement()
}

// IncorporateSignatureDomToParent places the ds:Signature element relative to the node the XPath
// location selected, honouring the configured XPathElementPlacement.
// Port of the overridden protected #incorporateSignatureDom(Node).
func (b *XPathPlacementSignatureBuilder) IncorporateSignatureDomToParent(parentNodeOfSignature *xmldom.Node) {
	if b.Params.XPathElementPlacement() == "" || utils.IsStringEmpty(b.Params.XPathLocationString()) {
		b.XAdESSignatureBuilder.IncorporateSignatureDomToParent(parentNodeOfSignature)
		return
	}

	switch b.Params.XPathElementPlacement() {
	case XPathElementPlacement_XPathAfter:
		// root element referenced by XPath
		// DEVIATION: Java compares with Node#isEqualNode, a deep structural comparison. The two
		// operands always come from the same document here, and no descendant of the document
		// element can be structurally equal to it (it contains itself), so pointer identity
		// decides the same branch while avoiding a whole-subtree walk.
		if parentNodeOfSignature == b.DocumentDom.DocumentElement() {
			// append signature at end of document
			parentNodeOfSignature.AppendChild(b.SignatureDom)

		} else {
			// insert signature before next sibling or as last child if no sibling exists
			parent := parentNodeOfSignature.Parent
			parent.InsertBefore(b.SignatureDom, parentNodeOfSignature.NextSibling)
		}

	case XPathElementPlacement_XPathFirstChildOf:
		parentNodeOfSignature.InsertBefore(b.SignatureDom, parentNodeOfSignature.FirstChild)

	default:
		parentNodeOfSignature.AppendChild(b.SignatureDom)
	}
}
