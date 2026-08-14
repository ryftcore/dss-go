package pdf

import (
	"bytes"
	"strings"
	"testing"
)

// TestScanRevisions_EOLLookahead pins the exact %%EOF + EOL arithmetic of
// PAdESUtils.extractRevisions. The boundary is the 1-based count of bytes
// consumed, extended by the EOL that follows: +1 for \n, +1 for \r, +2 for \r\n.
func TestScanRevisions_EOLLookahead(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   string
		want []int64
	}{
		{"bare", "%%EOF", []int64{5}},
		{"lf", "%%EOF\n", []int64{6}},
		{"cr", "%%EOF\r", []int64{6}},
		{"crlf", "%%EOF\r\n", []int64{7}},
		{"trailing junk", "%%EOFX", []int64{5}},
		{"two revisions", "a\n%%EOF\nb\n%%EOF\n", []int64{8, 16}},
		{"no eof", "no marker here", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ScanRevisions(strings.NewReader(tc.in))
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("revisions = %v, want %v", got, tc.want)
			}
			for i := range tc.want {
				if got[i].End != tc.want[i] {
					t.Errorf("revision %d ends at %d, want %d", i, got[i].End, tc.want[i])
				}
				if got[i].Index != i {
					t.Errorf("revision %d has index %d", i, got[i].Index)
				}
			}
		})
	}
}

// TestScanRevisions_CountsEOFInsideObjectData pins the quirk that must be
// reproduced: the scan is a byte scan, so a %%EOF inside a stream counts.
func TestScanRevisions_CountsEOFInsideObjectData(t *testing.T) {
	in := "1 0 obj\n<< >>\nstream\n%%EOF\nendstream\nendobj\n%%EOF\n"
	got, err := ScanRevisions(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("revisions = %v, want two (the one inside the stream counts)", got)
	}
}

// TestScanRevisions_BufferResets pins the two reset rules: any line-break byte
// resets the buffer, and so does exceeding five bytes. "%%%%EOF" therefore does
// not match.
func TestScanRevisions_BufferResets(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want int
	}{
		{"%%%%EOF", 0},
		{"x%%EOF", 0},   // six bytes without a line break: the buffer overflows first
		{"x\n%%EOF", 1}, // the line break resets, then %%EOF matches
		{"%%EO\nF", 0},
	} {
		got, err := ScanRevisions(strings.NewReader(tc.in))
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != tc.want {
			t.Errorf("ScanRevisions(%q) = %v, want %d revisions", tc.in, got, tc.want)
		}
	}
}

func TestDocumentRevisionsAreCached(t *testing.T) {
	data := buildReaderPDF("%PDF-1.4\n", catalogObjs(), "")
	d := mustOpen(t, data)
	a := d.Revisions()
	b := d.Revisions()
	if len(a) != 1 || len(b) != 1 {
		t.Fatalf("revisions = %v / %v", a, b)
	}
	if a[0].End != int64(len(data)) {
		t.Errorf("revision ends at %d, want %d", a[0].End, len(data))
	}
	// The returned slice is a copy: mutating it must not affect the document.
	a[0].End = 0
	if d.Revisions()[0].End == 0 {
		t.Error("Revisions() handed out its internal slice")
	}
}

func TestSignedRangesAndContentsRange(t *testing.T) {
	br := []int64{0, 840, 960, 240}
	spans, err := SignedRanges(br)
	if err != nil {
		t.Fatal(err)
	}
	if spans[0] != [2]int64{0, 840} || spans[1] != [2]int64{960, 1200} {
		t.Errorf("signed ranges = %v", spans)
	}
	c, err := ContentsRange(br)
	if err != nil {
		t.Fatal(err)
	}
	if c != [2]int64{840, 960} {
		t.Errorf("contents range = %v", c)
	}
	for _, bad := range [][]int64{nil, {0, 1, 2}, {-1, 1, 2, 3}} {
		if _, err := SignedRanges(bad); err == nil {
			t.Errorf("SignedRanges(%v) accepted", bad)
		}
	}
	if _, err := ContentsRange([]int64{0, 100, 50, 10}); err == nil {
		t.Error("ContentsRange accepted a negative span")
	}
}

// TestSignatureCoversWholeDocument pins the upstream formula, which is not the
// obvious one: (br[1]-br[0]) + (br[2]-br[1]-br[0]) + br[3].
func TestSignatureCoversWholeDocument(t *testing.T) {
	d := &Document{data: make([]byte, 1200)}
	// The obvious reading of a well-formed range: 0..840 and 960..1200.
	if !d.SignatureCoversWholeDocument(SignatureDictionary{ByteRange: []int64{0, 840, 960, 240}}) {
		t.Error("a well-formed range over the whole file must be covered")
	}
	// DESIGN.md §2.8 / checklist item 11: with a non-zero first element the
	// upstream formula subtracts br[0] *twice*, and reproducing that quirk —
	// rather than "fixing" it — is the whole point.
	//
	//	upstream: (br[1]-br[0]) + (br[2]-br[1]-br[0]) + br[3] = br[2]+br[3]-2*br[0]
	//	obvious:  (br[1]-br[0]) + (br[2]-br[1])       + br[3] = br[2]+br[3]-  br[0]
	//
	// For [10 840 960 240] those are 1180 and 1190, so a document of exactly
	// 1180 bytes separates them. The expectations below are literals on
	// purpose: deriving them from the same arithmetic the function uses would
	// make this test pass against either formula, which is how the quirk went
	// unpinned in the first place.
	quirk := SignatureDictionary{ByteRange: []int64{10, 840, 960, 240}}
	if !(&Document{data: make([]byte, 1180)}).SignatureCoversWholeDocument(quirk) {
		t.Error("upstream arithmetic must report a 1180-byte file as covered by [10 840 960 240]")
	}
	if (&Document{data: make([]byte, 1190)}).SignatureCoversWholeDocument(quirk) {
		t.Error("1190 bytes is the *obvious* formula's answer; the port must not use it")
	}
	if d.SignatureCoversWholeDocument(SignatureDictionary{ByteRange: []int64{0, 1}}) {
		t.Error("a short /ByteRange cannot cover the document")
	}
}

func TestScanRevisionsMatchesDocumentBytes(t *testing.T) {
	data := bytes.Join([][]byte{
		[]byte("%PDF-1.4\nrev one\n%%EOF\r\n"),
		[]byte("rev two\n%%EOF\n"),
	}, nil)
	revs, err := ScanRevisions(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(revs) != 2 {
		t.Fatalf("revisions = %v", revs)
	}
	if revs[0].End != 24 {
		t.Errorf("first revision ends at %d, want 24", revs[0].End)
	}
	if revs[1].End != int64(len(data)) {
		t.Errorf("second revision ends at %d, want %d", revs[1].End, len(data))
	}
}
