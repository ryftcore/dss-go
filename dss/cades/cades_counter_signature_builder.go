// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CAdESCounterSignatureBuilder.java (DSS 6.5.RC1).
//
// BC-bridge note: upstream builds the counter-signature SignerInfo and merges it into the
// target's unsigned attributes with org.bouncycastle.cms.CMSSignedDataGenerator
// (#generateCounterSigners) and org.bouncycastle.cms.SignerInformation (#addCounterSigners,
// #getCounterSignatures, #toASN1Structure). None of those BC types exist on this port, so this
// file reproduces their observable contract directly against cmscore.SignerInfo/Attributes:
//   - generateCounterSigners(SignerInformation) ~ SignerInfoGenerator.Generate(nil) - exactly one
//     SignerInfo is produced per call here (DSS only ever registers one generator), so no
//     SignerInformationStore type is needed - a []*cmscore.SignerInfo of length 1 stands in for it.
//   - SignerInformation#addCounterSigners(SignerInformation, SignerInformationStore) appends the
//     new counter-signer(s) to any existing id-countersignature unsigned attribute (RFC 5652
//     5.3: the attribute value is a SET OF SignerInfo, so multiple counter-signers share one
//     attribute rather than getting one attribute apiece); cadesCounterSignaturesOf
//     (cades_signature.go, same package) already ports the read side
//     (SignerInformation#getCounterSignatures()).
//   - SignerInformation#toASN1Structure() ~ (*cmscore.SignerInfo).DER().
package cades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/signature/resources"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// CAdESCounterSignatureBuilder builds a CAdES counter signature.
type CAdESCounterSignatureBuilder struct {
	// certificateVerifier is the CertificateVerifier to use.
	certificateVerifier validation.CertificateVerifier

	// manifestFile is a signature-signed manifest. Used for ASiC.
	manifestFile *model.ManifestFile

	// resourcesHandlerBuilder is used to create data container objects such as an OutputStream or
	// a DSSDocument.
	resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder
}

// NewCAdESCounterSignatureBuilder is the default constructor.
func NewCAdESCounterSignatureBuilder(certificateVerifier validation.CertificateVerifier) *CAdESCounterSignatureBuilder {
	return &CAdESCounterSignatureBuilder{certificateVerifier: certificateVerifier}
}

// SetManifestFile sets a signed manifest file. NOTE: ASiC only. Port of #setManifestFile.
func (b *CAdESCounterSignatureBuilder) SetManifestFile(manifestFile *model.ManifestFile) {
	b.manifestFile = manifestFile
}

// SetResourcesHandlerBuilder sets a DSSResourcesHandlerBuilder to be used for operating with
// internal objects during the signature creation procedure. Port of
// #setResourcesHandlerBuilder.
func (b *CAdESCounterSignatureBuilder) SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) {
	b.resourcesHandlerBuilder = resourcesHandlerBuilder
}

// AddCounterSignature adds a counter signature to the provided CMS. Port of
// #addCounterSignature(CMS, CAdESCounterSignatureParameters, SignatureValue).
func (b *CAdESCounterSignatureBuilder) AddCounterSignature(originalCMS *cms.CMS, parameters *CAdESCounterSignatureParameters,
	signatureValue *model.SignatureValue) (model.DSSDocument, error) {

	updatedSignerInfo, err := b.getUpdatedSignerInformations(originalCMS, originalCMS.SignerInfos(), parameters, signatureValue, nil)
	if err != nil {
		return nil, err
	}

	updatedCMS, err := cms.CMSUtilsReplaceSigners(originalCMS, updatedSignerInfo)
	if err != nil {
		return nil, err
	}
	updatedCMS, err = cms.CMSUtilsPopulateDigestAlgorithmSet(updatedCMS, originalCMS.DigestAlgorithmIDs())
	if err != nil {
		return nil, err
	}
	updatedCMS, err = b.addNewCertificates(updatedCMS, parameters)
	if err != nil {
		return nil, err
	}
	return cms.CMSUtilsWriteToDSSDocument(updatedCMS, b.resourcesHandlerBuilder)
}

// getUpdatedSignerInformations ports the private
// getUpdatedSignerInformations(CMS, SignerInformationStore, CAdESCounterSignatureParameters,
// SignatureValue, CAdESSignature).
func (b *CAdESCounterSignatureBuilder) getUpdatedSignerInformations(originalCMS *cms.CMS, signerInformationStore []*cmscore.SignerInfo,
	parameters *CAdESCounterSignatureParameters, signatureValue *model.SignatureValue, masterSignature *CAdESSignature) ([]*cmscore.SignerInfo, error) {

	var result []*cmscore.SignerInfo
	for _, signerInformation := range signerInformationStore {
		cadesSignature := NewCAdESSignature(originalCMS, signerInformation)
		cadesSignature.SetMasterSignature(cadesCounterSignatureBuilderAsAdvancedSignature(masterSignature))
		cadesSignature.SetDetachedContents(parameters.DetachedContents())
		cadesSignature.SetManifestFile(b.manifestFile)

		if cadesSignature.ID() == parameters.SignatureIdToCounterSign() {
			if masterSignature != nil {
				panic("Cannot recursively add a counter-signature")
			}
			if err := b.assertCounterSignaturePossible(signerInformation); err != nil {
				return nil, err
			}

			counterSignatureSignerInfos, err := b.generateCounterSignatureFromSignatureValue(signerInformation, parameters, signatureValue)
			if err != nil {
				return nil, err
			}

			updated, err := cadesCounterSignatureBuilderAddCounterSigners(signerInformation, counterSignatureSignerInfos)
			if err != nil {
				return nil, err
			}
			result = append(result, updated)

		} else if existingCounterSigners := cadesCounterSignaturesOf(signerInformation); len(existingCounterSigners) > 0 {
			updatedCounterSigners, err := b.getUpdatedSignerInformations(originalCMS, existingCounterSigners, parameters, signatureValue, cadesSignature)
			if err != nil {
				return nil, err
			}
			updated, err := b.replaceCounterSigners(signerInformation, updatedCounterSigners)
			if err != nil {
				return nil, err
			}
			result = append(result, updated)

		} else {
			result = append(result, signerInformation)
		}
	}

	return result, nil
}

// cadesCounterSignatureBuilderAsAdvancedSignature converts a nilable *CAdESSignature into a
// validation.AdvancedSignature, mapping a nil pointer to a true nil interface value (avoiding
// the classic Go typed-nil-in-interface trap: a naive interface conversion of a nil
// *CAdESSignature would make MasterSignature() != nil true even though there is no signature).
func cadesCounterSignatureBuilderAsAdvancedSignature(s *CAdESSignature) validation.AdvancedSignature {
	if s == nil {
		return nil
	}
	return s
}

// cadesCounterSignatureBuilderAddCounterSigners ports SignerInformation.addCounterSigners(
// SignerInformation, SignerInformationStore): it copies signerInformation, replacing the
// unsignedAttrs so that the (possibly non-existent) id-countersignature attribute now also
// carries every SignerInfo in newCounterSigners.
func cadesCounterSignatureBuilderAddCounterSigners(signerInformation *cmscore.SignerInfo, newCounterSigners []*cmscore.SignerInfo) (*cmscore.SignerInfo, error) {
	existing := cadesCounterSignaturesOf(signerInformation)
	allCounterSigners := append(append([]*cmscore.SignerInfo{}, existing...), newCounterSigners...)
	counterSignatureAttribute := cadesCounterSignatureBuilderCounterSignatureAttribute(allCounterSigners)

	var attrs cmscore.Attributes
	replaced := false
	for _, attribute := range signerInformation.UnsignedAttributes {
		if attribute.Type.Equal(cmscore.OIDCounterSignature) {
			if !replaced {
				attrs = append(attrs, counterSignatureAttribute)
				replaced = true
			}
			continue
		}
		attrs = append(attrs, attribute)
	}
	if !replaced {
		attrs = append(attrs, counterSignatureAttribute)
	}

	return cms.CMSUtilsReplaceUnsignedAttributes(signerInformation, attrs)
}

// replaceCounterSigners ports the private
// replaceCounterSigners(SignerInformation, List<SignerInformation>): it replaces the whole
// id-countersignature unsigned attribute's value set with updatedCounterSigners (used on the
// recursive branch, where the attribute already exists and every existing counter-signer -
// updated or not - is passed back in).
func (b *CAdESCounterSignatureBuilder) replaceCounterSigners(signerInformation *cmscore.SignerInfo, updatedCounterSigners []*cmscore.SignerInfo) (*cmscore.SignerInfo, error) {
	counterSignatureAttribute := cadesCounterSignatureBuilderCounterSignatureAttribute(updatedCounterSigners)

	var attrs cmscore.Attributes
	for _, attribute := range signerInformation.UnsignedAttributes {
		if attribute.Type.Equal(cmscore.OIDCounterSignature) {
			attrs = append(attrs, counterSignatureAttribute)
		} else {
			attrs = append(attrs, attribute)
		}
	}

	return cms.CMSUtilsReplaceUnsignedAttributes(signerInformation, attrs)
}

// cadesCounterSignatureBuilderCounterSignatureAttribute ports the private
// getUpdatedCounterSignatureAttribute(List<SignerInformation>).
func cadesCounterSignatureBuilderCounterSignatureAttribute(counterSigners []*cmscore.SignerInfo) *cmscore.Attribute {
	values := make([][]byte, 0, len(counterSigners))
	for _, counterSigner := range counterSigners {
		values = append(values, counterSigner.DER())
	}
	return cmscore.NewAttribute(cmscore.OIDCounterSignature, values...)
}

// addNewCertificates ports the private addNewCertificates(CMS, CAdESCounterSignatureParameters).
func (b *CAdESCounterSignatureBuilder) addNewCertificates(updatedCMS *cms.CMS, parameters *CAdESCounterSignatureParameters) (*cms.CMS, error) {
	selector := spi.NewBaselineBCertificateSelector(parameters.SigningCertificate(), parameters.CertificateChain()).
		SetTrustedCertificateSource(b.certificateVerifier.TrustedCertSources()).
		SetTrustAnchorBPPolicy(parameters.BLevel().IsTrustAnchorBPPolicy())
	newCertificates, err := selector.Certificates()
	if err != nil {
		return nil, err
	}

	cmsBuilder := cms.NewCMSBuilder().SetOriginalCMS(updatedCMS)
	return cmsBuilder.ExtendCMSSignedData(newCertificates, nil, nil)
}

// generateCounterSignatureFromSignatureValue ports the private
// generateCounterSignature(SignerInformation, CAdESCounterSignatureParameters, SignatureValue).
func (b *CAdESCounterSignatureBuilder) generateCounterSignatureFromSignatureValue(signerInformation *cmscore.SignerInfo,
	parameters *CAdESCounterSignatureParameters, signatureValue *model.SignatureValue) ([]*cmscore.SignerInfo, error) {
	signatureAlgorithm := parameters.SignatureAlgorithm()
	customContentSigner, err := cms.NewCustomContentSignerWithSignature(signatureAlgorithm.JCEID(), signatureValue.Value())
	if err != nil {
		return nil, err
	}
	return b.GenerateCounterSignature(signerInformation, &parameters.CAdESSignatureParameters, customContentSigner)
}

// GenerateCounterSignature generates a counter-signature SignerInformationStore (here: the
// []*cmscore.SignerInfo it would contain - see the file header). Port of
// #generateCounterSignature(SignerInformation, CAdESSignatureParameters, CustomContentSigner).
func (b *CAdESCounterSignatureBuilder) GenerateCounterSignature(signerInformation *cmscore.SignerInfo,
	parameters *CAdESSignatureParameters, customContentSigner *cms.CustomContentSigner) ([]*cmscore.SignerInfo, error) {

	toSignDocument := model.NewInMemoryDocument(signerInformation.Signature)
	cmsBuilderHelper := b.InitCMSBuilderHelper(toSignDocument, parameters, customContentSigner)
	signerInfoGenerator, err := cmsBuilderHelper.CreateSignerInfoGenerator()
	if err != nil {
		return nil, err
	}

	// NOTE: use a simplified SignerInfo generation to only create the required counter-signer.
	counterSigner, err := signerInfoGenerator.Generate(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to generate counter-signature: %w", err)
	}
	if counterSigner == nil {
		return nil, nil
	}
	return []*cmscore.SignerInfo{counterSigner}, nil
}

// GetSignerInformationToBeCounterSigned returns a SignerInformation to be counter-signed. Port
// of #getSignerInformationToBeCounterSigned(DSSDocument, CAdESCounterSignatureParameters).
func (b *CAdESCounterSignatureBuilder) GetSignerInformationToBeCounterSigned(signatureDocument model.DSSDocument,
	parameters *CAdESCounterSignatureParameters) (*cmscore.SignerInfo, error) {
	cadesSignature, err := b.getSignatureById(signatureDocument, parameters)
	if err != nil {
		return nil, err
	}
	if cadesSignature == nil {
		return nil, fmt.Errorf("CAdESSignature not found with the given dss id '%s'", parameters.SignatureIdToCounterSign())
	}
	return cadesSignature.SignerInformation(), nil
}

// getSignatureById ports the private getSignatureById(DSSDocument, CAdESCounterSignatureParameters).
//
// Panics when SignatureIdToCounterSign is empty (Java Objects.requireNonNull).
func (b *CAdESCounterSignatureBuilder) getSignatureById(signatureDocument model.DSSDocument,
	parameters *CAdESCounterSignatureParameters) (*CAdESSignature, error) {
	if parameters.SignatureIdToCounterSign() == "" {
		panic("The Id of a signature to be counter signed shall be defined! " +
			"Please use SerializableCounterSignatureParameters.setSignatureIdToCounterSign(signatureId) method.")
	}

	analyzer, err := NewCMSDocumentAnalyzerFromDocument(signatureDocument)
	if err != nil {
		return nil, err
	}
	analyzer.SetDetachedContents(parameters.DetachedContents())
	analyzer.SetManifestFile(b.manifestFile)

	signatures := analyzer.Signatures()
	return b.findSignatureRecursive(signatures, parameters.SignatureIdToCounterSign())
}

// findSignatureRecursive ports the private
// findSignatureRecursive(List<AdvancedSignature>, String).
func (b *CAdESCounterSignatureBuilder) findSignatureRecursive(signatures []validation.AdvancedSignature, signatureId string) (*CAdESSignature, error) {
	for _, advancedSignature := range signatures {
		if signatureId == advancedSignature.ID() {
			cadesSignature, ok := advancedSignature.(*CAdESSignature)
			if !ok {
				return nil, nil
			}
			if err := b.assertCounterSignaturePossible(cadesSignature.SignerInformation()); err != nil {
				return nil, err
			}
			return cadesSignature, nil
		}

		counterSignatureById, err := b.findSignatureRecursive(advancedSignature.CounterSignatures(), signatureId)
		if err != nil {
			return nil, err
		}
		if counterSignatureById != nil {
			// TODO : add a nested counter signature support + check if a master signature is not timestamped
			panic("Nested counter signatures are not supported with CAdES!")
		}
	}
	return nil, nil
}

// assertCounterSignaturePossible ports the private
// assertCounterSignaturePossible(SignerInformation).
func (b *CAdESCounterSignatureBuilder) assertCounterSignaturePossible(signerInformation *cmscore.SignerInfo) error {
	if CAdESUtilsContainsATSTv2(signerInformation) {
		return exception.NewIllegalInputException("Cannot add a counter signature to a CAdES containing an archiveTimestampV2")
	}
	if CAdESUtilsContainsEvidenceRecord(signerInformation) {
		return exception.NewIllegalInputException("Cannot add a counter signature to a CMS containing an evidence record unsigned attribute.")
	}
	return nil
}

// InitCMSBuilderHelper instantiates a CMSForCAdESBuilderHelper. Port of
// #initCMSBuilderHelper(DSSDocument, CAdESSignatureParameters, ContentSigner).
func (b *CAdESCounterSignatureBuilder) InitCMSBuilderHelper(contentToSign model.DSSDocument, signatureParameters *CAdESSignatureParameters,
	contentSigner cms.ContentSigner) *CMSForCAdESBuilderHelper {
	helper := NewCMSForCAdESBuilderHelper(contentToSign, signatureParameters, contentSigner)
	helper.SetTrustedCertificateSource(b.certificateVerifier.TrustedCertSources())
	return helper
}
