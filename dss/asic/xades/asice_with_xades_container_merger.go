// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/merge/ASiCEWithXAdESContainerMerger.java (DSS 6.5.RC1).
//
// Package flattening: Java's eu.europa.esig.dss.asic.xades.merge lands in this same Go
// package (dss/asic/xades).
package xades

import (
	"fmt"
	"slices"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	dssxades "github.com/ryftcore/dss-go/dss/xades"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// ASiCEWithXAdESContainerMerger is used to merge ASiC-E with XAdES containers.
type ASiCEWithXAdESContainerMerger struct {
	AbstractASiCWithXAdESContainerMerger
}

var _ asic.DefaultContainerMergerOverrides = (*ASiCEWithXAdESContainerMerger)(nil)
var _ asic.ASiCContainerMerger = (*ASiCEWithXAdESContainerMerger)(nil)

// newASiCEWithXAdESContainerMerger is the empty constructor. Port of the package-private empty
// constructor.
func newASiCEWithXAdESContainerMerger() *ASiCEWithXAdESContainerMerger {
	m := &ASiCEWithXAdESContainerMerger{
		AbstractASiCWithXAdESContainerMerger: NewAbstractASiCWithXAdESContainerMergerBase(),
	}
	m.InitDefaultContainerMerger(m)
	return m
}

// NewASiCEWithXAdESContainerMerger creates an ASiC-E With XAdES container merger from provided
// container documents. Ports ASiCEWithXAdESContainerMerger(DSSDocument...).
func NewASiCEWithXAdESContainerMerger(containers ...model.DSSDocument) *ASiCEWithXAdESContainerMerger {
	m := newASiCEWithXAdESContainerMerger()
	m.InitFromDocuments(containers...)
	return m
}

// NewASiCEWithXAdESContainerMergerFromContents creates an ASiC-E With XAdES from to given
// ASiCContents. Ports ASiCEWithXAdESContainerMerger(ASiCContent...).
func NewASiCEWithXAdESContainerMergerFromContents(asicContents ...*asic.ASiCContent) *ASiCEWithXAdESContainerMerger {
	m := newASiCEWithXAdESContainerMerger()
	m.InitFromASiCContents(asicContents...)
	return m
}

// IsSupportedDocument ports the @Override protected isSupported(DSSDocument).
func (m *ASiCEWithXAdESContainerMerger) IsSupportedDocument(container model.DSSDocument) bool {
	if !m.AbstractASiCWithXAdESContainerMerger.IsSupportedDocument(container) {
		return false
	}
	isASiCS, err := asic.ASiCUtilsIsASiCSContainer(container)
	if err != nil {
		panic(err)
	}
	return !isASiCS || (m.doesNotContainSignaturesDocument(container) && m.doesNotContainEvidenceRecordsDocument(container))
}

func (m *ASiCEWithXAdESContainerMerger) doesNotContainSignaturesDocument(container model.DSSDocument) bool {
	entryNames, err := asic.ZipUtilsInstance().ExtractEntryNames(container)
	if err != nil {
		panic(err)
	}
	return !asic.ASiCUtilsFilesContainSignatures(entryNames)
}

func (m *ASiCEWithXAdESContainerMerger) doesNotContainEvidenceRecordsDocument(container model.DSSDocument) bool {
	entryNames, err := asic.ZipUtilsInstance().ExtractEntryNames(container)
	if err != nil {
		panic(err)
	}
	return !asic.ASiCUtilsFilesContainEvidenceRecords(entryNames)
}

// IsSupportedContent ports the @Override protected isSupported(ASiCContent).
func (m *ASiCEWithXAdESContainerMerger) IsSupportedContent(asicContent *asic.ASiCContent) bool {
	if !m.AbstractASiCWithXAdESContainerMerger.IsSupportedContent(asicContent) {
		return false
	}
	isASiCS, err := asic.ASiCUtilsIsASiCSContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	return !isASiCS || (m.doesNotContainSignaturesContent(asicContent) && m.doesNotContainEvidenceRecordsContent(asicContent))
}

func (m *ASiCEWithXAdESContainerMerger) doesNotContainSignaturesContent(asicContent *asic.ASiCContent) bool {
	return len(asicContent.SignatureDocuments()) == 0
}

func (m *ASiCEWithXAdESContainerMerger) doesNotContainEvidenceRecordsContent(asicContent *asic.ASiCContent) bool {
	return len(asicContent.EvidenceRecordDocuments()) == 0
}

// GetTargetASiCContainerType ports the @Override protected getTargetASiCContainerType().
func (m *ASiCEWithXAdESContainerMerger) GetTargetASiCContainerType() enumerations.ASiCContainerType {
	return enumerations.ASiCContainerTypeASiCE
}

// EnsureContainerContentAllowMerge ports the @Override protected
// ensureContainerContentAllowMerge().
func (m *ASiCEWithXAdESContainerMerger) EnsureContainerContentAllowMerge() {
	allEmpty := true
	for _, asicContent := range m.AsicContents {
		if len(asicContent.SignatureDocuments()) != 0 || len(asicContent.EvidenceRecordDocuments()) != 0 {
			allEmpty = false
			break
		}
	}
	if allEmpty {
		// no signatures and evidence records -> can merge
		return
	}

	for _, asicContent := range m.AsicContents {
		if len(asicContent.TimestampDocuments()) != 0 {
			panic("Unable to merge ASiC-E with XAdES containers. One of the containers contains a detached timestamp!")
		}
	}

	for _, asicContent := range m.AsicContents {
		m.assertEvidenceRecordDocumentNameValid(asicContent.EvidenceRecordDocuments())
	}
}

func (m *ASiCEWithXAdESContainerMerger) assertEvidenceRecordDocumentNameValid(evidenceRecordDocuments []model.DSSDocument) {
	if len(evidenceRecordDocuments) == 0 {
		return
	}
	var evidenceRecordDocumentName string
	for _, evidenceRecordDocument := range evidenceRecordDocuments {
		if asic.ASiCUtilsEvidenceRecordXML != evidenceRecordDocument.Name() && asic.ASiCUtilsEvidenceRecordERS != evidenceRecordDocument.Name() {
			panic("Unable to merge ASiC-E with XAdES containers. The evidence record document in one of the containers has invalid naming!")
		}
		if evidenceRecordDocumentName == "" {
			evidenceRecordDocumentName = evidenceRecordDocument.Name()
		} else if evidenceRecordDocumentName != evidenceRecordDocument.Name() {
			panic("Unable to merge ASiC-E with XAdES containers. The evidence record documents have conflicting names within containers!")
		}
	}
}

// EnsureSignaturesAllowMerge ports the @Override protected ensureSignaturesAllowMerge().
func (m *ASiCEWithXAdESContainerMerger) EnsureSignaturesAllowMerge() {
	containersWithSignaturesOrERsCount := 0
	for _, asicContent := range m.AsicContents {
		if len(asicContent.SignatureDocuments()) != 0 || len(asicContent.EvidenceRecordDocuments()) != 0 {
			containersWithSignaturesOrERsCount++
		}
	}
	if containersWithSignaturesOrERsCount <= 1 {
		// no signatures or evidence records in all containers except maximum one. Can merge.
		return
	}

	coveredDocumentNames := m.getCoveredDocumentNames()

	coversManifest := false
	for range m.AsicContents {
		if m.doCoverManifest(coveredDocumentNames) {
			coversManifest = true
			break
		}
	}
	if coversManifest && !m.sameSignedDocuments() {
		panic("Unable to merge ASiC-E with XAdES containers. " +
			"manifest.xml is signed or covered and the signer data does not match between containers!")
	}

	signatureNames := m.getAllSignatureDocumentNames()
	conflictingSignatureDocumentNames := m.getConflictingDocumentNames(signatureNames)
	if len(conflictingSignatureDocumentNames) != 0 {
		for _, signatureDocumentName := range conflictingSignatureDocumentNames {
			if slices.Contains(coveredDocumentNames, signatureDocumentName) && !m.isSameDocumentContent(signatureDocumentName) {
				panic("Unable to merge ASiC-E with XAdES containers. " +
					"A signature is covered by another document, while having same signature names in both containers!")
			}
		}
		m.ensureSignatureNamesDiffer()
	}

	evidenceRecordManifestDocumentNames := m.getAllEvidenceRecordManifestDocumentNames()
	conflictingEvidenceRecordManifestNames := m.getConflictingDocumentNames(evidenceRecordManifestDocumentNames)
	if len(conflictingEvidenceRecordManifestNames) != 0 {
		for _, evidenceRecordManifestName := range conflictingEvidenceRecordManifestNames {
			if slices.Contains(coveredDocumentNames, evidenceRecordManifestName) && !m.isSameDocumentContent(evidenceRecordManifestName) {
				panic("Unable to merge ASiC-E with XAdES containers. " +
					"An evidence record manifest is covered by another document, while having same signature names in both containers!")
			}
		}
		m.ensureEvidenceRecordManifestNamesDiffer()
	}

	// Create a merged manifest.xml file
	newManifest := m.createNewManifest()
	for _, asicContent := range m.AsicContents {
		asicContent.SetManifestDocuments([]model.DSSDocument{newManifest})
	}
}

func (m *ASiCEWithXAdESContainerMerger) isSameDocumentContent(documentName string) bool {
	var digestValue []byte
	var digestSet bool
	for _, asicContent := range m.AsicContents {
		document := spi.DSSUtilsDocumentWithName(asicContent.AllDocuments(), documentName)
		currentDigestValue, err := document.DigestValue(asic.DefaultContainerMergerDefaultDigestAlgorithm)
		if err != nil {
			panic(err)
		}
		if !digestSet {
			digestValue = currentDigestValue
			digestSet = true
		} else if !bytesEqualXades(digestValue, currentDigestValue) {
			return false
		}
	}
	return true
}

func (m *ASiCEWithXAdESContainerMerger) getCoveredDocumentNames() []string {
	result := make([]string, 0)
	for _, asicContent := range m.AsicContents {
		for _, signatureDocument := range asicContent.SignatureDocuments() {
			documentAnalyzer, err := dssxades.NewXMLDocumentAnalyzer(signatureDocument)
			if err != nil {
				panic(err)
			}
			for _, signature := range documentAnalyzer.Signatures() {
				xadesSignature, ok := signature.(*dssxades.XAdESSignature)
				if !ok {
					continue
				}
				result = append(result, m.getCoveredDocumentNamesFromSignature(xadesSignature)...)
			}
		}
		for _, manifestDocument := range asicContent.EvidenceRecordManifestDocuments() {
			manifestFile := asic.ASiCManifestParserGetManifestFile(manifestDocument)
			if manifestFile != nil {
				for _, manifestEntry := range manifestFile.Entries() {
					result = append(result, manifestEntry.Uri())
				}
			}
		}
	}
	return result
}

func (m *ASiCEWithXAdESContainerMerger) getCoveredDocumentNamesFromSignature(signature *dssxades.XAdESSignature) []string {
	result := make([]string, 0)
	for _, reference := range signature.References() {
		referenceURI := dssxades.DSSXMLUtilsGetReferenceURI(reference)
		if !xmlutils.DomUtilsStartsFromHash(referenceURI) && !xmlutils.DomUtilsIsXPointerQuery(referenceURI) {
			result = append(result, referenceURI)
		}
	}
	return result
}

func (m *ASiCEWithXAdESContainerMerger) doCoverManifest(documentNames []string) bool {
	return slices.Contains(documentNames, asic.ASiCUtilsASiCEMetaInfManifest)
}

func (m *ASiCEWithXAdESContainerMerger) sameSignedDocuments() bool {
	var signedDocumentNames map[string]struct{}
	var set bool
	for _, asicContent := range m.AsicContents {
		currentSignedDocumentNames := make(map[string]struct{})
		for _, name := range spi.DSSUtilsDocumentNames(asicContent.SignedDocuments()) {
			currentSignedDocumentNames[name] = struct{}{}
		}
		if !set {
			signedDocumentNames = currentSignedDocumentNames
			set = true
		} else if !stringSetsEqualXades(signedDocumentNames, currentSignedDocumentNames) {
			return false
		}
	}
	return true
}

func (m *ASiCEWithXAdESContainerMerger) getAllSignatureDocumentNames() []string {
	signatureDocumentNames := make([]string, 0)
	for _, asicContent := range m.AsicContents {
		signatureDocumentNames = append(signatureDocumentNames, spi.DSSUtilsDocumentNames(asicContent.SignatureDocuments())...)
	}
	return signatureDocumentNames
}

func (m *ASiCEWithXAdESContainerMerger) getAllEvidenceRecordManifestDocumentNames() []string {
	erManifestDocumentNames := make([]string, 0)
	for _, asicContent := range m.AsicContents {
		erManifestDocumentNames = append(erManifestDocumentNames, spi.DSSUtilsDocumentNames(asicContent.EvidenceRecordManifestDocuments())...)
	}
	return erManifestDocumentNames
}

func (m *ASiCEWithXAdESContainerMerger) getConflictingDocumentNames(documentNames []string) []string {
	result := make([]string, 0)
	for _, name := range documentNames {
		count := 0
		for _, other := range documentNames {
			if other == name {
				count++
			}
		}
		if count > 1 {
			result = append(result, name)
		}
	}
	return result
}

func (m *ASiCEWithXAdESContainerMerger) createNewManifest() model.DSSDocument {
	manifestEntries := make([]*model.ManifestEntry, 0)
	addedFileNames := make([]string, 0)

	mergedContent := m.createEmptyContainer()
	for _, asicContent := range m.AsicContents {
		manifestDocuments := asicContent.ManifestDocuments()
		mergedContent.SetManifestDocuments(append(mergedContent.ManifestDocuments(), manifestDocuments...))

		for _, entry := range m.getManifestFileEntries(manifestDocuments) {
			if !slices.Contains(addedFileNames, entry.Uri()) {
				manifestEntries = append(manifestEntries, entry)
				addedFileNames = append(addedFileNames, entry.Uri())
			}
		}
		signedDocuments := asicContent.SignedDocuments()
		for _, entry := range asic.ASiCUtilsToSimpleManifestEntries(signedDocuments) {
			if !slices.Contains(addedFileNames, entry.Uri()) {
				manifestEntries = append(manifestEntries, entry)
				addedFileNames = append(addedFileNames, entry.Uri())
			}
		}
	}

	return m.createNewManifestXml(manifestEntries, mergedContent)
}

func (m *ASiCEWithXAdESContainerMerger) getManifestFileEntries(manifestDocuments []model.DSSDocument) []*model.ManifestEntry {
	if len(manifestDocuments) == 0 {
		return []*model.ManifestEntry{}
	} else if len(manifestDocuments) > 1 {
		panic(exception.NewIllegalInputException("One of the containers contain multiple manifest files!"))
	}

	manifestDocument := manifestDocuments[0]
	if asic.ASiCUtilsASiCEMetaInfManifest != manifestDocument.Name() {
		panic(exception.NewIllegalInputException(fmt.Sprintf("A manifest file shall have a name '%s'.", asic.ASiCUtilsASiCEMetaInfManifest)))
	}

	parser := NewASiCEWithXAdESManifestParser(manifestDocument)
	manifest := parser.Manifest()
	return manifest.Entries()
}

func (m *ASiCEWithXAdESContainerMerger) createNewManifestXml(manifestEntries []*model.ManifestEntry, asicContent *asic.ASiCContent) model.DSSDocument {
	document, err := NewASiCEWithXAdESManifestBuilder().
		SetEntries(manifestEntries).
		SetManifestFilename(m.asicFilenameFactory.ManifestFilename(asicContent)).
		Build()
	if err != nil {
		panic(err)
	}
	return document
}

func (m *ASiCEWithXAdESContainerMerger) ensureSignatureNamesDiffer() {
	usedSignatureNames := make(map[string]struct{})
	mergedASiCContent := m.createEmptyContainer()
	for _, asicContent := range m.AsicContents {
		mergedASiCContent.SetSignatureDocuments(append(mergedASiCContent.SignatureDocuments(), asicContent.SignatureDocuments()...))
	}

	for _, asicContent := range m.AsicContents {
		for _, signatureDocument := range asicContent.SignatureDocuments() {
			if _, used := usedSignatureNames[signatureDocument.Name()]; used {
				newSignatureName := m.asicFilenameFactory.SignatureFilename(mergedASiCContent)
				signatureDocument.SetName(newSignatureName)
			}
			usedSignatureNames[signatureDocument.Name()] = struct{}{}
		}
	}
}

func (m *ASiCEWithXAdESContainerMerger) ensureEvidenceRecordManifestNamesDiffer() {
	usedNames := make(map[string]struct{})
	mergedASiCContent := m.createEmptyContainer()
	for _, asicContent := range m.AsicContents {
		mergedASiCContent.SetEvidenceRecordManifestDocuments(append(mergedASiCContent.EvidenceRecordManifestDocuments(), asicContent.EvidenceRecordManifestDocuments()...))
	}

	for _, asicContent := range m.AsicContents {
		for _, evidenceRecordManifest := range asicContent.EvidenceRecordManifestDocuments() {
			if _, used := usedNames[evidenceRecordManifest.Name()]; used {
				newName := m.asicFilenameFactory.EvidenceRecordManifestFilename(mergedASiCContent)
				evidenceRecordManifest.SetName(newName)
			}
			usedNames[evidenceRecordManifest.Name()] = struct{}{}
		}
	}
}

// createEmptyContainer exposes the embedded DefaultContainerMerger's unexported
// createEmptyContainer via its already-exported constructor path: DefaultContainerMerger has no
// exported equivalent, so this file builds the same shape directly (NewASiCContent +
// SetContainerType(getContainerType())), mirroring the asic/cades precedent.
func (m *ASiCEWithXAdESContainerMerger) createEmptyContainer() *asic.ASiCContent {
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

// stringSetsEqualXades ports the repeated Set<String>.equals(Set<String>) idiom used above.
func stringSetsEqualXades(a, b map[string]struct{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if _, ok := b[k]; !ok {
			return false
		}
	}
	return true
}

// bytesEqualXades ports Arrays.equals(byte[], byte[]) for the isSameDocumentContent digest
// comparison above.
func bytesEqualXades(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
