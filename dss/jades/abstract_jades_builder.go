// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/AbstractJAdESBuilder.java (DSS 6.5.RC1).
//
// # assertConfigurationValidity, and how Java's abstract call survives the port
//
// Java's AbstractJAdESBuilder#buildDataToBeSigned() calls the abstract
// assertConfigurationValidity(parameters), which JAdESCompactBuilder and JAdESSerializationBuilder
// implement differently. Go has no virtual dispatch across embedding, so - per the
// TokenBase.InitToken(self) convention of PORTING.md, and exactly as xades.XAdESLevelBaselineT
// does - the concrete builder registers itself through InitAbstractJAdESBuilder(self, ...) and the
// base dispatches into it via b.overrides.
//
// # Errors
//
// Objects.requireNonNull becomes a panic carrying the Java message; the IllegalArgumentException
// raised for an empty document list, and every failure raised further down (payload reading,
// protected-header serialization), become returned errors. slf4j tracing of the payload is
// dropped.
package jades

import (
	"errors"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESBuilderOverrides declares the operation Java leaves abstract in AbstractJAdESBuilder and
// that the base implementation itself calls back into. Both concrete builders satisfy it.
type JAdESBuilderOverrides interface {
	// AssertConfigurationValidity verifies that the given signaturePackaging type is supported,
	// returning an error when the configuration is not valid. Port of the protected abstract
	// #assertConfigurationValidity.
	AssertConfigurationValidity(signatureParameters *JAdESSignatureParameters) error
}

// AbstractJAdESBuilder is the abstract JAdES signature builder.
type AbstractJAdESBuilder struct {
	// Parameters holds the signature parameters. Port of the protected final `parameters`.
	Parameters *JAdESSignatureParameters

	// JadesLevelBaselineB is the instance of the B-level generator class. Port of the protected
	// final `jadesLevelBaselineB`.
	JadesLevelBaselineB *JAdESLevelBaselineB

	// overrides points back at the concrete builder; see InitAbstractJAdESBuilder.
	overrides JAdESBuilderOverrides
}

// InitAbstractJAdESBuilder registers the concrete builder with this base and builds the B-level
// generator. Port of the protected AbstractJAdESBuilder(CertificateVerifier,
// JAdESSignatureParameters, List<DSSDocument>) constructor.
func (b *AbstractJAdESBuilder) InitAbstractJAdESBuilder(self JAdESBuilderOverrides,
	certificateVerifier validation.CertificateVerifier, parameters *JAdESSignatureParameters,
	documentsToSign []model.DSSDocument) error {
	if certificateVerifier == nil {
		panic("CertificateVerifier must be defined!")
	}
	if parameters == nil {
		panic("SignatureParameters must be defined!")
	}
	if utils.IsCollectionEmpty(documentsToSign) {
		return errors.New("Documents to sign must be provided!")
	}
	b.overrides = self
	b.Parameters = parameters
	levelBaselineB, err := NewJAdESLevelBaselineB(certificateVerifier, parameters, documentsToSign)
	if err != nil {
		return err
	}
	b.JadesLevelBaselineB = levelBaselineB
	return nil
}

// BuildDataToBeSigned builds the data to be signed. Port of #buildDataToBeSigned.
func (b *AbstractJAdESBuilder) BuildDataToBeSigned() (*model.ToBeSigned, error) {
	if err := b.overrides.AssertConfigurationValidity(b.Parameters); err != nil {
		return nil, err
	}

	jws := NewJWS()
	if err := b.IncorporateHeader(jws); err != nil {
		return nil, err
	}
	if err := b.IncorporatePayload(jws); err != nil {
		return nil, err
	}

	dataToSign := DSSJsonUtilsSigningInputBytes(jws)
	return model.NewToBeSignedWithBytes(dataToSign), nil
}

// IncorporateHeader incorporates the signed header into the given JWS. The insertion order of
// the JAdESLevelBaselineB properties - which is what the serialized protected header, and hence
// the signature, depends on - is preserved by the jose.Object the B-level generator builds, the
// Go counterpart of Java's LinkedHashMap. Port of the protected #incorporateHeader.
func (b *AbstractJAdESBuilder) IncorporateHeader(jws *JWS) error {
	signedProperties, err := b.JadesLevelBaselineB.SignedProperties()
	if err != nil {
		return err
	}
	for _, headerName := range signedProperties.Keys() {
		jws.SetHeader(headerName, signedProperties.Value(headerName))
	}
	return nil
}

// IncorporatePayload incorporates the payload into the given JWS.
// Port of the protected #incorporatePayload.
func (b *AbstractJAdESBuilder) IncorporatePayload(jws *JWS) error {
	payloadBytes, err := b.JadesLevelBaselineB.PayloadBytes()
	if err != nil {
		return err
	}
	if utils.IsArrayNotEmpty(payloadBytes) {
		// Upstream traces the payload and its base64 encoding here.
		jws.SetPayloadOctets(payloadBytes)
	}
	return nil
}
