// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/FileArchiveEntry.java
// (DSS 6.5.RC1).
//
// DEVIATION 1: upstream extends eu.europa.esig.dss.model.CommonDocument; that plumbing's Go port
// is unexported (see container_entry_document.go's header for the full note), so it is carried
// here too.
//
// DEVIATION 2: upstream keeps the original java.util.zip.ZipEntry instance and resolves the entry
// through ZipFile#getInputStream(ZipEntry), which - because the instance came from a
// ZipInputStream pass and not from that ZipFile - degrades to a lookup by entry name. This port
// therefore stores the name and looks the entry up by name in the central directory, taking the
// first match, which is the same observable resolution.
//
// hashCode() is dropped (nothing in the ported tree keys a hash container on this type).
package asic

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
)

// FileArchiveEntry is an internal type that is used for performance purposes, accessing
// ZIP-archive entries on request, instead of loading all files into memory.
type FileArchiveEntry struct {
	// zipArchive is the file system document representing a ZIP-container.
	zipArchive *model.FileDocument

	// entryName represents the original entry name : important to access the relevant entry
	// from the archive (see DEVIATION 2).
	entryName string

	// dssZipEntry contains metadata about the extracted entry.
	dssZipEntry *DSSZipEntry

	name      string
	mimeType  enumerations.MimeType
	digestMap *utils.OrderedMap[enumerations.DigestAlgorithm, []byte]
}

var _ DSSZipEntryDocument = (*FileArchiveEntry)(nil)

// newFileArchiveEntry is the default constructor. Port of the protected
// FileArchiveEntry(FileDocument, ZipEntry).
//
// Panics with the Java messages when either argument is missing (Objects.requireNonNull).
func newFileArchiveEntry(zipArchive *model.FileDocument, fileHeader *zip.FileHeader) *FileArchiveEntry {
	if zipArchive == nil {
		panic("ZIP Archive cannot be null!")
	}
	if fileHeader == nil {
		panic("ZIP Entry cannot be null!")
	}
	dssZipEntry := NewDSSZipEntryFromFileHeader(fileHeader)
	return &FileArchiveEntry{
		zipArchive:  zipArchive,
		entryName:   fileHeader.Name,
		dssZipEntry: dssZipEntry,
		name:        dssZipEntry.Name(),
		mimeType:    enumerations.MimeTypeFromFileName(dssZipEntry.Name()),
	}
}

// OpenStream ports openStream(): the returned stream reads the entry straight out of the archive
// file and owns the archive handle, closing it on Close - the Go equivalent of upstream's
// ZipFileEntryInputStream inner class.
func (d *FileArchiveEntry) OpenStream() (io.ReadCloser, error) {
	readCloser, err := zip.OpenReader(d.zipArchive.Path())
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to create an InputStream", err)
	}
	for _, file := range readCloser.File {
		if file.Name == d.entryName {
			entryReader, err := file.Open()
			if err != nil {
				_ = readCloser.Close()
				return nil, model.NewDSSErrorMessageCause("Unable to create an InputStream", err)
			}
			return &zipFileEntryInputStream{entryInputStream: entryReader, zipFile: readCloser}, nil
		}
	}
	_ = readCloser.Close()
	return nil, model.NewDSSErrorMessageCause("Unable to create an InputStream",
		fmt.Errorf("no entry with name '%s' within the archive", d.entryName))
}

// zipFileEntryInputStream creates an InputStream for a ZipEntry from the provided archive file
// and handles closing of the archive. Port of the FileArchiveEntry.ZipFileEntryInputStream inner
// class.
type zipFileEntryInputStream struct {
	// entryInputStream is the InputStream for the given ZIP entry.
	entryInputStream io.ReadCloser

	// zipFile reads the ZIP file in the file system.
	zipFile *zip.ReadCloser
}

// Read ports ZipFileEntryInputStream#read.
func (s *zipFileEntryInputStream) Read(p []byte) (int, error) {
	return s.entryInputStream.Read(p)
}

// Close ports ZipFileEntryInputStream#close.
func (s *zipFileEntryInputStream) Close() error {
	if err := s.entryInputStream.Close(); err != nil {
		_ = s.zipFile.Close()
		return err
	}
	return s.zipFile.Close()
}

// Name ports CommonDocument#getName.
func (d *FileArchiveEntry) Name() string { return d.name }

// SetName ports setName(String): the name of the wrapped ZIP entry is kept in sync.
func (d *FileArchiveEntry) SetName(name string) {
	d.name = name
	d.dssZipEntry.SetName(name)
}

// MimeType ports CommonDocument#getMimeType.
func (d *FileArchiveEntry) MimeType() enumerations.MimeType { return d.mimeType }

// SetMimeType ports CommonDocument#setMimeType.
func (d *FileArchiveEntry) SetMimeType(mimeType enumerations.MimeType) { d.mimeType = mimeType }

// ZipEntry ports getZipEntry().
func (d *FileArchiveEntry) ZipEntry() *DSSZipEntry { return d.dssZipEntry }

// String ports CommonDocument#toString.
func (d *FileArchiveEntry) String() string {
	mimeTypeString := ""
	if d.mimeType != nil {
		mimeTypeString = d.mimeType.MimeTypeString()
	}
	return "Name: " + d.name + " / MimeType: " + mimeTypeString
}

// WriteTo ports CommonDocument#writeTo for FileArchiveEntry.
func (d *FileArchiveEntry) WriteTo(w io.Writer) (int64, error) {
	rc, err := d.OpenStream()
	if err != nil {
		return 0, err
	}
	defer rc.Close()
	return io.Copy(w, rc)
}

// Save ports CommonDocument#save for FileArchiveEntry.
func (d *FileArchiveEntry) Save(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = d.WriteTo(f)
	return err
}

// Digest ports CommonDocument#getDigest for FileArchiveEntry.
func (d *FileArchiveEntry) Digest(digestAlgorithm enumerations.DigestAlgorithm) (model.Digest, error) {
	v, err := d.DigestValue(digestAlgorithm)
	if err != nil {
		return model.Digest{}, err
	}
	return model.NewDigest(digestAlgorithm, v), nil
}

// DigestValue ports CommonDocument#getDigestValue for FileArchiveEntry.
func (d *FileArchiveEntry) DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
	if d.digestMap == nil {
		d.digestMap = utils.NewOrderedMap[enumerations.DigestAlgorithm, []byte]()
	}
	if digest, ok := d.digestMap.Get(digestAlgorithm); ok {
		return digest, nil
	}
	h, err := spi.DSSUtilsMessageDigest(digestAlgorithm)
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to compute the digest", err)
	}
	rc, err := d.OpenStream()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to compute the digest", err)
	}
	defer rc.Close()
	if _, err := io.Copy(h, rc); err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to compute the digest", err)
	}
	digest := h.Sum(nil)
	d.digestMap.Set(digestAlgorithm, digest)
	return digest, nil
}

// Equals ports equals(Object).
func (d *FileArchiveEntry) Equals(other *FileArchiveEntry) bool {
	if d == other {
		return true
	}
	if d == nil || other == nil {
		return false
	}
	if d.mimeType != other.mimeType || d.name != other.name {
		return false
	}
	return d.zipArchive.Equals(other.zipArchive) && d.dssZipEntry.Equals(other.dssZipEntry)
}
