// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/DSSXmlErrorListener.java (DSS 6.5.RC1).
//
// DEVIATION (flagged for integrator): upstream implements javax.xml.transform.ErrorListener,
// hooked onto a javax.xml.transform.Transformer to log-then-rethrow warnings/errors/fatal
// errors raised during serialization. Go's serializer (internal/xmldom.Node.Serialize) has no
// pluggable error-listener hook - it just returns an error - and slf4j logging is dropped
// package-wide. DSSXmlErrorListener is kept as a documented, functionally
// inert type so call sites ported unchanged from Java compile: each method stands in for "log
// then rethrow" by returning err unchanged (there is no logger to log to, and nothing in this
// port currently raises through an ErrorListener-shaped hook).
package utils

// DSSXmlErrorListener is DSS's ErrorListener implementation.
type DSSXmlErrorListener struct{}

// NewDSSXmlErrorListener creates a DSSXmlErrorListener. Ports the default constructor.
func NewDSSXmlErrorListener() *DSSXmlErrorListener {
	return &DSSXmlErrorListener{}
}

// Warning ports warning(TransformerException): log then rethrow. There is no logger, so this
// returns err unchanged.
func (l *DSSXmlErrorListener) Warning(err error) error { return err }

// Error ports error(TransformerException): log then rethrow.
func (l *DSSXmlErrorListener) Error(err error) error { return err }

// FatalError ports fatalError(TransformerException): log then rethrow.
func (l *DSSXmlErrorListener) FatalError(err error) error { return err }
