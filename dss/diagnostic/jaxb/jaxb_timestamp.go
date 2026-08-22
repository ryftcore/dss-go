// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. The generated classes are grouped into
// schema-area files rather than one file per class; every Java class keeps its
// name and its exact field order so that encoding/xml reproduces the JAXB element
// sequence byte for byte.

package jaxb

// XmlFoundTimestamp is the Go form of the generated JAXB class XmlFoundTimestamp
// (complexType FoundTimestamp).
type XmlFoundTimestamp struct {
	Timestamp *XmlTimestamp `xml:"Timestamp,attr,omitempty"`
}

// XmlTSAGeneralName is the Go form of the generated JAXB class XmlTSAGeneralName
// (complexType TSAGeneralName).
type XmlTSAGeneralName struct {
	Value        string `xml:",chardata"`
	ContentMatch bool   `xml:"contentMatch,attr"`
	OrderMatch   bool   `xml:"orderMatch,attr"`
}

// XmlTimestampedObject is the Go form of the generated JAXB class XmlTimestampedObject
// (complexType TimestampedObject).
type XmlTimestampedObject struct {
	Token    *XmlTokenRef                `xml:"Token,attr,omitempty"`
	Category *TimestampedObjectTypeValue `xml:"Category,attr,omitempty"`
}

// XmlTimestamp is the Go form of the generated JAXB class XmlTimestamp
// (complexType Timestamp).
type XmlTimestamp struct {
	TimestampFilename           *string                           `xml:"TimestampFilename,omitempty"`
	ArchiveTimestampType        *ArchiveTimestampTypeValue        `xml:"ArchiveTimestampType,omitempty"`
	EvidenceRecordTimestampType *EvidenceRecordTimestampTypeValue `xml:"EvidenceRecordTimestampType,omitempty"`
	ArchiveTimestampHashIndex   *XmlArchiveTimestampHashIndex     `xml:"ArchiveTimestampHashIndex,omitempty"`
	ProductionTime              *XSDateTime                       `xml:"ProductionTime,omitempty"`
	DigestMatcher               []*XmlDigestMatcher               `xml:"DigestMatcher"`
	BasicSignature              *XmlBasicSignature                `xml:"BasicSignature,omitempty"`
	SigningCertificate          *XmlSigningCertificate            `xml:"SigningCertificate,omitempty"`
	CertificateChain            *CertificateChainWrapper          `xml:"CertificateChain"`
	SignerInformationStore      *SignerInformationStoreWrapper    `xml:"SignerInformationStore"`
	TSAGeneralName              *XmlTSAGeneralName                `xml:"TSAGeneralName,omitempty"`
	PDFRevision                 *XmlPDFRevision                   `xml:"PDFRevision,omitempty"`
	FoundCertificates           *XmlFoundCertificates             `xml:"FoundCertificates,omitempty"`
	FoundRevocations            *XmlFoundRevocations              `xml:"FoundRevocations,omitempty"`
	FoundEvidenceRecords        *FoundEvidenceRecordsWrapper      `xml:"FoundEvidenceRecords"`
	TimestampedObjects          *TimestampedObjectsWrapper        `xml:"TimestampedObjects"`
	TimestampScopes             *TimestampScopesWrapper           `xml:"TimestampScopes"`
	Base64Encoded               *Base64Binary                     `xml:"Base64Encoded,omitempty"`
	DigestAlgoAndValue          *XmlDigestAlgoAndValue            `xml:"DigestAlgoAndValue,omitempty"`
	Type                        *TimestampTypeValue               `xml:"Type,attr,omitempty"`
	XmlAbstractTokenAttrs
}
