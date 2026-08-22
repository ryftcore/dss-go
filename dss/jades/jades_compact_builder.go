// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESCompactBuilder.java (DSS 6.5.RC1).
//
// Note the asymmetry upstream deliberately keeps and this port reproduces: build() skips
// incorporatePayload for DETACHED packaging (so the middle segment of the compact serialization
// is empty), while the inherited buildDataToBeSigned() incorporates it unconditionally - the
// detached payload contributes to the signature but is not carried in the JWS.
//
// Objects.requireNonNull becomes a panic carrying the Java message; every
// IllegalArgumentException becomes a returned error.
package jades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// CompactBuilder builds a JWS Compact Signature.
type CompactBuilder struct {
	AbstractBuilder
}

// NewCompactBuilder is the default constructor.
// Port of JAdESCompactBuilder(CertificateVerifier, JAdESSignatureParameters, List<DSSDocument>).
func NewCompactBuilder(certificateVerifier validation.CertificateVerifier,
	parameters *SignatureParameters, documentsToSign []model.DSSDocument) (*CompactBuilder, error) {
	builder := &CompactBuilder{}
	if err := builder.InitAbstractJAdESBuilder(builder, certificateVerifier, parameters,
		documentsToSign); err != nil {
		return nil, err
	}
	return builder, nil
}

// Build builds the concatenation of signed header and payload (the dataToBeSigned string) in the
// way: BASE64URL(UTF8(JWS Protected Header)) || '.' || BASE64URL(JWS Payload), followed by the
// signature value. Port of #build(SignatureValue).
func (b *CompactBuilder) Build(signatureValue *model.SignatureValue) (model.DSSDocument, error) {
	if err := b.AssertConfigurationValidity(b.Parameters); err != nil {
		return nil, err
	}

	jws := NewJWS()
	if err := b.IncorporateHeader(jws); err != nil {
		return nil, err
	}
	if enumerations.SignaturePackagingDetached != b.Parameters.SignaturePackaging() {
		if err := b.IncorporatePayload(jws); err != nil {
			return nil, err
		}
	}
	var payload string
	if b.Parameters.IsBase64UrlEncodedPayload() {
		payload = jws.EncodedPayload()
	} else {
		payload = jws.UnverifiedPayload()
	}
	signatureValueBytes, err := spi.DSSASN1UtilsEnsurePlainSignatureValue(
		b.Parameters.EncryptionAlgorithm(), signatureValue.Value())
	if err != nil {
		return nil, err
	}

	signatureString := DSSJsonUtilsConcatenate(jws.EncodedHeader(), payload,
		DSSJsonUtilsToBase64Url(signatureValueBytes))
	return model.NewInMemoryDocument([]byte(signatureString)), nil
}

// MimeType returns the MimeType of the produced signature. Port of #getMimeType.
func (b *CompactBuilder) MimeType() enumerations.MimeType {
	return enumerations.MimeTypeEnumJOSE
}

// AssertConfigurationValidity verifies that the configured signature packaging and level are
// supported by a JAdES Compact Signature.
// Port of the protected, overridden #assertConfigurationValidity.
func (b *CompactBuilder) AssertConfigurationValidity(
	signatureParameters *SignatureParameters) error {
	if signatureParameters.SignaturePackaging() == "" {
		panic("SignaturePackaging shall be defined!")
	}
	if signatureParameters.SignatureLevel() == "" {
		panic("SignatureLevel shall be defined!")
	}

	packaging := signatureParameters.SignaturePackaging()
	if enumerations.SignaturePackagingEnveloping != packaging &&
		enumerations.SignaturePackagingDetached != packaging {
		return fmt.Errorf("Unsupported signature packaging for JAdES Compact Signature: %s", packaging)
	}
	signatureLevel := signatureParameters.SignatureLevel()
	if enumerations.SignatureLevelJAdESBaselineB != signatureLevel {
		return fmt.Errorf("Only JAdES_BASELINE_B level is allowed for JAdES Compact Signature! " +
			"Change JwsSerializationType in JAdESSignatureParameters in order to support extension!")
	}
	return nil
}

// Compile-time assertion that *CompactBuilder satisfies the builder contract.
var _ Builder = (*CompactBuilder)(nil)
