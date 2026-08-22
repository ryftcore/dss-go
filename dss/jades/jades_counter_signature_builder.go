// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESCounterSignatureBuilder.java (DSS 6.5.RC1).
//
// The counter-signature is embedded as a 'cSig' component of the master signature's 'etsiU'
// unsigned header: as a plain string for a compact counter-signature, as a parsed JSON object for
// a flattened one - and the ordered JsonObject the parse yields is what keeps the member order of
// the embedded signature intact, which matters because an outer archive time-stamp may cover it.
//
// updateMasterSignatureRecursively re-serializes the whole nested chain from the innermost
// signature outwards, replacing each 'cSig' component in place (keeping its EtsiUComponent
// identifier), so nesting depth is preserved.
//
// Objects.requireNonNull becomes a panic carrying the Java message; every
// IllegalArgumentException / IllegalInputException becomes a returned error.
package jades

import (
	"errors"
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// CounterSignatureBuilder creates a JAdES counter signature.
type CounterSignatureBuilder struct {
	ExtensionBuilder
}

// NewJAdESCounterSignatureBuilder is the default constructor.
func NewJAdESCounterSignatureBuilder() *CounterSignatureBuilder {
	return &CounterSignatureBuilder{}
}

// GetSignatureValueToBeSigned extracts the SignatureValue binaries from the provided JAdES
// signature. Port of #getSignatureValueToBeSigned.
func (b *CounterSignatureBuilder) GetSignatureValueToBeSigned(signatureDocument model.DSSDocument,
	parameters *CounterSignatureParameters) (model.DSSDocument, error) {

	documentAnalyzerFactory := NewJWSDocumentAnalyzerFactory()
	documentAnalyzer := jwsDocumentAnalyzerBase(documentAnalyzerFactory.Create(signatureDocument))

	jwsJsonSerializationObject := documentAnalyzer.JwsJsonSerializationObject()
	if err := b.AssertJSONSerializationObjectMayBeExtended(jwsJsonSerializationObject); err != nil {
		return nil, err
	}

	signatures := documentAnalyzer.Signatures()

	signature, err := b.extractSignatureById(signatures, parameters.SignatureIdToCounterSign())
	if err != nil {
		return nil, err
	}
	jadesSignature, ok := signature.(*Signature)
	if !ok {
		return nil, fmt.Errorf("unexpected signature type %T", signature)
	}
	return model.NewInMemoryDocument(jadesSignature.SignatureValue()), nil
}

// BuildEmbeddedCounterSignature embeds the counter signature into the original JAdES signature
// and returns the result. Port of #buildEmbeddedCounterSignature.
func (b *CounterSignatureBuilder) BuildEmbeddedCounterSignature(signatureDocument,
	counterSignature model.DSSDocument,
	parameters *CounterSignatureParameters) (model.DSSDocument, error) {

	documentAnalyzerFactory := NewJWSDocumentAnalyzerFactory()
	documentAnalyzer := jwsDocumentAnalyzerBase(documentAnalyzerFactory.Create(signatureDocument))

	jwsJsonSerializationObject := documentAnalyzer.JwsJsonSerializationObject()
	if err := b.AssertJSONSerializationObjectMayBeExtended(jwsJsonSerializationObject); err != nil {
		return nil, err
	}

	signatures := documentAnalyzer.Signatures()

	signature, err := b.extractSignatureById(signatures, parameters.SignatureIdToCounterSign())
	if err != nil {
		return nil, err
	}
	jadesSignature, ok := signature.(*Signature)
	if !ok {
		return nil, fmt.Errorf("unexpected signature type %T", signature)
	}
	if err := b.AssertEtsiUComponentsConsistent(jadesSignature.Jws(),
		&parameters.SignatureParameters); err != nil {
		return nil, err
	}

	cSig, err := jadesCounterSignatureBuilderCSig(counterSignature, parameters.JwsSerializationType())
	if err != nil {
		return nil, err
	}

	etsiUHeader := jadesSignature.EtsiUHeader()
	if err := etsiUHeader.AddComponent(JAdESHeaderParameterNamesCSig, cSig,
		utils.IsTrue(parameters.IsBase64UrlEncodedEtsiUComponents())); err != nil {
		return nil, err
	}

	if err := b.updateMasterSignatureRecursively(jadesSignature); err != nil {
		return nil, err
	}

	generator := NewJWSJsonSerializationGenerator(jwsJsonSerializationObject,
		jwsJsonSerializationObject.JWSSerializationType())
	return generator.Generate()
}

// updateMasterSignatureRecursively ports the private updateMasterSignatureRecursively.
func (b *CounterSignatureBuilder) updateMasterSignatureRecursively(
	jadesSignature *Signature) error {
	masterSignature, ok := jadesSignature.MasterSignature().(*Signature)
	if !ok || masterSignature == nil {
		return nil
	}
	masterCSigAttribute := jadesSignature.MasterCSigComponent()

	jwsJsonSerializationObject := jadesSignature.Jws().JwsJsonSerializationObject()
	generator := NewJWSJsonSerializationGenerator(jwsJsonSerializationObject,
		jwsJsonSerializationObject.JWSSerializationType())

	generated, err := generator.Generate()
	if err != nil {
		return err
	}
	cSig, err := jadesCounterSignatureBuilderCSig(generated, jwsJsonSerializationObject.JWSSerializationType())
	if err != nil {
		return err
	}
	// EtsiUComponentBuildFromValue wants the concrete *AttributeIdentifier; going through
	// the exported Identifier() would unwrap to the generic identifier.SignatureAttributeIdentifier
	// interface (see jades_attribute.go's Identifier()), so this reaches the unexported
	// identifier field directly, same-package, exactly as EtsiUComponent's own constructor does.
	updatedCSigAttribute := EtsiUComponentBuildFromValue(JAdESHeaderParameterNamesCSig, cSig,
		masterCSigAttribute.IsBase64UrlEncoded(), masterCSigAttribute.identifier)
	if err := b.replaceCSigComponent(jadesSignature, updatedCSigAttribute); err != nil {
		return err
	}

	return b.updateMasterSignatureRecursively(masterSignature)
}

// jadesCounterSignatureBuilderCSig ports the private getCSig.
func jadesCounterSignatureBuilderCSig(counterSignature model.DSSDocument,
	jwsSerializationType enumerations.JWSSerializationType) (any, error) {
	binaries, err := spi.DSSUtilsToByteArrayOfDocument(counterSignature)
	if err != nil {
		return nil, err
	}
	signatureString := string(binaries)
	switch jwsSerializationType {
	case enumerations.JWSSerializationTypeCompactSerialization:
		return signatureString, nil
	case enumerations.JWSSerializationTypeFlattenedJSONSerialization:
		parsed, err := DSSJsonUtilsParseJSONStringToMap(signatureString)
		if err != nil {
			return nil, exception.NewIllegalInputExceptionWithCause(
				fmt.Sprintf("Unable to parse a counter signature. Reason : %s", err.Error()), err)
		}
		return parsed, nil
	default:
		return nil, fmt.Errorf(
			"The JWSSerializarionType '%s' is not supported for a Counter Signature!", jwsSerializationType)
	}
}

// replaceCSigComponent ports the private replaceCSigComponent.
func (b *CounterSignatureBuilder) replaceCSigComponent(jadesSignature *Signature,
	cSigAttribute *EtsiUComponent) error {
	masterSignature, ok := jadesSignature.MasterSignature().(*Signature)
	if !ok {
		return fmt.Errorf("unexpected master signature type %T", jadesSignature.MasterSignature())
	}
	etsiUHeader := masterSignature.EtsiUHeader()
	if err := etsiUHeader.ReplaceComponent(cSigAttribute); err != nil {
		return err
	}

	jadesSignature.SetMasterCSigComponent(cSigAttribute)
	return nil
}

// extractSignatureById ports the private extractSignatureById.
func (b *CounterSignatureBuilder) extractSignatureById(signatures []validation.AdvancedSignature,
	signatureId string) (validation.AdvancedSignature, error) {
	if signatureId == "" {
		panic("The Id of a signature to be counter signed shall be defined! " +
			"Please use SerializableCounterSignatureParameters.setSignatureIdToCounterSign(signatureId) method.")
	}

	if utils.IsCollectionEmpty(signatures) {
		return nil, errors.New("The provided signatureDocument does not contain JAdES Signatures!")
	}
	for _, signature := range signatures {
		jadesSignature, ok := signature.(*Signature)
		if !ok {
			return nil, fmt.Errorf("unexpected signature type %T", signature)
		}
		signatureById, err := b.signatureOrItsCounterSignature(jadesSignature, signatureId)
		if err != nil {
			return nil, err
		}
		if signatureById != nil {
			return signatureById, nil
		}
	}
	return nil, fmt.Errorf(
		"The requested JAdES Signature with id '%s' has not been found in the provided file!", signatureId)
}

// signatureOrItsCounterSignature ports the private getSignatureOrItsCounterSignature.
func (b *CounterSignatureBuilder) signatureOrItsCounterSignature(signature *Signature,
	signatureId string) (*Signature, error) {
	if signatureId == "" || signatureId == signature.ID() {
		return signature, nil
	}

	cSigComponents := DSSJsonUtilsUnsignedPropertiesWithHeaderName(signature.EtsiUHeader(),
		JAdESHeaderParameterNamesCSig)

	if utils.IsCollectionNotEmpty(cSigComponents) {
		for _, cSigComponent := range cSigComponents {

			// check timestamp before incorporating a new property
			if signature.TimestampSource().IsTimestamped(signatureId,
				enumerations.TimestampedObjectTypeSignature) {
				return nil, exception.NewIllegalInputException(fmt.Sprintf(
					"Unable to counter sign a signature with Id '%s'. "+
						"The signature is timestamped by a master signature!", signatureId))
			}

			counterSignature, err := DSSJsonUtilsExtractJAdESCounterSignature(cSigComponent, signature)
			if err != nil {
				return nil, err
			}
			signatureById, err := b.signatureOrItsCounterSignature(counterSignature, signatureId)
			if err != nil {
				return nil, err
			}
			if signatureById != nil {
				if _, isString := cSigComponent.Value().(string); isString {
					return nil, exception.NewIllegalInputException(
						"Unable to extend a Compact JAdES Signature with id '" + signatureId + "'")
				}
				return signatureById, nil
			}

		}
	}

	return nil, nil
}
