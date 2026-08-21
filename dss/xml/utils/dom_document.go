// Ported from dss-xml-utils/src/main/java/eu/europa/esig/dss/xml/utils/DOMDocument.java (DSS 6.5.RC1).
//
// DEVIATION: upstream embeds eu.europa.esig.dss.model.CommonDocument for its default
// openStream/getName/getMimeType/save/writeTo/getDigest/getDigestValue plumbing. That
// plumbing's Go port (model.CommonDocument's commonDocumentWriteTo/commonDocumentDigest/...)
// is unexported, because every existing model.DSSDocument implementation
// (InMemoryDocument/FileDocument/DigestDocument) lives inside package model itself and calls
// it directly. DOMDocument is the first implementation outside package model, so this file
// carries its own copy of the same logic (digest caching via utils.OrderedMap, hashing via
// spi.DSSUtilsMessageDigest) instead of exporting model's internals for a single external
// caller. Flagged for integrator: consider exporting a CommonDocument helper surface from
// model if further DSSDocument implementations accumulate outside that package.
package utils

import (
	"bytes"
	"io"
	"os"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// DOMDocument allows handling of an *xmldom.Node as a model.DSSDocument. The class handles
// the node in memory, and reads its data only on request.
type DOMDocument struct {
	// node is the Node defining the document.
	node *xmldom.Node

	// bytesData caches the bytes of the node's serialization.
	bytesData []byte

	name      string
	mimeType  enumerations.MimeType
	digestMap *utils.OrderedMap[enumerations.DigestAlgorithm, []byte]
}

var _ model.DSSDocument = (*DOMDocument)(nil)

// NewDOMDocument creates a DOMDocument from node. NOTE: uses enumerations.MimeTypeEnum_XML by
// default. Panics if node is nil (Java Objects.requireNonNull("Element cannot be null")).
// Ports DOMDocument(Node).
func NewDOMDocument(node *xmldom.Node) *DOMDocument {
	return NewDOMDocumentWithName(node, "")
}

// NewDOMDocumentWithName creates a DOMDocument from node with the given name. NOTE: uses
// enumerations.MimeTypeEnum_XML by default. Panics if node is nil. Ports
// DOMDocument(Node, String).
func NewDOMDocumentWithName(node *xmldom.Node, name string) *DOMDocument {
	if node == nil {
		panic("Element cannot be null")
	}
	return &DOMDocument{node: node, name: name, mimeType: enumerations.MimeTypeEnum_XML}
}

// Node gets the Node used to define the document. Ports getNode().
func (d *DOMDocument) Node() *xmldom.Node {
	return d.node
}

// OpenStream ports openStream().
func (d *DOMDocument) OpenStream() (io.ReadCloser, error) {
	b, err := d.getBytes()
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(b)), nil
}

// getBytes gets the cached binary array as the result of the Node's serialization. Ports the
// protected getBytes().
func (d *DOMDocument) getBytes() ([]byte, error) {
	if d.bytesData == nil {
		b, err := DomUtilsSerializeNode(d.node)
		if err != nil {
			return nil, err
		}
		d.bytesData = b
	}
	return d.bytesData, nil
}

// Name ports CommonDocument#getName.
func (d *DOMDocument) Name() string { return d.name }

// SetName ports CommonDocument#setName.
func (d *DOMDocument) SetName(name string) { d.name = name }

// MimeType ports CommonDocument#getMimeType.
func (d *DOMDocument) MimeType() enumerations.MimeType { return d.mimeType }

// SetMimeType ports CommonDocument#setMimeType.
func (d *DOMDocument) SetMimeType(mimeType enumerations.MimeType) { d.mimeType = mimeType }

// String ports CommonDocument#toString.
func (d *DOMDocument) String() string {
	mimeTypeString := ""
	if d.mimeType != nil {
		mimeTypeString = d.mimeType.MimeTypeString()
	}
	return "Name: " + d.name + " / MimeType: " + mimeTypeString
}

// WriteTo ports CommonDocument#writeTo for DOMDocument.
func (d *DOMDocument) WriteTo(w io.Writer) (int64, error) {
	b, err := d.getBytes()
	if err != nil {
		return 0, err
	}
	n, err := w.Write(b)
	return int64(n), err
}

// Save ports CommonDocument#save for DOMDocument.
func (d *DOMDocument) Save(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = d.WriteTo(f)
	return err
}

// Digest ports CommonDocument#getDigest for DOMDocument.
func (d *DOMDocument) Digest(digestAlgorithm enumerations.DigestAlgorithm) (model.Digest, error) {
	v, err := d.DigestValue(digestAlgorithm)
	if err != nil {
		return model.Digest{}, err
	}
	return model.NewDigest(digestAlgorithm, v), nil
}

// DigestValue ports CommonDocument#getDigestValue for DOMDocument.
func (d *DOMDocument) DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error) {
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
	b, err := d.getBytes()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause("Unable to compute the digest", err)
	}
	h.Write(b)
	digest := h.Sum(nil)
	d.digestMap.Set(digestAlgorithm, digest)
	return digest, nil
}

// Equals ports DOMDocument#equals. Panics if serialization fails (Java lets the
// RuntimeException from getBytes() propagate out of equals() uncaught).
func (d *DOMDocument) Equals(other *DOMDocument) bool {
	if d == other {
		return true
	}
	if other == nil {
		return false
	}
	if d.mimeType != other.mimeType || d.name != other.name {
		return false
	}
	b1, err := d.getBytes()
	if err != nil {
		panic(err)
	}
	b2, err := other.getBytes()
	if err != nil {
		panic(err)
	}
	return bytes.Equal(b1, b2)
}
