// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfSignatureCache.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is the one Java package of dss-pades that landed in no s5b manifest
// (see pdf_object.go's header). Shape confirmed against the already-landed call sites
// (native_pdf_signature_service.go, pades_profile_parameters.go, pades_timestamp_parameters.go,
// pades_signature_parameters.go): MessageDigest/SetMessageDigest carry a model.DSSMessageDigest
// by value (its zero value stands in for Java's null - DSSMessageDigest.Value() returns nil for
// the zero value, matching native_pdf_signature_service.go's ".MessageDigest().Value() == nil"
// check), ToBeSignedDocument/SetToBeSignedDocument carry the model.DSSDocument interface, whose
// nil is Java's null.
package pades

import "github.com/ryftcore/dss-go/dss/model"

// PdfSignatureCache is used as a DTO containing cached data to accelerate the signature creation
// process. Port of the PdfSignatureCache class.
type PdfSignatureCache struct {
	// messageDigest is the cached digest value of the covered ByteRange.
	messageDigest model.DSSMessageDigest

	// toBeSignedDocument is a pre-generated PDF document, used for digest computation,
	// preserving a /Contents space for CMS Signed Data inclusion.
	toBeSignedDocument model.DSSDocument
}

// NewPdfSignatureCache instantiates an object with null-equivalent (zero) values.
// Port of the default constructor.
func NewPdfSignatureCache() *PdfSignatureCache {
	return &PdfSignatureCache{}
}

// MessageDigest gets message-digest computed in the prepared PDF revision ByteRange.
// Port of #getMessageDigest.
func (c *PdfSignatureCache) MessageDigest() model.DSSMessageDigest {
	return c.messageDigest
}

// SetMessageDigest sets message-digest of the ByteRange. Port of #setMessageDigest.
func (c *PdfSignatureCache) SetMessageDigest(messageDigest model.DSSMessageDigest) {
	c.messageDigest = messageDigest
}

// ToBeSignedDocument gets the ToBeSigned document. Port of #getToBeSignedDocument.
func (c *PdfSignatureCache) ToBeSignedDocument() model.DSSDocument {
	return c.toBeSignedDocument
}

// SetToBeSignedDocument sets the ToBeSigned document. Port of #setToBeSignedDocument.
func (c *PdfSignatureCache) SetToBeSignedDocument(toBeSignedDocument model.DSSDocument) {
	c.toBeSignedDocument = toBeSignedDocument
}
