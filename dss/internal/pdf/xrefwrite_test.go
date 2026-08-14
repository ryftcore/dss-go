// Byte-level tests for cross-reference emission: R12 (table), R13 (trailer),
// R14 (tail) and §3.4's stream form.

package pdf

import (
	"bytes"
	"compress/zlib"
	"io"
	"strings"
	"testing"
)

func TestXRefRanges_R12(t *testing.T) {
	// COSWriter.getXRefRanges' own doc example: 0 1 2 5 6 7 8 10 -> 0 3 5 4 10 1.
	entries := []xrefEntry{}
	for _, n := range []int64{0, 1, 2, 5, 6, 7, 8, 10} {
		entries = append(entries, xrefEntry{Key: ObjectKey{Num: n}})
	}
	got := xrefRanges(entries)
	want := [][2]int64{{0, 3}, {5, 4}, {10, 1}}
	if len(got) != len(want) {
		t.Fatalf("ranges = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ranges = %v, want %v", got, want)
		}
	}
	if r := xrefRanges(nil); r != nil {
		t.Errorf("empty ranges = %v", r)
	}
}

func TestXRefTableBytes_R12(t *testing.T) {
	entries := []xrefEntry{
		{Key: ObjectKey{Num: 12}, Offset: 1234},
		{Key: ObjectKey{Num: 13}, Offset: 5678},
		freeHeadEntry(),
	}
	var buf bytes.Buffer
	w := NewWriter(&buf)
	if err := writeXRefTable(w, entries); err != nil {
		t.Fatal(err)
	}
	want := "xref\n" +
		"0 1\n" +
		"0000000000 65535 f\r\n" +
		"12 2\n" +
		"0000001234 00000 n\r\n" +
		"0000005678 00000 n\r\n"
	if buf.String() != want {
		t.Errorf("xref table\n got %q\nwant %q", buf.String(), want)
	}
	// Every entry is exactly 20 bytes, which is what makes an xref table
	// randomly addressable: each "\r\n" must be preceded by 18 entry bytes.
	s := buf.String()
	for i := 0; i+1 < len(s); i++ {
		if s[i] != '\r' || s[i+1] != '\n' {
			continue
		}
		start := i - 18
		if start < 0 {
			t.Fatalf("entry ending at %d is shorter than 20 bytes", i)
		}
		entry := s[start : i+2]
		if len(entry) != 20 {
			t.Fatalf("entry %q is %d bytes", entry, len(entry))
		}
		if entry[10] != ' ' || entry[16] != ' ' || (entry[17] != 'n' && entry[17] != 'f') {
			t.Errorf("entry %q is not %%010d SP %%05d SP [nf] CRLF", entry)
		}
	}
}

func TestXRefTableContiguousWithFreeHead(t *testing.T) {
	// Objects 1 and 2 are contiguous with the object-0 free head, so they share
	// one subsection.
	entries := []xrefEntry{
		{Key: ObjectKey{Num: 1}, Offset: 10},
		{Key: ObjectKey{Num: 2}, Offset: 20},
		freeHeadEntry(),
	}
	var buf bytes.Buffer
	if err := writeXRefTable(NewWriter(&buf), entries); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(buf.String(), "xref\n0 3\n") {
		t.Errorf("expected a single 0 3 subsection, got %q", buf.String())
	}
}

func TestXRefEntryPadding(t *testing.T) {
	if got := string(pad(0, 10)); got != "0000000000" {
		t.Errorf("pad(0,10) = %q", got)
	}
	if got := string(pad(65535, 5)); got != "65535" {
		t.Errorf("pad(65535,5) = %q", got)
	}
	if got := string(pad(12345678901, 10)); got != "12345678901" {
		t.Errorf("an over-wide offset must not be truncated: %q", got)
	}
}

func TestWriteTail_R14(t *testing.T) {
	var buf bytes.Buffer
	if err := writeTail(NewWriter(&buf), 4242); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "startxref\n4242\n%%EOF\n" {
		t.Errorf("tail = %q", buf.String())
	}
}

func TestXRefStreamStructure(t *testing.T) {
	entries := []xrefEntry{
		freeHeadEntry(),
		{Key: ObjectKey{Num: 1}, Offset: 300},
		{Key: ObjectKey{Num: 2}, Offset: 70000},
	}
	trailer := DictOf(
		Name("Root"), Ref{Num: 1},
		Name("DocChecksum"), Name("ignored"),
	)
	stm, err := buildXRefStream(entries, trailer, 4, 999)
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := stm.Dict.GetRaw("Type").(Name); got != "XRef" {
		t.Errorf("/Type = %v", got)
	}
	if got, _ := stm.Dict.GetRaw("Size").(Integer); got != 4 {
		t.Errorf("/Size = %v, want 4", got)
	}
	if got, _ := stm.Dict.GetRaw("Prev").(Integer); got != 999 {
		t.Errorf("/Prev = %v", got)
	}
	if stm.Dict.Has("DocChecksum") {
		t.Error("/DocChecksum must not travel into the xref stream")
	}
	w, _ := stm.Dict.GetRaw("W").(Array)
	if len(w) != 3 || w[0] != Integer(1) || w[2] != Integer(2) {
		t.Fatalf("/W = %v, want [1 n 2]", w)
	}
	w2 := int(w[1].(Integer))
	if w2 != 3 { // 70000 needs three bytes
		t.Errorf("/W[1] = %d, want 3", w2)
	}
	idx, _ := stm.Dict.GetRaw("Index").(Array)
	if len(idx) != 2 || idx[0] != Integer(0) || idx[1] != Integer(3) {
		t.Errorf("/Index = %v, want [0 3]", idx)
	}
	parms, _ := stm.Dict.GetRaw("DecodeParms").(*Dict)
	if parms == nil || parms.GetRaw("Predictor") != Object(Integer(12)) {
		t.Fatalf("/DecodeParms = %v", parms)
	}
	if parms.GetRaw("Columns") != Object(Integer(1+w2+2)) {
		t.Errorf("/Columns = %v, want %d", parms.GetRaw("Columns"), 1+w2+2)
	}

	// The payload must inflate, un-predict and describe the three entries.
	zr, err := zlib.NewReader(bytes.NewReader(stm.Raw))
	if err != nil {
		t.Fatalf("xref stream payload is not zlib: %v", err)
	}
	predicted, err := io.ReadAll(zr)
	if err != nil {
		t.Fatal(err)
	}
	columns := 1 + w2 + 2
	rows := len(predicted) / (columns + 1)
	if rows != 3 {
		t.Fatalf("decoded %d rows, want 3", rows)
	}
	raw := make([]byte, 0, rows*columns)
	prev := make([]byte, columns)
	for r := 0; r < rows; r++ {
		row := predicted[r*(columns+1) : (r+1)*(columns+1)]
		if row[0] != 2 {
			t.Fatalf("row %d filter tag = %d, want 2 (Up)", r, row[0])
		}
		cur := make([]byte, columns)
		for i := 0; i < columns; i++ {
			cur[i] = row[1+i] + prev[i]
		}
		raw = append(raw, cur...)
		prev = cur
	}
	get := func(row, off, width int) int64 {
		var v int64
		for i := 0; i < width; i++ {
			v = v<<8 | int64(raw[row*columns+off+i])
		}
		return v
	}
	if get(0, 0, 1) != 0 || get(0, 1+w2, 2) != 65535 {
		t.Errorf("row 0 is not the free head: type=%d gen=%d", get(0, 0, 1), get(0, 1+w2, 2))
	}
	if get(1, 0, 1) != 1 || get(1, 1, w2) != 300 {
		t.Errorf("row 1 = type %d offset %d", get(1, 0, 1), get(1, 1, w2))
	}
	if get(2, 1, w2) != 70000 {
		t.Errorf("row 2 offset = %d", get(2, 1, w2))
	}
}

func TestPNGUpEncodeRoundTrip(t *testing.T) {
	raw := []byte{1, 2, 3, 10, 20, 30, 100, 200, 255}
	enc := pngUpEncode(raw, 3)
	if len(enc) != 12 {
		t.Fatalf("encoded length = %d, want 12", len(enc))
	}
	prev := make([]byte, 3)
	out := make([]byte, 0, 9)
	for r := 0; r < 3; r++ {
		row := enc[r*4 : (r+1)*4]
		if row[0] != 2 {
			t.Fatalf("row %d tag = %d", r, row[0])
		}
		cur := make([]byte, 3)
		for i := 0; i < 3; i++ {
			cur[i] = row[1+i] + prev[i]
		}
		out = append(out, cur...)
		prev = cur
	}
	if !bytes.Equal(out, raw) {
		t.Errorf("round trip = %v, want %v", out, raw)
	}
}
