// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/signature/asice/ASiCEWithXAdESManifestBuilder.java
// (DSS 6.5.RC1).
//
// Java's ASiCEWithXAdESManifestBuilder does NOT extend the ASiCManifest*.xml base used by
// dss-asic-cades (asic.AbstractASiCManifestBuilder) - it is a standalone class that produces a
// different, OpenDocument-flavoured manifest.xml. This port therefore does not embed that base
// either.
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// ASiCEWithXAdESManifestBuilder is used to build the manifest.xml file (ASiC-E).
//
// Sample:
//
//	<manifest:manifest xmlns:manifest="urn:oasis:names:tc:opendocument:xmlns:manifest:1.0" manifest:version="1.2">
//		<manifest:file-entry manifest:full-path="/" manifest:media-type="application/vnd.etsi.asic-e+zip"/>
//		<manifest:file-entry manifest:full-path="test.txt" manifest:media-type="text/plain"/>
//		<manifest:file-entry manifest:full-path="test-data-file.bin" manifest:media-type="application/octet-stream"/>
//	</manifest:manifest>
type ASiCEWithXAdESManifestBuilder struct {
	// documents is the list of documents to be included into the Manifest.
	//
	// WARN: shall not be used together with entries.
	documents []model.DSSDocument

	// entries is the list of manifest entries to be included into the Manifest.
	//
	// WARN: shall not be used together with documents.
	entries []*model.ManifestEntry

	// manifestFilename is the name of the created manifest document.
	manifestFilename string
}

// NewASiCEWithXAdESManifestBuilder is the empty constructor.
func NewASiCEWithXAdESManifestBuilder() *ASiCEWithXAdESManifestBuilder {
	return &ASiCEWithXAdESManifestBuilder{}
}

// SetDocuments sets documents to be included into the Manifest.
//
// WARN: shall not be used together with SetEntries(entries).
//
// Ports setDocuments(List<DSSDocument>).
func (b *ASiCEWithXAdESManifestBuilder) SetDocuments(documents []model.DSSDocument) *ASiCEWithXAdESManifestBuilder {
	b.documents = documents
	return b
}

// SetEntries sets manifest entries to be included into the Manifest.
//
// WARN: shall not be used together with SetDocuments(documents).
//
// Ports setEntries(List<ManifestEntry>).
func (b *ASiCEWithXAdESManifestBuilder) SetEntries(entries []*model.ManifestEntry) *ASiCEWithXAdESManifestBuilder {
	b.entries = entries
	return b
}

// SetManifestFilename sets the target name of the XML Manifest file to be created. Ports
// setManifestFilename(String).
func (b *ASiCEWithXAdESManifestBuilder) SetManifestFilename(manifestFilename string) *ASiCEWithXAdESManifestBuilder {
	b.manifestFilename = manifestFilename
	return b
}

// Build builds the XML manifest. Ports build().
//
// Panics with Java's NullPointerException/IllegalArgumentException messages
// (Objects.requireNonNull) on invalid entries/state.
func (b *ASiCEWithXAdESManifestBuilder) Build() (model.DSSDocument, error) {
	documentDom := xmlutils.DomUtilsBuildDOMEmpty()
	manifestDom := xmlutils.DomUtilsCreateElementNS(documentDom, ManifestNS, ManifestElement_MANIFEST)
	xmlutils.DomUtilsSetAttributeNS(manifestDom, ManifestNS, ManifestAttribute_VERSION, "1.2")
	documentDom.AppendChild(manifestDom)

	rootDom := xmlutils.DomUtilsAddElement(documentDom, manifestDom, ManifestNS, ManifestElement_FILE_ENTRY)
	xmlutils.DomUtilsSetAttributeNS(rootDom, ManifestNS, ManifestAttribute_FULL_PATH, "/")
	xmlutils.DomUtilsSetAttributeNS(rootDom, ManifestNS, ManifestAttribute_MEDIA_TYPE, enumerations.MimeTypeEnum_ASICE.MimeTypeString())

	for _, entry := range b.getEntries() {
		if entry == nil {
			panic("ManifestEntry cannot be null!")
		}
		if entry.Uri() == "" {
			panic("ManifestEntry#Uri cannot be null! Please define the document's name.")
		}

		fileDom := xmlutils.DomUtilsAddElement(documentDom, manifestDom, ManifestNS, ManifestElement_FILE_ENTRY)
		xmlutils.DomUtilsSetAttributeNS(fileDom, ManifestNS, ManifestAttribute_FULL_PATH, entry.Uri())
		mimeType := entry.MimeType()
		if mimeType == nil {
			// Upstream logs a warning naming the manifest entry's URI here (dropped per
			// PORTING.md).
			mimeType = enumerations.MimeTypeEnum_BINARY
		}
		xmlutils.DomUtilsSetAttributeNS(fileDom, ManifestNS, ManifestAttribute_MEDIA_TYPE, mimeType.MimeTypeString())
	}

	return xmlutils.DomUtilsCreateDssDocumentFromDomDocument(documentDom, b.manifestFilename)
}

// getEntries ports the private getEntries().
//
// Panics with Java's IllegalArgumentException/NullPointerException messages - both share the
// exact same message text upstream - when neither or both of documents/entries are provided.
func (b *ASiCEWithXAdESManifestBuilder) getEntries() []*model.ManifestEntry {
	documentsNotEmpty := utils.IsCollectionNotEmpty(b.documents)
	entriesNotEmpty := utils.IsCollectionNotEmpty(b.entries)
	switch {
	case documentsNotEmpty:
		if entriesNotEmpty {
			panic("Either DSSDocuments or ManifestEntries shall be provided!")
		}
		return asic.ASiCUtilsToSimpleManifestEntries(b.documents)
	case entriesNotEmpty:
		return b.entries
	default:
		panic("Either DSSDocuments or ManifestEntries shall be provided!")
	}
}
