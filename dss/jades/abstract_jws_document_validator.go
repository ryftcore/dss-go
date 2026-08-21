// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/AbstractJWSDocumentValidator.java (DSS 6.5.RC1).
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project (see the Java Javadoc precedent).
//
// Registration deviation: Java's initializeDiagnosticDataBuilder() override lives on
// AbstractJWSDocumentValidator itself, one level above the two concrete leaf validators. A Go
// value embedded in a not-yet-final struct cannot safely register its own address with
// InitSignedDocumentValidator (the pointer would target the temporary, not the copy that ends up
// embedded in the leaf) - see PORTING.md's Init<Base> registration lesson - so registration is
// deferred to the two leaf constructors (jws_compact_document_validator.go,
// jws_serialization_document_validator.go), each passing its own outer *JWSCompactDocumentValidator /
// *JWSSerializationDocumentValidator (which promotes AbstractJWSDocumentValidator's
// InitializeDiagnosticDataBuilder) as the SignedDocumentValidatorOverrides.
package jades

import (
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	dssvalidation "github.com/ryftcore/dss-go/dss/validation"
	dssdiagnostic "github.com/ryftcore/dss-go/dss/validation/reports/diagnostic"
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
// constructor. Does not call InitSignedDocumentValidator - see the file header.
func newAbstractJWSDocumentValidator(analyzer jwsDocumentAnalyzer) AbstractJWSDocumentValidator {
	return AbstractJWSDocumentValidator{
		SignedDocumentValidatorBase: dssvalidation.NewSignedDocumentValidatorBase(analyzer),
	}
}

// InitializeDiagnosticDataBuilder is the port of the initializeDiagnosticDataBuilder() override.
func (v *AbstractJWSDocumentValidator) InitializeDiagnosticDataBuilder() *dssdiagnostic.SignedDocumentDiagnosticDataBuilder {
	builder := NewJAdESDiagnosticDataBuilder()
	return &builder.SignedDocumentDiagnosticDataBuilder
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
