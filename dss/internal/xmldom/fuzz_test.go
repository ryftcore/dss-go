package xmldom

import (
	"testing"
)

// FuzzParse is the invariant gate from the design note's test strategy, restricted to
// what this package owns: Parse must never panic on arbitrary input, and whenever it
// succeeds the tree must be internally consistent and must survive Serialize followed
// by Parse unchanged. The canonicalization leg is added by internal/xmlc14n.
func FuzzParse(f *testing.F) {
	seeds := []string{
		"", "<", "<r", "<r>", "<r/>", "<r></r>", "<r></s>", "<a/><b/>", "<r/>tail",
		`<?xml version="1.0" encoding="UTF-8"?><r/>`,
		`<?xml version="1.1"?><r/>`,
		"\xef\xbb\xbf<r/>",
		"\xfe\xff\x00<\x00r\x00/\x00>",
		`<!DOCTYPE r><r/>`,
		`<!DOCTYPE r [<!ENTITY e "v">]><r>&e;</r>`,
		`<r a="1" a="2"/>`,
		`<r xmlns:p=""/>`,
		`<r xmlns:xml="urn:bogus"/>`,
		`<p:r/>`,
		`<r>&foo;</r>`,
		`<r>&#0;</r>`,
		`<r>&#xD800;</r>`,
		`<r>]]></r>`,
		`<r a="<"/>`,
		`<r a=">" b="]]>"/>`,
		"<r a=\"x\ty\nz\r\nw\rv\" b=\"&#9;&#10;&#13;\"/>",
		`<r><![CDATA[a < b]]]]><![CDATA[>c]]></r>`,
		`<!--a--><?p x?><r><x/></r><?q y?><!--b-->`,
		`<r xmlns:b="urn:b" xmlns:a="urn:a" b:z="1" a:z="2" z="3" a="4" xmlns="urn:d"/>`,
		`<r xmlns="urn:a"><c xmlns=""><d/></c></r>`,
		`<r xml:lang="en" xml:space="preserve" xml:id="i" xml:base="b/"><t/></r>`,
		"<r>aé€\U0001F600</r>",
		`<?xml version="1.0" encoding="ISO-8859-1"?><r>` + "\xe9" + `</r>`,
		`<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="s"><ds:SignedInfo/></ds:Signature>`,
	}
	for _, s := range seeds {
		f.Add([]byte(s))
	}

	f.Fuzz(func(t *testing.T, src []byte) {
		doc, err := Parse(src, &ParseOptions{MaxBytes: 1 << 20})
		if err != nil {
			return
		}
		checkInvariants(t, doc)

		// Every accessor must be safe on whatever came out.
		doc.RegisterIDs()
		doc.DuplicateIDs()
		doc.Walk(func(n *Node) bool {
			n.TextContent()
			n.Ancestors()
			n.Depth()
			n.Children()
			n.Elements()
			n.LookupNamespaceURI("")
			n.LookupPrefix("urn:x")
			for _, a := range n.Attrs {
				a.LookupNamespaceURI(a.Name.Prefix)
				doc.IDAttrs(n)
			}
			return true
		})

		out, err := doc.Bytes(nil)
		if err != nil {
			t.Fatalf("Serialize failed for a parsed document: %v", err)
		}
		again, err := Parse(out, nil)
		if err != nil {
			t.Fatalf("reparse of %q failed: %v", out, err)
		}
		if got, want := dump(again), dump(doc); got != want {
			t.Fatalf("round trip changed the tree\ninput:      %q\nserialized: %s\ngot:\n%s\nwant:\n%s", src, out, got, want)
		}
		out2, err := again.Bytes(nil)
		if err != nil {
			t.Fatalf("second Serialize failed: %v", err)
		}
		if string(out2) != string(out) {
			t.Fatalf("serialization is not idempotent\nfirst:  %s\nsecond: %s", out, out2)
		}
	})
}

// checkInvariants asserts the structural guarantees every parsed tree must hold.
func checkInvariants(t *testing.T, doc *Node) {
	t.Helper()
	if doc.Kind != Document {
		t.Fatalf("root node kind = %s, want Document", doc.Kind)
	}
	if doc.DocumentElement() == nil {
		t.Fatal("a parsed document must have a document element")
	}
	roots := 0
	doc.Walk(func(n *Node) bool {
		if n.Parent == nil && n != doc {
			t.Fatalf("%s node has no parent", n.Kind)
		}
		if n.Kind == Element && n.Parent == doc {
			roots++
		}
		if n.FirstChild != nil && n.FirstChild.PrevSibling != nil {
			t.Fatal("FirstChild has a previous sibling")
		}
		if n.LastChild != nil && n.LastChild.NextSibling != nil {
			t.Fatal("LastChild has a next sibling")
		}
		var last *Node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Parent != n {
				t.Fatal("child does not point back at its parent")
			}
			if c.PrevSibling != last {
				t.Fatal("PrevSibling link is broken")
			}
			if !canContain(n.Kind, c.Kind) {
				t.Fatalf("%s node contains a %s child", n.Kind, c.Kind)
			}
			last = c
		}
		if last != n.LastChild {
			t.Fatal("LastChild does not match the end of the child list")
		}
		for _, a := range n.Attrs {
			if a.Kind != Attribute {
				t.Fatalf("Attrs contains a %s node", a.Kind)
			}
			if a.Parent != n {
				t.Fatal("attribute does not point back at its element")
			}
			if n.Kind != Element {
				t.Fatalf("%s node carries attributes", n.Kind)
			}
			if a.FirstChild != nil {
				t.Fatal("attribute has children")
			}
		}
		return true
	})
	if roots != 1 {
		t.Fatalf("document has %d element children, want 1", roots)
	}
}
