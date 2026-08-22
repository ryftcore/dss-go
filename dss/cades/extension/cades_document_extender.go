// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/extension/CAdESDocumentExtender.java (DSS 6.5.RC1).
//
// Java's `extends AbstractDocumentExtender<CAdESSignatureParameters, CAdESTimestampParameters>`
// becomes embedding plus the InitAbstractDocumentExtender(self) registration documented in
// dss-document's abstract_document_extender.go: Go has no method overriding across embedding,
// so the base dispatches into the five methods below through
// document.AbstractDocumentExtenderOverrides.
package extension

import (
	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// CAdESDocumentExtender is the CAdES specific implementation of a
// eu.europa.esig.dss.spi.augmentation.DocumentExtender.
type CAdESDocumentExtender struct {
	document.AbstractDocumentExtender[*cades.SignatureParameters, *cades.TimestampParameters]
}

// newCAdESDocumentExtender is the package-private empty constructor, used by
// CAdESDocumentExtenderFactory#isSupported.
func newCAdESDocumentExtender() *CAdESDocumentExtender {
	extender := &CAdESDocumentExtender{}
	extender.InitAbstractDocumentExtender(extender)
	return extender
}

// NewCAdESDocumentExtender is the default constructor, taking the document to be extended.
// Port of CAdESDocumentExtender(DSSDocument); panics with the Java message when the document is
// nil (Objects.requireNonNull).
func NewCAdESDocumentExtender(doc model.DSSDocument) *CAdESDocumentExtender {
	if doc == nil {
		panic("Document to be extended cannot be null!")
	}
	extender := newCAdESDocumentExtender()
	extender.Document = doc
	return extender
}

// CreateSignatureService ports the overridden protected createSignatureService(). Panics with
// the Java message when no CertificateVerifier was provided (Objects.requireNonNull).
func (e *CAdESDocumentExtender) CreateSignatureService() document.SignatureService[*cades.SignatureParameters, *cades.TimestampParameters] {
	if e.CertificateVerifier == nil {
		panic("Please provide CertificateVerifier or corresponding CAdESService!")
	}
	service := cades.NewService(e.CertificateVerifier)
	service.SetTspSource(e.TspSource)
	return service
}

// IsSupported ports the overridden isSupported(DSSDocument).
func (e *CAdESDocumentExtender) IsSupported(dssDocument model.DSSDocument) bool {
	return cades.NewCMSDocumentAnalyzerFactory().IsSupported(dssDocument)
}

// EmptySignatureParameters ports the overridden protected emptySignatureParameters().
func (e *CAdESDocumentExtender) EmptySignatureParameters() *cades.SignatureParameters {
	return cades.NewSignatureParameters()
}

// IsSupportedParameters ports the overridden protected
// isSupportedParameters(SerializableSignatureParameters).
func (e *CAdESDocumentExtender) IsSupportedParameters(parameters model.SerializableSignatureParameters) bool {
	_, ok := parameters.(*cades.SignatureParameters)
	return ok
}

// IsSupportedService ports the overridden protected
// isSupportedService(SignatureService<?, ?>).
func (e *CAdESDocumentExtender) IsSupportedService(service any) bool {
	_, ok := service.(*cades.Service)
	return ok
}

// SignatureForm ports the overridden getSignatureForm().
func (e *CAdESDocumentExtender) SignatureForm() enumerations.SignatureForm {
	return enumerations.SignatureFormCAdES
}

// compile-time assertion that the extender satisfies the abstract base's contract.
var _ document.AbstractDocumentExtenderOverrides[*cades.SignatureParameters, *cades.TimestampParameters] = (*CAdESDocumentExtender)(nil)
