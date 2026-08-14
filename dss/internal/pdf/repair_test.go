package pdf

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
)

// TestLenient_X2_BruteForce pins the repair that saves
// validation/pdf-signed-corrupted.pdf: when an xref entry does not point at the
// object it claims, the whole table is replaced by a brute-force scan.
func TestLenient_X2_BruteForce(t *testing.T) {
	data := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "")
	// Shift every recorded offset by inserting bytes after the header.
	data = bytes.Replace(data, []byte("%PDF-1.4\n"),
		[]byte("%PDF-1.4\n% padding padding padding\n"), 1)
	d := mustOpen(t, data)
	if !hasWarning(d, WarnXRefBruteForce) {
		t.Errorf("expected a brute-force warning, got %+v", d.Warnings())
	}
	if got := d.NumberOfPages(); got != 1 {
		t.Errorf("pages = %d, want 1", got)
	}
	if len(d.ObjectKeys()) != 3 {
		t.Errorf("objects = %v, want the three scanned objects", d.ObjectKeys())
	}
}

// TestLenient_X3_BadStartxrefIsRepaired pins the repair that
// validation/pades_infinite_loop.pdf needs: a startxref past EOF is replaced by
// the nearest real xref found by brute force.
func TestLenient_X3_BadStartxrefIsRepaired(t *testing.T) {
	data := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "")
	idx := bytes.LastIndex(data, []byte("startxref\n"))
	end := bytes.Index(data[idx:], []byte("\n%%EOF"))
	fixed := append([]byte{}, data[:idx]...)
	fixed = append(fixed, []byte("startxref\n999999")...)
	fixed = append(fixed, data[idx+end:]...)

	d := mustOpen(t, fixed)
	if got := d.NumberOfPages(); got != 1 {
		t.Errorf("pages = %d, want 1", got)
	}
	if !hasWarning(d, WarnXRefOffsetRepaired) && !hasWarning(d, WarnXRefBruteForce) {
		t.Errorf("expected a repair warning, got %+v", d.Warnings())
	}
}

// TestLenient_T3_NoStartxref pins that a missing startxref triggers a full
// reconstruction rather than a failure.
func TestLenient_T3_NoStartxref(t *testing.T) {
	data := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "")
	data = bytes.Replace(data, []byte("startxref"), []byte("stortxref"), 1)
	d := mustOpen(t, data)
	if !hasWarning(d, WarnXRefBruteForce) {
		t.Errorf("expected a brute-force warning, got %+v", d.Warnings())
	}
	if got := d.NumberOfPages(); got != 1 {
		t.Errorf("pages = %d, want 1", got)
	}
}

// TestLenient_T1_MissingEOF pins that the startxref search runs to end of file
// when there is no %%EOF.
func TestLenient_T1_MissingEOF(t *testing.T) {
	data := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "")
	data = bytes.Replace(data, []byte("%%EOF\n"), []byte(""), 1)
	d := mustOpen(t, data)
	if !hasWarning(d, WarnMissingEOF) {
		t.Errorf("expected a missing-eof warning, got %+v", d.Warnings())
	}
	if got := d.NumberOfPages(); got != 1 {
		t.Errorf("pages = %d, want 1", got)
	}
}

// TestBruteForceLastOccurrenceWins pins that a later definition of an object
// number replaces an earlier one: later revisions are later in the file.
func TestBruteForceLastOccurrenceWins(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")
	buf.WriteString("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 1 1] >>\nendobj\n")
	buf.WriteString("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 2 2] >>\nendobj\n")
	// no xref at all: everything comes from the scan
	buf.WriteString("%%EOF\n")

	d := mustOpen(t, buf.Bytes())
	box, err := d.PageBox(1)
	if err != nil {
		t.Fatal(err)
	}
	if box != (Rect{0, 0, 2, 2}) {
		t.Errorf("page box = %v, want the last definition's", box)
	}
}

func TestBruteForceFindsTrailerItems(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")
	buf.WriteString("3 0 obj\n<< /Type /Page /Parent 2 0 R >>\nendobj\n")
	buf.WriteString("%%EOF\n")
	d := mustOpen(t, buf.Bytes())
	// The catalog was found by /Type /Catalog, with no trailer anywhere.
	cat, err := d.Catalog()
	if err != nil {
		t.Fatalf("catalog: %v", err)
	}
	if n, _ := d.GetName(cat, "Type"); n != "Catalog" {
		t.Errorf("catalog = %v", cat)
	}
}

func TestBruteForceRegistersObjectStreams(t *testing.T) {
	// An object stream whose contents are unreachable from any xref: the rebuild
	// must register its members (bfSearchForObjStreams).
	inner := "<< /Type /Catalog /Pages 2 0 R >> << /Type /Pages /Kids [] /Count 0 >>"
	header := "1 0 2 34 "
	payload := header + inner
	enc := FlateEncode([]byte(payload))
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.5\n")
	fmt.Fprintf(&buf, "9 0 obj\n<< /Type /ObjStm /N 2 /First %d /Filter /FlateDecode /Length %d >>\nstream\n",
		len(header), len(enc))
	buf.Write(enc)
	buf.WriteString("\nendstream\nendobj\n%%EOF\n")

	d, err := OpenBytes(buf.Bytes(), nil)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if _, err := d.Object(ObjectKey{Num: 1}); err != nil {
		t.Errorf("object 1 from the recovered object stream: %v", err)
	}
}

func TestFindXRefNearPicksTheClosestCandidate(t *testing.T) {
	data := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "")
	d := mustOpen(t, data)
	want := bytes.LastIndex(data, []byte("xref\n0 4"))
	got, ok := d.findXRefNear(int64(want) + 20)
	if !ok {
		t.Fatal("no candidate found")
	}
	if got != int64(want) {
		t.Errorf("findXRefNear = %d, want %d", got, want)
	}
}

func TestBruteForceGenerationIsBounded(t *testing.T) {
	// A generation number above 65535 cannot be one, so the scan must not record
	// it (and must not panic converting it).
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")
	buf.WriteString("1 999999 obj\n<< >>\nendobj\n")
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [] /Count 0 >>\nendobj\n")
	buf.WriteString("%%EOF\n")
	d := mustOpen(t, buf.Bytes())
	for _, k := range d.ObjectKeys() {
		if k.Gen > 65535 {
			t.Errorf("recorded an impossible generation: %v", k)
		}
	}
}

func TestBruteForceScanIsBoundedByTheLastEOF(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [] /Count 0 >>\nendobj\n")
	buf.WriteString("%%EOF\n")
	tail := buf.Len()
	buf.WriteString("7 0 obj\n<< /Beyond /TheLastEOF >>\nendobj\n")

	d := mustOpen(t, buf.Bytes())
	b := d.bf()
	if _, found := b.offsets[ObjectKey{Num: 7}]; found {
		t.Errorf("the scan went past the last %%%%EOF at %d", tail)
	}
}

func TestObjectKeyStringAndZero(t *testing.T) {
	if (ObjectKey{}).IsZero() != true {
		t.Error("the zero key must report IsZero")
	}
	if (ObjectKey{Num: 12}).String() != "12 0" {
		t.Errorf("ObjectKey.String = %q", (ObjectKey{Num: 12}).String())
	}
	if got := (Ref{Num: 3, Gen: 7}).Key(); got != (ObjectKey{Num: 3, Gen: 7}) {
		t.Errorf("Ref.Key = %v", got)
	}
	if got := strconv.Quote((ObjectKey{Num: 1, Gen: 2}).String()); got != `"1 2"` {
		t.Errorf("ObjectKey.String = %s", got)
	}
}
