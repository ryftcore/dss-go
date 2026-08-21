// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/item/XPathQueryAnyItem.java (DSS 6.5.RC1).
package common

import "github.com/ryftcore/dss-go/dss/internal/xmldom"

// xPathQueryAnyPath is the "*" XPath any-element wildcard.
const xPathQueryAnyPath = "*"

// XPathQueryAnyItem is an XPath expression chain item matching any element within the path.
type XPathQueryAnyItem struct {
	*AbstractXPathQueryItem
}

// NewXPathQueryAnyItem creates an XPathQueryAnyItem.
func NewXPathQueryAnyItem() *XPathQueryAnyItem {
	return &XPathQueryAnyItem{AbstractXPathQueryItem: &AbstractXPathQueryItem{}}
}

// process reports whether node is an Element node. Ports process(Node).
func (i *XPathQueryAnyItem) process(node *xmldom.Node) bool {
	return node.Kind == xmldom.Element
}

// MatchNode implements XPathQueryItem.
func (i *XPathQueryAnyItem) MatchNode(node *xmldom.Node) bool {
	return i.process(node) && i.matchParameters(node)
}

// IsElementRelated implements XPathQueryItem.
func (i *XPathQueryAnyItem) IsElementRelated() bool {
	return true
}

// IsAttributeRelated implements XPathQueryItem.
func (i *XPathQueryAnyItem) IsAttributeRelated() bool {
	return false
}

// QueryString implements XPathQueryItem.
func (i *XPathQueryAnyItem) QueryString() string {
	return xPathQueryAnyPath
}

var _ XPathQueryItem = (*XPathQueryAnyItem)(nil)
