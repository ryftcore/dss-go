// Native implementation of eu.europa.esig.dss.pdf.IPdfObjFactory on internal/pdf.
//
// It merges the behaviour of two upstream classes that have no separate reason to exist once
// there is a single backend: eu.europa.esig.dss.pdf.AbstractPdfObjFactory (the configuration
// holder and its configure() propagation) and
// eu.europa.esig.dss.pdf.pdfbox.PdfBoxDefaultObjectFactory (the four newXxxService factories,
// each one instantiating the signature service in a different PDFServiceMode).
// eu.europa.esig.dss.pdf.ServiceLoaderPdfObjFactory has no counterpart: the Go port has exactly
// one backend, so the ServiceLoader lookup collapses into NewDefaultPdfObjFactory
// (internal/pdf/DESIGN.md §0.2).
package pades

import "github.com/ryftcore/dss-go/dss/spi/signature/resources"

// NativePdfObjFactory produces PDFSignatureService instances backed by the native internal/pdf
// engine. A nil configuration field means "leave the service's own default", which is what
// AbstractPdfObjFactory#configure expresses with its null checks.
type NativePdfObjFactory struct {
	resourcesHandlerBuilder          resources.DSSResourcesHandlerBuilder
	pdfDifferencesFinder             PdfDifferencesFinder
	pdfObjectModificationsFinder     PdfObjectModificationsFinder
	pdfPermissionsChecker            *PdfPermissionsChecker
	pdfSignatureFieldPositionChecker *PdfSignatureFieldPositionChecker
	pdfMemoryUsageSetting            *PdfMemoryUsageSetting
}

// NewDefaultPdfObjFactory returns the default IPdfObjFactory of the Go port. It is what
// upstream's "new ServiceLoaderPdfObjFactory()" resolves to.
func NewDefaultPdfObjFactory() *NativePdfObjFactory {
	return &NativePdfObjFactory{}
}

// NewPAdESSignatureService returns the service used for a signature creation.
// Port of PdfBoxDefaultObjectFactory#newPAdESSignatureService.
func (f *NativePdfObjFactory) NewPAdESSignatureService() PDFSignatureService {
	return f.configure(NewNativePDFSignatureService(PDFServiceMode_SIGNATURE, f.signatureDrawerFactory()))
}

// NewContentTimestampService returns the service used for a content timestamp creation.
// Port of PdfBoxDefaultObjectFactory#newContentTimestampService.
func (f *NativePdfObjFactory) NewContentTimestampService() PDFSignatureService {
	return f.configure(NewNativePDFSignatureService(PDFServiceMode_CONTENT_TIMESTAMP, f.signatureDrawerFactory()))
}

// NewSignatureTimestampService returns the service used for a signature timestamp creation.
// Port of PdfBoxDefaultObjectFactory#newSignatureTimestampService.
func (f *NativePdfObjFactory) NewSignatureTimestampService() PDFSignatureService {
	return f.configure(NewNativePDFSignatureService(PDFServiceMode_SIGNATURE_TIMESTAMP, f.signatureDrawerFactory()))
}

// NewArchiveTimestampService returns the service used for an archive timestamp creation.
// Port of PdfBoxDefaultObjectFactory#newArchiveTimestampService.
func (f *NativePdfObjFactory) NewArchiveTimestampService() PDFSignatureService {
	return f.configure(NewNativePDFSignatureService(PDFServiceMode_ARCHIVE_TIMESTAMP, f.signatureDrawerFactory()))
}

// signatureDrawerFactory ports PdfBoxDefaultObjectFactory#getSignatureDrawerFactory.
func (f *NativePdfObjFactory) signatureDrawerFactory() SignatureDrawerFactory {
	return NewNativeSignatureDrawerFactory()
}

// SetResourcesHandlerBuilder ports AbstractPdfObjFactory#setResourcesHandlerBuilder.
func (f *NativePdfObjFactory) SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) {
	f.resourcesHandlerBuilder = resourcesHandlerBuilder
}

// SetPdfDifferencesFinder ports AbstractPdfObjFactory#setPdfDifferencesFinder.
func (f *NativePdfObjFactory) SetPdfDifferencesFinder(pdfDifferencesFinder PdfDifferencesFinder) {
	f.pdfDifferencesFinder = pdfDifferencesFinder
}

// SetPdfObjectModificationsFinder ports AbstractPdfObjFactory#setPdfObjectModificationsFinder.
func (f *NativePdfObjFactory) SetPdfObjectModificationsFinder(pdfObjectModificationsFinder PdfObjectModificationsFinder) {
	f.pdfObjectModificationsFinder = pdfObjectModificationsFinder
}

// SetPdfPermissionsChecker ports AbstractPdfObjFactory#setPdfPermissionsChecker.
func (f *NativePdfObjFactory) SetPdfPermissionsChecker(pdfPermissionsChecker *PdfPermissionsChecker) {
	f.pdfPermissionsChecker = pdfPermissionsChecker
}

// SetPdfSignatureFieldPositionChecker ports
// AbstractPdfObjFactory#setPdfSignatureFieldPositionChecker.
func (f *NativePdfObjFactory) SetPdfSignatureFieldPositionChecker(pdfSignatureFieldPositionChecker *PdfSignatureFieldPositionChecker) {
	f.pdfSignatureFieldPositionChecker = pdfSignatureFieldPositionChecker
}

// SetPdfMemoryUsageSetting ports AbstractPdfObjFactory#setPdfMemoryUsageSetting. The setting is
// a no-op knob kept for API compatibility: the native engine always works from an io.ReaderAt
// (internal/pdf/DESIGN.md §0.2).
func (f *NativePdfObjFactory) SetPdfMemoryUsageSetting(pdfMemoryUsageSetting PdfMemoryUsageSetting) {
	f.pdfMemoryUsageSetting = &pdfMemoryUsageSetting
}

// configure propagates the factory's configuration onto a newly created service.
// Port of the protected AbstractPdfObjFactory#configure.
func (f *NativePdfObjFactory) configure(pdfSignatureService PDFSignatureService) PDFSignatureService {
	if f.resourcesHandlerBuilder != nil {
		pdfSignatureService.SetResourcesHandlerBuilder(f.resourcesHandlerBuilder)
	}
	if f.pdfDifferencesFinder != nil {
		pdfSignatureService.SetPdfDifferencesFinder(f.pdfDifferencesFinder)
	}
	if f.pdfObjectModificationsFinder != nil {
		pdfSignatureService.SetPdfObjectModificationsFinder(f.pdfObjectModificationsFinder)
	}
	if f.pdfPermissionsChecker != nil {
		pdfSignatureService.SetPdfPermissionsChecker(f.pdfPermissionsChecker)
	}
	if f.pdfSignatureFieldPositionChecker != nil {
		pdfSignatureService.SetPdfSignatureFieldPositionChecker(f.pdfSignatureFieldPositionChecker)
	}
	if f.pdfMemoryUsageSetting != nil {
		pdfSignatureService.SetPdfMemoryUsageSetting(*f.pdfMemoryUsageSetting)
	}
	return pdfSignatureService
}

// Compile-time assertion standing in for Java's "extends AbstractPdfObjFactory".
var _ IPdfObjFactory = (*NativePdfObjFactory)(nil)
