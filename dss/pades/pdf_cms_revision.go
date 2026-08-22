// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfCMSRevision.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is implemented by this file and others (see pdf_object.go's header).
// Its interface shape (embedding PdfRevision, plus ByteRange(), AreAllOriginalBytesCovered() and
// SetModificationDetection()) is exercised by every call site across the package; see
// pdf_revision_scope_finder.go.
//
// STRUCTURE DEVIATION: Java's PdfCMSRevision is an abstract class carrying the state and
// behaviour shared by its two subclasses, PdfSignatureRevision and PdfDocTimestampRevision. Go
// has no implementation inheritance, so that shared state+behaviour becomes pdfCMSRevisionBase,
// a value type embedded by both subclasses' Go ports (see pdf_signature_revision.go and
// pdf_doc_timestamp_revision.go) - the same "specialised type embeds/re-implements what it
// needs" convention already used throughout this port (e.g. cms_for_pades_builder_helper.go).
// PdfDocTimestampRevision overrides SigningDate() in Java; nothing inside pdfCMSRevisionBase's
// own methods calls SigningDate() on itself, so ordinary Go method shadowing (the embedder
// declaring its own SigningDate()) reproduces the override with no virtual-dispatch gap.
package pades

import (
	"time"

	"github.com/ryftcore/dss-go/dss/cms"
	"github.com/ryftcore/dss-go/dss/model"
)

// PdfCMSRevision defines a PDF revision containing CMS data (signature/timestamp).
// Port of the abstract PdfCMSRevision class.
type PdfCMSRevision interface {
	PdfRevision

	// SignedData gets the current signature revision. Port of #getSignedData.
	SignedData() model.DSSDocument

	// PreviousRevision gets the PDF revision preceding the current signature revision.
	// Port of #getPreviousRevision.
	PreviousRevision() model.DSSDocument

	// ByteRange gets the signed byte range. Port of #getByteRange.
	ByteRange() *ByteRange

	// SigningDate gets the claimed signing time. Port of #getSigningDate.
	SigningDate() time.Time

	// AreAllOriginalBytesCovered gets whether all of the PDF's content is signed.
	// Port of #areAllOriginalBytesCovered.
	AreAllOriginalBytesCovered() bool

	// CMS gets the CMSSignedData. Port of #getCMS.
	CMS() *cms.CMS

	// SetModificationDetection sets the PdfModificationDetection result.
	// Port of #setModificationDetection.
	SetModificationDetection(modificationDetection *PdfModificationDetection)
}

// pdfCMSRevisionBase carries the state and behaviour Java's abstract PdfCMSRevision class
// implements once for both its subclasses; see this file's header.
type pdfCMSRevisionBase struct {
	// signatureDictionary is the PDF Signature Dictionary.
	signatureDictionary *PdfSignatureDictionary

	// signatureFields is a list of signed fields by the corresponding signature.
	signatureFields []*PdfSignatureField

	// signedContent is the signed data binaries.
	signedContent model.DSSDocument

	// previousRevision is the original signed revision content.
	previousRevision model.DSSDocument

	// coverAllOriginalBytes defines if the revision covers all document bytes.
	coverAllOriginalBytes bool

	// modificationDetection detects modification in the PDF content.
	modificationDetection *PdfModificationDetection
}

// newPdfCMSRevisionBase builds the shared PdfCMSRevision state. Port of the protected
// PdfCMSRevision constructor.
//
// DEVIATION: Java also requires signatureFields to be non-null (Objects.requireNonNull); a nil
// Go slice is the idiomatic empty list, so that particular check is not reproduced here -
// signatureDictionary/signedContent/previousRevision (reference values whose Java null truly
// signals a caller bug) keep theirs.
func newPdfCMSRevisionBase(signatureDictionary *PdfSignatureDictionary, signatureFields []*PdfSignatureField,
	signedContent, previousRevision model.DSSDocument, coverAllOriginalBytes bool) pdfCMSRevisionBase {
	if signatureDictionary == nil {
		panic("The signature dictionary cannot be null!")
	}
	if signedContent == nil {
		panic("The signed content cannot be null!")
	}
	if previousRevision == nil {
		panic("The previous revision cannot be null!")
	}
	return pdfCMSRevisionBase{
		signatureDictionary:   signatureDictionary,
		signatureFields:       signatureFields,
		signedContent:         signedContent,
		previousRevision:      previousRevision,
		coverAllOriginalBytes: coverAllOriginalBytes,
	}
}

// SignedData gets the current signature revision. Port of #getSignedData.
func (r *pdfCMSRevisionBase) SignedData() model.DSSDocument { return r.signedContent }

// PreviousRevision gets the PDF revision preceding the current signature revision.
// Port of #getPreviousRevision.
func (r *pdfCMSRevisionBase) PreviousRevision() model.DSSDocument { return r.previousRevision }

// PdfSigDictInfo returns the PDF Signature Dictionary. Port of #getPdfSigDictInfo.
func (r *pdfCMSRevisionBase) PdfSigDictInfo() *PdfSignatureDictionary { return r.signatureDictionary }

// ByteRange gets the signed byte range. Port of #getByteRange.
func (r *pdfCMSRevisionBase) ByteRange() *ByteRange { return r.signatureDictionary.ByteRange() }

// SigningDate gets the claimed signing time. Port of #getSigningDate.
func (r *pdfCMSRevisionBase) SigningDate() time.Time { return r.signatureDictionary.SigningDate() }

// AreAllOriginalBytesCovered gets whether all of the PDF's content is signed.
// Port of #areAllOriginalBytesCovered.
func (r *pdfCMSRevisionBase) AreAllOriginalBytesCovered() bool { return r.coverAllOriginalBytes }

// Fields returns the list of signature fields that refer the current object. Port of #getFields.
func (r *pdfCMSRevisionBase) Fields() []*PdfSignatureField { return r.signatureFields }

// CMS gets the CMSSignedData. Port of #getCMS.
func (r *pdfCMSRevisionBase) CMS() *cms.CMS { return r.signatureDictionary.CMS() }

// ModificationDetection returns information about changes made in the document.
// Port of #getModificationDetection.
func (r *pdfCMSRevisionBase) ModificationDetection() *PdfModificationDetection {
	return r.modificationDetection
}

// SetModificationDetection sets the PdfModificationDetection result. Port of #setModificationDetection.
func (r *pdfCMSRevisionBase) SetModificationDetection(modificationDetection *PdfModificationDetection) {
	r.modificationDetection = modificationDetection
}
