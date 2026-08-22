// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/definition/ASiCManifestElement.java (DSS 6.5.RC1).
package asic

import "github.com/ryftcore/dss-go/dss/xml/common"

// ASiCManifestElement is an ASiC XSD element. Ports the Java enum per PORTING.md's enum
// convention: a typed string whose value is the Java name(), with the wire tag name held in a
// lookup table.
type ManifestElement string

const (
	ASiCManifestElementXAdESSignatures               ManifestElement = "XADES_SIGNATURES"
	ASiCManifestElementASiCManifest                  ManifestElement = "ASIC_MANIFEST"
	ASiCManifestElementSigReference                  ManifestElement = "SIG_REFERENCE"
	ASiCManifestElementExtension                     ManifestElement = "EXTENSION"
	ASiCManifestElementDataObjectReference           ManifestElement = "DATA_OBJECT_REFERENCE"
	ASiCManifestElementASiCManifestExtensions        ManifestElement = "ASIC_MANIFEST_EXTENSIONS"
	ASiCManifestElementDataObjectReferenceExtensions ManifestElement = "DATA_OBJECT_REFERENCE_EXTENSIONS"
)

// asicManifestElementTagNames maps each constant to its wire tag name (getTagName()).
var asicManifestElementTagNames = map[ManifestElement]string{
	ASiCManifestElementXAdESSignatures:               "XAdESSignatures",
	ASiCManifestElementASiCManifest:                  "ASiCManifest",
	ASiCManifestElementSigReference:                  "SigReference",
	ASiCManifestElementExtension:                     "Extension",
	ASiCManifestElementDataObjectReference:           "DataObjectReference",
	ASiCManifestElementASiCManifestExtensions:        "ASiCManifestExtensions",
	ASiCManifestElementDataObjectReferenceExtensions: "DataObjectReferenceExtensions",
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
