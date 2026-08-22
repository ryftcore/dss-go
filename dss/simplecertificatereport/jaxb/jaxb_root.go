// Ported from SimpleCertificateReport.xsd (DSS 6.5.RC1) via the JAXB class
// generated into eu.europa.esig.dss.simplecertificatereport.jaxb.XmlSimpleCertificateReport.
// The generated classes are grouped into schema-area files rather than one
// file per class.
//
// Unlike dss/simplereport/jaxb's XmlSimpleReport, this root element's
// sequence has no choice group, so the ordinary declarative struct-tag
// encoding (XMLName plus per-field xml tags) is enough; no custom
// MarshalXML/UnmarshalXML is needed.

package jaxb

import "encoding/xml"

// XmlSimpleCertificateReport is the Go form of the generated JAXB class
// XmlSimpleCertificateReport (complexType SimpleCertificateReport, the
// document element).
type XmlSimpleCertificateReport struct {
	XMLName           xml.Name              `xml:"http://dss.esig.europa.eu/validation/simple-certificate-report SimpleCertificateReport"`
	ValidationPolicy  *XmlValidationPolicy  `xml:"ValidationPolicy"`
	Certificate       *XmlChainItem         `xml:"Certificate,omitempty"`
	ConnectionDetails *XmlConnectionDetails `xml:"ConnectionDetails,omitempty"`
	ValidationTime    *XSDateTime           `xml:"ValidationTime,attr,omitempty"`
}
