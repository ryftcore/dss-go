// Ported from org.apache.xml.security.c14n.implementations.CanonicalizerPhysical
// (Apache Santuario xmlsec 3.0.6): attribute emission for the Santuario "physical" method.
package xmlc14n

import "github.com/utain/esig/dss/internal/xmldom"

// outputAttributesSubtreePhysical ports CanonicalizerPhysical.outputAttributesSubtree.
//
// "Physical" means exactly that: every attribute the element actually carries is emitted,
// sorted by the same comparator as everywhere else, and nothing else happens. No namespace
// inheritance from the apex's ancestors, no suppression of a superfluous redeclaration, no
// synthesized xmlns="", no relative-URI check, and - unlike every other algorithm - the
// declaration of the xml prefix is emitted rather than skipped, because the loop has no
// special case for it. The output can therefore contain a prefix with no declaration in
// scope, which is by design: the method exists to hash a subtree's literal bytes, not to
// produce an independently parseable document.
func (e *engine) outputAttributesSubtreePhysical(el *xmldom.Node, _ *nsStack) error {
	if len(el.Attrs) == 0 {
		return nil
	}
	var result attrSet
	for _, attr := range el.Attrs {
		result.add(attrOf(attr))
	}
	result.writeTo(e.w)
	return nil
}

// outputAttributesPhysical ports CanonicalizerPhysical.outputAttributes, which throws
// CanonicalizationException("c14n.Canonicalizer.UnsupportedOperation"). The physical method
// has no node-set form; Canonicalize refuses Input.Subset for it up front, so this is only
// the backstop.
func (e *engine) outputAttributesPhysical(*xmldom.Node, *nsStack) error {
	return ErrPhysicalNodeSet
}
