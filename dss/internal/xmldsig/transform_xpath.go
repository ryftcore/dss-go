// Ported from org.apache.xml.security.transforms.implementations.TransformXPath and its
// XPathNodeFilter (Apache Santuario xmlsec 3.0.6).
package xmldsig

import (
	"errors"
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xpath10"
)

// xpathTransform is the ds:XPath transform of XMLDSIG 6.6.3. Port of TransformXPath.
//
// The evaluation model is per-node boolean filtering, not "select these nodes": the expression
// is compiled once and then evaluated with EVERY node the canonicalizer visits as the context
// node, and the node survives when the expression's boolean value is true. That is why
// "not(ancestor-or-self::ds:Signature)" works - it is asked of each node in turn - and why the
// transform sets the input's node-set flag: from here on the canonicalizer must take the
// document-subset path so that a node can be invisible while its descendants are not.
//
// Prefixes in the expression resolve against the declarations in scope on the ds:XPath ELEMENT,
// which is DOMNamespaceContext(xpathElement) in JDKXPathAPI. Nothing else is consulted: a
// prefix the ds:XPath element does not declare is bound to no namespace, so the name it
// qualifies matches only unnamespaced nodes, exactly as in upstream.
//
// Santuario's needsCircumvent/circumventBug2650 is deliberately not ported. It copies every
// ancestor namespace declaration onto every descendant element of the DOM before evaluating an
// expression that mentions "namespace" or "name()", to work around Xalan's non-conformant
// namespace axis (Apache bug 2650). internal/xpath10 has no namespace axis to be wrong about,
// and its name() reads Name.Prefix from the node itself rather than from a declaration, so the
// workaround has nothing to fix - and porting it would mean MUTATING the document under
// validation, which is exactly the kind of surprise a validator must not spring.
type xpathTransform struct{}

func (xpathTransform) Algorithm() string { return TransformXPath }

func (xpathTransform) Perform(in *Data, element *xmldom.Node, baseURI string, secureValidation bool) (*Data, error) {
	xpathElement := selectDSNode(element, "XPath", 0)
	if xpathElement == nil {
		return nil, errors.New("xmldsig: the XPath transform has no ds:XPath child")
	}
	first := xpathElement.FirstChild
	if first == nil {
		// DOMException(HIERARCHY_REQUEST_ERR, "Text must be in ds:Xpath").
		return nil, errors.New("xmldsig: the ds:XPath element is empty")
	}
	expr, err := xpath10.CompileTransform(stringFromNode(first), xpath10.NamespaceContextOf(xpathElement))
	if err != nil {
		return nil, fmt.Errorf("xmldsig: XPath transform: %w", err)
	}
	if err := in.AddNodeFilter(xpathNodeFilter{expr: expr}); err != nil {
		return nil, err
	}
	in.SetNodeSet(true)
	return in, nil
}

// xpathNodeFilter is XPathNodeFilter: one boolean evaluation per node, and the same answer
// whether the canonicalizer asks in document order or not - the expression is stateless, so
// isNodeIncludeDO simply delegates to isNodeInclude, as it does upstream.
//
// Note the answers: true is 1 and false is 0, never -1. A ds:XPath transform can therefore
// exclude an element while keeping its children, which is what makes the node-set traversal
// (rather than a subtree walk) necessary.
type xpathNodeFilter struct{ expr *xpath10.Expr }

func (f xpathNodeFilter) IsNodeInclude(n *xmldom.Node) (int, error) {
	ok, err := f.expr.EvaluateBoolean(n)
	if err != nil {
		return 0, err
	}
	if ok {
		return 1, nil
	}
	return 0, nil
}

func (f xpathNodeFilter) IsNodeIncludeDO(n *xmldom.Node, level int) (int, error) {
	return f.IsNodeInclude(n)
}
