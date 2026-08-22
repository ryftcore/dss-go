// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/AbstractSignatureWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"math/big"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// AbstractSignatureWrapperOverrides declares the operations Java's abstract class
// AbstractSignatureWrapper leaves abstract (getFilename(), getPDFRevision()), on top of the
// AbstractTokenProxy overrides it also requires. Concrete wrappers (SignatureWrapper,
// TimestampWrapper, in DIAGWRAP_B) implement the full interface and register themselves with
// InitSignatureWrapper.
type AbstractSignatureWrapperOverrides interface {
	AbstractTokenProxyOverrides

	// Filename gets name of the signature or timestamp file, when applicable. Port of the
	// abstract getFilename().
	Filename() string
	// PDFRevision returns a PAdES-specific PDF Revision info. NOTE: applicable only for
	// PAdES. Port of the abstract getPDFRevision().
	PDFRevision() *PDFRevisionWrapper
}

// AbstractSignatureWrapperBase contains common code for signature tokens (signature or
// timestamps). Port of the Java abstract class AbstractSignatureWrapper. Concrete wrappers
// embed it and register themselves with InitSignatureWrapper.
type AbstractSignatureWrapperBase struct {
	AbstractTokenProxyBase

	overrides AbstractSignatureWrapperOverrides
}

// InitSignatureWrapper registers the concrete wrapper with its base, and with the embedded
// AbstractTokenProxyBase, so that both bases can dispatch to the operations Java would reach
// through virtual dispatch. It must be called exactly once, by the concrete wrapper's
// constructor, before any other method.
func (a *AbstractSignatureWrapperBase) InitSignatureWrapper(overrides AbstractSignatureWrapperOverrides) {
	a.overrides = overrides
	a.InitTokenProxy(overrides)
}

func (a *AbstractSignatureWrapperBase) signatureWrapperOverrides() AbstractSignatureWrapperOverrides {
	if a.overrides == nil {
		panic("AbstractSignatureWrapper was not initialised: the concrete wrapper must call InitSignatureWrapper in its constructor")
	}
	return a.overrides
}

// ArePdfModificationsDetected checks if any visual modifications detected in the PDF. Port of
// arePdfModificationsDetected().
func (a *AbstractSignatureWrapperBase) ArePdfModificationsDetected() bool {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.ArePdfModificationsDetected()
	}
	return false
}

// PdfAnnotationsOverlapConcernedPages returns a list of PDF annotation overlap concerned
// pages. Port of getPdfAnnotationsOverlapConcernedPages().
func (a *AbstractSignatureWrapperBase) PdfAnnotationsOverlapConcernedPages() []*big.Int {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.PdfAnnotationsOverlapConcernedPages()
	}
	return nil
}

// PdfVisualDifferenceConcernedPages returns a list of PDF visual difference concerned
// pages. Port of getPdfVisualDifferenceConcernedPages().
func (a *AbstractSignatureWrapperBase) PdfVisualDifferenceConcernedPages() []*big.Int {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.PdfVisualDifferenceConcernedPages()
	}
	return nil
}

// PdfPageDifferenceConcernedPages returns a list of pages missing/added to the final
// revision in a comparison with a signed one. Port of getPdfPageDifferenceConcernedPages().
func (a *AbstractSignatureWrapperBase) PdfPageDifferenceConcernedPages() []*big.Int {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.PdfPageDifferenceConcernedPages()
	}
	return nil
}

// ArePdfObjectModificationsDetected checks whether object modifications are present after the
// current PDF revisions. Port of arePdfObjectModificationsDetected().
func (a *AbstractSignatureWrapperBase) ArePdfObjectModificationsDetected() bool {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.ArePdfObjectModificationsDetected()
	}
	return false
}

// PdfExtensionChanges returns a list of changes occurred in a PDF after the current
// signature's revision associated with a signature/document extension. Port of
// getPdfExtensionChanges().
func (a *AbstractSignatureWrapperBase) PdfExtensionChanges() []*jaxb.XmlObjectModification {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.PdfExtensionChanges()
	}
	return nil
}

// PdfSignatureOrFormFillChanges returns a list of changes occurred in a PDF after the
// current signature's revision associated with a signature creation, form filling. Port of
// getPdfSignatureOrFormFillChanges().
func (a *AbstractSignatureWrapperBase) PdfSignatureOrFormFillChanges() []*jaxb.XmlObjectModification {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.PdfSignatureOrFormFillChanges()
	}
	return nil
}

// PdfAnnotationChanges returns a list of changes occurred in a PDF after the current
// signature's revision associated with annotation(s) modification. Port of
// getPdfAnnotationChanges().
func (a *AbstractSignatureWrapperBase) PdfAnnotationChanges() []*jaxb.XmlObjectModification {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.PdfAnnotationChanges()
	}
	return nil
}

// PdfUndefinedChanges returns a list of undefined changes occurred in a PDF after the
// current signature's revision. Port of getPdfUndefinedChanges().
func (a *AbstractSignatureWrapperBase) PdfUndefinedChanges() []*jaxb.XmlObjectModification {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.PdfUndefinedChanges()
	}
	return nil
}

// ModifiedFieldNames returns a list of field names modified after the current signature's
// revision. Port of getModifiedFieldNames().
func (a *AbstractSignatureWrapperBase) ModifiedFieldNames() []string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.ModifiedFieldNames()
	}
	return nil
}

// FirstFieldName returns the first signature field name. Port of getFirstFieldName().
func (a *AbstractSignatureWrapperBase) FirstFieldName() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.FirstFieldName()
	}
	return ""
}

// SignatureFieldNames returns a list of signature field names, where the signature is
// referenced from. Port of getSignatureFieldNames().
func (a *AbstractSignatureWrapperBase) SignatureFieldNames() []string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.SignatureFieldNames()
	}
	return nil
}

// SignerName returns the signer's name. Port of getSignerName().
func (a *AbstractSignatureWrapperBase) SignerName() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.SignerName()
	}
	return ""
}

// SignatureDictionaryType returns the PDF signature dictionary /Type value. Port of
// getSignatureDictionaryType().
func (a *AbstractSignatureWrapperBase) SignatureDictionaryType() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.SignatureDictionaryType()
	}
	return ""
}

// Filter returns the PDF signature dictionary /Filter value. Port of getFilter().
func (a *AbstractSignatureWrapperBase) Filter() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.Filter()
	}
	return ""
}

// SubFilter returns the PDF signature dictionary /SubFilter value. Port of getSubFilter().
func (a *AbstractSignatureWrapperBase) SubFilter() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.SubFilter()
	}
	return ""
}

// ContactInfo returns the PDF signature dictionary /ContactInfo value. Port of
// getContactInfo().
func (a *AbstractSignatureWrapperBase) ContactInfo() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.ContactInfo()
	}
	return ""
}

// Location returns the PDF signature dictionary /Location value. Port of getLocation().
func (a *AbstractSignatureWrapperBase) Location() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.Location()
	}
	return ""
}

// Reason returns the PDF signature dictionary /Reason value. Port of getReason().
func (a *AbstractSignatureWrapperBase) Reason() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.Reason()
	}
	return ""
}

// SignatureByteRange returns the PDF signature dictionary /ByteRange value. Port of
// getSignatureByteRange().
func (a *AbstractSignatureWrapperBase) SignatureByteRange() []*big.Int {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.SignatureByteRange()
	}
	return nil
}

// IsSignatureByteRangeValid returns whether the PDF signature dictionary /ByteRange is found
// and valid. Port of isSignatureByteRangeValid().
func (a *AbstractSignatureWrapperBase) IsSignatureByteRangeValid() bool {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.IsSignatureByteRangeValid()
	}
	return false
}

// IsPdfSignatureDictionaryConsistent returns whether the PDF signature dictionary is
// consistent across PDF revisions. Port of isPdfSignatureDictionaryConsistent().
func (a *AbstractSignatureWrapperBase) IsPdfSignatureDictionaryConsistent() bool {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.IsPdfSignatureDictionaryConsistent()
	}
	return false
}

// DocMDPPermissions returns a CertificationPermission value of a /DocMDP dictionary, when
// present. Port of getDocMDPPermissions().
func (a *AbstractSignatureWrapperBase) DocMDPPermissions() enumerations.CertificationPermission {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.DocMDPPermissions()
	}
	return ""
}

// FieldMDP returns a /FieldMDP dictionary content, when present. Port of getFieldMDP().
func (a *AbstractSignatureWrapperBase) FieldMDP() *jaxb.XmlPDFLockDictionary {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.FieldMDP()
	}
	return nil
}

// SigFieldLock returns a /SigFieldLock dictionary, when present. Port of getSigFieldLock().
func (a *AbstractSignatureWrapperBase) SigFieldLock() *jaxb.XmlPDFLockDictionary {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.SigFieldLock()
	}
	return nil
}
