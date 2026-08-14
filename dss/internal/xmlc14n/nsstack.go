// Ported from org.apache.xml.security.c14n.implementations.NameSpaceSymbTable and its
// NameSpaceSymbEntry (Apache Santuario xmlsec 3.0.6).
package xmlc14n

import "sort"

// nsEntry is NameSpaceSymbEntry. lastRendered is the field that suppresses superfluous
// declarations and the redundant xmlns="": it remembers the URI last actually written for the
// prefix, across frames, so redeclaring a prefix to a URI that is already in force renders
// nothing. It is load-bearing and must not be "simplified" away.
type nsEntry struct {
	prefix string
	uri    string

	// hasNode is Santuario's "n != null": every mapping added during canonicalization
	// carries the declaration attribute that would be written for it. Only the initial
	// xmlns entry has none, and it can therefore never be emitted.
	hasNode bool

	lastRendered    string
	hasLastRendered bool // Java's lastrendered == null
	rendered        bool
}

func (e *nsEntry) clone() *nsEntry {
	c := *e
	return &c
}

// attr is the declaration this entry would render.
func (e *nsEntry) attr() outAttr { return nsAttr(e.prefix, e.uri) }

// nsStack is NameSpaceSymbTable: a stack of frames over one prefix -> entry map.
//
// Santuario stores the frames as copy-on-write clones of an open-addressing table; that is an
// optimization, and copying it would make output depend on hash order. This port keeps a plain
// map plus an explicit undo journal per frame, which is the same semantics - every write made
// inside a frame is reverted by its pop - with deterministic iteration. Nothing derived from
// map order reaches the output: getUnrenderedNodes feeds a sorted attrSet, and it walks its
// prefixes in sorted order anyway.
type nsStack struct {
	symb   map[string]*nsEntry
	frames []map[string]*nsEntry // prefix -> value before this frame first wrote it
}

func newNSStack() *nsStack {
	// The initial map holds the default binding for xmlns: uri "", already rendered, with
	// lastrendered "". That entry is why an apex which undeclares the default namespace emits
	// no xmlns="" (design section 2.4).
	initial := &nsEntry{
		prefix:          xmlnsPrefix,
		uri:             "",
		hasNode:         false,
		rendered:        true,
		lastRendered:    "",
		hasLastRendered: true,
	}
	return &nsStack{symb: map[string]*nsEntry{xmlnsPrefix: initial}}
}

// put records the pre-image of prefix in the current frame and then writes value. A nil value
// is Santuario's symb.put(prefix, null), which is indistinguishable from "absent" to every
// reader, so the undo journal restores by deleting.
func (s *nsStack) put(prefix string, value *nsEntry) {
	if n := len(s.frames); n > 0 {
		frame := s.frames[n-1]
		if _, recorded := frame[prefix]; !recorded {
			frame[prefix] = s.symb[prefix] // nil when absent
		}
	}
	if value == nil {
		delete(s.symb, prefix)
		return
	}
	s.symb[prefix] = value
}

func (s *nsStack) get(prefix string) *nsEntry { return s.symb[prefix] }

// push and pop are NameSpaceSymbTable.push/pop; outputNodePush and outputNodePop are the
// aliases the inclusive traversal uses for a visible element.
func (s *nsStack) push() {
	s.frames = append(s.frames, make(map[string]*nsEntry, 4))
}

func (s *nsStack) pop() {
	n := len(s.frames) - 1
	if n < 0 {
		return
	}
	for prefix, prev := range s.frames[n] {
		if prev == nil {
			delete(s.symb, prefix)
			continue
		}
		s.symb[prefix] = prev
	}
	s.frames = s.frames[:n]
}

func (s *nsStack) outputNodePush() { s.push() }
func (s *nsStack) outputNodePop()  { s.pop() }

func (s *nsStack) level() int { return len(s.frames) }

// getUnrenderedNodes adds every mapping that has not been rendered yet and marks it rendered.
// It is how inclusive c14n flushes the ancestors' declarations onto the apex.
//
// Santuario's SymbMap.entrySet() skips entries whose URI is empty, so an inherited xmlns=""
// is never flushed here; that filter is reproduced below and is the second half of the
// apex-undeclares-the-default-namespace rule.
func (s *nsStack) getUnrenderedNodes(result *attrSet) {
	prefixes := make([]string, 0, len(s.symb))
	for prefix := range s.symb {
		prefixes = append(prefixes, prefix)
	}
	sort.Strings(prefixes)
	for _, prefix := range prefixes {
		e := s.symb[prefix]
		if e == nil || e.uri == "" {
			continue
		}
		if e.rendered || !e.hasNode {
			continue
		}
		c := e.clone()
		c.lastRendered = c.uri
		c.hasLastRendered = true
		c.rendered = true
		s.put(prefix, c)
		result.add(c.attr())
	}
}

// getMapping returns the declaration to render for prefix and marks it rendered, or reports
// false when there is nothing to render - either the prefix is unbound (an unknown prefix in
// an InclusiveNamespaces PrefixList lands here) or it is already rendered with this URI.
func (s *nsStack) getMapping(prefix string) (outAttr, bool) {
	e := s.get(prefix)
	if e == nil || e.rendered {
		return outAttr{}, false
	}
	c := e.clone()
	c.rendered = true
	c.lastRendered = c.uri
	c.hasLastRendered = true
	s.put(prefix, c)
	if !c.hasNode {
		return outAttr{}, false
	}
	return c.attr(), true
}

// getMappingWithoutRendered is getMapping without the side effect. It also reports the bound
// URI, which getParentNameSpaces needs to detect an inherited xmlns="".
func (s *nsStack) getMappingWithoutRendered(prefix string) (uri string, ok bool) {
	e := s.get(prefix)
	if e == nil || e.rendered {
		return "", false
	}
	return e.uri, true
}

// addMapping records a binding without rendering it - exclusive c14n's path - and reports
// whether it was a new definition. The lastRendered inheritance is what makes a redeclaration
// of an already-rendered binding count as rendered.
func (s *nsStack) addMapping(prefix, uri string, hasNode bool) bool {
	ob := s.get(prefix)
	if ob != nil && ob.uri == uri {
		return false
	}
	ne := &nsEntry{prefix: prefix, uri: uri, hasNode: hasNode}
	if ob != nil {
		ne.lastRendered, ne.hasLastRendered = ob.lastRendered, ob.hasLastRendered
		if ob.hasLastRendered && ob.lastRendered == uri {
			ne.rendered = true
		}
	}
	s.put(prefix, ne)
	return true
}

// addMappingAndRender is inclusive c14n's path: record the binding and report the declaration
// to emit, or false when emitting it would be superfluous.
func (s *nsStack) addMappingAndRender(prefix, uri string, hasNode bool) (outAttr, bool) {
	ob := s.get(prefix)
	if ob != nil && ob.uri == uri {
		if !ob.rendered {
			c := ob.clone()
			c.lastRendered = uri
			c.hasLastRendered = true
			c.rendered = true
			s.put(prefix, c)
			if !c.hasNode {
				return outAttr{}, false
			}
			return c.attr(), true
		}
		return outAttr{}, false
	}
	ne := &nsEntry{
		prefix:          prefix,
		uri:             uri,
		hasNode:         hasNode,
		rendered:        true,
		lastRendered:    uri,
		hasLastRendered: true,
	}
	s.put(prefix, ne)
	if ob != nil && ob.hasLastRendered && ob.lastRendered == uri {
		// Already in force from an outer, still-rendered declaration: nothing to write.
		return outAttr{}, false
	}
	if !ne.hasNode {
		return outAttr{}, false
	}
	return ne.attr(), true
}

// removeMapping, removeMappingIfNotRender and removeMappingIfRender are the node-set paths for
// a declaration that the XPath selection excludes. removeMappingIfRender always reports false
// upstream - the return value is a Santuario bug that call sites depend on - so it is
// reproduced verbatim.
func (s *nsStack) removeMapping(prefix string) {
	if s.get(prefix) != nil {
		s.put(prefix, nil)
	}
}

func (s *nsStack) removeMappingIfNotRender(prefix string) {
	if e := s.get(prefix); e != nil && !e.rendered {
		s.put(prefix, nil)
	}
}

func (s *nsStack) removeMappingIfRender(prefix string) bool {
	if e := s.get(prefix); e != nil && e.rendered {
		s.put(prefix, nil)
	}
	return false
}
