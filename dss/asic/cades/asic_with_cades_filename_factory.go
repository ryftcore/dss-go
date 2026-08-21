// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/ASiCWithCAdESFilenameFactory.java (DSS 6.5.RC1).
package cades

import "github.com/ryftcore/dss-go/dss/asic"

// ASiCWithCAdESFilenameFactory is used to provide filenames for newly created ZIP-entries during
// a signature creation or extension for an ASiC with CAdES container.
//
// NOTE: Names of signature, timestamp or manifest files shall be defined with leading "META-INF/"
// string, specifying the target folder of the signature file within a container.
//
// As the same factory is used for ASiC-S and ASiC-E container types, it shall implement logic for
// both container types, when applicable. The type of the container can be obtained from
// asicContent.ContainerType().
type ASiCWithCAdESFilenameFactory interface {
	asic.ASiCFilenameFactory
	asic.ASiCEvidenceRecordFilenameFactory

	// TimestampFilename returns a filename for a timestamp file to be created. Port of
	// getTimestampFilename(ASiCContent).
	TimestampFilename(asicContent *asic.ASiCContent) string

	// ArchiveManifestFilename returns a new filename of an archive manifest file to be moved.
	//
	// NOTE: ASiC-E with CAdES shall always create a new archive manifest with the same name,
	// while moving the last existing archive manifest.
	//
	// Port of getArchiveManifestFilename(ASiCContent).
	ArchiveManifestFilename(asicContent *asic.ASiCContent) string
}
