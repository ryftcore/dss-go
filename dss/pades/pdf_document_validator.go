// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PDFDocumentValidator.java
// (DSS 6.5.RC1).
package pades

import (
	"github.com/ryftcore/dss-go/dss/model"
	dssvalidation "github.com/ryftcore/dss-go/dss/validation"
	dssdiagnostic "github.com/ryftcore/dss-go/dss/validation/reports/diagnostic"
)

// PDFDocumentValidator is the validation of a PDF document. Port of the class
// PDFDocumentValidator, extending validation.SignedDocumentValidator.
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the file header).
type PDFDocumentValidator struct {
	dssvalidation.SignedDocumentValidatorBase
}

// compile-time interface assertion.
var (
	_ dssvalidation.SignedDocumentValidator          = (*PDFDocumentValidator)(nil)
	_ dssvalidation.SignedDocumentValidatorOverrides = (*PDFDocumentValidator)(nil)
)

// newPDFDocumentValidator is the port of the protected empty constructor.
func newPDFDocumentValidator() *PDFDocumentValidator {
	v := &PDFDocumentValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(newPDFDocumentAnalyzer()),
	}
	v.InitSignedDocumentValidator(v)
	return v
}

// newPDFDocumentValidatorWithAnalyzer is the port of the protected
// PDFDocumentValidator(PDFDocumentAnalyzer) constructor.
func newPDFDocumentValidatorWithAnalyzer(pdfDocumentAnalyzer *PDFDocumentAnalyzer) *PDFDocumentValidator {
	v := &PDFDocumentValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(pdfDocumentAnalyzer),
	}
	v.InitSignedDocumentValidator(v)
	return v
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
func (v *PDFDocumentValidator) InitializeDiagnosticDataBuilder() *dssdiagnostic.SignedDocumentDiagnosticDataBuilder {
	builder := NewPAdESDiagnosticDataBuilder()
	return &builder.SignedDocumentDiagnosticDataBuilder
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
