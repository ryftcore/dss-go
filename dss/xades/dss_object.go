// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/DSSObject.java (DSS 6.5.RC1).
//
// java.io.Serializable and serialVersionUID are dropped (no Go counterpart).
package xades

import (
	"fmt"

	"github.com/utain/esig/dss/model"
)

// DSSObject allows creation of a custom ds:Object element.
type DSSObject struct {
	// content represents the content of the ds:Object element. Can be XML or any other format
	// (e.g. base64 encoded).
	content model.DSSDocument

	// id represents a value for the "Id" attribute.
	id string

	// mimeType represents a value for the "MimeType" attribute.
	mimeType string

	// encodingAlgorithm represents a value for the "Encoding" attribute.
	encodingAlgorithm string
}

// NewDSSObject is the default constructor.
func NewDSSObject() *DSSObject {
	return &DSSObject{}
}

// Content gets the content of the ds:Object element to be created. Ports getContent().
func (o *DSSObject) Content() model.DSSDocument {
	return o.content
}

// SetContent sets the content of ds:Object element to be created. Can be XML or any other
// format (e.g. base64 encoded). Ports setContent(DSSDocument).
func (o *DSSObject) SetContent(content model.DSSDocument) {
	o.content = content
}

// Id gets the Id. Ports getId().
func (o *DSSObject) Id() string {
	return o.id
}

// SetId sets the value for the "Id" attribute. Ports setId(String).
func (o *DSSObject) SetId(id string) {
	o.id = id
}

// MimeType gets the MimeType. Ports getMimeType().
func (o *DSSObject) MimeType() string {
	return o.mimeType
}

// SetMimeType sets the value for the "MimeType" attribute. Ports setMimeType(String).
func (o *DSSObject) SetMimeType(mimeType string) {
	o.mimeType = mimeType
}

// EncodingAlgorithm gets the encoding algorithm. Ports getEncodingAlgorithm().
func (o *DSSObject) EncodingAlgorithm() string {
	return o.encodingAlgorithm
}

// SetEncodingAlgorithm sets the value for the "encoding" attribute. Ports
// setEncodingAlgorithm(String).
func (o *DSSObject) SetEncodingAlgorithm(encodingAlgorithm string) {
	o.encodingAlgorithm = encodingAlgorithm
}

// String ports toString().
func (o *DSSObject) String() string {
	return fmt.Sprintf("DSSObject{content=%v, id='%s', mimeType=%v, encodingAlgorithm='%s'}",
		o.content, o.id, o.mimeType, o.encodingAlgorithm)
}
