package xmldom

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"unicode/utf16"
)

func utf16Bytes(t *testing.T, s string, bigEndian, bom bool) []byte {
	t.Helper()
	var out []byte
	units := utf16.Encode([]rune(s))
	if bom {
		units = append([]uint16{0xFEFF}, units...)
	}
	for _, u := range units {
		hi, lo := byte(u>>8), byte(u)
		if bigEndian {
			out = append(out, hi, lo)
		} else {
			out = append(out, lo, hi)
		}
	}
	return out
}

func TestParseEncodings(t *testing.T) {
	tests := []struct {
		name string
		src  []byte
		want string // text content of the root
	}{{
		name: "plain UTF-8",
		src:  []byte("<r>aé€\U0001F600</r>"),
		want: "aé€\U0001F600",
	}, {
		name: "UTF-8 BOM is consumed",
		src:  append([]byte{0xEF, 0xBB, 0xBF}, []byte(`<r>a</r>`)...),
		want: "a",
	}, {
		name: "UTF-8 BOM with a declaration",
		src:  append([]byte{0xEF, 0xBB, 0xBF}, []byte(`<?xml version="1.0" encoding="UTF-8"?><r>a</r>`)...),
		want: "a",
	}, {
		name: "declared UTF-8",
		src:  []byte(`<?xml version="1.0" encoding="utf-8"?><r>é</r>`),
		want: "é",
	}, {
		name: "declared US-ASCII",
		src:  []byte(`<?xml version="1.0" encoding="US-ASCII"?><r>abc</r>`),
		want: "abc",
	}, {
		name: "declared ISO-8859-1 with high bytes",
		src:  []byte("<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><r>a\xe9\xff</r>"),
		want: "aéÿ",
	}, {
		name: "declared latin1 alias",
		src:  []byte("<?xml version=\"1.0\" encoding=\"latin1\"?><r>\xe9</r>"),
		want: "é",
	}, {
		name: "UTF-16BE with BOM",
		src:  utf16Bytes(t, `<r>aé\U0001F600</r>`, true, true),
		want: `aé\U0001F600`,
	}, {
		name: "UTF-16LE with BOM",
		src:  utf16Bytes(t, `<r>aé</r>`, false, true),
		want: "aé",
	}, {
		name: "UTF-16BE with BOM and a declaration",
		src:  utf16Bytes(t, `<?xml version="1.0" encoding="UTF-16"?><r>é</r>`, true, true),
		want: "é",
	}, {
		name: "declared UTF-16BE without a BOM",
		src:  append([]byte(nil), utf16Bytes(t, `<?xml version="1.0" encoding="UTF-16BE"?><r>é</r>`, true, false)...),
		want: "é",
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := Parse(tc.src, nil)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			if got := doc.DocumentElement().TextContent(); got != tc.want {
				t.Errorf("text = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestDeclaredEncodingWinsOverActualBytes pins the Java behaviour: the declaration
// decides, so UTF-8 bytes declared as ISO-8859-1 produce mojibake rather than an
// error. Reproducing it is what keeps character-level parity with Xerces.
func TestDeclaredEncodingWinsOverActualBytes(t *testing.T) {
	src := []byte(`<?xml version="1.0" encoding="ISO-8859-1"?><r>aé</r>`) // é is C3 A9 in UTF-8
	doc, err := Parse(src, nil)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got, want := doc.DocumentElement().TextContent(), "aÃ©"; got != want {
		t.Errorf("text = %q, want the mojibake %q", got, want)
	}
}

func TestEncodingErrors(t *testing.T) {
	tests := []struct {
		name string
		src  []byte
		want string
	}{
		{"unsupported encoding", []byte(`<?xml version="1.0" encoding="Shift_JIS"?><r/>`), "unsupported encoding"},
		{"UTF-16 declared without a BOM", []byte(`<?xml version="1.0" encoding="UTF-16"?><r/>`), "without a byte-order mark"},
		{"non-ASCII under US-ASCII", []byte("<?xml version=\"1.0\" encoding=\"US-ASCII\"?><r>\xe9</r>"), "not valid US-ASCII"},
		{"invalid UTF-8", []byte("<r>\xff\xfe\xfd</r>"), "invalid UTF-8"},
		{"odd UTF-16 length", []byte{0xFE, 0xFF, 0x00}, "odd byte count"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.src, nil)
			if err == nil {
				t.Fatal("Parse succeeded, want an error")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error = %q, want it to contain %q", err, tc.want)
			}
		})
	}
}

// TestCharsetReaderOption exercises the escape hatch for encodings outside the
// built-in table.
func TestCharsetReaderOption(t *testing.T) {
	// A toy single-byte code page: ASCII below 0x80 - which is what lets the
	// declaration be read in the first place - and 0x80+n mapped to U+2600+n above.
	codepage := func(charset string, in io.Reader) (io.Reader, error) {
		if charset != "x-cp-test" {
			t.Errorf("CharsetReader charset = %q, want %q", charset, "x-cp-test")
		}
		b, err := io.ReadAll(in)
		if err != nil {
			return nil, err
		}
		var out []byte
		for _, c := range b {
			if c < 0x80 {
				out = append(out, c)
				continue
			}
			out = append(out, []byte(string(rune(0x2600+int(c-0x80))))...)
		}
		return bytes.NewReader(out), nil
	}
	src := []byte("<?xml version=\"1.0\" encoding=\"x-cp-test\"?><r a=\"\x81\">\x80\x82</r>")
	doc, err := Parse(src, &ParseOptions{CharsetReader: codepage})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got := doc.DocumentElement().Name.Local; got != "r" {
		t.Errorf("root = %q, want %q", got, "r")
	}
	if got, want := doc.DocumentElement().TextContent(), "☀☂"; got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
	if got, want := doc.DocumentElement().AttrValue("", "a"), "☁"; got != want {
		t.Errorf("attribute a = %q, want %q", got, want)
	}

	// A CharsetReader that fails must surface as a parse error, not a panic.
	if _, err := Parse([]byte(`<?xml version="1.0" encoding="x-bad"?><r/>`), &ParseOptions{
		CharsetReader: func(string, io.Reader) (io.Reader, error) { return nil, io.ErrUnexpectedEOF },
	}); err == nil {
		t.Error("a failing CharsetReader was ignored")
	}
	if _, err := Parse([]byte(`<?xml version="1.0" encoding="x-bad"?><r/>`), &ParseOptions{
		CharsetReader: func(string, io.Reader) (io.Reader, error) { return nil, nil },
	}); err == nil {
		t.Error("a nil reader from CharsetReader was ignored")
	}
}

func TestParseXMLDecl(t *testing.T) {
	tests := []struct {
		src                           string
		present                       bool
		version, encoding, standalone string
	}{
		{src: `<r/>`},
		{src: `<?xml version="1.0"?><r/>`, present: true, version: "1.0"},
		{src: `<?xml version='1.0' encoding='UTF-8' standalone='yes'?><r/>`, present: true, version: "1.0", encoding: "UTF-8", standalone: "yes"},
		{src: "<?xml\tversion = \"1.0\"\n  encoding\t=\t\"ISO-8859-1\" ?><r/>", present: true, version: "1.0", encoding: "ISO-8859-1"},
		{src: `<?xml-stylesheet href="a"?><r/>`}, // not a declaration
		{src: `<?xml?><r/>`, present: true},
	}
	for _, tc := range tests {
		got := parseXMLDecl([]byte(tc.src))
		if got.present != tc.present || got.version != tc.version ||
			got.encoding != tc.encoding || got.standalone != tc.standalone {
			t.Errorf("parseXMLDecl(%q) = %+v, want present=%v version=%q encoding=%q standalone=%q",
				tc.src, got, tc.present, tc.version, tc.encoding, tc.standalone)
		}
	}
}

func TestDecodeLatin1(t *testing.T) {
	in := make([]byte, 256)
	for i := range in {
		in[i] = byte(i)
	}
	out := decodeLatin1(in)
	runes := []rune(string(out))
	if len(runes) != 256 {
		t.Fatalf("got %d runes, want 256", len(runes))
	}
	for i, r := range runes {
		if r != rune(i) {
			t.Fatalf("byte %d decoded to U+%04X, want U+%04X", i, r, i)
		}
	}
}
