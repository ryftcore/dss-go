// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/evidencerecord/AbstractASiCContainerEvidenceRecordBuilder.java (DSS 6.5.RC1).
package asic

import (
	"fmt"
	"slices"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	"github.com/ryftcore/dss-go/dss/utils"
)

// AbstractASiCContainerEvidenceRecordBuilderOverrides declares the operation
// AbstractASiCContainerEvidenceRecordBuilder calls back into virtually from Build(). Every
// concrete builder must call InitAbstractASiCContainerEvidenceRecordBuilder with itself before
// use, or the base's virtual calls will not reach the override.
type AbstractASiCContainerEvidenceRecordBuilderOverrides interface {
	// GetASiCContentBuilder gets an instance of AbstractASiCContentBuilder. Port of the
	// protected abstract getASiCContentBuilder().
	GetASiCContentBuilder() *AbstractASiCContentBuilder

	// AssertEvidenceRecordFilenameValid verifies validity of the evidence record filename to
	// the ASiC container convention. Java declares this protected (not abstract) with a default
	// body, and CAdES/XAdES subclasses override it to layer format-specific checks on top via
	// super.assertEvidenceRecordFilenameValid(...); ported through the overrides interface (not
	// a plain self-call) so Build() reaches the leaf override. Port of the protected
	// assertEvidenceRecordFilenameValid(String, EvidenceRecordTypeEnum, Content).
	AssertEvidenceRecordFilenameValid(evidenceRecordFilename string, evidenceRecordType enumerations.EvidenceRecordTypeEnum, asicContent *Content)
}

// AbstractASiCContainerEvidenceRecordBuilder incorporates an existing evidence record document
// within an ASiC container.
type AbstractASiCContainerEvidenceRecordBuilder struct {
	// overrides points back at the concrete builder; see
	// InitAbstractASiCContainerEvidenceRecordBuilder.
	overrides AbstractASiCContainerEvidenceRecordBuilderOverrides

	// CertificateVerifier is used to verify the evidence record against the provided data.
	// Java declares the field protected final.
	CertificateVerifier validation.CertificateVerifier

	// AsicFilenameFactory is the filename factory. Java declares the field protected final.
	AsicFilenameFactory EvidenceRecordFilenameFactory
}

// NewAbstractASiCContainerEvidenceRecordBuilderBase is the default constructor. Port of the
// protected constructor(CertificateVerifier, EvidenceRecordFilenameFactory). The subclass
// constructor must follow it with InitAbstractASiCContainerEvidenceRecordBuilder.
func NewAbstractASiCContainerEvidenceRecordBuilderBase(certificateVerifier validation.CertificateVerifier, asicFilenameFactory EvidenceRecordFilenameFactory) AbstractASiCContainerEvidenceRecordBuilder {
	return AbstractASiCContainerEvidenceRecordBuilder{
		CertificateVerifier: certificateVerifier,
		AsicFilenameFactory: asicFilenameFactory,
	}
}

// InitAbstractASiCContainerEvidenceRecordBuilder registers the concrete builder with its base
// so the base can dispatch GetASiCContentBuilder. Every concrete builder constructor must call
// this once.
func (b *AbstractASiCContainerEvidenceRecordBuilder) InitAbstractASiCContainerEvidenceRecordBuilder(overrides AbstractASiCContainerEvidenceRecordBuilderOverrides) {
	b.overrides = overrides
}

func (b *AbstractASiCContainerEvidenceRecordBuilder) requireOverrides() AbstractASiCContainerEvidenceRecordBuilderOverrides {
	if b.overrides == nil {
		panic("AbstractASiCContainerEvidenceRecordBuilder was not initialised: the concrete builder must call InitAbstractASiCContainerEvidenceRecordBuilder in its constructor")
	}
	return b.overrides
}

// Build builds an ASiCContent containing the evidence record file document. Ports
// build(List, DSSDocument, ContainerEvidenceRecordParameters).
//
// Cross-chunk assumption (ZIPCORE): UtilsEnsureMimeTypeAndZipComment(*Content,
// *ASiCContainerEvidenceRecordParameters) *ASiCContent mirrors
// ASiCUtils.ensureMimeTypeAndZipComment(Content, Parameters).
func (b *AbstractASiCContainerEvidenceRecordBuilder) Build(documents []model.DSSDocument, evidenceRecordDocument model.DSSDocument, parameters *ContainerEvidenceRecordParameters) (*Content, error) {
	asicContent := b.initASiCContent(documents, parameters)
	b.assertASiCContentValid(asicContent, parameters)

	evidenceRecordManifest := b.getASiCEvidenceRecordManifest(parameters)
	manifestFile := b.parseManifestFile(evidenceRecordManifest, asicContent) // may be nil
	b.assertManifestFileValid(manifestFile, asicContent)

	evidenceRecord := b.getEvidenceRecord(evidenceRecordDocument, manifestFile, asicContent)
	b.assertEvidenceRecordValid(evidenceRecord, manifestFile)

	coveredDocuments := b.getDocumentsCoveredByEvidenceRecord(evidenceRecord, asicContent)
	b.assertSignedDataCovered(asicContent, spi.DSSUtilsDocumentNames(coveredDocuments))

	evidenceRecordFilename := b.getEvidenceRecordFilename(evidenceRecord, manifestFile, asicContent)
	b.requireOverrides().AssertEvidenceRecordFilenameValid(evidenceRecordFilename, evidenceRecord.EvidenceRecordType(), asicContent)

	evidenceRecordDocument.SetName(evidenceRecordFilename)
	asicContent.SetEvidenceRecordDocuments(append(asicContent.EvidenceRecordDocuments(), evidenceRecordDocument))

	if evidenceRecordManifest == nil {
		evidenceRecordManifest = b.buildEvidenceRecordManifest(
			asicContent, coveredDocuments, evidenceRecord.OriginalDigestAlgorithm(), evidenceRecordFilename)
	}
	// NOTE: can be nil (e.g. for ASiC-S)
	if evidenceRecordManifest != nil {
		asicContent.SetEvidenceRecordManifestDocuments(append(asicContent.EvidenceRecordManifestDocuments(), evidenceRecordManifest))
	}

	return UtilsEnsureMimeTypeAndZipComment(asicContent, &parameters.Parameters)
}

// initASiCContent initializes an ASiCContent from the given list of documents. Ports the
// protected initASiCContent(List, Parameters).
func (b *AbstractASiCContainerEvidenceRecordBuilder) initASiCContent(documents []model.DSSDocument, parameters *ContainerEvidenceRecordParameters) *Content {
	return b.requireOverrides().GetASiCContentBuilder().Build(documents, parameters.ContainerType())
}

// getASiCEvidenceRecordManifest gets the provided ASiCEvidenceRecordManifest file. Ports the
// protected getASiCEvidenceRecordManifest(ContainerEvidenceRecordParameters). Logging
// (LOG.info) is dropped per PORTING.md.
func (b *AbstractASiCContainerEvidenceRecordBuilder) getASiCEvidenceRecordManifest(parameters *ContainerEvidenceRecordParameters) model.DSSDocument {
	if parameters.AsicEvidenceRecordManifest() != nil {
		if enumerations.ASiCContainerTypeASiCE == parameters.ContainerType() {
			return parameters.AsicEvidenceRecordManifest()
		}
	}
	return nil
}

// getEvidenceRecord creates an EvidenceRecord from a provided evidenceRecordDocument. Ports
// the protected getEvidenceRecord(DSSDocument, ManifestFile, Content).
func (b *AbstractASiCContainerEvidenceRecordBuilder) getEvidenceRecord(evidenceRecordDocument model.DSSDocument, manifestFile *model.ManifestFile, asicContent *Content) validation.EvidenceRecord {
	evidenceRecordAnalyzer, err := analyzer.EvidenceRecordAnalyzerFromDocument(evidenceRecordDocument)
	if err == nil {
		evidenceRecordAnalyzer.SetManifestFile(manifestFile)
		evidenceRecordAnalyzer.SetDetachedContents(asicContent.AllDocuments())
		evidenceRecordAnalyzer.SetEvidenceRecordOrigin(enumerations.EvidenceRecordOriginContainer)
		return evidenceRecordAnalyzer.EvidenceRecord()
	}
	panic(exception.NewIllegalInputException(fmt.Sprintf(
		"Unable to build evidence record document. Reason : %s", err.Error())))
}

// parseManifestFile attempts to parse an evidenceRecordManifest document as an
// ASiCEvidenceRecordManifest file. Ports the protected parseManifestFile(DSSDocument,
// Content).
func (b *AbstractASiCContainerEvidenceRecordBuilder) parseManifestFile(evidenceRecordManifest model.DSSDocument, asicContent *Content) *model.ManifestFile {
	if evidenceRecordManifest == nil {
		return nil
	}

	if evidenceRecordManifest.Name() != "" {
		b.assertASiCEvidenceRecordManifestValid(evidenceRecordManifest.Name(), asicContent)
	} else {
		evidenceRecordManifest.SetName(b.AsicFilenameFactory.EvidenceRecordManifestFilename(asicContent))
	}

	manifestFile := ManifestParserGetManifestFile(evidenceRecordManifest)
	if manifestFile == nil {
		panic(exception.NewIllegalInputException("Unable to parse the provided ASiCEvidenceRecordManifest document! More detail in logs."))
	}
	manifestFile.SetManifestType(enumerations.ASiCManifestTypeEnumEvidenceRecord)
	return manifestFile
}

// assertASiCEvidenceRecordManifestValid verifies whether the ASiCEvidenceRecordManifest
// filename is valid. Ports the protected assertASiCEvidenceRecordManifestValid(String,
// Content).
func (b *AbstractASiCContainerEvidenceRecordBuilder) assertASiCEvidenceRecordManifestValid(manifestFilename string, asicContent *Content) {
	asicDocumentNames := spi.DSSUtilsDocumentNames(asicContent.AllDocuments())
	if slices.Contains(asicDocumentNames, manifestFilename) {
		panic(exception.NewIllegalInputException(fmt.Sprintf("The manifest filename '%s' is already present "+
			"within the ASiC container!", manifestFilename)))
	}
	if !UtilsIsEvidenceRecordManifest(manifestFilename) {
		panic(fmt.Sprintf("The manifest filename '%s' is not compliant "+
			"to the ASiCEvidenceRecordManifest filename convention!", manifestFilename))
	}
}

// assertManifestFileValid verifies the validity of the ASiCEvidenceRecordManifest file. Ports
// the protected assertManifestFileValid(ManifestFile, Content).
func (b *AbstractASiCContainerEvidenceRecordBuilder) assertManifestFileValid(manifestFile *model.ManifestFile, asicContent *Content) {
	if manifestFile == nil {
		return
	}

	manifestValidator := NewASiCManifestValidator(manifestFile, asicContent.AllDocuments())
	manifestValidator.ValidateEntries()

	for _, manifestEntry := range manifestFile.Entries() {
		if !manifestEntry.IsFound() || !manifestEntry.IsIntact() {
			panic(exception.NewIllegalInputException(fmt.Sprintf("The manifest entry '%s' was not found or digest does not intact! "+
				"Please provide a valid ASiCEvidenceRecordManifest document.", manifestEntry.Uri())))
		}
	}

	manifestCoveredFilenames := make([]string, 0, len(manifestFile.Entries()))
	for _, entry := range manifestFile.Entries() {
		if entry.Document() != nil {
			manifestCoveredFilenames = append(manifestCoveredFilenames, entry.Document().Name())
		}
	}
	b.assertSignedDataCovered(asicContent, manifestCoveredFilenames)
}

// getDocumentsCoveredByEvidenceRecord ports the private getDocumentsCoveredByEvidenceRecord(
// EvidenceRecord, Content).
func (b *AbstractASiCContainerEvidenceRecordBuilder) getDocumentsCoveredByEvidenceRecord(evidenceRecord validation.EvidenceRecord, asicContent *Content) []model.DSSDocument {
	coveredDocuments := make([]model.DSSDocument, 0)
	allDocuments := asicContent.AllDocuments()
	allDocumentFilenames := spi.DSSUtilsDocumentNames(allDocuments)
	for _, referenceValidation := range evidenceRecord.ReferenceValidation() {
		if referenceValidation.Document() != nil && slices.Contains(allDocumentFilenames, referenceValidation.Document().Name()) {
			coveredDocuments = append(coveredDocuments, spi.DSSUtilsDocumentWithName(allDocuments, referenceValidation.Document().Name()))
		}
	}
	return coveredDocuments
}

// getEvidenceRecordFilename gets the filename for the evidence record to be incorporated.
// Ports the protected getEvidenceRecordFilename(EvidenceRecord, ManifestFile, ASiCContent).
func (b *AbstractASiCContainerEvidenceRecordBuilder) getEvidenceRecordFilename(evidenceRecord validation.EvidenceRecord, manifestFile *model.ManifestFile, asicContent *Content) string {
	if manifestFile != nil {
		return manifestFile.SignatureFilename()
	}
	return b.AsicFilenameFactory.EvidenceRecordFilename(asicContent, evidenceRecord.EvidenceRecordType())
}

// buildEvidenceRecordManifest builds an ASiCEvidenceRecordManifest for the evidence record
// based on a list of coveredDocuments when required. Ports the protected
// buildEvidenceRecordManifest(Content, List, DigestAlgorithm, String).
func (b *AbstractASiCContainerEvidenceRecordBuilder) buildEvidenceRecordManifest(asicContent *Content, coveredDocuments []model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm, evidenceRecordFilename string) model.DSSDocument {
	if enumerations.ASiCContainerTypeASiCE == asicContent.ContainerType() {
		names := spi.DSSUtilsDocumentNames(coveredDocuments)
		manifestDocument, err := NewASiCEvidenceRecordManifestBuilder(asicContent, digestAlgorithm, evidenceRecordFilename).
			SetAsicContentDocumentFilter(AllowedFilenamesFilter(names...)).
			SetEvidenceRecordFilenameFactory(b.AsicFilenameFactory).
			Build()
		if err != nil {
			panic(err)
		}
		return manifestDocument
	}
	// skip for ASiC-S
	return nil
}

// assertASiCContentValid verifies whether the provided Content is valid and can be
// successfully protected by a new evidence record. Ports the protected
// assertASiCContentValid(Content, Parameters).
//
// Cross-chunk assumption (ZIPCORE): ASiCUtilsIsASiCE(*ASiCParameters) bool mirrors
// ASiCUtils.isASiCE(Parameters).
func (b *AbstractASiCContainerEvidenceRecordBuilder) assertASiCContentValid(asicContent *Content, parameters *ContainerEvidenceRecordParameters) {
	currentContainerType := asicContent.ContainerType()

	asice := UtilsIsASiCE(&parameters.Parameters)
	switch {
	case asice && enumerations.ASiCContainerTypeASiCE == currentContainerType:
		// ok

	case !asice && enumerations.ASiCContainerTypeASiCS == currentContainerType:
		if utils.CollectionSize(asicContent.SignedDocuments()) != 1 {
			panic("Only one original document is expected for the ASiC-S container type! If required, " +
				"please create a 'package.zip' and provide it directly as a parameter. " +
				"Otherwise, please switch to the ASiC-E type.")
		}
		if utils.IsCollectionNotEmpty(asicContent.SignatureDocuments()) ||
			utils.IsCollectionNotEmpty(asicContent.TimestampDocuments()) ||
			utils.IsCollectionNotEmpty(asicContent.EvidenceRecordDocuments()) {
			panic(exception.NewIllegalInputException(
				"Only one of the signature, timestamp or evidence record document types is allowed " +
					"within an ASiC-S container type!"))
		}

	default:
		panic(fmt.Sprintf("Original container type '%s' vs parameter : '%s'", currentContainerType,
			parameters.ContainerType()))
	}
}

// assertSignedDataCovered verifies whether the original or signed documents are successfully
// covered by the evidence record. Ports the protected assertSignedDataCovered(ASiCContent,
// List).
func (b *AbstractASiCContainerEvidenceRecordBuilder) assertSignedDataCovered(asicContent *Content, coveredDocumentFilenames []string) {
	signedDocumentNames := spi.DSSUtilsDocumentNames(asicContent.SignedDocuments())
	for _, signedDocumentFilename := range signedDocumentNames {
		if !slices.Contains(coveredDocumentFilenames, signedDocumentFilename) {
			panic(exception.NewIllegalInputException(fmt.Sprintf("The original document with name '%s' is not covered "+
				"by the evidence record!", signedDocumentFilename)))
		}
	}

	for _, documentName := range coveredDocumentFilenames {
		linkedManifest := ManifestParserGetLinkedManifest(asicContent.AllManifestDocuments(), documentName)
		b.assertManifestSignedDataCoveredRecursively(linkedManifest, coveredDocumentFilenames, asicContent)
	}
}

// assertManifestSignedDataCoveredRecursively ports the private
// assertManifestSignedDataCoveredRecursively(DSSDocument, List, Content).
func (b *AbstractASiCContainerEvidenceRecordBuilder) assertManifestSignedDataCoveredRecursively(manifestDocument model.DSSDocument, coveredDocumentNames []string, asicContent *Content) {
	if manifestDocument == nil {
		return
	}
	if !slices.Contains(coveredDocumentNames, manifestDocument.Name()) {
		panic(exception.NewIllegalInputException(fmt.Sprintf("Digest of a signed ASiC Manifest with name '%s' "+
			"has not been found in the evidence record's covered objects!", manifestDocument.Name())))
	}
	manifestFile := ManifestParserGetManifestFile(manifestDocument)
	if manifestFile != nil {
		for _, entry := range manifestFile.Entries() {
			if !slices.Contains(coveredDocumentNames, entry.Uri()) {
				panic(exception.NewIllegalInputException(fmt.Sprintf("Digest for a document referenced from "+
					"a covered ASiC Manifest with name '%s' has not been found in the evidence record's covered objects!",
					entry.Uri())))
			}
			linkedManifest := ManifestParserGetLinkedManifest(asicContent.AllManifestDocuments(), entry.Uri())
			b.assertManifestSignedDataCoveredRecursively(linkedManifest, coveredDocumentNames, asicContent)
		}
	}
}

// assertEvidenceRecordValid verifies whether the provided EvidenceRecord covers the original
// data files. Ports the protected assertEvidenceRecordValid(EvidenceRecord, ManifestFile).
func (b *AbstractASiCContainerEvidenceRecordBuilder) assertEvidenceRecordValid(evidenceRecord validation.EvidenceRecord, manifestFile *model.ManifestFile) {
	if manifestFile != nil {
		for _, manifestEntry := range manifestFile.Entries() {
			digest := manifestEntry.Digest()
			if !digest.IsEmpty() && evidenceRecord.OriginalDigestAlgorithm() != digest.Algorithm() {
				panic(exception.NewIllegalInputException(fmt.Sprintf("Digest algorithm '%s' found in the ASiCEvidenceRecordManifest document "+
					"does not correspond to the Digest Algorithm '%s' used for the first data object group of evidence record generation!",
					digest.Algorithm(), evidenceRecord.OriginalDigestAlgorithm())))
			}
		}
	}

	errorMessage := "The digest covered by the evidence record do not correspond to " +
		"the digest computed on the provided content!"
	signedDataFound := false
	for _, referenceValidation := range evidenceRecord.ReferenceValidation() {
		if enumerations.DigestMatcherTypeEvidenceRecordOrphanReference != referenceValidation.Type() {
			if !referenceValidation.IsIntact() {
				if referenceValidation.Document() != nil {
					panic(exception.NewIllegalInputException(fmt.Sprintf("The digest of document '%s' has not been found "+
						"within the manifest file or/and evidence record!", referenceValidation.Document().Name())))
				}
				panic(exception.NewIllegalInputException(errorMessage))
			}
			signedDataFound = true
		}
	}
	if !signedDataFound {
		panic(exception.NewIllegalInputException(errorMessage))
	}
	b.validateTimestamps(evidenceRecord)
}

// AssertEvidenceRecordFilenameValid verifies validity of the evidence record filename to the
// ASiC container convention. This is the default body Java's protected (non-abstract) method
// provides; CAdES/XAdES leaf builders call this via
// b.AbstractASiCContainerEvidenceRecordBuilder.AssertEvidenceRecordFilenameValid(...) before
// layering their own checks, matching Java's super.assertEvidenceRecordFilenameValid(...). Ports
// the protected assertEvidenceRecordFilenameValid(String, EvidenceRecordTypeEnum, Content).
func (b *AbstractASiCContainerEvidenceRecordBuilder) AssertEvidenceRecordFilenameValid(evidenceRecordFilename string, evidenceRecordType enumerations.EvidenceRecordTypeEnum, asicContent *Content) {
	asicDocumentNames := spi.DSSUtilsDocumentNames(asicContent.AllDocuments())
	if slices.Contains(asicDocumentNames, evidenceRecordFilename) {
		panic(exception.NewIllegalInputException(fmt.Sprintf("The evidence record filename '%s' is already present "+
			"within the ASiC container!", evidenceRecordFilename)))
	}
}

// validateTimestamps ports the private validateTimestamps(EvidenceRecord).
func (b *AbstractASiCContainerEvidenceRecordBuilder) validateTimestamps(evidenceRecord validation.EvidenceRecord) {
	validationContext := validation.NewSignatureValidationContext()
	validationContext.Initialize(b.CertificateVerifier)

	validationContext.AddDocumentCertificateSource(evidenceRecord.CertificateSource())
	for _, timestampToken := range evidenceRecord.Timestamps() {
		validationContext.AddTimestampTokenForVerification(timestampToken)
	}

	validationContext.Validate()

	signatureValidationAlerter := validation.NewSignatureValidationAlerter(validationContext)
	signatureValidationAlerter.SetSigningOperation(enumerations.SigningOperationAddEvidenceRecord)
	signatureValidationAlerter.AssertAllTimestampsValid()
}
