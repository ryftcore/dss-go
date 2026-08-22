// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/definition/ASiCManifestElement.java (DSS 6.5.RC1).
package asic

import "github.com/ryftcore/dss-go/dss/xml/common"

// ASiCManifestElement is an ASiC XSD element. Ports the Java enum per PORTING.md's enum
// convention: a typed string whose value is the Java name(), with the wire tag name held in a
// lookup table.
type ManifestElement string

const (
	ManifestElementXAdESSignatures               ManifestElement = "XADES_SIGNATURES"
	ManifestElementASiCManifest                  ManifestElement = "ASIC_MANIFEST"
	ManifestElementSigReference                  ManifestElement = "SIG_REFERENCE"
	ManifestElementExtension                     ManifestElement = "EXTENSION"
	ManifestElementDataObjectReference           ManifestElement = "DATA_OBJECT_REFERENCE"
	ManifestElementASiCManifestExtensions        ManifestElement = "ASIC_MANIFEST_EXTENSIONS"
	ManifestElementDataObjectReferenceExtensions ManifestElement = "DATA_OBJECT_REFERENCE_EXTENSIONS"
)

// asicManifestElementTagNames maps each constant to its wire tag name (getTagName()).
var asicManifestElementTagNames = map[ManifestElement]string{
	ManifestElementXAdESSignatures:               "XAdESSignatures",
	ManifestElementASiCManifest:                  "ASiCManifest",
	ManifestElementSigReference:                  "SigReference",
	ManifestElementExtension:                     "Extension",
	ManifestElementDataObjectReference:           "DataObjectReference",
	ManifestElementASiCManifestExtensions:        "ASiCManifestExtensions",
	ManifestElementDataObjectReferenceExtensions: "DataObjectReferenceExtensions",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e ManifestElement) TagName() string {
	return asicManifestElementTagNames[e]
}

// Namespace implements common.DSSElement. Every ManifestElement carries ASiCManifestNS.
// Ports getNamespace().
func (e ManifestElement) Namespace() *common.DSSNamespace {
	return ASiCManifestNS
}

// URI implements common.DSSElement. Ports getURI().
func (e ManifestElement) URI() string {
	return ASiCManifestNS.Uri()
}

// IsSameTagName implements common.DSSElement. Ports isSameTagName(String).
func (e ManifestElement) IsSameTagName(value string) bool {
	return e.TagName() == value
}

var _ common.DSSElement = ManifestElement("")
