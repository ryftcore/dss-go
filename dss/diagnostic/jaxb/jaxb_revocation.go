// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. The generated classes are grouped into
// schema-area files rather than one file per class; every Java class keeps its
// name and its exact field order so that encoding/xml reproduces the JAXB element
// sequence byte for byte.

package jaxb

import ()

// XmlFoundRevocationContent carries the element content of the XmlFoundRevocation base type. JAXB emits base
// content before extension content, so derived types embed it first.
type XmlFoundRevocationContent struct {
	Type          *RevocationTypeValue    `xml:"Type,omitempty"`
	Origin        []RevocationOriginValue `xml:"Origin"`
	RevocationRef []*XmlRevocationRef     `xml:"RevocationRef"`
}

// XmlFoundRevocation is the Go form of the generated JAXB class XmlFoundRevocation
// (complexType FoundRevocation).
type XmlFoundRevocation struct {
	XmlFoundRevocationContent
}

// XmlFoundRevocations is the Go form of the generated JAXB class XmlFoundRevocations
// (complexType FoundRevocations).
type XmlFoundRevocations struct {
	RelatedRevocation []*XmlRelatedRevocation `xml:"RelatedRevocation"`
	OrphanRevocation  []*XmlOrphanRevocation  `xml:"OrphanRevocation"`
}

// XmlRevocationRef is the Go form of the generated JAXB class XmlRevocationRef
// (complexType RevocationRef).
type XmlRevocationRef struct {
	Origin             []RevocationRefOriginValue `xml:"Origin"`
	DigestAlgoAndValue *XmlDigestAlgoAndValue     `xml:"DigestAlgoAndValue,omitempty"`
	ProducedAt         *XSDateTime                `xml:"ProducedAt,omitempty"`
	ResponderId        *XmlSignerInfo             `xml:"ResponderId,omitempty"`
	Issuer             *string                    `xml:"Issuer,omitempty"`
	IssueTime          *XSDateTime                `xml:"IssueTime,omitempty"`
	CRLNumber          *BigInteger                `xml:"CRLNumber,omitempty"`
	Uri                *string                    `xml:"Uri,omitempty"`
}

// XmlOrphanRevocation is the Go form of the generated JAXB class XmlOrphanRevocation
// (complexType OrphanRevocation).
type XmlOrphanRevocation struct {
	XmlFoundRevocationContent
	Token *XmlOrphanRevocationToken `xml:"Token,attr,omitempty"`
}

// XmlRelatedRevocation is the Go form of the generated JAXB class XmlRelatedRevocation
// (complexType RelatedRevocation).
type XmlRelatedRevocation struct {
	XmlFoundRevocationContent
	Revocation *XmlRevocation `xml:"Revocation,attr,omitempty"`
}

// XmlRevocation is the Go form of the generated JAXB class XmlRevocation
// (complexType Revocation).
type XmlRevocation struct {
	Origin                   *RevocationOriginValue   `xml:"Origin,omitempty"`
	Type                     *RevocationTypeValue     `xml:"Type,omitempty"`
	SourceAddress            *string                  `xml:"SourceAddress,omitempty"`
	ProductionDate           *XSDateTime              `xml:"ProductionDate,omitempty"`
	ThisUpdate               *XSDateTime              `xml:"ThisUpdate,omitempty"`
	NextUpdate               *XSDateTime              `xml:"NextUpdate,omitempty"`
	CRLNumber                *BigInteger              `xml:"CRLNumber,omitempty"`
	ExpiredCertsOnCRL        *XSDateTime              `xml:"ExpiredCertsOnCRL,omitempty"`
	ArchiveCutOff            *XSDateTime              `xml:"ArchiveCutOff,omitempty"`
	CertHashExtensionPresent *bool                    `xml:"CertHashExtensionPresent,omitempty"`
	CertHashExtensionMatch   *bool                    `xml:"CertHashExtensionMatch,omitempty"`
	BasicSignature           *XmlBasicSignature       `xml:"BasicSignature,omitempty"`
	SigningCertificate       *XmlSigningCertificate   `xml:"SigningCertificate,omitempty"`
	CertificateChain         *CertificateChainWrapper `xml:"CertificateChain"`
	FoundCertificates        *XmlFoundCertificates    `xml:"FoundCertificates,omitempty"`
	Base64Encoded            *Base64Binary            `xml:"Base64Encoded,omitempty"`
	DigestAlgoAndValue       *XmlDigestAlgoAndValue   `xml:"DigestAlgoAndValue,omitempty"`
	XmlAbstractTokenAttrs
}
