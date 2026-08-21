// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/XPathTransform.java (DSS 6.5.RC1).
//
// org.apache.xml.security.transforms.Transforms.TRANSFORM_XPATH is xmldsig.TransformXPath, per
// internal/xmldsig's doc.go mapping table; the transform itself is executed by that package's
// registry through ComplexTransform.
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// XPathTransform is the XPath transform.
type XPathTransform struct {
	ComplexTransform

	// xPathExpression is the XPath expression to use.
	xPathExpression string
}

// NewXPathTransform ports XPathTransform(String), which delegates with XMLDSigNamespace.NS and
// the ds:XPath algorithm URI.
func NewXPathTransform(xPathExpression string) *XPathTransform {
	return NewXPathTransformWithNamespace(common.XMLDSigNS, xPathExpression)
}

// NewXPathTransformWithNamespace ports XPathTransform(DSSNamespace, String).
func NewXPathTransformWithNamespace(xmlDSigNamespace *common.DSSNamespace, xPathExpression string) *XPathTransform {
	t := new(XPathTransform)
	*t = newXPathTransform(xmlDSigNamespace, xmldsig.TransformXPath, xPathExpression)
	t.InitAbstractTransform(t)
	return t
}

// newXPathTransform ports the protected XPathTransform(DSSNamespace, String, String), the
// constructor a subclass reaches for a custom algorithm URI. It returns a value so that a
// subclass can embed it; the caller must still call InitAbstractTransform with itself.
//
// Java's Objects.requireNonNull(xPathExpression) has no Go analogue: a Go string is never nil.
func newXPathTransform(xmlDSigNamespace *common.DSSNamespace, algorithm, xPathExpression string) XPathTransform {
	return XPathTransform{
		ComplexTransform: newComplexTransform(xmlDSigNamespace, algorithm),
		xPathExpression:  xPathExpression,
	}
}

// xPathTransform returns the embedded XPathTransform. Promoted through embedding, it is how a
// subclass's Equals reaches the field Java's XPathTransform#equals compares directly.
func (t *XPathTransform) xPathTransform() *XPathTransform { return t }

// CreateTransform appends the ds:Transform element plus its ds:XPath child carrying the
// expression. Ports createTransform(Document, Element); the returned element is the ds:XPath
// one, as upstream.
func (t *XPathTransform) CreateTransform(document, parentNode *xmldom.Node) *xmldom.Node {
	transform := t.AbstractTransform.CreateTransform(document, parentNode)
	return xmlutils.DomUtilsAddTextElement(document, transform, t.namespace,
		common.XMLDSigElement_XPATH, t.xPathExpression)
}

// Equals ports equals(Object): everything AbstractTransform#equals compares, plus the XPath
// expression. hashCode() has no Go counterpart.
func (t *XPathTransform) Equals(obj DSSTransform) bool {
	if !t.AbstractTransform.Equals(obj) {
		return false
	}
	other := xPathTransformOf(obj)
	if other == nil {
		return false
	}
	return t.xPathExpression == other.xPathExpression
}

// String ports toString().
func (t *XPathTransform) String() string {
	return "XPathTransform [xPathExpression=" + t.xPathExpression + ", algorithm=" + t.algorithm +
		", namespace=" + abstractTransformNamespaceString(t.namespace) + "]"
}

// xPathTransformHolder is the accessor xPathTransformOf looks for; it is promoted to every
// type embedding XPathTransform.
type xPathTransformHolder interface {
	xPathTransform() *XPathTransform
}

// xPathTransformOf extracts the embedded XPathTransform of any transform.
func xPathTransformOf(transform DSSTransform) *XPathTransform {
	if holder, ok := transform.(xPathTransformHolder); ok {
		return holder.xPathTransform()
	}
	return nil
}
