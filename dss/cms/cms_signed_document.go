// Ported from dss-cms/src/main/java/eu/europa/esig/dss/cms/CMSSignedDocument.java (DSS 6.5.RC1).
//
// Java extends CommonDocument and holds a raw org.bouncycastle.cms.CMSSignedData, recomputing
// its bytes on every call to writeTo/getBytes by re-serialising signedData.toASN1Structure()
// with ASN1Encoding.DL - "keep DL to ensure the original order of elements" - which forces
// every length to the definite form while leaving BER vs. DER structural choices (BERSet vs.
// DERSet etc.) untouched. That is exactly what *cms.CMS.DEREncoded() (internal/cmscore's
// element.DEREncoded(), or a built CMS's already-DER encoding) already computes, so this port
// wraps *CMS instead of a BouncyCastle CMSSignedData and delegates to it.
//
// DEVIATION: bytes are computed once at construction (via model.InMemoryDocument) rather than
// recomputed on every call as Java's getBytes()/writeTo() do; both signedData and *CMS are
// immutable once built or parsed, so the two are observationally identical.
package cms

import (
	"io"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// CMSSignedDocument is a document composed by a CMS. Port of the CMSSignedDocument class,
// which extends CommonDocument; the Go port composes model.InMemoryDocument instead, which
// gives it the OpenStream/WriteTo/Save/Digest/DigestValue/Name/MimeType behaviour Java inherits
// from CommonDocument.
type CMSSignedDocument struct {
	*model.InMemoryDocument

	// signedData is the CMS representing the document.
	signedData *CMS
}

var _ model.DSSDocument = (*CMSSignedDocument)(nil)

// NewCMSSignedDocument is the default constructor for CMSSignedDocument.
// Port of CMSSignedDocument(CMS), i.e. CMSSignedDocument(data, null).
//
// Panics with the Java message when data is nil (Objects.requireNonNull).
func NewCMSSignedDocument(data *CMS) *CMSSignedDocument {
	return NewCMSSignedDocumentWithName(data, "")
}

// NewCMSSignedDocumentWithName is the constructor for CMSSignedDocument with a custom document
// name. Port of CMSSignedDocument(CMS, String).
//
// Panics with the Java message when data is nil (Objects.requireNonNull).
func NewCMSSignedDocumentWithName(data *CMS, name string) *CMSSignedDocument {
	if data == nil {
		panic("The CMSSignedData cannot be null")
	}
	inMemory := model.NewInMemoryDocumentWithMimeType(data.DEREncoded(), name, enumerations.MimeTypeEnumPKCS7)
	return &CMSSignedDocument{InMemoryDocument: inMemory, signedData: data}
}

// CMSSignedData gets the CMS. Port of #getCMSSignedData.
func (d *CMSSignedDocument) CMSSignedData() *CMS { return d.signedData }

// Bytes returns the encoded binaries of the CMS. Port of #getBytes.
func (d *CMSSignedDocument) Bytes() []byte { return d.InMemoryDocument.Bytes() }

// OpenStream ports #openStream (inherited effect: opens a stream on Bytes()).
func (d *CMSSignedDocument) OpenStream() (io.ReadCloser, error) {
	return d.InMemoryDocument.OpenStream()
}

// Equals ports #equals: same mimeType/name (CommonDocument#equals) and the same
// SignedData.toASN1Structure(), which for this port means the same DER encoding.
func (d *CMSSignedDocument) Equals(other *CMSSignedDocument) bool {
	if d == other {
		return true
	}
	if other == nil {
		return false
	}
	if d.Name() != other.Name() || d.MimeType() != other.MimeType() {
		return false
	}
	if d.signedData == nil || other.signedData == nil {
		return d.signedData == other.signedData
	}
	return string(d.signedData.DEREncoded()) == string(other.signedData.DEREncoded())
}
