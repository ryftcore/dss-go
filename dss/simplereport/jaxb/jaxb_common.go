// Ported from SimpleReport.xsd (DSS 6.5.RC1) via the JAXB classes generated
// into eu.europa.esig.dss.simplereport.jaxb. Per the phase-8a generated-JAXB
// rule the generated classes are grouped into schema-area files rather than
// one file per class; every Java class keeps its name and its exact field
// order so that encoding/xml reproduces the JAXB element sequence byte for
// byte.

package jaxb

// XmlValidationPolicy is the Go form of the generated JAXB class
// XmlValidationPolicy (complexType ValidationPolicy).
type XmlValidationPolicy struct {
	PolicyName        *string `xml:"PolicyName,omitempty"`
	PolicyDescription *string `xml:"PolicyDescription,omitempty"`
}

// XmlPDFAInfo is the Go form of the generated JAXB class XmlPDFAInfo
// (complexType PDFAInfo).
type XmlPDFAInfo struct {
	PDFAProfile        *string                `xml:"PDFAProfile,omitempty"`
	ValidationMessages *XmlValidationMessages `xml:"ValidationMessages,omitempty"`
	// Valid mirrors the generated isValid()'s Boolean return: the schema
	// attribute is optional and JAXB's isPDFACompliant() unboxes it without a
	// null check (a genuine upstream NPE risk on a document that carries a
	// PDFAInfo element without the valid attribute); this port preserves that
	// behaviour rather than papering over it, see SimpleReport.IsPDFACompliant.
	Valid *bool `xml:"valid,attr,omitempty"`
}

// XmlValidationMessages is the Go form of the generated JAXB anonymous
// complex type embedded as XmlPDFAInfo.ValidationMessages.
type XmlValidationMessages struct {
	Error []string `xml:"Error,omitempty"`
}

// XmlDetails is the Go form of the generated JAXB class XmlDetails
// (complexType Details).
type XmlDetails struct {
	Error   []*XmlMessage `xml:"Error,omitempty"`
	Warning []*XmlMessage `xml:"Warning,omitempty"`
	Info    []*XmlMessage `xml:"Info,omitempty"`
}

// XmlMessage is the Go form of the generated JAXB class XmlMessage
// (complexType Message).
type XmlMessage struct {
	Value string  `xml:",chardata"`
	Key   *string `xml:"Key,attr,omitempty"`
}

// XmlSemantic is the Go form of the generated JAXB class XmlSemantic
// (complexType Semantic).
type XmlSemantic struct {
	Value string `xml:",chardata"`
	Key   string `xml:"Key,attr"`
}

// XmlSignatureScope is the Go form of the generated JAXB class
// XmlSignatureScope (complexType SignatureScope), reused for the
// SignatureScope, TimestampScope and EvidenceRecordScope elements.
type XmlSignatureScope struct {
	Value string                   `xml:",chardata"`
	Id    string                   `xml:"Id,attr"`
	Name  *string                  `xml:"name,attr,omitempty"`
	Scope *SignatureScopeTypeValue `xml:"scope,attr,omitempty"`
}

// XmlSignatureLevel is the Go form of the generated JAXB class
// XmlSignatureLevel (complexType SignatureLevel, simpleContent over the
// SignatureQualification simple type - not to be confused with the
// enumerations.SignatureLevel Go type, which the SignatureFormat attribute
// of XmlSignature/XmlTimestamp binds through SignatureLevelValue instead).
type XmlSignatureLevel struct {
	Value       SignatureQualificationValue `xml:",chardata"`
	Description *string                     `xml:"description,attr,omitempty"`
}

// XmlTimestampLevel is the Go form of the generated JAXB class
// XmlTimestampLevel (complexType TimestampLevel).
type XmlTimestampLevel struct {
	Value       TimestampQualificationValue `xml:",chardata"`
	Description *string                     `xml:"description,attr,omitempty"`
}

// XmlEAALevel is the Go form of the generated JAXB class XmlEAALevel
// (complexType EAALevel).
type XmlEAALevel struct {
	Value       EAAQualificationValue `xml:",chardata"`
	Description *string               `xml:"description,attr,omitempty"`
}

// XmlCertificateChain is the Go form of the generated JAXB class
// XmlCertificateChain (complexType CertificateChain).
type XmlCertificateChain struct {
	Certificate []*XmlCertificate `xml:"Certificate,omitempty"`
}

// XmlCertificate is the Go form of the generated JAXB class XmlCertificate
// (complexType Certificate).
type XmlCertificate struct {
	QualifiedName string           `xml:"QualifiedName"`
	TrustAnchors  *XmlTrustAnchors `xml:"TrustAnchors,omitempty"`
	Id            string           `xml:"Id,attr"`
	// Trusted mirrors the generated isTrusted(), which defaults a nil
	// attribute to false rather than leaving it unset; Trusted() below
	// reproduces that default for callers that only have the pointer.
	Trusted    *bool       `xml:"trusted,attr,omitempty"`
	SunsetDate *XSDateTime `xml:"sunsetDate,attr,omitempty"`
}

// IsTrusted returns the trusted attribute, defaulting to false when absent.
// Port of XmlCertificate.isTrusted().
func (c *XmlCertificate) IsTrusted() bool {
	if c == nil || c.Trusted == nil {
		return false
	}
	return *c.Trusted
}

// XmlTrustAnchors is the Go form of the generated JAXB class
// XmlTrustAnchors (complexType TrustAnchors).
type XmlTrustAnchors struct {
	TrustAnchor []*XmlTrustAnchor `xml:"TrustAnchor,omitempty"`
}

// XmlTrustAnchor is the Go form of the generated JAXB class XmlTrustAnchor
// (complexType TrustAnchor).
type XmlTrustAnchor struct {
	TSLType                            *string  `xml:"TSLType,omitempty"`
	TrustServiceProvider               *string  `xml:"TrustServiceProvider,omitempty"`
	TrustServiceProviderRegistrationId *string  `xml:"TrustServiceProviderRegistrationId,omitempty"`
	TrustServiceName                   []string `xml:"TrustServiceName,omitempty"`
	CountryCode                        *string  `xml:"countryCode,attr,omitempty"`
}
