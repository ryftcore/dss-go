// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/xpath/NativeDOMXPathQueryExecutor.java (DSS 6.5.RC1).
//
// ASSUMPTION/DEVIATION (flagged prominently for integrator review): upstream implements this
// class - marked "(Experimental)" in its own Javadoc - as a hand-rolled DOM walk over the
// XPathQueryItem chain (XPathQueryItem#matchNode/isElementRelated/isAttributeRelated), never
// compiling an XPath string at all. xml/common/doc.go, written by this phase's xml/common
// implementer, explicitly anticipates that design: "consumed by dss-xml-utils's
// NativeDOMXPathQueryExecutor - a later phase - which is why this package already depends on
// internal/xmldom rather than leaving matchNode unported."
//
// This phase's task brief instead directs collapsing JavaXmlXPathQueryExecutor and
// NativeDOMXPathQueryExecutor into "a single native executor on xpath10", so this type is a
// thin embedding of JavaXmlXPathQueryExecutor rather than an independent MatchNode-based
// walker: Go has no second XPath backend to mirror Xalan-vs-hand-rolled-walker with, and
// internal/xpath10 was purpose-built (per its doc.go) to evaluate exactly the expressions
// XPathQueryBuilder emits, so routing both executors through it avoids maintaining two
// query engines for one behaviour. The XPathQueryItem.MatchNode chain in package
// xml/common remains fully ported and unused by this package; it is available intact should
// a later phase prefer a native (non-xpath10) walker - e.g. to avoid a full XPath grammar on
// a hot path. Reconcile with the tech lead if the MatchNode-based design was intended
// specifically for this file.
package utils

import (
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/xml/common"
)

// NativeDOMXPathQueryExecutor is DSS's "(Experimental)" implementation of XPathQueryExecutor
// based on native XML DOM Node processing. In this port it shares
// JavaXmlXPathQueryExecutor's xpath10-backed implementation; see the file header for why.
type NativeDOMXPathQueryExecutor struct {
	JavaXmlXPathQueryExecutor
}

// NewNativeDOMXPathQueryExecutor creates a NativeDOMXPathQueryExecutor. Ports the default
// constructor.
func NewNativeDOMXPathQueryExecutor() *NativeDOMXPathQueryExecutor {
	return &NativeDOMXPathQueryExecutor{}
}

var _ XPathQueryExecutor = (*NativeDOMXPathQueryExecutor)(nil)

// GetNodeList implements XPathQueryExecutor. Ports getNodeList(Node, XPathQuery); see the
// file header for the deviation from upstream's hand-rolled DOM walk.
func (e *NativeDOMXPathQueryExecutor) GetNodeList(xmlNode *xmldom.Node, xPathQuery common.XPathQuery) ([]*xmldom.Node, error) {
	return e.JavaXmlXPathQueryExecutor.GetNodeList(xmlNode, xPathQuery)
}
