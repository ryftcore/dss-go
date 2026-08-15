// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/DSSZipEntry.java
// (DSS 6.5.RC1).
//
// TYPE MAPPING: java.util.zip.ZipEntry has no direct Go counterpart. archive/zip.FileHeader is the
// closest equivalent and is what this port uses for both the DSSZipEntry(ZipEntry) constructor
// (NewDSSZipEntryFromFileHeader) and createZipEntry(). java.nio.file.attribute.FileTime and
// java.util.Date both become time.Time, with Java's null spelled as the zero Time (epoch 0 is a
// legitimate value and is NOT the Go zero Time, so no aliasing occurs).
//
// DEVIATION 1 (creation/access times): java.util.zip.ZipEntry decodes the Info-ZIP extended
// timestamp (0x5455) and NTFS (0x000a) extra fields into its mtime/atime/ctime fields; Go's
// archive/zip decodes only the modification time and drops the other two. Since DSSZipEntry
// exposes getCreationTime()/getLastAccessTime(), the same extra-field decoding
// (java.util.zip.ZipEntry#setExtra0) is reproduced here, over the preserved raw extra bytes.
//
// DEVIATION 2 (createZipEntry): archive/zip.FileHeader has no creation-time field, so
// createZipEntry()'s `if (creationTime != null) zipEntry.setCreationTime(creationTime)` has no
// counterpart. In Java that call makes ZipOutputStream append a *synthesized* extended-timestamp
// extra field to the entry it writes - on top of the extra bytes copied by setExtra(extra), from
// which the creation time was decoded in the first place. This port copies those extra bytes
// verbatim and synthesizes nothing, so the creation time still round-trips, just without Java's
// duplicated encoding. Neither form is covered by any signature.
//
// hashCode() is dropped (nothing in the ported tree keys a hash container on DSSZipEntry);
// java.io.Serializable is dropped silently.
package asic

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"time"
)

// dssZipEntryWindowsTimeNotAvailable is java.util.zip.ZipUtils.WINDOWS_TIME_NOT_AVAILABLE.
const dssZipEntryWindowsTimeNotAvailable uint64 = 0

// dssZipEntryExtIDNTFS / dssZipEntryExtIDEXTT are java.util.zip.ZipConstants64.EXTID_NTFS and
// EXTID_EXTT.
const (
	dssZipEntryExtIDNTFS = 0x000a
	dssZipEntryExtIDEXTT = 0x5455
)

// DSSZipEntry contains metadata for a ZIP-container entry.
type DSSZipEntry struct {
	// name is the ZIP entry name.
	name string

	// comment is the comment for the ZIP entry. Java's null is spelled "" here: a zero-length
	// ZIP comment and an absent one are indistinguishable both on the wire and after a
	// read-back, so the distinction Java's `if (comment != null)` guard draws is not
	// observable.
	comment string

	// compressionMethod is the compression method to be used for the current document within a
	// ZIP-container to be created.
	//
	// Default : DEFLATED (8) - the entry will be compressed
	compressionMethod int

	// creationTime is the time indicating when the document has been created.
	creationTime time.Time

	// extra contains extra metadata for the ZIP entry.
	extra []byte

	// modificationTime is the time of the last entry modification.
	modificationTime time.Time

	// lastAccessTime is the time of the last entry access.
	lastAccessTime time.Time

	// size is the size of the document.
	size int64

	// compressedSize is the size of the document after compression.
	compressedSize int64

	// crc is the CRC-32 hash of the uncompressed document.
	crc int64
}

// NewDSSZipEntry is the default constructor. Port of DSSZipEntry(String).
//
// Java's Objects.requireNonNull(name, "Name cannot be null!") has no Go trigger (a string is
// never nil), so no panic is possible here.
func NewDSSZipEntry(name string) *DSSZipEntry {
	return &DSSZipEntry{
		name:              name,
		compressionMethod: int(zip.Deflate),
	}
}

// NewDSSZipEntryFromFileHeader constructs a DSSZipEntry from an existing archive/zip.FileHeader.
// Port of DSSZipEntry(ZipEntry).
//
// Panics with the Java message when fileHeader is nil (Objects.requireNonNull).
func NewDSSZipEntryFromFileHeader(fileHeader *zip.FileHeader) *DSSZipEntry {
	if fileHeader == nil {
		panic("ZipEntry cannot be null!")
	}
	e := &DSSZipEntry{
		name:              fileHeader.Name,
		compressionMethod: int(fileHeader.Method),
	}
	if fileHeader.Comment != "" {
		e.comment = fileHeader.Comment
	}
	if len(fileHeader.Extra) != 0 {
		e.extra = fileHeader.Extra
	}
	// java.util.zip.ZipEntry decodes these three from the extra field; archive/zip only
	// decodes the modification time, so the other two are recovered here (DEVIATION 1).
	mtime, atime, ctime := dssZipEntryTimesFromExtra(fileHeader.Extra)
	if !ctime.IsZero() {
		e.creationTime = ctime
	}
	// ZipEntry#getLastModifiedTime() returns the decoded mtime when the extra field carries
	// one, and otherwise falls back to the MS-DOS timestamp - which is exactly what
	// archive/zip already folds into FileHeader.Modified.
	if !mtime.IsZero() {
		e.modificationTime = mtime
	} else if !fileHeader.Modified.IsZero() {
		e.modificationTime = fileHeader.Modified
	}
	if !atime.IsZero() {
		e.lastAccessTime = atime
	}
	e.size = int64(fileHeader.UncompressedSize64)
	e.compressedSize = int64(fileHeader.CompressedSize64)
	e.crc = int64(fileHeader.CRC32)
	return e
}

// Name gets the name of the ZIP entry. Port of getName().
func (e *DSSZipEntry) Name() string { return e.name }

// SetName sets the name of the ZIP entry. Port of setName(String).
func (e *DSSZipEntry) SetName(name string) { e.name = name }

// Comment gets the comment defined for the ZIP entry. Port of getComment().
func (e *DSSZipEntry) Comment() string { return e.comment }

// SetComment sets the comment defined for the ZIP entry. Port of setComment(String).
func (e *DSSZipEntry) SetComment(comment string) { e.comment = comment }

// CompressionMethod gets the compression method for the ZIP entry. Port of
// getCompressionMethod().
func (e *DSSZipEntry) CompressionMethod() int { return e.compressionMethod }

// SetCompressionMethod sets the compression method for the ZIP entry. Port of
// setCompressionMethod(int).
func (e *DSSZipEntry) SetCompressionMethod(compressionMethod int) {
	e.compressionMethod = compressionMethod
}

// CreationTime gets the creation time of the document. Port of getCreationTime().
func (e *DSSZipEntry) CreationTime() time.Time { return e.creationTime }

// SetCreationTime sets the creation time of the document. Ports both
// setCreationTime(FileTime) and setCreationTime(Date), which collapse onto one method in Go.
func (e *DSSZipEntry) SetCreationTime(creationTime time.Time) { e.creationTime = creationTime }

// Extra gets the extra field for the document. Port of getExtra().
func (e *DSSZipEntry) Extra() []byte { return e.extra }

// SetExtra sets the extra field for the document. Port of setExtra(byte[]).
func (e *DSSZipEntry) SetExtra(extra []byte) { e.extra = extra }

// ModificationTime gets the last modification time of the document. Port of
// getModificationTime().
func (e *DSSZipEntry) ModificationTime() time.Time { return e.modificationTime }

// LastAccessTime gets the last access time of the document. Port of getLastAccessTime().
func (e *DSSZipEntry) LastAccessTime() time.Time { return e.lastAccessTime }

// Size gets the size of the uncompressed document. Port of getSize().
func (e *DSSZipEntry) Size() int64 { return e.size }

// CompressedSize gets the size of the compressed document. Port of getCompressedSize().
func (e *DSSZipEntry) CompressedSize() int64 { return e.compressedSize }

// Crc gets the CRC-32 checksum of the uncompressed document. Port of getCrc().
func (e *DSSZipEntry) Crc() int64 { return e.crc }

// CreateZipEntry creates a new copy of the underlying ZIP entry header.
//
// NOTE: some fields are not copied, as they can be changed during container creation (i.e.
// modification time, size, crc, etc.).
//
// Port of createZipEntry(); see DEVIATION 2 in this file's header regarding the creation time.
func (e *DSSZipEntry) CreateZipEntry() *zip.FileHeader {
	fileHeader := &zip.FileHeader{Name: e.name}
	if e.comment != "" {
		fileHeader.Comment = e.comment
	}
	if e.compressionMethod != -1 {
		fileHeader.Method = uint16(e.compressionMethod)
	}
	if e.extra != nil {
		fileHeader.Extra = e.extra
	}
	return fileHeader
}

// Equals ports equals(Object).
func (e *DSSZipEntry) Equals(other *DSSZipEntry) bool {
	if e == other {
		return true
	}
	if e == nil || other == nil {
		return false
	}
	return e.compressionMethod == other.compressionMethod &&
		e.size == other.size &&
		e.compressedSize == other.compressedSize &&
		e.crc == other.crc &&
		e.name == other.name &&
		e.comment == other.comment &&
		e.creationTime.Equal(other.creationTime) &&
		bytes.Equal(e.extra, other.extra) &&
		e.modificationTime.Equal(other.modificationTime) &&
		e.lastAccessTime.Equal(other.lastAccessTime)
}

// dssZipEntryTimesFromExtra decodes the modification, access and creation times a ZIP extra field
// carries, returning the zero Time for any that is absent. It reproduces the EXTID_NTFS and
// EXTID_EXTT branches of java.util.zip.ZipEntry#setExtra0 - the decoding DSSZipEntry(ZipEntry)
// relies on having already happened - including their length guards and their
// WINDOWS_TIME_NOT_AVAILABLE / per-flag-bit skips. The ZIP64 branch is not reproduced: the sizes
// it patches up are already resolved by the reader before a FileHeader reaches this file. One
// guard is added: a zero-length EXTT payload makes Java read one byte past the field and throw
// ArrayIndexOutOfBoundsException, which here would panic, so such a field is skipped instead.
func dssZipEntryTimesFromExtra(extra []byte) (mtime, atime, ctime time.Time) {
	off := 0
	length := len(extra)
	for off+4 <= length {
		tag := int(binary.LittleEndian.Uint16(extra[off:]))
		sz := int(binary.LittleEndian.Uint16(extra[off+2:]))
		off += 4
		if off+sz > length {
			break
		}
		switch tag {
		case dssZipEntryExtIDNTFS:
			if sz < 32 {
				break
			}
			pos := off + 4 // reserved 4 bytes
			if binary.LittleEndian.Uint16(extra[pos:]) != 0x0001 ||
				binary.LittleEndian.Uint16(extra[pos+2:]) != 24 {
				break
			}
			if wtime := binary.LittleEndian.Uint64(extra[pos+4:]); wtime != dssZipEntryWindowsTimeNotAvailable {
				mtime = dssZipEntryWinTimeToTime(wtime)
			}
			if wtime := binary.LittleEndian.Uint64(extra[pos+12:]); wtime != dssZipEntryWindowsTimeNotAvailable {
				atime = dssZipEntryWinTimeToTime(wtime)
			}
			if wtime := binary.LittleEndian.Uint64(extra[pos+20:]); wtime != dssZipEntryWindowsTimeNotAvailable {
				ctime = dssZipEntryWinTimeToTime(wtime)
			}
		case dssZipEntryExtIDEXTT:
			if sz < 1 {
				break
			}
			flag := int(extra[off])
			sz0 := 1
			// The CEN-header extra field contains the modification time only, or no
			// timestamp at all.
			if flag&0x1 != 0 && sz0+4 <= sz {
				mtime = dssZipEntryUnixTimeToTime(int32(binary.LittleEndian.Uint32(extra[off+sz0:])))
				sz0 += 4
			}
			if flag&0x2 != 0 && sz0+4 <= sz {
				atime = dssZipEntryUnixTimeToTime(int32(binary.LittleEndian.Uint32(extra[off+sz0:])))
				sz0 += 4
			}
			if flag&0x4 != 0 && sz0+4 <= sz {
				ctime = dssZipEntryUnixTimeToTime(int32(binary.LittleEndian.Uint32(extra[off+sz0:])))
			}
		}
		off += sz
	}
	return mtime, atime, ctime
}

// dssZipEntryWinTimeToTime converts a Windows FILETIME (100ns ticks since 1601-01-01 UTC) to a
// time.Time. Port of java.util.zip.ZipUtils#winTimeToFileTime.
func dssZipEntryWinTimeToTime(wtime uint64) time.Time {
	const ticksPerSecond = 10000000
	// 11644473600 = seconds between 1601-01-01 and 1970-01-01.
	const epochDiffSeconds = 11644473600
	sec := int64(wtime/ticksPerSecond) - epochDiffSeconds
	nsec := int64(wtime%ticksPerSecond) * 100
	return time.Unix(sec, nsec).UTC()
}

// dssZipEntryUnixTimeToTime converts a signed 32-bit Unix timestamp to a time.Time. Port of
// java.util.zip.ZipUtils#unixTimeToFileTime.
func dssZipEntryUnixTimeToTime(utime int32) time.Time {
	return time.Unix(int64(utime), 0).UTC()
}
