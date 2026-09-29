// Ported from org.apache.xml.security.c14n.implementations.XmlAttrStack
// (Apache Santuario xmlsec 3.0.6).
package xmlc14n

import (
	"errors"
	"sort"
)

// xmlAttrLevel is XmlAttrStack.XmlsStackElement.
type xmlAttrLevel struct {
	level    int
	rendered bool
	nodes    []outAttr
}

// xmlAttrStack collects the xml:* attributes of elements that are not themselves output -
// the ancestors of a subtree apex, or the unselected ancestors of a node set - and flushes
// them onto the first element that is. Only inclusive c14n has one; exclusive and physical
// never inherit xml:* attributes.
type xmlAttrStack struct {
	currentLevel int
	lastLevel    int
	cur          *xmlAttrLevel
	levels       []*xmlAttrLevel

	// c14n11 selects the Canonical XML 1.1 behaviour: xml:id is not inherited (handled by the
	// emitters and by handleParent) and the unrendered ancestors' xml:base values are joined
	// with joinURI rather than first-one-wins.
	c14n11 bool
}

func (s *xmlAttrStack) push(level int) {
	s.currentLevel = level
	if s.currentLevel == -1 {
		return
	}
	s.cur = nil
	for s.lastLevel >= s.currentLevel {
		s.levels = s.levels[:len(s.levels)-1]
		if len(s.levels) == 0 {
			s.lastLevel = 0
			return
		}
		s.lastLevel = s.levels[len(s.levels)-1].level
	}
}

func (s *xmlAttrStack) addXmlnsAttr(a outAttr) {
	if s.cur == nil {
		s.cur = &xmlAttrLevel{level: s.currentLevel}
		s.levels = append(s.levels, s.cur)
		s.lastLevel = s.currentLevel
	}
	s.cur.nodes = append(s.cur.nodes, a)
}

// getXmlnsAttr floats the collected xml:* attributes onto the element being emitted.
//
// The C14N 1.0 rule is not "nearest ancestor wins": handleParent pushes level -1 for every
// ancestor of a subtree apex, so all of them share one level, their attributes are appended
// outermost-first, and the loa map keyed by QName keeps the FIRST occurrence. The outermost
// xml:lang therefore wins, which the xmlattrs-lang-space known-answer test pins.
func (s *xmlAttrStack) getXmlnsAttr(col *attrSet) error {
	size := len(s.levels) - 1
	if s.cur == nil {
		s.cur = &xmlAttrLevel{level: s.currentLevel}
		s.lastLevel = s.currentLevel
		s.levels = append(s.levels, s.cur)
	}
	parentRendered := false
	if size == -1 {
		parentRendered = true
	} else {
		e := s.levels[size]
		if e.rendered && e.level+1 == s.currentLevel {
			parentRendered = true
		}
	}
	if parentRendered {
		for _, a := range s.cur.nodes {
			col.add(a)
		}
		s.cur.rendered = true
		return nil
	}

	loa := make(map[string]outAttr)
	if s.c14n11 {
		// C14N 1.1 walks the omitted levels from innermost outwards and stops at the first
		// one already rendered: successiveOmitted goes false at the top of that level and
		// gates the inner loop, so a rendered level contributes nothing at all and neither
		// does anything outside it. xml:base is collected separately, to be joined; every
		// other xml:* attribute keeps the first-occurrence-wins rule.
		var baseAttrs []outAttr
		successiveOmitted := true
		for ; size >= 0; size-- {
			e := s.levels[size]
			if e.rendered {
				successiveOmitted = false
			}
			for _, a := range e.nodes {
				if !successiveOmitted {
					break
				}
				if a.local == "base" && !e.rendered {
					baseAttrs = append(baseAttrs, a)
				} else if _, seen := loa[a.qname]; !seen {
					loa[a.qname] = a
				}
			}
		}
		if len(baseAttrs) > 0 {
			// The element's own xml:base, already in col, seeds the accumulator; otherwise
			// the outermost ancestor's does. Each further-in ancestor is then folded in as
			// joinURI(inner, accumulated) - the accumulated value is the RELATIVE argument.
			// Santuario searches col by local name alone, so a plain, unprefixed base="..."
			// attribute is picked up as though it were xml:base; that is reproduced, not
			// fixed, and pinned by the xmlbase-plain-base known-answer tests.
			base, hasBase := "", false
			baseIdx := -1
			var baseAttr outAttr
			items := col.sorted()
			for i, a := range items {
				if a.local == "base" {
					base, hasBase, baseIdx, baseAttr = a.value, true, i, a
					break
				}
			}
			for _, n := range baseAttrs {
				if !hasBase {
					base, hasBase, baseAttr = n.value, true, n
					continue
				}
				joined, err := joinURI(n.value, true, base)
				switch {
				case err == nil:
					base = joined
				case errors.Is(err, ErrXMLBaseUnjoinable):
					// Unchecked in Java: it escapes the catch and kills the run.
					return err
				default:
					// URISyntaxException: logged at debug, previous value kept.
				}
			}
			if hasBase && base != "" {
				baseAttr.value = base
				if baseIdx >= 0 {
					items[baseIdx].value = base // Java mutates the Attr in place.
				} else {
					col.add(baseAttr)
				}
			}
		}
	} else {
		for ; size >= 0; size-- {
			for _, a := range s.levels[size].nodes {
				if _, seen := loa[a.qname]; !seen {
					loa[a.qname] = a
				}
			}
		}
	}
	s.cur.rendered = true
	// Java hands a HashMap's values to a TreeSet, so its iteration order cannot reach the
	// output; sorting the keys here keeps that true without relying on it.
	names := make([]string, 0, len(loa))
	for name := range loa {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		col.add(loa[name])
	}
	return nil
}
