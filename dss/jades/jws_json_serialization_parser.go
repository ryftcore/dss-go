// Ported from
// dss-jades/src/main/java/eu/europa/esig/dss/jades/JWSJsonSerializationParser.java
// (DSS 6.5.RC1).
package jades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/jades/specs"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JWSJsonSerializationParser creates a JWSJsonSerializationObject from a document. Port of the
// class JWSJsonSerializationParser.
type JWSJsonSerializationParser struct {
	// document is the document to be parsed. Port of the private final field of the same name.
	document model.DSSDocument
}

// NewJWSJsonSerializationParser is the default constructor for a parser that extracts a list of
// signatures and the payload. Port of JWSJsonSerializationParser(DSSDocument).
func NewJWSJsonSerializationParser(document model.DSSDocument) *JWSJsonSerializationParser {
	return &JWSJsonSerializationParser{document: document}
}

// Parse parses the provided document and returns the JWSJsonSerializationObject, if applicable.
// Port of parse().
//
// Which serialization type a document gets is decided solely by the presence of a "signatures"
// member: with it the document is a complete JWS JSON Serialization, without it a flattened one.
// That is upstream's rule and it is deliberately not tightened, because the type it records
// travels on into JWS.getJwsSerializationType and back out of the generator on an extension.
func (p *JWSJsonSerializationParser) Parse() (*JWSJsonSerializationObject, error) {
	binaries, err := spi.DSSUtilsToByteArrayOfDocument(p.document)
	if err != nil {
		return nil, exception.NewIllegalInputExceptionWithCause(fmt.Sprintf(
			"Unable to parse document with name '%s'. Reason : %s", p.document.Name(), err.Error()), err)
	}
	jsonDocument := string(binaries)

	jwsJsonSerializationObject := NewJWSJsonSerializationObject()

	structureValidationErrors := p.validateJWSStructure(jsonDocument)
	if utils.IsCollectionNotEmpty(structureValidationErrors) {
		jwsJsonSerializationObject.SetStructuralValidationErrors(structureValidationErrors)
	}

	rootStructure, err := jose.ParseJSON(jsonDocument)
	if err != nil {
		return nil, exception.NewIllegalInputExceptionWithCause(fmt.Sprintf(
			"Unable to parse document with name '%s'. Reason : %s", p.document.Name(), err.Error()), err)
	}

	if payload, ok := rootStructure.Value(JWSConstantsPayload).(string); ok {
		jwsJsonSerializationObject.SetPayload(payload)
	}

	// try to extract complete JWS JSON Serialization signatures
	signaturesObject := rootStructure.Value(JWSConstantsSignatures)
	if signaturesObject != nil {
		jwsJsonSerializationObject.SetJWSSerializationType(enumerations.JWSSerializationTypeJSONSerialization)
		if err := p.extractSignatures(jwsJsonSerializationObject, signaturesObject); err != nil {
			return nil, err
		}
	} else {
		// otherwise extract flattened JWS JSON Serialization signature
		jwsJsonSerializationObject.SetJWSSerializationType(enumerations.JWSSerializationTypeFlattenedJSONSerialization)
		if err := p.extractAndAddJWSSignature(jwsJsonSerializationObject, rootStructure); err != nil {
			return nil, err
		}
	}

	return jwsJsonSerializationObject, nil
}

// IsSupported verifies whether the given document is supported by the parser. Port of
// isSupported().
func (p *JWSJsonSerializationParser) IsSupported() (bool, error) {
	return DSSJsonUtilsIsJsonDocument(p.document)
}

// extractSignatures walks the "signatures" array. Port of the private
// extractSignatures(JWSJsonSerializationObject, Object).
func (p *JWSJsonSerializationParser) extractSignatures(
	jwsJsonSerializationObject *JWSJsonSerializationObject, signaturesObject any) error {
	signaturesObjectList, ok := signaturesObject.([]any)
	if !ok {
		return nil
	}
	for _, signatureObject := range signaturesObjectList {
		signatureMap, ok := signatureObject.(*jose.Object)
		if !ok {
			continue
		}
		if err := p.extractAndAddJWSSignature(jwsJsonSerializationObject, signatureMap); err != nil {
			return err
		}
	}
	return nil
}

// extractAndAddJWSSignature builds one JWS from a signature object and appends it. Port of the
// private extractAndAddJWSSignature(JWSJsonSerializationObject, Map).
//
// The order of the steps is load-bearing and is upstream's: the protected header is installed
// before the payload, because it carries 'b64' and therefore decides whether the shared payload
// is taken as raw UTF-8 bytes or base64url-decoded first.
func (p *JWSJsonSerializationParser) extractAndAddJWSSignature(
	jwsJsonSerializationObject *JWSJsonSerializationObject, signatureMap *jose.Object) error {

	fail := func(err error) error {
		return exception.NewIllegalInputExceptionWithCause(
			fmt.Sprintf("Unable to parse a JWS signature. Reason : [%s]", err.Error()), err)
	}

	signature := NewJWS()

	signatureObject := signatureMap.Value(JWSConstantsSignature)
	if signatureObject == nil {
		return nil
	}
	if signatureBase64Url, ok := signatureObject.(string); ok {
		if utils.IsStringBlank(signatureBase64Url) {
			return nil
		}
		signature.SetSignature(DSSJsonUtilsFromBase64Url(signatureBase64Url))
	}

	if protectedBase64Url, ok := signatureMap.Value(JWSConstantsProtected).(string); ok {
		if err := signature.SetProtected(protectedBase64Url); err != nil {
			return fail(err)
		}
	}

	if header, ok := signatureMap.Value(JWSConstantsHeader).(*jose.Object); ok {
		signature.SetUnprotected(header)
	}

	if signature.IsRfc7797UnencodedPayload() {
		signature.SetPayloadBytes([]byte(jwsJsonSerializationObject.Payload()))
	} else {
		signature.SetPayloadBytes(DSSJsonUtilsFromBase64Url(jwsJsonSerializationObject.Payload()))
	}

	signature.SetJwsJsonSerializationObject(jwsJsonSerializationObject)
	jwsJsonSerializationObject.AddSignature(signature)
	return nil
}

// validateJWSStructure validates the document against the JWS schema. Port of the private
// validateJWSStructure(String).
func (p *JWSJsonSerializationParser) validateJWSStructure(jsonDocument string) []string {
	return specs.JAdESUtilsInstance().ValidateAgainstSchema(jsonDocument)
}
