// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/dataobject/DataObjectFormatBuilder.java (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/utils"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// DataObjectFormatBuilder builds DSSDataObjectFormat objects. The class handles the cases when
// a DataObjectFormat is not required to be created (e.g. in case of a counter-signature).
type DataObjectFormatBuilder struct {
	// references is the collection of references to build DataObjectFormat objects based on.
	references []*DSSReference
}

// NewDataObjectFormatBuilder is the empty constructor.
func NewDataObjectFormatBuilder() *DataObjectFormatBuilder {
	return &DataObjectFormatBuilder{}
}

// SetReferences sets references to be used as a base for DataObjectFormat objects building.
// Ports setReferences(Collection<DSSReference>).
func (b *DataObjectFormatBuilder) SetReferences(references []*DSSReference) *DataObjectFormatBuilder {
	b.references = references
	return b
}

// Build builds a collection of DSSDataObjectFormats corresponding to the provided
// configuration. Ports build().
func (b *DataObjectFormatBuilder) Build() []*DSSDataObjectFormat {
	if utils.IsCollectionEmpty(b.references) {
		return []*DSSDataObjectFormat{}
	}
	result := make([]*DSSDataObjectFormat, 0, len(b.references))
	for _, reference := range b.references {
		if DSSXMLUtilsIsCounterSignatureReferenceType(reference.Type()) {
			// 6.3 Requirements on XAdES signature's elements, qualifying properties and services
			//
			// k) Requirement for DataObjectFormat. One DataObjectFormat shall be generated for each
			// signed data object, except the SignedProperties element, and except if the signature
			// is a baseline signature countersigning a signature. If the signature is a baseline
			// signature countersigning another signature, and if it only signs its own signed
			// properties and the countersigned signature, then it shall not include any
			// DataObjectFormat signed property. If the signature is a baseline signature
			// countersigning another signature and if it signs its own signed properties, the
			// countersigned signature, and other data object(s), then it shall include one
			// DataObjectFormat signed property for each of these other signed data object(s)
			// aforementioned.
			continue
		}
		result = append(result, b.toDataObjectFormat(reference))
	}
	return result
}

// toDataObjectFormat creates a DSSDataObjectFormat based on the given DSSReference object to be
// incorporated to the signature. Panics with the Java message when reference is nil
// (Objects.requireNonNull upstream). Ports the protected toDataObjectFormat(DSSReference).
func (b *DataObjectFormatBuilder) toDataObjectFormat(reference *DSSReference) *DSSDataObjectFormat {
	if reference == nil {
		panic("Reference cannot be null!")
	}
	dataObjectFormat := NewDSSDataObjectFormat()
	if utils.IsStringNotEmpty(reference.Id()) {
		dataObjectFormat.SetObjectReference(xmlutils.DomUtilsToElementReference(reference.Id()))
	}
	dataObjectFormat.SetMimeType(b.getReferenceMimeType(reference))
	return dataObjectFormat
}

// getReferenceMimeType returns the mimetype String of the given reference. Ports the private
// getReferenceMimeType(DSSReference).
func (b *DataObjectFormatBuilder) getReferenceMimeType(reference *DSSReference) string {
	dataObjectFormatMimeType := enumerations.MimeType(enumerations.MimeTypeEnumBinary)
	content := reference.Contents()
	if content != nil && content.MimeType() != nil {
		dataObjectFormatMimeType = content.MimeType()
	}
	return dataObjectFormatMimeType.MimeTypeString()
}
