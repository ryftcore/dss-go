// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. Per the phase-8a generated-JAXB rule the
// generated classes are grouped into schema-area files rather than one file per
// class; every Java class keeps its name and its exact field order so that
// encoding/xml reproduces the JAXB element sequence byte for byte.

package jaxb

import (
	"math/big"
)

// XmlCertificateContentEquivalence is the Go form of the generated JAXB class XmlCertificateContentEquivalence
// (complexType CertificateContentEquivalence).
type XmlCertificateContentEquivalence struct {
	Uri     *string `xml:"uri,attr,omitempty"`
	Enacted bool    `xml:"enacted,attr"`
}

// XmlCertificateRef is the Go form of the generated JAXB class XmlCertificateRef
// (complexType CertificateRef).
type XmlCertificateRef struct {
	Origin             *CertificateRefOriginValue `xml:"Origin,omitempty"`
	IssuerSerial       *XmlIssuerSerial           `xml:"IssuerSerial,omitempty"`
	DigestAlgoAndValue *XmlDigestAlgoAndValue     `xml:"DigestAlgoAndValue,omitempty"`
	SerialInfo         *XmlSignerInfo             `xml:"SerialInfo,omitempty"`
	KID                *string                    `xml:"KID,omitempty"`
	X509Url            *string                    `xml:"X509Url,omitempty"`
}

// XmlCertificateRevocation is the Go form of the generated JAXB class XmlCertificateRevocation
// (complexType CertificateRevocation).
type XmlCertificateRevocation struct {
	Status         *CertificateStatusValue `xml:"Status,omitempty"`
	Reason         *RevocationReasonValue  `xml:"Reason,omitempty"`
	RevocationDate *XSDateTime             `xml:"RevocationDate,omitempty"`
	Revocation     *XmlRevocation          `xml:"Revocation,attr,omitempty"`
}

// XmlDistinguishedName is the Go form of the generated JAXB class XmlDistinguishedName
// (complexType DistinguishedName).
type XmlDistinguishedName struct {
	Value  string  `xml:",chardata"`
	Format *string `xml:"Format,attr,omitempty"`
}

// XmlX509Certificate is the Go form of the generated JAXB class XmlX509Certificate
// (complexType anonymous).
type XmlX509Certificate struct {
	Certificate *XmlCertificate `xml:"Certificate,attr,omitempty"`
}

// XmlCertificate is the Go form of the generated JAXB class XmlCertificate
// (complexType Certificate).
type XmlCertificate struct {
	SubjectDistinguishedName []*XmlDistinguishedName       `xml:"SubjectDistinguishedName"`
	IssuerDistinguishedName  []*XmlDistinguishedName       `xml:"IssuerDistinguishedName"`
	SerialNumber             *big.Int                      `xml:"SerialNumber,omitempty"`
	SubjectSerialNumber      *string                       `xml:"SubjectSerialNumber,omitempty"`
	CommonName               *string                       `xml:"CommonName,omitempty"`
	Locality                 *string                       `xml:"Locality,omitempty"`
	State                    *string                       `xml:"State,omitempty"`
	CountryName              *string                       `xml:"CountryName,omitempty"`
	OrganizationIdentifier   *string                       `xml:"OrganizationIdentifier,omitempty"`
	OrganizationName         *string                       `xml:"OrganizationName,omitempty"`
	OrganizationalUnit       *string                       `xml:"OrganizationalUnit,omitempty"`
	Title                    *string                       `xml:"Title,omitempty"`
	GivenName                *string                       `xml:"GivenName,omitempty"`
	Surname                  *string                       `xml:"Surname,omitempty"`
	Pseudonym                *string                       `xml:"Pseudonym,omitempty"`
	Email                    *string                       `xml:"Email,omitempty"`
	Sources                  *SourcesWrapper               `xml:"Sources"`
	NotAfter                 *XSDateTime                   `xml:"NotAfter,omitempty"`
	NotBefore                *XSDateTime                   `xml:"NotBefore,omitempty"`
	PublicKeySize            int                           `xml:"PublicKeySize"`
	PublicKeyEncryptionAlgo  *EncryptionAlgorithmValue     `xml:"PublicKeyEncryptionAlgo,omitempty"`
	EntityKey                *string                       `xml:"EntityKey,omitempty"`
	IssuerEntityKey          *XmlIssuerEntityKey           `xml:"IssuerEntityKey,omitempty"`
	BasicSignature           *XmlBasicSignature            `xml:"BasicSignature,omitempty"`
	SigningCertificate       *XmlSigningCertificate        `xml:"SigningCertificate,omitempty"`
	CertificateChain         *CertificateChainWrapper      `xml:"CertificateChain"`
	Trusted                  *XmlTrusted                   `xml:"Trusted,omitempty"`
	SelfSigned               bool                          `xml:"SelfSigned"`
	CertificateExtensions    *CertificateExtensionsWrapper `xml:"CertificateExtensions"`
	TrustServiceProviders    *TrustServiceProvidersWrapper `xml:"TrustServiceProviders"`
	TrustedEntities          *TrustedEntitiesWrapper       `xml:"TrustedEntities"`
	Revocations              *RevocationsWrapper           `xml:"Revocations"`
	Base64Encoded            *Base64Binary                 `xml:"Base64Encoded,omitempty"`
	DigestAlgoAndValue       *XmlDigestAlgoAndValue        `xml:"DigestAlgoAndValue,omitempty"`
	XmlAbstractTokenAttrs
}
