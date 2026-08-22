// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/extension/PAdESDocumentExtender.java
// (DSS 6.5.RC1).
//
// Java's `extends AbstractDocumentExtender<PAdESSignatureParameters, PAdESTimestampParameters>`
// becomes embedding plus the InitAbstractDocumentExtender(self) registration documented in
// dss-document's abstract_document_extender.go, as in
// cades/extension/cades_document_extender.go.
//
// eu.europa.esig.dss.pades.extension flattens into the Go package pades
// (root+signature+timestamp+validation+dss+scope+timestamp+extension all flatten into one
// package), unlike CAdES/XAdES which keep a separate extension subpackage.
package pades

import (
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// PAdESDocumentExtender is the PAdES specific implementation of a
// eu.europa.esig.dss.spi.augmentation.DocumentExtender.
type PAdESDocumentExtender struct {
	document.AbstractDocumentExtender[*PAdESSignatureParameters, *PAdESTimestampParameters]
}

// newPAdESDocumentExtender is the package-private empty constructor, used by
// PAdESDocumentExtenderFactory#isSupported.
func newPAdESDocumentExtender() *PAdESDocumentExtender {
	extender := &PAdESDocumentExtender{}
	extender.InitAbstractDocumentExtender(extender)
	return extender
}

// NewPAdESDocumentExtender is the default constructor, taking the document to be extended. Port
// of PAdESDocumentExtender(DSSDocument); panics with the Java message when the document is nil
// (Objects.requireNonNull).
func NewPAdESDocumentExtender(doc model.DSSDocument) *PAdESDocumentExtender {
	if doc == nil {
		panic("Document to be extended cannot be null!")
	}
	extender := newPAdESDocumentExtender()
	extender.Document = doc
	return extender
}

// CreateSignatureService ports the overridden protected createSignatureService(). Panics with
// the Java message when no CertificateVerifier was provided (Objects.requireNonNull).
func (e *PAdESDocumentExtender) CreateSignatureService() document.DocumentSignatureService[*PAdESSignatureParameters, *PAdESTimestampParameters] {
	if e.CertificateVerifier == nil {
		panic("Please provide CertificateVerifier or corresponding PAdESService!")
	}
	service := NewPAdESService(e.CertificateVerifier)
	service.SetTspSource(e.TspSource)
	return service
}

// IsSupported ports the overridden isSupported(DSSDocument).
func (e *PAdESDocumentExtender) IsSupported(dssDocument model.DSSDocument) bool {
	return NewPDFDocumentAnalyzerFactory().IsSupported(dssDocument)
}

// EmptySignatureParameters ports the overridden protected emptySignatureParameters().
func (e *PAdESDocumentExtender) EmptySignatureParameters() *PAdESSignatureParameters {
	return NewPAdESSignatureParameters()
}

// IsSupportedParameters ports the overridden protected
// isSupportedParameters(SerializableSignatureParameters).
func (e *PAdESDocumentExtender) IsSupportedParameters(parameters model.SerializableSignatureParameters) bool {
	_, ok := parameters.(*PAdESSignatureParameters)
	return ok
}

// IsSupportedService ports the overridden protected
// isSupportedService(DocumentSignatureService<?, ?>).
func (e *PAdESDocumentExtender) IsSupportedService(service any) bool {
	_, ok := service.(*PAdESService)
	return ok
}

// SignatureForm ports the overridden getSignatureForm().
func (e *PAdESDocumentExtender) SignatureForm() enumerations.SignatureForm {
	return enumerations.SignatureFormPAdES
}

// compile-time assertion that the extender satisfies the abstract base's contract.
var _ document.AbstractDocumentExtenderOverrides[*PAdESSignatureParameters, *PAdESTimestampParameters] = (*PAdESDocumentExtender)(nil)
