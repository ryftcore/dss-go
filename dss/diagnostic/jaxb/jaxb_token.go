// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. The generated classes are grouped into
// schema-area files rather than one file per class; every Java class keeps its
// name and its exact field order so that encoding/xml reproduces the JAXB element
// sequence byte for byte.

package jaxb

import ()

// XmlAbstractTokenAttrs carries the attributes of the XmlAbstractToken base type. The JAXB RI writes
// extension attributes before base attributes, so derived types embed it last.
type XmlAbstractTokenAttrs struct {
	Id         *CollapsedString `xml:"Id,attr,omitempty"`
	Duplicated *bool            `xml:"Duplicated,attr,omitempty"`
}

// XmlAbstractToken is the Go form of the generated JAXB class XmlAbstractToken
// (complexType AbstractToken).
type XmlAbstractToken struct {
	XmlAbstractTokenAttrs
}

// XmlBasicSignature is the Go form of the generated JAXB class XmlBasicSignature
// (complexType BasicSignature).
type XmlBasicSignature struct {
	EncryptionAlgoUsedToSignThisToken *EncryptionAlgorithmValue `xml:"EncryptionAlgoUsedToSignThisToken,omitempty"`
	KeyLengthUsedToSignThisToken      *string                   `xml:"KeyLengthUsedToSignThisToken,omitempty"`
	DigestAlgoUsedToSignThisToken     *DigestAlgorithmValue     `xml:"DigestAlgoUsedToSignThisToken,omitempty"`
	SignatureIntact                   *bool                     `xml:"SignatureIntact,omitempty"`
	SignatureValid                    *bool                     `xml:"SignatureValid,omitempty"`
}

// XmlChainItem is the Go form of the generated JAXB class XmlChainItem
// (complexType anonymous).
type XmlChainItem struct {
	Certificate *XmlCertificate `xml:"Certificate,attr,omitempty"`
}

// XmlIssuerEntityKey is the Go form of the generated JAXB class XmlIssuerEntityKey
// (complexType IssuerEntityKey).
type XmlIssuerEntityKey struct {
	Value       string `xml:",chardata"`
	Key         bool   `xml:"key,attr"`
	SubjectName bool   `xml:"subjectName,attr"`
}

// XmlIssuerSerial is the Go form of the generated JAXB class XmlIssuerSerial
// (complexType IssuerSerial).
type XmlIssuerSerial struct {
	Value Base64Binary `xml:",chardata"`
	Match *bool        `xml:"match,attr,omitempty"`
}

// XmlSignerInfo is the Go form of the generated JAXB class XmlSignerInfo
// (complexType SignerInfo).
type XmlSignerInfo struct {
	IssuerName   *string       `xml:"IssuerName,omitempty"`
	SerialNumber *BigInteger   `xml:"SerialNumber,omitempty"`
	Ski          *Base64Binary `xml:"Ski,omitempty"`
	Current      *bool         `xml:"Current,attr,omitempty"`
}

// XmlSigningCertificate is the Go form of the generated JAXB class XmlSigningCertificate
// (complexType SigningCertificate).
type XmlSigningCertificate struct {
	PublicKey   *Base64Binary   `xml:"PublicKey,omitempty"`
	Certificate *XmlCertificate `xml:"Certificate,attr,omitempty"`
}

// XmlStructuralValidationContent carries the element content of the XmlStructuralValidation base type. JAXB emits base
// content before extension content, so derived types embed it first.
type XmlStructuralValidationContent struct {
	Message []string `xml:"Message"`
}

// XmlStructuralValidationAttrs carries the attributes of the XmlStructuralValidation base type. The JAXB RI writes
// extension attributes before base attributes, so derived types embed it last.
type XmlStructuralValidationAttrs struct {
	Valid bool `xml:"valid,attr"`
}

// XmlStructuralValidation is the Go form of the generated JAXB class XmlStructuralValidation
// (complexType StructuralValidation).
type XmlStructuralValidation struct {
	XmlStructuralValidationContent
	XmlStructuralValidationAttrs
}

// XmlTrusted is the Go form of the generated JAXB class XmlTrusted
// (complexType Trusted).
type XmlTrusted struct {
	Value      bool        `xml:",chardata"`
	StartDate  *XSDateTime `xml:"startDate,attr,omitempty"`
	SunsetDate *XSDateTime `xml:"sunsetDate,attr,omitempty"`
}

// XmlArchiveTimestampHashIndex is the Go form of the generated JAXB class XmlArchiveTimestampHashIndex
// (complexType ArchiveTimestampHashIndex).
type XmlArchiveTimestampHashIndex struct {
	XmlStructuralValidationContent
	Version *ArchiveTimestampHashIndexVersionValue `xml:"version,attr,omitempty"`
	XmlStructuralValidationAttrs
}

// XmlOrphanTokenAttrs carries the attributes of the XmlOrphanToken base type. The JAXB RI writes
// extension attributes before base attributes, so derived types embed it last.
type XmlOrphanTokenAttrs struct {
	EncapsulationType *XmlEncapsulationType `xml:"EncapsulationType,attr,omitempty"`
}

// XmlOrphanToken is the Go form of the generated JAXB class XmlOrphanToken
// (complexType OrphanToken).
type XmlOrphanToken struct {
	XmlOrphanTokenAttrs
	XmlAbstractTokenAttrs
}

// XmlSignerData is the Go form of the generated JAXB class XmlSignerData
// (complexType SignerData).
type XmlSignerData struct {
	ReferencedName     *string                `xml:"ReferencedName,omitempty"`
	DigestAlgoAndValue *XmlDigestAlgoAndValue `xml:"DigestAlgoAndValue,omitempty"`
	Parent             *XmlSignerData         `xml:"Parent,attr,omitempty"`
	XmlAbstractTokenAttrs
}

// XmlOrphanCertificateToken is the Go form of the generated JAXB class XmlOrphanCertificateToken
// (complexType OrphanCertificateToken).
type XmlOrphanCertificateToken struct {
	SubjectDistinguishedName []*XmlDistinguishedName `xml:"SubjectDistinguishedName"`
	IssuerDistinguishedName  []*XmlDistinguishedName `xml:"IssuerDistinguishedName"`
	SerialNumber             *BigInteger             `xml:"SerialNumber,omitempty"`
	NotAfter                 *XSDateTime             `xml:"NotAfter,omitempty"`
	NotBefore                *XSDateTime             `xml:"NotBefore,omitempty"`
	EntityKey                *string                 `xml:"EntityKey,omitempty"`
	Trusted                  *bool                   `xml:"Trusted,omitempty"`
	SelfSigned               *bool                   `xml:"SelfSigned,omitempty"`
	Base64Encoded            *Base64Binary           `xml:"Base64Encoded,omitempty"`
	DigestAlgoAndValue       *XmlDigestAlgoAndValue  `xml:"DigestAlgoAndValue,omitempty"`
	XmlOrphanTokenAttrs
	XmlAbstractTokenAttrs
}

// XmlOrphanRevocationToken is the Go form of the generated JAXB class XmlOrphanRevocationToken
// (complexType OrphanRevocationToken).
type XmlOrphanRevocationToken struct {
	RevocationType     *RevocationTypeValue   `xml:"RevocationType,omitempty"`
	Base64Encoded      *Base64Binary          `xml:"Base64Encoded,omitempty"`
	DigestAlgoAndValue *XmlDigestAlgoAndValue `xml:"DigestAlgoAndValue,omitempty"`
	XmlOrphanTokenAttrs
	XmlAbstractTokenAttrs
}
