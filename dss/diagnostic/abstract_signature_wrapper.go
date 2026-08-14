// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/AbstractSignatureWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"math/big"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
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

// GetPdfAnnotationsOverlapConcernedPages returns a list of PDF annotation overlap concerned
// pages. Port of getPdfAnnotationsOverlapConcernedPages().
func (a *AbstractSignatureWrapperBase) GetPdfAnnotationsOverlapConcernedPages() []*big.Int {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetPdfAnnotationsOverlapConcernedPages()
	}
	return nil
}

// GetPdfVisualDifferenceConcernedPages returns a list of PDF visual difference concerned
// pages. Port of getPdfVisualDifferenceConcernedPages().
func (a *AbstractSignatureWrapperBase) GetPdfVisualDifferenceConcernedPages() []*big.Int {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetPdfVisualDifferenceConcernedPages()
	}
	return nil
}

// GetPdfPageDifferenceConcernedPages returns a list of pages missing/added to the final
// revision in a comparison with a signed one. Port of getPdfPageDifferenceConcernedPages().
func (a *AbstractSignatureWrapperBase) GetPdfPageDifferenceConcernedPages() []*big.Int {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetPdfPageDifferenceConcernedPages()
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

// GetPdfExtensionChanges returns a list of changes occurred in a PDF after the current
// signature's revision associated with a signature/document extension. Port of
// getPdfExtensionChanges().
func (a *AbstractSignatureWrapperBase) GetPdfExtensionChanges() []*jaxb.XmlObjectModification {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetPdfExtensionChanges()
	}
	return nil
}

// GetPdfSignatureOrFormFillChanges returns a list of changes occurred in a PDF after the
// current signature's revision associated with a signature creation, form filling. Port of
// getPdfSignatureOrFormFillChanges().
func (a *AbstractSignatureWrapperBase) GetPdfSignatureOrFormFillChanges() []*jaxb.XmlObjectModification {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetPdfSignatureOrFormFillChanges()
	}
	return nil
}

// GetPdfAnnotationChanges returns a list of changes occurred in a PDF after the current
// signature's revision associated with annotation(s) modification. Port of
// getPdfAnnotationChanges().
func (a *AbstractSignatureWrapperBase) GetPdfAnnotationChanges() []*jaxb.XmlObjectModification {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetPdfAnnotationChanges()
	}
	return nil
}

// GetPdfUndefinedChanges returns a list of undefined changes occurred in a PDF after the
// current signature's revision. Port of getPdfUndefinedChanges().
func (a *AbstractSignatureWrapperBase) GetPdfUndefinedChanges() []*jaxb.XmlObjectModification {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetPdfUndefinedChanges()
	}
	return nil
}

// GetModifiedFieldNames returns a list of field names modified after the current signature's
// revision. Port of getModifiedFieldNames().
func (a *AbstractSignatureWrapperBase) GetModifiedFieldNames() []string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetModifiedFieldNames()
	}
	return nil
}

// GetFirstFieldName returns the first signature field name. Port of getFirstFieldName().
func (a *AbstractSignatureWrapperBase) GetFirstFieldName() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetFirstFieldName()
	}
	return ""
}

// GetSignatureFieldNames returns a list of signature field names, where the signature is
// referenced from. Port of getSignatureFieldNames().
func (a *AbstractSignatureWrapperBase) GetSignatureFieldNames() []string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetSignatureFieldNames()
	}
	return nil
}

// GetSignerName returns the signer's name. Port of getSignerName().
func (a *AbstractSignatureWrapperBase) GetSignerName() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetSignerName()
	}
	return ""
}

// GetSignatureDictionaryType returns the PDF signature dictionary /Type value. Port of
// getSignatureDictionaryType().
func (a *AbstractSignatureWrapperBase) GetSignatureDictionaryType() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetSignatureDictionaryType()
	}
	return ""
}

// GetFilter returns the PDF signature dictionary /Filter value. Port of getFilter().
func (a *AbstractSignatureWrapperBase) GetFilter() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetFilter()
	}
	return ""
}

// GetSubFilter returns the PDF signature dictionary /SubFilter value. Port of getSubFilter().
func (a *AbstractSignatureWrapperBase) GetSubFilter() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetSubFilter()
	}
	return ""
}

// GetContactInfo returns the PDF signature dictionary /ContactInfo value. Port of
// getContactInfo().
func (a *AbstractSignatureWrapperBase) GetContactInfo() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetContactInfo()
	}
	return ""
}

// GetLocation returns the PDF signature dictionary /Location value. Port of getLocation().
func (a *AbstractSignatureWrapperBase) GetLocation() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetLocation()
	}
	return ""
}

// GetReason returns the PDF signature dictionary /Reason value. Port of getReason().
func (a *AbstractSignatureWrapperBase) GetReason() string {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetReason()
	}
	return ""
}

// GetSignatureByteRange returns the PDF signature dictionary /ByteRange value. Port of
// getSignatureByteRange().
func (a *AbstractSignatureWrapperBase) GetSignatureByteRange() []*big.Int {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetSignatureByteRange()
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

// GetDocMDPPermissions returns a CertificationPermission value of a /DocMDP dictionary, when
// present. Port of getDocMDPPermissions().
func (a *AbstractSignatureWrapperBase) GetDocMDPPermissions() enumerations.CertificationPermission {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetDocMDPPermissions()
	}
	return ""
}

// GetFieldMDP returns a /FieldMDP dictionary content, when present. Port of getFieldMDP().
func (a *AbstractSignatureWrapperBase) GetFieldMDP() *jaxb.XmlPDFLockDictionary {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetFieldMDP()
	}
	return nil
}

// GetSigFieldLock returns a /SigFieldLock dictionary, when present. Port of getSigFieldLock().
func (a *AbstractSignatureWrapperBase) GetSigFieldLock() *jaxb.XmlPDFLockDictionary {
	pdfRevision := a.signatureWrapperOverrides().PDFRevision()
	if pdfRevision != nil {
		return pdfRevision.GetSigFieldLock()
	}
	return nil
}
