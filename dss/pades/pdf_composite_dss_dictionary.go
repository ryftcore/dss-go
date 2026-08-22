// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/dss/PdfCompositeDssDictionary.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pades.validation.dss flattens into the Go package pades,
// so the type keeps its Java name unqualified.
package pades

// PdfCompositeDssDictionary represents a merged result of all /DSS dictionaries' content
// extracted from a PDF document.
//
// Java's Serializable has no Go counterpart.
type PdfCompositeDssDictionary struct {
	// certificateSource represents a merged result of certificate sources extracted from the
	// PDF document.
	certificateSource *PdfCompositeDssDictCertificateSource

	// crlSource represents a merged result of CRL sources extracted from the PDF document.
	crlSource *PdfCompositeDssDictCRLSource

	// ocspSource represents a merged result of OCSP sources extracted from the PDF document.
	ocspSource *PdfCompositeDssDictOCSPSource
}

// NewPdfCompositeDssDictionary builds an empty composite dictionary, with its three empty
// composite sources. Port of the default constructor.
func NewPdfCompositeDssDictionary() *PdfCompositeDssDictionary {
	return &PdfCompositeDssDictionary{
		certificateSource: NewPdfCompositeDssDictCertificateSource(),
		crlSource:         NewPdfCompositeDssDictCRLSource(),
		ocspSource:        NewPdfCompositeDssDictOCSPSource(),
	}
}

// CertificateSource gets the composite certificate source. Port of getCertificateSource().
func (d *PdfCompositeDssDictionary) CertificateSource() *PdfCompositeDssDictCertificateSource {
	return d.certificateSource
}

// CrlSource gets the composite CRL source. Port of getCrlSource().
func (d *PdfCompositeDssDictionary) CrlSource() *PdfCompositeDssDictCRLSource {
	return d.crlSource
}

// OcspSource gets the composite OCSP source. Port of getOcspSource().
func (d *PdfCompositeDssDictionary) OcspSource() *PdfCompositeDssDictOCSPSource {
	return d.ocspSource
}

// PopulateFromDssDictionary populates the certificate and revocation sources with the data
// extracted from a /DSS revision. Port of populateFromDssDictionary(PdfDssDict).
func (d *PdfCompositeDssDictionary) PopulateFromDssDictionary(dssDict PdfDssDict) {
	if dssDict != nil {
		d.certificateSource.PopulateFromDssDictionary(dssDict)
		d.crlSource.PopulateFromDssDictionary(dssDict)
		d.ocspSource.PopulateFromDssDictionary(dssDict)
	}
}
