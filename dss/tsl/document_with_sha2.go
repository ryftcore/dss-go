// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/sha2/DocumentWithSha2.java (DSS 6.5.RC1).
//
// DEVIATION: upstream extends eu.europa.esig.dss.model.CommonDocument for its default
// getName/getMimeType/save/writeTo/getDigest/getDigestValue plumbing, whose Go port
// (model.CommonDocument's commonDocumentWriteTo/commonDocumentDigest/...) is unexported because
// every model.DSSDocument implementation inside package model calls it directly. This file
// follows the precedent set by asic/container_entry_document.go and
// pades/pdf_byte_range_document.go and carries its own copy of that plumbing.
//
// hashCode() is dropped (nothing in the ported tree keys a hash container on this type).
package tsl

import (
	"io"
	"os"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
)

// DocumentWithSha2 represents a downloaded model.DSSDocument together with its corresponding
// ".sha2" file.
type DocumentWithSha2 struct {
	// document is the original downloaded document.
	document model.DSSDocument

	// sha2Document is the corresponding sha2 document, containing digests of the document.
	sha2Document model.DSSDocument

	// errors is the list of errors occurred during .sha2 document processing.
	errors []string

	name      string
	mimeType  enumerations.MimeType
	digestMap *utils.OrderedMap[enumerations.DigestAlgorithm, []byte]
}

var _ model.DSSDocument = (*DocumentWithSha2)(nil)

// NewDocumentWithSha2 is the default constructor. Port of the protected
// DocumentWithSha2(DSSDocument, DSSDocument); Go has no protected visibility, so the
// constructor is exported - Sha2FileCacheDataLoader#mergeDocumentWithSha2, upstream's only
// caller, lives in the same Go package but a caller outside it can now build one too.
func NewDocumentWithSha2(document, sha2Document model.DSSDocument) *DocumentWithSha2 {
	return &DocumentWithSha2{document: document, sha2Document: sha2Document}
}

// Document gets the original document. Port of getDocument().
func (d *DocumentWithSha2) Document() model.DSSDocument {
	return d.document
}

// Sha2Document gets the downloaded sha2 document corresponding to the document. Port of
// getSha2Document().
func (d *DocumentWithSha2) Sha2Document() model.DSSDocument {
	return d.sha2Document
}

// AddErrorMessage adds an error message occurred during the .sha2 file validation. Port of the
// protected addErrorMessage(String); exported for the same reason as the constructor.
func (d *DocumentWithSha2) AddErrorMessage(errorMessage string) {
	d.errors = append(d.errors, errorMessage)
}

// Errors returns the list of errors occurred during processing of the .sha2 document. Port of
// getErrors(), which lazily instantiates the list and therefore never answers null; Go's nil
// slice is indistinguishable from an empty one for len/range/append, so no lazy instantiation
// is needed and callers observe an empty list either way.
//
// NOTE: the Javadoc claims "Returns NULL if validation succeeded", which the implementation
// contradicts (it always materialises the list). The implementation is what is ported.
func (d *DocumentWithSha2) Errors() []string {
	return d.errors
}

// OpenStream ports openStream(), which delegates to the wrapped document.
//
// Panics with the Java message when the document is nil (Objects.requireNonNull).
func (d *DocumentWithSha2) OpenStream() (io.ReadCloser, error) {
	if d.document == nil {
		panic("Document is null! Unable to open InputStream.")
	}
	return d.document.OpenStream()
}

// Name ports CommonDocument#getName.
func (d *DocumentWithSha2) Name() string { return d.name }

// SetName ports CommonDocument#setName.
func (d *DocumentWithSha2) SetName(name string) { d.name = name }

// MimeType ports CommonDocument#getMimeType.
func (d *DocumentWithSha2) MimeType() enumerations.MimeType { return d.mimeType }

// SetMimeType ports CommonDocument#setMimeType.
func (d *DocumentWithSha2) SetMimeType(mimeType enumerations.MimeType) { d.mimeType = mimeType }

// String ports CommonDocument#toString.
func (d *DocumentWithSha2) String() string {
	mimeTypeString := ""
	if d.mimeType != nil {
		mimeTypeString = d.mimeType.MimeTypeString()
	}
	return "Name: " + d.name + " / MimeType: " + mimeTypeString
}

// WriteTo ports CommonDocument#writeTo for DocumentWithSha2.
func (d *DocumentWithSha2) WriteTo(w io.Writer) (int64, error) {
	rc, err := d.OpenStream()
	if err != nil {
		return 0, err
	}
	defer rc.Close()
	return io.Copy(w, rc)
}

// Save ports CommonDocument#save for DocumentWithSha2.
func (d *DocumentWithSha2) Save(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = d.WriteTo(f)
	return err
}

// Digest ports CommonDocument#getDigest for DocumentWithSha2.
func (d *DocumentWithSha2) Digest(digestAlgorithm enumerations.DigestAlgorithm) (model.Digest, error) {
	v, err := d.DigestValue(digestAlgorithm)
	if err != nil {
		return model.Digest{}, err
	}
	return model.NewDigest(digestAlgorithm, v), nil
}

// DigestValue ports CommonDocument#getDigestValue for DocumentWithSha2.
func (d *DocumentWithSha2) DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
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

// Equals ports equals(Object): CommonDocument's name/mimeType, then the two wrapped documents
// and the error list.
//
// Java's document.equals(that.document) dispatches to the concrete DSSDocument implementation;
// model.DSSDocument declares no Equals in its interface (see the note in
// asic/container_entry_document.go), so reference identity is compared for the two wrapped
// documents instead.
func (d *DocumentWithSha2) Equals(other *DocumentWithSha2) bool {
	if d == other {
		return true
	}
	if d == nil || other == nil {
		return false
	}
	if d.mimeType != other.mimeType || d.name != other.name {
		return false
	}
	if d.document != other.document || d.sha2Document != other.sha2Document {
		return false
	}
	if (d.errors == nil) != (other.errors == nil) || len(d.errors) != len(other.errors) {
		return false
	}
	for index := range d.errors {
		if d.errors[index] != other.errors[index] {
			return false
		}
	}
	return true
}
