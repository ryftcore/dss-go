package xmldom

import (
	"sort"
	"strings"
)

// RegisterIDs indexes ID attributes over the whole document, reproducing
// XAdESDOMDocument.recursiveIdBrowse: per element, in attribute order, the first
// attribute whose local name equals "Id" case-insensitively; plus any xml:id.
// It must be called on a document node.
//
// The DSS rule is a scan-and-break, exactly as setIDIdentifier writes it: the first
// case-insensitive "Id" match ends the scan for that element, so a later Id-like
// attribute on the same element is not an ID. Note that xml:id has local name "id"
// and therefore also satisfies the DSS rule; when it comes first it consumes the
// element's one DSS registration.
func (n *Node) RegisterIDs() {
	d := n.mustDocState("RegisterIDs")
	d.scanAll = true
	n.rebuildIDs()
}

// RegisterIDAttr indexes a single attribute of elem as an ID attribute.
//
// This is the equivalent of Element.setIdAttribute(name, true). Manual registrations
// survive index invalidation: they are recorded and replayed after any lazy rebuild.
func (n *Node) RegisterIDAttr(elem, attr *Node) {
	d := n.mustDocState("RegisterIDAttr")
	if elem == nil || attr == nil {
		return
	}
	d.manual = append(d.manual, manualID{elem: elem, attr: attr})
	if !d.built || d.idGen != d.gen {
		n.rebuildIDs()
		return
	}
	d.register(elem, attr)
}

// ElementByID returns the element carrying the given ID value, or nil. A leading '#'
// is not accepted; strip the fragment delimiter before calling.
func (n *Node) ElementByID(id string) *Node {
	d := n.syncIDs()
	if d == nil {
		return nil
	}
	return d.ids[id]
}

// IDAttrs returns the ID attribute nodes registered for elem, in attribute order.
func (n *Node) IDAttrs(elem *Node) []*Node {
	d := n.syncIDs()
	if d == nil {
		return nil
	}
	got := d.idAttrs[elem]
	if len(got) == 0 {
		return nil
	}
	out := make([]*Node, len(got))
	copy(out, got)
	return out
}

// DuplicateIDs returns the ID values registered more than once, sorted.
//
// Registration itself never fails: DSSXMLUtils.isDuplicateIdsDetected reports the
// condition and leaves the caller to decide whether it is fatal, because in XAdES
// validation a duplicate Id is an attack indicator rather than a parse error.
func (n *Node) DuplicateIDs() []string {
	d := n.syncIDs()
	if d == nil || len(d.dups) == 0 {
		return nil
	}
	out := make([]string, 0, len(d.dups))
	for v := range d.dups {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

// mustDocState returns the document state, panicking when n is not a document node.
func (n *Node) mustDocState(who string) *docState {
	if n == nil || n.Kind != Document || n.doc == nil {
		panic("xmldom: " + who + " must be called on a document node")
	}
	return n.doc
}

// syncIDs returns the document state with the index up to date, or nil when n is not
// a document node.
func (n *Node) syncIDs() *docState {
	if n == nil || n.Kind != Document || n.doc == nil {
		return nil
	}
	if !n.doc.built || n.doc.idGen != n.doc.gen {
		n.rebuildIDs()
	}
	return n.doc
}

// rebuildIDs discards the index and recomputes it from the current tree.
func (n *Node) rebuildIDs() {
	d := n.doc
	d.ids = make(map[string]*Node)
	d.idAttrs = make(map[*Node][]*Node)
	d.dups = make(map[string]struct{})
	if d.scanAll {
		n.Walk(func(x *Node) bool {
			if x.Kind == Element {
				d.scanElement(x)
			}
			return true
		})
	}
	for _, m := range d.manual {
		d.register(m.elem, m.attr)
	}
	d.built = true
	d.idGen = d.gen
}

// scanElement applies the DSS rule and the xml:id rule to one element.
func (d *docState) scanElement(el *Node) {
	for _, a := range el.Attrs {
		if strings.EqualFold(a.Name.Local, "Id") {
			d.register(el, a)
			break
		}
	}
	for _, a := range el.Attrs {
		if a.Name.Space == XMLNamespace && a.Name.Local == "id" {
			d.register(el, a)
		}
	}
}

// register indexes attr as an ID attribute of elem. Registering the same attribute
// twice is a no-op, so the DSS rule and the xml:id rule agreeing on one attribute
// does not manufacture a duplicate.
func (d *docState) register(elem, attr *Node) {
	for _, have := range d.idAttrs[elem] {
		if have == attr {
			return
		}
	}
	d.idAttrs[elem] = insertInAttrOrder(elem, d.idAttrs[elem], attr)
	if _, seen := d.ids[attr.Value]; seen {
		d.dups[attr.Value] = struct{}{}
		return // first registration in document order wins
	}
	d.ids[attr.Value] = elem
}

// insertInAttrOrder inserts attr into got so that got stays ordered by the attribute's
// position in elem.Attrs.
func insertInAttrOrder(elem *Node, got []*Node, attr *Node) []*Node {
	pos := attrIndex(elem, attr)
	i := sort.Search(len(got), func(k int) bool { return attrIndex(elem, got[k]) > pos })
	got = append(got, nil)
	copy(got[i+1:], got[i:])
	got[i] = attr
	return got
}

// attrIndex returns attr's position in elem.Attrs, or len(elem.Attrs) when it is not
// one of them (a manually registered foreign attribute sorts last).
func attrIndex(elem, attr *Node) int {
	for i, a := range elem.Attrs {
		if a == attr {
			return i
		}
	}
	return len(elem.Attrs)
}
