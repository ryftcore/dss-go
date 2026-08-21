// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/definition/ASiCManifestAttribute.java (DSS 6.5.RC1).
package asic

import "github.com/ryftcore/dss-go/dss/xml/common"

// ASiCManifestAttribute represents a collection of attributes defined in ASiC XSD schema.
// Ports the Java enum per PORTING.md's enum convention: a typed string whose value is the
// Java name(), with the wire attribute name held in a lookup table.
type ASiCManifestAttribute string

const (
	ASiCManifestAttribute_URI       ASiCManifestAttribute = "URI"
	ASiCManifestAttribute_MIME_TYPE ASiCManifestAttribute = "MIME_TYPE"
	ASiCManifestAttribute_ROOTFILE  ASiCManifestAttribute = "ROOTFILE"
	ASiCManifestAttribute_CRITICAL  ASiCManifestAttribute = "CRITICAL"
)

// asicManifestAttributeNames maps each constant to its wire attribute name (getAttributeName()).
var asicManifestAttributeNames = map[ASiCManifestAttribute]string{
	ASiCManifestAttribute_URI:       "URI",
	ASiCManifestAttribute_MIME_TYPE: "MimeType",
	ASiCManifestAttribute_ROOTFILE:  "Rootfile",
	ASiCManifestAttribute_CRITICAL:  "Critical",
}

// AttributeName implements common.DSSAttribute. Ports getAttributeName().
func (a ASiCManifestAttribute) AttributeName() string {
	return asicManifestAttributeNames[a]
}

var _ common.DSSAttribute = ASiCManifestAttribute("")
