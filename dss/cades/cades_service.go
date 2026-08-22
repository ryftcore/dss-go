// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CAdESService.java (DSS 6.5.RC1).
//
// Java extends AbstractSignatureService<CAdESSignatureParameters, CAdESTimestampParameters> and
// implements SignatureService, CounterSignatureService and
// EvidenceRecordIncorporationService; the Go port embeds
// document.AbstractSignatureService[*SignatureParameters, *TimestampParameters] and
// satisfies the three interfaces with the methods below - the compile-time assertions at the end
// of the file check that it does.
//
// # Errors
//
// The three service interfaces of dss-document return bare values, because none of the Java
// methods they declare has a checked exception: every failure here is a DSSException, an
// IllegalArgumentException or an IllegalInputException, all unchecked. This class is therefore
// the place where the (T, error) of the layers below turns back into Java's propagating
// exception, i.e. into a panic; the panic value is the error itself, so a recovering caller can
// still inspect it with errors.As.
//
// java.io.Serializable and the serialVersionUID are dropped (no Go counterpart), as is slf4j.
package cades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/signature/resources"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// SignatureExtender is the slice of the abstract SignatureExtension that
// getExtensionProfile's local variable is typed with upstream. Java assigns a
// LevelBaselineT, a LevelBaselineLT or a LevelBaselineLTA to a
// SignatureExtension variable; the three are distinct Go types embedding
// SignatureExtension, so the switch below needs an interface value, and Go's structural
// interfaces let this one name exactly the three operations Service then performs. Every
// CAdES extension satisfies it through its embedded SignatureExtension.
type SignatureExtender interface {
	// SetResourcesHandlerBuilder sets the DSSResourcesHandlerBuilder used while extending.
	// Port of CAdESSignatureExtension#setResourcesHandlerBuilder.
	SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder)

	// ExtendSignatures extends every signature of the given document.
	// Port of CAdESSignatureExtension#extendSignatures(DSSDocument, CAdESSignatureParameters).
	ExtendSignatures(signatureToExtend model.DSSDocument,
		parameters *SignatureParameters) (model.DSSDocument, error)

	// ExtendCMSSignaturesWithSigner extends the given SignerInformation of a CMS, leaving the
	// other signers untouched. Port of
	// CAdESSignatureExtension#extendCMSSignatures(CMS, SignerInformation, CAdESSignatureParameters).
	ExtendCMSSignaturesWithSigner(cmsToExtend *cms.CMS, signerInformation *cmscore.SignerInfo,
		parameters *SignatureParameters) (*cms.CMS, error)
}

// Service is the CAdES implementation of SignatureService.
type Service struct {
	document.AbstractSignatureService[*SignatureParameters, *TimestampParameters]

	// ResourcesHandlerBuilder is used to create data container objects such as an OutputStream
	// or a DSSDocument.
	ResourcesHandlerBuilder resources.DSSResourcesHandlerBuilder
}

// NewService creates an instance of the Service. A certificate verifier must be
// provided. Port of CAdESService(CertificateVerifier).
func NewService(certificateVerifier validation.CertificateVerifier) *Service {
	// Upstream logs "+ CAdESService created".
	return &Service{
		AbstractSignatureService: document.NewAbstractSignatureService[*SignatureParameters, *TimestampParameters](certificateVerifier),
		ResourcesHandlerBuilder:  CAdESUtilsDefaultResourcesHandlerBuilder,
	}
}

// SetResourcesHandlerBuilder sets a DSSResourcesHandlerBuilder to be used for operating with
// internal objects during the signature creation procedure.
// NOTE: The DSSResourcesHandlerBuilder is supported only within the 'dss-cms-stream' module!
// Port of #setResourcesHandlerBuilder.
func (s *Service) SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) {
	s.ResourcesHandlerBuilder = cms.UtilsResourcesHandlerBuilder(resourcesHandlerBuilder)
}

// GetContentTimestamp requests a content time-stamp for the document to be signed.
// Port of #getContentTimestamp.
func (s *Service) GetContentTimestamp(toSignDocument model.DSSDocument,
	parameters *SignatureParameters) *validation.TimestampToken {
	if s.TspSource == nil {
		panic("A TSPSource is required !")
	}

	digestAlgorithm := parameters.GetContentTimestampParameters().DigestAlgorithm()
	digestValue, err := toSignDocument.DigestValue(digestAlgorithm)
	if err != nil {
		panic(err)
	}
	timeStampResponse, err := s.TspSource.TimeStampResponse(digestAlgorithm, digestValue)
	if err != nil {
		panic(err)
	}
	timestampToken, err := validation.NewTimestampToken(timeStampResponse.Bytes(),
		enumerations.TimestampTypeContentTimestamp)
	if err != nil {
		panic(model.NewDSSErrorMessageCause("Cannot create a content TimestampToken", err))
	}
	return timestampToken
}

// GetDataToSign retrieves the data to be signed. Port of #getDataToSign.
func (s *Service) GetDataToSign(toSignDocument model.DSSDocument,
	parameters *SignatureParameters) *model.ToBeSigned {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}

	s.AssertSigningCertificateValid(&parameters.AbstractSignatureParameters)
	packaging := parameters.SignaturePackaging()
	cadesServiceAssertSignaturePackaging(packaging)

	originalCMS := s.originalCMS(toSignDocument, parameters)
	contentToSign := cadesServiceContentToSign(toSignDocument, parameters, originalCMS)

	contentSigner, err := cms.NewCustomContentSignerBuilder().Build(parameters.SignatureAlgorithm())
	if err != nil {
		panic(err)
	}
	cmsBuilderHelper := s.InitCMSBuilderHelper(contentToSign, parameters, contentSigner).
		SetOriginalCMS(originalCMS)
	if _, err := cmsBuilderHelper.CreateCMS(); err != nil {
		panic(err)
	}
	return model.NewToBeSignedWithBytes(contentSigner.OutputStream().Bytes())
}

// SignDocument signs the document with the provided signature value. Port of #signDocument.
func (s *Service) SignDocument(toSignDocument model.DSSDocument, parameters *SignatureParameters,
	signatureValue *model.SignatureValue) model.DSSDocument {
	if toSignDocument == nil {
		panic("toSignDocument cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	if signatureValue == nil {
		panic("SignatureValue cannot be null!")
	}

	s.AssertSigningCertificateValid(&parameters.AbstractSignatureParameters)
	packaging := parameters.SignaturePackaging()
	cadesServiceAssertSignaturePackaging(packaging)
	signatureAlgorithm := parameters.SignatureAlgorithm()
	signatureValue, err := s.EnsureSignatureValue(signatureAlgorithm, signatureValue)
	if err != nil {
		panic(err)
	}

	originalCMS := s.originalCMS(toSignDocument, parameters)
	if originalCMS == nil && enumerations.SignaturePackagingDetached == packaging {
		parameters.GetContext().SetDetachedContents([]model.DSSDocument{toSignDocument})
	}
	contentToSign := cadesServiceContentToSign(toSignDocument, parameters, originalCMS)

	contentSigner, err := cms.NewCustomContentSignerBuilder().BuildWithSignatureValue(
		parameters.SignatureAlgorithm(), signatureValue)
	if err != nil {
		panic(err)
	}
	cmsBuilderHelper := s.InitCMSBuilderHelper(contentToSign, parameters, contentSigner).
		SetIncludeUnsignedAttributes(true).
		SetOriginalCMS(originalCMS)

	signedCMS, err := cmsBuilderHelper.CreateCMS()
	if err != nil {
		panic(err)
	}

	signatureLevel := parameters.SignatureLevel()
	if enumerations.SignatureLevelCAdESBaselineB != signatureLevel {
		// Only the last signature will be extended
		newSignerInformation := cadesServiceNewSignerInformation(originalCMS, signedCMS)
		extension := s.extensionProfile(parameters)
		if signedCMS, err = extension.ExtendCMSSignaturesWithSigner(signedCMS, newSignerInformation, parameters); err != nil {
			panic(err)
		}
	}

	signature, err := cms.UtilsWriteToDSSDocument(signedCMS, s.ResourcesHandlerBuilder)
	if err != nil {
		panic(err)
	}
	name, err := s.GetFinalFileNameWithPackaging(toSignDocument, enumerations.SigningOperationSign,
		parameters.SignatureLevel(), parameters.SignaturePackaging())
	if err != nil {
		panic(err)
	}
	signature.SetName(name)
	parameters.Reinit()
	return signature
}

// ExtendDocument extends the signatures of the given document. All signatures are extended.
// Port of #extendDocument.
func (s *Service) ExtendDocument(toExtendDocument model.DSSDocument,
	parameters *SignatureParameters) model.DSSDocument {
	if toExtendDocument == nil {
		panic("toExtendDocument is not defined!")
	}
	if parameters == nil {
		panic("Cannot extend the signature. SignatureParameters are not defined!")
	}
	// false: All signature are extended
	extension := s.extensionProfile(parameters)
	dssDocument, err := extension.ExtendSignatures(toExtendDocument, parameters)
	if err != nil {
		panic(err)
	}
	name, err := s.GetFinalFileNameWithLevel(toExtendDocument, enumerations.SigningOperationExtend,
		parameters.SignatureLevel())
	if err != nil {
		panic(err)
	}
	dssDocument.SetName(name)
	return dssDocument
}

// cadesServiceContentToSign retrieves the data to be signed. If this data is located within a
// signature then it is extracted. Port of the private getContentToSign.
func cadesServiceContentToSign(toSignDocument model.DSSDocument, parameters *SignatureParameters,
	originalCMS *cms.CMS) model.DSSDocument {
	detachedContents := parameters.DetachedContents()
	if originalCMS == nil {
		return toSignDocument
	} else if utils.IsCollectionNotEmpty(detachedContents) {
		// * CAdES only can sign one document
		// * ASiC-S -> the document to sign or package.zip
		// * ASiC-E -> ASiCManifest
		return detachedContents[0]
	}
	return cadesServiceSignedContent(originalCMS)
}

// cadesServiceSignedContent returns the signed content of the already signed CMS.
// Port of the private getSignedContent.
func cadesServiceSignedContent(signedCMS *cms.CMS) model.DSSDocument {
	if signedCMS.IsDetachedSignature() {
		panic("Detached content shall be provided on parallel signing of a detached signature! " +
			"Please use cadesSignatureParameters#setDetachedContents method to provide original files.")
	}
	return signedCMS.SignedContent()
}

// cadesServiceNewSignerInformation ports the private getNewSignerInformation.
func cadesServiceNewSignerInformation(originalCMS, newCMS *cms.CMS) *cmscore.SignerInfo {
	signers := newCMS.SignerInfos()
	if originalCMS != nil {
		for _, signerInformation := range signers {
			if !cadesServiceContainsSignerInfo(originalCMS, signerInformation) {
				return signerInformation
			}
		}
	}
	// return the first one if originalSignedData is null (single signature creation)
	if len(signers) == 0 {
		// Java's iterator().next() raises a NoSuchElementException; a CMS built by
		// CMSForCAdESBuilderHelper always carries exactly one signer.
		return nil
	}
	return signers[0]
}

// cadesServiceContainsSignerInfo ports the private containsSignerInfo, whose Java body compares
// the SignerInfo structures by reference identity.
func cadesServiceContainsSignerInfo(signedCMS *cms.CMS, signerInformationToFind *cmscore.SignerInfo) bool {
	for _, signerInformation := range signedCMS.SignerInfos() {
		if signerInformationToFind == signerInformation {
			return true
		}
	}
	return false
}

// extensionProfile returns the extension profile to be used for a CAdES signature augmentation.
// Port of the private getExtensionProfile.
func (s *Service) extensionProfile(parameters *SignatureParameters) SignatureExtender {
	signatureLevel := parameters.SignatureLevel()
	if signatureLevel == "" {
		panic("SignatureLevel must be defined!")
	}
	var cadesSignatureExtension SignatureExtender
	switch signatureLevel {
	case enumerations.SignatureLevelCAdESBaselineT:
		cadesSignatureExtension = NewLevelBaselineT(s.TspSource, s.CertificateVerifier)
	case enumerations.SignatureLevelCAdESBaselineLT:
		cadesSignatureExtension = NewLevelBaselineLT(s.TspSource, s.CertificateVerifier)
	case enumerations.SignatureLevelCAdESBaselineLTA:
		cadesSignatureExtension = NewLevelBaselineLTA(s.TspSource, s.CertificateVerifier)
	default:
		panic(fmt.Sprintf("Unsupported signature format '%s' for extension.", signatureLevel))
	}
	cadesSignatureExtension.SetResourcesHandlerBuilder(s.ResourcesHandlerBuilder)
	return cadesSignatureExtension
}

// originalCMS returns the CMS of the document when it is a CMS signed message and a parallel
// signature is requested, nil otherwise. Port of the private getOriginalCMS.
func (s *Service) originalCMS(dssDocument model.DSSDocument, parameters *SignatureParameters) *cms.CMS {
	var originalCMS *cms.CMS
	_, isDigestDocument := dssDocument.(*model.DigestDocument)
	if parameters.IsParallelSignature() && !isDigestDocument && cadesServiceStartsWithSequenceTag(dssDocument) {
		if parsed, err := cms.UtilsParseToCMS(dssDocument); err == nil {
			originalCMS = parsed
		}
		// otherwise: not a parallel signature
		if originalCMS != nil {
			cadesServiceAssertSignaturePossible(originalCMS, parameters)
		}
	}
	return originalCMS
}

// cadesServiceStartsWithSequenceTag ports
// DSSASN1Utils.isASN1SequenceTag(DSSUtils.readFirstByte(dssDocument)); an unreadable document
// makes DSSUtils#readFirstByte raise a DSSException, which propagates here as a panic.
func cadesServiceStartsWithSequenceTag(dssDocument model.DSSDocument) bool {
	firstByte, err := spi.DSSUtilsReadFirstByte(dssDocument)
	if err != nil {
		panic(err)
	}
	return spi.DSSASN1UtilsIsASN1SequenceTag(firstByte)
}

// cadesServiceAssertSignaturePossible ports the private assertSignaturePossible.
func cadesServiceAssertSignaturePossible(originalCMS *cms.CMS, parameters *SignatureParameters) {
	if originalCMS.IsDetachedSignature() != (enumerations.SignaturePackagingDetached == parameters.SignaturePackaging()) {
		panic(fmt.Sprintf("Unable to create a parallel signature with packaging '%s'"+
			" which is different than the one used in the original signature!", parameters.SignaturePackaging()))
	}
	for _, signerInformation := range originalCMS.SignerInfos() {
		if UtilsContainsEvidenceRecord(signerInformation) {
			panic(exception.NewIllegalInputException(
				"Signature is not possible due to the CMS containing an evidence record unsigned attribute."))
		}
	}
}

// InitCMSBuilderHelper instantiates a CMSForCAdESBuilderHelper.
// Port of the protected #initCMSBuilderHelper.
func (s *Service) InitCMSBuilderHelper(contentToSign model.DSSDocument,
	signatureParameters *SignatureParameters, contentSigner cms.ContentSigner) *CMSForCAdESBuilderHelper {
	return NewCMSForCAdESBuilderHelper(contentToSign, signatureParameters, contentSigner).
		SetTrustedCertificateSource(s.CertificateVerifier.TrustedCertSources())
}

// cadesServiceAssertSignaturePackaging checks that the packaging is supported for this kind of
// signature. Port of the private assertSignaturePackaging.
func cadesServiceAssertSignaturePackaging(packaging enumerations.SignaturePackaging) {
	if packaging != enumerations.SignaturePackagingEnveloping && packaging != enumerations.SignaturePackagingDetached {
		panic("Unsupported signature packaging: " + string(packaging))
	}
}

// AddSignaturePolicyStore incorporates a Signature Policy Store as an unsigned property into the
// CAdES Signature. Port of #addSignaturePolicyStore.
func (s *Service) AddSignaturePolicyStore(doc model.DSSDocument,
	signaturePolicyStore *model.SignaturePolicyStore) model.DSSDocument {
	if doc == nil {
		panic("The document cannot be null")
	}
	if signaturePolicyStore == nil {
		panic("The signaturePolicyStore cannot be null")
	}

	builder := s.CAdESSignaturePolicyStoreBuilder()
	documentWithPolicyStore, err := builder.AddSignaturePolicyStore(doc, signaturePolicyStore)
	if err != nil {
		panic(err)
	}
	name, err := s.GetFinalFileNameWithLevel(doc, enumerations.SigningOperationExtend, "")
	if err != nil {
		panic(err)
	}
	documentWithPolicyStore.SetName(name)
	return documentWithPolicyStore
}

// SignaturePolicyStoreBuilder loads the relevant SignaturePolicyStoreBuilder.
// Port of the protected #getCAdESSignaturePolicyStoreBuilder.
func (s *Service) CAdESSignaturePolicyStoreBuilder() *SignaturePolicyStoreBuilder {
	builder := NewSignaturePolicyStoreBuilder()
	builder.SetResourcesHandlerBuilder(s.ResourcesHandlerBuilder)
	return builder
}

// GetDataToBeCounterSigned returns the data to be counter-signed for the signature identified by
// the parameters.
// Port of the (DSSDocument, CAdESCounterSignatureParameters) #getDataToBeCounterSigned.
func (s *Service) GetDataToBeCounterSigned(signatureDocument model.DSSDocument,
	parameters *CounterSignatureParameters) *model.ToBeSigned {
	if signatureDocument == nil {
		panic("signatureDocument cannot be null!")
	}
	if parameters == nil {
		panic("parameters cannot be null!")
	}
	if parameters.SignatureIdToCounterSign() == "" {
		panic("The signature to be counter-signed must be specified")
	}
	s.AssertSigningCertificateValid(&parameters.AbstractSignatureParameters)
	cadesServiceAssertCounterSignaturePossible(parameters)

	counterSignatureBuilder := s.CAdESCounterSignatureBuilder()
	signerInfoToCounterSign, err := counterSignatureBuilder.GetSignerInformationToBeCounterSigned(signatureDocument, parameters)
	if err != nil {
		panic(err)
	}

	return s.GetDataToBeCounterSignedForSigner(signerInfoToCounterSign, &parameters.SignatureParameters)
}

// GetDataToBeCounterSignedForSigner returns the data toBeSigned for a counter signature on the
// given signerInfoToCounterSign.
// Port of the (SignerInformation, CAdESSignatureParameters) #getDataToBeCounterSigned.
func (s *Service) GetDataToBeCounterSignedForSigner(signerInfoToCounterSign *cmscore.SignerInfo,
	parameters *SignatureParameters) *model.ToBeSigned {
	signatureAlgorithm := parameters.SignatureAlgorithm()
	customContentSigner, err := cms.NewCustomContentSignerBuilder().Build(signatureAlgorithm)
	if err != nil {
		panic(err)
	}

	counterSignatureBuilder := s.CAdESCounterSignatureBuilder()
	if _, err := counterSignatureBuilder.GenerateCounterSignature(signerInfoToCounterSign, parameters,
		customContentSigner); err != nil {
		panic(err)
	}

	return model.NewToBeSignedWithBytes(customContentSigner.OutputStream().Bytes())
}

// CounterSignSignature counter-signs the signature identified by the parameters.
// Port of #counterSignSignature.
func (s *Service) CounterSignSignature(signatureDocument model.DSSDocument,
	parameters *CounterSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	if signatureDocument == nil {
		panic("signatureDocument cannot be null!")
	}
	if parameters == nil {
		panic("parameters cannot be null!")
	}
	if parameters.SignatureIdToCounterSign() == "" {
		panic("The signature to be counter-signed must be specified")
	}
	if signatureValue == nil {
		panic("signatureValue cannot be null!")
	}
	s.AssertSigningCertificateValid(&parameters.AbstractSignatureParameters)
	cadesServiceAssertCounterSignaturePossible(parameters)
	signatureValue, err := s.EnsureSignatureValue(parameters.SignatureAlgorithm(), signatureValue)
	if err != nil {
		panic(err)
	}

	originalCMS, err := cms.UtilsParseToCMS(signatureDocument)
	if err != nil {
		panic(err)
	}

	counterSignatureBuilder := s.CAdESCounterSignatureBuilder()
	counterSigned, err := counterSignatureBuilder.AddCounterSignature(originalCMS, parameters, signatureValue)
	if err != nil {
		panic(err)
	}
	name, err := s.GetFinalFileNameWithLevel(signatureDocument, enumerations.SigningOperationCounterSign,
		parameters.SignatureLevel())
	if err != nil {
		panic(err)
	}
	counterSigned.SetName(name)
	counterSigned.SetMimeType(signatureDocument.MimeType())

	return counterSigned
}

// CounterSignatureBuilder loads the relevant CounterSignatureBuilder.
// Port of the protected #getCAdESCounterSignatureBuilder.
func (s *Service) CAdESCounterSignatureBuilder() *CounterSignatureBuilder {
	counterSignatureBuilder := NewCounterSignatureBuilder(s.CertificateVerifier)
	counterSignatureBuilder.SetResourcesHandlerBuilder(s.ResourcesHandlerBuilder)
	return counterSignatureBuilder
}

// AddSignatureEvidenceRecord incorporates an evidence record into the signature document.
// Port of #addSignatureEvidenceRecord.
func (s *Service) AddSignatureEvidenceRecord(signatureDocument, evidenceRecordDocument model.DSSDocument,
	parameters *EvidenceRecordIncorporationParameters) model.DSSDocument {
	if signatureDocument == nil {
		panic("The signature document cannot be null")
	}
	if evidenceRecordDocument == nil {
		panic("The evidence record document cannot be null")
	}

	builder := NewEmbeddedEvidenceRecordBuilder(s.CertificateVerifier)
	signatureWithEvidenceRecord, err := builder.AddEvidenceRecord(signatureDocument, evidenceRecordDocument, parameters)
	if err != nil {
		panic(err)
	}
	name, err := s.GetFinalFileName(signatureDocument, enumerations.SigningOperationAddEvidenceRecord)
	if err != nil {
		panic(err)
	}
	signatureWithEvidenceRecord.SetName(name)
	signatureWithEvidenceRecord.SetMimeType(signatureDocument.MimeType())
	return signatureWithEvidenceRecord
}

// cadesServiceAssertCounterSignaturePossible ports the private assertCounterSignaturePossible.
func cadesServiceAssertCounterSignaturePossible(parameters *CounterSignatureParameters) {
	if enumerations.SignatureLevelCAdESBaselineB != parameters.SignatureLevel() {
		panic(fmt.Sprintf("A counter signature with a level '%s' is not supported! "+
			"Please, use CAdES-BASELINE-B", parameters.SignatureLevel()))
	}
}

// Compile-time interface assertions, standing in for Java's "extends AbstractSignatureService
// ... implements CounterSignatureService, EvidenceRecordIncorporationService".
var (
	_ document.SignatureService[*SignatureParameters, *TimestampParameters]               = (*Service)(nil)
	_ document.CounterSignatureService[*CounterSignatureParameters]                       = (*Service)(nil)
	_ document.EvidenceRecordIncorporationService[*EvidenceRecordIncorporationParameters] = (*Service)(nil)
)
