// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/analyzer/evidencerecord/EvidenceRecordAnalyzerFactory.java (DSS 6.5.RC1).
package analyzer

import "github.com/utain/esig/dss/model"

// EvidenceRecordAnalyzerFactory is used to load a corresponding implementation of
// EvidenceRecordAnalyzer for processing of an evidence record document.
type EvidenceRecordAnalyzerFactory interface {
	// IsSupported tests if the current implementation of EvidenceRecordAnalyzer supports the
	// given document. Port of isSupported(DSSDocument).
	IsSupported(document model.DSSDocument) bool

	// Create instantiates an EvidenceRecordAnalyzer with the given document. Port of
	// create(DSSDocument).
	Create(document model.DSSDocument) EvidenceRecordAnalyzer
}

// evidenceRecordAnalyzerFactoryRegistry holds the EvidenceRecordAnalyzerFactory implementations
// registered via RegisterEvidenceRecordAnalyzerFactory, consulted in registration order — the
// Go equivalent of Java's ServiceLoader.load(EvidenceRecordAnalyzerFactory.class) iteration (Go
// has no runtime service-provider discovery). Evidence-record-format-specific implementations
// register themselves here in later phases.
var evidenceRecordAnalyzerFactoryRegistry []EvidenceRecordAnalyzerFactory

// RegisterEvidenceRecordAnalyzerFactory registers an EvidenceRecordAnalyzerFactory to be
// consulted by EvidenceRecordAnalyzerIsSupportedDocument and EvidenceRecordAnalyzerFromDocument.
func RegisterEvidenceRecordAnalyzerFactory(f EvidenceRecordAnalyzerFactory) {
	evidenceRecordAnalyzerFactoryRegistry = append(evidenceRecordAnalyzerFactoryRegistry, f)
}

// EvidenceRecordAnalyzerIsSupportedDocument verifies if document is supported by one of the
// registered EvidenceRecordAnalyzerFactory implementations. Port of the static
// isSupportedDocument(DSSDocument).
//
// Panics when document is nil (Objects.requireNonNull("DSSDocument is null")).
func EvidenceRecordAnalyzerIsSupportedDocument(document model.DSSDocument) bool {
	if document == nil {
		panic("DSSDocument is null")
	}
	for _, factory := range evidenceRecordAnalyzerFactoryRegistry {
		if factory.IsSupported(document) {
			return true
		}
	}
	return false
}

// EvidenceRecordAnalyzerFromDocument creates an EvidenceRecordAnalyzer by finding a
// corresponding registered implementation. Port of the static fromDocument(DSSDocument).
//
// Panics when document is nil, matching EvidenceRecordAnalyzerIsSupportedDocument. Java's
// UnsupportedOperationException("Document format not recognized/handled"), thrown when no
// registered implementation supports the document, is returned as an error instead: whether a
// document format is recognized is data-dependent on which analyzer implementations happen to
// be registered (loaded modules), matching the "not supported in the Go port" precedent used
// for similarly optional runtime capabilities elsewhere in this port.
func EvidenceRecordAnalyzerFromDocument(document model.DSSDocument) (EvidenceRecordAnalyzer, error) {
	if document == nil {
		panic("DSSDocument is null")
	}
	for _, factory := range evidenceRecordAnalyzerFactoryRegistry {
		if factory.IsSupported(document) {
			return factory.Create(document), nil
		}
	}
	return nil, errDocumentFormatNotRecognized
}
