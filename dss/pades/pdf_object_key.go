// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PdfObjectKey.java (DSS 6.5.RC1).
//
// Java expresses this as a backend-abstraction interface: pdfbox and openpdf each supply their
// own implementation carrying whatever native object-reference representation that backend
// uses, with equals()/hashCode() over the (object number, generation) pair so the type works as
// a HashMap key. This port has a single native PDF engine (internal/pdf; see its doc.go), so
// there is exactly one implementation: NativePdfObjectKey (native_pdf_object_key.go).
//
// INTEGRATION CLEANUP: this file originally also declared a second, unused implementation
// (nativePdfObjectKey + NewPdfObjectKey, returning *PdfObjectKey) speculatively, before
// native_pdf_array.go / native_pdf_dict.go landed and settled on NativePdfObjectKey /
// NewNativePdfObjectKey instead (both files' ObjectKey() methods return the bare PdfObjectKey
// interface, nil for "no indirect reference" - not the *PdfObjectKey this file assumed). That
// second implementation was dead code and has been removed; only the interface, the one shape
// every file in this package actually depends on, remains here.
//
// PdfObjectKey is used directly (not by pointer) as a Go map key - interfaces are valid,
// comparable map keys as long as the dynamic value is comparable, which a (number, generation)
// pair is.
package pades

// PdfObjectKey represents a PDF object identifier within a PDF document.
type PdfObjectKey interface {
	// Value gets the format specific object reference value. Port of getValue().
	Value() any

	// Number gets object's key number. Port of getNumber().
	Number() int64

	// Generation gets generation of the Pdf object. Port of getGeneration().
	Generation() int
}
