// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESLevelBaselineLTA.java (DSS 6.5.RC1).
//
// This level embeds JAdESLevelBaselineLT and overrides its virtual ExtendSignatures; see the
// header of jades_level_baseline_t.go for how the Java override chain is expressed. The "super"
// call is the explicit lta.JAdESLevelBaselineLT.ExtendSignatures.
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

// JAdESLevelBaselineLTA creates an LTA-level of a JAdES signature.
type JAdESLevelBaselineLTA struct {
	JAdESLevelBaselineLT
}

// NewJAdESLevelBaselineLTA is the default constructor.
// Port of JAdESLevelBaselineLTA(CertificateVerifier).
func NewJAdESLevelBaselineLTA(certificateVerifier validation.CertificateVerifier) *JAdESLevelBaselineLTA {
	extension := &JAdESLevelBaselineLTA{}
	extension.InitJAdESLevelBaselineT(extension, certificateVerifier)
	return extension
}

// ExtendSignatures extends the signatures to the -LTA level.
// Port of the protected, overridden #extendSignatures(List, JAdESSignatureParameters).
func (lta *JAdESLevelBaselineLTA) ExtendSignatures(signatures []validation.AdvancedSignature,
	params *JAdESSignatureParameters) error {
	if err := lta.JAdESLevelBaselineLT.ExtendSignatures(signatures, params); err != nil {
		return err
	}

	signatureRequirementsChecker := lta.SignatureRequirementsChecker(params)
	signatureRequirementsChecker.AssertSignaturesValid(signatures)

	addTimestampValidationData := false

	for _, signature := range signatures {
		jadesSignature, ok := signature.(*JAdESSignature)
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
	var validationDataContainer *validation.ValidationDataContainer
	if addTimestampValidationData {
		container, err := lta.DocumentAnalyzer.GetValidationData(signatures)
		if err != nil {
			return err
		}
		validationDataContainer = container
	}

	for _, signature := range signatures {
		jadesSignature, ok := signature.(*JAdESSignature)
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
// AdvancedSignature, JAdESEtsiUHeader, JAdESSignatureParameters).
func (lta *JAdESLevelBaselineLTA) ltaIncorporateValidationDataForTimestamps(
	validationDataContainer *validation.ValidationDataContainer, signature validation.AdvancedSignature,
	etsiUHeader *JAdESEtsiUHeader,
	signatureParameters *JAdESSignatureParameters) (*validation.ValidationData, error) {
	var validationData *validation.ValidationData
	validationDataEncapsulationStrategy := signatureParameters.ValidationDataEncapsulationStrategy()
	switch validationDataEncapsulationStrategy {
	case enumerations.ValidationDataEncapsulationStrategy_CERTIFICATE_REVOCATION_VALUES_AND_TIMESTAMP_VALIDATION_DATA,
		enumerations.ValidationDataEncapsulationStrategy_CERTIFICATE_REVOCATION_VALUES_AND_TIMESTAMP_VALIDATION_DATA_LT_SEPARATED:
		validationData = validationDataContainer.AllValidationDataForSignatureForInclusion(signature)
		if err := lta.IncorporateTstValidationData(etsiUHeader, validationData,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents())); err != nil {
			return nil, err
		}
	case enumerations.ValidationDataEncapsulationStrategy_CERTIFICATE_REVOCATION_VALUES_AND_TIMESTAMP_VALIDATION_DATA_AND_ANY_VALIDATION_DATA:
		validationData = validationDataContainer.ValidationDataForSignatureTimestampsForInclusion(signature)
		if err := lta.IncorporateTstValidationData(etsiUHeader, validationData,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents())); err != nil {
			return nil, err
		}

	case enumerations.ValidationDataEncapsulationStrategy_CERTIFICATE_REVOCATION_VALUES_AND_ANY_VALIDATION_DATA,
		enumerations.ValidationDataEncapsulationStrategy_ANY_VALIDATION_DATA_ONLY:
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
// JAdESEtsiUHeader, JAdESSignatureParameters, ValidationData).
func (lta *JAdESLevelBaselineLTA) ltaIncorporateAnyValidationData(
	validationDataContainer *validation.ValidationDataContainer, signature validation.AdvancedSignature,
	etsiUHeader *JAdESEtsiUHeader, signatureParameters *JAdESSignatureParameters,
	validationDataToExclude *validation.ValidationData) error {
	var validationData *validation.ValidationData
	validationDataEncapsulationStrategy := signatureParameters.ValidationDataEncapsulationStrategy()
	switch validationDataEncapsulationStrategy {
	case enumerations.ValidationDataEncapsulationStrategy_CERTIFICATE_REVOCATION_VALUES_AND_TIMESTAMP_VALIDATION_DATA_AND_ANY_VALIDATION_DATA:
		validationData = validationDataContainer.ValidationDataForSignatureForInclusion(signature)
		validationData.AddValidationData(
			validationDataContainer.ValidationDataForCounterSignaturesForInclusion(signature))
		validationData.AddValidationData(
			validationDataContainer.ValidationDataForCounterSignatureTimestampsForInclusion(signature))
		validationData.ExcludeValidationData(validationDataToExclude)
		return lta.IncorporateAnyValidationData(etsiUHeader, validationData,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents()))

	case enumerations.ValidationDataEncapsulationStrategy_CERTIFICATE_REVOCATION_VALUES_AND_ANY_VALIDATION_DATA,
		enumerations.ValidationDataEncapsulationStrategy_ANY_VALIDATION_DATA_ONLY:
		validationData = validationDataContainer.AllValidationDataForSignatureForInclusion(signature)
		validationData.ExcludeValidationData(validationDataToExclude)
		// skip
		return lta.IncorporateAnyValidationData(etsiUHeader, validationData,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents()))

	case enumerations.ValidationDataEncapsulationStrategy_CERTIFICATE_REVOCATION_VALUES_AND_TIMESTAMP_VALIDATION_DATA,
		enumerations.ValidationDataEncapsulationStrategy_CERTIFICATE_REVOCATION_VALUES_AND_TIMESTAMP_VALIDATION_DATA_LT_SEPARATED:
		// skip
		return nil

	default:
		return fmt.Errorf("The ValidationDataEncapsulationStrategy '%s' is not supported!",
			validationDataEncapsulationStrategy)
	}
}

// incorporateArcTst ports the private incorporateArcTst.
func (lta *JAdESLevelBaselineLTA) incorporateArcTst(signature *JAdESSignature,
	etsiUHeader *JAdESEtsiUHeader, signatureParameters *JAdESSignatureParameters) error {
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
func (lta *JAdESLevelBaselineLTA) archiveTimestamp(jadesSignature *JAdESSignature,
	params *JAdESSignatureParameters) (*model.TimestampBinary, error) {
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
func (lta *JAdESLevelBaselineLTA) assertExtendSignatureToLTAPossible(jadesSignature *JAdESSignature,
	params *JAdESSignatureParameters) error {
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
func jadesLevelBaselineLTACheckArchiveTimestampParameters(params *JAdESSignatureParameters) error {
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
func jadesLevelBaselineLTAAssertDetachedDocumentsContainBinaries(params *JAdESSignatureParameters) error {
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
func jadesLevelBaselineLTACheckEtsiUContentUnicity(jadesSignature *JAdESSignature) error {
	etsiU := DSSJsonUtilsEtsiU(jadesSignature.Jws())
	if !DSSJsonUtilsCheckComponentsUnicity(etsiU) {
		return exception.NewIllegalInputException(
			"Unsupported 'etsiU' container structure! Extension is not possible.")
	}
	return nil
}

// Compile-time assertion that *JAdESLevelBaselineLTA satisfies the extension contract.
var _ JAdESLevelBaselineExtension = (*JAdESLevelBaselineLTA)(nil)
