// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. Per the phase-8a generated-JAXB rule the
// generated classes are grouped into schema-area files rather than one file per
// class; every Java class keeps its name and its exact field order so that
// encoding/xml reproduces the JAXB element sequence byte for byte.

package jaxb

import ()

// XmlByteRange is the Go form of the generated JAXB class XmlByteRange
// (complexType ByteRange).
type XmlByteRange struct {
	Value BigIntegerList `xml:",chardata"`
	Valid bool           `xml:"valid,attr"`
}

// XmlDocMDP is the Go form of the generated JAXB class XmlDocMDP
// (complexType anonymous).
type XmlDocMDP struct {
	Permissions *CertificationPermissionValue `xml:"Permissions,omitempty"`
}

// XmlModification is the Go form of the generated JAXB class XmlModification
// (complexType Modification).
type XmlModification struct {
	Page *BigInteger `xml:"Page,attr,omitempty"`
}

// XmlModificationDetection is the Go form of the generated JAXB class XmlModificationDetection
// (complexType ModificationDetection).
type XmlModificationDetection struct {
	AnnotationOverlap   []*XmlModification      `xml:"AnnotationOverlap"`
	VisualDifference    []*XmlModification      `xml:"VisualDifference"`
	PageDifference      []*XmlModification      `xml:"PageDifference"`
	ObjectModifications *XmlObjectModifications `xml:"ObjectModifications,omitempty"`
}

// XmlObjectModification is the Go form of the generated JAXB class XmlObjectModification
// (complexType ObjectModification).
type XmlObjectModification struct {
	Value     string                          `xml:",chardata"`
	Action    *PdfObjectModificationTypeValue `xml:"action,attr,omitempty"`
	FieldName *string                         `xml:"fieldName,attr,omitempty"`
	Type      *string                         `xml:"type,attr,omitempty"`
}

// XmlObjectModifications is the Go form of the generated JAXB class XmlObjectModifications
// (complexType ObjectModifications).
type XmlObjectModifications struct {
	ExtensionChange     []*XmlObjectModification `xml:"ExtensionChange"`
	SignatureOrFormFill []*XmlObjectModification `xml:"SignatureOrFormFill"`
	AnnotationChange    []*XmlObjectModification `xml:"AnnotationChange"`
	Undefined           []*XmlObjectModification `xml:"Undefined"`
}

// XmlPDFLockDictionary is the Go form of the generated JAXB class XmlPDFLockDictionary
// (complexType PDFLockDictionary).
type XmlPDFLockDictionary struct {
	Action      *PdfLockActionValue           `xml:"Action,omitempty"`
	Field       []string                      `xml:"Field"`
	Permissions *CertificationPermissionValue `xml:"Permissions,omitempty"`
}

// XmlPDFRevision is the Go form of the generated JAXB class XmlPDFRevision
// (complexType PDFRevision).
type XmlPDFRevision struct {
	SignatureField         []*XmlPDFSignatureField    `xml:"SignatureField"`
	PDFSignatureDictionary *XmlPDFSignatureDictionary `xml:"PDFSignatureDictionary,omitempty"`
	ModificationDetection  *XmlModificationDetection  `xml:"ModificationDetection,omitempty"`
}

// XmlPDFSignatureDictionary is the Go form of the generated JAXB class XmlPDFSignatureDictionary
// (complexType PDFSignatureDictionary).
type XmlPDFSignatureDictionary struct {
	SignerName         *string               `xml:"SignerName,omitempty"`
	Type               *string               `xml:"Type,omitempty"`
	Filter             *string               `xml:"Filter,omitempty"`
	SubFilter          *string               `xml:"SubFilter,omitempty"`
	ContactInfo        *string               `xml:"ContactInfo,omitempty"`
	Location           *string               `xml:"Location,omitempty"`
	Reason             *string               `xml:"Reason,omitempty"`
	SignatureByteRange *XmlByteRange         `xml:"SignatureByteRange,omitempty"`
	DocMDP             *XmlDocMDP            `xml:"DocMDP,omitempty"`
	FieldMDP           *XmlPDFLockDictionary `xml:"FieldMDP,omitempty"`
	Consistent         bool                  `xml:"consistent,attr"`
}

// XmlPDFSignatureField is the Go form of the generated JAXB class XmlPDFSignatureField
// (complexType PDFSignatureField).
type XmlPDFSignatureField struct {
	SigFieldLock *XmlPDFLockDictionary `xml:"SigFieldLock,omitempty"`
	Name         *string               `xml:"name,attr,omitempty"`
}
