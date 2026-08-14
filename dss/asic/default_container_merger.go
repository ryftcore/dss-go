// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/merge/DefaultContainerMerger.java (DSS 6.5.RC1).
package asic

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/exception"
	"github.com/utain/esig/dss/utils"
)

// DefaultContainerMergerDefaultDigestAlgorithm is the digest algorithm used for internal
// documents comparison. Port of the protected static final DEFAULT_DIGEST_ALGORITHM.
const DefaultContainerMergerDefaultDigestAlgorithm = enumerations.DigestAlgorithm_SHA256

// DefaultContainerMergerOverrides declares the operations DefaultContainerMerger calls back
// into virtually - the base's Merge()/MergeToASiCContent() dispatch to these, the way
// model.TokenBase dispatches to model.TokenOverrides via InitToken. Per S7_BRIEF.md's
// virtual-dispatch warning, every concrete merger (ASiCWithCAdESContainerMerger,
// ASiCWithXAdESContainerMerger, ASiCSContainerMerger; CADSIGN/XADSIGN chunks) must call
// InitDefaultContainerMerger with itself before use.
type DefaultContainerMergerOverrides interface {
	// GetContainerExtractor returns a relevant ASiC container extractor. Port of the protected
	// abstract getContainerExtractor(DSSDocument).
	//
	// Cross-chunk assumption (ZIPCORE): DefaultASiCContainerExtractor exposes an Extract()
	// (*ASiCContent, error) method.
	GetContainerExtractor(container model.DSSDocument) *DefaultASiCContainerExtractor

	// IsSupportedDocument verifies whether the provided container is supported by the current
	// class. Port of the protected abstract isSupported(DSSDocument).
	IsSupportedDocument(container model.DSSDocument) bool

	// IsSupportedContent verifies whether the provided ASiCContent is supported by the current
	// class. Port of the protected abstract isSupported(ASiCContent).
	IsSupportedContent(asicContent *ASiCContent) bool

	// EnsureContainerContentAllowMerge verifies whether containers can be merged. Port of the
	// protected abstract ensureContainerContentAllowMerge().
	EnsureContainerContentAllowMerge()

	// EnsureSignaturesAllowMerge ensures that the entry names between the containers' entries
	// are different. Port of the protected abstract ensureSignaturesAllowMerge().
	EnsureSignaturesAllowMerge()

	// GetTargetASiCContainerType returns a target ASiC Container Type of the current merger
	// class. Port of the protected abstract getTargetASiCContainerType().
	GetTargetASiCContainerType() enumerations.ASiCContainerType
}

// DefaultContainerMerger loads a relevant ASiCContainerMerger in order to merge content of
// given containers. Ports the abstract class implementing ASiCContainerMerger.
type DefaultContainerMerger struct {
	// overrides points back at the concrete merger; see InitDefaultContainerMerger.
	overrides DefaultContainerMergerOverrides

	// AsicContents is an array of ASiC contents representing containers to be merged. Java
	// declares the field protected; exported here since Go subclasses in the cades/xades ASiC
	// packages live in different packages.
	AsicContents []*ASiCContent

	// creationTime defines creation time of the merged container.
	creationTime    time.Time
	creationTimeSet bool
}

// NewDefaultContainerMergerBase builds the empty base state a subclass embeds. Port of the
// protected empty constructor. The subclass constructor must follow it with
// InitDefaultContainerMerger, and then either InitFromDocuments or InitFromASiCContents.
func NewDefaultContainerMergerBase() *DefaultContainerMerger {
	return &DefaultContainerMerger{}
}

// InitDefaultContainerMerger registers the concrete merger with its base so the base can
// dispatch GetContainerExtractor/IsSupportedDocument/IsSupportedContent/
// EnsureContainerContentAllowMerge/EnsureSignaturesAllowMerge/GetTargetASiCContainerType. Every
// concrete merger constructor must call this once, before InitFromDocuments/
// InitFromASiCContents (both call back into the overrides).
func (m *DefaultContainerMerger) InitDefaultContainerMerger(overrides DefaultContainerMergerOverrides) {
	m.overrides = overrides
}

func (m *DefaultContainerMerger) requireOverrides() DefaultContainerMergerOverrides {
	if m.overrides == nil {
		panic("DefaultContainerMerger was not initialised: the concrete merger must call InitDefaultContainerMerger in its constructor")
	}
	return m.overrides
}

// InitFromDocuments is used to create an ASiCContainerMerger from provided container
// documents. Port of the protected DefaultContainerMerger(DSSDocument...) constructor, split
// out because it calls back into GetContainerExtractor (an overrides method) - see
// InitDefaultContainerMerger's doc comment.
//
// Panics with Java's messages on invalid input (Objects.requireNonNull / NullPointerException).
func (m *DefaultContainerMerger) InitFromDocuments(containers ...model.DSSDocument) {
	defaultContainerMergerAssertDocumentsNotNull(containers)
	overrides := m.requireOverrides()
	asicContents := make([]*ASiCContent, len(containers))
	for i, container := range containers {
		extractor := overrides.GetContainerExtractor(container)
		content, err := extractor.Extract()
		if err != nil {
			panic(err)
		}
		asicContents[i] = content
	}
	m.AsicContents = asicContents
}

// InitFromASiCContents is used to create an ASiCContainerMerger from the given ASiCContents.
// Port of the protected DefaultContainerMerger(ASiCContent...) constructor.
func (m *DefaultContainerMerger) InitFromASiCContents(asicContents ...*ASiCContent) {
	defaultContainerMergerAssertContentsNotNull(asicContents)
	m.AsicContents = asicContents
}

// GetCreationTime gets the merged container result's creation time. Ports getCreationTime().
func (m *DefaultContainerMerger) GetCreationTime() time.Time {
	if !m.creationTimeSet {
		m.creationTime = time.Now()
		m.creationTimeSet = true
	}
	return m.creationTime
}

// SetCreationTime sets the creation time of the merged container result (optional). Ports
// setCreationTime(Date).
func (m *DefaultContainerMerger) SetCreationTime(creationTime time.Time) {
	m.creationTime = creationTime
	m.creationTimeSet = true
}

// defaultContainerMergerFactoryRegistry holds the ASiCContainerMergerFactory implementations
// registered via RegisterASiCContainerMergerFactory, consulted in registration order - the Go
// equivalent of Java's ServiceLoader.load(ASiCContainerMergerFactory.class) iteration.
var defaultContainerMergerFactoryRegistry []ASiCContainerMergerFactory

// RegisterASiCContainerMergerFactory registers an ASiCContainerMergerFactory to be consulted by
// DefaultContainerMergerFromDocuments and DefaultContainerMergerFromASiCContents.
func RegisterASiCContainerMergerFactory(f ASiCContainerMergerFactory) {
	defaultContainerMergerFactoryRegistry = append(defaultContainerMergerFactoryRegistry, f)
}

// DefaultContainerMergerFromDocuments loads a relevant ASiCContainerMerger to be used to merge
// given container documents. Ports the static fromDocuments(DSSDocument...).
//
// Java's UnsupportedOperationException("Document format not recognized/handled") is returned
// as an error instead, matching the EvidenceRecordAnalyzerFromDocument precedent (data-
// dependent on which merger factories happen to be registered).
func DefaultContainerMergerFromDocuments(containers ...model.DSSDocument) (ASiCContainerMerger, error) {
	defaultContainerMergerAssertDocumentsNotNull(containers)
	for _, mergerFactory := range defaultContainerMergerFactoryRegistry {
		if mergerFactory.IsSupportedDocuments(containers...) {
			return mergerFactory.CreateFromDocuments(containers...), nil
		}
	}
	return nil, fmt.Errorf("document format not recognized/handled")
}

// DefaultContainerMergerFromASiCContents loads a relevant ASiCContainerMerger to be used to
// merge given ASiCContents. Ports the static fromASiCContents(ASiCContent...).
func DefaultContainerMergerFromASiCContents(asicContents ...*ASiCContent) (ASiCContainerMerger, error) {
	defaultContainerMergerAssertContentsNotNull(asicContents)
	for _, mergerFactory := range defaultContainerMergerFactoryRegistry {
		if mergerFactory.IsSupportedContents(asicContents...) {
			return mergerFactory.CreateFromContents(asicContents...), nil
		}
	}
	return nil, fmt.Errorf("document format not recognized/handled")
}

// IsSupportedDocuments implements ASiCContainerMerger. Ports the @Override
// isSupported(DSSDocument...).
func (m *DefaultContainerMerger) IsSupportedDocuments(containers ...model.DSSDocument) bool {
	defaultContainerMergerAssertDocumentsNotNull(containers)
	overrides := m.requireOverrides()
	for _, containerDocument := range containers {
		isZip, err := ASiCUtilsIsZip(containerDocument)
		if err != nil {
			panic(err)
		}
		if !isZip {
			panic(exception.NewIllegalInputException(
				fmt.Sprintf("The document with name '%s' is not a ZIP archive!", containerDocument.Name())))
		}
		if !overrides.IsSupportedDocument(containerDocument) {
			return false
		}
	}
	return true
}

// IsSupportedContents implements ASiCContainerMerger. Ports the @Override
// isSupported(ASiCContent...).
func (m *DefaultContainerMerger) IsSupportedContents(asicContents ...*ASiCContent) bool {
	defaultContainerMergerAssertContentsNotNull(asicContents)
	overrides := m.requireOverrides()
	for _, asicContent := range asicContents {
		if !overrides.IsSupportedContent(asicContent) {
			return false
		}
	}
	return true
}

// Merge implements ASiCContainerMerger. Ports the @Override merge().
func (m *DefaultContainerMerger) Merge() model.DSSDocument {
	mergeResult := m.MergeToASiCContent()
	containerDocument, err := ZipUtilsInstance().CreateZipArchiveAt(mergeResult, m.GetCreationTime())
	if err != nil {
		panic(err)
	}
	containerDocument.SetName(m.getFinalContainerName(mergeResult.ContainerType()))
	if enumerations.ASiCContainerType_ASiC_S == mergeResult.ContainerType() {
		containerDocument.SetMimeType(enumerations.MimeTypeEnum_ASICS)
	} else {
		containerDocument.SetMimeType(enumerations.MimeTypeEnum_ASICE)
	}
	return containerDocument
}

// MergeToASiCContent implements ASiCContainerMerger. Ports the @Override
// mergeToASiCContent().
//
// Panics with Java's NullPointerException message when no container was provided.
func (m *DefaultContainerMerger) MergeToASiCContent() *ASiCContent {
	if len(m.AsicContents) == 0 {
		panic("At least one container shall be provided!")
	}

	overrides := m.requireOverrides()
	overrides.EnsureContainerContentAllowMerge()
	overrides.EnsureSignaturesAllowMerge()
	return m.createMergedResult()
}

// createMergedResult creates a new ASiCContent by merging the given containers. Ports the
// protected createMergedResult().
//
// Cross-chunk assumption (ZIPCORE): ASiCContent exposes the getter/setter pairs named below
// (Go-cased from the Java getX/setX accessors) for SignedDocuments, SignatureDocuments,
// ManifestDocuments, ArchiveManifestDocuments, EvidenceRecordManifestDocuments,
// TimestampDocuments, EvidenceRecordDocuments, UnsupportedDocuments and Folders, all
// []model.DSSDocument.
func (m *DefaultContainerMerger) createMergedResult() *ASiCContent {
	asicContent := m.createEmptyContainer()

	asicContent.SetZipComment(m.getZipComment())
	asicContent.SetMimeTypeDocument(m.getMimeTypeDocument())

	asicContent.SetSignedDocuments(m.mergeDocumentLists(collectDocumentLists(m.AsicContents, (*ASiCContent).SignedDocuments)))
	asicContent.SetSignatureDocuments(m.mergeDocumentLists(collectDocumentLists(m.AsicContents, (*ASiCContent).SignatureDocuments)))
	asicContent.SetManifestDocuments(m.mergeDocumentLists(collectDocumentLists(m.AsicContents, (*ASiCContent).ManifestDocuments)))
	asicContent.SetArchiveManifestDocuments(m.mergeDocumentLists(collectDocumentLists(m.AsicContents, (*ASiCContent).ArchiveManifestDocuments)))
	asicContent.SetEvidenceRecordManifestDocuments(m.mergeDocumentLists(collectDocumentLists(m.AsicContents, (*ASiCContent).EvidenceRecordManifestDocuments)))
	asicContent.SetTimestampDocuments(m.mergeDocumentLists(collectDocumentLists(m.AsicContents, (*ASiCContent).TimestampDocuments)))
	asicContent.SetEvidenceRecordDocuments(m.mergeDocumentLists(collectDocumentLists(m.AsicContents, (*ASiCContent).EvidenceRecordDocuments)))
	asicContent.SetUnsupportedDocuments(m.mergeDocumentLists(collectDocumentLists(m.AsicContents, (*ASiCContent).UnsupportedDocuments)))
	asicContent.SetFolders(m.mergeDocumentLists(collectDocumentLists(m.AsicContents, (*ASiCContent).Folders)))

	return asicContent
}

// collectDocumentLists ports the repeated
// Arrays.stream(asicContents).map(ASiCContent::getX).collect(Collectors.toList()) idiom as a
// small generic helper local to this file (not a cross-file shared helper, per PORTING.md).
func collectDocumentLists(asicContents []*ASiCContent, getter func(*ASiCContent) []model.DSSDocument) [][]model.DSSDocument {
	result := make([][]model.DSSDocument, len(asicContents))
	for i, c := range asicContents {
		result[i] = getter(c)
	}
	return result
}

// createEmptyContainer creates an empty container. Ports the protected
// createEmptyContainer().
//
// Cross-chunk assumption (ZIPCORE): NewASiCContent() *ASiCContent and
// SetContainerType(enumerations.ASiCContainerType).
func (m *DefaultContainerMerger) createEmptyContainer() *ASiCContent {
	asicContent := NewASiCContent()
	asicContent.SetContainerType(m.getContainerType())
	return asicContent
}

// getContainerType ports the private getContainerType().
func (m *DefaultContainerMerger) getContainerType() enumerations.ASiCContainerType {
	for _, asicContent := range m.AsicContents {
		if asicContent.ContainerType() != "" {
			return asicContent.ContainerType()
		}
	}
	return m.requireOverrides().GetTargetASiCContainerType()
}

// getZipComment ports the private getZipComment().
//
// Panics with Java's UnsupportedOperationException message on conflicting zip comments.
func (m *DefaultContainerMerger) getZipComment() string {
	var zipComment string
	for _, asicContent := range m.AsicContents {
		currentZipComment := asicContent.ZipComment()
		if utils.IsStringNotEmpty(currentZipComment) {
			if utils.IsStringEmpty(zipComment) {
				zipComment = currentZipComment
			} else if zipComment != currentZipComment {
				panic(fmt.Sprintf("Unable to merge containers. "+
					"Containers contain different zip comments : '%s' and '%s'!", zipComment, currentZipComment))
			}
		}
	}
	return zipComment
}

// getMimeTypeDocument ports the private getMimeTypeDocument().
//
// Panics with Java's UnsupportedOperationException message on conflicting mimetype documents.
func (m *DefaultContainerMerger) getMimeTypeDocument() model.DSSDocument {
	var mimeTypeDoc model.DSSDocument
	for _, asicContent := range m.AsicContents {
		currentMimeTypeDocument := asicContent.MimeTypeDocument()
		if currentMimeTypeDocument != nil {
			if mimeTypeDoc == nil {
				mimeTypeDoc = currentMimeTypeDocument
			} else {
				a, errA := mimeTypeDoc.DigestValue(DefaultContainerMergerDefaultDigestAlgorithm)
				b, errB := currentMimeTypeDocument.DigestValue(DefaultContainerMergerDefaultDigestAlgorithm)
				if errA != nil || errB != nil || !bytes.Equal(a, b) {
					panic("Unable to merge containers. Containers contain different mimetype documents!")
				}
			}
		}
	}
	return mimeTypeDoc
}

// mergeDocumentLists ports the private mergeDocumentLists(Collection).
//
// Panics with Java's UnsupportedOperationException message when two containers carry different
// documents under the same name.
func (m *DefaultContainerMerger) mergeDocumentLists(documentsLists [][]model.DSSDocument) []model.DSSDocument {
	result := make([]model.DSSDocument, 0)
	addedDocumentNames := make([]string, 0)
	for _, documentsList := range documentsLists {
		for _, document := range documentsList {
			if !containsString(addedDocumentNames, document.Name()) {
				result = append(result, document)
				addedDocumentNames = append(addedDocumentNames, document.Name())
			} else {
				originalListDocument := spi.DSSUtilsDocumentWithName(result, document.Name())
				a, errA := originalListDocument.DigestValue(DefaultContainerMergerDefaultDigestAlgorithm)
				b, errB := document.DigestValue(DefaultContainerMergerDefaultDigestAlgorithm)
				if errA != nil || errB != nil || !bytes.Equal(a, b) {
					panic(fmt.Sprintf("Unable to merge containers. "+
						"Containers contain different documents under the same name : %s!", document.Name()))
				}
				// continue, no document to be added
			}
		}
	}
	return result
}

// getFinalContainerName returns a filename for the merged container. Ports the protected
// getFinalContainerName(ASiCContainerType).
func (m *DefaultContainerMerger) getFinalContainerName(asicContainerType enumerations.ASiCContainerType) string {
	originalFilename := m.getOriginalContainerFilename()
	originalExtension := utils.GetFileNameExtension(originalFilename)
	if utils.IsStringNotEmpty(originalExtension) {
		// remove extension
		originalFilename = originalFilename[:len(originalFilename)-len(originalExtension)-1]
	}

	var sb strings.Builder
	sb.WriteString(originalFilename)
	sb.WriteString("-merged")
	sb.WriteString(".")

	finalExtension := m.getFinalExtension(asicContainerType, originalExtension)
	sb.WriteString(finalExtension)

	return sb.String()
}

// getOriginalContainerFilename ports the private getOriginalContainerFilename().
//
// Cross-chunk assumption (ZIPCORE): ASiCContent exposes AsicContainer() model.DSSDocument.
func (m *DefaultContainerMerger) getOriginalContainerFilename() string {
	for _, asicContent := range m.AsicContents {
		if asicContent.AsicContainer() != nil && asicContent.AsicContainer().Name() != "" {
			return asicContent.AsicContainer().Name()
		}
	}
	return "container"
}

// getFinalExtension ports the private getFinalExtension(ASiCContainerType, String).
func (m *DefaultContainerMerger) getFinalExtension(asicContainerType enumerations.ASiCContainerType, originalExtension string) string {
	if utils.IsStringNotEmpty(originalExtension) {
		return originalExtension
	} else if asicContainerType != "" {
		if enumerations.ASiCContainerType_ASiC_S == asicContainerType {
			return enumerations.MimeTypeEnum_ASICS.Extension()
		}
		return enumerations.MimeTypeEnum_ASICE.Extension()
	}
	return "zip"
}

// defaultContainerMergerAssertDocumentsNotNull ports the private static
// assertNotNull(DSSDocument...).
//
// Panics with Java's NullPointerException message when no document was provided, or a
// contained document is nil.
func defaultContainerMergerAssertDocumentsNotNull(containers []model.DSSDocument) {
	if len(containers) == 0 {
		panic("At least one document shall be provided!")
	}
	for _, containerDocument := range containers {
		if containerDocument == nil {
			panic("DSSDocument cannot be null!")
		}
	}
}

// defaultContainerMergerAssertContentsNotNull ports the private static
// assertNotNull(ASiCContent...).
func defaultContainerMergerAssertContentsNotNull(asicContents []*ASiCContent) {
	if len(asicContents) == 0 {
		panic("At least one ASiCContent shall be provided!")
	}
	for _, asicContent := range asicContents {
		if asicContent == nil {
			panic("ASiCContent cannot be null!")
		}
	}
}
