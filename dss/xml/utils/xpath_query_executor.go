// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/xpath/XPathQueryExecutor.java (DSS 6.5.RC1).
package utils

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/xml/common"
)

// XPathQueryExecutor executes the given common.XPathQuery.
type XPathQueryExecutor interface {
	// SetNamespaceContext sets the namespace context map containing a declaration of
	// namespaces defined within the used XPath expressions. Ports
	// setNamespaceContext(NamespaceContextMap).
	SetNamespaceContext(namespaceContext *NamespaceContextMap)

	// GetNodeList returns the nodes matching xPathQuery, in document order, without
	// duplicates. Ports getNodeList(Node, XPathQuery), whose NodeList is empty rather than
	// nil when nothing matches.
	GetNodeList(xmlNode *xmldom.Node, xPathQuery common.XPathQuery) ([]*xmldom.Node, error)
}
