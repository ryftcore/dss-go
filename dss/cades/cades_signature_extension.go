// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CAdESSignatureExtension.java (DSS 6.5.RC1).
//
// Java's abstract class becomes an embeddable base plus the CAdESSignatureExtensionOverrides
// interface its subclasses (LevelBaselineT, -LT, -LTA) satisfy and register with
// InitCAdESSignatureExtension, the TokenBase.InitToken(self) convention of PORTING.md: Go has no
// method overriding across embedding, so the base dispatches back into the concrete extension
// through that interface.
//
// Java's four extendCMSSignatures overloads get four Go names, since Go has no overloading:
//
//	extendSignatures(DSSDocument, SignatureParameters)                    -> ExtendSignatures
//	extendCMSSignatures(CMS, SignatureParameters)                         -> ExtendCMSSignatures
//	extendCMSSignatures(CMS, SignerInformation, SignatureParameters)      -> ExtendCMSSignaturesForSigner
//	extendCMSSignatures(CMS, Collection<SignerInformation>, ...)               -> ExtendCMSSignaturesWithSigners
//	extendCMSSignatures(CMS, SignatureParameters, List<String>)           -> ExtendCMSSignaturesWithIds (abstract)
//
// BouncyCastle replacements (see PORTING.md): SignerInformation -> *cmscore.SignerInfo, and
// SignerInformationStore -> the []*cmscore.SignerInfo the CMS layer already exchanges;
// AttributeTable -> cmscore.Attributes.
//
// slf4j logging is dropped (PORTING.md).
package cades

import (
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/signature/resources"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/executor"
	"github.com/ryftcore/dss-go/dss/utils"
)

// CAdESSignatureExtensionOverrides declares the operation Java's abstract
// SignatureExtension leaves to its subclasses and calls back into.
type SignatureExtensionOverrides interface {
	// ExtendCMSSignaturesWithIds extends the signatures in the cms whose ids are listed within
	// signatureIdsToExtend. Port of the protected abstract
	// extendCMSSignatures(CMS, SignatureParameters, List<String>).
	ExtendCMSSignaturesWithIds(cmsToExtend *cms.CMS, parameters *SignatureParameters,
		signatureIdsToExtend []string) (*cms.CMS, error)
}

// SignatureExtension is the base class for extending a Signature. It implements
// document.SignatureExtension[*SignatureParameters] through the concrete extension
// embedding it.
type SignatureExtension struct {
	// TspSource is the TSPSource to request a timestamp (T- and LTA-levels).
	TspSource validation.TSPSource

	// CertificateVerifier is the CertificateVerifier to use.
	CertificateVerifier validation.CertificateVerifier

	// ResourcesHandlerBuilder is used to create data container objects such as an OutputStream
	// or a DSSDocument.
	ResourcesHandlerBuilder resources.DSSResourcesHandlerBuilder

	// overrides points back at the concrete extension; see InitCAdESSignatureExtension.
	overrides SignatureExtensionOverrides
}

// InitCAdESSignatureExtension registers the concrete extension with its base, and applies the
// constructor's checks. Port of the protected
// CAdESSignatureExtension(TSPSource, CertificateVerifier) constructor; panics with the Java
// message when tspSource is nil (Objects.requireNonNull).
func (e *SignatureExtension) InitCAdESSignatureExtension(self SignatureExtensionOverrides,
	tspSource validation.TSPSource, certificateVerifier validation.CertificateVerifier) {
	if tspSource == nil {
		panic("The TSPSource cannot be null")
	}
	e.overrides = self
	e.TspSource = tspSource
	e.CertificateVerifier = certificateVerifier
}

// SetResourcesHandlerBuilder sets a DSSResourcesHandlerBuilder to be used for operating with
// internal objects during the signature creation procedure.
// Port of setResourcesHandlerBuilder(DSSResourcesHandlerBuilder).
func (e *SignatureExtension) SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) {
	e.ResourcesHandlerBuilder = resourcesHandlerBuilder
}

// ExtendSignatures extends the CMS signatures provided within the signatureToExtend document
// and returns a new extended document. Port of
// extendSignatures(DSSDocument, SignatureParameters).
func (e *SignatureExtension) ExtendSignatures(signatureToExtend model.DSSDocument,
	parameters *SignatureParameters) (model.DSSDocument, error) {
	cmsToExtend, err := e.getCMS(signatureToExtend)
	if err != nil {
		return nil, err
	}
	extendedCMS, err := e.ExtendCMSSignatures(cmsToExtend, parameters)
	if err != nil {
		return nil, err
	}
	return cms.UtilsWriteToDSSDocument(extendedCMS, e.ResourcesHandlerBuilder)
}

// getCMS ports the private getCMS(DSSDocument).
func (e *SignatureExtension) getCMS(document model.DSSDocument) (*cms.CMS, error) {
	return cms.UtilsParseToCMS(document)
}

// ExtendCMSSignatures extends every signer of a CMS.
// Port of extendCMSSignatures(CMS, CAdESSignatureParameters).
func (e *SignatureExtension) ExtendCMSSignatures(cmsToExtend *cms.CMS,
	parameters *SignatureParameters) (*cms.CMS, error) {
	return e.ExtendCMSSignaturesWithSigners(cmsToExtend, cmsToExtend.SignerInfos(), parameters)
}

// ExtendCMSSignaturesForSigner extends a CMS with a specified SignerInformation.
// NOTE: does not modify other SignerInformations.
// Port of extendCMSSignatures(CMS, SignerInformation, CAdESSignatureParameters).
func (e *SignatureExtension) ExtendCMSSignaturesWithSigner(cmsToExtend *cms.CMS,
	signerInformation *cmscore.SignerInfo, parameters *SignatureParameters) (*cms.CMS, error) {
	return e.ExtendCMSSignaturesWithSigners(cmsToExtend, []*cmscore.SignerInfo{signerInformation}, parameters)
}

// ExtendCMSSignaturesWithSigners loops on each signerInformation of the cms and extends the ones
// defined in the collection signerInformationsToExtend.
// Port of the protected extendCMSSignatures(CMS, Collection<SignerInformation>, CAdESSignatureParameters).
func (e *SignatureExtension) ExtendCMSSignaturesWithSigners(cmsToExtend *cms.CMS,
	signerInformationsToExtend []*cmscore.SignerInfo, parameters *SignatureParameters) (*cms.CMS, error) {
	// extract signerInformations before pre-extension
	signerInformationCollection := cmsToExtend.SignerInfos()
	if utils.IsCollectionEmpty(signerInformationCollection) {
		return nil, exception.NewIllegalInputException("Unable to extend the document! No signatures found.")
	}

	signatureIdsToExtend := make([]string, 0)

	analyzer, err := e.DocumentAnalyzer(cmsToExtend, parameters)
	if err != nil {
		return nil, err
	}
	for _, signature := range analyzer.Signatures() {
		cadesSignature, ok := signature.(*Signature)
		if !ok {
			// Java's cast; a CMSDocumentAnalyzer only ever yields CAdESSignatures.
			continue
		}
		if cadesLTAContainsSigner(signerInformationsToExtend, cadesSignature.SignerInformation()) {
			signatureIdsToExtend = append(signatureIdsToExtend, cadesSignature.ID())
		}
	}

	return e.overrides.ExtendCMSSignaturesWithIds(cmsToExtend, parameters, signatureIdsToExtend)
}

// cadesLTAContainsSigner is Collection#contains(SignerInformation). BouncyCastle's
// SignerInformation does not override equals(), so Java compares by identity; the pointer
// comparison below is that same comparison, and the SignerInfos of one parsed CMS are shared
// rather than re-parsed, so the two agree.
func cadesLTAContainsSigner(signerInformations []*cmscore.SignerInfo, signerInformation *cmscore.SignerInfo) bool {
	for _, candidate := range signerInformations {
		if candidate == signerInformation {
			return true
		}
	}
	return false
}

// ReplaceSigners replaces the signers within the provided originalCMS.
// Port of the protected replaceSigners(CMS, List<SignerInformation>).
func (e *SignatureExtension) ReplaceSigners(originalCMS *cms.CMS,
	newSignerInformationList []*cmscore.SignerInfo) (*cms.CMS, error) {
	updatedCmsSignedData, err := cms.UtilsReplaceSigners(originalCMS, newSignerInformationList)
	if err != nil {
		return nil, err
	}
	return cms.UtilsPopulateDigestAlgorithmSet(updatedCmsSignedData, originalCMS.DigestAlgorithmIDs())
}

// TimeStampAttributeValue generates and returns a TimeStamp attribute value: the DER encoding of
// a TimeStampToken whose unsigned attributes carry attributesForTimestampToken.
// Port of the protected getTimeStampAttributeValue(DSSMessageDigest, DigestAlgorithm, Attribute...).
func (e *SignatureExtension) TimeStampAttributeValue(timestampMessageDigest model.DSSMessageDigest,
	timestampDigestAlgorithm enumerations.DigestAlgorithm,
	attributesForTimestampToken ...*cmscore.Attribute) ([]byte, error) {
	timeStampToken, err := e.TspSource.TimeStampResponse(timestampDigestAlgorithm, timestampMessageDigest.Value())
	if err != nil {
		return nil, err
	}
	timestampCMS, err := cms.UtilsParseToCMSBinaries(timeStampToken.Bytes())
	if err != nil {
		return nil, err
	}

	// TODO (upstream, 27/08/2014): attributesForTimestampToken cannot be null: to be modified
	if attributesForTimestampToken != nil {
		// timeStampToken contains one and only one signer
		signerInformation := timestampCMS.SignerInfos()[0]
		unsignedAttributes := UtilsUnsignedAttributes(signerInformation)
		for _, attributeToAdd := range attributesForTimestampToken {
			attrType := attributeToAdd.Type
			objectAt := attributeToAdd.ValueEncodings()[0]
			unsignedAttributes = UtilsAddAttribute(unsignedAttributes, attrType, objectAt)
		}
		// Unsigned attributes cannot be empty (RFC 5652 5.3)
		if len(unsignedAttributes) == 0 {
			unsignedAttributes = nil
		}
		newSignerInformation, err := cms.UtilsReplaceUnsignedAttributes(signerInformation, unsignedAttributes)
		if err != nil {
			return nil, err
		}
		timestampCMS, err = cms.UtilsReplaceSigners(timestampCMS, []*cmscore.SignerInfo{newSignerInformation})
		if err != nil {
			return nil, err
		}
	}
	return timestampCMS.DEREncoded(), nil
}

// DocumentAnalyzer returns a document analyzer for a CMS.
// Port of the protected getDocumentAnalyzer(CMS, CAdESSignatureParameters).
func (e *SignatureExtension) DocumentAnalyzer(cmsToAnalyze *cms.CMS,
	parameters *SignatureParameters) (*CMSDocumentAnalyzer, error) {
	documentValidator := NewCMSDocumentAnalyzer(cmsToAnalyze)
	documentValidator.SetCertificateVerifier(e.CertificateVerifier)
	documentValidator.SetDetachedContents(parameters.DetachedContents())
	documentValidator.SetValidationContextExecutor(executor.CompleteValidationContextExecutorInstance)
	return documentValidator, nil
}
