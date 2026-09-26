package asn1ber

import (
	"bytes"
	"strings"
	"testing"
)

// nestedIndefinite returns depth indefinite-length SEQUENCEs nested inside one another,
// terminated by their end-of-contents octets.
func nestedIndefinite(depth int) []byte {
	return append(bytes.Repeat([]byte{0x30, 0x80}, depth), bytes.Repeat([]byte{0x00, 0x00}, depth)...)
}

// TestParseRejectsExcessiveNesting: Parse used to recurse once per nesting level with no
// bound, so tens of megabytes of "30 80" overflowed the goroutine stack - a fatal,
// unrecoverable runtime error - and a few hundred kilobytes already made DEREncoded run for
// tens of seconds (its cost is depth x size).
func TestParseRejectsExcessiveNesting(t *testing.T) {
	if _, rest, err := Parse(nestedIndefinite(MaxNestingDepth + 1)); err != nil || len(rest) != 0 {
		t.Fatalf("nesting at the limit: err %v, %d trailing bytes", err, len(rest))
	}
	if _, _, err := Parse(nestedIndefinite(MaxNestingDepth + 2)); err == nil {
		t.Fatal("nesting beyond the limit was accepted")
	}
	// Definite lengths nest the same way.
	deep := []byte{0x05, 0x00}
	for range MaxNestingDepth + 1 {
		deep = WriteSequence(deep)
	}
	if _, _, err := Parse(deep); err == nil {
		t.Fatal("definite-length nesting beyond the limit was accepted")
	}
	if _, _, err := Parse(nestedIndefinite(200000)); err == nil {
		t.Fatal("200000 nested SEQUENCEs were accepted")
	}
}

// TestValueToStringEscapesLongValues: the escaping used to insert one rune at a time into the
// buffer, which is quadratic - 200,000 commas took minutes. The output must stay what the
// insert loops produced.
func TestValueToStringEscapesLongValues(t *testing.T) {
	const count = 200000
	element, _, err := Parse(WriteTLV(TagUTF8String, bytes.Repeat([]byte{','}, count)))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := ValueToString(element), strings.Repeat(`\,`, count); got != want {
		t.Fatalf("got %d runes, want %d", len(got), len(want))
	}
	element, _, err = Parse(WriteTLV(TagUTF8String, bytes.Repeat([]byte{' '}, count)))
	if err != nil {
		t.Fatal(err)
	}
	// Every leading space is escaped, then the one trailing space of the result once more.
	if got, want := ValueToString(element), strings.Repeat(`\ `, count-1)+`\\ `; got != want {
		t.Fatalf("got %q..., want %q...", got[len(got)-8:], want[len(want)-8:])
	}

	testCases := []struct{ value, expected string }{
		{"  a,b  ", `\ \ a\,b\ \ `},
		{" ", `\\ `},
		{"#x ", `\#x\ `},
		{`a"b+c=d<e>f;g\h`, `a\"b\+c\=d\<e\>f\;g\\h`},
	}
	for _, testCase := range testCases {
		element, _, err := Parse(WriteTLV(TagUTF8String, []byte(testCase.value)))
		if err != nil {
			t.Fatal(err)
		}
		if got := ValueToString(element); got != testCase.expected {
			t.Errorf("ValueToString(%q) = %q, want %q", testCase.value, got, testCase.expected)
		}
	}
}

// TestParseBoundsTagNumbers: tag numbers used to be accepted up to 63 bits, so a
// context-specific [2^64-1] converted to int by the CMS layer read as -1, the tag number
// cmscore reserves for a plain Certificate. BouncyCastle's ASN1InputStream refuses any tag
// number beyond 31 bits ("Tag number more than 31 bits").
func TestParseBoundsTagNumbers(t *testing.T) {
	element, _, err := Parse([]byte{0x9F, 0x87, 0xFF, 0xFF, 0xFF, 0x7F, 0x00})
	if err != nil {
		t.Fatalf("tag number 2^31-1: %v", err)
	}
	if element.TagNumber() != 1<<31-1 {
		t.Fatalf("tag number = %d, want 2^31-1", element.TagNumber())
	}
	for _, input := range [][]byte{
		{0x9F, 0x88, 0x80, 0x80, 0x80, 0x00, 0x00},                               // 2^31
		{0x9F, 0x81, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F, 0x00}, // 2^64-1
	} {
		if _, _, err := Parse(input); err == nil {
			t.Errorf("% x: a tag number beyond 31 bits was accepted", input)
		}
	}
}
