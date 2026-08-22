// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/signature/JAdESLevelBaselineLT.java (DSS 6.5.RC1).
//
// This level embeds JAdESLevelBaselineT and overrides its virtual ExtendSignatures; see the
// header of jades_level_baseline_t.go for how the Java override chain is expressed. The "super"
// call is the explicit lt.JAdESLevelBaselineT.ExtendSignatures.
//
// Every JSON object this level builds is HashMap-ordered upstream and therefore built here with
// NewJsonObject(), whose wrapped object reproduces java.util.HashMap's bucket order: the 'val'
// and 'x509Cert' wrappers are json_simple JSONObjects (a HashMap subclass), and 'rVals' and
// 'tstVD' come from DSS's own no-argument `new JsonObject()`, which wraps a bare HashMap. The
// wrappers carry a single member each, but 'rVals' (crlVals, ocspVals) and 'tstVD' (xVals,
// rVals) carry two, and an archive time-stamp covers exactly those bytes - so the ordering is
// not a detail that can be normalized away here.
//
// Java's Set<CertificateToken>/Set<CRLToken>/Set<OCSPToken> become the slices the frozen
// ValidationData answers, which already fixes their iteration order.
//
// Every Java throw becomes a returned error.
package jades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESLevelBaselineLT creates an LT-level of a JAdES signature.
type JAdESLevelBaselineLT struct {
	JAdESLevelBaselineT
}

// NewJAdESLevelBaselineLT is the default constructor.
// Port of JAdESLevelBaselineLT(CertificateVerifier).
func NewJAdESLevelBaselineLT(certificateVerifier validation.CertificateVerifier) *JAdESLevelBaselineLT {
	extension := &JAdESLevelBaselineLT{}
	extension.InitJAdESLevelBaselineT(extension, certificateVerifier)
	return extension
}

// ExtendSignatures extends the signatures to the -LT level.
// Port of the protected, overridden #extendSignatures(List, JAdESSignatureParameters).
func (lt *JAdESLevelBaselineLT) ExtendSignatures(signatures []validation.AdvancedSignature,
	params *JAdESSignatureParameters) error {
	if err := lt.JAdESLevelBaselineT.ExtendSignatures(signatures, params); err != nil {
		return err
	}

	signaturesToExtend := jadesLevelBaselineLTExtendToLTLevelSignatures(signatures, params)
	if utils.IsCollectionEmpty(signaturesToExtend) {
		return nil
	}

	// Reset sources
	for _, signature := range signaturesToExtend {
		jadesSignature, ok := signature.(*JAdESSignature)
		if !ok {
			return fmt.Errorf("unexpected signature type %T", signature)
		}

		// Data sources can already be loaded in memory (force reload)
		jadesSignature.ResetCertificateSource()
		jadesSignature.ResetRevocationSources()
		jadesSignature.ResetTimestampSource()
	}

	signatureRequirementsChecker := lt.SignatureRequirementsChecker(params)
	if enumerations.SignatureLevelJAdESBaselineLT == params.SignatureLevel() {
		signatureRequirementsChecker.AssertExtendToLTLevelPossible(signaturesToExtend)
	}
	signatureRequirementsChecker.AssertSignaturesValid(signaturesToExtend)
	signatureRequirementsChecker.AssertCertificateChainValidForLTLevel(signaturesToExtend)

	// Perform signature validation
	validationDataContainer, err := lt.DocumentAnalyzer.GetValidationData(signatures)
	if err != nil {
		return err
	}

	// Append ValidationData
	for _, signature := range signaturesToExtend {
		jadesSignature, ok := signature.(*JAdESSignature)
		if !ok {
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		if jadesSignature.HasLTAProfile() {
			continue
		}

		if err := lt.AssertEtsiUComponentsConsistent(jadesSignature.Jws(), params); err != nil {
			return err
		}

		etsiUHeader := jadesSignature.EtsiUHeader()

		if err := lt.removeOldCertificateValues(jadesSignature, etsiUHeader); err != nil {
			return err
		}
		if err := lt.removeOldRevocationValues(jadesSignature, etsiUHeader); err != nil {
			return err
		}
		if err := lt.RemoveLastTimestampAndAnyValidationData(jadesSignature, etsiUHeader); err != nil {
			return err
		}

		includedValidationData, err := lt.incorporateValidationDataForSignature(validationDataContainer,
			signature, etsiUHeader, params)
		if err != nil {
			return err
		}
		if err := lt.incorporateValidationDataForTimestamps(validationDataContainer, signature,
			etsiUHeader, params, includedValidationData); err != nil {
			return err
		}
	}
	return nil
}

// removeOldCertificateValues ports the private removeOldCertificateValues.
func (lt *JAdESLevelBaselineLT) removeOldCertificateValues(jadesSignature *JAdESSignature,
	etsiUHeader *JAdESEtsiUHeader) error {
	if err := etsiUHeader.RemoveComponent(JAdESHeaderParameterNamesXVals); err != nil {
		return err
	}
	jadesSignature.ResetCertificateSource()
	return nil
}

// removeOldRevocationValues ports the private removeOldRevocationValues.
func (lt *JAdESLevelBaselineLT) removeOldRevocationValues(jadesSignature *JAdESSignature,
	etsiUHeader *JAdESEtsiUHeader) error {
	if err := etsiUHeader.RemoveComponent(JAdESHeaderParameterNamesRVals); err != nil {
		return err
	}
	jadesSignature.ResetRevocationSources()
	return nil
}

// RemoveLastTimestampAndAnyValidationData removes the 'tstVd' and 'anyValData' header parameters
// appearing at the end of the 'etsiU' unsigned property array.
// Port of the protected #removeLastTimestampAndAnyValidationData.
func (lt *JAdESLevelBaselineLT) RemoveLastTimestampAndAnyValidationData(jadesSignature *JAdESSignature,
	etsiUHeader *JAdESEtsiUHeader) error {
	resetSources := false
	for {
		removed, err := etsiUHeader.RemoveLastComponent(JAdESHeaderParameterNamesTstVD,
			JAdESHeaderParameterNamesAnyValData)
		if err != nil {
			return err
		}
		if !removed {
			break
		}
		resetSources = true
	}
	if resetSources {
		jadesSignature.ResetCertificateSource()
		jadesSignature.ResetRevocationSources()
	}
	return nil
}

// incorporateValidationDataForSignature incorporates the validation data for the signature
// validation, according to the chosen validation data encapsulation mechanism, and returns the
// incorporated validation data. Port of the private incorporateValidationDataForSignature.
func (lt *JAdESLevelBaselineLT) incorporateValidationDataForSignature(
	validationDataContainer *validation.ValidationDataContainer, signature validation.AdvancedSignature,
	etsiUHeader *JAdESEtsiUHeader,
	signatureParameters *JAdESSignatureParameters) (*validation.ValidationData, error) {
	var validationDataForInclusion *validation.ValidationData
	validationDataEncapsulationStrategy := signatureParameters.ValidationDataEncapsulationStrategy()
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
		if err := lt.IncorporateXVals(etsiUHeader, certificateValuesToAdd,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents())); err != nil {
			return nil, err
		}
		if err := lt.IncorporateRVals(etsiUHeader, crlsToAdd, ocspsToAdd,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents())); err != nil {
			return nil, err
		}

	case enumerations.ValidationDataEncapsulationStrategyAnyValidationDataOnly:
		if err := lt.IncorporateAnyValidationData(etsiUHeader, validationDataForInclusion,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents())); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("The ValidationDataEncapsulationStrategy '%s' is not supported!",
			validationDataEncapsulationStrategy)
	}
	return validationDataForInclusion, nil
}

// incorporateValidationDataForTimestamps incorporates the validation data for the signature
// timestamps validation, according to the chosen validation data encapsulation mechanism.
// Port of the private incorporateValidationDataForTimestamps.
func (lt *JAdESLevelBaselineLT) incorporateValidationDataForTimestamps(
	validationDataContainer *validation.ValidationDataContainer, signature validation.AdvancedSignature,
	etsiUHeader *JAdESEtsiUHeader, signatureParameters *JAdESSignatureParameters,
	validationDataToExclude *validation.ValidationData) error {
	var validationData *validation.ValidationData
	validationDataEncapsulationStrategy := signatureParameters.ValidationDataEncapsulationStrategy()
	switch validationDataEncapsulationStrategy {
	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataLTSeparated:
		validationData = validationDataContainer.ValidationDataForSignatureTimestampsForInclusion(signature)
		validationData.AddValidationData(
			validationDataContainer.ValidationDataForCounterSignatureTimestampsForInclusion(signature))
		validationData.ExcludeValidationData(validationDataToExclude)
		return lt.IncorporateTstValidationData(etsiUHeader, validationData,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents()))

	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationDataAndAnyValidationData:
		validationData = validationDataContainer.ValidationDataForSignatureTimestampsForInclusion(signature)
		validationData.ExcludeValidationData(validationDataToExclude)
		if err := lt.IncorporateTstValidationData(etsiUHeader, validationData,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents())); err != nil {
			return err
		}

		// incorporate validation data for counter-signature timestamps within
		// AnyValidationData element
		counterSigTstValidationData :=
			validationDataContainer.ValidationDataForCounterSignatureTimestampsForInclusion(signature)
		counterSigTstValidationData.ExcludeValidationData(validationData)
		counterSigTstValidationData.ExcludeValidationData(validationDataToExclude)
		return lt.IncorporateAnyValidationData(etsiUHeader, counterSigTstValidationData,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents()))

	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndAnyValidationData:
		validationData = validationDataContainer.ValidationDataForSignatureTimestampsForInclusion(signature)
		validationData.AddValidationData(
			validationDataContainer.ValidationDataForCounterSignatureTimestampsForInclusion(signature))
		validationData.ExcludeValidationData(validationDataToExclude)
		return lt.IncorporateAnyValidationData(etsiUHeader, validationData,
			utils.IsTrue(signatureParameters.IsBase64UrlEncodedEtsiUComponents()))

	case enumerations.ValidationDataEncapsulationStrategyCertificateRevocationValuesAndTimestampValidationData,
		enumerations.ValidationDataEncapsulationStrategyAnyValidationDataOnly:
		// skip
		return nil

	default:
		return fmt.Errorf("The ValidationDataEncapsulationStrategy '%s' is not supported!",
			validationDataEncapsulationStrategy)
	}
}

// XVals builds and returns the 'xVals' array. Port of the protected #getXVals.
func (lt *JAdESLevelBaselineLT) XVals(certificateValuesToAdd []*model.CertificateToken) []any {
	xValsArray := make([]any, 0, len(certificateValuesToAdd))
	for _, certificateToken := range certificateValuesToAdd {
		xValsArray = append(xValsArray, jadesLevelBaselineLTX509CertObject(certificateToken))
	}
	return xValsArray
}

// jadesLevelBaselineLTX509CertObject ports the private getX509CertObject.
func jadesLevelBaselineLTX509CertObject(certificateToken *model.CertificateToken) *JsonObject {
	pkiOb := NewJsonObject()
	pkiOb.Put(JAdESHeaderParameterNamesVal, utils.ToBase64(certificateToken.Encoded()))

	x509Cert := NewJsonObject()
	x509Cert.Put(JAdESHeaderParameterNamesX509Cert, pkiOb)
	return x509Cert
}

// IncorporateXVals incorporates the provided set of certificates into etsiUHeader.
// Port of the protected #incorporateXVals.
func (lt *JAdESLevelBaselineLT) IncorporateXVals(etsiUHeader *JAdESEtsiUHeader,
	certificateValuesToAdd []*model.CertificateToken, base64UrlEncoded bool) error {
	if utils.IsCollectionNotEmpty(certificateValuesToAdd) {
		xVals := lt.XVals(certificateValuesToAdd)
		return etsiUHeader.AddComponent(JAdESHeaderParameterNamesXVals, xVals, base64UrlEncoded)
	}
	return nil
}

// RVals builds and returns the 'rVals' object. Port of the protected #getRVals.
func (lt *JAdESLevelBaselineLT) RVals(crlsToAdd []*spi.CRLToken, ocspsToAdd []*spi.OCSPToken) *JsonObject {
	rValsObject := NewJsonObject()
	if utils.IsCollectionNotEmpty(crlsToAdd) {
		rValsObject.Put(JAdESHeaderParameterNamesCrlVals, jadesLevelBaselineLTCrlVals(crlsToAdd))
	}
	if utils.IsCollectionNotEmpty(ocspsToAdd) {
		rValsObject.Put(JAdESHeaderParameterNamesOcspVals, jadesLevelBaselineLTOcspVals(ocspsToAdd))
	}
	return rValsObject
}

// jadesLevelBaselineLTCrlVals ports the private getCrlVals.
func jadesLevelBaselineLTCrlVals(crlsToAdd []*spi.CRLToken) []any {
	array := make([]any, 0, len(crlsToAdd))
	for _, crlToken := range crlsToAdd {
		pkiOb := NewJsonObject()
		pkiOb.Put(JAdESHeaderParameterNamesVal, utils.ToBase64(crlToken.Encoded()))
		array = append(array, pkiOb)
	}
	return array
}

// jadesLevelBaselineLTOcspVals ports the private getOcspVals.
func jadesLevelBaselineLTOcspVals(ocspsToAdd []*spi.OCSPToken) []any {
	array := make([]any, 0, len(ocspsToAdd))
	for _, ocspToken := range ocspsToAdd {
		pkiOb := NewJsonObject()
		pkiOb.Put(JAdESHeaderParameterNamesVal, utils.ToBase64(ocspToken.Encoded()))
		array = append(array, pkiOb)
	}
	return array
}

// IncorporateRVals incorporates the provided revocation data into etsiUHeader.
// Port of the protected #incorporateRVals.
func (lt *JAdESLevelBaselineLT) IncorporateRVals(etsiUHeader *JAdESEtsiUHeader, crlsToAdd []*spi.CRLToken,
	ocspsToAdd []*spi.OCSPToken, base64UrlEncoded bool) error {
	if utils.IsCollectionNotEmpty(crlsToAdd) || utils.IsCollectionNotEmpty(ocspsToAdd) {
		rVals := lt.RVals(crlsToAdd, ocspsToAdd)
		return etsiUHeader.AddComponent(JAdESHeaderParameterNamesRVals, rVals, base64UrlEncoded)
	}
	return nil
}

// IncorporateTstValidationData incorporates the 'tstVD' dictionary in the signature.
// Port of the protected #incorporateTstValidationData.
func (lt *JAdESLevelBaselineLT) IncorporateTstValidationData(etsiUHeader *JAdESEtsiUHeader,
	validationDataForInclusion *validation.ValidationData, base64UrlEncoded bool) error {
	return lt.IncorporateValidationData(etsiUHeader, validationDataForInclusion,
		JAdESHeaderParameterNamesTstVD, base64UrlEncoded)
}

// IncorporateAnyValidationData incorporates the 'anyValData' dictionary in the signature.
// Port of the protected #incorporateAnyValidationData.
func (lt *JAdESLevelBaselineLT) IncorporateAnyValidationData(etsiUHeader *JAdESEtsiUHeader,
	validationDataForInclusion *validation.ValidationData, base64UrlEncoded bool) error {
	return lt.IncorporateValidationData(etsiUHeader, validationDataForInclusion,
		JAdESHeaderParameterNamesAnyValData, base64UrlEncoded)
}

// IncorporateValidationData incorporates the validation data container in the signature under the
// given header name. Port of the protected #incorporateValidationData.
func (lt *JAdESLevelBaselineLT) IncorporateValidationData(etsiUHeader *JAdESEtsiUHeader,
	validationDataForInclusion *validation.ValidationData, headerName string,
	base64UrlEncoded bool) error {
	if !validationDataForInclusion.IsEmpty() {
		tstVd := lt.tstVd(validationDataForInclusion)
		return etsiUHeader.AddComponent(headerName, tstVd, base64UrlEncoded)
	}
	return nil
}

// tstVd ports the private getTstVd.
func (lt *JAdESLevelBaselineLT) tstVd(validationDataForInclusion *validation.ValidationData) *JsonObject {
	certificateTokens := validationDataForInclusion.CertificateTokens()
	crlTokens := validationDataForInclusion.CrlTokens()
	ocspTokens := validationDataForInclusion.OcspTokens()

	tstVd := NewJsonObject()
	if utils.IsCollectionNotEmpty(certificateTokens) {
		tstVd.Put(JAdESHeaderParameterNamesXVals, lt.XVals(certificateTokens))
	}
	if utils.IsCollectionNotEmpty(crlTokens) || utils.IsCollectionNotEmpty(ocspTokens) {
		tstVd.Put(JAdESHeaderParameterNamesRVals, lt.RVals(crlTokens, ocspTokens))
	}
	return tstVd
}

// jadesLevelBaselineLTExtendToLTLevelSignatures ports the private getExtendToLTLevelSignatures.
func jadesLevelBaselineLTExtendToLTLevelSignatures(signatures []validation.AdvancedSignature,
	parameters *JAdESSignatureParameters) []validation.AdvancedSignature {
	toBeExtended := make([]validation.AdvancedSignature, 0)
	for _, signature := range signatures {
		if jadesLevelBaselineLTLtLevelExtensionRequired(signature, parameters) {
			toBeExtended = append(toBeExtended, signature)
		}
	}
	return toBeExtended
}

// jadesLevelBaselineLTLtLevelExtensionRequired ports the private ltLevelExtensionRequired.
func jadesLevelBaselineLTLtLevelExtensionRequired(signature validation.AdvancedSignature,
	parameters *JAdESSignatureParameters) bool {
	return enumerations.SignatureLevelJAdESBaselineLT == parameters.SignatureLevel() ||
		!signature.HasLTAProfile()
}

// Compile-time assertion that *JAdESLevelBaselineLT satisfies the extension contract.
var _ JAdESLevelBaselineExtension = (*JAdESLevelBaselineLT)(nil)
