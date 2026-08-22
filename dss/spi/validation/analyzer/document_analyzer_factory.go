// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/analyzer/DocumentAnalyzerFactory.java (DSS 6.5.RC1).
package analyzer

import (
	"errors"

	"github.com/ryftcore/dss-go/dss/model"
)

// errDocumentFormatNotRecognized is shared by DocumentAnalyzerFromDocument,
// EvidenceRecordAnalyzerFromDocument and (dss/spi/validation/analyzer/eaa)
// EAAPresentationAnalyzerFromDocument as the port of Java's
// UnsupportedOperationException("Document format not recognized/handled").
var errDocumentFormatNotRecognized = errors.New("document format not recognized/handled")

// DocumentAnalyzerFactory is used to analyze the format of the given DSSDocument and create a
// corresponding implementation of DocumentAnalyzer.
type DocumentAnalyzerFactory interface {
	// IsSupported tests if the current implementation of DocumentAnalyzer supports the given
	// document. Port of isSupported(DSSDocument).
	IsSupported(document model.DSSDocument) bool

	// Create instantiates a DocumentAnalyzer with the given document. Port of
	// create(DSSDocument).
	Create(document model.DSSDocument) DocumentAnalyzer
}

// documentAnalyzerFactoryRegistry holds the DocumentAnalyzerFactory implementations registered
// via RegisterDocumentAnalyzerFactory, consulted in registration order — the Go equivalent of
// Java's ServiceLoader.load(DocumentAnalyzerFactory.class) iteration (Go has no runtime
// service-provider discovery). Format-specific implementations (CAdES/XAdES/JAdES/PAdES/ASiC)
// register themselves here in later phases.
var documentAnalyzerFactoryRegistry []DocumentAnalyzerFactory

// RegisterDocumentAnalyzerFactory registers a DocumentAnalyzerFactory to be consulted by
// DocumentAnalyzerFromDocument.
func RegisterDocumentAnalyzerFactory(f DocumentAnalyzerFactory) {
	documentAnalyzerFactoryRegistry = append(documentAnalyzerFactoryRegistry, f)
}
