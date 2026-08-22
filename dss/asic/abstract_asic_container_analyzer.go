// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/validation/AbstractASiCContainerAnalyzer.java (DSS 6.5.RC1).
package asic

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	analyzertimestamp "github.com/ryftcore/dss-go/dss/spi/validation/analyzer/timestamp"
	"github.com/ryftcore/dss-go/dss/utils"
)

// AbstractASiCContainerAnalyzerOverrides declares the ASiC-specific abstract operations
// AbstractASiCContainerAnalyzer calls back into virtually, distinct from (but composed
// alongside) analyzer.DefaultDocumentAnalyzerOverrides - see this type's doc comment for the
// two-tier composition contract every concrete analyzer (CADSIGN/XADSIGN chunks) must follow.
type AbstractASiCContainerAnalyzerOverrides interface {
	// IsSupportedASiCContent checks if the ASiCContent is supported by the current validator.
	// Port of the public abstract isSupported(ASiCContent).
	IsSupportedASiCContent(asicContent *ASiCContent) bool

	// GetContainerExtractor returns the relevant container extractor. Port of the protected
	// abstract getContainerExtractor().
	GetContainerExtractor() *DefaultASiCContainerExtractor

	// GetManifestFilesDescriptions returns a list of parser ManifestFiles. Port of the
	// protected abstract getManifestFilesDescriptions().
	GetManifestFilesDescriptions() []*model.ManifestFile

	// GetSignatureAnalyzers returns a list of analyzers for signature documents embedded into
	// the container. Port of the protected abstract getSignatureAnalyzers().
	GetSignatureAnalyzers() []analyzer.DocumentAnalyzer

	// AttachExternalTimestamps attaches existing external timestamps to allSignatures. Port of
	// the protected attachExternalTimestamps(List), which - unlike the four methods above - is
	// NOT abstract in Java: the base supplies an empty body and only ASiCContainerWithCAdESAnalyzer
	// overrides it. It is listed here because GetAllSignatures self-calls it, so it needs the
	// same virtual dispatch the shadowed methods get; leaf analyzers that do not override it
	// satisfy this method through the promoted default body below, exactly as a Java subclass
	// inherits the base implementation.
	AttachExternalTimestamps(allSignatures []validation.AdvancedSignature) []*validation.TimestampToken
}

// AbstractASiCContainerAnalyzer is the abstract class for an ASiC container validation. Ports
// the Java class extending analyzer.DefaultDocumentAnalyzer.
//
// Two-tier virtual dispatch (per S7_BRIEF.md's warning): AbstractASiCContainerAnalyzer itself
// overrides 5 of analyzer.DefaultDocumentAnalyzerOverrides' methods (BuildSignatures,
// BuildDetachedEvidenceRecords, CoversSignature, AddReference, GetAllSignatures - defined
// directly below, shadowing the embedded DefaultDocumentAnalyzer's default bodies) while
// leaving 2 REQUIRED methods of that same interface unimplemented (IsSupported,
// OriginalDocumentsForSignature - never overridden by the Java class either) plus its OWN 4
// abstract methods (AbstractASiCContainerAnalyzerOverrides above). The concrete leaf analyzer
// (e.g. ASiCContainerWithCAdESAnalyzer) must therefore:
//  1. embed *AbstractASiCContainerAnalyzer;
//  2. implement IsSupported(model.DSSDocument) bool and
//     OriginalDocumentsForSignature(validation.AdvancedSignature) []model.DSSDocument itself
//     (Go promotes AbstractASiCContainerAnalyzer's 5 overridden methods automatically, so the
//     leaf value as a whole then satisfies the full analyzer.DefaultDocumentAnalyzerOverrides
//     interface);
//  3. implement AbstractASiCContainerAnalyzerOverrides' 4 abstract methods itself (its fifth
//     member, AttachExternalTimestamps, is optional: the base supplies a default body that
//     Go promotes, so only ASiCContainerWithCAdESAnalyzer declares its own);
//  4. call InitAbstractASiCContainerAnalyzer(leaf) on the embedded base, then
//     embeddedBase.InitDefaultDocumentAnalyzer(leaf) on its further-embedded
//     analyzer.DefaultDocumentAnalyzer, then InitFromDocument/InitFromContent (in that order -
//     the latter two call back into GetContainerExtractor, an overrides method).
type AbstractASiCContainerAnalyzer struct {
	analyzer.DefaultDocumentAnalyzer

	// overrides points back at the concrete analyzer; see InitAbstractASiCContainerAnalyzer.
	overrides AbstractASiCContainerAnalyzerOverrides

	// AsicContent is the container extraction result. Java declares the field protected;
	// exported here since Go subclasses live in different packages.
	AsicContent *ASiCContent

	// SignatureValidators is the list of signature document analyzers. Java declares the
	// field protected.
	SignatureValidators []analyzer.DocumentAnalyzer

	// TimestampAnalyzers is the list of timestamp document analyzers. Java declares the field
	// protected.
	TimestampAnalyzers []analyzertimestamp.TimestampAnalyzer

	// EvidenceRecordAnalyzers is the list of evidence record document analyzers. Java declares
	// the field protected.
	EvidenceRecordAnalyzers []analyzer.EvidenceRecordAnalyzer

	// manifestFiles is the list of manifest files, cached lazily by ManifestFiles.
	manifestFiles    []*model.ManifestFile
	manifestFilesSet bool
}

// NewAbstractASiCContainerAnalyzerBase builds the empty base state a subclass embeds. Port of
// the protected empty constructor. The subclass constructor must follow it with
// InitAbstractASiCContainerAnalyzer, InitDefaultDocumentAnalyzer (on the further-embedded
// DefaultDocumentAnalyzer) and then InitFromDocument or InitFromContent - see this type's doc
// comment.
func NewAbstractASiCContainerAnalyzerBase() *AbstractASiCContainerAnalyzer {
	return &AbstractASiCContainerAnalyzer{
		DefaultDocumentAnalyzer: analyzer.NewDefaultDocumentAnalyzerBase(),
	}
}

// InitAbstractASiCContainerAnalyzer registers the concrete analyzer with its base so the base
// can dispatch IsSupportedASiCContent/GetContainerExtractor/GetManifestFilesDescriptions/
// GetSignatureAnalyzers.
func (a *AbstractASiCContainerAnalyzer) InitAbstractASiCContainerAnalyzer(overrides AbstractASiCContainerAnalyzerOverrides) {
	a.overrides = overrides
}

func (a *AbstractASiCContainerAnalyzer) requireOverrides() AbstractASiCContainerAnalyzerOverrides {
	if a.overrides == nil {
		panic("AbstractASiCContainerAnalyzer was not initialised: the concrete analyzer must call InitAbstractASiCContainerAnalyzer in its constructor")
	}
	return a.overrides
}

// AbstractASiCContainerAnalyzerBase returns the receiver itself. Cross-package accessor added
// during phase 8f un-gating: AbstractASiCContainerValidator (abstract_asic_container_validator.go,
// same package) needs to recover this base pointer from the analyzer.DocumentAnalyzer interface
// value SignedDocumentValidatorBase.DocumentAnalyzer() returns, whose dynamic type is a concrete
// leaf analyzer (e.g. asic/cades.ASiCContainerWithCAdESAnalyzer) embedding
// *AbstractASiCContainerAnalyzer by pointer - a type assertion straight to
// *AbstractASiCContainerAnalyzer fails for that dynamic type (different concrete struct), so an
// exported method every embedder promotes automatically is needed instead, following the same
// pattern as jades.jwsDocumentAnalyzerBase. Purely additive.
func (a *AbstractASiCContainerAnalyzer) AbstractASiCContainerAnalyzerBase() *AbstractASiCContainerAnalyzer {
	return a
}

// InitFromDocument ports the protected AbstractASiCContainerAnalyzer(DSSDocument) constructor,
// split out because it calls back into GetContainerExtractor (an overrides method).
func (a *AbstractASiCContainerAnalyzer) InitFromDocument(document model.DSSDocument) {
	a.SetDocument(document)
	a.AsicContent = a.extractEntries()
}

// InitFromContent ports the protected AbstractASiCContainerAnalyzer(ASiCContent) constructor.
func (a *AbstractASiCContainerAnalyzer) InitFromContent(asicContent *ASiCContent) {
	a.SetDocument(asicContent.AsicContainer())
	a.AsicContent = asicContent
}

// extractEntries extracts documents from a container. Ports the private extractEntries().
func (a *AbstractASiCContainerAnalyzer) extractEntries() *ASiCContent {
	extractor := a.requireOverrides().GetContainerExtractor()
	content, err := extractor.Extract()
	if err != nil {
		panic(err)
	}
	return content
}

// GetContainerInfo allows retrieving the container information (ASiC Container). Ports the
// protected getContainerInfo().
func (a *AbstractASiCContainerAnalyzer) GetContainerInfo() *model.ContainerInfo {
	containerInfo := model.NewContainerInfo()
	containerInfo.SetContainerType(a.AsicContent.ContainerType())
	containerInfo.SetZipComment(a.AsicContent.ZipComment())

	mimeTypeDocument := a.AsicContent.MimeTypeDocument()
	if mimeTypeDocument != nil {
		content, err := spi.DSSUtilsToByteArrayOfDocument(mimeTypeDocument)
		if err == nil {
			containerInfo.SetMimeTypeContent(string(content))
		}
	}

	originalSignedDocuments := a.AsicContent.SignedDocuments()
	if utils.IsCollectionNotEmpty(originalSignedDocuments) {
		signedDocumentFilenames := make([]string, 0, len(originalSignedDocuments))
		for _, dssDocument := range originalSignedDocuments {
			signedDocumentFilenames = append(signedDocumentFilenames, dssDocument.Name())
		}
		containerInfo.SetSignedDocumentFilenames(signedDocumentFilenames)
	}

	containerInfo.SetManifestFiles(a.ManifestFiles())

	return containerInfo
}

// AttachExternalTimestamps attaches existing external timestamps to the list of
// AdvancedSignatures. Default: not applicable (used only in ASiC CAdES). Ports the protected
// attachExternalTimestamps(List).
func (a *AbstractASiCContainerAnalyzer) AttachExternalTimestamps(allSignatures []validation.AdvancedSignature) []*validation.TimestampToken {
	return []*validation.TimestampToken{}
}

// GetAllSignatures ports the @Override getAllSignatures(). Shadows the embedded
// DefaultDocumentAnalyzer.GetAllSignatures.
//
// The AttachExternalTimestamps call goes through requireOverrides(): Java's is a virtual call,
// and ASiCContainerWithCAdESAnalyzer overrides it to attach container-level (ASiC-S container /
// ASiC-E archive) timestamps to the signatures they cover. Calling a.AttachExternalTimestamps
// directly would bind to the empty default below and silently drop those timestamps - the
// return value is discarded here exactly as upstream discards it, because the work the override
// does is the side effect on the AdvancedSignatures in allSignatureList.
func (a *AbstractASiCContainerAnalyzer) GetAllSignatures() []validation.AdvancedSignature {
	allSignatureList := a.DefaultDocumentAnalyzer.GetAllSignatures()
	a.requireOverrides().AttachExternalTimestamps(allSignatureList)
	return allSignatureList
}

// BuildSignatures ports the @Override protected buildSignatures(). Shadows the embedded
// DefaultDocumentAnalyzer.BuildSignatures.
func (a *AbstractASiCContainerAnalyzer) BuildSignatures() []validation.AdvancedSignature {
	signatureList := make([]validation.AdvancedSignature, 0)
	for _, an := range a.requireOverrides().GetSignatureAnalyzers() {
		signatureList = append(signatureList, an.Signatures()...)
	}
	return signatureList
}

// GetContainerType returns a container type. Ports getContainerType().
func (a *AbstractASiCContainerAnalyzer) GetContainerType() enumerations.ASiCContainerType {
	return a.AsicContent.ContainerType()
}

// GetAllDocuments returns a list of all embedded documents. Ports getAllDocuments().
func (a *AbstractASiCContainerAnalyzer) GetAllDocuments() []model.DSSDocument {
	return a.AsicContent.AllDocuments()
}

// GetSignatureDocuments returns a list of embedded signature documents. Ports
// getSignatureDocuments().
func (a *AbstractASiCContainerAnalyzer) GetSignatureDocuments() []model.DSSDocument {
	return a.AsicContent.SignatureDocuments()
}

// GetSignedDocuments returns a list of embedded signed documents. Ports getSignedDocuments().
func (a *AbstractASiCContainerAnalyzer) GetSignedDocuments() []model.DSSDocument {
	return a.AsicContent.SignedDocuments()
}

// GetManifestDocuments returns a list of embedded signature manifest documents. Ports
// getManifestDocuments().
func (a *AbstractASiCContainerAnalyzer) GetManifestDocuments() []model.DSSDocument {
	return a.AsicContent.ManifestDocuments()
}

// GetTimestampDocuments returns a list of embedded timestamp documents. Ports
// getTimestampDocuments().
func (a *AbstractASiCContainerAnalyzer) GetTimestampDocuments() []model.DSSDocument {
	return a.AsicContent.TimestampDocuments()
}

// GetEvidenceRecordDocuments returns a list of embedded evidence record documents. Ports
// getEvidenceRecordDocuments().
func (a *AbstractASiCContainerAnalyzer) GetEvidenceRecordDocuments() []model.DSSDocument {
	return a.AsicContent.EvidenceRecordDocuments()
}

// GetArchiveManifestDocuments returns a list of embedded archive manifest documents. Ports
// getArchiveManifestDocuments().
func (a *AbstractASiCContainerAnalyzer) GetArchiveManifestDocuments() []model.DSSDocument {
	return a.AsicContent.ArchiveManifestDocuments()
}

// GetEvidenceRecordManifestDocuments returns a list of embedded evidence record manifest
// documents. Ports getEvidenceRecordManifestDocuments().
func (a *AbstractASiCContainerAnalyzer) GetEvidenceRecordManifestDocuments() []model.DSSDocument {
	return a.AsicContent.EvidenceRecordManifestDocuments()
}

// GetAllManifestDocuments returns a list of all embedded manifest documents. Ports
// getAllManifestDocuments().
func (a *AbstractASiCContainerAnalyzer) GetAllManifestDocuments() []model.DSSDocument {
	return a.AsicContent.AllManifestDocuments()
}

// GetArchiveDocuments returns a list of archive documents embedded the container. Ports
// getArchiveDocuments().
func (a *AbstractASiCContainerAnalyzer) GetArchiveDocuments() []model.DSSDocument {
	return a.AsicContent.ContainerDocuments()
}

// GetMimeTypeDocument returns a mimetype document. Ports getMimeTypeDocument().
func (a *AbstractASiCContainerAnalyzer) GetMimeTypeDocument() model.DSSDocument {
	return a.AsicContent.MimeTypeDocument()
}

// GetUnsupportedDocuments returns a list of unsupported documents from the container. Ports
// getUnsupportedDocuments().
func (a *AbstractASiCContainerAnalyzer) GetUnsupportedDocuments() []model.DSSDocument {
	return a.AsicContent.UnsupportedDocuments()
}

// ManifestFiles returns a list of parser Manifest files, cached after the first call. Ports
// getManifestFiles().
func (a *AbstractASiCContainerAnalyzer) ManifestFiles() []*model.ManifestFile {
	if !a.manifestFilesSet {
		a.manifestFiles = a.requireOverrides().GetManifestFilesDescriptions()
		a.manifestFilesSet = true
	}
	return a.manifestFiles
}

// GetSignedDocumentsASiCS returns a list of "package.zip" documents. Ports the protected
// getSignedDocumentsASiCS(List).
func (a *AbstractASiCContainerAnalyzer) GetSignedDocumentsASiCS(retrievedDocs []model.DSSDocument) []model.DSSDocument {
	containerDocuments := a.AsicContent.ContainerDocuments()
	if utils.IsCollectionNotEmpty(containerDocuments) {
		return containerDocuments
	}
	return retrievedDocs
}

// BuildDetachedEvidenceRecords ports the @Override protected buildDetachedEvidenceRecords().
// Shadows the embedded DefaultDocumentAnalyzer.BuildDetachedEvidenceRecords. Logging (LOG.warn)
// is dropped per PORTING.md.
func (a *AbstractASiCContainerAnalyzer) BuildDetachedEvidenceRecords() []validation.EvidenceRecord {
	embeddedEvidenceRecords := make([]validation.EvidenceRecord, 0)
	for _, evidenceRecordAnalyzer := range a.GetEvidenceRecordAnalyzers() {
		evidenceRecord := evidenceRecordAnalyzer.EvidenceRecord()
		if evidenceRecord != nil {
			embeddedEvidenceRecords = append(embeddedEvidenceRecords, evidenceRecord)
		}
	}
	detachedEvidenceRecords := append([]validation.EvidenceRecord{}, a.DefaultDocumentAnalyzer.BuildDetachedEvidenceRecords()...)
	a.AttachExternalEvidenceRecords(embeddedEvidenceRecords, detachedEvidenceRecords)
	// return all
	detachedEvidenceRecords = append(detachedEvidenceRecords, embeddedEvidenceRecords...)
	return detachedEvidenceRecords
}

// AttachExternalEvidenceRecords appends detached evidence record provided to the validator to
// the evidence records covered by the corresponding evidence records. Ports the protected
// attachExternalEvidenceRecords(List, List).
func (a *AbstractASiCContainerAnalyzer) AttachExternalEvidenceRecords(embeddedEvidenceRecords, detachedEvidenceRecords []validation.EvidenceRecord) {
	if utils.IsCollectionNotEmpty(embeddedEvidenceRecords) {
		for _, coveredEvidenceRecord := range embeddedEvidenceRecords {
			for _, coveringEvidenceRecord := range embeddedEvidenceRecords {
				if a.coversEvidenceRecord(coveredEvidenceRecord, coveringEvidenceRecord) {
					coveredEvidenceRecord.AddExternalEvidenceRecord(coveringEvidenceRecord)
				}
			}
			// assert all detached evidence records cover embedded data
			for _, coveringEvidenceRecord := range detachedEvidenceRecords {
				coveredEvidenceRecord.AddExternalEvidenceRecord(coveringEvidenceRecord)
			}
		}
	}
}

// GetEvidenceRecordAnalyzers builds and returns a list of evidence record analyzers, cached
// after the first call. Ports the protected getEvidenceRecordAnalyzers().
func (a *AbstractASiCContainerAnalyzer) GetEvidenceRecordAnalyzers() []analyzer.EvidenceRecordAnalyzer {
	if a.EvidenceRecordAnalyzers == nil {
		a.EvidenceRecordAnalyzers = make([]analyzer.EvidenceRecordAnalyzer, 0)
		for _, evidenceRecordDocument := range a.GetEvidenceRecordDocuments() {
			evidenceRecordAnalyzer := a.getEvidenceRecordAnalyzer(evidenceRecordDocument)
			if evidenceRecordAnalyzer != nil {
				a.EvidenceRecordAnalyzers = append(a.EvidenceRecordAnalyzers, evidenceRecordAnalyzer)
			}
		}
	}
	return a.EvidenceRecordAnalyzers
}

// getEvidenceRecordAnalyzer ports the private getEvidenceRecordAnalyzer(DSSDocument). Logging
// (LOG.warn) is dropped per PORTING.md; the caught exception is swallowed exactly as Java
// swallows it (returns nil), matching the original control flow.
func (a *AbstractASiCContainerAnalyzer) getEvidenceRecordAnalyzer(evidenceRecordDocument model.DSSDocument) (result analyzer.EvidenceRecordAnalyzer) {
	defer func() {
		if recover() != nil {
			result = nil
		}
	}()

	var manifestFile *model.ManifestFile
	detachedContents := a.GetAllDocuments()

	evidenceRecordManifest := ASiCManifestParserGetLinkedManifest(a.GetEvidenceRecordManifestDocuments(), evidenceRecordDocument.Name())
	if evidenceRecordManifest != nil {
		manifestFile = a.GetValidatedManifestFile(evidenceRecordManifest)
	}

	isASiCSContainer, err := ASiCUtilsIsASiCSContainerContent(a.AsicContent)
	if err != nil {
		return nil
	}
	if isASiCSContainer {
		if manifestFile != nil {
			manifestFile = nil
		}
		rootLevelSignedDocuments := ASiCUtilsRootLevelSignedDocuments(a.AsicContent)
		if len(rootLevelSignedDocuments) == 1 {
			detachedContents = rootLevelSignedDocuments
		} else {
			detachedContents = []model.DSSDocument{}
		}

	} else {
		if manifestFile == nil {
			detachedContents = []model.DSSDocument{}
			manifestFile = model.NewManifestFile() // empty manifest
		}
	}

	evidenceRecordAnalyzer, err := analyzer.EvidenceRecordAnalyzerFromDocument(evidenceRecordDocument)
	if err != nil {
		return nil
	}
	a.assertEvidenceRecordDocumentExtensionMatch(evidenceRecordDocument, evidenceRecordAnalyzer.EvidenceRecordType())
	evidenceRecordAnalyzer.SetDetachedContents(detachedContents)
	evidenceRecordAnalyzer.SetManifestFile(manifestFile)
	evidenceRecordAnalyzer.SetCertificateVerifier(a.CertificateVerifier())
	evidenceRecordAnalyzer.SetEvidenceRecordOrigin(enumerations.EvidenceRecordOrigin_CONTAINER)
	return evidenceRecordAnalyzer
}

// assertEvidenceRecordDocumentExtensionMatch verifies whether the extension of
// evidenceRecordDocument is conformant to the applicable standard for the given
// evidenceRecordTypeEnum. Ports the protected assertEvidenceRecordDocumentExtensionMatch(
// DSSDocument, EvidenceRecordTypeEnum).
//
// Cross-chunk assumption (ZIPCORE): the ".xml"/".ers" extension literals below match
// ASiCUtils.XML_EXTENSION / ASiCUtils.ER_ASN1_EXTENSION verbatim (both are plain ".xml"/".ers"
// suffixes per the ASiC filename convention); if ZIPCORE names these differently the literals
// still produce identical behavior.
//
// Panics with a *model.DSSError on a mismatched extension, or Java's
// UnsupportedOperationException message for an unrecognized evidenceRecordTypeEnum.
func (a *AbstractASiCContainerAnalyzer) assertEvidenceRecordDocumentExtensionMatch(evidenceRecordDocument model.DSSDocument, evidenceRecordTypeEnum enumerations.EvidenceRecordTypeEnum) {
	switch evidenceRecordTypeEnum {
	case enumerations.EvidenceRecordTypeEnum_XML_EVIDENCE_RECORD:
		if evidenceRecordDocument.Name() != "" && !hasSuffix(evidenceRecordDocument.Name(), ".xml") {
			panic(model.NewDSSError("Document containing an XMLERS evidence record shall end with '.xml' extension!"))
		}
	case enumerations.EvidenceRecordTypeEnum_ASN1_EVIDENCE_RECORD:
		if evidenceRecordDocument.Name() != "" && !hasSuffix(evidenceRecordDocument.Name(), ".ers") {
			panic(model.NewDSSError("Document containing an ERS evidence record shall end with '.ers' extension!"))
		}
	default:
		panic("The evidence record type '" + string(evidenceRecordTypeEnum) + "' is not supported!")
	}
}

// hasSuffix is a tiny local helper avoiding an extra "strings" import for a single call site
// used twice in this file (not a cross-file shared helper, per PORTING.md).
func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// CoversSignature ports the @Override protected coversSignature(AdvancedSignature,
// EvidenceRecord). Shadows the embedded DefaultDocumentAnalyzer.CoversSignature.
func (a *AbstractASiCContainerAnalyzer) CoversSignature(signature validation.AdvancedSignature, evidenceRecord validation.EvidenceRecord) bool {
	documentName := ""
	if a.Document() != nil {
		documentName = a.Document().Name()
	}
	return a.coversFile(evidenceRecord, documentName) || a.coversFile(evidenceRecord, signature.Filename())
}

// coversEvidenceRecord ports the private coversEvidenceRecord(EvidenceRecord, EvidenceRecord).
func (a *AbstractASiCContainerAnalyzer) coversEvidenceRecord(coveredEvidenceRecord, coveringEvidenceRecord validation.EvidenceRecord) bool {
	return a.coversFile(coveringEvidenceRecord, coveredEvidenceRecord.Filename())
}

// coversFile ports the private coversFile(EvidenceRecord, String).
func (a *AbstractASiCContainerAnalyzer) coversFile(evidenceRecord validation.EvidenceRecord, filename string) bool {
	for _, referenceValidation := range evidenceRecord.ReferenceValidation() {
		referenceDocument := referenceValidation.Document()
		if filename == "" || (referenceDocument != nil && filename == referenceDocument.Name()) {
			return true
		}
	}
	return evidenceRecord.ManifestFile() != nil && a.coversManifestFile(evidenceRecord.ManifestFile(), filename)
}

// coversManifestFile ports the private coversFile(ManifestFile, String) overload.
func (a *AbstractASiCContainerAnalyzer) coversManifestFile(manifestFile *model.ManifestFile, filename string) bool {
	if manifestFile != nil {
		for _, manifestEntry := range manifestFile.Entries() {
			if filename == manifestEntry.Uri() {
				return true
			}
		}
	}
	return false
}

// GetValidatedManifestFile returns a validated ManifestFile for the given manifest document.
// Ports the protected getValidatedManifestFile(DSSDocument).
func (a *AbstractASiCContainerAnalyzer) GetValidatedManifestFile(manifest model.DSSDocument) *model.ManifestFile {
	allManifestFiles := a.ManifestFiles()
	if utils.IsCollectionNotEmpty(allManifestFiles) {
		for _, manifestFile := range allManifestFiles {
			if manifest.Name() == manifestFile.Filename() {
				return manifestFile
			}
		}
	}
	return nil
}

// AddReference ports the @Override protected addReference(SignatureScope). Shadows the
// embedded DefaultDocumentAnalyzer.AddReference.
//
// Cross-chunk assumption (ZIPCORE): ASiCUtilsIsSignature/ASiCUtilsIsTimestamp/
// ASiCUtilsIsEvidenceRecord take a filename string.
func (a *AbstractASiCContainerAnalyzer) AddReference(signatureScope scope.SignatureScope) bool {
	fileName := signatureScope.DocumentName()
	return fileName == "" || (!ASiCUtilsIsSignature(fileName) && !ASiCUtilsIsTimestamp(fileName) && !ASiCUtilsIsEvidenceRecord(fileName))
}
