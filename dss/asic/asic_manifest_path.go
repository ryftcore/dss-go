// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/definition/ASiCManifestPath.java (DSS 6.5.RC1).
package asic

import "github.com/utain/esig/dss/xml/common"

// Contains paths for an ASiC scheme. Ports the Java class ASiCManifestPath (a concrete
// AbstractPath subclass with no other subclasses) as package-level vars, per the
// xmldsig_path.go precedent.
var (
	// ASiCManifestPath_ASIC_MANIFEST_PATH = "./asic:ASiCManifest"
	ASiCManifestPath_ASIC_MANIFEST_PATH = common.FromCurrentPosition(ASiCManifestElement_ASIC_MANIFEST)

	// ASiCManifestPath_DATA_OBJECT_REFERENCE_PATH = "./asic:DataObjectReference"
	ASiCManifestPath_DATA_OBJECT_REFERENCE_PATH = common.FromCurrentPosition(ASiCManifestElement_DATA_OBJECT_REFERENCE)

	// ASiCManifestPath_SIG_REFERENCE_PATH = "./asic:SigReference"
	ASiCManifestPath_SIG_REFERENCE_PATH = common.FromCurrentPosition(ASiCManifestElement_SIG_REFERENCE)

	// ASiCManifestPath_SIG_REFERENCE_URI_PATH = "./asic:SigReference@URI"
	ASiCManifestPath_SIG_REFERENCE_URI_PATH = common.FromCurrentPositionAttribute(ASiCManifestElement_SIG_REFERENCE, ASiCManifestAttribute_URI)
)
