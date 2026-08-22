package xmldom

import (
	"reflect"
	"strings"
	"testing"
)

// idsOf renders the registered ID attributes of every element, so the DSS scan rule
// can be asserted as a whole.
func idsOf(doc *Node) []string {
	var out []string
	doc.Walk(func(n *Node) bool {
		if n.Kind != Element {
			return true
		}
		for _, a := range doc.IDAttrs(n) {
			out = append(out, n.Name.QName()+"/"+a.Name.QName()+"="+a.Value)
		}
		return true
	})
	return out
}

// TestRegisterIDsFollowsTheDSSRule pins DOMDocument.setIDIdentifier: per element,
// in attribute order, the FIRST attribute whose local name equals "Id" case
// insensitively, then break.
func TestRegisterIDsFollowsTheDSSRule(t *testing.T) {
	tests := []struct {
		name, src string
		want      []string
	}{{
		name: "case insensitive spellings all match",
		src:  `<r><a Id="1"/><b ID="2"/><c id="3"/><d iD="4"/></r>`,
		want: []string{"a/Id=1", "b/ID=2", "c/id=3", "d/iD=4"},
	}, {
		name: "a prefixed Id matches too",
		src:  `<r xmlns:p="urn:1"><a p:Id="1"/></r>`,
		want: []string{"a/p:Id=1"},
	}, {
		name: "only the first Id-like attribute on an element is registered",
		src:  `<r><a Id="1" ID="2" id="3"/></r>`,
		want: []string{"a/Id=1"},
	}, {
		name: "the scan respects attribute order",
		src:  `<r><a other="x" ID="2" Id="1"/></r>`,
		want: []string{"a/ID=2"},
	}, {
		name: "attributes that merely contain Id do not match",
		src:  `<r><a Ident="1" MyId="2" IdX="3"/></r>`,
		want: nil,
	}, {
		name: "xml:id is registered as well",
		src:  `<r><a Id="1" xml:id="2"/></r>`,
		want: []string{"a/Id=1", "a/xml:id=2"},
	}, {
		name: "xml:id first consumes the one DSS registration",
		src:  `<r><a xml:id="1" Id="2"/></r>`,
		want: []string{"a/xml:id=1"},
	}, {
		name: "xml:id alone",
		src:  `<r><a xml:id="1"/></r>`,
		want: []string{"a/xml:id=1"},
	}, {
		name: "the document element itself is scanned",
		src:  `<r Id="top"><a/></r>`,
		want: []string{"r/Id=top"},
	}, {
		name: "nested elements are all scanned",
		src:  `<r Id="1"><a Id="2"><b Id="3"/></a></r>`,
		want: []string{"r/Id=1", "a/Id=2", "b/Id=3"},
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc := mustParse(t, tc.src)
			doc.RegisterIDs()
			if got := idsOf(doc); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("IDs = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestElementByID(t *testing.T) {
	doc := mustParse(t, `<r Id="root"><a ID="one"/><b xml:id="two"/><c/></r>`)
	doc.RegisterIDs()

	for id, wantLocal := range map[string]string{"root": "r", "one": "a", "two": "b"} {
		el := doc.ElementByID(id)
		if el == nil {
			t.Errorf("ElementByID(%q) = nil", id)
			continue
		}
		if el.Name.Local != wantLocal {
			t.Errorf("ElementByID(%q) = <%s>, want <%s>", id, el.Name.Local, wantLocal)
		}
	}
	if doc.ElementByID("missing") != nil {
		t.Error("ElementByID for an unknown value should be nil")
	}
	// A leading '#' is not accepted; the caller strips the fragment delimiter.
	if doc.ElementByID("#root") != nil {
		t.Error(`ElementByID("#root") should be nil`)
	}
	// Without RegisterIDs nothing is indexed.
	fresh := mustParse(t, `<r Id="root"/>`)
	if fresh.ElementByID("root") != nil {
		t.Error("IDs should not be indexed until RegisterIDs is called")
	}
	// A non-document node yields nothing rather than panicking.
	if doc.DocumentElement().ElementByID("root") != nil {
		t.Error("ElementByID on a non-document node should be nil")
	}
}

// TestDuplicateIDs mirrors DSSXMLUtils.isDuplicateIdsDetected: registration never fails,
// the collision is reported, and the LAST occurrence in document order is the one
// getElementById answers with.
//
// Last, not first: CoreDocumentImpl.putIdentifier is identifiers.put(id, element) into a
// HashMap and the recursive Id browse walks the document in order, so each duplicate
// replaces the one before it. Probed against OpenJDK 21 with DSS's own registration loop -
// <a Id="x">first</a><b Id="x">second</b><c Id="x">third</c> answers <c>.
//
// It is worth stating why this is not cosmetic. Two elements sharing an Id is the shape of
// an XML signature wrapping attack: the attacker leaves the signed element in place and
// appends a second one with the same Id carrying their content. Resolving "#x" to the first
// element makes the digest match and the forgery verify;
// internal/xmldsig/testdata/corpus/validation/dss2329/xades-with-manifest-with-duplicated-reference.xml
// is exactly that document, and Santuario answers false on it.
func TestDuplicateIDs(t *testing.T) {
	doc := mustParse(t, `<r><a Id="dup"/><b Id="dup"/><c Id="other"/><d Id="dup"/><e Id="z"/><f Id="z"/></r>`)
	doc.RegisterIDs()

	if got, want := doc.DuplicateIDs(), []string{"dup", "z"}; !reflect.DeepEqual(got, want) {
		t.Errorf("DuplicateIDs = %v, want %v (sorted)", got, want)
	}
	el := doc.ElementByID("dup")
	if el == nil || el.Name.Local != "d" {
		t.Errorf("ElementByID(dup) = %v, want the last in document order <d>", el)
	}
	if el := doc.ElementByID("z"); el == nil || el.Name.Local != "f" {
		t.Errorf("ElementByID(z) = %v, want the last in document order <f>", el)
	}
	if got := mustParse(t, `<r Id="x"/>`); len(func() []string { got.RegisterIDs(); return got.DuplicateIDs() }()) != 0 {
		t.Error("a document without collisions reported duplicates")
	}
}

// TestIDIndexIsInvalidatedByMutation covers the lazy-rebuild contract: every mutating
// entry point bumps the document generation and the next query rebuilds.
func TestIDIndexIsInvalidatedByMutation(t *testing.T) {
	doc := mustParse(t, `<r Id="root"><a Id="one"/></r>`)
	doc.RegisterIDs()
	root := doc.DocumentElement()
	a := root.FirstElementChild()

	if doc.ElementByID("one") != a {
		t.Fatal("initial index is wrong")
	}

	t.Run("SetAttr", func(t *testing.T) {
		a.SetAttr(Name{Local: "Id"}, "changed")
		if doc.ElementByID("one") != nil {
			t.Error("the stale ID is still resolvable")
		}
		if doc.ElementByID("changed") != a {
			t.Error("the new ID was not indexed")
		}
	})

	t.Run("RemoveAttr", func(t *testing.T) {
		a.RemoveAttr("", "Id")
		if doc.ElementByID("changed") != nil {
			t.Error("the ID survived RemoveAttr")
		}
	})

	t.Run("AppendChild", func(t *testing.T) {
		n := NewElement(Name{Local: "n"})
		n.SetAttr(Name{Local: "Id"}, "fresh")
		root.AppendChild(n)
		if doc.ElementByID("fresh") != n {
			t.Error("an appended element's ID was not indexed")
		}
	})

	t.Run("RemoveChild", func(t *testing.T) {
		n := doc.ElementByID("fresh")
		root.RemoveChild(n)
		if doc.ElementByID("fresh") != nil {
			t.Error("a removed element's ID is still resolvable")
		}
	})

	t.Run("InsertBefore", func(t *testing.T) {
		n := NewElement(Name{Local: "n2"})
		n.SetAttr(Name{Local: "Id"}, "ins")
		root.InsertBefore(n, root.FirstChild)
		if doc.ElementByID("ins") != n {
			t.Error("an inserted element's ID was not indexed")
		}
	})

	t.Run("ReplaceChild", func(t *testing.T) {
		n := NewElement(Name{Local: "n3"})
		n.SetAttr(Name{Local: "Id"}, "rep")
		root.ReplaceChild(n, doc.ElementByID("ins"))
		if doc.ElementByID("ins") != nil {
			t.Error("the replaced element's ID is still resolvable")
		}
		if doc.ElementByID("rep") != n {
			t.Error("the replacement element's ID was not indexed")
		}
	})

	t.Run("SetTextContent", func(t *testing.T) {
		before := doc.ElementByID("root")
		root.SetTextContent("wiped")
		if doc.ElementByID("rep") != nil {
			t.Error("an ID under a wiped subtree is still resolvable")
		}
		if doc.ElementByID("root") != before {
			t.Error("the root's own ID should have survived")
		}
	})
}

// TestRegisterIDAttr covers the manual, single-attribute path and its survival across
// a lazy rebuild.
func TestRegisterIDAttr(t *testing.T) {
	doc := mustParse(t, `<r><a Ref="manual"/></r>`)
	a := doc.DocumentElement().FirstElementChild()
	doc.RegisterIDAttr(a, a.Attr("", "Ref"))

	if doc.ElementByID("manual") != a {
		t.Fatal("a manually registered attribute is not resolvable")
	}
	// A structural change forces a rebuild; the manual registration must survive it.
	doc.DocumentElement().AppendChild(NewElement(Name{Local: "z"}))
	if doc.ElementByID("manual") != a {
		t.Error("the manual registration was lost on rebuild")
	}
	// Registering the same attribute twice is a no-op, not a duplicate.
	doc.RegisterIDAttr(a, a.Attr("", "Ref"))
	if got := doc.DuplicateIDs(); len(got) != 0 {
		t.Errorf("re-registering the same attribute reported duplicates %v", got)
	}
	if got := doc.IDAttrs(a); len(got) != 1 {
		t.Errorf("IDAttrs = %v, want one entry", got)
	}
	// Manual and scanned registrations coexist.
	doc.RegisterIDs()
	if doc.ElementByID("manual") != a {
		t.Error("RegisterIDs discarded the manual registration")
	}
	doc.RegisterIDAttr(nil, nil) // must not panic
}

func TestIDAttrsOrderAndCopy(t *testing.T) {
	doc := mustParse(t, `<r><a Id="1" xml:id="2"/></r>`)
	doc.RegisterIDs()
	a := doc.DocumentElement().FirstElementChild()

	got := doc.IDAttrs(a)
	if len(got) != 2 {
		t.Fatalf("IDAttrs = %v, want two entries", got)
	}
	if got[0].Name.QName() != "Id" || got[1].Name.QName() != "xml:id" {
		t.Errorf("IDAttrs order = %q, %q, want attribute order Id then xml:id",
			got[0].Name.QName(), got[1].Name.QName())
	}
	// The returned slice is a snapshot: mutating it must not corrupt the index.
	got[0] = nil
	if again := doc.IDAttrs(a); again[0] == nil {
		t.Error("IDAttrs handed out its internal slice")
	}
	if doc.IDAttrs(NewElement(Name{Local: "unknown"})) != nil {
		t.Error("IDAttrs for an unregistered element should be nil")
	}
}

// TestRegisterIDsOnRealisticXAdES exercises the rule on the attribute shapes that
// actually appear in a signature.
func TestRegisterIDsOnRealisticXAdES(t *testing.T) {
	const src = `<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="sig-1">
  <ds:SignedInfo Id="si-1"><ds:Reference URI="#sp-1"/></ds:SignedInfo>
  <ds:Object><xades:QualifyingProperties xmlns:xades="http://uri.etsi.org/01903/v1.3.2#" Target="#sig-1">
    <xades:SignedProperties Id="sp-1"/>
  </xades:QualifyingProperties></ds:Object>
</ds:Signature>`
	doc := mustParse(t, src)
	doc.RegisterIDs()

	for _, id := range []string{"sig-1", "si-1", "sp-1"} {
		if doc.ElementByID(id) == nil {
			t.Errorf("ElementByID(%q) = nil", id)
		}
	}
	// URI and Target are not ID attributes and must not be indexed.
	if doc.ElementByID("#sp-1") != nil || doc.ElementByID("#sig-1") != nil {
		t.Error("a URI or Target value was indexed as an ID")
	}
	sp := doc.ElementByID("sp-1")
	if sp == nil || sp.Name.Local != "SignedProperties" {
		t.Fatalf("sp-1 resolved to %v", sp)
	}
	if got := strings.TrimPrefix(sp.Name.Space, ""); got != "http://uri.etsi.org/01903/v1.3.2#" {
		t.Errorf("SignedProperties namespace = %q", got)
	}
}
