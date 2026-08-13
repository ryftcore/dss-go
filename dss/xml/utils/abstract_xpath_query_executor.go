// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/xpath/AbstractXPathQueryExecutor.java (DSS 6.5.RC1).
package utils

// AbstractXPathQueryExecutor holds the fields and methods shared by every XPathQueryExecutor
// implementation. Ports the abstract class of the same name; concrete executors embed it.
type AbstractXPathQueryExecutor struct {
	// namespaceContext is the map containing the defined namespaces.
	namespaceContext *NamespaceContextMap
}

// SetNamespaceContext implements XPathQueryExecutor. Ports setNamespaceContext(NamespaceContextMap).
func (e *AbstractXPathQueryExecutor) SetNamespaceContext(namespaceContext *NamespaceContextMap) {
	e.namespaceContext = namespaceContext
}
