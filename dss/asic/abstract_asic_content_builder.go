// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/signature/AbstractASiCContentBuilder.java (DSS 6.5.RC1).
package asic

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/exception"
	"github.com/utain/esig/dss/utils"
)

// abstractASiCContentBuilderZipEntryDetachedFile is the default name for a detached file if
// one is not defined. Port of the private static final ZIP_ENTRY_DETACHED_FILE.
const abstractASiCContentBuilderZipEntryDetachedFile = "detached-file"

// AbstractASiCContentBuilderOverrides declares the operation AbstractASiCContentBuilder calls
// back into virtually from Build(). Per S7_BRIEF.md's virtual-dispatch warning, every concrete
// builder must call InitAbstractASiCContentBuilder with itself before use.
type AbstractASiCContentBuilderOverrides interface {
	// GetContainerExtractor returns an instance of a corresponding container extractor class.
	// Port of the protected abstract getContainerExtractor(DSSDocument).
	GetContainerExtractor(archiveDocument model.DSSDocument) ASiCContainerExtractor
}

// AbstractASiCContentBuilder builds an instance of ASiCContent. As input, an ASiC Container can
// be used or documents to be signed.
type AbstractASiCContentBuilder struct {
	// overrides points back at the concrete builder; see InitAbstractASiCContentBuilder.
	overrides AbstractASiCContentBuilderOverrides
}

// NewAbstractASiCContentBuilderBase builds the base state a subclass embeds. Port of the
// protected empty constructor. The subclass constructor must follow it with
// InitAbstractASiCContentBuilder.
func NewAbstractASiCContentBuilderBase() *AbstractASiCContentBuilder {
	return &AbstractASiCContentBuilder{}
}

// InitAbstractASiCContentBuilder registers the concrete builder with its base so the base can
// dispatch GetContainerExtractor. Every concrete builder constructor must call this once.
func (b *AbstractASiCContentBuilder) InitAbstractASiCContentBuilder(overrides AbstractASiCContentBuilderOverrides) {
	b.overrides = overrides
}

func (b *AbstractASiCContentBuilder) requireOverrides() AbstractASiCContentBuilderOverrides {
	if b.overrides == nil {
		panic("AbstractASiCContentBuilder was not initialised: the concrete builder must call InitAbstractASiCContentBuilder in its constructor")
	}
	return b.overrides
}

// Build builds the ASiCContent from the given documents, representing an ASiC Container or a
// list of documents to be signed, and the target asicContainerType. Ports
// build(List, ASiCContainerType).
func (b *AbstractASiCContentBuilder) Build(documents []model.DSSDocument, asicContainerType enumerations.ASiCContainerType) *ASiCContent {
	if utils.IsCollectionNotEmpty(documents) && len(documents) == 1 {
		archiveDocument := documents[0]
		isASiC, err := ASiCUtilsIsASiC(archiveDocument)
		if err != nil {
			panic(err)
		}
		if isASiC {
			extractor := b.requireOverrides().GetContainerExtractor(archiveDocument)
			if extractor.IsSupportedContainerFormat() {
				return b.fromZipArchive(extractor, asicContainerType)
			}
		}
	}
	return b.fromFiles(documents, asicContainerType)
}

// fromZipArchive ports the private fromZipArchive(ASiCContainerExtractor, ASiCContainerType).
func (b *AbstractASiCContentBuilder) fromZipArchive(extractor ASiCContainerExtractor, asicContainerType enumerations.ASiCContainerType) *ASiCContent {
	asicContent, err := extractor.Extract()
	if err != nil {
		panic(err)
	}
	b.assertContainerTypeValid(asicContent, asicContainerType)
	return asicContent
}

// fromFiles ports the private fromFiles(List, ASiCContainerType).
func (b *AbstractASiCContentBuilder) fromFiles(documents []model.DSSDocument, asicContainerType enumerations.ASiCContainerType) *ASiCContent {
	b.assertDocumentNamesDefined(documents)

	asicContent := NewASiCContent()
	asicContent.SetContainerType(asicContainerType)
	asicContent.SetSignedDocuments(documents)

	return asicContent
}

// assertContainerTypeValid ports the private assertContainerTypeValid(ASiCContent,
// ASiCContainerType).
func (b *AbstractASiCContentBuilder) assertContainerTypeValid(result *ASiCContent, asicContainerType enumerations.ASiCContainerType) {
	if ASiCUtilsFilesContainSignatures(spi.DSSUtilsDocumentNames(result.AllDocuments())) &&
		utils.IsCollectionEmpty(result.SignatureDocuments()) {
		panic("Container type doesn't match! The same container type shall be chosen.")
	}
	if asicContainerType != result.ContainerType() {
		panic(exception.NewIllegalInputException(fmt.Sprintf(
			"The provided container of type '%s' does not correspond the expected format '%s'!",
			result.ContainerType(), asicContainerType)))
	}
}

// assertDocumentNamesDefined checks if the document names are defined and adds them if
// needed. Ports the private assertDocumentNamesDefined(List).
func (b *AbstractASiCContentBuilder) assertDocumentNamesDefined(documents []model.DSSDocument) {
	unnamedDocuments := b.getDocumentsWithoutNames(documents)
	if len(unnamedDocuments) == 1 {
		unnamedDocuments[0].SetName(abstractASiCContentBuilderZipEntryDetachedFile)
	} else {
		for ii, dssDocument := range unnamedDocuments {
			dssDocument.SetName(fmt.Sprintf("%s-%d", abstractASiCContentBuilderZipEntryDetachedFile, ii))
		}
	}
}

// getDocumentsWithoutNames ports the private getDocumentsWithoutNames(List).
func (b *AbstractASiCContentBuilder) getDocumentsWithoutNames(documents []model.DSSDocument) []model.DSSDocument {
	result := make([]model.DSSDocument, 0)
	for _, document := range documents {
		if utils.IsStringBlank(document.Name()) {
			result = append(result, document)
		}
	}
	return result
}
