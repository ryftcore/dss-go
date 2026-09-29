// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/xpath/AbstractXPathQueryExecutor.java (DSS 6.5.RC1).
package utils

import "sync/atomic"

// AbstractXPathQueryExecutor holds the fields and methods shared by every XPathQueryExecutor
// implementation. Ports the abstract class of the same name; concrete executors embed it.
type AbstractXPathQueryExecutor struct {
	// namespaceContext is the map containing the defined namespaces. It is an atomic pointer
	// because XPathUtilsGetXPathQueryExecutor re-sets it on the shared, process-global
	// executor for every query, while other goroutines read it to compile theirs.
	namespaceContext atomic.Pointer[NamespaceContextMap]
}

// SetNamespaceContext implements XPathQueryExecutor. Ports setNamespaceContext(NamespaceContextMap).
func (e *AbstractXPathQueryExecutor) SetNamespaceContext(namespaceContext *NamespaceContextMap) {
	e.namespaceContext.Store(namespaceContext)
}

// getNamespaceContext returns the namespace context set by SetNamespaceContext, or nil.
func (e *AbstractXPathQueryExecutor) getNamespaceContext() *NamespaceContextMap {
	return e.namespaceContext.Load()
}
