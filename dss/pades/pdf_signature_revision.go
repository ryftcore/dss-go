// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfSignatureRevision.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is implemented by this file and others (see pdf_object.go's header).
// Its shape (embeds/satisfies PdfCMSRevision, CompositeDssDictionary(), DssDictionary()) is used
// by pades_certificate_source.go and pades_signature.go, and its constructor's exact argument
// order matches native_pdf_signature_service.go's call site.
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
