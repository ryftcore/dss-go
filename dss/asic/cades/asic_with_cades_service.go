// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/ASiCWithCAdESService.java (DSS 6.5.RC1).
//
// slf4j logging is dropped (PORTING.md), as are java.io.Serializable / serialVersionUID.
//
// Java's `extends AbstractASiCSignatureService<...>` becomes embedding plus the
// InitAbstractASiCSignatureService(self) registration; the base's single-document convenience
// wrappers (GetContentTimestamp, GetDataToSign, SignDocument, SignatureTimestamp - promoted here
// as this type's plain-named methods) dispatch into the *Multiple methods below
// (S7_BRIEF.md's virtual-dispatch warning), which carry the real, multi-document implementation.
//
// Go has no overloading, so this type cannot itself satisfy both
// document.DocumentSignatureService[SP,TP] (single-document GetContentTimestamp/GetDataToSign/
// SignDocument/Timestamp) and document.MultipleDocumentsSignatureService[SP,TP] (the same method
// names, list-typed) the way Java's class does. This port keeps the plain names for the
// single-document shape - so *ASiCWithCAdESService satisfies document.DocumentSignatureService -
// and MultipleDocumentsService() below returns a thin adapter satisfying
// document.MultipleDocumentsSignatureService by forwarding to the *Multiple methods, matching the
// precedent in xades/xades_service.go and jades/jades_service.go.
package cades

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/asic"
	dsscades "github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/signature/resources"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ASiCWithCAdESService contains the main methods for ASiC with CAdES signature
// creation/extension.
type ASiCWithCAdESService struct {
	asic.AbstractASiCSignatureService[*ASiCWithCAdESSignatureParameters, *ASiCWithCAdESTimestampParameters,
		*dsscades.CAdESCounterSignatureParameters, *dsscades.CAdESEvidenceRecordIncorporationParameters]

	// asicFilenameFactory defines rules for filename creation for new ZIP entries (e.g.
	// signature files, etc.).
	asicFilenameFactory ASiCWithCAdESFilenameFactory

	// ResourcesHandlerBuilder is used to create data container objects such as an OutputStream
	// or a DSSDocument. Java declares the field protected.
	ResourcesHandlerBuilder resources.DSSResourcesHandlerBuilder
}

var (
	_ asic.AbstractASiCSignatureServiceOverrides[*ASiCWithCAdESSignatureParameters, *ASiCWithCAdESTimestampParameters] = (*ASiCWithCAdESService)(nil)
	_ document.DocumentSignatureService[*ASiCWithCAdESSignatureParameters, *ASiCWithCAdESTimestampParameters]          = (*ASiCWithCAdESService)(nil)
	_ asic.EvidenceRecordIncorporationService[*dsscades.CAdESEvidenceRecordIncorporationParameters]                    = (*ASiCWithCAdESService)(nil)
)

// NewASiCWithCAdESService is the default constructor to instantiate the service. Ports
// ASiCWithCAdESService(CertificateVerifier).
func NewASiCWithCAdESService(certificateVerifier validation.CertificateVerifier) *ASiCWithCAdESService {
	// Upstream logs "+ ASiCService with CAdES created".
	service := &ASiCWithCAdESService{
		AbstractASiCSignatureService: asic.NewAbstractASiCSignatureService[*ASiCWithCAdESSignatureParameters,
			*ASiCWithCAdESTimestampParameters, *dsscades.CAdESCounterSignatureParameters,
			*dsscades.CAdESEvidenceRecordIncorporationParameters](certificateVerifier),
		asicFilenameFactory:     NewDefaultASiCWithCAdESFilenameFactory(),
		ResourcesHandlerBuilder: dsscades.CAdESUtilsDefaultResourcesHandlerBuilder,
	}
	service.InitAbstractASiCSignatureService(service)
	return service
}

// SetResourcesHandlerBuilder sets a DSSResourcesHandlerBuilder to be used for operating with
// internal objects during the signature creation procedure.
// NOTE: The DSSResourcesHandlerBuilder is supported only within the 'dss-cms-stream' module!
// Ports setResourcesHandlerBuilder(DSSResourcesHandlerBuilder).
func (s *ASiCWithCAdESService) SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) {
	s.ResourcesHandlerBuilder = cms.CMSUtilsResourcesHandlerBuilder(resourcesHandlerBuilder)
}

// SetAsicFilenameFactory sets the ASiCWithCAdESFilenameFactory defining a set of rules for
// naming of newly created ZIP entries, such as signature files. Ports
// setAsicFilenameFactory(ASiCWithCAdESFilenameFactory).
//
// Panics with Java's message when asicFilenameFactory is nil (Objects.requireNonNull).
func (s *ASiCWithCAdESService) SetAsicFilenameFactory(asicFilenameFactory ASiCWithCAdESFilenameFactory) {
	if asicFilenameFactory == nil {
		panic("ASiCWithCAdESFilenameFactory cannot be null!")
	}
	s.asicFilenameFactory = asicFilenameFactory
}

// assertSigningCertificateValid raises a panic if the signing rules forbid the use of the
// certificate. Ports the inherited protected
// AbstractSignatureService#assertSigningCertificateValid(AbstractSignatureParameters<?>).
//
// Java declares the parameter with a wildcard, which matters here: this service's TP is
// ASiCWithCAdESTimestampParameters, while the parameters objects it receives extend
// CAdESSignatureParameters, i.e. AbstractSignatureParameters<CAdESTimestampParameters>. The Go
// port of the base method is bound to the service's own TP and so cannot accept them; this
// wrapper re-instantiates the frozen base with the CAdES timestamp-parameter type and delegates
// to it, rather than duplicating the check.
func (s *ASiCWithCAdESService) assertSigningCertificateValid(
	parameters *document.AbstractSignatureParameters[*dsscades.CAdESTimestampParameters]) {
	checker := document.NewAbstractSignatureService[*dsscades.CAdESSignatureParameters,
		*dsscades.CAdESTimestampParameters](s.CertificateVerifier)
	checker.AssertSigningCertificateValid(parameters)
}

// GetContentTimestampMultiple ports the @Override
// getContentTimestamp(List, ASiCWithCAdESSignatureParameters). Named Multiple (not the plain
// Java name) because document.DocumentSignatureService[SP,TP] and
// document.MultipleDocumentsSignatureService[SP,TP] both declare a method literally named
// GetContentTimestamp with different (single-document vs list) signatures - Go cannot overload,
// so this type keeps the plain name reserved for the single-document convenience wrapper
// promoted from asic.AbstractASiCSignatureService, and exposes this multi-document shape via
// MultipleDocumentsService()'s adapter (see the bottom of this file), matching the precedent in
// xades/xades_service.go and jades/jades_service.go.
//
// Panics with Java's message when parameters is nil (Objects.requireNonNull).
func (s *ASiCWithCAdESService) GetContentTimestampMultiple(toSignDocuments []model.DSSDocument,
	parameters *ASiCWithCAdESSignatureParameters) *validation.TimestampToken {
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	s.AssertSignaturePossible(toSignDocuments)

	asicContent := NewASiCWithCAdESASiCContentBuilder().
		Build(toSignDocuments, parameters.ASiC().ContainerType())
	dataToSignHelper := NewASiCWithCAdESSignatureDataToSignHelperBuilder(s.asicFilenameFactory).
		Build(asicContent, parameters)
	toBeSigned := dataToSignHelper.ToBeSigned()
	return s.GetCAdESService().GetContentTimestamp(toBeSigned, &parameters.CAdESSignatureParameters)
}

// GetDataToSignMultiple ports the @Override getDataToSign(List, ASiCWithCAdESSignatureParameters).
// Named Multiple for the same reason as GetContentTimestampMultiple above.
//
// Panics with Java's message when parameters is nil (Objects.requireNonNull).
func (s *ASiCWithCAdESService) GetDataToSignMultiple(toSignDocuments []model.DSSDocument,
	parameters *ASiCWithCAdESSignatureParameters) *model.ToBeSigned {
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	s.AssertSignaturePossible(toSignDocuments)
	s.assertSigningCertificateValid(&parameters.AbstractSignatureParameters)

	asicContent := NewASiCWithCAdESASiCContentBuilder().
		Build(toSignDocuments, parameters.ASiC().ContainerType())
	dataToSignHelper := NewASiCWithCAdESSignatureDataToSignHelperBuilder(s.asicFilenameFactory).
		Build(asicContent, parameters)
	s.assertSignaturePossibleWithTimestamps(asicContent.TimestampDocuments(), parameters.ASiC())

	cadesParameters := s.GetCAdESParametersWithHelper(parameters, dataToSignHelper)

	toBeSigned := dataToSignHelper.ToBeSigned()
	return s.GetCAdESService().GetDataToSign(toBeSigned, cadesParameters)
}

// SignDocumentMultiple ports the @Override
// signDocument(List, ASiCWithCAdESSignatureParameters, SignatureValue). Named Multiple for the
// same reason as GetContentTimestampMultiple above.
//
// Panics with Java's messages when a required argument is nil (Objects.requireNonNull).
func (s *ASiCWithCAdESService) SignDocumentMultiple(toSignDocuments []model.DSSDocument,
	parameters *ASiCWithCAdESSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	if toSignDocuments == nil {
		panic("toSignDocument cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	if signatureValue == nil {
		panic("SignatureValue cannot be null!")
	}
	s.AssertSignaturePossible(toSignDocuments)
	s.assertSigningCertificateValid(&parameters.AbstractSignatureParameters)

	asicContent := NewASiCWithCAdESASiCContentBuilder().
		Build(toSignDocuments, parameters.ASiC().ContainerType())
	dataToSignHelper := NewASiCWithCAdESSignatureDataToSignHelperBuilder(s.asicFilenameFactory).
		Build(asicContent, parameters)
	asicParameters := parameters.ASiC()
	s.assertSignaturePossibleWithTimestamps(asicContent.TimestampDocuments(), asicParameters)

	cadesParameters := s.GetCAdESParametersWithHelper(parameters, dataToSignHelper)

	// Archive Timestamp in case of ASiC-E is not embedded into the CAdES signature
	addASiCArchiveManifest := s.isAddASiCEArchiveManifest(parameters.SignatureLevel(), parameters.ASiC().ContainerType())
	if addASiCArchiveManifest {
		cadesParameters.SetSignatureLevel(enumerations.SignatureLevel_CAdES_BASELINE_LT)
	}

	toBeSigned := dataToSignHelper.ToBeSigned()
	if enumerations.ASiCContainerType_ASiC_E == asicParameters.ContainerType() {
		// XML Document in case of ASiC-E container
		asicContent.SetManifestDocuments(append(asicContent.ManifestDocuments(), toBeSigned))
	}

	signature := s.GetCAdESService().SignDocument(toBeSigned, cadesParameters, signatureValue)
	signature.SetName(s.asicFilenameFactory.SignatureFilename(asicContent))

	asicContent.SetSignatureDocuments(asic.ASiCUtilsAddOrReplaceDocument(asicContent.SignatureDocuments(), signature))

	if addASiCArchiveManifest {
		// Java's try/finally restores the LTA level even when the extension fails.
		func() {
			defer cadesParameters.SetSignatureLevel(enumerations.SignatureLevel_CAdES_BASELINE_LTA)
			extensionProfile := s.GetLTALevelExtensionProfile()
			asicContent = extensionProfile.Extend(asicContent, cadesParameters)
		}()
	}

	asicContainer := s.BuildASiCContainerAt(asicContent, parameters.ZipCreationDate())
	name, err := s.GetFinalDocumentNameWithMimeType(asicContainer, enumerations.SigningOperation_SIGN,
		parameters.SignatureLevel(), asicContainer.MimeType())
	if err != nil {
		panic(err)
	}
	asicContainer.SetName(name)
	parameters.Reinit()
	return asicContainer
}

// TimestampMultiple ports the @Override timestamp(List, ASiCWithCAdESTimestampParameters).
// Named Multiple for the same reason as GetContentTimestampMultiple above.
//
// Panics with Java's messages when parameters is nil or the document list is empty.
func (s *ASiCWithCAdESService) TimestampMultiple(toTimestampDocuments []model.DSSDocument,
	parameters *ASiCWithCAdESTimestampParameters) model.DSSDocument {
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	if utils.IsCollectionEmpty(toTimestampDocuments) {
		panic("List of documents to be timestamped cannot be empty!")
	}

	asicContent := NewASiCWithCAdESASiCContentBuilder().
		Build(toTimestampDocuments, parameters.ASiC().ContainerType())

	asicParameters := parameters.ASiC()
	signatureDocuments := asicContent.SignatureDocuments()

	s.assertTimestampPossible(signatureDocuments, asicParameters)

	signatures := asicContent.SignatureDocuments()
	timestampDocuments := asicContent.TimestampDocuments()
	if (utils.IsCollectionNotEmpty(signatures) || utils.IsCollectionNotEmpty(timestampDocuments)) &&
		s.isLtaExtensionPossible(asicContent) {
		toTimestampDocument := toTimestampDocuments[0]

		extensionProfile := s.GetLTALevelExtensionProfile()
		asicContent = extensionProfile.ExtendWithDigestAlgorithm(asicContent, parameters.DigestAlgorithm())

		extensionResult := s.BuildASiCContainerAt(asicContent, parameters.ZipCreationDate())
		name, err := s.GetFinalDocumentNameWithMimeType(toTimestampDocument,
			enumerations.SigningOperation_TIMESTAMP, "", toTimestampDocument.MimeType())
		if err != nil {
			panic(err)
		}
		extensionResult.SetName(name)
		return extensionResult
	}

	timestampService := NewASiCWithCAdESTimestampServiceWithFilenameFactory(s.TspSource, s.asicFilenameFactory)
	asicContent = timestampService.Timestamp(asicContent, parameters)

	asicContainer := s.BuildASiCContainerAt(asicContent, parameters.ZipCreationDate())
	name, err := s.GetFinalDocumentNameWithMimeType(asicContainer,
		enumerations.SigningOperation_TIMESTAMP, "", asicContainer.MimeType())
	if err != nil {
		panic(err)
	}
	asicContainer.SetName(name)
	return asicContainer
}

// ExtendDocument ports the @Override
// extendDocument(DSSDocument, ASiCWithCAdESSignatureParameters).
//
// Panics with Java's messages when a required argument is nil, and with an
// *exception.IllegalInputException when the container type cannot be extracted.
func (s *ASiCWithCAdESService) ExtendDocument(toExtendDocument model.DSSDocument,
	parameters *ASiCWithCAdESSignatureParameters) model.DSSDocument {
	if toExtendDocument == nil {
		panic("toExtendDocument is not defined!")
	}
	if parameters == nil {
		panic("Cannot extend the signature. SignatureParameters are not defined!")
	}

	s.assertExtensionSupported(toExtendDocument)
	asicContent := s.ExtractCurrentArchive(toExtendDocument)

	s.assertValidSignaturesToExtendFound(asicContent.SignatureDocuments())

	containerType := asicContent.ContainerType()
	if containerType == "" {
		panic(exception.NewIllegalInputException(
			"The container type of the provided document is not supported or cannot be extracted!"))
	}

	extensionProfile := s.GetExtensionProfile(parameters.SignatureLevel(), containerType)
	asicContent = extensionProfile.Extend(asicContent, &parameters.CAdESSignatureParameters)

	extensionResult := s.BuildASiCContainerAt(asicContent, parameters.ZipCreationDate())
	name, err := s.GetFinalDocumentNameWithMimeType(toExtendDocument, enumerations.SigningOperation_EXTEND,
		parameters.SignatureLevel(), toExtendDocument.MimeType())
	if err != nil {
		panic(err)
	}
	extensionResult.SetName(name)
	return extensionResult
}

// assertExtensionSupported ports the private assertExtensionSupported(DSSDocument).
func (s *ASiCWithCAdESService) assertExtensionSupported(toExtendDocument model.DSSDocument) {
	isASiC, err := asic.ASiCUtilsIsASiC(toExtendDocument)
	if err != nil {
		panic(err)
	}
	if !isASiC {
		panic(exception.NewIllegalInputException("The provided file is not ASiC document!"))
	}
}

// assertValidSignaturesToExtendFound ports the private assertValidSignaturesToExtendFound(List).
func (s *ASiCWithCAdESService) assertValidSignaturesToExtendFound(signatureDocuments []model.DSSDocument) {
	if utils.IsCollectionEmpty(signatureDocuments) {
		panic(exception.NewIllegalInputException("No supported signature documents found! Unable to extend the container."))
	}
}

// GetArchiveExtractor ports the @Override protected getArchiveExtractor(DSSDocument).
//
// Java returns the DefaultASiCContainerExtractor supertype; the Go port hands back the base
// embedded in ASiCWithCAdESContainerExtractor, which already carries the concrete extractor as
// its registered overrides, so Extract() dispatches exactly as Java's does.
func (s *ASiCWithCAdESService) GetArchiveExtractor(archive model.DSSDocument) *asic.DefaultASiCContainerExtractor {
	extractor := NewASiCWithCAdESContainerExtractor(archive)
	return &extractor.DefaultASiCContainerExtractor
}

// GetCAdESService returns the CAdESService to be used for signature/timestamp creation. Ports
// the protected getCAdESService().
func (s *ASiCWithCAdESService) GetCAdESService() *dsscades.CAdESService {
	cadesService := dsscades.NewCAdESService(s.CertificateVerifier)
	cadesService.SetTspSource(s.TspSource)
	if dsscades.CAdESUtilsDefaultResourcesHandlerBuilder != s.ResourcesHandlerBuilder {
		cadesService.SetResourcesHandlerBuilder(s.ResourcesHandlerBuilder)
	}
	return cadesService
}

// GetCAdESParameters returns CAdESSignatureParameters from the given
// ASiCWithCAdESSignatureParameters. Ports the protected
// getCAdESParameters(ASiCWithCAdESSignatureParameters).
//
// Java returns the very same object upcast to its superclass; the Go port returns the address of
// the embedded CAdESSignatureParameters, preserving that aliasing.
func (s *ASiCWithCAdESService) GetCAdESParameters(
	parameters *ASiCWithCAdESSignatureParameters) *dsscades.CAdESSignatureParameters {
	parameters.SetSignaturePackaging(enumerations.SignaturePackaging_DETACHED)
	parameters.GetContext().SetDetachedContents(nil)
	return &parameters.CAdESSignatureParameters
}

// GetCAdESParametersWithHelper returns CAdESSignatureParameters from the given
// ASiCWithCAdESSignatureParameters and GetDataToSignASiCWithCAdESHelper. Ports the protected
// getCAdESParameters(ASiCWithCAdESSignatureParameters, GetDataToSignASiCWithCAdESHelper).
func (s *ASiCWithCAdESService) GetCAdESParametersWithHelper(parameters *ASiCWithCAdESSignatureParameters,
	dataToSignHelper GetDataToSignASiCWithCAdESHelper) *dsscades.CAdESSignatureParameters {
	cadesParameters := s.GetCAdESParameters(parameters)
	cadesParameters.GetContext().SetDetachedContents(dataToSignHelper.DetachedContents())
	// no detached content is used in case of file signing
	cadesParameters.SetParallelSignature(utils.IsCollectionNotEmpty(dataToSignHelper.DetachedContents()))
	return cadesParameters
}

// isAddASiCEArchiveManifest ports the private
// isAddASiCEArchiveManifest(SignatureLevel, ASiCContainerType).
func (s *ASiCWithCAdESService) isAddASiCEArchiveManifest(signatureLevel enumerations.SignatureLevel,
	containerType enumerations.ASiCContainerType) bool {
	return enumerations.SignatureLevel_CAdES_BASELINE_LTA == signatureLevel &&
		enumerations.ASiCContainerType_ASiC_E == containerType
}

// AddSignaturePolicyStore incorporates a Signature Policy Store as an unsigned property into the
// ASiC with CAdES Signature. Ports
// addSignaturePolicyStore(DSSDocument, SignaturePolicyStore).
//
// Panics with Java's messages when a required argument is nil (Objects.requireNonNull).
func (s *ASiCWithCAdESService) AddSignaturePolicyStore(asicContainer model.DSSDocument,
	signaturePolicyStore *model.SignaturePolicyStore) model.DSSDocument {
	if asicContainer == nil {
		panic("The asicContainer cannot be null")
	}
	if signaturePolicyStore == nil {
		panic("The signaturePolicyStore cannot be null")
	}

	asicContent := s.ExtractCurrentArchive(asicContainer)
	s.AssertAddSignaturePolicyStorePossible(asicContent)

	cadesService := s.GetCAdESService()

	signatureDocuments := asicContent.SignatureDocuments()
	// Ensure iteration not over original list
	iterated := make([]model.DSSDocument, len(signatureDocuments))
	copy(iterated, signatureDocuments)
	for _, signature := range iterated {
		signatureWithPolicyStore := cadesService.AddSignaturePolicyStore(signature, signaturePolicyStore)
		signatureWithPolicyStore.SetName(signature.Name())
		signatureDocuments = asic.ASiCUtilsAddOrReplaceDocument(signatureDocuments, signatureWithPolicyStore)
		asicContent.SetSignatureDocuments(signatureDocuments)
	}

	// Java passes a null Date here; the zero time is what the ZIP container handler treats as
	// "not specified" (see ZipUtils#CreateZipArchiveFromEntriesAt).
	resultArchive := s.BuildASiCContainerAt(asicContent, time.Time{})
	name, err := s.GetFinalArchiveName(asicContainer, enumerations.SigningOperation_ADD_SIG_POLICY_STORE,
		asicContainer.MimeType())
	if err != nil {
		panic(err)
	}
	resultArchive.SetName(name)
	return resultArchive
}

// AssertAddSignaturePolicyStorePossible ports the @Override protected
// assertAddSignaturePolicyStorePossible(ASiCContent).
//
// Panics with an *exception.IllegalInputException when a signature is covered by another
// manifest.
func (s *ASiCWithCAdESService) AssertAddSignaturePolicyStorePossible(asicContent *asic.ASiCContent) {
	s.AbstractASiCSignatureService.AssertAddSignaturePolicyStorePossible(asicContent)

	signatureDocuments := asicContent.SignatureDocuments()
	for _, signature := range signatureDocuments {
		if asic.ASiCUtilsIsCoveredByManifest(asicContent.AllManifestDocuments(), signature.Name()) {
			panic(exception.NewIllegalInputException(fmt.Sprintf("Not possible to add a signature policy store! "+
				"Reason : a signature with a filename '%s' is covered by another manifest.", signature.Name())))
		}
	}
}

// GetDataToBeCounterSigned ports the @Override
// getDataToBeCounterSigned(DSSDocument, CAdESCounterSignatureParameters).
//
// Panics with Java's messages when a required argument is nil (Objects.requireNonNull).
func (s *ASiCWithCAdESService) GetDataToBeCounterSigned(asicContainer model.DSSDocument,
	parameters *dsscades.CAdESCounterSignatureParameters) *model.ToBeSigned {
	if asicContainer == nil {
		panic("asicContainer cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	s.assertSigningCertificateValid(&parameters.AbstractSignatureParameters)
	s.AssertCounterSignatureParametersValid(parameters)

	counterSignatureHelper := NewASiCWithCAdESSignatureExtensionHelper(asicContainer)
	signatureDocument := counterSignatureHelper.ExtractSignatureDocument(parameters.SignatureIdToCounterSign())

	counterSignatureBuilder := s.GetCAdESCounterSignatureBuilder()
	counterSignatureBuilder.SetManifestFile(counterSignatureHelper.GetManifestFile(signatureDocument.Name()))
	counterSignatureBuilder.SetResourcesHandlerBuilder(s.ResourcesHandlerBuilder)

	signerInfoToCounterSign, err := counterSignatureBuilder.GetSignerInformationToBeCounterSigned(signatureDocument, parameters)
	if err != nil {
		panic(err)
	}

	cadesService := s.GetCAdESService()
	return cadesService.GetDataToBeCounterSignedForSigner(signerInfoToCounterSign, &parameters.CAdESSignatureParameters)
}

// CounterSignSignature ports the @Override
// counterSignSignature(DSSDocument, CAdESCounterSignatureParameters, SignatureValue).
//
// Panics with Java's messages when a required argument is nil (Objects.requireNonNull).
func (s *ASiCWithCAdESService) CounterSignSignature(asicContainer model.DSSDocument,
	parameters *dsscades.CAdESCounterSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	if asicContainer == nil {
		panic("asicContainer cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	if signatureValue == nil {
		panic("signatureValue cannot be null!")
	}
	s.AssertCounterSignatureParametersValid(parameters)

	counterSignatureHelper := NewASiCWithCAdESSignatureExtensionHelper(asicContainer)
	asicContent := counterSignatureHelper.GetAsicContent()

	signatureDocument := counterSignatureHelper.ExtractSignatureDocument(parameters.SignatureIdToCounterSign())
	originalCMS, err := cms.CMSUtilsParseToCMS(signatureDocument)
	if err != nil {
		panic(err)
	}

	counterSignatureBuilder := s.GetCAdESCounterSignatureBuilder()
	counterSignatureBuilder.SetManifestFile(counterSignatureHelper.GetManifestFile(signatureDocument.Name()))

	counterSignedSignature, err := counterSignatureBuilder.AddCounterSignature(originalCMS, parameters, signatureValue)
	if err != nil {
		panic(err)
	}
	counterSignedSignature.SetName(signatureDocument.Name())
	asicContent.SetSignatureDocuments(asic.ASiCUtilsAddOrReplaceDocument(asicContent.SignatureDocuments(), counterSignedSignature))

	// Java passes bLevel().getSigningDate() straight through; BLevelParameters#SigningDate
	// returns a *time.Time in this port, and a nil one keeps the zero time the ZIP container
	// handler treats as "not specified" (see ZipUtils#CreateZipArchiveFromEntriesAt).
	var signingDate time.Time
	if bLevelSigningDate := parameters.BLevel().SigningDate(); bLevelSigningDate != nil {
		signingDate = *bLevelSigningDate
	}
	resultArchive := s.BuildASiCContainerAt(asicContent, signingDate)
	name, err := s.GetFinalDocumentNameWithMimeType(asicContainer, enumerations.SigningOperation_COUNTER_SIGN,
		parameters.SignatureLevel(), asicContainer.MimeType())
	if err != nil {
		panic(err)
	}
	resultArchive.SetName(name)
	return resultArchive
}

// AddSignatureEvidenceRecord ports the @Override
// addSignatureEvidenceRecord(DSSDocument, DSSDocument, CAdESEvidenceRecordIncorporationParameters).
//
// Panics with Java's messages when a required argument is nil (Objects.requireNonNull).
func (s *ASiCWithCAdESService) AddSignatureEvidenceRecord(asicContainer model.DSSDocument,
	evidenceRecordDocument model.DSSDocument,
	parameters *dsscades.CAdESEvidenceRecordIncorporationParameters) model.DSSDocument {
	if asicContainer == nil {
		panic("The ASiC container cannot be null!")
	}
	if evidenceRecordDocument == nil {
		panic("The evidence record document cannot be null!")
	}
	if parameters == nil {
		panic("Parameters cannot be null!")
	}

	asicContainerHelper := NewASiCWithCAdESSignatureExtensionHelper(asicContainer)
	asicContent := asicContainerHelper.GetAsicContent()

	signatureDocument := asicContainerHelper.ExtractSignatureDocument(parameters.SignatureId())
	parameters.SetDetachedContents(asicContainerHelper.GetDetachedDocuments(signatureDocument.Name()))

	builder := dsscades.NewCAdESEmbeddedEvidenceRecordBuilder(s.CertificateVerifier)
	builder.SetManifestFile(asicContainerHelper.GetManifestFile(signatureDocument.Name()))
	signatureWithEvidenceRecord, err := builder.AddEvidenceRecord(signatureDocument, evidenceRecordDocument, parameters)
	if err != nil {
		panic(err)
	}
	signatureWithEvidenceRecord.SetName(signatureDocument.Name())
	asicContent.SetSignatureDocuments(asic.ASiCUtilsAddOrReplaceDocument(asicContent.SignatureDocuments(), signatureWithEvidenceRecord))

	resultArchive := s.BuildASiCContainer(asicContent)
	name, err := s.GetFinalArchiveName(asicContainer, enumerations.SigningOperation_ADD_EVIDENCE_RECORD,
		asicContainer.MimeType())
	if err != nil {
		panic(err)
	}
	resultArchive.SetName(name)
	return resultArchive
}

// AddContainerEvidenceRecord ports the @Override
// addContainerEvidenceRecord(List, DSSDocument, ASiCContainerEvidenceRecordParameters).
//
// Panics with Java's messages when a required argument is nil or the document list is empty.
func (s *ASiCWithCAdESService) AddContainerEvidenceRecord(documents []model.DSSDocument,
	evidenceRecordDocument model.DSSDocument, parameters *asic.ASiCContainerEvidenceRecordParameters) model.DSSDocument {
	if evidenceRecordDocument == nil {
		panic("The evidence record document cannot be null!")
	}
	if parameters == nil {
		panic("Parameters cannot be null!")
	}
	if parameters.ContainerType() == "" {
		panic("ASiCContainerType must be defined!")
	}
	if utils.IsCollectionEmpty(documents) {
		panic("List of documents cannot be empty!")
	}

	asicContent, err := NewASiCWithCAdESContainerEvidenceRecordBuilder(s.CertificateVerifier, s.asicFilenameFactory).
		Build(documents, evidenceRecordDocument, parameters)
	if err != nil {
		panic(err)
	}
	asicContainer := s.BuildASiCContainer(asicContent)
	name, err := s.GetFinalArchiveName(asicContainer, enumerations.SigningOperation_ADD_EVIDENCE_RECORD,
		asicContainer.MimeType())
	if err != nil {
		panic(err)
	}
	asicContainer.SetName(name)
	return asicContainer
}

// AddContainerEvidenceRecordMultiple delegates to AddContainerEvidenceRecord; it exists to
// satisfy asic.AbstractASiCSignatureServiceOverrides, whose name disambiguates Java's overload.
func (s *ASiCWithCAdESService) AddContainerEvidenceRecordMultiple(documents []model.DSSDocument,
	evidenceRecordDocument model.DSSDocument, parameters *asic.ASiCContainerEvidenceRecordParameters) model.DSSDocument {
	return s.AddContainerEvidenceRecord(documents, evidenceRecordDocument, parameters)
}

// GetExtensionProfile returns the extension profile to be used for the current signature. Ports
// the protected getExtensionProfile(SignatureLevel, ASiCContainerType).
//
// Panics with Java's messages when the signature level is undefined or unsupported.
func (s *ASiCWithCAdESService) GetExtensionProfile(signatureLevel enumerations.SignatureLevel,
	containerType enumerations.ASiCContainerType) ASiCWithCAdESSignatureExtender {
	if signatureLevel == "" {
		panic("SignatureLevel must be defined!")
	}
	var extensionProfile ASiCWithCAdESSignatureExtender
	switch signatureLevel {
	case enumerations.SignatureLevel_CAdES_BASELINE_T, enumerations.SignatureLevel_CAdES_BASELINE_LT:
		extensionProfile = NewASiCWithCAdESSignatureExtension(s.CertificateVerifier, s.TspSource)
	case enumerations.SignatureLevel_CAdES_BASELINE_LTA:
		if enumerations.ASiCContainerType_ASiC_E == containerType {
			extensionProfile = NewASiCWithCAdESLevelBaselineLTAWithFilenameFactory(
				s.CertificateVerifier, s.TspSource, s.asicFilenameFactory)
		} else {
			extensionProfile = NewASiCWithCAdESSignatureExtension(s.CertificateVerifier, s.TspSource)
		}
	default:
		panic(fmt.Sprintf("Unsupported signature format '%s' for extension.", signatureLevel))
	}
	extensionProfile.SetResourcesHandlerBuilder(s.ResourcesHandlerBuilder)
	return extensionProfile
}

// GetLTALevelExtensionProfile gets an extension profile for an LTA-level signature augmentation.
// Ports the protected getLTALevelExtensionProfile().
func (s *ASiCWithCAdESService) GetLTALevelExtensionProfile() *ASiCWithCAdESLevelBaselineLTA {
	extensionProfile := NewASiCWithCAdESLevelBaselineLTAWithFilenameFactory(
		s.CertificateVerifier, s.TspSource, s.asicFilenameFactory)
	extensionProfile.SetResourcesHandlerBuilder(s.ResourcesHandlerBuilder)
	return extensionProfile
}

// GetCAdESCounterSignatureBuilder creates a new instance of CAdESCounterSignatureBuilder. Ports
// the protected getCAdESCounterSignatureBuilder().
func (s *ASiCWithCAdESService) GetCAdESCounterSignatureBuilder() *dsscades.CAdESCounterSignatureBuilder {
	counterSignatureBuilder := dsscades.NewCAdESCounterSignatureBuilder(s.CertificateVerifier)
	counterSignatureBuilder.SetResourcesHandlerBuilder(s.ResourcesHandlerBuilder)
	return counterSignatureBuilder
}

// AssertCounterSignatureParametersValid ports the @Override protected
// assertCounterSignatureParametersValid(CAdESCounterSignatureParameters).
//
// Panics with Java's UnsupportedOperationException message on an unsupported level.
func (s *ASiCWithCAdESService) AssertCounterSignatureParametersValid(
	parameters *dsscades.CAdESCounterSignatureParameters) {
	s.AbstractASiCSignatureService.AssertCounterSignatureParametersValid(parameters)

	if enumerations.SignatureLevel_CAdES_BASELINE_B != parameters.SignatureLevel() {
		panic(fmt.Sprintf("A counter signature with a level '%s' is not supported! "+
			"Please, use CAdES-BASELINE-B", parameters.SignatureLevel()))
	}
}

// assertSignaturePossibleWithTimestamps ports the private
// assertSignaturePossible(List, ASiCParameters); the name disambiguates it from the base's
// assertSignaturePossible(List) overload.
//
// Panics with an *exception.IllegalInputException when time assertion files are present.
func (s *ASiCWithCAdESService) assertSignaturePossibleWithTimestamps(timestampDocuments []model.DSSDocument,
	asicParameters *asic.ASiCParameters) {
	if asic.ASiCUtilsIsASiCS(asicParameters) && utils.IsCollectionNotEmpty(timestampDocuments) {
		panic(exception.NewIllegalInputException(
			"Unable to sign an ASiC-S with CAdES container containing time assertion files!"))
	}
}

// assertTimestampPossible ports the private assertTimestampPossible(List, ASiCParameters).
//
// Panics with an *exception.IllegalInputException when signature files are present.
func (s *ASiCWithCAdESService) assertTimestampPossible(signatureDocuments []model.DSSDocument,
	asicParameters *asic.ASiCParameters) {
	if asic.ASiCUtilsIsASiCS(asicParameters) && utils.IsCollectionNotEmpty(signatureDocuments) {
		panic(exception.NewIllegalInputException(
			"Unable to timestamp an ASiC-S with CAdES container containing signature files! " +
				"Use extendDocument(...) method for signature extension."))
	}
}

// MultipleDocumentsService adapts this service to document.MultipleDocumentsSignatureService.
// See the package-level note: Java's ASiCWithCAdESService implements that interface directly,
// which Go cannot express on the same type because document.DocumentSignatureService declares
// the same four method names with single-document parameters.
func (s *ASiCWithCAdESService) MultipleDocumentsService() document.MultipleDocumentsSignatureService[
	*ASiCWithCAdESSignatureParameters, *ASiCWithCAdESTimestampParameters] {
	return &asicWithCAdESServiceMultipleDocumentsAdapter{service: s}
}

// asicWithCAdESServiceMultipleDocumentsAdapter forwards MultipleDocumentsSignatureService to the
// *Multiple methods of ASiCWithCAdESService. It holds no state and adds no behaviour.
type asicWithCAdESServiceMultipleDocumentsAdapter struct {
	service *ASiCWithCAdESService
}

func (a *asicWithCAdESServiceMultipleDocumentsAdapter) GetContentTimestamp(toSignDocuments []model.DSSDocument,
	parameters *ASiCWithCAdESSignatureParameters) *validation.TimestampToken {
	return a.service.GetContentTimestampMultiple(toSignDocuments, parameters)
}

func (a *asicWithCAdESServiceMultipleDocumentsAdapter) GetDataToSign(toSignDocuments []model.DSSDocument,
	parameters *ASiCWithCAdESSignatureParameters) *model.ToBeSigned {
	return a.service.GetDataToSignMultiple(toSignDocuments, parameters)
}

func (a *asicWithCAdESServiceMultipleDocumentsAdapter) IsValidSignatureValue(toBeSigned *model.ToBeSigned,
	signatureValue *model.SignatureValue, signingCertificate *model.CertificateToken) bool {
	return a.service.IsValidSignatureValue(toBeSigned, signatureValue, signingCertificate)
}

func (a *asicWithCAdESServiceMultipleDocumentsAdapter) SignDocument(toSignDocuments []model.DSSDocument,
	parameters *ASiCWithCAdESSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	return a.service.SignDocumentMultiple(toSignDocuments, parameters, signatureValue)
}

func (a *asicWithCAdESServiceMultipleDocumentsAdapter) ExtendDocument(toExtendDocument model.DSSDocument,
	parameters *ASiCWithCAdESSignatureParameters) model.DSSDocument {
	return a.service.ExtendDocument(toExtendDocument, parameters)
}

func (a *asicWithCAdESServiceMultipleDocumentsAdapter) Timestamp(toTimestampDocuments []model.DSSDocument,
	parameters *ASiCWithCAdESTimestampParameters) model.DSSDocument {
	return a.service.TimestampMultiple(toTimestampDocuments, parameters)
}

var _ document.MultipleDocumentsSignatureService[*ASiCWithCAdESSignatureParameters, *ASiCWithCAdESTimestampParameters] = (*asicWithCAdESServiceMultipleDocumentsAdapter)(nil)
