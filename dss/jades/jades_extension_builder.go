// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESExtensionBuilder.java (DSS 6.5.RC1).
//
// Java overloads assertEtsiUComponentsConsistent on (JWS, JAdESSignatureParameters) and
// (JWS, Boolean); Go cannot overload, so the parameter-taking one keeps the plain name and the
// Boolean-taking one - the one JAdESSignaturePolicyStoreBuilder calls with a plain boolean -
// becomes AssertEtsiUComponentsConsistentWithEncoding.
//
// Java's `Boolean` (a nullable tri-state: TRUE, FALSE, or "not configured") becomes *bool here,
// since the null case drives real behaviour: it means "adopt whatever the existing etsiU uses".
//
// Every IllegalInputException becomes a returned error; slf4j logging is dropped.
package jades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESExtensionBuilder is the abstract base allowing the signature extension.
type JAdESExtensionBuilder struct {
}

// AssertEtsiUComponentsConsistent checks that the type of etsiU components is consistent and
// writes the resolved encoding back into the signature parameters.
// Port of the protected #assertEtsiUComponentsConsistent(JWS, JAdESSignatureParameters).
func (b *JAdESExtensionBuilder) AssertEtsiUComponentsConsistent(jws *JWS,
	signatureParameters *JAdESSignatureParameters) error {
	isBase64UrlEtsiUComponents := signatureParameters.IsBase64UrlEncodedEtsiUComponents()
	resolved, err := b.AssertEtsiUComponentsConsistentWithEncoding(jws, isBase64UrlEtsiUComponents)
	if err != nil {
		return err
	}
	signatureParameters.SetBase64UrlEncodedEtsiUComponents(&resolved)
	return nil
}

// AssertEtsiUComponentsConsistentWithEncoding checks that the type of etsiU components is
// consistent and returns the target encoding: TRUE if the etsiU parameters shall be base64url
// encoded, FALSE otherwise.
// Port of the protected #assertEtsiUComponentsConsistent(JWS, Boolean).
func (b *JAdESExtensionBuilder) AssertEtsiUComponentsConsistentWithEncoding(jws *JWS,
	isBase64UrlEtsiUComponents *bool) (bool, error) {
	etsiU := DSSJsonUtilsEtsiU(jws)
	if utils.IsCollectionNotEmpty(etsiU) {
		if !DSSJsonUtilsCheckComponentsUnicity(etsiU) {
			return false, exception.NewIllegalInputException(
				"Extension is not possible, because components of the 'etsiU' header have " +
					"not common format! Shall be all Strings or Objects.")
		}

		isEtsiUInBase64UrlForm := DSSJsonUtilsAreAllBase64UrlComponents(etsiU)
		if isBase64UrlEtsiUComponents == nil {
			// Upstream logs "base64UrlEtsiUComponents parameter is not defined. The check of
			// etsiU unsigned header structure is skipped. Use the current value."

		} else if *isBase64UrlEtsiUComponents != isEtsiUInBase64UrlForm {
			return false, exception.NewIllegalInputException(fmt.Sprintf(
				"Extension is not possible! The encoding of 'etsiU' "+
					"components shall match! Use jadesSignatureParameters.setBase64UrlEncodedEtsiUComponents(%t)",
				!*isBase64UrlEtsiUComponents))
		}
		return isEtsiUInBase64UrlForm, nil

	} else if isBase64UrlEtsiUComponents == nil {
		// Upstream traces "base64UrlEtsiUComponents parameters is not defined. Use the default
		// value (true)."
		return true, nil
	}

	return *isBase64UrlEtsiUComponents, nil
}

// AssertJWSJsonSerializationObjectValid checks that the jwsJsonSerializationObject is valid and
// can be extended. Port of the protected #assertJWSJsonSerializationObjectValid.
func (b *JAdESExtensionBuilder) AssertJWSJsonSerializationObjectValid(
	jwsJsonSerializationObject *JWSJsonSerializationObject) error {
	if jwsJsonSerializationObject == nil {
		return exception.NewIllegalInputException(
			"The provided document is not a valid JAdES signature! Unable to extend.")
	}
	if utils.IsCollectionEmpty(jwsJsonSerializationObject.Signatures()) {
		return exception.NewIllegalInputException("No signatures found to be extended!")
	}
	if !jwsJsonSerializationObject.IsValid() {
		return exception.NewIllegalInputException(fmt.Sprintf(
			"Signature extension is not supported for invalid RFC 7515 files "+
				"(shall be a Serializable JAdES signature). Reason(s) : %s",
			jwsJsonSerializationObject.StructuralValidationErrors()))
	}
	return nil
}

// AssertJSONSerializationObjectMayBeExtended checks that the given jwsJsonSerializationObject can
// be extended. Port of the protected #assertJSONSerializationObjectMayBeExtended.
func (b *JAdESExtensionBuilder) AssertJSONSerializationObjectMayBeExtended(
	jwsJsonSerializationObject *JWSJsonSerializationObject) error {
	if err := b.AssertJWSJsonSerializationObjectValid(jwsJsonSerializationObject); err != nil {
		return err
	}

	jwsSerializationType := jwsJsonSerializationObject.JWSSerializationType()
	if enumerations.JWSSerializationType_JSON_SERIALIZATION != jwsSerializationType &&
		enumerations.JWSSerializationType_FLATTENED_JSON_SERIALIZATION != jwsSerializationType {
		return exception.NewIllegalInputException(
			"The extended signature shall have JSON Serialization (or Flattened) type! " +
				"Use JWSConverter to convert the signature.")
	}
	return nil
}
