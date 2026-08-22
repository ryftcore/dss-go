// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/item/XPathQueryNotChildOfParameter.java (DSS 6.5.RC1).
package common

import "github.com/ryftcore/dss-go/dss/internal/xmldom"

const (
	xPathQueryNotChildOfConditionStart = "not(parent::"
	xPathQueryNotChildOfConditionEnd   = ")"
)

// XPathQueryNotChildOfParameter represents an item within an XPath expression filtering out
// elements with a particular parent element.
type XPathQueryNotChildOfParameter struct {
	*AbstractXPathQueryParameter
	elementItem *XPathQueryElementItem
}

// NewXPathQueryNotChildOfParameter creates a parameter rejecting nodes whose parent matches
// parentElement. Panics if parentElement is nil (Java Objects.requireNonNull(parentElement,
// "Parent element cannot be null!")).
func NewXPathQueryNotChildOfParameter(parentElement DSSElement) *XPathQueryNotChildOfParameter {
	if parentElement == nil {
		panic("Parent element cannot be null!")
	}
	return &XPathQueryNotChildOfParameter{
		AbstractXPathQueryParameter: newAbstractXPathQueryParameter(),
		elementItem:                 NewXPathQueryElementItem(parentElement),
	}
}

// ParentElement returns the parent element to be avoided. Ports getParentElement().
func (p *XPathQueryNotChildOfParameter) ParentElement() DSSElement {
	return p.elementItem.Element()
}

// process reports whether node's parent does not match the excluded parent element. Ports
// process(Node).
//
// Judgment call: Java's process only ever runs on an Element node, whose getParentNode() is
// never null in a well-formed tree (an element's parent is either another element or the
// owning Document); a genuinely detached element - not something this codebase ever passes -
// would NPE. This port makes the same assumption rather than adding a defensive nil check:
// node.Parent == nil for a detached element flows straight into elementItem.MatchNode(nil),
// whose process dereferences node.Kind and panics, mirroring Java's NPE on the same
// unreachable input.
func (p *XPathQueryNotChildOfParameter) process(node *xmldom.Node) bool {
	if node.Kind != xmldom.Element {
		return false
	}
	return !p.elementItem.MatchNode(node.Parent)
}

// MatchNode implements XPathQueryItem.
func (p *XPathQueryNotChildOfParameter) MatchNode(node *xmldom.Node) bool {
	return p.process(node)
}

// IsElementRelated implements XPathQueryItem.
func (p *XPathQueryNotChildOfParameter) IsElementRelated() bool {
	return true
}

// IsAttributeRelated implements XPathQueryItem.
func (p *XPathQueryNotChildOfParameter) IsAttributeRelated() bool {
	return false
}

// QueryString implements XPathQueryItem.
func (p *XPathQueryNotChildOfParameter) QueryString() string {
	return xPathQueryNotChildOfConditionStart + p.elementItem.QueryString() + xPathQueryNotChildOfConditionEnd
}

var _ XPathQueryParameter = (*XPathQueryNotChildOfParameter)(nil)
