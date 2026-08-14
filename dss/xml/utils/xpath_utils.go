// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/xpath/XPathUtils.java (DSS 6.5.RC1).
//
// NAMING: dss-xml-utils merges two upstream static-utility classes with overlapping method
// names (DomUtils and XPathUtils both declare GetValue/GetNodeList/GetNode/GetElement/
// GetNodesAmount/GetChildrenNames/GetElementById/RegisterNamespace, with different
// semantics) into one Go package, following the DomUtilsXxx/XPathUtilsXxx-style class-name
// prefixing already established for merged "Utils" classes in package spi (DSSUtilsXxx,
// DSSASN1UtilsXxx). This file's exported functions are therefore prefixed XPathUtils.
package utils

import (
	"fmt"

	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/xml/common"
)

// xPathUtilsXMLNS is the default namespace prefix, which RegisterNamespace refuses to bind.
const xPathUtilsXMLNS = "xmlns"

var (
	// xPathUtilsNamespacePrefixMapper is the map containing the defined namespaces.
	xPathUtilsNamespacePrefixMapper = NewNamespaceContextMap()

	// xPathUtilsExecutorLoader is the XPath query executor loader.
	xPathUtilsExecutorLoader = NewXPathQueryExecutorLoader()
)

// XPathUtilsGetXPathQueryExecutor gets the XPath Query executor. Ports getXPathQueryExecutor().
func XPathUtilsGetXPathQueryExecutor() XPathQueryExecutor {
	executor := xPathUtilsExecutorLoader.GetXPathQueryExecutor()
	executor.SetNamespaceContext(xPathUtilsNamespacePrefixMapper)
	return executor
}

// XPathUtilsSetXPathQueryExecutor sets the XPathQueryExecutor to be used by the
// implementation. Ports setXPathQueryExecutor(XPathQueryExecutor).
func XPathUtilsSetXPathQueryExecutor(xPathQueryExecutor XPathQueryExecutor) {
	xPathUtilsExecutorLoader.SetXPathQueryExecutor(xPathQueryExecutor)
}

// XPathUtilsRegisterNamespace registers a namespace and associated prefix. If the prefix
// exists already it is replaced. It reports whether the map did not already contain the
// specified element. Ports registerNamespace(DSSNamespace); the two
// UnsupportedOperationException cases become panics per PORTING.md.
func XPathUtilsRegisterNamespace(namespace *common.DSSNamespace) bool {
	prefix := namespace.Prefix()
	uri := namespace.Uri()
	if utils.IsStringEmpty(prefix) {
		panic("The empty namespace cannot be registered!")
	}
	if xPathUtilsXMLNS == prefix {
		panic(fmt.Sprintf("The default namespace '%s' cannot be registered!", xPathUtilsXMLNS))
	}
	return xPathUtilsNamespacePrefixMapper.RegisterNamespace(prefix, uri)
}

// XPathUtilsGetNamespaceContextMap returns the stored namespace definitions map. Ports
// getNamespaceContextMap().
func XPathUtilsGetNamespaceContextMap() *NamespaceContextMap {
	return xPathUtilsNamespacePrefixMapper
}

// XPathUtilsGetValue returns the String value of the corresponding to the XPath query: the
// unique matched node's text content, trimmed, or "" (Java null) if no node matches. Ports
// getValue(Node, XPathQuery).
//
// NOTE the different semantics from DomUtilsGetValue: this requires a single match (erroring
// on more than one, via GetNode) and reads getTextContent(), whereas DomUtilsGetValue
// evaluates the query directly as an XPath STRING result (first node's string value, no
// "more than one" check). Both are faithful ports of their respective, differently-behaved
// upstream methods.
func XPathUtilsGetValue(xmlNode *xmldom.Node, xPathQuery common.XPathQuery) (string, error) {
	node, err := XPathUtilsGetNode(xmlNode, xPathQuery)
	if err != nil {
		return "", err
	}
	if node == nil {
		return "", nil
	}
	return utils.Trim(node.TextContent()), nil
}

// XPathUtilsGetNodeList returns the nodes corresponding to the XPath query, in document
// order. Ports getNodeList(Node, XPathQuery).
func XPathUtilsGetNodeList(xmlNode *xmldom.Node, xPathQuery common.XPathQuery) ([]*xmldom.Node, error) {
	return XPathUtilsGetXPathQueryExecutor().GetNodeList(xmlNode, xPathQuery)
}

// XPathUtilsGetNode returns the Node corresponding to the XPath query, or nil if none match.
// Ports getNode(Node, XPathQuery).
func XPathUtilsGetNode(xmlNode *xmldom.Node, xPathQuery common.XPathQuery) (*xmldom.Node, error) {
	list, err := XPathUtilsGetNodeList(xmlNode, xPathQuery)
	if err != nil {
		return nil, err
	}
	if len(list) > 1 {
		return nil, model.NewDSSError(fmt.Sprintf("More than one result for XPath: %s", xPathQuery.QueryString()))
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

// XPathUtilsGetElement returns the Element corresponding to the XPath query. Ports
// getElement(Node, XPathQuery).
func XPathUtilsGetElement(xmlNode *xmldom.Node, xPathQuery common.XPathQuery) (*xmldom.Node, error) {
	return XPathUtilsGetNode(xmlNode, xPathQuery)
}

// XPathUtilsGetNodesAmount returns the amount of nodes matching the XPath query. Ports
// getNodesAmount(Node, XPathQuery).
func XPathUtilsGetNodesAmount(xmlNode *xmldom.Node, xPathQuery common.XPathQuery) (int, error) {
	list, err := XPathUtilsGetNodeList(xmlNode, xPathQuery)
	if err != nil {
		return 0, err
	}
	return len(list), nil
}

// XPathUtilsGetChildrenNames returns the list of children's names for the element matched by
// xPathQuery. Ports getChildrenNames(Node, XPathQuery), which only collects children whose
// getLocalName() is non-null - i.e., Element/Attribute/ProcInst children, never Text/CDATA/
// Comment.
func XPathUtilsGetChildrenNames(xmlNode *xmldom.Node, xPathQuery common.XPathQuery) ([]string, error) {
	var childrenNames []string
	element, err := XPathUtilsGetElement(xmlNode, xPathQuery)
	if err != nil {
		return nil, err
	}
	if element != nil {
		for c := element.FirstChild; c != nil; c = c.NextSibling {
			if c.Name.Local != "" {
				childrenNames = append(childrenNames, c.Name.Local)
			}
		}
	}
	return childrenNames, nil
}

// XPathUtilsGetElementById extracts an element from the given document node with the given
// Id, namespace independently. Ports getElementById(Node, String).
func XPathUtilsGetElementById(node *xmldom.Node, id string) *xmldom.Node {
	return XPathUtilsGetElementByIdWithQuery(node, common.XPathQueryBuilderAllFromCurrentPosition().Build(), id)
}

// XPathUtilsGetElementByIdWithQuery extracts an element from the given node according to
// xPathQuery with the matching id, if any, normalizing id safely (ensuring it is not a URI or
// XPointer value). It returns nil if no unique result is found, or if the lookup errors -
// mirroring Java's catch(Exception) { LOG.warn(...); return null; }, with logging dropped.
// Ports getElementById(Node, XPathQuery, String).
func XPathUtilsGetElementByIdWithQuery(node *xmldom.Node, xPathQuery common.XPathQuery, id string) (result *xmldom.Node) {
	defer func() {
		if recover() != nil {
			result = nil
		}
	}()
	built := common.XPathQueryBuilderFromXPathQuery(xPathQuery).IdValue(DomUtilsGetId(id)).Build()
	element, err := XPathUtilsGetElement(node, built)
	if err != nil {
		return nil
	}
	return element
}
