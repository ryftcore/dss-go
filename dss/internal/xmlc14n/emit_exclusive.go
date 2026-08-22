// Ported from org.apache.xml.security.c14n.implementations.Canonicalizer20010315Excl
// (Apache Santuario xmlsec 3.0.6): attribute emission for Exclusive XML Canonicalization 1.0.
package xmlc14n

import (
	"sort"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
)

// outputAttributesSubtreeExclusive ports Canonicalizer20010315Excl.outputAttributesSubtree.
//
// Nothing is inherited. At each output element the visibly-utilized prefix set is computed -
// the element's own prefix, or "xmlns" when the element name is unprefixed; the prefix of
// every output attribute except xml and xmlns; and every entry of the InclusiveNamespaces
// PrefixList - and only those bindings are rendered. An unprefixed attribute is in no
// namespace and therefore does not utilize the default namespace, which is why a default
// declaration migrates down to the element that actually uses it.
func (e *engine) outputAttributesSubtreeExclusive(el *xmldom.Node, ns *nsStack) error {
	var result attrSet
	visiblyUtilized := make(map[string]struct{}, len(el.Attrs)+len(e.inclusivePrefixes)+1)
	for _, prefix := range e.inclusivePrefixes {
		visiblyUtilized[prefix] = struct{}{}
	}

	for _, attr := range el.Attrs {
		name, value := attr.Name.Local, attr.Value
		if attr.Name.Space != xmldom.XMLNSNamespace {
			if prefix := attr.Name.Prefix; prefix != "" && prefix != xmlPrefix && prefix != xmlnsPrefix {
				visiblyUtilized[prefix] = struct{}{}
			}
			result.add(attrOf(attr))
			continue
		}
		// The short-circuit order is Santuario's: addMapping runs only for a declaration that
		// is not the xml one, and the relative-URI check runs only when addMapping reports a
		// new definition. A relative declaration inherited from an ancestor is rendered below
		// through getMapping and is never checked - see the ns-relative KATs.
		if name == xmlPrefix && value == xmldom.XMLNamespace {
			continue
		}
		if ns.addMapping(name, value, true) && namespaceIsRelative(value) {
			return &RelativeNamespaceError{Element: el.Name.QName(), Prefix: name, URI: value}
		}
	}

	// Santuario's propagateDefaultNamespace branch is deliberately not ported: it exists for
	// XML Encryption, DSS never sets it, and Canonicalize has no knob for it (design 2.3).

	prefix := xmlnsPrefix
	if el.Name.Space != "" && el.Name.Prefix != "" {
		prefix = el.Name.Prefix
	}
	visiblyUtilized[prefix] = struct{}{}

	for _, s := range sortedKeys(visiblyUtilized) {
		if rendered, ok := ns.getMapping(s); ok {
			result.add(rendered)
		}
	}
	result.writeTo(e.w)
	return nil
}

// outputAttributesExclusive ports Canonicalizer20010315Excl.outputAttributes, the node-set
// form: an element outside the node set renders nothing itself but still updates the symbol
// table for its visible descendants, and a PrefixList entry can be rendered from such an
// element.
func (e *engine) outputAttributesExclusive(el *xmldom.Node, ns *nsStack) error {
	var result attrSet
	// isVisibleDO, not isVisible: see the note in outputAttributesInclusive.
	isOutputElement := e.isVisibleDO(el, ns.level()) == 1
	var visiblyUtilized map[string]struct{}
	if isOutputElement {
		visiblyUtilized = make(map[string]struct{}, len(el.Attrs)+len(e.inclusivePrefixes)+1)
		for _, prefix := range e.inclusivePrefixes {
			visiblyUtilized[prefix] = struct{}{}
		}
	}

	for _, attr := range el.Attrs {
		name, value := attr.Name.Local, attr.Value
		switch {
		case attr.Name.Space != xmldom.XMLNSNamespace:
			if e.isVisible(attr) && isOutputElement {
				if prefix := attr.Name.Prefix; prefix != "" && prefix != xmlPrefix && prefix != xmlnsPrefix {
					visiblyUtilized[prefix] = struct{}{}
				}
				result.add(attrOf(attr))
			}
		case isOutputElement && !e.isVisible(attr) && name != xmlnsPrefix:
			ns.removeMappingIfNotRender(name)
		default:
			if !isOutputElement && e.isVisible(attr) && e.inPrefixList(name) && !ns.removeMappingIfRender(name) {
				if rendered, ok := ns.addMappingAndRender(name, value, true); ok {
					result.add(rendered)
					if namespaceIsRelative(value) {
						return &RelativeNamespaceError{Element: el.Name.QName(), Prefix: name, URI: value}
					}
				}
			}
			if ns.addMapping(name, value, true) && namespaceIsRelative(value) {
				return &RelativeNamespaceError{Element: el.Name.QName(), Prefix: name, URI: value}
			}
		}
	}

	if isOutputElement {
		if xmlns := el.Attr(xmldom.XMLNSNamespace, xmlnsPrefix); xmlns != nil && !e.isVisible(xmlns) {
			ns.addMapping(xmlnsPrefix, "", true)
		}
		prefix := xmlnsPrefix
		if el.Name.Space != "" && el.Name.Prefix != "" {
			prefix = el.Name.Prefix
		}
		visiblyUtilized[prefix] = struct{}{}
		for _, s := range sortedKeys(visiblyUtilized) {
			if rendered, ok := ns.getMapping(s); ok {
				result.add(rendered)
			}
		}
	}

	result.writeTo(e.w)
	return nil
}

func (e *engine) inPrefixList(prefix string) bool {
	for _, p := range e.inclusivePrefixes {
		if p == prefix {
			return true
		}
	}
	return false
}

// sortedKeys iterates the visibly-utilized set in the order Santuario's TreeSet<String> does.
// The rendered declarations land in an attrSet that sorts them again, so the order only
// decides which of two equal-comparing declarations wins - it cannot differ - but iterating a
// Go map here would still be a source of nondeterminism, so it is sorted.
func sortedKeys(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
