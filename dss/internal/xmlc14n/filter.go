// Ported from org.apache.xml.security.signature.NodeFilter and the isVisible/isVisibleDO/
// isVisibleInt trio of org.apache.xml.security.c14n.implementations.CanonicalizerBase
// (Apache Santuario xmlsec 3.0.6).
package xmlc14n

import "github.com/ryftcore/dss-go/dss/internal/xmldom"

// NodeFilter decides node-set membership during a document-subset canonicalization. It is the
// port of org.apache.xml.security.signature.NodeFilter, the interface the XML-DSig transform
// pipeline attaches to an XMLSignatureInput: the enveloped-signature transform contributes one
// filter, each ds:XPath transform another, and an XPath Filter 2.0 transform a third.
//
// The three-valued result is Santuario's and is load-bearing:
//
//	 1  include the node
//	 0  exclude the node, but keep walking into its subtree and keep its namespace
//	    declarations in scope for whatever below it IS included
//	-1  exclude the node AND its whole subtree, without any namespace bookkeeping
//
// Filters are consulted in order and the first non-1 answer wins, so a filter never sees a
// node an earlier filter has already rejected. A filter must be pure: the canonicalizer calls
// it once per node per traversal and relies on nothing else.
//
// IsNodeIncludeDO ("DO" = document order) is asked for elements and comments and additionally
// receives the current namespace-symbol-table level, which XPath Filter 2.0 uses to remember
// how deep the subtree it is currently inside began. IsNodeInclude is asked for everything
// else - text, CDATA, processing instructions, attributes, and the parent element again when
// its end tag is written - and must answer from the node alone.
type NodeFilter interface {
	IsNodeInclude(n *xmldom.Node) (int, error)
	IsNodeIncludeDO(n *xmldom.Node, level int) (int, error)
}

// isVisibleDO ports CanonicalizerBase.isVisibleDO.
func (e *engine) isVisibleDO(n *xmldom.Node, level int) int {
	for _, f := range e.filters {
		i, err := f.IsNodeIncludeDO(n, level)
		if err != nil {
			e.noteFilterErr(err)
			return 0
		}
		if i != 1 {
			return i
		}
	}
	if e.subset != nil && !e.subset.Has(n) {
		return 0
	}
	return 1
}

// isVisibleInt ports CanonicalizerBase.isVisibleInt: the three-valued form asked once, of the
// traversal root, before anything is written.
func (e *engine) isVisibleInt(n *xmldom.Node) int {
	for _, f := range e.filters {
		i, err := f.IsNodeInclude(n)
		if err != nil {
			e.noteFilterErr(err)
			return 0
		}
		if i != 1 {
			return i
		}
	}
	if e.subset != nil && !e.subset.Has(n) {
		return 0
	}
	return 1
}

// isVisible ports CanonicalizerBase.isVisible, the boolean form. In subtree mode there is no
// node set and no filter, so everything is visible - Santuario's xpathNodeSet == null.
func (e *engine) isVisible(n *xmldom.Node) bool {
	for _, f := range e.filters {
		i, err := f.IsNodeInclude(n)
		if err != nil {
			e.noteFilterErr(err)
			return false
		}
		if i != 1 {
			return false
		}
	}
	return e.subset == nil || e.subset.Has(n)
}

// noteFilterErr records the first filter failure. Santuario wraps it in a
// CanonicalizationException and unwinds immediately; unwinding out of the traversal here would
// mean returning an error from every visibility test and from both emitters, so the error is
// made sticky instead and Canonicalize reports it in place of the (discarded) output. The
// traversal that keeps running after a filter has failed writes to a buffer nobody reads.
func (e *engine) noteFilterErr(err error) {
	if e.filterErr == nil {
		e.filterErr = err
	}
}
