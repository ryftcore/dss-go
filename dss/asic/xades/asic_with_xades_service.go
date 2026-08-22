// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/signature/ASiCWithXAdESService.java
// (DSS 6.5.RC1).
//
// slf4j logging is dropped (PORTING.md), as is java.io.Serializable/serialVersionUID.
//
// Java's `extends AbstractASiCSignatureService<...>` becomes embedding plus the
// InitAbstractASiCSignatureService(self) registration; the base's single-document convenience
// wrappers (promoted here as this type's plain-named GetContentTimestamp/GetDataToSign/
// SignDocument/Timestamp) dispatch into the *Multiple methods below, which carry the real,
// multi-document implementation. ExtendDocument
// has only ever one Java overload (single DSSDocument, shared by DocumentSignatureService and
// MultipleDocumentsSignatureService) so it needs no Multiple companion.
//
// Go has no overloading, so this type cannot itself satisfy both
// document.SignatureService[SP,TP] and document.MultipleDocumentsSignatureService[SP,TP]
// (they declare the same method names, single-document vs list-typed) the way Java's class does.
// This port keeps the plain names for the single-document shape - so *ASiCWithXAdESService
// satisfies document.SignatureService - and MultipleDocumentsService() below returns a
// thin adapter satisfying document.MultipleDocumentsSignatureService by forwarding to the
// *Multiple methods.
package xades

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
	dssxades "github.com/ryftcore/dss-go/dss/xades"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// init ports the static initializer block of ASiCWithXAdESService.
func init() {
	xmlutils.XPathUtilsRegisterNamespace(asic.ASiCManifestNS)
	xmlutils.XPathUtilsRegisterNamespace(ManifestNS)
}

// ASiCWithXAdESService contains the main methods for ASiC with XAdES signature
// creation/extension.
type ASiCWithXAdESService struct {
	asic.AbstractASiCSignatureService[*ASiCWithXAdESSignatureParameters, *dssxades.TimestampParameters,
		*dssxades.CounterSignatureParameters, *dssxades.EvidenceRecordIncorporationParameters]

	// asicFilenameFactory defines rules for filename creation for new ZIP entries (e.g.
	// signature files, etc.).
	asicFilenameFactory ASiCWithXAdESFilenameFactory
}

var (
	_ asic.AbstractASiCSignatureServiceOverrides[*ASiCWithXAdESSignatureParameters, *dssxades.TimestampParameters] = (*ASiCWithXAdESService)(nil)
	_ document.SignatureService[*ASiCWithXAdESSignatureParameters, *dssxades.TimestampParameters]                  = (*ASiCWithXAdESService)(nil)
	_ asic.EvidenceRecordIncorporationService[*dssxades.EvidenceRecordIncorporationParameters]                     = (*ASiCWithXAdESService)(nil)
)

// NewASiCWithXAdESService is the default constructor to instantiate the service. Ports
// ASiCWithXAdESService(CertificateVerifier).
func NewASiCWithXAdESService(certificateVerifier validation.CertificateVerifier) *ASiCWithXAdESService {
	// Upstream logs "+ ASiCService with XAdES created".
	service := &ASiCWithXAdESService{
		AbstractASiCSignatureService: asic.NewAbstractASiCSignatureService[*ASiCWithXAdESSignatureParameters,
			*dssxades.TimestampParameters, *dssxades.CounterSignatureParameters,
			*dssxades.EvidenceRecordIncorporationParameters](certificateVerifier),
		asicFilenameFactory: NewDefaultASiCWithXAdESFilenameFactory(),
	}
	service.InitAbstractASiCSignatureService(service)
	return service
}

// SetAsicFilenameFactory sets the ASiCWithXAdESFilenameFactory defining a set of rules for
// naming of newly created ZIP entries, such as signature files. Ports
// setAsicFilenameFactory(ASiCWithXAdESFilenameFactory).
//
// Panics with Java's message when asicFilenameFactory is nil (Objects.requireNonNull).
func (s *ASiCWithXAdESService) SetAsicFilenameFactory(asicFilenameFactory ASiCWithXAdESFilenameFactory) {
	if asicFilenameFactory == nil {
		panic("ASiCWithXAdESFilenameFactory cannot be null!")
	}
	s.asicFilenameFactory = asicFilenameFactory
}

// assertSigningCertificateValid raises a panic if the signing rules forbid the use of the
// certificate. Ports the inherited protected
// AbstractSignatureService#assertSigningCertificateValid(AbstractSignatureParameters<?>).
//
// Java declares the parameter with a wildcard, which matters here: this service's TP is
// TimestampParameters, while the parameters objects it receives extend
// SignatureParameters, i.e. AbstractSignatureParameters<TimestampParameters>. The Go
// port of the base method is bound to the service's own TP and so cannot accept them; this
// wrapper re-instantiates the frozen base with the XAdES timestamp-parameter type and delegates
// to it, rather than duplicating the check.
func (s *ASiCWithXAdESService) assertSigningCertificateValid(
	parameters *document.AbstractSignatureParameters[*dssxades.TimestampParameters]) {
	checker := document.NewAbstractSignatureService[*dssxades.SignatureParameters,
		*dssxades.TimestampParameters](s.CertificateVerifier)
	checker.AssertSigningCertificateValid(parameters)
}

// GetContentTimestampMultiple ports the @Override
// getContentTimestamp(List, ASiCWithXAdESSignatureParameters). Named Multiple (not the plain
// Java name) since document.DocumentSignatureService[SP,TP] and
// document.MultipleDocumentsSignatureService[SP,TP] both declare a method literally named
// GetContentTimestamp with different (single-document vs list) signatures - see the package doc
// comment.
//
// Panics with Java's message when parameters is nil (Objects.requireNonNull).
func (s *ASiCWithXAdESService) GetContentTimestampMultiple(toSignDocuments []model.DSSDocument,
	parameters *ASiCWithXAdESSignatureParameters) *validation.TimestampToken {
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	s.AssertSignaturePossible(toSignDocuments)

	asicContent := NewASiCWithXAdESASiCContentBuilder().
		Build(toSignDocuments, parameters.ASiC().ContainerType())
	dataToSignHelper := NewASiCWithXAdESDataToSignHelperBuilder(s.asicFilenameFactory).
		Build(asicContent, parameters)
	xadesParameters := s.getXAdESParameters(parameters, asicContent.SignatureDocuments(), dataToSignHelper.IsOpenDocument())
	return s.GetXAdESService().GetContentTimestampForDocuments(dataToSignHelper.ToBeSigned(), xadesParameters)
}

// GetDataToSignMultiple ports the @Override getDataToSign(List, ASiCWithXAdESSignatureParameters).
// Named Multiple for the same reason as GetContentTimestampMultiple above.
//
// Panics with Java's message when parameters is nil (Objects.requireNonNull).
func (s *ASiCWithXAdESService) GetDataToSignMultiple(toSignDocuments []model.DSSDocument,
	parameters *ASiCWithXAdESSignatureParameters) *model.ToBeSigned {
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	s.AssertSignaturePossible(toSignDocuments)
	s.assertSigningCertificateValid(&parameters.AbstractSignatureParameters)

	asicContent := NewASiCWithXAdESASiCContentBuilder().
		Build(toSignDocuments, parameters.ASiC().ContainerType())
	dataToSignHelper := NewASiCWithXAdESDataToSignHelperBuilder(s.asicFilenameFactory).
		Build(asicContent, parameters)
	xadesParameters := s.getXAdESParameters(parameters, asicContent.SignatureDocuments(), dataToSignHelper.IsOpenDocument())
	return s.GetXAdESService().GetDataToSignForDocuments(dataToSignHelper.ToBeSigned(), xadesParameters)
}

// SignDocumentMultiple ports the @Override
// signDocument(List, ASiCWithXAdESSignatureParameters, SignatureValue). Named Multiple for the
// same reason as GetContentTimestampMultiple above.
//
// Panics with Java's messages when a required argument is nil (Objects.requireNonNull).
func (s *ASiCWithXAdESService) SignDocumentMultiple(toSignDocuments []model.DSSDocument,
	parameters *ASiCWithXAdESSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
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

	asicContent := NewASiCWithXAdESASiCContentBuilder().
		Build(toSignDocuments, parameters.ASiC().ContainerType())
	dataToSignHelper := NewASiCWithXAdESDataToSignHelperBuilder(s.asicFilenameFactory).
		Build(asicContent, parameters)

	xadesParameters := s.getXAdESParameters(parameters, asicContent.SignatureDocuments(), dataToSignHelper.IsOpenDocument())
	newSignature := s.GetXAdESService().SignDocuments(dataToSignHelper.ToBeSigned(), xadesParameters, signatureValue)
	newSignature.SetName(s.asicFilenameFactory.SignatureFilename(asicContent))

	asicContent.SetSignatureDocuments(asic.UtilsAddOrReplaceDocument(asicContent.SignatureDocuments(), newSignature))

	var signingDate time.Time
	if bLevelSigningDate := parameters.BLevel().SigningDate(); bLevelSigningDate != nil {
		signingDate = *bLevelSigningDate
	}
	asicSignature := s.BuildASiCContainerAt(asicContent, signingDate)
	name, err := s.GetFinalDocumentNameWithMimeType(asicSignature, enumerations.SigningOperationSign,
		parameters.SignatureLevel(), asicSignature.MimeType())
	if err != nil {
		panic(err)
	}
	asicSignature.SetName(name)
	parameters.Reinit()
	return asicSignature
}

// TimestampMultiple ports the @Override timestamp(List, XAdESTimestampParameters). Named
// Multiple for the same reason as GetContentTimestampMultiple above.
//
// Panics unconditionally with Java's UnsupportedOperationException message: ASiC-S/E with
// XAdES does not support adding a detached timestamp file.
func (s *ASiCWithXAdESService) TimestampMultiple(toTimestampDocuments []model.DSSDocument,
	parameters *dssxades.TimestampParameters) model.DSSDocument {
	panic("Timestamp file cannot be added with ASiC-S/E + XAdES")
}

// ExtendDocument ports the @Override
// extendDocument(DSSDocument, ASiCWithXAdESSignatureParameters).
//
// Panics with Java's messages when a required argument is nil (Objects.requireNonNull), and
// with an *exception.IllegalInputException when the provided document is not an ASiC document.
func (s *ASiCWithXAdESService) ExtendDocument(toExtendDocument model.DSSDocument,
	parameters *ASiCWithXAdESSignatureParameters) model.DSSDocument {
	if toExtendDocument == nil {
		panic("toExtendDocument is not defined!")
	}
	if parameters == nil {
		panic("Cannot extend the signature. SignatureParameters are not defined!")
	}

	s.assertExtensionSupported(toExtendDocument)

	extensionHelper := NewASiCWithXAdESSignatureExtensionHelper(toExtendDocument)
	asicContent := extensionHelper.GetAsicContent()

	signatureDocuments := extensionHelper.GetSignatureDocuments()

	isOpenDocument, err := asic.UtilsIsOpenDocument(asicContent.MimeTypeDocument())
	if err != nil {
		panic(err)
	}

	parameters.SetSignaturePackaging(enumerations.SignaturePackagingDetached)
	parameters.GetContext().SetDetachedContents(s.GetDetachedContents(asicContent, isOpenDocument))

	for _, signature := range signatureDocuments {
		extendedDocument := s.GetXAdESService().ExtendDocument(signature, &parameters.SignatureParameters)
		extendedDocument.SetName(signature.Name())
		signatureDocuments = asic.UtilsAddOrReplaceDocument(signatureDocuments, extendedDocument)
		asicContent.SetSignatureDocuments(signatureDocuments)
	}

	var signingDate time.Time
	if bLevelSigningDate := parameters.BLevel().SigningDate(); bLevelSigningDate != nil {
		signingDate = *bLevelSigningDate
	}
	extensionResult := s.BuildASiCContainerAt(asicContent, signingDate)
	name, err := s.GetFinalDocumentNameWithMimeType(toExtendDocument, enumerations.SigningOperationExtend,
		parameters.SignatureLevel(), toExtendDocument.MimeType())
	if err != nil {
		panic(err)
	}
	extensionResult.SetName(name)

	parameters.Reinit()
	return extensionResult
}

// assertExtensionSupported ports the private assertExtensionSupported(DSSDocument).
func (s *ASiCWithXAdESService) assertExtensionSupported(toExtendDocument model.DSSDocument) {
	isASiC, err := asic.UtilsIsASiC(toExtendDocument)
	if err != nil {
		panic(err)
	}
	if !isASiC {
		panic(exception.NewIllegalInputException("The provided file is not ASiC document!"))
	}
}

// GetDetachedContents returns a detached contents to be used for a signature validation. Ports
// the protected getDetachedContents(Content, boolean).
func (s *ASiCWithXAdESService) GetDetachedContents(asicContent *asic.Content, isOpenDocument bool) []model.DSSDocument {
	if isOpenDocument {
		return OpenDocumentSupportUtilsGetOpenDocumentCoverage(asicContent)
	}
	return asicContent.SignedDocuments()
}

// GetXAdESService returns the XAdESService to be used for signing. Ports the protected
// getXAdESService().
func (s *ASiCWithXAdESService) GetXAdESService() *dssxades.Service {
	xadesService := dssxades.NewService(s.CertificateVerifier)
	xadesService.SetTspSource(s.TspSource)
	return xadesService
}

// getXAdESParameters returns an instance of SignatureParameters to be used for a signature
// file creation. Ports the private getXAdESParameters(ASiCWithXAdESSignatureParameters, List,
// boolean).
//
// Panics with an *exception.IllegalInputException when more than one signature file is present
// for an ASiC-S/OpenDocument container, or when the existing signature file is not valid XML.
func (s *ASiCWithXAdESService) getXAdESParameters(parameters *ASiCWithXAdESSignatureParameters,
	signatureDocuments []model.DSSDocument, openDocument bool) *dssxades.SignatureParameters {
	parameters.SetSignaturePackaging(enumerations.SignaturePackagingDetached)

	var rootDocument *xmldom.Node
	// If already existing signature file and ASiC-S OR OpenDocument type, we re-use the same
	// signature file
	if utils.IsCollectionNotEmpty(signatureDocuments) &&
		(enumerations.ASiCContainerTypeASiCS == parameters.ASiC().ContainerType() || openDocument) {
		if utils.CollectionSize(signatureDocuments) > 1 {
			panic(exception.NewIllegalInputException("Unable to choose signature file to add a new signature into! " +
				"Only one signature file shall be present for the particular container format."))
		}
		existingXAdESSignature := signatureDocuments[0]
		if !xmlutils.DomUtilsIsDOM(existingXAdESSignature) {
			panic(exception.NewIllegalInputException(fmt.Sprintf("The provided signature file '%s' is not a valid XML! "+
				"Unable to sign.", existingXAdESSignature.Name())))
		}
		var err error
		rootDocument, err = xmlutils.DomUtilsBuildDOMFromDocument(existingXAdESSignature)
		if err != nil {
			panic(err)
		}

	} else {
		// No signatures or ASiC-E
		rootDocument = s.buildDomRoot(openDocument)
	}

	parameters.SetRootDocument(rootDocument)
	return &parameters.SignatureParameters
}

// buildDomRoot ports the private buildDomRoot(boolean).
func (s *ASiCWithXAdESService) buildDomRoot(openDocument bool) *xmldom.Node {
	rootDocument := xmlutils.DomUtilsBuildDOMEmpty()

	var xadesSignatures *xmldom.Node
	if openDocument {
		xadesSignatures = xmldom.NewElement(xmldom.Name{
			Space: asic.ASiCManifestNamespaceLibreOfficeNS,
			Local: asic.ASiCManifestNamespaceLibreOfficeSignatures,
		})
	} else {
		xadesSignatures = xmlutils.DomUtilsCreateElementNS(rootDocument, asic.ASiCManifestNS, asic.ManifestElementXAdESSignatures)
	}
	rootDocument.AppendChild(xadesSignatures)
	return rootDocument
}

// GetArchiveExtractor ports the @Override protected getArchiveExtractor(DSSDocument).
//
// Java returns the DefaultASiCContainerExtractor supertype; the Go port hands back the base
// embedded (by value) in ASiCWithXAdESContainerExtractor, which already carries the concrete
// extractor as its registered overrides, so Extract() dispatches exactly as Java's does. The
// address-of is required here - unlike the CADSIGN chunk's analogous method, which returns the
// embedded value directly and does not type-check against its own *asic.DefaultContainerExtractor
// return type; not fixed here since asic/cades is a sibling chunk's file, flagged in notes
// instead per PORTING.md.
func (s *ASiCWithXAdESService) GetArchiveExtractor(archive model.DSSDocument) *asic.DefaultContainerExtractor {
	extractor := NewASiCWithXAdESContainerExtractor(archive)
	return &extractor.DefaultContainerExtractor
}

// AddSignaturePolicyStore incorporates a Signature Policy Store as an unsigned property into the
// ASiC with XAdES Signature. Ports addSignaturePolicyStore(DSSDocument, SignaturePolicyStore).
//
// Panics with Java's messages when a required argument is nil (Objects.requireNonNull).
func (s *ASiCWithXAdESService) AddSignaturePolicyStore(asicContainer model.DSSDocument,
	signaturePolicyStore *model.SignaturePolicyStore) model.DSSDocument {
	if asicContainer == nil {
		panic("The asicContainer cannot be null")
	}
	if signaturePolicyStore == nil {
		panic("The signaturePolicyStore cannot be null")
	}

	extensionHelper := NewASiCWithXAdESSignatureExtensionHelper(asicContainer)
	asicContent := extensionHelper.GetAsicContent()
	s.AssertAddSignaturePolicyStorePossible(asicContent)

	xadesService := s.GetXAdESService()

	signatureDocuments := extensionHelper.GetSignatureDocuments()
	for _, signature := range signatureDocuments {
		signatureWithPolicyStore := xadesService.AddSignaturePolicyStore(signature, signaturePolicyStore)
		signatureWithPolicyStore.SetName(signature.Name())
		signatureDocuments = asic.UtilsAddOrReplaceDocument(signatureDocuments, signatureWithPolicyStore)
		asicContent.SetSignatureDocuments(signatureDocuments)
	}

	// Java passes a null Date here; the zero time is what the ZIP container handler treats as
	// "not specified" (see ZipUtils#CreateZipArchiveFromEntriesAt).
	resultArchive := s.BuildASiCContainerAt(asicContent, time.Time{})
	name, err := s.GetFinalArchiveName(asicContainer, enumerations.SigningOperationAddSigPolicyStore,
		asicContainer.MimeType())
	if err != nil {
		panic(err)
	}
	resultArchive.SetName(name)
	return resultArchive
}

// GetDataToBeCounterSigned ports the @Override
// getDataToBeCounterSigned(DSSDocument, CounterSignatureParameters).
//
// Panics with Java's messages when a required argument is nil (Objects.requireNonNull).
func (s *ASiCWithXAdESService) GetDataToBeCounterSigned(asicContainer model.DSSDocument,
	parameters *dssxades.CounterSignatureParameters) *model.ToBeSigned {
	if asicContainer == nil {
		panic("asicContainer cannot be null!")
	}
	if parameters == nil {
		panic("SignatureParameters cannot be null!")
	}
	s.AssertCounterSignatureParametersValid(parameters)

	counterSignatureHelper := NewASiCWithXAdESSignatureExtensionHelper(asicContainer)
	signatureDocument := counterSignatureHelper.ExtractSignatureDocument(parameters.SignatureIdToCounterSign())

	xadesService := s.GetXAdESService()
	return xadesService.GetDataToBeCounterSigned(signatureDocument, parameters)
}

// CounterSignSignature ports the @Override
// counterSignSignature(DSSDocument, CounterSignatureParameters, SignatureValue).
//
// Panics with Java's messages when a required argument is nil (Objects.requireNonNull).
func (s *ASiCWithXAdESService) CounterSignSignature(asicContainer model.DSSDocument,
	parameters *dssxades.CounterSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
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

	counterSignatureHelper := NewASiCWithXAdESSignatureExtensionHelper(asicContainer)
	asicContent := counterSignatureHelper.GetAsicContent()

	signatureDocument := counterSignatureHelper.ExtractSignatureDocument(parameters.SignatureIdToCounterSign())

	xadesService := s.GetXAdESService()
	counterSignedSignature := xadesService.CounterSignSignature(signatureDocument, parameters, signatureValue)
	counterSignedSignature.SetName(signatureDocument.Name())
	asicContent.SetSignatureDocuments(asic.UtilsAddOrReplaceDocument(asicContent.SignatureDocuments(), counterSignedSignature))

	var signingDate time.Time
	if bLevelSigningDate := parameters.BLevel().SigningDate(); bLevelSigningDate != nil {
		signingDate = *bLevelSigningDate
	}
	resultArchive := s.BuildASiCContainerAt(asicContent, signingDate)
	name, err := s.GetFinalDocumentNameWithMimeType(asicContainer, enumerations.SigningOperationCounterSign,
		parameters.SignatureLevel(), asicContainer.MimeType())
	if err != nil {
		panic(err)
	}
	resultArchive.SetName(name)
	return resultArchive
}

// AddSignatureEvidenceRecord ports the @Override
// addSignatureEvidenceRecord(DSSDocument, DSSDocument, EvidenceRecordIncorporationParameters).
//
// Panics with Java's messages when a required argument is nil (Objects.requireNonNull).
func (s *ASiCWithXAdESService) AddSignatureEvidenceRecord(asicContainer model.DSSDocument,
	evidenceRecordDocument model.DSSDocument,
	parameters *dssxades.EvidenceRecordIncorporationParameters) model.DSSDocument {
	if asicContainer == nil {
		panic("The ASiC container cannot be null!")
	}
	if evidenceRecordDocument == nil {
		panic("The evidence record document cannot be null!")
	}
	if parameters == nil {
		panic("Parameters cannot be null!")
	}

	asicContainerHelper := NewASiCWithXAdESSignatureExtensionHelper(asicContainer)
	asicContent := asicContainerHelper.GetAsicContent()

	signatureDocument := asicContainerHelper.ExtractSignatureDocument(parameters.SignatureId())
	parameters.SetDetachedContents(asicContainerHelper.GetDetachedDocuments(signatureDocument.Name()))

	xadesService := s.GetXAdESService()
	signatureWithEvidenceRecord := xadesService.AddSignatureEvidenceRecord(signatureDocument, evidenceRecordDocument, parameters)
	signatureWithEvidenceRecord.SetName(signatureDocument.Name())
	asicContent.SetSignatureDocuments(asic.UtilsAddOrReplaceDocument(asicContent.SignatureDocuments(), signatureWithEvidenceRecord))

	resultArchive := s.BuildASiCContainer(asicContent)
	name, err := s.GetFinalArchiveName(asicContainer, enumerations.SigningOperationAddEvidenceRecord,
		asicContainer.MimeType())
	if err != nil {
		panic(err)
	}
	resultArchive.SetName(name)
	return resultArchive
}

// AddContainerEvidenceRecord ports the @Override
// addContainerEvidenceRecord(List, DSSDocument, ContainerEvidenceRecordParameters).
//
// Panics with Java's messages when a required argument is nil, undefined, or the document list
// is empty.
func (s *ASiCWithXAdESService) AddContainerEvidenceRecord(documents []model.DSSDocument,
	evidenceRecordDocument model.DSSDocument, parameters *asic.ContainerEvidenceRecordParameters) model.DSSDocument {
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

	asicContent, err := NewASiCWithXAdESContainerEvidenceRecordBuilder(s.CertificateVerifier, s.asicFilenameFactory).
		Build(documents, evidenceRecordDocument, parameters)
	if err != nil {
		panic(err)
	}
	asicContainer := s.BuildASiCContainer(asicContent)
	name, err := s.GetFinalArchiveName(asicContainer, enumerations.SigningOperationAddEvidenceRecord,
		asicContainer.MimeType())
	if err != nil {
		panic(err)
	}
	asicContainer.SetName(name)
	return asicContainer
}

// AddContainerEvidenceRecordMultiple delegates to AddContainerEvidenceRecord; it exists to
// satisfy asic.AbstractASiCSignatureServiceOverrides, whose name disambiguates Java's overload.
func (s *ASiCWithXAdESService) AddContainerEvidenceRecordMultiple(documents []model.DSSDocument,
	evidenceRecordDocument model.DSSDocument, parameters *asic.ContainerEvidenceRecordParameters) model.DSSDocument {
	return s.AddContainerEvidenceRecord(documents, evidenceRecordDocument, parameters)
}

// MultipleDocumentsService adapts this service to document.MultipleDocumentsSignatureService.
// See the package-level note: Java's ASiCWithXAdESService implements that interface directly,
// which Go cannot express on the same type because document.SignatureService declares
// the same four method names with single-document parameters.
func (s *ASiCWithXAdESService) MultipleDocumentsService() document.MultipleDocumentsSignatureService[
	*ASiCWithXAdESSignatureParameters, *dssxades.TimestampParameters] {
	return &asicWithXAdESServiceMultipleDocumentsAdapter{service: s}
}

// asicWithXAdESServiceMultipleDocumentsAdapter forwards MultipleDocumentsSignatureService to the
// *Multiple methods of ASiCWithXAdESService. It holds no state and adds no behaviour.
type asicWithXAdESServiceMultipleDocumentsAdapter struct {
	service *ASiCWithXAdESService
}

func (a *asicWithXAdESServiceMultipleDocumentsAdapter) GetContentTimestamp(toSignDocuments []model.DSSDocument,
	parameters *ASiCWithXAdESSignatureParameters) *validation.TimestampToken {
	return a.service.GetContentTimestampMultiple(toSignDocuments, parameters)
}

func (a *asicWithXAdESServiceMultipleDocumentsAdapter) GetDataToSign(toSignDocuments []model.DSSDocument,
	parameters *ASiCWithXAdESSignatureParameters) *model.ToBeSigned {
	return a.service.GetDataToSignMultiple(toSignDocuments, parameters)
}

func (a *asicWithXAdESServiceMultipleDocumentsAdapter) IsValidSignatureValue(toBeSigned *model.ToBeSigned,
	signatureValue *model.SignatureValue, signingCertificate *model.CertificateToken) bool {
	return a.service.IsValidSignatureValue(toBeSigned, signatureValue, signingCertificate)
}

func (a *asicWithXAdESServiceMultipleDocumentsAdapter) SignDocument(toSignDocuments []model.DSSDocument,
	parameters *ASiCWithXAdESSignatureParameters, signatureValue *model.SignatureValue) model.DSSDocument {
	return a.service.SignDocumentMultiple(toSignDocuments, parameters, signatureValue)
}

func (a *asicWithXAdESServiceMultipleDocumentsAdapter) ExtendDocument(toExtendDocument model.DSSDocument,
	parameters *ASiCWithXAdESSignatureParameters) model.DSSDocument {
	return a.service.ExtendDocument(toExtendDocument, parameters)
}

func (a *asicWithXAdESServiceMultipleDocumentsAdapter) Timestamp(toTimestampDocuments []model.DSSDocument,
	parameters *dssxades.TimestampParameters) model.DSSDocument {
	return a.service.TimestampMultiple(toTimestampDocuments, parameters)
}

var _ document.MultipleDocumentsSignatureService[*ASiCWithXAdESSignatureParameters, *dssxades.TimestampParameters] = (*asicWithXAdESServiceMultipleDocumentsAdapter)(nil)
