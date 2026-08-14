// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/analyzer/eaa/EAAPresentationAnalyzerFactory.java (DSS 6.5.RC1).
package eaa

import (
	"errors"

	"github.com/utain/esig/dss/model"
)

// EAAPresentationAnalyzerFactory is used to load a specific EAAPresentationAnalyzer based on
// the format of the provided presentation of Electronic Attestation of Attributes.
type EAAPresentationAnalyzerFactory interface {
	// IsSupported tests if the current implementation of EAAPresentationAnalyzer supports the
	// given document. Port of isSupported(DSSDocument).
	IsSupported(document model.DSSDocument) bool

	// Create instantiates an EAAPresentationAnalyzer with the given document. Port of
	// create(DSSDocument).
	Create(document model.DSSDocument) EAAPresentationAnalyzer
}

// eaaPresentationAnalyzerFactoryRegistry holds the EAAPresentationAnalyzerFactory
// implementations registered via RegisterEAAPresentationAnalyzerFactory, consulted in
// registration order — the Go equivalent of Java's
// ServiceLoader.load(EAAPresentationAnalyzerFactory.class) iteration (Go has no runtime
// service-provider discovery). EAA-presentation-format-specific implementations register
// themselves here in later phases.
var eaaPresentationAnalyzerFactoryRegistry []EAAPresentationAnalyzerFactory

// RegisterEAAPresentationAnalyzerFactory registers an EAAPresentationAnalyzerFactory to be
// consulted by EAAPresentationAnalyzerIsSupportedDocument and
// EAAPresentationAnalyzerFromDocument.
func RegisterEAAPresentationAnalyzerFactory(f EAAPresentationAnalyzerFactory) {
	eaaPresentationAnalyzerFactoryRegistry = append(eaaPresentationAnalyzerFactoryRegistry, f)
}

// errDocumentFormatNotRecognized is the port of Java's
// UnsupportedOperationException("Document format not recognized/handled").
var errDocumentFormatNotRecognized = errors.New("document format not recognized/handled")

// EAAPresentationAnalyzerIsSupportedDocument verifies if document is supported by one of the
// registered EAAPresentationAnalyzerFactory implementations. Port of the static
// isSupportedDocument(DSSDocument).
//
// Panics when document is nil (Objects.requireNonNull("DSSDocument is null")).
func EAAPresentationAnalyzerIsSupportedDocument(document model.DSSDocument) bool {
	if document == nil {
		panic("DSSDocument is null")
	}
	for _, factory := range eaaPresentationAnalyzerFactoryRegistry {
		if factory.IsSupported(document) {
			return true
		}
	}
	return false
}

// EAAPresentationAnalyzerFromDocument creates an EAAPresentationAnalyzer by finding a
// corresponding registered implementation. Port of the static fromDocument(DSSDocument).
//
// Panics when document is nil, matching EAAPresentationAnalyzerIsSupportedDocument. Java's
// UnsupportedOperationException("Document format not recognized/handled"), thrown when no
// registered implementation supports the document, is returned as an error: whether a document
// format is recognized is data-dependent on which analyzer implementations happen to be
// registered (loaded modules).
func EAAPresentationAnalyzerFromDocument(document model.DSSDocument) (EAAPresentationAnalyzer, error) {
	if document == nil {
		panic("DSSDocument is null")
	}
	for _, factory := range eaaPresentationAnalyzerFactoryRegistry {
		if factory.IsSupported(document) {
			return factory.Create(document), nil
		}
	}
	return nil, errDocumentFormatNotRecognized
}
