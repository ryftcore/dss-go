package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestXRefTableChain(t *testing.T) {
	// A two-revision document: the second revision replaces object 3.
	base := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "")
	var buf bytes.Buffer
	buf.Write(base)
	objOff := buf.Len()
	buf.WriteString("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 10 10] >>\nendobj\n")
	xrefOff := buf.Len()
	fmt.Fprintf(&buf, "xref\n0 1\n0000000000 65535 f\r\n3 1\n%010d 00000 n\r\n", objOff)
	fmt.Fprintf(&buf, "trailer\n<< /Size 4 /Root 1 0 R /Prev %d >>\nstartxref\n%d\n%%%%EOF\n",
		bytes.Index(base, []byte("xref\n0 4")), xrefOff)

	d := mustOpen(t, buf.Bytes())
	secs := d.XRefSections()
	if len(secs) != 2 {
		t.Fatalf("sections = %d, want 2", len(secs))
	}
	if secs[0].Style != XRefTable || secs[1].Style != XRefTable {
		t.Errorf("styles = %v %v", secs[0].Style, secs[1].Style)
	}
	if secs[0].Offset != int64(xrefOff) {
		t.Errorf("newest section offset = %d, want %d", secs[0].Offset, xrefOff)
	}
	if secs[0].Prev != secs[1].Offset {
		t.Errorf("/Prev = %d, want %d", secs[0].Prev, secs[1].Offset)
	}
	// The newer object 3 must win.
	box, err := d.PageBox(1)
	if err != nil {
		t.Fatal(err)
	}
	if box != (Rect{0, 0, 10, 10}) {
		t.Errorf("page box = %v, want the newest revision's", box)
	}
	// The resolved trailer is the union of the chain's trailers, newest values
	// winning (XrefTrailerResolver.setStartxref).
	if !d.Trailer().Has("Prev") || !d.Trailer().Has("Root") || !d.Trailer().Has("Size") {
		t.Errorf("trailer = %v", d.Trailer().Keys())
	}
	if d.StartXref() != int64(xrefOff) {
		t.Errorf("StartXref = %d, want %d", d.StartXref(), xrefOff)
	}
}

// TestXRefPrevZeroEndsTheChain pins the behaviour that keeps
// validation/pades-5-signatures-and-1-document-timestamp.pdf at 8 objects:
// pdfbox's loop is `while (prev > 0)`.
func TestXRefPrevZeroEndsTheChain(t *testing.T) {
	data := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "/Prev 0\n")
	d := mustOpen(t, data)
	if got := len(d.XRefSections()); got != 1 {
		t.Errorf("sections = %d, want 1: /Prev 0 is not an offset", got)
	}
}

func TestXRefEntriesWithOffsetZeroAreSkipped(t *testing.T) {
	data := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "")
	// Rewrite object 3's entry to offset 0 while keeping it marked in-use.
	lines := bytes.Split(data, []byte("\r\n"))
	for i, l := range lines {
		if bytes.HasSuffix(l, []byte(" n")) && bytes.Contains(l, []byte("00000 n")) {
			if i == 3 { // 0 free, then objects 1, 2, 3
				lines[i] = []byte("0000000000 00000 n")
			}
		}
	}
	data = bytes.Join(lines, []byte("\r\n"))
	d := mustOpen(t, data)
	for _, k := range d.ObjectKeys() {
		if k.Num == 3 {
			// Object 3 may still be found by brute force, but never at offset 0.
			if e := d.xref[k]; e.typ == 1 && e.offset == 0 {
				t.Errorf("an xref entry with offset 0 was recorded")
			}
		}
	}
}

// TestLenient_PDFBOX474_ShortEntryLine pins the tolerance for `XXXX XXX XX n`
// style corruption: a line with fewer than three fields aborts the subsection.
func TestLenient_PDFBOX474_ShortEntryLine(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")
	off1 := buf.Len()
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	off2 := buf.Len()
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [] /Count 0 >>\nendobj\n")
	xrefOff := buf.Len()
	fmt.Fprintf(&buf, "xref\n0 3\n0000000000 65535 f\r\n%010d 00000 n\r\nBROKEN\r\n", off1)
	fmt.Fprintf(&buf, "trailer\n<< /Size 3 /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", xrefOff)
	_ = off2

	d, err := OpenBytes(buf.Bytes(), nil)
	if err != nil {
		t.Fatalf("a corrupt xref line must not fail the document: %v", err)
	}
	if !hasWarning(d, WarnXRefEntryInvalid) && !hasWarning(d, WarnXRefBruteForce) {
		t.Errorf("expected a recovery warning, got %+v", d.Warnings())
	}
	// The document still resolves through the brute-force repair.
	if _, err := d.Catalog(); err != nil {
		t.Errorf("catalog: %v", err)
	}
}

// --- xref streams -----------------------------------------------------------

// buildXRefStreamPDF emits a document whose cross-reference section is a stream.
func buildXRefStreamPDF(t *testing.T, extraDictEntries string) []byte {
	t.Helper()
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.5\n")
	offsets := map[int64]int{}
	writeObj := func(num int64, body string) {
		offsets[num] = buf.Len()
		fmt.Fprintf(&buf, "%d 0 obj\n%s\nendobj\n", num, body)
	}
	writeObj(1, "<< /Type /Catalog /Pages 2 0 R >>")
	writeObj(2, "<< /Type /Pages /Kids [3 0 R] /Count 1 >>")
	writeObj(3, "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 100 200] >>")

	xrefOff := buf.Len()
	// /W [1 4 2]: type, 4-byte offset, 2-byte generation.
	var raw bytes.Buffer
	put := func(typ byte, a uint32, b uint16) {
		raw.WriteByte(typ)
		raw.Write([]byte{byte(a >> 24), byte(a >> 16), byte(a >> 8), byte(a)})
		raw.Write([]byte{byte(b >> 8), byte(b)})
	}
	put(0, 0, 65535)
	put(1, uint32(offsets[1]), 0)
	put(1, uint32(offsets[2]), 0)
	put(1, uint32(offsets[3]), 0)
	put(1, uint32(xrefOff), 0)
	enc := FlateEncode(raw.Bytes())
	fmt.Fprintf(&buf, "4 0 obj\n<< /Type /XRef /Size 5 /W [1 4 2] /Root 1 0 R "+
		"/Filter /FlateDecode /Length %d %s>>\nstream\n", len(enc), extraDictEntries)
	buf.Write(enc)
	buf.WriteString("\nendstream\nendobj\n")
	fmt.Fprintf(&buf, "startxref\n%d\n%%%%EOF\n", xrefOff)
	return buf.Bytes()
}

func TestXRefStream(t *testing.T) {
	d := mustOpen(t, buildXRefStreamPDF(t, ""))
	secs := d.XRefSections()
	if len(secs) != 1 || secs[0].Style != XRefStream {
		t.Fatalf("sections = %+v", secs)
	}
	if got := len(d.ObjectKeys()); got != 4 {
		t.Errorf("objects = %d, want 4 (the free entry is not recorded)", got)
	}
	box, err := d.PageBox(1)
	if err != nil {
		t.Fatal(err)
	}
	if box != (Rect{0, 0, 100, 200}) {
		t.Errorf("page box = %v", box)
	}
	if !d.Trailer().Has("Root") {
		t.Errorf("the xref stream dictionary is the trailer: %v", d.Trailer().Keys())
	}
}

func TestXRefStreamDefaultIndexAndTypeWidthZero(t *testing.T) {
	// /W [0 4 2] means "type 1 for every entry".
	var raw bytes.Buffer
	put := func(a uint32, b uint16) {
		raw.Write([]byte{byte(a >> 24), byte(a >> 16), byte(a >> 8), byte(a)})
		raw.Write([]byte{byte(b >> 8), byte(b)})
	}
	put(10, 0)
	put(20, 1)
	d := &Document{opts: Options{}.withDefaults()}
	dict := DictOf(Name("W"), Array{Integer(0), Integer(4), Integer(2)},
		Name("Size"), Integer(2))
	entries, err := d.decodeXRefStream(dict, raw.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("entries = %+v", entries)
	}
	if entries[0].key != (ObjectKey{Num: 0}) || entries[0].e.offset != 10 {
		t.Errorf("entry 0 = %+v", entries[0])
	}
	if entries[1].key != (ObjectKey{Num: 1, Gen: 1}) {
		t.Errorf("entry 1 = %+v", entries[1])
	}
}

func TestXRefStreamRejectsBadW(t *testing.T) {
	d := &Document{opts: Options{}.withDefaults()}
	for _, w := range []Array{
		{Integer(1), Integer(2)},
		{Integer(1), Integer(-2), Integer(3)},
		{Integer(9), Integer(9), Integer(9)},
		{Integer(0), Integer(0), Integer(0)},
	} {
		if _, err := d.decodeXRefStream(DictOf(Name("W"), w, Name("Size"), Integer(1)), nil); err == nil {
			t.Errorf("/W %v was accepted", w)
		}
	}
}

func TestXRefStreamTruncatedDataStops(t *testing.T) {
	d := &Document{opts: Options{}.withDefaults()}
	dict := DictOf(Name("W"), Array{Integer(1), Integer(4), Integer(2)},
		Name("Size"), Integer(10))
	entries, err := d.decodeXRefStream(dict, []byte{1, 0, 0, 0, 5, 0, 0, 1, 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("entries = %d, want 1 from the one complete row", len(entries))
	}
}

// --- hybrid -----------------------------------------------------------------

// TestHybridXRefStmMergesBeneathTheTable builds a table whose trailer carries
// /XRefStm, and checks both that the stream's entries are picked up and that a
// table entry for the same object wins.
func TestHybridXRefStmMergesBeneathTheTable(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.5\n")
	offsets := map[int64]int{}
	writeObj := func(num int64, body string) {
		offsets[num] = buf.Len()
		fmt.Fprintf(&buf, "%d 0 obj\n%s\nendobj\n", num, body)
	}
	writeObj(1, "<< /Type /Catalog /Pages 2 0 R >>")
	writeObj(2, "<< /Type /Pages /Kids [3 0 R] /Count 1 >>")
	writeObj(3, "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 1 1] >>")
	// object 4 is reachable only through the /XRefStm
	writeObj(4, "<< /Only /InXRefStm >>")
	// a decoy for object 3 that only the /XRefStm points at: the table must win
	decoy := buf.Len()
	fmt.Fprintf(&buf, "3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 999 999] >>\nendobj\n")

	stmOff := buf.Len()
	var raw bytes.Buffer
	put := func(typ byte, a uint32, b uint16) {
		raw.WriteByte(typ)
		raw.Write([]byte{byte(a >> 24), byte(a >> 16), byte(a >> 8), byte(a)})
		raw.Write([]byte{byte(b >> 8), byte(b)})
	}
	put(1, uint32(decoy), 0)      // object 3 -> decoy
	put(1, uint32(offsets[4]), 0) // object 4
	enc := FlateEncode(raw.Bytes())
	fmt.Fprintf(&buf, "5 0 obj\n<< /Type /XRef /Size 6 /W [1 4 2] /Index [3 2] "+
		"/Root 1 0 R /Filter /FlateDecode /Length %d >>\nstream\n", len(enc))
	buf.Write(enc)
	buf.WriteString("\nendstream\nendobj\n")

	xrefOff := buf.Len()
	buf.WriteString("xref\n0 4\n0000000000 65535 f\r\n")
	fmt.Fprintf(&buf, "%010d 00000 n\r\n%010d 00000 n\r\n%010d 00000 n\r\n",
		offsets[1], offsets[2], offsets[3])
	fmt.Fprintf(&buf, "trailer\n<< /Size 6 /Root 1 0 R /XRefStm %d >>\nstartxref\n%d\n%%%%EOF\n",
		stmOff, xrefOff)

	d := mustOpen(t, buf.Bytes())
	if !d.HasHybridXRef() {
		t.Error("HasHybridXRef = false")
	}
	box, err := d.PageBox(1)
	if err != nil {
		t.Fatal(err)
	}
	if box != (Rect{0, 0, 1, 1}) {
		t.Errorf("page box = %v: the table entry must win over the /XRefStm", box)
	}
	if _, err := d.Object(ObjectKey{Num: 4}); err != nil {
		t.Errorf("object 4 (only in the /XRefStm) not found: %v", err)
	}
}

func TestHybridBadXRefStmIsSkipped(t *testing.T) {
	data := buildReaderPDF("%PDF-1.5\n", catalogObjs(), "/XRefStm 0\n")
	d := mustOpen(t, data)
	if !hasWarning(d, WarnXRefStmSkipped) {
		t.Errorf("expected an xrefstm-skipped warning, got %+v", d.Warnings())
	}
}

func TestSplitOnSpace(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want []string
	}{
		{"0000000016 00000 n", []string{"0000000016", "00000", "n"}},
		{"  0   1  ", []string{"0", "1"}},
		{"", nil},
		{"single", []string{"single"}},
	} {
		got := splitOnSpace(tc.in)
		if len(got) != len(tc.want) {
			t.Errorf("splitOnSpace(%q) = %v, want %v", tc.in, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("splitOnSpace(%q)[%d] = %q, want %q", tc.in, i, got[i], tc.want[i])
			}
		}
	}
}

func TestReadLineAt(t *testing.T) {
	l := newLexer([]byte("a\r\nb\nc\rd"), nil)
	for _, want := range []string{"a", "b", "c", "d"} {
		if got := readLineAt(l); got != want {
			t.Errorf("readLineAt = %q, want %q", got, want)
		}
	}
}

func TestXRefStyleString(t *testing.T) {
	if XRefTable.String() != "table" || XRefStream.String() != "stream" {
		t.Error("XRefStyle.String is part of the oracle's field format")
	}
	if XRefStyle(0).String() != "?" {
		t.Error("the zero style must render as ?")
	}
}

func TestSectionCapIsBounded(t *testing.T) {
	// A self-referential /Prev must not spin: the cycle guard ends the walk.
	data := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "")
	idx := bytes.LastIndex(data, []byte("startxref\n"))
	off := strings.TrimSpace(string(data[idx+len("startxref\n") : idx+len("startxref\n")+10]))
	data = bytes.Replace(data, []byte("/Root 1 0 R\n"),
		[]byte("/Root 1 0 R\n/Prev "+off+"\n"), 1)
	d := mustOpen(t, data)
	if len(d.XRefSections()) > 2 {
		t.Errorf("cycle guard failed: %d sections", len(d.XRefSections()))
	}
}

// An xref stream's /DecodeParms is decoded at Open, before anything else: a
// /Columns of 2^45 made the predictor allocate a 2^42-byte row and the process
// died with "out of memory". It is now a hard ErrLimitExceeded.
func TestXRefStreamHostilePredictorFailsOpen(t *testing.T) {
	data := buildXRefStreamPDF(t, "/DecodeParms << /Predictor 12 /Columns 35184372088832 >> ")
	_, err := OpenBytes(data, nil)
	if !errors.Is(err, ErrLimitExceeded) {
		t.Fatalf("err = %v, want ErrLimitExceeded", err)
	}
}

// Rows decoded from an xref stream are bounded by MaxObjects before they are
// allocated, not after mergeSections has built a map of them.
func TestXRefStreamEntriesBoundedByMaxObjects(t *testing.T) {
	d := &Document{opts: Options{MaxObjects: 10}.withDefaults()}
	raw := bytes.Repeat([]byte{1, 9}, 100) // /W [1 1 0]: 100 in-use entries
	dict := DictOf(Name("W"), Array{Integer(1), Integer(1), Integer(0)},
		Name("Size"), Integer(100))
	if _, err := d.decodeXRefStream(dict, raw); !errors.Is(err, ErrLimitExceeded) {
		t.Fatalf("err = %v, want ErrLimitExceeded", err)
	}
	if !errors.Is(d.limitErr, ErrLimitExceeded) {
		t.Error("the guard must be recorded so Open fails instead of repairing around it")
	}
}
