// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfSignatureField.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf flattens into the Go package pades together with the rest of the phase
// 5b layout (see internal/pdf/doc.go's Layering section: "pades (ported Java classes: PAdESUtils,
// PdfSigDictWrapper, SingleDssDict, ByteRange, …)"), so the type keeps its Java name unqualified.
//
// java.io.Serializable is dropped (no Go counterpart).
//
// FORWARD DEPENDENCIES (not in this chunk's manifest):
//   - PdfDict (eu.europa.esig.dss.pdf.PdfDict), already landed by the SIGN chunk
//     (pades/native_pdf_dict.go implements it); this file uses its StringValue(name string)
//     string and AsDict(name string) PdfDict methods.
//   - SigFieldPermissions (eu.europa.esig.dss.pdf.SigFieldPermissions) - a struct with
//     Action() enumerations.PdfLockAction, Fields() []string and
//     CertificationPermission() enumerations.CertificationPermission.
//   - PAdESUtilsExtractPermissionsDictionary(lock PdfDict) *SigFieldPermissions - the flattened
//     static PAdESUtils.extractPermissionsDictionary(PdfDict), following the
//     PAdESUtilsVRIsWithName precedent (pdf_composite_dss_dict_certificate_source.go) for how a
//     flattened PAdESUtils static method is named in this port.
//   - PAdESConstantsFieldNameName / PAdESConstantsLockName (eu.europa.esig.dss.pdf.PAdESConstants
//     .FIELD_NAME_NAME / .LOCK_NAME) - matching the PAdESConstants<FieldName> naming already
//     observed in the landed pades/native_pdf_signature_service.go and
//     pades/native_pdf_document_reader.go (e.g. PAdESConstantsTimestampType,
//     PAdESConstantsSignatureDefaultSubFilter, PAdESConstantsURName).
package pades

import "fmt"

// PdfSignatureField represents a PDF Signature field.
type PdfSignatureField struct {
	// fieldName is the name of the signature field.
	fieldName string

	// lockDictionary is the /Lock dictionary content, nil when absent.
	lockDictionary *SigFieldPermissions
}

// NewPdfSignatureField is the default constructor. Port of the constructor
// PdfSignatureField(PdfDict).
//
// Panics with the Java message when sigFieldDict is nil (Objects.requireNonNull).
func NewPdfSignatureField(sigFieldDict PdfDict) *PdfSignatureField {
	if sigFieldDict == nil {
		panic("sigFieldDict cannot be null!")
	}
	return &PdfSignatureField{
		fieldName:      pdfSignatureFieldExtractFieldName(sigFieldDict),
		lockDictionary: pdfSignatureFieldExtractLockDictionary(sigFieldDict),
	}
}

// pdfSignatureFieldExtractFieldName ports the private static extractFieldName(PdfDict).
func pdfSignatureFieldExtractFieldName(sigFieldDict PdfDict) string {
	return sigFieldDict.StringValue(PAdESConstantsFieldNameName)
}

// pdfSignatureFieldExtractLockDictionary ports the private static
// extractLockDictionary(PdfDict).
func pdfSignatureFieldExtractLockDictionary(sigFieldDict PdfDict) *SigFieldPermissions {
	lock := sigFieldDict.AsDict(PAdESConstantsLockName)
	if lock != nil {
		return PAdESUtilsExtractPermissionsDictionary(lock)
	}
	return nil
}

// FieldName returns a signature field's name. Port of getFieldName().
func (f *PdfSignatureField) FieldName() string {
	return f.fieldName
}

// LockDictionary returns a /Lock dictionary content, when present. Port of getLockDictionary().
func (f *PdfSignatureField) LockDictionary() *SigFieldPermissions {
	return f.lockDictionary
}

// String ports toString().
func (f *PdfSignatureField) String() string {
	return fmt.Sprintf("PdfSignatureField {name=%s}", f.FieldName())
}
