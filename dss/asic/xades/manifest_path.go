// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/definition/ManifestPath.java
// (DSS 6.5.RC1).
//
// Ports the Java class ManifestPath (a concrete AbstractPath subclass with no other
// subclasses) as a package-level var/functions, per the xmldsig_path.go / asic_manifest_path.go
// precedent.
package xades

import (
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xml/common"
)

// ManifestPath_FILE_ENTRY_PATH is the XPath expression to return manifest:file-entry entries
// from the current position.
//
// KAT-verified against a Java oracle run over dss-asic-xades 6.5.RC1 + dss-xml-common 6.5.RC1
// (java.lang.reflect over ManifestPath.FILE_ENTRY_PATH.getQueryString()) - see
// manifest_path_test.go, which asserts the value below against the same string the oracle
// printed. NOTE: the upstream field's javadoc comment claims the value is
// "//manifest/file-entry" - that is a stale doc comment; the actual code path
// (fromCurrentPosition, not all) and the oracle both agree on the value below.
//
// FILE_ENTRY_PATH.getQueryString() = "./manifest:manifest/manifest:file-entry"
var ManifestPath_FILE_ENTRY_PATH = common.FromCurrentPosition(ManifestElement_MANIFEST, ManifestElement_FILE_ENTRY)

// ManifestPathGetFullPathAttribute returns "manifest:full-path" with the given
// manifestNamespace's prefix. Ports the static getFullPathAttribute(DSSNamespace).
func ManifestPathGetFullPathAttribute(manifestNamespace *common.DSSNamespace) string {
	return manifestPathAddPrefixIfNeeded(ManifestAttribute_FULL_PATH.AttributeName(), manifestNamespace)
}

// ManifestPathGetMediaTypeAttribute returns "manifest:media-type" with the given
// manifestNamespace's prefix. Ports the static getMediaTypeAttribute(DSSNamespace).
func ManifestPathGetMediaTypeAttribute(manifestNamespace *common.DSSNamespace) string {
	return manifestPathAddPrefixIfNeeded(ManifestAttribute_MEDIA_TYPE.AttributeName(), manifestNamespace)
}

// manifestPathAddPrefixIfNeeded ports the private static addPrefixIfNeeded(String,
// DSSNamespace).
func manifestPathAddPrefixIfNeeded(attributeName string, manifestNamespace *common.DSSNamespace) string {
	if utils.IsStringEmpty(manifestNamespace.Prefix()) {
		return attributeName
	}
	return manifestNamespace.Prefix() + ":" + attributeName
}
