// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/item/XPathQueryAttributeItem.java (DSS 6.5.RC1).
package common

import "github.com/ryftcore/dss-go/dss/internal/xmldom"

// xPathQueryAttributePrefix is the "@" XPath attribute-value marker.
const xPathQueryAttributePrefix = "@"

// XPathQueryAttributeItem is an XPath expression item allowing access to an attribute value
// of the current element, independent of any expected value (contrast
// XPathQueryAttributeParameter, which requires a value).
type XPathQueryAttributeItem struct {
	*AbstractXPathQueryItem
	attribute DSSAttribute
}

// NewXPathQueryAttributeItem creates an item extracting an element containing the given
// attribute (any value is accepted). Panics if attribute is nil (Java
// Objects.requireNonNull(attribute, "Attribute cannot be null!")).
func NewXPathQueryAttributeItem(attribute DSSAttribute) *XPathQueryAttributeItem {
	if attribute == nil {
		panic("Attribute cannot be null!")
	}
	return &XPathQueryAttributeItem{AbstractXPathQueryItem: &AbstractXPathQueryItem{}, attribute: attribute}
}

// Attribute returns the corresponding DSSAttribute. Ports getAttribute().
func (i *XPathQueryAttributeItem) Attribute() DSSAttribute {
	return i.attribute
}

// process reports whether node is an Attribute node with the expected local name. Ports
// process(Node).
func (i *XPathQueryAttributeItem) process(node *xmldom.Node) bool {
	if node.Kind == xmldom.Attribute {
		return i.attribute.AttributeName() == localName(node)
	}
	return false
}

// MatchNode implements XPathQueryItem.
func (i *XPathQueryAttributeItem) MatchNode(node *xmldom.Node) bool {
	return i.process(node) && i.matchParameters(node)
}

// IsElementRelated implements XPathQueryItem.
func (i *XPathQueryAttributeItem) IsElementRelated() bool {
	return false
}

// IsAttributeRelated implements XPathQueryItem.
func (i *XPathQueryAttributeItem) IsAttributeRelated() bool {
	return true
}

// QueryString implements XPathQueryItem.
func (i *XPathQueryAttributeItem) QueryString() string {
	return xPathQueryAttributePrefix + i.attribute.AttributeName()
}

var _ XPathQueryItem = (*XPathQueryAttributeItem)(nil)
