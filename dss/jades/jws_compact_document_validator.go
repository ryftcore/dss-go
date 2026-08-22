// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JWSCompactDocumentValidator.java (DSS 6.5.RC1).
package jades

import "github.com/ryftcore/dss-go/dss/model"

// JWSCompactDocumentValidator validates a JWS Compact signature. Port of the class
// JWSCompactDocumentValidator, extending AbstractJWSDocumentValidator.
//
// NOTE: in order to perform the validation process, please ensure the dss/validation package is
// available within the dependencies list of your project.
type JWSCompactDocumentValidator struct {
	AbstractJWSDocumentValidator
}

// NewJWSCompactDocumentValidator is the port of the empty constructor.
func NewJWSCompactDocumentValidator() *JWSCompactDocumentValidator {
	v := &JWSCompactDocumentValidator{
		AbstractJWSDocumentValidator: newAbstractJWSDocumentValidator(NewJWSCompactDocumentAnalyzer()),
	}
	v.InitSignedDocumentValidator(v)
	return v
}

// NewJWSCompactDocumentValidatorFromDocument is the port of the (DSSDocument) constructor.
func NewJWSCompactDocumentValidatorFromDocument(document model.DSSDocument) *JWSCompactDocumentValidator {
	v := &JWSCompactDocumentValidator{
		AbstractJWSDocumentValidator: newAbstractJWSDocumentValidator(NewJWSCompactDocumentAnalyzerFromDocument(document)),
	}
	v.InitSignedDocumentValidator(v)
	return v
}

// DocumentAnalyzer returns the JWSCompactDocumentAnalyzer of this validator. Port of the
// getDocumentAnalyzer() override, narrowing the base's return type.
func (v *JWSCompactDocumentValidator) DocumentAnalyzer() *JWSCompactDocumentAnalyzer {
	return v.AbstractJWSDocumentValidator.DocumentAnalyzer().(*JWSCompactDocumentAnalyzer)
}
