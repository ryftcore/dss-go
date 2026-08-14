//go:build phase8

// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESDocumentValidatorFactory.java (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (phase 6 integration): gated behind the `phase8` build tag; see
// abstract_jws_document_validator.go's file header for the assumed dss/validation shape
// (DocumentValidatorFactory).
package jades

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/exception"
	dssvalidation "github.com/utain/esig/dss/validation"
)

// JAdESDocumentValidatorFactory loads the relevant Validator to process a given JAdES signature.
// Port of the class JAdESDocumentValidatorFactory, implementing validation.DocumentValidatorFactory.
type JAdESDocumentValidatorFactory struct{}

// NewJAdESDocumentValidatorFactory is the port of the default constructor.
func NewJAdESDocumentValidatorFactory() *JAdESDocumentValidatorFactory {
	return &JAdESDocumentValidatorFactory{}
}

// IsSupported is the port of isSupported(DSSDocument).
func (f *JAdESDocumentValidatorFactory) IsSupported(document model.DSSDocument) bool {
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
func (f *JAdESDocumentValidatorFactory) Create(document model.DSSDocument) dssvalidation.SignedDocumentValidator {
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
var _ dssvalidation.DocumentValidatorFactory = (*JAdESDocumentValidatorFactory)(nil)

func init() {
	dssvalidation.RegisterDocumentValidatorFactory(NewJAdESDocumentValidatorFactory())
}
