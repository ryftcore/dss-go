// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CAdESSignaturePolicyStoreBuilder.java (DSS 6.5.RC1).
package cades

import (
	"fmt"

	"github.com/utain/esig/dss/cms"
	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/internal/cmscore"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/exception"
	"github.com/utain/esig/dss/spi/policy"
	"github.com/utain/esig/dss/spi/signature/resources"
)

// CAdESSignaturePolicyStoreBuilder builds a SignaturePolicyStore for a CAdES signature.
type CAdESSignaturePolicyStoreBuilder struct {
	// resourcesHandlerBuilder is used to create data container objects such as an OutputStream or
	// a DSSDocument.
	resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder
}

// NewCAdESSignaturePolicyStoreBuilder is the default constructor.
func NewCAdESSignaturePolicyStoreBuilder() *CAdESSignaturePolicyStoreBuilder {
	return &CAdESSignaturePolicyStoreBuilder{}
}

// SetResourcesHandlerBuilder sets a DSSResourcesHandlerBuilder to be used for operating with
// internal objects during the signature creation procedure. Port of
// #setResourcesHandlerBuilder.
func (b *CAdESSignaturePolicyStoreBuilder) SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) {
	b.resourcesHandlerBuilder = resourcesHandlerBuilder
}

// AddSignaturePolicyStore extends all signatures within the given document, matching the
// provided policy in signaturePolicyStore. Port of
// #addSignaturePolicyStore(DSSDocument, SignaturePolicyStore).
//
// Panics when signatureDocument is nil (Java Objects.requireNonNull("Signature document must be
// provided!")).
func (b *CAdESSignaturePolicyStoreBuilder) AddSignaturePolicyStore(signatureDocument model.DSSDocument,
	signaturePolicyStore *model.SignaturePolicyStore) (model.DSSDocument, error) {
	if signatureDocument == nil {
		panic("Signature document must be provided!")
	}

	originalCmsSignedData, err := cms.CMSUtilsParseToCMS(signatureDocument)
	if err != nil {
		return nil, err
	}
	newCmsSignedData, err := b.ExtendCMS(originalCmsSignedData, signaturePolicyStore)
	if err != nil {
		return nil, err
	}
	newCmsSignedData, err = cms.CMSUtilsPopulateDigestAlgorithmSet(newCmsSignedData, originalCmsSignedData.DigestAlgorithmIDs())
	if err != nil {
		return nil, err
	}
	return cms.CMSUtilsWriteToDSSDocument(newCmsSignedData, b.resourcesHandlerBuilder)
}

// ExtendCMS creates a new CMS with a SignaturePolicyStore for matching signatures. Port of
// #extendCMS(CMS, SignaturePolicyStore).
//
// Panics when cmsObj is nil (Java Objects.requireNonNull("CMS must be provided!")).
func (b *CAdESSignaturePolicyStoreBuilder) ExtendCMS(cmsObj *cms.CMS, signaturePolicyStore *model.SignaturePolicyStore) (*cms.CMS, error) {
	if cmsObj == nil {
		panic("CMS must be provided!")
	}
	if err := b.assertConfigurationValid(signaturePolicyStore); err != nil {
		return nil, err
	}

	documentAnalyzer := NewCMSDocumentAnalyzer(cmsObj)
	signatures := documentAnalyzer.Signatures()
	if len(signatures) == 0 {
		return nil, exception.NewIllegalInputException("Unable to extend the document! No signatures found.")
	}

	var newSignerInformationList []*cmscore.SignerInfo

	signaturePolicyStoreAdded := false
	for _, sig := range signatures {
		cadesSignature, ok := sig.(*CAdESSignature)
		if !ok {
			return nil, fmt.Errorf("unexpected signature type %T", sig)
		}
		newSignerInformation, err := b.addSignaturePolicyStoreIfDigestMatch(cadesSignature, signaturePolicyStore)
		if err != nil {
			return nil, err
		}
		if cadesSignature.SignerInformation() != newSignerInformation {
			signaturePolicyStoreAdded = true
		}
		newSignerInformationList = append(newSignerInformationList, newSignerInformation)
	}
	if !signaturePolicyStoreAdded {
		return nil, exception.NewIllegalInputException("The process did not find a signature to add SignaturePolicyStore!")
	}
	return cms.CMSUtilsReplaceSigners(cmsObj, newSignerInformationList)
}

// AddSignaturePolicyStoreForSignature adds a signaturePolicyStore to a signature with the given
// signatureId, if the signature policy identifier matches the policy provided within
// signaturePolicyStore. Port of #addSignaturePolicyStore(DSSDocument, SignaturePolicyStore,
// String).
//
// Panics when signatureDocument is nil (Java Objects.requireNonNull("Signature document must be
// provided!")).
func (b *CAdESSignaturePolicyStoreBuilder) AddSignaturePolicyStoreForSignature(signatureDocument model.DSSDocument,
	signaturePolicyStore *model.SignaturePolicyStore, signatureId string) (model.DSSDocument, error) {
	if signatureDocument == nil {
		panic("Signature document must be provided!")
	}

	originalCmsSignedData, err := cms.CMSUtilsParseToCMS(signatureDocument)
	if err != nil {
		return nil, err
	}
	newCmsSignedData, err := b.ExtendCMSForSignature(originalCmsSignedData, signaturePolicyStore, signatureId)
	if err != nil {
		return nil, err
	}
	newCmsSignedData, err = cms.CMSUtilsPopulateDigestAlgorithmSet(newCmsSignedData, originalCmsSignedData.DigestAlgorithmIDs())
	if err != nil {
		return nil, err
	}
	return cms.CMSUtilsWriteToDSSDocument(newCmsSignedData, b.resourcesHandlerBuilder)
}

// ExtendCMSForSignature creates a new CMS with a SignaturePolicyStore for a signature with
// signatureId. Port of #extendCMS(CMS, SignaturePolicyStore, String).
//
// Panics when cmsObj is nil (Java Objects.requireNonNull("CMS must be provided!")).
func (b *CAdESSignaturePolicyStoreBuilder) ExtendCMSForSignature(cmsObj *cms.CMS, signaturePolicyStore *model.SignaturePolicyStore,
	signatureId string) (*cms.CMS, error) {
	if cmsObj == nil {
		panic("CMS must be provided!")
	}
	if err := b.assertConfigurationValid(signaturePolicyStore); err != nil {
		return nil, err
	}

	documentAnalyzer := NewCMSDocumentAnalyzer(cmsObj)
	sig := documentAnalyzer.SignatureByID(signatureId)
	if sig == nil {
		return nil, exception.NewIllegalInputException(fmt.Sprintf("Unable to find a signature with Id : %s!", signatureId))
	}

	var newSignerInformationList []*cmscore.SignerInfo
	for _, currentSignature := range documentAnalyzer.Signatures() {
		cadesSignature, ok := currentSignature.(*CAdESSignature)
		if !ok {
			return nil, fmt.Errorf("unexpected signature type %T", currentSignature)
		}
		if sig.ID() == cadesSignature.ID() {
			newSignerInformation, err := b.addSignaturePolicyStoreIfDigestMatch(cadesSignature, signaturePolicyStore)
			if err != nil {
				return nil, err
			}
			if cadesSignature.SignerInformation() == newSignerInformation {
				return nil, exception.NewIllegalInputException(fmt.Sprintf(
					"The process was not able to add SignaturePolicyStore to a signature with Id : %s!", signatureId))
			}
			newSignerInformationList = append(newSignerInformationList, newSignerInformation)

		} else {
			newSignerInformationList = append(newSignerInformationList, cadesSignature.SignerInformation())
		}
	}
	return cms.CMSUtilsReplaceSigners(cmsObj, newSignerInformationList)
}

// addSignaturePolicyStoreIfDigestMatch adds SignaturePolicyStore to cadesSignature if required.
// Port of the protected #addSignaturePolicyStoreIfDigestMatch(CAdESSignature,
// SignaturePolicyStore).
func (b *CAdESSignaturePolicyStoreBuilder) addSignaturePolicyStoreIfDigestMatch(cadesSignature *CAdESSignature,
	signaturePolicyStore *model.SignaturePolicyStore) (*cmscore.SignerInfo, error) {
	signerInformation := cadesSignature.SignerInformation()

	if err := b.assertSignaturePolicyStoreExtensionPossible(signerInformation); err != nil {
		return nil, err
	}
	newSignerInformation := signerInformation

	match, err := b.checkDigest(cadesSignature, signaturePolicyStore)
	if err != nil {
		return nil, err
	}
	if match {
		newSignerInformation, err = b.addSignaturePolicyStoreToSignerInformation(signerInformation, signaturePolicyStore)
		if err != nil {
			return nil, err
		}
	}
	return newSignerInformation, nil
}

// checkDigest verifies if the digests computed in the provided signaturePolicyStore match the
// digest defined in the incorporated signature policy identifier. Port of the protected
// #checkDigest(CAdESSignature, SignaturePolicyStore).
func (b *CAdESSignaturePolicyStoreBuilder) checkDigest(cadesSignature *CAdESSignature, signaturePolicyStore *model.SignaturePolicyStore) (bool, error) {
	signaturePolicy := cadesSignature.SignaturePolicy()
	if signaturePolicy == nil {
		// signature-policy-identifier is not defined for a signature.
		return false, nil
	}
	expectedDigest := signaturePolicy.Digest()
	if expectedDigest.IsEmpty() {
		// signature-policy-identifier digest is not found for the signature.
		return false, nil
	}

	signaturePolicyContent := signaturePolicyStore.SignaturePolicyContent()
	if signaturePolicyContent == nil {
		// No policy document has been provided. Digests are not checked.
		return true, nil
	}
	signaturePolicy.SetPolicyContent(signaturePolicyContent)

	// The input to hash computation of sigPolicyHash depends on the technical specification of
	// the signature policy.
	signaturePolicyValidatorLoader := policy.DefaultSignaturePolicyValidatorLoaderPolicyBased()
	validator := signaturePolicyValidatorLoader.LoadValidator(signaturePolicy)
	computedDigest := validator.GetComputedDigest(signaturePolicyContent, expectedDigest.Algorithm())

	return expectedDigest.Equals(computedDigest), nil
}

// addSignaturePolicyStoreToSignerInformation ports the private
// addSignaturePolicyStore(SignerInformation, SignaturePolicyStore).
func (b *CAdESSignaturePolicyStoreBuilder) addSignaturePolicyStoreToSignerInformation(signerInformation *cmscore.SignerInfo,
	signaturePolicyStore *model.SignaturePolicyStore) (*cmscore.SignerInfo, error) {
	unsignedAttributes := CAdESUtilsUnsignedAttributes(signerInformation)
	sigPolicyStore, err := b.signaturePolicyStoreValue(signaturePolicyStore)
	if err != nil {
		return nil, err
	}
	unsignedAttributesWithPolicyStore := append(append(cmscore.Attributes{}, unsignedAttributes...),
		cmscore.NewAttribute(spi.OID_id_aa_ets_sigPolicyStore, sigPolicyStore))
	return cms.CMSUtilsReplaceUnsignedAttributes(signerInformation, unsignedAttributesWithPolicyStore)
}

// signaturePolicyStoreValue ports the private getSignaturePolicyStore(SignaturePolicyStore):
//
//	SignaturePolicyStore ::= SEQUENCE {
//	 spDocSpec SPDocSpecification ,
//	 spDocument SignaturePolicyDocument
//	}
//	SignaturePolicyDocument ::= CHOICE {
//	 sigPolicyEncoded OCTET STRING,
//	 sigPolicyLocalURI IA5String
//	}
func (b *CAdESSignaturePolicyStoreBuilder) signaturePolicyStoreValue(signaturePolicyStore *model.SignaturePolicyStore) ([]byte, error) {
	var body []byte

	// spDocSpec
	spDocSpec, err := spi.DSSASN1UtilsBuildSPDocSpecificationID(signaturePolicyStore.SpDocSpecification().Id())
	if err != nil {
		return nil, err
	}
	body = append(body, spDocSpec...)

	// spDocument
	signaturePolicyContent := signaturePolicyStore.SignaturePolicyContent()
	if signaturePolicyContent != nil {
		content, err := spi.DSSUtilsToByteArrayOfDocument(signaturePolicyContent)
		if err != nil {
			return nil, err
		}
		body = append(body, asn1ber.WriteTLV(asn1ber.TagOctetString, content)...)
	}
	sigPolDocLocalURI := signaturePolicyStore.SigPolDocLocalURI()
	if sigPolDocLocalURI != "" {
		body = append(body, asn1ber.WriteTLV(asn1ber.TagIA5String, []byte(sigPolDocLocalURI))...)
	}

	return asn1ber.WriteSequence(body), nil
}

// assertConfigurationValid ports the private #assertConfigurationValid(SignaturePolicyStore).
//
// Panics with the Java message when signaturePolicyStore, its SpDocSpecification or the
// SpDocSpecification's Id is missing (Java Objects.requireNonNull).
func (b *CAdESSignaturePolicyStoreBuilder) assertConfigurationValid(signaturePolicyStore *model.SignaturePolicyStore) error {
	if signaturePolicyStore == nil {
		panic("SignaturePolicyStore must be provided")
	}
	if signaturePolicyStore.SpDocSpecification() == nil {
		panic("SpDocSpecification must be provided")
	}
	if signaturePolicyStore.SpDocSpecification().Id() == "" {
		panic("ID (OID or URI) for SpDocSpecification must be provided")
	}

	signaturePolicyContentPresent := signaturePolicyStore.SignaturePolicyContent() != nil
	sigPolDocLocalURIPresent := signaturePolicyStore.SigPolDocLocalURI() != ""
	if signaturePolicyContentPresent == sigPolDocLocalURIPresent {
		panic("SignaturePolicyStore shall contain either " +
			"SignaturePolicyContent document or sigPolDocLocalURI!")
	}
	return nil
}

// assertSignaturePolicyStoreExtensionPossible ports the private
// assertSignaturePolicyStoreExtensionPossible(SignerInformation).
func (b *CAdESSignaturePolicyStoreBuilder) assertSignaturePolicyStoreExtensionPossible(signerInformation *cmscore.SignerInfo) error {
	if CAdESUtilsContainsATSTv2(signerInformation) {
		return exception.NewIllegalInputException("Cannot add signature policy store to a CAdES containing an archiveTimestampV2")
	}
	if CAdESUtilsContainsEvidenceRecord(signerInformation) {
		return exception.NewIllegalInputException("Cannot add signature policy store to a CMS containing an evidence record unsigned attribute.")
	}
	return nil
}
