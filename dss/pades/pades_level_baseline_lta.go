// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/signature/PAdESLevelBaselineLTA.java (DSS 6.5.RC1).
//
// Java extends PAdESLevelBaselineLT and overrides the public two-argument extendSignatures. The
// Go type embeds LevelBaselineLT and shadows ExtendSignatures; the super call goes to the
// embedded LevelBaselineT.ExtendSignatures, which dispatches back through the overrides
// interface into PAdESLevelBaselineLT.ExtendSignaturesWithAnalyzer - exactly Java's chain.
package pades

import (
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// LevelBaselineLTA holds the PAdES Baseline LTA signature profile.
type LevelBaselineLTA struct {
	LevelBaselineLT
}

// NewLevelBaselineLTA is the default constructor.
// Port of PAdESLevelBaselineLTA(TSPSource, CertificateVerifier, IPdfObjFactory).
func NewLevelBaselineLTA(tspSource validation.TSPSource, certificateVerifier validation.CertificateVerifier,
	pdfObjectFactory IPdfObjFactory) *LevelBaselineLTA {
	extension := &LevelBaselineLTA{}
	extension.InitPAdESLevelBaselineT(extension, tspSource, certificateVerifier, pdfObjectFactory)
	return extension
}

// ExtendSignatures extends the document to LT-level when needed and then adds a document
// timestamp (not CMS). Port of the
// extendSignatures(DSSDocument, SignatureParameters) override.
func (lta *LevelBaselineLTA) ExtendSignatures(signedDocument model.DSSDocument,
	parameters *SignatureParameters) model.DSSDocument {
	// check if needed to extend with LevelBaselineLT
	signedDocument = lta.LevelBaselineT.ExtendSignatures(signedDocument, parameters)

	// Will add a Document TimeStamp (not CMS)
	return lta.TimestampDocument(signedDocument, parameters.ArchiveTimestampParameters(),
		parameters.PasswordProtection(), lta.archiveTimestampService())
}

// archiveTimestampService returns a PDFSignatureService to be used for an archive timestamp
// creation. Port of the private #getArchiveTimestampService.
func (lta *LevelBaselineLTA) archiveTimestampService() PDFSignatureService {
	return lta.PdfObjectFactory.NewArchiveTimestampService()
}

// Compile-time assertion standing in for Java's "extends PAdESLevelBaselineLT".
var _ document.SignatureExtension[*SignatureParameters] = (*LevelBaselineLTA)(nil)
