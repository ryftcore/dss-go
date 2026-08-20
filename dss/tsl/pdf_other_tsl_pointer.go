// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/PDFOtherTSLPointer.java (DSS 6.5.RC1).
package tsl

// pdfOtherTSLPointerExpectedMimetype is the private static EXPECTED_MIMETYPE.
const pdfOtherTSLPointerExpectedMimetype = "application/pdf"

// PDFOtherTSLPointer allows filtering of TSL pointers by a PDF MimeType.
type PDFOtherTSLPointer struct {
	MimetypeOtherTSLPointer
}

// NewPDFOtherTSLPointer is the default constructor. Port of PDFOtherTSLPointer().
func NewPDFOtherTSLPointer() *PDFOtherTSLPointer {
	return &PDFOtherTSLPointer{MimetypeOtherTSLPointer: *NewMimetypeOtherTSLPointer(pdfOtherTSLPointerExpectedMimetype)}
}
