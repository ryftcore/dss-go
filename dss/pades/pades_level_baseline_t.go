// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/signature/PAdESLevelBaselineT.java (DSS 6.5.RC1).
//
// Java's class is package-private and implements SignatureExtension<PAdESSignatureParameters>.
// The Go type is exported (Go has no package-private-yet-subclassable visibility, and the whole
// PAdES port lives in one package anyway).
//
// Java's two extendSignatures overloads cannot share one Go name:
//
//	extendSignatures(DSSDocument, SignatureParameters)                        -> ExtendSignatures
//	extendSignatures(DSSDocument, PDFDocumentAnalyzer, SignatureParameters)   -> ExtendSignaturesWithAnalyzer
//
// The second is overridden by LevelBaselineLT and called from the first, so the base
// dispatches through LevelBaselineTOverrides, registered by InitPAdESLevelBaselineT - the
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

// LevelBaselineTOverrides declares the operation LevelBaselineT calls back into and
// LevelBaselineLT overrides.
type LevelBaselineTOverrides interface {
	// ExtendSignaturesWithAnalyzer performs the document extension for the signatures the given
	// analyzer found. Port of the protected
	// extendSignatures(DSSDocument, PDFDocumentAnalyzer, SignatureParameters).
	ExtendSignaturesWithAnalyzer(signedDocument model.DSSDocument, pdfDocumentAnalyzer *PDFDocumentAnalyzer,
		parameters *SignatureParameters) model.DSSDocument
}

// LevelBaselineT holds the PAdES Baseline T signature profile.
type LevelBaselineT struct {
	// tspSource obtains a timestamp.
	tspSource validation.TSPSource

	// CertificateVerifier is the used CertificateVerifier.
	CertificateVerifier validation.CertificateVerifier

	// PdfObjectFactory is the used implementation for processing of a PDF document.
	PdfObjectFactory IPdfObjFactory

	// overrides points back at the concrete extension; see InitPAdESLevelBaselineT.
	overrides LevelBaselineTOverrides
}

// NewLevelBaselineT is the default constructor.
// Port of PAdESLevelBaselineT(TSPSource, CertificateVerifier, IPdfObjFactory).
func NewLevelBaselineT(tspSource validation.TSPSource, certificateVerifier validation.CertificateVerifier,
	pdfObjectFactory IPdfObjFactory) *LevelBaselineT {
	extension := &LevelBaselineT{}
	extension.InitPAdESLevelBaselineT(extension, tspSource, certificateVerifier, pdfObjectFactory)
	return extension
}

// InitPAdESLevelBaselineT registers the concrete extension with its base and applies the
// constructor's checks. Port of the protected
// LevelBaselineT(TSPSource, CertificateVerifier, IPdfObjFactory) constructor.
func (t *LevelBaselineT) InitPAdESLevelBaselineT(self LevelBaselineTOverrides,
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
// extendSignatures(DSSDocument, SignatureParameters).
func (t *LevelBaselineT) ExtendSignatures(signedDocument model.DSSDocument,
	params *SignatureParameters) model.DSSDocument {
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
// extendSignatures(DSSDocument, PDFDocumentAnalyzer, SignatureParameters).
func (t *LevelBaselineT) ExtendSignaturesWithAnalyzer(signedDocument model.DSSDocument,
	pdfDocumentAnalyzer *PDFDocumentAnalyzer, parameters *SignatureParameters) model.DSSDocument {
	signatures := pdfDocumentAnalyzer.Signatures()
	if utils.IsCollectionEmpty(signatures) {
		panic(exception.NewIllegalInputException("No signatures found to be extended!"))
	}

	if padesLevelBaselineTIsTLevelExtensionRequired(parameters, signatures) {
		signatureRequirementsChecker := NewSignatureRequirementsChecker(t.CertificateVerifier, parameters)

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
func (t *LevelBaselineT) signatureTimestampService() PDFSignatureService {
	return t.PdfObjectFactory.NewSignatureTimestampService()
}

// TimestampDocument timestamps the document with the given PDFSignatureService.
// Port of the protected #timestampDocument.
func (t *LevelBaselineT) TimestampDocument(signedDocument model.DSSDocument,
	timestampParameters *TimestampParameters, pwd []byte,
	pdfSignatureService PDFSignatureService) model.DSSDocument {
	padesTimestampService := NewTimestampServiceWithPDFService(t.tspSource, pdfSignatureService)
	timestampParameters.SetPasswordProtection(pwd)
	return padesTimestampService.TimestampDocument(signedDocument, timestampParameters)
}

// PDFDocumentValidator returns a document analyzer instance configured for the extension.
// Port of the protected #getPDFDocumentValidator.
func (t *LevelBaselineT) PDFDocumentValidator(signedDocument model.DSSDocument,
	parameters *SignatureParameters) *PDFDocumentAnalyzer {
	pdfDocumentValidator := NewPDFDocumentAnalyzer(signedDocument)
	pdfDocumentValidator.SetCertificateVerifier(t.CertificateVerifier)
	pdfDocumentValidator.SetValidationContextExecutor(executor.CompleteValidationContextExecutorInstance)
	pdfDocumentValidator.SetPasswordProtection(parameters.PasswordProtection())
	pdfDocumentValidator.SetPdfObjFactory(t.PdfObjectFactory)
	return pdfDocumentValidator
}

// padesLevelBaselineTIsTLevelExtensionRequired ports the private #isTLevelExtensionRequired.
func padesLevelBaselineTIsTLevelExtensionRequired(parameters *SignatureParameters,
	signatures []validation.AdvancedSignature) bool {
	tLevelExtensionRequired := false
	for _, signature := range signatures {
		padesSignature := signature.(*Signature)
		if padesLevelBaselineTRequiresDocumentTimestamp(padesSignature, parameters) {
			tLevelExtensionRequired = true
		}
	}
	return tLevelExtensionRequired
}

// padesLevelBaselineTRequiresDocumentTimestamp ports the private #requiresDocumentTimestamp.
func padesLevelBaselineTRequiresDocumentTimestamp(signature *Signature,
	signatureParameters *SignatureParameters) bool {
	return enumerations.SignatureLevelPAdESBaselineT == signatureParameters.SignatureLevel() ||
		!signature.HasTProfile()
}

// Compile-time assertion standing in for Java's "implements SignatureExtension<PAdESSignatureParameters>".
var _ document.SignatureExtension[*SignatureParameters] = (*LevelBaselineT)(nil)
