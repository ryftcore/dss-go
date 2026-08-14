// Tests for the incremental serializer: §3.1's invariant, R15/R16 allocation
// and ordering, §3.4's style continuation, R13's trailer and R18's /ByteRange
// arithmetic and placeholder sizing.

package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"
)

func newFixtureUpdater(t *testing.T, spec fixtureSpec) (*Updater, []byte) {
	t.Helper()
	doc, raw := openFixture(t, spec)
	u, err := NewUpdater(doc)
	if err != nil {
		t.Fatalf("NewUpdater: %v", err)
	}
	return u, raw
}

// TestPrefixPreserved is §3.1's invariant, asserted on every shape of update.
func TestPrefixPreserved(t *testing.T) {
	for _, tc := range []struct {
		name string
		spec fixtureSpec
		run  func(u *Updater) error
	}{
		{"bare object", simpleSpec(), func(u *Updater) error {
			u.Add(Integer(1))
			return nil
		}},
		{"catalog edit", simpleSpec(), func(u *Updater) error {
			u.Catalog().Set("Extensions", DictOf(Name("ESIC"), DictOf(
				Name("BaseVersion"), Name("1.7"),
				Name("ExtensionLevel"), Integer(1),
			)))
			return nil
		}},
		{"signature", simpleSpec(), func(u *Updater) error {
			_, err := u.AddSignature(SignatureOptions{SubFilter: "ETSI.CAdES.detached"})
			return err
		}},
		{"xref stream source", fixtureStreamSpec(), func(u *Updater) error {
			_, err := u.AddSignature(SignatureOptions{SubFilter: "ETSI.CAdES.detached"})
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			u, raw := newFixtureUpdater(t, tc.spec)
			if err := tc.run(u); err != nil {
				t.Fatal(err)
			}
			res, err := u.Write()
			if err != nil {
				t.Fatal(err)
			}
			if int64(len(res.Bytes)) <= res.OriginalLength {
				t.Fatalf("output (%d) is not longer than the input (%d)", len(res.Bytes), res.OriginalLength)
			}
			if !bytes.Equal(res.Bytes[:len(raw)], raw) {
				t.Error("the original bytes were not preserved verbatim")
			}
			if res.OriginalLength != int64(len(raw)) {
				t.Errorf("OriginalLength = %d, want %d", res.OriginalLength, len(raw))
			}
		})
	}
}

// fixtureStreamSpec is simpleSpec with a cross-reference stream.
func fixtureStreamSpec() fixtureSpec {
	s := simpleSpec()
	s.Style = XRefStream
	return s
}

// TestDeterminism runs the identical script twice and demands identical bytes.
// R21: no clock, no PRNG, no map iteration, no locale.
func TestDeterminism(t *testing.T) {
	run := func() []byte {
		u, _ := newFixtureUpdater(t, emptyFieldSpec())
		if _, err := u.AddSignature(SignatureOptions{
			SubFilter:   "ETSI.CAdES.detached",
			FieldID:     "Signature1",
			SignerName:  "Alice",
			Reason:      "test",
			ContentSize: 256,
		}); err != nil {
			t.Fatal(err)
		}
		if err := u.SetDSSDictionary(DSSDictionary{
			Certs: []TokenRef{{Data: []byte{1, 2, 3}}, {Data: []byte{4, 5}}},
			CRLs:  []TokenRef{{Key: ObjectKey{Num: 3}}},
			VRI: []VRIEntry{
				{Name: "AAAA", Certs: []TokenRef{{Data: []byte{9}}}},
				{Name: "BBBB", OCSPs: []TokenRef{{Data: []byte{8}}}},
			},
		}); err != nil {
			t.Fatal(err)
		}
		res, err := u.Write()
		if err != nil {
			t.Fatal(err)
		}
		return res.Bytes
	}
	a, b := run(), run()
	if !bytes.Equal(a, b) {
		t.Error("two identical runs produced different bytes")
	}
}

// TestObjectAllocation is R16: numbers start at HighestObjectNumber()+1 and
// increase by one per Alloc/Add, in call order; generation is always 0.
func TestObjectAllocation_R16(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec()) // objects 1..3
	if got := u.doc.HighestObjectNumber(); got != 3 {
		t.Fatalf("fixture highest object = %d, want 3", got)
	}
	k1 := u.Alloc()
	k2 := u.Add(Integer(1))
	k3 := u.Add(Integer(2))
	if k1.Num != 4 || k2.Num != 5 || k3.Num != 6 {
		t.Errorf("allocation = %v %v %v, want 4 5 6", k1, k2, k3)
	}
	if k1.Gen != 0 || k2.Gen != 0 {
		t.Error("generation must always be 0")
	}
	// Alloc reserves without scheduling.
	sched := u.Scheduled()
	if len(sched) != 2 || sched[0] != k2 || sched[1] != k3 {
		t.Errorf("Scheduled() = %v", sched)
	}
}

// TestWriteOrderAscending is R15's deliberate deviation from pdfbox.
func TestWriteOrderAscending_R15(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	// Schedule out of order.
	u.Put(ObjectKey{Num: 9}, Integer(9))
	u.Put(ObjectKey{Num: 5}, Integer(5))
	u.Put(ObjectKey{Num: 7}, Integer(7))
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	inc := string(res.Bytes[res.OriginalLength:])
	i5, i7, i9 := strings.Index(inc, "5 0 obj"), strings.Index(inc, "7 0 obj"), strings.Index(inc, "9 0 obj")
	if i5 < 0 || i7 < 0 || i9 < 0 {
		t.Fatalf("objects missing from the increment:\n%s", inc)
	}
	if !(i5 < i7 && i7 < i9) {
		t.Errorf("objects were not written in ascending order: 5@%d 7@%d 9@%d", i5, i7, i9)
	}
	// And no object stream is ever emitted on write (R15).
	if strings.Contains(inc, "/ObjStm") {
		t.Error("the increment contains an object stream")
	}
}

// TestXRefStyleContinuation is §3.4: a table stays a table, a stream stays a
// stream, and a hybrid file degrades to a table.
func TestXRefStyleContinuation(t *testing.T) {
	cases := []struct {
		name string
		spec fixtureSpec
		want XRefStyle
	}{
		{"table", simpleSpec(), XRefTable},
		{"stream", fixtureStreamSpec(), XRefStream},
		{"hybrid degrades", func() fixtureSpec { s := simpleSpec(); s.Hybrid = true; return s }(), XRefTable},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			u, _ := newFixtureUpdater(t, c.spec)
			if got := u.doc.HasHybridXRef(); got != c.spec.Hybrid {
				t.Fatalf("HasHybridXRef() = %v, want %v", got, c.spec.Hybrid)
			}
			u.Add(Integer(1))
			res, err := u.Write()
			if err != nil {
				t.Fatal(err)
			}
			if res.XRefStyle != c.want {
				t.Fatalf("style = %v, want %v", res.XRefStyle, c.want)
			}
			inc := res.Bytes[res.OriginalLength:]
			hasTable := bytes.Contains(inc, []byte("\nxref\n")) || bytes.HasPrefix(inc, []byte("xref\n"))
			hasTrailer := bytes.Contains(inc, []byte("trailer\n"))
			hasXRefStm := bytes.Contains(inc, []byte("/Type /XRef"))
			switch c.want {
			case XRefTable:
				if !hasTable || !hasTrailer || hasXRefStm {
					t.Errorf("expected a table+trailer increment, got:\n%s", inc)
				}
			case XRefStream:
				if hasTable || hasTrailer || !hasXRefStm {
					t.Errorf("expected an xref-stream increment, got:\n%s", inc)
				}
			}
			// The section we just wrote must be where startxref says it is.
			if res.StartXref < res.OriginalLength || res.StartXref >= int64(len(res.Bytes)) {
				t.Errorf("StartXref %d is outside the increment", res.StartXref)
			}
		})
	}
}

// TestStartXrefPointsAtTheNewSection reads the emitted startxref back out of
// the bytes and checks it names the section we wrote.
func TestStartXrefPointsAtTheNewSection(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	u.Add(Integer(1))
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`startxref\n(\d+)\n%%EOF\n$`).FindSubmatch(res.Bytes)
	if m == nil {
		t.Fatalf("tail not found at the end of:\n%s", res.Bytes[res.OriginalLength:])
	}
	var declared int64
	fmt.Sscanf(string(m[1]), "%d", &declared)
	if declared != res.StartXref {
		t.Fatalf("startxref = %d, Result.StartXref = %d", declared, res.StartXref)
	}
	if !bytes.HasPrefix(res.Bytes[declared:], []byte("xref\n")) {
		t.Errorf("startxref does not point at an xref keyword: %q", res.Bytes[declared:declared+16])
	}
}

// TestTrailerRules is R13.
func TestTrailerRules_R13(t *testing.T) {
	spec := simpleSpec()
	u, raw := newFixtureUpdater(t, spec)
	// A /DocChecksum in the source must not survive into our trailer.
	u.doc.Trailer().Set("DocChecksum", Name("stale"))
	u.doc.Trailer().Set("XRefStm", Integer(17))
	u.Add(Integer(1))
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	inc := string(res.Bytes[len(raw):])
	trailer := inc[strings.Index(inc, "trailer"):]
	if strings.Contains(trailer, "/DocChecksum") {
		t.Error("/DocChecksum was not dropped")
	}
	if strings.Contains(trailer, "/XRefStm") {
		t.Error("/XRefStm was not dropped from a table trailer")
	}
	if !strings.Contains(trailer, "/Prev 0") && !regexp.MustCompile(`/Prev \d+`).MatchString(trailer) {
		t.Errorf("/Prev missing from the trailer:\n%s", trailer)
	}
	// /Size is the highest object number carried or written, plus one.
	if !strings.Contains(trailer, "/Size 5") {
		t.Errorf("/Size should be 5 (highest written = 4):\n%s", trailer)
	}
	// /ID is forced direct, and R21 repeats the first element when no
	// DocumentID was supplied.
	if !strings.Contains(trailer, "/ID [<DEADBEEF> <DEADBEEF>]") {
		t.Errorf("/ID not written as a direct pair:\n%s", trailer)
	}
}

func TestTrailerDocumentID_R21(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	if _, err := u.AddSignature(SignatureOptions{
		SubFilter:  "ETSI.CAdES.detached",
		DocumentID: []byte{0x01, 0x02},
	}); err != nil {
		t.Fatal(err)
	}
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	inc := string(res.Bytes[res.OriginalLength:])
	if !strings.Contains(inc, "/ID [<DEADBEEF> <0102>]") {
		t.Errorf("second /ID element should be the supplied DocumentID:\n%s", inc)
	}
}

// TestByteRangeArithmetic is R18: the four values, their relationship to the
// /Contents span, and the fact that the signed spans exclude exactly the hex
// string and nothing else.
func TestByteRangeArithmetic_R18(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	if _, err := u.AddSignature(SignatureOptions{
		SubFilter:   "ETSI.CAdES.detached",
		ContentSize: 128,
	}); err != nil {
		t.Fatal(err)
	}
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	br := res.ByteRange
	total := int64(len(res.Bytes))

	if br[0] != 0 {
		t.Errorf("ByteRange[0] = %d, want 0", br[0])
	}
	if br[1] != res.ContentsOffset {
		t.Errorf("ByteRange[1] = %d, want ContentsOffset %d", br[1], res.ContentsOffset)
	}
	if br[2] != res.ContentsOffset+res.ContentsLength {
		t.Errorf("ByteRange[2] = %d, want %d", br[2], res.ContentsOffset+res.ContentsLength)
	}
	if br[3] != total-br[2] {
		t.Errorf("ByteRange[3] = %d, want %d", br[3], total-br[2])
	}
	if br[1]+br[3]+res.ContentsLength != total {
		t.Errorf("the two spans plus /Contents (%d+%d+%d) do not cover the file (%d)",
			br[1], br[3], res.ContentsLength, total)
	}
	// The excluded span is exactly the hex string, angle brackets included.
	if res.Bytes[res.ContentsOffset] != '<' {
		t.Errorf("ContentsOffset does not point at '<': %q", res.Bytes[res.ContentsOffset])
	}
	if res.Bytes[br[2]-1] != '>' {
		t.Errorf("the byte before ByteRange[2] is not '>': %q", res.Bytes[br[2]-1])
	}
	if res.ContentsLength != int64(2*128+2) {
		t.Errorf("ContentsLength = %d, want %d", res.ContentsLength, 2*128+2)
	}

	// SignedData must be the concatenation of the two spans and nothing else.
	got, err := io.ReadAll(res.SignedData())
	if err != nil {
		t.Fatal(err)
	}
	want := append(append([]byte{}, res.Bytes[:br[1]]...), res.Bytes[br[2]:br[2]+br[3]]...)
	if !bytes.Equal(got, want) {
		t.Errorf("SignedData() is %d bytes, want %d", len(got), len(want))
	}
	if bytes.Contains(got, []byte("0000000000000000")) {
		t.Error("SignedData covers part of the /Contents placeholder")
	}
}

// TestByteRangePlaceholderSizing is R18's 35 reserved bytes: the placeholder
// occupies exactly them, and the formatted value is padded with 0x20.
func TestByteRangePlaceholderSizing_R18(t *testing.T) {
	// The reserved form, before patching, is exactly 36 bytes on the wire:
	// '[' plus 34 content bytes plus ']'.
	var buf bytes.Buffer
	w := NewWriter(&buf)
	if err := w.WriteObject(reservedByteRange()); err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSuffix(buf.String(), "\n"); len(got) != 36 {
		t.Fatalf("placeholder %q is %d bytes, want 36", got, len(got))
	}

	u, _ := newFixtureUpdater(t, simpleSpec())
	if _, err := u.AddSignature(SignatureOptions{SubFilter: "ETSI.CAdES.detached", ContentSize: 64}); err != nil {
		t.Fatal(err)
	}
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`/ByteRange \[(.{35})`).FindSubmatch(res.Bytes)
	if m == nil {
		t.Fatalf("no /ByteRange in the increment:\n%s", res.Bytes[res.OriginalLength:])
	}
	field := string(m[1])
	want := fmt.Sprintf("0 %d %d %d]", res.ByteRange[1], res.ByteRange[2], res.ByteRange[3])
	if !strings.HasPrefix(field, want) {
		t.Errorf("patched field %q does not start with %q", field, want)
	}
	if pad := field[len(want):]; strings.Trim(pad, " ") != "" {
		t.Errorf("padding is %q, want only spaces", pad)
	}
	if len(field) != 35 {
		t.Errorf("patched field is %d bytes, want 35", len(field))
	}
}

// TestByteRangeOverflow drives patchByteRange straight, because provoking a
// 36-digit /ByteRange from a real document would need a petabyte of input.
func TestByteRangeOverflow_R18(t *testing.T) {
	res := &Result{Bytes: make([]byte, 64)}
	watch := &sigWatch{
		contentsSeen:    true,
		byteRangeSeen:   true,
		contentsOffset:  1_000_000_000_000,
		contentsLength:  1_000_000_000_000,
		byteRangeOffset: 0,
		byteRangeLength: 10, // deliberately too small
	}
	err := patchByteRange(res, watch)
	if !errors.Is(err, ErrByteRangeTooLarge) {
		t.Fatalf("err = %v, want ErrByteRangeTooLarge", err)
	}
}

// TestUpdateIsIdempotent: two features editing one object must compose.
func TestUpdateReturnsOneInstance(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	a, err := u.Update(ObjectKey{Num: 3})
	if err != nil {
		t.Fatal(err)
	}
	b, err := u.Update(ObjectKey{Num: 3})
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Error("Update returned two different clones for one key")
	}
	if _, err := u.Update(ObjectKey{Num: 99}); !errors.Is(err, ErrNoSuchObject) {
		t.Errorf("Update of an absent key: %v, want ErrNoSuchObject", err)
	}
	// The clone must not be the document's own object.
	orig, _ := u.doc.Object(ObjectKey{Num: 3})
	if orig == a {
		t.Error("Update handed out the document's own dictionary")
	}
}

func TestCatalogIsScheduledOnce(t *testing.T) {
	u, _ := newFixtureUpdater(t, simpleSpec())
	c1 := u.Catalog()
	c2 := u.Catalog()
	if c1 != c2 {
		t.Fatal("Catalog() returned two instances")
	}
	c1.Set("Marker", Bool(true))
	res, err := u.Write()
	if err != nil {
		t.Fatal(err)
	}
	inc := string(res.Bytes[res.OriginalLength:])
	if strings.Count(inc, "1 0 obj") != 1 {
		t.Errorf("catalog written %d times", strings.Count(inc, "1 0 obj"))
	}
	if !strings.Contains(inc, "/Marker true") {
		t.Errorf("catalog edit missing:\n%s", inc)
	}
}
