// Ported from the generated JAXB classes:
//   - ValidationObjectType.java
//   - ValidationObjectListType.java
//   - ValidationObjectRepresentationType.java
//   - SignersDocumentType.java
//
// (specs-validation-report, DSS 6.5.RC1).
//
// ValidationObjectRepresentationType and SignersDocumentType are both
// xs:choice complexTypes marshalled by hand via the choiceItem machinery
// jaxb_sa_attributes.go defines; see that file's header for the technique.
package jaxb

import (
	"encoding/xml"
	"reflect"
)

// ValidationObjectListType is the Go form of the generated JAXB class
// ValidationObjectListType (complexType ValidationObjectListType).
type ValidationObjectListType struct {
	ValidationObject []*ValidationObjectType `xml:"ValidationObject"`
}

// ValidationObjectType is the Go form of the generated JAXB class
// ValidationObjectType (complexType ValidationObjectType). Id is
// round-tripped as a plain xs:ID string - see xml.go's header.
type ValidationObjectType struct {
	ObjectType                     ObjectType                         `xml:"ObjectType"`
	ValidationObjectRepresentation ValidationObjectRepresentationType `xml:"ValidationObjectRepresentation"`
	POE                            *POEType                           `xml:"POE,omitempty"`
	POEProvisioning                *POEProvisioningType               `xml:"POEProvisioning,omitempty"`
	ValidationReport               *SignatureValidationReportType     `xml:"ValidationReport,omitempty"`
	Id                             string                             `xml:"id,attr"`
}

// validationObjectRepresentationElements is the name->type table of
// ValidationObjectRepresentationType's 4-member choice, in schema order.
var validationObjectRepresentationElements = []choiceElement{
	{"direct", reflect.TypeOf(RawContent{})},
	{"base64", reflect.TypeOf(Base64Binary{})},
	{"DigestAlgAndValue", reflect.TypeOf(DigestAlgAndValueType{})},
	{"URI", reflect.TypeOf("")},
}

// ValidationObjectRepresentationType is the Go form of the generated JAXB
// class ValidationObjectRepresentationType (complexType
// ValidationObjectRepresentationType).
type ValidationObjectRepresentationType struct {
	Items []choiceItem
}

// MarshalXML writes each representation item under the element name it was
// recorded with.
func (t *ValidationObjectRepresentationType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if err := encodeChoiceItems(e, t.Items); err != nil {
		return err
	}
	return e.EncodeToken(start.End())
}

// UnmarshalXML builds each representation item from the type its element
// name is bound to.
func (t *ValidationObjectRepresentationType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	items, err := decodeChoiceItems(d, validationObjectRepresentationElements)
	if err != nil {
		return err
	}
	t.Items = items
	return nil
}

// signersDocumentElements is the name->type table of SignersDocumentType's
// 3-member choice, in schema order. SignersDocumentRepresentation appears at
// two positions of the schema's nested choice/sequence (the javadoc of the
// upstream generated class notes the resulting duplicate accessor name) but
// both bind to the same type, so the table has one unambiguous entry for it.
var signersDocumentElements = []choiceElement{
	{"DigestAlgAndValue", reflect.TypeOf(DigestAlgAndValueType{})},
	{"SignersDocumentRepresentation", reflect.TypeOf(VOReferenceType{})},
	{"SignersDocumentRef", reflect.TypeOf(VOReferenceType{})},
}

// SignersDocumentType is the Go form of the generated JAXB class
// SignersDocumentType (complexType SignersDocumentType).
type SignersDocumentType struct {
	Items []choiceItem
}

// MarshalXML writes each content item under the element name it was
// recorded with.
func (t *SignersDocumentType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if err := encodeChoiceItems(e, t.Items); err != nil {
		return err
	}
	return e.EncodeToken(start.End())
}

// UnmarshalXML builds each content item from the type its element name is
// bound to.
func (t *SignersDocumentType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	items, err := decodeChoiceItems(d, signersDocumentElements)
	if err != nil {
		return err
	}
	t.Items = items
	return nil
}
