// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/definition/ManifestElement.java
// (DSS 6.5.RC1).
package xades

import "github.com/utain/esig/dss/xml/common"

// ManifestElement is the Manifest element enumeration. Ports the Java enum per PORTING.md's
// enum convention: a typed string whose value is the Java name(), with the wire tag name held
// in a lookup table.
//
// KAT-verified against a Java oracle run over dss-asic-xades 6.5.RC1 + dss-xml-common 6.5.RC1
// (java.lang.reflect over ManifestElement.class.getEnumConstants(), each constant's
// getTagName()/getURI()/getNamespace().getPrefix() printed) - see manifest_element_test.go.
type ManifestElement string

const (
	// ManifestElement_MANIFEST is the "manifest" element.
	ManifestElement_MANIFEST ManifestElement = "MANIFEST"

	// ManifestElement_FILE_ENTRY is the "file-entry" element.
	ManifestElement_FILE_ENTRY ManifestElement = "FILE_ENTRY"
)

// manifestElementTagNames maps each constant to its wire tag name (getTagName()).
var manifestElementTagNames = map[ManifestElement]string{
	ManifestElement_MANIFEST:   "manifest",
	ManifestElement_FILE_ENTRY: "file-entry",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e ManifestElement) TagName() string {
	return manifestElementTagNames[e]
}

// Namespace implements common.DSSElement. Every ManifestElement carries ManifestNS. Ports
// getNamespace().
func (e ManifestElement) Namespace() *common.DSSNamespace {
	return ManifestNS
}

// URI implements common.DSSElement. Ports getURI().
func (e ManifestElement) URI() string {
	return ManifestNS.Uri()
}

// IsSameTagName implements common.DSSElement. Ports isSameTagName(String).
func (e ManifestElement) IsSameTagName(value string) bool {
	return e.TagName() == value
}

var _ common.DSSElement = ManifestElement("")
