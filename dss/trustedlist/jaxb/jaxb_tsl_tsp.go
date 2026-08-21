// Ported from ts_119612v020401_xsd.xsd (DSS 6.5.RC1) via the JAXB classes
// generated into eu.europa.esig.trustedlist.jaxb.tsl. Per the phase-8a
// generated-JAXB rule the generated classes are grouped into schema-area
// files rather than one file per class; every Java class keeps its name and
// its exact field order so that encoding/xml reproduces the JAXB element
// sequence byte for byte. This file holds the TSP-information and
// scheme-information types.
package jaxb

// TSPInformationType is the Go form of the generated JAXB class
// TSPInformationType (complexType TSPInformationType).
type TSPInformationType struct {
	TSPName                  *InternationalNamesType       `xml:"TSPName"`
	TSPTradeName             *InternationalNamesType       `xml:"TSPTradeName,omitempty"`
	TSPAddress               *AddressType                  `xml:"TSPAddress"`
	TSPInformationURI        *NonEmptyMultiLangURIListType `xml:"TSPInformationURI"`
	TSPInformationExtensions *ExtensionsListType           `xml:"TSPInformationExtensions,omitempty"`
}

// TSPType is the Go form of the generated JAXB class TSPType (complexType
// TSPType).
type TSPType struct {
	TSPInformation *TSPInformationType  `xml:"TSPInformation"`
	TSPServices    *TSPServicesListType `xml:"TSPServices"`
}

// TrustServiceProviderListType is the Go form of the generated JAXB class
// TrustServiceProviderListType (complexType TrustServiceProviderListType).
type TrustServiceProviderListType struct {
	TrustServiceProvider []*TSPType `xml:"TrustServiceProvider"`
}

// TSLSchemeInformationType is the Go form of the generated JAXB class
// TSLSchemeInformationType (complexType TSLSchemeInformationType).
// ListIssueDateTime is the raw xs:dateTime lexical form, round-tripped
// verbatim - see jaxb_tsl_root.go's header.
type TSLSchemeInformationType struct {
	TSLVersionIdentifier        *BigInteger                   `xml:"TSLVersionIdentifier"`
	TSLSequenceNumber           *BigInteger                   `xml:"TSLSequenceNumber"`
	TSLType                     string                        `xml:"TSLType"`
	SchemeOperatorName          *InternationalNamesType       `xml:"SchemeOperatorName"`
	SchemeOperatorAddress       *AddressType                  `xml:"SchemeOperatorAddress"`
	SchemeName                  *InternationalNamesType       `xml:"SchemeName"`
	SchemeInformationURI        *NonEmptyMultiLangURIListType `xml:"SchemeInformationURI"`
	StatusDeterminationApproach string                        `xml:"StatusDeterminationApproach"`
	SchemeTypeCommunityRules    *NonEmptyMultiLangURIListType `xml:"SchemeTypeCommunityRules,omitempty"`
	SchemeTerritory             *string                       `xml:"SchemeTerritory,omitempty"`
	PolicyOrLegalNotice         *PolicyOrLegalnoticeType      `xml:"PolicyOrLegalNotice,omitempty"`
	HistoricalInformationPeriod *BigInteger                   `xml:"HistoricalInformationPeriod"`
	PointersToOtherTSL          *OtherTSLPointersType         `xml:"PointersToOtherTSL,omitempty"`
	ListIssueDateTime           *string                       `xml:"ListIssueDateTime,omitempty"`
	NextUpdate                  *NextUpdateType               `xml:"NextUpdate"`
	DistributionPoints          *NonEmptyURIListType          `xml:"DistributionPoints,omitempty"`
	SchemeExtensions            *ExtensionsListType           `xml:"SchemeExtensions,omitempty"`
}
