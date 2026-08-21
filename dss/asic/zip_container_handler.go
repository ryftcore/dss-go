// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/ZipContainerHandler.java
// (DSS 6.5.RC1).
package asic

import (
	"time"

	"github.com/ryftcore/dss-go/dss/model"
)

// ZipContainerHandler provides utilities for data extraction/creation of ZIP-archives.
//
// Java's methods declare no checked exceptions and signal failure through
// DSSException/IllegalInputException (RuntimeExceptions); per PORTING.md those become returned
// errors, so each method gains an error result.
type ZipContainerHandler interface {
	// ExtractContainerContent extracts a list of DSSDocument from the given ZIP-archive. Port
	// of extractContainerContent(DSSDocument).
	ExtractContainerContent(zipArchive model.DSSDocument) ([]model.DSSDocument, error)

	// ExtractEntryNames returns a list of ZIP archive entry names. Port of
	// extractEntryNames(DSSDocument).
	ExtractEntryNames(zipArchive model.DSSDocument) ([]string, error)

	// CreateZipArchive creates a ZIP-Archive with the given containerEntries.
	//
	// creationTime is optional: it defines the time of an archive creation and will be set for
	// all embedded files; the zero Time (Java null) means the local current time will be used.
	// zipComment is optional.
	//
	// Port of createZipArchive(List<DSSDocument>, Date, String).
	CreateZipArchive(containerEntries []model.DSSDocument, creationTime time.Time, zipComment string) (model.DSSDocument, error)
}
