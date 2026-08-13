// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/xpath/AllXPathQuery.java (DSS 6.5.RC1).
package common

// xPathQueryAllPreamble is the "//" search-the-whole-document XPath preamble.
const xPathQueryAllPreamble = "//"

// AllXPathQuery gets all elements from the XML document matching the given XPath
// expression, starting evaluation from the root document element.
type AllXPathQuery struct {
	*AbstractXPathQuery
}

// NewAllXPathQuery creates an AllXPathQuery.
func NewAllXPathQuery() *AllXPathQuery {
	return &AllXPathQuery{AbstractXPathQuery: newAbstractXPathQuery(xPathQueryAllPreamble, true, false)}
}

var _ XPathQuery = (*AllXPathQuery)(nil)
