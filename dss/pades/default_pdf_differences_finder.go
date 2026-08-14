// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/modifications/DefaultPdfDifferencesFinder.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf.modifications is part of the eu.europa.esig.dss.pdf module that landed
// in no s5b manifest (see pdf_object.go's header). slf4j is dropped, per PORTING.md.
//
// DEVIATION: VisualDifferences always returns nil in practice - the native PDF engine has no
// rasteriser (internal/pdf/DESIGN.md §0.2), so GenerateImageScreenshot always returns
// ErrRasterisationNotSupported, which this method's error handling (faithfully ported from
// upstream's IOException catch) turns into "no differences reported for this page", exactly as
// native_pdf_signature_service.go's VisualDifferences doc comment already documents.
package pades

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// DefaultPdfDifferencesFinder is the default implementation used to find differences in pages
// between two PDF revisions. Port of the DefaultPdfDifferencesFinder class.
type DefaultPdfDifferencesFinder struct {
	// maximalPagesAmountForVisualComparison sets the maximal amount of pages in a PDF to execute
	// visual screenshot comparison for. Default: 10 pages.
	maximalPagesAmountForVisualComparison int
}

// NewDefaultPdfDifferencesFinder instantiates an object with the default configuration.
// Port of the default constructor.
func NewDefaultPdfDifferencesFinder() *DefaultPdfDifferencesFinder {
	return &DefaultPdfDifferencesFinder{maximalPagesAmountForVisualComparison: 10}
}

// SetMaximalPagesAmountForVisualComparison sets the maximal pages amount in a PDF to process a
// visual screenshot comparison for. Set to 0 to disable the visual comparison check.
// Default: 10 pages. Port of #setMaximalPagesAmountForVisualComparison.
func (f *DefaultPdfDifferencesFinder) SetMaximalPagesAmountForVisualComparison(pagesAmount int) {
	f.maximalPagesAmountForVisualComparison = pagesAmount
}

// AnnotationOverlaps returns a list of found annotation overlaps. Port of #getAnnotationOverlaps.
func (f *DefaultPdfDifferencesFinder) AnnotationOverlaps(reader PdfDocumentReader) []PdfModification {
	var annotationOverlaps []PdfModification

	for pageNumber := 1; pageNumber <= reader.NumberOfPages(); pageNumber++ {
		pdfAnnotations := f.pdfAnnotations(reader, pageNumber)
		for i, annotation := range pdfAnnotations {
			// The remaining annotations (i.e. every entry but this one), matching Java's
			// "remove the annotation being checked from the comparison list" iterator dance.
			remaining := make([]*PdfAnnotation, 0, len(pdfAnnotations)-1)
			remaining = append(remaining, pdfAnnotations[:i]...)
			remaining = append(remaining, pdfAnnotations[i+1:]...)
			if f.IsAnnotationBoxOverlapping(annotation.AnnotationBox(), remaining) {
				annotationOverlaps = append(annotationOverlaps, NewCommonPdfModification(pageNumber))
				break
			}
		}
	}

	return annotationOverlaps
}

// pdfAnnotations ports the private getPdfAnnotations(PdfDocumentReader, int).
func (f *DefaultPdfDifferencesFinder) pdfAnnotations(reader PdfDocumentReader, pageNumber int) []*PdfAnnotation {
	annotations, err := reader.PdfAnnotations(pageNumber)
	if err != nil {
		// Upstream logs "Unable to extract annotations from a PDF document for a page number :
		// {}. Reason : {}".
		return nil
	}
	return annotations
}

// IsAnnotationBoxOverlapping checks if annotationBox overlaps with pdfAnnotations.
// Port of #isAnnotationBoxOverlapping.
func (f *DefaultPdfDifferencesFinder) IsAnnotationBoxOverlapping(annotationBox AnnotationBox, pdfAnnotations []*PdfAnnotation) bool {
	if annotationBox.Width() == 0 || annotationBox.Height() == 0 {
		// invisible field
		return false
	}
	for _, pdfAnnotation := range pdfAnnotations {
		if annotationBox.IsOverlapAnnotation(pdfAnnotation) {
			return true
		}
	}
	return false
}

// PagesDifferences returns a list of missing/added pages between the signed and final revisions.
// Port of #getPagesDifferences.
func (f *DefaultPdfDifferencesFinder) PagesDifferences(signedRevisionReader, finalRevisionReader PdfDocumentReader) []PdfModification {
	signedPages := signedRevisionReader.NumberOfPages()
	finalPages := finalRevisionReader.NumberOfPages()

	maxNumberOfPages := signedPages
	if finalPages > maxNumberOfPages {
		maxNumberOfPages = finalPages
	}
	minNumberOfPages := signedPages
	if finalPages < minNumberOfPages {
		minNumberOfPages = finalPages
	}

	var missingPages []PdfModification
	for ii := maxNumberOfPages; ii > minNumberOfPages; ii-- {
		missingPages = append(missingPages, NewCommonPdfModification(ii))
	}

	if len(missingPages) > 0 {
		// Upstream logs "The provided PDF file contains {} additional pages against the signed
		// revision!" at WARN.
	}

	return missingPages
}

// VisualDifferences returns a list of visual differences found between the signed and final
// revisions, excluding newly created annotations. Port of #getVisualDifferences.
func (f *DefaultPdfDifferencesFinder) VisualDifferences(signedRevisionReader, finalRevisionReader PdfDocumentReader) []PdfModification {
	pagesAmount := finalRevisionReader.NumberOfPages()
	if f.maximalPagesAmountForVisualComparison < pagesAmount {
		// Upstream logs "The provided document contains {} pages, while the limit for a visual
		// comparison is set to {}. Visual differences comparison is skipped." at debug level.
		return nil
	}

	var visualDifferences []PdfModification
	for pageNumber := 1; pageNumber <= signedRevisionReader.NumberOfPages() && pageNumber <= finalRevisionReader.NumberOfPages(); pageNumber++ {
		signedScreenshot, err := signedRevisionReader.GenerateImageScreenshot(pageNumber)
		if err != nil {
			// Upstream logs "Unable to get visual differences for a page number : {}. Reason :
			// {}" (with XML-invalid characters cleaned from the message).
			continue
		}

		signedAnnotations, err := signedRevisionReader.PdfAnnotations(pageNumber)
		if err != nil {
			continue
		}
		finalAnnotations, err := finalRevisionReader.PdfAnnotations(pageNumber)
		if err != nil {
			continue
		}

		addedAnnotations := f.updatedAnnotations(signedAnnotations, finalAnnotations)
		finalScreenshot, err := finalRevisionReader.GenerateImageScreenshotWithoutAnnotations(pageNumber, addedAnnotations)
		if err != nil {
			continue
		}

		equal, err := f.screenshotsEqual(signedScreenshot, finalScreenshot)
		if err != nil {
			continue
		}
		if !equal {
			// Upstream logs "A visual difference found on page {} between a signed revision and
			// the final document!" at WARN.
			visualDifferences = append(visualDifferences, NewCommonPdfModification(pageNumber))
		}
	}
	return visualDifferences
}

// updatedAnnotations ports the private getUpdatedAnnotations.
func (f *DefaultPdfDifferencesFinder) updatedAnnotations(signedAnnotations, finalAnnotations []*PdfAnnotation) []*PdfAnnotation {
	var updatedAnnotations []*PdfAnnotation
	for _, annotationBox := range finalAnnotations {
		found := false
		for _, signed := range signedAnnotations {
			if signed.Equals(annotationBox) {
				found = true
				break
			}
		}
		if !found {
			updatedAnnotations = append(updatedAnnotations, annotationBox)
		}
	}
	return updatedAnnotations
}

// screenshotsEqual compares two rendered page screenshots for equality. Java compares decoded
// BufferedImage pixel data (ImageUtils#imagesEqual, eu.europa.esig.dss.pdf.visible - part of the
// "not ported" visible SPI stack per internal/pdf/DESIGN.md §0.2); since the native engine never
// produces a screenshot (GenerateImageScreenshot/GenerateImageScreenshotWithoutAnnotations both
// return ErrRasterisationNotSupported), this code path is unreachable in practice, so a byte-
// level comparison of the two documents' encodings is a faithful-enough stand-in.
func (f *DefaultPdfDifferencesFinder) screenshotsEqual(signed, final model.DSSDocument) (bool, error) {
	signedStream, err := signed.OpenStream()
	if err != nil {
		return false, err
	}
	defer func() { _ = signedStream.Close() }()
	finalStream, err := final.OpenStream()
	if err != nil {
		return false, err
	}
	defer func() { _ = finalStream.Close() }()
	return utils.CompareInputStreams(signedStream, finalStream)
}

// Compile-time assertion standing in for Java's "implements PdfDifferencesFinder".
var _ PdfDifferencesFinder = (*DefaultPdfDifferencesFinder)(nil)
