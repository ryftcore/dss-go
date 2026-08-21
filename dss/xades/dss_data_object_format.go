// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/dataobject/DSSDataObjectFormat.java (DSS 6.5.RC1).
//
// java.io.Serializable and serialVersionUID are dropped (no Go counterpart).
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// DSSDataObjectFormat represents a <xades:DataObjectFormat> element as part of
// <xades:SignedDataObjectProperties>.
type DSSDataObjectFormat struct {
	// description describes the data object.
	description string

	// objectIdentifier provides an identifier to the data object.
	objectIdentifier enumerations.ObjectIdentifier

	// mimeType defines the MimeType of the data object.
	mimeType string

	// encoding defines the encoding of the data object.
	encoding string

	// objectReference provides a reference to the data object.
	objectReference string
}

// NewDSSDataObjectFormat is the empty constructor.
func NewDSSDataObjectFormat() *DSSDataObjectFormat {
	return &DSSDataObjectFormat{}
}

// Description gets description of the data object. Ports getDescription().
func (f *DSSDataObjectFormat) Description() string {
	return f.description
}

// SetDescription sets description of the data object. Ports setDescription(String).
func (f *DSSDataObjectFormat) SetDescription(description string) {
	f.description = description
}

// ObjectIdentifier gets object identifier (reference) of the data object. Ports
// getObjectIdentifier().
func (f *DSSDataObjectFormat) ObjectIdentifier() enumerations.ObjectIdentifier {
	return f.objectIdentifier
}

// SetObjectIdentifier sets object identifier (reference) of the data object. Ports
// setObjectIdentifier(ObjectIdentifier).
func (f *DSSDataObjectFormat) SetObjectIdentifier(objectIdentifier enumerations.ObjectIdentifier) {
	f.objectIdentifier = objectIdentifier
}

// MimeType gets MimeType of the data object. Ports getMimeType().
func (f *DSSDataObjectFormat) MimeType() string {
	return f.mimeType
}

// SetMimeType sets MimeType of the data object. Ports setMimeType(String).
func (f *DSSDataObjectFormat) SetMimeType(mimeType string) {
	f.mimeType = mimeType
}

// Encoding gets encoding of the data object. Ports getEncoding().
func (f *DSSDataObjectFormat) Encoding() string {
	return f.encoding
}

// SetEncoding sets encoding of the data object. Ports setEncoding(String).
func (f *DSSDataObjectFormat) SetEncoding(encoding string) {
	f.encoding = encoding
}

// ObjectReference gets reference to the data object. Ports getObjectReference().
func (f *DSSDataObjectFormat) ObjectReference() string {
	return f.objectReference
}

// SetObjectReference sets reference to the data object. Ports setObjectReference(String).
func (f *DSSDataObjectFormat) SetObjectReference(objectReference string) {
	f.objectReference = objectReference
}

// Equals ports equals(Object).
func (f *DSSDataObjectFormat) Equals(other *DSSDataObjectFormat) bool {
	if f == other {
		return true
	}
	if other == nil {
		return false
	}
	return f.description == other.description &&
		f.objectIdentifier == other.objectIdentifier &&
		f.mimeType == other.mimeType &&
		f.encoding == other.encoding &&
		f.objectReference == other.objectReference
}

// String ports toString().
func (f *DSSDataObjectFormat) String() string {
	return fmt.Sprintf("DSSDataObjectFormat{description='%s', objectIdentifier=%v, mimeType='%s', "+
		"encoding='%s', objectReference='%s'}", f.description, f.objectIdentifier, f.mimeType, f.encoding, f.objectReference)
}
