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
	// ManifestAttribute_VERSION is the "version" attribute.
	ManifestAttribute_VERSION ManifestAttribute = "VERSION"

	// ManifestAttribute_FULL_PATH is the "full-path" attribute.
	ManifestAttribute_FULL_PATH ManifestAttribute = "FULL_PATH"

	// ManifestAttribute_MEDIA_TYPE is the "media-type" attribute.
	ManifestAttribute_MEDIA_TYPE ManifestAttribute = "MEDIA_TYPE"
)

// manifestAttributeNames maps each constant to its wire attribute name (getAttributeName()).
var manifestAttributeNames = map[ManifestAttribute]string{
	ManifestAttribute_VERSION:    "version",
	ManifestAttribute_FULL_PATH:  "full-path",
	ManifestAttribute_MEDIA_TYPE: "media-type",
}

// AttributeName implements common.DSSAttribute. Ports getAttributeName().
func (a ManifestAttribute) AttributeName() string {
	return manifestAttributeNames[a]
}

var _ common.DSSAttribute = ManifestAttribute("")
