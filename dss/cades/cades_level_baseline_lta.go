// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CAdESLevelBaselineLTA.java (DSS 6.5.RC1).
//
// Java's `extends CAdESLevelBaselineLT` becomes embedding, and `super.extendCMSSignatures(...)`
// an explicit call on the embedded base; see cades_signature_extension.go for the convention.
package cades

import (
	"encoding/asn1"

	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// LevelBaselineLTA holds the CAdES-A signature profiles; it supports the later, over time
// _extension_ of a signature with id-aa-ets-archiveTimestampV2 attributes as defined in
// ETSI TS 101 733 V1.8.1, clause 6.4.1.
//
// "If the certificate-values and revocation-values attributes are not present in the CAdES-BES
// or CAdES-EPES, then they shall be added to the electronic signature prior to computing the
// archive time-stamp token." is the reason we extend from the XL profile.
type LevelBaselineLTA struct {
	LevelBaselineLT
}

// NewCAdESLevelBaselineLTA is the default constructor, taking the TSPSource to request a
// timestamp and the CertificateVerifier. Port of
// LevelBaselineLTA(TSPSource, CertificateVerifier).
func NewCAdESLevelBaselineLTA(tspSource validation.TSPSource,
	certificateVerifier validation.CertificateVerifier) *LevelBaselineLTA {
	extension := &LevelBaselineLTA{}
	extension.InitCAdESSignatureExtension(extension, tspSource, certificateVerifier)
	return extension
}

// ExtendCMSSignaturesWithIds ports the overridden protected
// extendCMSSignatures(CMS, SignatureParameters, List<String>).
func (e *LevelBaselineLTA) ExtendCMSSignaturesWithIds(cmsToExtend *cms.CMS,
	parameters *SignatureParameters, signatureIdsToExtend []string) (*cms.CMS, error) {
	cmsToExtend, err := e.LevelBaselineLT.ExtendCMSSignaturesWithIds(cmsToExtend, parameters, signatureIdsToExtend)
	if err != nil {
		return nil, err
	}

	var newSignerInformationList []*cmscore.SignerInfo

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
	signatureRequirementsChecker.AssertExtendToLTALevelPossible(signaturesToExtend)

	// signature validity is checked within -LT augmentation

	for _, sig := range signatures {
		cadesSignature, ok := sig.(*Signature)
		if !ok {
			continue
		}
		signerInformation := cadesSignature.SignerInformation()
		newSignerInformation := signerInformation

		if cadesLTAContainsSignature(signaturesToExtend, cadesSignature) {
			unsignedAttributes := UtilsUnsignedAttributes(signerInformation)
			unsignedAttributes, err = e.addArchiveTimestampV3Attribute(
				cadesSignature, signerInformation, parameters, unsignedAttributes)
			if err != nil {
				return nil, err
			}
			newSignerInformation, err = cms.UtilsReplaceUnsignedAttributes(signerInformation, unsignedAttributes)
			if err != nil {
				return nil, err
			}

			/*
			 * Add a DigestAlgorithm used by an Archive TimeStamp to SignedData.digestAlgorithms
			 * set, when required.
			 * See ETSI EN 319 122-1, ch. "5.5.3 The archive-time-stamp-v3 attribute"
			 */
			timestampDigestAlgorithm := parameters.GetArchiveTimestampParameters().DigestAlgorithm()
			algorithmIdentifier, err := spi.DSSASN1UtilsAlgorithmIdentifierForDigest(timestampDigestAlgorithm)
			if err != nil {
				return nil, err
			}
			cmsToExtend, err = cms.UtilsPopulateDigestAlgorithmSet(cmsToExtend, []*spi.AlgorithmIdentifier{algorithmIdentifier})
			if err != nil {
				return nil, err
			}
		}
		newSignerInformationList = append(newSignerInformationList, newSignerInformation)
	}
	return e.ReplaceSigners(cmsToExtend, newSignerInformationList)
}

// addArchiveTimestampV3Attribute computes and appends the archive-time-stamp-v3 attribute.
//
// The input for the archive-time-stamp-v3's message imprint computation shall be the
// concatenation (in the order shown by the list below) of the signed data hash (see bullet 2
// below) and certain fields in their binary encoded form without any modification and including
// the tag, length and value octets:
//
//  1. The SignedData.encapContentInfo.eContentType.
//  2. The octets representing the hash of the signed data. The hash is computed on the same
//     content that was used for computing the hash value that is encapsulated within the
//     message-digest signed attribute of the CAdES signature being archive-time-stamped. The
//     hash algorithm applied shall be the same as the hash algorithm used for computing the
//     archive time-stamp's message imprint. The inclusion of the hash algorithm in the
//     SignedData.digestAlgorithms set is recommended.
//  3. Fields version, sid, digestAlgorithm, signedAttrs, signatureAlgorithm, and signature
//     within the SignedData.signerInfos's item corresponding to the signature being archive
//     time-stamped, in their order of appearance.
//  4. A single instance of ATSHashIndex type (created as specified in clause 6.4.2).
//
// Port of the private addArchiveTimestampV3Attribute.
func (e *LevelBaselineLTA) addArchiveTimestampV3Attribute(cadesSignature *Signature,
	signerInformation *cmscore.SignerInfo, parameters *SignatureParameters,
	unsignedAttributes cmscore.Attributes) (cmscore.Attributes, error) {

	timestampExtractor := NewCadesLevelBaselineLTATimestampExtractor(cadesSignature)
	timestampDigestAlgorithm := parameters.GetArchiveTimestampParameters().DigestAlgorithm()
	originalDocument, err := cadesSignature.OriginalDocument()
	if err != nil {
		return nil, err
	}

	atsHashIndexTableIdentifier := cadesLTAAtsHashIndexTableIdentifier(parameters)
	atsHashIndexAttribute, err := timestampExtractor.AtsHashIndex(
		signerInformation, timestampDigestAlgorithm, atsHashIndexTableIdentifier)
	if err != nil {
		return nil, err
	}

	timestampMessageDigest, err := timestampExtractor.ArchiveTimestampV3MessageImprint(
		signerInformation, atsHashIndexAttribute, originalDocument, timestampDigestAlgorithm)
	if err != nil {
		return nil, err
	}

	timeStampAttributeValue, err := e.TimeStampAttributeValue(
		timestampMessageDigest, timestampDigestAlgorithm, atsHashIndexAttribute)
	if err != nil {
		return nil, err
	}

	return UtilsAddAttribute(unsignedAttributes,
		spi.OIDIdAaEtsArchiveTimestampV3, timeStampAttributeValue), nil
}

// cadesLTAAtsHashIndexTableIdentifier ports the private
// getAtsHashIndexTableIdentifier(SignatureParameters).
func cadesLTAAtsHashIndexTableIdentifier(signatureParameters *SignatureParameters) asn1.ObjectIdentifier {
	if !signatureParameters.IsEn319122() {
		return spi.OIDIdAaATSHashIndex
	}
	return spi.OIDIdAaATSHashIndexV3
}

// compile-time assertion that the LTA profile satisfies the abstract base's contract.
var _ SignatureExtensionOverrides = (*LevelBaselineLTA)(nil)
