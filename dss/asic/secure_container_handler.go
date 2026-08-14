// Ported from
// dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/SecureContainerHandler.java
// (DSS 6.5.RC1).
//
// ZIP ENGINE (the central deviation of this file, please read before changing anything below):
//
// Upstream reads containers two different ways, and the difference is observable:
//   - java.util.zip.ZipInputStream, a strictly sequential walk over LOCAL file headers. This is
//     what extractZipEntries()/getCurrentEntryDocument() use, so it decides the entry list, the
//     entry ORDER, the per-entry `extra` bytes and the malformed-entry accounting.
//   - java.util.zip.ZipFile, a central-directory reader. This is used only by extractComments()
//     and by FileArchiveEntry's per-entry random access.
//
// Go's archive/zip offers only the second of the two: there is no ZipInputStream equivalent. Using
// the central directory for enumeration instead was measured against a Java oracle over the 188
// ASiC fixture containers of dss-asic-{common,cades,xades}/src/test/resources and is NOT
// equivalent:
//   - of the 1218 entries both Java views can see, 281 carry a different `extra` field in the
//     central directory than in their local header, and 264 of those have no local extra field at
//     all while the central directory carries an NTFS (0x000a) one - so a central-directory read
//     would resurrect metadata upstream drops, and write it back out on the next createZipArchive;
//   - validation/false-zip-comment.odp cannot be opened by archive/zip at all (nor by Java's
//     ZipFile), yet ZipInputStream reads all 23 of its entries;
//   - validation/malformed-container.{asics,asice} yield entries through archive/zip but none
//     through ZipInputStream, which is what makes upstream reject them.
//
// This file therefore carries a sequential local-file-header reader (zipLocalHeaderReader /
// zipEntryStream at the bottom) as the port of ZipInputStream, and uses archive/zip exactly where
// upstream uses ZipFile: comment extraction, isInFileProcessingSupported, FileArchiveEntry (in
// file_archive_entry.go) and archive WRITING (createZipArchive). With that split, all 188 fixtures
// produce the same entry inventory, order, metadata and per-entry content as the Java oracle.
//
// Deflate output bytes differ from java.util.zip.Deflater at the same level; that is accepted (see
// zip_utils.go's header, which documents the full write-side deviation list).
//
// MALFORMED ENTRIES: Java's readLOC consumes the fixed 30-byte header and then the entry name
// before it can fail on any later field, so getNextValidEntry's retry resumes just past the broken
// name, finds no further local header signature there and ends the walk cleanly with the entries
// collected so far - upstream returns 5 documents for validation/cp852encoded_signature.asice and
// 0 for signable/CP-852encoded.zip, with no exception in either case. nextEntry below advances the
// read position by exactly the same amount on every failure path so that the retry behaves
// identically (and, as a side effect, always terminates).
//
// slf4j LOG calls are dropped per PORTING.md; the branches they sit in are preserved.
package asic

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/flate"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"time"
	"unicode/utf8"

	"github.com/utain/esig/dss/document"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/exception"
	"github.com/utain/esig/dss/spi/signature/resources"
	"github.com/utain/esig/dss/utils"
)

// SecureContainerHandlerMimetype is the mimetype filename. Port of
// SecureContainerHandler.MIMETYPE (a second spelling of ASiCUtils.MIME_TYPE, kept because
// upstream's ensureCompressionMethod reads this one).
const SecureContainerHandlerMimetype = "mimetype"

// SecureContainerHandlerOverrides captures the protected methods of Java's
// SecureContainerHandler that a subclass may override and that the class calls on itself. Static
// Go dispatch would silently drop such an override, so every self-call is routed through this
// interface (the Overrides + Init pattern used throughout the ported tree).
type SecureContainerHandlerOverrides interface {
	// InstantiateResourcesHandler instantiates a new DSSResourcesHandler. Port of the
	// protected instantiateResourcesHandler().
	InstantiateResourcesHandler() (resources.DSSResourcesHandler, error)

	// BuildZip stores all containerEntries in a given order to a zip.Writer with the given
	// parameters. Port of the protected buildZip(List, Date, String, ZipOutputStream).
	BuildZip(containerEntries []model.DSSDocument, creationTime time.Time, zipComment string, zos *zip.Writer) error

	// ZipEntry creates a new zip.FileHeader for the given entry at creationTime. Port of the
	// protected getZipEntry(DSSDocument, Date).
	ZipEntry(entry model.DSSDocument, creationTime time.Time) (*zip.FileHeader, error)

	// SecureCopy reads and copies a stream in a secure way. Port of the protected
	// secureCopy(InputStream, OutputStream, long).
	SecureCopy(is io.Reader, os io.Writer, allowedSize int64) error

	// SecureSkip skips a stream securely without caching the content. Port of the protected
	// secureSkip(InputStream, long).
	SecureSkip(is io.Reader, allowedSize int64) error
}

// SecureContainerHandler is the default implementation of ZipContainerHandler, providing utilities
// to prevent a denial of service attacks, such as zip-bombing.
type SecureContainerHandler struct {
	// overrides points back at the concrete handler; see InitSecureContainerHandler.
	overrides SecureContainerHandlerOverrides

	// threshold is the minimum file size to be analyzed on zip bombing.
	threshold int64

	// maxCompressionRatio is the maximum compression ratio.
	maxCompressionRatio int64

	// maxAllowedFilesAmount defines the maximal amount of files that can be inside a ZIP
	// container.
	maxAllowedFilesAmount int

	// maxMalformedFiles is the max iteration over the zip entries.
	maxMalformedFiles int

	// extractComments defines whether comments of ZIP entries shall be extracted.
	// Default : false (not extracted)
	extractComments bool

	// byteCounter is an internal variable used to calculate the extracted entries size.
	// NOTE: shall be reset on every use
	byteCounter int64

	// malformedFilesCounter is an internal variable used to count a number of malformed ZIP
	// entries.
	// NOTE: shall be reset on every use
	malformedFilesCounter int

	// resourcesHandlerBuilder is the builder to be used to create a new DSSResourcesHandler for
	// each internal call, defining a way working with internal resources (e.g. in memory or by
	// using temporary files). The resources are used on a document creation.
	//
	// Default : document.InMemoryResourcesHandler, working with data in memory
	resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder
}

var _ ZipContainerHandler = (*SecureContainerHandler)(nil)
var _ SecureContainerHandlerOverrides = (*SecureContainerHandler)(nil)

// NewSecureContainerHandler instantiates the handler with default configuration. Port of the
// default constructor (field initializers included).
func NewSecureContainerHandler() *SecureContainerHandler {
	h := &SecureContainerHandler{
		threshold:               1000000, // 1 MB
		maxCompressionRatio:     100,
		maxAllowedFilesAmount:   1000,
		maxMalformedFiles:       100,
		extractComments:         false,
		resourcesHandlerBuilder: document.NewInMemoryResourcesHandlerBuilder(),
	}
	h.overrides = h
	return h
}

// InitSecureContainerHandler re-points the self-call dispatch at a subclass embedding this
// struct. It has no Java counterpart (Java gets the dispatch from the JVM).
func (h *SecureContainerHandler) InitSecureContainerHandler(overrides SecureContainerHandlerOverrides) {
	h.overrides = overrides
}

// SetThreshold sets the maximum allowed threshold after exceeding each the security checks are
// enforced.
//
// Default : 1000000 (1 MB)
//
// Port of setThreshold(long).
func (h *SecureContainerHandler) SetThreshold(threshold int64) {
	h.threshold = threshold
}

// SetMaxCompressionRatio sets the maximum allowed compression ratio. If the container compression
// ratio exceeds the value, an error is being returned.
//
// Default : 100
//
// Port of setMaxCompressionRatio(long).
func (h *SecureContainerHandler) SetMaxCompressionRatio(maxCompressionRatio int64) {
	h.maxCompressionRatio = maxCompressionRatio
}

// SetMaxAllowedFilesAmount sets the maximum allowed amount of files inside a container.
//
// Default : 1000
//
// Port of setMaxAllowedFilesAmount(int).
func (h *SecureContainerHandler) SetMaxAllowedFilesAmount(maxAllowedFilesAmount int) {
	h.maxAllowedFilesAmount = maxAllowedFilesAmount
}

// SetMaxMalformedFiles sets the maximum allowed amount of malformed files.
//
// Default : 100
//
// Port of setMaxMalformedFiles(int).
func (h *SecureContainerHandler) SetMaxMalformedFiles(maxMalformedFiles int) {
	h.maxMalformedFiles = maxMalformedFiles
}

// SetExtractComments sets whether comments of ZIP entries shall be extracted.
//
// Enabling of the feature can be useful when editing an existing archive, in order to preserve the
// existing data (i.e. comments).
//
// Reason : all ZIP entries are extracted through this port's ZipInputStream equivalent, which -
// like java.util.zip.ZipInputStream - reads local file headers and therefore cannot see entry
// comments (they live in the central directory only). To extract comments the archive has to be
// read a second time through archive/zip, which, like Java's ZipFile, needs a real file; the
// feature is consequently limited to FileDocument inputs, exactly as upstream.
//
// Default : false (not extracted)
//
// Port of setExtractComments(boolean).
func (h *SecureContainerHandler) SetExtractComments(extractComments bool) {
	h.extractComments = extractComments
}

// SetResourcesHandlerBuilder sets the DSSResourcesHandlerBuilder to be used for a
// DSSResourcesHandler creation in internal methods. DSSResourcesHandler defines a way to operate
// with OutputStreams and create DSSDocuments.
//
// Default : document.InMemoryResourcesHandler. Works with data in memory.
//
// Panics with the Java message when resourcesHandlerBuilder is nil (Objects.requireNonNull).
//
// Port of setResourcesHandlerBuilder(DSSResourcesHandlerBuilder).
func (h *SecureContainerHandler) SetResourcesHandlerBuilder(resourcesHandlerBuilder resources.DSSResourcesHandlerBuilder) {
	if resourcesHandlerBuilder == nil {
		panic("DSSResourcesHandlerBuilder cannot be null!")
	}
	h.resourcesHandlerBuilder = resourcesHandlerBuilder
}

// ExtractContainerContent extracts a list of DSSDocument from the given ZIP-archive. Port of
// extractContainerContent(DSSDocument).
func (h *SecureContainerHandler) ExtractContainerContent(zipArchive model.DSSDocument) ([]model.DSSDocument, error) {
	h.resetCounters()

	result := make([]model.DSSDocument, 0)
	if h.isInFileProcessingSupported(zipArchive) {
		zipFileDocument, _ := zipArchive.(*model.FileDocument)
		zipEntries, err := h.extractZipEntries(zipFileDocument)
		if err != nil {
			return nil, err
		}
		if !h.malformedEntriesDetected() {
			for _, zipEntry := range zipEntries {
				result = append(result, newFileArchiveEntry(zipFileDocument, zipEntry))
			}
			return result, nil
		}
		// The archive contains malformed entries and cannot be parsed through the central
		// directory. Continue with the sequential reader.
	}

	containerSize, err := spi.DSSUtilsFileByteSize(zipArchive)
	if err != nil {
		return nil, exception.NewIllegalInputExceptionWithCause("Unable to extract content from zip archive", err)
	}
	zis, closer, err := h.openZipStream(zipArchive)
	if err != nil {
		return nil, exception.NewIllegalInputExceptionWithCause("Unable to extract content from zip archive", err)
	}
	defer closer()
	for {
		currentDocument, err := h.getNextDocument(zis, containerSize)
		if err != nil {
			return nil, err
		}
		if currentDocument == nil {
			break
		}
		result = append(result, currentDocument)
		if err := h.assertCollectionSizeValid(len(result)); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// isInFileProcessingSupported verifies whether the provided archive container is supported by the
// archive/zip central-directory reader (upstream: java.util.zip.ZipFile). Port of
// isInFileProcessingSupported(DSSDocument).
func (h *SecureContainerHandler) isInFileProcessingSupported(zipArchive model.DSSDocument) bool {
	if fileDocument, ok := zipArchive.(*model.FileDocument); ok {
		readCloser, err := zip.OpenReader(fileDocument.Path())
		if err != nil {
			return false
		}
		_ = readCloser.Close()
		return true
	}
	return false
}

// malformedEntriesDetected reports whether the ZIP archive contains malformed ZIP entries.
//
// The central-directory reader is not able to work with malformed archives. Therefore, we need to
// continue with the sequential reader when encountering a malformed ZIP archive.
//
// Port of malformedEntriesDetected().
func (h *SecureContainerHandler) malformedEntriesDetected() bool {
	return h.malformedFilesCounter > 0
}

// getNextDocument returns the next document from zis, or nil when the end of the archive is
// reached. Port of getNextDocument(ZipInputStream, long).
func (h *SecureContainerHandler) getNextDocument(zis *zipLocalHeaderReader, containerSize int64) (model.DSSDocument, error) {
	entry, err := h.getNextValidEntry(zis)
	if err != nil {
		return nil, err
	}
	if entry != nil {
		return h.getCurrentEntryDocument(zis, entry, containerSize)
	}
	return nil, nil
}

// ExtractEntryNames returns a list of ZIP archive entry names. Port of
// extractEntryNames(DSSDocument).
func (h *SecureContainerHandler) ExtractEntryNames(zipArchive model.DSSDocument) ([]string, error) {
	zipEntries, err := h.extractZipEntries(zipArchive)
	if err != nil {
		return nil, err
	}
	if utils.IsCollectionNotEmpty(zipEntries) {
		names := make([]string, 0, len(zipEntries))
		for _, zipEntry := range zipEntries {
			names = append(names, zipEntry.Name)
		}
		return names, nil
	}
	return []string{}, nil
}

// extractZipEntries reads the archive sequentially, collecting every entry header. Port of
// extractZipEntries(DSSDocument).
func (h *SecureContainerHandler) extractZipEntries(zipArchive model.DSSDocument) ([]*zip.FileHeader, error) {
	h.resetCounters()

	containerSize, err := spi.DSSUtilsFileByteSize(zipArchive)
	if err != nil {
		return nil, errors.New("Unable to extract entries from zip archive")
	}
	allowedSize := containerSize * h.maxCompressionRatio

	/*
	 * Read with the sequential local-header reader in order to extract ZipEntry dates
	 */
	result := make([]*zip.FileHeader, 0)
	zis, closer, err := h.openZipStream(zipArchive)
	if err != nil {
		return nil, errors.New("Unable to extract entries from zip archive")
	}
	defer closer()
	for {
		entry, err := h.getNextValidEntry(zis)
		if err != nil {
			return nil, err
		}
		if entry == nil {
			break
		}
		result = append(result, entry)
		if err := h.assertCollectionSizeValid(len(result)); err != nil {
			return nil, err
		}
		// read securely before accessing the next entry
		if err := h.overrides.SecureSkip(zis, allowedSize); err != nil {
			if _, ok := err.(*exception.IllegalInputException); ok {
				return nil, err
			}
			return nil, errors.New("Unable to extract entries from zip archive")
		}
	}
	h.extractCommentsInto(zipArchive, result)
	return result, nil
}

// extractCommentsInto fills in the entry comments, when enabled.
//
// When reading a zip file sequentially, the comment is not available (it only exists in the
// central directory - see https://bugs.openjdk.java.net/browse/JDK-4201267). Therefore, we need to
// read comments through archive/zip, when possible.
//
// Port of extractComments(DSSDocument, List<ZipEntry>).
func (h *SecureContainerHandler) extractCommentsInto(zipArchive model.DSSDocument, zipEntries []*zip.FileHeader) {
	if h.extractComments {
		fileDocument, ok := zipArchive.(*model.FileDocument)
		if !ok {
			return
		}
		readCloser, err := zip.OpenReader(fileDocument.Path())
		if err != nil {
			return
		}
		defer readCloser.Close()
		for _, zipEntry := range zipEntries {
			for _, zipFileEntry := range readCloser.File {
				if zipFileEntry.Name == zipEntry.Name {
					zipEntry.Comment = zipFileEntry.Comment
					break
				}
			}
		}
	}
}

// CreateZipArchive creates a ZIP-Archive with the given containerEntries. Port of
// createZipArchive(List<DSSDocument>, Date, String).
func (h *SecureContainerHandler) CreateZipArchive(containerEntries []model.DSSDocument, creationTime time.Time, zipComment string) (model.DSSDocument, error) {
	dssResourcesHandler, err := h.overrides.InstantiateResourcesHandler()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to create an ASiC container. Reason : %s", err.Error()), err)
	}
	defer dssResourcesHandler.Close()
	outputStream, err := dssResourcesHandler.CreateOutputStream()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to create an ASiC container. Reason : %s", err.Error()), err)
	}
	zos := zip.NewWriter(outputStream)
	if err := h.overrides.BuildZip(containerEntries, creationTime, zipComment, zos); err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to create an ASiC container. Reason : %s", err.Error()), err)
	}
	dssDocument, err := dssResourcesHandler.WriteToDSSDocument()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to create an ASiC container. Reason : %s", err.Error()), err)
	}
	return dssDocument, nil
}

// InstantiateResourcesHandler instantiates a new DSSResourcesHandler. Port of the protected
// instantiateResourcesHandler().
func (h *SecureContainerHandler) InstantiateResourcesHandler() (resources.DSSResourcesHandler, error) {
	return h.resourcesHandlerBuilder.CreateResourcesHandler(), nil
}

// BuildZip stores all containerEntries in a given order to zos with the given parameters. Port of
// the protected buildZip(List<DSSDocument>, Date, String, ZipOutputStream).
//
// The STORED/DEFLATED split below is what keeps the "mimetype" entry spec-compliant: archive/zip's
// CreateHeader unconditionally sets the data-descriptor flag and moves the sizes out of the local
// header, which EN 319 162-1 A.1 readers that expect the mimetype payload at a fixed offset cannot
// follow, and which java.util.zip.ZipOutputStream does not do for STORED entries either.
// CreateRaw writes the sizes in the local header instead, matching Java byte layout.
func (h *SecureContainerHandler) BuildZip(containerEntries []model.DSSDocument, creationTime time.Time, zipComment string, zos *zip.Writer) error {
	for _, entry := range containerEntries {
		zipEntry, err := h.overrides.ZipEntry(entry, creationTime)
		if err != nil {
			return err
		}
		entryWriter, err := secureContainerHandlerCreateEntry(zos, zipEntry)
		if err != nil {
			return err
		}
		entryIS, err := entry.OpenStream()
		if err != nil {
			return err
		}
		err = h.overrides.SecureCopy(entryIS, entryWriter, -1)
		_ = entryIS.Close()
		if err != nil {
			return err
		}
	}
	if utils.IsStringNotEmpty(zipComment) {
		if err := zos.SetComment(zipComment); err != nil {
			return err
		}
	}
	return zos.Close() // ZipOutputStream#finish
}

// secureContainerHandlerCreateEntry starts a new ZIP entry, choosing the archive/zip API that
// reproduces java.util.zip.ZipOutputStream's header layout for the entry's compression method.
// It has no Java counterpart - Java has a single putNextEntry that branches internally.
func secureContainerHandlerCreateEntry(zos *zip.Writer, zipEntry *zip.FileHeader) (io.Writer, error) {
	if zipEntry.Method != zip.Store {
		// DEFLATED: java.util.zip.ZipOutputStream sets the data-descriptor flag and writes
		// crc/sizes after the payload, which is exactly what CreateHeader does.
		return zos.CreateHeader(zipEntry)
	}
	// STORED: crc and sizes are already known (addStoredContent computed them) and must
	// precede the payload, with no data descriptor.
	zipEntry.ReaderVersion = 20
	zipEntry.CreatorVersion = 20
	if !utf8ASCIIOnly(zipEntry.Name) || !utf8ASCIIOnly(zipEntry.Comment) {
		if utf8.ValidString(zipEntry.Name) && utf8.ValidString(zipEntry.Comment) {
			zipEntry.Flags |= 0x800
		}
	}
	return zos.CreateRaw(zipEntry)
}

// utf8ASCIIOnly reports whether s is pure ASCII, i.e. whether the UTF-8 general-purpose flag bit
// may be left clear. It mirrors the CP-437-compatibility test archive/zip's CreateHeader performs
// for the DEFLATED path so that both paths in this file agree.
func utf8ASCIIOnly(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

// ZipEntry creates a new zip.FileHeader for the given entry at creationTime. Port of the protected
// getZipEntry(DSSDocument, Date).
func (h *SecureContainerHandler) ZipEntry(entry model.DSSDocument, creationTime time.Time) (*zip.FileHeader, error) {
	var zipEntryWrapper *DSSZipEntry
	if dssZipEntry, ok := entry.(DSSZipEntryDocument); ok {
		zipEntryWrapper = dssZipEntry.ZipEntry()
	} else {
		zipEntryWrapper = NewDSSZipEntry(entry.Name())
	}
	zipEntry := zipEntryWrapper.CreateZipEntry()
	if err := h.ensureCompressionMethod(zipEntry, entry); err != nil {
		return nil, err
	}
	h.ensureTime(zipEntry, creationTime)
	return zipEntry, nil
}

// ensureCompressionMethod is the port of ensureCompressionMethod(ZipEntry, DSSDocument).
func (h *SecureContainerHandler) ensureCompressionMethod(zipEntry *zip.FileHeader, content model.DSSDocument) error {
	if SecureContainerHandlerMimetype == zipEntry.Name {
		/*
		 * EN 319 162-1 "A.1 The mimetype file":
		 * "mimetype" shall not be compressed (i.e. compression method in its ZIP header at
		 * offset 8 shall be set to zero);
		 */
		// (upstream logs a warning here when the method was set to something other than
		// STORED; slf4j is dropped)
		zipEntry.Method = zip.Store
	}
	/*
	 * If you switch to STORED note that you'll have to set the size (or compressed size; they
	 * must be the same, but it's okay to only set one) and CRC yourself because they must
	 * appear before the user data in the resulting zip file.
	 */
	if zipEntry.Method == zip.Store {
		return h.addStoredContent(zipEntry, content)
	}
	/*
	 * The default is DEFLATED, which will cause the size, compressed size, and CRC to be set
	 * automatically, and the entry's data to be compressed.
	 */
	return nil
}

// addStoredContent is the port of addStoredContent(ZipEntry, DSSDocument).
func (h *SecureContainerHandler) addStoredContent(zipEntry *zip.FileHeader, content model.DSSDocument) error {
	var size int64
	crc := crc32.NewIEEE()
	is, err := content.OpenStream()
	if err != nil {
		// upstream's try-with-resources lets a DSSException from openStream() propagate;
		// only IOException from the read loop is swallowed with a warning.
		return err
	}
	buffer := make([]byte, 8192)
	for {
		nbRead, readErr := is.Read(buffer)
		if nbRead > 0 {
			crc.Write(buffer[:nbRead])
			size += int64(nbRead)
		}
		if readErr != nil {
			// upstream: catch (IOException e) { LOG.warn("Unable to process CRC32
			// computation", e); } - the partially accumulated size/crc are kept.
			break
		}
	}
	_ = is.Close()
	zipEntry.UncompressedSize64 = uint64(size)
	zipEntry.CompressedSize64 = uint64(size)
	zipEntry.CRC32 = crc.Sum32()
	return nil
}

// ensureTime is the port of ensureTime(ZipEntry, Date).
//
// The MS-DOS date/time fields are set directly rather than through archive/zip's FileHeader.
// Modified, because CreateHeader appends a synthesized Info-ZIP extended-timestamp (0x5455) extra
// field whenever Modified is set - metadata java.util.zip.ZipOutputStream does not add for an
// entry whose FileTime-valued mtime is null, which is exactly the state ZipEntry#setTime leaves
// the entry in.
func (h *SecureContainerHandler) ensureTime(zipEntry *zip.FileHeader, creationTime time.Time) {
	// if not set, the local current time will be used
	if !creationTime.IsZero() {
		zipEntry.ModifiedDate, zipEntry.ModifiedTime = secureContainerHandlerToMsDosTime(creationTime)
		zipEntry.Modified = time.Time{}
	} else {
		zipEntry.ModifiedDate, zipEntry.ModifiedTime = secureContainerHandlerToMsDosTime(time.Now())
		zipEntry.Modified = time.Time{}
	}
}

// secureContainerHandlerToMsDosTime converts t to the MS-DOS date/time pair, in t's own location.
// Port of java.util.zip.ZipUtils#javaToDosTime, including its pre-1980 clamp to 1980-01-01
// 00:00:00 - which archive/zip's own timeToMsDosTime lacks.
func secureContainerHandlerToMsDosTime(t time.Time) (fDate, fTime uint16) {
	if t.Year() < 1980 {
		// Java's (1 << 21) | (1 << 16) combined dostime, i.e. date = month 1, day 1.
		return (1 << 5) | 1, 0
	}
	fDate = uint16(t.Day() + int(t.Month())<<5 + (t.Year()-1980)<<9)
	fTime = uint16(t.Second()/2 + t.Minute()<<5 + t.Hour()<<11)
	return fDate, fTime
}

// resetCounters is the port of resetCounters().
func (h *SecureContainerHandler) resetCounters() {
	h.byteCounter = 0
	h.malformedFilesCounter = 0
}

// getNextValidEntry returns the next entry from the given reader by skipping corrupted or not
// accessible files.
//
// NOTE: returns nil only when the end of the archive is reached.
//
// Port of getNextValidEntry(ZipInputStream); returns a DSSError if too many tries failed.
func (h *SecureContainerHandler) getNextValidEntry(zis *zipLocalHeaderReader) (*zip.FileHeader, error) {
	for h.malformedFilesCounter < h.maxMalformedFiles {
		entry, err := zis.nextEntry()
		if err == nil {
			return entry, nil
		}
		// ZIP container contains a malformed, corrupted or not accessible entry! The entry
		// is skipped. (upstream logs the reason; slf4j is dropped)
		h.malformedFilesCounter++
		if closeErr := h.closeEntry(zis); closeErr != nil {
			return nil, closeErr
		}
	}
	return nil, model.NewDSSError(fmt.Sprintf("Unable to retrieve a valid ZipEntry (%d tries)", h.maxMalformedFiles))
}

// closeEntry closes the current ZIP entry. Port of closeEntry(ZipInputStream).
func (h *SecureContainerHandler) closeEntry(zis *zipLocalHeaderReader) error {
	if err := zis.closeEntry(); err != nil {
		return model.NewDSSErrorMessageCause("Unable to close entry", err)
	}
	return nil
}

// getCurrentEntryDocument returns the current file from the given reader. Port of
// getCurrentEntryDocument(ZipInputStream, ZipEntry, long).
func (h *SecureContainerHandler) getCurrentEntryDocument(zis *zipLocalHeaderReader, entry *zip.FileHeader, containerSize int64) (model.DSSDocument, error) {
	allowedSize := containerSize * h.maxCompressionRatio
	dssResourcesHandler, err := h.overrides.InstantiateResourcesHandler()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to read an entry binaries. Reason : %s", err.Error()), err)
	}
	defer dssResourcesHandler.Close()
	outputStream, err := dssResourcesHandler.CreateOutputStream()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to read an entry binaries. Reason : %s", err.Error()), err)
	}
	if err := h.overrides.SecureCopy(zis, outputStream, allowedSize); err != nil {
		if _, ok := err.(*exception.IllegalInputException); ok {
			return nil, err
		}
		_ = h.closeEntry(zis)
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to read an entry binaries. Reason : %s", err.Error()), err)
	}

	currentDocument, err := dssResourcesHandler.WriteToDSSDocument()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(fmt.Sprintf("Unable to read an entry binaries. Reason : %s", err.Error()), err)
	}
	fileName := entry.Name
	currentDocument.SetName(entry.Name)
	currentDocument.SetMimeType(enumerations.MimeTypeFromFileName(fileName))

	containerEntryDocument, err := NewContainerEntryDocumentWithZipEntry(currentDocument, NewDSSZipEntryFromFileHeader(entry))
	if err != nil {
		return nil, err
	}
	return containerEntryDocument, nil
}

// SecureCopy reads and copies a stream in a secure way to a writer. Detects "ZipBombing" (large
// files inside a zip container) depending on the provided container size.
//
// allowedSize defines an allowed size of the ZIP container entries; -1 skips the validation.
//
// Port of the protected secureCopy(InputStream, OutputStream, long).
func (h *SecureContainerHandler) SecureCopy(is io.Reader, os io.Writer, allowedSize int64) error {
	data := make([]byte, 8192)
	for {
		nRead, err := is.Read(data)
		if nRead > 0 {
			h.byteCounter += int64(nRead)
			if assertErr := h.assertExtractEntryLengthValid(allowedSize); assertErr != nil {
				return assertErr
			}
			if _, writeErr := os.Write(data[:nRead]); writeErr != nil {
				return writeErr
			}
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// SecureSkip skips a stream securely without caching the content.
//
// Port of the protected secureSkip(InputStream, long).
func (h *SecureContainerHandler) SecureSkip(is io.Reader, allowedSize int64) error {
	for {
		nRead, err := io.CopyN(io.Discard, is, 8192)
		if nRead > 0 {
			h.byteCounter += nRead
			if assertErr := h.assertExtractEntryLengthValid(allowedSize); assertErr != nil {
				return assertErr
			}
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if nRead == 0 {
			return nil
		}
	}
}

// assertExtractEntryLengthValid is the port of assertExtractEntryLengthValid(long).
func (h *SecureContainerHandler) assertExtractEntryLengthValid(allowedSize int64) error {
	if allowedSize != -1 && h.byteCounter > h.threshold && h.byteCounter > allowedSize {
		return exception.NewIllegalInputException("Zip Bomb detected in the ZIP container. Validation is interrupted.")
	}
	return nil
}

// assertCollectionSizeValid is the port of assertCollectionSizeValid(Collection).
func (h *SecureContainerHandler) assertCollectionSizeValid(size int) error {
	if size > h.maxAllowedFilesAmount {
		return exception.NewIllegalInputException("Too many files detected. Cannot extract ASiC content from the file.")
	}
	return nil
}

// openZipStream opens the sequential reader over zipArchive. It has no dedicated Java counterpart:
// upstream simply wraps document.openStream() in a ZipInputStream. Random access is needed here to
// locate a data descriptor precisely (see zipEntryStream.finish), so a FileDocument is read
// through its file handle and any other document is materialized in memory - which, for the
// in-memory documents this is reached with in practice, costs nothing extra.
func (h *SecureContainerHandler) openZipStream(zipArchive model.DSSDocument) (*zipLocalHeaderReader, func(), error) {
	if fileDocument, ok := zipArchive.(*model.FileDocument); ok {
		file, err := os.Open(fileDocument.Path())
		if err != nil {
			return nil, func() {}, err
		}
		info, err := file.Stat()
		if err != nil {
			_ = file.Close()
			return nil, func() {}, err
		}
		return &zipLocalHeaderReader{ra: file, size: info.Size()}, func() { _ = file.Close() }, nil
	}
	stream, err := zipArchive.OpenStream()
	if err != nil {
		return nil, func() {}, err
	}
	defer stream.Close()
	data, err := io.ReadAll(stream)
	if err != nil {
		return nil, func() {}, err
	}
	return &zipLocalHeaderReader{ra: bytes.NewReader(data), size: int64(len(data))}, func() {}, nil
}

// ---------------------------------------------------------------------------------------------
// java.util.zip.ZipInputStream equivalent
// ---------------------------------------------------------------------------------------------

// ZIP structural constants, named after java.util.zip.ZipConstants.
const (
	zipLocalHeaderSignature      uint32 = 0x04034b50 // LOCSIG
	zipDataDescriptorSignature   uint32 = 0x08074b50 // EXTSIG
	zipLocalHeaderLen                   = 30         // LOCHDR
	zipDataDescriptorLen                = 16         // EXTHDR
	zipZip64DataDescriptorLen           = 24         // ZIP64_EXTHDR
	zipZip64MagicValue           uint64 = 0xFFFFFFFF // ZIP64_MAGICVAL
	zipExtID64                          = 0x0001     // EXTID_ZIP64
	zipFlagEncrypted             uint16 = 0x1
	zipFlagHasDataDescriptor     uint16 = 0x8
	zipCompressionMethodStored   uint16 = 0
	zipCompressionMethodDeflated uint16 = 8
)

// zipLocalHeaderReader walks a ZIP archive's local file headers in order, exactly as
// java.util.zip.ZipInputStream does, and reads the current entry's uncompressed content through
// its own io.Reader implementation (which is how upstream passes `zis` to secureCopy/secureSkip).
type zipLocalHeaderReader struct {
	// ra is the archive's random-access view; see openZipStream for why random access is
	// required.
	ra io.ReaderAt

	// size is the archive's total byte size.
	size int64

	// pos is the offset of the next local file header.
	pos int64

	// cur is the stream over the entry returned by the last nextEntry call.
	cur *zipEntryStream
}

// nextEntry positions the reader on the next entry and returns its header, or (nil, nil) at the
// end of the archive. Port of ZipInputStream#getNextEntry / #readLOC.
func (r *zipLocalHeaderReader) nextEntry() (*zip.FileHeader, error) {
	if err := r.closeEntry(); err != nil {
		return nil, err
	}

	header := make([]byte, zipLocalHeaderLen)
	if _, err := r.readAt(header, r.pos); err != nil {
		return nil, nil // ZipInputStream#readLOC: catch (EOFException e) { return null; }
	}
	if binary.LittleEndian.Uint32(header) != zipLocalHeaderSignature {
		return nil, nil
	}
	flag := binary.LittleEndian.Uint16(header[6:])
	method := binary.LittleEndian.Uint16(header[8:])
	modTime := binary.LittleEndian.Uint16(header[10:])
	modDate := binary.LittleEndian.Uint16(header[12:])
	crc := binary.LittleEndian.Uint32(header[14:])
	compressedSize := uint64(binary.LittleEndian.Uint32(header[18:]))
	uncompressedSize := uint64(binary.LittleEndian.Uint32(header[22:]))
	nameLen := int(binary.LittleEndian.Uint16(header[26:]))
	extraLen := int(binary.LittleEndian.Uint16(header[28:]))

	// From here on a failure leaves the read position where java.util.zip.ZipInputStream's own
	// readLOC would have left it - it consumes the fixed header, then the name, before it can
	// fail on any later field. That is what makes getNextValidEntry's retry land past the
	// broken entry and find no further local header signature, ending the walk cleanly with the
	// entries collected so far, exactly as upstream does for
	// validation/cp852encoded_signature.asice (5 entries, no exception).
	nameBuffer := make([]byte, nameLen)
	if _, err := r.readAt(nameBuffer, r.pos+zipLocalHeaderLen); err != nil {
		r.pos += zipLocalHeaderLen
		return nil, err
	}
	name := string(nameBuffer)
	r.pos += zipLocalHeaderLen + int64(nameLen)
	// DEVIATION: java.util.zip decodes entry names with a reporting UTF-8 decoder and throws
	// IllegalArgumentException("malformed input ...") on non-UTF-8 bytes (which is how upstream
	// stops early on signable/CP-852encoded.zip and validation/cp852encoded_signature.asice). Go
	// strings hold arbitrary bytes, so the check is made explicitly to keep that behaviour.
	if !utf8.ValidString(name) {
		return nil, errors.New("malformed input in entry name")
	}
	if flag&zipFlagEncrypted != 0 {
		return nil, errors.New("encrypted ZIP entry not supported")
	}
	var extra []byte
	if extraLen > 0 {
		extra = make([]byte, extraLen)
		if _, err := r.readAt(extra, r.pos); err != nil {
			return nil, err
		}
	}
	if flag&zipFlagHasDataDescriptor != 0 {
		if method != zipCompressionMethodDeflated {
			return nil, errors.New("only DEFLATED entries can have EXT descriptor")
		}
		crc, compressedSize, uncompressedSize = 0, 0, 0
	} else if compressedSize == zipZip64MagicValue || uncompressedSize == zipZip64MagicValue {
		uncompressedSize, compressedSize = zipZip64SizesFromExtra(extra, uncompressedSize, compressedSize)
	}

	fileHeader := &zip.FileHeader{
		Name:               name,
		Method:             method,
		Flags:              flag,
		ModifiedTime:       modTime,
		ModifiedDate:       modDate,
		Modified:           zipDosToTime(modDate, modTime),
		Extra:              extra,
		CRC32:              crc,
		CompressedSize64:   compressedSize,
		UncompressedSize64: uncompressedSize,
	}
	dataStart := r.pos + int64(extraLen)
	r.cur = newZipEntryStream(r, fileHeader, dataStart)
	return fileHeader, nil
}

// closeEntry drains the current entry so that the reader lands on the next local file header. Port
// of ZipInputStream#closeEntry; like Java's, it is a no-op when no entry is open (which is the
// state a failed header parse leaves the reader in).
func (r *zipLocalHeaderReader) closeEntry() error {
	if r.cur == nil {
		return nil
	}
	_, err := io.Copy(io.Discard, r.cur)
	r.cur = nil
	return err
}

// Read reads the current entry's uncompressed content. Port of ZipInputStream#read.
func (r *zipLocalHeaderReader) Read(p []byte) (int, error) {
	if r.cur == nil {
		return 0, io.EOF
	}
	return r.cur.Read(p)
}

// readAt fills buf from off, reporting io.EOF-style truncation as an error.
func (r *zipLocalHeaderReader) readAt(buf []byte, off int64) (int, error) {
	if off < 0 || off+int64(len(buf)) > r.size {
		return 0, io.ErrUnexpectedEOF
	}
	return io.ReadFull(io.NewSectionReader(r.ra, off, int64(len(buf))), buf)
}

// zipEntryStream streams one entry's uncompressed content, verifying its CRC and sizes at the end
// and, for entries carrying one, consuming the trailing data descriptor. It is the port of the
// per-entry half of ZipInputStream (#read's DEFLATED/STORED branches plus #readEnd).
type zipEntryStream struct {
	reader     *zipLocalHeaderReader
	fileHeader *zip.FileHeader
	dataStart  int64

	// stored is set for uncompressed entries, which are read straight out of the archive.
	stored bool

	// remaining counts the bytes left of a STORED entry.
	remaining int64

	// src is the decompressed byte source (the raw section for STORED entries, a flate reader
	// for DEFLATED ones).
	src io.Reader

	// counter and buffered together locate the end of a DEFLATED stream exactly: the flate
	// reader consumes whole bytes only, so the compressed length is the number of bytes handed
	// to the bufio.Reader minus what is still buffered in it.
	counter  *zipCountingReader
	buffered *bufio.Reader

	crc     hash.Hash32
	written int64
	eof     bool
	err     error
}

// newZipEntryStream prepares the content stream for the entry starting at dataStart.
func newZipEntryStream(reader *zipLocalHeaderReader, fileHeader *zip.FileHeader, dataStart int64) *zipEntryStream {
	s := &zipEntryStream{
		reader:     reader,
		fileHeader: fileHeader,
		dataStart:  dataStart,
		crc:        crc32.NewIEEE(),
	}
	// Upstream's readLOC accepts any compression method and any truncation; the failure only
	// surfaces once the entry is read (ZipInputStream#read throws "invalid compression method"
	// / "unexpected EOF"), which is why these are recorded as a pending stream error rather
	// than a header error - a header error would instead be counted as a malformed entry.
	if fileHeader.Method != zipCompressionMethodStored && fileHeader.Method != zipCompressionMethodDeflated {
		s.err = fmt.Errorf("invalid compression method %d", fileHeader.Method)
		return s
	}
	if dataStart > reader.size {
		s.err = io.ErrUnexpectedEOF
		return s
	}
	if fileHeader.Method == zipCompressionMethodStored {
		s.stored = true
		s.remaining = int64(fileHeader.UncompressedSize64)
		if dataStart+s.remaining > reader.size {
			s.err = io.ErrUnexpectedEOF
			return s
		}
		s.src = io.NewSectionReader(reader.ra, dataStart, s.remaining)
		if s.remaining == 0 {
			s.reader.pos = dataStart
		}
		return s
	}
	s.counter = &zipCountingReader{r: io.NewSectionReader(reader.ra, dataStart, reader.size-dataStart)}
	s.buffered = bufio.NewReader(s.counter)
	s.src = flate.NewReader(s.buffered)
	return s
}

// Read ports the body of ZipInputStream#read for the current entry.
func (s *zipEntryStream) Read(p []byte) (int, error) {
	if s.err != nil {
		return 0, s.err
	}
	if s.eof {
		return 0, io.EOF
	}
	if s.stored {
		if s.remaining <= 0 {
			s.eof = true
			return 0, io.EOF
		}
		if int64(len(p)) > s.remaining {
			p = p[:s.remaining]
		}
		n, err := s.src.Read(p)
		if n == 0 && err != nil {
			s.err = errors.New("unexpected EOF")
			return 0, s.err
		}
		s.crc.Write(p[:n])
		s.remaining -= int64(n)
		s.written += int64(n)
		if s.remaining == 0 {
			if s.crc.Sum32() != s.fileHeader.CRC32 {
				s.err = fmt.Errorf("invalid entry CRC (expected 0x%x but got 0x%x)", s.fileHeader.CRC32, s.crc.Sum32())
				return n, s.err
			}
			s.reader.pos = s.dataStart + s.written
		}
		return n, nil
	}
	n, err := s.src.Read(p)
	if n > 0 {
		s.crc.Write(p[:n])
		s.written += int64(n)
	}
	if err != nil {
		if err == io.EOF {
			s.eof = true
			if finishErr := s.finish(); finishErr != nil {
				s.err = finishErr
				return n, finishErr
			}
			if n > 0 {
				return n, nil
			}
			return 0, io.EOF
		}
		s.err = err
		return n, err
	}
	return n, nil
}

// finish consumes the data descriptor, when present, and validates the entry. Port of
// ZipInputStream#readEnd.
func (s *zipEntryStream) finish() error {
	consumed := s.counter.n - int64(s.buffered.Buffered())
	position := s.dataStart + consumed

	if s.fileHeader.Flags&zipFlagHasDataDescriptor != 0 {
		if uint64(s.written) > zipZip64MagicValue || uint64(consumed) > zipZip64MagicValue {
			descriptor := make([]byte, zipZip64DataDescriptorLen)
			if _, err := s.reader.readAt(descriptor, position); err != nil {
				return err
			}
			if signature := binary.LittleEndian.Uint32(descriptor); signature != zipDataDescriptorSignature {
				s.fileHeader.CRC32 = signature
				s.fileHeader.CompressedSize64 = binary.LittleEndian.Uint64(descriptor[4:])
				s.fileHeader.UncompressedSize64 = binary.LittleEndian.Uint64(descriptor[12:])
				position += zipZip64DataDescriptorLen - 4
			} else {
				s.fileHeader.CRC32 = binary.LittleEndian.Uint32(descriptor[4:])
				s.fileHeader.CompressedSize64 = binary.LittleEndian.Uint64(descriptor[8:])
				s.fileHeader.UncompressedSize64 = binary.LittleEndian.Uint64(descriptor[16:])
				position += zipZip64DataDescriptorLen
			}
		} else {
			descriptor := make([]byte, zipDataDescriptorLen)
			if _, err := s.reader.readAt(descriptor, position); err != nil {
				return err
			}
			if signature := binary.LittleEndian.Uint32(descriptor); signature != zipDataDescriptorSignature {
				s.fileHeader.CRC32 = signature
				s.fileHeader.CompressedSize64 = uint64(binary.LittleEndian.Uint32(descriptor[4:]))
				s.fileHeader.UncompressedSize64 = uint64(binary.LittleEndian.Uint32(descriptor[8:]))
				position += zipDataDescriptorLen - 4
			} else {
				s.fileHeader.CRC32 = binary.LittleEndian.Uint32(descriptor[4:])
				s.fileHeader.CompressedSize64 = uint64(binary.LittleEndian.Uint32(descriptor[8:]))
				s.fileHeader.UncompressedSize64 = uint64(binary.LittleEndian.Uint32(descriptor[12:]))
				position += zipDataDescriptorLen
			}
		}
	}
	if int64(s.fileHeader.UncompressedSize64) != s.written {
		return fmt.Errorf("invalid entry size (expected %d but got %d bytes)", s.fileHeader.UncompressedSize64, s.written)
	}
	if int64(s.fileHeader.CompressedSize64) != consumed {
		return fmt.Errorf("invalid entry compressed size (expected %d but got %d bytes)", s.fileHeader.CompressedSize64, consumed)
	}
	if s.fileHeader.CRC32 != s.crc.Sum32() {
		return fmt.Errorf("invalid entry CRC (expected 0x%x but got 0x%x)", s.fileHeader.CRC32, s.crc.Sum32())
	}
	s.reader.pos = position
	return nil
}

// zipCountingReader counts the bytes delivered by the wrapped reader.
type zipCountingReader struct {
	r io.Reader
	n int64
}

// Read implements io.Reader.
func (c *zipCountingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += int64(n)
	return n, err
}

// zipZip64SizesFromExtra resolves the ZIP64 sizes a local header defers to its extra field. Port
// of the EXTID_ZIP64 branch of java.util.zip.ZipEntry#setExtra0 for a local header.
func zipZip64SizesFromExtra(extra []byte, uncompressedSize, compressedSize uint64) (uint64, uint64) {
	off := 0
	length := len(extra)
	for off+4 <= length {
		tag := int(binary.LittleEndian.Uint16(extra[off:]))
		sz := int(binary.LittleEndian.Uint16(extra[off+2:]))
		off += 4
		if off+sz > length {
			break
		}
		if tag == zipExtID64 {
			// LOC extra ZIP64 entry MUST include BOTH original and compressed file size
			// fields.
			if sz >= 16 {
				uncompressedSize = binary.LittleEndian.Uint64(extra[off:])
				compressedSize = binary.LittleEndian.Uint64(extra[off+8:])
			}
			break
		}
		off += sz
	}
	return uncompressedSize, compressedSize
}

// zipDosToTime converts an MS-DOS date/time pair to a time.Time in the local zone. Port of
// java.util.zip.ZipUtils#dosToJavaTime, which - unlike archive/zip's msDosTimeToTime - resolves
// the wall-clock fields against ZoneId.systemDefault(); the resulting value is therefore
// machine-dependent in exactly the way upstream's is.
func zipDosToTime(fDate, fTime uint16) time.Time {
	return time.Date(
		int(fDate>>9)+1980,
		time.Month(fDate>>5&0xf),
		int(fDate&0x1f),
		int(fTime>>11),
		int(fTime>>5&0x3f),
		int(fTime&0x1f)*2,
		0,
		time.Local,
	)
}
