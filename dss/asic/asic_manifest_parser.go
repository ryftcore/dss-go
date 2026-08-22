// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/validation/ASiCManifestParser.java (DSS 6.5.RC1).
package asic

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
	xmlcommon "github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// asicManifestParserNamespacesRegistered mirrors the Java static initializer block, which
// registers the XMLDSig and ASiCManifest namespaces once at class-load time.
var asicManifestParserNamespacesRegistered = func() bool {
	xmlutils.XPathUtilsRegisterNamespace(xmlcommon.XMLDSigNS)
	xmlutils.XPathUtilsRegisterNamespace(ASiCManifestNS)
	return true
}()

// ASiCManifestParserGetManifestFile parses and converts a DSSDocument to a ManifestFile. Ports
// the static getManifestFile(DSSDocument). Returns nil when the document isn't a valid ASiC
// manifest, matching Java's null return.
func ManifestParserGetManifestFile(manifestDocument model.DSSDocument) *model.ManifestFile {
	_ = asicManifestParserNamespacesRegistered
	root := asicManifestParserGetManifestRootElement(manifestDocument)
	if root == nil {
		return nil
	}
	manifest := model.NewManifestFile()
	manifest.SetDocument(manifestDocument)
	manifest.SetSignatureFilename(asicManifestParserGetLinkedSignatureName(root))
	manifest.SetEntries(asicManifestParserParseManifestEntries(root))
	manifest.SetManifestType(asicManifestParserGetManifestType(manifestDocument.Name(), root))
	return manifest
}

// ManifestParserGetLinkedManifest returns the relative manifest for the given signature
// name. Ports the static getLinkedManifest(List, String).
func ManifestParserGetLinkedManifest(manifestDocuments []model.DSSDocument, signatureName string) model.DSSDocument {
	for _, manifest := range manifestDocuments {
		manifestRoot := asicManifestParserGetManifestRootElement(manifest)
		if manifestRoot != nil {
			linkedSignatureName := spi.DSSUtilsDecodeURI(asicManifestParserGetLinkedSignatureName(manifestRoot))
			if signatureName == linkedSignatureName {
				return manifest
			}
		}
	}
	return nil
}

// asicManifestParserGetManifestRootElement ports the private static
// getManifestRootElement(DSSDocument). Logging (LOG.warn) is dropped per PORTING.md
// (slf4j dropped).
func asicManifestParserGetManifestRootElement(manifestDocument model.DSSDocument) *xmldom.Node {
	if !xmlutils.DomUtilsIsDOM(manifestDocument) {
		return nil
	}
	manifestDom, err := xmlutils.DomUtilsBuildDOMFromDocument(manifestDocument)
	if err != nil {
		return nil
	}
	element, err := xmlutils.XPathUtilsGetElement(manifestDom, ASiCManifestPathASiCManifestPath)
	if err != nil {
		return nil
	}
	return element
}

// asicManifestParserGetLinkedSignatureName ports the private static
// getLinkedSignatureName(Element).
func asicManifestParserGetLinkedSignatureName(root *xmldom.Node) string {
	value, err := xmlutils.XPathUtilsGetValue(root, ASiCManifestPathSigReferenceURIPath)
	if err != nil {
		return ""
	}
	return value
}

// asicManifestParserGetMimeType ports the private static getMimeType(Element). Java's
// try/catch around MimeType.fromMimeTypeString(String) exists to swallow a DSSException on an
// unrecognized string; the ported MimeTypeFromMimeTypeString never errors (it falls back to
// MimeTypeEnumBinary), so there is no error path left to catch here - a judgment call
// documented since it is an observable behavior difference (Java returns nil for a malformed
// mimetype string in some MimeType loader implementations that choose to throw; the Go port's
// registered loaders do not).
func asicManifestParserGetMimeType(element *xmldom.Node) enumerations.MimeType {
	mimeTypeString := element.AttrValue("", ManifestAttributeMIMEType.AttributeName())
	if utils.IsStringNotBlank(mimeTypeString) {
		return enumerations.MimeTypeFromMimeTypeString(mimeTypeString)
	}
	return nil
}

// asicManifestParserGetDigestAlgorithm ports the private static
// getDigestAlgorithm(Element).
func asicManifestParserGetDigestAlgorithm(dataObjectReference *xmldom.Node) enumerations.DigestAlgorithm {
	value, err := xmlutils.XPathUtilsGetValue(dataObjectReference, xmlcommon.XMLDSigPathDigestMethodAlgorithmPath)
	if err != nil {
		return ""
	}
	algorithm, err := enumerations.DigestAlgorithmForXML(value)
	if err != nil {
		return ""
	}
	return algorithm
}

// asicManifestParserGetDigestValue ports the private static getDigestValue(Element).
func asicManifestParserGetDigestValue(dataObjectReference *xmldom.Node) []byte {
	digestValueElement, err := xmlutils.XPathUtilsGetElement(dataObjectReference, xmlcommon.XMLDSigPathDigestValuePath)
	if err != nil || digestValueElement == nil {
		return nil
	}
	digestText := digestValueElement.TextContent()
	if utils.IsBase64Encoded(digestText) {
		return utils.FromBase64(digestText)
	}
	return nil
}

// asicManifestParserGetManifestType ports the private static getManifestType(String, Element).
//
// Cross-chunk assumption (ZIPCORE): ASiCUtils exposes IsArchiveManifest,
// IsEvidenceRecordManifest and IsManifest package functions taking a filename string.
func asicManifestParserGetManifestType(manifestFilename string, root *xmldom.Node) enumerations.ASiCManifestTypeEnum {
	if UtilsIsArchiveManifest(manifestFilename) {
		return enumerations.ASiCManifestTypeEnumArchiveManifest
	} else if UtilsIsEvidenceRecordManifest(manifestFilename) {
		return enumerations.ASiCManifestTypeEnumEvidenceRecord
	} else if UtilsIsManifest(manifestFilename) {
		sigReference, err := xmlutils.XPathUtilsGetElement(root, ASiCManifestPathSigReferencePath)
		if err == nil && sigReference != nil {
			mimeType := asicManifestParserGetMimeType(sigReference)
			if mimeType == enumerations.MimeTypeEnumTST {
				return enumerations.ASiCManifestTypeEnumTimestamp
			}
			return enumerations.ASiCManifestTypeEnumSignature
		}
	}
	return ""
}

// asicManifestParserParseManifestEntries ports the private static
// parseManifestEntries(Element).
func asicManifestParserParseManifestEntries(root *xmldom.Node) []*model.ManifestEntry {
	entries := make([]*model.ManifestEntry, 0)
	dataObjectReferences, err := xmlutils.XPathUtilsGetNodeList(root, ASiCManifestPathDataObjectReferencePath)
	if err != nil || len(dataObjectReferences) == 0 {
		return entries
	}
	for _, dataObjectReference := range dataObjectReferences {
		entry := model.NewManifestEntry()
		entry.SetUri(spi.DSSUtilsDecodeURI(dataObjectReference.AttrValue("", ManifestAttributeURI.AttributeName())))
		entry.SetMimeType(asicManifestParserGetMimeType(dataObjectReference))

		digestAlgorithm := asicManifestParserGetDigestAlgorithm(dataObjectReference)
		digestValueBinary := asicManifestParserGetDigestValue(dataObjectReference)
		if digestAlgorithm != "" && digestValueBinary != nil {
			entry.SetDigest(model.NewDigest(digestAlgorithm, digestValueBinary))
		}

		attribute := dataObjectReference.AttrValue("", ManifestAttributeRootFile.AttributeName())
		if utils.AreStringsEqualIgnoreCase("true", attribute) {
			entry.SetRootfile(true)
		}

		entries = append(entries, entry)
	}
	return entries
}
