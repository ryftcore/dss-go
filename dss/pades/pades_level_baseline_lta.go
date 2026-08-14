// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/signature/PAdESLevelBaselineLTA.java (DSS 6.5.RC1).
//
// Java extends PAdESLevelBaselineLT and overrides the public two-argument extendSignatures. The
// Go type embeds PAdESLevelBaselineLT and shadows ExtendSignatures; the super call goes to the
// embedded PAdESLevelBaselineT.ExtendSignatures, which dispatches back through the overrides
// interface into PAdESLevelBaselineLT.ExtendSignaturesWithAnalyzer - exactly Java's chain.
package pades

import (
	"github.com/utain/esig/dss/document"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
)

// PAdESLevelBaselineLTA holds the PAdES Baseline LTA signature profile.
type PAdESLevelBaselineLTA struct {
	PAdESLevelBaselineLT
}

// NewPAdESLevelBaselineLTA is the default constructor.
// Port of PAdESLevelBaselineLTA(TSPSource, CertificateVerifier, IPdfObjFactory).
func NewPAdESLevelBaselineLTA(tspSource validation.TSPSource, certificateVerifier validation.CertificateVerifier,
	pdfObjectFactory IPdfObjFactory) *PAdESLevelBaselineLTA {
	extension := &PAdESLevelBaselineLTA{}
	extension.InitPAdESLevelBaselineT(extension, tspSource, certificateVerifier, pdfObjectFactory)
	return extension
}

// ExtendSignatures extends the document to LT-level when needed and then adds a document
// timestamp (not CMS). Port of the
// extendSignatures(DSSDocument, PAdESSignatureParameters) override.
func (lta *PAdESLevelBaselineLTA) ExtendSignatures(signedDocument model.DSSDocument,
	parameters *PAdESSignatureParameters) model.DSSDocument {
	// check if needed to extend with PAdESLevelBaselineLT
	signedDocument = lta.PAdESLevelBaselineT.ExtendSignatures(signedDocument, parameters)

	// Will add a Document TimeStamp (not CMS)
	return lta.TimestampDocument(signedDocument, parameters.ArchiveTimestampParameters(),
		parameters.PasswordProtection(), lta.archiveTimestampService())
}

// archiveTimestampService returns a PDFSignatureService to be used for an archive timestamp
// creation. Port of the private #getArchiveTimestampService.
func (lta *PAdESLevelBaselineLTA) archiveTimestampService() PDFSignatureService {
	return lta.PdfObjectFactory.NewArchiveTimestampService()
}

// Compile-time assertion standing in for Java's "extends PAdESLevelBaselineLT".
var _ document.SignatureExtension[*PAdESSignatureParameters] = (*PAdESLevelBaselineLTA)(nil)
