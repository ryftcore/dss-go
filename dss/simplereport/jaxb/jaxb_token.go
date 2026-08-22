// Ported from SimpleReport.xsd (DSS 6.5.RC1) via the JAXB classes generated
// into eu.europa.esig.dss.simplereport.jaxb. The generated classes are
// grouped into schema-area files rather than one file per class; every Java
// class keeps its name and its exact field order so that encoding/xml
// reproduces the JAXB element sequence byte for byte.
//
// Indication/BestSignatureTime/ProductionTime/POETime/SignatureFormat are
// declared without minOccurs="0" ("required") in the schema, but the RI
// does not enforce "required" on marshal/unmarshal - only XSD validation
// does, and this port's KAT runs with validateXml=false throughout, matching
// AbstractJaxbFacade's default. A null field is simply omitted like any
// optional one; testdata/oracle's sr-* fixtures (reserialised from upstream
// dss-simple-report-jaxb test resources, plus dss/simplecertificatereport's
// analogous notBefore/notAfter finding) exercise this. Every field that
// binds to a struct-kind Go type (XSDateTime) is a pointer here rather than
// a plain value, since encoding/xml's omitempty cannot recognise a
// struct-kind zero value as empty; the string-kind enum adapter fields
// (Indication, SignatureFormat) only need the omitempty tag added, since
// their zero value ("") never collides with a real lexical form.

package jaxb

// XmlTokenContent carries the element content of the abstract Token base
// type. JAXB emits base content before extension content, so derived types
// embed it first.
type XmlTokenContent struct {
	Filename              *string              `xml:"Filename,omitempty"`
	CertificateChain      *XmlCertificateChain `xml:"CertificateChain,omitempty"`
	Indication            IndicationValue      `xml:"Indication,omitempty"`
	SubIndication         *SubIndicationValue  `xml:"SubIndication,omitempty"`
	AdESValidationDetails *XmlDetails          `xml:"AdESValidationDetails,omitempty"`
	QualificationDetails  *XmlDetails          `xml:"QualificationDetails,omitempty"`
	ExtensionPeriodMin    *XSDateTime          `xml:"ExtensionPeriodMin,omitempty"`
	ExtensionPeriodMax    *XSDateTime          `xml:"ExtensionPeriodMax,omitempty"`
}

// XmlTokenAttrs carries the attributes of the abstract Token base type. The
// JAXB RI writes extension attributes before base attributes, so derived
// types embed it last.
type XmlTokenAttrs struct {
	Id string `xml:"Id,attr"`
}

// XmlTokenItem is implemented by every alternative of the choice group
// SimpleReport.xsd binds to Signature/Timestamp/EvidenceRecord/EAA: the Go
// stand-in for List<XmlToken> in
// XmlSimpleReport.SignatureOrTimestampOrEvidenceRecord (Java's XmlToken is
// the abstract Token base class; a like-named Go interface would collide
// with the XmlToken... struct types above, so the choice-group marker keeps
// the "Item" suffix instead). Package callers type-switch on the concrete
// *XmlSignature / *XmlTimestamp / *XmlEvidenceRecord / *XmlEAA the way
// SimpleReport.java's instanceof checks do.
type XmlTokenItem interface {
	xmlTokenElementName() string
}

func (*XmlSignature) xmlTokenElementName() string      { return "Signature" }
func (*XmlTimestamp) xmlTokenElementName() string      { return "Timestamp" }
func (*XmlEvidenceRecord) xmlTokenElementName() string { return "EvidenceRecord" }
func (*XmlEAA) xmlTokenElementName() string            { return "EAA" }

// tokenItemFor builds a zero-valued alternative for the given element name,
// or nil if the name is not one of the choice group's four alternatives.
func tokenItemFor(name string) XmlTokenItem {
	switch name {
	case "Signature":
		return &XmlSignature{}
	case "Timestamp":
		return &XmlTimestamp{}
	case "EvidenceRecord":
		return &XmlEvidenceRecord{}
	case "EAA":
		return &XmlEAA{}
	default:
		return nil
	}
}

// XmlSignature is the Go form of the generated JAXB class XmlSignature
// (complexType Signature, extension of Token).
type XmlSignature struct {
	XmlTokenContent
	SigningTime       *XSDateTime          `xml:"SigningTime,omitempty"`
	BestSignatureTime *XSDateTime          `xml:"BestSignatureTime,omitempty"`
	SignedBy          *string              `xml:"SignedBy,omitempty"`
	SignatureLevel    *XmlSignatureLevel   `xml:"SignatureLevel,omitempty"`
	SignatureScope    []*XmlSignatureScope `xml:"SignatureScope,omitempty"`
	Timestamps        *XmlTimestamps       `xml:"Timestamps,omitempty"`
	EvidenceRecords   *XmlEvidenceRecords  `xml:"EvidenceRecords,omitempty"`
	CounterSignature  *bool                `xml:"CounterSignature,attr,omitempty"`
	ParentId          *string              `xml:"ParentId,attr,omitempty"`
	SignatureFormat   SignatureLevelValue  `xml:"SignatureFormat,attr,omitempty"`
	XmlTokenAttrs
}

// XmlTimestamp is the Go form of the generated JAXB class XmlTimestamp
// (complexType Timestamp, extension of Token).
type XmlTimestamp struct {
	XmlTokenContent
	ProductionTime  *XSDateTime          `xml:"ProductionTime,omitempty"`
	ProducedBy      *string              `xml:"ProducedBy,omitempty"`
	TimestampLevel  *XmlTimestampLevel   `xml:"TimestampLevel,omitempty"`
	TimestampScope  []*XmlSignatureScope `xml:"TimestampScope,omitempty"`
	EvidenceRecords *XmlEvidenceRecords  `xml:"EvidenceRecords,omitempty"`
	XmlTokenAttrs
}

// XmlEvidenceRecord is the Go form of the generated JAXB class
// XmlEvidenceRecord (complexType EvidenceRecord, extension of Token).
type XmlEvidenceRecord struct {
	XmlTokenContent
	POETime             *XSDateTime          `xml:"POETime,omitempty"`
	EvidenceRecordScope []*XmlSignatureScope `xml:"EvidenceRecordScope,omitempty"`
	Timestamps          *XmlTimestamps       `xml:"Timestamps,omitempty"`
	Embedded            *bool                `xml:"Embedded,attr,omitempty"`
	ParentId            *string              `xml:"ParentId,attr,omitempty"`
	XmlTokenAttrs
}

// XmlEAA is the Go form of the generated JAXB class XmlEAA (complexType EAA,
// extension of Token).
type XmlEAA struct {
	XmlTokenContent
	EAALevel            []*XmlEAALevel  `xml:"EAALevel,omitempty"`
	EAASignature        []*XmlSignature `xml:"EAASignature,omitempty"`
	KeyBindingSignature *XmlSignature   `xml:"KeyBindingSignature,omitempty"`
	EAAPayload          *XmlEAAPayload  `xml:"EAAPayload,omitempty"`
	XmlTokenAttrs
}

// XmlTimestamps is the Go form of the generated JAXB class XmlTimestamps
// (complexType Timestamps).
type XmlTimestamps struct {
	Timestamp []*XmlTimestamp `xml:"Timestamp,omitempty"`
}

// XmlEvidenceRecords is the Go form of the generated JAXB class
// XmlEvidenceRecords (complexType EvidenceRecords).
type XmlEvidenceRecords struct {
	EvidenceRecord []*XmlEvidenceRecord `xml:"EvidenceRecord,omitempty"`
}
