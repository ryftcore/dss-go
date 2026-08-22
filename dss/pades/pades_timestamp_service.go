// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/timestamp/PAdESTimestampService.java (DSS 6.5.RC1).
//
// Java's timestamp package is flattened into the single pades package.
// Java's two constructors become two constructor funcs, since Go has no overloading:
//
//	TimestampService(TSPSource)                       -> NewTimestampService
//	TimestampService(TSPSource, PDFSignatureService)  -> NewTimestampServiceWithPDFService
//
// ServiceLoaderPdfObjFactory has no Go counterpart - there is exactly one native backend, so the
// factory collapses to NewDefaultPdfObjFactory (internal/pdf/DESIGN.md §0.2).
package pades

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// TimestampService timestamps a PDF, i.e. creates a document time-stamp revision.
type TimestampService struct {
	// tspSource obtains the timestamp.
	tspSource validation.TSPSource

	// pdfSignatureService is the signature service implementation to use.
	pdfSignatureService PDFSignatureService
}

// NewTimestampService instantiates the service with a default PDFSignatureService for an
// archive (document) timestamp creation. Port of PAdESTimestampService(TSPSource).
func NewTimestampService(tspSource validation.TSPSource) *TimestampService {
	return NewTimestampServiceWithPDFService(tspSource, NewDefaultPdfObjFactory().NewArchiveTimestampService())
}

// NewTimestampServiceWithPDFService is the default constructor.
// Port of PAdESTimestampService(TSPSource, PDFSignatureService).
func NewTimestampServiceWithPDFService(tspSource validation.TSPSource,
	pdfSignatureService PDFSignatureService) *TimestampService {
	if tspSource == nil {
		panic("TSPSource shall be provided!")
	}
	if pdfSignatureService == nil {
		panic("PDFSignatureService shall be provided!")
	}
	return &TimestampService{tspSource: tspSource, pdfSignatureService: pdfSignatureService}
}

// TimestampDocument timestamps the document. Port of #timestampDocument.
func (s *TimestampService) TimestampDocument(document model.DSSDocument,
	params *TimestampParameters) model.DSSDocument {
	if document == nil {
		panic("DSSDocument shall be provided!")
	}
	if params == nil {
		panic("PAdESTimestampParameters cannot be null!")
	}
	UtilsAssertPdfDocument(document)

	messageDigest := s.pdfSignatureService.MessageDigest(document, params)
	timeStampToken, err := s.tspSource.TimeStampResponse(messageDigest.Algorithm(), messageDigest.Value())
	if err != nil {
		panic(err)
	}
	encoded, err := spi.DSSASN1UtilsDEREncodedTimestampBinary(timeStampToken)
	if err != nil {
		panic(err)
	}
	return s.pdfSignatureService.Sign(document, encoded, params)
}
