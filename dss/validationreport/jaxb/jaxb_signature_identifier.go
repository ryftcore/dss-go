// Ported from the generated JAXB classes:
//   - SignatureIdentifierType.java
//   - SignatureValidatorType.java
//   - SignatureReferenceType.java
//   - XAdESSignaturePtrType.java
//
// (specs-validation-report, DSS 6.5.RC1).
package jaxb

// SignatureIdentifierType is the Go form of the generated JAXB class
// SignatureIdentifierType (complexType SignatureIdentifierType).
type SignatureIdentifierType struct {
	DigestAlgAndValue *DigestAlgAndValueType `xml:"DigestAlgAndValue,omitempty"`
	SignatureValue    *SignatureValueType    `xml:"http://www.w3.org/2000/09/xmldsig# SignatureValue,omitempty"`
	HashOnly          bool                   `xml:"HashOnly"`
	DocHashOnly       bool                   `xml:"DocHashOnly"`
	DAIdentifier      *string                `xml:"DAIdentifier,omitempty"`
	OtherInformation  *AnyType               `xml:"OtherInformation,omitempty"`
	Id                string                 `xml:"id,attr,omitempty"`
}

// SignatureValidatorType is the Go form of the generated JAXB class
// SignatureValidatorType (complexType SignatureValidatorType).
type SignatureValidatorType struct {
	DigitalId        []*DigitalIdentityType `xml:"DigitalId"`
	TSPInformation   *TSPInformationType    `xml:"TSPInformation,omitempty"`
	OtherInformation *AnyType               `xml:"OtherInformation,omitempty"`
}

// SignatureReferenceType is the Go form of the generated JAXB class
// SignatureReferenceType (complexType SignatureReferenceType, an
// xs:choice). Like the Java class, this is not wrapped in a JAXBElement
// choice group - the alternatives use disjoint element names, so xjc bound
// each to a plain optional field and leaves selecting exactly one
// alternative to the caller, exactly as upstream does.
type SignatureReferenceType struct {
	CanonicalizationMethod *string                `xml:"CanonicalizationMethod,omitempty"`
	DigestMethod           *string                `xml:"DigestMethod,omitempty"`
	DigestValue            Base64Binary           `xml:"DigestValue,omitempty"`
	XAdESSignaturePtr      *XAdESSignaturePtrType `xml:"XAdESSignaturePtr,omitempty"`
	PAdESFieldName         *string                `xml:"PAdESFieldName,omitempty"`
	OtherInformation       *AnyType               `xml:"OtherInformation,omitempty"`
}

// XAdESSignaturePtrType is the Go form of the generated JAXB class
// XAdESSignaturePtrType (complexType XAdESSignaturePtrType). WhichDocument
// and SchemaRefs are round-tripped as plain IDREF/IDREFS strings - see
// xml.go's header.
type XAdESSignaturePtrType struct {
	NsPrefixMapping []*NsPrefixMappingType `xml:"NsPrefixMapping,omitempty"`
	WhichDocument   string                 `xml:"WhichDocument,attr,omitempty"`
	XPath           string                 `xml:"XPath,attr,omitempty"`
	SchemaRefs      IDREFS                 `xml:"SchemaRefs,attr"`
}
