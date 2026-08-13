// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/xpath/JavaXmlXPathQueryExecutor.java (DSS 6.5.RC1).
//
// Executes XPath expression queries. Upstream compiles javax.xml.xpath.XPathExpression
// objects (Xalan, as shipped in the JDK) against xPathQuery.getQueryString(); this port
// compiles the same query string with internal/xpath10, the Go XPath 1.0 subset engine built
// specifically to evaluate every expression eu.europa.esig.dss.xml.common.xpath.
// XPathQueryBuilder can produce (see internal/xpath10/doc.go).
package utils

import (
	"fmt"

	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/internal/xpath10"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/xml/common"
)

// JavaXmlXPathQueryExecutor is DSS's default XPathQueryExecutor / XPathStringExecutor
// implementation, backed by a real XPath 1.0 engine.
type JavaXmlXPathQueryExecutor struct {
	AbstractXPathQueryExecutor
}

// NewJavaXmlXPathQueryExecutor creates a JavaXmlXPathQueryExecutor. Ports the default
// constructor.
func NewJavaXmlXPathQueryExecutor() *JavaXmlXPathQueryExecutor {
	return &JavaXmlXPathQueryExecutor{}
}

var _ XPathQueryExecutor = (*JavaXmlXPathQueryExecutor)(nil)
var _ XPathStringExecutor = (*JavaXmlXPathQueryExecutor)(nil)

// GetNodeList implements XPathQueryExecutor. Ports getNodeList(Node, XPathQuery).
func (e *JavaXmlXPathQueryExecutor) GetNodeList(xmlNode *xmldom.Node, xPathQuery common.XPathQuery) ([]*xmldom.Node, error) {
	return e.GetNodeListByString(xmlNode, xPathQuery.QueryString())
}

// GetNodeListByString implements XPathStringExecutor. Ports getNodeList(Node, String).
func (e *JavaXmlXPathQueryExecutor) GetNodeListByString(xmlNode *xmldom.Node, xPathString string) ([]*xmldom.Node, error) {
	expr, err := e.createXPathExpression(xPathString)
	if err != nil {
		return nil, err
	}
	nodes, err := expr.Evaluate(xmlNode)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to find a NodeList by the given xPathString '%s'. Reason : %s", xPathString, err.Error()), err)
	}
	return nodes, nil
}

// createXPathExpression ports the protected createXPathExpression(String) helper.
func (e *JavaXmlXPathQueryExecutor) createXPathExpression(xpathString string) (*xpath10.Expr, error) {
	expr, err := xpath10.Compile(xpathString, namespaceContextMapToXPath10(e.namespaceContext))
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to create an XPath expression : %s", err.Error()), err)
	}
	return expr, nil
}
