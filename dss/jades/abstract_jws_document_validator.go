//go:build phase8

// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/AbstractJWSDocumentValidator.java (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (phase 6 integration): gated behind the `phase8` build tag so that
// `go build ./...` / `go vet ./...` / `go test ./...` are green for the rest of the module while
// dss/validation does not exist yet. Drop the tag once Phase 8 lands the package; the file needs
// no other change (see cades/cms_document_validator.go's identical, already-landed precedent for
// the assumed dss/validation shape this file relies on:
// validation.SignedDocumentValidator/SignedDocumentValidatorBase/NewSignedDocumentValidatorBase).
//
// FORWARD DEPENDENCY: JAdESDiagnosticDataBuilder is this same manifest's own
// jades_diagnostic_data_builder.go (also phase8-tagged).
package jades

import (
	"github.com/utain/esig/dss/spi/validation/analyzer"
	dssvalidation "github.com/utain/esig/dss/validation"
)

// jwsDocumentAnalyzer is the minimal surface AbstractJWSDocumentValidator needs from a JWS
// analyzer. Go type assertions require an exact concrete type match (unlike Java's upcast of a
// JWSCompactDocumentAnalyzer/JWSSerializationAnalyzerValidator instance to its
// AbstractJWSDocumentAnalyzer superclass), so DocumentAnalyzer below asserts to this interface -
// which every concrete analyzer's promoted AbstractJWSDocumentAnalyzer.JwsJsonSerializationObject
// method satisfies - rather than to the *AbstractJWSDocumentAnalyzer struct type itself.
type jwsDocumentAnalyzer interface {
	analyzer.DocumentAnalyzer

	// JwsJsonSerializationObject gets the JWSJsonSerializationObject.
	JwsJsonSerializationObject() *JWSJsonSerializationObject
}

// AbstractJWSDocumentValidator is the abstract class for a JWS signature validation. Port of the
// class AbstractJWSDocumentValidator, extending validation.SignedDocumentValidator.
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the file header).
type AbstractJWSDocumentValidator struct {
	dssvalidation.SignedDocumentValidatorBase
}

// newAbstractJWSDocumentValidator is the port of the protected (AbstractJWSDocumentAnalyzer)
// constructor.
func newAbstractJWSDocumentValidator(analyzer jwsDocumentAnalyzer) AbstractJWSDocumentValidator {
	return AbstractJWSDocumentValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(analyzer),
	}
}

// InitializeDiagnosticDataBuilder is the port of the initializeDiagnosticDataBuilder() override.
// Returns the concrete *JAdESDiagnosticDataBuilder (rather than the declared Java return type
// SignedDocumentDiagnosticDataBuilder) following the cades/cms_document_validator.go precedent's
// covariant-return convention, since Go has no virtual dispatch to reach
// JAdESDiagnosticDataBuilder.BuildDetachedXmlSignature's override through a base-typed return
// value otherwise.
func (v *AbstractJWSDocumentValidator) InitializeDiagnosticDataBuilder() *JAdESDiagnosticDataBuilder {
	return NewJAdESDiagnosticDataBuilder()
}

// DocumentAnalyzer returns the JWS analyzer of this validator. Port of the getDocumentAnalyzer()
// override, narrowing the base's return type. See jwsDocumentAnalyzer's doc comment for why the
// assertion targets that interface rather than the *AbstractJWSDocumentAnalyzer struct type.
func (v *AbstractJWSDocumentValidator) DocumentAnalyzer() jwsDocumentAnalyzer {
	return v.SignedDocumentValidatorBase.DocumentAnalyzer().(jwsDocumentAnalyzer)
}

// JwsJsonSerializationObject gets the JWSJsonSerializationObject. Port of
// getJwsJsonSerializationObject().
func (v *AbstractJWSDocumentValidator) JwsJsonSerializationObject() *JWSJsonSerializationObject {
	return v.DocumentAnalyzer().JwsJsonSerializationObject()
}
