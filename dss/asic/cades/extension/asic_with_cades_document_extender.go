// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/extension/ASiCWithCAdESDocumentExtender.java (DSS 6.5.RC1).
//
// Java's `extends AbstractDocumentExtender<ASiCWithCAdESSignatureParameters,
// ASiCWithCAdESTimestampParameters>` becomes embedding plus the InitAbstractDocumentExtender(self)
// registration documented in dss-document's abstract_document_extender.go: Go has no method
// overriding across embedding, so the base dispatches into the five methods below through
// document.AbstractDocumentExtenderOverrides.
package extension

import (
	asiccades "github.com/utain/esig/dss/asic/cades"
	"github.com/utain/esig/dss/document"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// ASiCWithCAdESDocumentExtender is the ASiC with CAdES container specific implementation of a
// eu.europa.esig.dss.spi.augmentation.DocumentExtender.
type ASiCWithCAdESDocumentExtender struct {
	document.AbstractDocumentExtender[*asiccades.ASiCWithCAdESSignatureParameters, *asiccades.ASiCWithCAdESTimestampParameters]
}

// newASiCWithCAdESDocumentExtender is the package-private empty constructor, used by
// ASiCWithCAdESDocumentExtenderFactory#isSupported.
func newASiCWithCAdESDocumentExtender() *ASiCWithCAdESDocumentExtender {
	extender := &ASiCWithCAdESDocumentExtender{}
	extender.InitAbstractDocumentExtender(extender)
	return extender
}

// NewASiCWithCAdESDocumentExtender is the default constructor, taking the document to be
// extended. Port of ASiCWithCAdESDocumentExtender(DSSDocument); panics with the Java message
// when the document is nil (Objects.requireNonNull).
func NewASiCWithCAdESDocumentExtender(doc model.DSSDocument) *ASiCWithCAdESDocumentExtender {
	if doc == nil {
		panic("Document to be extended cannot be null!")
	}
	extender := newASiCWithCAdESDocumentExtender()
	extender.Document = doc
	return extender
}

// CreateSignatureService ports the overridden protected createSignatureService(). Panics with
// the Java message when no CertificateVerifier was provided (Objects.requireNonNull).
func (e *ASiCWithCAdESDocumentExtender) CreateSignatureService() document.DocumentSignatureService[*asiccades.ASiCWithCAdESSignatureParameters, *asiccades.ASiCWithCAdESTimestampParameters] {
	if e.CertificateVerifier == nil {
		panic("Please provide CertificateVerifier or corresponding ASiCWithCAdESService!")
	}
	service := asiccades.NewASiCWithCAdESService(e.CertificateVerifier)
	service.SetTspSource(e.TspSource)
	return service
}

// IsSupported ports the overridden isSupported(DSSDocument).
func (e *ASiCWithCAdESDocumentExtender) IsSupported(dssDocument model.DSSDocument) bool {
	return asiccades.NewASiCWithCAdESFormatDetector().IsSupportedASiC(dssDocument)
}

// EmptySignatureParameters ports the overridden protected emptySignatureParameters().
func (e *ASiCWithCAdESDocumentExtender) EmptySignatureParameters() *asiccades.ASiCWithCAdESSignatureParameters {
	return asiccades.NewASiCWithCAdESSignatureParameters()
}

// IsSupportedParameters ports the overridden protected
// isSupportedParameters(SerializableSignatureParameters).
func (e *ASiCWithCAdESDocumentExtender) IsSupportedParameters(parameters model.SerializableSignatureParameters) bool {
	_, ok := parameters.(*asiccades.ASiCWithCAdESSignatureParameters)
	return ok
}

// IsSupportedService ports the overridden protected
// isSupportedService(DocumentSignatureService<?, ?>).
func (e *ASiCWithCAdESDocumentExtender) IsSupportedService(service any) bool {
	_, ok := service.(*asiccades.ASiCWithCAdESService)
	return ok
}

// SignatureForm ports the overridden getSignatureForm().
func (e *ASiCWithCAdESDocumentExtender) SignatureForm() enumerations.SignatureForm {
	return enumerations.SignatureForm_CAdES
}

// IsASiC ports the overridden isASiC().
func (e *ASiCWithCAdESDocumentExtender) IsASiC() bool {
	return false
}

// compile-time assertion that the extender satisfies the abstract base's contract.
var _ document.AbstractDocumentExtenderOverrides[*asiccades.ASiCWithCAdESSignatureParameters, *asiccades.ASiCWithCAdESTimestampParameters] = (*ASiCWithCAdESDocumentExtender)(nil)
