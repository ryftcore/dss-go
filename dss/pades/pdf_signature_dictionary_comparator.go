// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfSignatureDictionaryComparator.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is the one Java package of dss-pades that landed in no s5b manifest
// (see pdf_object.go's header). Java's java.util.Comparator<PdfSignatureDictionary> becomes a
// Compare(a, b *PdfSignatureDictionary) int method, matching native_pdf_signature_service.go's
// already-landed call site (an explicit stable insertion sort, since Go's sort package is not
// guaranteed stable the way Java's Stream.sorted() is - see that call site's own comment).
// slf4j is dropped, per PORTING.md; the two WARN messages are kept as comments where they fired.
package pades

import "time"

// PdfSignatureDictionaryComparator is used to sort signatures by ByteRange.
// Port of the PdfSignatureDictionaryComparator class.
type PdfSignatureDictionaryComparator struct{}

// NewPdfSignatureDictionaryComparator is the default constructor.
func NewPdfSignatureDictionaryComparator() *PdfSignatureDictionaryComparator {
	return &PdfSignatureDictionaryComparator{}
}

// Compare ports #compare.
func (c *PdfSignatureDictionaryComparator) Compare(o1, o2 *PdfSignatureDictionary) int {
	byteRange1 := o1.ByteRange()
	byteRange2 := o2.ByteRange()

	begin1 := byteRange1.FirstPartStart()
	begin2 := byteRange2.FirstPartStart()

	// length = (before signature value) + (signature value) + (after signature value)
	length1 := byteRange1.Length()
	length2 := byteRange2.Length()

	end1 := byteRange1.FirstPartEnd()
	end2 := byteRange2.FirstPartEnd()

	switch {
	case begin1 >= begin2 && length1 < end2:
		// 2nd byterange envelops the whole 1st byterange
		return -1
	case begin2 >= begin1 && length2 < end1:
		// 1st byterange envelops the whole 2nd byterange
		return 1
	case byteRange1.Equals(byteRange2):
		// Upstream logs "More than one signature with the same byte range !" at WARN.
		return compareTime(o1.SigningDate(), o2.SigningDate())
	default:
		// Upstream logs "Strange byte ranges (ByteRange : {} / ByteRange : {})" at WARN.
		if end1 < end2 {
			return -1
		} else if end1 > end2 {
			return 1
		}
		return compareTime(o1.SigningDate(), o2.SigningDate())
	}
}

// compareTime ports java.util.Date#compareTo for the two SigningDate() comparisons above.
func compareTime(a, b time.Time) int {
	switch {
	case a.Before(b):
		return -1
	case a.After(b):
		return 1
	default:
		return 0
	}
}
