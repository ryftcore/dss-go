//go:build phase8

// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XMLDocumentValidator.java
// (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (Phase 4d integration): gated behind the `phase8` build tag so that
// `go build ./...` / `go vet ./...` / `go test ./...` are green for the rest of the module while
// dss/validation does not exist yet. Drop the tag once Phase 8 lands the package; the file needs
// no other change (see below). Mirrors the cades.CMSDocumentValidator precedent
// (cades/cms_document_validator.go) exactly, including its assumed dss/validation shape.
//
// BLOCKED FORWARD DEPENDENCY (flagged per S4D_BRIEF.md's "flag needs in notes" rule): this
// class's Java base, eu.europa.esig.dss.validation.SignedDocumentValidator, belongs to
// dss-validation, which PORTING_PLAN.md assigns to the not-yet-ported `validation` package
// (Phase 8: "validation engine + reports" - unstarted). There is no Go type to embed here yet,
// and the upstream Javadoc itself flags the same optionality ("In order to perform
// validation-process, please ensure the `dss-validation` module is loaded").
//
// The method bodies below are ported 1:1 against the package path and shape
// cades/cms_document_validator.go already established (github.com/utain/esig/dss/validation,
// holding SignedDocumentValidator/SignedDocumentValidatorBase), so that this file needs no
// further changes once Phase 8 lands the package - only its imports need to resolve.
package xades

import (
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/model"
	dssvalidation "github.com/utain/esig/dss/validation"
	"github.com/utain/esig/dss/xades/definition"
)

// XMLDocumentValidator is the validator of an XML Signed document.
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the file header). Port of the
// class XMLDocumentValidator, extending validation.SignedDocumentValidator.
type XMLDocumentValidator struct {
	dssvalidation.SignedDocumentValidatorBase
}

// newXMLDocumentValidator is the port of the package-private default constructor.
func newXMLDocumentValidator() *XMLDocumentValidator {
	return &XMLDocumentValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(newXMLDocumentAnalyzer()),
	}
}

// NewXMLDocumentValidator is the default constructor for XMLDocumentValidator. The created
// instance is initialised with default XAdESPaths, allowing support of XAdES v1.1.1, v1.2.2 and
// v1.3.2. Port of the public XMLDocumentValidator(DSSDocument) constructor.
func NewXMLDocumentValidator(dssDocument model.DSSDocument) (*XMLDocumentValidator, error) {
	documentAnalyzer, err := NewXMLDocumentAnalyzer(dssDocument)
	if err != nil {
		return nil, err
	}
	return &XMLDocumentValidator{SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(documentAnalyzer)}, nil
}

// NewXMLDocumentValidatorWithPathHolders allows providing a custom list of XAdES Path holders.
// Can be used to enforce signature validation of a certain XAdES version(s) only. Port of the
// public XMLDocumentValidator(DSSDocument, List<XAdESPath>) constructor.
func NewXMLDocumentValidatorWithPathHolders(dssDocument model.DSSDocument, xadesPathHolders []definition.XAdESPath) (*XMLDocumentValidator, error) {
	documentAnalyzer, err := NewXMLDocumentAnalyzerWithPathHolders(dssDocument, xadesPathHolders)
	if err != nil {
		return nil, err
	}
	return &XMLDocumentValidator{SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(documentAnalyzer)}, nil
}

// DocumentAnalyzer returns the XMLDocumentAnalyzer of this validator. Port of the
// getDocumentAnalyzer() override, narrowing the base's return type.
func (v *XMLDocumentValidator) DocumentAnalyzer() *XMLDocumentAnalyzer {
	return v.SignedDocumentValidatorBase.DocumentAnalyzer().(*XMLDocumentAnalyzer)
}

// SetDisableXSWProtection: NOT RECOMMENDED - this parameter allows disabling protection against
// XML Signature wrapping attacks (XSW). It disables the research by XPath expression for defined
// Type attributes. Port of setDisableXSWProtection(boolean).
func (v *XMLDocumentValidator) SetDisableXSWProtection(disableXSWProtection bool) {
	v.DocumentAnalyzer().SetDisableXSWProtection(disableXSWProtection)
}

// XAdESPathsHolder returns the XAdESPaths. Port of the deprecated getXAdESPathsHolder().
func (v *XMLDocumentValidator) XAdESPathsHolder() []definition.XAdESPath {
	return v.DocumentAnalyzer().XAdESPathsHolder()
}

// AddXAdESPathsHolder adds a XAdESPaths. This is useful when the signature follows a particular
// schema. Port of the deprecated addXAdESPathsHolder(XAdESPath).
func (v *XMLDocumentValidator) AddXAdESPathsHolder(xadesPathsHolder definition.XAdESPath) {
	v.DocumentAnalyzer().AddXAdESPathsHolder(xadesPathsHolder)
}

// ClearQueryHolders removes all elements from the list of query holders. The list will be empty
// after this call returns. Port of the deprecated clearQueryHolders().
func (v *XMLDocumentValidator) ClearQueryHolders() {
	v.DocumentAnalyzer().ClearQueryHolders()
}

// RootElement returns the root element of the validating document. Port of the public
// getRootElement().
func (v *XMLDocumentValidator) RootElement() *xmldom.Node {
	return v.DocumentAnalyzer().RootElement()
}
