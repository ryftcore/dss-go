// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. The generated classes are grouped into
// schema-area files rather than one file per class; every Java class keeps its
// name and its exact field order so that encoding/xml reproduces the JAXB element
// sequence byte for byte.

package jaxb

import (
	"encoding/xml"
)

// XmlConnectionInfo is the Go form of the generated JAXB class XmlConnectionInfo
// (complexType ConnectionInfo).
type XmlConnectionInfo struct {
	TLSCertificate                 *XmlTLSCertificate                 `xml:"TLSCertificate,omitempty"`
	TLSCertificateBindingUrl       *string                            `xml:"TLSCertificateBindingUrl,omitempty"`
	TLSCertificateBindingSignature *XmlTLSCertificateBindingSignature `xml:"TLSCertificateBindingSignature,omitempty"`
	Url                            *string                            `xml:"url,attr,omitempty"`
}

// XmlContainerInfo is the Go form of the generated JAXB class XmlContainerInfo
// (complexType ContainerInfo).
type XmlContainerInfo struct {
	ContainerType       *ASiCContainerTypeValue `xml:"ContainerType,omitempty"`
	ZipComment          *string                 `xml:"ZipComment,omitempty"`
	MimeTypeFilePresent *bool                   `xml:"MimeTypeFilePresent,omitempty"`
	MimeTypeContent     *string                 `xml:"MimeTypeContent,omitempty"`
	ManifestFiles       *ManifestFilesWrapper   `xml:"ManifestFiles"`
	ContentFiles        *ContentFilesWrapper    `xml:"ContentFiles"`
}

// XmlDiagnosticData is the Go form of the generated JAXB class XmlDiagnosticData
// (complexType Data).
type XmlDiagnosticData struct {
	XMLName                 xml.Name                        `xml:"http://dss.esig.europa.eu/validation/diagnostic DiagnosticData"`
	DocumentName            *string                         `xml:"DocumentName,omitempty"`
	ValidationDate          *XSDateTime                     `xml:"ValidationDate,omitempty"`
	ContainerInfo           *XmlContainerInfo               `xml:"ContainerInfo,omitempty"`
	PDFAInfo                *XmlPDFAInfo                    `xml:"PDFAInfo,omitempty"`
	ConnectionInfo          *XmlConnectionInfo              `xml:"ConnectionInfo,omitempty"`
	EAAPresentationInfo     *XmlEAAPresentationInfo         `xml:"EAAPresentationInfo,omitempty"`
	EAAs                    *EAAsWrapper                    `xml:"EAAs"`
	Signatures              *SignaturesWrapper              `xml:"Signatures"`
	EvidenceRecords         *EvidenceRecordsWrapper         `xml:"EvidenceRecords"`
	UsedEAARevocationTokens *UsedEAARevocationTokensWrapper `xml:"UsedEAARevocationTokens"`
	UsedCertificates        *UsedCertificatesWrapper        `xml:"UsedCertificates"`
	UsedRevocations         *UsedRevocationsWrapper         `xml:"UsedRevocations"`
	UsedTimestamps          *UsedTimestampsWrapper          `xml:"UsedTimestamps"`
	OrphanTokens            *XmlOrphanTokens                `xml:"OrphanTokens,omitempty"`
	OriginalDocuments       *OriginalDocumentsWrapper       `xml:"OriginalDocuments"`
	TrustedLists            *TrustedListsWrapper            `xml:"TrustedLists"`
	ListsOfTrustedEntities  *ListsOfTrustedEntitiesWrapper  `xml:"ListsOfTrustedEntities"`
}

// XmlManifestFile is the Go form of the generated JAXB class XmlManifestFile
// (complexType ManifestFile).
type XmlManifestFile struct {
	Filename          *string         `xml:"Filename,omitempty"`
	SignatureFilename *string         `xml:"SignatureFilename,omitempty"`
	Entries           *EntriesWrapper `xml:"Entries"`
}

// XmlOrphanTokens is the Go form of the generated JAXB class XmlOrphanTokens
// (complexType anonymous).
type XmlOrphanTokens struct {
	OrphanCertificate []*XmlOrphanCertificateToken `xml:"OrphanCertificate"`
	OrphanRevocation  []*XmlOrphanRevocationToken  `xml:"OrphanRevocation"`
}

// XmlPDFAInfo is the Go form of the generated JAXB class XmlPDFAInfo
// (complexType PDFAInfo).
type XmlPDFAInfo struct {
	ProfileId          *string                    `xml:"ProfileId,omitempty"`
	ValidationMessages *ValidationMessagesWrapper `xml:"ValidationMessages"`
	Compliant          bool                       `xml:"compliant,attr"`
}

// XmlTLSCertificate is the Go form of the generated JAXB class XmlTLSCertificate
// (complexType anonymous).
type XmlTLSCertificate struct {
	Certificate *XmlCertificate `xml:"Certificate,attr,omitempty"`
}

// XmlTLSCertificateBindingSignature is the Go form of the generated JAXB class XmlTLSCertificateBindingSignature
// (complexType anonymous).
type XmlTLSCertificateBindingSignature struct {
	Signature *XmlSignature `xml:"Signature,attr,omitempty"`
}
