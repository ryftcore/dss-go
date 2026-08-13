package xmldom

import (
	"bytes"
	"testing"
)

// serializerNormalForm documents exactly which inputs Parse+Serialize reproduces byte for
// byte. Serialize is not canonicalization and makes no byte-parity promise in general (see
// the package comment), but a document already written the way Serialize writes one MUST come
// back unchanged - otherwise a XAdES document that needed no DOM surgery would still be
// rewritten, and every reference digest over an unmodified part of it would move.
//
// A document is in serializer normal form when all of the following hold:
//
//   - it is UTF-8 with no byte-order mark;
//   - it opens with exactly <?xml version="1.0" encoding="UTF-8"?> and nothing before it;
//   - every element is written <e></e>, never <e/>;
//   - attribute values are delimited by '"';
//   - the only escapes used are &amp; &lt; &gt; in text, &amp; &lt; &quot; in attribute
//     values, and &#x9; &#xA; &#xD; for whitespace that attribute-value normalization would
//     otherwise destroy;
//   - no literal CR or CRLF appears anywhere (clause 2.11 would fold it to LF);
//   - no literal TAB or LF appears inside an attribute value (clause 3.3.3 would fold it to
//     a space).
const normalFormDecl = `<?xml version="1.0" encoding="UTF-8"?>`

func TestParseSerializeIsByteExactInNormalForm(t *testing.T) {
	bodies := []string{
		`<r></r>`,
		`<r a="1" b="2"></r>`,
		`<r xmlns="urn:d" xmlns:p="urn:1"><p:c q="v"></p:c></r>`,
		`<r>text</r>`,
		`<r>a&amp;b&lt;c&gt;d</r>`,
		`<r a="&amp;&lt;&quot;"></r>`,
		`<r a="&#x9;&#xA;&#xD;"></r>`,
		`<r>tab	and newline
here</r>`,
		`<r><![CDATA[a < b & c]]></r>`,
		`<r><!--comment--><?pi data?><?bare?></r>`,
		`<r><a></a><b>t</b><c><d></d></c></r>`,
		`<r xml:lang="en" xml:space="preserve" xml:base="http://x/a/"></r>`,
		`<r a="caf&#xD;é" b="日本語">中文 𝄞</r>`,
		`<r>&#xD;</r>`,
		`<r>  <a></a>
	<b> </b>
</r>`,
		`<café x="ü"></café>`,
		`<r xmlns:unused="urn:u" z="26" a="1"></r>`,
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

// TestSerializeIsIdempotent is the property that holds for EVERY parseable document, normal
// form or not: the first Serialize maps the tree into normal form, and every later one is the
// identity. A serializer that emits something it cannot itself re-parse to the same tree
// fails here even when no hand-written expectation covers the case.
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
		"<r> </r>",
		`<r xml:id="x" Id="y" ID="z"/>`,
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
			// And the tree itself must be stable, not merely its bytes: dump covers
			// kind, expanded name, prefix, value and attribute order.
			if a, b := dump(doc), dump(reparsed); a != b {
				t.Errorf("round trip changed the tree\nbefore: %s\nafter:  %s", a, b)
			}
		})
	}
}
