// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfSignatureRevision.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is the one Java package of dss-pades that landed in no s5b manifest
// (see pdf_object.go's header). Shape (embeds/satisfies PdfCMSRevision, CompositeDssDictionary(),
// DssDictionary()) confirmed against pades_certificate_source.go's / pades_signature.go's already-
// landed forward-dependency headers, and the constructor's exact argument order against
// native_pdf_signature_service.go's already-landed call site.
package pades

import "github.com/ryftcore/dss-go/dss/model"

// PdfSignatureRevision represents a PDF revision for an electronic signature.
// Port of the PdfSignatureRevision class, extending PdfCMSRevision.
type PdfSignatureRevision struct {
	pdfCMSRevisionBase

	// compositeDssDictionary is the composite DSS dictionary combined from all /DSS revisions' content.
	compositeDssDictionary *PdfCompositeDssDictionary

	// dssDictionary is the corresponding DSS dictionary.
	dssDictionary PdfDssDict
}

// NewPdfSignatureRevision is the default constructor. Port of the
// PdfSignatureRevision(PdfSignatureDictionary, PdfCompositeDssDictionary, PdfDssDict,
// List<PdfSignatureField>, DSSDocument, DSSDocument, boolean) constructor.
func NewPdfSignatureRevision(signatureDictionary *PdfSignatureDictionary, compositeDssDictionary *PdfCompositeDssDictionary,
	dssDictionary PdfDssDict, signatureFields []*PdfSignatureField, signedContent, previousRevision model.DSSDocument,
	coverCompleteRevision bool) *PdfSignatureRevision {
	return &PdfSignatureRevision{
		pdfCMSRevisionBase:     newPdfCMSRevisionBase(signatureDictionary, signatureFields, signedContent, previousRevision, coverCompleteRevision),
		compositeDssDictionary: compositeDssDictionary,
		dssDictionary:          dssDictionary,
	}
}

// CompositeDssDictionary gets the composite DSS dictionary. Port of #getCompositeDssDictionary.
func (r *PdfSignatureRevision) CompositeDssDictionary() *PdfCompositeDssDictionary {
	return r.compositeDssDictionary
}

// DssDictionary gets the DSS dictionary. Port of #getDssDictionary.
func (r *PdfSignatureRevision) DssDictionary() PdfDssDict { return r.dssDictionary }

// Compile-time assertion standing in for Java's "extends PdfCMSRevision".
var _ PdfCMSRevision = (*PdfSignatureRevision)(nil)
