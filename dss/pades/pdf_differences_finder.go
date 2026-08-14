// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/modifications/PdfDifferencesFinder.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf.modifications is part of the eu.europa.esig.dss.pdf module that landed
// in no s5b manifest (see pdf_object.go's header). Method names (AnnotationOverlaps/
// PagesDifferences/VisualDifferences) confirmed against native_pdf_signature_service.go's
// already-landed call sites.
package pades

// PdfDifferencesFinder is used to encounter differences in pages between given PDF revisions.
// Port of the PdfDifferencesFinder interface.
type PdfDifferencesFinder interface {
	// AnnotationOverlaps returns a list of found annotation overlaps. Port of getAnnotationOverlaps.
	AnnotationOverlaps(reader PdfDocumentReader) []PdfModification

	// IsAnnotationBoxOverlapping checks if the given annotationBox overlaps with pdfAnnotations.
	// Port of isAnnotationBoxOverlapping.
	IsAnnotationBoxOverlapping(annotationBox AnnotationBox, pdfAnnotations []*PdfAnnotation) bool

	// PagesDifferences returns a list of missing/added pages between the signed and final
	// revisions. Port of getPagesDifferences.
	PagesDifferences(signedRevisionReader, finalRevisionReader PdfDocumentReader) []PdfModification

	// VisualDifferences returns a list of visual differences found between the signed and final
	// revisions, excluding newly created annotations. Port of getVisualDifferences.
	VisualDifferences(signedRevisionReader, finalRevisionReader PdfDocumentReader) []PdfModification
}
