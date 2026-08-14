// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/item/XPathQueryItem.java (DSS 6.5.RC1).
package common

import "github.com/utain/esig/dss/internal/xmldom"

// XPathQueryItem represents a single XPath expression chain item. MatchNode evaluates one
// chain link directly against an *xmldom.Node (the org.w3c.dom.Node counterpart), which is
// what lets dss-xml-utils's native DOM query executor (a later phase) walk a document
// without a full XPath engine - see doc.go.
type XPathQueryItem interface {
	// NextItem returns the next XPath chain item, or nil if none. Ports nextItem().
	NextItem() XPathQueryItem

	// SetNextItem sets the next chain item and returns it. Ports setNextItem(XPathQueryItem).
	SetNextItem(nextItem XPathQueryItem) XPathQueryItem

	// MatchNode reports whether node satisfies this XPathQueryItem.
	MatchNode(node *xmldom.Node) bool

	// AddParameter adds a parameter to this XPathQueryItem.
	AddParameter(parameter XPathQueryParameter)

	// Parameters returns the parameters related to this XPathQueryItem, if any.
	Parameters() []XPathQueryParameter

	// IsElementRelated reports whether this item processes Element nodes.
	IsElementRelated() bool

	// IsAttributeRelated reports whether this item processes Attribute nodes.
	IsAttributeRelated() bool

	// IsEmpty reports whether this XPathQueryItem is empty.
	IsEmpty() bool

	// QueryString returns the string representation of this chain item.
	QueryString() string
}
