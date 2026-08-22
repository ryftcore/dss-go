// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/signature/PAdESLevelBaselineT.java (DSS 6.5.RC1).
//
// Java's class is package-private and implements SignatureExtension<PAdESSignatureParameters>.
// The Go type is exported (Go has no package-private-yet-subclassable visibility, and the whole
// PAdES port lives in one package anyway).
//
// Java's two extendSignatures overloads cannot share one Go name:
//
//	extendSignatures(DSSDocument, PAdESSignatureParameters)                        -> ExtendSignatures
//	extendSignatures(DSSDocument, PDFDocumentAnalyzer, PAdESSignatureParameters)   -> ExtendSignaturesWithAnalyzer
//
// The second is overridden by PAdESLevelBaselineLT and called from the first, so the base
// dispatches through PAdESLevelBaselineTOverrides, registered by InitPAdESLevelBaselineT - the
// convention cades/cades_signature_extension.go established.
//
// slf4j is dropped (PORTING.md).
package pades

import (
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/executor"
	"github.com/ryftcore/dss-go/dss/utils"
)

// PAdESLevelBaselineTOverrides declares the operation PAdESLevelBaselineT calls back into and
// PAdESLevelBaselineLT overrides.
type PAdESLevelBaselineTOverrides interface {
	// ExtendSignaturesWithAnalyzer performs the document extension for the signatures the given
	// analyzer found. Port of the protected
	// extendSignatures(DSSDocument, PDFDocumentAnalyzer, PAdESSignatureParameters).
	ExtendSignaturesWithAnalyzer(signedDocument model.DSSDocument, pdfDocumentAnalyzer *PDFDocumentAnalyzer,
		parameters *PAdESSignatureParameters) model.DSSDocument
}

// PAdESLevelBaselineT holds the PAdES Baseline T signature profile.
type PAdESLevelBaselineT struct {
	// tspSource obtains a timestamp.
	tspSource validation.TSPSource

	// CertificateVerifier is the used CertificateVerifier.
	CertificateVerifier validation.CertificateVerifier

	// PdfObjectFactory is the used implementation for processing of a PDF document.
	PdfObjectFactory IPdfObjFactory

	// overrides points back at the concrete extension; see InitPAdESLevelBaselineT.
	overrides PAdESLevelBaselineTOverrides
}

// NewPAdESLevelBaselineT is the default constructor.
// Port of PAdESLevelBaselineT(TSPSource, CertificateVerifier, IPdfObjFactory).
func NewPAdESLevelBaselineT(tspSource validation.TSPSource, certificateVerifier validation.CertificateVerifier,
	pdfObjectFactory IPdfObjFactory) *PAdESLevelBaselineT {
	extension := &PAdESLevelBaselineT{}
	extension.InitPAdESLevelBaselineT(extension, tspSource, certificateVerifier, pdfObjectFactory)
	return extension
}

// InitPAdESLevelBaselineT registers the concrete extension with its base and applies the
// constructor's checks. Port of the protected
// PAdESLevelBaselineT(TSPSource, CertificateVerifier, IPdfObjFactory) constructor.
func (t *PAdESLevelBaselineT) InitPAdESLevelBaselineT(self PAdESLevelBaselineTOverrides,
	tspSource validation.TSPSource, certificateVerifier validation.CertificateVerifier,
	pdfObjectFactory IPdfObjFactory) {
	if tspSource == nil {
		panic("The TSPSource cannot be null")
	}
	if certificateVerifier == nil {
		panic("CertificateVerifier shall be defined!")
	}
	if pdfObjectFactory == nil {
		panic("pdfObjectFactory shall be defined!")
	}
	t.overrides = self
	t.tspSource = tspSource
	t.CertificateVerifier = certificateVerifier
	t.PdfObjectFactory = pdfObjectFactory
}

// ExtendSignatures adds a DocumentTimeStamp to the document; a signature-timestamp (CMS) is
// impossible to add while extending. Port of
// extendSignatures(DSSDocument, PAdESSignatureParameters).
func (t *PAdESLevelBaselineT) ExtendSignatures(signedDocument model.DSSDocument,
	params *PAdESSignatureParameters) model.DSSDocument {
	if signedDocument == nil {
		panic("DSSDocument cannot be null!")
	}
	if params == nil {
		panic("SignatureParameters cannot be null!")
	}
	// Will add a DocumentTimeStamp. signature-timestamp (CMS) is impossible to add while extending
	pdfDocumentValidator := t.PDFDocumentValidator(signedDocument, params)
	return t.overrides.ExtendSignaturesWithAnalyzer(signedDocument, pdfDocumentValidator, params)
}

// ExtendSignaturesWithAnalyzer performs a document extension. Port of the protected
// extendSignatures(DSSDocument, PDFDocumentAnalyzer, PAdESSignatureParameters).
func (t *PAdESLevelBaselineT) ExtendSignaturesWithAnalyzer(signedDocument model.DSSDocument,
	pdfDocumentAnalyzer *PDFDocumentAnalyzer, parameters *PAdESSignatureParameters) model.DSSDocument {
	signatures := pdfDocumentAnalyzer.Signatures()
	if utils.IsCollectionEmpty(signatures) {
		panic(exception.NewIllegalInputException("No signatures found to be extended!"))
	}

	if padesLevelBaselineTIsTLevelExtensionRequired(parameters, signatures) {
		signatureRequirementsChecker := NewPAdESSignatureRequirementsChecker(t.CertificateVerifier, parameters)

		signatureRequirementsChecker.AssertExtendToTLevelPossible(signatures)

		signatureRequirementsChecker.AssertSignaturesValid(signatures)
		signatureRequirementsChecker.AssertSigningCertificatesAreValid(signatures)

		// Will add a DocumentTimeStamp. signature-timestamp (CMS) is impossible to add while extending
		return t.TimestampDocument(signedDocument, parameters.SignatureTimestampParameters(),
			parameters.PasswordProtection(), t.signatureTimestampService())
	}

	return signedDocument
}

// signatureTimestampService returns a PDFSignatureService to be used for a signature timestamp
// creation. Port of the private #getSignatureTimestampService.
func (t *PAdESLevelBaselineT) signatureTimestampService() PDFSignatureService {
	return t.PdfObjectFactory.NewSignatureTimestampService()
}

// TimestampDocument timestamps the document with the given PDFSignatureService.
// Port of the protected #timestampDocument.
func (t *PAdESLevelBaselineT) TimestampDocument(signedDocument model.DSSDocument,
	timestampParameters *PAdESTimestampParameters, pwd []byte,
	pdfSignatureService PDFSignatureService) model.DSSDocument {
	padesTimestampService := NewPAdESTimestampServiceWithPDFService(t.tspSource, pdfSignatureService)
	timestampParameters.SetPasswordProtection(pwd)
	return padesTimestampService.TimestampDocument(signedDocument, timestampParameters)
}

// PDFDocumentValidator returns a document analyzer instance configured for the extension.
// Port of the protected #getPDFDocumentValidator.
func (t *PAdESLevelBaselineT) PDFDocumentValidator(signedDocument model.DSSDocument,
	parameters *PAdESSignatureParameters) *PDFDocumentAnalyzer {
	pdfDocumentValidator := NewPDFDocumentAnalyzer(signedDocument)
	pdfDocumentValidator.SetCertificateVerifier(t.CertificateVerifier)
	pdfDocumentValidator.SetValidationContextExecutor(executor.CompleteValidationContextExecutorInstance)
	pdfDocumentValidator.SetPasswordProtection(parameters.PasswordProtection())
	pdfDocumentValidator.SetPdfObjFactory(t.PdfObjectFactory)
	return pdfDocumentValidator
}

// padesLevelBaselineTIsTLevelExtensionRequired ports the private #isTLevelExtensionRequired.
func padesLevelBaselineTIsTLevelExtensionRequired(parameters *PAdESSignatureParameters,
	signatures []validation.AdvancedSignature) bool {
	tLevelExtensionRequired := false
	for _, signature := range signatures {
		padesSignature := signature.(*PAdESSignature)
		if padesLevelBaselineTRequiresDocumentTimestamp(padesSignature, parameters) {
			tLevelExtensionRequired = true
		}
	}
	return tLevelExtensionRequired
}

// padesLevelBaselineTRequiresDocumentTimestamp ports the private #requiresDocumentTimestamp.
func padesLevelBaselineTRequiresDocumentTimestamp(signature *PAdESSignature,
	signatureParameters *PAdESSignatureParameters) bool {
	return enumerations.SignatureLevel_PAdES_BASELINE_T == signatureParameters.SignatureLevel() ||
		!signature.HasTProfile()
}

// Compile-time assertion standing in for Java's "implements SignatureExtension<PAdESSignatureParameters>".
var _ document.SignatureExtension[*PAdESSignatureParameters] = (*PAdESLevelBaselineT)(nil)
