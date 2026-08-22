// Native implementation of eu.europa.esig.dss.pades.validation.PdfObjectKey on internal/pdf.
//
// Behavioural reference: eu.europa.esig.dss.pdf.pdfbox.PdfBoxObjectKey, which wraps pdfbox's
// COSObjectKey. Here the wrapped value is internal/pdf's pdf.ObjectKey, whose (Num, Gen) pair is
// exactly a COSObjectKey. The type is a comparable struct, not a pointer type, because upstream
// uses PdfObjectKey as a HashMap key over its equals/hashCode - the ported sources key Go maps
// with it (see the FORWARD DEPENDENCIES note in pdf_composite_dss_dict_certificate_source.go).
package pades

import "github.com/ryftcore/dss-go/dss/internal/pdf"

// NativePdfObjectKey identifies an indirect PDF object by number and generation.
type NativePdfObjectKey struct {
	// value is the wrapped internal/pdf key.
	value pdf.ObjectKey
}

// NewNativePdfObjectKey wraps an internal/pdf object key.
// Port of the PdfBoxObjectKey(COSObjectKey) constructor.
func NewNativePdfObjectKey(value pdf.ObjectKey) NativePdfObjectKey {
	return NativePdfObjectKey{value: value}
}

// Value returns the wrapped key. Port of #getValue.
func (k NativePdfObjectKey) Value() any {
	return k.value
}

// Number returns the object number. Port of #getNumber.
func (k NativePdfObjectKey) Number() int64 {
	return k.value.Num
}

// Generation returns the generation number. Port of #getGeneration.
func (k NativePdfObjectKey) Generation() int {
	return int(k.value.Gen)
}

// String ports COSObjectKey#toString, i.e. "12 0".
func (k NativePdfObjectKey) String() string {
	return k.value.String()
}

// Compile-time assertion standing in for Java's "implements PdfObjectKey".
var _ PdfObjectKey = NativePdfObjectKey{}
