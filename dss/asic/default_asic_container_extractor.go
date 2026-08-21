// Ported from
// dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/extract/DefaultASiCContainerExtractor.java
// (DSS 6.5.RC1).
//
// The Java extract sub-package flattens into this Go package per the phase-7 package layout.
//
// Java's six abstract isAllowed*(String) predicates - the whole reason this class is abstract - are
// carried by DefaultASiCContainerExtractorOverrides. Routing zipParsing's calls through that
// interface is what keeps a subclass's predicates in play; plain Go embedding would bind them
// statically to this file's (non-existent) implementations.
//
// slf4j LOG calls are dropped per PORTING.md; the branches they sit in are preserved.
package asic

import (
	"fmt"
	"strings"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/utils"
)

// DefaultASiCContainerExtractorOverrides captures the abstract methods of Java's
// DefaultASiCContainerExtractor that the class calls on itself.
type DefaultASiCContainerExtractorOverrides interface {
	// IsAllowedManifest checks if the given file name represents an allowed manifest name for
	// the current ASiC container format. Port of the protected abstract
	// isAllowedManifest(String).
	IsAllowedManifest(entryName string) bool

	// IsAllowedArchiveManifest checks if the given file name represents an allowed archive
	// manifest name for the current ASiC container format. Port of the protected abstract
	// isAllowedArchiveManifest(String).
	IsAllowedArchiveManifest(entryName string) bool

	// IsAllowedEvidenceRecordManifest checks if the given file name represents an allowed
	// evidence record manifest name for the current ASiC container format. Port of the
	// protected abstract isAllowedEvidenceRecordManifest(String).
	IsAllowedEvidenceRecordManifest(entryName string) bool

	// IsAllowedSignature checks if the given file name represents an allowed signature document
	// name for the current ASiC container format. Port of the protected abstract
	// isAllowedSignature(String).
	IsAllowedSignature(entryName string) bool

	// IsAllowedTimestamp checks if the given file name represents an allowed timestamp document
	// name for the current ASiC container format. Port of the protected abstract
	// isAllowedTimestamp(String).
	IsAllowedTimestamp(entryName string) bool

	// IsAllowedEvidenceRecord checks if the given file name represents an allowed evidence
	// record document name for the current ASiC container format. Port of the protected
	// abstract isAllowedEvidenceRecord(String).
	IsAllowedEvidenceRecord(entryName string) bool
}

// DefaultASiCContainerExtractor is used to read an ASiC Container and to retrieve its content files.
type DefaultASiCContainerExtractor struct {
	// overrides points back at the concrete extractor; see InitDefaultASiCContainerExtractor.
	overrides DefaultASiCContainerExtractorOverrides

	// AsicContainer represents an ASiC container. Port of the protected final asicContainer.
	AsicContainer model.DSSDocument
}

// InitDefaultASiCContainerExtractor is the port of the protected
// DefaultASiCContainerExtractor(DSSDocument) constructor, extended with the overrides argument Go
// needs to keep the subclass predicates reachable.
func (e *DefaultASiCContainerExtractor) InitDefaultASiCContainerExtractor(overrides DefaultASiCContainerExtractorOverrides, asicContainer model.DSSDocument) {
	e.overrides = overrides
	e.AsicContainer = asicContainer
}

// DefaultASiCContainerExtractorFromDocument loads an implementation of ASiCContainerExtractor
// corresponding to the asicContainer type.
//
// Panics with the Java message when asicContainer is nil (Objects.requireNonNull); returns an error
// when no registered factory supports it (UnsupportedOperationException).
//
// Port of the static fromDocument(DSSDocument); the ServiceLoader iteration becomes the registry in
// asic_container_extractor_factory.go.
func DefaultASiCContainerExtractorFromDocument(asicContainer model.DSSDocument) (ASiCContainerExtractor, error) {
	if asicContainer == nil {
		panic("ASiC container cannot be null!")
	}

	for _, factory := range asicContainerExtractorFactoryRegistry {
		if factory.IsSupported(asicContainer) {
			return factory.Create(asicContainer), nil
		}
	}
	return nil, fmt.Errorf("Document format not recognized/handled")
}

// Extract extracts the content (documents) embedded into the ASiC container. Port of extract().
func (e *DefaultASiCContainerExtractor) Extract() (*ASiCContent, error) {
	result, err := e.zipParsing(e.AsicContainer)
	if err != nil {
		return nil, err
	}
	zipComment, err := ASiCUtilsZipCommentFromArchiveContainer(e.AsicContainer)
	if err != nil {
		return nil, err
	}
	result.SetZipComment(zipComment)
	containerType, err := ASiCUtilsContainerTypeOfContent(result)
	if err != nil {
		return nil, err
	}
	result.SetContainerType(containerType)
	containerDocuments, err := e.containerDocuments(result)
	if err != nil {
		return nil, err
	}
	result.SetContainerDocuments(containerDocuments)
	return result, nil
}

// zipParsing is the port of the private zipParsing(DSSDocument).
func (e *DefaultASiCContainerExtractor) zipParsing(asicContainer model.DSSDocument) (*ASiCContent, error) {
	result := NewASiCContent()
	result.SetAsicContainer(asicContainer)

	documents, err := ZipUtilsInstance().ExtractContainerContent(asicContainer)
	if err != nil {
		return nil, err
	}
	if utils.IsCollectionEmpty(documents) {
		return nil, exception.NewIllegalInputException(fmt.Sprintf(
			"The provided file with name '%s' does not contain documents inside. "+
				"Probably file has an unsupported format or has been corrupted. "+
				"The signature validation is not possible", asicContainer.Name()))
	}

	for _, currentDocument := range documents {
		entryName := currentDocument.Name()

		if e.isMetaInfFolder(entryName) {
			if e.overrides.IsAllowedSignature(entryName) {
				result.SetSignatureDocuments(append(result.SignatureDocuments(), currentDocument))
			} else if e.overrides.IsAllowedManifest(entryName) {
				result.SetManifestDocuments(append(result.ManifestDocuments(), currentDocument))
			} else if e.overrides.IsAllowedArchiveManifest(entryName) {
				result.SetArchiveManifestDocuments(append(result.ArchiveManifestDocuments(), currentDocument))
			} else if e.overrides.IsAllowedEvidenceRecordManifest(entryName) {
				result.SetEvidenceRecordManifestDocuments(append(result.EvidenceRecordManifestDocuments(), currentDocument))
			} else if e.overrides.IsAllowedTimestamp(entryName) {
				result.SetTimestampDocuments(append(result.TimestampDocuments(), currentDocument))
			} else if e.overrides.IsAllowedEvidenceRecord(entryName) {
				result.SetEvidenceRecordDocuments(append(result.EvidenceRecordDocuments(), currentDocument))
			} else if !e.isFolder(entryName) {
				result.SetUnsupportedDocuments(append(result.UnsupportedDocuments(), currentDocument))
			}

		} else if !e.isFolder(entryName) {
			if ASiCUtilsIsMimetype(entryName) {
				result.SetMimeTypeDocument(currentDocument)
			} else {
				result.SetSignedDocuments(append(result.SignedDocuments(), currentDocument))
			}

		} else {
			result.SetFolders(append(result.Folders(), currentDocument))
		}
	}

	// (upstream warns "Unsupported files : {}" when any were collected)
	return result, nil
}

// containerDocuments is the port of the private getContainerDocuments(ASiCContent).
func (e *DefaultASiCContainerExtractor) containerDocuments(asicContent *ASiCContent) ([]model.DSSDocument, error) {
	containerDocuments := make([]model.DSSDocument, 0)
	isASiCS, err := ASiCUtilsIsASiCSContainerContent(asicContent)
	if err != nil {
		return nil, err
	}
	if isASiCS {
		for _, signerDocument := range asicContent.RootLevelSignedDocuments() {
			if utils.IsCollectionNotEmpty(containerDocuments) {
				// More than one ZIP archive found on a root level of the ASiC-S
				// container! Extraction of embedded documents not possible.
				return []model.DSSDocument{}, nil
			}
			isZip, err := ASiCUtilsIsZip(signerDocument)
			if err != nil {
				return nil, err
			}
			if isZip {
				extracted, err := ZipUtilsInstance().ExtractContainerContent(signerDocument)
				if err != nil {
					return nil, err
				}
				containerDocuments = append(containerDocuments, extracted...)
			}
		}
	}
	return containerDocuments, nil
}

// isMetaInfFolder is the port of the private isMetaInfFolder(String).
func (e *DefaultASiCContainerExtractor) isMetaInfFolder(entryName string) bool {
	return strings.HasPrefix(entryName, ASiCUtilsMetaInfFolder)
}

// isFolder is the port of the private isFolder(String).
func (e *DefaultASiCContainerExtractor) isFolder(entryName string) bool {
	return strings.HasSuffix(entryName, "/")
}
