// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/FromCurrentPositionXPathQuery.java (DSS 6.5.RC1).
package common

// xPathQueryFromCurrentPositionPreamble is the "./" start-from-current-element XPath preamble.
const xPathQueryFromCurrentPositionPreamble = "./"

// FromCurrentPositionXPathQuery gets elements matching the XPath expression starting from
// the current element. The XPath expression path must match completely to succeed.
type FromCurrentPositionXPathQuery struct {
	*AbstractXPathQuery
}

// NewFromCurrentPositionXPathQuery creates a FromCurrentPositionXPathQuery.
func NewFromCurrentPositionXPathQuery() *FromCurrentPositionXPathQuery {
	return &FromCurrentPositionXPathQuery{AbstractXPathQuery: newAbstractXPathQuery(xPathQueryFromCurrentPositionPreamble, false, true)}
}

var _ XPathQuery = (*FromCurrentPositionXPathQuery)(nil)
