// Ported from org.apache.xml.security.c14n.implementations.Canonicalizer20010315
// (Apache Santuario xmlsec 3.0.6): attribute emission for Canonical XML 1.0 and 1.1.
package xmlc14n

import "github.com/utain/esig/dss/internal/xmldom"

// outputAttributesSubtreeInclusive ports Canonicalizer20010315.outputAttributesSubtree.
//
// Every declaration in scope at the apex is inherited, and each element below it emits only
// the bindings newly rendered there: addMappingAndRender reports nothing to write when the
// prefix is already bound to the same URI and already rendered, which is the
// superfluous-declaration rule. Declared-but-unused prefixes ARE emitted - that is the
// difference exclusive c14n exists to remove.
func (e *engine) outputAttributesSubtreeInclusive(el *xmldom.Node, ns *nsStack) error {
	if len(el.Attrs) == 0 && !e.firstCall {
		return nil
	}
	var result attrSet
	for _, attr := range el.Attrs {
		name, value := attr.Name.Local, attr.Value
		if attr.Name.Space != xmldom.XMLNSNamespace {
			result.add(attrOf(attr))
			continue
		}
		// Omit the namespace node with local name xml that defines the xml prefix, if its
		// value is http://www.w3.org/XML/1998/namespace.
		if name == xmlPrefix && value == xmldom.XMLNamespace {
			continue
		}
		if rendered, ok := ns.addMappingAndRender(name, value, true); ok {
			result.add(rendered)
			// The relative-URI check fires only for a declaration that is both physically
			// present on this element and actually rendered. An inherited relative
			// declaration flushed by getUnrenderedNodes below is never checked, and neither
			// is an unrendered one. That is Santuario's behaviour, pinned by the ns-relative
			// KATs; do not hoist the check.
			if namespaceIsRelative(value) {
				return &RelativeNamespaceError{Element: el.Name.QName(), Prefix: name, URI: value}
			}
		}
	}
	if e.firstCall {
		// The apex: flush the ancestors' declarations and their xml:* attributes.
		ns.getUnrenderedNodes(&result)
		if err := e.xmlAttrs.getXmlnsAttr(&result); err != nil {
			return err
		}
		e.firstCall = false
	}
	result.writeTo(e.w)
	return nil
}

// outputAttributesInclusive ports Canonicalizer20010315.outputAttributes, the node-set form.
// An element that is not itself in the node set still contributes its declarations and its
// xml:* attributes to its visible descendants, and a declaration excluded from the node set
// can force an xmlns="" onto a visible element.
func (e *engine) outputAttributesInclusive(el *xmldom.Node, ns *nsStack) error {
	e.xmlAttrs.push(ns.level())
	isRealVisible := e.isVisible(el)
	var result attrSet

	for _, attr := range el.Attrs {
		name, value := attr.Name.Local, attr.Value
		if attr.Name.Space != xmldom.XMLNSNamespace {
			if attr.Name.Space == xmldom.XMLNamespace {
				if e.c14n11 && name == "id" {
					// C14N 1.1 treats xml:id as an ordinary attribute: emitted when the
					// element carries it, never inherited.
					if isRealVisible {
						result.add(attrOf(attr))
					}
				} else {
					e.xmlAttrs.addXmlnsAttr(attrOf(attr))
				}
			} else if isRealVisible {
				result.add(attrOf(attr))
			}
			continue
		}
		if name == xmlPrefix && value == xmldom.XMLNamespace {
			continue
		}
		if e.isVisible(attr) {
			if isRealVisible || !ns.removeMappingIfRender(name) {
				if rendered, ok := ns.addMappingAndRender(name, value, true); ok {
					result.add(rendered)
					if namespaceIsRelative(value) {
						return &RelativeNamespaceError{Element: el.Name.QName(), Prefix: name, URI: value}
					}
				}
			}
		} else if isRealVisible && name != xmlnsPrefix {
			ns.removeMapping(name)
		} else {
			ns.addMapping(name, value, true)
		}
	}

	if isRealVisible {
		xmlns := el.Attr(xmldom.XMLNSNamespace, xmlnsPrefix)
		var rendered outAttr
		ok := false
		if xmlns == nil {
			rendered, ok = ns.getMapping(xmlnsPrefix)
		} else if !e.isVisible(xmlns) {
			// The element declares a default namespace that the node set excludes, so the
			// canonical form has to undeclare it.
			rendered, ok = ns.addMappingAndRender(xmlnsPrefix, "", true)
		}
		if ok {
			result.add(rendered)
		}
		if err := e.xmlAttrs.getXmlnsAttr(&result); err != nil {
			return err
		}
		ns.getUnrenderedNodes(&result)
	}

	result.writeTo(e.w)
	return nil
}
