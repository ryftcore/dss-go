// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/XPathQuery.java (DSS 6.5.RC1).
package common

// XPathQuery represents an XPath expression query, ready for evaluation against an XML DOM.
type XPathQuery interface {
	// FirstXPathQueryItem returns the first part of the XPath expression (usually the first
	// element of the XPath string). Ports getFirstXPathQueryItem().
	FirstXPathQueryItem() XPathQueryItem

	// IsAll reports whether all descendants matching the expression are returned; if false,
	// only children of the current element are matched. Ports isAll().
	IsAll() bool

	// IsFromCurrentPosition reports whether evaluation starts from the current position
	// rather than the document root. Ports isFromCurrentPosition().
	IsFromCurrentPosition() bool

	// IsEmpty reports whether the expression contains no item definitions. Ports isEmpty().
	IsEmpty() bool

	// SetNextItem appends nextItem to the expression chain and returns the query. Ports
	// setNextItem(XPathQueryItem).
	SetNextItem(nextItem XPathQueryItem) XPathQuery

	// QueryString returns the string representation of the XPath expression. Ports
	// getQueryString().
	QueryString() string
}
