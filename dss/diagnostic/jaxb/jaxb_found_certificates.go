// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. The generated classes are grouped into
// schema-area files rather than one file per class; every Java class keeps its
// name and its exact field order so that encoding/xml reproduces the JAXB element
// sequence byte for byte.

package jaxb

// XmlFoundCertificateContent carries the element content of the XmlFoundCertificate base type. JAXB emits base
// content before extension content, so derived types embed it first.
type XmlFoundCertificateContent struct {
	Origin         []CertificateOriginValue `xml:"Origin"`
	CertificateRef []*XmlCertificateRef     `xml:"CertificateRef"`
}

// XmlFoundCertificate is the Go form of the generated JAXB class XmlFoundCertificate
// (complexType FoundCertificate).
type XmlFoundCertificate struct {
	XmlFoundCertificateContent
}

// XmlFoundCertificates is the Go form of the generated JAXB class XmlFoundCertificates
// (complexType FoundCertificates).
type XmlFoundCertificates struct {
	RelatedCertificate []*XmlRelatedCertificate `xml:"RelatedCertificate"`
	OrphanCertificate  []*XmlOrphanCertificate  `xml:"OrphanCertificate"`
}

// XmlOrphanCertificate is the Go form of the generated JAXB class XmlOrphanCertificate
// (complexType OrphanCertificate).
type XmlOrphanCertificate struct {
	XmlFoundCertificateContent
	Token *XmlOrphanCertificateToken `xml:"Token,attr,omitempty"`
}

// XmlRelatedCertificate is the Go form of the generated JAXB class XmlRelatedCertificate
// (complexType RelatedCertificate).
type XmlRelatedCertificate struct {
	XmlFoundCertificateContent
	Certificate *XmlCertificate `xml:"Certificate,attr,omitempty"`
}
