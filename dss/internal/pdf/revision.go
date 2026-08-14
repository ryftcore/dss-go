// Revisions and /ByteRange geometry. See DESIGN.md §2.8.
//
// Provenance: eu.europa.esig.dss.pades.PAdESUtils.extractRevisions (DSS 6.5.RC1)
// and eu.europa.esig.dss.pdf.pdfbox.PdfBoxDocumentReader.isSignatureCoversWholeDocument.
//
// Two different notions of "revision" live here and must not be confused: the
// %%EOF scan below, which is the only one DSS uses to find previous revisions,
// and Document.XRefSections, which is the /Prev chain.

package pdf

import (
	"bufio"
	"errors"
	"io"
)

// Revision is one %%EOF-delimited prefix: the revision is bytes [0, End).
type Revision struct {
	Index int
	End   int64
}

var eofMarker = []byte("%%EOF")

// ScanRevisions reproduces PAdESUtils.extractRevisions byte for byte, including
// its %%EOF + EOL lookahead. It does not walk /Prev, and it deliberately counts
// %%EOF sequences that occur inside object data: PAdESUtils.getPreviousRevision
// picks the candidate whose length is the largest below byteRange[0]+byteRange[1],
// so a different revision list changes which document DSS reports as the signed
// original.
func ScanRevisions(r io.Reader) ([]Revision, error) {
	br := bufio.NewReader(r)
	var out []Revision
	var line []byte
	var position int64

	readByte := func() (byte, bool, error) {
		b, err := br.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return 0, false, nil
			}
			return 0, false, err
		}
		return b, true, nil
	}

	for {
		b, ok, err := readByte()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		position++
		line = append(line, b)

		if len(line) == len(eofMarker) && string(line) == string(eofMarker) {
			line = line[:0]
			eofPosition := position
			c, ok, err := readByte()
			if err != nil {
				return nil, err
			}
			if ok {
				position++
			}
			if ok && c == '\n' {
				eofPosition++
			} else if ok && c == '\r' {
				eofPosition++
				dch, ok2, err := readByte()
				if err != nil {
					return nil, err
				}
				if ok2 {
					position++
				}
				if ok2 && dch == '\n' {
					eofPosition++
				}
			}
			out = append(out, Revision{Index: len(out), End: eofPosition})
		} else if b == '\n' || b == '\r' || len(line) > len(eofMarker) {
			line = line[:0]
		}
	}
	return out, nil
}

// Revisions returns the %%EOF revision list of the source document. The result
// is computed once and cached.
func (d *Document) Revisions() []Revision {
	if d.revisions == nil {
		revs, err := ScanRevisions(bytesReader(d.data))
		if err != nil {
			revs = []Revision{}
		}
		if revs == nil {
			revs = []Revision{}
		}
		d.revisions = revs
	}
	out := make([]Revision, len(d.revisions))
	copy(out, d.revisions)
	return out
}

// SignedRanges returns the two covered spans as [start,end) pairs:
//
//	signed = [br[0], br[0]+br[1]) u [br[2], br[2]+br[3])
func SignedRanges(br []int64) ([2][2]int64, error) {
	var out [2][2]int64
	if len(br) < 4 {
		return out, errors.New("pdf: /ByteRange needs four values")
	}
	if br[0] < 0 || br[1] < 0 || br[2] < 0 || br[3] < 0 {
		return out, errors.New("pdf: /ByteRange has a negative value")
	}
	out[0] = [2]int64{br[0], br[0] + br[1]}
	out[1] = [2]int64{br[2], br[2] + br[3]}
	return out, nil
}

// ContentsRange returns the span of the /Contents hex string including its < >:
//
//	contents = [br[0]+br[1], br[2])
func ContentsRange(br []int64) ([2]int64, error) {
	var out [2]int64
	if len(br) < 4 {
		return out, errors.New("pdf: /ByteRange needs four values")
	}
	start := br[0] + br[1]
	end := br[2]
	if end < start {
		return out, errors.New("pdf: /ByteRange /Contents span is negative")
	}
	return [2]int64{start, end}, nil
}

// SignatureCoversWholeDocument reproduces
// PdfBoxDocumentReader.isSignatureCoversWholeDocument including its arithmetic,
// which is not the obvious one:
//
//	(br[1]-br[0]) + (br[2]-br[1]-br[0]) + br[3] == fileLength
//
// Do not "fix" the formula, or documents upstream reports as fully covered will
// stop matching.
func (d *Document) SignatureCoversWholeDocument(sd SignatureDictionary) bool {
	br := sd.ByteRange
	if len(br) < 4 {
		return false
	}
	before := br[1] - br[0]
	expectedCMS := br[2] - br[1] - br[0]
	after := br[3]
	return d.Size() == before+expectedCMS+after
}
