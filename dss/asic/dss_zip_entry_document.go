// Ported from
// dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/DSSZipEntryDocument.java
// (DSS 6.5.RC1).
package asic

import "github.com/utain/esig/dss/model"

// DSSZipEntryDocument is a DSSDocument carrying metadata for a ZIP-container entry.
type DSSZipEntryDocument interface {
	model.DSSDocument

	// ZipEntry returns the ZIP entry wrapper containing metadata about a file within a
	// ZIP-container. Port of getZipEntry().
	ZipEntry() *DSSZipEntry
}
