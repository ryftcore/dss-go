// Byte-level tests for the serialization primitives. Every expectation here is
// a literal, because R1..R11 are pinned: if one of these changes, a golden
// changed, and that needs justifying in the PR that did it.

package pdf

import (
	"bytes"
	"math"
	"strings"
	"testing"
)

func writeToString(t *testing.T, f func(w *Writer) error) string {
	t.Helper()
	var buf bytes.Buffer
	w := NewWriter(&buf)
	if err := f(w); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := w.Err(); err != nil {
		t.Fatalf("writer error: %v", err)
	}
	return buf.String()
}

func TestFormatReal_R8(t *testing.T) {
	// Java Float.toString semantics, then BigDecimal.stripTrailingZeros().
	// toPlainString() whenever that used an exponent (COSFloat.formatString).
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0.0"},
		{1, "1.0"},
		{-1, "-1.0"},
		{0.5, "0.5"},
		{-0.5, "-0.5"},
		{0.001, "0.001"},
		{0.0001, "0.0001"},  // Float.toString gives 1.0E-4
		{1e-7, "0.0000001"}, // 1.0E-7
		{1e7, "10000000"},   // 1.0E7
		{9999999, "9999999.0"},
		{1e8, "100000000"},
		{1e20, "100000000000000000000"},
		{595.276, "595.276"},
		{841.89, "841.89"},
		{-0.0, "0.0"}, //nolint:staticcheck // deliberate: the case pins that Go's untyped -0.0 constant is +0, which is why the real negative zero is exercised below through math.Copysign.
		{math.NaN(), "0.0"},
		{math.Inf(1), "0.0"},
		{math.Inf(-1), "0.0"},
	}
	for _, c := range cases {
		if got := FormatReal(c.in); got != c.want {
			t.Errorf("FormatReal(%v) = %q, want %q", c.in, got, c.want)
		}
	}
	if got := FormatReal(math.Copysign(0, -1)); got != "-0.0" {
		t.Errorf("FormatReal(-0.0) = %q, want %q", got, "-0.0")
	}
}

func TestWriteRealPrefersRaw_R8(t *testing.T) {
	got := writeToString(t, func(w *Writer) error {
		return w.WriteObject(Real{Val: 595.276, Raw: "595.27600"})
	})
	if got != "595.27600" {
		t.Errorf("Real.Raw not echoed verbatim: %q", got)
	}
}

func TestEncodeName_R5(t *testing.T) {
	cases := []struct {
		in   Name
		want string
	}{
		{"Type", "/Type"},
		{"Adobe.PPKLite", "/Adobe.PPKLite"},
		{"ETSI.CAdES.detached", "/ETSI.CAdES.detached"},
		{"Prop_Build", "/Prop_Build"},
		{"A B", "/A#20B"},
		{"Sig!", "/Sig#21"},
		{"a~b", "/a#7Eb"},
		{"a'b", "/a#27b"},
		{"a,b", "/a#2Cb"},
		{"1.2", "/1.2"},
		{"a#b", "/a#23b"},
		{"", "/"},
		{Name([]byte{0xce, 0xa9}), "/#CE#A9"}, // U+03A9, UTF-8
	}
	for _, c := range cases {
		if got := string(EncodeName(c.in)); got != c.want {
			t.Errorf("EncodeName(%q) = %q, want %q", string(c.in), got, c.want)
		}
	}
}

func TestEncodeString_R6(t *testing.T) {
	cases := []struct {
		in   String
		want string
	}{
		{String{Bytes: []byte("hello")}, "(hello)"},
		{String{Bytes: []byte("a(b)c")}, `(a\(b\)c)`},
		{String{Bytes: []byte(`a\b`)}, `(a\\b)`},
		{String{Bytes: []byte("D:20240101120000+00'00'")}, "(D:20240101120000+00'00')"},
		{String{Bytes: []byte("tab\there")}, "(tab\there)"}, // tab stays literal
		{String{Bytes: []byte("a\nb")}, "<610A62>"},         // LF forces hex
		{String{Bytes: []byte("a\rb")}, "<610D62>"},         // CR forces hex
		{String{Bytes: []byte{0x80}}, "<80>"},               // >= 0x80 forces hex
		{String{Bytes: []byte("ok"), Hex: true}, "<6F6B>"},  // Hex flag forces hex
		{String{Bytes: nil}, "()"},
		{String{Bytes: nil, Hex: true}, "<>"},
	}
	for _, c := range cases {
		if got := string(EncodeString(c.in)); got != c.want {
			t.Errorf("EncodeString(%q, hex=%v) = %q, want %q", c.in.Bytes, c.in.Hex, got, c.want)
		}
	}
}

func TestWriteEOLSuppressed_R2(t *testing.T) {
	got := writeToString(t, func(w *Writer) error {
		if err := w.WriteRaw([]byte("x")); err != nil {
			return err
		}
		for i := 0; i < 5; i++ {
			if err := w.WriteEOL(); err != nil {
				return err
			}
		}
		return nil
	})
	if got != "x\n" {
		t.Errorf("repeated WriteEOL emitted %q, want %q", got, "x\n")
	}
	// writeCRLF must NOT set the flag, exactly as pdfbox's write(byte[]) does:
	// that is what puts "trailer" hard against the last xref entry.
	got = writeToString(t, func(w *Writer) error {
		if err := w.writeCRLF(); err != nil {
			return err
		}
		return w.WriteEOL()
	})
	if got != "\r\n\n" {
		t.Errorf("writeCRLF then WriteEOL = %q, want %q", got, "\r\n\n")
	}
}

func TestWriteDict_R3(t *testing.T) {
	d := DictOf(
		Name("Type"), Name("Sig"),
		Name("Filter"), Name("Adobe.PPKLite"),
		Name("Count"), Integer(3),
	)
	// A nil value is skipped entirely (visitFromDictionary's dangling branch).
	d.Set("Missing", nil)
	got := writeToString(t, func(w *Writer) error { return w.WriteObject(d) })
	want := "<<\n/Type /Sig\n/Filter /Adobe.PPKLite\n/Count 3\n>>\n"
	if got != want {
		t.Errorf("dict bytes\n got %q\nwant %q", got, want)
	}
}

func TestWriteDictNested(t *testing.T) {
	d := DictOf(
		Name("A"), DictOf(Name("B"), Integer(1)),
		Name("C"), Integer(2),
	)
	got := writeToString(t, func(w *Writer) error { return w.WriteObject(d) })
	// The inner dictionary's trailing EOL suppresses the entry's own EOL (R2).
	want := "<<\n/A <<\n/B 1\n>>\n/C 2\n>>\n"
	if got != want {
		t.Errorf("nested dict\n got %q\nwant %q", got, want)
	}
}

func TestWriteArrayTenPerLine_R4(t *testing.T) {
	var a Array
	for i := 1; i <= 25; i++ {
		a = append(a, Integer(i))
	}
	got := writeToString(t, func(w *Writer) error { return w.WriteObject(a) })
	lines := strings.Split(strings.TrimSuffix(got, "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("25 items produced %d lines: %q", len(lines), got)
	}
	if lines[0] != "[1 2 3 4 5 6 7 8 9 10" {
		t.Errorf("line 1 = %q", lines[0])
	}
	if lines[1] != "11 12 13 14 15 16 17 18 19 20" {
		t.Errorf("line 2 = %q", lines[1])
	}
	if lines[2] != "21 22 23 24 25]" {
		t.Errorf("line 3 = %q", lines[2])
	}
	// Exactly ten items: pdfbox writes no separator after the last item, so no
	// trailing EOL sneaks in before the ']'.
	got = writeToString(t, func(w *Writer) error { return w.WriteObject(a[:10]) })
	if got != "[1 2 3 4 5 6 7 8 9 10]\n" {
		t.Errorf("ten items = %q", got)
	}
}

func TestWriteStream_R11(t *testing.T) {
	s := NewStream(DictOf(Name("Type"), Name("Test"), Name("Length"), Integer(999)), []byte("ABC"))
	got := writeToString(t, func(w *Writer) error { return w.WriteObject(s) })
	want := "<<\n/Type /Test\n/Length 3\n>>\nstream\r\nABC\r\nendstream\n"
	if got != want {
		t.Errorf("stream bytes\n got %q\nwant %q", got, want)
	}
	// /Length is rewritten in place, keeping its position in the dictionary,
	// and the caller's dictionary is not mutated.
	if v, ok := s.Dict.GetRaw("Length").(Integer); !ok || v != 999 {
		t.Errorf("writing a stream mutated the caller's /Length: %v", s.Dict.GetRaw("Length"))
	}
}

func TestWriteIndirect_R10(t *testing.T) {
	got := writeToString(t, func(w *Writer) error {
		return w.WriteIndirect(ObjectKey{Num: 12, Gen: 0}, Integer(7))
	})
	if got != "12 0 obj\n7\nendobj\n" {
		t.Errorf("indirect object = %q", got)
	}
}

func TestWriteReference_R9(t *testing.T) {
	got := writeToString(t, func(w *Writer) error { return w.WriteObject(Ref{Num: 6, Gen: 2}) })
	if got != "6 2 R" {
		t.Errorf("reference = %q", got)
	}
}

func TestWriteScalars(t *testing.T) {
	cases := []struct {
		in   Object
		want string
	}{
		{nil, "null"},
		{Null{}, "null"},
		{Bool(true), "true"},
		{Bool(false), "false"},
		{Integer(-42), "-42"},
		{Name("Sig"), "/Sig"},
	}
	for _, c := range cases {
		got := writeToString(t, func(w *Writer) error { return w.WriteObject(c.in) })
		if got != c.want {
			t.Errorf("WriteObject(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestWriterPosIsAbsolute(t *testing.T) {
	var buf bytes.Buffer
	w := newWriterAt(&buf, 1000)
	if w.Pos() != 1000 {
		t.Fatalf("base position = %d", w.Pos())
	}
	if err := w.WriteRaw([]byte("abc")); err != nil {
		t.Fatal(err)
	}
	if w.Pos() != 1003 {
		t.Errorf("position after 3 bytes = %d, want 1003", w.Pos())
	}
}
