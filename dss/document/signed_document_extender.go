// Ported from dss-document/src/main/java/eu/europa/esig/dss/extension/SignedDocumentExtender.java (DSS 6.5.RC1).
//
// DEVIATION: Java's static #fromDocument uses ServiceLoader<SignedDocumentExtenderFactory> to
// scan the classpath for META-INF/services registrations. Go has no classpath/ServiceLoader
// equivalent, so this is replaced by an explicit package-level registry
// (RegisterSignedDocumentExtenderFactory/FromDocument below); format packages (CAdES, XAdES, ...
// ported in later phases) call RegisterSignedDocumentExtenderFactory from an init() func to
// participate, preserving registration-order iteration the way ServiceLoader preserves
// declaration order.
package document

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/extension"
	"github.com/utain/esig/dss/spi/validation"
)

// SignedDocumentExtender contains common code for signature augmentation utilities.
type SignedDocumentExtender interface {
	extension.DocumentExtender

	// SignatureForm gets the signature form for the current implementation. Port of the
	// abstract #getSignatureForm.
	SignatureForm() enumerations.SignatureForm

	// IsASiC gets whether the document to be extended represents an ASiC container. Port of
	// #isASiC.
	IsASiC() bool
}

// SignedDocumentExtenderBase is an embeddable base implementing the common, non-abstract parts
// of SignedDocumentExtender: the CertificateVerifier/TspSource/Services fields and their
// setters, plus the default IsASiC() (false). Go has no method overriding across an embedded
// base (see PORTING.md), so unlike Java's abstract class this base does not itself satisfy
// SignedDocumentExtender - concrete extenders (CAdES, XAdES, ... ported in later phases) embed
// it and additionally supply SignatureForm() plus the DocumentExtender ExtendDocument* methods.
type SignedDocumentExtenderBase struct {
	// CertificateVerifier is the reference to the certificate verifier. The current DSS
	// implementation proposes validation.CommonCertificateVerifier. This verifier encapsulates
	// the references to different sources used in the signature validation process.
	CertificateVerifier validation.CertificateVerifier

	// TspSource is the source to be used for a timestamp token request, when applicable.
	TspSource validation.TSPSource

	// Services are (optional) document signature services. When defined, the applicable
	// instance of a corresponding service will be used. If no suitable service found, a new
	// service instance will be created. Each element is a DocumentSignatureService[SP, TP] for
	// some format-specific SP/TP - Go generics have no wildcard/existential type to express
	// Java's DocumentSignatureService<?, ?>[], so this is untyped here and type-asserted by
	// AbstractDocumentExtender's InitSignatureService.
	Services []any
}

// SetCertificateVerifier ports #setCertificateVerifier.
func (e *SignedDocumentExtenderBase) SetCertificateVerifier(certificateVerifier validation.CertificateVerifier) {
	e.CertificateVerifier = certificateVerifier
}

// SetTspSource ports #setTspSource.
func (e *SignedDocumentExtenderBase) SetTspSource(tspSource validation.TSPSource) {
	e.TspSource = tspSource
}

// SetServices (optional) sets document signature services. When defined, the applicable
// instance of a corresponding service will be used. If no suitable service found, a new service
// instance will be created. Port of #setServices.
func (e *SignedDocumentExtenderBase) SetServices(services ...any) {
	e.Services = services
}

// IsASiC ports #isASiC: the default implementation returns false.
func (e *SignedDocumentExtenderBase) IsASiC() bool {
	return false
}

// signedDocumentExtenderFactories is the package-level registry backing FromDocument; see the
// file DEVIATION above.
var signedDocumentExtenderFactories []SignedDocumentExtenderFactory

// RegisterSignedDocumentExtenderFactory registers a SignedDocumentExtenderFactory to be
// consulted by FromDocument, in registration order. Stands in for a
// META-INF/services/...SignedDocumentExtenderFactory ServiceLoader entry; see the file
// DEVIATION above.
func RegisterSignedDocumentExtenderFactory(factory SignedDocumentExtenderFactory) {
	signedDocumentExtenderFactories = append(signedDocumentExtenderFactories, factory)
}

// FromDocument guesses the document format and returns an appropriate document reader. Port of
// the static #fromDocument. Panics if dssDocument is nil (Java's
// Objects.requireNonNull(dssDocument, "DSSDocument is null")) or if no registered factory
// supports it (Java's unchecked UnsupportedOperationException).
func FromDocument(dssDocument model.DSSDocument) SignedDocumentExtender {
	if dssDocument == nil {
		panic("DSSDocument is null")
	}
	for _, factory := range signedDocumentExtenderFactories {
		if factory.IsSupported(dssDocument) {
			return factory.Create(dssDocument)
		}
	}
	panic("Document format not recognized/handled")
}
