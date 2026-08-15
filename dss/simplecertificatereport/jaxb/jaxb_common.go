// Ported from SimpleCertificateReport.xsd (DSS 6.5.RC1) via the JAXB
// classes generated into eu.europa.esig.dss.simplecertificatereport.jaxb.
// Per the phase-8a generated-JAXB rule the generated classes are grouped
// into schema-area files rather than one file per class; every Java class
// keeps its name and its exact field order so that encoding/xml reproduces
// the JAXB element sequence byte for byte.

package jaxb

// XmlValidationPolicy is the Go form of the generated JAXB class
// XmlValidationPolicy (complexType ValidationPolicy).
type XmlValidationPolicy struct {
	PolicyName        *string `xml:"PolicyName,omitempty"`
	PolicyDescription *string `xml:"PolicyDescription,omitempty"`
}

// XmlConnectionDetails is the Go form of the generated JAXB class
// XmlConnectionDetails (complexType ConnectionDetails).
type XmlConnectionDetails struct {
	Url                       *string `xml:"Url,omitempty"`
	TLSCertificateBindingLink *string `xml:"TLSCertificateBindingLink,omitempty"`
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

// XmlSignatureScope is the Go form of the generated JAXB class
// XmlSignatureScope (complexType SignatureScope).
type XmlSignatureScope struct {
	Value string                   `xml:",chardata"`
	Id    string                   `xml:"Id,attr"`
	Name  *string                  `xml:"name,attr,omitempty"`
	Scope *SignatureScopeTypeValue `xml:"scope,attr,omitempty"`
}

// XmlSubject is the Go form of the generated JAXB class XmlSubject
// (complexType Subject).
type XmlSubject struct {
	CommonName       *string `xml:"commonName,omitempty"`
	Surname          *string `xml:"surname,omitempty"`
	GivenName        *string `xml:"givenName,omitempty"`
	Pseudonym        *string `xml:"pseudonym,omitempty"`
	OrganizationName *string `xml:"organizationName,omitempty"`
	OrganizationUnit *string `xml:"organizationUnit,omitempty"`
	Email            *string `xml:"email,omitempty"`
	Locality         *string `xml:"locality,omitempty"`
	State            *string `xml:"state,omitempty"`
	Country          *string `xml:"country,omitempty"`
}

// XmlRevocation is the Go form of the generated JAXB class XmlRevocation
// (complexType Revocation).
type XmlRevocation struct {
	ThisUpdate       *XSDateTime            `xml:"thisUpdate,omitempty"`
	RevocationDate   *XSDateTime            `xml:"revocationDate,omitempty"`
	RevocationReason *RevocationReasonValue `xml:"revocationReason,omitempty"`
}

// XmlKeyUsages is the Go form of the @XmlElementWrapper collection
// XmlChainItem.keyUsages. xjc has no standalone Java class for the KeyUsages
// complexType (it only appears wrapping a repeating element), but the
// wrapper element genuinely exists in the marshalled XML, so this port
// gives it the schema's own complexType name to marshal/unmarshal it
// explicitly. See jaxb_chain_item.go's header for why this - unlike Chain -
// is a pointer, omitted rather than self-closed when there is nothing to
// report.
type XmlKeyUsages struct {
	KeyUsage []KeyUsageBitValue `xml:"keyUsage"`
}

// XmlExtendedKeyUsages is the Go form of the @XmlElementWrapper collection
// XmlChainItem.extendedKeyUsages (complexType ExtendedKeyUsages).
type XmlExtendedKeyUsages struct {
	ExtendedKeyUsage []string `xml:"extendedKeyUsage"`
}

// XmlOcspUrls is the Go form of the @XmlElementWrapper collection
// XmlChainItem.ocspUrls (anonymous complexType).
type XmlOcspUrls struct {
	OcspUrl []string `xml:"ocspUrl"`
}

// XmlCrlUrls is the Go form of the @XmlElementWrapper collection
// XmlChainItem.crlUrls (anonymous complexType).
type XmlCrlUrls struct {
	CrlUrl []string `xml:"crlUrl"`
}

// XmlAiaUrls is the Go form of the @XmlElementWrapper collection
// XmlChainItem.aiaUrls (anonymous complexType).
type XmlAiaUrls struct {
	AiaUrl []string `xml:"aiaUrl"`
}

// XmlCpsUrls is the Go form of the @XmlElementWrapper collection
// XmlChainItem.cpsUrls (anonymous complexType).
type XmlCpsUrls struct {
	CpsUrl []string `xml:"cpsUrl"`
}

// XmlPdsUrls is the Go form of the @XmlElementWrapper collection
// XmlChainItem.pdsUrls (anonymous complexType).
type XmlPdsUrls struct {
	PdsUrl []string `xml:"pdsUrl"`
}

// XmlTrustAnchors is the Go form of the @XmlElementWrapper collection
// XmlChainItem.trustAnchors (complexType TrustAnchors).
type XmlTrustAnchors struct {
	TrustAnchor []*XmlTrustAnchor `xml:"trustAnchor"`
}

// XmlTrustAnchor is the Go form of the generated JAXB class XmlTrustAnchor
// (complexType TrustAnchor).
type XmlTrustAnchor struct {
	CountryCode                        string  `xml:"countryCode"`
	TslType                            *string `xml:"tslType,omitempty"`
	TrustServiceProvider               string  `xml:"trustServiceProvider"`
	TrustServiceProviderRegistrationId *string `xml:"trustServiceProviderRegistrationId,omitempty"`
	TrustServiceName                   string  `xml:"trustServiceName"`
}

// XmlCertificateApprovalStatusAtTime carries the shared content of the
// CertificateApprovalStatusAtIssuanceTime/AtValidationTime base type
// (complexType CertificateApprovalStatusAtTime).
type XmlCertificateApprovalStatusAtTime struct {
	CertificateApprovalStatus []*XmlCertificateApprovalStatus `xml:"certificateApprovalStatus,omitempty"`
}

// XmlCertificateApprovalStatusAtIssuanceTime is the Go form of the
// generated JAXB class XmlCertificateApprovalStatusAtIssuanceTime
// (complexType CertificateApprovalStatusAtIssuanceTime, an unextended
// extension of CertificateApprovalStatusAtTime).
type XmlCertificateApprovalStatusAtIssuanceTime struct {
	XmlCertificateApprovalStatusAtTime
}

// XmlCertificateApprovalStatusAtValidationTime is the Go form of the
// generated JAXB class XmlCertificateApprovalStatusAtValidationTime
// (complexType CertificateApprovalStatusAtValidationTime, an unextended
// extension of CertificateApprovalStatusAtTime).
type XmlCertificateApprovalStatusAtValidationTime struct {
	XmlCertificateApprovalStatusAtTime
}

// XmlCertificateApprovalStatus is the Go form of the generated JAXB class
// XmlCertificateApprovalStatus (complexType CertificateApprovalStatus).
// ListType/ServiceTypeIdentifier/ServiceStatus are declared xs:string in
// the schema but bound through XmlJavaTypeAdapter to the interface-typed
// enumerations ListType/LoTEServiceTypeIdentifier/LoTEServiceStatus; see
// jaxb_adapters.go.
type XmlCertificateApprovalStatus struct {
	ListType              *ListTypeValue                  `xml:"ListType,omitempty"`
	ServiceTypeIdentifier *LoTEServiceTypeIdentifierValue `xml:"ServiceTypeIdentifier,omitempty"`
	ServiceStatus         *LoTEServiceStatusValue         `xml:"ServiceStatus,omitempty"`
	Details               *XmlDetails                     `xml:"Details,omitempty"`
	Label                 *string                         `xml:"label,attr,omitempty"`
}
