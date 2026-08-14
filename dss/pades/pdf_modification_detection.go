// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/modifications/PdfModificationDetection.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf.modifications is part of the eu.europa.esig.dss.pdf module that landed
// in no s5b manifest (see pdf_object.go's header). Shape (SetAnnotationOverlaps/SetPageDifferences/
// SetVisualDifferences/SetObjectModifications, NewPdfModificationDetection()) confirmed against
// native_pdf_signature_service.go's already-landed call sites, and pdf_revision.go's
// ModificationDetection() *PdfModificationDetection getter.
package pades

// PdfModificationDetection contains necessary information about a PDF visual or structure
// modification. Port of the PdfModificationDetection class.
type PdfModificationDetection struct {
	// annotationOverlaps is the list of annotation overlaps.
	annotationOverlaps []PdfModification

	// visualDifferences is the list of visual differences between revisions.
	visualDifferences []PdfModification

	// pageDifferences is the list of page amount differences between revisions.
	pageDifferences []PdfModification

	// objectModifications is the filtered collection of ObjectModifications between the signed
	// and final revisions.
	objectModifications PdfObjectModifications
}

// NewPdfModificationDetection instantiates an object with null-equivalent (zero) values.
// Port of the default constructor.
func NewPdfModificationDetection() *PdfModificationDetection {
	return &PdfModificationDetection{}
}

// AnnotationOverlaps returns information about annotations overlapping. Port of #getAnnotationOverlaps.
func (d *PdfModificationDetection) AnnotationOverlaps() []PdfModification {
	return d.annotationOverlaps
}

// SetAnnotationOverlaps sets the annotation overlaps. Port of #setAnnotationOverlaps.
func (d *PdfModificationDetection) SetAnnotationOverlaps(annotationOverlaps []PdfModification) {
	d.annotationOverlaps = annotationOverlaps
}

// PageDifferences returns information about missing/added pages between the signed and final
// revisions. Port of #getPageDifferences.
func (d *PdfModificationDetection) PageDifferences() []PdfModification { return d.pageDifferences }

// SetPageDifferences sets the page differences (for missing/added pages). Port of #setPageDifferences.
func (d *PdfModificationDetection) SetPageDifferences(pageDifferences []PdfModification) {
	d.pageDifferences = pageDifferences
}

// VisualDifferences returns information about pages with visual differences between signed and
// final revisions. Port of #getVisualDifferences.
func (d *PdfModificationDetection) VisualDifferences() []PdfModification { return d.visualDifferences }

// SetVisualDifferences sets the visual differences. Port of #setVisualDifferences.
func (d *PdfModificationDetection) SetVisualDifferences(visualDifferences []PdfModification) {
	d.visualDifferences = visualDifferences
}

// ObjectModifications returns a filtered collection of modified objects between the signed and
// final document revisions. Port of #getObjectModifications.
func (d *PdfModificationDetection) ObjectModifications() PdfObjectModifications {
	return d.objectModifications
}

// SetObjectModifications sets the collection of filtered object modifications.
// Port of #setObjectModifications.
func (d *PdfModificationDetection) SetObjectModifications(objectModifications PdfObjectModifications) {
	d.objectModifications = objectModifications
}

// AreModificationsDetected returns whether any modifications have been detected.
// Port of #areModificationsDetected.
func (d *PdfModificationDetection) AreModificationsDetected() bool {
	return len(d.annotationOverlaps) > 0 || len(d.visualDifferences) > 0 || len(d.pageDifferences) > 0 ||
		!d.objectModifications.IsEmpty()
}
