// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PDFDocumentValidatorFactory.java
// (DSS 6.5.RC1).
package pades

import (
	"github.com/utain/esig/dss/model"
	dssvalidation "github.com/utain/esig/dss/validation"
)

// PDFDocumentValidatorFactory loads a relevant validator for a PDF document.
type PDFDocumentValidatorFactory struct{}

// NewPDFDocumentValidatorFactory is the port of the default constructor.
func NewPDFDocumentValidatorFactory() *PDFDocumentValidatorFactory {
	return &PDFDocumentValidatorFactory{}
}

// IsSupported is the port of isSupported(DSSDocument).
func (f *PDFDocumentValidatorFactory) IsSupported(document model.DSSDocument) bool {
	validator := newPDFDocumentValidator()
	return validator.IsSupported(document)
}

// Create is the port of create(DSSDocument).
func (f *PDFDocumentValidatorFactory) Create(document model.DSSDocument) dssvalidation.SignedDocumentValidator {
	return NewPDFDocumentValidator(document)
}

// compile-time interface assertion.
var _ dssvalidation.DocumentValidatorFactory = (*PDFDocumentValidatorFactory)(nil)

func init() {
	dssvalidation.RegisterDocumentValidatorFactory(NewPDFDocumentValidatorFactory())
}
