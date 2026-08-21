// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XMLDocumentValidator.java
// (DSS 6.5.RC1).
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the Java Javadoc precedent).
//
// Java's XMLDocumentValidator overrides no other SignedDocumentValidator member besides the
// covariant getDocumentAnalyzer() - unlike dss-cades, dss-xades ships no
// XAdES-specific DiagnosticDataBuilder subclass, so this validator relies entirely on the base's
// default initializeDiagnosticDataBuilder() and never registers itself via
// InitSignedDocumentValidator.
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	dssvalidation "github.com/ryftcore/dss-go/dss/validation"
	"github.com/ryftcore/dss-go/dss/xades/definition"
)

// XMLDocumentValidator is the validator of an XML Signed document. Port of the class
// XMLDocumentValidator, extending validation.SignedDocumentValidator.
type XMLDocumentValidator struct {
	dssvalidation.SignedDocumentValidatorBase
}

// compile-time interface assertion.
var _ dssvalidation.SignedDocumentValidator = (*XMLDocumentValidator)(nil)

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
