// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/definition/ASiCManifestPath.java (DSS 6.5.RC1).
package asic

import "github.com/ryftcore/dss-go/dss/xml/common"

// Contains paths for an ASiC scheme. Ports the Java class ASiCManifestPath (a concrete
// AbstractPath subclass with no other subclasses) as package-level vars, per the
// xmldsig_path.go precedent.
var (
	// ASiCManifestPathASiCManifestPath = "./asic:ASiCManifest"
	ASiCManifestPathASiCManifestPath = common.FromCurrentPosition(ManifestElementASiCManifest)

	// ASiCManifestPathDataObjectReferencePath = "./asic:DataObjectReference"
	ASiCManifestPathDataObjectReferencePath = common.FromCurrentPosition(ManifestElementDataObjectReference)

	// ASiCManifestPathSigReferencePath = "./asic:SigReference"
	ASiCManifestPathSigReferencePath = common.FromCurrentPosition(ManifestElementSigReference)

	// ASiCManifestPathSigReferenceURIPath = "./asic:SigReference@URI"
	ASiCManifestPathSigReferenceURIPath = common.FromCurrentPositionAttribute(ManifestElementSigReference, ManifestAttributeURI)
)
