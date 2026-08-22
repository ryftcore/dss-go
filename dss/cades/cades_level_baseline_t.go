// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CAdESLevelBaselineT.java (DSS 6.5.RC1).
//
// Java extends the abstract CAdESSignatureExtension and overrides its abstract
// extendCMSSignatures(CMS, CAdESSignatureParameters, List<String>). Go has neither inheritance
// nor method overriding, so - as PORTING.md prescribes and document/abstract_document_extender.go
// already demonstrates - the base is embedded and the concrete extension registers itself with it
// through InitCAdESSignatureExtension(self, ...), satisfying CAdESSignatureExtensionOverrides.
// Java's overloaded extendCMSSignatures cannot share one Go name; the abstract one is
// ExtendCMSSignaturesWithIds, see cades_signature_extension.go for the whole mapping.
package cades

import (
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// CAdESLevelBaselineT holds the CAdES-T signature profile; it supports the inclusion of the
// mandatory unsigned id-aa-signatureTimeStampToken attribute as specified in ETSI TS 101 733
// V1.8.1, clause 6.1.1.
type CAdESLevelBaselineT struct {
	CAdESSignatureExtension
}

// NewCAdESLevelBaselineT is the default constructor with a CertificateVerifier.
// Port of CAdESLevelBaselineT(TSPSource, CertificateVerifier).
func NewCAdESLevelBaselineT(tspSource validation.TSPSource,
	certificateVerifier validation.CertificateVerifier) *CAdESLevelBaselineT {
	extension := &CAdESLevelBaselineT{}
	extension.InitCAdESSignatureExtension(extension, tspSource, certificateVerifier)
	return extension
}

// ExtendCMSSignaturesWithIds extends the signatures of the CMS whose ids are listed in
// signatureIdsToExtend with a signature time-stamp.
// Port of the overridden protected
// extendCMSSignatures(CMS, CAdESSignatureParameters, List<String>).
func (t *CAdESLevelBaselineT) ExtendCMSSignaturesWithIds(cmsToExtend *cms.CMS,
	parameters *CAdESSignatureParameters, signatureIdsToExtend []string) (*cms.CMS, error) {
	newSignerInformationList := make([]*cmscore.SignerInfo, 0)

	documentAnalyzer, err := t.DocumentAnalyzer(cmsToExtend, parameters)
	if err != nil {
		return nil, err
	}
	signatures := documentAnalyzer.Signatures()
	if utils.IsCollectionEmpty(signatures) {
		return nil, exception.NewIllegalInputException("No signatures found to be extended!")
	}

	signaturesToExtend := t.ExtendToTLevelSignatures(signatures, signatureIdsToExtend, parameters)
	if utils.IsCollectionEmpty(signaturesToExtend) {
		return cmsToExtend, nil
	}

	signatureRequirementsChecker := t.SignatureRequirementsChecker(parameters)
	signatureRequirementsChecker.AssertExtendToTLevelPossible(signaturesToExtend)

	signatureRequirementsChecker.AssertSignaturesValid(signaturesToExtend)
	signatureRequirementsChecker.AssertSigningCertificatesAreValid(signaturesToExtend)

	for _, signature := range signatures {
		cadesSignature, ok := signature.(*CAdESSignature)
		if !ok {
			// Java's cast; a CMSDocumentAnalyzer only ever yields CAdESSignatures.
			continue
		}
		signerInformation := cadesSignature.SignerInformation()
		newSignerInformation := signerInformation
		if cadesLevelBaselineTContains(signaturesToExtend, signature) {
			if newSignerInformation, err = t.extendSignerInformation(signerInformation, parameters); err != nil {
				return nil, err
			}
		}
		newSignerInformationList = append(newSignerInformationList, newSignerInformation)
	}

	return t.ReplaceSigners(cmsToExtend, newSignerInformationList)
}

// extendSignerInformation ports the private extendSignerInformation.
func (t *CAdESLevelBaselineT) extendSignerInformation(signerInformation *cmscore.SignerInfo,
	parameters *CAdESSignatureParameters) (*cmscore.SignerInfo, error) {
	unsignedAttributes := CAdESUtilsUnsignedAttributes(signerInformation)
	unsignedAttributes, err := t.addSignatureTimestampAttribute(signerInformation, unsignedAttributes, parameters)
	if err != nil {
		return nil, err
	}
	return cms.CMSUtilsReplaceUnsignedAttributes(signerInformation, unsignedAttributes)
}

// SignatureRequirementsChecker instantiates a SignatureRequirementsChecker.
// Port of the protected #getSignatureRequirementsChecker.
func (t *CAdESLevelBaselineT) SignatureRequirementsChecker(
	parameters *CAdESSignatureParameters) *document.SignatureRequirementsChecker[*CAdESTimestampParameters] {
	return document.NewSignatureRequirementsChecker[*CAdESTimestampParameters](t.CertificateVerifier,
		&parameters.AbstractSignatureParameters)
}

// ExtendToTLevelSignatures returns the signatures of the list that are named in
// signatureIdsToExtend and still need a T-level extension.
// Port of the private getExtendToTLevelSignatures; exported because CAdESLevelBaselineLT and
// -LTA embed this type, and Go has no visibility level between a type and the types embedding it.
func (t *CAdESLevelBaselineT) ExtendToTLevelSignatures(signatures []validation.AdvancedSignature,
	signatureIdsToExtend []string, parameters *CAdESSignatureParameters) []validation.AdvancedSignature {
	toBeExtended := make([]validation.AdvancedSignature, 0)
	for _, signature := range signatures {
		if cadesLevelBaselineTContainsID(signatureIdsToExtend, signature.ID()) &&
			cadesLevelBaselineTLevelExtensionRequired(signature, parameters) {
			toBeExtended = append(toBeExtended, signature)
		}
	}
	return toBeExtended
}

// cadesLevelBaselineTLevelExtensionRequired ports the private tLevelExtensionRequired.
func cadesLevelBaselineTLevelExtensionRequired(cadesSignature validation.AdvancedSignature,
	parameters *CAdESSignatureParameters) bool {
	return enumerations.SignatureLevelCAdESBaselineT == parameters.SignatureLevel() || !cadesSignature.HasTProfile()
}

// addSignatureTimestampAttribute ports the private addSignatureTimestampAttribute.
func (t *CAdESLevelBaselineT) addSignatureTimestampAttribute(signerInformation *cmscore.SignerInfo,
	unsignedAttributes cmscore.Attributes, parameters *CAdESSignatureParameters) (cmscore.Attributes, error) {
	timestampDigestAlgorithm := parameters.GetSignatureTimestampParameters().DigestAlgorithm()
	digestValue, err := spi.DSSUtilsDigest(timestampDigestAlgorithm, signerInformation.Signature)
	if err != nil {
		return nil, err
	}
	messageDigest := model.NewDSSMessageDigestWithValue(timestampDigestAlgorithm, digestValue)
	signatureTimeStamp, err := t.TimeStampAttributeValue(messageDigest, timestampDigestAlgorithm)
	if err != nil {
		return nil, err
	}
	return CAdESUtilsAddAttribute(unsignedAttributes, OIDIdAaSignatureTimeStampToken, signatureTimeStamp), nil
}

// cadesLevelBaselineTContains ports List#contains for a signature list, whose Java semantics are
// AdvancedSignature identity (DefaultAdvancedSignature does not override equals).
func cadesLevelBaselineTContains(signatures []validation.AdvancedSignature, signature validation.AdvancedSignature) bool {
	for _, candidate := range signatures {
		if candidate == signature {
			return true
		}
	}
	return false
}

// cadesLevelBaselineTContainsID ports List<String>#contains.
func cadesLevelBaselineTContainsID(ids []string, id string) bool {
	for _, candidate := range ids {
		if candidate == id {
			return true
		}
	}
	return false
}
