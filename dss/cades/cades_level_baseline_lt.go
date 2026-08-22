// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CAdESLevelBaselineLT.java (DSS 6.5.RC1).
//
// Java's `extends CAdESLevelBaselineT` becomes embedding, and the `super.extendCMSSignatures(...)`
// call becomes an explicit call on the embedded base. The subclass registers itself with
// InitCAdESSignatureExtension so that the abstract base dispatches into
// LevelBaselineLT.ExtendCMSSignaturesWithIds rather than into LevelBaselineT's.
//
// slf4j logging is dropped (PORTING.md).
package cades

import (
	"bytes"

	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// LevelBaselineLT holds the CAdES-LT signature profiles.
type LevelBaselineLT struct {
	LevelBaselineT
}

// NewCAdESLevelBaselineLT is the default constructor, taking the TSPSource for a timestamp
// creation and the CertificateVerifier. Port of
// LevelBaselineLT(TSPSource, CertificateVerifier).
func NewCAdESLevelBaselineLT(tspSource validation.TSPSource,
	certificateVerifier validation.CertificateVerifier) *LevelBaselineLT {
	extension := &LevelBaselineLT{}
	extension.InitCAdESSignatureExtension(extension, tspSource, certificateVerifier)
	return extension
}

// ExtendCMSSignaturesWithIds ports the overridden protected
// extendCMSSignatures(CMS, SignatureParameters, List<String>).
func (e *LevelBaselineLT) ExtendCMSSignaturesWithIds(cmsToExtend *cms.CMS,
	parameters *SignatureParameters, signatureIdsToExtend []string) (*cms.CMS, error) {
	cmsToExtend, err := e.LevelBaselineT.ExtendCMSSignaturesWithIds(cmsToExtend, parameters, signatureIdsToExtend)
	if err != nil {
		return nil, err
	}

	documentAnalyzer, err := e.DocumentAnalyzer(cmsToExtend, parameters)
	if err != nil {
		return nil, err
	}
	signatures := documentAnalyzer.Signatures()

	signaturesToExtend := e.ExtendToLTLevelSignatures(signatures, signatureIdsToExtend)
	if utils.IsCollectionEmpty(signaturesToExtend) {
		return cmsToExtend, nil
	}

	signatureRequirementsChecker := e.SignatureRequirementsChecker(parameters)
	if enumerations.SignatureLevelCAdESBaselineLT == parameters.SignatureLevel() {
		signatureRequirementsChecker.AssertExtendToLTLevelPossible(signaturesToExtend)
	}

	signatureRequirementsChecker.AssertSignaturesValid(signaturesToExtend)
	signatureRequirementsChecker.AssertCertificateChainValidForLTLevel(signaturesToExtend)

	// Perform signatures validation
	validationDataContainer, err := documentAnalyzer.GetValidationData(signaturesToExtend)
	if err != nil {
		return nil, err
	}

	/*
	 * ETSI EN 319 122-1 V1.1.1 (2016-04), chapter "5.5.3 The archive-time-stamp-v3 attribute":
	 *
	 * The present document specifies two strategies for the inclusion of validation data,
	 * depending on whether attributes for long term availability, as defined in different
	 * versions of ETSI TS 101 733 [1], have already been added to the SignedData:
	 */
	if e.IncludesATSv2(cmsToExtend) {
		if err := cms.UtilsAssertATSv2AugmentationSupported(); err != nil {
			return nil, err
		}
		/*
		 * - If an ATSv2, or other earlier form of archive time-stamp or a long-term-validation
		 *   attribute, is present in any SignerInfo of the root SignedData then the root
		 *   SignedData.certificates and SignedData.crls contents shall not be modified. The new
		 *   validation material shall be provided within the TimeStampToken of the latest archive
		 *   time-stamp (which can be an ATSv2 as defined in ETSI TS 101 733 [1], or an ATSv3) or
		 *   within the latest long-term-validation attribute (defined in ETSI TS 101 733 [1])
		 *   already contained in the SignerInfo ...
		 */
		newSignerInformationList := make([]*cmscore.SignerInfo, 0, len(signatures))
		for _, sig := range signatures {
			cadesSignature, ok := sig.(*Signature)
			if !ok {
				continue
			}
			signerInformation := cadesSignature.SignerInformation()
			newSignerInformation := signerInformation
			if cadesLTAContainsSignature(signaturesToExtend, cadesSignature) {
				validationData := validationDataContainer.AllValidationDataForSignatureForInclusion(cadesSignature)
				newSignerInformation, err = e.extendSignerInformation(signerInformation, validationData)
				if err != nil {
					return nil, err
				}
			}
			newSignerInformationList = append(newSignerInformationList, newSignerInformation)
		}
		cmsToExtend, err = e.ReplaceSigners(cmsToExtend, newSignerInformationList)
		if err != nil {
			return nil, err
		}

	} else {
		/*
		 * - If none of ATSv2 attributes (see clause A.2.4), or an earlier form of archive
		 *   time-stamp as defined in ETSI TS 101 733 [1] or long-term-validation (see clause
		 *   A.2.5) attributes is already present in any SignerInfo of the root SignedData, then
		 *   the new validation material shall be included within the root
		 *   SignedData.certificates, or SignedData.crls as applicable.
		 */
		allValidationData := validationDataContainer.AllValidationData()
		for _, sig := range signaturesToExtend {
			allValidationData.ExcludeCertificateTokens(sig.CertificateSource().Certificates())
			allValidationData.ExcludeCRLTokens(cadesLTARevocationIdentifiers(sig.CRLSource().AllRevocationBinaries()))
			allValidationData.ExcludeOCSPTokens(cadesLTARevocationIdentifiers(sig.OCSPSource().AllRevocationBinaries()))
		}

		cmsToExtend, err = e.extendWithValidationData(cmsToExtend, allValidationData)
		if err != nil {
			return nil, err
		}
	}

	return cmsToExtend, nil
}

// extendSignerInformation ports the private extendSignerInformation(SignerInformation, ValidationData).
func (e *LevelBaselineLT) extendSignerInformation(signerInformation *cmscore.SignerInfo,
	validationData *validation.Data) (*cmscore.SignerInfo, error) {
	unsignedAttributes := UtilsUnsignedAttributes(signerInformation)
	unsignedAttributes, err := e.addValidationData(unsignedAttributes, validationData)
	if err != nil {
		return nil, err
	}
	return cms.UtilsReplaceUnsignedAttributes(signerInformation, unsignedAttributes)
}

// addValidationData ports the private addValidationData(AttributeTable, ValidationData).
func (e *LevelBaselineLT) addValidationData(unsignedAttributes cmscore.Attributes,
	validationData *validation.Data) (cmscore.Attributes, error) {
	timestampTokenToExtend, err := e.lastArchiveTimestamp(unsignedAttributes)
	if err != nil {
		return nil, err
	}
	if timestampTokenToExtend != nil {
		timestampCMS, err := cms.UtilsToCMS(timestampTokenToExtend)
		if err != nil {
			return nil, err
		}
		extendedTimestampCMS, err := e.extendWithValidationData(timestampCMS, validationData)
		if err != nil {
			return nil, err
		}

		unsignedAttributes = e.replaceTimeStampAttribute(unsignedAttributes, timestampCMS, extendedTimestampCMS)
	}
	return unsignedAttributes, nil
}

// lastArchiveTimestamp ports the private getLastArchiveTimestamp(AttributeTable).
func (e *LevelBaselineLT) lastArchiveTimestamp(unsignedAttributes cmscore.Attributes) (*cmscore.TimeStampToken, error) {
	var lastTimeStampToken *cmscore.TimeStampToken
	comparator := NewTimeStampTokenProductionComparator()
	timeStampTokens, err := UtilsFindArchiveTimeStampTokens(unsignedAttributes)
	if err != nil {
		return nil, err
	}
	for _, timeStampToken := range timeStampTokens {
		if lastTimeStampToken == nil || comparator.After(timeStampToken, lastTimeStampToken) {
			lastTimeStampToken = timeStampToken
		}
	}
	return lastTimeStampToken, nil
}

// replaceTimeStampAttribute returns a new attribute table with attributeToReplace replaced by
// attributeToAdd. Port of the private
// replaceTimeStampAttribute(AttributeTable, CMS, CMS).
func (e *LevelBaselineLT) replaceTimeStampAttribute(attributeTable cmscore.Attributes,
	attributeToReplace, attributeToAdd *cms.CMS) cmscore.Attributes {
	newAsn1EncodableVector := make(cmscore.Attributes, 0, len(attributeTable))
	for _, attribute := range cadesLTAAttributeTableOrder(attributeTable) {
		newAttribute := attribute
		if UtilsIsArchiveTimeStampToken(attribute) {
			// ContentInfo binaries have to be compared, therefore CMS creation is required
			attributeValue, err := UtilsEncodedValue(attribute)
			if err == nil && bytes.Equal(attributeToReplace.DEREncoded(), attributeValue) {
				newAttribute = cmscore.NewAttribute(attribute.Type, attributeToAdd.DEREncoded())
			}
			// Upstream logs "Unable to build a CMS object from an unsigned attribute. Reason :
			// {}" and continues with the original object, because it would not be possible to
			// extend the attribute anyway.
		}
		newAsn1EncodableVector = append(newAsn1EncodableVector, newAttribute)
	}
	return newAsn1EncodableVector
}

// extendWithValidationData extends the cms with the LT-level (validation data).
// Port of the private extendWithValidationData(CMS, ValidationData).
func (e *LevelBaselineLT) extendWithValidationData(cmsToExtend *cms.CMS,
	validationDataForInclusion *validation.Data) (*cms.CMS, error) {
	cmsBuilder := cms.NewCMSBuilder().SetOriginalCMS(cmsToExtend)
	return cmsBuilder.ExtendCMSSignedData(validationDataForInclusion.CertificateTokens(),
		validationDataForInclusion.CrlTokens(), validationDataForInclusion.OcspTokens())
}

// IncludesATSv2 verifies whether the CMS contains an ATSTv2.
// Port of the protected includesATSv2(CMS).
func (e *LevelBaselineLT) IncludesATSv2(cmsToCheck *cms.CMS) bool {
	for _, signerInformation := range cmsToCheck.SignerInfos() {
		if UtilsContainsATSTv2(signerInformation) {
			return true
		}
	}
	return false
}

// ExtendToLTLevelSignatures returns the signatures to be extended according to the list of
// signatureIdsToExtend. Port of the protected
// getExtendToLTLevelSignatures(List<AdvancedSignature>, List<String>).
func (e *LevelBaselineLT) ExtendToLTLevelSignatures(signatures []validation.AdvancedSignature,
	signatureIdsToExtend []string) []validation.AdvancedSignature {
	toBeExtended := make([]validation.AdvancedSignature, 0, len(signatures))
	for _, sig := range signatures {
		if cadesLTAContainsID(signatureIdsToExtend, sig.ID()) {
			toBeExtended = append(toBeExtended, sig)
		}
	}
	return toBeExtended
}

// cadesLTARevocationIdentifiers narrows a list of revocation binaries to the Identifiers
// ValidationData#excludeCRLTokens / #excludeOCSPTokens take, which is what Java's
// `Collection<? extends Identifier>` parameter accepts directly.
func cadesLTARevocationIdentifiers[R revocation.Revocation](
	binaries []spi.EncapsulatedRevocationTokenIdentifier[R]) []model.Identifier {
	identifiers := make([]model.Identifier, 0, len(binaries))
	for _, binary := range binaries {
		identifiers = append(identifiers, binary.DSSID())
	}
	return identifiers
}

// cadesLTAContainsID is List<String>#contains.
func cadesLTAContainsID(ids []string, id string) bool {
	for _, candidate := range ids {
		if candidate == id {
			return true
		}
	}
	return false
}

// cadesLTAContainsSignature is List<AdvancedSignature>#contains, which AdvancedSignature answers
// by identity: neither the Java class nor this port overrides equals.
func cadesLTAContainsSignature(signatures []validation.AdvancedSignature, signature validation.AdvancedSignature) bool {
	for _, candidate := range signatures {
		if candidate == signature {
			return true
		}
	}
	return false
}

// compile-time assertion that the LT profile satisfies the abstract base's contract.
var _ SignatureExtensionOverrides = (*LevelBaselineLT)(nil)
