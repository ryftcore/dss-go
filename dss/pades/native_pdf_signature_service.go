// Native implementation of eu.europa.esig.dss.pdf.PDFSignatureService on internal/pdf.
//
// It merges the two upstream classes that only exist separately because upstream has two
// interchangeable backends: eu.europa.esig.dss.pdf.AbstractPDFSignatureService (the
// format-independent flow: the message-digest cache, the revision extraction, the /ByteRange
// validation, the developer extensions, the modification analysis) and
// eu.europa.esig.dss.pdf.pdfbox.PdfBoxSignatureService (the library-specific half: creating the
// signature dictionary, saving incrementally, building the /DSS dictionary, adding a signature
// field). Each method below names the upstream method it reproduces.
//
// The PDF mechanics are internal/pdf's and are never re-implemented here: pdf.Updater.AddSignature
// already performs PDDocument.addSignature + PdfBoxSignatureService.createSignatureDictionary,
// including the /Contents and /ByteRange reservation, the field/AcroForm/page wiring, the
// FieldMDP built from a field /Lock and the DocMDP permission (internal/pdf/DESIGN.md §3.3), and
// pdf.Updater.SetDSSDictionary writes /Root /DSS (§3.5). Re-encryption of an encrypted source
// document happens automatically inside the Updater (DESIGN.md Amendment 1), which is why
// PdfBoxSignatureService's SecureRandomProvider / checkEncryptedAndSaveIncrementally seam has no
// counterpart here.
//
// Not supported, per internal/pdf/DESIGN.md §0.2 (no rasteriser):
// previewPageWithVisualSignature, previewSignatureField and the visual half of a visible
// signature - the widget, its /Rect and an empty appearance XObject are produced, the image/text
// content is not.
package pades

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/pdf"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/pades/exception"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/signature/resources"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// NativePDFSignatureService signs, timestamps and extends PDF documents with the native
// internal/pdf engine.
type NativePDFSignatureService struct {
	// serviceMode tells whether this instance produces a signature or a document-timestamp
	// revision.
	serviceMode PDFServiceMode

	// signatureDrawerFactory builds the drawer used for a visual signature/timestamp creation.
	signatureDrawerFactory SignatureDrawerFactory

	// ResourcesHandlerBuilder creates a DSSResourcesHandler for each internal call, defining how
	// the produced documents are materialised.
	//
	// Default: an in-memory handler.
	ResourcesHandlerBuilder resources.DSSResourcesHandlerBuilder

	// PdfDifferencesFinder finds the differences occurred between PDF revisions (e.g. visible
	// changes). Default: DefaultPdfDifferencesFinder.
	PdfDifferencesFinder PdfDifferencesFinder

	// PdfObjectModificationsFinder finds the differences within internal PDF objects occurred
	// between PDF revisions. Default: DefaultPdfObjectModificationsFinder.
	PdfObjectModificationsFinder PdfObjectModificationsFinder

	// PdfPermissionsChecker verifies the PDF document permissions regarding a new signature
	// creation.
	PdfPermissionsChecker *PdfPermissionsChecker

	// PdfSignatureFieldPositionChecker verifies the signature field position placement validity.
	PdfSignatureFieldPositionChecker *PdfSignatureFieldPositionChecker

	// PdfMemoryUsageSetting specifies the load mode of the PDF document. It is accepted and
	// ignored; see native_pdf_document_reader.go.
	PdfMemoryUsageSetting PdfMemoryUsageSetting
}

// NewNativePDFSignatureService creates the service for the given mode and drawer factory.
// Port of the AbstractPDFSignatureService(PDFServiceMode, SignatureDrawerFactory) constructor and
// of PdfBoxSignatureService(PDFServiceMode, PdfBoxSignatureDrawerFactory).
func NewNativePDFSignatureService(serviceMode PDFServiceMode,
	signatureDrawerFactory SignatureDrawerFactory) *NativePDFSignatureService {
	if serviceMode == "" {
		panic("The PDFServiceMode shall be defined!")
	}
	if signatureDrawerFactory == nil {
		panic("The SignatureDrawerFactory shall be defined!")
	}
	return &NativePDFSignatureService{
		serviceMode:                      serviceMode,
		signatureDrawerFactory:           signatureDrawerFactory,
		ResourcesHandlerBuilder:          PAdESUtilsDefaultResourcesHandlerBuilder,
		PdfDifferencesFinder:             NewDefaultPdfDifferencesFinder(),
		PdfObjectModificationsFinder:     NewDefaultPdfObjectModificationsFinder(),
		PdfPermissionsChecker:            NewPdfPermissionsChecker(),
		PdfSignatureFieldPositionChecker: NewPdfSignatureFieldPositionChecker(),
		PdfMemoryUsageSetting:            PAdESUtilsDefaultPdfMemoryUsageSetting,
	}
}

// SetResourcesHandlerBuilder ports #setResourcesHandlerBuilder.
func (s *NativePDFSignatureService) SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) {
	if resourcesHandlerBuilder == nil {
		panic("DSSResourcesFactoryBuilder cannot be null!")
	}
	s.ResourcesHandlerBuilder = resourcesHandlerBuilder
}

// SetPdfDifferencesFinder ports #setPdfDifferencesFinder.
func (s *NativePDFSignatureService) SetPdfDifferencesFinder(pdfDifferencesFinder PdfDifferencesFinder) {
	if pdfDifferencesFinder == nil {
		panic("PdfDifferencesFinder cannot be null!")
	}
	s.PdfDifferencesFinder = pdfDifferencesFinder
}

// SetPdfObjectModificationsFinder ports #setPdfObjectModificationsFinder.
func (s *NativePDFSignatureService) SetPdfObjectModificationsFinder(pdfObjectModificationsFinder PdfObjectModificationsFinder) {
	if pdfObjectModificationsFinder == nil {
		panic("PdfObjectModificationsFinder cannot be null!")
	}
	s.PdfObjectModificationsFinder = pdfObjectModificationsFinder
}

// SetPdfPermissionsChecker ports #setPdfPermissionsChecker.
func (s *NativePDFSignatureService) SetPdfPermissionsChecker(pdfPermissionsChecker *PdfPermissionsChecker) {
	if pdfPermissionsChecker == nil {
		panic("PdfPermissionsChecker cannot be null!")
	}
	s.PdfPermissionsChecker = pdfPermissionsChecker
}

// SetPdfSignatureFieldPositionChecker ports #setPdfSignatureFieldPositionChecker.
func (s *NativePDFSignatureService) SetPdfSignatureFieldPositionChecker(pdfSignatureFieldPositionChecker *PdfSignatureFieldPositionChecker) {
	if pdfSignatureFieldPositionChecker == nil {
		panic("PdfSignatureFieldPositionChecker cannot be null!")
	}
	s.PdfSignatureFieldPositionChecker = pdfSignatureFieldPositionChecker
}

// SetPdfMemoryUsageSetting ports #setPdfMemoryUsageSetting.
func (s *NativePDFSignatureService) SetPdfMemoryUsageSetting(pdfMemoryUsageSetting PdfMemoryUsageSetting) {
	s.PdfMemoryUsageSetting = pdfMemoryUsageSetting
}

// LoadSignatureDrawer returns a SignatureDrawer initialised from the configured factory.
// Port of the protected #loadSignatureDrawer.
func (s *NativePDFSignatureService) LoadSignatureDrawer(imageParameters *SignatureImageParameters) SignatureDrawer {
	signatureDrawer := s.signatureDrawerFactory.SignatureDrawer(imageParameters)
	if signatureDrawer == nil {
		panic("SignatureDrawer shall be defined for the used SignatureDrawerFactory!")
	}
	return signatureDrawer
}

// MessageDigest returns the message-digest computed on the PDF signature revision's ByteRange,
// caching it in the parameters. Port of #messageDigest.
func (s *NativePDFSignatureService) MessageDigest(toSignDocument model.DSSDocument,
	parameters CommonParameters) model.DSSMessageDigest {
	if toSignDocument == nil {
		panic("DSSDocument shall be provided!")
	}
	if parameters == nil {
		panic("PAdESCommonParameters cannot be null!")
	}

	pdfSignatureCache := parameters.PdfSignatureCache()
	if pdfSignatureCache.MessageDigest().Value() == nil {
		pdfSignatureCache.SetMessageDigest(s.ComputeDigest(toSignDocument, parameters))
	}
	return pdfSignatureCache.MessageDigest()
}

// ComputeDigest computes the digest of the to-be-signed data of the revision the given parameters
// describe, and caches the laid-out revision in the parameters' PdfSignatureCache.
// Port of the abstract #computeDigest, implemented by PdfBoxSignatureService#computeDigest.
func (s *NativePDFSignatureService) ComputeDigest(toSignDocument model.DSSDocument,
	parameters CommonParameters) model.DSSMessageDigest {
	toBeSignedDocument, messageDigest := s.buildRevision(toSignDocument, nil, parameters, true)
	// cache the computed document
	parameters.PdfSignatureCache().SetToBeSignedDocument(toBeSignedDocument)
	return messageDigest
}

// Sign signs a PDF document, reusing the cached to-be-signed revision when there is one.
// Port of #sign.
func (s *NativePDFSignatureService) Sign(toSignDocument model.DSSDocument, cmsSignedData []byte,
	parameters CommonParameters) model.DSSDocument {
	if toSignDocument == nil {
		panic("DSSDocument shall be provided!")
	}
	if cmsSignedData == nil {
		panic("CMSSignedData cannot be null!")
	}
	if parameters == nil {
		panic("PAdESCommonParameters cannot be null!")
	}

	pdfSignatureCache := parameters.PdfSignatureCache()
	var signedDocument model.DSSDocument
	if pdfSignatureCache.ToBeSignedDocument() != nil {
		replaced, err := UtilsReplaceSignature(pdfSignatureCache.ToBeSignedDocument(), cmsSignedData,
			s.ResourcesHandlerBuilder)
		if err == nil {
			signedDocument = replaced
		}
		// Upstream logs "Unable to sign document using a resources caching! Reason : '{}'. Sign
		// using a complete processing..." and falls through to the complete processing.
	}
	parameters.Reinit()

	if signedDocument == nil {
		signedDocument = s.SignDocument(toSignDocument, cmsSignedData, parameters)
	}
	signedDocument.SetMimeType(enumerations.MimeTypeEnumPDF)
	return signedDocument
}

// SignDocument creates a signed document from the original toSignDocument, incorporating a new
// revision enveloping the provided CMS signed data.
// Port of the abstract #signDocument, implemented by PdfBoxSignatureService#signDocument.
func (s *NativePDFSignatureService) SignDocument(toSignDocument model.DSSDocument, cmsSignedData []byte,
	parameters CommonParameters) model.DSSDocument {
	s.AssertContentSizeSufficient(cmsSignedData, parameters)
	signedDocument, _ := s.buildRevision(toSignDocument, cmsSignedData, parameters, false)
	return signedDocument
}

// buildRevision lays out the new signature revision and, when computeDigest is set, digests the
// two spans named by /ByteRange. It is PdfBoxSignatureService#signDocumentAndReturnDigest, with
// PDDocument.addSignature + saveIncremental replaced by pdf.Updater.
func (s *NativePDFSignatureService) buildRevision(toSignDocument model.DSSDocument, cmsSignedData []byte,
	parameters CommonParameters, computeDigest bool) (model.DSSDocument, model.DSSMessageDigest) {
	reader, err := NewNativePdfDocumentReader(toSignDocument, parameters.PasswordProtection())
	if err != nil {
		panic(err)
	}
	defer func() { _ = reader.Close() }()

	fieldParameters := parameters.ImageParameters().FieldParameters()
	s.CheckPdfPermissions(reader, fieldParameters)

	updater, err := pdf.NewUpdater(reader.Document())
	if err != nil {
		panic(err)
	}

	options := s.signatureOptions(reader, parameters, fieldParameters)
	if _, err = updater.AddSignature(options); err != nil {
		panic(err)
	}

	s.applyPendingVersion(updater, reader)
	s.DigitalSignatureEnhancement(updater, reader, parameters)

	result, err := updater.Write()
	if err != nil {
		panic(model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to compute digest for a PDF : %s", err.Error()), err))
	}

	messageDigest := model.CreateEmptyDSSMessageDigest()
	if computeDigest {
		digestAlgorithm := parameters.DigestAlgorithm()
		digest, err := spi.DSSUtilsMessageDigest(digestAlgorithm)
		if err != nil {
			panic(err)
		}
		if _, err = io.Copy(digest, result.SignedData()); err != nil {
			panic(err)
		}
		messageDigest = model.NewDSSMessageDigestWithValue(digestAlgorithm, digest.Sum(nil))
	}

	if len(cmsSignedData) > 0 {
		if err = result.InsertContents(cmsSignedData); err != nil {
			panic(model.NewDSSErrorMessageCause(
				fmt.Sprintf("Unable to save a document. Reason : %s", err.Error()), err))
		}
	}

	signedDocument := model.NewInMemoryDocument(result.Bytes)
	signedDocument.SetMimeType(enumerations.MimeTypeEnumPDF)
	return signedDocument, messageDigest
}

// signatureOptions translates the PAdES parameters into internal/pdf's SignatureOptions. It is
// PdfBoxSignatureService#createSignatureDictionary plus the SignatureOptions configuration of
// #signDocumentAndReturnDigest; the /Reference (FieldMDP, DocMDP) and /V wiring the two perform
// by hand is internal/pdf's AddSignature (DESIGN.md §3.3).
func (s *NativePDFSignatureService) signatureOptions(reader *NativePdfDocumentReader,
	parameters CommonParameters, fieldParameters *SignatureFieldParameters) pdf.SignatureOptions {
	options := pdf.SignatureOptions{
		Type:        pdf.Name(s.Type()),
		ContentSize: parameters.ContentSize(),
		FieldID:     fieldParameters.FieldId(),
		Page:        fieldParameters.Page(),
	}
	if utils.IsStringNotEmpty(parameters.Filter()) {
		options.Filter = pdf.Name(parameters.Filter())
	}
	// sub-filter for basic and PAdES Part 2 signatures
	if utils.IsStringNotEmpty(parameters.SubFilter()) {
		options.SubFilter = pdf.Name(parameters.SubFilter())
	}
	if utils.IsStringNotEmpty(parameters.AppName()) {
		options.AppName = parameters.AppName()
	}

	if !s.IsDocumentTimestampLayer() {
		if signatureParameters, ok := parameters.(*SignatureParameters); ok {
			if utils.IsStringNotEmpty(signatureParameters.SignerName()) {
				options.SignerName = signatureParameters.SignerName()
			}
			if utils.IsStringNotEmpty(signatureParameters.ContactInfo()) {
				options.ContactInfo = signatureParameters.ContactInfo()
			}
			if utils.IsStringNotEmpty(signatureParameters.Location()) {
				options.Location = signatureParameters.Location()
			}
			if utils.IsStringNotEmpty(signatureParameters.Reason()) {
				options.Reason = signatureParameters.Reason()
			}
			// A document can contain only one signature field that contains a DocMDP transform
			// method; it shall be the first signed field in the document. AddSignature enforces
			// the "first signed field" half.
			if permission := signatureParameters.Permission(); permission != "" {
				options.DocMDP = permission.Code()
			}
			// the signing date, needed for a valid signature
			options.SigningTime = signatureParameters.SigningDate().In(signatureParameters.SigningTimeZone())
		}
	}

	// the document needs to have an ID, if not the current system time is used, and then the
	// digest of the signed data will be different
	if id := reader.Document().ID(); id[0] == nil && id[1] == nil {
		options.DocumentID = reader.GenerateDocumentID(parameters)
	}

	imageParameters := parameters.ImageParameters()
	if !imageParameters.IsEmpty() {
		signatureDrawer := s.LoadSignatureDrawer(imageParameters)
		annotationBox := AnnotationBox{}
		if utils.IsStringEmpty(fieldParameters.FieldId()) {
			// check the signature field position only for new annotations
			box := s.VisibleSignatureFieldBoxPositionForDrawer(signatureDrawer, reader, fieldParameters)
			if box != nil {
				annotationBox = *box
			}
		}
		options.Rect = pdf.Rect{
			MinX: float64(annotationBox.MinX()), MinY: float64(annotationBox.MinY()),
			MaxX: float64(annotationBox.MaxX()), MaxY: float64(annotationBox.MaxY()),
		}
		options.Appearance = nativePDFSignatureServiceAppearance(options.Rect)
		// Upstream calls signatureDrawer.draw() here; the native drawer cannot paint (see
		// native_signature_drawer.go), so the field is created with an empty appearance stream
		// instead of failing the whole signing operation.
	}

	return options
}

// nativePDFSignatureServiceAppearance builds the minimal /AP /N form XObject of a visible
// signature field: a correctly sized, empty appearance. It is the native stand-in for what
// PdfBoxUtils.createSignatureAppearanceDictionary and the drawers produce.
func nativePDFSignatureServiceAppearance(rect pdf.Rect) *pdf.Stream {
	dictionary := pdf.DictOf(
		pdf.Name("Type"), pdf.Name("XObject"),
		pdf.Name("Subtype"), pdf.Name("Form"),
		pdf.Name("FormType"), pdf.Integer(1),
		pdf.Name("BBox"), pdf.Array{
			pdf.Integer(0), pdf.Integer(0),
			pdf.Real{Val: rect.Width()}, pdf.Real{Val: rect.Height()},
		},
		pdf.Name("Resources"), pdf.NewDict(),
	)
	return pdf.NewStream(dictionary, nil)
}

// applyPendingVersion writes the version recorded by PdfDocumentReader#setVersion into the
// updated catalog. Upstream's PdfBoxDocumentReader#setVersion mutates the live catalog, which an
// append-only writer cannot do; see NativePdfDocumentReader.SetVersion.
func (s *NativePDFSignatureService) applyPendingVersion(updater *pdf.Updater, reader *NativePdfDocumentReader) {
	if reader.pendingVersion == nil {
		return
	}
	// pdfDocument.getDocumentCatalog().setVersion(Float.toString(version)) - a name object.
	updater.Catalog().Set(pdf.Name(PAdESConstantsVersionName),
		pdf.Name(pdf.FormatReal(float64(*reader.pendingVersion))))
}

// IsDocumentTimestampLayer reports whether a DocumentTimestamp has to be added in the current
// mode; a CONTENT_TIMESTAMP is part of the signature. Port of the protected
// #isDocumentTimestampLayer.
func (s *NativePDFSignatureService) IsDocumentTimestampLayer() bool {
	return PDFServiceModeSignatureTimestamp == s.serviceMode ||
		PDFServiceModeArchiveTimestamp == s.serviceMode
}

// Type returns the /Type of the signature dictionary. Port of the protected #getType.
func (s *NativePDFSignatureService) Type() string {
	if s.IsDocumentTimestampLayer() {
		return PAdESConstantsTimestampType
	}
	return PAdESConstantsSignatureType
}

// DigitalSignatureEnhancement ensures the PDF document structure carries the developer
// extensions the signature's functionalities require.
// Port of the protected #digitalSignatureEnhancement.
func (s *NativePDFSignatureService) DigitalSignatureEnhancement(updater *pdf.Updater,
	reader *NativePdfDocumentReader, parameters CommonParameters) {
	if s.IsDocumentTimestampLayer() {
		s.EnsureESICDeveloperExtension1(updater, reader)
	}
	if s.IsCAdESDetached(parameters) {
		s.EnsureESICDeveloperExtension2(updater, reader)
	}
	if s.IsISO32001(parameters) {
		s.EnsureISO32001DeveloperExtension(updater, reader)
	}
	if s.IsISO32002(parameters) {
		s.EnsureISO32002DeveloperExtension(updater, reader)
	}
}

// EnsureESICDeveloperExtension1 adds the ESIC/ADBE developer extension of level 1 when it is
// missing. Port of the protected #ensureESICDeveloperExtension1.
func (s *NativePDFSignatureService) EnsureESICDeveloperExtension1(updater *pdf.Updater,
	reader *NativePdfDocumentReader) {
	// the standard says the /BaseVersion shall be less than or equal to both the document header
	// and catalog version, so skip adding the extension when the document's version is lower
	if reader.PdfHeaderVersion() < 1.7 || reader.Version() < 1.7 {
		return
	}
	// Skip inclusion of the dictionary, as the properties are already defined in PDF 2.0
	// (ISO 32000-2:2020)
	if reader.Version() >= 2.0 {
		return
	}
	esicExtension := nativePDFSignatureServiceDeveloperExtension("1.7", 1, "", "", "")
	adbeExtension := nativePDFSignatureServiceDeveloperExtension("1.7", 8, "", "", "")
	if !s.IsDeveloperExtensionPresent(reader, "ESIC", esicExtension) &&
		!s.IsDeveloperExtensionPresent(reader, "ADBE", adbeExtension) {
		s.AddDeveloperExtension(updater, reader, "ADBE", adbeExtension)
	}
}

// EnsureESICDeveloperExtension2 adds the ESIC/ADBE developer extension of level 2 when it is
// missing. Port of the protected #ensureESICDeveloperExtension2.
func (s *NativePDFSignatureService) EnsureESICDeveloperExtension2(updater *pdf.Updater,
	reader *NativePdfDocumentReader) {
	if reader.PdfHeaderVersion() < 1.7 || reader.Version() < 1.7 {
		return
	}
	if reader.Version() >= 2.0 {
		return
	}
	esicExtension := nativePDFSignatureServiceDeveloperExtension("1.7", 2, "", "", "")
	adbeExtension := nativePDFSignatureServiceDeveloperExtension("1.7", 8, "", "", "")
	if !s.IsDeveloperExtensionPresent(reader, "ESIC", esicExtension) &&
		!s.IsDeveloperExtensionPresent(reader, "ADBE", adbeExtension) {
		s.AddDeveloperExtension(updater, reader, "ADBE", adbeExtension)
	}
}

// EnsureISO32001DeveloperExtension adds the ISO 32001 developer extension when it is missing.
// Port of the protected #ensureISO_32001DeveloperExtension.
func (s *NativePDFSignatureService) EnsureISO32001DeveloperExtension(updater *pdf.Updater,
	reader *NativePdfDocumentReader) {
	if reader.PdfHeaderVersion() < 2.0 || reader.Version() < 2.0 {
		return
	}
	developerExtension := nativePDFSignatureServiceDeveloperExtension("2.0", 32001, ":2022",
		"DeveloperExtensions", "https://www.iso.org/standard/45874.html")
	if !s.IsDeveloperExtensionPresent(reader, "ISO_", developerExtension) {
		s.AddDeveloperExtension(updater, reader, "ISO_", developerExtension)
	}
}

// EnsureISO32002DeveloperExtension adds the ISO 32002 developer extension when it is missing.
// Port of the protected #ensureISO_32002DeveloperExtension.
func (s *NativePDFSignatureService) EnsureISO32002DeveloperExtension(updater *pdf.Updater,
	reader *NativePdfDocumentReader) {
	if reader.PdfHeaderVersion() < 2.0 || reader.Version() < 2.0 {
		return
	}
	developerExtension := nativePDFSignatureServiceDeveloperExtension("2.0", 32002, ":2022",
		"DeveloperExtensions", "https://www.iso.org/standard/45875.html")
	if !s.IsDeveloperExtensionPresent(reader, "ISO_", developerExtension) {
		s.AddDeveloperExtension(updater, reader, "ISO_", developerExtension)
	}
}

// nativePDFSignatureServiceDeveloperExtension creates a developer extension dictionary with the
// given configuration; an empty string or a zero extension level omits the entry, which is what
// upstream's null checks do. Port of the protected #createDeveloperExtensionDict, which builds it
// through PdfDocumentReader#createPdfDict; the native flow works on pdf.Dict directly, since the
// dictionary is written by the Updater and never read back through the SPI.
func nativePDFSignatureServiceDeveloperExtension(baseVersion string, extensionLevel int,
	extensionRevision, extensionType, url string) *pdf.Dict {
	dictionary := pdf.NewDict()
	if baseVersion != "" {
		dictionary.Set(pdf.Name(PAdESConstantsBaseVersionName), pdf.Name(baseVersion))
	}
	if extensionLevel != 0 {
		dictionary.Set(pdf.Name(PAdESConstantsExtensionLevelName), pdf.Integer(extensionLevel))
	}
	if extensionRevision != "" {
		dictionary.Set(pdf.Name(PAdESConstantsExtensionRevisionName),
			pdf.String{Bytes: []byte(extensionRevision)})
	}
	if extensionType != "" {
		dictionary.Set(pdf.Name(PAdESConstantsTypeName), pdf.Name(extensionType))
	}
	if url != "" {
		dictionary.Set(pdf.Name(PAdESConstantsURLName), pdf.String{Bytes: []byte(url)})
	}
	return dictionary
}

// IsCAdESDetached verifies whether the signature is created with the "ETSI.CAdES.detached"
// SubFilter. Port of the protected #isCAdESDetached.
func (s *NativePDFSignatureService) IsCAdESDetached(parameters CommonParameters) bool {
	return PAdESConstantsSignatureDefaultSubFilter == parameters.SubFilter()
}

// IsISO32001 verifies whether the ISO 32001 developer extension shall be included.
// Port of the protected #isISO_32001.
func (s *NativePDFSignatureService) IsISO32001(parameters CommonParameters) bool {
	subFilter := parameters.SubFilter()
	digestAlgorithm := parameters.DigestAlgorithm()
	return (PAdESConstantsSignaturePKCS7SubFilter == subFilter ||
		PAdESConstantsSignatureDefaultSubFilter == subFilter ||
		PAdESConstantsTimestampDefaultSubFilter == subFilter) &&
		(enumerations.DigestAlgorithmSHA3256 == digestAlgorithm ||
			enumerations.DigestAlgorithmSHA3384 == digestAlgorithm ||
			enumerations.DigestAlgorithmSHA3512 == digestAlgorithm ||
			enumerations.DigestAlgorithmSHAKE256 == digestAlgorithm)
}

// IsISO32002 verifies whether the ISO 32002 developer extension shall be included.
// Port of the protected #isISO_32002.
//
// Upstream's notes are kept: ECDSA elliptic curves are not covered yet, and ISO 32002 mistakenly
// refers to id-shake256 instead of id-shake256-len for Ed448
// (https://github.com/pdf-association/pdf-issues/issues/404); the extension is not enforced for
// id-shake256-len, to stay compliant with the current version of ISO 32002.
func (s *NativePDFSignatureService) IsISO32002(parameters CommonParameters) bool {
	subFilter := parameters.SubFilter()
	digestAlgorithm := parameters.DigestAlgorithm()
	return (PAdESConstantsSignaturePKCS7SubFilter == subFilter ||
		PAdESConstantsSignatureDefaultSubFilter == subFilter ||
		PAdESConstantsTimestampDefaultSubFilter == subFilter) &&
		enumerations.EncryptionAlgorithmEDDSA == parameters.EncryptionAlgorithm() &&
		(enumerations.DigestAlgorithmSHA512 == digestAlgorithm ||
			enumerations.DigestAlgorithmSHAKE256 == digestAlgorithm)
}

// IsDeveloperExtensionPresent verifies whether the specified developer extension is present in
// the document's catalog; the extension shall fully match the defined parameters.
// Port of the protected #isDeveloperExtensionPresent.
func (s *NativePDFSignatureService) IsDeveloperExtensionPresent(reader *NativePdfDocumentReader,
	prefix string, developerExtension *pdf.Dict) bool {
	document := reader.Document()
	catalog, err := document.Catalog()
	if err != nil {
		return false
	}
	extensionsDict, ok := document.GetDict(catalog, pdf.Name(PAdESConstantsExtensionsName))
	if !ok {
		return false
	}
	// can be an array or a dictionary
	if extensionDictArray, ok := document.GetArray(extensionsDict, pdf.Name(prefix)); ok {
		for i := range extensionDictArray {
			extensionDict, ok := document.Resolve(extensionDictArray[i]).(*pdf.Dict)
			if ok && nativePDFSignatureServiceMatch(document, extensionDict, developerExtension) {
				return true
			}
		}
		return false
	}
	if extensionDict, ok := document.GetDict(extensionsDict, pdf.Name(prefix)); ok {
		return nativePDFSignatureServiceMatch(document, extensionDict, developerExtension)
	}
	return false
}

// nativePDFSignatureServiceMatch is PdfDict#match: every entry of the expected dictionary must be
// present in the candidate with an equal value.
func nativePDFSignatureServiceMatch(document *pdf.Document, candidate, expected *pdf.Dict) bool {
	if candidate == nil {
		return false
	}
	for _, key := range expected.Keys() {
		targetObject := document.Resolve(expected.GetRaw(key))
		currentObject := document.Resolve(candidate.GetRaw(key))
		if !nativePdfIsNull(targetObject) && !nativePdfObjectEquals(targetObject, currentObject) {
			return false
		}
	}
	return true
}

// AddDeveloperExtension adds the given developer extension under the given prefix.
// Port of the protected #addDeveloperExtension.
//
// DEVIATION: upstream toggles COSBase#setDirect to control whether the extension dictionary is
// written inline; internal/pdf derives that from the object graph, so the /Extensions dictionary
// and its entries are written as direct objects inside the updated catalog - which is exactly the
// shape upstream aims for with its setDirect(true) calls ("add directly to ensure better
// compatibility (PDF 1.7)").
func (s *NativePDFSignatureService) AddDeveloperExtension(updater *pdf.Updater,
	reader *NativePdfDocumentReader, prefix string, developerExtension *pdf.Dict) {
	document := reader.Document()
	catalog := updater.Catalog()

	extensionsDict := pdf.NewDict()
	if existing, ok := document.GetDict(catalog, pdf.Name(PAdESConstantsExtensionsName)); ok {
		extensionsDict = existing.Clone()
	}
	catalog.Set(pdf.Name(PAdESConstantsExtensionsName), extensionsDict)

	if existingArray, ok := document.GetArray(extensionsDict, pdf.Name(prefix)); ok {
		extensionsDict.Set(pdf.Name(prefix),
			append(append(pdf.Array{}, existingArray...), developerExtension))
		return
	}
	if existingDictionary, ok := document.GetDict(extensionsDict, pdf.Name(prefix)); ok {
		extensionsDict.Set(pdf.Name(prefix), pdf.Array{existingDictionary, developerExtension})
		return
	}
	// add directly to ensure better compatibility (PDF 1.7)
	extensionsDict.Set(pdf.Name(prefix), developerExtension)
}

// GetRevisions retrieves the revisions of a PDF document: one entry per signature,
// document-timestamp and DSS dictionary update, ordered from the latest to the first.
// Port of #getRevisions.
func (s *NativePDFSignatureService) GetRevisions(document model.DSSDocument, pwd []byte) []PdfRevision {
	if document == nil {
		panic("DSSDocument shall be provided!")
	}

	revisions := make([]PdfRevision, 0)
	revisionDocuments := UtilsExtractRevisions(document)

	reader, err := NewNativePdfDocumentReader(document, pwd)
	if err != nil {
		panic(model.NewDSSErrorMessageCause(fmt.Sprintf(
			"The document with name [%s] is either not accessible or not PDF compatible. Reason : [%s]",
			document.Name(), err.Error()), err))
	}
	defer func() { _ = reader.Close() }()

	compositeDssDictionary := NewPdfCompositeDssDictionary()

	dssDictionary := reader.DSSDictionary()
	lastDSSDictionary := dssDictionary // the last created DSS dictionary
	compositeDssDictionary.PopulateFromDssDictionary(lastDSSDictionary)

	sigDictionaries, err := reader.ExtractSigDictionaries()
	if err != nil {
		panic(model.NewDSSErrorMessageCause("Cannot analyze signatures : "+err.Error(), err))
	}
	sigDictionaries = nativePDFSignatureServiceSortSignatureDictionaries(sigDictionaries)

	for _, sigDictEntry := range sigDictionaries {
		signatureDictionary := sigDictEntry.SignatureDictionary
		fields := sigDictEntry.Fields
		fieldNames := nativePDFSignatureServiceToStringNames(fields)
		// Upstream logs "Signature fields: {}".

		// Upstream wraps this whole per-dictionary body in try { ... } catch (Exception e) and
		// logs "Unable to parse signature {} . Reason : {}": one malformed revision (say, a
		// document time-stamp whose token does not match its signed data, whose
		// NewPdfDocTimestampRevision panics as PdfDocTimestampRevision throws a DSSException) is
		// skipped and the remaining signatures are still analysed. Whatever the body had
		// already added to revisions or lastDSSDictionary before failing is kept, as upstream
		// keeps it: the closure updates both in place.
		func() {
			defer func() {
				// Upstream logs "Unable to parse signature {} . Reason : {}".
				_ = recover()
			}()
			byteRange := signatureDictionary.ByteRange()
			cms := signatureDictionary.Contents()
			byteRangeValid := s.ValidateByteRange(byteRange, document, cms)
			byteRange.SetValid(byteRangeValid)

			var signedContent model.DSSDocument
			if byteRange.IsValid() {
				signedContent = NewPdfByteRangeDocument(document, byteRange)
				if !nativePDFSignatureServiceIsSignedContentComplete(byteRange, signedContent) {
					byteRange.SetValid(false)
				}
			}

			if !byteRange.IsValid() {
				signedContent = model.CreateEmptyDocument()
				// Upstream logs "The signature '{}' has an invalid /ByteRange! The validation will
				// result to a broken signature.".
			}

			signatureCoversWholeDocument := reader.IsSignatureCoversWholeDocument(signatureDictionary)

			revisionContent := UtilsGetRevisionContent(document, byteRange)
			if revisionReader, err := NewNativePdfDocumentReader(revisionContent, pwd); err == nil {
				// detect a modification within the signature dictionary itself (spoofing attack)
				nativePDFSignatureServiceVerifyPdfSignatureDictionary(signatureDictionary, fieldNames, revisionReader)

				// create a DSS revision if updated
				revisions, lastDSSDictionary = nativePDFSignatureServicePreviousDssDictAndUpdateIfNeeded(
					revisions, compositeDssDictionary, lastDSSDictionary, revisionReader.DSSDictionary())
				_ = revisionReader.Close()
			}
			// Upstream logs "Cannot read signature revision '{}' : {}" on failure and continues.

			previousRevision := UtilsGetPreviousRevision(byteRange, revisionDocuments)
			var newRevision PdfRevision
			if s.IsDocTimestamp(signatureDictionary) {
				newRevision = NewPdfDocTimestampRevision(signatureDictionary, fields, signedContent,
					previousRevision, signatureCoversWholeDocument)

			} else if s.IsSignature(signatureDictionary) {
				// a signature contains all the DSS dictionaries present after it
				var signatureDssDictionary PdfDssDict
				if nativePDFSignatureServiceContainsDSSRevisions(revisions) {
					signatureDssDictionary = dssDictionary
				}
				newRevision = NewPdfSignatureRevision(signatureDictionary, compositeDssDictionary,
					signatureDssDictionary, fields, signedContent, previousRevision, signatureCoversWholeDocument)

			}
			// Upstream logs "The entry {} is skipped. A signature dictionary entry with a type '{}'
			// and subFilter '{}' is not acceptable configuration!" otherwise.

			// add the signature/timestamp revision
			if newRevision != nil {
				revisions = append(revisions, newRevision)
			}

			// No nil guard on previousRevision: upstream has none either, and
			// PAdESUtilsGetPreviousRevision, like PAdESUtils#getPreviousRevision, falls back to an
			// empty document rather than returning nil.
			if revisionReader, err := NewNativePdfDocumentReader(previousRevision, pwd); err == nil {
				// check whether there is a previous update of the DSS dictionary and create a new
				// revision if needed
				revisions, lastDSSDictionary = nativePDFSignatureServicePreviousDssDictAndUpdateIfNeeded(
					revisions, compositeDssDictionary, lastDSSDictionary, revisionReader.DSSDictionary())
				_ = revisionReader.Close()
			}
		}()
	}

	return revisions
}

// nativePDFSignatureServiceSortSignatureDictionaries sorts the entries from the latest revision
// to the first. Port of the private #sortSignatureDictionaries, whose comparator is
// PdfSignatureDictionaryComparator, reversed.
func nativePDFSignatureServiceSortSignatureDictionaries(entries []PdfSignatureDictionaryFields) []PdfSignatureDictionaryFields {
	sorted := make([]PdfSignatureDictionaryFields, len(entries))
	copy(sorted, entries)
	comparator := NewPdfSignatureDictionaryComparator()
	// A stable insertion sort keeps Java's Stream.sorted() stability, which the LinkedHashMap
	// collector relies on for dictionaries that compare equal.
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && comparator.Compare(sorted[j-1].SignatureDictionary,
			sorted[j].SignatureDictionary) < 0; j-- {
			sorted[j-1], sorted[j] = sorted[j], sorted[j-1]
		}
	}
	return sorted
}

// nativePDFSignatureServiceVerifyPdfSignatureDictionary ports the private
// #verifyPdfSignatureDictionary together with #getSignatureDictionaryForFieldNames.
func nativePDFSignatureServiceVerifyPdfSignatureDictionary(signatureDictionary *PdfSignatureDictionary,
	fieldNames []string, revisionReader *NativePdfDocumentReader) {
	entries, err := revisionReader.ExtractSigDictionaries()
	if err != nil {
		return
	}
	var signatureDictionaryToCompare *PdfSignatureDictionary
	for _, entry := range entries {
		if nativePDFSignatureServiceEqualNames(fieldNames, nativePDFSignatureServiceToStringNames(entry.Fields)) {
			signatureDictionaryToCompare = entry.SignatureDictionary
			break
		}
	}
	// The only effect upstream is the warning "The signature dictionary for signature {} is not
	// consistent!"; slf4j is dropped, so the check is kept for its (side-effect free) contract.
	_ = signatureDictionary.CheckConsistency(signatureDictionaryToCompare)
}

// nativePDFSignatureServiceToStringNames ports the private #toStringNames.
func nativePDFSignatureServiceToStringNames(signatureFields []*PdfSignatureField) []string {
	names := make([]string, 0, len(signatureFields))
	for _, signatureField := range signatureFields {
		names = append(names, signatureField.FieldName())
	}
	return names
}

// nativePDFSignatureServiceEqualNames is List#equals for the field-name lists
// #getSignatureDictionaryForFieldNames compares.
func nativePDFSignatureServiceEqualNames(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

// nativePDFSignatureServicePreviousDssDictAndUpdateIfNeeded ports the private
// #getPreviousDssDictAndUpdateIfNeeded; Go returns the revisions slice too, because appending to
// it is not visible to the caller otherwise.
func nativePDFSignatureServicePreviousDssDictAndUpdateIfNeeded(revisions []PdfRevision,
	compositeDssDictionary *PdfCompositeDssDictionary, lastDSSDictionary PdfDssDict,
	currentDssDict PdfDssDict) ([]PdfRevision, PdfDssDict) {
	if lastDSSDictionary != nil && !PdfDssDictEquals(lastDSSDictionary, currentDssDict) {
		compositeDssDictionary.PopulateFromDssDictionary(lastDSSDictionary)
		revisions = append(revisions, NewPdfDocDssRevision(compositeDssDictionary, lastDSSDictionary))
	}
	return revisions, currentDssDict
}

// nativePDFSignatureServiceContainsDSSRevisions ports the private #containsDSSRevisions.
func nativePDFSignatureServiceContainsDSSRevisions(revisions []PdfRevision) bool {
	for _, revision := range revisions {
		if _, ok := revision.(*PdfDocDssRevision); ok {
			return true
		}
	}
	return false
}

// ValidateByteRange verifies the validity of the /ByteRange field against the CMS extracted from
// /Contents and the current PDF document. Port of the protected #validateByteRange.
//
// Upstream wraps its ENTIRE body - not just ByteRange#validate() - in "catch (Exception e)": a
// deliberately/incidentally malformed /ByteRange can make PAdESUtils#getSignatureValue's
// Utils.fromHex(...) throw an unchecked IllegalArgumentException (e.g. the extracted /Contents
// span is not valid hex at all - non-hex bytes, an odd length, or plain garbage a SIWA-style
// /ByteRange attack points at), and upstream still returns a graceful "false" rather than letting
// that propagate. PAdESUtilsGetSignatureValue below is a literal port of that unchecked
// exception as a panic (per PORTING.md), so this boundary needs the same recover() Java's
// catch(Exception e) provides - confirmed by pades/testdata/upstream/validation/
// pdf-byterange-overlap.pdf in the PAdES cross-validation harness, whose second signature has
// exactly such a /ByteRange and used to crash the whole parse here instead of reporting a broken
// signature.
func (s *NativePDFSignatureService) ValidateByteRange(byteRange *ByteRange, document model.DSSDocument,
	cms []byte) (valid bool) {
	defer func() {
		if r := recover(); r != nil {
			// Upstream logs "/ByteRange validation ended with error : %s. Reason : %s".
			valid = false
		}
	}()
	if err := byteRange.Validate(); err != nil {
		// Upstream logs "/ByteRange validation ended with error : %s. Reason : %s".
		return false
	}
	if !nativePDFSignatureServiceIsContentValueEqualsByteRangeExtraction(byteRange, document, cms) {
		// Upstream logs "Signature with the /ByteRange '{}' is invalid. SIWA detected!".
		return false
	}
	return true
}

// nativePDFSignatureServiceIsContentValueEqualsByteRangeExtraction checks whether the value
// incorporated into /Contents matches the range defined by the /ByteRange. It is used for SIWA
// detection. Port of the private #isContentValueEqualsByteRangeExtraction.
func nativePDFSignatureServiceIsContentValueEqualsByteRangeExtraction(byteRange *ByteRange,
	document model.DSSDocument, cms []byte) bool {
	cmsWithByteRange := UtilsGetSignatureValue(document, byteRange)
	match := bytes.Equal(cms, cmsWithByteRange)
	if !match {
		// Upstream logs "The value extracted according to /ByteRange '{}' does not match the
		// signature present in /Contents field!".
		return false
	}
	return true
}

// nativePDFSignatureServiceIsSignedContentComplete verifies whether the extracted signed content
// corresponds to the /ByteRange. Port of the private #isSignedContentComplete.
func nativePDFSignatureServiceIsSignedContentComplete(byteRange *ByteRange,
	signedContent model.DSSDocument) bool {
	expectedSignedContentLength := (byteRange.FirstPartEnd() - byteRange.FirstPartStart()) +
		byteRange.SecondPartEnd()
	signedContentLength, err := spi.DSSUtilsFileByteSize(signedContent)
	if err != nil || int64(expectedSignedContentLength) != signedContentLength {
		// Upstream logs "The length of the extracted signed content '{}' does not correspond to
		// the content length defined by the ByteRange {} : {}!".
		return false
	}
	return true
}

// IsDocTimestamp checks whether the given signature dictionary represents a DocTimeStamp.
// Both the historical TS 102 778-4 and the new EN 319 142-1 shapes are supported.
// Port of the protected #isDocTimestamp.
func (s *NativePDFSignatureService) IsDocTimestamp(pdfSigDict *PdfSignatureDictionary) bool {
	sigType := pdfSigDict.Type()
	subFilter := pdfSigDict.SubFilter()
	return (sigType == "" || PAdESConstantsTimestampType == sigType) &&
		PAdESConstantsTimestampDefaultSubFilter == subFilter
}

// IsSignature checks whether the given signature dictionary represents a signature.
// Both the historical TS 102 778-4 and the new EN 319 142-1 shapes are supported.
// Port of the protected #isSignature.
func (s *NativePDFSignatureService) IsSignature(pdfSigDict *PdfSignatureDictionary) bool {
	sigType := pdfSigDict.Type()
	subFilter := pdfSigDict.SubFilter()
	return (sigType == "" || PAdESConstantsSignatureType == sigType) &&
		PAdESConstantsTimestampDefaultSubFilter != subFilter
}

// AddDssDictionaryDefault adds the DSS dictionary to a document without password-protection and
// without a VRI dictionary. Port of addDssDictionary(DSSDocument, PdfValidationDataContainer).
func (s *NativePDFSignatureService) AddDssDictionaryDefault(document model.DSSDocument,
	validationDataForInclusion *PdfValidationDataContainer) model.DSSDocument {
	return s.AddDssDictionaryWithPassword(document, validationDataForInclusion, nil)
}

// AddDssDictionaryWithPassword adds the DSS dictionary to a password-protected document without a
// VRI dictionary. Port of addDssDictionary(DSSDocument, PdfValidationDataContainer, char[]).
func (s *NativePDFSignatureService) AddDssDictionaryWithPassword(document model.DSSDocument,
	validationDataForInclusion *PdfValidationDataContainer, pwd []byte) model.DSSDocument {
	return s.AddDssDictionary(document, validationDataForInclusion, pwd, false)
}

// AddDssDictionary adds the DSS dictionary revision to the given document.
// Port of PdfBoxSignatureService#addDssDictionary.
func (s *NativePDFSignatureService) AddDssDictionary(document model.DSSDocument,
	validationDataForInclusion *PdfValidationDataContainer, pwd []byte, includeVRIDict bool) model.DSSDocument {
	reader, err := NewNativePdfDocumentReader(document, pwd)
	if err != nil {
		panic(model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to add a new dss dictionary revision : %s", err.Error()), err))
	}
	defer func() { _ = reader.Close() }()

	updater, err := pdf.NewUpdater(reader.Document())
	if err != nil {
		panic(model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to add a new dss dictionary revision : %s", err.Error()), err))
	}

	if !validationDataForInclusion.IsEmpty() {
		if err = updater.SetDSSDictionary(s.BuildDSSDictionary(reader, validationDataForInclusion,
			includeVRIDict)); err != nil {
			panic(model.NewDSSErrorMessageCause(
				fmt.Sprintf("Unable to add a new dss dictionary revision : %s", err.Error()), err))
		}
	}
	s.EnsureESICDeveloperExtension1(updater, reader)

	// encryption is not required (no signature/timestamp is added on this step) - the Updater
	// re-encrypts anyway when the source document is encrypted (DESIGN.md Amendment 1)
	result, err := updater.Write()
	if err != nil {
		panic(model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to save a document. Reason : %s", err.Error()), err))
	}

	extendedDocument := model.NewInMemoryDocument(result.Bytes)
	extendedDocument.SetMimeType(enumerations.MimeTypeEnumPDF)
	return extendedDocument
}

// BuildDSSDictionary assembles the /DSS dictionary content from the validation data, in
// internal/pdf's DSSDictionary/TokenRef/VRIEntry shape. The dedup that upstream performs with
// COSArray#indexOf lives here, as internal/pdf never sorts and never deduplicates (DESIGN.md
// §3.5). Port of the private PdfBoxSignatureService#buildDSSDictionary together with
// #getPdfObjectForToken.
func (s *NativePDFSignatureService) BuildDSSDictionary(reader *NativePdfDocumentReader,
	validationDataForInclusion *PdfValidationDataContainer, includeVRIDict bool) pdf.DSSDictionary {
	dssDictionary := pdf.DSSDictionary{}
	knownObjects := make(map[string]pdf.TokenRef)

	signatures := validationDataForInclusion.Signatures()
	if utils.IsCollectionNotEmpty(signatures) {
		var vriEntries []pdf.VRIEntry
		for _, signature := range signatures {
			vriEntry := pdf.VRIEntry{}

			validationDataToAdd := validation.NewData()
			validationDataToAdd.AddValidationData(validationDataForInclusion.AllValidationDataForSignature(signature))

			if !validationDataToAdd.IsEmpty() {
				for _, certificateToken := range validationDataToAdd.CertificateTokens() {
					tokenRef := nativePDFSignatureServiceTokenRef(reader, validationDataForInclusion,
						knownObjects, certificateToken)
					// ensure there are no duplicated references
					if !nativePDFSignatureServiceContainsToken(vriEntry.Certs, tokenRef) {
						vriEntry.Certs = append(vriEntry.Certs, tokenRef)
						if !nativePDFSignatureServiceContainsToken(dssDictionary.Certs, tokenRef) {
							dssDictionary.Certs = append(dssDictionary.Certs, tokenRef)
						}
					}
				}
				for _, crlToken := range validationDataToAdd.CrlTokens() {
					tokenRef := nativePDFSignatureServiceTokenRef(reader, validationDataForInclusion,
						knownObjects, crlToken)
					if !nativePDFSignatureServiceContainsToken(vriEntry.CRLs, tokenRef) {
						vriEntry.CRLs = append(vriEntry.CRLs, tokenRef)
						if !nativePDFSignatureServiceContainsToken(dssDictionary.CRLs, tokenRef) {
							dssDictionary.CRLs = append(dssDictionary.CRLs, tokenRef)
						}
					}
				}
				for _, ocspToken := range validationDataToAdd.OcspTokens() {
					tokenRef := nativePDFSignatureServiceTokenRef(reader, validationDataForInclusion,
						knownObjects, ocspToken)
					if !nativePDFSignatureServiceContainsToken(vriEntry.OCSPs, tokenRef) {
						vriEntry.OCSPs = append(vriEntry.OCSPs, tokenRef)
						if !nativePDFSignatureServiceContainsToken(dssDictionary.OCSPs, tokenRef) {
							dssDictionary.OCSPs = append(dssDictionary.OCSPs, tokenRef)
						}
					}
				}

				// We can't use the CMS signed data, the pdSignature content is trimmed (000000)
				vriEntry.Name = signature.(*Signature).VRIKey()
				vriEntries = append(vriEntries, vriEntry)
			}
		}

		// optional
		if includeVRIDict {
			dssDictionary.VRI = vriEntries
		}
	}

	detachedTimestamps := validationDataForInclusion.DetachedTimestamps()
	if utils.IsCollectionNotEmpty(detachedTimestamps) { // for detached timestamps
		validationDataToAdd := validationDataForInclusion.AllValidationData()
		for _, certificateToken := range validationDataToAdd.CertificateTokens() {
			tokenRef := nativePDFSignatureServiceTokenRef(reader, validationDataForInclusion,
				knownObjects, certificateToken)
			if !nativePDFSignatureServiceContainsToken(dssDictionary.Certs, tokenRef) {
				dssDictionary.Certs = append(dssDictionary.Certs, tokenRef)
			}
		}
		for _, crlToken := range validationDataToAdd.CrlTokens() {
			tokenRef := nativePDFSignatureServiceTokenRef(reader, validationDataForInclusion,
				knownObjects, crlToken)
			if !nativePDFSignatureServiceContainsToken(dssDictionary.CRLs, tokenRef) {
				dssDictionary.CRLs = append(dssDictionary.CRLs, tokenRef)
			}
		}
		for _, ocspToken := range validationDataToAdd.OcspTokens() {
			tokenRef := nativePDFSignatureServiceTokenRef(reader, validationDataForInclusion,
				knownObjects, ocspToken)
			if !nativePDFSignatureServiceContainsToken(dssDictionary.OCSPs, tokenRef) {
				dssDictionary.OCSPs = append(dssDictionary.OCSPs, tokenRef)
			}
		}
	}

	return dssDictionary
}

// nativePDFSignatureServiceTokenRef returns the reference to use for the given token: the object
// already present in the document when the validation data container knows one, a new stream
// carrying the token's DER otherwise. Port of the private #getPdfObjectForToken.
func nativePDFSignatureServiceTokenRef(reader *NativePdfDocumentReader,
	validationDataContainer *PdfValidationDataContainer, knownObjects map[string]pdf.TokenRef,
	token model.Token) pdf.TokenRef {
	tokenKey := validationDataContainer.TokenKey(token)
	if object, known := knownObjects[tokenKey]; known {
		return object
	}

	var object pdf.TokenRef
	if objectKey := validationDataContainer.TokenReference(token); objectKey == nil {
		object = pdf.TokenRef{Data: token.Encoded()}
	} else {
		object = pdf.TokenRef{Key: pdf.ObjectKey{Num: objectKey.Number(), Gen: uint16(objectKey.Generation())}}
	}

	knownObjects[tokenKey] = object
	return object
}

// nativePDFSignatureServiceContainsToken is COSArray#indexOf(...) != -1 for the token references
// the /DSS builder accumulates.
func nativePDFSignatureServiceContainsToken(tokens []pdf.TokenRef, token pdf.TokenRef) bool {
	for _, candidate := range tokens {
		if !candidate.Key.IsZero() && candidate.Key == token.Key {
			return true
		}
		if candidate.Key.IsZero() && token.Key.IsZero() && bytes.Equal(candidate.Data, token.Data) {
			return true
		}
	}
	return false
}

// GetAvailableSignatureFieldsDefault returns the not signed signature fields.
// Port of getAvailableSignatureFields(DSSDocument).
func (s *NativePDFSignatureService) GetAvailableSignatureFieldsDefault(document model.DSSDocument) []string {
	return s.GetAvailableSignatureFields(document, nil)
}

// GetAvailableSignatureFields returns the not signed signature fields of a (possibly encrypted)
// document. Port of PdfBoxSignatureService#getAvailableSignatureFields.
func (s *NativePDFSignatureService) GetAvailableSignatureFields(document model.DSSDocument, pwd []byte) []string {
	reader, err := NewNativePdfDocumentReader(document, pwd)
	if err != nil {
		var invalidPassword *exception.InvalidPasswordException
		if errors.As(err, &invalidPassword) {
			panic(invalidPassword)
		}
		panic(model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to retrieve signature fields. Reason : %s", err.Error()), err))
	}
	defer func() { _ = reader.Close() }()

	result := make([]string, 0)
	fields, err := reader.Document().SignatureFields()
	if err != nil {
		panic(model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to retrieve signature fields. Reason : %s", err.Error()), err))
	}
	for _, field := range fields {
		if field.Value == nil {
			// PDSignatureField#getPartialName is the field's own /T; internal/pdf reports the
			// fully-qualified, dot-joined name, which coincides with it for the flat AcroForms
			// PAdES produces and is the more useful identifier for a nested one.
			result = append(result, field.Name)
		}
	}
	return result
}

// AddNewSignatureFieldDefault adds a new signature field to an existing PDF document.
// Port of addNewSignatureField(DSSDocument, SignatureFieldParameters).
func (s *NativePDFSignatureService) AddNewSignatureFieldDefault(document model.DSSDocument,
	parameters *SignatureFieldParameters) model.DSSDocument {
	return s.AddNewSignatureField(document, parameters, nil)
}

// AddNewSignatureField adds a new, empty signature field to an existing (possibly encrypted) PDF
// document. Port of PdfBoxSignatureService#addNewSignatureField.
func (s *NativePDFSignatureService) AddNewSignatureField(document model.DSSDocument,
	parameters *SignatureFieldParameters, pwd []byte) model.DSSDocument {
	reader, err := NewNativePdfDocumentReader(document, pwd)
	if err != nil {
		panic(model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to add a new signature field. Reason : %s", err.Error()), err))
	}
	defer func() { _ = reader.Close() }()

	s.CheckPdfPermissions(reader, parameters)

	if reader.NumberOfPages() < parameters.Page() {
		panic(fmt.Sprintf("The page number '%d' does not exist in the file!", parameters.Page()))
	}

	targetFieldId := parameters.FieldId()
	if utils.IsStringNotBlank(targetFieldId) {
		fields, err := reader.Document().SignatureFields()
		if err != nil {
			panic(err)
		}
		for _, field := range fields {
			if field.Name == targetFieldId {
				panic(fmt.Sprintf("The field '%s' already exists within the PDF document!", targetFieldId))
			}
		}
	}

	annotationBox := s.VisibleSignatureFieldBoxPosition(reader, parameters)

	updater, err := pdf.NewUpdater(reader.Document())
	if err != nil {
		panic(err)
	}
	rect := pdf.Rect{
		MinX: float64(annotationBox.MinX()), MinY: float64(annotationBox.MinY()),
		MaxX: float64(annotationBox.MaxX()), MaxY: float64(annotationBox.MaxY()),
	}
	if err = s.addEmptySignatureField(updater, reader, targetFieldId, parameters.Page(), rect); err != nil {
		panic(model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to add a new signature field. Reason : %s", err.Error()), err))
	}

	result, err := updater.Write()
	if err != nil {
		panic(model.NewDSSErrorMessageCause(
			fmt.Sprintf("Unable to add a new signature field. Reason : %s", err.Error()), err))
	}

	updatedDocument := model.NewInMemoryDocument(result.Bytes)
	updatedDocument.SetName("new-document.pdf")
	updatedDocument.SetMimeType(enumerations.MimeTypeEnumPDF)
	return updatedDocument
}

// addEmptySignatureField creates an unsigned signature field with the given name, box and page,
// and wires it into the AcroForm and the page's /Annots. It is the object-model half of
// PdfBoxSignatureService#addNewSignatureField, i.e. what "new PDSignatureField(acroForm)" plus
// PdfBoxUtils.createSignatureAppearanceDictionary do; internal/pdf exposes no AddSignatureField
// (its AddSignature always installs a signature dictionary), so the wiring is done here through
// the Updater's public primitives - no byte-level PDF work happens outside internal/pdf.
//
// The shape follows internal/pdf/DESIGN.md §3.3 steps 3-5: a merged field+widget dictionary,
// /SigFlags 3 on the AcroForm, and /Annots re-emitted as a direct array on the updated page.
func (s *NativePDFSignatureService) addEmptySignatureField(updater *pdf.Updater,
	reader *NativePdfDocumentReader, fieldID string, page int, rect pdf.Rect) error {
	document := reader.Document()

	pageDict, pageKey, err := document.Page(page)
	if err != nil {
		return err
	}
	if fieldID == "" {
		fieldID, err = nativePDFSignatureServiceGenerateFieldName(document)
		if err != nil {
			return err
		}
	}

	appearanceKey := updater.Add(nativePDFSignatureServiceAppearance(rect))
	fieldKey := updater.Alloc()
	field := pdf.DictOf(
		pdf.Name("FT"), pdf.Name("Sig"),
		pdf.Name("Type"), pdf.Name("Annot"),
		pdf.Name("Subtype"), pdf.Name("Widget"),
		pdf.Name("F"), pdf.Integer(4), // Print
		pdf.Name("T"), pdf.String{Bytes: []byte(fieldID)},
		pdf.Name("Rect"), rect.Array(),
		pdf.Name("P"), pdf.Ref{Num: pageKey.Num, Gen: pageKey.Gen},
		pdf.Name("AP"), pdf.DictOf(pdf.Name("N"),
			pdf.Ref{Num: appearanceKey.Num, Gen: appearanceKey.Gen}),
	)
	updater.Put(fieldKey, field)
	fieldRef := pdf.Ref{Num: fieldKey.Num, Gen: fieldKey.Gen}

	// --- AcroForm: create when absent, append to /Fields, set /SigFlags 3
	catalog := updater.Catalog()
	var acroForm *pdf.Dict
	if acroFormKey, ok := document.RefAt(catalog, pdf.Name(PAdESConstantsAcroFormName)); ok {
		updated, err := updater.Update(acroFormKey)
		if err != nil {
			return err
		}
		acroForm, ok = updated.(*pdf.Dict)
		if !ok {
			return fmt.Errorf("pades: /AcroForm is not a dictionary")
		}
	} else if existing, ok := document.GetDict(catalog, pdf.Name(PAdESConstantsAcroFormName)); ok {
		acroForm = existing.Clone()
		catalog.Set(pdf.Name(PAdESConstantsAcroFormName), acroForm)
	} else {
		acroForm = pdf.NewDict()
		catalog.Set(pdf.Name(PAdESConstantsAcroFormName), acroForm)
	}
	fields, _ := document.GetArray(acroForm, pdf.Name(PAdESConstantsFieldsName))
	acroForm.Set(pdf.Name(PAdESConstantsFieldsName), append(append(pdf.Array{}, fields...), fieldRef))
	acroForm.Set(pdf.Name(PAdESConstantsSigFlagsName), pdf.Integer(3)) // SignaturesExist|AppendOnly

	// --- Page: append the widget to /Annots, re-emitted direct
	updatedPage, err := updater.Update(pageKey)
	if err != nil {
		return err
	}
	updatedPageDict, ok := updatedPage.(*pdf.Dict)
	if !ok {
		return fmt.Errorf("pades: page %d is not a dictionary", page)
	}
	annots, _ := document.GetArray(pageDict, pdf.Name(PAdESConstantsAnnotsName))
	updatedPageDict.Set(pdf.Name(PAdESConstantsAnnotsName), append(append(pdf.Array{}, annots...), fieldRef))
	return nil
}

// nativePDFSignatureServiceGenerateFieldName builds the "Signature<N>" partial name PDFBox's
// PDSignatureField generates for an unnamed field, skipping the names already taken.
func nativePDFSignatureServiceGenerateFieldName(document *pdf.Document) (string, error) {
	fields, err := document.SignatureFields()
	if err != nil {
		return "", err
	}
	taken := make(map[string]bool, len(fields))
	for _, field := range fields {
		taken[field.Name] = true
	}
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("Signature%d", i)
		if !taken[candidate] {
			return candidate, nil
		}
	}
}

// VisibleSignatureFieldBoxPositionForDrawer checks the validity of the signature field position
// and returns the calculated signature field box, nil when the drawer cannot build one.
// Port of the protected
// getVisibleSignatureFieldBoxPosition(SignatureDrawer, PdfDocumentReader, SignatureFieldParameters).
func (s *NativePDFSignatureService) VisibleSignatureFieldBoxPositionForDrawer(signatureDrawer SignatureDrawer,
	documentReader *NativePdfDocumentReader, fieldParameters *SignatureFieldParameters) *AnnotationBox {
	signatureFieldAnnotation := s.BuildSignatureFieldBox(signatureDrawer)
	if signatureFieldAnnotation != nil {
		s.AssertSignatureFieldPositionValid(documentReader, *signatureFieldAnnotation, fieldParameters.Page())
	}
	return signatureFieldAnnotation
}

// BuildSignatureFieldBox returns the signature field box of the given drawer, nil when the drawer
// cannot build one. Port of the protected #buildSignatureFieldBox.
func (s *NativePDFSignatureService) BuildSignatureFieldBox(signatureDrawer SignatureDrawer) *AnnotationBox {
	if signatureFieldBoxBuilder, ok := signatureDrawer.(SignatureFieldBoxBuilder); ok {
		signatureFieldBox, err := signatureFieldBoxBuilder.BuildSignatureFieldBox()
		if err == nil && signatureFieldBox != nil {
			annotationBox := signatureFieldBox.AnnotationBox()
			return &annotationBox
		}
	}
	// Upstream logs "The used SignatureDrawer shall be an instance of
	// VisibleSignatureFieldBoxBuilder in order to verify a SignatureField position!".
	return nil
}

// VisibleSignatureFieldBoxPosition computes the signature field box of a new field, rotating it
// relatively to the page and checking that it does not overlap with existing annotations.
// Port of the protected
// getVisibleSignatureFieldBoxPosition(PdfDocumentReader, SignatureFieldParameters).
func (s *NativePDFSignatureService) VisibleSignatureFieldBoxPosition(reader *NativePdfDocumentReader,
	parameters *SignatureFieldParameters) AnnotationBox {
	pageRotation := reader.PageRotation(parameters.Page())
	globalRotation := ImageRotationUtilsRotation(parameters.Rotation(), pageRotation)

	originalPageBox := reader.PageBox(parameters.Page())
	pageBox := originalPageBox
	annotationBox := NewAnnotationBoxFromFieldParameters(parameters)
	if ImageRotationUtilsIsSwapOfDimensionsRequired(globalRotation) {
		pageBox = ImageRotationUtilsSwapDimensions(pageBox)
	}

	annotationBox = ImageRotationUtilsRotateRelativelyWrappingBox(annotationBox, pageBox, 360-globalRotation)

	annotationBox = s.ToPdfPageCoordinates(annotationBox, originalPageBox)

	s.AssertSignatureFieldPositionValid(reader, annotationBox, parameters.Page())
	return annotationBox
}

// AssertSignatureFieldPositionValid verifies the validity of the signature field box
// configuration. Port of the protected #assertSignatureFieldPositionValid.
func (s *NativePDFSignatureService) AssertSignatureFieldPositionValid(documentReader PdfDocumentReader,
	annotationBox AnnotationBox, pageNumber int) {
	s.PdfSignatureFieldPositionChecker.AssertSignatureFieldPositionValid(documentReader, annotationBox, pageNumber)
}

// ToPdfPageCoordinates transforms the field box's position and dimensions to the given page.
// Port of the protected #toPdfPageCoordinates.
func (s *NativePDFSignatureService) ToPdfPageCoordinates(fieldAnnotationBox, pageBox AnnotationBox) AnnotationBox {
	return fieldAnnotationBox.ToPdfPageCoordinates(pageBox)
}

// AnalyzePdfModifications analyzes the PDF revisions of the given signatures and their document
// timestamps for modifications (shadow attacks). Port of #analyzePdfModifications.
func (s *NativePDFSignatureService) AnalyzePdfModifications(document model.DSSDocument,
	signatures []validation.AdvancedSignature, pwd []byte) {
	if utils.IsCollectionEmpty(signatures) {
		return
	}

	finalRevisionReader, err := NewNativePdfDocumentReader(document, pwd)
	if err != nil {
		// Upstream logs "Unable to proceed PDF modification detection. Reason : {}".
		return
	}
	defer func() { _ = finalRevisionReader.Close() }()

	for _, signature := range signatures {
		padesSignature := signature.(*Signature)
		s.AnalyzeRevisionModifications(document, padesSignature.PdfRevision(), finalRevisionReader, pwd)
	}
	for _, timestampToken := range nativePDFSignatureServiceUniqueTimestamps(signatures) {
		if pdfTimestampToken, ok := PdfTimestampTokenOf(timestampToken); ok {
			s.AnalyzeRevisionModifications(document, pdfTimestampToken.PdfRevision(), finalRevisionReader, pwd)
		}
	}
}

// nativePDFSignatureServiceUniqueTimestamps ports the private #getUniqueTimestamps.
func nativePDFSignatureServiceUniqueTimestamps(signatures []validation.AdvancedSignature) []*validation.TimestampToken {
	var timestampTokens []*validation.TimestampToken
	for _, signature := range signatures {
		timestampTokens = append(timestampTokens, signature.DocumentTimestamps()...)
	}
	return timestampTokens
}

// AnalyzeTimestampPdfModifications analyzes the PDF revisions of the given detached document
// timestamps for modifications. Port of #analyzeTimestampPdfModifications.
func (s *NativePDFSignatureService) AnalyzeTimestampPdfModifications(document model.DSSDocument,
	timestamps []*validation.TimestampToken, pwd []byte) {
	if utils.IsCollectionEmpty(timestamps) {
		return
	}

	finalRevisionReader, err := NewNativePdfDocumentReader(document, pwd)
	if err != nil {
		// Upstream logs "Unable to proceed PDF modification detection. Reason : {}".
		return
	}
	defer func() { _ = finalRevisionReader.Close() }()

	for _, timestampToken := range timestamps {
		// Java's "timestampToken instanceof PdfTimestampToken"; a Go *validation.TimestampToken
		// cannot be asserted back to the embedding PdfTimestampToken, so the recovery goes
		// through PdfTimestampTokenOf (see validation/timestamp).
		if pdfTimestampToken, ok := PdfTimestampTokenOf(timestampToken); ok {
			s.AnalyzeRevisionModifications(document, pdfTimestampToken.PdfRevision(), finalRevisionReader, pwd)
		}
	}
}

// AnalyzeRevisionModifications performs the modification analysis for a single revision.
// Port of the protected
// analyzePdfModifications(DSSDocument, PdfCMSRevision, PdfDocumentReader, char[]).
func (s *NativePDFSignatureService) AnalyzeRevisionModifications(document model.DSSDocument,
	pdfRevision PdfCMSRevision, finalRevisionReader *NativePdfDocumentReader, pwd []byte) {
	revisionContent := UtilsGetRevisionContent(document, pdfRevision.ByteRange())
	modificationDetection := s.modificationDetection(finalRevisionReader, revisionContent, pwd)
	if modificationDetection != nil {
		pdfRevision.SetModificationDetection(modificationDetection)
	}
}

// modificationDetection ports the private #getModificationDetection.
func (s *NativePDFSignatureService) modificationDetection(finalRevisionReader *NativePdfDocumentReader,
	originalDocument model.DSSDocument, pwd []byte) *PdfModificationDetection {
	signedRevisionReader, err := NewNativePdfDocumentReader(originalDocument, pwd)
	if err != nil {
		return nil
	}
	defer func() { _ = signedRevisionReader.Close() }()

	pdfModificationDetection := NewPdfModificationDetection()
	pdfModificationDetection.SetAnnotationOverlaps(s.PdfDifferencesFinder.AnnotationOverlaps(finalRevisionReader))
	pdfModificationDetection.SetPageDifferences(s.PdfDifferencesFinder.PagesDifferences(signedRevisionReader,
		finalRevisionReader))
	pdfModificationDetection.SetVisualDifferences(s.VisualDifferences(signedRevisionReader, finalRevisionReader))
	pdfModificationDetection.SetObjectModifications(s.PdfObjectModificationsFinder.Find(signedRevisionReader,
		finalRevisionReader))
	return pdfModificationDetection
}

// VisualDifferences returns the visual differences between the signed revision and the final
// document. Port of the protected #getVisualDifferences.
//
// NOTE: the native engine has no rasteriser, so DefaultPdfDifferencesFinder yields no
// VISUAL_DIFFERENCE modification (internal/pdf/DESIGN.md §0.2). The geometric checks -
// annotation overlaps and page differences - are unaffected.
func (s *NativePDFSignatureService) VisualDifferences(signedRevisionReader,
	finalRevisionReader PdfDocumentReader) []PdfModification {
	return s.PdfDifferencesFinder.VisualDifferences(signedRevisionReader, finalRevisionReader)
}

// CheckPdfPermissions verifies the PDF permissions dictionaries.
// Port of the protected #checkPdfPermissions.
func (s *NativePDFSignatureService) CheckPdfPermissions(documentReader PdfDocumentReader,
	fieldParameters *SignatureFieldParameters) {
	s.PdfPermissionsChecker.CheckDocumentPermissions(documentReader, fieldParameters)
	if !s.IsDocumentTimestampLayer() {
		s.PdfPermissionsChecker.CheckSignatureRestrictionDictionaries(documentReader, fieldParameters)
	}
}

// AssertContentSizeSufficient verifies whether the assigned /Contents size is sufficient to
// encapsulate the CMS signed data. Port of the protected #assertContentSizeSufficient.
func (s *NativePDFSignatureService) AssertContentSizeSufficient(cmsSignedData []byte,
	parameters CommonParameters) {
	csize := parameters.ContentSize()
	if csize < len(cmsSignedData) {
		panic(fmt.Sprintf("Unable to save a document. Reason : The signature size [%d] is too small "+
			"for the signature value with a length [%d]. Use setContentSize(...) method "+
			"to define a bigger length.", csize, len(cmsSignedData)))
	}
}

// PreviewPageWithVisualSignature is not supported; see ErrRasterisationNotSupported.
// Port of PdfBoxSignatureService#previewPageWithVisualSignature.
func (s *NativePDFSignatureService) PreviewPageWithVisualSignature(toSignDocument model.DSSDocument,
	parameters CommonParameters) model.DSSDocument {
	panic(ErrRasterisationNotSupported)
}

// PreviewSignatureField is not supported; see ErrRasterisationNotSupported.
// Port of PdfBoxSignatureService#previewSignatureField.
func (s *NativePDFSignatureService) PreviewSignatureField(toSignDocument model.DSSDocument,
	parameters CommonParameters) model.DSSDocument {
	panic(ErrRasterisationNotSupported)
}

// Compile-time assertion standing in for Java's "extends AbstractPDFSignatureService".
var _ PDFSignatureService = (*NativePDFSignatureService)(nil)
