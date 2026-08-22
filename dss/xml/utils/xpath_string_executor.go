// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/xpath/XPathStringExecutor.java (DSS 6.5.RC1).
package utils

import "github.com/ryftcore/dss-go/dss/internal/xmldom"

// XPathStringExecutor executes the given XPath String expression.
//
// NAMING: Java overloads getNodeList(Node, XPathQuery) (XPathQueryExecutor) and
// getNodeList(Node, String) (this interface) on the same simple name; Go has no overloading,
// so the string-based query method is named GetNodeListByString here.
type XPathStringExecutor interface {
	// SetNamespaceContext sets the namespace context map containing a declaration of
	// namespaces defined within the used XPath expressions. Ports
	// setNamespaceContext(NamespaceContextMap).
	SetNamespaceContext(namespaceContext *NamespaceContextMap)

	// GetNodeListByString returns the nodes matching xPathString, in document order,
	// without duplicates. Ports getNodeList(Node, String).
	GetNodeListByString(xmlNode *xmldom.Node, xPathString string) ([]*xmldom.Node, error)
}
