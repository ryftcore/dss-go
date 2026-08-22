// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/XsltTransform.java (DSS 6.5.RC1).
//
// org.apache.xml.security.transforms.Transforms.TRANSFORM_XSLT is xmldsig.TransformXSLT, per
// internal/xmldsig's doc.go mapping table.
//
// # Accepted divergence: executing the transform
//
// WRITING this transform into a ds:Transforms element is byte-identical to upstream - the KAT
// case "tf-xslt" pins it. EXECUTING it is not: internal/xmldsig registers the XSLT algorithm as
// a transform that always refuses (see its DefaultRegistry doc comment - running an XSLT
// stylesheet that arrives inside a signature is not a behaviour worth reproducing, and no XAdES
// profile uses the transform), so PerformTransform returns xmldsig.ErrForbiddenTransform where
// upstream returns the stylesheet's output.
//
// The refusal is NOT what upstream does, and the earlier reading of this file said otherwise:
// DSS calls Transform#performTransform(input, true), but that secureValidation flag only
// hardens Santuario's TransformerFactory (no external DTD, no external stylesheet) - the
// outright ForbiddenTransform refusal lives in Transforms#checkSecureValidation, on the
// validation-side chain runner, which ComplexTransform never goes through. The Java oracle
// confirms it: gen/RefsOracle.java's "apply-xslt" case records upstream transforming the
// fixture into "<?xml version=\"1.0\" encoding=\"UTF-8\"?><out/>". The KAT asserts the refusal
// deliberately, naming this comment; the decision belongs to internal/xmldsig, which is frozen.
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/xml/common"
)

// XsltTransform is the XSLT transform.
type XsltTransform struct {
	ComplexTransform

	// content is the document to be added.
	content *xmldom.Node
}

// NewXsltTransform ports XsltTransform(Document), which delegates with XMLDSigNamespace.NS.
func NewXsltTransform(content *xmldom.Node) *XsltTransform {
	return NewXsltTransformWithNamespace(common.XMLDSigNS, content)
}

// NewXsltTransformWithNamespace ports XsltTransform(DSSNamespace, Document). Java's
// Objects.requireNonNull(content, "The content cannot be null!") becomes a panic carrying the
// same message (PORTING.md).
func NewXsltTransformWithNamespace(xmlDSigNamespace *common.DSSNamespace, content *xmldom.Node) *XsltTransform {
	if content == nil {
		panic("The content cannot be null!")
	}
	t := &XsltTransform{
		ComplexTransform: newComplexTransform(xmlDSigNamespace, xmldsig.TransformXSLT),
		content:          content,
	}
	t.InitAbstractTransform(t)
	return t
}

// CreateTransform appends the ds:Transform element and moves a deep copy of the stylesheet's
// document element into it. Ports createTransform(Document, Element): cloneNode(true),
// getDocumentElement(), adoptNode() and appendChild(), where adoptNode's detach-and-reown is
// the RemoveChild below (the node is already a copy, so nothing shared is disturbed).
func (t *XsltTransform) CreateTransform(document, parentNode *xmldom.Node) *xmldom.Node {
	transform := t.AbstractTransform.CreateTransform(document, parentNode)
	clonedNode := t.content.Clone(true)
	contextDocumentElement := clonedNode.DocumentElement()
	clonedNode.RemoveChild(contextDocumentElement)
	return transform.AppendChild(contextDocumentElement)
}
