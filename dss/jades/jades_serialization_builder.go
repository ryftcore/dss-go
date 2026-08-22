// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESSerializationBuilder.java (DSS 6.5.RC1).
//
// Java offers two constructors, the second of which computes its super() argument from an
// existing JWSJsonSerializationObject through the static extractDocumentToBeSigned. Go has no
// constructor chaining, so the two become NewJAdESSerializationBuilder and
// NewJAdESSerializationBuilderFromSignature, both funnelling into InitAbstractJAdESBuilder;
// extractDocumentToBeSigned stays a package-level function.
//
// Every IllegalArgumentException becomes a returned error.
package jades

import (
	"errors"
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESSerializationBuilder builds a JWS JSON Serialization signature.
type JAdESSerializationBuilder struct {
	AbstractJAdESBuilder

	// jwsJsonSerializationObject is the JWS signature container.
	jwsJsonSerializationObject *JWSJsonSerializationObject
}

// NewJAdESSerializationBuilder is the default constructor.
// Port of JAdESSerializationBuilder(CertificateVerifier, JAdESSignatureParameters,
// List<DSSDocument>).
func NewJAdESSerializationBuilder(certificateVerifier validation.CertificateVerifier,
	parameters *JAdESSignatureParameters,
	documentsToSign []model.DSSDocument) (*JAdESSerializationBuilder, error) {
	builder := &JAdESSerializationBuilder{}
	if err := builder.InitAbstractJAdESBuilder(builder, certificateVerifier, parameters,
		documentsToSign); err != nil {
		return nil, err
	}
	return builder, nil
}

// NewJAdESSerializationBuilderFromSignature is the constructor from an existing signature.
// Port of JAdESSerializationBuilder(CertificateVerifier, JAdESSignatureParameters,
// JWSJsonSerializationObject).
func NewJAdESSerializationBuilderFromSignature(certificateVerifier validation.CertificateVerifier,
	parameters *JAdESSignatureParameters,
	jwsJsonSerializationObject *JWSJsonSerializationObject) (*JAdESSerializationBuilder, error) {
	documentsToSign, err := jadesSerializationBuilderExtractDocumentToBeSigned(parameters,
		jwsJsonSerializationObject)
	if err != nil {
		return nil, err
	}
	builder := &JAdESSerializationBuilder{}
	if err := builder.InitAbstractJAdESBuilder(builder, certificateVerifier, parameters,
		documentsToSign); err != nil {
		return nil, err
	}
	builder.jwsJsonSerializationObject = jwsJsonSerializationObject
	return builder, nil
}

// jadesSerializationBuilderExtractDocumentToBeSigned ports the private static
// extractDocumentToBeSigned.
func jadesSerializationBuilderExtractDocumentToBeSigned(parameters *JAdESSignatureParameters,
	jwsJsonSerializationObject *JWSJsonSerializationObject) ([]model.DSSDocument, error) {
	if utils.IsStringNotBlank(jwsJsonSerializationObject.Payload()) {
		// enveloping signature
		jws := jwsJsonSerializationObject.Signatures()[0]

		var payloadBytes []byte
		if jws.IsRfc7797UnencodedPayload() {
			payloadBytes = []byte(jwsJsonSerializationObject.Payload())
		} else {
			payloadBytes = DSSJsonUtilsFromBase64Url(jwsJsonSerializationObject.Payload())
		}
		return []model.DSSDocument{model.NewInMemoryDocument(payloadBytes)}, nil

	} else if utils.IsCollectionNotEmpty(parameters.DetachedContents()) {
		// detached signature
		return parameters.DetachedContents(), nil

	}
	return nil, errors.New("The payload or detached content must be provided!")
}

// Build builds the JWS JSON Serialization signature. Port of #build(SignatureValue).
func (b *JAdESSerializationBuilder) Build(signatureValue *model.SignatureValue) (model.DSSDocument, error) {
	if err := b.AssertConfigurationValidity(b.Parameters); err != nil {
		return nil, err
	}

	jws, err := b.jws()
	if err != nil {
		return nil, err
	}

	if b.jwsJsonSerializationObject == nil {
		b.jwsJsonSerializationObject = NewJWSJsonSerializationObject()
		if enumerations.SignaturePackagingDetached != b.Parameters.SignaturePackaging() {
			// do not include payload for detached case
			b.jwsJsonSerializationObject.SetPayload(jws.SignedPayload())
		}
	} else {
		if err := b.assertB64ConfigurationConsistent(); err != nil {
			return nil, err
		}
	}

	signatureValueBytes, err := spi.DSSASN1UtilsEnsurePlainSignatureValue(
		b.Parameters.EncryptionAlgorithm(), signatureValue.Value())
	if err != nil {
		return nil, err
	}
	jws.SetSignature(signatureValueBytes)

	// Java mutates the live List returned by getSignatures(); Go's slice is a value, so the
	// append goes through JWSJsonSerializationObject#AddSignature - same resulting container.
	b.jwsJsonSerializationObject.AddSignature(jws)

	generator := NewJWSJsonSerializationGenerator(b.jwsJsonSerializationObject,
		b.Parameters.JwsSerializationType())
	return generator.Generate()
}

// assertB64ConfigurationConsistent checks that all not detached signatures have the same 'b64'
// value. Port of the private assertB64ConfigurationConsistent.
func (b *JAdESSerializationBuilder) assertB64ConfigurationConsistent() error {
	// verify only for non-detached cases
	if enumerations.SignaturePackagingDetached != b.Parameters.SignaturePackaging() {
		base64UrlEncodedPayload := b.Parameters.IsBase64UrlEncodedPayload()
		for _, jws := range b.jwsJsonSerializationObject.Signatures() {
			if base64UrlEncodedPayload != !jws.IsRfc7797UnencodedPayload() {
				return errors.New("'b64' value shall be the same for all signatures! " +
					"Change 'Base64UrlEncodedPayload' signature parameter or sign another file!")
			}
		}
	}
	return nil
}

// jws ports the private getJWS.
func (b *JAdESSerializationBuilder) jws() (*JWS, error) {
	jws := NewJWS()
	if err := b.IncorporateHeader(jws); err != nil {
		return nil, err
	}
	if err := b.IncorporatePayload(jws); err != nil {
		return nil, err
	}
	return jws, nil
}

// MimeType returns the MimeType of the produced signature. Port of #getMimeType.
func (b *JAdESSerializationBuilder) MimeType() enumerations.MimeType {
	return enumerations.MimeTypeEnumJOSEJSON
}

// AssertConfigurationValidity verifies that the configured signature packaging and serialization
// type are supported by a JWS JSON Serialization signature.
// Port of the protected, overridden #assertConfigurationValidity.
func (b *JAdESSerializationBuilder) AssertConfigurationValidity(
	signatureParameters *JAdESSignatureParameters) error {
	packaging := signatureParameters.SignaturePackaging()
	if packaging != enumerations.SignaturePackagingEnveloping &&
		packaging != enumerations.SignaturePackagingDetached {
		return fmt.Errorf("Unsupported signature packaging for JSON Serialization Signature: %s", packaging)
	}
	if enumerations.JWSSerializationTypeJSONSerialization != signatureParameters.JwsSerializationType() &&
		b.jwsJsonSerializationObject != nil {
		return fmt.Errorf("The '%s' type is not supported for a parallel signing!",
			signatureParameters.JwsSerializationType())
	}
	return nil
}

// Compile-time assertion that *JAdESSerializationBuilder satisfies the builder contract.
var _ JAdESBuilder = (*JAdESSerializationBuilder)(nil)
