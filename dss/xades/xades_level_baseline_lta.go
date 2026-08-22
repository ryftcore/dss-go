// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XAdESLevelBaselineLTA.java (DSS 6.5.RC1).
//
// Java extends XAdESLevelBaselineLT and overrides extendSignatures(List); the Go port embeds the
// -LT level, and "super.extendSignatures(signatures)" is the explicit
// lta.XAdESLevelBaselineLT.ExtendSignatures call.
//
// Java's two private helpers here overload names already taken further up the hierarchy
// (incorporateValidationDataForTimestamps in XAdESLevelBaselineLT, incorporateAnyValidationData
// in XAdESLevelBaselineT). Go has no overloading and both would shadow the inherited member, so
// each carries the FromContainer suffix, naming the ValidationDataContainer argument that tells
// them apart from their same-named ancestors.
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// XAdESLevelBaselineLTA holds the level LTA aspects of XAdES.
type XAdESLevelBaselineLTA struct {
	XAdESLevelBaselineLT
}

// NewXAdESLevelBaselineLTA is the default constructor for XAdESLevelBaselineLTA.
// Port of XAdESLevelBaselineLTA(CertificateVerifier).
func NewXAdESLevelBaselineLTA(certificateVerifier validation.CertificateVerifier) *XAdESLevelBaselineLTA {
	extension := &XAdESLevelBaselineLTA{}
	extension.InitXAdESLevelBaselineLTA(extension, certificateVerifier)
	return extension
}

// InitXAdESLevelBaselineLTA registers the concrete extension level with this base and forwards to
// the -LT level. Port of the super(certificateVerifier) call of
// XAdESLevelBaselineLTA(CertificateVerifier).
func (lta *XAdESLevelBaselineLTA) InitXAdESLevelBaselineLTA(self XAdESSignatureExtensionOverrides,
	certificateVerifier validation.CertificateVerifier) {
	lta.InitXAdESLevelBaselineLT(self, certificateVerifier)
}

// ExtendSignatures adds the ArchiveTimeStamp element which is an unsigned property
// qualifying the signature. The hash sent to the TSA (messageImprint) is computed on the
// XAdES-LT form of the electronic signature and the signed data objects. A XAdES-LTA form MAY
// contain several ArchiveTimeStamp elements.
// Port of the overridden protected #extendSignatures(List).
func (lta *XAdESLevelBaselineLTA) ExtendSignatures(signatures []validation.AdvancedSignature) error {
	if err := lta.XAdESLevelBaselineLT.ExtendSignatures(signatures); err != nil {
		return err
	}

	signatureRequirementsChecker := lta.SignatureRequirementsChecker()
	signatureRequirementsChecker.AssertExtendToLTALevelPossible(signatures)

	signatureRequirementsChecker.AssertSignaturesValid(signatures)

	addTimestampValidationData := false
	for _, signature := range signatures {
		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			// Java's (XAdESSignature) cast; a non-XAdES signature would raise a ClassCastException.
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		if _, err := lta.InitializeSignatureBuilder(xadesSignature); err != nil {
			return err
		}

		if lta.XadesSignature.HasLTAProfile() {
			addTimestampValidationData = true
		}
	}

	// Perform signature validation
	var validationDataContainer *validation.ValidationDataContainer
	if addTimestampValidationData {
		container, err := lta.DocumentAnalyzer.GetValidationData(signatures)
		if err != nil {
			return err
		}
		validationDataContainer = container
	}

	// Append LTA-level (+ ValidationData)
	for _, signature := range signatures {
		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			// Java's (XAdESSignature) cast; a non-XAdES signature would raise a ClassCastException.
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		if _, err := lta.InitializeSignatureBuilder(xadesSignature); err != nil {
			return err
		}

		if err := lta.assertExtendSignatureToLTAPossible(); err != nil {
			return err
		}

		levelLTUnsignedProperties := lta.UnsignedSignaturePropertiesDom.Clone(true)

		if lta.XadesSignature.HasLTAProfile() && addTimestampValidationData {
			indent, err := lta.RemoveLastTimestampAndAnyValidationData()
			if err != nil {
				return err
			}
			includedValidationData, err := lta.incorporateValidationDataForTimestampsFromContainer(
				validationDataContainer, signature, indent)
			if err != nil {
				return err
			}
			if err := lta.incorporateAnyValidationDataFromContainer(validationDataContainer, signature,
				indent, includedValidationData); err != nil {
				return err
			}
		}

		if err := lta.IncorporateArchiveTimestamp(); err != nil {
			return err
		}
		indented, err := lta.IndentIfPrettyPrint(lta.UnsignedSignaturePropertiesDom, levelLTUnsignedProperties)
		if err != nil {
			return err
		}
		lta.UnsignedSignaturePropertiesDom = indented
	}
	return nil
}

// incorporateValidationDataForTimestampsFromContainer incorporates the validation data for the
// signature timestamps validation, according to the chosen validation data encapsulation
// mechanism, and returns the validation data incorporated within the TimeStampValidationData
// element. Port of the private incorporateValidationDataForTimestamps.
func (lta *XAdESLevelBaselineLTA) incorporateValidationDataForTimestampsFromContainer(
	validationDataContainer *validation.ValidationDataContainer,
	signature validation.AdvancedSignature, indent string) (*validation.ValidationData, error) {
	var validationData *validation.ValidationData
	validationDataEncapsulationStrategy := lta.ValidationDataEncapsulationStrategy()
	switch validationDataEncapsulationStrategy {
	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationData,
		enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataLTSeparated:
		validationData = validationDataContainer.AllValidationDataForSignatureForInclusion(signature)
		if err := lta.IncorporateTimestampValidationData(validationData, indent); err != nil {
			return nil, err
		}

	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataAndAnyValidationData:
		validationData = validationDataContainer.ValidationDataForSignatureTimestampsForInclusion(signature)
		if err := lta.IncorporateTimestampValidationData(validationData, indent); err != nil {
			return nil, err
		}

	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndAnyValidationData,
		enumerations.ValidationDataEncapsulationStrategyAnyValidationDataOnly:
		validationData = validation.NewValidationData()

	default:
		return nil, fmt.Errorf("The ValidationDataEncapsulationStrategy '%s' is not supported!",
			validationDataEncapsulationStrategy)
	}
	return validationData, nil
}

// incorporateAnyValidationDataFromContainer incorporates the validation data for the signature
// validation, according to the chosen validation data encapsulation mechanism, excluding
// validationDataToExclude to avoid duplicates.
// Port of the private incorporateAnyValidationData(ValidationDataContainer, ...).
func (lta *XAdESLevelBaselineLTA) incorporateAnyValidationDataFromContainer(
	validationDataContainer *validation.ValidationDataContainer,
	signature validation.AdvancedSignature, indent string,
	validationDataToExclude *validation.ValidationData) error {
	var validationData *validation.ValidationData
	validationDataEncapsulationStrategy := lta.ValidationDataEncapsulationStrategy()
	switch validationDataEncapsulationStrategy {
	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataAndAnyValidationData:
		validationData = validationDataContainer.ValidationDataForSignatureForInclusion(signature)
		validationData.AddValidationData(
			validationDataContainer.ValidationDataForCounterSignaturesForInclusion(signature))
		validationData.AddValidationData(
			validationDataContainer.ValidationDataForCounterSignatureTimestampsForInclusion(signature))
		validationData.ExcludeValidationData(validationDataToExclude)
		return lta.IncorporateAnyValidationData(validationData, indent)

	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndAnyValidationData,
		enumerations.ValidationDataEncapsulationStrategyAnyValidationDataOnly:
		validationData = validationDataContainer.AllValidationDataForSignatureForInclusion(signature)
		validationData.ExcludeValidationData(validationDataToExclude)
		return lta.IncorporateAnyValidationData(validationData, indent)

	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationData,
		enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataLTSeparated:
		// skip
		return nil

	default:
		return fmt.Errorf("The ValidationDataEncapsulationStrategy '%s' is not supported!",
			validationDataEncapsulationStrategy)
	}
}

// assertExtendSignatureToLTAPossible ports the private assertExtendSignatureToLTAPossible.
func (lta *XAdESLevelBaselineLTA) assertExtendSignatureToLTAPossible() error {
	if enumerations.SignatureLevelXAdESBaselineLTA == lta.Params.SignatureLevel() {
		return lta.AssertDetachedDocumentsContainBinaries()
	}
	return nil
}
