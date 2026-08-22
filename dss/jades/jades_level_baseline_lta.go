// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESLevelBaselineLTA.java (DSS 6.5.RC1).
//
// This level embeds LevelBaselineLT and overrides its virtual ExtendSignatures; see the
// header of jades_level_baseline_t.go for how the Java override chain is expressed. The "super"
// call is the explicit lta.LevelBaselineLT.ExtendSignatures.
//
// # The archive time-stamp message imprint
//
// The 'arcTst' message imprint is NOT recomputed here: it is whatever
// JAdESTimestampSource#getArchiveTimestampData(DigestAlgorithm, canonicalizationMethod) answers
// for the signature being extended, exactly as upstream, so the ETSI TS 119 182-1 clause 5.3.6
// computation lives in one place (the VAL chunk's timestamp source) and cannot drift between the
// signing and the validating side. The canonicalization method threaded through is always the
// null/empty one: JAdESTimestampParameters#setCanonicalizationMethod is an unsupported operation
// upstream, so the field is never populated.
//
// Java's two overloaded private helpers (incorporateValidationDataForTimestamps, and the
// incorporateAnyValidationData whose signature differs from the inherited one) get distinct Go
// names, since Go has no overloading and the inherited LT methods keep the plain ones:
//
//	incorporateValidationDataForTimestamps(container, sig, header, params)          -> ltaIncorporateValidationDataForTimestamps
//	incorporateAnyValidationData(container, sig, header, params, toExclude)         -> ltaIncorporateAnyValidationData
//
// Every Java throw becomes a returned error.
package jades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// LevelBaselineLTA creates an LTA-level of a JAdES signature.
type LevelBaselineLTA struct {
	LevelBaselineLT
}

// NewJAdESLevelBaselineLTA is the default constructor.
// Port of JAdESLevelBaselineLTA(CertificateVerifier).
func NewJAdESLevelBaselineLTA(certificateVerifier validation.CertificateVerifier) *LevelBaselineLTA {
	extension := &LevelBaselineLTA{}
	extension.InitJAdESLevelBaselineT(extension, certificateVerifier)
	return extension
}

// ExtendSignatures extends the signatures to the -LTA level.
// Port of the protected, overridden #extendSignatures(List, JAdESSignatureParameters).
func (lta *LevelBaselineLTA) ExtendSignatures(signatures []validation.AdvancedSignature,
	params *SignatureParameters) error {
	if err := lta.LevelBaselineLT.ExtendSignatures(signatures, params); err != nil {
		return err
	}

	signatureRequirementsChecker := lta.SignatureRequirementsChecker(params)
	signatureRequirementsChecker.AssertSignaturesValid(signatures)

	addTimestampValidationData := false

	for _, signature := range signatures {
		jadesSignature, ok := signature.(*Signature)
		if !ok {
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		if err := lta.AssertEtsiUComponentsConsistent(jadesSignature.Jws(), params); err != nil {
			return err
		}
		if err := lta.assertExtendSignatureToLTAPossible(jadesSignature, params); err != nil {
			return err
		}

		if jadesSignature.HasLTAProfile() {
			addTimestampValidationData = true
		}
	}

	// Perform signature validation
	var validationDataContainer *validation.DataContainer
	if addTimestampValidationData {
		container, err := lta.DocumentAnalyzer.GetValidationData(signatures)
		if err != nil {
			return err
		}
		validationDataContainer = container
	}

	for _, signature := range signatures {
		jadesSignature, ok := signature.(*Signature)
		if !ok {
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		etsiUHeader := jadesSignature.EtsiUHeader()

		if jadesSignature.HasLTAProfile() && addTimestampValidationData {
			if err := lta.RemoveLastTimestampAndAnyValidationData(jadesSignature, etsiUHeader); err != nil {
				return err
			}

			includedValidationData, err := lta.ltaIncorporateValidationDataForTimestamps(
				validationDataContainer, signature, etsiUHeader, params)
			if err != nil {
				return err
			}
			if err := lta.ltaIncorporateAnyValidationData(validationDataContainer, signature,
				etsiUHeader, params, includedValidationData); err != nil {
				return err
			}
		}

		if err := lta.incorporateArcTst(jadesSignature, etsiUHeader, params); err != nil {
			return err
		}
	}
	return nil
}

// ltaIncorporateValidationDataForTimestamps incorporates the validation data for the signature
// timestamps validation, according to the chosen validation data encapsulation mechanism, and
// returns the incorporated validation data.
// Port of the private incorporateValidationDataForTimestamps(ValidationDataContainer,
// AdvancedSignature, EtsiUHeader, SignatureParameters).
func (lta *LevelBaselineLTA) ltaIncorporateValidationDataForTimestamps(
	validationDataContainer *validation.DataContainer, signature validation.AdvancedSignature,
	etsiUHeader *EtsiUHeader,
	signatureParameters *SignatureParameters) (*validation.Data, error) {
	var validationData *validation.Data
	validationDataEncapsulationStrategy := signatureParameters.ValidationDataEncapsulationStrategy()
	switch validationDataEncapsulationStrategy {
	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationData,
		enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataLTSeparated:
		validationData = validationDataContainer.AllValidationDataForSignatureForInclusion(signature)
		if err := lta.IncorporateTstValidationData(etsiUHeader, validationData,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents())); err != nil {
			return nil, err
		}
	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataAndAnyValidationData:
		validationData = validationDataContainer.ValidationDataForSignatureTimestampsForInclusion(signature)
		if err := lta.IncorporateTstValidationData(etsiUHeader, validationData,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents())); err != nil {
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

// ltaIncorporateAnyValidationData incorporates the validation data for the signature validation,
// according to the chosen validation data encapsulation mechanism.
// Port of the private incorporateAnyValidationData(ValidationDataContainer, AdvancedSignature,
// EtsiUHeader, SignatureParameters, Data).
func (lta *LevelBaselineLTA) ltaIncorporateAnyValidationData(
	validationDataContainer *validation.DataContainer, signature validation.AdvancedSignature,
	etsiUHeader *EtsiUHeader, signatureParameters *SignatureParameters,
	validationDataToExclude *validation.Data) error {
	var validationData *validation.Data
	validationDataEncapsulationStrategy := signatureParameters.ValidationDataEncapsulationStrategy()
	switch validationDataEncapsulationStrategy {
	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataAndAnyValidationData:
		validationData = validationDataContainer.ValidationDataForSignatureForInclusion(signature)
		validationData.AddValidationData(
			validationDataContainer.ValidationDataForCounterSignaturesForInclusion(signature))
		validationData.AddValidationData(
			validationDataContainer.ValidationDataForCounterSignatureTimestampsForInclusion(signature))
		validationData.ExcludeValidationData(validationDataToExclude)
		return lta.IncorporateAnyValidationData(etsiUHeader, validationData,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents()))

	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndAnyValidationData,
		enumerations.ValidationDataEncapsulationStrategyAnyValidationDataOnly:
		validationData = validationDataContainer.AllValidationDataForSignatureForInclusion(signature)
		validationData.ExcludeValidationData(validationDataToExclude)
		// skip
		return lta.IncorporateAnyValidationData(etsiUHeader, validationData,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents()))

	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationData,
		enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataLTSeparated:
		// skip
		return nil

	default:
		return fmt.Errorf("The ValidationDataEncapsulationStrategy '%s' is not supported!",
			validationDataEncapsulationStrategy)
	}
}

// incorporateArcTst ports the private incorporateArcTst.
func (lta *LevelBaselineLTA) incorporateArcTst(signature *Signature,
	etsiUHeader *EtsiUHeader, signatureParameters *SignatureParameters) error {
	timestampBinary, err := lta.archiveTimestamp(signature, signatureParameters)
	if err != nil {
		return err
	}
	arcTst, err := DSSJsonUtilsTstContainer([]*model.TimestampBinary{timestampBinary},
		signatureParameters.GetArchiveTimestampParameters().CanonicalizationMethod())
	if err != nil {
		return err
	}
	return etsiUHeader.AddComponent(JAdESHeaderParameterNamesArcTst, arcTst,
		utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents()))
}

// archiveTimestamp ports the private getArchiveTimestamp.
func (lta *LevelBaselineLTA) archiveTimestamp(jadesSignature *Signature,
	params *SignatureParameters) (*model.TimestampBinary, error) {
	archiveTimestampParameters := params.GetArchiveTimestampParameters()
	digestAlgorithmForTimestampRequest := archiveTimestampParameters.DigestAlgorithm()
	// TODO : Support canonicalization
	canonicalizationMethod := archiveTimestampParameters.CanonicalizationMethod()

	timestampSource, err := jadesLevelBaselineTTimestampSource(jadesSignature.TimestampSource())
	if err != nil {
		return nil, err
	}
	messageDigest := timestampSource.GetArchiveTimestampData(
		digestAlgorithmForTimestampRequest, canonicalizationMethod)
	return lta.TspSource.TimeStampResponse(digestAlgorithmForTimestampRequest, messageDigest.Value())
}

// assertExtendSignatureToLTAPossible checks that the extension is possible.
// Port of the private assertExtendSignatureToLTAPossible.
func (lta *LevelBaselineLTA) assertExtendSignatureToLTAPossible(jadesSignature *Signature,
	params *SignatureParameters) error {
	if err := jadesLevelBaselineLTACheckArchiveTimestampParameters(params); err != nil {
		return err
	}
	if err := jadesLevelBaselineLTAAssertDetachedDocumentsContainBinaries(params); err != nil {
		return err
	}
	return jadesLevelBaselineLTACheckEtsiUContentUnicity(jadesSignature)
}

// jadesLevelBaselineLTACheckArchiveTimestampParameters ports the private
// checkArchiveTimestampParameters.
func jadesLevelBaselineLTACheckArchiveTimestampParameters(params *SignatureParameters) error {
	archiveTimestampParameters := params.GetArchiveTimestampParameters()
	if !utils.IsTrue(params.IsBase64UrlEncodedEtsiUComponents()) &&
		utils.IsStringEmpty(archiveTimestampParameters.CanonicalizationMethod()) {
		return exception.NewIllegalInputException(
			"Unable to extend JAdES-LTA level. Clear 'etsiU' incorporation requires a canonicalization method!")
	}
	return nil
}

// jadesLevelBaselineLTAAssertDetachedDocumentsContainBinaries ports the private
// assertDetachedDocumentsContainBinaries.
func jadesLevelBaselineLTAAssertDetachedDocumentsContainBinaries(params *SignatureParameters) error {
	detachedContents := params.DetachedContents()
	if utils.IsCollectionNotEmpty(detachedContents) {
		for _, detachedDocument := range detachedContents {
			if _, isDigestDocument := detachedDocument.(*model.DigestDocument); isDigestDocument {
				return fmt.Errorf("JAdES-LTA requires complete binaries of signed documents! " +
					"Extension with a DigestDocument is not possible.")
			}
		}
	}
	return nil
}

// jadesLevelBaselineLTACheckEtsiUContentUnicity ports the private checkEtsiUContentUnicity.
func jadesLevelBaselineLTACheckEtsiUContentUnicity(jadesSignature *Signature) error {
	etsiU := DSSJsonUtilsEtsiU(jadesSignature.Jws())
	if !DSSJsonUtilsCheckComponentsUnicity(etsiU) {
		return exception.NewIllegalInputException(
			"Unsupported 'etsiU' container structure! Extension is not possible.")
	}
	return nil
}

// Compile-time assertion that *LevelBaselineLTA satisfies the extension contract.
var _ LevelBaselineExtension = (*LevelBaselineLTA)(nil)
