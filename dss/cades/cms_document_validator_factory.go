//go:build phase8

// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CMSDocumentValidatorFactory.java (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (Phase 3 integration): gated behind the `phase8` build tag, same as
// cms_document_validator.go - see that file's header. Drop the tag once Phase 8 lands
// dss/validation.
//
// BLOCKED FORWARD DEPENDENCY: see cms_document_validator.go's header - the same applies here.
package cades

import (
	"github.com/utain/esig/dss/model"
	dssvalidation "github.com/utain/esig/dss/validation"
)

// CMSDocumentValidatorFactory checks if the document is supported and creates a relevant
// validator for the provided document. Port of the class CMSDocumentValidatorFactory,
// implementing validation.DocumentValidatorFactory.
type CMSDocumentValidatorFactory struct{}

// NewCMSDocumentValidatorFactory is the port of the default constructor.
func NewCMSDocumentValidatorFactory() *CMSDocumentValidatorFactory {
	return &CMSDocumentValidatorFactory{}
}

// IsSupported is the port of isSupported(DSSDocument).
func (f *CMSDocumentValidatorFactory) IsSupported(document model.DSSDocument) bool {
	validator := newCMSDocumentValidator()
	return validator.IsSupported(document)
}

// Create is the port of create(DSSDocument).
func (f *CMSDocumentValidatorFactory) Create(document model.DSSDocument) dssvalidation.SignedDocumentValidator {
	validator, err := NewCMSDocumentValidatorFromDocument(document)
	if err != nil {
		panic(err)
	}
	return validator
}

// compile-time interface assertion.
var _ dssvalidation.DocumentValidatorFactory = (*CMSDocumentValidatorFactory)(nil)

func init() {
	dssvalidation.RegisterDocumentValidatorFactory(NewCMSDocumentValidatorFactory())
}
