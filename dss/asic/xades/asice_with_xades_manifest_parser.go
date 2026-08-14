// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/validation/ASiCEWithXAdESManifestParser.java (DSS 6.5.RC1).
//
// Package flattening: the Java package eu.europa.esig.dss.asic.xades.validation lands in this
// same Go package (dss/asic/xades) per S7_BRIEF.md's package layout table.
//
// FORWARD CROSS-CHUNK DEPENDENCY (flagged per S7_BRIEF.md's "flag needs in notes" rule):
// dss-asic-xades/.../definition/ManifestPath.java is NOT in this manifest (S7_XADVAL) - it is
// needed by no other file in this manifest either, so its ownership is left to whichever chunk
// ports the `definition` package (ManifestElement.java, needed by the XADSIGN manifest builder,
// most likely lands there too). Assumed shape, following the asic/asic_manifest_path.go and
// manifest_attribute.go precedents already in this package:
//
//	var ManifestPath_FILE_ENTRY_PATH common.XPathQuery                                       // ManifestPath.FILE_ENTRY_PATH
//	func ManifestPathGetFullPathAttribute(manifestNamespace *common.DSSNamespace) string      // ManifestPath.getFullPathAttribute(DSSNamespace)
//	func ManifestPathGetMediaTypeAttribute(manifestNamespace *common.DSSNamespace) string     // ManifestPath.getMediaTypeAttribute(DSSNamespace)
package xades

import (
	"strings"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/xml/common"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// asiceWithXAdESManifestParserNamespacesRegistered mirrors the Java static initializer block,
// which registers the ManifestNamespace once at class-load time.
var asiceWithXAdESManifestParserNamespacesRegistered = func() bool {
	xmlutils.XPathUtilsRegisterNamespace(ManifestNS)
	return true
}()

// ASiCEWithXAdESManifestParser parses the ASiC with XAdES manifest document and produces a
// model.ManifestFile.
type ASiCEWithXAdESManifestParser struct {
	// signatureDocument is the related signature document.
	signatureDocument model.DSSDocument

	// manifestDocument is the manifest document to be parsed.
	manifestDocument model.DSSDocument
}

// NewASiCEWithXAdESManifestParser is the constructor with a manifest document only (no assigned
// signature). Ports ASiCEWithXAdESManifestParser(DSSDocument).
func NewASiCEWithXAdESManifestParser(manifestDocument model.DSSDocument) *ASiCEWithXAdESManifestParser {
	return NewASiCEWithXAdESManifestParserWithSignature(nil, manifestDocument)
}

// NewASiCEWithXAdESManifestParserWithSignature is the default constructor. Ports
// ASiCEWithXAdESManifestParser(DSSDocument, DSSDocument).
func NewASiCEWithXAdESManifestParserWithSignature(signatureDocument, manifestDocument model.DSSDocument) *ASiCEWithXAdESManifestParser {
	_ = asiceWithXAdESManifestParserNamespacesRegistered
	return &ASiCEWithXAdESManifestParser{
		signatureDocument: signatureDocument,
		manifestDocument:  manifestDocument,
	}
}

// Manifest returns a parsed model.ManifestFile. Ports getManifest().
func (p *ASiCEWithXAdESManifestParser) Manifest() *model.ManifestFile {
	manifest := model.NewManifestFile()
	manifest.SetDocument(p.manifestDocument)
	if p.signatureDocument != nil {
		manifest.SetSignatureFilename(p.signatureDocument.Name())
	}
	manifest.SetEntries(p.getEntries())
	return manifest
}

// getEntries ports the private getEntries(). Logging (LOG.warn) is dropped per PORTING.md; the
// caught-exception branch returns the same empty result as its logged counterpart.
func (p *ASiCEWithXAdESManifestParser) getEntries() (result []*model.ManifestEntry) {
	if !xmlutils.DomUtilsIsDOM(p.manifestDocument) {
		return []*model.ManifestEntry{}
	}

	result = []*model.ManifestEntry{}
	defer func() {
		if recover() != nil {
			result = []*model.ManifestEntry{}
		}
	}()

	manifestDom, err := xmlutils.DomUtilsBuildDOMFromDocument(p.manifestDocument)
	if err != nil {
		return []*model.ManifestEntry{}
	}
	manifestNamespace := p.getManifestNamespace(manifestDom)

	nodeList, err := xmlutils.XPathUtilsGetNodeList(manifestDom, ManifestPath_FILE_ENTRY_PATH)
	if err != nil {
		return []*model.ManifestEntry{}
	}

	entries := make([]*model.ManifestEntry, 0, len(nodeList))
	for _, fileEntryElement := range nodeList {
		fullPathValue := fileEntryElement.AttrValue("", ManifestPathGetFullPathAttribute(manifestNamespace))
		if !p.isFolder(fullPathValue) {
			manifestEntry := model.NewManifestEntry()
			manifestEntry.SetUri(fullPathValue)
			manifestEntry.SetMimeType(p.getMimeType(fileEntryElement, manifestNamespace))
			entries = append(entries, manifestEntry)
		}
	}
	return entries
}

// getManifestNamespace ports the private getManifestNamespace(Document).
func (p *ASiCEWithXAdESManifestParser) getManifestNamespace(manifestDom *xmldom.Node) *common.DSSNamespace {
	manifestNamespace := xmlutils.DomUtilsBrowseRecursivelyForNamespaceWithUri(manifestDom.DocumentElement(), ManifestNS.Uri())
	if manifestNamespace != nil {
		return manifestNamespace
	}
	return ManifestNS
}

// getMimeType ports the private getMimeType(Element, DSSNamespace).
func (p *ASiCEWithXAdESManifestParser) getMimeType(fileEntryElement *xmldom.Node, manifestNamespace *common.DSSNamespace) enumerations.MimeType {
	mediaType := fileEntryElement.AttrValue("", ManifestPathGetMediaTypeAttribute(manifestNamespace))
	if strings.TrimSpace(mediaType) != "" {
		return enumerations.MimeTypeFromMimeTypeString(mediaType)
	}
	return nil
}

// isFolder ports the private isFolder(String).
func (p *ASiCEWithXAdESManifestParser) isFolder(fullpathValue string) bool {
	return strings.HasSuffix(fullpathValue, "/")
}
