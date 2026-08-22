// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESDocumentValidatorFactory.java (DSS 6.5.RC1).
package jades

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	dssvalidation "github.com/ryftcore/dss-go/dss/validation"
)

// DocumentValidatorFactory loads the relevant Validator to process a given JAdES signature.
// Port of the class JAdESDocumentValidatorFactory, implementing validation.DocumentValidatorFactory.
type DocumentValidatorFactory struct{}

// NewJAdESDocumentValidatorFactory is the port of the default constructor.
func NewJAdESDocumentValidatorFactory() *DocumentValidatorFactory {
	return &DocumentValidatorFactory{}
}

// IsSupported is the port of isSupported(DSSDocument).
func (f *DocumentValidatorFactory) IsSupported(document model.DSSDocument) bool {
	compactValidator := NewJWSCompactDocumentValidator()
	if compactValidator.IsSupported(document) {
		return true
	}

	serializationValidator := NewJWSSerializationDocumentValidator()
	return serializationValidator.IsSupported(document)
}

// Create is the port of create(DSSDocument).
//
// Panics with the Java message when document is not supported (Java throws
// IllegalArgumentException, unchecked).
func (f *DocumentValidatorFactory) Create(document model.DSSDocument) dssvalidation.SignedDocumentValidator {
	compactValidator := NewJWSCompactDocumentValidator()
	if compactValidator.IsSupported(document) {
		return NewJWSCompactDocumentValidatorFromDocument(document)
	}

	serializationValidator := NewJWSSerializationDocumentValidator()
	if serializationValidator.IsSupported(document) {
		return NewJWSSerializationDocumentValidatorFromDocument(document)
	}

	panic(exception.NewIllegalInputException("Not supported document"))
}

// compile-time interface assertion.
var _ dssvalidation.DocumentValidatorFactory = (*DocumentValidatorFactory)(nil)

func init() {
	dssvalidation.RegisterDocumentValidatorFactory(NewJAdESDocumentValidatorFactory())
}
