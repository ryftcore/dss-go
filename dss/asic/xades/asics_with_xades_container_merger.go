// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/merge/ASiCSWithXAdESContainerMerger.java (DSS 6.5.RC1).
//
// Package flattening: Java's eu.europa.esig.dss.asic.xades.merge lands in this same Go
// package (dss/asic/xades).
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	dssxades "github.com/ryftcore/dss-go/dss/xades"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// ASiCSWithXAdESContainerMerger is used to merge ASiC-S with XAdES containers.
type ASiCSWithXAdESContainerMerger struct {
	AbstractASiCWithXAdESContainerMerger
}

var _ asic.DefaultContainerMergerOverrides = (*ASiCSWithXAdESContainerMerger)(nil)
var _ asic.ContainerMerger = (*ASiCSWithXAdESContainerMerger)(nil)

// newASiCSWithXAdESContainerMerger is the empty constructor. Port of the package-private empty
// constructor.
func newASiCSWithXAdESContainerMerger() *ASiCSWithXAdESContainerMerger {
	m := &ASiCSWithXAdESContainerMerger{
		AbstractASiCWithXAdESContainerMerger: NewAbstractASiCWithXAdESContainerMergerBase(),
	}
	m.InitDefaultContainerMerger(m)
	return m
}

// NewASiCSWithXAdESContainerMerger creates an ASiC-S With XAdES container merger from provided
// container documents. Ports ASiCSWithXAdESContainerMerger(DSSDocument...).
func NewASiCSWithXAdESContainerMerger(containers ...model.DSSDocument) *ASiCSWithXAdESContainerMerger {
	m := newASiCSWithXAdESContainerMerger()
	m.InitFromDocuments(containers...)
	return m
}

// NewASiCSWithXAdESContainerMergerFromContents creates an ASiC-S With XAdES from to given
// ASiCContents. Ports ASiCSWithXAdESContainerMerger(ASiCContent...).
func NewASiCSWithXAdESContainerMergerFromContents(asicContents ...*asic.Content) *ASiCSWithXAdESContainerMerger {
	m := newASiCSWithXAdESContainerMerger()
	m.InitFromASiCContents(asicContents...)
	return m
}

// IsSupportedDocument ports the @Override public isSupported(DSSDocument).
func (m *ASiCSWithXAdESContainerMerger) IsSupportedDocument(container model.DSSDocument) bool {
	if !m.AbstractASiCWithXAdESContainerMerger.IsSupportedDocument(container) {
		return false
	}
	isASiCE, err := asic.UtilsIsASiCEContainer(container)
	if err != nil {
		panic(err)
	}
	return !isASiCE
}

// IsSupportedContent ports the @Override public isSupported(ASiCContent).
func (m *ASiCSWithXAdESContainerMerger) IsSupportedContent(asicContent *asic.Content) bool {
	if !m.AbstractASiCWithXAdESContainerMerger.IsSupportedContent(asicContent) {
		return false
	}
	isASiCE, err := asic.UtilsIsASiCEContainerContent(asicContent)
	if err != nil {
		panic(err)
	}
	return !isASiCE
}

// GetTargetASiCContainerType ports the @Override protected getTargetASiCContainerType().
func (m *ASiCSWithXAdESContainerMerger) GetTargetASiCContainerType() enumerations.ASiCContainerType {
	return enumerations.ASiCContainerTypeASiCS
}

// EnsureContainerContentAllowMerge ports the @Override protected
// ensureContainerContentAllowMerge().
func (m *ASiCSWithXAdESContainerMerger) EnsureContainerContentAllowMerge() {
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
		if len(asicContent.SignatureDocuments()) > 1 {
			panic("Unable to merge ASiC-S with XAdES containers. One of the containers has more than one signature documents!")
		}
	}
	for _, asicContent := range m.AsicContents {
		if len(asicContent.TimestampDocuments()) != 0 {
			panic("Unable to merge ASiC-S with XAdES containers. One of the containers contains a detached timestamp!")
		}
	}
	for _, asicContent := range m.AsicContents {
		if len(asicContent.EvidenceRecordDocuments()) > 1 {
			panic("Unable to merge ASiC-S with XAdES containers. One of the containers has more than one evidence record documents!")
		}
	}
	containersWithERCount := 0
	for _, asicContent := range m.AsicContents {
		if len(asicContent.EvidenceRecordDocuments()) != 0 {
			containersWithERCount++
		}
	}
	if containersWithERCount > 1 {
		panic("Unable to merge ASiC-S with XAdES containers. Evidence record containers cannot be merged with the given container type!")
	}
	for _, asicContent := range m.AsicContents {
		if len(asicContent.RootLevelSignedDocuments()) > 1 {
			panic("Unable to merge ASiC-S with XAdES containers. One of the containers has more than one signer documents!")
		}
	}
	if len(m.getSignerDocumentNameSet()) > 1 {
		panic("Unable to merge ASiC-S with XAdES containers. Signer documents have different names!")
	}
	hasSignatures := false
	hasEvidenceRecords := false
	for _, asicContent := range m.AsicContents {
		if len(asicContent.SignatureDocuments()) != 0 {
			hasSignatures = true
		}
		if len(asicContent.EvidenceRecordDocuments()) != 0 {
			hasEvidenceRecords = true
		}
	}
	if hasSignatures && hasEvidenceRecords {
		panic("Unable to merge ASiC-S with XAdES containers. Only one type of a container is allowed (signature or evidence record)!")
	}

	for _, asicContent := range m.AsicContents {
		m.assertSignatureDocumentNameValid(asicContent.SignatureDocuments())
	}
	for _, asicContent := range m.AsicContents {
		m.assertEvidenceRecordDocumentNameValid(asicContent.EvidenceRecordDocuments())
	}
}

func (m *ASiCSWithXAdESContainerMerger) assertSignatureDocumentNameValid(signatureDocuments []model.DSSDocument) {
	for _, signatureDocument := range signatureDocuments {
		if asic.ASiCUtilsSignaturesXML != signatureDocument.Name() {
			panic("Unable to merge ASiC-S with XAdES containers. The signature document in one of the containers has invalid naming!")
		}
	}
}

func (m *ASiCSWithXAdESContainerMerger) assertEvidenceRecordDocumentNameValid(evidenceRecordDocuments []model.DSSDocument) {
	if len(evidenceRecordDocuments) == 0 {
		return
	}
	var evidenceRecordDocumentName string
	for _, evidenceRecordDocument := range evidenceRecordDocuments {
		if asic.ASiCUtilsEvidenceRecordXML != evidenceRecordDocument.Name() && asic.ASiCUtilsEvidenceRecordERS != evidenceRecordDocument.Name() {
			panic("Unable to merge ASiC-S with XAdES containers. The evidence record document in one of the containers has invalid naming!")
		}
		if evidenceRecordDocumentName == "" {
			evidenceRecordDocumentName = evidenceRecordDocument.Name()
		} else if evidenceRecordDocumentName != evidenceRecordDocument.Name() {
			panic("Unable to merge ASiC-S with XAdES containers. The evidence record documents have conflicting names within containers!")
		}
	}
}

func (m *ASiCSWithXAdESContainerMerger) getSignerDocumentNameSet() map[string]struct{} {
	result := make(map[string]struct{})
	for _, asicContent := range m.AsicContents {
		for _, name := range spi.DSSUtilsDocumentNames(asicContent.RootLevelSignedDocuments()) {
			result[name] = struct{}{}
		}
	}
	return result
}

// EnsureSignaturesAllowMerge ports the @Override protected ensureSignaturesAllowMerge().
func (m *ASiCSWithXAdESContainerMerger) EnsureSignaturesAllowMerge() {
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

	m.mergeSignatureDocuments()
}

func (m *ASiCSWithXAdESContainerMerger) mergeSignatureDocuments() {
	documentAnalyzers := m.getAllDocumentAnalyzers()
	allSignatures := m.getAllSignatures(documentAnalyzers)
	if len(allSignatures) == 0 {
		return
	}

	if !m.checkNoCommonIdsBetweenSignatures(allSignatures) {
		panic(exception.NewIllegalInputException("Signature documents contain signatures with the same identifiers!"))
	}
	if !m.checkNoCommonIdsBetweenSignedData(allSignatures) {
		panic(exception.NewIllegalInputException("Signature documents contain signatures signed enveloped objects with the same identifiers!"))
	}
	if !m.checkNoCommonIdsBetweenSignatureValues(allSignatures) {
		panic(exception.NewIllegalInputException("Signature documents contain signatures with SignatureValue elements sharing the same ids!"))
	}
	m.assertSameRootElement(documentAnalyzers)

	signaturesXml := m.getMergedSignaturesXml(documentAnalyzers)
	for _, asicContent := range m.AsicContents {
		asicContent.SetSignatureDocuments([]model.DSSDocument{signaturesXml})
	}
}

func (m *ASiCSWithXAdESContainerMerger) getAllDocumentAnalyzers() []*dssxades.XMLDocumentAnalyzer {
	analyzers := make([]*dssxades.XMLDocumentAnalyzer, 0)
	for _, asicContent := range m.AsicContents {
		for _, signatureDocument := range asicContent.SignatureDocuments() {
			analyzer, err := dssxades.NewXMLDocumentAnalyzer(signatureDocument)
			if err != nil {
				panic(err)
			}
			analyzers = append(analyzers, analyzer)
		}
	}
	return analyzers
}

func (m *ASiCSWithXAdESContainerMerger) getAllSignatures(analyzers []*dssxades.XMLDocumentAnalyzer) []*dssxades.Signature {
	signatures := make([]*dssxades.Signature, 0)
	for _, analyzer := range analyzers {
		for _, signature := range analyzer.Signatures() {
			if xadesSignature, ok := signature.(*dssxades.Signature); ok {
				signatures = append(signatures, xadesSignature)
			}
		}
	}
	return signatures
}

func (m *ASiCSWithXAdESContainerMerger) checkNoCommonIdsBetweenSignatures(signatures []*dssxades.Signature) bool {
	signatureIds := m.getSignatureIds(signatures)
	return !m.checkDuplicatesPresent(signatureIds)
}

func (m *ASiCSWithXAdESContainerMerger) getSignatureIds(signatures []*dssxades.Signature) []string {
	ids := make([]string, 0, len(signatures))
	for _, signature := range signatures {
		ids = append(ids, signature.DAIdentifier())
	}
	return ids
}

func (m *ASiCSWithXAdESContainerMerger) checkNoCommonIdsBetweenSignedData(signatures []*dssxades.Signature) bool {
	signedDataObjectIdsOne := m.getSignedDataObjectIds(signatures)
	return !m.checkDuplicatesPresent(signedDataObjectIdsOne)
}

func (m *ASiCSWithXAdESContainerMerger) getSignedDataObjectIds(signatures []*dssxades.Signature) []string {
	ids := make([]string, 0)
	for _, xadesSignature := range signatures {
		for _, reference := range xadesSignature.References() {
			// Java's outer null-check (referenceURI != null) and inner empty-check
			// (Utils.EMPTY_STRING.equals(referenceURI)) both collapse onto Go's "" sentinel for
			// a Java null string (per PORTING.md), so the inner branch is unreachable here -
			// ported verbatim rather than simplified away, matching upstream's structure.
			referenceURI := dssxades.DSSXMLUtilsGetReferenceURI(reference)
			if referenceURI != "" {
				if referenceURI == "" {
					panic(exception.NewIllegalInputException(
						"Unable to merge signatures, as one of them covers the whole signature file document!"))
				}
				if xmlutils.DomUtilsStartsFromHash(referenceURI) || xmlutils.DomUtilsIsXPointerQuery(referenceURI) {
					// identifiers referencing objects within the document should be analyzed
					ids = append(ids, referenceURI)
				}
			}
		}
	}
	return ids
}

func (m *ASiCSWithXAdESContainerMerger) checkNoCommonIdsBetweenSignatureValues(signatures []*dssxades.Signature) bool {
	signatureValueIds := m.getSignatureValueIds(signatures)
	return !m.checkDuplicatesPresent(signatureValueIds)
}

func (m *ASiCSWithXAdESContainerMerger) getSignatureValueIds(signatures []*dssxades.Signature) []string {
	ids := make([]string, 0, len(signatures))
	for _, xadesSignature := range signatures {
		ids = append(ids, xadesSignature.SignatureValueId())
	}
	return ids
}

func (m *ASiCSWithXAdESContainerMerger) checkDuplicatesPresent(strs []string) bool {
	for _, s := range strs {
		count := 0
		for _, other := range strs {
			if other == s {
				count++
			}
		}
		if count > 1 {
			return true
		}
	}
	return false
}

func (m *ASiCSWithXAdESContainerMerger) assertSameRootElement(documentAnalyzers []*dssxades.XMLDocumentAnalyzer) {
	var rootElement *xmldom.Node
	for _, documentAnalyzer := range documentAnalyzers {
		currentRootElement := documentAnalyzer.RootElement().DocumentElement()
		if rootElement == nil {
			rootElement = currentRootElement
			continue
		}
		if rootElement.Name.Local != currentRootElement.Name.Local {
			panic(exception.NewIllegalInputException("Signature containers have different root elements!"))
		}
		if (rootElement.Name.Space != "") != (currentRootElement.Name.Space != "") {
			panic(exception.NewIllegalInputException("Signature containers have different namespaces!"))
		}
		if rootElement.Name.Space != "" && rootElement.Name.Space != currentRootElement.Name.Space {
			panic(exception.NewIllegalInputException("Signature containers have different namespaces!"))
		}
		if rootElement.Name.Prefix != currentRootElement.Name.Prefix {
			panic(exception.NewIllegalInputException("Signature containers have different namespace prefixes!"))
		}
	}
}

func (m *ASiCSWithXAdESContainerMerger) getMergedSignaturesXml(documentAnalyzers []*dssxades.XMLDocumentAnalyzer) model.DSSDocument {
	var documentElement *xmldom.Node

	for _, documentAnalyzer := range documentAnalyzers {
		if documentElement == nil {
			documentElement = documentAnalyzer.RootElement().DocumentElement()
		} else {
			toBeAdopted := documentAnalyzer.RootElement().DocumentElement()
			xmlutils.DomUtilsAdoptChildren(documentElement, toBeAdopted)
		}
	}

	bytes, err := xmlutils.DomUtilsSerializeNode(documentElement)
	if err != nil {
		panic(err)
	}
	return model.NewInMemoryDocumentWithMimeType(bytes, asic.ASiCUtilsSignaturesXML, enumerations.MimeTypeEnumXML)
}

// compile-time assertion that the AdvancedSignature interface stays wired to the validation
// package, matching this package's other files that type-assert *dssxades.Signature out of
// analyzer.Signatures().
var _ validation.AdvancedSignature = (*dssxades.Signature)(nil)
