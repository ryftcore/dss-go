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
		{"empty element becomes a pair", `<r/>`, `<r></r>`},
		{"element pair stays a pair", `<r></r>`, `<r></r>`},
		{"attributes keep document order", `<r z="1" a="2" xmlns:p="urn:1" p:b="3"/>`, `<r z="1" a="2" xmlns:p="urn:1" p:b="3"></r>`},
		{"prefixes are verbatim", `<a:r xmlns:a="urn:a"><a:c/></a:r>`, `<a:r xmlns:a="urn:a"><a:c></a:c></a:r>`},
		{"text escapes", `<r>a&amp;b&lt;c&gt;d"e'f</r>`, `<r>a&amp;b&lt;c&gt;d"e'f</r>`},
		{"attribute escapes", `<r a="&amp;&lt;&quot;&apos;>"/>`, `<r a="&amp;&lt;&quot;'>"></r>`},
		{"cdata is preserved as cdata", `<r><![CDATA[a < b]]></r>`, `<r><![CDATA[a < b]]></r>`},
		{"comment", `<r><!-- x --></r>`, `<r><!-- x --></r>`},
		{"pi with data", `<r><?t d?></r>`, `<r><?t d?></r>`},
		{"pi without data", `<r><?t?></r>`, `<r><?t?></r>`},
		{"prolog and epilog", `<!--a--><?p x?><r/><?q y?><!--b-->`, `<!--a--><?p x?><r></r><?q y?><!--b-->`},
		{"whitespace is preserved", "<r>  <c>\t</c>\n</r>", "<r>  <c>\t</c>\n</r>"},
		{"default namespace undeclaration", `<r xmlns="urn:a"><c xmlns=""/></r>`, `<r xmlns="urn:a"><c xmlns=""></c></r>`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := mustSerialize(t, mustParse(t, tc.src)); got != tc.want {
				t.Errorf("Serialize(%q) = %q, want %q", tc.src, got, tc.want)
			}
		})
	}
}

// TestSerializeEscapesWhitespaceForRoundTripping documents the deliberate departure
// from the design note's escaping table. Unescaped, a TAB/LF/CR in an attribute value
// would be normalized to a space by clause 3.3.3 on the next parse, and a CR in text
// would become LF under clause 2.11, so Serialize would be lossy. Java's Transformer
// escapes these for the same reason.
func TestSerializeEscapesWhitespaceForRoundTripping(t *testing.T) {
	doc := NewDocument()
	root := NewElement(Name{Local: "r"})
	doc.AppendChild(root)
	root.SetAttr(Name{Local: "a"}, "x\ty\nz\rw")
	root.AppendChild(NewText("p\rq\nr"))

	got := mustSerialize(t, doc)
	const want = `<r a="x&#x9;y&#xA;z&#xD;w">p&#xD;q` + "\n" + `r</r>`
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

	if got, want := string(mustBytes(t, doc, nil)), `<?xml version="1.0" encoding="UTF-8"?><r></r>`; got != want {
		t.Errorf("nil options = %q, want %q", got, want)
	}
	if got, want := string(mustBytes(t, doc, &SerializeOptions{XMLDeclaration: true})), `<?xml version="1.0" encoding="UTF-8"?><r></r>`; got != want {
		t.Errorf("empty Encoding should default to UTF-8, got %q", got)
	}
	if got, want := string(mustBytes(t, doc, &SerializeOptions{XMLDeclaration: true, Encoding: "ISO-8859-1"})),
		`<?xml version="1.0" encoding="ISO-8859-1"?><r></r>`; got != want {
		t.Errorf("Encoding = %q, want %q", got, want)
	}
	if got, want := string(mustBytes(t, doc, &SerializeOptions{})), `<r></r>`; got != want {
		t.Errorf("XMLDeclaration false = %q, want %q", got, want)
	}
	// An element is a fragment, so it never carries a declaration.
	if got, want := string(mustBytes(t, doc.DocumentElement(), nil)), `<r></r>`; got != want {
		t.Errorf("element with default options = %q, want %q", got, want)
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

func TestSerializeStandaloneNodes(t *testing.T) {
	if got, want := mustSerialize(t, NewAttr(Name{Prefix: "p", Local: "a"}, `v"w`)), `p:a="v&quot;w"`; got != want {
		t.Errorf("attribute node = %q, want %q", got, want)
	}
	if got, want := mustSerialize(t, NewText("a<b")), `a&lt;b`; got != want {
		t.Errorf("text node = %q, want %q", got, want)
	}
	if got := mustSerialize(t, nil); got != "" {
		t.Errorf("nil node = %q, want %q", got, "")
	}
}

// TestDOMRoundTrip is the gate from the design note, minus the canonicalization leg
// that phase 4a's xmlc14n will add: Parse -> Serialize -> Parse must yield an
// identical tree, so the serializer cannot lose or invent detail.
func TestDOMRoundTrip(t *testing.T) {
	corpus := []string{
		`<r/>`,
		`<!--a--><?p1 x?><r><x/></r><?p2 y?><!--b-->`,
		`<r xmlns:b="urn:b" xmlns:a="urn:a" b:z="1" a:z="2" z="3" a="4" xmlns="urn:d"/>`,
		"<r a=\"x\ty\nz\r\nw\rv\" b=\"&#9;&#10;&#13;\"/>",
		`<r a="it's" b='say "hi"'/>`,
		`<r>a&amp;b&lt;c&gt;d</r>`,
		`<r>a&#13;b</r>`,
		`<r><![CDATA[a < b & c]]>tail<![CDATA[]]><![CDATA[x]]]]><![CDATA[>y]]></r>`,
		`<r><!--c1--><a><!--c2--><?pi d?></a></r>`,
		`<r xmlns:p="urn:1"><p:c xmlns:p="urn:1"><p:d xmlns:p="urn:2"/></p:c></r>`,
		`<r xmlns:unused="urn:u"><c/></r>`,
		`<r xmlns="urn:a"><c xmlns=""><d/></c></r>`,
		`<r xmlns:xml="http://www.w3.org/XML/1998/namespace" xml:lang="en" xml:space="preserve" xml:id="i" xml:base="b/"><t/></r>`,
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
