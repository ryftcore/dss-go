package xmldom

import (
	"errors"
	"strings"
	"testing"
)

// dump renders a tree as one indented line per node, attributes included, so that
// structural expectations can be written as literals.
func dump(n *Node) string {
	var b strings.Builder
	var walk func(*Node, int)
	walk = func(x *Node, depth int) {
		b.WriteString(strings.Repeat("  ", depth))
		b.WriteString(x.Kind.String())
		switch x.Kind {
		case Element, Attribute:
			b.WriteString(" " + x.Name.QName() + " {" + x.Name.Space + "}" + x.Name.Local)
		case ProcInst:
			b.WriteString(" " + x.Name.Local)
		}
		switch x.Kind {
		case Attribute, Text, CDATA, Comment, ProcInst:
			b.WriteString(" = " + quote(x.Value))
		}
		b.WriteString("\n")
		for _, a := range x.Attrs {
			walk(a, depth+1)
		}
		for c := x.FirstChild; c != nil; c = c.NextSibling {
			walk(c, depth+1)
		}
	}
	walk(n, 0)
	return b.String()
}

func quote(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '\t':
			b.WriteString(`\t`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '"':
			b.WriteString(`\"`)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}

func TestParseTreeShape(t *testing.T) {
	tests := []struct {
		name, src, want string
	}{{
		name: "minimal",
		src:  `<r/>`,
		want: `Document
  Element r {}r
`,
	}, {
		name: "prolog and epilog",
		src:  `<!--a--><?p1 x?><r><x/></r><?p2 y?><!--b-->`,
		want: `Document
  Comment = "a"
  ProcInst p1 = "x"
  Element r {}r
    Element x {}x
  ProcInst p2 = "y"
  Comment = "b"
`,
	}, {
		name: "declaration is not a node",
		src:  `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><r/>`,
		want: `Document
  Element r {}r
`,
	}, {
		name: "namespace declarations are attributes",
		src:  `<a:r xmlns:a="urn:a" xmlns="urn:d" a:q="1" p="2"><c/></a:r>`,
		want: `Document
  Element a:r {urn:a}r
    Attribute xmlns:a {http://www.w3.org/2000/xmlns/}a = "urn:a"
    Attribute xmlns {http://www.w3.org/2000/xmlns/}xmlns = "urn:d"
    Attribute a:q {urn:a}q = "1"
    Attribute p {}p = "2"
    Element c {urn:d}c
`,
	}, {
		name: "text cdata comment pi interleaved",
		src:  `<r>a<![CDATA[b<c]]>d<!--e--><?f g?></r>`,
		want: `Document
  Element r {}r
    Text = "a"
    CDATA = "b<c"
    Text = "d"
    Comment = "e"
    ProcInst f = "g"
`,
	}, {
		name: "adjacent cdata sections are separate nodes",
		src:  `<r><![CDATA[p]]><![CDATA[q]]></r>`,
		want: `Document
  Element r {}r
    CDATA = "p"
    CDATA = "q"
`,
	}, {
		name: "empty cdata section",
		src:  `<r><![CDATA[]]></r>`,
		want: `Document
  Element r {}r
    CDATA = ""
`,
	}, {
		name: "whitespace is never stripped",
		src:  "<r>  <c>\t</c>\n</r>",
		want: `Document
  Element r {}r
    Text = "  "
    Element c {}c
      Text = "\t"
    Text = "\n"
`,
	}, {
		name: "whitespace outside the root produces no node",
		src:  "  <r/>\n\n",
		want: `Document
  Element r {}r
`,
	}, {
		name: "xml prefix is always in scope",
		src:  `<r xml:lang="en" xml:space="preserve"/>`,
		want: `Document
  Element r {}r
    Attribute xml:lang {http://www.w3.org/XML/1998/namespace}lang = "en"
    Attribute xml:space {http://www.w3.org/XML/1998/namespace}space = "preserve"
`,
	}, {
		name: "unprefixed attribute is in no namespace even under a default ns",
		src:  `<r xmlns="urn:d" a="1"/>`,
		want: `Document
  Element r {urn:d}r
    Attribute xmlns {http://www.w3.org/2000/xmlns/}xmlns = "urn:d"
    Attribute a {}a = "1"
`,
	}, {
		name: "default namespace undeclaration",
		src:  `<r xmlns="urn:a"><c xmlns=""><d/></c></r>`,
		want: `Document
  Element r {urn:a}r
    Attribute xmlns {http://www.w3.org/2000/xmlns/}xmlns = "urn:a"
    Element c {}c
      Attribute xmlns {http://www.w3.org/2000/xmlns/}xmlns = ""
      Element d {}d
`,
	}, {
		name: "prefix rebound then restored",
		src:  `<p:a xmlns:p="urn:1"><p:b xmlns:p="urn:2"><p:c xmlns:p="urn:1"/></p:b></p:a>`,
		want: `Document
  Element p:a {urn:1}a
    Attribute xmlns:p {http://www.w3.org/2000/xmlns/}p = "urn:1"
    Element p:b {urn:2}b
      Attribute xmlns:p {http://www.w3.org/2000/xmlns/}p = "urn:2"
      Element p:c {urn:1}c
        Attribute xmlns:p {http://www.w3.org/2000/xmlns/}p = "urn:1"
`,
	}, {
		name: "explicit xml prefix declaration is allowed",
		src:  `<r xmlns:xml="http://www.w3.org/XML/1998/namespace" xml:id="x"/>`,
		want: `Document
  Element r {}r
    Attribute xmlns:xml {http://www.w3.org/2000/xmlns/}xml = "http://www.w3.org/XML/1998/namespace"
    Attribute xml:id {http://www.w3.org/XML/1998/namespace}id = "x"
`,
	}, {
		name: "processing instruction without data",
		src:  `<?a?><r><?b?></r>`,
		want: `Document
  ProcInst a = ""
  Element r {}r
    ProcInst b = ""
`,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := dump(mustParse(t, tc.src))
			if got != tc.want {
				t.Errorf("Parse(%q) tree =\n%s\nwant\n%s", tc.src, got, tc.want)
			}
		})
	}
}

func TestParseAttributeParentIsItsElement(t *testing.T) {
	el := mustParseRoot(t, `<r a="1"/>`)
	a := el.Attr("", "a")
	if a == nil {
		t.Fatal("attribute a not found")
	}
	if a.Parent != el {
		t.Errorf("attribute parent = %v, want the element", a.Parent)
	}
	if a.Document() == nil {
		t.Error("attribute should reach its document by walking Parent")
	}
	if !el.Contains(a) {
		t.Error("element should contain its attribute")
	}
}

func TestParseDoctypeAllowedByOption(t *testing.T) {
	const src = `<!DOCTYPE r SYSTEM "x.dtd"><r/>`
	if _, err := Parse([]byte(src), nil); err == nil {
		t.Fatal("DOCTYPE accepted with default options")
	}
	doc, err := Parse([]byte(src), &ParseOptions{AllowDoctype: true})
	if err != nil {
		t.Fatalf("AllowDoctype: %v", err)
	}
	if got := dump(doc); got != "Document\n  Element r {}r\n" {
		t.Errorf("DOCTYPE should be tolerated but produce no node, got\n%s", got)
	}
	// Tolerating the declaration must not start processing the subset.
	if _, err := Parse([]byte(`<!DOCTYPE r [<!ENTITY e "v">]><r>&e;</r>`), &ParseOptions{AllowDoctype: true}); err == nil {
		t.Error("entity reference declared in an internal subset was expanded; it must stay a hard error")
	}
}

func TestParseMaxDepth(t *testing.T) {
	build := func(n int) string {
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString("<e>")
		}
		for i := 0; i < n; i++ {
			b.WriteString("</e>")
		}
		return b.String()
	}
	if _, err := Parse([]byte(build(500)), nil); err != nil {
		t.Errorf("500 levels should fit the default limit: %v", err)
	}
	if _, err := Parse([]byte(build(501)), nil); err == nil {
		t.Error("501 levels accepted at the default MaxDepth of 500")
	}
	if _, err := Parse([]byte(build(4)), &ParseOptions{MaxDepth: 3}); err == nil {
		t.Error("MaxDepth 3 accepted 4 levels")
	}
	if _, err := Parse([]byte(build(3)), &ParseOptions{MaxDepth: 3}); err != nil {
		t.Errorf("MaxDepth 3 rejected 3 levels: %v", err)
	}
}

func TestParseMaxBytes(t *testing.T) {
	src := []byte(`<r>0123456789</r>`)
	if _, err := Parse(src, &ParseOptions{MaxBytes: int64(len(src))}); err != nil {
		t.Errorf("exactly MaxBytes rejected: %v", err)
	}
	if _, err := Parse(src, &ParseOptions{MaxBytes: int64(len(src)) - 1}); err == nil {
		t.Error("over MaxBytes accepted")
	}
	if _, err := ParseReader(strings.NewReader(string(src)), &ParseOptions{MaxBytes: int64(len(src)) - 1}); err == nil {
		t.Error("ParseReader over MaxBytes accepted")
	}
	if _, err := ParseReader(strings.NewReader(string(src)), nil); err != nil {
		t.Errorf("ParseReader unlimited: %v", err)
	}
}

func TestSyntaxErrorCarriesPosition(t *testing.T) {
	_, err := Parse([]byte("<r>\n  <a>\n  </b>\n</r>"), nil)
	var se *SyntaxError
	if !errors.As(err, &se) {
		t.Fatalf("error is %T (%v), want *SyntaxError", err, err)
	}
	if se.Line != 3 {
		t.Errorf("Line = %d, want 3", se.Line)
	}
	if se.Column != 3 {
		t.Errorf("Column = %d, want 3", se.Column)
	}
	if !strings.Contains(se.Error(), "does not match start tag") {
		t.Errorf("Error() = %q", se.Error())
	}
}
