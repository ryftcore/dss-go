// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/AbstractJWSDocumentAnalyzer.java (DSS 6.5.RC1).
//
// package jades (the Java eu.europa.esig.dss.jades.validation package is flattened into this one
// Go package).
package jades

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/policy"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
)

// AbstractJWSDocumentAnalyzerOverrides declares the operations AbstractJWSDocumentAnalyzer calls
// back into virtually, dispatched the way analyzer.DefaultDocumentAnalyzer dispatches to
// analyzer.DefaultDocumentAnalyzerOverrides: it embeds that interface (every method a concrete
// JWS analyzer must supply) and adds the abstract buildJwsJsonSerializationObject() this base
// additionally declares.
type AbstractJWSDocumentAnalyzerOverrides interface {
	analyzer.DefaultDocumentAnalyzerOverrides

	// BuildJwsJsonSerializationObject builds a JWSJsonSerializationObject. Port of the abstract
	// protected buildJwsJsonSerializationObject().
	BuildJwsJsonSerializationObject() *JWSJsonSerializationObject
}

// AbstractJWSDocumentAnalyzer is the abstract class for a JWS signature validation. Port of the
// class AbstractJWSDocumentAnalyzer, extending analyzer.DefaultDocumentAnalyzer.
type AbstractJWSDocumentAnalyzer struct {
	analyzer.DefaultDocumentAnalyzer

	// overrides points back at the concrete analyzer; see InitAbstractJWSDocumentAnalyzer.
	overrides AbstractJWSDocumentAnalyzerOverrides

	// jwsJsonSerializationObject is a cached copy of JWS Json Serialization object.
	jwsJsonSerializationObject *JWSJsonSerializationObject
}

// NewAbstractJWSDocumentAnalyzerBase builds the base state a subclass embeds, for the empty
// (no-document) constructor case. Port of the protected empty AbstractJWSDocumentAnalyzer()
// constructor; the subclass constructor must follow it with InitAbstractJWSDocumentAnalyzer.
func NewAbstractJWSDocumentAnalyzerBase() AbstractJWSDocumentAnalyzer {
	return AbstractJWSDocumentAnalyzer{DefaultDocumentAnalyzer: analyzer.NewDefaultDocumentAnalyzerBase()}
}

// InitAbstractJWSDocumentAnalyzer registers the concrete analyzer (JWSCompactDocumentAnalyzer or
// JWSSerializationAnalyzerValidator) with its base, including the embedded
// analyzer.DefaultDocumentAnalyzer. Every concrete subclass constructor must call this once, in
// place of calling InitDefaultDocumentAnalyzer directly.
func (a *AbstractJWSDocumentAnalyzer) InitAbstractJWSDocumentAnalyzer(overrides AbstractJWSDocumentAnalyzerOverrides) {
	a.InitDefaultDocumentAnalyzer(overrides)
	a.overrides = overrides
}

// InitFromDocument finishes construction for the (DSSDocument) constructor case. Port of the
// protected AbstractJWSDocumentAnalyzer(DSSDocument) constructor's body
// (`this.document = document; this.jwsJsonSerializationObject = buildJwsJsonSerializationObject();`).
// The concrete subclass's (DSSDocument) constructor must call InitAbstractJWSDocumentAnalyzer
// first, then this method - buildJwsJsonSerializationObject() is only dispatchable once overrides
// is set, unlike Java where `this` is already the concrete subtype throughout the constructor
// chain.
//
// Panics when document is nil (Objects.requireNonNull("Document to be validated cannot be
// null!")).
func (a *AbstractJWSDocumentAnalyzer) InitFromDocument(document model.DSSDocument) {
	if document == nil {
		panic("Document to be validated cannot be null!")
	}
	a.SetDocument(document)
	a.jwsJsonSerializationObject = a.overrides.BuildJwsJsonSerializationObject()
}

// OriginalDocumentsForSignature implements analyzer.DefaultDocumentAnalyzerOverrides. Port of
// getOriginalDocuments(AdvancedSignature).
//
// The DSSException Java catches (logging "Cannot retrieve a list of original documents") is
// swallowed the same way here, since slf4j logging is dropped per PORTING.md.
func (a *AbstractJWSDocumentAnalyzer) OriginalDocumentsForSignature(advancedSignature validation.AdvancedSignature) []model.DSSDocument {
	jadesSignature := advancedSignature.(*JAdESSignature)
	documents, err := jadesSignature.OriginalDocuments()
	if err != nil {
		return []model.DSSDocument{}
	}
	return documents
}

// JwsJsonSerializationObject gets the JWSJsonSerializationObject. Port of
// getJwsJsonSerializationObject().
func (a *AbstractJWSDocumentAnalyzer) JwsJsonSerializationObject() *JWSJsonSerializationObject {
	return a.jwsJsonSerializationObject
}

// GetDefaultSignaturePolicyValidator implements analyzer.DefaultDocumentAnalyzerOverrides. Port
// of the protected getDefaultSignaturePolicyValidator() override.
func (a *AbstractJWSDocumentAnalyzer) GetDefaultSignaturePolicyValidator() policy.SignaturePolicyValidator {
	return policy.NewNonASN1SignaturePolicyValidator()
}

// abstractJWSDocumentAnalyzer is implemented by both concrete JWS analyzers
// (*JWSCompactDocumentAnalyzer and *JWSSerializationAnalyzerValidator, each embedding
// AbstractJWSDocumentAnalyzer by value). analyzer.DocumentAnalyzerFactory.Create - the frozen
// spi/validation/analyzer interface JWSDocumentAnalyzerFactory.Create implements - can only
// return the frozen analyzer.DocumentAnalyzer interface, which does not expose
// JwsJsonSerializationObject/SignatureByID/SignaturePolicyValidatorLoader. Call sites that need
// that JAdES-specific surface recover it with jwsDocumentAnalyzerBase below, rather than
// asserting to either concrete type individually.
type abstractJWSDocumentAnalyzer interface {
	analyzer.DocumentAnalyzer
	abstractJWSDocumentAnalyzerBase() *AbstractJWSDocumentAnalyzer
}

// abstractJWSDocumentAnalyzerBase implements abstractJWSDocumentAnalyzer for
// *JWSCompactDocumentAnalyzer (see jws_compact_document_analyzer.go).
func (a *JWSCompactDocumentAnalyzer) abstractJWSDocumentAnalyzerBase() *AbstractJWSDocumentAnalyzer {
	return &a.AbstractJWSDocumentAnalyzer
}

// abstractJWSDocumentAnalyzerBase implements abstractJWSDocumentAnalyzer for
// *JWSSerializationAnalyzerValidator (see jws_serialization_analyzer_validator.go).
func (a *JWSSerializationAnalyzerValidator) abstractJWSDocumentAnalyzerBase() *AbstractJWSDocumentAnalyzer {
	return &a.AbstractJWSDocumentAnalyzer
}

// jwsDocumentAnalyzerBase recovers the *AbstractJWSDocumentAnalyzer embedded in whichever
// concrete analyzer JWSDocumentAnalyzerFactory.Create produced. Panics (via the failed type
// assertion) if given a DocumentAnalyzer that isn't one of the two JWS analyzers, which would be
// a programming error at every call site - they all obtain their argument from
// NewJWSDocumentAnalyzerFactory().Create.
func jwsDocumentAnalyzerBase(a analyzer.DocumentAnalyzer) *AbstractJWSDocumentAnalyzer {
	return a.(abstractJWSDocumentAnalyzer).abstractJWSDocumentAnalyzerBase()
}
