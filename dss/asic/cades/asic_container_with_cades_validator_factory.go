// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/validation/ASiCContainerWithCAdESValidatorFactory.java (DSS 6.5.RC1).
package cades

import (
	"github.com/utain/esig/dss/asic"
	"github.com/utain/esig/dss/model"
	dssvalidation "github.com/utain/esig/dss/validation"
)

// ASiCContainerWithCAdESValidatorFactory returns a relevant validator for an ASiC with CAdES
// container validation. Port of the class ASiCContainerWithCAdESValidatorFactory, implementing
// validation.DocumentValidatorFactory.
type ASiCContainerWithCAdESValidatorFactory struct{}

var _ dssvalidation.DocumentValidatorFactory = (*ASiCContainerWithCAdESValidatorFactory)(nil)

// NewASiCContainerWithCAdESValidatorFactory is the default constructor.
func NewASiCContainerWithCAdESValidatorFactory() *ASiCContainerWithCAdESValidatorFactory {
	return &ASiCContainerWithCAdESValidatorFactory{}
}

// IsSupported ports the @Override isSupported(DSSDocument).
func (f *ASiCContainerWithCAdESValidatorFactory) IsSupported(document model.DSSDocument) bool {
	validator := newASiCContainerWithCAdESValidator()
	return validator.IsSupported(document)
}

// IsSupportedContent verifies whether the provided ASiCContent is supported by the underlying
// validator's class. Ports isSupported(ASiCContent).
func (f *ASiCContainerWithCAdESValidatorFactory) IsSupportedContent(asicContent *asic.ASiCContent) bool {
	validator := newASiCContainerWithCAdESValidator()
	return validator.IsSupportedContent(asicContent)
}

// Create ports the @Override create(DSSDocument).
func (f *ASiCContainerWithCAdESValidatorFactory) Create(document model.DSSDocument) dssvalidation.SignedDocumentValidator {
	return NewASiCContainerWithCAdESValidator(document)
}

// CreateFromContent creates a SignedDocumentValidator for the given asicContent. Ports
// create(ASiCContent).
func (f *ASiCContainerWithCAdESValidatorFactory) CreateFromContent(asicContent *asic.ASiCContent) dssvalidation.SignedDocumentValidator {
	return NewASiCContainerWithCAdESValidatorFromContent(asicContent)
}

func init() {
	dssvalidation.RegisterDocumentValidatorFactory(NewASiCContainerWithCAdESValidatorFactory())
}
