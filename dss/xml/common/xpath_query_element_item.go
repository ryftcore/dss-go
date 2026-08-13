// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/item/XPathQueryElementItem.java (DSS 6.5.RC1).
package common

import "github.com/utain/esig/dss/internal/xmldom"

// xPathQueryElementColon is the namespace-prefix separator.
const xPathQueryElementColon = ":"

// XPathQueryElementItem represents a single element to look for within an XPath expression.
type XPathQueryElementItem struct {
	*AbstractXPathQueryItem
	element DSSElement
}

// NewXPathQueryElementItem creates an XPathQueryElementItem. Panics if element is nil (Java
// Objects.requireNonNull(element, "Element cannot be null!")).
func NewXPathQueryElementItem(element DSSElement) *XPathQueryElementItem {
	if element == nil {
		panic("Element cannot be null!")
	}
	return &XPathQueryElementItem{AbstractXPathQueryItem: &AbstractXPathQueryItem{}, element: element}
}

// Element returns the DSSElement. Ports getElement().
func (i *XPathQueryElementItem) Element() DSSElement {
	return i.element
}

// process reports whether node is an Element node with a matching tag name and namespace.
// Ports process(Node); see DSSElement.URI's doc comment for the null-vs-"" argument.
func (i *XPathQueryElementItem) process(node *xmldom.Node) bool {
	if node.Kind == xmldom.Element {
		return i.element.IsSameTagName(node.Name.Local) && (i.element.URI() == "" || i.element.URI() == node.Name.Space)
	}
	return false
}

// MatchNode implements XPathQueryItem.
func (i *XPathQueryElementItem) MatchNode(node *xmldom.Node) bool {
	return i.process(node) && i.matchParameters(node)
}

// IsElementRelated implements XPathQueryItem.
func (i *XPathQueryElementItem) IsElementRelated() bool {
	return true
}

// IsAttributeRelated implements XPathQueryItem.
func (i *XPathQueryElementItem) IsAttributeRelated() bool {
	return false
}

// QueryString implements XPathQueryItem.
func (i *XPathQueryElementItem) QueryString() string {
	if namespace := i.element.Namespace(); namespace != nil {
		return namespace.Prefix() + xPathQueryElementColon + i.element.TagName()
	}
	return i.element.TagName()
}

var _ XPathQueryItem = (*XPathQueryElementItem)(nil)
