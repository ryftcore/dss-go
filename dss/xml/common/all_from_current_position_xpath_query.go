// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/AllFromCurrentPositionXPathQuery.java (DSS 6.5.RC1).
package common

// xPathQueryAllFromCurrentPositionPreamble is the ".//" search-all-descendants-from-current-element
// XPath preamble.
const xPathQueryAllFromCurrentPositionPreamble = ".//"

// AllFromCurrentPositionXPathQuery gets all elements within the current parent matching the
// XPath expression, searching both descendants and direct children.
type AllFromCurrentPositionXPathQuery struct {
	*AbstractXPathQuery
}

// NewAllFromCurrentPositionXPathQuery creates an AllFromCurrentPositionXPathQuery.
func NewAllFromCurrentPositionXPathQuery() *AllFromCurrentPositionXPathQuery {
	return &AllFromCurrentPositionXPathQuery{AbstractXPathQuery: newAbstractXPathQuery(xPathQueryAllFromCurrentPositionPreamble, true, true)}
}

var _ XPathQuery = (*AllFromCurrentPositionXPathQuery)(nil)
