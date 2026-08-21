// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PdfByteRangeDocument.java
// (DSS 6.5.RC1).
//
// DEVIATION: upstream extends eu.europa.esig.dss.model.CommonDocument for its default
// getName/getMimeType/save/writeTo/getDigest/getDigestValue plumbing. That plumbing's Go port
// (model.CommonDocument's commonDocumentWriteTo/commonDocumentDigest/...) is unexported, because
// every existing model.DSSDocument implementation (InMemoryDocument/FileDocument/DigestDocument)
// lives inside package model itself and calls it directly. xml/utils/dom_document.go was the
// first implementation outside package model and carries its own copy of the same logic for the
// same reason (see its header); this file follows the same precedent, streaming through
// OpenStream() rather than buffering into memory (unlike DOMDocument's getBytes()-based digest),
// matching this type's own stated purpose of "reduc[ing] memory overloading during the
// execution".
//
// FORWARD DEPENDENCY (not in this chunk's manifest): ByteRangeInputStream
// (eu.europa.esig.dss.pades.validation.ByteRangeInputStream), an io.Reader-shaped wrapper that
// reads only the two spans a ByteRange covers out of an underlying stream. Assumed shape,
// inferred from the single call this file makes to it:
//
//	func NewByteRangeInputStream(wrapped io.ReadCloser, byteRange *ByteRange) io.ReadCloser
package pades

import (
	"bytes"
	"io"
	"os"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// PdfByteRangeDocument is an internal representation of a PDF document. Used to reduce memory
// overloading during the execution.
type PdfByteRangeDocument struct {
	// pdfDocument is the input PDF document to read.
	pdfDocument model.DSSDocument

	// byteRange is the ByteRange of the revision to be read.
	byteRange *ByteRange

	name      string
	mimeType  enumerations.MimeType
	digestMap *utils.OrderedMap[enumerations.DigestAlgorithm, []byte]
}

var _ model.DSSDocument = (*PdfByteRangeDocument)(nil)

// NewPdfByteRangeDocument creates a PdfByteRangeDocument. Port of the constructor
// PdfByteRangeDocument(DSSDocument, ByteRange).
//
// Panics with the Java messages when either argument is missing (Objects.requireNonNull).
func NewPdfByteRangeDocument(pdfDocument model.DSSDocument, byteRange *ByteRange) *PdfByteRangeDocument {
	if pdfDocument == nil {
		panic("PdfDocument cannot be null!")
	}
	if byteRange == nil {
		panic("ByteRange cannot be null!")
	}
	return &PdfByteRangeDocument{pdfDocument: pdfDocument, byteRange: byteRange}
}

// ByteRange returns the ByteRange of the document. Port of getByteRange().
func (d *PdfByteRangeDocument) ByteRange() *ByteRange {
	return d.byteRange
}

// OpenStream ports openStream(): only the two spans byteRange covers are readable.
func (d *PdfByteRangeDocument) OpenStream() (io.ReadCloser, error) {
	wrapped, err := d.pdfDocument.OpenStream()
	if err != nil {
		return nil, err
	}
	return NewByteRangeInputStream(wrapped, d.byteRange), nil
}

// Name ports CommonDocument#getName.
func (d *PdfByteRangeDocument) Name() string { return d.name }

// SetName ports CommonDocument#setName.
func (d *PdfByteRangeDocument) SetName(name string) { d.name = name }

// MimeType ports CommonDocument#getMimeType.
func (d *PdfByteRangeDocument) MimeType() enumerations.MimeType { return d.mimeType }

// SetMimeType ports CommonDocument#setMimeType.
func (d *PdfByteRangeDocument) SetMimeType(mimeType enumerations.MimeType) { d.mimeType = mimeType }

// String ports CommonDocument#toString.
func (d *PdfByteRangeDocument) String() string {
	mimeTypeString := ""
	if d.mimeType != nil {
		mimeTypeString = d.mimeType.MimeTypeString()
	}
	return "Name: " + d.name + " / MimeType: " + mimeTypeString
}

// WriteTo ports CommonDocument#writeTo for PdfByteRangeDocument.
func (d *PdfByteRangeDocument) WriteTo(w io.Writer) (int64, error) {
	rc, err := d.OpenStream()
	if err != nil {
		return 0, err
	}
	defer rc.Close()
	return io.Copy(w, rc)
}

// Save ports CommonDocument#save for PdfByteRangeDocument.
func (d *PdfByteRangeDocument) Save(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = d.WriteTo(f)
	return err
}

// Digest ports CommonDocument#getDigest for PdfByteRangeDocument.
func (d *PdfByteRangeDocument) Digest(digestAlgorithm enumerations.DigestAlgorithm) (model.Digest, error) {
	v, err := d.DigestValue(digestAlgorithm)
	if err != nil {
		return model.Digest{}, err
	}
	return model.NewDigest(digestAlgorithm, v), nil
}

// DigestValue ports CommonDocument#getDigestValue for PdfByteRangeDocument.
func (d *PdfByteRangeDocument) DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
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
func (d *PdfByteRangeDocument) Equals(other *PdfByteRangeDocument) bool {
	if d == other {
		return true
	}
	if other == nil {
		return false
	}
	if d.mimeType != other.mimeType || d.name != other.name {
		return false
	}
	return pdfByteRangeDocumentContentsEqual(d.pdfDocument, other.pdfDocument) && d.byteRange.Equals(other.byteRange)
}

// pdfByteRangeDocumentContentsEqual stands in for Java's polymorphic DSSDocument#equals()
// dispatch used by pdfDocument.equals(that.pdfDocument): unlike model.Token, model.DSSDocument
// declares no Equals method in its interface, since every concrete implementation's Equals lives
// on its own concrete type (InMemoryDocument.Equals, FileDocument.Equals, DigestDocument.Equals,
// ...) rather than being expressible generically, and pdfDocument here is typically an arbitrary
// caller-supplied document this package cannot exhaustively type-switch over. Byte-for-byte
// stream comparison reproduces the same observable result equals() is used for.
func pdfByteRangeDocumentContentsEqual(a, b model.DSSDocument) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	ar, err := a.OpenStream()
	if err != nil {
		return false
	}
	defer ar.Close()
	br, err := b.OpenStream()
	if err != nil {
		return false
	}
	defer br.Close()
	aBytes, err := io.ReadAll(ar)
	if err != nil {
		return false
	}
	bBytes, err := io.ReadAll(br)
	if err != nil {
		return false
	}
	return bytes.Equal(aBytes, bBytes)
}
