// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/merge/ASiCEWithCAdESContainerMerger.java (DSS 6.5.RC1).
package cades

import (
	"bytes"
	"fmt"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// ASiCEWithCAdESContainerMerger is used to merge ASiC-E with CAdES containers.
type ASiCEWithCAdESContainerMerger struct {
	AbstractASiCWithCAdESContainerMerger
}

var _ asic.DefaultContainerMergerOverrides = (*ASiCEWithCAdESContainerMerger)(nil)
var _ asic.ASiCContainerMerger = (*ASiCEWithCAdESContainerMerger)(nil)

// newASiCEWithCAdESContainerMerger is the empty constructor. Port of the package-private empty
// constructor.
func newASiCEWithCAdESContainerMerger() *ASiCEWithCAdESContainerMerger {
	m := &ASiCEWithCAdESContainerMerger{
		AbstractASiCWithCAdESContainerMerger: NewAbstractASiCWithCAdESContainerMergerBase(),
	}
	m.InitDefaultContainerMerger(m)
	return m
}

// NewASiCEWithCAdESContainerMerger creates an ASiC-E with CAdES container merger from provided
// container documents. Ports ASiCEWithCAdESContainerMerger(DSSDocument...).
func NewASiCEWithCAdESContainerMerger(containers ...model.DSSDocument) *ASiCEWithCAdESContainerMerger {
	m := newASiCEWithCAdESContainerMerger()
	m.InitFromDocuments(containers...)
	return m
}

// NewASiCEWithCAdESContainerMergerFromContents creates an ASiC-E with CAdES container merger
// from the given ASiCContents. Ports ASiCEWithCAdESContainerMerger(ASiCContent...).
func NewASiCEWithCAdESContainerMergerFromContents(asicContents ...*asic.ASiCContent) *ASiCEWithCAdESContainerMerger {
	m := newASiCEWithCAdESContainerMerger()
	m.InitFromASiCContents(asicContents...)
	return m
}

// IsSupportedDocument ports the @Override protected isSupported(DSSDocument).
func (m *ASiCEWithCAdESContainerMerger) IsSupportedDocument(container model.DSSDocument) bool {
	if !m.AbstractASiCWithCAdESContainerMerger.IsSupportedDocument(container) {
		return false
	}
	isASiCS, err := asic.ASiCUtilsIsASiCSContainer(container)
	if err != nil {
		panic(err)
	}
	return !isASiCS || (m.doesNotContainSignaturesDocument(container) &&
		m.doesNotContainTimestampsDocument(container) && m.doesNotContainEvidenceRecordsDocument(container))
}

func (m *ASiCEWithCAdESContainerMerger) doesNotContainSignaturesDocument(container model.DSSDocument) bool {
	entryNames, err := asic.ZipUtilsInstance().ExtractEntryNames(container)
	if err != nil {
		panic(err)
	}
	return !asic.ASiCUtilsFilesContainSignatures(entryNames)
}

func (m *ASiCEWithCAdESContainerMerger) doesNotContainTimestampsDocument(container model.DSSDocument) bool {
	entryNames, err := asic.ZipUtilsInstance().ExtractEntryNames(container)
	if err != nil {
		panic(err)
	}
	return !asic.ASiCUtilsFilesContainTimestamps(entryNames)
}

func (m *ASiCEWithCAdESContainerMerger) doesNotContainEvidenceRecordsDocument(container model.DSSDocument) bool {
	entryNames, err := asic.ZipUtilsInstance().ExtractEntryNames(container)
	if err != nil {
		panic(err)
	}
	return !asic.ASiCUtilsFilesContainEvidenceRecords(entryNames)
}

// IsSupportedContent ports the @Override protected isSupported(ASiCContent).
func (m *ASiCEWithCAdESContainerMerger) IsSupportedContent(asicContent *asic.ASiCContent) bool {
	if !m.AbstractASiCWithCAdESContainerMerger.IsSupportedContent(asicContent) {
		return false
	}
	isASiCS, err := asic.ASiCUtilsIsASiCSContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	return !isASiCS || (utils.IsCollectionEmpty(asicContent.SignatureDocuments()) &&
		utils.IsCollectionEmpty(asicContent.TimestampDocuments()) && utils.IsCollectionEmpty(asicContent.EvidenceRecordDocuments()))
}

// GetTargetASiCContainerType ports the @Override protected getTargetASiCContainerType().
func (m *ASiCEWithCAdESContainerMerger) GetTargetASiCContainerType() enumerations.ASiCContainerType {
	return enumerations.ASiCContainerTypeASiCE
}

// EnsureContainerContentAllowMerge ports the @Override protected
// ensureContainerContentAllowMerge().
func (m *ASiCEWithCAdESContainerMerger) EnsureContainerContentAllowMerge() {
	// no checks available
}

// EnsureSignaturesAllowMerge ports the @Override protected ensureSignaturesAllowMerge().
func (m *ASiCEWithCAdESContainerMerger) EnsureSignaturesAllowMerge() {
	containersWithSignaturesCount := 0
	for _, asicContent := range m.AsicContents {
		if utils.IsCollectionNotEmpty(asicContent.SignatureDocuments()) ||
			utils.IsCollectionNotEmpty(asicContent.TimestampDocuments()) ||
			utils.IsCollectionNotEmpty(asicContent.EvidenceRecordDocuments()) {
			containersWithSignaturesCount++
		}
	}
	if containersWithSignaturesCount <= 1 {
		// no signature, timestamp nor evidence record documents in all containers except maximum one. Can merge.
		return
	}

	m.ensureSignatureDocumentsValid()
	m.ensureEvidenceRecordDocumentsValid()
	m.ensureManifestDocumentsValid()
}

func (m *ASiCEWithCAdESContainerMerger) ensureSignatureDocumentsValid() {
	mergedSignatureNames := make([]string, 0)
	asicContentsToProcess := append([]*asic.ASiCContent{}, m.AsicContents...)

	for len(asicContentsToProcess) > 0 {
		asicContent := asicContentsToProcess[0]
		asicContentsToProcess = asicContentsToProcess[1:] // remove entry to avoid recursive comparison

		signatureDocumentList := append([]model.DSSDocument{}, asicContent.SignatureDocuments()...)
		for _, signatureDocument := range signatureDocumentList {
			if containsString(mergedSignatureNames, signatureDocument.Name()) {
				continue
			}

			signaturesToMerge := m.getSignatureDocumentsToBeMerged(asicContent, signatureDocument, asicContentsToProcess)
			if utils.IsCollectionNotEmpty(signaturesToMerge) {
				signaturesToMerge = append(signaturesToMerge, signatureDocument)
				mergedSignatureNames = append(mergedSignatureNames, signatureDocument.Name())

				signaturesCms := m.MergeCmsSignatures(signaturesToMerge)
				m.updateMergedSignatureInContainers(signaturesCms)
			}
		}
	}
}

func (m *ASiCEWithCAdESContainerMerger) getSignatureDocumentsToBeMerged(currentASiCContent *asic.ASiCContent,
	currentSignatureDocument model.DSSDocument, asicContentList []*asic.ASiCContent) []model.DSSDocument {
	if currentSignatureDocument.Name() == "" {
		panic("Name shall be provided for a document!")
	}
	manifest := asic.ASiCManifestParserGetLinkedManifest(currentASiCContent.AllManifestDocuments(), currentSignatureDocument.Name())
	if manifest == nil {
		panic(fmt.Sprintf("Unable to merge ASiC-E with CAdES containers. "+
			"A signature with filename '%s' does not have a corresponding manifest file!", currentSignatureDocument.Name()))
	}

	result := make([]model.DSSDocument, 0)

	for _, asicContentToCompare := range asicContentList {
		signatureToCompare := spi.DSSUtilsDocumentWithName(asicContentToCompare.SignatureDocuments(), currentSignatureDocument.Name())
		if signatureToCompare != nil {
			manifestToCompare := asic.ASiCManifestParserGetLinkedManifest(asicContentToCompare.AllManifestDocuments(), signatureToCompare.Name())
			if manifestToCompare == nil {
				panic(fmt.Sprintf("Unable to merge ASiC-E with CAdES containers. "+
					"A signature with filename '%s' does not have a corresponding manifest file!", signatureToCompare.Name()))

			} else if asic.ASiCUtilsIsCoveredByManifest(currentASiCContent.AllManifestDocuments(), currentSignatureDocument.Name()) ||
				asic.ASiCUtilsIsCoveredByManifest(asicContentToCompare.AllManifestDocuments(), signatureToCompare.Name()) {
				panic(fmt.Sprintf("Unable to merge ASiC-E with CAdES containers. "+
					"A signature with name '%s' in a container is covered by a manifest!", currentSignatureDocument.Name()))

			} else if manifest.Name() == manifestToCompare.Name() && digestEquals(manifest, manifestToCompare) {
				result = append(result, signatureToCompare)

			} else {
				panic(fmt.Sprintf("Unable to merge ASiC-E with CAdES containers. "+
					"Signatures with filename '%s' sign different manifests!", currentSignatureDocument.Name()))
			}
		}
	}

	return result
}

func (m *ASiCEWithCAdESContainerMerger) updateMergedSignatureInContainers(mergedCmsSignature model.DSSDocument) {
	for _, asicContent := range m.AsicContents {
		if containsString(spi.DSSUtilsDocumentNames(asicContent.SignatureDocuments()), mergedCmsSignature.Name()) {
			asicContent.SetSignatureDocuments(asic.ASiCUtilsAddOrReplaceDocument(asicContent.SignatureDocuments(), mergedCmsSignature))
		}
	}
}

func (m *ASiCEWithCAdESContainerMerger) ensureManifestDocumentsValid() {
	mergedASiCContent := m.createEmptyContainer()
	for _, asicContent := range m.AsicContents {
		mergedASiCContent.SetManifestDocuments(append(mergedASiCContent.ManifestDocuments(), asicContent.ManifestDocuments()...))
		mergedASiCContent.SetArchiveManifestDocuments(append(mergedASiCContent.ArchiveManifestDocuments(), asicContent.ArchiveManifestDocuments()...))
		mergedASiCContent.SetEvidenceRecordManifestDocuments(append(mergedASiCContent.EvidenceRecordManifestDocuments(), asicContent.EvidenceRecordManifestDocuments()...))
	}

	asicContentsToProcess := append([]*asic.ASiCContent{}, m.AsicContents...)
	for len(asicContentsToProcess) > 0 {
		asicContent := asicContentsToProcess[0]
		asicContentsToProcess = asicContentsToProcess[1:]
		m.ensureSimpleManifestDocumentsValid(mergedASiCContent, asicContentsToProcess, asicContent)
		m.ensureArchiveManifestDocumentsValid(mergedASiCContent, asicContentsToProcess, asicContent)
		m.ensureEvidenceRecordManifestDocumentsValid(mergedASiCContent, asicContentsToProcess, asicContent)
	}
}

func (m *ASiCEWithCAdESContainerMerger) ensureSimpleManifestDocumentsValid(mergedASiCContent *asic.ASiCContent, asicContentsToProcess []*asic.ASiCContent, asicContent *asic.ASiCContent) {
	for _, manifest := range asicContent.ManifestDocuments() {
		for _, currentASiCContent := range asicContentsToProcess {
			for _, currentManifest := range currentASiCContent.ManifestDocuments() {
				if manifest.Name() != "" && manifest.Name() == currentManifest.Name() {
					if digestEquals(manifest, currentManifest) {
						// continue

					} else if asic.ASiCUtilsIsCoveredByManifest(asicContent.AllManifestDocuments(), manifest.Name()) ||
						asic.ASiCUtilsIsCoveredByManifest(currentASiCContent.AllManifestDocuments(), currentManifest.Name()) {
						panic(fmt.Sprintf("Unable to merge ASiC-E with CAdES containers. "+
							"A manifest with name '%s' in a container is covered by another manifest!", currentManifest.Name()))

					} else {
						newManifestName := m.AsicFilenameFactory.ManifestFilename(mergedASiCContent)
						currentManifest.SetName(newManifestName)
					}
				}
			}
		}
	}
}

func (m *ASiCEWithCAdESContainerMerger) ensureArchiveManifestDocumentsValid(mergedASiCContent *asic.ASiCContent, asicContentsToProcess []*asic.ASiCContent, asicContent *asic.ASiCContent) {
	for _, manifest := range asicContent.ArchiveManifestDocuments() {
		for _, currentASiCContent := range asicContentsToProcess {
			for _, currentManifest := range currentASiCContent.ArchiveManifestDocuments() {
				if manifest.Name() != "" && manifest.Name() == currentManifest.Name() {
					if digestEquals(manifest, currentManifest) {
						// continue

					} else if asic.ASiCUtilsIsCoveredByManifest(asicContent.AllManifestDocuments(), manifest.Name()) ||
						asic.ASiCUtilsIsCoveredByManifest(currentASiCContent.AllManifestDocuments(), currentManifest.Name()) {
						panic(fmt.Sprintf("Unable to merge ASiC-E with CAdES containers. "+
							"A manifest with name '%s' in a container is covered by another manifest!", currentManifest.Name()))

					} else {
						newManifestName := m.AsicFilenameFactory.ArchiveManifestFilename(mergedASiCContent)
						currentManifest.SetName(newManifestName)
					}
				}
			}
		}
	}
}

func (m *ASiCEWithCAdESContainerMerger) ensureEvidenceRecordManifestDocumentsValid(mergedASiCContent *asic.ASiCContent, asicContentsToProcess []*asic.ASiCContent, asicContent *asic.ASiCContent) {
	for _, manifest := range asicContent.EvidenceRecordManifestDocuments() {
		for _, currentASiCContent := range asicContentsToProcess {
			for _, currentManifest := range currentASiCContent.EvidenceRecordManifestDocuments() {
				if manifest.Name() != "" && manifest.Name() == currentManifest.Name() {
					if digestEquals(manifest, currentManifest) {
						// continue

					} else if asic.ASiCUtilsIsCoveredByManifest(asicContent.AllManifestDocuments(), manifest.Name()) ||
						asic.ASiCUtilsIsCoveredByManifest(currentASiCContent.AllManifestDocuments(), currentManifest.Name()) {
						panic(fmt.Sprintf("Unable to merge ASiC-E with CAdES containers. "+
							"A manifest with name '%s' in a container is covered by another manifest!", currentManifest.Name()))

					} else {
						newManifestName := m.AsicFilenameFactory.EvidenceRecordManifestFilename(mergedASiCContent)
						currentManifest.SetName(newManifestName)
					}
				}
			}
		}
	}
}

func (m *ASiCEWithCAdESContainerMerger) ensureEvidenceRecordDocumentsValid() {
	mergedASiCContent := m.createEmptyContainer()
	for _, asicContent := range m.AsicContents {
		mergedASiCContent.SetEvidenceRecordDocuments(append(mergedASiCContent.EvidenceRecordDocuments(), asicContent.EvidenceRecordDocuments()...))
		mergedASiCContent.SetEvidenceRecordManifestDocuments(append(mergedASiCContent.EvidenceRecordManifestDocuments(), asicContent.EvidenceRecordManifestDocuments()...))
	}

	asicContentsToProcess := append([]*asic.ASiCContent{}, m.AsicContents...)
	for len(asicContentsToProcess) > 0 {
		asicContent := asicContentsToProcess[0]
		asicContentsToProcess = asicContentsToProcess[1:]
		for _, evidenceRecord := range asicContent.EvidenceRecordDocuments() {
			for _, currentASiCContent := range asicContentsToProcess {
				for _, currentEvidenceRecord := range currentASiCContent.EvidenceRecordDocuments() {
					if evidenceRecord.Name() != "" && evidenceRecord.Name() == currentEvidenceRecord.Name() {
						if digestEquals(evidenceRecord, currentEvidenceRecord) {
							// continue

						} else if asic.ASiCUtilsIsCoveredByManifest(asicContent.AllManifestDocuments(), evidenceRecord.Name()) ||
							asic.ASiCUtilsIsCoveredByManifest(currentASiCContent.AllManifestDocuments(), currentEvidenceRecord.Name()) {
							panic(fmt.Sprintf("Unable to merge ASiC-E with CAdES containers. "+
								"An evidence record with name '%s' in a container is covered by a manifest!", currentEvidenceRecord.Name()))

						} else {
							currentEvidenceRecordManifest := asic.ASiCManifestParserGetLinkedManifest(
								currentASiCContent.EvidenceRecordManifestDocuments(), currentEvidenceRecord.Name())
							if currentEvidenceRecordManifest == nil {
								panic(fmt.Sprintf(
									"No linked evidence record manifest for an evidence record with filename '%s' has been found!",
									currentEvidenceRecord.Name()))
							}

							evidenceRecordType := m.getEvidenceRecordType(currentEvidenceRecord.Name())
							newEvidenceRecordName := m.AsicFilenameFactory.EvidenceRecordFilename(mergedASiCContent, evidenceRecordType)
							currentEvidenceRecord.SetName(newEvidenceRecordName)

							currentEvidenceRecordManifest = m.replaceSigReferenceDocumentName(currentEvidenceRecordManifest, newEvidenceRecordName)
							currentASiCContent.SetEvidenceRecordManifestDocuments(
								asic.ASiCUtilsAddOrReplaceDocument(currentASiCContent.EvidenceRecordManifestDocuments(), currentEvidenceRecordManifest))
						}
					}
				}
			}
		}
	}
}

func (m *ASiCEWithCAdESContainerMerger) getEvidenceRecordType(evidenceRecordFilename string) enumerations.EvidenceRecordTypeEnum {
	if asic.ASiCUtilsIsXmlEvidenceRecord(evidenceRecordFilename) {
		return enumerations.EvidenceRecordTypeEnumXMLEvidenceRecord
	} else if asic.ASiCUtilsIsAsn1EvidenceRecord(evidenceRecordFilename) {
		return enumerations.EvidenceRecordTypeEnumASN1EvidenceRecord
	}
	panic(fmt.Sprintf("The evidence record with filename '%s' is not supported!", evidenceRecordFilename))
}

func (m *ASiCEWithCAdESContainerMerger) replaceSigReferenceDocumentName(evidenceRecordManifest model.DSSDocument, newEvidenceRecordName string) model.DSSDocument {
	manifestDocumentDom, err := xmlutils.DomUtilsBuildDOMFromDocument(evidenceRecordManifest)
	if err != nil {
		panic(err)
	}
	manifestRoot, err := xmlutils.XPathUtilsGetElement(manifestDocumentDom, asic.ASiCManifestPathASiCManifestPath)
	if err != nil || manifestRoot == nil {
		panic(fmt.Sprintf("Invalid structure of ASiCEvidenceRecordManifest with name '%s'.", evidenceRecordManifest.Name()))
	}
	sigReferenceElement, err := xmlutils.XPathUtilsGetElement(manifestRoot, asic.ASiCManifestPathSigReferencePath)
	if err != nil || sigReferenceElement == nil {
		panic(fmt.Sprintf("Invalid structure of ASiCEvidenceRecordManifest with name '%s'.", evidenceRecordManifest.Name()))
	}
	sigReferenceElement.SetAttr(xmldom.Name{Local: asic.ASiCManifestAttributeURI.AttributeName()}, newEvidenceRecordName)
	serializedBytes, err := xmlutils.DomUtilsSerializeNode(manifestDocumentDom)
	if err != nil {
		panic(err)
	}
	return model.NewInMemoryDocumentWithMimeType(serializedBytes, evidenceRecordManifest.Name(), evidenceRecordManifest.MimeType())
}

// digestEquals compares two documents' digests under DefaultContainerMergerDefaultDigestAlgorithm,
// porting the repeated Arrays.equals(doc.getDigestValue(DEFAULT_DIGEST_ALGORITHM), ...) idiom.
// Not a cross-file shared helper per PORTING.md - kept local to this file (used four times above).
func digestEquals(a, b model.DSSDocument) bool {
	da, errA := a.DigestValue(asic.DefaultContainerMergerDefaultDigestAlgorithm)
	db, errB := b.DigestValue(asic.DefaultContainerMergerDefaultDigestAlgorithm)
	return errA == nil && errB == nil && bytes.Equal(da, db)
}

// containsString is a tiny local helper (List.contains equivalent); not a cross-file shared
// helper per PORTING.md.
func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

// createEmptyContainer exposes the embedded DefaultContainerMerger's unexported
// createEmptyContainer via its already-exported constructor path: DefaultContainerMerger has no
// exported equivalent, so this file builds the same shape directly (NewASiCContent +
// SetContainerType(getContainerType())), mirroring createMergedResult's own use of it.
func (m *ASiCEWithCAdESContainerMerger) createEmptyContainer() *asic.ASiCContent {
	asicContent := asic.NewASiCContent()
	asicContent.SetContainerType(m.GetTargetASiCContainerType())
	for _, ac := range m.AsicContents {
		if ac.ContainerType() != "" {
			asicContent.SetContainerType(ac.ContainerType())
			break
		}
	}
	return asicContent
}
