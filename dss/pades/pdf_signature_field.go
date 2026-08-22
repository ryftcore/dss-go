// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfSignatureField.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf flattens into the Go package pades (see internal/pdf/doc.go's Layering
// section: "pades (ported Java classes: PAdESUtils, PdfSigDictWrapper, SingleDssDict, ByteRange,
// …)"), so the type keeps its Java name unqualified.
//
// java.io.Serializable is dropped (no Go counterpart).
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
