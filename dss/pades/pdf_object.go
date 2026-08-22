// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfObject.java,
// dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfDict.java,
// dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfArray.java and
// dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfSimpleObject.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf is implemented by this file, together with the other pdf_*.go /
// *_checker.go / annotation_box.go / image_rotation_utils.go / pades_constants.go /
// pdf_memory_usage_setting.go / pdf_service_mode.go / sig_field_permissions.go /
// pdf_signature_cache.go files. Where usage and this package's shape diverge only in
// inessential ways (return-by-value vs return-by-pointer, etc.), the shape actually exercised by
// native_pdf_dict.go / native_pdf_array.go wins, since those are the sole concrete implementers.
package pades

import (
	"io"
	"time"
)

// PdfObject represents a PDF internal object. Port of the PdfObject interface.
type PdfObject interface {
	// Value gets the value of the PDF object. Port of getValue().
	Value() any

	// Parent returns the parent of the current PdfObject if applicable. Port of getParent().
	Parent() PdfObject
}

// PdfDict permits the user to choose the underlying PDF library used to create PDF signatures.
// Port of the PdfDict interface.
type PdfDict interface {
	PdfObject

	// AsDict gets an embedded dictionary by name. Port of getAsDict(String).
	AsDict(name string) PdfDict

	// AsArray gets the PdfArray by name. Port of getAsArray(String).
	AsArray(name string) PdfArray

	// BinariesValue gets binaries by dictionary name. Port of getBinariesValue(String).
	BinariesValue(name string) ([]byte, error)

	// List lists all encapsulated dictionary names. Port of list().
	List() []string

	// StringValue gets a string value by property name. Port of getStringValue(String).
	StringValue(name string) string

	// NameValue gets a name of the dictionary. Port of getNameValue(String).
	NameValue(name string) string

	// DateValue gets a date. The zero time.Time stands for Java's null. Port of getDateValue(String).
	DateValue(name string) time.Time

	// NumberValue returns a number value, nil when absent. Port of getNumberValue(String).
	NumberValue(name string) *float64

	// Object returns a PDF object. Port of getObject(String).
	Object(name string) PdfObject

	// ObjectKey returns the Pdf object key of an indirect reference to an object, when
	// applicable, nil otherwise. Port of getObjectKey(String).
	ObjectKey(name string) PdfObjectKey

	// StreamBytes returns the object's stream binaries, when available. Port of getStreamBytes().
	StreamBytes() ([]byte, error)

	// CreateRawInputStream creates a new raw input stream (raw, not decoded data).
	// Port of createRawInputStream().
	CreateRawInputStream() (io.ReadCloser, error)

	// RawStreamSize returns the size of the raw stream, if present, -1 if not applicable.
	// Port of getRawStreamSize().
	RawStreamSize() (int64, error)

	// SetPdfObjectValue sets the Dictionary pdfObject with the given key.
	// Port of setPdfObjectValue(String, PdfObject).
	SetPdfObjectValue(key string, pdfObject PdfObject)

	// SetNameValue sets the Name value with the given key. Port of setNameValue(String, String).
	SetNameValue(key string, value string)

	// SetStringValue sets the String value with the given key. Port of setStringValue(String, String).
	SetStringValue(key string, value string)

	// SetIntegerValue sets the Integer value with the given key. Port of setIntegerValue(String, Integer).
	SetIntegerValue(key string, value int)

	// SetDirect sets whether the object shall be written directly to its parent. Port of setDirect(boolean).
	SetDirect(direct bool)

	// Match verifies if the content of pdfDict matches the corresponding attributes of the
	// current dictionary. Unlike Equals, this does not ensure full equality. Port of match(PdfDict).
	Match(pdfDict PdfDict) bool
}

// PdfArray permits the user to choose the underlying PDF library used to create PDF signatures.
// Port of the PdfArray interface.
type PdfArray interface {
	PdfObject

	// Size retrieves the array size. Port of size().
	Size() int

	// StreamBytes retrieves the stream byte array at the position i. Port of getStreamBytes(int).
	StreamBytes(i int) ([]byte, error)

	// ObjectKey retrieves the object key for the position i. Port of getObjectKey(int).
	ObjectKey(i int) PdfObjectKey

	// Number retrieves the number at the position i, nil when it is not one. Port of getNumber(int).
	Number(i int) *float64

	// String returns a String entry at the position i. Port of getString(int).
	String(i int) string

	// AsDict returns a dictionary entry at the position i. Port of getAsDict(int).
	AsDict(i int) PdfDict

	// Object returns an object entry at the position i. Port of getObject(int).
	Object(i int) PdfObject

	// AddObject adds pdfObject. Port of addObject(PdfObject).
	AddObject(pdfObject PdfObject)

	// SetDirect sets whether the array shall be written directly to its parent. Port of setDirect(boolean).
	SetDirect(direct bool)
}

// pdfSimpleObject is a wrapper for a simple value (Integer, String, etc.), extracted from a PDF.
// Port of PdfSimpleObject.
type pdfSimpleObject struct {
	// value is the value of the object.
	value any

	// parent is the parent of the object.
	parent PdfObject
}

// NewPdfSimpleObject wraps a simple value with an optional parent. Java splits this into a
// 1-arg constructor (parent nil) and a 2-arg constructor; Go collapses them, matching every
// call site's signature (native_pdf_dict.go, native_pdf_array.go).
// Port of PdfSimpleObject(Object) / PdfSimpleObject(Object, PdfObject).
func NewPdfSimpleObject(value any, parent PdfObject) *pdfSimpleObject {
	return &pdfSimpleObject{value: value, parent: parent}
}

// Value returns the wrapped value. Port of #getValue.
func (o *pdfSimpleObject) Value() any {
	return o.value
}

// Parent returns the parent of the object. Port of #getParent.
func (o *pdfSimpleObject) Parent() PdfObject {
	return o.parent
}

// Compile-time assertion standing in for Java's "implements PdfObject".
var _ PdfObject = (*pdfSimpleObject)(nil)
