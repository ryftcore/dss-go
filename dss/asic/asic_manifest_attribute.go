// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/definition/ASiCManifestAttribute.java (DSS 6.5.RC1).
package asic

import "github.com/ryftcore/dss-go/dss/xml/common"

// ManifestAttribute represents a collection of attributes defined in ASiC XSD schema.
// Ports the Java enum per PORTING.md's enum convention: a typed string whose value is the
// Java name(), with the wire attribute name held in a lookup table.
type ManifestAttribute string

const (
	ASiCManifestAttributeURI      ManifestAttribute = "URI"
	ASiCManifestAttributeMIMEType ManifestAttribute = "MIME_TYPE"
	ASiCManifestAttributeRootFile ManifestAttribute = "ROOTFILE"
	ASiCManifestAttributeCritical ManifestAttribute = "CRITICAL"
)

// asicManifestAttributeNames maps each constant to its wire attribute name (getAttributeName()).
var asicManifestAttributeNames = map[ManifestAttribute]string{
	ASiCManifestAttributeURI:      "URI",
	ASiCManifestAttributeMIMEType: "MimeType",
	ASiCManifestAttributeRootFile: "Rootfile",
	ASiCManifestAttributeCritical: "Critical",
}

// AttributeName implements common.DSSAttribute. Ports getAttributeName().
func (a ManifestAttribute) AttributeName() string {
	return asicManifestAttributeNames[a]
}

var _ common.DSSAttribute = ManifestAttribute("")
