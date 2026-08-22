// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XMLDocumentValidatorFactory.java
// (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/model"
	dssvalidation "github.com/ryftcore/dss-go/dss/validation"
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
