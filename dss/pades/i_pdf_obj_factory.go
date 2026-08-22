// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/IPdfObjFactory.java (DSS 6.5.RC1).
//
// Upstream ships two interchangeable SPI backends (dss-pades-pdfbox, dss-pades-openpdf) that
// ServiceLoaderPdfObjFactory discovers at runtime; the Go port has exactly one native backend, so
// ServiceLoaderPdfObjFactory / AbstractPdfObjFactory collapse into NewDefaultPdfObjFactory
// (internal/pdf/DESIGN.md §0.2). The interface itself is kept: it is the documented extension
// point PAdESService#setPdfObjFactory exposes.
package pades

import "github.com/ryftcore/dss-go/dss/spi/signature/resources"

// IPdfObjFactory loads the relevant implementations of PDFSignatureService.
type IPdfObjFactory interface {
	// NewContentTimestampService returns the service used for a content timestamp creation.
	// Port of #newContentTimestampService.
	NewContentTimestampService() PDFSignatureService

	// NewPAdESSignatureService returns the service used for a signature creation.
	// Port of #newPAdESSignatureService.
	NewPAdESSignatureService() PDFSignatureService

	// NewSignatureTimestampService returns the service used for a signature timestamp creation.
	// Port of #newSignatureTimestampService.
	NewSignatureTimestampService() PDFSignatureService

	// NewArchiveTimestampService returns the service used for an archive timestamp creation.
	// Port of #newArchiveTimestampService.
	NewArchiveTimestampService() PDFSignatureService

	// SetResourcesHandlerBuilder sets a DSSResourcesHandlerBuilder to be used for operating with
	// internal objects during the signature creation procedure.
	// Port of #setResourcesHandlerBuilder.
	SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder)

	// SetPdfDifferencesFinder sets a custom PdfDifferencesFinder to detect differences between
	// signed and final PDF document revisions. Port of #setPdfDifferencesFinder.
	SetPdfDifferencesFinder(pdfDifferencesFinder PdfDifferencesFinder)

	// SetPdfObjectModificationsFinder sets a custom PdfObjectModificationsFinder to detect
	// modifications occurred within internal PDF objects between signed and final PDF document
	// revisions. Port of #setPdfObjectModificationsFinder.
	SetPdfObjectModificationsFinder(pdfObjectModificationsFinder PdfObjectModificationsFinder)

	// SetPdfPermissionsChecker sets a custom PdfPermissionsChecker to verify the PDF document
	// encryption dictionary permission rules for a new signature creation, when applicable.
	// Port of #setPdfPermissionsChecker.
	SetPdfPermissionsChecker(pdfPermissionsChecker *PdfPermissionsChecker)

	// SetPdfSignatureFieldPositionChecker sets a custom PdfSignatureFieldPositionChecker to
	// verify the validity of a new signature field placement.
	// Port of #setPdfSignatureFieldPositionChecker.
	SetPdfSignatureFieldPositionChecker(pdfSignatureFieldPositionChecker *PdfSignatureFieldPositionChecker)

	// SetPdfMemoryUsageSetting sets a custom PdfMemoryUsageSetting specifying the load mode of
	// the PDF document. Port of #setPdfMemoryUsageSetting.
	SetPdfMemoryUsageSetting(pdfMemoryUsageSetting PdfMemoryUsageSetting)
}
