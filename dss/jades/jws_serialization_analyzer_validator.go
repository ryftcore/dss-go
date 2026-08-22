// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JWSSerializationAnalyzerValidator.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: *JAdESSignature - see abstract_jws_document_analyzer.go's and
// jws_compact_document_analyzer.go's file headers.
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

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// JWSSerializationAnalyzerValidator validates a JWS Serialization or Flattened signature format.
// Port of the class JWSSerializationAnalyzerValidator, extending AbstractJWSDocumentAnalyzer.
type JWSSerializationAnalyzerValidator struct {
	AbstractJWSDocumentAnalyzer
}

// NewJWSSerializationAnalyzerValidator is the port of the empty constructor.
func NewJWSSerializationAnalyzerValidator() *JWSSerializationAnalyzerValidator {
	a := &JWSSerializationAnalyzerValidator{AbstractJWSDocumentAnalyzer: NewAbstractJWSDocumentAnalyzerBase()}
	a.InitAbstractJWSDocumentAnalyzer(a)
	return a
}

// NewJWSSerializationAnalyzerValidatorFromDocument is the port of the (DSSDocument) constructor.
func NewJWSSerializationAnalyzerValidatorFromDocument(document model.DSSDocument) *JWSSerializationAnalyzerValidator {
	a := NewJWSSerializationAnalyzerValidator()
	a.InitFromDocument(document)
	return a
}

// IsSupported implements AbstractJWSDocumentAnalyzerOverrides. Port of isSupported(DSSDocument).
func (a *JWSSerializationAnalyzerValidator) IsSupported(document model.DSSDocument) bool {
	jwsJsonSerializationParser := NewJWSJsonSerializationParser(document)
	supported, err := jwsJsonSerializationParser.IsSupported()
	return err == nil && supported
}

// BuildSignatures implements analyzer.DefaultDocumentAnalyzerOverrides (through
// AbstractJWSDocumentAnalyzerOverrides). Port of the protected buildSignatures() override.
func (a *JWSSerializationAnalyzerValidator) BuildSignatures() []validation.AdvancedSignature {
	var signatures []validation.AdvancedSignature
	jwsJsonSerializationObject := a.JwsJsonSerializationObject()
	foundSignatures := jwsJsonSerializationObject.Signatures()
	// Upstream logs "{} signature(s) found".
	for _, jws := range foundSignatures {
		jadesSignature := NewJAdESSignature(jws)
		jadesSignature.SetFilename(a.Document().Name())
		jadesSignature.SetSigningCertificateSource(a.SigningCertificateSource())
		jadesSignature.SetDetachedContents(a.DetachedContents())
		jadesSignature.InitBaselineRequirementsChecker(a.CertificateVerifier())
		a.ValidateSignaturePolicy(jadesSignature)
		signatures = append(signatures, jadesSignature)
	}
	return signatures
}

// BuildJwsJsonSerializationObject implements AbstractJWSDocumentAnalyzerOverrides. Port of the
// protected buildJwsJsonSerializationObject() override.
//
// Panics with an IllegalInputException when the document is not supported, and with the wrapped
// parse error otherwise, matching Java's unchecked throws.
func (a *JWSSerializationAnalyzerValidator) BuildJwsJsonSerializationObject() *JWSJsonSerializationObject {
	jwsJsonSerializationParser := NewJWSJsonSerializationParser(a.Document())
	supported, err := jwsJsonSerializationParser.IsSupported()
	if err == nil && supported {
		jwsJsonSerializationObject, err := jwsJsonSerializationParser.Parse()
		if err != nil {
			panic(model.NewDSSErrorWithCause(err))
		}
		return jwsJsonSerializationObject
	}
	panic(exception.NewIllegalInputException("The given document is not supported by JWSSerializationDocumentValidator!"))
}

// compile-time assertion: *JWSSerializationAnalyzerValidator implements AbstractJWSDocumentAnalyzerOverrides.
var _ AbstractJWSDocumentAnalyzerOverrides = (*JWSSerializationAnalyzerValidator)(nil)
