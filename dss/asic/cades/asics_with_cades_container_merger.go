// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/merge/ASiCSWithCAdESContainerMerger.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ASiCSWithCAdESContainerMerger is used to merge ASiC-S with CAdES containers.
type ASiCSWithCAdESContainerMerger struct {
	AbstractASiCWithCAdESContainerMerger
}

var _ asic.DefaultContainerMergerOverrides = (*ASiCSWithCAdESContainerMerger)(nil)
var _ asic.ContainerMerger = (*ASiCSWithCAdESContainerMerger)(nil)

// newASiCSWithCAdESContainerMerger is the empty constructor. Port of the package-private empty
// constructor.
func newASiCSWithCAdESContainerMerger() *ASiCSWithCAdESContainerMerger {
	m := &ASiCSWithCAdESContainerMerger{
		AbstractASiCWithCAdESContainerMerger: NewAbstractASiCWithCAdESContainerMergerBase(),
	}
	m.InitDefaultContainerMerger(m)
	return m
}

// NewASiCSWithCAdESContainerMerger creates an ASiC-S with CAdES container merger from provided
// container documents. Ports ASiCSWithCAdESContainerMerger(DSSDocument...).
func NewASiCSWithCAdESContainerMerger(containers ...model.DSSDocument) *ASiCSWithCAdESContainerMerger {
	m := newASiCSWithCAdESContainerMerger()
	m.InitFromDocuments(containers...)
	return m
}

// NewASiCSWithCAdESContainerMergerFromContents creates an ASiC-S with CAdES container merger
// from the given ASiCContents. Ports ASiCSWithCAdESContainerMerger(ASiCContent...).
func NewASiCSWithCAdESContainerMergerFromContents(asicContents ...*asic.Content) *ASiCSWithCAdESContainerMerger {
	m := newASiCSWithCAdESContainerMerger()
	m.InitFromASiCContents(asicContents...)
	return m
}

// IsSupportedDocument ports the @Override public isSupported(DSSDocument).
func (m *ASiCSWithCAdESContainerMerger) IsSupportedDocument(container model.DSSDocument) bool {
	if !m.AbstractASiCWithCAdESContainerMerger.IsSupportedDocument(container) {
		return false
	}
	isASiCE, err := asic.UtilsIsASiCEContainer(container)
	if err != nil {
		panic(err)
	}
	return !isASiCE
}

// IsSupportedContent ports the @Override public isSupported(ASiCContent).
func (m *ASiCSWithCAdESContainerMerger) IsSupportedContent(asicContent *asic.Content) bool {
	if !m.AbstractASiCWithCAdESContainerMerger.IsSupportedContent(asicContent) {
		return false
	}
	isASiCE, err := asic.UtilsIsASiCEContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	return !isASiCE
}

// GetTargetASiCContainerType ports the @Override protected getTargetASiCContainerType().
func (m *ASiCSWithCAdESContainerMerger) GetTargetASiCContainerType() enumerations.ASiCContainerType {
	return enumerations.ASiCContainerTypeASiCS
}

// EnsureContainerContentAllowMerge ports the @Override protected
// ensureContainerContentAllowMerge().
func (m *ASiCSWithCAdESContainerMerger) EnsureContainerContentAllowMerge() {
	allEmpty := true
	for _, asicContent := range m.AsicContents {
		if utils.IsCollectionNotEmpty(asicContent.SignatureDocuments()) ||
			utils.IsCollectionNotEmpty(asicContent.TimestampDocuments()) ||
			utils.IsCollectionNotEmpty(asicContent.EvidenceRecordDocuments()) {
			allEmpty = false
			break
		}
	}
	if allEmpty {
		return // no signatures, timestamps and evidence records -> can merge
	}

	for _, asicContent := range m.AsicContents {
		if utils.CollectionSize(asicContent.SignatureDocuments())+utils.CollectionSize(asicContent.TimestampDocuments())+
			utils.CollectionSize(asicContent.EvidenceRecordDocuments()) > 1 {
			panic("Unable to merge ASiC-S with CAdES containers. " +
				"One of the containers has more than one signature, timestamp or evidence record documents!")
		}
	}

	hasSignatures := false
	hasTimestampsOrEvidenceRecords := false
	for _, asicContent := range m.AsicContents {
		if utils.IsCollectionNotEmpty(asicContent.SignatureDocuments()) {
			hasSignatures = true
		}
		if utils.IsCollectionNotEmpty(asicContent.TimestampDocuments()) || utils.IsCollectionNotEmpty(asicContent.EvidenceRecordDocuments()) {
			hasTimestampsOrEvidenceRecords = true
		}
	}
	if hasSignatures && hasTimestampsOrEvidenceRecords {
		panic("Unable to merge ASiC-S with CAdES containers. " +
			"Only one type of a container is allowed (signature, timestamp or evidence record)!")
	}

	timestampOrEvidenceRecordContainers := 0
	for _, asicContent := range m.AsicContents {
		if utils.IsCollectionNotEmpty(asicContent.TimestampDocuments()) || utils.IsCollectionNotEmpty(asicContent.EvidenceRecordDocuments()) {
			timestampOrEvidenceRecordContainers++
		}
	}
	if timestampOrEvidenceRecordContainers > 1 {
		panic("Unable to merge ASiC-S with CAdES containers. " +
			"Timestamp or evidence record containers cannot be merged with the given container type!")
	}

	for _, asicContent := range m.AsicContents {
		m.assertSignatureDocumentNameValid(asicContent.SignatureDocuments())
	}
	for _, asicContent := range m.AsicContents {
		m.assertTimestampDocumentNameValid(asicContent.TimestampDocuments())
	}
	for _, asicContent := range m.AsicContents {
		m.assertEvidenceRecordDocumentNameValid(asicContent.EvidenceRecordDocuments())
	}

	for _, asicContent := range m.AsicContents {
		if utils.CollectionSize(asicContent.RootLevelSignedDocuments()) > 1 {
			panic("Unable to merge ASiC-S with CAdES containers. " +
				"One of the containers has more than one signer documents!")
		}
	}

	if !m.checkRootSignerDocumentsNames() {
		panic("Unable to merge ASiC-S with CAdES containers. " +
			"Signer documents have different names!")
	}
}

func (m *ASiCSWithCAdESContainerMerger) assertSignatureDocumentNameValid(signatureDocuments []model.DSSDocument) {
	if utils.IsCollectionNotEmpty(signatureDocuments) {
		for _, signatureDocument := range signatureDocuments {
			if asic.ASiCUtilsSignatureP7S != signatureDocument.Name() {
				panic("Unable to merge ASiC-S with CAdES containers. " +
					"The signature document in one of the containers has invalid naming!")
			}
		}
	}
}

func (m *ASiCSWithCAdESContainerMerger) assertTimestampDocumentNameValid(timestampDocuments []model.DSSDocument) {
	if utils.IsCollectionNotEmpty(timestampDocuments) {
		for _, tstDocument := range timestampDocuments {
			if asic.ASiCUtilsTimestampTST != tstDocument.Name() {
				panic("Unable to merge ASiC-S with CAdES containers. " +
					"The timestamp document in one of the containers has invalid naming!")
			}
		}
	}
}

func (m *ASiCSWithCAdESContainerMerger) assertEvidenceRecordDocumentNameValid(evidenceRecordDocuments []model.DSSDocument) {
	if utils.IsCollectionNotEmpty(evidenceRecordDocuments) {
		evidenceRecordDocumentName := ""
		for _, evidenceRecordDocument := range evidenceRecordDocuments {
			if asic.ASiCUtilsEvidenceRecordXML != evidenceRecordDocument.Name() &&
				asic.ASiCUtilsEvidenceRecordERS != evidenceRecordDocument.Name() {
				panic("Unable to merge ASiC-S with CAdES containers. " +
					"The evidence record document in one of the containers has invalid naming!")
			}
			if evidenceRecordDocumentName == "" {
				evidenceRecordDocumentName = evidenceRecordDocument.Name()
			} else if evidenceRecordDocumentName != evidenceRecordDocument.Name() {
				panic("Unable to merge ASiC-S with CAdES containers. " +
					"The evidence record documents have conflicting names within containers!")
			}
		}
	}
}

func (m *ASiCSWithCAdESContainerMerger) checkRootSignerDocumentsNames() bool {
	rootSignedDocumentName := ""
	for _, asicContent := range m.AsicContents {
		rootLevelSignedDocuments := asicContent.RootLevelSignedDocuments()
		if utils.IsCollectionNotEmpty(rootLevelSignedDocuments) {
			currentSignedDocument := rootLevelSignedDocuments[0] // only one shall be present
			if rootSignedDocumentName == "" {
				rootSignedDocumentName = currentSignedDocument.Name()
			} else {
				return rootSignedDocumentName == currentSignedDocument.Name()
			}
		}
	}
	return true
}

// EnsureSignaturesAllowMerge ports the @Override protected ensureSignaturesAllowMerge().
func (m *ASiCSWithCAdESContainerMerger) EnsureSignaturesAllowMerge() {
	containersWithSignaturesCount := 0
	for _, asicContent := range m.AsicContents {
		if utils.IsCollectionNotEmpty(asicContent.SignatureDocuments()) ||
			utils.IsCollectionNotEmpty(asicContent.TimestampDocuments()) ||
			utils.IsCollectionNotEmpty(asicContent.EvidenceRecordDocuments()) {
			containersWithSignaturesCount++
		}
	}
	if containersWithSignaturesCount <= 1 {
		// one of the containers does not contain a signature document. Can merge.
		return
	}

	m.mergeSignatureDocuments()
}

func (m *ASiCSWithCAdESContainerMerger) mergeSignatureDocuments() {
	allSignatureDocuments := m.GetAllSignatureDocuments(m.AsicContents...)
	mergedCMSSignaturesDocument := m.MergeCmsSignatures(allSignatureDocuments)
	for _, asicContent := range m.AsicContents {
		asicContent.SetSignatureDocuments([]model.DSSDocument{mergedCMSSignaturesDocument})
	}
}
