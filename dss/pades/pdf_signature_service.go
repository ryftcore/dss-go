// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PDFSignatureService.java (DSS 6.5.RC1).
//
// The eu.europa.esig.dss.pdf SPI is flattened into the pades package together with the rest of
// dss-pades. Upstream ships two interchangeable backends (pdfbox, openpdf)
// selected by a ServiceLoader; the Go port has exactly one, the native internal/pdf engine, so
// this interface stays as the seam a caller can substitute while IPdfObjFactory collapses to a
// constructor (internal/pdf/DESIGN.md §0.2).
//
// Java's overloads become distinct Go names:
//
//	addDssDictionary(doc, data)                          -> AddDssDictionaryDefault
//	addDssDictionary(doc, data, pwd)                     -> AddDssDictionaryWithPassword
//	addDssDictionary(doc, data, pwd, includeVRIDict)     -> AddDssDictionary
//	getAvailableSignatureFields(doc)                     -> GetAvailableSignatureFieldsDefault
//	getAvailableSignatureFields(doc, pwd)                -> GetAvailableSignatureFields
//	addNewSignatureField(doc, params)                    -> AddNewSignatureFieldDefault
//	addNewSignatureField(doc, params, pwd)               -> AddNewSignatureField
//
// char[] is []byte (token/password_protection.go's convention); the Java methods declare no
// checked exception, so every failure is a panic carrying the error, as elsewhere in the
// signature services.
package pades

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/signature/resources"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// PDFSignatureService lets the user choose the underlying PDF implementation used to create PDF
// signatures.
type PDFSignatureService interface {
	// MessageDigest returns the message-digest computed on the PDF signature revision's
	// ByteRange. Port of #messageDigest.
	MessageDigest(toSignDocument model.DSSDocument, parameters CommonParameters) model.DSSMessageDigest

	// Sign signs a PDF document, enveloping the encoded CMS signed data in a new revision.
	// Port of #sign.
	Sign(toSignDocument model.DSSDocument, cmsSignedData []byte,
		parameters CommonParameters) model.DSSDocument

	// GetRevisions retrieves the revisions from a PDF document. pwd is nil for a document that
	// is not encrypted. Port of #getRevisions.
	GetRevisions(document model.DSSDocument, pwd []byte) []PdfRevision

	// AddDssDictionaryDefault adds the DSS dictionary (Baseline-LT) to a document without
	// password-protection and without a VRI dictionary.
	// Port of addDssDictionary(DSSDocument, PdfValidationDataContainer).
	AddDssDictionaryDefault(document model.DSSDocument,
		validationDataForInclusion *PdfValidationDataContainer) model.DSSDocument

	// AddDssDictionaryWithPassword adds the DSS dictionary (Baseline-LT) to a
	// password-protected document without a VRI dictionary.
	// Port of addDssDictionary(DSSDocument, PdfValidationDataContainer, char[]).
	AddDssDictionaryWithPassword(document model.DSSDocument,
		validationDataForInclusion *PdfValidationDataContainer, pwd []byte) model.DSSDocument

	// AddDssDictionary adds the DSS dictionary (Baseline-LT) to a password-protected document,
	// with a VRI dictionary when includeVRIDict is set.
	// Port of addDssDictionary(DSSDocument, PdfValidationDataContainer, char[], boolean).
	AddDssDictionary(document model.DSSDocument, validationDataForInclusion *PdfValidationDataContainer,
		pwd []byte, includeVRIDict bool) model.DSSDocument

	// GetAvailableSignatureFieldsDefault returns the not signed signature fields.
	// Port of getAvailableSignatureFields(DSSDocument).
	GetAvailableSignatureFieldsDefault(document model.DSSDocument) []string

	// GetAvailableSignatureFields returns the not signed signature fields of an encrypted
	// document. Port of getAvailableSignatureFields(DSSDocument, char[]).
	GetAvailableSignatureFields(document model.DSSDocument, pwd []byte) []string

	// AddNewSignatureFieldDefault adds a new signature field to an existing PDF document.
	// Port of addNewSignatureField(DSSDocument, SignatureFieldParameters).
	AddNewSignatureFieldDefault(document model.DSSDocument,
		parameters *SignatureFieldParameters) model.DSSDocument

	// AddNewSignatureField adds a new signature field to an existing encrypted PDF document.
	// Port of addNewSignatureField(DSSDocument, SignatureFieldParameters, char[]).
	AddNewSignatureField(document model.DSSDocument, parameters *SignatureFieldParameters,
		pwd []byte) model.DSSDocument

	// AnalyzePdfModifications analyzes the PDF revisions and tries to detect any modification
	// (shadow attacks) for the given signatures. Port of #analyzePdfModifications.
	AnalyzePdfModifications(document model.DSSDocument, signatures []validation.AdvancedSignature, pwd []byte)

	// AnalyzeTimestampPdfModifications analyzes the PDF revisions and tries to detect any
	// modification (shadow attacks) for the given document timestamps.
	// Port of #analyzeTimestampPdfModifications.
	AnalyzeTimestampPdfModifications(document model.DSSDocument, timestamps []*validation.TimestampToken, pwd []byte)

	// PreviewPageWithVisualSignature returns a page preview with the visual signature, as a
	// document containing a PNG picture. Port of #previewPageWithVisualSignature.
	PreviewPageWithVisualSignature(toSignDocument model.DSSDocument,
		parameters CommonParameters) model.DSSDocument

	// PreviewSignatureField returns a preview of the signature field, as a document containing a
	// PNG picture. Port of #previewSignatureField.
	PreviewSignatureField(toSignDocument model.DSSDocument,
		parameters CommonParameters) model.DSSDocument

	// SetResourcesHandlerBuilder sets the DSSResourcesHandlerBuilder used to create a
	// DSSResourcesHandler in internal methods, which defines how OutputStreams are operated and
	// DSSDocuments created.
	//
	// Default: document.InMemoryResourcesHandlerBuilder, working with data in memory.
	// Port of #setResourcesHandlerBuilder.
	SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder)

	// SetPdfDifferencesFinder sets the PdfDifferencesFinder used to find the differences on
	// pages between given PDF revisions.
	//
	// Default: DefaultPdfDifferencesFinder. Port of #setPdfDifferencesFinder.
	SetPdfDifferencesFinder(pdfDifferencesFinder PdfDifferencesFinder)

	// SetPdfObjectModificationsFinder sets the PdfObjectModificationsFinder used to find the
	// differences between internal PDF objects occurred between given PDF revisions.
	//
	// Default: DefaultPdfObjectModificationsFinder. Port of #setPdfObjectModificationsFinder.
	SetPdfObjectModificationsFinder(pdfObjectModificationsFinder PdfObjectModificationsFinder)

	// SetPdfPermissionsChecker sets the PdfPermissionsChecker used to verify the PDF document
	// rules for a new signature creation. Port of #setPdfPermissionsChecker.
	SetPdfPermissionsChecker(pdfPermissionsChecker *PdfPermissionsChecker)

	// SetPdfSignatureFieldPositionChecker sets the PdfSignatureFieldPositionChecker used to
	// verify the validity of a new signature field placement, e.g. that it lies within the PDF
	// page borders and/or does not overlap with existing signature fields.
	// Port of #setPdfSignatureFieldPositionChecker.
	SetPdfSignatureFieldPositionChecker(pdfSignatureFieldPositionChecker *PdfSignatureFieldPositionChecker)

	// SetPdfMemoryUsageSetting sets the PdfMemoryUsageSetting specifying the load mode of the
	// PDF document. Port of #setPdfMemoryUsageSetting.
	SetPdfMemoryUsageSetting(pdfMemoryUsageSetting PdfMemoryUsageSetting)
}
