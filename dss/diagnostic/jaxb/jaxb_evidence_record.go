// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. Per the phase-8a generated-JAXB rule the
// generated classes are grouped into schema-area files rather than one file per
// class; every Java class keeps its name and its exact field order so that
// encoding/xml reproduces the JAXB element sequence byte for byte.

package jaxb

// XmlFoundEvidenceRecord is the Go form of the generated JAXB class XmlFoundEvidenceRecord
// (complexType FoundEvidenceRecord).
type XmlFoundEvidenceRecord struct {
	EvidenceRecord *XmlEvidenceRecord `xml:"EvidenceRecord,attr,omitempty"`
}

// XmlEvidenceRecord is the Go form of the generated JAXB class XmlEvidenceRecord
// (complexType EvidenceRecord).
type XmlEvidenceRecord struct {
	DocumentName             *string                               `xml:"DocumentName,omitempty"`
	Type                     *EvidenceRecordTypeEnumValue          `xml:"Type,omitempty"`
	Origin                   *EvidenceRecordOriginValue            `xml:"Origin,omitempty"`
	IncorporationType        *EvidenceRecordIncorporationTypeValue `xml:"IncorporationType,omitempty"`
	StructuralValidation     *XmlStructuralValidation              `xml:"StructuralValidation,omitempty"`
	DigestMatchers           *DigestMatchersWrapper                `xml:"DigestMatchers"`
	EvidenceRecordTimestamps *EvidenceRecordTimestampsWrapper      `xml:"EvidenceRecordTimestamps"`
	FoundCertificates        *XmlFoundCertificates                 `xml:"FoundCertificates,omitempty"`
	FoundRevocations         *XmlFoundRevocations                  `xml:"FoundRevocations,omitempty"`
	TimestampedObjects       *TimestampedObjectsWrapper            `xml:"TimestampedObjects"`
	EvidenceRecordScopes     *EvidenceRecordScopesWrapper          `xml:"EvidenceRecordScopes"`
	Base64Encoded            *Base64Binary                         `xml:"Base64Encoded,omitempty"`
	DigestAlgoAndValue       *XmlDigestAlgoAndValue                `xml:"DigestAlgoAndValue,omitempty"`
	Embedded                 *bool                                 `xml:"Embedded,attr,omitempty"`
	Parent                   *XmlSignature                         `xml:"Parent,attr,omitempty"`
	XmlAbstractTokenAttrs
}
