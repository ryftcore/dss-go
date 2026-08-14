// Native implementation of eu.europa.esig.dss.pdf.PdfArray on internal/pdf.
//
// Behavioural reference: eu.europa.esig.dss.pdf.pdfbox.PdfBoxArray, method for method. A COSArray
// becomes a pdf.Array; because a Go slice is a value, the wrapper keeps a pointer to it so that
// AddObject is visible to whoever handed the array over, which is what COSArray#add gives
// upstream.
package pades

import (
	"fmt"

	"github.com/utain/esig/dss/internal/pdf"
	"github.com/utain/esig/dss/model"
)

// nativePdfArray wraps a pdf.Array as a PdfArray.
type nativePdfArray struct {
	// document resolves indirect references.
	document *pdf.Document

	// wrapped points at the array itself.
	wrapped *pdf.Array

	// parent is the dictionary or array this one was reached through.
	parent PdfObject
}

// newNativePdfArray wraps an array reached directly. Port of PdfBoxArray(COSArray, PDDocument).
func newNativePdfArray(document *pdf.Document, wrapped pdf.Array) *nativePdfArray {
	return newNativePdfArrayWithParent(document, wrapped, nil)
}

// newNativePdfArrayWithParent wraps an array reached through a parent object.
// Port of PdfBoxArray(COSArray, PDDocument, PdfObject).
func newNativePdfArrayWithParent(document *pdf.Document, wrapped pdf.Array, parent PdfObject) *nativePdfArray {
	if wrapped == nil {
		wrapped = pdf.Array{}
	}
	if document == nil {
		panic("Pdf document shall be provided!")
	}
	return &nativePdfArray{document: document, wrapped: &wrapped, parent: parent}
}

// Value returns the wrapped array. Port of #getValue.
func (a *nativePdfArray) Value() any {
	return *a.wrapped
}

// Parent returns the object this array was reached through. Port of #getParent.
func (a *nativePdfArray) Parent() PdfObject {
	return a.parent
}

// Size returns the number of elements. Port of #size.
func (a *nativePdfArray) Size() int {
	return len(*a.wrapped)
}

// StreamBytes returns the decoded content of the stream at the given position.
// Port of #getStreamBytes together with the private toBytes.
func (a *nativePdfArray) StreamBytes(i int) ([]byte, error) {
	value := a.object(i)
	stream, ok := value.(*pdf.Stream)
	if !ok {
		return nil, model.NewDSSError(fmt.Sprintf("Cannot find value for %v of class %T", value, value))
	}
	return a.document.StreamData(stream)
}

// ObjectKey returns the key of the indirect reference at the given position, nil when the entry
// is direct. Port of #getObjectKey.
func (a *nativePdfArray) ObjectKey(i int) PdfObjectKey {
	if key, ok := a.document.IndexRef(*a.wrapped, i); ok {
		return NewNativePdfObjectKey(key)
	}
	return nil
}

// Number returns the numeric entry at the given position, nil when it is not a number. Java
// returns a Long for a COSInteger and a Float for any other COSNumber; Go returns a *float64,
// which covers both. Port of #getNumber.
func (a *nativePdfArray) Number(i int) *float64 {
	switch value := a.object(i).(type) {
	case pdf.Integer:
		number := float64(value)
		return &number
	case pdf.Real:
		number := value.Val
		return &number
	default:
		return nil
	}
}

// String returns the text of the string entry at the given position, "" when it is not a string.
// Port of #getString.
func (a *nativePdfArray) String(i int) string {
	if value, ok := a.object(i).(pdf.String); ok {
		return nativePdfDictDecodeText(value.Bytes)
	}
	return ""
}

// AsDict returns the entry at the given position as a dictionary, nil when it is not one.
// Port of #getAsDict.
func (a *nativePdfArray) AsDict(i int) PdfDict {
	key, _ := a.document.IndexRef(*a.wrapped, i)
	switch value := a.object(i).(type) {
	case *pdf.Dict:
		return newNativePdfDictWithParent(a.document, value, key, a)
	case *pdf.Stream:
		return newNativePdfStreamDict(a.document, value, key, a)
	default:
		// Upstream logs "Unable to extract array entry as dictionary!".
		return nil
	}
}

// Object returns the entry at the given position as a PdfObject, nil when absent, null or of an
// unsupported type. Port of #getObject.
func (a *nativePdfArray) Object(i int) PdfObject {
	switch value := a.object(i).(type) {
	case nil:
		return nil
	case *pdf.Dict, *pdf.Stream:
		return a.AsDict(i)
	case pdf.Array:
		return newNativePdfArrayWithParent(a.document, value, a)
	case pdf.String:
		return NewPdfSimpleObject(a.String(i), a)
	case pdf.Name:
		return NewPdfSimpleObject(string(value), a)
	case pdf.Integer, pdf.Real:
		return NewPdfSimpleObject(a.Number(i), a)
	case pdf.Bool:
		return NewPdfSimpleObject(bool(value), a)
	default:
		// pdf.Null, or an entry of a type upstream logs
		// "Unable to process an entry on position '{}' of type '{}'." for.
		return nil
	}
}

// AddObject appends the given object. Port of #addObject.
func (a *nativePdfArray) AddObject(pdfObject PdfObject) {
	value, ok := pdfObject.Value().(pdf.Object)
	if !ok {
		panic("The object to be added shall be of type COSBase!")
	}
	*a.wrapped = append(*a.wrapped, value)
}

// SetDirect marks the array as direct. Port of #setDirect.
//
// DEVIATION: see nativePdfDict.SetDirect - internal/pdf derives direct/indirect from the object
// graph, so there is no flag to set.
func (a *nativePdfArray) SetDirect(direct bool) {
	// no-op; see the doc comment
}

// object resolves the element at the given position, nil when out of range.
// It is COSArray#getObject(int).
func (a *nativePdfArray) object(i int) pdf.Object {
	if i < 0 || i >= len(*a.wrapped) {
		return nil
	}
	return a.document.Resolve((*a.wrapped)[i])
}

// Compile-time assertion standing in for Java's "implements PdfArray".
var _ PdfArray = (*nativePdfArray)(nil)
