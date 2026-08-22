// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/definition/ManifestAttribute.java
// (DSS 6.5.RC1).
package xades

import "github.com/ryftcore/dss-go/dss/xml/common"

// ManifestAttribute is the enumeration of Manifest attributes. Ports the Java enum per
// PORTING.md's enum convention: a typed string whose value is the Java name(), with the wire
// attribute name held in a lookup table.
//
// KAT-verified against a Java oracle run over dss-asic-xades 6.5.RC1 (java.lang.reflect over
// ManifestAttribute.class.getEnumConstants(), each constant's getAttributeName() printed) - see
// manifest_attribute_test.go.
type ManifestAttribute string

const (
	// ManifestAttributeVersion is the "version" attribute.
	ManifestAttributeVersion ManifestAttribute = "VERSION"

	// ManifestAttributeFullPath is the "full-path" attribute.
	ManifestAttributeFullPath ManifestAttribute = "FULL_PATH"

	// ManifestAttributeMediaType is the "media-type" attribute.
	ManifestAttributeMediaType ManifestAttribute = "MEDIA_TYPE"
)

// manifestAttributeNames maps each constant to its wire attribute name (getAttributeName()).
var manifestAttributeNames = map[ManifestAttribute]string{
	ManifestAttributeVersion:   "version",
	ManifestAttributeFullPath:  "full-path",
	ManifestAttributeMediaType: "media-type",
}

// AttributeName implements common.DSSAttribute. Ports getAttributeName().
func (a ManifestAttribute) AttributeName() string {
	return manifestAttributeNames[a]
}

var _ common.DSSAttribute = ManifestAttribute("")
