// Ported from SimpleCertificateReport.xsd (DSS 6.5.RC1) via the JAXB class
// generated into eu.europa.esig.dss.simplecertificatereport.jaxb.XmlSignature.
// Per the phase-8a generated-JAXB rule the generated classes are grouped
// into schema-area files rather than one file per class.
//
// Unlike dss/simplereport/jaxb.XmlSignature, this Signature complexType is
// not an extension of an abstract Token base (the certificate-report schema
// has no equivalent), so it is a single flat struct.

package jaxb

// XmlSignature is the Go form of the generated JAXB class XmlSignature
// (complexType Signature).
type XmlSignature struct {
	Url                   *string              `xml:"Url,omitempty"`
	Indication            IndicationValue      `xml:"Indication"`
	SubIndication         *SubIndicationValue  `xml:"SubIndication,omitempty"`
	AdESValidationDetails *XmlDetails          `xml:"AdESValidationDetails,omitempty"`
	SigningTime           *XSDateTime          `xml:"SigningTime,omitempty"`
	SignatureScope        []*XmlSignatureScope `xml:"SignatureScope,omitempty"`
	Chain                 []*XmlChainItem      `xml:"Chain>ChainItem"`
	Id                    string               `xml:"Id,attr"`
	SignatureFormat       SignatureLevelValue  `xml:"SignatureFormat,attr"`
}
