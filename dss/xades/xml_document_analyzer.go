// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XMLDocumentAnalyzer.java
// (DSS 6.5.RC1).
//
// Extends spi/validation/analyzer.DefaultDocumentAnalyzer, the same
// DocumentAnalyzer/DefaultDocumentAnalyzer split cades.CMSDocumentAnalyzer already establishes
// for CAdES; this file mirrors that precedent's shape (newXMLDocumentAnalyzer +
// NewXMLDocumentAnalyzer + BuildSignatures + IsSupported + OriginalDocumentsForSignature +
// GetDefaultSignaturePolicyValidator).
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/policy"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// XMLDocumentAnalyzer is the validator of an XML Signed document. Port of the class
// XMLDocumentAnalyzer, extending analyzer.DefaultDocumentAnalyzer.
type XMLDocumentAnalyzer struct {
	analyzer.DefaultDocumentAnalyzer

	// domDocument is the document to validate.
	domDocument *XAdESDOMDocument

	// disableXSWProtection defines if the XSW protection shall be disabled (false by default).
	disableXSWProtection bool
}

// newXMLDocumentAnalyzer is the port of the package-private empty constructor.
func newXMLDocumentAnalyzer() *XMLDocumentAnalyzer {
	a := &XMLDocumentAnalyzer{DefaultDocumentAnalyzer: analyzer.NewDefaultDocumentAnalyzerBase()}
	a.InitDefaultDocumentAnalyzer(a)
	return a
}

// xmlDocumentAnalyzerInitXAdESPathsHolders ports the private static initXAdESPathsHolders().
func xmlDocumentAnalyzerInitXAdESPathsHolders() []definition.XAdESPath {
	return []definition.XAdESPath{
		&definition.XAdES111Path{},
		&definition.XAdES122Path{},
		&definition.XAdES132Path{},
	}
}

// NewXMLDocumentAnalyzer is the default constructor for XMLDocumentAnalyzer. The created
// instance is initialised with default XAdESPaths, allowing support of XAdES v1.1.1, v1.2.2 and
// v1.3.2. Port of the public XMLDocumentAnalyzer(DSSDocument) constructor.
func NewXMLDocumentAnalyzer(dssDocument model.DSSDocument) (*XMLDocumentAnalyzer, error) {
	return NewXMLDocumentAnalyzerWithPathHolders(dssDocument, xmlDocumentAnalyzerInitXAdESPathsHolders())
}

// NewXMLDocumentAnalyzerWithPathHolders allows providing a custom list of XAdES Path holders.
// Can be used to enforce signature validation of a certain XAdES version(s) only. Port of the
// public XMLDocumentAnalyzer(DSSDocument, List<XAdESPath>) constructor.
//
// Panics with the Java messages when dssDocument or xadesPathHolders is missing
// (Objects.requireNonNull).
func NewXMLDocumentAnalyzerWithPathHolders(dssDocument model.DSSDocument, xadesPathHolders []definition.XAdESPath) (*XMLDocumentAnalyzer, error) {
	if dssDocument == nil {
		panic("Document to be validated cannot be null!")
	}
	if xadesPathHolders == nil {
		panic("XAdES Path holders cannot be null!")
	}
	a := newXMLDocumentAnalyzer()
	a.SetDocument(dssDocument)
	domDocument, err := xmlDocumentAnalyzerToDomDocument(dssDocument, xadesPathHolders)
	if err != nil {
		return nil, err
	}
	a.domDocument = domDocument
	return a, nil
}

// xmlDocumentAnalyzerToDomDocument ports the private toDomDocument(DSSDocument, List<XAdESPath>).
func xmlDocumentAnalyzerToDomDocument(document model.DSSDocument, xadesPathsHolders []definition.XAdESPath) (*XAdESDOMDocument, error) {
	dom, err := xmlutils.DomUtilsBuildDOMFromDocument(document)
	if err != nil {
		return nil, exception.NewIllegalInputExceptionWithCause(fmt.Sprintf("An XML file is expected : %s", err.Error()), err)
	}
	return NewXAdESDOMDocument(dom, xadesPathsHolders), nil
}

// IsSupported checks if the document is supported by the current validator. Port of
// isSupported(DSSDocument).
func (a *XMLDocumentAnalyzer) IsSupported(dssDocument model.DSSDocument) bool {
	startsWithXmlPreamble, err := xmlutils.DomUtilsStartsWithXmlPreambleDocument(dssDocument)
	if err != nil {
		return false
	}
	return startsWithXmlPreamble && !analyzer.EvidenceRecordAnalyzerIsSupportedDocument(dssDocument)
}

// SetDisableXSWProtection: NOT RECOMMENDED - this parameter allows disabling protection against
// XML Signature wrapping attacks (XSW). It disables the research by XPath expression for defined
// Type attributes. Port of setDisableXSWProtection(boolean).
func (a *XMLDocumentAnalyzer) SetDisableXSWProtection(disableXSWProtection bool) {
	a.disableXSWProtection = disableXSWProtection
}

// BuildSignatures builds a list of signatures to be extracted from a document. Port of the
// protected buildSignatures() override.
func (a *XMLDocumentAnalyzer) BuildSignatures() []validation.AdvancedSignature {
	signatures := []validation.AdvancedSignature{}
	for _, signatureEl := range a.domDocument.SignatureNodes() {
		parent := signatureEl.Parent
		if parent != nil {
			nodeName := parent.Name.QName()
			ns := parent.Name.Space
			if nodeName == "saml2:Assertion" && DSSXMLUtilsSAMLNamespace.IsSameUri(ns) {
				// Skip signed assertions.
				continue
			}
		}

		signatureDomElement := NewXAdESDOMElement(signatureEl, a.domDocument)
		xadesSignature := NewXAdESSignatureFromDOMElement(signatureDomElement)
		if a.HasDocument() {
			xadesSignature.SetFilename(a.Document().Name())
		}
		xadesSignature.SetDetachedContents(a.DetachedContents())
		xadesSignature.SetContainerContents(a.ContainerContents())
		xadesSignature.SetSigningCertificateSource(a.SigningCertificateSource())
		xadesSignature.SetDisableXSWProtection(a.disableXSWProtection)
		xadesSignature.InitBaselineRequirementsChecker(a.CertificateVerifier())
		a.ValidateSignaturePolicy(xadesSignature)
		signatures = append(signatures, xadesSignature)
	}
	return signatures
}

// OriginalDocumentsForSignature returns the signed document(s) without their signature(s). Port
// of the getOriginalDocuments(AdvancedSignature) override.
func (a *XMLDocumentAnalyzer) OriginalDocumentsForSignature(advancedSignature validation.AdvancedSignature) []model.DSSDocument {
	xadesSignature := advancedSignature.(*XAdESSignature)
	return XAdESSignatureUtilsGetSignerDocuments(xadesSignature)
}

// XAdESPathsHolder returns the XAdESPaths. Port of the deprecated getXAdESPathsHolder().
func (a *XMLDocumentAnalyzer) XAdESPathsHolder() []definition.XAdESPath {
	return a.domDocument.XAdESPathHolders()
}

// AddXAdESPathsHolder adds a XAdESPaths. This is useful when the signature follows a particular
// schema. Port of the deprecated addXAdESPathsHolder(XAdESPath).
func (a *XMLDocumentAnalyzer) AddXAdESPathsHolder(xadesPathsHolder definition.XAdESPath) {
	a.domDocument.AddXAdESPathHolder(xadesPathsHolder)
}

// ClearQueryHolders removes all elements from the list of query holders. The list will be empty
// after this call returns. Port of the deprecated clearQueryHolders().
func (a *XMLDocumentAnalyzer) ClearQueryHolders() {
	a.domDocument.ClearXAdESPathHolders()
}

// RootElement returns the root element of the validating document. Port of the public
// getRootElement().
func (a *XMLDocumentAnalyzer) RootElement() *xmldom.Node {
	return a.domDocument.Document()
}

// GetDefaultSignaturePolicyValidator returns a signature format dependent
// SignaturePolicyValidator, based on the specification. Port of the protected
// getDefaultSignaturePolicyValidator() override.
func (a *XMLDocumentAnalyzer) GetDefaultSignaturePolicyValidator() policy.SignaturePolicyValidator {
	return NewXMLSignaturePolicyValidator()
}
