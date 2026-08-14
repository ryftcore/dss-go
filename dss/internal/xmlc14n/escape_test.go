package xmlc14n

import (
	"bufio"
	"bytes"
	"testing"
)

func escaped(t *testing.T, fn func(*bufio.Writer, string), s string) string {
	t.Helper()
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	fn(w, s)
	if err := w.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	return buf.String()
}

func TestWriteAttrValueEscaped(t *testing.T) {
	// CanonicalizerBase.outputAttrToWriter: & < " #x9 #xA #xD, uppercase hex, lowercase x,
	// no leading zeroes. '>' and '\'' are NOT escaped in attribute values.
	for _, tc := range [][2]string{
		{"", ""},
		{"plain", "plain"},
		{"a&b", "a&amp;b"},
		{"a<b", "a&lt;b"},
		{`a"b`, "a&quot;b"},
		{"a>b", "a>b"},
		{"a'b", "a'b"},
		{"a\tb", "a&#x9;b"},
		{"a\nb", "a&#xA;b"},
		{"a\rb", "a&#xD;b"},
		{"\t\n\r", "&#x9;&#xA;&#xD;"},
		{"café \U0001D11E", "café \U0001D11E"},
	} {
		if got := escaped(t, writeAttrValueEscaped, tc[0]); got != tc[1] {
			t.Errorf("writeAttrValueEscaped(%q) = %q, want %q", tc[0], got, tc[1])
		}
	}
}

func TestWriteTextEscaped(t *testing.T) {
	// CanonicalizerBase.outputTextToWriter: & < > #xD. Tab and LF stay literal in text, which
	// is the mirror image of the attribute-value table.
	for _, tc := range [][2]string{
		{"a&b", "a&amp;b"},
		{"a<b", "a&lt;b"},
		{"a>b", "a&gt;b"},
		{`a"b`, `a"b`},
		{"a'b", "a'b"},
		{"a\tb", "a\tb"},
		{"a\nb", "a\nb"},
		{"a\rb", "a&#xD;b"},
		{"]]>", "]]&gt;"},
		{"a < b & c", "a &lt; b &amp; c"}, // a CDATA section canonicalizes as escaped text
	} {
		if got := escaped(t, writeTextEscaped, tc[0]); got != tc[1] {
			t.Errorf("writeTextEscaped(%q) = %q, want %q", tc[0], got, tc[1])
		}
	}
}

func TestWriteCarriageReturnEscaped(t *testing.T) {
	// Comments and processing instructions escape #xD and nothing else: they recognize no
	// markup and no references, so "&#13;" written in a comment is six literal characters.
	for _, tc := range [][2]string{
		{"a&#13;b", "a&#13;b"},
		{"a<b>&c", "a<b>&c"},
		{"a\rb", "a&#xD;b"},
		{"a\nb", "a\nb"},
		{"café", "café"},
	} {
		if got := escaped(t, writeCarriageReturnEscaped, tc[0]); got != tc[1] {
			t.Errorf("writeCarriageReturnEscaped(%q) = %q, want %q", tc[0], got, tc[1])
		}
	}
}

func TestWriteRuneInvalidUTF8(t *testing.T) {
	// UtfHelpper writes '?' for a code point it cannot encode - in Java an unpaired
	// surrogate. A well-formed U+FFFD must survive as itself.
	if got := escaped(t, writeTextEscaped, "a\xffb"); got != "a?b" {
		t.Errorf("invalid byte = %q, want %q", got, "a?b")
	}
	if got := escaped(t, writeTextEscaped, "a�b"); got != "a�b" {
		t.Errorf("U+FFFD = %q, want it unchanged", got)
	}
	// A surrogate half encoded as WTF-8 is three bytes Go's decoder rejects one at a time.
	if got := escaped(t, writeTextEscaped, "a\xed\xa0\x80b"); got != "a???b" {
		t.Errorf("WTF-8 surrogate = %q, want %q", got, "a???b")
	}
}
