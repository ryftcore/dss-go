// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/analyzer/eaa/EAAPresentationAnalyzerFactory.java (DSS 6.5.RC1).
package eaa

import (
	"errors"

	"github.com/ryftcore/dss-go/dss/model"
)

// PresentationAnalyzerFactory is used to load a specific PresentationAnalyzer based on
// the format of the provided presentation of Electronic Attestation of Attributes.
type PresentationAnalyzerFactory interface {
	// IsSupported tests if the current implementation of PresentationAnalyzer supports the
	// given document. Port of isSupported(DSSDocument).
	IsSupported(document model.DSSDocument) bool

	// Create instantiates an EAAPresentationAnalyzer with the given document. Port of
	// create(DSSDocument).
	Create(document model.DSSDocument) PresentationAnalyzer
}

// eaaPresentationAnalyzerFactoryRegistry holds the PresentationAnalyzerFactory
// implementations registered via RegisterPresentationAnalyzerFactory, consulted in
// registration order — the Go equivalent of Java's
// ServiceLoader.load(PresentationAnalyzerFactory.class) iteration (Go has no runtime
// service-provider discovery). EAA-presentation-format-specific implementations register
// themselves here in later phases.
var eaaPresentationAnalyzerFactoryRegistry []PresentationAnalyzerFactory

// RegisterPresentationAnalyzerFactory registers an PresentationAnalyzerFactory to be
// consulted by PresentationAnalyzerIsSupportedDocument and
// PresentationAnalyzerFromDocument.
func RegisterPresentationAnalyzerFactory(f PresentationAnalyzerFactory) {
	eaaPresentationAnalyzerFactoryRegistry = append(eaaPresentationAnalyzerFactoryRegistry, f)
}

// errDocumentFormatNotRecognized is the port of Java's
// UnsupportedOperationException("Document format not recognized/handled").
var errDocumentFormatNotRecognized = errors.New("document format not recognized/handled")

// PresentationAnalyzerIsSupportedDocument verifies if document is supported by one of the
// registered EAAPresentationAnalyzerFactory implementations. Port of the static
// isSupportedDocument(DSSDocument).
//
// Panics when document is nil (Objects.requireNonNull("DSSDocument is null")).
func PresentationAnalyzerIsSupportedDocument(document model.DSSDocument) bool {
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

// PresentationAnalyzerFromDocument creates an PresentationAnalyzer by finding a
// corresponding registered implementation. Port of the static fromDocument(DSSDocument).
//
// Panics when document is nil, matching EAAPresentationAnalyzerIsSupportedDocument. Java's
// UnsupportedOperationException("Document format not recognized/handled"), thrown when no
// registered implementation supports the document, is returned as an error: whether a document
// format is recognized is data-dependent on which analyzer implementations happen to be
// registered (loaded modules).
func PresentationAnalyzerFromDocument(document model.DSSDocument) (PresentationAnalyzer, error) {
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
