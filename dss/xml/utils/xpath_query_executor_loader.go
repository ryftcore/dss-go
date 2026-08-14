// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/xpath/XPathQueryExecutorLoader.java (DSS 6.5.RC1).
//
// DEVIATION: upstream discovers its executor via java.util.ServiceLoader against
// META-INF/services entries. Both of dss-xml-utils's META-INF/services files list
// JavaXmlXPathQueryExecutor, and the XPathQueryExecutor one lists it before
// NativeDOMXPathQueryExecutor - so it is the practical default (ServiceLoader's iterator
// yields providers in the order their declaring files list them). Go has no
// plugin/ServiceLoader mechanism, so loadXPathQueryExecutor/loadXPathStringExecutor are
// hardwired to return NewJavaXmlXPathQueryExecutor(), matching that practical default
// exactly, rather than reproducing Java's "no implementation found in classpath"
// ExceptionInInitializerError branch, which cannot occur here.
package utils

// XPathQueryExecutorLoader loads an implementation of the corresponding XPath executor.
type XPathQueryExecutorLoader struct {
	// xPathQueryExecutor is the cached version of the executor.
	xPathQueryExecutor XPathQueryExecutor

	// xPathStringExecutor is the cached version of the XPath string executor.
	xPathStringExecutor XPathStringExecutor
}

// NewXPathQueryExecutorLoader creates an XPathQueryExecutorLoader. Ports the default
// constructor.
func NewXPathQueryExecutorLoader() *XPathQueryExecutorLoader {
	return &XPathQueryExecutorLoader{}
}

// GetXPathQueryExecutor returns a cached or provided XPathQueryExecutor, loading the default
// implementation on first use. Ports getXPathQueryExecutor().
func (l *XPathQueryExecutorLoader) GetXPathQueryExecutor() XPathQueryExecutor {
	if l.xPathQueryExecutor == nil {
		l.xPathQueryExecutor = l.loadXPathQueryExecutor()
	}
	return l.xPathQueryExecutor
}

// SetXPathQueryExecutor sets the XPathQueryExecutor to be used. Ports setXPathQueryExecutor(XPathQueryExecutor).
func (l *XPathQueryExecutorLoader) SetXPathQueryExecutor(xPathQueryExecutor XPathQueryExecutor) {
	l.xPathQueryExecutor = xPathQueryExecutor
}

// loadXPathQueryExecutor loads the default implementation of XPathQueryExecutor. Ports
// loadXPathQueryExecutor(); see the file header for the ServiceLoader deviation.
func (l *XPathQueryExecutorLoader) loadXPathQueryExecutor() XPathQueryExecutor {
	return NewJavaXmlXPathQueryExecutor()
}

// GetXPathStringExecutor returns a cached or provided XPathStringExecutor, loading the
// default implementation on first use. Ports getXPathStringExecutor().
func (l *XPathQueryExecutorLoader) GetXPathStringExecutor() XPathStringExecutor {
	if l.xPathStringExecutor == nil {
		l.xPathStringExecutor = l.loadXPathStringExecutor()
	}
	return l.xPathStringExecutor
}

// SetXPathStringExecutor sets the XPathStringExecutor to be used. Ports setXPathStringExecutor(XPathStringExecutor).
func (l *XPathQueryExecutorLoader) SetXPathStringExecutor(xPathStringExecutor XPathStringExecutor) {
	l.xPathStringExecutor = xPathStringExecutor
}

// loadXPathStringExecutor loads the default implementation of XPathStringExecutor. Ports
// loadXPathStringExecutor(); see the file header for the ServiceLoader deviation.
func (l *XPathQueryExecutorLoader) loadXPathStringExecutor() XPathStringExecutor {
	return NewJavaXmlXPathQueryExecutor()
}
