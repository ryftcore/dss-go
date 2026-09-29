// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/validation/ASiCEWithXAdESManifestParser.java (DSS 6.5.RC1).
//
// Package flattening: Java's eu.europa.esig.dss.asic.xades.validation lands in this
// same Go package (dss/asic/xades).
package xades

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
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
//
// As upstream (which only logs a warning), a manifest that cannot be read - not XML, malformed,
// truncated - yields a ManifestFile without entries, so an empty entry list means either "no
// entries" or "unreadable manifest"; callers must not treat it as a well-formedness signal.
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

	nodeList, err := xmlutils.XPathUtilsGetNodeList(manifestDom, ManifestPathFileEntryPath)
	if err != nil {
		return []*model.ManifestEntry{}
	}

	entries := make([]*model.ManifestEntry, 0, len(nodeList))
	for _, fileEntryElement := range nodeList {
		fullPathValue := asiceWithXAdESManifestParserAttrValue(
			fileEntryElement, ManifestPathGetFullPathAttribute(manifestNamespace))
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
	mediaType := asiceWithXAdESManifestParserAttrValue(
		fileEntryElement, ManifestPathGetMediaTypeAttribute(manifestNamespace))
	if strings.TrimSpace(mediaType) != "" {
		return enumerations.MimeTypeFromMimeTypeString(mediaType)
	}
	return nil
}

// asiceWithXAdESManifestParserAttrValue returns element's attribute with the given QUALIFIED
// name ("prefix:local"), or "" when it has none.
//
// Upstream reads these two attributes with org.w3c.dom.Element#getAttribute(String), which
// matches an attribute's qualified name (its nodeName) rather than its expanded
// {namespace-uri}local-name pair - and ASiCEWithXAdESManifestParser leans on that deliberately,
// building the name from whichever prefix the document itself binds to the manifest namespace
// (getManifestNamespace + ManifestPath.getFullPathAttribute produce e.g. "manifest:full-path").
//
// xmldom's Attr/AttrValue match on the expanded name, so passing them a qualified name with an
// empty namespace - as this file previously did - never matches: every entry came back with an
// empty URI and no mime type, which also defeated isFolder() and let the "/" root entry through.
// The unprefixed XMLDSig and ASiCManifest attributes read elsewhere in the port are unaffected,
// since for them qualified and expanded lookups coincide.
func asiceWithXAdESManifestParserAttrValue(element *xmldom.Node, qualifiedName string) string {
	if element == nil {
		return ""
	}
	for _, attr := range element.Attrs {
		if attr.Name.QName() == qualifiedName {
			return attr.Value
		}
	}
	return ""
}

// isFolder ports the private isFolder(String).
func (p *ASiCEWithXAdESManifestParser) isFolder(fullpathValue string) bool {
	return strings.HasSuffix(fullpathValue, "/")
}
