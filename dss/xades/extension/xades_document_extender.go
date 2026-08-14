// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/extension/XAdESDocumentExtender.java (DSS 6.5.RC1).
//
// Java's `extends AbstractDocumentExtender<XAdESSignatureParameters, XAdESTimestampParameters>`
// becomes embedding plus the InitAbstractDocumentExtender(self) registration documented in
// dss-document's abstract_document_extender.go: Go has no method overriding across embedding,
// so the base dispatches into the five methods below through
// document.AbstractDocumentExtenderOverrides.
package extension

import (
	"github.com/utain/esig/dss/document"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/xades"
)

// XAdESDocumentExtender is the XAdES specific implementation of a
// eu.europa.esig.dss.spi.augmentation.DocumentExtender.
type XAdESDocumentExtender struct {
	document.AbstractDocumentExtender[*xades.XAdESSignatureParameters, *xades.XAdESTimestampParameters]
}

// newXAdESDocumentExtender is the package-private empty constructor, used by
// XAdESDocumentExtenderFactory#isSupported.
func newXAdESDocumentExtender() *XAdESDocumentExtender {
	extender := &XAdESDocumentExtender{}
	extender.InitAbstractDocumentExtender(extender)
	return extender
}

// NewXAdESDocumentExtender is the default constructor, taking the document to be extended.
// Port of XAdESDocumentExtender(DSSDocument); panics with the Java message when the document is
// nil (Objects.requireNonNull).
func NewXAdESDocumentExtender(doc model.DSSDocument) *XAdESDocumentExtender {
	if doc == nil {
		panic("Document to be extended cannot be null!")
	}
	extender := newXAdESDocumentExtender()
	extender.Document = doc
	return extender
}

// CreateSignatureService ports the overridden protected createSignatureService(). Panics with
// the Java message when no CertificateVerifier was provided (Objects.requireNonNull).
func (e *XAdESDocumentExtender) CreateSignatureService() document.DocumentSignatureService[*xades.XAdESSignatureParameters, *xades.XAdESTimestampParameters] {
	if e.CertificateVerifier == nil {
		panic("Please provide CertificateVerifier or corresponding XAdESService!")
	}
	service := xades.NewXAdESService(e.CertificateVerifier)
	service.SetTspSource(e.TspSource)
	return service
}

// IsSupported ports the overridden isSupported(DSSDocument).
func (e *XAdESDocumentExtender) IsSupported(dssDocument model.DSSDocument) bool {
	return xades.NewXMLDocumentAnalyzerFactory().IsSupported(dssDocument)
}

// EmptySignatureParameters ports the overridden protected emptySignatureParameters().
func (e *XAdESDocumentExtender) EmptySignatureParameters() *xades.XAdESSignatureParameters {
	return xades.NewXAdESSignatureParameters()
}

// IsSupportedParameters ports the overridden protected
// isSupportedParameters(SerializableSignatureParameters).
func (e *XAdESDocumentExtender) IsSupportedParameters(parameters model.SerializableSignatureParameters) bool {
	_, ok := parameters.(*xades.XAdESSignatureParameters)
	return ok
}

// IsSupportedService ports the overridden protected
// isSupportedService(DocumentSignatureService<?, ?>).
func (e *XAdESDocumentExtender) IsSupportedService(service any) bool {
	_, ok := service.(*xades.XAdESService)
	return ok
}

// SignatureForm ports the overridden getSignatureForm().
func (e *XAdESDocumentExtender) SignatureForm() enumerations.SignatureForm {
	return enumerations.SignatureForm_XAdES
}

// compile-time assertion that the extender satisfies the abstract base's contract.
var _ document.AbstractDocumentExtenderOverrides[*xades.XAdESSignatureParameters, *xades.XAdESTimestampParameters] = (*XAdESDocumentExtender)(nil)
