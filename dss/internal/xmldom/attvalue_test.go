package xmldom

import (
	"strings"
	"testing"
)

// TestAttributeValueNormalization is the clause 3.3.3 conformance table. The "want"
// column is what Xerces produces, which is what Santuario then canonicalizes.
func TestAttributeValueNormalization(t *testing.T) {
	tests := []struct {
		name, src, want string
	}{
		{"literal tab becomes a space", "<r a=\"x\ty\"/>", "x y"},
		{"literal LF becomes a space", "<r a=\"x\ny\"/>", "x y"},
		{"literal CR becomes a space", "<r a=\"x\ry\"/>", "x y"},
		{"literal CRLF becomes one space", "<r a=\"x\r\ny\"/>", "x y"},
		{"literal space is untouched", `<r a="x y"/>`, "x y"},
		{"runs are not collapsed", "<r a=\"x\t\t\ty\"/>", "x   y"},
		{"the design note example", "<r a=\"x\ty\nz\r\nw\rv\"/>", "x y z w v"},

		{"tab reference survives", `<r a="&#9;"/>`, "\t"},
		{"LF reference survives", `<r a="&#10;"/>`, "\n"},
		{"CR reference survives", `<r a="&#13;"/>`, "\r"},
		{"hex CR reference survives", `<r a="&#xD;"/>`, "\r"},
		{"the design note reference example", `<r a="&#9;&#10;&#13;"/>`, "\t\n\r"},
		{"space reference survives as a space", `<r a="&#32;"/>`, " "},

		{"predefined entities", `<r a="&amp;&lt;&gt;&quot;&apos;"/>`, `&<>"'`},
		{"decimal and hex references", `<r a="&#65;&#x42;&#x1F600;"/>`, "AB\U0001F600"},
		{"leading zeroes in a reference", `<r a="&#x0041;"/>`, "A"},

		{"gt is legal unescaped", `<r a=">"/>`, ">"},
		{"cdata end marker is legal", `<r a="]]>"/>`, "]]>"},
		{"apostrophe inside double quotes", `<r a="it's"/>`, "it's"},
		{"quote inside single quotes", `<r a='say "hi"'/>`, `say "hi"`},
		{"empty value", `<r a=""/>`, ""},
		{"multibyte value", `<r a="éюあ"/>`, "éюあ"},

		{"mixed literal and referenced whitespace", "<r a=\"\t&#9;\"/>", " \t"},
		{"reference next to a literal newline", "<r a=\"&#10;\n\"/>", "\n "},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			el := mustParseRoot(t, tc.src)
			if got := el.AttrValue("", "a"); got != tc.want {
				t.Errorf("Parse(%q) a = %q, want %q", tc.src, got, tc.want)
			}
		})
	}
}

// TestAttributeValueNormalizationSurvivesWhitespaceAroundEquals checks the scanner's
// handling of the optional whitespace the XML grammar allows around Eq, and of the two
// quote characters.
func TestAttributeValueNormalizationSurvivesWhitespaceAroundEquals(t *testing.T) {
	el := mustParseRoot(t, "<r\n  a = \"1\"\n  b\t=\t'2'\n  c='3'/>")
	for _, tc := range []struct{ name, want string }{{"a", "1"}, {"b", "2"}, {"c", "3"}} {
		if got := el.AttrValue("", tc.name); got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, got, tc.want)
		}
	}
	if len(el.Attrs) != 3 {
		t.Errorf("got %d attributes, want 3", len(el.Attrs))
	}
}

// TestAttributeOrderIsDocumentOrder pins that Attrs is document order, which is what
// the physical canonicalization method and the c14n tie-breaks depend on.
func TestAttributeOrderIsDocumentOrder(t *testing.T) {
	el := mustParseRoot(t, `<r xmlns:b="urn:b" xmlns:a="urn:a" b:z="1" a:z="2" z="3" a="4" xmlns="urn:d"/>`)
	want := []string{"xmlns:b", "xmlns:a", "b:z", "a:z", "z", "a", "xmlns"}
	if len(el.Attrs) != len(want) {
		t.Fatalf("got %d attributes, want %d", len(el.Attrs), len(want))
	}
	for i, w := range want {
		if got := el.Attrs[i].Name.QName(); got != w {
			t.Errorf("Attrs[%d] = %q, want %q", i, got, w)
		}
	}
}

func TestSplitQName(t *testing.T) {
	tests := []struct {
		in             string
		prefix, local  string
		wantErrMessage string
	}{
		{in: "a", local: "a"},
		{in: "p:a", prefix: "p", local: "a"},
		{in: "xmlns", local: "xmlns"},
		{in: "xmlns:p", prefix: "xmlns", local: "p"},
		{in: "a:b:c", wantErrMessage: "qualified name"},
		{in: ":a", wantErrMessage: "qualified name"},
		{in: "a:", wantErrMessage: "qualified name"},
		{in: "", wantErrMessage: "empty name"},
	}
	for _, tc := range tests {
		prefix, local, err := splitQName(tc.in)
		if tc.wantErrMessage != "" {
			if err == nil || !strings.Contains(err.Error(), tc.wantErrMessage) {
				t.Errorf("splitQName(%q) err = %v, want it to contain %q", tc.in, err, tc.wantErrMessage)
			}
			continue
		}
		if err != nil {
			t.Errorf("splitQName(%q) = %v", tc.in, err)
			continue
		}
		if prefix != tc.prefix || local != tc.local {
			t.Errorf("splitQName(%q) = (%q, %q), want (%q, %q)", tc.in, prefix, local, tc.prefix, tc.local)
		}
	}
}

func TestDecodeRefRejectsNonXMLCharacters(t *testing.T) {
	bad := []string{"#0", "#x0", "#x1", "#x8", "#xB", "#xC", "#xE", "#x1F",
		"#xD800", "#xDBFF", "#xDC00", "#xDFFF", "#xFFFE", "#xFFFF", "#x110000", "#xFFFFFFFF",
		// #4294967295 is the decimal form of #xFFFFFFFF: v > utf8.MaxRune must be
		// rejected before rune(v) is taken (CodeQL go/incorrect-integer-conversion
		// alert 36 — the check now runs before the narrowing, not after).
		"#4294967295",
		"#", "#x", "#xZZ", "#-1", "#+1", "#X41", "nbsp", "foo"}
	for _, b := range bad {
		if r, err := decodeRef(b); err == nil {
			t.Errorf("decodeRef(%q) = %q, want an error", b, r)
		}
	}
	good := map[string]rune{
		"amp": '&', "lt": '<', "gt": '>', "quot": '"', "apos": '\'',
		"#9": '\t', "#10": '\n', "#13": '\r', "#x20": ' ',
		"#xD7FF": 0xD7FF, "#xE000": 0xE000, "#xFFFD": 0xFFFD, "#x10FFFF": 0x10FFFF,
	}
	for b, want := range good {
		got, err := decodeRef(b)
		if err != nil {
			t.Errorf("decodeRef(%q) = %v", b, err)
			continue
		}
		if got != want {
			t.Errorf("decodeRef(%q) = %q, want %q", b, got, want)
		}
	}
}

// TestScanStartTagAgreesWithEncodingXML is the self-check described in the design
// note, run over a wide set of tags: the rescan must agree with encoding/xml on the
// attribute count and on every QName, or Parse fails closed.
func TestScanStartTagAgreesWithEncodingXML(t *testing.T) {
	tags := []string{
		`<r/>`,
		`<r >`,
		`<r a="1"/>`,
		`<r a='1' b="2" />`,
		`<a:r xmlns:a="urn:a" a:b="c"/>`,
		`<r a=">" b="]]>" c="&gt;"/>`,
		"<r\n\ta = \"1\"\n\tb\t=\t'2'\n/>",
		`<r xmlns="urn:d" xmlns:p="urn:1"/>`,
		`<r a="éюあ𐀀"/>`,
		`<r a="" b=''/>`,
	}
	for _, tag := range tags {
		src := tag
		if !strings.HasSuffix(tag, "/>") {
			src = tag + "</r>"
		}
		if _, err := Parse([]byte(src), nil); err != nil {
			t.Errorf("Parse(%q) = %v", src, err)
		}
	}
}

func TestCheckTextRefs(t *testing.T) {
	if err := checkTextRefs([]byte("a&#xD800;b")); err == nil {
		t.Error("surrogate reference in text was accepted")
	}
	if err := checkTextRefs([]byte("a&#xFFFD;b")); err != nil {
		t.Errorf("U+FFFD reference rejected: %v", err)
	}
	if err := checkTextRefs([]byte("plain &amp; text")); err != nil {
		t.Errorf("entity reference rejected: %v", err)
	}
	if err := checkTextRefs([]byte("no refs at all")); err != nil {
		t.Errorf("plain text rejected: %v", err)
	}
}
