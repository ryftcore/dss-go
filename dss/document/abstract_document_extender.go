// Ported from dss-document/src/main/java/eu/europa/esig/dss/extension/AbstractDocumentExtender.java (DSS 6.5.RC1).
//
// This class provides an abstract implementation of spi.extension.DocumentExtender.
package document

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// AbstractDocumentExtenderTarget is the subset of AbstractSignatureParameters[TP]'s exported API
// (independent of TP) AbstractDocumentExtender needs, mirroring Java's
// `SP extends AbstractSignatureParameters<?>` bound - the same technique
// AbstractSignatureParametersBuilderTarget uses for AbstractSignatureParametersBuilder's `SP
// extends AbstractSignatureParameters` bound (see abstract_signature_parameters_builder.go).
type AbstractDocumentExtenderTarget interface {
	model.SerializableSignatureParameters

	SignatureLevel() enumerations.SignatureLevel
	SetSignatureLevel(signatureLevel enumerations.SignatureLevel)
	DetachedContents() []model.DSSDocument
	SetDetachedContents(detachedContents []model.DSSDocument)
}

// AbstractDocumentExtenderOverrides declares the operations Java's abstract
// AbstractDocumentExtender class declares abstract, or expects a subclass to override, and that
// the base implementation itself calls back into. Concrete format extenders (CAdES, XAdES, ...
// ported in later phases) satisfy this and pass themselves to InitAbstractDocumentExtender,
// mirroring the TokenBase.InitToken(self) convention documented in PORTING.md for base types
// that must call back into the concrete subclass across an embedded base - Go has no
// method overriding across embedding, so this "overrides" interface stands in for Java's virtual
// dispatch onto the abstract methods below.
type AbstractDocumentExtenderOverrides[SP AbstractDocumentExtenderTarget, TP model.SerializableTimestampParameters] interface {
	// SignatureForm gets the signature form for the current implementation; needed by
	// AbstractDocumentExtender.GetSignatureLevel. Port of the abstract #getSignatureForm
	// (SignedDocumentExtender).
	SignatureForm() enumerations.SignatureForm

	// CreateSignatureService creates a new instance of DocumentSignatureService. Port of the
	// protected abstract #createSignatureService.
	CreateSignatureService() SignatureService[SP, TP]

	// IsSupportedService verifies whether the provided document signature service is supported
	// by the current implementation. Port of the protected abstract #isSupportedService; service
	// is any element of SignedDocumentExtenderBase.Services (see that field's doc comment for
	// why it is untyped).
	IsSupportedService(service any) bool

	// EmptySignatureParameters returns a new instance of empty signature parameters, according
	// to the given format implementation. Port of the protected abstract
	// #emptySignatureParameters.
	EmptySignatureParameters() SP

	// IsSupportedParameters verifies whether the provided signature parameters are supported by
	// the current implementation. Port of the protected abstract #isSupportedParameters.
	IsSupportedParameters(parameters model.SerializableSignatureParameters) bool
}

// AbstractDocumentExtender provides an abstract implementation of
// spi/extension.DocumentExtender, generic over the SP implementation of
// AbstractDocumentExtenderTarget specifying the signature creation parameters and TP
// implementation of SerializableTimestampParameters specifying the timestamp creation
// parameters, when applicable.
type AbstractDocumentExtender[SP AbstractDocumentExtenderTarget, TP model.SerializableTimestampParameters] struct {
	SignedDocumentExtenderBase

	// Document is the document to be augmented (with the signatures).
	Document model.DSSDocument

	// overrides points back at the concrete extender; see InitAbstractDocumentExtender. Left
	// nil (bare zero value) panics on first use, matching the TokenBase.InitToken(self)
	// convention.
	overrides AbstractDocumentExtenderOverrides[SP, TP]
}

// InitAbstractDocumentExtender wires the base to the concrete extender's overridden behaviour.
// Must be called by every concrete extender's constructor before any other method, mirroring the
// TokenBase.InitToken(self) convention documented in PORTING.md.
func (e *AbstractDocumentExtender[SP, TP]) InitAbstractDocumentExtender(self AbstractDocumentExtenderOverrides[SP, TP]) {
	e.overrides = self
}

// ExtendDocument ports the SignatureProfile-only #extendDocument overload.
func (e *AbstractDocumentExtender[SP, TP]) ExtendDocument(signatureProfile enumerations.SignatureProfile) model.DSSDocument {
	return e.extendDocument(signatureProfile, nil)
}

// ExtendDocumentDetached ports the SignatureProfile+List #extendDocument overload.
func (e *AbstractDocumentExtender[SP, TP]) ExtendDocumentDetached(signatureProfile enumerations.SignatureProfile, detachedContents []model.DSSDocument) model.DSSDocument {
	return e.extendDocument(signatureProfile, detachedContents)
}

// ExtendDocumentWithParameters ports the SignatureProfile+varargs #extendDocument overload.
func (e *AbstractDocumentExtender[SP, TP]) ExtendDocumentWithParameters(signatureProfile enumerations.SignatureProfile, extensionParameters ...model.SerializableSignatureParameters) model.DSSDocument {
	return e.extendDocument(signatureProfile, nil, extensionParameters...)
}

// ExtendDocumentDetachedWithParameters ports the full SignatureProfile+List+varargs
// #extendDocument overload.
func (e *AbstractDocumentExtender[SP, TP]) ExtendDocumentDetachedWithParameters(signatureProfile enumerations.SignatureProfile, detachedContents []model.DSSDocument, extensionParameters ...model.SerializableSignatureParameters) model.DSSDocument {
	return e.extendDocument(signatureProfile, detachedContents, extensionParameters...)
}

// extendDocument is the private implementation shared by all four ExtendDocument* overloads.
func (e *AbstractDocumentExtender[SP, TP]) extendDocument(signatureProfile enumerations.SignatureProfile, detachedContents []model.DSSDocument, extensionParameters ...model.SerializableSignatureParameters) model.DSSDocument {
	if e.Document == nil {
		panic("Document is not provided to the extender")
	}
	if signatureProfile == "" {
		panic("SignatureProfile cannot be null!")
	}

	service := e.InitSignatureService()
	parameters := e.InitSignatureParameters(signatureProfile, detachedContents, extensionParameters...)
	return service.ExtendDocument(e.Document, parameters)
}

// InitSignatureService initializes a new DocumentSignatureService. Port of the protected
// #initSignatureService.
func (e *AbstractDocumentExtender[SP, TP]) InitSignatureService() SignatureService[SP, TP] {
	for _, service := range e.Services {
		if e.overrides.IsSupportedService(service) {
			if typed, ok := service.(SignatureService[SP, TP]); ok {
				return typed
			}
		}
	}
	return e.overrides.CreateSignatureService()
}

// InitSignatureParameters initializes signature parameters to be used on the signature
// augmentation. Port of the protected #initSignatureParameters.
func (e *AbstractDocumentExtender[SP, TP]) InitSignatureParameters(signatureProfile enumerations.SignatureProfile, detachedContents []model.DSSDocument, extensionParameters ...model.SerializableSignatureParameters) SP {
	signatureParameters := e.getFromProvidedParameters(extensionParameters...)
	return e.FillSignatureParameters(signatureParameters, signatureProfile, detachedContents)
}

// getFromProvidedParameters ports the private #getFromProvidedParameters.
func (e *AbstractDocumentExtender[SP, TP]) getFromProvidedParameters(explicitParameters ...model.SerializableSignatureParameters) SP {
	for _, parameters := range explicitParameters {
		if e.overrides.IsSupportedParameters(parameters) {
			if typed, ok := parameters.(SP); ok {
				return typed
			}
		}
	}
	return e.overrides.EmptySignatureParameters()
}

// FillSignatureParameters fills signatureParameters with the parameters from the
// augmentationParameters. Port of the protected #fillSignatureParameters.
//
// NOTE: Java's else-if branch here only logs a level/profile mismatch via slf4j - dropped as
// not load-bearing, per PORTING.md.
func (e *AbstractDocumentExtender[SP, TP]) FillSignatureParameters(signatureParameters SP, signatureProfile enumerations.SignatureProfile, detachedContents []model.DSSDocument) SP {
	if signatureParameters.SignatureLevel() == "" {
		signatureParameters.SetSignatureLevel(e.GetSignatureLevel(signatureProfile))
	}
	if len(signatureParameters.DetachedContents()) == 0 {
		signatureParameters.SetDetachedContents(detachedContents)
	}
	return signatureParameters
}

// GetSignatureLevel gets the target SignatureLevel for the given SignatureProfile relatively to
// the signature format. Port of the protected #getSignatureLevel. Panics if signatureProfile is
// empty (Java Objects.requireNonNull(...)) or if no SignatureLevel is found (Java's unchecked
// IllegalArgumentException).
func (e *AbstractDocumentExtender[SP, TP]) GetSignatureLevel(signatureProfile enumerations.SignatureProfile) enumerations.SignatureLevel {
	if signatureProfile == "" {
		panic("SignatureProfile cannot be null!")
	}
	signatureForm := e.overrides.SignatureForm()
	signatureLevel, err := enumerations.GetSignatureLevel(signatureForm, signatureProfile)
	if err != nil {
		panic(fmt.Sprintf("No SignatureLevel found for the given SignatureForm '%s' and SignatureProfile '%s'.", signatureForm, signatureProfile))
	}
	return signatureLevel
}
