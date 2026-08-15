// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/ASiCFilenameFactory.java
// (DSS 6.5.RC1).
//
// java.io.Serializable is dropped silently (no Go counterpart).
package asic

// ASiCFilenameFactory provides filenames for newly created ZIP-entries during a signature
// creation or extension for ASiC containers.
//
// NOTE: Names of signature or manifest files shall be defined with leading "META-INF/" string,
// specifying the target folder of the signature file within a container.
//
// As the same factory is used for ASiC-S and ASiC-E container types, it shall implement logic for
// both container types, when applicable. The type of the container can be obtained from
// asicContent.ContainerType().
type ASiCFilenameFactory interface {
	// SignatureFilename returns a filename for a signature file to be created. Port of
	// getSignatureFilename(ASiCContent).
	SignatureFilename(asicContent *ASiCContent) string

	// ManifestFilename returns a filename of a manifest file to be created. Port of
	// getManifestFilename(ASiCContent).
	ManifestFilename(asicContent *ASiCContent) string

	// DataPackageFilename returns a name of a ZIP archive containing signed documents in case
	// of an ASiC-E signature, when applicable. Port of getDataPackageFilename(ASiCContent).
	DataPackageFilename(asicContent *ASiCContent) string
}
