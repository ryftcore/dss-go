// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/definition/ASiCManifestElement.java (DSS 6.5.RC1).
package asic

import "github.com/ryftcore/dss-go/dss/xml/common"

// ASiCManifestElement is an ASiC XSD element. Ports the Java enum per PORTING.md's enum
// convention: a typed string whose value is the Java name(), with the wire tag name held in a
// lookup table.
type ASiCManifestElement string

const (
	ASiCManifestElement_XADES_SIGNATURES                 ASiCManifestElement = "XADES_SIGNATURES"
	ASiCManifestElement_ASIC_MANIFEST                    ASiCManifestElement = "ASIC_MANIFEST"
	ASiCManifestElement_SIG_REFERENCE                    ASiCManifestElement = "SIG_REFERENCE"
	ASiCManifestElement_EXTENSION                        ASiCManifestElement = "EXTENSION"
	ASiCManifestElement_DATA_OBJECT_REFERENCE            ASiCManifestElement = "DATA_OBJECT_REFERENCE"
	ASiCManifestElement_ASIC_MANIFEST_EXTENSIONS         ASiCManifestElement = "ASIC_MANIFEST_EXTENSIONS"
	ASiCManifestElement_DATA_OBJECT_REFERENCE_EXTENSIONS ASiCManifestElement = "DATA_OBJECT_REFERENCE_EXTENSIONS"
)

// asicManifestElementTagNames maps each constant to its wire tag name (getTagName()).
var asicManifestElementTagNames = map[ASiCManifestElement]string{
	ASiCManifestElement_XADES_SIGNATURES:                 "XAdESSignatures",
	ASiCManifestElement_ASIC_MANIFEST:                    "ASiCManifest",
	ASiCManifestElement_SIG_REFERENCE:                    "SigReference",
	ASiCManifestElement_EXTENSION:                        "Extension",
	ASiCManifestElement_DATA_OBJECT_REFERENCE:            "DataObjectReference",
	ASiCManifestElement_ASIC_MANIFEST_EXTENSIONS:         "ASiCManifestExtensions",
	ASiCManifestElement_DATA_OBJECT_REFERENCE_EXTENSIONS: "DataObjectReferenceExtensions",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e ASiCManifestElement) TagName() string {
	return asicManifestElementTagNames[e]
}

// Namespace implements common.DSSElement. Every ASiCManifestElement carries ASiCManifestNS.
// Ports getNamespace().
func (e ASiCManifestElement) Namespace() *common.DSSNamespace {
	return ASiCManifestNS
}

// URI implements common.DSSElement. Ports getURI().
func (e ASiCManifestElement) URI() string {
	return ASiCManifestNS.Uri()
}

// IsSameTagName implements common.DSSElement. Ports isSameTagName(String).
func (e ASiCManifestElement) IsSameTagName(value string) bool {
	return e.TagName() == value
}

var _ common.DSSElement = ASiCManifestElement("")
