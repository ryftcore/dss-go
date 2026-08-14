// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/XPath2FilterTransform.java (DSS 6.5.RC1).
//
// org.apache.xml.security.transforms.Transforms.TRANSFORM_XPATH2FILTER is
// xmldsig.TransformXPath2Filter, per internal/xmldsig's doc.go mapping table; the transform
// itself is executed by that package's registry through ComplexTransform.
package xades

import (
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/internal/xmldsig"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/xades/definition"
	"github.com/utain/esig/dss/xml/common"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// xPath2FilterTransformFilterAttribute is the filter attribute name. Port of the private
// FILTER_ATTRIBUTE constant.
const xPath2FilterTransformFilterAttribute = "Filter"

// XPath2FilterTransform represents a Filter 2.0 transform.
type XPath2FilterTransform struct {
	XPathTransform

	// filter is the filter.
	filter string
}

// NewXPath2FilterTransform ports XPath2FilterTransform(String, String), which delegates with
// XMLDSigNamespace.NS.
func NewXPath2FilterTransform(xPathExpression, filter string) *XPath2FilterTransform {
	return NewXPath2FilterTransformWithNamespace(common.XMLDSigNS, xPathExpression, filter)
}

// NewXPath2FilterTransformWithNamespace ports
// XPath2FilterTransform(DSSNamespace, String, String).
func NewXPath2FilterTransformWithNamespace(xmlDSigNamespace *common.DSSNamespace,
	xPathExpression, filter string) *XPath2FilterTransform {
	t := new(XPath2FilterTransform)
	*t = newXPath2FilterTransform(xmlDSigNamespace, xPathExpression, filter)
	t.InitAbstractTransform(t)
	return t
}

// newXPath2FilterTransform builds the embeddable value the constructors and the
// XPath2FilterEnvelopedSignatureTransform subclass share; the caller must still call
// InitAbstractTransform with itself.
//
// Java's Objects.requireNonNull(filter) has no Go analogue: a Go string is never nil.
func newXPath2FilterTransform(xmlDSigNamespace *common.DSSNamespace,
	xPathExpression, filter string) XPath2FilterTransform {
	return XPath2FilterTransform{
		XPathTransform: newXPathTransform(xmlDSigNamespace, xmldsig.TransformXPath2Filter, xPathExpression),
		filter:         filter,
	}
}

// xPath2FilterTransform returns the embedded XPath2FilterTransform. Promoted through embedding,
// it is how a subclass's Equals reaches the field Java's XPath2FilterTransform#equals compares.
func (t *XPath2FilterTransform) xPath2FilterTransform() *XPath2FilterTransform { return t }

// CreateTransform appends the ds:Transform element plus the XPath element, which - unlike the
// plain XPath transform's - lives in the XMLDSIG Filter 2.0 namespace and re-declares both that
// namespace and the xmldsig one on itself. Ports createTransform(Document, Element).
//
// Note that this override does NOT delegate to AbstractTransform#createTransform: it builds the
// ds:Transform element itself, exactly as upstream, and returns the XPath element.
func (t *XPath2FilterTransform) CreateTransform(document, parentNode *xmldom.Node) *xmldom.Node {
	transform := xmlutils.DomUtilsAddElement(document, parentNode, t.namespace, common.XMLDSigElement_TRANSFORM)
	transform.SetAttr(xmldom.Name{Local: common.XMLDSigAttribute_ALGORITHM.AttributeName()}, t.algorithm)

	// XPath element must have a specific namespace
	xPathElement := xmlutils.DomUtilsAddTextElement(document, transform, definition.XAdESNamespace_XMLDSIG_FILTER2,
		common.XMLDSigElement_XPATH, t.xPathExpression)

	// Java calls Element#setPrefix here; DomUtilsAddTextElement already created the element
	// with the Filter 2.0 prefix, so re-applying it would be a no-op and is dropped.
	xPathElement.SetAttr(xmldom.Name{
		Space:  xmldom.XMLNSNamespace,
		Local:  definition.XAdESNamespace_XMLDSIG_FILTER2.Prefix(),
		Prefix: "xmlns",
	}, definition.XAdESNamespace_XMLDSIG_FILTER2.Uri())

	// "xmlns" + (prefix empty ? "" : ":" + prefix): a non-empty prefix declares xmlns:<prefix>,
	// an empty one declares the default namespace.
	xmldsigDeclaration := xmldom.Name{Space: xmldom.XMLNSNamespace, Local: "xmlns"}
	if utils.IsStringNotEmpty(t.namespace.Prefix()) {
		xmldsigDeclaration = xmldom.Name{
			Space:  xmldom.XMLNSNamespace,
			Local:  t.namespace.Prefix(),
			Prefix: "xmlns",
		}
	}
	xPathElement.SetAttr(xmldsigDeclaration, t.namespace.Uri())

	xPathElement.SetAttr(xmldom.Name{Local: xPath2FilterTransformFilterAttribute}, t.filter)
	return xPathElement
}

// Equals ports equals(Object): everything XPathTransform#equals compares, plus the filter.
// hashCode() has no Go counterpart.
func (t *XPath2FilterTransform) Equals(obj DSSTransform) bool {
	if !t.XPathTransform.Equals(obj) {
		return false
	}
	other := xPath2FilterTransformOf(obj)
	if other == nil {
		return false
	}
	return t.filter == other.filter
}

// String ports toString().
func (t *XPath2FilterTransform) String() string {
	return "XPath2FilterTransform [filter=" + t.filter + ", xPathExpression=" + t.xPathExpression +
		", algorithm=" + t.algorithm + ", namespace=" + abstractTransformNamespaceString(t.namespace) + "]"
}

// xPath2FilterTransformHolder is the accessor xPath2FilterTransformOf looks for; it is promoted
// to every type embedding XPath2FilterTransform.
type xPath2FilterTransformHolder interface {
	xPath2FilterTransform() *XPath2FilterTransform
}

// xPath2FilterTransformOf extracts the embedded XPath2FilterTransform of any transform.
func xPath2FilterTransformOf(transform DSSTransform) *XPath2FilterTransform {
	if holder, ok := transform.(xPath2FilterTransformHolder); ok {
		return holder.xPath2FilterTransform()
	}
	return nil
}
