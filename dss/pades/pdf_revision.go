// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PdfRevision.java
// (DSS 6.5.RC1).
//
// java.io.Serializable is dropped (no Go counterpart), as elsewhere in this port.
//
// FORWARD DEPENDENCIES (sibling chunks; already referenced elsewhere in this package):
//   - PdfSignatureField (pdf_signature_field.go's header already documents its own shape).
//   - PdfModificationDetection - used already by native_pdf_signature_service.go
//     (pdfModificationDetection.SetAnnotationOverlaps/SetPageDifferences/SetVisualDifferences/
//     SetObjectModifications and NewPdfModificationDetection()); the getter this interface
//     needs is the matching ModificationDetection() *PdfModificationDetection.
package pades

// PdfRevision permits the user to choose the underlying PDF library used to create PDF
// signatures.
type PdfRevision interface {
	// PdfSigDictInfo returns a PDF Signature Dictionary info container.
	// Port of getPdfSigDictInfo().
	PdfSigDictInfo() *PdfSignatureDictionary

	// Fields returns a list of signature fields that refer the current object.
	// Port of getFields().
	Fields() []*PdfSignatureField

	// ModificationDetection returns an information about changes made in the document.
	// Port of getModificationDetection().
	ModificationDetection() *PdfModificationDetection
}
