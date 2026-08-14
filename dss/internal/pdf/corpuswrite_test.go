// Writer verification against real corpus documents.
//
// Every other writer test in this package drives the Updater from a hand-built
// fixture, which proves the formatting rules R1-R21 but says nothing about what
// happens when the input is a real-world PDF with an object-stream xref, a
// hybrid xref, an existing AcroForm, 54 prior revisions or AES encryption. This
// file closes that gap: it appends a signature increment to twelve
// representative corpus documents and asserts, per document,
//
//	(a) the output's prefix is byte-identical to the input (DESIGN.md §3.1) and
//	    every prior revision boundary survives unchanged;
//	(b) the Go reader parses the result, sees exactly one more signature field
//	    than the input had, and the new field's /ByteRange is self-consistent,
//	    covers the whole document and brackets the /Contents placeholder;
//	(c) the increment's xref style follows the §3.4 contract and its /Prev is
//	    the input's startxref.
//
// Set PDF_WRITE_OUT=<dir> to also drop each signed result there, which is what
// feeds the pdfbox side of the check (testdata/gen/PdfOracle.java over that
// directory; a Java-side parse plus field/ByteRange agreement is KAT-C's
// Go -> Java direction).

package pdf

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// writerCorpusCase names an upstream document and what makes it worth signing.
type writerCorpusCase struct {
	path string
	why  string
	pass string // password, "" for none
}

var writerCorpus = []writerCorpusCase{
	{"dss-pades/src/test/resources/EmptyPage.pdf", "xref stream, minimal, no AcroForm", ""},
	{"dss-pades/src/test/resources/doc.pdf", "xref table, 1.4", ""},
	{"dss-pades/src/test/resources/pdf-2.0.pdf", "header version 2.0", ""},
	{"dss-asic-xades/src/test/resources/bdoc-spec21.pdf", "hybrid xref: style must degrade to a table", ""},
	{"dss-pades/src/test/resources/validation/51sigs.pdf", "54 prior revisions, /DSS, many signature fields", ""},
	{"dss-pades/src/test/resources/pdf-xref-streams.pdf", "xref streams throughout", ""},
	{"dss-cookbook/src/main/resources/hello-world.pdf", "carries an empty field named ExistingSignatureField", ""},
	{"dss-pades/src/test/resources/protected/open_protected.pdf", "AES-128 (/AESV2) encrypted source, R20", " "},
	// /AESV3, and the user password is empty: the space is the owner password
	// only, so opening with it fails on both sides. The oracle records
	// password=- for this file.
	{"dss-pades/src/test/resources/protected/restricted_fields.pdf", "AES-256 (/AESV3) encrypted source, R20", ""},
	{"dss-pades/src/test/resources/validation/encrypted.pdf", "RC4-128 (V2/R3) encrypted source, R20", ""},
	{"dss-pades/src/test/resources/validation/doc-firmado-LT.pdf", "DSS-produced, already signed", ""},
	{"dss-cookbook/src/test/resources/snippets/25sigs.pdf", "50 prior revisions, object streams", ""},
}

func TestWriterOverCorpus(t *testing.T) {
	dir := os.Getenv("PDF_CORPUS_DIR")
	if dir == "" {
		t.Skip("set PDF_CORPUS_DIR to the upstream checkout to run the corpus writer sweep")
	}
	outDir := os.Getenv("PDF_WRITE_OUT")
	if outDir != "" {
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", outDir, err)
		}
	}

	for _, tc := range writerCorpus {
		t.Run(tc.path, func(t *testing.T) {
			in, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(tc.path)))
			if err != nil {
				t.Skipf("not in this checkout: %v", err)
			}
			t.Logf("why: %s", tc.why)

			// R20/R21: Options.Random is a stream, so determinism means "the
			// same seed", not "the same reader". Each run gets its own.
			newOpts := func() *Options {
				o := &Options{Random: fixedRandom()}
				if tc.pass != "" {
					o.Password = []byte(tc.pass)
				}
				return o
			}
			opts := newOpts()

			doc, err := OpenBytes(in, opts)
			if err != nil {
				t.Fatalf("OpenBytes: %v", err)
			}
			beforeFields, err := doc.SignatureFields()
			if err != nil {
				t.Fatalf("SignatureFields (input): %v", err)
			}
			beforeRevs := doc.Revisions()
			beforeSections := doc.XRefSections()
			if len(beforeSections) == 0 {
				t.Fatal("input has no xref sections")
			}
			inputStartXref := doc.StartXref()
			inputStyle := beforeSections[0].Style
			inputHybrid := doc.HasHybridXRef()

			out, res := signCorpusDocument(t, doc, in)

			// --- (a) prefix preservation and revision survival ---------------
			if len(out) <= len(in) {
				t.Fatalf("output %d bytes is not longer than the input's %d", len(out), len(in))
			}
			if !bytes.Equal(out[:len(in)], in) {
				t.Fatalf("prefix not preserved: first difference at byte %d",
					firstDifference(out[:len(in)], in))
			}
			if res.OriginalLength != int64(len(in)) {
				t.Errorf("Result.OriginalLength %d, input is %d bytes", res.OriginalLength, len(in))
			}

			afterRevs, err := scanRevisionBytes(out)
			if err != nil {
				t.Fatalf("ScanRevisions(out): %v", err)
			}
			if len(afterRevs) != len(beforeRevs)+1 {
				t.Errorf("revision count %d, want %d (input) + 1", len(afterRevs), len(beforeRevs))
			}
			for i := range beforeRevs {
				if i < len(afterRevs) && afterRevs[i].End != beforeRevs[i].End {
					t.Errorf("revision %d boundary moved: %d -> %d",
						i, beforeRevs[i].End, afterRevs[i].End)
				}
			}

			// --- (b) the result re-parses and the new field is there ---------
			out2, err := OpenBytes(out, newOpts())
			if err != nil {
				t.Fatalf("re-open signed output: %v", err)
			}
			afterFields, err := out2.SignatureFields()
			if err != nil {
				t.Fatalf("SignatureFields (output): %v", err)
			}
			if len(afterFields) != len(beforeFields)+1 {
				t.Fatalf("signature fields %d, want %d (input) + 1",
					len(afterFields), len(beforeFields))
			}
			checkNewSignature(t, out2, out, res, beforeFields)

			// --- (c) the increment's shape ----------------------------------
			// §3.4: a table stays a table and a hybrid degrades to one; only a
			// pure xref-stream input keeps the stream form.
			wantStyle := XRefStream
			if inputStyle == XRefTable || inputHybrid {
				wantStyle = XRefTable
			}
			if res.XRefStyle != wantStyle {
				t.Errorf("increment xref style %v, §3.4 requires %v (input style %v, hybrid %v)",
					res.XRefStyle, wantStyle, inputStyle, inputHybrid)
			}
			sections := out2.XRefSections()
			if len(sections) < 2 {
				t.Fatalf("signed output has %d xref sections, want at least 2", len(sections))
			}
			if sections[0].Style != wantStyle {
				t.Errorf("newest section style %v, want %v", sections[0].Style, wantStyle)
			}
			if sections[0].Prev != inputStartXref {
				t.Errorf("/Prev %d, input startxref was %d", sections[0].Prev, inputStartXref)
			}
			if out2.StartXref() != res.StartXref {
				t.Errorf("startxref %d, Result says %d", out2.StartXref(), res.StartXref)
			}
			// The whole prior chain must still be reachable and unchanged.
			if len(sections)-1 != len(beforeSections) {
				t.Errorf("chain length %d after the increment, input had %d",
					len(sections)-1, len(beforeSections))
			}
			for i, s := range beforeSections {
				if i+1 < len(sections) && sections[i+1].Offset != s.Offset {
					t.Errorf("prior section %d moved: offset %d -> %d",
						i, s.Offset, sections[i+1].Offset)
				}
			}

			// Determinism (R21): the same script twice, byte for byte.
			doc2, err := OpenBytes(in, newOpts())
			if err != nil {
				t.Fatalf("re-open input: %v", err)
			}
			again, _ := signCorpusDocument(t, doc2, in)
			if !bytes.Equal(out, again) {
				t.Errorf("non-deterministic: two runs differ at byte %d",
					firstDifference(out, again))
			}

			if outDir != "" {
				name := strings.ReplaceAll(tc.path, "/", "_")
				if err := os.WriteFile(filepath.Join(outDir, name), out, 0o644); err != nil {
					t.Fatalf("write output: %v", err)
				}
			}
		})
	}
}

// signCorpusDocument appends one invisible-signature increment with a fixed CMS.
func signCorpusDocument(t *testing.T, doc *Document, in []byte) ([]byte, *Result) {
	t.Helper()
	u, err := NewUpdater(doc)
	if err != nil {
		t.Fatalf("NewUpdater: %v", err)
	}
	if _, err := u.AddSignature(SignatureOptions{
		SubFilter:   "ETSI.CAdES.detached",
		SignerName:  "Corpus Writer KAT",
		Reason:      "audit",
		SigningTime: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		ContentSize: 4096,
	}); err != nil {
		t.Fatalf("AddSignature: %v", err)
	}
	res, err := u.Write()
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	cms := bytes.Repeat([]byte{0x30, 0x82, 0xAB, 0xCD}, 64)
	if err := res.InsertContents(cms); err != nil {
		t.Fatalf("InsertContents: %v", err)
	}
	return res.Bytes, res
}

// checkNewSignature validates the signature dictionary the increment added: it
// must be the only one that was not already in the input, its /ByteRange must
// be the four values Result promised, those values must bracket the /Contents
// string exactly, and the signature must cover the whole document.
func checkNewSignature(t *testing.T, out2 *Document, out []byte, res *Result, before []SignatureField) {
	t.Helper()
	known := make(map[string]bool, len(before))
	for _, f := range before {
		known[f.Name] = true
	}
	dicts, err := out2.SignatureDictionaries()
	if err != nil {
		t.Fatalf("SignatureDictionaries: %v", err)
	}
	fields, err := out2.SignatureFields()
	if err != nil {
		t.Fatalf("SignatureFields: %v", err)
	}
	var found *SignatureDictionary
	for i := range dicts {
		for _, fi := range dicts[i].Fields {
			if fi < len(fields) && !known[fields[fi].Name] {
				found = &dicts[i]
			}
		}
	}
	if found == nil {
		t.Fatal("the appended signature dictionary is not reachable from /AcroForm /Fields")
	}
	if found.SubFilter != "ETSI.CAdES.detached" {
		t.Errorf("/SubFilter %q, want ETSI.CAdES.detached", found.SubFilter)
	}
	if found.Filter != "Adobe.PPKLite" {
		t.Errorf("/Filter %q, want Adobe.PPKLite", found.Filter)
	}
	if len(found.ByteRange) != 4 {
		t.Fatalf("/ByteRange has %d values, want 4: %v", len(found.ByteRange), found.ByteRange)
	}
	br := found.ByteRange
	if [4]int64{br[0], br[1], br[2], br[3]} != res.ByteRange {
		t.Errorf("parsed /ByteRange %v, Result promised %v", br, res.ByteRange)
	}
	// The gap between the two signed spans is exactly the /Contents string,
	// angle brackets included (DESIGN.md §2.8).
	gapStart, gapEnd := br[0]+br[1], br[2]
	if gapStart != res.ContentsOffset {
		t.Errorf("/ByteRange gap starts at %d, /Contents at %d", gapStart, res.ContentsOffset)
	}
	if gapEnd-gapStart != res.ContentsLength {
		t.Errorf("/ByteRange gap is %d bytes, /Contents is %d", gapEnd-gapStart, res.ContentsLength)
	}
	if gapStart < 0 || gapEnd > int64(len(out)) {
		t.Fatalf("/ByteRange gap [%d,%d) is outside the %d-byte file", gapStart, gapEnd, len(out))
	}
	if out[gapStart] != '<' || out[gapEnd-1] != '>' {
		t.Errorf("the /ByteRange gap does not bracket a hex string: %q ... %q",
			out[gapStart], out[gapEnd-1])
	}
	if br[0] != 0 {
		t.Errorf("/ByteRange[0] is %d, an appended signature must start at 0", br[0])
	}
	if br[2]+br[3] != int64(len(out)) {
		t.Errorf("/ByteRange covers up to %d, the file is %d bytes", br[2]+br[3], len(out))
	}
	if !out2.SignatureCoversWholeDocument(*found) {
		t.Error("SignatureCoversWholeDocument is false for a freshly appended signature")
	}
}

// scanRevisionBytes is ScanRevisions over an in-memory document.
func scanRevisionBytes(b []byte) ([]Revision, error) {
	return ScanRevisions(bytesReader(b))
}

func firstDifference(a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return n
}

// fixedRandom is a deterministic byte source for AES initialisation vectors, so
// signing an encrypted source is reproducible (R20/R21).
func fixedRandom() *fixedReader { return &fixedReader{} }

type fixedReader struct{ n byte }

func (r *fixedReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = r.n
		r.n++
	}
	return len(p), nil
}
