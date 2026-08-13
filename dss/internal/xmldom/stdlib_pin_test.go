package xmldom

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
	"testing"
)

// This file pins the behaviour of encoding/xml that xmldom is built on. Every
// assertion here is a fact the parser relies on, so if a future Go release changes one
// of them this test fails loudly rather than the parser failing quietly. Each case
// notes whether the behaviour is something we exploit or something we compensate for.

// rawTokens drives RawToken and returns each token together with its exact source span.
func rawTokens(t *testing.T, src string) []struct {
	tok xml.Token
	raw string
} {
	t.Helper()
	d := xml.NewDecoder(bytes.NewReader([]byte(src)))
	d.CharsetReader = func(_ string, in io.Reader) (io.Reader, error) { return in, nil }
	var out []struct {
		tok xml.Token
		raw string
	}
	prev := int64(0)
	for {
		tok, err := d.RawToken()
		if err == io.EOF {
			return out
		}
		if err != nil {
			t.Fatalf("RawToken(%q): %v", src, err)
		}
		off := d.InputOffset()
		out = append(out, struct {
			tok xml.Token
			raw string
		}{xml.CopyToken(tok), src[prev:off]})
		prev = off
	}
}

// TestPinRawTokenKeepsPrefixes is why the parser uses RawToken and not Token: Token
// rewrites Name.Space to the resolved URI and destroys the literal prefix, which
// canonicalization must emit verbatim.
func TestPinRawTokenKeepsPrefixes(t *testing.T) {
	const src = `<a:root xmlns:a="urn:a" a:q="3"/>`

	dr := xml.NewDecoder(strings.NewReader(src))
	raw, err := dr.RawToken()
	if err != nil {
		t.Fatal(err)
	}
	se := raw.(xml.StartElement)
	if se.Name.Space != "a" || se.Name.Local != "root" {
		t.Errorf("RawToken name = {%q,%q}, want {%q,%q}", se.Name.Space, se.Name.Local, "a", "root")
	}

	dc := xml.NewDecoder(strings.NewReader(src))
	cooked, err := dc.Token()
	if err != nil {
		t.Fatal(err)
	}
	ce := cooked.(xml.StartElement)
	if ce.Name.Space != "urn:a" {
		t.Errorf("Token name.Space = %q, want the resolved URI %q", ce.Name.Space, "urn:a")
	}
}

// TestPinAttributeOrderAndNamespaceDeclShape pins that attributes arrive in document
// order and that namespace declarations arrive as ordinary attributes shaped exactly
// as org.w3c.dom shapes them.
func TestPinAttributeOrderAndNamespaceDeclShape(t *testing.T) {
	const src = `<a:root xmlns:z="urn:z" b="2" xmlns:a="urn:a" a="1" xmlns="urn:d" a:q="3" xml:lang="en"/>`
	toks := rawTokens(t, src)
	se := toks[0].tok.(xml.StartElement)

	want := []struct{ space, local string }{
		{"xmlns", "z"},
		{"", "b"},
		{"xmlns", "a"},
		{"", "a"},
		{"", "xmlns"}, // the default declaration, local name "xmlns", no prefix
		{"a", "q"},
		{"xml", "lang"},
	}
	if len(se.Attr) != len(want) {
		t.Fatalf("got %d attributes, want %d", len(se.Attr), len(want))
	}
	for i, w := range want {
		if se.Attr[i].Name.Space != w.space || se.Attr[i].Name.Local != w.local {
			t.Errorf("attr[%d] = {%q,%q}, want {%q,%q}",
				i, se.Attr[i].Name.Space, se.Attr[i].Name.Local, w.space, w.local)
		}
	}
}

// TestPinInputOffsetSpans pins the bracketing that recovers raw start-tag text: the
// StartElement span is the whole tag and a self-closing tag's synthetic EndElement
// span is empty. Both are load-bearing.
func TestPinInputOffsetSpans(t *testing.T) {
	toks := rawTokens(t, `<r a="1"/><s b="2"></s>`)
	wantRaw := []string{`<r a="1"/>`, ``, `<s b="2">`, `</s>`}
	if len(toks) != len(wantRaw) {
		t.Fatalf("got %d tokens, want %d", len(toks), len(wantRaw))
	}
	for i, w := range wantRaw {
		if toks[i].raw != w {
			t.Errorf("token %d raw = %q, want %q", i, toks[i].raw, w)
		}
	}
}

// TestPinCDATAIsOnlyVisibleInTheRawSpan pins the fact that a CDATA section and a text
// run are the same token type with the same value, so the raw span is the only thing
// that tells them apart. xmldom keeps the distinction for Serialize and for node-count
// parity with Xerces, which also emits one node per section.
func TestPinCDATAIsOnlyVisibleInTheRawSpan(t *testing.T) {
	toks := rawTokens(t, `<r>a<![CDATA[b]]>c</r>`)
	if len(toks) != 5 {
		t.Fatalf("got %d tokens, want 5", len(toks))
	}
	for i, want := range []struct{ val, raw string }{
		{"a", "a"},
		{"b", "<![CDATA[b]]>"},
		{"c", "c"},
	} {
		cd, ok := toks[i+1].tok.(xml.CharData)
		if !ok {
			t.Fatalf("token %d is %T, want xml.CharData", i+1, toks[i+1].tok)
		}
		if string(cd) != want.val || toks[i+1].raw != want.raw {
			t.Errorf("token %d = (%q, raw %q), want (%q, raw %q)", i+1, string(cd), toks[i+1].raw, want.val, want.raw)
		}
	}
}

// TestPinCharDataBufferIsReused pins the aliasing trap: the byte slice behind a
// CharData is overwritten by the next RawToken call, so the parser must copy the value
// out. It does, via the []byte-to-string conversion.
func TestPinCharDataBufferIsReused(t *testing.T) {
	d := xml.NewDecoder(strings.NewReader(`<r>aaaaaaaa<x/>bbbbbbbb</r>`))
	var held []xml.CharData
	for {
		tok, err := d.RawToken()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if cd, ok := tok.(xml.CharData); ok {
			held = append(held, cd) // deliberately not copied
		}
	}
	if len(held) != 2 {
		t.Fatalf("got %d CharData tokens, want 2", len(held))
	}
	if string(held[0]) == "aaaaaaaa" {
		t.Skip("encoding/xml no longer reuses the CharData buffer; the copy in charData is now redundant but harmless")
	}
	if string(held[0]) != string(held[1]) {
		t.Errorf("expected both retained slices to alias one buffer, got %q and %q", held[0], held[1])
	}
}

// TestPinAttributeValueNormalizationIsNotPerformed is the single most important pin:
// it is why attvalue.go exists. Go applies clause 2.11 line-ending normalization but
// not clause 3.3.3 whitespace normalization, and it decodes references in place, so
// xml.Attr.Value can no longer distinguish a literal tab from &#9;.
func TestPinAttributeValueNormalizationIsNotPerformed(t *testing.T) {
	toks := rawTokens(t, "<r a=\"x\ty\nz\r\nw\rv\" b=\"&#9;&#10;&#13;\"/>")
	se := toks[0].tok.(xml.StartElement)

	if got, want := se.Attr[0].Value, "x\ty\nz\nw\nv"; got != want {
		t.Errorf("Go attr a = %q, want %q (clause 2.11 applied, clause 3.3.3 not)", got, want)
	}
	if got, want := se.Attr[1].Value, "\t\n\r"; got != want {
		t.Errorf("Go attr b = %q, want %q (references decoded verbatim)", got, want)
	}
	// What Xerces produces for the same input, and what xmldom must produce.
	el := mustParseRoot(t, "<r a=\"x\ty\nz\r\nw\rv\" b=\"&#9;&#10;&#13;\"/>")
	if got, want := el.AttrValue("", "a"), "x y z w v"; got != want {
		t.Errorf("xmldom attr a = %q, want %q", got, want)
	}
	if got, want := el.AttrValue("", "b"), "\t\n\r"; got != want {
		t.Errorf("xmldom attr b = %q, want %q", got, want)
	}
}

// TestPinPredefinedEntitiesAndLineEndings pins the two decoding behaviours we rely on
// in character data: the five predefined entities are decoded, and clause 2.11 line
// ending normalization is applied.
func TestPinPredefinedEntitiesAndLineEndings(t *testing.T) {
	for _, tc := range []struct{ src, want string }{
		{`<r>&amp;&lt;&gt;&quot;&apos;</r>`, `&<>"'`},
		{`<r>&#65;&#x42;</r>`, "AB"},
		{"<r>a\r\nb\rc\nd</r>", "a\nb\nc\nd"},
	} {
		toks := rawTokens(t, tc.src)
		cd, ok := toks[1].tok.(xml.CharData)
		if !ok {
			t.Fatalf("%q: token 1 is %T", tc.src, toks[1].tok)
		}
		if string(cd) != tc.want {
			t.Errorf("%q -> %q, want %q", tc.src, string(cd), tc.want)
		}
	}
}

// TestPinWhatGoRejectsAndWhatItDoesNot separates the well-formedness errors we inherit
// from the ones xmldom has to add itself. Everything in the "accepted" list is a check
// implemented in parse.go.
func TestPinWhatGoRejectsAndWhatItDoesNot(t *testing.T) {
	rejected := map[string]string{
		"unknown entity":         `<r>&foo;</r>`,
		"]]> outside CDATA":      `<r>]]></r>`,
		"< in attribute value":   `<r a="<"/>`,
		"NUL character ref":      `<r>&#0;</r>`,
		"C0 character ref":       `<r>&#x1;</r>`,
		"U+FFFE character ref":   `<r>&#xFFFE;</r>`,
		"beyond U+10FFFF":        `<r>&#x110000;</r>`,
		"double dash in comment": `<r><!-- a -- b --></r>`,
		"unquoted attribute":     `<r a=1/>`,
		"bare ampersand":         `<r>a & b</r>`,
		"XML 1.1 declaration":    `<?xml version="1.1"?><r/>`,
	}
	for name, src := range rejected {
		d := xml.NewDecoder(strings.NewReader(src))
		var err error
		for err == nil {
			_, err = d.RawToken()
		}
		if err == io.EOF {
			t.Errorf("%s: encoding/xml accepted %q; xmldom now carries this check alone", name, src)
		}
	}

	accepted := map[string]string{
		"start/end tag mismatch":     `<r></s>`,
		"unclosed element at EOF":    `<r>`,
		"two root elements":          `<a/><b/>`,
		"trailing character data":    `<r/>tail`,
		"no root element":            `<!--only a comment-->`,
		"duplicate attribute":        `<r a="1" a="2"/>`,
		"undeclared element prefix":  `<p:r/>`,
		"undeclared attr prefix":     `<r x:y="1"/>`,
		"empty prefixed binding":     `<r xmlns:p=""/>`,
		"bogus xml prefix binding":   `<r xmlns:xml="urn:bogus"/>`,
		"surrogate character ref":    `<r>&#xD800;</r>`,
		"reserved PI target":         `<r><?xml v?></r>`,
		"DOCTYPE declaration":        `<!DOCTYPE r SYSTEM "x.dtd"><r/>`,
		"CDATA outside the root":     `<![CDATA[x]]><r/>`,
		"markup decl outside a DTD":  `<!ELEMENT r (#PCDATA)><r/>`,
		"xmlns bound to xmlns URI":   `<r xmlns:p="http://www.w3.org/2000/xmlns/"/>`,
		"other prefix on xml URI":    `<r xmlns:p="http://www.w3.org/XML/1998/namespace"/>`,
		"default ns bound to xml ns": `<r xmlns="http://www.w3.org/XML/1998/namespace"/>`,
	}
	for name, src := range accepted {
		d := xml.NewDecoder(strings.NewReader(src))
		d.CharsetReader = func(_ string, in io.Reader) (io.Reader, error) { return in, nil }
		var err error
		for err == nil {
			_, err = d.RawToken()
		}
		if err != io.EOF {
			t.Errorf("%s: encoding/xml now rejects %q with %v; the corresponding xmldom check may be dead code", name, src, err)
			continue
		}
		if _, perr := Parse([]byte(src), nil); perr == nil {
			t.Errorf("%s: xmldom accepted %q, want a SyntaxError", name, src)
		}
	}
}

// TestPinGreaterThanIsLegalInAttributeValues pins that '>' and even ']]>' are legal
// inside an attribute value and that both parsers accept them. It is why the start-tag
// rescanner is quote-aware instead of searching forward for '>'.
func TestPinGreaterThanIsLegalInAttributeValues(t *testing.T) {
	const src = `<r a=">" b="]]>"/>`
	el := mustParseRoot(t, src)
	if got := el.AttrValue("", "a"); got != ">" {
		t.Errorf("a = %q, want %q", got, ">")
	}
	if got := el.AttrValue("", "b"); got != "]]>" {
		t.Errorf("b = %q, want %q", got, "]]>")
	}
}

// TestPinIdentityCharsetReaderKeepsOffsets pins that supplying a pass-through
// CharsetReader for an already-transcoded buffer leaves InputOffset consistent, which
// is what lets decodeSource run before the tokenizer.
func TestPinIdentityCharsetReaderKeepsOffsets(t *testing.T) {
	const src = `<?xml version="1.0" encoding="ISO-8859-1"?><r a="é">é</r>`
	toks := rawTokens(t, src)
	want := []string{
		`<?xml version="1.0" encoding="ISO-8859-1"?>`,
		`<r a="é">`,
		"é",
		`</r>`,
	}
	if len(toks) != len(want) {
		t.Fatalf("got %d tokens, want %d", len(toks), len(want))
	}
	for i, w := range want {
		if toks[i].raw != w {
			t.Errorf("token %d raw = %q, want %q", i, toks[i].raw, w)
		}
	}
}
