// Ported from SimpleCertificateReport.xsd (DSS 6.5.RC1) via the JAXB class
// generated into eu.europa.esig.dss.simplecertificatereport.jaxb.XmlChainItem.
// The generated classes are grouped into schema-area files rather than one
// file per class.
//
// xjc eagerly initialises every @XmlElementWrapper collection field of
// XmlChainItem (keyUsages, extendedKeyUsages, ocspUrls, crlUrls, aiaUrls,
// cpsUrls, pdsUrls, trustAnchors, chain) to an empty ArrayList at
// declaration, but that only pins the Java default - what actually reaches
// the marshaller depends on whether the (out-of-scope, dss-validation)
// SimpleCertificateReportBuilder calls the setter. The marshal-parity
// oracle corpus shows the builder only ever sets keyUsages/.../trustAnchors
// when it has something to report, leaving the field null - and therefore
// omitted, not self-closed - otherwise; Chain is the one exception, always
// set to a real (possibly empty) list. The seven omit-when-empty fields are
// therefore pointers to the small wrapper structs in jaxb_common.go, with
// the ordinary nil-omits/non-nil-writes handling used for every other
// optional complex element in this model; Chain keeps the always-marshalled
// nested-path slice.
//
// notBefore/notAfter are declared without minOccurs="0" (i.e. "required")
// in the schema, but "required" is a schema-validation constraint the RI
// does not enforce on marshal/unmarshal (validateXml is false throughout
// this port's KAT, matching AbstractJaxbFacade's default): a null Date
// field is simply omitted, same as any optional one. testdata/oracle's
// simple-cert-report2.xml - reserialised from an upstream test fixture that
// never populates them - exercises exactly this, so both fields are
// pointers here rather than plain XSDateTime.
package jaxb

// XmlChainItem is the Go form of the generated JAXB class XmlChainItem
// (complexType ChainItem).
type XmlChainItem struct {
	Subject                                   *XmlSubject                                   `xml:"subject"`
	IssuerId                                  *string                                       `xml:"issuerId,omitempty"`
	NotBefore                                 *XSDateTime                                   `xml:"notBefore,omitempty"`
	NotAfter                                  *XSDateTime                                   `xml:"notAfter,omitempty"`
	KeyUsages                                 *XmlKeyUsages                                 `xml:"keyUsages,omitempty"`
	ExtendedKeyUsages                         *XmlExtendedKeyUsages                         `xml:"extendedKeyUsages,omitempty"`
	OcspUrls                                  *XmlOcspUrls                                  `xml:"ocspUrls,omitempty"`
	CrlUrls                                   *XmlCrlUrls                                   `xml:"crlUrls,omitempty"`
	AiaUrls                                   *XmlAiaUrls                                   `xml:"aiaUrls,omitempty"`
	CpsUrls                                   *XmlCpsUrls                                   `xml:"cpsUrls,omitempty"`
	PdsUrls                                   *XmlPdsUrls                                   `xml:"pdsUrls,omitempty"`
	QualificationAtIssuance                   *CertificateQualificationValue                `xml:"qualificationAtIssuance,omitempty"`
	QualificationDetailsAtIssuance            *XmlDetails                                   `xml:"qualificationDetailsAtIssuance,omitempty"`
	QualificationAtValidation                 *CertificateQualificationValue                `xml:"qualificationAtValidation,omitempty"`
	QualificationDetailsAtValidation          *XmlDetails                                   `xml:"qualificationDetailsAtValidation,omitempty"`
	QwacProfile                               *QWACProfileValue                             `xml:"qwacProfile,omitempty"`
	QwacDetails                               *XmlDetails                                   `xml:"qwacDetails,omitempty"`
	CertificateApprovalStatusAtIssuanceTime   *XmlCertificateApprovalStatusAtIssuanceTime   `xml:"certificateApprovalStatusAtIssuanceTime,omitempty"`
	CertificateApprovalStatusAtValidationTime *XmlCertificateApprovalStatusAtValidationTime `xml:"certificateApprovalStatusAtValidationTime,omitempty"`
	EnactedMRA                                *bool                                         `xml:"enactedMRA,omitempty"`
	Revocation                                *XmlRevocation                                `xml:"revocation,omitempty"`
	TrustAnchors                              *XmlTrustAnchors                              `xml:"trustAnchors,omitempty"`
	TrustStartDate                            *XSDateTime                                   `xml:"trustStartDate,omitempty"`
	TrustSunsetDate                           *XSDateTime                                   `xml:"trustSunsetDate,omitempty"`
	Indication                                IndicationValue                               `xml:"Indication,omitempty"`
	SubIndication                             *SubIndicationValue                           `xml:"SubIndication,omitempty"`
	X509ValidationDetails                     *XmlDetails                                   `xml:"X509ValidationDetails,omitempty"`
	Chain                                     []*XmlChainItem                               `xml:"Chain>ChainItem"`
	TLSBindingSignature                       *XmlSignature                                 `xml:"TLSBindingSignature,omitempty"`
	Id                                        string                                        `xml:"Id,attr"`
}
