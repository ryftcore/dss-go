// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XAdESLevelBaselineLT.java (DSS 6.5.RC1).
//
// Java extends XAdESLevelBaselineT and overrides extendSignatures(List); the Go port embeds the
// -T level and registers itself with InitXAdESLevelBaselineLT so that the public
// ExtendSignatures entry point dispatches here. "super.extendSignatures(signatures)" is the
// explicit lt.XAdESLevelBaselineT.ExtendSignatures call.
//
// Java's private incorporateValidationDataForTimestamps(ValidationDataContainer, ...) collides
// by name with the differently-shaped private method XAdESLevelBaselineLTA declares; since the
// LTA type embeds this one, the Go port keeps them apart by name
// (incorporateValidationDataForTimestamps here, ...FromContainer there) instead of relying on
// shadowing.
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// XAdESLevelBaselineLT is the LT profile of a XAdES signature.
type XAdESLevelBaselineLT struct {
	XAdESLevelBaselineT
}

// NewXAdESLevelBaselineLT is the default constructor for XAdESLevelBaselineLT.
// Port of XAdESLevelBaselineLT(CertificateVerifier).
func NewXAdESLevelBaselineLT(certificateVerifier validation.CertificateVerifier) *XAdESLevelBaselineLT {
	extension := &XAdESLevelBaselineLT{}
	extension.InitXAdESLevelBaselineLT(extension, certificateVerifier)
	return extension
}

// InitXAdESLevelBaselineLT registers the concrete extension level with this base and forwards to
// the -T level. Port of the super(certificateVerifier) call of
// XAdESLevelBaselineLT(CertificateVerifier).
func (lt *XAdESLevelBaselineLT) InitXAdESLevelBaselineLT(self XAdESSignatureExtensionOverrides,
	certificateVerifier validation.CertificateVerifier) {
	lt.InitXAdESLevelBaselineT(self, certificateVerifier)
}

// ExtendSignatures adds CertificateValues and RevocationValues segments to
// UnsignedSignatureProperties. An XML electronic signature MAY contain at most one
// CertificateValues element and at most one RevocationValues element.
// Port of the overridden protected #extendSignatures(List).
func (lt *XAdESLevelBaselineLT) ExtendSignatures(signatures []validation.AdvancedSignature) error {
	if err := lt.XAdESLevelBaselineT.ExtendSignatures(signatures); err != nil {
		return err
	}

	signaturesToExtend := lt.extendToLTLevelSignatures(signatures)
	if utils.IsCollectionEmpty(signaturesToExtend) {
		return nil
	}

	// Reset sources
	for _, signature := range signaturesToExtend {
		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			// Java's (XAdESSignature) cast; a non-XAdES signature would raise a ClassCastException.
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		if _, err := lt.InitializeSignatureBuilder(xadesSignature); err != nil {
			return err
		}

		// Data sources can already be loaded in memory (force reload)
		lt.XadesSignature.ResetCertificateSource()
		lt.XadesSignature.ResetRevocationSources()
		lt.XadesSignature.ResetTimestampSource()
	}

	signatureRequirementsChecker := lt.SignatureRequirementsChecker()
	if enumerations.SignatureLevelXAdESBaselineLT == lt.Params.SignatureLevel() {
		signatureRequirementsChecker.AssertExtendToLTLevelPossible(signaturesToExtend)
	}

	signatureRequirementsChecker.AssertSignaturesValid(signaturesToExtend)
	signatureRequirementsChecker.AssertCertificateChainValidForLTLevel(signaturesToExtend)

	// Perform signature validation
	validationDataContainer, err := lt.DocumentAnalyzer.GetValidationData(signaturesToExtend)
	if err != nil {
		return err
	}

	// Append ValidationData
	for _, signature := range signaturesToExtend {
		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			// Java's (XAdESSignature) cast; a non-XAdES signature would raise a ClassCastException.
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		if _, err := lt.InitializeSignatureBuilder(xadesSignature); err != nil {
			return err
		}
		if signatureRequirementsChecker.HasLTALevelOrHigher(signature) {
			// avoid overriding of elements, when covered by an ArchiveTimeStamp
			continue
		}

		indent, err := lt.RemoveOldCertificateValues()
		if err != nil {
			return err
		}
		if err := lt.RemoveOldRevocationValues(); err != nil {
			return err
		}
		anyDataIndent, err := lt.RemoveLastTimestampAndAnyValidationData()
		if err != nil {
			return err
		}
		if indent == "" {
			indent = anyDataIndent
		}

		levelTUnsignedProperties := lt.UnsignedSignaturePropertiesDom.Clone(true)

		includedValidationData, err := lt.incorporateValidationDataForSignature(validationDataContainer,
			signature, indent)
		if err != nil {
			return err
		}
		if err := lt.incorporateValidationDataForTimestamps(validationDataContainer, signature, indent,
			includedValidationData); err != nil {
			return err
		}

		indented, err := lt.IndentIfPrettyPrint(lt.UnsignedSignaturePropertiesDom, levelTUnsignedProperties)
		if err != nil {
			return err
		}
		lt.UnsignedSignaturePropertiesDom = indented
	}
	return nil
}

// ValidationDataEncapsulationStrategy returns the ValidationDataEncapsulationStrategy to be used.
// Port of the protected #getValidationDataEncapsulationStrategy.
func (lt *XAdESLevelBaselineLT) ValidationDataEncapsulationStrategy() enumerations.ValidationDataEncapsulationStrategy {
	if lt.Params.IsEn319132() {
		return lt.Params.ValidationDataEncapsulationStrategy()
	}
	// AnyValidationData is not supported in old XAdES definition
	return enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationData
}

// incorporateValidationDataForSignature incorporates the validation data for the signature
// validation, according to the chosen validation data encapsulation mechanism, and returns the
// incorporated validation data. Port of the private incorporateValidationDataForSignature.
func (lt *XAdESLevelBaselineLT) incorporateValidationDataForSignature(
	validationDataContainer *validation.ValidationDataContainer,
	signature validation.AdvancedSignature, indent string) (*validation.ValidationData, error) {
	var validationDataForInclusion *validation.ValidationData
	validationDataEncapsulationStrategy := lt.ValidationDataEncapsulationStrategy()
	switch validationDataEncapsulationStrategy {
	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationData,
		enumerations.ValidationDataEncapsulationStrategyAnyValidationDataOnly:
		validationDataForInclusion = validationDataContainer.AllValidationDataForSignatureForInclusion(signature)

	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataLTSeparated,
		enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataAndAnyValidationData,
		enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndAnyValidationData:
		validationDataForInclusion = validationDataContainer.ValidationDataForSignatureForInclusion(signature)
		validationDataForInclusion.AddValidationData(
			validationDataContainer.ValidationDataForCounterSignaturesForInclusion(signature))

	default:
		return nil, fmt.Errorf("The ValidationDataEncapsulationStrategy '%s' is not supported!",
			validationDataEncapsulationStrategy)
	}

	certificateValuesToAdd := validationDataForInclusion.CertificateTokens()
	crlsToAdd := validationDataForInclusion.CrlTokens()
	ocspsToAdd := validationDataForInclusion.OcspTokens()

	switch validationDataEncapsulationStrategy {
	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationData,
		enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataLTSeparated,
		enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataAndAnyValidationData,
		enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndAnyValidationData:
		if err := lt.IncorporateCertificateValuesWithIndent(lt.UnsignedSignaturePropertiesDom,
			certificateValuesToAdd, indent); err != nil {
			return nil, err
		}
		if err := lt.IncorporateRevocationValuesWithIndent(lt.UnsignedSignaturePropertiesDom,
			crlsToAdd, ocspsToAdd, indent); err != nil {
			return nil, err
		}

	case enumerations.ValidationDataEncapsulationStrategyAnyValidationDataOnly:
		if err := lt.IncorporateAnyValidationData(validationDataForInclusion, indent); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("The ValidationDataEncapsulationStrategy '%s' is not supported!",
			validationDataEncapsulationStrategy)
	}
	return validationDataForInclusion, nil
}

// incorporateValidationDataForTimestamps incorporates the validation data for the signature
// timestamps validation, according to the chosen validation data encapsulation mechanism,
// excluding validationDataToExclude to avoid duplicates.
// Port of the private incorporateValidationDataForTimestamps.
func (lt *XAdESLevelBaselineLT) incorporateValidationDataForTimestamps(
	validationDataContainer *validation.ValidationDataContainer,
	signature validation.AdvancedSignature, indent string,
	validationDataToExclude *validation.ValidationData) error {
	var validationData *validation.ValidationData
	validationDataEncapsulationStrategy := lt.ValidationDataEncapsulationStrategy()
	switch validationDataEncapsulationStrategy {
	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataLTSeparated:
		validationData = validationDataContainer.ValidationDataForSignatureTimestampsForInclusion(signature)
		validationData.AddValidationData(
			validationDataContainer.ValidationDataForCounterSignatureTimestampsForInclusion(signature))
		validationData.ExcludeValidationData(validationDataToExclude)
		return lt.IncorporateTimestampValidationData(validationData, indent)

	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataAndAnyValidationData:
		validationData = validationDataContainer.ValidationDataForSignatureTimestampsForInclusion(signature)
		validationData.ExcludeValidationData(validationDataToExclude)
		if err := lt.IncorporateTimestampValidationData(validationData, indent); err != nil {
			return err
		}

		// incorporate validation data for counter-signature timestamps within AnyValidationData element
		counterSigTstValidationData := validationDataContainer.ValidationDataForCounterSignatureTimestampsForInclusion(signature)
		counterSigTstValidationData.ExcludeValidationData(validationData)
		counterSigTstValidationData.ExcludeValidationData(validationDataToExclude)
		return lt.IncorporateAnyValidationData(counterSigTstValidationData, indent)

	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndAnyValidationData:
		validationData = validationDataContainer.ValidationDataForSignatureTimestampsForInclusion(signature)
		validationData.AddValidationData(
			validationDataContainer.ValidationDataForCounterSignatureTimestampsForInclusion(signature))
		validationData.ExcludeValidationData(validationDataToExclude)
		return lt.IncorporateAnyValidationData(validationData, indent)

	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationData,
		enumerations.ValidationDataEncapsulationStrategyAnyValidationDataOnly:
		// skip
		return nil

	default:
		return fmt.Errorf("The ValidationDataEncapsulationStrategy '%s' is not supported!",
			validationDataEncapsulationStrategy)
	}
}

// extendToLTLevelSignatures ports the private getExtendToLTLevelSignatures.
func (lt *XAdESLevelBaselineLT) extendToLTLevelSignatures(
	signatures []validation.AdvancedSignature) []validation.AdvancedSignature {
	toBeExtended := make([]validation.AdvancedSignature, 0)
	for _, signature := range signatures {
		if lt.ltLevelExtensionRequired(signature) {
			toBeExtended = append(toBeExtended, signature)
		}
	}
	return toBeExtended
}

// ltLevelExtensionRequired ports the private ltLevelExtensionRequired.
func (lt *XAdESLevelBaselineLT) ltLevelExtensionRequired(signature validation.AdvancedSignature) bool {
	return enumerations.SignatureLevelXAdESBaselineLT == lt.Params.SignatureLevel() || !signature.HasLTAProfile()
}
