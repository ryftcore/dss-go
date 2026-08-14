//go:build phase8

// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JWSCompactDocumentValidator.java (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (phase 6 integration): gated behind the `phase8` build tag; see
// abstract_jws_document_validator.go's file header.
package jades

import "github.com/utain/esig/dss/model"

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
	return &JWSCompactDocumentValidator{
		AbstractJWSDocumentValidator: newAbstractJWSDocumentValidator(NewJWSCompactDocumentAnalyzer()),
	}
}

// NewJWSCompactDocumentValidatorFromDocument is the port of the (DSSDocument) constructor.
func NewJWSCompactDocumentValidatorFromDocument(document model.DSSDocument) *JWSCompactDocumentValidator {
	return &JWSCompactDocumentValidator{
		AbstractJWSDocumentValidator: newAbstractJWSDocumentValidator(NewJWSCompactDocumentAnalyzerFromDocument(document)),
	}
}

// DocumentAnalyzer returns the JWSCompactDocumentAnalyzer of this validator. Port of the
// getDocumentAnalyzer() override, narrowing the base's return type.
func (v *JWSCompactDocumentValidator) DocumentAnalyzer() *JWSCompactDocumentAnalyzer {
	return v.AbstractJWSDocumentValidator.DocumentAnalyzer().(*JWSCompactDocumentAnalyzer)
}
