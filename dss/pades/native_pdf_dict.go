// Native implementation of eu.europa.esig.dss.pdf.PdfDict on internal/pdf.
//
// Behavioural reference: eu.europa.esig.dss.pdf.pdfbox.PdfBoxDict, method for method. A
// COSDictionary becomes a *pdf.Dict, a COSStream becomes a *pdf.Stream (the wrapper carries both,
// as PdfBoxDict does through COSStream extending COSDictionary), and "the document" - needed to
// resolve indirect references, which pdfbox does inside COSObject - becomes the *pdf.Document.
package pades

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/utain/esig/dss/internal/pdf"
)

// nativePdfDict wraps a pdf.Dict as a PdfDict.
type nativePdfDict struct {
	// document resolves indirect references.
	document *pdf.Document

	// wrapped is the dictionary itself. For a stream object it is the stream's dictionary.
	wrapped *pdf.Dict

	// stream is non-nil when the wrapped dictionary belongs to a stream object, which is what
	// PdfBoxDict's "wrapped instanceof COSStream" tests.
	stream *pdf.Stream

	// key is the indirect object this dictionary was loaded from; the zero value means the
	// dictionary is direct or freshly created.
	key pdf.ObjectKey

	// parent is the dictionary or array this one was reached through.
	parent PdfObject
}

// newNativePdfDict wraps a dictionary reached directly. Port of
// PdfBoxDict(COSDictionary, PDDocument).
func newNativePdfDict(document *pdf.Document, wrapped *pdf.Dict, key pdf.ObjectKey) *nativePdfDict {
	return newNativePdfDictWithParent(document, wrapped, key, nil)
}

// newNativePdfDictWithParent wraps a dictionary reached through a parent object. Port of
// PdfBoxDict(COSDictionary, PDDocument, PdfObject).
func newNativePdfDictWithParent(document *pdf.Document, wrapped *pdf.Dict, key pdf.ObjectKey,
	parent PdfObject) *nativePdfDict {
	if wrapped == nil {
		panic("Pdf dictionary shall be provided!")
	}
	if document == nil {
		panic("Pdf document shall be provided!")
	}
	return &nativePdfDict{document: document, wrapped: wrapped, key: key, parent: parent}
}

// newNativePdfStreamDict wraps a stream object. There is no separate upstream constructor: a
// COSStream is a COSDictionary, which Go's distinct pdf.Stream / pdf.Dict types cannot express.
func newNativePdfStreamDict(document *pdf.Document, stream *pdf.Stream, key pdf.ObjectKey,
	parent PdfObject) *nativePdfDict {
	d := newNativePdfDictWithParent(document, stream.Dict, key, parent)
	d.stream = stream
	return d
}

// Value returns the wrapped dictionary. Port of #getValue.
func (d *nativePdfDict) Value() any {
	if d.stream != nil {
		return d.stream
	}
	return d.wrapped
}

// Parent returns the object this dictionary was reached through. Port of #getParent.
func (d *nativePdfDict) Parent() PdfObject {
	return d.parent
}

// AsDict returns the entry with the given name as a dictionary, nil when absent or of another
// type. Port of #getAsDict.
func (d *nativePdfDict) AsDict(name string) PdfDict {
	raw := d.wrapped.GetRaw(pdf.Name(name))
	if raw == nil {
		return nil
	}
	key, _ := d.document.RefAt(d.wrapped, pdf.Name(name))
	switch resolved := d.document.Resolve(raw).(type) {
	case *pdf.Dict:
		return newNativePdfDictWithParent(d.document, resolved, key, d)
	case *pdf.Stream:
		return newNativePdfStreamDict(d.document, resolved, key, d)
	default:
		// Upstream logs "Unable to extract entry with name '{}' as dictionary!".
		return nil
	}
}

// AsArray returns the entry with the given name as an array, nil when absent or of another type.
// Port of #getAsArray.
func (d *nativePdfDict) AsArray(name string) PdfArray {
	if array, ok := d.document.GetArray(d.wrapped, pdf.Name(name)); ok {
		return newNativePdfArrayWithParent(d.document, array, d)
	}
	return nil
}

// BinariesValue returns the entry with the given name as the bytes of a PDF string.
// Port of #getBinariesValue.
func (d *nativePdfDict) BinariesValue(name string) ([]byte, error) {
	if value, ok := d.document.GetString(d.wrapped, pdf.Name(name)); ok {
		return value, nil
	}
	return nil, fmt.Errorf("%s was expected to be a COSString element but was : %v",
		name, d.wrapped.GetRaw(pdf.Name(name)))
}

// List returns the dictionary's keys in insertion order. Port of #list.
func (d *nativePdfDict) List() []string {
	keys := d.wrapped.Keys()
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, string(key))
	}
	return result
}

// StringValue returns the text of the string entry with the given name, "" when absent.
// Port of #getStringValue, i.e. of COSDictionary#getString.
func (d *nativePdfDict) StringValue(name string) string {
	value, ok := d.document.GetString(d.wrapped, pdf.Name(name))
	if !ok {
		return ""
	}
	return nativePdfDictDecodeText(value)
}

// NameValue returns the name entry with the given name as a string, "" when absent.
// Port of #getNameValue, i.e. of COSDictionary#getNameAsString.
func (d *nativePdfDict) NameValue(name string) string {
	if value, ok := d.document.GetName(d.wrapped, pdf.Name(name)); ok {
		return string(value)
	}
	return ""
}

// DateValue returns the date entry with the given name; the zero time means absent.
// Port of #getDateValue.
func (d *nativePdfDict) DateValue(name string) time.Time {
	if value, ok := d.document.GetDate(d.wrapped, pdf.Name(name)); ok {
		return value
	}
	return time.Time{}
}

// NumberValue returns the numeric entry with the given name, nil when absent. Java returns a
// Number (a Float for a COSFloat, a Long for any other COSNumber); Go returns a *float64, which
// covers both and keeps "absent" expressible.
// Port of #getNumberValue.
func (d *nativePdfDict) NumberValue(name string) *float64 {
	if value, ok := d.document.GetInt(d.wrapped, pdf.Name(name)); ok {
		number := float64(value)
		return &number
	}
	if value, ok := d.document.GetReal(d.wrapped, pdf.Name(name)); ok {
		number := value
		return &number
	}
	return nil
}

// Object returns the entry with the given name as a PdfObject, nil when absent, null or of an
// unsupported type. Port of #getObject.
func (d *nativePdfDict) Object(name string) PdfObject {
	raw := d.wrapped.GetRaw(pdf.Name(name))
	if raw == nil {
		return nil
	}
	switch resolved := d.document.Resolve(raw).(type) {
	case *pdf.Dict, *pdf.Stream:
		_ = resolved
		return d.AsDict(name)
	case pdf.Array:
		return d.AsArray(name)
	case pdf.String:
		return NewPdfSimpleObject(d.StringValue(name), d)
	case pdf.Name:
		return NewPdfSimpleObject(d.NameValue(name), d)
	case pdf.Integer, pdf.Real:
		return NewPdfSimpleObject(d.NumberValue(name), d)
	case pdf.Bool:
		return NewPdfSimpleObject(bool(resolved), d)
	default:
		// pdf.Null, or an entry of a type upstream logs
		// "Unable to process an entry with name '{}' of type '{}'." for.
		return nil
	}
}

// ObjectKey returns the key of the indirect reference stored under the given name, nil when the
// entry is direct or absent. Port of #getObjectKey.
func (d *nativePdfDict) ObjectKey(name string) PdfObjectKey {
	if key, ok := d.document.RefAt(d.wrapped, pdf.Name(name)); ok {
		return NewNativePdfObjectKey(key)
	}
	return nil
}

// StreamBytes returns the decoded stream content, nil when the dictionary is not a stream.
// Port of #getStreamBytes.
func (d *nativePdfDict) StreamBytes() ([]byte, error) {
	if d.stream == nil {
		return nil, nil
	}
	return d.document.StreamData(d.stream)
}

// CreateRawInputStream returns the still-encoded stream content, nil when the dictionary is not
// a stream. Port of #createRawInputStream.
func (d *nativePdfDict) CreateRawInputStream() (io.ReadCloser, error) {
	if d.stream == nil {
		return nil, nil
	}
	raw, err := d.document.RawStreamData(d.stream)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(raw)), nil
}

// RawStreamSize returns the size of the still-encoded stream content, -1 when the dictionary is
// not a stream. Port of #getRawStreamSize.
func (d *nativePdfDict) RawStreamSize() (int64, error) {
	if d.stream == nil {
		return -1, nil
	}
	return d.document.RawStreamSize(d.stream), nil
}

// SetPdfObjectValue sets the given object under the given key. Port of #setPdfObjectValue.
func (d *nativePdfDict) SetPdfObjectValue(key string, pdfObject PdfObject) {
	value, ok := pdfObject.Value().(pdf.Object)
	if !ok {
		panic("pdfObject argument shall be of COSBase type!")
	}
	d.wrapped.Set(pdf.Name(key), value)
}

// SetNameValue sets a name entry. Port of #setNameValue.
func (d *nativePdfDict) SetNameValue(key string, value string) {
	d.wrapped.Set(pdf.Name(key), pdf.Name(value))
}

// SetStringValue sets a string entry. Port of #setStringValue.
func (d *nativePdfDict) SetStringValue(key string, value string) {
	d.wrapped.Set(pdf.Name(key), pdf.String{Bytes: []byte(value)})
}

// SetIntegerValue sets an integer entry. Port of #setIntegerValue.
func (d *nativePdfDict) SetIntegerValue(key string, value int) {
	d.wrapped.Set(pdf.Name(key), pdf.Integer(value))
}

// SetDirect marks the dictionary as direct, i.e. to be written inside its parent rather than as
// an indirect object. Port of #setDirect.
//
// DEVIATION: internal/pdf decides direct/indirect from the object graph the Updater is handed -
// a dictionary reached by value is written inline, one reached through a pdf.Ref is written as
// an indirect object - so there is no flag to set and this is a no-op. The one place upstream
// depends on it (AbstractPDFSignatureService#addDeveloperExtension) is implemented natively in
// native_pdf_signature_service.go, where the distinction is expressed by construction.
func (d *nativePdfDict) SetDirect(direct bool) {
	// no-op; see the doc comment
}

// Match reports whether every entry of the given dictionary is present here with an equal value.
// Port of #match.
func (d *nativePdfDict) Match(pdfDict PdfDict) bool {
	other, ok := pdfDict.(*nativePdfDict)
	if !ok {
		panic("pdfDict argument shall be of PdfBoxDict type!")
	}
	for _, key := range other.wrapped.Keys() {
		targetObject := other.document.Resolve(other.wrapped.GetRaw(key))
		currentObject := d.document.Resolve(d.wrapped.GetRaw(key))
		if targetObject != nil && !nativePdfObjectEquals(targetObject, currentObject) {
			return false
		}
	}
	return true
}

// String ports #toString.
func (d *nativePdfDict) String() string {
	return d.wrapped.String()
}

// nativePdfDictDecodeText reproduces COSString#getString: a UTF-16 byte-order mark selects
// UTF-16 (BE or LE), anything else is PDFDocEncoding. Delegated to internal/pdf's own
// DecodeTextString so the PAdES layer and the parser layer cannot drift apart; it used to
// approximate PDFDocEncoding with Latin-1, which is wrong for every code in 0x18-0x1F and
// 0x80-0xA0 - see that function's doc comment for the fixture this fixes.
func nativePdfDictDecodeText(value []byte) string {
	return pdf.DecodeTextString(value)
}

// nativePdfObjectEquals compares two resolved PDF objects by value, which is what
// COSBase#equals does for the scalar kinds developer-extension matching compares.
func nativePdfObjectEquals(a, b pdf.Object) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	switch left := a.(type) {
	case pdf.Name:
		right, ok := b.(pdf.Name)
		return ok && left == right
	case pdf.Integer:
		right, ok := b.(pdf.Integer)
		return ok && left == right
	case pdf.Real:
		right, ok := b.(pdf.Real)
		return ok && left.Val == right.Val
	case pdf.Bool:
		right, ok := b.(pdf.Bool)
		return ok && left == right
	case pdf.String:
		right, ok := b.(pdf.String)
		return ok && bytes.Equal(left.Bytes, right.Bytes)
	case *pdf.Dict:
		right, ok := b.(*pdf.Dict)
		if !ok || left.Len() != right.Len() {
			return false
		}
		for _, key := range left.Keys() {
			if !nativePdfObjectEquals(left.GetRaw(key), right.GetRaw(key)) {
				return false
			}
		}
		return true
	case pdf.Array:
		right, ok := b.(pdf.Array)
		if !ok || len(left) != len(right) {
			return false
		}
		for i := range left {
			if !nativePdfObjectEquals(left[i], right[i]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}

// Compile-time assertion standing in for Java's "implements PdfDict".
var _ PdfDict = (*nativePdfDict)(nil)
