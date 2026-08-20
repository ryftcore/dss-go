// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JWSSerializationDocumentValidator.java (DSS 6.5.RC1).
//
// This class performs validation of a JWS Serialization or Flattened signature format.
//
// {
//
// "payload":"payload contents",
//
// "signatures":[
//
// {"protected":"integrity-protected header 1 contents",
// "header":non-integrity-protected header 1 contents,
// "signature":"signature 1 contents"},
//
// ...
//
// {"protected":"integrity-protected header N contents",
// "header":non-integrity-protected header N contents,
// "signature":"signature N contents"}
//
// ]
//
// }
package jades

import "github.com/utain/esig/dss/model"

// JWSSerializationDocumentValidator validates a JWS Serialization or Flattened signature. Port
// of the class JWSSerializationDocumentValidator, extending AbstractJWSDocumentValidator.
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project.
type JWSSerializationDocumentValidator struct {
	AbstractJWSDocumentValidator
}

// NewJWSSerializationDocumentValidator is the port of the empty constructor.
func NewJWSSerializationDocumentValidator() *JWSSerializationDocumentValidator {
	v := &JWSSerializationDocumentValidator{
		AbstractJWSDocumentValidator: newAbstractJWSDocumentValidator(NewJWSSerializationAnalyzerValidator()),
	}
	v.InitSignedDocumentValidator(v)
	return v
}

// NewJWSSerializationDocumentValidatorFromDocument is the port of the (DSSDocument) constructor.
func NewJWSSerializationDocumentValidatorFromDocument(document model.DSSDocument) *JWSSerializationDocumentValidator {
	v := &JWSSerializationDocumentValidator{
		AbstractJWSDocumentValidator: newAbstractJWSDocumentValidator(NewJWSSerializationAnalyzerValidatorFromDocument(document)),
	}
	v.InitSignedDocumentValidator(v)
	return v
}

// DocumentAnalyzer returns the JWSSerializationAnalyzerValidator of this validator. Port of the
// getDocumentAnalyzer() override, narrowing the base's return type.
func (v *JWSSerializationDocumentValidator) DocumentAnalyzer() *JWSSerializationAnalyzerValidator {
	return v.AbstractJWSDocumentValidator.DocumentAnalyzer().(*JWSSerializationAnalyzerValidator)
}
