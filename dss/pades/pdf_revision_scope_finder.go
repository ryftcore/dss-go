// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/scope/PdfRevisionScopeFinder.java
// (DSS 6.5.RC1).
package pades

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	spiscope "github.com/ryftcore/dss-go/dss/spi/validation/scope"
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
