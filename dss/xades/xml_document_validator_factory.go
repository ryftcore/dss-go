//go:build phase8

// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XMLDocumentValidatorFactory.java
// (DSS 6.5.RC1).
//
// INTEGRATOR NOTE: gated behind the `phase8` build tag, same as cades/cms_document_validator_factory.go -
// it imports dss/validation.SignedDocumentValidator/DocumentValidatorFactory, which do not exist
// until Phase 8 lands. Drop the tag once Phase 8 lands dss/validation.
package xades

import (
	"github.com/utain/esig/dss/model"
	dssvalidation "github.com/utain/esig/dss/validation"
)

// XMLDocumentValidatorFactory loads the relevant validator for an XML document validation. Port
// of the class XMLDocumentValidatorFactory, implementing validation.DocumentValidatorFactory.
type XMLDocumentValidatorFactory struct{}

// NewXMLDocumentValidatorFactory is the port of the default constructor.
func NewXMLDocumentValidatorFactory() *XMLDocumentValidatorFactory {
	return &XMLDocumentValidatorFactory{}
}

// IsSupported is the port of isSupported(DSSDocument).
func (f *XMLDocumentValidatorFactory) IsSupported(document model.DSSDocument) bool {
	validator := newXMLDocumentValidator()
	return validator.IsSupported(document)
}

// Create is the port of create(DSSDocument).
func (f *XMLDocumentValidatorFactory) Create(document model.DSSDocument) dssvalidation.SignedDocumentValidator {
	validator, err := NewXMLDocumentValidator(document)
	if err != nil {
		panic(err)
	}
	return validator
}

// compile-time interface assertion.
var _ dssvalidation.DocumentValidatorFactory = (*XMLDocumentValidatorFactory)(nil)

func init() {
	dssvalidation.RegisterDocumentValidatorFactory(NewXMLDocumentValidatorFactory())
}
