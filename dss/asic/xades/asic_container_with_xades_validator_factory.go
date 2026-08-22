// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/validation/ASiCContainerWithXAdESValidatorFactory.java (DSS 6.5.RC1).
//
// Package flattening: Java's eu.europa.esig.dss.asic.xades.validation lands in this
// same Go package (dss/asic/xades).
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
	dssvalidation "github.com/ryftcore/dss-go/dss/validation"
)

// ASiCContainerWithXAdESValidatorFactory returns a relevant validator for an ASiC with XAdES
// container validation.
type ASiCContainerWithXAdESValidatorFactory struct{}

var _ dssvalidation.DocumentValidatorFactory = (*ASiCContainerWithXAdESValidatorFactory)(nil)

// NewASiCContainerWithXAdESValidatorFactory is the default constructor.
func NewASiCContainerWithXAdESValidatorFactory() *ASiCContainerWithXAdESValidatorFactory {
	return &ASiCContainerWithXAdESValidatorFactory{}
}

// IsSupported ports the @Override isSupported(DSSDocument).
func (f *ASiCContainerWithXAdESValidatorFactory) IsSupported(document model.DSSDocument) bool {
	validator := newASiCContainerWithXAdESValidator()
	return validator.IsSupported(document)
}

// IsSupportedContent verifies whether the provided ASiCContent is supported by the underlying
// validator's class. Ports isSupported(ASiCContent).
func (f *ASiCContainerWithXAdESValidatorFactory) IsSupportedContent(asicContent *asic.ASiCContent) bool {
	validator := newASiCContainerWithXAdESValidator()
	return validator.IsSupportedContent(asicContent)
}

// Create ports the @Override create(DSSDocument).
func (f *ASiCContainerWithXAdESValidatorFactory) Create(document model.DSSDocument) dssvalidation.SignedDocumentValidator {
	return NewASiCContainerWithXAdESValidator(document)
}

// CreateFromContent creates a SignedDocumentValidator for the given asicContent. Ports
// create(ASiCContent).
func (f *ASiCContainerWithXAdESValidatorFactory) CreateFromContent(asicContent *asic.ASiCContent) dssvalidation.SignedDocumentValidator {
	return NewASiCContainerWithXAdESValidatorFromContent(asicContent)
}

func init() {
	dssvalidation.RegisterDocumentValidatorFactory(NewASiCContainerWithXAdESValidatorFactory())
}
