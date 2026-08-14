// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PdfValidationDataContainer.java
// (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY (sibling chunk of this same phase; assumed shape, matching the accessor
// names PdfDssDictCertificateSource/PdfDssDictCRLSource/PdfDssDictOCSPSource already expose in
// this package - pdf_dss_dict_certificate_source.go, pdf_dss_dict_crl_source.go,
// pdf_dss_dict_ocsp_source.go):
//
//	type PdfDocDssRevision struct { ... }
//	func (r *PdfDocDssRevision) CertificateSource() *PdfDssDictCertificateSource
//	func (r *PdfDocDssRevision) CRLSource() *PdfDssDictCRLSource
//	func (r *PdfDocDssRevision) OCSPSource() *PdfDssDictOCSPSource
//
// java.util.HashMap<String, PdfObjectKey> -> a plain Go map[string]PdfObjectKey: the only
// operation performed against it is containsKey/put/get, order never observed.
package pades

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
)

// PdfValidationDataContainer is a PDF implementation of ValidationDataContainer containing a
// validation data to be incorporated within a PDF document.
type PdfValidationDataContainer struct {
	validation.ValidationDataContainer

	// pdfDssRevisions is a list of PDF DSS revisions.
	pdfDssRevisions []*PdfDocDssRevision

	// knownObjects is the cached map of known object references from a PDF.
	knownObjects map[string]PdfObjectKey
}

// NewPdfValidationDataContainer is the default constructor.
// Port of the constructor PdfValidationDataContainer(Collection<PdfDocDssRevision>).
func NewPdfValidationDataContainer(pdfDssRevisions []*PdfDocDssRevision) *PdfValidationDataContainer {
	return &PdfValidationDataContainer{
		ValidationDataContainer: *validation.NewValidationDataContainer(),
		pdfDssRevisions:         pdfDssRevisions,
	}
}

// KnownObjectsMap builds a map of token identifiers and their unique references within a PDF
// document from a list of extracted PdfRevisions. Port of getKnownObjectsMap().
func (c *PdfValidationDataContainer) KnownObjectsMap() map[string]PdfObjectKey {
	if c.knownObjects == nil {
		c.knownObjects = make(map[string]PdfObjectKey)

		if utils.IsCollectionNotEmpty(c.pdfDssRevisions) {
			for _, dssRevision := range c.pdfDssRevisions {
				certificateSource := dssRevision.CertificateSource()
				for objectKey, certificateToken := range certificateSource.CertificateMap() {
					tokenKey := c.TokenKey(certificateToken)
					if _, known := c.knownObjects[tokenKey]; !known { // keeps the really first occurrence
						c.knownObjects[tokenKey] = objectKey
					}
				}

				crlSource := dssRevision.CRLSource()
				for objectKey, crlBinary := range crlSource.CrlMap() {
					tokenKey := crlBinary.DSSID().AsXmlID()
					if _, known := c.knownObjects[tokenKey]; !known { // keeps the really first occurrence
						c.knownObjects[tokenKey] = objectKey
					}
				}

				ocspSource := dssRevision.OCSPSource()
				for objectKey, ocspResponseBinary := range ocspSource.OcspMap() {
					tokenKey := ocspResponseBinary.DSSID().AsXmlID()
					if _, known := c.knownObjects[tokenKey]; !known { // keeps the really first occurrence
						c.knownObjects[tokenKey] = objectKey
					}
				}
			}
		}
	}

	return c.knownObjects
}

// TokenReference returns a reference corresponding to the given token from the PDF document, nil
// when absent. Port of getTokenReference(Token); matches the shape already established by the
// landed SIGN chunk's sole caller (native_pdf_signature_service.go's nativePDFSignatureServiceTokenRef:
// "if objectKey := validationDataContainer.TokenReference(token); objectKey == nil { ... }"),
// PdfObjectKey being an interface so its nil value is directly comparable.
func (c *PdfValidationDataContainer) TokenReference(token model.Token) PdfObjectKey {
	tokenKey := c.TokenKey(token)
	return c.KnownObjectsMap()[tokenKey]
}

// TokenKey gets a token key (DSS Id or EntityKey Id for a CertificateToken).
// Port of getTokenKey(Token).
func (c *PdfValidationDataContainer) TokenKey(token model.Token) string {
	if certificateToken, ok := token.(*model.CertificateToken); ok {
		return certificateToken.EntityKey().AsXmlID()
	}
	return token.DSSIDAsString()
}
