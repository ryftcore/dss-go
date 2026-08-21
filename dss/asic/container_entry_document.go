// Ported from
// dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/ContainerEntryDocument.java
// (DSS 6.5.RC1).
//
// DEVIATION: upstream extends eu.europa.esig.dss.model.CommonDocument for its default
// getName/getMimeType/save/writeTo/getDigest/getDigestValue plumbing, whose Go port
// (model.CommonDocument's commonDocumentWriteTo/commonDocumentDigest/...) is unexported because
// every model.DSSDocument implementation inside package model calls it directly. This file
// follows the precedent set by xml/utils/dom_document.go and pades/pdf_byte_range_document.go and
// carries its own copy of that plumbing.
//
// hashCode() is dropped (nothing in the ported tree keys a hash container on this type).
package asic

import (
	"io"
	"os"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ContainerEntryDocument represents an entry within a ZIP archive, containing its metadata and
// file's content. This class can be used to create a file entry to be incorporated within an ASiC
// container with customized ZipEntry metadata (e.g. creation time, compression method, etc.).
type ContainerEntryDocument struct {
	// content is the document representing content of a file to be embedded into ZIP container.
	content model.DSSDocument

	// zipEntry contains metadata about a file within ZIP archive.
	zipEntry *DSSZipEntry

	name      string
	mimeType  enumerations.MimeType
	digestMap *utils.OrderedMap[enumerations.DigestAlgorithm, []byte]
}

var _ DSSZipEntryDocument = (*ContainerEntryDocument)(nil)

// NewContainerEntryDocument is the default constructor. Port of
// ContainerEntryDocument(DSSDocument).
//
// Panics with the Java messages when content is nil or unnamed (Objects.requireNonNull).
func NewContainerEntryDocument(content model.DSSDocument) *ContainerEntryDocument {
	if content == nil {
		panic("Document content cannot be null!")
	}
	if content.Name() == "" {
		panic("Document shall contain name!")
	}
	return &ContainerEntryDocument{
		content:  content,
		zipEntry: NewDSSZipEntry(content.Name()),
		name:     content.Name(),
		mimeType: content.MimeType(),
	}
}

// NewContainerEntryDocumentWithZipEntry is the constructor with a provided DSSZipEntry. Port of
// ContainerEntryDocument(DSSDocument, DSSZipEntry).
//
// Panics with the Java messages when content is nil, unnamed, or zipEntry is nil
// (Objects.requireNonNull), and returns an error when the names disagree
// (IllegalArgumentException).
func NewContainerEntryDocumentWithZipEntry(content model.DSSDocument, zipEntry *DSSZipEntry) (*ContainerEntryDocument, error) {
	if content == nil {
		panic("Document content cannot be null!")
	}
	if content.Name() == "" {
		panic("Document shall contain name!")
	}
	if zipEntry == nil {
		panic("ZipEntry cannot be null!")
	}
	if content.Name() != zipEntry.Name() {
		return nil, model.NewDSSError("Name of the document shall match the name of ZipEntry!")
	}
	return &ContainerEntryDocument{
		content:  content,
		zipEntry: zipEntry,
		name:     content.Name(),
		mimeType: content.MimeType(),
	}, nil
}

// OpenStream ports openStream().
func (d *ContainerEntryDocument) OpenStream() (io.ReadCloser, error) {
	return d.content.OpenStream()
}

// Name ports CommonDocument#getName.
func (d *ContainerEntryDocument) Name() string { return d.name }

// SetName ports setName(String): the name of the wrapped ZIP entry is kept in sync.
func (d *ContainerEntryDocument) SetName(name string) {
	d.name = name
	d.zipEntry.SetName(name)
}

// MimeType ports CommonDocument#getMimeType.
func (d *ContainerEntryDocument) MimeType() enumerations.MimeType { return d.mimeType }

// SetMimeType ports CommonDocument#setMimeType.
func (d *ContainerEntryDocument) SetMimeType(mimeType enumerations.MimeType) {
	d.mimeType = mimeType
}

// ZipEntry ports getZipEntry().
func (d *ContainerEntryDocument) ZipEntry() *DSSZipEntry { return d.zipEntry }

// String ports CommonDocument#toString.
func (d *ContainerEntryDocument) String() string {
	mimeTypeString := ""
	if d.mimeType != nil {
		mimeTypeString = d.mimeType.MimeTypeString()
	}
	return "Name: " + d.name + " / MimeType: " + mimeTypeString
}

// WriteTo ports CommonDocument#writeTo for ContainerEntryDocument.
func (d *ContainerEntryDocument) WriteTo(w io.Writer) (int64, error) {
	rc, err := d.OpenStream()
	if err != nil {
		return 0, err
	}
	defer rc.Close()
	return io.Copy(w, rc)
}

// Save ports CommonDocument#save for ContainerEntryDocument.
func (d *ContainerEntryDocument) Save(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = d.WriteTo(f)
	return err
}

// Digest ports CommonDocument#getDigest for ContainerEntryDocument.
func (d *ContainerEntryDocument) Digest(digestAlgorithm enumerations.DigestAlgorithm) (model.Digest, error) {
	v, err := d.DigestValue(digestAlgorithm)
	if err != nil {
		return model.Digest{}, err
	}
	return model.NewDigest(digestAlgorithm, v), nil
}

// DigestValue ports CommonDocument#getDigestValue for ContainerEntryDocument.
func (d *ContainerEntryDocument) DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
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

// Content returns the wrapped content document. It has no Java counterpart (upstream keeps the
// field private and only reads it from equals()); it exists so that Equals below - and the
// ASiC merge/extension code in the format-specific packages, which needs to re-wrap an entry's
// payload - can reach the payload without re-reading the stream.
func (d *ContainerEntryDocument) Content() model.DSSDocument { return d.content }

// Equals ports equals(Object). Java's content.equals(that.content) dispatches to the concrete
// DSSDocument implementation; model.DSSDocument declares no Equals in its interface (see the
// note in pades/pdf_byte_range_document.go), so reference identity is compared instead.
func (d *ContainerEntryDocument) Equals(other *ContainerEntryDocument) bool {
	if d == other {
		return true
	}
	if d == nil || other == nil {
		return false
	}
	if d.mimeType != other.mimeType || d.name != other.name {
		return false
	}
	return d.content == other.content && d.zipEntry.Equals(other.zipEntry)
}
