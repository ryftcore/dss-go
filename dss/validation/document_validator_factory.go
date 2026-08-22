// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/DocumentValidatorFactory.java
// (DSS 6.5.RC1).

package validation

import (
	"github.com/ryftcore/dss-go/dss/model"
)

// DocumentValidatorFactory defines the factory to create a DocumentValidator
// for a given DSSDocument. Port of the DocumentValidatorFactory interface.
type DocumentValidatorFactory interface {

	// IsSupported tests if the current implementation of DocumentValidator
	// supports the given document. Port of isSupported(DSSDocument).
	IsSupported(document model.DSSDocument) bool

	// Create instantiates a DocumentValidator with the given document. Port
	// of create(DSSDocument).
	Create(document model.DSSDocument) SignedDocumentValidator
}

// documentValidatorFactoryRegistry holds the DocumentValidatorFactory
// implementations registered via RegisterDocumentValidatorFactory, consulted
// in registration order - the Go equivalent of Java's
// ServiceLoader.load(DocumentValidatorFactory.class) iteration performed by
// SignedDocumentValidator.fromDocument (Go has no runtime service-provider
// discovery). It mirrors the registry already used by
// dss/spi/validation/analyzer.RegisterDocumentAnalyzerFactory.
//
// Registration order, not init order, decides the winner: the format packages
// register at most once each and SignedDocumentValidatorFromDocument walks the
// slice front to back, exactly as ServiceLoader walks the provider list.
var documentValidatorFactoryRegistry []DocumentValidatorFactory

// RegisterDocumentValidatorFactory registers a DocumentValidatorFactory to be
// consulted by SignedDocumentValidatorFromDocument.
func RegisterDocumentValidatorFactory(f DocumentValidatorFactory) {
	documentValidatorFactoryRegistry = append(documentValidatorFactoryRegistry, f)
}

// DocumentValidatorFactories returns the registered factories in registration
// order. Exposed for tests and for diagnostics of the registration wiring; it
// has no Java counterpart (Java inspects the ServiceLoader directly).
func DocumentValidatorFactories() []DocumentValidatorFactory {
	result := make([]DocumentValidatorFactory, len(documentValidatorFactoryRegistry))
	copy(result, documentValidatorFactoryRegistry)
	return result
}
