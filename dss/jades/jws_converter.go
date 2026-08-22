// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/JWSConverter.java (DSS 6.5.RC1).
package jades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// The document names the converter stamps onto its results. Ports of the private static final
// String fields of the same names.
const (
	// jwsConverterFlattenedSerializationDocumentName is the name for a Flattened Serialization
	// signature. Port of FLATTENED_SERIALIZATION_DOCUMENT_NAME.
	jwsConverterFlattenedSerializationDocumentName = "json-flattened-serialization.json"

	// jwsConverterSerializationDocumentName is the name for a JSON Serialization signature. Port
	// of SERIALIZATION_DOCUMENT_NAME.
	jwsConverterSerializationDocumentName = "json-serialization.json"

	// jwsConverterClearEtsiUDocumentName is the name for a signature containing JSON components
	// in clear JSON form. Port of CLEAR_ETSIU_DOCUMENT_NAME.
	jwsConverterClearEtsiUDocumentName = "etsiU-clear-incorporation.json"

	// jwsConverterBase64UrlEtsiUDocumentName is the name for a signature containing JSON
	// components in their corresponding base64url encoded form. Port of
	// BASE64URL_ETSIU_DOCUMENT_NAME.
	jwsConverterBase64UrlEtsiUDocumentName = "etsiU-base64url-incorporation.json"
)

// jwsConverterTimestampHeaderNames lists the timestamp headers covering other 'etsiU' headers.
// Port of the private static timestampHeaderNames list.
//
// A component with one of these names cannot be re-encoded: the timestamp inside it was computed
// over the very octets of its siblings, so changing their incorporation form invalidates it.
var jwsConverterTimestampHeaderNames = []string{
	JAdESHeaderParameterNamesArcTst,
	JAdESHeaderParameterNamesRfsTst,
	JAdESHeaderParameterNamesSigRTst,
}

// JWSConverterFromJWSCompactToJSONFlattenedSerialization converts a JWS Compact Serialization to
// a JSON Flattened Serialization. Port of
// fromJWSCompactToJSONFlattenedSerialization(DSSDocument).
func JWSConverterFromJWSCompactToJSONFlattenedSerialization(document model.DSSDocument) (model.DSSDocument, error) {
	parser := NewJWSCompactSerializationParser(document)
	jws, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	jwsJsonSerializationObject := DSSJsonUtilsToJWSJsonSerializationObject(jws)
	generator := NewJWSJsonSerializationGenerator(jwsJsonSerializationObject,
		enumerations.JWSSerializationTypeFlattenedJSONSerialization)

	signatureDocument, err := generator.Generate()
	if err != nil {
		return nil, err
	}
	signatureDocument.SetName(jwsConverterFlattenedSerializationDocumentName)
	signatureDocument.SetMimeType(enumerations.MimeTypeEnumJSON)
	return signatureDocument, nil
}

// JWSConverterFromJWSCompactToJSONSerialization converts a JWS Compact Serialization to a JSON
// Serialization. Port of fromJWSCompactToJSONSerialization(DSSDocument).
func JWSConverterFromJWSCompactToJSONSerialization(document model.DSSDocument) (model.DSSDocument, error) {
	parser := NewJWSCompactSerializationParser(document)
	jws, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	jwsJsonSerializationObject := DSSJsonUtilsToJWSJsonSerializationObject(jws)
	generator := NewJWSJsonSerializationGenerator(jwsJsonSerializationObject,
		enumerations.JWSSerializationTypeJSONSerialization)

	signatureDocument, err := generator.Generate()
	if err != nil {
		return nil, err
	}
	signatureDocument.SetName(jwsConverterSerializationDocumentName)
	signatureDocument.SetMimeType(enumerations.MimeTypeEnumJSON)
	return signatureDocument, nil
}

// JWSConverterFromEtsiUWithBase64UrlToClearJSONIncorporation converts the unprotected content of
// the 'etsiU' header of the JAdES signatures inside a Serialization (or Flattened) document to
// its clear JSON incorporation form. Port of
// fromEtsiUWithBase64UrlToClearJsonIncorporation(DSSDocument).
func JWSConverterFromEtsiUWithBase64UrlToClearJSONIncorporation(document model.DSSDocument) (model.DSSDocument, error) {
	parser := NewJWSJsonSerializationParser(document)
	jwsJsonSerializationObject, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	for _, jws := range jwsJsonSerializationObject.Signatures() {
		etsiUContent := DSSJsonUtilsEtsiU(jws)
		if utils.IsCollectionEmpty(etsiUContent) {
			// do nothing
			continue
		}

		if err := jwsConverterAssertConvertPossible(etsiUContent); err != nil {
			return nil, err
		}

		clearEtsiUContent, err := jwsConverterToClearJSONIncorporation(etsiUContent)
		if err != nil {
			return nil, err
		}
		unprotected := jws.Unprotected()
		unprotected.Replace(JAdESHeaderParameterNamesEtsiU, clearEtsiUContent)
	}

	generator := NewJWSJsonSerializationGenerator(jwsJsonSerializationObject,
		jwsJsonSerializationObject.JWSSerializationType())

	signatureDocument, err := generator.Generate()
	if err != nil {
		return nil, err
	}
	signatureDocument.SetName(jwsConverterClearEtsiUDocumentName)
	signatureDocument.SetMimeType(enumerations.MimeTypeEnumJSON)
	return signatureDocument, nil
}

// jwsConverterAssertConvertPossible refuses a mixed-form 'etsiU' array. Port of the private
// assertConvertPossible(List).
func jwsConverterAssertConvertPossible(etsiUContent []any) error {
	if !DSSJsonUtilsCheckComponentsUnicity(etsiUContent) {
		return model.NewDSSError("Unable to convert the EtsiU content! All components shall have a common form.")
	}
	return nil
}

// jwsConverterToClearJSONIncorporation re-expresses every component as clear JSON. Port of the
// private toClearJsonIncorporation(List).
func jwsConverterToClearJSONIncorporation(etsiUContent []any) ([]any, error) {
	clearEtsiUContent := make([]any, 0, len(etsiUContent))
	for _, item := range etsiUContent {
		clearEtsiUComponent, ok := DSSJsonUtilsParseEtsiUComponent(item)
		if !ok {
			return nil, model.NewDSSError(fmt.Sprintf("Unable to parse 'etsiU' component : '%v'", item))
		}
		if err := jwsConverterAssertComponentSupportsConversion(clearEtsiUComponent); err != nil {
			return nil, err
		}
		clearEtsiUContent = append(clearEtsiUContent, NewJsonObjectFromMap(clearEtsiUComponent))
	}
	return clearEtsiUContent, nil
}

// JWSConverterFromEtsiUWithClearJSONToBase64UrlIncorporation converts the unprotected content of
// the 'etsiU' header of the JAdES signatures inside a Serialization (or Flattened) document to
// its base64url encoded incorporation form. Port of
// fromEtsiUWithClearJsonToBase64UrlIncorporation(DSSDocument).
func JWSConverterFromEtsiUWithClearJSONToBase64UrlIncorporation(document model.DSSDocument) (model.DSSDocument, error) {
	parser := NewJWSJsonSerializationParser(document)
	jwsJsonSerializationObject, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	for _, jws := range jwsJsonSerializationObject.Signatures() {
		etsiUContent := DSSJsonUtilsEtsiU(jws)
		if utils.IsCollectionEmpty(etsiUContent) {
			// do nothing
			continue
		}

		if err := jwsConverterAssertConvertPossible(etsiUContent); err != nil {
			return nil, err
		}

		base64UrlEtsiUContent, err := jwsConverterToBase64UrlIncorporation(etsiUContent)
		if err != nil {
			return nil, err
		}
		unprotected := jws.Unprotected()
		unprotected.Replace(JAdESHeaderParameterNamesEtsiU, base64UrlEtsiUContent)
	}

	generator := NewJWSJsonSerializationGenerator(jwsJsonSerializationObject,
		jwsJsonSerializationObject.JWSSerializationType())

	signatureDocument, err := generator.Generate()
	if err != nil {
		return nil, err
	}
	signatureDocument.SetName(jwsConverterBase64UrlEtsiUDocumentName)
	signatureDocument.SetMimeType(enumerations.MimeTypeEnumJSON)
	return signatureDocument, nil
}

// jwsConverterToBase64UrlIncorporation re-expresses every component as base64url text. Port of
// the private toBase64UrlIncorporation(List).
func jwsConverterToBase64UrlIncorporation(etsiUContent []any) ([]any, error) {
	base64UrlEtsiUContent := make([]any, 0, len(etsiUContent))
	for _, item := range etsiUContent {
		base64UrlEtsiUComponent, ok := DSSJsonUtilsParseEtsiUComponent(item)
		if !ok {
			return nil, model.NewDSSError(fmt.Sprintf("Unable to parse 'etsiU' component : '%v'", item))
		}
		if err := jwsConverterAssertComponentSupportsConversion(base64UrlEtsiUComponent); err != nil {
			return nil, err
		}
		base64UrlEtsiUContent = append(base64UrlEtsiUContent, DSSJsonUtilsToBase64UrlObject(base64UrlEtsiUComponent))
	}
	return base64UrlEtsiUContent, nil
}

// jwsConverterAssertComponentSupportsConversion refuses to re-encode a component whose name
// belongs to a timestamp covering its siblings. Port of the private
// assertComponentSupportsConversion(Map).
func jwsConverterAssertComponentSupportsConversion(etsiUComponent *jose.Object) error {
	// only one is allowed
	keys := etsiUComponent.Keys()
	if len(keys) == 0 {
		return nil
	}
	componentName := keys[0]
	for _, timestampHeaderName := range jwsConverterTimestampHeaderNames {
		if componentName == timestampHeaderName {
			return model.NewDSSError(fmt.Sprintf(
				"Unable to convert a signature! 'etsiU' contains a component with name '%s', "+
					"which is sensible to a format change.", componentName))
		}
	}
	return nil
}
