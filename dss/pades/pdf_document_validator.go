//go:build phase8

// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PDFDocumentValidator.java
// (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (Phase 3 integration precedent, cades/cms_document_validator.go): gated behind
// the `phase8` build tag so that `go build ./...` / `go vet ./...` / `go test ./...` stay green
// for the rest of the module while dss/validation does not exist yet. Drop the tag once Phase 8
// lands the package; the file needs no other change.
//
// BLOCKED FORWARD DEPENDENCY: this class's Java base, eu.europa.esig.dss.validation.
// SignedDocumentValidator, belongs to dss-validation (Phase 8, unstarted). Following
// cades/cms_document_validator.go's precedent, the assumed Phase 8 shape is an interface
// `validation.SignedDocumentValidator` plus an embeddable `validation.SignedDocumentValidatorBase`
// providing NewSignedDocumentValidatorBase(analyzer.DocumentAnalyzer) and DocumentAnalyzer().
//
// FORWARD DEPENDENCY (not in this chunk's manifest): PAdESDiagnosticDataBuilder
// (eu.europa.esig.dss.pades.validation.PAdESDiagnosticDataBuilder), the PAdES-specific
// diagnostic-data builder, following the CAdESDiagnosticDataBuilder precedent in
// cades_diagnostic_data_builder.go (also phase8-gated, same sibling situation).
package pades

import (
	"github.com/utain/esig/dss/model"
	dssvalidation "github.com/utain/esig/dss/validation"
)

// PDFDocumentValidator is the validation of a PDF document. Port of the class
// PDFDocumentValidator, extending validation.SignedDocumentValidator.
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the file header).
type PDFDocumentValidator struct {
	dssvalidation.SignedDocumentValidatorBase
}

// newPDFDocumentValidator is the port of the protected empty constructor.
func newPDFDocumentValidator() *PDFDocumentValidator {
	return &PDFDocumentValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(newPDFDocumentAnalyzer()),
	}
}

// newPDFDocumentValidatorWithAnalyzer is the port of the protected
// PDFDocumentValidator(PDFDocumentAnalyzer) constructor.
func newPDFDocumentValidatorWithAnalyzer(pdfDocumentAnalyzer *PDFDocumentAnalyzer) *PDFDocumentValidator {
	return &PDFDocumentValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(pdfDocumentAnalyzer),
	}
}

// NewPDFDocumentValidator creates a PDFDocumentValidator from a DSSDocument. Port of the
// default constructor PDFDocumentValidator(DSSDocument).
func NewPDFDocumentValidator(document model.DSSDocument) *PDFDocumentValidator {
	return newPDFDocumentValidatorWithAnalyzer(NewPDFDocumentAnalyzer(document))
}

// DocumentAnalyzer returns the PDFDocumentAnalyzer of this validator. Port of the
// getDocumentAnalyzer() override, narrowing the base's return type.
func (v *PDFDocumentValidator) DocumentAnalyzer() *PDFDocumentAnalyzer {
	return v.SignedDocumentValidatorBase.DocumentAnalyzer().(*PDFDocumentAnalyzer)
}

// SetPdfObjFactory sets the IPdfObjFactory. Allow to set the used implementation. Cannot be nil.
// Port of setPdfObjFactory(IPdfObjFactory).
func (v *PDFDocumentValidator) SetPdfObjFactory(pdfObjFactory IPdfObjFactory) {
	v.DocumentAnalyzer().SetPdfObjFactory(pdfObjFactory)
}

// SetPasswordProtection specifies the used password for the encrypted document.
// Port of setPasswordProtection(char[]).
func (v *PDFDocumentValidator) SetPasswordProtection(passwordProtection []byte) {
	v.DocumentAnalyzer().SetPasswordProtection(passwordProtection)
}

// InitializeDiagnosticDataBuilder is the port of the initializeDiagnosticDataBuilder() override.
func (v *PDFDocumentValidator) InitializeDiagnosticDataBuilder() *PAdESDiagnosticDataBuilder {
	return NewPAdESDiagnosticDataBuilder()
}

// DssDictionaries returns a list of found DSS Dictionaries across different revisions.
// Port of getDssDictionaries().
func (v *PDFDocumentValidator) DssDictionaries() []PdfDssDict {
	return v.DocumentAnalyzer().DssDictionaries()
}

// Revisions gets the list of PDF document revisions. Port of getRevisions().
func (v *PDFDocumentValidator) Revisions() []PdfRevision {
	return v.DocumentAnalyzer().Revisions()
}
