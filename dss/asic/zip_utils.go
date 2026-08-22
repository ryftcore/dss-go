// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/ZipUtils.java
// (DSS 6.5.RC1).
//
// SINGLETON: Java's lazily-created, mutable ZipUtils singleton becomes a package-level pointer
// guarded by a mutex. Java's own lazy initialization is not thread-safe either, but a Go data race
// is undefined behaviour rather than a benign double-instantiation, so the lock is added.
//
// WRITE-SIDE DEVIATIONS from java.util.zip (read-side deviations are documented in
// secure_container_handler.go's header):
//
//  1. DEFLATE OUTPUT BYTES. Go's compress/flate and java.util.zip.Deflater produce different
//     compressed bytes for the same input at the same level. This is accepted: no ASiC signature
//     covers container bytes - the signed payloads are the entry CONTENTS (which round-trip
//     byte-exactly) and the manifest XML documents.
//
//  2. STORED vs DEFLATED HEADER LAYOUT. archive/zip's Writer.CreateHeader always sets the
//     data-descriptor flag and zeroes crc/sizes in the local header. java.util.zip.ZipOutputStream
//     does that for DEFLATED entries only; for STORED entries it writes crc and sizes up front and
//     no descriptor - which is what EN 319 162-1 A.1 readers expect of the "mimetype" entry.
//     SecureContainerHandler.BuildZip therefore routes STORED entries through Writer.CreateRaw,
//     reproducing the Java layout byte for byte.
//
//  3. TIMESTAMP EXTRA FIELDS. CreateHeader appends a synthesized Info-ZIP extended-timestamp
//     (0x5455) extra field whenever FileHeader.Modified is set. Java appends one only when the
//     entry carries a FileTime-valued mtime/atime/ctime, which ZipEntry#setTime does not produce.
//     SecureContainerHandler.ensureTime therefore writes the MS-DOS date/time fields directly and
//     leaves Modified unset, so no extra field is synthesized and the entry's original `extra`
//     bytes are the only ones written.
//
//  4. GENERAL PURPOSE FLAG BIT 11. Java sets the UTF-8 name flag on every entry (its ZipCoder is
//     UTF-8 by default); this port follows archive/zip and sets it only for entries whose name or
//     comment is not pure ASCII. The two encodings are indistinguishable for ASCII names, which is
//     every filename ASiCUtils generates.
//
//  5. DIRECTORY ENTRIES. archive/zip's Writer forces a name ending in "/" to STORED with zero
//     sizes; java.util.zip.ZipOutputStream keeps whatever method the entry asked for, emitting a
//     2-byte empty deflate stream for a DEFLATED one. Across the 188-container fixture corpus this
//     changes only the framing of content-free entries (Content.folders), which no manifest
//     ever digest-references, and preserving it would require pre-compressing every entry to learn
//     its compressed size before the local header is written.
package asic

import (
	"sync"
	"time"

	"github.com/ryftcore/dss-go/dss/model"
)

// ZipUtils is used for processing (reading and creation) of ZIP archives. See ZipContainerHandler.
type ZipUtils struct {
	// zipContainerHandlerBuilder provides utils for ZIP-archive content extraction.
	zipContainerHandlerBuilder ZipContainerHandlerBuilder
}

// zipUtilsSingleton and zipUtilsMutex back ZipUtilsInstance.
var (
	zipUtilsSingleton *ZipUtils
	zipUtilsMutex     sync.Mutex
)

// ZipUtilsInstance returns the instance of the ZipUtils class. Port of the static
// getInstance().
func ZipUtilsInstance() *ZipUtils {
	zipUtilsMutex.Lock()
	defer zipUtilsMutex.Unlock()
	if zipUtilsSingleton == nil {
		zipUtilsSingleton = &ZipUtils{
			zipContainerHandlerBuilder: NewSecureContainerHandlerBuilder(),
		}
	}
	return zipUtilsSingleton
}

// SetZipContainerHandlerBuilder sets a builder to create an instance of a handler to process
// ZIP-content retrieving. The handler will be created on each call of the ZipUtils class.
//
// Default : SecureContainerHandlerBuilder
//
// Panics with the Java message when zipContainerHandlerBuilder is nil (Objects.requireNonNull).
//
// Port of setZipContainerHandlerBuilder(ZipContainerHandlerBuilder).
func (z *ZipUtils) SetZipContainerHandlerBuilder(zipContainerHandlerBuilder ZipContainerHandlerBuilder) {
	if zipContainerHandlerBuilder == nil {
		panic("ZipContainerHandlerBuilder shall be defined!")
	}
	zipUtilsMutex.Lock()
	defer zipUtilsMutex.Unlock()
	z.zipContainerHandlerBuilder = zipContainerHandlerBuilder
}

// ExtractContainerContent extracts a list of DSSDocument from the given ZIP-archive. Port of
// extractContainerContent(DSSDocument).
func (z *ZipUtils) ExtractContainerContent(zipPackage model.DSSDocument) ([]model.DSSDocument, error) {
	return z.zipContainerHandler().ExtractContainerContent(zipPackage)
}

// ExtractEntryNames returns a list of ZIP archive entry names. Port of
// extractEntryNames(DSSDocument).
func (z *ZipUtils) ExtractEntryNames(zipPackage model.DSSDocument) ([]string, error) {
	return z.zipContainerHandler().ExtractEntryNames(zipPackage)
}

// CreateZipArchiveFromEntries creates a ZIP-Archive with the given containerEntries. This method
// uses the current time and no zip comment for definition of the ZIP container.
//
// Port of createZipArchive(List<DSSDocument>).
func (z *ZipUtils) CreateZipArchiveFromEntries(containerEntries []model.DSSDocument) (model.DSSDocument, error) {
	return z.CreateZipArchiveFromEntriesAt(containerEntries, time.Now(), "")
}

// CreateZipArchiveFromEntriesAt creates a ZIP-Archive with the given containerEntries.
//
// creationTime is optional: it defines the time of an archive creation and will be set for all
// embedded files; the zero Time (Java null) means the local current time will be used. zipComment
// is optional.
//
// Port of createZipArchive(List<DSSDocument>, Date, String).
func (z *ZipUtils) CreateZipArchiveFromEntriesAt(containerEntries []model.DSSDocument, creationTime time.Time, zipComment string) (model.DSSDocument, error) {
	return z.zipContainerHandler().CreateZipArchive(containerEntries, creationTime, zipComment)
}

// CreateZipArchive creates a ZIP-Archive with the given asicContent, indicating the current
// creation time. Port of createZipArchive(ASiCContent).
func (z *ZipUtils) CreateZipArchive(asicContent *Content) (model.DSSDocument, error) {
	return z.CreateZipArchiveAt(asicContent, time.Now())
}

// CreateZipArchiveAt creates a ZIP-Archive with the given asicContent.
//
// creationTime is optional: see CreateZipArchiveFromEntriesAt.
//
// Port of createZipArchive(ASiCContent, Date).
func (z *ZipUtils) CreateZipArchiveAt(asicContent *Content, creationTime time.Time) (model.DSSDocument, error) {
	return z.CreateZipArchiveFromEntriesAt(asicContent.AllDocuments(), creationTime, asicContent.ZipComment())
}

// zipContainerHandler returns a new instance of ZipContainerHandler. Port of the private
// getZipContainerHandler(). The builder field is read under the same lock
// SetZipContainerHandlerBuilder writes it under; Java leaves this read unsynchronized, which is a
// benign field race there and undefined behaviour in Go.
func (z *ZipUtils) zipContainerHandler() ZipContainerHandler {
	zipUtilsMutex.Lock()
	builder := z.zipContainerHandlerBuilder
	zipUtilsMutex.Unlock()
	return builder.Build()
}
