// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CMSDocumentValidator.java (DSS 6.5.RC1).
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the Java Javadoc precedent).
package cades

import (
	"github.com/utain/esig/dss/cms"
	"github.com/utain/esig/dss/model"
	dssvalidation "github.com/utain/esig/dss/validation"
	dssdiagnostic "github.com/utain/esig/dss/validation/reports/diagnostic"
)

// CMSDocumentValidator is the validation of a CMS document. Port of the class
// CMSDocumentValidator, extending validation.SignedDocumentValidator.
type CMSDocumentValidator struct {
	dssvalidation.SignedDocumentValidatorBase

	// cmsSignedData is the CMSSignedData to be validated.
	cmsSignedData *cms.CMS
}

// compile-time interface assertions.
var (
	_ dssvalidation.SignedDocumentValidator          = (*CMSDocumentValidator)(nil)
	_ dssvalidation.SignedDocumentValidatorOverrides = (*CMSDocumentValidator)(nil)
)

// newCMSDocumentValidator is the port of the package-private default constructor.
func newCMSDocumentValidator() *CMSDocumentValidator {
	v := &CMSDocumentValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(newCMSDocumentAnalyzer()),
	}
	v.InitSignedDocumentValidator(v)
	return v
}

// NewCMSDocumentValidator creates a CMSDocumentValidator from a CMS. Port of the constructor
// CMSDocumentValidator(CMS).
func NewCMSDocumentValidator(cmsObj *cms.CMS) *CMSDocumentValidator {
	v := &CMSDocumentValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(NewCMSDocumentAnalyzer(cmsObj)),
	}
	v.InitSignedDocumentValidator(v)
	return v
}

// NewCMSDocumentValidatorFromDocument creates a CMSDocumentValidator from a DSSDocument. Port
// of the constructor CMSDocumentValidator(DSSDocument).
func NewCMSDocumentValidatorFromDocument(document model.DSSDocument) (*CMSDocumentValidator, error) {
	documentAnalyzer, err := NewCMSDocumentAnalyzerFromDocument(document)
	if err != nil {
		return nil, err
	}
	v := &CMSDocumentValidator{SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(documentAnalyzer)}
	v.InitSignedDocumentValidator(v)
	return v, nil
}

// DocumentAnalyzer returns the CMSDocumentAnalyzer of this validator. Port of the
// getDocumentAnalyzer() override, narrowing the base's return type.
func (v *CMSDocumentValidator) DocumentAnalyzer() *CMSDocumentAnalyzer {
	return v.SignedDocumentValidatorBase.DocumentAnalyzer().(*CMSDocumentAnalyzer)
}

// CmsSignedData returns a CMS. Port of getCmsSignedData().
func (v *CMSDocumentValidator) CmsSignedData() *cms.CMS {
	return v.DocumentAnalyzer().CMS()
}

// InitializeDiagnosticDataBuilder is the port of the @Override initializeDiagnosticDataBuilder().
//
// Returns the base (non-covariant) type as required by
// dssvalidation.SignedDocumentValidatorOverrides - Go has no covariant return types, so the
// concrete *CAdESDiagnosticDataBuilder wiring is instead reached through its own
// InitSignedDocumentDiagnosticDataBuilder virtual-dispatch registration (see
// cades_diagnostic_data_builder.go).
func (v *CMSDocumentValidator) InitializeDiagnosticDataBuilder() *dssdiagnostic.SignedDocumentDiagnosticDataBuilder {
	builder := NewCAdESDiagnosticDataBuilder()
	return &builder.SignedDocumentDiagnosticDataBuilder
}
