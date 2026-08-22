// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/signature/PAdESLevelBaselineLT.java (DSS 6.5.RC1).
//
// Java's class is package-private and extends PAdESLevelBaselineT, overriding the protected
// three-argument extendSignatures. The Go type embeds PAdESLevelBaselineT and registers itself
// with InitPAdESLevelBaselineT, so the base's ExtendSignatures dispatches into the override
// below.
package pades

import (
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// PAdESLevelBaselineLT holds the PAdES Baseline LT signature profile.
type PAdESLevelBaselineLT struct {
	PAdESLevelBaselineT
}

// NewPAdESLevelBaselineLT is the default constructor.
// Port of PAdESLevelBaselineLT(TSPSource, CertificateVerifier, IPdfObjFactory).
func NewPAdESLevelBaselineLT(tspSource validation.TSPSource, certificateVerifier validation.CertificateVerifier,
	pdfObjectFactory IPdfObjFactory) *PAdESLevelBaselineLT {
	extension := &PAdESLevelBaselineLT{}
	extension.InitPAdESLevelBaselineT(extension, tspSource, certificateVerifier, pdfObjectFactory)
	return extension
}

// ExtendSignaturesWithAnalyzer extends the document to LT-level, adding a DSS dictionary
// revision after the T-level extension the base performs. Port of the protected
// extendSignatures(DSSDocument, PDFDocumentAnalyzer, PAdESSignatureParameters) override.
func (lt *PAdESLevelBaselineLT) ExtendSignaturesWithAnalyzer(signedDocument model.DSSDocument,
	pdfDocumentAnalyzer *PDFDocumentAnalyzer, parameters *PAdESSignatureParameters) model.DSSDocument {
	extendedDocument := lt.PAdESLevelBaselineT.ExtendSignaturesWithAnalyzer(signedDocument,
		pdfDocumentAnalyzer, parameters)
	tLevelAdded := extendedDocument != signedDocument
	if tLevelAdded { // check if T-level has been added
		pdfDocumentAnalyzer = lt.PDFDocumentValidator(extendedDocument, parameters)
	}

	signatures := pdfDocumentAnalyzer.Signatures()

	signatureRequirementsChecker := NewPAdESSignatureRequirementsChecker(lt.CertificateVerifier, parameters)
	if !tLevelAdded && enumerations.SignatureLevelPAdESBaselineLT == parameters.SignatureLevel() {
		signatureRequirementsChecker.AssertExtendToLTLevelPossible(signatures)
	}
	signatureRequirementsChecker.AssertSignaturesValid(signatures)
	signatureRequirementsChecker.AssertCertificateChainValidForLTLevel(signatures)

	detachedTimestamps := pdfDocumentAnalyzer.DetachedTimestamps()
	validationData := pdfDocumentAnalyzer.ValidationData(signatures, detachedTimestamps)

	signatureService := lt.padesSignatureService()
	return signatureService.AddDssDictionary(extendedDocument, validationData,
		parameters.PasswordProtection(), parameters.IsIncludeVRIDictionary())
}

// padesSignatureService returns a PDFSignatureService to be used for a DSS dictionary addition.
// Port of the private #getPAdESSignatureService.
func (lt *PAdESLevelBaselineLT) padesSignatureService() PDFSignatureService {
	return lt.PdfObjectFactory.NewPAdESSignatureService()
}

// Compile-time assertion standing in for Java's "extends PAdESLevelBaselineT".
var _ document.SignatureExtension[*PAdESSignatureParameters] = (*PAdESLevelBaselineLT)(nil)
