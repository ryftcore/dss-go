// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/extension/ASiCWithXAdESDocumentExtender.java (DSS 6.5.RC1).
//
// Java's `extends AbstractDocumentExtender<ASiCWithXAdESSignatureParameters,
// TimestampParameters>` becomes embedding plus the InitAbstractDocumentExtender(self)
// registration documented in dss-document's abstract_document_extender.go: Go has no method
// overriding across embedding, so the base dispatches into the five methods below through
// document.AbstractDocumentExtenderOverrides. Mirrors the
// asic/cades/extension/asic_with_cades_document_extender.go precedent.
package extension

import (
	asicxades "github.com/ryftcore/dss-go/dss/asic/xades"
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	dssxades "github.com/ryftcore/dss-go/dss/xades"
)

// ASiCWithXAdESDocumentExtender is the ASiC with XAdES container specific implementation of a
// eu.europa.esig.dss.spi.augmentation.DocumentExtender.
type ASiCWithXAdESDocumentExtender struct {
	document.AbstractDocumentExtender[*asicxades.ASiCWithXAdESSignatureParameters, *dssxades.TimestampParameters]
}

// newASiCWithXAdESDocumentExtender is the package-private empty constructor, used by
// ASiCWithXAdESDocumentExtenderFactory#isSupported.
func newASiCWithXAdESDocumentExtender() *ASiCWithXAdESDocumentExtender {
	extender := &ASiCWithXAdESDocumentExtender{}
	extender.InitAbstractDocumentExtender(extender)
	return extender
}

// NewASiCWithXAdESDocumentExtender is the default constructor, taking the document to be
// extended. Port of ASiCWithXAdESDocumentExtender(DSSDocument); panics with the Java message
// when the document is nil (Objects.requireNonNull).
func NewASiCWithXAdESDocumentExtender(doc model.DSSDocument) *ASiCWithXAdESDocumentExtender {
	if doc == nil {
		panic("Document to be extended cannot be null!")
	}
	extender := newASiCWithXAdESDocumentExtender()
	extender.Document = doc
	return extender
}

// CreateSignatureService ports the overridden protected createSignatureService(). Panics with
// the Java message when no CertificateVerifier was provided (Objects.requireNonNull).
func (e *ASiCWithXAdESDocumentExtender) CreateSignatureService() document.SignatureService[*asicxades.ASiCWithXAdESSignatureParameters, *dssxades.TimestampParameters] {
	if e.CertificateVerifier == nil {
		panic("Please provide CertificateVerifier or corresponding ASiCWithXAdESService!")
	}
	service := asicxades.NewASiCWithXAdESService(e.CertificateVerifier)
	service.SetTspSource(e.TspSource)
	return service
}

// IsSupported ports the overridden isSupported(DSSDocument).
func (e *ASiCWithXAdESDocumentExtender) IsSupported(dssDocument model.DSSDocument) bool {
	return asicxades.NewASiCWithXAdESFormatDetector().IsSupportedASiC(dssDocument)
}

// EmptySignatureParameters ports the overridden protected emptySignatureParameters().
func (e *ASiCWithXAdESDocumentExtender) EmptySignatureParameters() *asicxades.ASiCWithXAdESSignatureParameters {
	return asicxades.NewASiCWithXAdESSignatureParameters()
}

// IsSupportedParameters ports the overridden protected
// isSupportedParameters(SerializableSignatureParameters).
func (e *ASiCWithXAdESDocumentExtender) IsSupportedParameters(parameters model.SerializableSignatureParameters) bool {
	_, ok := parameters.(*asicxades.ASiCWithXAdESSignatureParameters)
	return ok
}

// IsSupportedService ports the overridden protected
// isSupportedService(SignatureService<?, ?>).
func (e *ASiCWithXAdESDocumentExtender) IsSupportedService(service any) bool {
	_, ok := service.(*asicxades.ASiCWithXAdESService)
	return ok
}

// SignatureForm ports the overridden getSignatureForm().
func (e *ASiCWithXAdESDocumentExtender) SignatureForm() enumerations.SignatureForm {
	return enumerations.SignatureFormXAdES
}

// IsASiC ports the overridden isASiC().
func (e *ASiCWithXAdESDocumentExtender) IsASiC() bool {
	return true
}

// compile-time assertion that the extender satisfies the abstract base's contract.
var _ document.AbstractDocumentExtenderOverrides[*asicxades.ASiCWithXAdESSignatureParameters, *dssxades.TimestampParameters] = (*ASiCWithXAdESDocumentExtender)(nil)
