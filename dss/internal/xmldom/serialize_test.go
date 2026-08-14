package xmldom

import (
	"errors"
	"strings"
	"testing"
)

func TestSerialize(t *testing.T) {
	tests := []struct {
		name, src, want string
	}{
		{"empty element stays short", `<r/>`, `<r/>`},
		{"element pair becomes short", `<r></r>`, `<r/>`},
		{"attributes are sorted, declarations first", `<r z="1" a="2" xmlns:p="urn:1" p:b="3"/>`, `<r xmlns:p="urn:1" a="2" p:b="3" z="1"/>`},
		{"prefixes are verbatim", `<a:r xmlns:a="urn:a"><a:c/></a:r>`, `<a:r xmlns:a="urn:a"><a:c/></a:r>`},
		{"text escapes", `<r>a&amp;b&lt;c&gt;d"e'f</r>`, `<r>a&amp;b&lt;c&gt;d"e'f</r>`},
		{"attribute escapes", `<r a="&amp;&lt;&quot;&apos;>"/>`, `<r a="&amp;&lt;&quot;'&gt;"/>`},
		{"cdata is preserved as cdata", `<r><![CDATA[a < b]]></r>`, `<r><![CDATA[a < b]]></r>`},
		{"comment", `<r><!-- x --></r>`, `<r><!-- x --></r>`},
		{"pi with data", `<r><?t d?></r>`, `<r><?t d?></r>`},
		{"pi without data", `<r><?t?></r>`, `<r><?t?></r>`},
		{"prolog and epilog", `<!--a--><?p x?><r/><?q y?><!--b-->`, `<!--a--><?p x?><r/><?q y?><!--b-->`},
		{"whitespace is preserved", "<r>  <c>\t</c>\n</r>", "<r>  <c>\t</c>\n</r>"},
		{"default namespace undeclaration", `<r xmlns="urn:a"><c xmlns=""/></r>`, `<r xmlns="urn:a"><c xmlns=""/></r>`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := mustSerialize(t, mustParse(t, tc.src)); got != tc.want {
				t.Errorf("Serialize(%q) = %q, want %q", tc.src, got, tc.want)
			}
		})
	}
}

// TestSerializeEscapesWhitespaceForRoundTripping covers the whitespace CharInfo marks
// special without giving it an entity name: a TAB/LF/CR in an attribute value and a CR in
// text all become DECIMAL character references. Unescaped they would not survive a reparse
// - clause 3.3.3 folds attribute whitespace to a space and clause 2.11 folds a literal CR
// to LF - which is why Java escapes them too.
func TestSerializeEscapesWhitespaceForRoundTripping(t *testing.T) {
	doc := NewDocument()
	root := NewElement(Name{Local: "r"})
	doc.AppendChild(root)
	root.SetAttr(Name{Local: "a"}, "x\ty\nz\rw")
	root.AppendChild(NewText("p\rq\nr"))

	got := mustSerialize(t, doc)
	const want = `<r a="x&#9;y&#10;z&#13;w">p&#13;q` + "\n" + `r</r>`
	if got != want {
		t.Fatalf("Serialize = %q, want %q", got, want)
	}
	back := mustParseRoot(t, got)
	if v := back.AttrValue("", "a"); v != "x\ty\nz\rw" {
		t.Errorf("attribute did not survive the round trip: %q", v)
	}
	if v := back.TextContent(); v != "p\rq\nr" {
		t.Errorf("text did not survive the round trip: %q", v)
	}
}

func TestSerializeCDATASplitsOnEndMarker(t *testing.T) {
	doc := NewDocument()
	root := NewElement(Name{Local: "r"})
	doc.AppendChild(root)
	root.AppendChild(NewCDATA("a]]>b"))

	got := mustSerialize(t, doc)
	if want := `<r><![CDATA[a]]]]><![CDATA[>b]]></r>`; got != want {
		t.Errorf("Serialize = %q, want %q", got, want)
	}
	if v := mustParseRoot(t, got).TextContent(); v != "a]]>b" {
		t.Errorf("round trip = %q, want %q", v, "a]]>b")
	}
}

func TestSerializeOptions(t *testing.T) {
	doc := mustParse(t, `<r/>`)

	if got, want := string(mustBytes(t, doc, nil)),
		`<?xml version="1.0" encoding="UTF-8" standalone="no"?><r/>`; got != want {
		t.Errorf("nil options = %q, want %q", got, want)
	}
	if got, want := string(mustBytes(t, doc, &SerializeOptions{XMLDeclaration: true})),
		`<?xml version="1.0" encoding="UTF-8" standalone="no"?><r/>`; got != want {
		t.Errorf("empty Encoding should default to UTF-8, got %q", got)
	}
	if got, want := string(mustBytes(t, doc, &SerializeOptions{XMLDeclaration: true, Encoding: "ISO-8859-1"})),
		`<?xml version="1.0" encoding="ISO-8859-1" standalone="no"?><r/>`; got != want {
		t.Errorf("Encoding = %q, want %q", got, want)
	}
	if got, want := string(mustBytes(t, doc, &SerializeOptions{})), `<r/>`; got != want {
		t.Errorf("XMLDeclaration false = %q, want %q", got, want)
	}
	// A node that is not the Document never carries standalone, but it does carry the
	// declaration: DomUtils sets no OMIT_XML_DECLARATION for any node kind.
	if got, want := string(mustBytes(t, doc.DocumentElement(), nil)),
		`<?xml version="1.0" encoding="UTF-8"?><r/>`; got != want {
		t.Errorf("element with default options = %q, want %q", got, want)
	}
	// An encoding Java has no charset for fails on both sides.
	if _, err := doc.Bytes(&SerializeOptions{Encoding: "latin-1"}); err == nil {
		t.Error("expected an error for an encoding the Transformer cannot write")
	}
}

func mustBytes(t *testing.T, n *Node, opts *SerializeOptions) []byte {
	t.Helper()
	b, err := n.Bytes(opts)
	if err != nil {
		t.Fatalf("Bytes: %v", err)
	}
	return b
}

func TestSerializeWriteError(t *testing.T) {
	doc := mustParse(t, `<r><a/><b/></r>`)
	want := errors.New("boom")
	err := doc.Serialize(failingWriter{want}, nil)
	if !errors.Is(err, want) {
		t.Errorf("Serialize error = %v, want %v", err, want)
	}
}

type failingWriter struct{ err error }

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }

// TestSerializeStandaloneNodes pins the node kinds DOM2TO handles outside an element.
// An attribute produces NOTHING - DOM2TO's switch ignores ATTRIBUTE_NODE - and a bare text
// node is written at element depth 0, where accumDefaultEscape's "depth > 0" guard means a
// carriage return goes out literally instead of as &#13;.
func TestSerializeStandaloneNodes(t *testing.T) {
	if got := mustSerialize(t, NewAttr(Name{Prefix: "p", Local: "a"}, `v"w`)); got != "" {
		t.Errorf("attribute node = %q, want %q", got, "")
	}
	if got, want := mustSerialize(t, NewText("a<b")), `a&lt;b`; got != want {
		t.Errorf("text node = %q, want %q", got, want)
	}
	if got, want := mustSerialize(t, NewText("a\rb")), "a\rb"; got != want {
		t.Errorf("text node at depth 0 = %q, want %q", got, want)
	}
	if got := mustSerialize(t, nil); got != "" {
		t.Errorf("nil node = %q, want %q", got, "")
	}
}

// TestDOMRoundTrip is the gate from the design note, minus the canonicalization leg that
// xmlc14n adds: Parse -> Serialize -> Parse must yield an identical tree, so the
// serializer cannot lose or invent detail. The corpus is restricted to documents already
// in serializer normal form for attribute order and namespace declarations, since
// TestSerializeRewritesTheTree owns the cases where the first pass legitimately rewrites
// the tree.
func TestDOMRoundTrip(t *testing.T) {
	corpus := []string{
		`<r/>`,
		`<!--a--><?p1 x?><r><x/></r><?p2 y?><!--b-->`,
		`<r xmlns="urn:d" xmlns:a="urn:a" xmlns:b="urn:b" a="4" a:z="2" b:z="1" z="3"/>`,
		"<r a=\"x\ty\nz\r\nw\rv\" b=\"&#9;&#10;&#13;\"/>",
		`<r a="it's" b='say "hi"'/>`,
		`<r>a&amp;b&lt;c&gt;d</r>`,
		`<r>a&#13;b</r>`,
		`<r><![CDATA[a < b & c]]>tail<![CDATA[x]]]]><![CDATA[>y]]></r>`,
		`<r><!--c1--><a><!--c2--><?pi d?></a></r>`,
		`<r xmlns:p="urn:1"><p:c><p:d xmlns:p="urn:2"/></p:c></r>`,
		`<r xmlns:unused="urn:u"><c/></r>`,
		`<r xmlns="urn:a"><c xmlns=""><d/></c></r>`,
		`<r xml:base="b/" xml:id="i" xml:lang="en" xml:space="preserve"><t/></r>`,
		"<r>aé€\U0001F600́ </r>",
		`<r a="éюあ𐀀"/>`,
		"  <r>\n  <c>  </c>\n</r>\n",
		`<r Id="1"><a ID="2"><b xml:id="3"/></a></r>`,
		`<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="s"><ds:SignedInfo><ds:Reference URI="#x"/></ds:SignedInfo></ds:Signature>`,
	}
	for _, src := range corpus {
		t.Run(src, func(t *testing.T) {
			first := mustParse(t, src)
			out := mustBytes(t, first, nil)
			second, err := Parse(out, nil)
			if err != nil {
				t.Fatalf("reparse of %q failed: %v", out, err)
			}
			if got, want := dump(second), dump(first); got != want {
				t.Errorf("round trip changed the tree.\nserialized: %s\ngot:\n%s\nwant:\n%s", out, got, want)
			}
			// A second round must be byte-stable.
			out2 := mustBytes(t, second, nil)
			if string(out2) != string(out) {
				t.Errorf("serialization is not idempotent:\nfirst:  %s\nsecond: %s", out, out2)
			}
		})
	}
}

// TestDeepNestingRoundTrip keeps the recursive serializer honest at the depth limit.
func TestDeepNestingRoundTrip(t *testing.T) {
	const depth = 200
	src := strings.Repeat("<e>", depth) + "x" + strings.Repeat("</e>", depth)
	doc := mustParse(t, src)
	out := mustBytes(t, doc, &SerializeOptions{})
	if string(out) != src {
		t.Errorf("deep nesting round trip differs")
	}
	if got := doc.DocumentElement().TextContent(); got != "x" {
		t.Errorf("TextContent = %q, want %q", got, "x")
	}
}
