//go:build phase8

// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CMSDocumentValidator.java (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (Phase 3 integration): gated behind the `phase8` build tag so that
// `go build ./...` / `go vet ./...` / `go test ./...` are green for the rest of the module while
// dss/validation does not exist yet. Drop the tag once Phase 8 lands the package; the file needs
// no other change (see below).
//
// BLOCKED FORWARD DEPENDENCY (flagged per S3_BRIEF.md's "flag needs in notes" rule): this
// class's Java base, eu.europa.esig.dss.validation.SignedDocumentValidator, belongs to
// dss-validation, which PORTING_PLAN.md assigns to the not-yet-ported `validation` package
// (Phase 8: "validation engine + reports" - unstarted; see cades_diagnostic_data_builder.go for
// the sibling situation with dss-validation's diagnostic-data builder). There is no Go type to
// embed here yet, and the upstream Javadoc itself flags the same optionality ("In order to
// perform validation-process, please ensure the `dss-validation` module is loaded").
//
// The method bodies below are ported 1:1 against the package path and shape PORTING_PLAN.md's
// table implies (github.com/utain/esig/dss/validation, holding SignedDocumentValidator), so
// that this file needs no further changes once Phase 8 lands the package - only its imports need
// to resolve. Java's SignedDocumentValidator is itself an abstract class implementing the
// DocumentValidator interface (mirrored by DocumentValidatorFactory.Create's covariant return
// type in cms_document_validator_factory.go), so - following the DocumentAnalyzer/
// DefaultDocumentAnalyzer split spi/validation/analyzer already establishes for the sibling
// hierarchy - the assumed Phase 8 shape is an interface `validation.SignedDocumentValidator`
// plus an embeddable `validation.SignedDocumentValidatorBase` providing
// NewSignedDocumentValidatorBase(analyzer.DocumentAnalyzer) and DocumentAnalyzer().
package cades

import (
	"github.com/utain/esig/dss/cms"
	"github.com/utain/esig/dss/model"
	dssvalidation "github.com/utain/esig/dss/validation"
)

// CMSDocumentValidator is the validation of a CMS document. Port of the class
// CMSDocumentValidator, extending validation.SignedDocumentValidator.
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the file header).
type CMSDocumentValidator struct {
	dssvalidation.SignedDocumentValidatorBase

	// cmsSignedData is the CMSSignedData to be validated.
	cmsSignedData *cms.CMS
}

// newCMSDocumentValidator is the port of the package-private default constructor.
func newCMSDocumentValidator() *CMSDocumentValidator {
	return &CMSDocumentValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(newCMSDocumentAnalyzer()),
	}
}

// NewCMSDocumentValidator creates a CMSDocumentValidator from a CMS. Port of the constructor
// CMSDocumentValidator(CMS).
func NewCMSDocumentValidator(cmsObj *cms.CMS) *CMSDocumentValidator {
	return &CMSDocumentValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(NewCMSDocumentAnalyzer(cmsObj)),
	}
}

// NewCMSDocumentValidatorFromDocument creates a CMSDocumentValidator from a DSSDocument. Port
// of the constructor CMSDocumentValidator(DSSDocument).
func NewCMSDocumentValidatorFromDocument(document model.DSSDocument) (*CMSDocumentValidator, error) {
	documentAnalyzer, err := NewCMSDocumentAnalyzerFromDocument(document)
	if err != nil {
		return nil, err
	}
	return &CMSDocumentValidator{SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(documentAnalyzer)}, nil
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

// InitializeDiagnosticDataBuilder is the port of the initializeDiagnosticDataBuilder() override.
func (v *CMSDocumentValidator) InitializeDiagnosticDataBuilder() *CAdESDiagnosticDataBuilder {
	return NewCAdESDiagnosticDataBuilder()
}
