// Ported from DiagnosticData.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.diagnostic.jaxb. Per the phase-8a generated-JAXB rule the
// generated classes are grouped into schema-area files rather than one file per
// class; every Java class keeps its name and its exact field order so that
// encoding/xml reproduces the JAXB element sequence byte for byte.

package jaxb

// XmlMRATrustServiceMapping is the Go form of the generated JAXB class XmlMRATrustServiceMapping
// (complexType MRATrustServiceMapping).
type XmlMRATrustServiceMapping struct {
	TrustServiceLegalIdentifier   *string                                     `xml:"TrustServiceLegalIdentifier,omitempty"`
	EquivalenceStatusStartingTime *XSDateTime                                 `xml:"EquivalenceStatusStartingTime,omitempty"`
	EquivalenceStatusEndingTime   *XSDateTime                                 `xml:"EquivalenceStatusEndingTime,omitempty"`
	OriginalThirdCountryMapping   *XmlOriginalThirdCountryTrustServiceMapping `xml:"OriginalThirdCountryMapping,omitempty"`
}

// XmlOriginalThirdCountryTrustServiceMapping is the Go form of the generated JAXB class XmlOriginalThirdCountryTrustServiceMapping
// (complexType OriginalThirdCountryTrustServiceMapping).
type XmlOriginalThirdCountryTrustServiceMapping struct {
	ServiceType               *string                           `xml:"ServiceType,omitempty"`
	Status                    *string                           `xml:"Status,omitempty"`
	CapturedQualifiers        *CapturedQualifiersWrapper        `xml:"CapturedQualifiers"`
	AdditionalServiceInfoUris *AdditionalServiceInfoUrisWrapper `xml:"AdditionalServiceInfoUris"`
}

// XmlQualifier is the Go form of the generated JAXB class XmlQualifier
// (complexType Qualifier).
type XmlQualifier struct {
	Value    string `xml:",chardata"`
	Critical *bool  `xml:"critical,attr,omitempty"`
}

// XmlTrustServiceEquivalenceInformation is the Go form of the generated JAXB class XmlTrustServiceEquivalenceInformation
// (complexType TrustServiceEquivalenceInformation).
type XmlTrustServiceEquivalenceInformation struct {
	TrustServiceLegalIdentifier       *string                                   `xml:"TrustServiceLegalIdentifier,omitempty"`
	CertificateContentEquivalenceList *CertificateContentEquivalenceListWrapper `xml:"CertificateContentEquivalenceList"`
}

// XmlTrustServiceProvider is the Go form of the generated JAXB class XmlTrustServiceProvider
// (complexType TrustServiceProvider).
type XmlTrustServiceProvider struct {
	TSPNames                   *TSPNamesWrapper                   `xml:"TSPNames"`
	TSPTradeNames              *TSPTradeNamesWrapper              `xml:"TSPTradeNames"`
	TSPRegistrationIdentifiers *TSPRegistrationIdentifiersWrapper `xml:"TSPRegistrationIdentifiers"`
	TrustServices              *TrustServicesWrapper              `xml:"TrustServices"`
	TL                         *XmlTrustedList                    `xml:"TL,attr,omitempty"`
	LOTL                       *XmlTrustedList                    `xml:"LOTL,attr,omitempty"`
}

// XmlTrustSourceListContent carries the element content of the XmlTrustSourceList base type. JAXB emits base
// content before extension content, so derived types embed it first.
type XmlTrustSourceListContent struct {
	CountryCode          *string                  `xml:"CountryCode,omitempty"`
	Url                  *string                  `xml:"Url,omitempty"`
	Type                 *string                  `xml:"Type,omitempty"`
	SequenceNumber       *int                     `xml:"SequenceNumber,omitempty"`
	Version              *int                     `xml:"Version,omitempty"`
	LastLoading          *XSDateTime              `xml:"LastLoading,omitempty"`
	IssueDate            *XSDateTime              `xml:"IssueDate,omitempty"`
	NextUpdate           *XSDateTime              `xml:"NextUpdate,omitempty"`
	WellSigned           bool                     `xml:"WellSigned"`
	StructuralValidation *XmlStructuralValidation `xml:"StructuralValidation,omitempty"`
}

// XmlTrustSourceListAttrs carries the attributes of the XmlTrustSourceList base type. The JAXB RI writes
// extension attributes before base attributes, so derived types embed it last.
type XmlTrustSourceListAttrs struct {
	Id *CollapsedString `xml:"Id,attr,omitempty"`
}

// XmlTrustSourceList is the Go form of the generated JAXB class XmlTrustSourceList
// (complexType TrustSourceList).
type XmlTrustSourceList struct {
	XmlTrustSourceListContent
	XmlTrustSourceListAttrs
}

// XmlTrustedEntity is the Go form of the generated JAXB class XmlTrustedEntity
// (complexType TrustedEntity).
type XmlTrustedEntity struct {
	Names                   *NamesWrapper                   `xml:"Names"`
	TradeNames              *TradeNamesWrapper              `xml:"TradeNames"`
	RegistrationIdentifiers *RegistrationIdentifiersWrapper `xml:"RegistrationIdentifiers"`
	TrustedEntityServices   *TrustedEntityServicesWrapper   `xml:"TrustedEntityServices"`
	LoTE                    *XmlTrustSourceList             `xml:"LoTE,attr,omitempty"`
	LoLoTE                  *XmlTrustSourceList             `xml:"LoLoTE,attr,omitempty"`
}

// XmlTrustedEntityServiceContent carries the element content of the XmlTrustedEntityService base type. JAXB emits base
// content before extension content, so derived types embed it first.
type XmlTrustedEntityServiceContent struct {
	ServiceNames              *ServiceNamesWrapper              `xml:"ServiceNames"`
	ServiceType               *string                           `xml:"ServiceType,omitempty"`
	Status                    *string                           `xml:"Status,omitempty"`
	StartDate                 *XSDateTime                       `xml:"StartDate,omitempty"`
	EndDate                   *XSDateTime                       `xml:"EndDate,omitempty"`
	CapturedQualifiers        *CapturedQualifiersWrapper        `xml:"CapturedQualifiers"`
	AdditionalServiceInfoUris *AdditionalServiceInfoUrisWrapper `xml:"AdditionalServiceInfoUris"`
	ServiceSupplyPoints       *ServiceSupplyPointsWrapper       `xml:"ServiceSupplyPoints"`
}

// XmlTrustedEntityServiceAttrs carries the attributes of the XmlTrustedEntityService base type. The JAXB RI writes
// extension attributes before base attributes, so derived types embed it last.
type XmlTrustedEntityServiceAttrs struct {
	ServiceDigitalIdentifier *XmlCertificate `xml:"ServiceDigitalIdentifier,attr,omitempty"`
}

// XmlTrustedEntityService is the Go form of the generated JAXB class XmlTrustedEntityService
// (complexType TrustedEntityService).
type XmlTrustedEntityService struct {
	XmlTrustedEntityServiceContent
	XmlTrustedEntityServiceAttrs
}

// XmlListOfTrustedEntities is the Go form of the generated JAXB class XmlListOfTrustedEntities
// (complexType ListOfTrustedEntities).
type XmlListOfTrustedEntities struct {
	XmlTrustSourceListContent
	LoLoTE *bool                     `xml:"LoLoTE,attr,omitempty"`
	Parent *XmlListOfTrustedEntities `xml:"parent,attr,omitempty"`
	XmlTrustSourceListAttrs
}

// XmlTrustService is the Go form of the generated JAXB class XmlTrustService
// (complexType TrustService).
type XmlTrustService struct {
	XmlTrustedEntityServiceContent
	ExpiredCertsRevocationInfo *XSDateTime                `xml:"ExpiredCertsRevocationInfo,omitempty"`
	MRATrustServiceMapping     *XmlMRATrustServiceMapping `xml:"MRATrustServiceMapping,omitempty"`
	EnactedMRA                 *bool                      `xml:"enactedMRA,attr,omitempty"`
	XmlTrustedEntityServiceAttrs
}

// XmlTrustedList is the Go form of the generated JAXB class XmlTrustedList
// (complexType TrustedList).
type XmlTrustedList struct {
	XmlTrustSourceListContent
	LOTL   *bool           `xml:"LOTL,attr,omitempty"`
	Parent *XmlTrustedList `xml:"parent,attr,omitempty"`
	Mra    *bool           `xml:"mra,attr,omitempty"`
	XmlTrustSourceListAttrs
}
