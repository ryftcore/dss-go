// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/extension/JAdESDocumentExtender.java (DSS 6.5.RC1).
//
// Java's `extends AbstractDocumentExtender<JAdESSignatureParameters, JAdESTimestampParameters>`
// becomes embedding plus the InitAbstractDocumentExtender(self) registration documented in
// dss-document's abstract_document_extender.go: Go has no method overriding across embedding, so
// the base dispatches into the five methods below through
// document.AbstractDocumentExtenderOverrides.
//
// slf4j logging is dropped, per PORTING.md; the emptySignatureParameters() log statement carried
// no other behaviour.
package extension

import (
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/jades"
	"github.com/ryftcore/dss-go/dss/model"
)

// JAdESDocumentExtender is the JAdES specific implementation of a
// eu.europa.esig.dss.spi.augmentation.DocumentExtender.
type JAdESDocumentExtender struct {
	document.AbstractDocumentExtender[*jades.JAdESSignatureParameters, *jades.JAdESTimestampParameters]
}

// newJAdESDocumentExtender is the package-private empty constructor, used by
// JAdESDocumentExtenderFactory#isSupported.
func newJAdESDocumentExtender() *JAdESDocumentExtender {
	extender := &JAdESDocumentExtender{}
	extender.InitAbstractDocumentExtender(extender)
	return extender
}

// NewJAdESDocumentExtender is the default constructor, taking the document to be extended.
// Port of JAdESDocumentExtender(DSSDocument); panics with the Java message when the document is
// nil (Objects.requireNonNull).
func NewJAdESDocumentExtender(doc model.DSSDocument) *JAdESDocumentExtender {
	if doc == nil {
		panic("Document to be extended cannot be null!")
	}
	extender := newJAdESDocumentExtender()
	extender.Document = doc
	return extender
}

// CreateSignatureService ports the overridden protected createSignatureService(). Panics with
// the Java message when no CertificateVerifier was provided (Objects.requireNonNull).
func (e *JAdESDocumentExtender) CreateSignatureService() document.DocumentSignatureService[*jades.JAdESSignatureParameters, *jades.JAdESTimestampParameters] {
	if e.CertificateVerifier == nil {
		panic("Please provide CertificateVerifier or corresponding JAdESService!")
	}
	service := jades.NewJAdESService(e.CertificateVerifier)
	service.SetTspSource(e.TspSource)
	return service
}

// IsSupported ports the overridden isSupported(DSSDocument).
func (e *JAdESDocumentExtender) IsSupported(dssDocument model.DSSDocument) bool {
	return jades.NewJWSDocumentAnalyzerFactory().IsSupported(dssDocument)
}

// EmptySignatureParameters ports the overridden protected emptySignatureParameters(). Java falls
// back to JWSSerializationType.JSON_SERIALIZATION when no JAdES specific parameters were found,
// logging the fallback at INFO level; the log statement carried no other behaviour and is
// dropped per PORTING.md.
func (e *JAdESDocumentExtender) EmptySignatureParameters() *jades.JAdESSignatureParameters {
	emptyParameters := jades.NewJAdESSignatureParameters()
	emptyParameters.SetJwsSerializationType(enumerations.JWSSerializationTypeJSONSerialization)
	return emptyParameters
}

// IsSupportedParameters ports the overridden protected
// isSupportedParameters(SerializableSignatureParameters).
func (e *JAdESDocumentExtender) IsSupportedParameters(parameters model.SerializableSignatureParameters) bool {
	_, ok := parameters.(*jades.JAdESSignatureParameters)
	return ok
}

// IsSupportedService ports the overridden protected
// isSupportedService(DocumentSignatureService<?, ?>).
func (e *JAdESDocumentExtender) IsSupportedService(service any) bool {
	_, ok := service.(*jades.JAdESService)
	return ok
}

// SignatureForm ports the overridden getSignatureForm().
func (e *JAdESDocumentExtender) SignatureForm() enumerations.SignatureForm {
	return enumerations.SignatureFormJAdES
}

// compile-time assertion that the extender satisfies the abstract base's contract.
var _ document.AbstractDocumentExtenderOverrides[*jades.JAdESSignatureParameters, *jades.JAdESTimestampParameters] = (*JAdESDocumentExtender)(nil)
