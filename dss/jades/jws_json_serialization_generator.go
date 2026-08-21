// Ported from
// dss-jades/src/main/java/eu/europa/esig/dss/jades/JWSJsonSerializationGenerator.java
// (DSS 6.5.RC1).
package jades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JWSJsonSerializationGenerator creates a JWS Serialization signature. Port of the class
// JWSJsonSerializationGenerator.
type JWSJsonSerializationGenerator struct {
	// jwsJsonSerializationObject is the container for JWS signature elements. Port of the
	// private final field of the same name.
	jwsJsonSerializationObject *JWSJsonSerializationObject

	// output is the target signature's format. Port of the private final field of the same name.
	output enumerations.JWSSerializationType
}

// NewJWSJsonSerializationGenerator is the default constructor. Port of
// JWSJsonSerializationGenerator(JWSJsonSerializationObject, JWSSerializationType).
func NewJWSJsonSerializationGenerator(jwsJsonSerializationObject *JWSJsonSerializationObject,
	output enumerations.JWSSerializationType) *JWSJsonSerializationGenerator {
	return &JWSJsonSerializationGenerator{
		jwsJsonSerializationObject: jwsJsonSerializationObject,
		output:                     output,
	}
}

// Generate produces the JWS signature document. Port of generate().
//
// Every map below is insertion-ordered on purpose: the bytes this returns are the signature
// document, and while the members here are not themselves signed, an archive timestamp over the
// document is computed on exactly these octets.
func (g *JWSJsonSerializationGenerator) Generate() (model.DSSDocument, error) {
	var jsonSerialization *JsonObject
	var err error

	switch g.output {
	case enumerations.JWSSerializationType_JSON_SERIALIZATION:
		jsonSerialization = g.buildJWSJsonSerialization()
	case enumerations.JWSSerializationType_FLATTENED_JSON_SERIALIZATION:
		jsonSerialization, err = g.buildFlattenedJwsJsonSerialization()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(
			"The JWSJsonSerializationGenerator does not support the given JWS Serialization Type '%s'", g.output)
	}

	binaries := []byte(jsonSerialization.ToJSONString())
	return model.NewInMemoryDocument(binaries), nil
}

// buildJWSJsonSerialization builds the complete JWS JSON Serialization form. Port of the private
// buildJWSJsonSerialization().
func (g *JWSJsonSerializationGenerator) buildJWSJsonSerialization() *JsonObject {
	// if the source is flattened: "A flattened signature will be transformed to a Complete JWS
	// JSON Serialization Format!"

	jsonSerializationMap := jose.NewObject()

	payload := g.jwsJsonSerializationObject.Payload()
	if utils.IsStringNotBlank(payload) {
		jsonSerializationMap.Put(JWSConstantsPayload, payload)
	}

	signatureList := make([]any, 0, len(g.jwsJsonSerializationObject.Signatures()))
	for _, signature := range g.jwsJsonSerializationObject.Signatures() {
		signatureMap := g.signatureJSONMap(signature)
		signatureList = append(signatureList, NewJsonObjectFromMap(signatureMap))
	}
	jsonSerializationMap.Put(JWSConstantsSignatures, signatureList)

	return NewJsonObjectFromMap(jsonSerializationMap)
}

// buildFlattenedJwsJsonSerialization builds the flattened JWS JSON Serialization form. Port of
// the private buildFlattenedJwsJsonSerialization().
func (g *JWSJsonSerializationGenerator) buildFlattenedJwsJsonSerialization() (*JsonObject, error) {
	flattenedJwsMap := jose.NewObject()
	payload := g.jwsJsonSerializationObject.Payload()
	if utils.IsStringNotBlank(payload) {
		flattenedJwsMap.Put(JWSConstantsPayload, payload)
	}

	signatures := g.jwsJsonSerializationObject.Signatures()
	if utils.CollectionSize(signatures) != 1 {
		return nil, model.NewDSSError(fmt.Sprintf(
			"JSON Flattened Serialization can only contain 1 signature (current : %d)", utils.CollectionSize(signatures)))
	}

	jws := signatures[0]
	signatureJSONMap := g.signatureJSONMap(jws)
	flattenedJwsMap.PutAll(signatureJSONMap)

	return NewJsonObjectFromMap(flattenedJwsMap), nil
}

// signatureJSONMap builds the per-signature members: 'protected', optional 'header', 'signature'.
// Port of the private getSignatureJsonMap(JWS).
func (g *JWSJsonSerializationGenerator) signatureJSONMap(signature *JWS) *jose.Object {
	signatureMap := jose.NewObject()

	encodedProtected := signature.EncodedHeader()
	if utils.IsStringNotBlank(encodedProtected) {
		signatureMap.Put(JWSConstantsProtected, encodedProtected)
	}

	unprotected := signature.Unprotected()
	if unprotected.Size() != 0 {
		signatureMap.Put(JWSConstantsHeader, unprotected)
	}

	encodedSignatureValue := signature.EncodedSignature()
	signatureMap.Put(JWSConstantsSignature, encodedSignatureValue)

	return signatureMap
}
