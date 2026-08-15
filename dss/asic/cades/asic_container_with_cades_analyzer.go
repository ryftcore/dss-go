//go:build phase8

// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/validation/ASiCContainerWithCAdESAnalyzer.java (DSS 6.5.RC1).
//
// INTEGRATOR NOTE (Phase 7 integration): gated behind the `phase8` build tag because it
// constructs ASiCWithCAdESTimestampAnalyzer, itself gated (see that file's header) since its
// Java base belongs to the not-yet-ported dss-validation module. Drop the tag on both files
// together once Phase 8 lands the package.
//
// AttachExternalTimestamps below is reached virtually: asic.AbstractASiCContainerAnalyzer's
// GetAllSignatures self-calls it through AbstractASiCContainerAnalyzerOverrides, which the
// method is a member of precisely so this leaf's override is not lost to Go's static dispatch
// (the "virtual-dispatch warning" bug class S7_BRIEF.md calls out). Analyzers that do not
// override it - ASiCContainerWithXAdESAnalyzer - inherit the base's empty body by promotion,
// matching Java's non-abstract protected default.
//
// FLAGGED CROSS-CHUNK GAP: Java's getSignatureAnalyzers() forwards
// `this.getSignaturePolicyProvider()` (a protected accessor on the frozen
// analyzer.DefaultDocumentAnalyzer, unexported in the Go port as
// signaturePolicyProviderOrDefault and not reachable from another package) into each nested
// CMSDocumentAnalyzer. No exported equivalent exists on analyzer.DefaultDocumentAnalyzer today,
// so this propagation is dropped here: each nested CMSDocumentAnalyzer instead lazily
// instantiates its own default SignaturePolicyProvider. This only differs observably when a
// caller has set a *custom* SignaturePolicyProvider on the outer analyzer via
// SetSignaturePolicyProvider - flagging for the integrator to add an exported getter upstream.
package cades

import (
	"github.com/utain/esig/dss/asic"
	dsscades "github.com/utain/esig/dss/cades"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/spi/validation/analyzer"
	analyzertimestamp "github.com/utain/esig/dss/spi/validation/analyzer/timestamp"
	timestampsrc "github.com/utain/esig/dss/spi/validation/timestamp"
	"github.com/utain/esig/dss/utils"
)

// ASiCContainerWithCAdESAnalyzer is an implementation to validate ASiC containers with CAdES
// signature(s).
type ASiCContainerWithCAdESAnalyzer struct {
	*asic.AbstractASiCContainerAnalyzer
}

var _ asic.AbstractASiCContainerAnalyzerOverrides = (*ASiCContainerWithCAdESAnalyzer)(nil)
var _ analyzer.DocumentAnalyzer = (*ASiCContainerWithCAdESAnalyzer)(nil)

// newASiCContainerWithCAdESAnalyzer is the empty constructor. Port of the package-private empty
// constructor.
func newASiCContainerWithCAdESAnalyzer() *ASiCContainerWithCAdESAnalyzer {
	a := &ASiCContainerWithCAdESAnalyzer{
		AbstractASiCContainerAnalyzer: asic.NewAbstractASiCContainerAnalyzerBase(),
	}
	a.InitAbstractASiCContainerAnalyzer(a)
	a.InitDefaultDocumentAnalyzer(a)
	return a
}

// NewASiCContainerWithCAdESAnalyzer is the default constructor. Ports
// ASiCContainerWithCAdESAnalyzer(DSSDocument).
func NewASiCContainerWithCAdESAnalyzer(asicContainer model.DSSDocument) *ASiCContainerWithCAdESAnalyzer {
	a := newASiCContainerWithCAdESAnalyzer()
	a.InitFromDocument(asicContainer)
	return a
}

// NewASiCContainerWithCAdESAnalyzerFromContent is the constructor with ASiCContent. Ports
// ASiCContainerWithCAdESAnalyzer(ASiCContent).
func NewASiCContainerWithCAdESAnalyzerFromContent(asicContent *asic.ASiCContent) *ASiCContainerWithCAdESAnalyzer {
	a := newASiCContainerWithCAdESAnalyzer()
	a.InitFromContent(asicContent)
	return a
}

// IsSupported ports the @Override isSupported(DSSDocument).
func (a *ASiCContainerWithCAdESAnalyzer) IsSupported(dssDocument model.DSSDocument) bool {
	return NewASiCWithCAdESFormatDetector().IsSupportedASiC(dssDocument)
}

// IsSupportedASiCContent ports the @Override isSupported(ASiCContent), implementing
// asic.AbstractASiCContainerAnalyzerOverrides.
func (a *ASiCContainerWithCAdESAnalyzer) IsSupportedASiCContent(asicContent *asic.ASiCContent) bool {
	return NewASiCWithCAdESFormatDetector().IsSupportedASiCContent(asicContent)
}

// GetContainerExtractor ports the @Override protected getContainerExtractor(), implementing
// asic.AbstractASiCContainerAnalyzerOverrides.
func (a *ASiCContainerWithCAdESAnalyzer) GetContainerExtractor() *asic.DefaultASiCContainerExtractor {
	return &NewASiCWithCAdESContainerExtractor(a.Document()).DefaultASiCContainerExtractor
}

// GetSignatureAnalyzers ports the @Override protected getSignatureAnalyzers(), implementing
// asic.AbstractASiCContainerAnalyzerOverrides. See the file header's flagged gap regarding the
// dropped SignaturePolicyProvider propagation.
func (a *ASiCContainerWithCAdESAnalyzer) GetSignatureAnalyzers() []analyzer.DocumentAnalyzer {
	if a.SignatureValidators == nil {
		a.SignatureValidators = make([]analyzer.DocumentAnalyzer, 0)
		for _, signature := range a.GetSignatureDocuments() {
			cadesValidator, err := dsscades.NewCMSDocumentAnalyzerFromDocument(signature)
			if err != nil {
				panic(err)
			}
			cadesValidator.SetCertificateVerifier(a.CertificateVerifier())
			cadesValidator.SetContainerContents(a.GetArchiveDocuments())

			signedDocument := ASiCWithCAdESUtilsGetSignedDocument(a.AsicContent, signature.Name())
			if signedDocument != nil {
				cadesValidator.SetDetachedContents([]model.DSSDocument{signedDocument})
			}

			signatureManifest := asic.ASiCManifestParserGetLinkedManifest(a.GetAllManifestDocuments(), signature.Name())
			if signatureManifest != nil {
				manifestFile := a.GetValidatedManifestFile(signatureManifest)
				cadesValidator.SetManifestFile(manifestFile)
			}

			a.SignatureValidators = append(a.SignatureValidators, cadesValidator)
		}
	}
	return a.SignatureValidators
}

// GetTimestampAnalyzers returns a list of timestamp validators for timestamps embedded into the
// container. Ports the protected getTimestampAnalyzers().
func (a *ASiCContainerWithCAdESAnalyzer) GetTimestampAnalyzers() []analyzertimestamp.TimestampAnalyzer {
	if a.TimestampAnalyzers == nil {
		a.TimestampAnalyzers = make([]analyzertimestamp.TimestampAnalyzer, 0)
		for _, timestamp := range a.GetTimestampDocuments() {
			timestampValidator := a.getTimestampValidator(timestamp)
			if timestampValidator != nil {
				a.TimestampAnalyzers = append(a.TimestampAnalyzers, timestampValidator)
			}
		}
		comparator := analyzertimestamp.NewTimestampAnalyzerComparator()
		sortTimestampAnalyzers(a.TimestampAnalyzers, comparator)
	}
	return a.TimestampAnalyzers
}

// sortTimestampAnalyzers ports the .sort(new TimestampAnalyzerComparator()) call, local to this
// file per PORTING.md (no cross-file shared helpers) via a plain insertion sort over the
// comparator's Less.
func sortTimestampAnalyzers(items []analyzertimestamp.TimestampAnalyzer, comparator analyzertimestamp.TimestampAnalyzerComparator) {
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && comparator.Less(items[j], items[j-1]); j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
}

// getTimestampValidator ports the private getTimestampValidator(DSSDocument). Logging
// (LOG.warn) is dropped per PORTING.md; the branches it sits in are preserved.
func (a *ASiCContainerWithCAdESAnalyzer) getTimestampValidator(timestampDocument model.DSSDocument) *ASiCWithCAdESTimestampAnalyzer {
	var timestampedDocument model.DSSDocument
	var manifestFile *model.ManifestFile
	var archiveTimestampType enumerations.ArchiveTimestampType
	var archiveTimestampTypeSet bool

	archiveManifest := asic.ASiCManifestParserGetLinkedManifest(a.GetAllManifestDocuments(), timestampDocument.Name())
	if archiveManifest != nil {
		timestampedDocument = archiveManifest
		manifestFile = a.GetValidatedManifestFile(archiveManifest)
		if manifestFile != nil {
			if asic.ASiCUtilsCoversSignature(manifestFile) {
				archiveTimestampType = enumerations.ArchiveTimestampType_CAdES_DETACHED
				archiveTimestampTypeSet = true
			}
		}
		// (upstream warns "A linked manifest is not found for a timestamp with name [{}]!" when manifestFile is nil)

	} else {
		rootLevelSignedDocuments := a.AsicContent.RootLevelSignedDocuments()
		if utils.CollectionSize(rootLevelSignedDocuments) == 1 {
			timestampedDocument = rootLevelSignedDocuments[0]
		} else {
			// (upstream warns "Timestamp {} is skipped (no linked archive manifest found / unique file)")
			return nil
		}
	}

	timestampValidator := NewASiCWithCAdESTimestampAnalyzerWithType(timestampDocument, enumerations.TimestampType_CONTAINER_TIMESTAMP)
	timestampValidator.SetTimestampedData(timestampedDocument)
	timestampValidator.SetManifestFile(manifestFile)
	if archiveTimestampTypeSet {
		timestampValidator.SetArchiveTimestampType(archiveTimestampType)
	}
	timestampValidator.SetOriginalDocuments(a.GetAllDocuments())
	timestampValidator.SetArchiveDocuments(a.GetArchiveDocuments())
	timestampValidator.SetDetachedEvidenceRecords(a.DetachedEvidenceRecords())
	timestampValidator.SetCertificateVerifier(a.CertificateVerifier())
	return timestampValidator
}

// BuildDetachedTimestamps ports the @Override protected buildDetachedTimestamps(), implementing
// analyzer.DefaultDocumentAnalyzerOverrides (promoted from AbstractASiCContainerAnalyzer's
// embedded DefaultDocumentAnalyzer; shadowed here since ASiC-CAdES is the one format that
// supports container timestamps).
func (a *ASiCContainerWithCAdESAnalyzer) BuildDetachedTimestamps() []*validation.TimestampToken {
	detachedTimestampSource := timestampsrc.NewDetachedTimestampSource()
	for _, timestampAnalyzer := range a.GetTimestampAnalyzers() {
		_ = detachedTimestampSource.AddExternalTimestamp(timestampAnalyzer.Timestamp())
	}
	return detachedTimestampSource.DetachedTimestamps()
}

// GetArchiveDocuments ports the @Override getArchiveDocuments().
func (a *ASiCContainerWithCAdESAnalyzer) GetArchiveDocuments() []model.DSSDocument {
	archiveContents := a.AbstractASiCContainerAnalyzer.GetArchiveDocuments()
	// in case of Manifest file (ASiC-E CAdES signature) add signed documents
	if utils.IsCollectionNotEmpty(a.GetManifestDocuments()) {
		for _, document := range a.GetAllDocuments() {
			if !containsDocument(archiveContents, document) {
				archiveContents = append(archiveContents, document)
			}
		}
	}
	return archiveContents
}

// containsDocument ports the List.contains(document) check, local to this file per PORTING.md.
func containsDocument(documents []model.DSSDocument, target model.DSSDocument) bool {
	for _, document := range documents {
		if document == target {
			return true
		}
	}
	return false
}

// AttachExternalTimestamps ports the @Override protected attachExternalTimestamps(List). It is
// reached through AbstractASiCContainerAnalyzerOverrides when the base's GetAllSignatures runs;
// see this file's header.
func (a *ASiCContainerWithCAdESAnalyzer) AttachExternalTimestamps(allSignatures []validation.AdvancedSignature) []*validation.TimestampToken {
	externalTimestamps := make([]*validation.TimestampToken, 0)

	for _, tstAnalyzer := range a.GetTimestampAnalyzers() {
		timestamp := a.getExternalTimestamp(tstAnalyzer, allSignatures)
		if timestamp != nil {
			externalTimestamps = append(externalTimestamps, timestamp)
		}
	}

	return externalTimestamps
}

// getExternalTimestamp ports the private getExternalTimestamp(TimestampAnalyzer, List).
func (a *ASiCContainerWithCAdESAnalyzer) getExternalTimestamp(tstAnalyzer analyzertimestamp.TimestampAnalyzer, allSignatures []validation.AdvancedSignature) *validation.TimestampToken {
	timestampValidator, ok := tstAnalyzer.(*ASiCWithCAdESTimestampAnalyzer)
	if !ok {
		return nil
	}
	timestamp := timestampValidator.Timestamp()
	if timestamp.TimeStampType().IsContainerTimestamp() {
		coveredManifest := timestampValidator.GetCoveredManifest()
		if coveredManifest != nil {
			for _, entry := range coveredManifest.Entries() {
				cadesSignature := a.getCAdESSignatureFromFileName(allSignatures, entry.Uri())
				if cadesSignature != nil {
					cadesSignature.AddExternalTimestamp(timestamp)
				}
			}
		}
	}

	return timestamp
}

// getCAdESSignatureFromFileName ports the private getCAdESSignatureFromFileName(List, String).
func (a *ASiCContainerWithCAdESAnalyzer) getCAdESSignatureFromFileName(signatures []validation.AdvancedSignature, fileName string) *dsscades.CAdESSignature {
	for _, advancedSignature := range signatures {
		if utils.AreStringsEqual(fileName, advancedSignature.Filename()) && !advancedSignature.IsCounterSignature() {
			if cadesSignature, ok := advancedSignature.(*dsscades.CAdESSignature); ok {
				return cadesSignature
			}
			return nil
		}
	}
	return nil
}

// GetManifestFilesDescriptions ports the @Override protected getManifestFilesDescriptions(),
// implementing asic.AbstractASiCContainerAnalyzerOverrides.
func (a *ASiCContainerWithCAdESAnalyzer) GetManifestFilesDescriptions() []*model.ManifestFile {
	descriptions := make([]*model.ManifestFile, 0)

	for _, manifestDocument := range a.GetManifestDocuments() {
		manifestFile := asic.ASiCManifestParserGetManifestFile(manifestDocument)
		if manifestFile != nil {
			asiceWithCAdESManifestValidator := asic.NewASiCManifestValidator(manifestFile, a.GetAllDocuments())
			asiceWithCAdESManifestValidator.ValidateEntries()
			descriptions = append(descriptions, manifestFile)
		}
	}

	for _, manifestDocument := range a.GetArchiveManifestDocuments() {
		manifestFile := asic.ASiCManifestParserGetManifestFile(manifestDocument)
		if manifestFile != nil {
			manifestFile.SetManifestType(enumerations.ASiCManifestTypeEnum_ARCHIVE_MANIFEST)
			asiceWithCAdESManifestValidator := asic.NewASiCManifestValidator(manifestFile, a.GetAllDocuments())
			asiceWithCAdESManifestValidator.ValidateEntries()
			descriptions = append(descriptions, manifestFile)
		}
	}

	for _, manifestDocument := range a.GetEvidenceRecordManifestDocuments() {
		manifestFile := asic.ASiCManifestParserGetManifestFile(manifestDocument)
		if manifestFile != nil {
			manifestFile.SetManifestType(enumerations.ASiCManifestTypeEnum_EVIDENCE_RECORD)
			asiceWithCAdESManifestValidator := asic.NewASiCManifestValidator(manifestFile, a.GetAllDocuments())
			asiceWithCAdESManifestValidator.ValidateEntries()
			descriptions = append(descriptions, manifestFile)
		}
	}

	return descriptions
}

// OriginalDocumentsForSignature ports the @Override getOriginalDocuments(AdvancedSignature),
// implementing analyzer.DefaultDocumentAnalyzerOverrides (promoted from
// AbstractASiCContainerAnalyzer's embedded DefaultDocumentAnalyzer).
func (a *ASiCContainerWithCAdESAnalyzer) OriginalDocumentsForSignature(advancedSignature validation.AdvancedSignature) []model.DSSDocument {
	if advancedSignature.IsCounterSignature() {
		cadesSignature, ok := advancedSignature.(*dsscades.CAdESSignature)
		if ok {
			originalDocument, err := cadesSignature.OriginalDocument()
			if err == nil {
				return []model.DSSDocument{originalDocument}
			}
		}
		return []model.DSSDocument{}
	}
	retrievedDocs := advancedSignature.DetachedContents()
	if enumerations.ASiCContainerType_ASiC_S == a.GetContainerType() {
		return a.GetSignedDocumentsASiCS(retrievedDocs)
	}
	linkedManifest := asic.ASiCManifestParserGetLinkedManifest(a.GetManifestDocuments(), advancedSignature.Filename())
	if linkedManifest == nil {
		return []model.DSSDocument{}
	}
	manifestFile := asic.ASiCManifestParserGetManifestFile(linkedManifest)
	if manifestFile == nil {
		return []model.DSSDocument{}
	}
	return a.getManifestedDocuments(manifestFile)
}

// getManifestedDocuments ports the private getManifestedDocuments(ManifestFile).
func (a *ASiCContainerWithCAdESAnalyzer) getManifestedDocuments(manifestFile *model.ManifestFile) []model.DSSDocument {
	entries := manifestFile.Entries()
	signedDocuments := a.GetAllDocuments()

	result := make([]model.DSSDocument, 0)
	for _, entry := range entries {
		for _, signedDocument := range signedDocuments {
			if utils.AreStringsEqual(entry.Uri(), signedDocument.Name()) {
				result = append(result, signedDocument)
			}
		}
	}
	return result
}
