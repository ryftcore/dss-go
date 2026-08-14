// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/PDFRevisionWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"math/big"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
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

// GetPdfAnnotationsOverlapConcernedPages returns a list of PDF annotation overlap concerned
// pages. Port of getPdfAnnotationsOverlapConcernedPages().
func (w *PDFRevisionWrapper) GetPdfAnnotationsOverlapConcernedPages() []*big.Int {
	modificationDetection := w.pdfRevision.ModificationDetection
	if modificationDetection != nil {
		return getConcernedPages(modificationDetection.AnnotationOverlap)
	}
	return nil
}

// GetPdfVisualDifferenceConcernedPages returns a list of PDF visual difference concerned
// pages. Port of getPdfVisualDifferenceConcernedPages().
func (w *PDFRevisionWrapper) GetPdfVisualDifferenceConcernedPages() []*big.Int {
	modificationDetection := w.pdfRevision.ModificationDetection
	if modificationDetection != nil {
		return getConcernedPages(modificationDetection.VisualDifference)
	}
	return nil
}

// GetPdfPageDifferenceConcernedPages returns a list of pages missing/added to the final
// revision in a comparison with a signed one. Port of getPdfPageDifferenceConcernedPages().
func (w *PDFRevisionWrapper) GetPdfPageDifferenceConcernedPages() []*big.Int {
	modificationDetection := w.pdfRevision.ModificationDetection
	if modificationDetection != nil {
		return getConcernedPages(modificationDetection.PageDifference)
	}
	return nil
}

// ArePdfObjectModificationsDetected checks whether object modifications are present after the
// current PDF revisions. Port of arePdfObjectModificationsDetected().
func (w *PDFRevisionWrapper) ArePdfObjectModificationsDetected() bool {
	return w.getPdfObjectModifications() != nil
}

// GetPdfExtensionChanges returns a list of changes occurred in a PDF after the current
// signature's revision associated with a signature/document extension. Port of
// getPdfExtensionChanges().
func (w *PDFRevisionWrapper) GetPdfExtensionChanges() []*jaxb.XmlObjectModification {
	pdfObjectModifications := w.getPdfObjectModifications()
	if pdfObjectModifications != nil {
		return pdfObjectModifications.ExtensionChanges
	}
	return nil
}

// GetPdfSignatureOrFormFillChanges returns a list of changes occurred in a PDF after the
// current signature's revision associated with a signature creation, form filling. Port of
// getPdfSignatureOrFormFillChanges().
func (w *PDFRevisionWrapper) GetPdfSignatureOrFormFillChanges() []*jaxb.XmlObjectModification {
	pdfObjectModifications := w.getPdfObjectModifications()
	if pdfObjectModifications != nil {
		return pdfObjectModifications.SignatureOrFormFill
	}
	return nil
}

// GetPdfAnnotationChanges returns a list of changes occurred in a PDF after the current
// signature's revision associated with annotation(s) modification. Port of
// getPdfAnnotationChanges().
func (w *PDFRevisionWrapper) GetPdfAnnotationChanges() []*jaxb.XmlObjectModification {
	pdfObjectModifications := w.getPdfObjectModifications()
	if pdfObjectModifications != nil {
		return pdfObjectModifications.AnnotationChanges
	}
	return nil
}

// GetPdfUndefinedChanges returns a list of undefined changes occurred in a PDF after the
// current signature's revision. Port of getPdfUndefinedChanges().
func (w *PDFRevisionWrapper) GetPdfUndefinedChanges() []*jaxb.XmlObjectModification {
	pdfObjectModifications := w.getPdfObjectModifications()
	if pdfObjectModifications != nil {
		return pdfObjectModifications.Undefined
	}
	return nil
}

// GetModifiedFieldNames returns a list of field names modified after the current signature's
// revision. Port of getModifiedFieldNames().
func (w *PDFRevisionWrapper) GetModifiedFieldNames() []string {
	var names []string
	pdfObjectModifications := w.getPdfObjectModifications()
	if pdfObjectModifications != nil {
		names = append(names, getModifiedFieldNames(pdfObjectModifications.ExtensionChanges)...)
		names = append(names, getModifiedFieldNames(pdfObjectModifications.SignatureOrFormFill)...)
		names = append(names, getModifiedFieldNames(pdfObjectModifications.AnnotationChanges)...)
		names = append(names, getModifiedFieldNames(pdfObjectModifications.Undefined)...)
	}
	return names
}

// GetFirstFieldName returns the first signature field name. Port of getFirstFieldName().
func (w *PDFRevisionWrapper) GetFirstFieldName() string {
	fields := w.pdfRevision.Fields
	if len(fields) != 0 {
		return fields[0].Name
	}
	return ""
}

// GetSignatureFieldNames returns a list of signature field names, where the signature is
// referenced from. Port of getSignatureFieldNames().
func (w *PDFRevisionWrapper) GetSignatureFieldNames() []string {
	var names []string
	fields := w.pdfRevision.Fields
	if len(fields) != 0 {
		for _, signatureField := range fields {
			names = append(names, signatureField.Name)
		}
	}
	return names
}

// GetSignerName returns the signer's name. Port of getSignerName().
func (w *PDFRevisionWrapper) GetSignerName() string {
	return w.pdfRevision.PDFSignatureDictionary.SignerName
}

// GetSignatureDictionaryType returns the PDF signature dictionary /Type value. Port of
// getSignatureDictionaryType().
func (w *PDFRevisionWrapper) GetSignatureDictionaryType() string {
	return w.pdfRevision.PDFSignatureDictionary.Type
}

// GetFilter returns the PDF signature dictionary /Filter value. Port of getFilter().
func (w *PDFRevisionWrapper) GetFilter() string {
	return w.pdfRevision.PDFSignatureDictionary.Filter
}

// GetSubFilter returns the PDF signature dictionary /SubFilter value. Port of getSubFilter().
func (w *PDFRevisionWrapper) GetSubFilter() string {
	return w.pdfRevision.PDFSignatureDictionary.SubFilter
}

// GetContactInfo returns the PDF signature dictionary /ContactInfo value. Port of
// getContactInfo().
func (w *PDFRevisionWrapper) GetContactInfo() string {
	return w.pdfRevision.PDFSignatureDictionary.ContactInfo
}

// GetLocation returns the PDF signature dictionary /Location value. Port of getLocation().
func (w *PDFRevisionWrapper) GetLocation() string {
	return w.pdfRevision.PDFSignatureDictionary.Location
}

// GetReason returns the PDF signature dictionary /Reason value. Port of getReason().
func (w *PDFRevisionWrapper) GetReason() string {
	return w.pdfRevision.PDFSignatureDictionary.Reason
}

// GetSignatureByteRange returns the PDF signature dictionary /ByteRange value. Port of
// getSignatureByteRange().
func (w *PDFRevisionWrapper) GetSignatureByteRange() []*big.Int {
	byteRange := w.getXmlByteRange()
	if byteRange != nil {
		return byteRange.Value
	}
	return nil
}

// IsSignatureByteRangeValid returns whether the PDF signature dictionary /ByteRange is found
// and valid. Port of isSignatureByteRangeValid().
func (w *PDFRevisionWrapper) IsSignatureByteRangeValid() bool {
	byteRange := w.getXmlByteRange()
	if byteRange != nil {
		return byteRange.Valid
	}
	return false
}

func (w *PDFRevisionWrapper) getXmlByteRange() *jaxb.XmlByteRange {
	return w.pdfRevision.PDFSignatureDictionary.SignatureByteRange
}

// IsPdfSignatureDictionaryConsistent returns whether the PDF signature dictionary is
// consistent across PDF revisions. Port of isPdfSignatureDictionaryConsistent().
func (w *PDFRevisionWrapper) IsPdfSignatureDictionaryConsistent() bool {
	return w.pdfRevision.PDFSignatureDictionary.Consistent
}

// GetDocMDPPermissions returns a CertificationPermission value of a /DocMDP dictionary, when
// present. Port of getDocMDPPermissions().
func (w *PDFRevisionWrapper) GetDocMDPPermissions() enumerations.CertificationPermission {
	docMDP := w.pdfRevision.PDFSignatureDictionary.DocMDP
	if docMDP != nil {
		return docMDP.Permissions
	}
	return ""
}

// GetFieldMDP returns a /FieldMDP dictionary content, when present. Port of getFieldMDP().
func (w *PDFRevisionWrapper) GetFieldMDP() *jaxb.XmlPDFLockDictionary {
	return w.pdfRevision.PDFSignatureDictionary.FieldMDP
}

// GetSigFieldLock returns a /SigFieldLock dictionary, when present. Port of getSigFieldLock().
func (w *PDFRevisionWrapper) GetSigFieldLock() *jaxb.XmlPDFLockDictionary {
	for _, field := range w.pdfRevision.Fields {
		if field.SigFieldLock != nil {
			return field.SigFieldLock
		}
	}
	return nil
}

func getConcernedPages(xmlModifications []*jaxb.XmlModification) []*big.Int {
	var pages []*big.Int
	for _, modification := range xmlModifications {
		pages = append(pages, modification.Page)
	}
	return pages
}

func (w *PDFRevisionWrapper) getPdfObjectModifications() *jaxb.XmlObjectModifications {
	modificationDetection := w.pdfRevision.ModificationDetection
	if modificationDetection != nil {
		return modificationDetection.ObjectModifications
	}
	return nil
}

func getModifiedFieldNames(objectModifications []*jaxb.XmlObjectModification) []string {
	var names []string
	for _, objectModification := range objectModifications {
		fieldName := objectModification.FieldName
		if fieldName != "" {
			names = append(names, fieldName)
		}
	}
	return names
}
