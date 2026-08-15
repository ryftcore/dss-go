// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. Per the phase-8a generated-JAXB rule the
// generated classes are grouped into schema-area files rather than one file per
// class; every Java class keeps its name and its exact field order so that
// encoding/xml reproduces the JAXB element sequence byte for byte.

package jaxb

import (
	"math/big"
)

// XmlCertificateExtensionContent carries the element content of the XmlCertificateExtension base type. JAXB emits base
// content before extension content, so derived types embed it first.
type XmlCertificateExtensionContent struct {
	Octets *Base64Binary `xml:"octets,omitempty"`
}

// XmlCertificateExtensionAttrs carries the attributes of the XmlCertificateExtension base type. The JAXB RI writes
// extension attributes before base attributes, so derived types embed it last.
type XmlCertificateExtensionAttrs struct {
	OID         *string `xml:"OID,attr,omitempty"`
	Description *string `xml:"description,attr,omitempty"`
	Critical    *bool   `xml:"critical,attr,omitempty"`
}

// XmlCertificateExtension is the Go form of the generated JAXB class XmlCertificateExtension
// (complexType CertificateExtension).
type XmlCertificateExtension struct {
	XmlCertificateExtensionContent
	XmlCertificateExtensionAttrs
}

// XmlGeneralNameContent carries the element content of the XmlGeneralName base type. JAXB emits base
// content before extension content, so derived types embed it first.
type XmlGeneralNameContent struct {
	Value string `xml:",chardata"`
}

// XmlGeneralNameAttrs carries the attributes of the XmlGeneralName base type. The JAXB RI writes
// extension attributes before base attributes, so derived types embed it last.
type XmlGeneralNameAttrs struct {
	Type *GeneralNameTypeValue `xml:"type,attr,omitempty"`
}

// XmlGeneralName is the Go form of the generated JAXB class XmlGeneralName
// (complexType GeneralName).
type XmlGeneralName struct {
	XmlGeneralNameContent
	XmlGeneralNameAttrs
}

// XmlOIDContent carries the element content of the XmlOID base type. JAXB emits base
// content before extension content, so derived types embed it first.
type XmlOIDContent struct {
	Value string `xml:",chardata"`
}

// XmlOIDAttrs carries the attributes of the XmlOID base type. The JAXB RI writes
// extension attributes before base attributes, so derived types embed it last.
type XmlOIDAttrs struct {
	Description *string `xml:"Description,attr,omitempty"`
}

// XmlOID is the Go form of the generated JAXB class XmlOID
// (complexType OID).
type XmlOID struct {
	XmlOIDContent
	XmlOIDAttrs
}

// XmlAuthorityInformationAccess is the Go form of the generated JAXB class XmlAuthorityInformationAccess
// (complexType AuthorityInformationAccess).
type XmlAuthorityInformationAccess struct {
	XmlCertificateExtensionContent
	CaIssuersUrl []string `xml:"caIssuersUrl"`
	OcspUrl      []string `xml:"ocspUrl"`
	XmlCertificateExtensionAttrs
}

// XmlAuthorityKeyIdentifier is the Go form of the generated JAXB class XmlAuthorityKeyIdentifier
// (complexType AuthorityKeyIdentifier).
type XmlAuthorityKeyIdentifier struct {
	XmlCertificateExtensionContent
	KeyIdentifier             *Base64Binary `xml:"keyIdentifier,omitempty"`
	AuthorityCertIssuerSerial *Base64Binary `xml:"authorityCertIssuerSerial,omitempty"`
	XmlCertificateExtensionAttrs
}

// XmlBasicConstraints is the Go form of the generated JAXB class XmlBasicConstraints
// (complexType BasicConstraints).
type XmlBasicConstraints struct {
	XmlCertificateExtensionContent
	CA                bool `xml:"CA,attr"`
	PathLenConstraint *int `xml:"pathLenConstraint,attr,omitempty"`
	XmlCertificateExtensionAttrs
}

// XmlCRLDistributionPointsContent carries the element content of the XmlCRLDistributionPoints base type. JAXB emits base
// content before extension content, so derived types embed it first.
type XmlCRLDistributionPointsContent struct {
	CrlUrl []string `xml:"crlUrl"`
}

// XmlCRLDistributionPoints is the Go form of the generated JAXB class XmlCRLDistributionPoints
// (complexType CRLDistributionPoints).
type XmlCRLDistributionPoints struct {
	XmlCertificateExtensionContent
	XmlCRLDistributionPointsContent
	XmlCertificateExtensionAttrs
}

// XmlCertificatePolicies is the Go form of the generated JAXB class XmlCertificatePolicies
// (complexType CertificatePolicies).
type XmlCertificatePolicies struct {
	XmlCertificateExtensionContent
	CertificatePolicy []*XmlCertificatePolicy `xml:"certificatePolicy"`
	XmlCertificateExtensionAttrs
}

// XmlCertificatePolicy is the Go form of the generated JAXB class XmlCertificatePolicy
// (complexType CertificatePolicy).
type XmlCertificatePolicy struct {
	XmlOIDContent
	CpsUrl *string `xml:"cpsUrl,attr,omitempty"`
	XmlOIDAttrs
}

// XmlExtendedKeyUsages is the Go form of the generated JAXB class XmlExtendedKeyUsages
// (complexType ExtendedKeyUsages).
type XmlExtendedKeyUsages struct {
	XmlCertificateExtensionContent
	ExtendedKeyUsageOid []*XmlOID `xml:"extendedKeyUsageOid"`
	XmlCertificateExtensionAttrs
}

// XmlGeneralSubtree is the Go form of the generated JAXB class XmlGeneralSubtree
// (complexType GeneralSubtree).
type XmlGeneralSubtree struct {
	XmlGeneralNameContent
	Minimum *big.Int `xml:"minimum,attr,omitempty"`
	Maximum *big.Int `xml:"maximum,attr,omitempty"`
	XmlGeneralNameAttrs
}

// XmlIdPkixOcspNoCheck is the Go form of the generated JAXB class XmlIdPkixOcspNoCheck
// (complexType IdPkixOcspNoCheck).
type XmlIdPkixOcspNoCheck struct {
	XmlCertificateExtensionContent
	Present *bool `xml:"present,attr,omitempty"`
	XmlCertificateExtensionAttrs
}

// XmlInhibitAnyPolicy is the Go form of the generated JAXB class XmlInhibitAnyPolicy
// (complexType InhibitAnyPolicy).
type XmlInhibitAnyPolicy struct {
	XmlCertificateExtensionContent
	Value *int `xml:"value,attr,omitempty"`
	XmlCertificateExtensionAttrs
}

// XmlKeyUsages is the Go form of the generated JAXB class XmlKeyUsages
// (complexType KeyUsages).
type XmlKeyUsages struct {
	XmlCertificateExtensionContent
	KeyUsageBit []KeyUsageBitValue `xml:"keyUsageBit"`
	XmlCertificateExtensionAttrs
}

// XmlNameConstraints is the Go form of the generated JAXB class XmlNameConstraints
// (complexType NameConstraints).
type XmlNameConstraints struct {
	XmlCertificateExtensionContent
	PermittedSubtree []*XmlGeneralSubtree `xml:"PermittedSubtree"`
	ExcludedSubtree  []*XmlGeneralSubtree `xml:"ExcludedSubtree"`
	XmlCertificateExtensionAttrs
}

// XmlNoRevAvail is the Go form of the generated JAXB class XmlNoRevAvail
// (complexType NoRevAvail).
type XmlNoRevAvail struct {
	XmlCertificateExtensionContent
	Present *bool `xml:"present,attr,omitempty"`
	XmlCertificateExtensionAttrs
}

// XmlPolicyConstraints is the Go form of the generated JAXB class XmlPolicyConstraints
// (complexType PolicyConstraints).
type XmlPolicyConstraints struct {
	XmlCertificateExtensionContent
	RequireExplicitPolicy *int `xml:"requireExplicitPolicy,attr,omitempty"`
	InhibitPolicyMapping  *int `xml:"inhibitPolicyMapping,attr,omitempty"`
	XmlCertificateExtensionAttrs
}

// XmlSubjectAlternativeNames is the Go form of the generated JAXB class XmlSubjectAlternativeNames
// (complexType SubjectAlternativeNames).
type XmlSubjectAlternativeNames struct {
	XmlCertificateExtensionContent
	SubjectAlternativeName []*XmlGeneralName `xml:"subjectAlternativeName"`
	XmlCertificateExtensionAttrs
}

// XmlSubjectKeyIdentifier is the Go form of the generated JAXB class XmlSubjectKeyIdentifier
// (complexType SubjectKeyIdentifier).
type XmlSubjectKeyIdentifier struct {
	XmlCertificateExtensionContent
	Ski *Base64Binary `xml:"ski,omitempty"`
	XmlCertificateExtensionAttrs
}

// XmlValAssuredShortTermCertificate is the Go form of the generated JAXB class XmlValAssuredShortTermCertificate
// (complexType ValAssuredShortTermCertificate).
type XmlValAssuredShortTermCertificate struct {
	XmlCertificateExtensionContent
	Present *bool `xml:"present,attr,omitempty"`
	XmlCertificateExtensionAttrs
}

// XmlFreshestCRL is the Go form of the generated JAXB class XmlFreshestCRL
// (complexType FreshestCRL).
type XmlFreshestCRL struct {
	XmlCertificateExtensionContent
	XmlCRLDistributionPointsContent
	XmlCertificateExtensionAttrs
}
