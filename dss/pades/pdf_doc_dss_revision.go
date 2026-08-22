// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfDocDssRevision.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is implemented by this file and others (see pdf_object.go's header).
// Its shape (NewPdfDocDssRevision(*PdfCompositeDssDictionary, PdfDssDict)) is used by
// pdf_validation_data_container.go and native_pdf_signature_service.go's call sites.
package pades

// PdfDocDssRevision represents an LT-level PDF revision containing a DSS dictionary.
// Port of the PdfDocDssRevision class, implementing PdfRevision directly (not PdfCMSRevision -
// a DSS revision carries no signature dictionary, CMS or byte range).
type PdfDocDssRevision struct {
	// compositeDssDictionary is the composite DSS dictionary combined from all /DSS revisions' content.
	compositeDssDictionary *PdfCompositeDssDictionary

	// dssDictionary is the DSS dictionary from the revision.
	dssDictionary PdfDssDict

	// certificateSource is the cached certificate source.
	certificateSource *PdfDssDictCertificateSource

	// crlSource is the cached CRL source.
	crlSource *PdfDssDictCRLSource

	// ocspSource is the cached OCSP source.
	ocspSource *PdfDssDictOCSPSource
}

// NewPdfDocDssRevision is the default constructor. Port of the
// PdfDocDssRevision(PdfCompositeDssDictionary, PdfDssDict) constructor.
func NewPdfDocDssRevision(compositeDssDictionary *PdfCompositeDssDictionary, dssDictionary PdfDssDict) *PdfDocDssRevision {
	if compositeDssDictionary == nil {
		panic("Composite DSS dictionary cannot be null!")
	}
	if dssDictionary == nil {
		panic("The dssDictionary cannot be null!")
	}
	return &PdfDocDssRevision{compositeDssDictionary: compositeDssDictionary, dssDictionary: dssDictionary}
}

// DssDictionary returns the DSS dictionary. Port of #getDssDictionary.
func (r *PdfDocDssRevision) DssDictionary() PdfDssDict { return r.dssDictionary }

// PdfSigDictInfo is not applicable for a DSS revision; it returns nil. Port of #getPdfSigDictInfo.
func (r *PdfDocDssRevision) PdfSigDictInfo() *PdfSignatureDictionary { return nil }

// Fields is not applicable for a DSS revision; it returns nil (Java's Collections.emptyList()).
// Port of #getFields.
func (r *PdfDocDssRevision) Fields() []*PdfSignatureField { return nil }

// ModificationDetection is not applicable for a DSS revision; it returns nil.
// Port of #getModificationDetection.
func (r *PdfDocDssRevision) ModificationDetection() *PdfModificationDetection { return nil }

// CertificateSource returns the corresponding certificate source. Port of #getCertificateSource.
func (r *PdfDocDssRevision) CertificateSource() *PdfDssDictCertificateSource {
	if r.certificateSource == nil {
		r.certificateSource = NewPdfDssDictCertificateSource(r.compositeDssDictionary.CertificateSource(), r.dssDictionary)
	}
	return r.certificateSource
}

// CRLSource returns the corresponding CRL source. Port of #getCRLSource.
func (r *PdfDocDssRevision) CRLSource() *PdfDssDictCRLSource {
	if r.crlSource == nil {
		r.crlSource = NewPdfDssDictCRLSource(r.compositeDssDictionary.CrlSource(), r.dssDictionary)
	}
	return r.crlSource
}

// OCSPSource returns the corresponding OCSP source. Port of #getOCSPSource.
func (r *PdfDocDssRevision) OCSPSource() *PdfDssDictOCSPSource {
	if r.ocspSource == nil {
		r.ocspSource = NewPdfDssDictOCSPSource(r.compositeDssDictionary.OcspSource(), r.dssDictionary)
	}
	return r.ocspSource
}

// Compile-time assertion standing in for Java's "implements PdfRevision".
var _ PdfRevision = (*PdfDocDssRevision)(nil)
