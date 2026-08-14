//go:build phase8

// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PDFDocumentValidatorFactory.java
// (DSS 6.5.RC1).
//
// INTEGRATOR NOTE: gated behind the `phase8` build tag so that `go build ./...` / `go vet ./...`
// / `go test ./...` are green for the rest of the module while dss/validation does not exist
// yet, following the cades_diagnostic_data_builder.go / cms_document_validator_factory.go
// precedent. Drop the tag once Phase 8 lands the package.
//
// BLOCKED FORWARD DEPENDENCIES: this class's Java interface,
// eu.europa.esig.dss.validation.DocumentValidatorFactory, and return type,
// eu.europa.esig.dss.validation.SignedDocumentValidator, belong to dss-validation, which
// PORTING_PLAN.md assigns to the not-yet-ported `validation` package (Phase 8).
// PDFDocumentValidator (eu.europa.esig.dss.pades.validation.PDFDocumentValidator) is not in this
// chunk's manifest either, but has already landed (pdf_document_validator.go, same phase8 gate)
// with its real, confirmed API: the package-private no-arg constructor is newPDFDocumentValidator()
// (callable here unqualified - same package), and the document constructor is
// NewPDFDocumentValidator(document) (*PDFDocumentValidator, error). IsSupported(model.DSSDocument)
// bool is inherited via the embedded dssvalidation.SignedDocumentValidatorBase, matching the
// cades/cms_document_validator_factory.go precedent this file otherwise follows exactly.
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
	validator, err := NewPDFDocumentValidator(document)
	if err != nil {
		panic(err)
	}
	return validator
}

// compile-time interface assertion.
var _ dssvalidation.DocumentValidatorFactory = (*PDFDocumentValidatorFactory)(nil)

func init() {
	dssvalidation.RegisterDocumentValidatorFactory(NewPDFDocumentValidatorFactory())
}
