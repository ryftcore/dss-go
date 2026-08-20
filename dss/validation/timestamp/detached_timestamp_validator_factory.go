// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/timestamp/DetachedTimestampValidatorFactory.java (DSS 6.5.RC1).
//
// Implements dss/validation.DocumentValidatorFactory, mirroring
// analyzer.DocumentAnalyzerFactory - the same ServiceLoader-registry stand-in pattern this
// porting effort uses throughout (see validation.RegisterDocumentValidatorFactory).
package timestamp

import (
	"github.com/utain/esig/dss/model"
	dssvalidation "github.com/utain/esig/dss/validation"
)

// DetachedTimestampValidatorFactory returns a validator for a detached timestamp document.
type DetachedTimestampValidatorFactory struct{}

// compile-time interface assertion.
var _ dssvalidation.DocumentValidatorFactory = (*DetachedTimestampValidatorFactory)(nil)

// NewDetachedTimestampValidatorFactory is the default constructor.
func NewDetachedTimestampValidatorFactory() *DetachedTimestampValidatorFactory {
	return &DetachedTimestampValidatorFactory{}
}

// IsSupported checks if the document is supported by the current implementation of
// DetachedTimestampValidator. Port of isSupported(DSSDocument).
func (f *DetachedTimestampValidatorFactory) IsSupported(document model.DSSDocument) bool {
	validator := newDetachedTimestampValidator(newDetachedTimestampAnalyzer())
	return validator.IsSupported(document)
}

// Create instantiates a DetachedTimestampValidator with the given document. Port of
// create(DSSDocument).
func (f *DetachedTimestampValidatorFactory) Create(document model.DSSDocument) dssvalidation.SignedDocumentValidator {
	return NewDetachedTimestampValidator(document)
}

// init registers this factory with the document validator registry, replacing upstream's
// META-INF/services/eu.europa.esig.dss.validation.DocumentValidatorFactory entry (matching the
// self-registration convention already used by, e.g., cades.CMSDocumentAnalyzerFactory for the
// sibling analyzer registry).
func init() {
	dssvalidation.RegisterDocumentValidatorFactory(NewDetachedTimestampValidatorFactory())
}
