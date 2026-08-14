package pdf

import (
	"bytes"
	"compress/flate"
	"errors"
	"strings"
	"testing"
)

func TestFlateRoundTrip(t *testing.T) {
	for _, payload := range [][]byte{
		nil,
		[]byte("x"),
		[]byte(strings.Repeat("compress me ", 500)),
		bytes.Repeat([]byte{0}, 4096),
	} {
		enc := FlateEncode(payload)
		var warn []Warning
		got := FlateDecode(enc, &warn)
		if !bytes.Equal(got, payload) {
			t.Errorf("round trip lost data: %d in, %d out", len(payload), len(got))
		}
		if len(warn) != 0 {
			t.Errorf("clean round trip warned: %+v", warn)
		}
	}
}

// TestFlateDecodeSkipsTwoBytes is the load-bearing pdfbox quirk: the first two
// bytes are discarded unconditionally and the rest is inflated as raw DEFLATE,
// with no zlib header check and no Adler-32 validation. Do not "fix" this.
func TestFlateDecodeSkipsTwoBytes(t *testing.T) {
	payload := []byte("no zlib header here")
	var raw bytes.Buffer
	raw.WriteString("XX") // two bytes of garbage where the zlib header would be
	w, _ := flate.NewWriter(&raw, 6)
	w.Write(payload)
	w.Close()
	// no Adler-32 trailer at all
	var warn []Warning
	if got := FlateDecode(raw.Bytes(), &warn); !bytes.Equal(got, payload) {
		t.Errorf("decoded %q, want %q", got, payload)
	}
	if len(warn) != 0 {
		t.Errorf("a missing checksum must not warn: %+v", warn)
	}
}

func TestFlateDecodeShortInput(t *testing.T) {
	var warn []Warning
	for _, in := range [][]byte{nil, {}, {0x78}} {
		if got := FlateDecode(in, &warn); len(got) != 0 {
			t.Errorf("FlateDecode(%x) = %x, want empty", in, got)
		}
	}
	if len(warn) != 0 {
		t.Errorf("short input must not warn: %+v", warn)
	}
}

// TestFlateDecodeTruncatedKeepsPartialOutput pins rule 3 of §2.5: on a data
// format error the bytes decoded so far are returned with a warning, never an
// error.
func TestFlateDecodeTruncatedKeepsPartialOutput(t *testing.T) {
	// Deliberately poorly compressible, so half the deflate stream still yields
	// several hundred decoded bytes.
	payload := make([]byte, 4096)
	for i := range payload {
		payload[i] = byte(i*7 + i/13)
	}
	enc := FlateEncode(payload)
	truncated := enc[:len(enc)/2]
	var warn []Warning
	got := FlateDecode(truncated, &warn)
	if len(got) == 0 {
		t.Fatal("no partial output")
	}
	if !bytes.HasPrefix(payload, got) {
		t.Errorf("partial output is not a prefix of the input")
	}
	if !warningsContain(warn, WarnFlateTruncated) {
		t.Errorf("warnings = %+v", warn)
	}
}

func TestFlateDecodeGarbageDoesNotPanic(t *testing.T) {
	var warn []Warning
	FlateDecode([]byte{0x78, 0x9c, 0xff, 0xff, 0xff, 0xff}, &warn)
	FlateDecode(bytes.Repeat([]byte{0xAB}, 1000), &warn)
}

func TestASCIIHexDecode(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"48656C6C6F>", "Hello"},
		{"48 65\n6C6C6F>", "Hello"},
		{"48656C6C6F", "Hello"}, // missing EOD marker
		{"4>", "@"},             // odd digit padded with '0'
		{">", ""},
	} {
		if got := asciiHexDecode([]byte(tc.in)); string(got) != tc.want {
			t.Errorf("asciiHexDecode(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestASCII85Decode(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"87cURD]i,\"Ebo80~>", "Hello World!"},
		{"<~87cURD]i,\"Ebo80~>", "Hello World!"},
		{"z~>", "\x00\x00\x00\x00"},
		{"~>", ""},
	} {
		if got := ascii85Decode([]byte(tc.in)); string(got) != tc.want {
			t.Errorf("ascii85Decode(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRunLengthDecode(t *testing.T) {
	// 2 -> copy 3 literal bytes; 254 -> repeat the next byte 3 times; 128 -> EOD
	in := []byte{2, 'a', 'b', 'c', 254, 'z', 128, 'i', 'g', 'n'}
	if got := runLengthDecode(in); string(got) != "abczzz" {
		t.Errorf("runLengthDecode = %q, want %q", got, "abczzz")
	}
	if got := runLengthDecode([]byte{5, 'a'}); string(got) != "a" {
		t.Errorf("truncated literal run = %q", got)
	}
}

func TestLZWDecode(t *testing.T) {
	// The worked example from ISO 32000-1 §7.4.4.2 (Table 8), read backwards:
	// the encoded stream 80 0B 60 50 22 0C 0C 85 01 decodes to
	// 45 45 45 45 45 65 45 45 45 66, and exercises the KwKwK case at code 258.
	in := []byte{0x80, 0x0B, 0x60, 0x50, 0x22, 0x0C, 0x0C, 0x85, 0x01}
	want := []byte{45, 45, 45, 45, 45, 65, 45, 45, 45, 66}
	var warn []Warning
	if got := lzwDecode(in, 1, &warn); !bytes.Equal(got, want) {
		t.Errorf("lzwDecode = %v, want %v", got, want)
	}
}

func TestDecodeChainAndUnsupportedFilters(t *testing.T) {
	payload := []byte("chained")
	// Flate, then hex-encode the result: the chain decodes hex first, then flate.
	enc := FlateEncode(payload)
	var hexed strings.Builder
	const digits = "0123456789ABCDEF"
	for _, b := range enc {
		hexed.WriteByte(digits[b>>4])
		hexed.WriteByte(digits[b&0xF])
	}
	hexed.WriteByte('>')
	var warn []Warning
	got, err := Decode([]byte(hexed.String()), []Name{FilterASCIIHex, FilterFlate}, nil, &warn)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Errorf("chain decoded to %q", got)
	}

	for _, f := range []Name{"DCTDecode", "CCITTFaxDecode", "JPXDecode", "JBIG2Decode"} {
		_, err := Decode([]byte("whatever"), []Name{f}, nil, &warn)
		if !errors.Is(err, ErrUnsupportedFilter) {
			t.Errorf("Decode(%s) err = %v, want ErrUnsupportedFilter", f, err)
		}
		var fe *FilterError
		if !errors.As(err, &fe) || fe.Filter != f {
			t.Errorf("Decode(%s) error does not name the filter: %v", f, err)
		}
	}
}

func TestDecodeCryptFilter(t *testing.T) {
	var warn []Warning
	// /Crypt /Identity is a no-op; anything else is unsupported.
	got, err := Decode([]byte("data"), []Name{FilterCrypt},
		[]*Dict{DictOf(Name("Name"), Name("Identity"))}, &warn)
	if err != nil || string(got) != "data" {
		t.Errorf("identity crypt = %q %v", got, err)
	}
	if _, err := Decode([]byte("data"), []Name{FilterCrypt},
		[]*Dict{DictOf(Name("Name"), Name("StdCF"))}, &warn); !errors.Is(err, ErrUnsupportedFilter) {
		t.Errorf("non-identity crypt err = %v", err)
	}
}

func TestFilterAbbreviations(t *testing.T) {
	payload := []byte("abbrev")
	enc := FlateEncode(payload)
	var warn []Warning
	got, err := Decode(enc, []Name{"Fl"}, nil, &warn)
	if err != nil || !bytes.Equal(got, payload) {
		t.Errorf("abbreviated filter name failed: %q %v", got, err)
	}
}

// --- predictors -------------------------------------------------------------

func TestApplyPredictorPNG(t *testing.T) {
	// Three rows of four bytes, each row prefixed by its filter type.
	// Row 0: None; row 1: Up; row 2: Sub.
	rows := []byte{
		0, 1, 2, 3, 4,
		2, 1, 1, 1, 1, // Up: previous row + 1
		1, 5, 1, 1, 1, // Sub: left neighbour + value, bpp = 1
	}
	var warn []Warning
	got := ApplyPredictor(rows, 12, 1, 8, 4, &warn)
	want := []byte{
		1, 2, 3, 4,
		2, 3, 4, 5,
		5, 6, 7, 8,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("PNG predictor = %v, want %v", got, want)
	}
	if len(warn) != 0 {
		t.Errorf("warnings = %+v", warn)
	}
}

func TestApplyPredictorPNGAverageAndPaeth(t *testing.T) {
	rows := []byte{
		0, 10, 20, 30, 40,
		3, 0, 0, 0, 0, // Average of left and up
		4, 0, 0, 0, 0, // Paeth
	}
	var warn []Warning
	got := ApplyPredictor(rows, 15, 1, 8, 4, &warn)
	if len(got) != 12 {
		t.Fatalf("output = %d bytes, want 12", len(got))
	}
	// Average row: cur[i] = 0 + (left + up)/2
	wantAvg := []byte{5, 12, 21, 30}
	if !bytes.Equal(got[4:8], wantAvg) {
		t.Errorf("average row = %v, want %v", got[4:8], wantAvg)
	}
	// Paeth with a zero delta reproduces the predictor's choice
	wantPaeth := []byte{5, 12, 21, 30}
	if !bytes.Equal(got[8:12], wantPaeth) {
		t.Errorf("paeth row = %v, want %v", got[8:12], wantPaeth)
	}
}

func TestApplyPredictorTruncatedRowIsPadded(t *testing.T) {
	rows := []byte{0, 1, 2, 3, 4, 2, 1, 1} // second row is short
	var warn []Warning
	got := ApplyPredictor(rows, 12, 1, 8, 4, &warn)
	if len(got) != 8 {
		t.Errorf("output = %d bytes, want 8 (short row zero-padded)", len(got))
	}
	if !warningsContain(warn, WarnPredictorTruncated) {
		t.Errorf("warnings = %+v", warn)
	}
}

func TestApplyPredictorNoneAndTIFF(t *testing.T) {
	data := []byte{1, 2, 3, 4}
	var warn []Warning
	if got := ApplyPredictor(data, 1, 1, 8, 4, &warn); !bytes.Equal(got, data) {
		t.Errorf("predictor 1 changed the data: %v", got)
	}
	// TIFF predictor 2, 8 bits, one colour: each byte adds its left neighbour.
	tiff := []byte{1, 1, 1, 1}
	if got := ApplyPredictor(tiff, 2, 1, 8, 4, &warn); !bytes.Equal(got, []byte{1, 2, 3, 4}) {
		t.Errorf("TIFF predictor = %v, want [1 2 3 4]", got)
	}
	// Two colour components interleave.
	tiff2 := []byte{1, 2, 1, 2}
	if got := ApplyPredictor(tiff2, 2, 2, 8, 2, &warn); !bytes.Equal(got, []byte{1, 2, 2, 4}) {
		t.Errorf("TIFF predictor (2 colours) = %v, want [1 2 2 4]", got)
	}
	// Sub-byte components.
	if got := ApplyPredictor([]byte{0x55}, 2, 1, 4, 2, &warn); len(got) != 1 {
		t.Errorf("4-bit TIFF predictor = %v", got)
	}
}

func TestApplyPredictorDegenerateParameters(t *testing.T) {
	var warn []Warning
	data := []byte{1, 2, 3}
	for _, p := range []struct{ colors, bpc, columns int }{
		{0, 8, 1}, {1, 0, 1}, {1, 8, 0}, {-1, -1, -1},
	} {
		if got := ApplyPredictor(data, 12, p.colors, p.bpc, p.columns, &warn); got == nil {
			t.Errorf("ApplyPredictor(%v) returned nil", p)
		}
	}
}

// TestDecodeAppliesPredictorFromDecodeParms is the xref-stream path: /Predictor
// 12 with /Columns from /DecodeParms.
func TestDecodeAppliesPredictorFromDecodeParms(t *testing.T) {
	raw := []byte{2, 1, 1, 1, 1} // one Up row over an implicit zero row
	enc := FlateEncode(raw)
	parms := DictOf(Name("Predictor"), Integer(12), Name("Columns"), Integer(4))
	var warn []Warning
	got, err := Decode(enc, []Name{FilterFlate}, []*Dict{parms}, &warn)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, []byte{1, 1, 1, 1}) {
		t.Errorf("decode+predictor = %v", got)
	}
}

func TestStreamFiltersReadsAbbreviationsAndArrays(t *testing.T) {
	d := DictOf(Name("F"), Name("Fl"), Name("DP"), DictOf(Name("Predictor"), Integer(12)))
	names, parms := streamFilters(d, nil)
	if len(names) != 1 || names[0] != "Fl" || len(parms) != 1 {
		t.Errorf("abbreviated filter/parms = %v %v", names, parms)
	}
	d = DictOf(Name("Filter"), Array{Name("ASCII85Decode"), Name("FlateDecode")},
		Name("DecodeParms"), Array{Null{}, DictOf(Name("Predictor"), Integer(15))})
	names, parms = streamFilters(d, nil)
	if len(names) != 2 || len(parms) != 2 || parms[0] != nil || parms[1] == nil {
		t.Errorf("array filter/parms = %v %v", names, parms)
	}
}
