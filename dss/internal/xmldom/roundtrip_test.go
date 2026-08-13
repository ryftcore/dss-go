package xmldom

import (
	"bytes"
	"testing"
)

// serializer normal form documents exactly which inputs Parse+Serialize reproduces byte
// for byte. Serialize reproduces OpenJDK's identity Transformer (see serialize.go), so
// the normal form is the form that Transformer writes, and a document already in it MUST
// come back unchanged - otherwise a XAdES document that needed no DOM surgery would still
// be rewritten and every reference digest over an unmodified part of it would move.
//
// A document is in serializer normal form when all of the following hold:
//
//   - it is UTF-8 with no byte-order mark;
//   - it opens with exactly <?xml version="1.0" encoding="UTF-8" standalone="no"?> and
//     nothing before it (a document that declared standalone="yes" is NOT in normal form:
//     the Transformer drops the pseudo-attribute);
//   - every empty element is written <e/>, never <e></e>;
//   - each element's attributes are in ascending node-name order, with the namespace
//     declarations among them moved in front;
//   - no namespace declaration is redundant (a prefix rebound to the URI it already has)
//     and none declares a prefix beginning with "xml";
//   - attribute values are delimited by '"';
//   - the only escapes used are &amp; &lt; &gt; and &#13; in text, and &amp; &lt; &gt;
//     &quot; &#9; &#10; &#13; in attribute values;
//   - no literal CR or CRLF appears anywhere (clause 2.11 would fold it to LF);
//   - no literal TAB or LF appears inside an attribute value (clause 3.3.3 would fold it
//     to a space).
const normalFormDecl = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>`

func TestParseSerializeIsByteExactInNormalForm(t *testing.T) {
	bodies := []string{
		`<r/>`,
		`<r a="1" b="2"/>`,
		`<r xmlns="urn:d" xmlns:p="urn:1"><p:c q="v"/></r>`,
		`<r>text</r>`,
		`<r>a&amp;b&lt;c&gt;d</r>`,
		`<r a="&amp;&lt;&gt;&quot;"/>`,
		`<r a="&#9;&#10;&#13;"/>`,
		`<r>tab	and newline
here</r>`,
		`<r><![CDATA[a < b & c]]></r>`,
		`<r><!--comment--><?pi data?><?bare?></r>`,
		`<r><a/><b>t</b><c><d/></c></r>`,
		`<r xml:base="http://x/a/" xml:lang="en" xml:space="preserve"/>`,
		`<r a="caf&#13;é" b="日本語">中文 &#119070;</r>`,
		`<r>&#13;</r>`,
		`<r>  <a/>
	<b> </b>
</r>`,
		`<café x="ü"/>`,
		`<r xmlns:unused="urn:u" a="1" z="26"/>`,
	}
	for _, body := range bodies {
		src := normalFormDecl + body
		t.Run(body, func(t *testing.T) {
			doc, err := Parse([]byte(src), nil)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			got, err := doc.Bytes(nil)
			if err != nil {
				t.Fatalf("serialize: %v", err)
			}
			if string(got) != src {
				t.Errorf("round trip changed the bytes\n got: %q\nwant: %q", got, src)
			}
		})
	}
}

// TestSerializeIsIdempotent is the property that holds for EVERY parseable document,
// normal form or not: the first Serialize maps the tree into normal form and every later
// one is the identity. A serializer that emits something it cannot itself re-parse to the
// same tree fails here even when no hand-written expectation covers the case.
func TestSerializeIsIdempotent(t *testing.T) {
	inputs := []string{
		`<r/>`,
		"\xef\xbb\xbf<r>bom</r>",
		"<r a='single'/>",
		"<r a=\"x\ty\nz\r\nw\rv\"/>",
		`<?xml version="1.0" encoding="ISO-8859-1"?><r>x</r>`,
		"<r>a\r\nb\rc</r>",
		`<r><![CDATA[]]]]><![CDATA[>]]></r>`,
		`<r>&#13;&#10;&#9;</r>`,
		`<!--lead--><?pi x?><r/><?trail y?><!--tail-->`,
		`<r xmlns="urn:a"><c xmlns=""><d/></c></r>`,
		`<r a="&#38;&#60;&#62;&#34;&#39;"/>`,
		"<r> </r>",
		`<r xml:id="x" Id="y" ID="z"/>`,
		`<r xmlns:p="urn:1"><p:c xmlns:p="urn:1"/></r>`,
		`<r xmlns:xml="http://www.w3.org/XML/1998/namespace" xml:lang="en"/>`,
	}
	for _, src := range inputs {
		t.Run(src, func(t *testing.T) {
			doc, err := Parse([]byte(src), nil)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			first, err := doc.Bytes(nil)
			if err != nil {
				t.Fatalf("serialize: %v", err)
			}
			reparsed, err := Parse(first, nil)
			if err != nil {
				t.Fatalf("reparse %q: %v", first, err)
			}
			second, err := reparsed.Bytes(nil)
			if err != nil {
				t.Fatalf("re-serialize: %v", err)
			}
			if !bytes.Equal(first, second) {
				t.Errorf("Serialize is not idempotent\nfirst:  %q\nsecond: %q", first, second)
			}
			// And the tree, not merely its bytes, must be stable from the first
			// serialization onward: dump covers kind, expanded name, prefix, value and
			// attribute order. It is compared against the SECOND parse rather than the
			// original, because the first pass is where the Transformer's own rewriting
			// happens - see TestSerializeRewritesTheTree.
			again, err := Parse(second, nil)
			if err != nil {
				t.Fatalf("third parse: %v", err)
			}
			if a, b := dump(reparsed), dump(again); a != b {
				t.Errorf("round trip changed the tree\nbefore: %s\nafter:  %s", a, b)
			}
		})
	}
}

// TestSerializeRewritesTheTree pins the two ways a Parse/Serialize/Parse round trip does
// NOT give back the tree it started from. Both come straight from the Java serializer this
// one reproduces, both are invisible to canonicalization (which sorts attributes itself
// and ignores redundant declarations), and both are visible to the "physical" c14n
// algorithm, which is why internal/xmlc14n's round-trip property excludes it.
func TestSerializeRewritesTheTree(t *testing.T) {
	t.Run("attributes are reordered into node-name order", func(t *testing.T) {
		got := mustSerialize(t, mustParse(t, `<r z="1" a="2" xmlns:p="urn:1" p:b="3"/>`))
		if want := `<r xmlns:p="urn:1" a="2" p:b="3" z="1"/>`; got != want {
			t.Errorf("Serialize = %q, want %q", got, want)
		}
	})
	t.Run("a redundant redeclaration is dropped", func(t *testing.T) {
		got := mustSerialize(t, mustParse(t, `<r xmlns:p="urn:1"><p:c xmlns:p="urn:1"/></r>`))
		if want := `<r xmlns:p="urn:1"><p:c/></r>`; got != want {
			t.Errorf("Serialize = %q, want %q", got, want)
		}
	})
	t.Run("a declaration of the xml prefix is dropped", func(t *testing.T) {
		got := mustSerialize(t, mustParse(t,
			`<r xmlns:xml="http://www.w3.org/XML/1998/namespace" xml:lang="en"/>`))
		if want := `<r xml:lang="en"/>`; got != want {
			t.Errorf("Serialize = %q, want %q", got, want)
		}
	})
	// standalone="yes" is the one input for which Serialize is not idempotent, in Java
	// too: DOM2TO calls setStandalone only when getXmlStandalone() is false, so the
	// pseudo-attribute is dropped on the way out - and the document that comes back from
	// re-parsing that output is no longer standalone, so the next pass writes
	// standalone="no".
	t.Run("standalone yes is dropped and comes back as no", func(t *testing.T) {
		doc := mustParse(t, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><r/>`)
		first, err := doc.Bytes(nil)
		if err != nil {
			t.Fatalf("serialize: %v", err)
		}
		if want := `<?xml version="1.0" encoding="UTF-8"?><r/>`; string(first) != want {
			t.Fatalf("first pass = %q, want %q", first, want)
		}
		second, err := mustParse(t, string(first)).Bytes(nil)
		if err != nil {
			t.Fatalf("re-serialize: %v", err)
		}
		if want := `<?xml version="1.0" encoding="UTF-8" standalone="no"?><r/>`; string(second) != want {
			t.Fatalf("second pass = %q, want %q", second, want)
		}
	})
}
