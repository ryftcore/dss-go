// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/signature/PAdESExtensionService.java (DSS 6.5.RC1).
//
// Java's three incorporateValidationData overloads become three Go names, since Go has no
// overloading:
//
//	incorporateValidationData(DSSDocument)                       -> IncorporateValidationData
//	incorporateValidationData(DSSDocument, char[])               -> IncorporateValidationDataWithPassword
//	incorporateValidationData(DSSDocument, char[], boolean)      -> IncorporateValidationDataWithVRI
//
// char[] is []byte, the convention token/password_protection.go set. slf4j is dropped.
package pades

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/executor"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ExtensionService obtains the validation data for the signatures/timestamps within a PDF
// file and incorporates it on the LT-level of the document, i.e. creates a DSS dictionary
// revision.
type ExtensionService struct {
	// certificateVerifier processes the validation data.
	certificateVerifier validation.CertificateVerifier

	// pdfObjectFactory is used for the PDF document processing.
	pdfObjectFactory IPdfObjFactory
}

// NewPAdESExtensionService instantiates the service with the default IPdfObjFactory.
// Port of PAdESExtensionService(CertificateVerifier).
func NewPAdESExtensionService(certificateVerifier validation.CertificateVerifier) *ExtensionService {
	return NewPAdESExtensionServiceWithFactory(certificateVerifier, NewDefaultPdfObjFactory())
}

// NewPAdESExtensionServiceWithFactory is the default constructor.
// Port of PAdESExtensionService(CertificateVerifier, IPdfObjFactory).
func NewPAdESExtensionServiceWithFactory(certificateVerifier validation.CertificateVerifier,
	pdfObjectFactory IPdfObjFactory) *ExtensionService {
	if certificateVerifier == nil {
		panic("CertificateVerifier cannot be null!")
	}
	if pdfObjectFactory == nil {
		panic("PdfObjectFactory cannot be null!")
	}
	return &ExtensionService{certificateVerifier: certificateVerifier, pdfObjectFactory: pdfObjectFactory}
}

// IncorporateValidationData adds a DSS dictionary revision to the given document without
// password-protection, with the required validation data if needed and no VRI dictionary
// created.
//
// NOTE: this method does not check the validity of the provided signatures/timestamps (e.g. a
// T-level, ...). Port of incorporateValidationData(DSSDocument).
func (s *ExtensionService) IncorporateValidationData(document model.DSSDocument) model.DSSDocument {
	return s.IncorporateValidationDataWithPassword(document, nil)
}

// IncorporateValidationDataWithPassword adds a DSS dictionary revision to the given document
// protected by a passwordProtection, with the required validation data if needed and without a
// VRI dictionary. Port of incorporateValidationData(DSSDocument, char[]).
func (s *ExtensionService) IncorporateValidationDataWithPassword(document model.DSSDocument,
	passwordProtection []byte) model.DSSDocument {
	return s.IncorporateValidationDataWithVRI(document, passwordProtection, false)
}

// IncorporateValidationDataWithVRI adds a DSS dictionary revision to the given document
// protected by a passwordProtection, with the required validation data if needed and a VRI
// dictionary when includeVRIDict is set.
// Port of incorporateValidationData(DSSDocument, char[], boolean).
func (s *ExtensionService) IncorporateValidationDataWithVRI(document model.DSSDocument,
	passwordProtection []byte, includeVRIDict bool) model.DSSDocument {
	if document == nil {
		panic("The document to be extended shall be provided!")
	}

	pdfDocumentAnalyzer := s.pdfDocumentAnalyzer(document, passwordProtection)

	signatures := pdfDocumentAnalyzer.Signatures()
	detachedTimestamps := pdfDocumentAnalyzer.DetachedTimestamps()
	if utils.IsCollectionNotEmpty(signatures) {
		signatureTimestamps := padesExtensionServiceSignatureTimestamps(signatures)
		if utils.IsCollectionEmpty(signatureTimestamps) {
			// Upstream logs "The found signatures within the document with name '{}' do not have
			// a T-level. Validation data incorporation skipped."
			return document
		}

	} else if utils.IsCollectionNotEmpty(detachedTimestamps) {
		// continue

	} else {
		// Upstream logs "No signatures or timestamps found within a document with name '{}'."
		return document
	}

	validationData := pdfDocumentAnalyzer.ValidationData(signatures, detachedTimestamps)
	if validationData.IsEmpty() {
		// Upstream logs "No validation data has been obtained for the document with name '{}'.
		// Return original document."
		return document
	}

	signatureService := s.newPdfSignatureService()
	return signatureService.AddDssDictionary(document, validationData, passwordProtection, includeVRIDict)
}

// padesExtensionServiceSignatureTimestamps ports the private #getSignatureTimestamps.
func padesExtensionServiceSignatureTimestamps(signatures []validation.AdvancedSignature) []*validation.TimestampToken {
	var signatureTimestamps []*validation.TimestampToken
	for _, signature := range signatures {
		signatureTimestamps = append(signatureTimestamps, signature.AllTimestamps()...)
	}
	return signatureTimestamps
}

// pdfDocumentAnalyzer ports the private #getPDFDocumentAnalyzer.
func (s *ExtensionService) pdfDocumentAnalyzer(document model.DSSDocument,
	passwordProtection []byte) *PDFDocumentAnalyzer {
	pdfDocumentAnalyzer := NewPDFDocumentAnalyzer(document)
	pdfDocumentAnalyzer.SetCertificateVerifier(s.certificateVerifier)
	pdfDocumentAnalyzer.SetValidationContextExecutor(executor.CompleteValidationContextExecutorInstance)
	pdfDocumentAnalyzer.SetPasswordProtection(passwordProtection)
	pdfDocumentAnalyzer.SetPdfObjFactory(s.pdfObjectFactory)
	return pdfDocumentAnalyzer
}

// newPdfSignatureService ports the private #newPdfSignatureService.
func (s *ExtensionService) newPdfSignatureService() PDFSignatureService {
	return s.pdfObjectFactory.NewPAdESSignatureService()
}
