// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/scope/PdfRevisionScopeFinder.java
// (DSS 6.5.RC1).
//
// FORWARD DEPENDENCIES:
//
//   - PdfCMSRevision (eu.europa.esig.dss.pdf.PdfCMSRevision) is already used as a forward-
//     referenced Go interface type elsewhere in this package (native_pdf_signature_service.go's
//     AnalyzeRevisionModifications takes a `pdfRevision PdfCMSRevision` parameter and calls
//     pdfRevision.ByteRange()/SetModificationDetection(...)). This file additionally needs
//     AreAllOriginalBytesCovered() bool (Java's areAllOriginalBytesCovered()), so the type's
//     complete assumed shape, gathering every call site across the package, is:
//
//     type PdfCMSRevision interface {
//     PdfRevision
//     ByteRange() *ByteRange
//     AreAllOriginalBytesCovered() bool
//     SetModificationDetection(*PdfModificationDetection)
//     }
//
//   - PAdESUtilsGetOriginalPDFFromRevision(pdfRevision PdfCMSRevision) model.DSSDocument is the
//     flattened static PAdESUtils.getOriginalPDF(PdfCMSRevision) overload; named distinctly from
//     the already-assumed PAdESUtilsGetOriginalPDF(padesSignature *PAdESSignature) model.DSSDocument
//     (pdf_document_analyzer.go's header), which flattens the sibling
//     PAdESUtils.getOriginalPDF(PAdESSignature) overload - Go has no overloading, and upstream's
//     PAdESSignature overload only ever delegates to this one (getOriginalPDF(padesSignature.getPdfRevision())).
package pades

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/scope"
	spiscope "github.com/utain/esig/dss/spi/validation/scope"
)

// PdfRevisionScopeFinder is an abstract type to find a PdfRevision scope. Port of the abstract
// class PdfRevisionScopeFinder, extending spiscope.AbstractSignatureScopeFinder.
type PdfRevisionScopeFinder struct {
	spiscope.AbstractSignatureScopeFinder
}

// newPdfRevisionScopeFinder is the port of the protected default constructor.
func newPdfRevisionScopeFinder() PdfRevisionScopeFinder {
	return PdfRevisionScopeFinder{AbstractSignatureScopeFinder: spiscope.NewAbstractSignatureScopeFinder()}
}

// findSignatureScope finds signature scopes from a PdfCMSRevision. Port of the protected
// findSignatureScope(PdfCMSRevision).
func (f *PdfRevisionScopeFinder) findSignatureScope(pdfRevision PdfCMSRevision) scope.SignatureScope {
	if pdfRevision.AreAllOriginalBytesCovered() {
		return NewFullPdfByteRangeSignatureScope(pdfRevision.ByteRange(), f.getOriginalPdfRevision(pdfRevision))
	}
	return NewPartialPdfByteRangeSignatureScope(pdfRevision.ByteRange(), f.getOriginalPdfRevision(pdfRevision))
}

// getOriginalPdfRevision ports the private getOriginalPdfRevision(PdfCMSRevision).
func (f *PdfRevisionScopeFinder) getOriginalPdfRevision(pdfRevision PdfCMSRevision) model.DSSDocument {
	return PAdESUtilsGetOriginalPDFFromRevision(pdfRevision)
}
