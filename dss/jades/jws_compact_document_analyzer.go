// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JWSCompactDocumentAnalyzer.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY: *JAdESSignature - see abstract_jws_document_analyzer.go's file header. This
// file additionally needs:
//
//	func NewJAdESSignature(jws *JWS) *JAdESSignature // JAdESSignature(JWS)
//
// package jades (the Java eu.europa.esig.dss.jades.validation package folds into this one Go
// package per the phase-6 package layout).
package jades

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/exception"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
)

// JWSCompactDocumentAnalyzer validates a JWS Compact signature. Port of the class
// JWSCompactDocumentAnalyzer, extending AbstractJWSDocumentAnalyzer.
type JWSCompactDocumentAnalyzer struct {
	AbstractJWSDocumentAnalyzer
}

// NewJWSCompactDocumentAnalyzer is the port of the empty constructor.
func NewJWSCompactDocumentAnalyzer() *JWSCompactDocumentAnalyzer {
	a := &JWSCompactDocumentAnalyzer{AbstractJWSDocumentAnalyzer: NewAbstractJWSDocumentAnalyzerBase()}
	a.InitAbstractJWSDocumentAnalyzer(a)
	return a
}

// NewJWSCompactDocumentAnalyzerFromDocument is the port of the (DSSDocument) constructor.
func NewJWSCompactDocumentAnalyzerFromDocument(document model.DSSDocument) *JWSCompactDocumentAnalyzer {
	a := NewJWSCompactDocumentAnalyzer()
	a.InitFromDocument(document)
	return a
}

// IsSupported implements AbstractJWSDocumentAnalyzerOverrides. Port of isSupported(DSSDocument).
func (a *JWSCompactDocumentAnalyzer) IsSupported(dssDocument model.DSSDocument) bool {
	parser := NewJWSCompactSerializationParser(dssDocument)
	supported, err := parser.IsSupported()
	return err == nil && supported
}

// BuildSignatures implements analyzer.DefaultDocumentAnalyzerOverrides (through
// AbstractJWSDocumentAnalyzerOverrides). Port of the protected buildSignatures() override.
//
// Panics with the Java DSSException message when no signature is present, matching Java's
// unchecked throw.
func (a *JWSCompactDocumentAnalyzer) BuildSignatures() []validation.AdvancedSignature {
	jwsJsonSerializationObject := a.JwsJsonSerializationObject()
	foundSignatures := jwsJsonSerializationObject.Signatures()
	if utils.IsCollectionEmpty(foundSignatures) {
		panic(model.NewDSSError("No signatures is present in the document!"))
	}
	// only one signature is supported by compact serialization
	jws := foundSignatures[0]
	jadesSignature := NewJAdESSignature(jws)
	jadesSignature.SetFilename(a.Document().Name())
	jadesSignature.SetSigningCertificateSource(a.SigningCertificateSource())
	jadesSignature.SetDetachedContents(a.DetachedContents())
	jadesSignature.InitBaselineRequirementsChecker(a.CertificateVerifier())
	a.ValidateSignaturePolicy(jadesSignature)
	return []validation.AdvancedSignature{jadesSignature}
}

// BuildJwsJsonSerializationObject implements AbstractJWSDocumentAnalyzerOverrides. Port of the
// protected buildJwsJsonSerializationObject() override.
//
// Returns an error wrapped as an IllegalInputException-carrying panic when the document is not
// supported, matching Java's unchecked throw.
func (a *JWSCompactDocumentAnalyzer) BuildJwsJsonSerializationObject() *JWSJsonSerializationObject {
	jwsCompactSerializationParser := NewJWSCompactSerializationParser(a.Document())
	supported, err := jwsCompactSerializationParser.IsSupported()
	if err == nil && supported {
		jws, err := jwsCompactSerializationParser.Parse()
		if err != nil {
			panic(model.NewDSSErrorWithCause(err))
		}
		jwsJsonSerializationObject := DSSJsonUtilsToJWSJsonSerializationObject(jws)
		jwsJsonSerializationObject.SetJWSSerializationType(enumerations.JWSSerializationType_COMPACT_SERIALIZATION)
		return jwsJsonSerializationObject
	}
	panic(exception.NewIllegalInputException("The given document is not supported by JWSCompactDocumentValidator!"))
}

// compile-time assertion: *JWSCompactDocumentAnalyzer implements AbstractJWSDocumentAnalyzerOverrides.
var _ AbstractJWSDocumentAnalyzerOverrides = (*JWSCompactDocumentAnalyzer)(nil)
