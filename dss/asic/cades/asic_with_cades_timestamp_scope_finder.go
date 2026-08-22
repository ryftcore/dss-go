// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/validation/scope/ASiCWithCAdESTimestampScopeFinder.java (DSS 6.5.RC1).
//
// The Java `validation.scope` sub-package flattens into this Go package.
//
// Virtual-dispatch note: Java's findTimestampScope() and
// getTimestampSignatureScopeForDocument() are both @Override'd here, and findTimestampScope()
// self-calls getTimestampSignatureScopeForDocument() expecting the override to be reached. The
// embedded scope.DetachedTimestampScopeFinder has no overrides-field indirection (unlike the
// AbstractASiCContainerAnalyzer/DefaultContainerMerger precedents), so both methods are
// reimplemented directly on this leaf type rather than relying on embedding to reach the
// override - the leaf's FindTimestampScope calls its own GetTimestampSignatureScopeForDocument
// method by name, which Go resolves statically to this type's (shadowing) definition.
package cades

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
	mscope "github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/scope"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ASiCWithCAdESTimestampScopeFinder is used to find a timestamp source for a detached timestamp
// within an ASiC with CAdES container.
type ASiCWithCAdESTimestampScopeFinder struct {
	scope.DetachedTimestampScopeFinder

	// containerDocuments represents a list of documents encapsulated within an ASiC container.
	containerDocuments []model.DSSDocument

	// archiveDocuments represents a list of documents encapsulated within a package.zip archive,
	// when applicable (ASiC-S).
	archiveDocuments []model.DSSDocument
}

// NewASiCWithCAdESTimestampScopeFinder is the default constructor instantiating object with
// empty lists of documents.
func NewASiCWithCAdESTimestampScopeFinder() *ASiCWithCAdESTimestampScopeFinder {
	return &ASiCWithCAdESTimestampScopeFinder{
		DetachedTimestampScopeFinder: *scope.NewDetachedTimestampScopeFinder(),
	}
}

// SetContainerDocuments sets a list of container original documents. Port of
// setContainerDocuments(List).
func (f *ASiCWithCAdESTimestampScopeFinder) SetContainerDocuments(containerDocuments []model.DSSDocument) {
	f.containerDocuments = containerDocuments
}

// SetArchiveDocuments sets a list of package.zip archive documents. Port of
// setArchiveDocuments(List).
func (f *ASiCWithCAdESTimestampScopeFinder) SetArchiveDocuments(archiveDocuments []model.DSSDocument) {
	f.archiveDocuments = archiveDocuments
}

// FindTimestampScope ports the @Override findTimestampScope(TimestampToken).
func (f *ASiCWithCAdESTimestampScopeFinder) FindTimestampScope(timestampToken *validation.TimestampToken) []mscope.SignatureScope {
	if timestampToken.IsMessageImprintDataIntact() {
		if timestampToken.ManifestFile() != nil {
			return f.getTimestampSignatureScopeForManifest(timestampToken.ManifestFile())
		}
		return f.GetTimestampSignatureScopeForDocument(f.TimestampedData)
	}
	return []mscope.SignatureScope{}
}

// getTimestampSignatureScopeForManifest extracts timestamped signature scopes from a
// ManifestFile. Ports the private getTimestampSignatureScopeForManifest(ManifestFile).
func (f *ASiCWithCAdESTimestampScopeFinder) getTimestampSignatureScopeForManifest(manifestFile *model.ManifestFile) []mscope.SignatureScope {
	result := make([]mscope.SignatureScope, 0)
	result = append(result, scope.NewManifestSignatureScope(manifestFile))
	if utils.IsCollectionNotEmpty(f.containerDocuments) {
		rootLevelDocuments := asic.UtilsRootLevelDocuments(f.containerDocuments)
		for _, manifestEntry := range manifestFile.Entries() {
			result = append(result, f.getTimestampSignatureScopeForManifestEntry(manifestEntry, rootLevelDocuments)...)
		}
	}
	return result
}

// getTimestampSignatureScopeForManifestEntry ports the private
// getTimestampSignatureScopeForManifestEntry(ManifestEntry, List).
func (f *ASiCWithCAdESTimestampScopeFinder) getTimestampSignatureScopeForManifestEntry(manifestEntry *model.ManifestEntry, rootLevelDocuments []model.DSSDocument) []mscope.SignatureScope {
	if manifestEntry.IsIntact() {
		for _, document := range f.containerDocuments {
			if utils.AreStringsEqual(manifestEntry.Uri(), document.Name()) {
				if utils.CollectionSize(rootLevelDocuments) == 1 && f.isASiCSContainer(document) {
					return f.getTimestampSignatureScopeForZipPackage(document)
				}
				return f.DetachedTimestampScopeFinder.GetTimestampSignatureScopeForDocument(document)
			}
		}
	}
	return []mscope.SignatureScope{}
}

// GetTimestampSignatureScopeForDocument ports the @Override protected
// getTimestampSignatureScopeForDocument(DSSDocument).
func (f *ASiCWithCAdESTimestampScopeFinder) GetTimestampSignatureScopeForDocument(document model.DSSDocument) []mscope.SignatureScope {
	if f.isASiCSContainer(document) {
		return f.getTimestampSignatureScopeForZipPackage(document)
	}
	return f.DetachedTimestampScopeFinder.GetTimestampSignatureScopeForDocument(document)
}

// getTimestampSignatureScopeForZipPackage ports the private
// getTimestampSignatureScopeForZipPackage(DSSDocument).
func (f *ASiCWithCAdESTimestampScopeFinder) getTimestampSignatureScopeForZipPackage(document model.DSSDocument) []mscope.SignatureScope {
	result := make([]mscope.SignatureScope, 0)
	result = append(result, scope.NewContainerSignatureScope(document))
	if utils.IsCollectionNotEmpty(f.archiveDocuments) {
		for _, archivedDocument := range f.archiveDocuments {
			result = append(result, scope.NewContainerContentSignatureScope(archivedDocument))
		}
	}
	return result
}

// isASiCSContainer verifies whether the given document is an ASiC-S ZIP container. Ports the
// private isASiCSContainer(DSSDocument).
func (f *ASiCWithCAdESTimestampScopeFinder) isASiCSContainer(document model.DSSDocument) bool {
	if document.Name() == "" || strings.Contains(document.Name(), "/") {
		return false
	}
	isZip, err := asic.UtilsIsZip(document)
	return err == nil && isZip
}
