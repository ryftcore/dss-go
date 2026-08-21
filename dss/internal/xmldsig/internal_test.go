package xmldsig

import (
	"bytes"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
)

// decodeBase64 is java.util.Base64.getMimeDecoder().decode, and the corpus only exercises the
// happy path plus one truncated ds:DigestValue. These are the cases that decide whether a
// mangled signature is rejected or silently accepted with a shorter digest.
func TestDecodeBase64MatchesTheJavaMimeDecoder(t *testing.T) {
	for _, tc := range []struct {
		name  string
		in    string
		want  string
		fails bool
	}{
		{name: "plain", in: "aGVsbG8=", want: "hello"},
		{name: "no padding needed", in: "aGVsbG8h", want: "hello!"},
		{name: "two pad characters", in: "aGVsbG9vbw==", want: "hellooo"},
		{name: "unpadded final quantum of two", in: "aGVsbG9vbw", want: "hellooo"},
		{name: "unpadded final quantum of three", in: "aGVsbG8", want: "hello"},

		// The MIME decoder skips everything outside the alphabet, which is what makes an
		// indented, line-wrapped ds:DigestValue decode at all.
		{name: "newlines and indentation", in: "\n   aGVs\r\n\tbG8=\n  ", want: "hello"},
		{name: "junk between characters", in: "aG!Vs*bG8=", want: "hello"},

		// ... but a malformed final quantum is still an error. "aGVsbG8" is 7 characters, so a
		// '=' after it closes a quantum holding three characters and is fine; after two it
		// demands a second '=' and after one there are not enough bits at all.
		{name: "single character then padding", in: "aGVsbG8ha=", fails: true},
		{name: "two characters then one padding", in: "aGVsaG=", fails: true},
		{name: "two characters then two padding", in: "aGVsaG==", want: "helh"},
		{name: "padding at a quantum boundary", in: "aGVsbG8h=", fails: true},
		{name: "lone padding", in: "=", fails: true},

		{name: "empty", in: "", want: ""},
		{name: "non-alphabet junk after padding is ignored", in: "aGVsbG8=!!! \n", want: "hello"},
		{name: "base64 data after padding is an error", in: "aGVsbG8=!!!zzz", fails: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeBase64(tc.in)
			if tc.fails {
				if err == nil {
					t.Fatalf("decodeBase64(%q) = %q, want an error", tc.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("decodeBase64(%q): %v", tc.in, err)
			}
			if !bytes.Equal(got, []byte(tc.want)) {
				t.Fatalf("decodeBase64(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// stringFromNode is XMLUtils#getStrFromNode: given a text node it joins the character data of
// every text-node sibling, not just the run around it, so a ds:XPath expression interrupted by
// a comment is still read whole.
func TestStringFromNodeJoinsAllTextSiblings(t *testing.T) {
	doc, err := xmldom.Parse([]byte(`<x>not(<!--c-->ancestor-or-self::ds:Signature)</x>`), nil)
	if err != nil {
		t.Fatal(err)
	}
	first := doc.DocumentElement().FirstChild
	if got := stringFromNode(first); got != "not(ancestor-or-self::ds:Signature)" {
		t.Fatalf("stringFromNode = %q", got)
	}
}
