// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/item/XPathQueryEndItem.java (DSS 6.5.RC1).
package common

import "github.com/ryftcore/dss-go/dss/internal/xmldom"

// XPathQueryEndItem marks the end of an XPath expression chain. Its use is optional and
// does not affect processing, but it prevents further extension of the chain. Unlike every
// other XPathQueryItem it implements the interface directly rather than embedding
// AbstractXPathQueryItem, exactly as Java's XPathQueryEndItem implements XPathQueryItem
// directly rather than extending AbstractXPathQueryItem.
type XPathQueryEndItem struct{}

// NewXPathQueryEndItem creates an XPathQueryEndItem.
func NewXPathQueryEndItem() *XPathQueryEndItem {
	return &XPathQueryEndItem{}
}

// NextItem always returns nil. Ports nextItem().
func (i *XPathQueryEndItem) NextItem() XPathQueryItem {
	return nil
}

// SetNextItem always panics. Ports setNextItem(XPathQueryItem), which Java throws
// UnsupportedOperationException from.
func (i *XPathQueryEndItem) SetNextItem(nextItem XPathQueryItem) XPathQueryItem {
	panic("Unable to continue XPath query after the XPathQueryEnd item.")
}

// MatchNode always returns true. Ports matchNode(Node).
func (i *XPathQueryEndItem) MatchNode(node *xmldom.Node) bool {
	return true
}

// AddParameter always panics. Ports addParameter(XPathQueryParameter), which Java throws
// UnsupportedOperationException from.
func (i *XPathQueryEndItem) AddParameter(parameter XPathQueryParameter) {
	panic("Unable to add parameters to the XPathQueryEnd item.")
}

// Parameters always returns nil. Ports getParameters().
func (i *XPathQueryEndItem) Parameters() []XPathQueryParameter {
	return nil
}

// IsElementRelated always returns true. Ports isElementRelated().
func (i *XPathQueryEndItem) IsElementRelated() bool {
	return true
}

// IsAttributeRelated always returns true. Ports isAttributeRelated().
func (i *XPathQueryEndItem) IsAttributeRelated() bool {
	return true
}

// IsEmpty always returns true. Ports isEmpty().
func (i *XPathQueryEndItem) IsEmpty() bool {
	return true
}

// QueryString always returns "". Ports getQueryString().
func (i *XPathQueryEndItem) QueryString() string {
	return ""
}

var _ XPathQueryItem = (*XPathQueryEndItem)(nil)
