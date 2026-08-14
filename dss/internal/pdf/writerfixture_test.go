// Hand-built PDF fixtures for the writer tests.
//
// The fixtures are emitted by the writer's own primitives, which is legitimate
// here for two reasons: those primitives are pinned byte for byte by R1..R14
// and tested directly against expected bytes in writer_test.go, and every
// fixture is a genuinely valid PDF that the reader (and pdfbox, via the KAT-C
// oracle) parses. What the fixtures give the writer tests is a document whose
// every byte is known in advance, so an assertion can name an exact offset.

package pdf

import (
	"bytes"
	"testing"
)

// fixtureObj is one indirect object of a fixture document.
type fixtureObj struct {
	Key ObjectKey
	Obj Object
}

// fixtureSpec describes a whole fixture document.
type fixtureSpec struct {
	Objs  []fixtureObj // ascending object number
	Root  ObjectKey
	Style XRefStyle // cross-reference style of the fixture's own last section
	ID    []byte    // /ID, both elements; nil omits the key
	// Hybrid marks the fixture as carrying an /XRefStm, which forces the
	// writer to degrade to a table (§3.4). The fixture itself stays a plain
	// table; only the flag matters to the writer.
	Hybrid bool
}

// build emits the fixture and returns its bytes together with the trailer and
// the startxref the writer will chain onto.
func (f fixtureSpec) build(t *testing.T) (raw []byte, trailer *Dict, startxref int64) {
	t.Helper()
	var buf bytes.Buffer
	w := NewWriter(&buf)
	// %PDF-1.4 plus pdfbox's four high-byte comment bytes, so tools treat the
	// file as binary.
	if err := w.WriteRaw([]byte("%PDF-1.4\n%\xf6\xe4\xfc\xdf\n")); err != nil {
		t.Fatalf("fixture header: %v", err)
	}

	entries := make([]xrefEntry, 0, len(f.Objs)+2)
	highest := int64(0)
	for _, o := range f.Objs {
		if o.Key.Num > highest {
			highest = o.Key.Num
		}
		entries = append(entries, xrefEntry{Key: o.Key, Offset: w.Pos()})
		if err := w.WriteIndirect(o.Key, o.Obj); err != nil {
			t.Fatalf("fixture object %s: %v", o.Key, err)
		}
	}
	entries = append(entries, freeHeadEntry())

	trailer = NewDict()
	trailer.Set("Size", Integer(highest+1))
	trailer.Set("Root", Ref{Num: f.Root.Num, Gen: f.Root.Gen})
	if f.ID != nil {
		trailer.Set("ID", Array{
			String{Bytes: f.ID, Hex: true},
			String{Bytes: f.ID, Hex: true},
		})
	}
	if f.Hybrid {
		// A genuine hybrid file: a table trailer whose /XRefStm names a
		// parallel cross-reference stream describing the same objects. Only a
		// table can be hybrid, which is why fixtureSpec rejects the
		// combination with XRefStream below.
		xrefKey := ObjectKey{Num: highest + 1}
		stmOffset := w.Pos()
		hybridEntries := append(append([]xrefEntry{}, entries...),
			xrefEntry{Key: xrefKey, Offset: stmOffset})
		stm, err := buildXRefStream(hybridEntries, trailer, xrefKey.Num+1, -1)
		if err != nil {
			t.Fatalf("fixture hybrid xref stream: %v", err)
		}
		if err := w.WriteIndirect(xrefKey, stm); err != nil {
			t.Fatalf("fixture hybrid xref stream write: %v", err)
		}
		entries = append(entries, xrefEntry{Key: xrefKey, Offset: stmOffset})
		trailer.Set("XRefStm", Integer(stmOffset))
		trailer.Set("Size", Integer(xrefKey.Num+1))
	}

	switch f.Style {
	case XRefStream:
		xrefKey := ObjectKey{Num: highest + 1}
		startxref = w.Pos()
		entries = append(entries, xrefEntry{Key: xrefKey, Offset: startxref})
		stm, err := buildXRefStream(entries, trailer, xrefKey.Num+1, -1)
		if err != nil {
			t.Fatalf("fixture xref stream: %v", err)
		}
		if err := w.WriteIndirect(xrefKey, stm); err != nil {
			t.Fatalf("fixture xref stream write: %v", err)
		}
	default:
		startxref = w.Pos()
		if err := writeXRefTable(w, entries); err != nil {
			t.Fatalf("fixture xref: %v", err)
		}
		if err := writeTrailer(w, trailer); err != nil {
			t.Fatalf("fixture trailer: %v", err)
		}
	}
	if err := writeTail(w, startxref); err != nil {
		t.Fatalf("fixture tail: %v", err)
	}
	if err := w.Err(); err != nil {
		t.Fatalf("fixture: %v", err)
	}
	return buf.Bytes(), trailer, startxref
}

// openFixture emits the fixture and parses it back with the reader, so every
// writer test is also a reader/writer round trip over bytes we control.
func openFixture(t *testing.T, f fixtureSpec) (*Document, []byte) {
	t.Helper()
	raw, _, _ := f.build(t)
	d, err := OpenBytes(raw, nil)
	if err != nil {
		t.Fatalf("OpenBytes(fixture): %v", err)
	}
	return d, raw
}

// simpleSpec is a one-page 1.4 document with a table cross-reference:
//
//	1 = catalog, 2 = page tree root, 3 = the page.
func simpleSpec() fixtureSpec {
	return fixtureSpec{
		Root:  ObjectKey{Num: 1},
		Style: XRefTable,
		ID:    []byte{0xde, 0xad, 0xbe, 0xef},
		Objs: []fixtureObj{
			{Key: ObjectKey{Num: 1}, Obj: DictOf(
				Name("Type"), Name("Catalog"),
				Name("Pages"), Ref{Num: 2},
			)},
			{Key: ObjectKey{Num: 2}, Obj: DictOf(
				Name("Type"), Name("Pages"),
				Name("Kids"), Array{Ref{Num: 3}},
				Name("Count"), Integer(1),
			)},
			{Key: ObjectKey{Num: 3}, Obj: DictOf(
				Name("Type"), Name("Page"),
				Name("Parent"), Ref{Num: 2},
				Name("MediaBox"), Array{Integer(0), Integer(0), Integer(595), Integer(842)},
			)},
		},
	}
}

// twoPageSpec adds a second page, for the widget-placement tests.
func twoPageSpec() fixtureSpec {
	s := simpleSpec()
	s.Objs[1].Obj = DictOf(
		Name("Type"), Name("Pages"),
		Name("Kids"), Array{Ref{Num: 3}, Ref{Num: 4}},
		Name("Count"), Integer(2),
	)
	s.Objs = append(s.Objs, fixtureObj{Key: ObjectKey{Num: 4}, Obj: DictOf(
		Name("Type"), Name("Page"),
		Name("Parent"), Ref{Num: 2},
		Name("MediaBox"), Array{Integer(0), Integer(0), Integer(595), Integer(842)},
	)})
	return s
}

// emptyFieldSpec is a document that already carries one empty signature field
// named "Signature1", with a /Lock, so the fill-in and FieldMDP branches have
// something to bite on.
func emptyFieldSpec() fixtureSpec {
	s := simpleSpec()
	s.Objs[0].Obj = DictOf(
		Name("Type"), Name("Catalog"),
		Name("Pages"), Ref{Num: 2},
		Name("AcroForm"), Ref{Num: 4},
	)
	s.Objs[2].Obj = DictOf(
		Name("Type"), Name("Page"),
		Name("Parent"), Ref{Num: 2},
		Name("MediaBox"), Array{Integer(0), Integer(0), Integer(595), Integer(842)},
		Name("Annots"), Array{Ref{Num: 5}},
	)
	s.Objs = append(s.Objs,
		fixtureObj{Key: ObjectKey{Num: 4}, Obj: DictOf(
			Name("Fields"), Array{Ref{Num: 5}},
			Name("SigFlags"), Integer(3),
		)},
		fixtureObj{Key: ObjectKey{Num: 5}, Obj: DictOf(
			Name("FT"), Name("Sig"),
			Name("Type"), Name("Annot"),
			Name("Subtype"), Name("Widget"),
			Name("T"), String{Bytes: []byte("Signature1")},
			Name("Rect"), Array{Integer(0), Integer(0), Integer(0), Integer(0)},
			Name("P"), Ref{Num: 3},
			Name("Lock"), DictOf(
				Name("Action"), Name("All"),
			),
		)},
	)
	return s
}

// signedSpec is emptyFieldSpec with the field filled, so "already signed" and
// "a DocMDP already exists" can be exercised.
func signedSpec() fixtureSpec {
	s := emptyFieldSpec()
	s.Objs[4].Obj = DictOf(
		Name("FT"), Name("Sig"),
		Name("Type"), Name("Annot"),
		Name("Subtype"), Name("Widget"),
		Name("T"), String{Bytes: []byte("Signature1")},
		Name("Rect"), Array{Integer(0), Integer(0), Integer(0), Integer(0)},
		Name("P"), Ref{Num: 3},
		Name("V"), Ref{Num: 6},
	)
	s.Objs = append(s.Objs, fixtureObj{Key: ObjectKey{Num: 6}, Obj: DictOf(
		Name("Type"), Name("Sig"),
		Name("Filter"), Name("Adobe.PPKLite"),
		Name("SubFilter"), Name("ETSI.CAdES.detached"),
		Name("Contents"), String{Bytes: []byte{0x30, 0x00}, Hex: true},
		Name("ByteRange"), Array{Integer(0), Integer(10), Integer(20), Integer(30)},
	)})
	return s
}
