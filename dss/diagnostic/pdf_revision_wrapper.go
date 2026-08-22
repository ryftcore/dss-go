// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/PDFRevisionWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"math/big"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// PDFRevisionWrapper contains user-friendly methods to extract information from a
// jaxb.XmlPDFRevision.
type PDFRevisionWrapper struct {
	// pdfRevision is the wrapped XmlPDFRevision.
	pdfRevision *jaxb.XmlPDFRevision
}

// NewPDFRevisionWrapper is the default constructor. Port of PDFRevisionWrapper(XmlPDFRevision);
// panics per Objects.requireNonNull(pdfRevision, "XmlPDFRevision cannot be null!").
func NewPDFRevisionWrapper(pdfRevision *jaxb.XmlPDFRevision) *PDFRevisionWrapper {
	if pdfRevision == nil {
		panic("XmlPDFRevision cannot be null!")
	}
	return &PDFRevisionWrapper{pdfRevision: pdfRevision}
}

// ArePdfModificationsDetected checks if any visual modifications detected in the PDF. Port of
// arePdfModificationsDetected().
func (w *PDFRevisionWrapper) ArePdfModificationsDetected() bool {
	modificationDetection := w.pdfRevision.ModificationDetection
	if modificationDetection != nil {
		return len(modificationDetection.AnnotationOverlap) != 0 ||
			len(modificationDetection.VisualDifference) != 0 ||
			len(modificationDetection.PageDifference) != 0
	}
	return false
}

// PdfAnnotationsOverlapConcernedPages returns a list of PDF annotation overlap concerned
// pages. Port of getPdfAnnotationsOverlapConcernedPages().
func (w *PDFRevisionWrapper) PdfAnnotationsOverlapConcernedPages() []*big.Int {
	modificationDetection := w.pdfRevision.ModificationDetection
	if modificationDetection != nil {
		return concernedPages(modificationDetection.AnnotationOverlap)
	}
	return nil
}

// PdfVisualDifferenceConcernedPages returns a list of PDF visual difference concerned
// pages. Port of getPdfVisualDifferenceConcernedPages().
func (w *PDFRevisionWrapper) PdfVisualDifferenceConcernedPages() []*big.Int {
	modificationDetection := w.pdfRevision.ModificationDetection
	if modificationDetection != nil {
		return concernedPages(modificationDetection.VisualDifference)
	}
	return nil
}

// PdfPageDifferenceConcernedPages returns a list of pages missing/added to the final
// revision in a comparison with a signed one. Port of getPdfPageDifferenceConcernedPages().
func (w *PDFRevisionWrapper) PdfPageDifferenceConcernedPages() []*big.Int {
	modificationDetection := w.pdfRevision.ModificationDetection
	if modificationDetection != nil {
		return concernedPages(modificationDetection.PageDifference)
	}
	return nil
}

// ArePdfObjectModificationsDetected checks whether object modifications are present after the
// current PDF revisions. Port of arePdfObjectModificationsDetected().
func (w *PDFRevisionWrapper) ArePdfObjectModificationsDetected() bool {
	return w.pdfObjectModifications() != nil
}

// PdfExtensionChanges returns a list of changes occurred in a PDF after the current
// signature's revision associated with a signature/document extension. Port of
// getPdfExtensionChanges().
func (w *PDFRevisionWrapper) PdfExtensionChanges() []*jaxb.XmlObjectModification {
	pdfObjectModifications := w.pdfObjectModifications()
	if pdfObjectModifications != nil {
		return pdfObjectModifications.ExtensionChange
	}
	return nil
}

// PdfSignatureOrFormFillChanges returns a list of changes occurred in a PDF after the
// current signature's revision associated with a signature creation, form filling. Port of
// getPdfSignatureOrFormFillChanges().
func (w *PDFRevisionWrapper) PdfSignatureOrFormFillChanges() []*jaxb.XmlObjectModification {
	pdfObjectModifications := w.pdfObjectModifications()
	if pdfObjectModifications != nil {
		return pdfObjectModifications.SignatureOrFormFill
	}
	return nil
}

// PdfAnnotationChanges returns a list of changes occurred in a PDF after the current
// signature's revision associated with annotation(s) modification. Port of
// getPdfAnnotationChanges().
func (w *PDFRevisionWrapper) PdfAnnotationChanges() []*jaxb.XmlObjectModification {
	pdfObjectModifications := w.pdfObjectModifications()
	if pdfObjectModifications != nil {
		return pdfObjectModifications.AnnotationChange
	}
	return nil
}

// PdfUndefinedChanges returns a list of undefined changes occurred in a PDF after the
// current signature's revision. Port of getPdfUndefinedChanges().
func (w *PDFRevisionWrapper) PdfUndefinedChanges() []*jaxb.XmlObjectModification {
	pdfObjectModifications := w.pdfObjectModifications()
	if pdfObjectModifications != nil {
		return pdfObjectModifications.Undefined
	}
	return nil
}

// ModifiedFieldNames returns a list of field names modified after the current signature's
// revision. Port of getModifiedFieldNames().
func (w *PDFRevisionWrapper) ModifiedFieldNames() []string {
	var names []string
	pdfObjectModifications := w.pdfObjectModifications()
	if pdfObjectModifications != nil {
		names = append(names, modifiedFieldNames(pdfObjectModifications.ExtensionChange)...)
		names = append(names, modifiedFieldNames(pdfObjectModifications.SignatureOrFormFill)...)
		names = append(names, modifiedFieldNames(pdfObjectModifications.AnnotationChange)...)
		names = append(names, modifiedFieldNames(pdfObjectModifications.Undefined)...)
	}
	return names
}

// FirstFieldName returns the first signature field name. Port of getFirstFieldName().
func (w *PDFRevisionWrapper) FirstFieldName() string {
	fields := w.pdfRevision.SignatureField
	if len(fields) != 0 && fields[0].Name != nil {
		return *fields[0].Name
	}
	return ""
}

// SignatureFieldNames returns a list of signature field names, where the signature is
// referenced from. Port of getSignatureFieldNames().
func (w *PDFRevisionWrapper) SignatureFieldNames() []string {
	var names []string
	fields := w.pdfRevision.SignatureField
	if len(fields) != 0 {
		for _, signatureField := range fields {
			if signatureField.Name != nil {
				names = append(names, *signatureField.Name)
			}
		}
	}
	return names
}

// SignerName returns the signer's name. Port of getSignerName().
func (w *PDFRevisionWrapper) SignerName() string {
	if w.pdfRevision.PDFSignatureDictionary.SignerName != nil {
		return *w.pdfRevision.PDFSignatureDictionary.SignerName
	}
	return ""
}

// SignatureDictionaryType returns the PDF signature dictionary /Type value. Port of
// getSignatureDictionaryType().
func (w *PDFRevisionWrapper) SignatureDictionaryType() string {
	if w.pdfRevision.PDFSignatureDictionary.Type != nil {
		return *w.pdfRevision.PDFSignatureDictionary.Type
	}
	return ""
}

// Filter returns the PDF signature dictionary /Filter value. Port of getFilter().
func (w *PDFRevisionWrapper) Filter() string {
	if w.pdfRevision.PDFSignatureDictionary.Filter != nil {
		return *w.pdfRevision.PDFSignatureDictionary.Filter
	}
	return ""
}

// SubFilter returns the PDF signature dictionary /SubFilter value. Port of getSubFilter().
func (w *PDFRevisionWrapper) SubFilter() string {
	if w.pdfRevision.PDFSignatureDictionary.SubFilter != nil {
		return *w.pdfRevision.PDFSignatureDictionary.SubFilter
	}
	return ""
}

// ContactInfo returns the PDF signature dictionary /ContactInfo value. Port of
// getContactInfo().
func (w *PDFRevisionWrapper) ContactInfo() string {
	if w.pdfRevision.PDFSignatureDictionary.ContactInfo != nil {
		return *w.pdfRevision.PDFSignatureDictionary.ContactInfo
	}
	return ""
}

// Location returns the PDF signature dictionary /Location value. Port of getLocation().
func (w *PDFRevisionWrapper) Location() string {
	if w.pdfRevision.PDFSignatureDictionary.Location != nil {
		return *w.pdfRevision.PDFSignatureDictionary.Location
	}
	return ""
}

// Reason returns the PDF signature dictionary /Reason value. Port of getReason().
func (w *PDFRevisionWrapper) Reason() string {
	if w.pdfRevision.PDFSignatureDictionary.Reason != nil {
		return *w.pdfRevision.PDFSignatureDictionary.Reason
	}
	return ""
}

// SignatureByteRange returns the PDF signature dictionary /ByteRange value. Port of
// getSignatureByteRange().
func (w *PDFRevisionWrapper) SignatureByteRange() []*big.Int {
	byteRange := w.xmlByteRange()
	if byteRange != nil {
		return byteRange.Value
	}
	return nil
}

// IsSignatureByteRangeValid returns whether the PDF signature dictionary /ByteRange is found
// and valid. Port of isSignatureByteRangeValid().
func (w *PDFRevisionWrapper) IsSignatureByteRangeValid() bool {
	byteRange := w.xmlByteRange()
	if byteRange != nil {
		return byteRange.Valid
	}
	return false
}

func (w *PDFRevisionWrapper) xmlByteRange() *jaxb.XmlByteRange {
	return w.pdfRevision.PDFSignatureDictionary.SignatureByteRange
}

// IsPdfSignatureDictionaryConsistent returns whether the PDF signature dictionary is
// consistent across PDF revisions. Port of isPdfSignatureDictionaryConsistent().
func (w *PDFRevisionWrapper) IsPdfSignatureDictionaryConsistent() bool {
	return w.pdfRevision.PDFSignatureDictionary.Consistent
}

// DocMDPPermissions returns a CertificationPermission value of a /DocMDP dictionary, when
// present. Port of getDocMDPPermissions().
func (w *PDFRevisionWrapper) DocMDPPermissions() enumerations.CertificationPermission {
	docMDP := w.pdfRevision.PDFSignatureDictionary.DocMDP
	if docMDP != nil && docMDP.Permissions != nil {
		return enumerations.CertificationPermission(*docMDP.Permissions)
	}
	return ""
}

// FieldMDP returns a /FieldMDP dictionary content, when present. Port of getFieldMDP().
func (w *PDFRevisionWrapper) FieldMDP() *jaxb.XmlPDFLockDictionary {
	return w.pdfRevision.PDFSignatureDictionary.FieldMDP
}

// SigFieldLock returns a /SigFieldLock dictionary, when present. Port of getSigFieldLock().
func (w *PDFRevisionWrapper) SigFieldLock() *jaxb.XmlPDFLockDictionary {
	for _, field := range w.pdfRevision.SignatureField {
		if field.SigFieldLock != nil {
			return field.SigFieldLock
		}
	}
	return nil
}

func concernedPages(xmlModifications []*jaxb.XmlModification) []*big.Int {
	var pages []*big.Int
	for _, modification := range xmlModifications {
		pages = append(pages, modification.Page.BigInt())
	}
	return pages
}

func (w *PDFRevisionWrapper) pdfObjectModifications() *jaxb.XmlObjectModifications {
	modificationDetection := w.pdfRevision.ModificationDetection
	if modificationDetection != nil {
		return modificationDetection.ObjectModifications
	}
	return nil
}

func modifiedFieldNames(objectModifications []*jaxb.XmlObjectModification) []string {
	var names []string
	for _, objectModification := range objectModifications {
		fieldName := objectModification.FieldName
		// Java guards on null only: an empty field name is a name, and it reaches
		// AbstractPdfLockDictionaryCheck#process(), whose Level.FAIL branches turn an
		// empty modified-field list into a PASS. Dropping "" here made SigFieldLockCheck
		// and FieldMDPCheck answer OK where upstream answers NOT OK.
		if fieldName != nil {
			names = append(names, *fieldName)
		}
	}
	return names
}
