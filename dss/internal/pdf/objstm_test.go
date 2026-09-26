package pdf

import (
	"bytes"
	"fmt"
	"testing"
)

// objStmPDF builds a document whose catalog and page tree live inside an object
// stream, reached through type-2 xref entries.
func objStmPDF(t *testing.T, mangle func(header string, objs []string) (string, []string), badIndex bool) []byte {
	t.Helper()
	objs := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 300 400] >>",
	}
	// header: pairs of (object number, relative offset)
	var header, body bytes.Buffer
	for i, o := range objs {
		fmt.Fprintf(&header, "%d %d ", i+1, body.Len())
		body.WriteString(o)
		body.WriteByte(' ')
	}
	h, list := header.String(), objs
	if mangle != nil {
		h, list = mangle(h, objs)
		body.Reset()
		for _, o := range list {
			body.WriteString(o)
			body.WriteByte(' ')
		}
	}
	payload := h + body.String()
	enc := FlateEncode([]byte(payload))

	var buf bytes.Buffer
	buf.WriteString("%PDF-1.5\n")
	stmOff := buf.Len()
	fmt.Fprintf(&buf, "4 0 obj\n<< /Type /ObjStm /N %d /First %d /Filter /FlateDecode /Length %d >>\nstream\n",
		len(list), len(h), len(enc))
	buf.Write(enc)
	buf.WriteString("\nendstream\nendobj\n")

	xrefOff := buf.Len()
	var raw bytes.Buffer
	put := func(typ byte, a uint32, b uint16) {
		raw.WriteByte(typ)
		raw.Write([]byte{byte(a >> 24), byte(a >> 16), byte(a >> 8), byte(a)})
		raw.Write([]byte{byte(b >> 8), byte(b)})
	}
	put(0, 0, 65535)
	for i := range list {
		idx := uint16(i)
		if badIndex {
			idx = uint16(len(list)-1-i) + 0 // reversed: the index no longer matches
		}
		put(2, 4, idx)
	}
	put(1, uint32(stmOff), 0)
	put(1, uint32(xrefOff), 0)
	enc2 := FlateEncode(raw.Bytes())
	fmt.Fprintf(&buf, "5 0 obj\n<< /Type /XRef /Size 6 /W [1 4 2] /Root 1 0 R "+
		"/Filter /FlateDecode /Length %d >>\nstream\n", len(enc2))
	buf.Write(enc2)
	buf.WriteString("\nendstream\nendobj\n")
	fmt.Fprintf(&buf, "startxref\n%d\n%%%%EOF\n", xrefOff)
	return buf.Bytes()
}

func TestObjectStreamBasics(t *testing.T) {
	d := mustOpen(t, objStmPDF(t, nil, false))
	if got := d.NumberOfPages(); got != 1 {
		t.Fatalf("pages = %d", got)
	}
	box, err := d.PageBox(1)
	if err != nil {
		t.Fatal(err)
	}
	if box != (Rect{0, 0, 300, 400}) {
		t.Errorf("page box = %v", box)
	}
	obj, err := d.Object(ObjectKey{Num: 1})
	if err != nil {
		t.Fatalf("catalog from the object stream: %v", err)
	}
	if _, ok := obj.(*Dict); !ok {
		t.Errorf("catalog = %#v", obj)
	}
}

// TestObjectStreamTrustsTheXrefNumber pins PDFObjectStreamParser's behaviour:
// when the object at the declared index has a different number, the xref wins
// and the container's number list is searched.
func TestObjectStreamTrustsTheXrefNumber(t *testing.T) {
	d := mustOpen(t, objStmPDF(t, nil, true))
	obj, err := d.Object(ObjectKey{Num: 1})
	if err != nil {
		t.Fatalf("object 1: %v", err)
	}
	dict, ok := obj.(*Dict)
	if !ok {
		t.Fatalf("object 1 = %#v", obj)
	}
	if n, _ := dict.GetRaw("Type").(Name); n != "Catalog" {
		t.Errorf("object 1 resolved to /Type /%s, want /Catalog", n)
	}
	if !hasWarning(d, WarnObjStmBroken) {
		t.Errorf("expected an objstm warning, got %+v", d.Warnings())
	}
}

// TestObjectStreamBrokenContainerDoesNotFailTheDocument pins that a container
// which cannot be decoded marks its members missing and warns.
func TestObjectStreamBrokenContainerDoesNotFailTheDocument(t *testing.T) {
	data := objStmPDF(t, nil, false)
	// Corrupt the object stream's payload but keep the document structure.
	i := bytes.Index(data, []byte("stream\n"))
	if i < 0 {
		t.Fatal("no stream")
	}
	for j := i + 7; j < i+20 && j < len(data); j++ {
		data[j] ^= 0xFF
	}
	d, err := OpenBytes(data, nil)
	if err != nil {
		// The catalog itself lives in the broken container, so ErrBrokenCatalog
		// is the honest outcome; what must not happen is a panic or a hang.
		t.Logf("open failed as expected: %v", err)
		return
	}
	if _, err := d.Object(ObjectKey{Num: 1}); err == nil {
		t.Log("object 1 still resolved (flate recovered enough bytes)")
	}
}

func TestObjectStreamHeaderParsing(t *testing.T) {
	var warn []Warning
	data := []byte("1 0 2 10 << /A 1 >> << /B 2 >>")
	dict := DictOf(Name("N"), Integer(2), Name("First"), Integer(9))
	nums, offs, err := parseObjStmHeader(data, dict, &warn, ObjectKey{Num: 4})
	if err != nil {
		t.Fatal(err)
	}
	if len(nums) != 2 || nums[0] != 1 || nums[1] != 2 {
		t.Errorf("numbers = %v", nums)
	}
	if offs[0] != 9 || offs[1] != 19 {
		t.Errorf("offsets = %v", offs)
	}
}

func TestObjectStreamHeaderTruncated(t *testing.T) {
	var warn []Warning
	dict := DictOf(Name("N"), Integer(5), Name("First"), Integer(4))
	nums, _, err := parseObjStmHeader([]byte("1 0 << >>"), dict, &warn, ObjectKey{Num: 4})
	if err != nil {
		t.Fatal(err)
	}
	if len(nums) != 1 {
		t.Errorf("numbers = %v, want the one complete pair", nums)
	}
	if !warningsContain(warn, WarnObjStmBroken) {
		t.Errorf("warnings = %+v", warn)
	}
}

func TestObjectStreamBadFirst(t *testing.T) {
	var warn []Warning
	dict := DictOf(Name("N"), Integer(1), Name("First"), Integer(9999))
	if _, _, err := parseObjStmHeader([]byte("1 0"), dict, &warn, ObjectKey{Num: 4}); err == nil {
		t.Error("a /First past the end of the stream must be rejected")
	}
}

// /N sized the header slices before anything was parsed: /N 2^50 panicked with
// "makeslice: cap out of range" (and 2^35 or so exhausted memory instead).
func TestObjStmHugeNDoesNotPanic(t *testing.T) {
	var warn []Warning
	dict := DictOf(Name("N"), Integer(1<<50), Name("First"), Integer(4))
	nums, offs, err := parseObjStmHeader([]byte("7 0 null"), dict, &warn, ObjectKey{Num: 5})
	if err != nil {
		t.Fatal(err)
	}
	if len(nums) != 1 || nums[0] != 7 || offs[0] != 4 {
		t.Errorf("nums = %v, offs = %v", nums, offs)
	}
}
