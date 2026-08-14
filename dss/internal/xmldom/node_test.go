package xmldom

import (
	"strings"
	"testing"
)

func TestKindString(t *testing.T) {
	for k, want := range map[Kind]string{
		Document: "Document", Element: "Element", Attribute: "Attribute",
		Text: "Text", CDATA: "CDATA", Comment: "Comment", ProcInst: "ProcInst",
		Kind(0): "Kind(0)", Kind(99): "Kind(99)",
	} {
		if got := k.String(); got != want {
			t.Errorf("Kind(%d).String() = %q, want %q", k, got, want)
		}
	}
}

func TestQName(t *testing.T) {
	for _, tc := range []struct {
		n    Name
		want string
	}{
		{Name{Local: "a"}, "a"},
		{Name{Prefix: "p", Local: "a"}, "p:a"},
		{Name{Space: XMLNSNamespace, Local: "xmlns"}, "xmlns"},
		{Name{Space: XMLNSNamespace, Local: "p", Prefix: "xmlns"}, "xmlns:p"},
	} {
		if got := tc.n.QName(); got != tc.want {
			t.Errorf("%+v.QName() = %q, want %q", tc.n, got, tc.want)
		}
	}
}

func TestBuildTreeByHand(t *testing.T) {
	doc := NewDocument()
	root := NewElement(Name{Space: "urn:a", Local: "r", Prefix: "a"})
	doc.AppendChild(NewComment("lead"))
	doc.AppendChild(root)
	doc.AppendChild(NewProcInst("trail", "x"))
	root.SetAttr(Name{Space: XMLNSNamespace, Local: "a", Prefix: "xmlns"}, "urn:a")
	root.SetAttr(Name{Local: "b"}, "1")
	root.AppendChild(NewText("t"))
	root.AppendChild(NewCDATA("c"))

	want := `Document
  Comment = "lead"
  Element a:r {urn:a}r
    Attribute xmlns:a {http://www.w3.org/2000/xmlns/}a = "urn:a"
    Attribute b {}b = "1"
    Text = "t"
    CDATA = "c"
  ProcInst trail = "x"
`
	if got := dump(doc); got != want {
		t.Errorf("tree =\n%s\nwant\n%s", got, want)
	}
	if got, want := mustSerialize(t, doc), `<!--lead--><a:r xmlns:a="urn:a" b="1">t<![CDATA[c]]></a:r><?trail x?>`; got != want {
		t.Errorf("serialized = %q, want %q", got, want)
	}
}

func TestInsertBeforeAndRemove(t *testing.T) {
	root := NewElement(Name{Local: "r"})
	a, b, c := NewText("a"), NewText("b"), NewText("c")
	root.AppendChild(a)
	root.AppendChild(c)
	root.InsertBefore(b, c)

	if got := root.TextContent(); got != "abc" {
		t.Errorf("after InsertBefore = %q, want %q", got, "abc")
	}
	if root.FirstChild != a || root.LastChild != c {
		t.Error("FirstChild/LastChild are wrong")
	}
	if b.PrevSibling != a || b.NextSibling != c || c.PrevSibling != b {
		t.Error("sibling links are wrong")
	}

	if got := root.RemoveChild(b); got != b {
		t.Error("RemoveChild returned the wrong node")
	}
	if b.Parent != nil || b.PrevSibling != nil || b.NextSibling != nil {
		t.Error("removed node still has links")
	}
	if got := root.TextContent(); got != "ac" {
		t.Errorf("after RemoveChild = %q, want %q", got, "ac")
	}
	root.RemoveChild(a)
	root.RemoveChild(c)
	if root.FirstChild != nil || root.LastChild != nil {
		t.Error("emptied element still has child links")
	}
}

func TestReplaceChild(t *testing.T) {
	root := NewElement(Name{Local: "r"})
	a, b, c := NewText("a"), NewText("b"), NewText("c")
	root.AppendChild(a)
	root.AppendChild(b)
	root.AppendChild(c)

	nb := NewText("B")
	if got := root.ReplaceChild(nb, b); got != b {
		t.Error("ReplaceChild should return the old node")
	}
	if got := root.TextContent(); got != "aBc" {
		t.Errorf("after ReplaceChild = %q, want %q", got, "aBc")
	}
	if b.Parent != nil {
		t.Error("replaced node still has a parent")
	}

	// Replacing at the edges must maintain FirstChild/LastChild.
	root.ReplaceChild(NewText("A"), a)
	root.ReplaceChild(NewText("C"), c)
	if got := root.TextContent(); got != "ABC" {
		t.Errorf("after edge replacements = %q, want %q", got, "ABC")
	}

	// Replacing a document's root element must not trip the one-element-child rule.
	doc := NewDocument()
	old := NewElement(Name{Local: "old"})
	doc.AppendChild(old)
	doc.ReplaceChild(NewElement(Name{Local: "new"}), old)
	if got := doc.DocumentElement().Name.Local; got != "new" {
		t.Errorf("root = %q, want %q", got, "new")
	}
}

func TestMutationPanics(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
	}{
		{"attribute as a child", func() { NewElement(Name{Local: "r"}).AppendChild(NewAttr(Name{Local: "a"}, "1")) }},
		{"text under a document", func() { NewDocument().AppendChild(NewText("x")) }},
		{"element under a text node", func() { NewText("x").AppendChild(NewElement(Name{Local: "r"})) }},
		{"document under a document", func() { NewDocument().AppendChild(NewDocument()) }},
		{"two root elements", func() {
			d := NewDocument()
			d.AppendChild(NewElement(Name{Local: "a"}))
			d.AppendChild(NewElement(Name{Local: "b"}))
		}},
		{"node that already has a parent", func() {
			p := NewElement(Name{Local: "p"})
			c := NewElement(Name{Local: "c"})
			p.AppendChild(c)
			NewElement(Name{Local: "q"}).AppendChild(c)
		}},
		{"cycle", func() {
			p := NewElement(Name{Local: "p"})
			c := NewElement(Name{Local: "c"})
			p.AppendChild(c)
			c.AppendChild(p)
		}},
		{"self insertion", func() {
			p := NewElement(Name{Local: "p"})
			p.AppendChild(p)
		}},
		{"ref is not a child", func() {
			p := NewElement(Name{Local: "p"})
			p.InsertBefore(NewText("x"), NewText("y"))
		}},
		{"remove a non-child", func() {
			NewElement(Name{Local: "p"}).RemoveChild(NewText("x"))
		}},
		{"replace a non-child", func() {
			NewElement(Name{Local: "p"}).ReplaceChild(NewText("x"), NewText("y"))
		}},
		{"SetAttr on a non-element", func() { NewText("x").SetAttr(Name{Local: "a"}, "1") }},
		{"SetTextContent on a non-element", func() { NewText("x").SetTextContent("y") }},
		{"RegisterIDs on a non-document", func() { NewElement(Name{Local: "r"}).RegisterIDs() }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Error("expected a panic")
				}
			}()
			tc.fn()
		})
	}
}

func TestAttrAccessors(t *testing.T) {
	el := mustParseRoot(t, `<r xmlns:p="urn:1" a="1" p:b="2" xmlns="urn:d"/>`)

	if got := el.AttrValue("", "a"); got != "1" {
		t.Errorf("a = %q, want %q", got, "1")
	}
	if got := el.AttrValue("urn:1", "b"); got != "2" {
		t.Errorf("p:b = %q, want %q", got, "2")
	}
	if got := el.AttrValue(XMLNSNamespace, "p"); got != "urn:1" {
		t.Errorf("xmlns:p = %q, want %q", got, "urn:1")
	}
	if got := el.AttrValue(XMLNSNamespace, "xmlns"); got != "urn:d" {
		t.Errorf("xmlns = %q, want %q", got, "urn:d")
	}
	if got := el.AttrValue("", "missing"); got != "" {
		t.Errorf("missing = %q, want %q", got, "")
	}
	if el.Attr("", "missing") != nil {
		t.Error("Attr for a missing attribute should be nil")
	}

	// SetAttr replaces in place and keeps the position.
	el.SetAttr(Name{Local: "a"}, "9")
	if got := el.AttrValue("", "a"); got != "9" {
		t.Errorf("after SetAttr a = %q, want %q", got, "9")
	}
	if got := el.Attrs[1].Name.QName(); got != "a" {
		t.Errorf("SetAttr moved the attribute; Attrs[1] = %q, want %q", got, "a")
	}
	if len(el.Attrs) != 4 {
		t.Errorf("SetAttr changed the attribute count to %d, want 4", len(el.Attrs))
	}

	// A new attribute is appended.
	el.SetAttr(Name{Local: "z"}, "0")
	if got := el.Attrs[len(el.Attrs)-1].Name.QName(); got != "z" {
		t.Errorf("new attribute was not appended; last = %q", got)
	}

	if !el.RemoveAttr("", "a") {
		t.Error("RemoveAttr reported nothing removed")
	}
	if el.Attr("", "a") != nil {
		t.Error("attribute survived RemoveAttr")
	}
	if el.RemoveAttr("", "a") {
		t.Error("RemoveAttr reported a second removal")
	}
}

func TestNavigation(t *testing.T) {
	doc := mustParse(t, `<r><a/>text<b><c/></b><!--x--></r>`)
	root := doc.DocumentElement()

	if got := len(root.Children()); got != 4 {
		t.Errorf("Children = %d, want 4", got)
	}
	if got := len(root.Elements()); got != 2 {
		t.Errorf("Elements = %d, want 2", got)
	}
	a := root.FirstElementChild()
	if a == nil || a.Name.Local != "a" {
		t.Fatalf("FirstElementChild = %v", a)
	}
	b := a.NextElementSibling()
	if b == nil || b.Name.Local != "b" {
		t.Fatalf("NextElementSibling = %v", b)
	}
	if b.NextElementSibling() != nil {
		t.Error("NextElementSibling past the last element should be nil")
	}
	c := b.FirstElementChild()

	if got, want := c.Depth(), 3; got != want {
		t.Errorf("Depth(c) = %d, want %d", got, want)
	}
	if got, want := doc.Depth(), 0; got != want {
		t.Errorf("Depth(document) = %d, want %d", got, want)
	}
	if got := len(c.Ancestors()); got != 3 {
		t.Errorf("Ancestors(c) = %d, want 3", got)
	}
	if c.Ancestors()[0] != b {
		t.Error("Ancestors must be nearest-first")
	}
	if !root.Contains(c) || !root.Contains(root) || root.Contains(doc) {
		t.Error("Contains is wrong")
	}
	if doc.DocumentElement() != root || root.DocumentElement() != nil {
		t.Error("DocumentElement is wrong")
	}
	if c.Document() != doc {
		t.Error("Document should reach the owning document")
	}
	if NewElement(Name{Local: "detached"}).Document() != nil {
		t.Error("a detached subtree has no document")
	}
}

func TestWalkOrderAndPruning(t *testing.T) {
	doc := mustParse(t, `<r><a><a1/></a><b><b1/></b></r>`)
	var seen []string
	doc.Walk(func(n *Node) bool {
		if n.Kind == Element {
			seen = append(seen, n.Name.Local)
		}
		return true
	})
	if got, want := strings.Join(seen, ","), "r,a,a1,b,b1"; got != want {
		t.Errorf("Walk order = %q, want %q", got, want)
	}

	seen = nil
	doc.Walk(func(n *Node) bool {
		if n.Kind == Element {
			seen = append(seen, n.Name.Local)
		}
		return !(n.Kind == Element && n.Name.Local == "a")
	})
	if got, want := strings.Join(seen, ","), "r,a,b,b1"; got != want {
		t.Errorf("pruned Walk = %q, want %q", got, want)
	}

	// Returning false at the root visits nothing else.
	count := 0
	doc.Walk(func(*Node) bool { count++; return false })
	if count != 1 {
		t.Errorf("Walk visited %d nodes after a root prune, want 1", count)
	}
}

func TestTextContentAndSetTextContent(t *testing.T) {
	root := mustParseRoot(t, `<r>a<b>c<![CDATA[d]]></b><!--ignored--><?pi ignored?>e</r>`)
	if got, want := root.TextContent(), "acde"; got != want {
		t.Errorf("TextContent = %q, want %q", got, want)
	}
	if got := root.Attrs; got != nil {
		t.Errorf("unexpected attributes %v", got)
	}

	root.SetTextContent("only")
	if got, want := dump(root), "Element r {}r\n  Text = \"only\"\n"; got != want {
		t.Errorf("after SetTextContent =\n%s\nwant\n%s", got, want)
	}
	root.SetTextContent("")
	if root.FirstChild != nil || root.LastChild != nil {
		t.Error(`SetTextContent("") should leave the element childless`)
	}

	// Leaf kinds report their own value.
	if got := NewText("t").TextContent(); got != "t" {
		t.Errorf("Text.TextContent = %q", got)
	}
	if got := NewCDATA("c").TextContent(); got != "c" {
		t.Errorf("CDATA.TextContent = %q", got)
	}
	if got := NewAttr(Name{Local: "a"}, "v").TextContent(); got != "v" {
		t.Errorf("Attribute.TextContent = %q", got)
	}
}

func TestLookupNamespaceURI(t *testing.T) {
	doc := mustParse(t, `<r xmlns:p="urn:1" xmlns="urn:d"><m xmlns:p="urn:2"><t xmlns=""/></m></r>`)
	root := doc.DocumentElement()
	m := root.FirstElementChild()
	tn := m.FirstElementChild()

	for _, tc := range []struct {
		node   *Node
		prefix string
		uri    string
		ok     bool
	}{
		{root, "p", "urn:1", true},
		{m, "p", "urn:2", true},
		{tn, "p", "urn:2", true},
		{root, "", "urn:d", true},
		{m, "", "urn:d", true},
		{tn, "", "", true}, // explicitly undeclared, and that is a binding
		{root, "xml", XMLNamespace, true},
		{root, "xmlns", XMLNSNamespace, true},
		{root, "nope", "", false},
	} {
		uri, ok := tc.node.LookupNamespaceURI(tc.prefix)
		if uri != tc.uri || ok != tc.ok {
			t.Errorf("<%s>.LookupNamespaceURI(%q) = (%q, %v), want (%q, %v)",
				tc.node.Name.Local, tc.prefix, uri, ok, tc.uri, tc.ok)
		}
	}

	// An element with no declaration anywhere reports no default binding.
	if uri, ok := mustParseRoot(t, `<r/>`).LookupNamespaceURI(""); ok || uri != "" {
		t.Errorf(`LookupNamespaceURI("") on an undecorated root = (%q, %v), want ("", false)`, uri, ok)
	}
	// Attributes resolve through their element.
	a := mustParseRoot(t, `<r xmlns:p="urn:1" a="1"/>`).Attr("", "a")
	if uri, ok := a.LookupNamespaceURI("p"); !ok || uri != "urn:1" {
		t.Errorf("attribute lookup = (%q, %v), want (%q, true)", uri, ok, "urn:1")
	}
}

func TestLookupPrefix(t *testing.T) {
	doc := mustParse(t, `<r xmlns:p="urn:1" xmlns:q="urn:2" xmlns="urn:d"><m xmlns:p="urn:9"><t/></m></r>`)
	root := doc.DocumentElement()
	tn := root.FirstElementChild().FirstElementChild()

	for _, tc := range []struct {
		node   *Node
		uri    string
		prefix string
		ok     bool
	}{
		{root, "urn:1", "p", true},
		{root, "urn:2", "q", true},
		{tn, "urn:2", "q", true},
		{tn, "urn:9", "p", true},
		{tn, "urn:1", "", false},   // p is rebound closer to t, so urn:1 is unreachable
		{root, "urn:d", "", false}, // the default declaration is never a candidate
		{root, XMLNamespace, "xml", true},
		{root, "", "", false},
		{root, "urn:nope", "", false},
	} {
		prefix, ok := tc.node.LookupPrefix(tc.uri)
		if prefix != tc.prefix || ok != tc.ok {
			t.Errorf("<%s>.LookupPrefix(%q) = (%q, %v), want (%q, %v)",
				tc.node.Name.Local, tc.uri, prefix, ok, tc.prefix, tc.ok)
		}
	}
}

func TestCloneAndImport(t *testing.T) {
	doc := mustParse(t, `<r xmlns:p="urn:1" a="1"><c p:b="2">t</c></r>`)
	root := doc.DocumentElement()

	shallow := root.Clone(false)
	if shallow.FirstChild != nil {
		t.Error("a shallow clone must have no children")
	}
	if len(shallow.Attrs) != 2 {
		t.Errorf("a shallow clone has %d attributes, want 2", len(shallow.Attrs))
	}
	if shallow.Parent != nil || shallow.Document() != nil {
		t.Error("a clone must have no parent and no document")
	}
	if shallow.Attrs[0] == root.Attrs[0] {
		t.Error("a clone must not share attribute nodes with the original")
	}
	if shallow.Attrs[0].Parent != shallow {
		t.Error("a cloned attribute must point at its clone")
	}

	deep := root.Clone(true)
	if got, want := dump(deep), dump(root); got != want {
		t.Errorf("deep clone =\n%s\nwant\n%s", got, want)
	}
	// Mutating the clone must not touch the original.
	deep.FirstElementChild().SetTextContent("changed")
	if got := root.TextContent(); got != "t" {
		t.Errorf("the original changed to %q; the clone is not independent", got)
	}

	// A cloned document keeps its own ID index.
	docClone := doc.Clone(true)
	if docClone.doc == nil || docClone.doc == doc.doc {
		t.Error("a cloned document must get fresh document state")
	}

	// Import is Clone under this ownership model, and the result is insertable.
	imported := doc.Import(root.FirstElementChild(), true)
	if imported.Parent != nil {
		t.Error("an imported node must be detached")
	}
	root.AppendChild(imported)
	if got := len(root.Elements()); got != 2 {
		t.Errorf("after importing, root has %d element children, want 2", got)
	}

	if NewText("x").Clone(true) == nil {
		t.Error("cloning a leaf returned nil")
	}
	var nilNode *Node
	if nilNode.Clone(true) != nil {
		t.Error("cloning nil should return nil")
	}
}

func TestNodeSet(t *testing.T) {
	doc := mustParse(t, `<r a="1"><c b="2">t</c></r>`)
	root := doc.DocumentElement()
	c := root.FirstElementChild()

	s := NewNodeSet(root, c)
	if s.Len() != 2 || !s.Has(root) || !s.Has(c) {
		t.Errorf("NewNodeSet = %v", s)
	}
	s.Remove(c)
	if s.Has(c) || s.Len() != 1 {
		t.Error("Remove failed")
	}
	s.Add(nil)
	if s.Len() != 1 {
		t.Error("a nil node should be ignored")
	}

	full := NewNodeSet()
	full.AddSubtree(root)
	// r, its attribute, c, its attribute, the text node.
	if full.Len() != 5 {
		t.Errorf("AddSubtree produced %d members, want 5", full.Len())
	}
	if !full.Has(root.Attr("", "a")) || !full.Has(c.Attr("", "b")) {
		t.Error("AddSubtree must include attribute nodes")
	}
	if !full.Has(c.FirstChild) {
		t.Error("AddSubtree must include text nodes")
	}
	if full.Has(doc) {
		t.Error("AddSubtree must not walk upwards")
	}
}
