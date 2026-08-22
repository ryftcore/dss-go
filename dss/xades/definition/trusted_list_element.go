// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/tsl/TrustedListElement.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// TrustedListElement is a list of TS 119 612 XSD Trusted List elements.
type TrustedListElement string

// TrustedListElement constants, one per TS 119 612 Trusted List XSD element name.
const (
	TrustedListElementAdditionalInformation        TrustedListElement = "ADDITIONAL_INFORMATION"
	TrustedListElementAdditionalServiceInformation TrustedListElement = "ADDITIONAL_SERVICE_INFORMATION"
	TrustedListElementCountryName                  TrustedListElement = "COUNTRY_NAME"
	TrustedListElementDateTime                     TrustedListElement = "DATE_TIME"
	TrustedListElementDigitalID                    TrustedListElement = "DIGITAL_ID"
	TrustedListElementDistributionPoints           TrustedListElement = "DISTRIBUTION_POINTS"
	TrustedListElementElectronicAddress            TrustedListElement = "ELECTRONIC_ADDRESS"
	TrustedListElementExpiredCertsRevocationInfo   TrustedListElement = "EXPIRED_CERTS_REVOCATION_INFO"
	TrustedListElementExtension                    TrustedListElement = "EXTENSION"
	TrustedListElementHistoricalInformationPeriod  TrustedListElement = "HISTORICAL_INFORMATION_PERIOD"
	TrustedListElementInformationValue             TrustedListElement = "INFORMATION_VALUE"
	TrustedListElementListIssueDateTime            TrustedListElement = "LIST_ISSUE_DATE_TIME"
	TrustedListElementLocality                     TrustedListElement = "LOCALITY"
	TrustedListElementName                         TrustedListElement = "NAME"
	TrustedListElementNextUpdate                   TrustedListElement = "NEXT_UPDATE"
	TrustedListElementOther                        TrustedListElement = "OTHER"
	TrustedListElementOtherInformation             TrustedListElement = "OTHER_INFORMATION"
	TrustedListElementOtherTSLPointer              TrustedListElement = "OTHER_TSL_POINTER"
	TrustedListElementPointersToOtherTSL           TrustedListElement = "POINTERS_TO_OTHER_TSL"
	TrustedListElementPolicyOrLegalNotice          TrustedListElement = "POLICY_OR_LEGAL_NOTICE"
	TrustedListElementPostalAddress                TrustedListElement = "POSTAL_ADDRESS"
	TrustedListElementPostalAddresses              TrustedListElement = "POSTAL_ADDRESSES"
	TrustedListElementPostalCode                   TrustedListElement = "POSTAL_CODE"
	TrustedListElementSchemeExtension              TrustedListElement = "SCHEME_EXTENSION"
	TrustedListElementSchemeInformation            TrustedListElement = "SCHEME_INFORMATION"
	TrustedListElementSchemeInformationURI         TrustedListElement = "SCHEME_INFORMATION_URI"
	TrustedListElementSchemeName                   TrustedListElement = "SCHEME_NAME"
	TrustedListElementSchemeOperatorAddress        TrustedListElement = "SCHEME_OPERATOR_ADDRESS"
	TrustedListElementSchemeOperatorName           TrustedListElement = "SCHEME_OPERATOR_NAME"
	TrustedListElementSchemeServiceDefinitionURI   TrustedListElement = "SCHEME_SERVICE_DEFINITION_URI"
	TrustedListElementSchemeTerritory              TrustedListElement = "SCHEME_TERRITORY"
	TrustedListElementSchemeTypeCommunityRules     TrustedListElement = "SCHEME_TYPE_COMMUNITY_RULES"
	TrustedListElementServiceDigitalIdentity       TrustedListElement = "SERVICE_DIGITAL_IDENTITY"
	TrustedListElementServiceDigitalIdentities     TrustedListElement = "SERVICE_DIGITAL_IDENTITIES"
	TrustedListElementServiceHistory               TrustedListElement = "SERVICE_HISTORY"
	TrustedListElementServiceHistoryInstance       TrustedListElement = "SERVICE_HISTORY_INSTANCE"
	TrustedListElementServiceInformation           TrustedListElement = "SERVICE_INFORMATION"
	TrustedListElementServiceInformationExtensions TrustedListElement = "SERVICE_INFORMATION_EXTENSIONS"
	TrustedListElementServiceTypeIdentifier        TrustedListElement = "SERVICE_TYPE_IDENTIFIER"
	TrustedListElementServiceName                  TrustedListElement = "SERVICE_NAME"
	TrustedListElementServiceStatus                TrustedListElement = "SERVICE_STATUS"
	TrustedListElementServiceSupplyPoint           TrustedListElement = "SERVICE_SUPPLY_POINT"
	TrustedListElementServiceSupplyPoints          TrustedListElement = "SERVICE_SUPPLY_POINTS"
	TrustedListElementStateOrProvince              TrustedListElement = "STATE_OR_PROVINCE"
	TrustedListElementStatusDeterminationApproach  TrustedListElement = "STATUS_DETERMINATION_APPROACH"
	TrustedListElementStatusStartingTime           TrustedListElement = "STATUS_STARTING_TIME"
	TrustedListElementStreetAddress                TrustedListElement = "STREET_ADDRESS"
	TrustedListElementTextualInformation           TrustedListElement = "TEXTUAL_INFORMATION"
	TrustedListElementTrustServiceProvider         TrustedListElement = "TRUST_SERVICE_PROVIDER"
	TrustedListElementTrustServiceProviderList     TrustedListElement = "TRUST_SERVICE_PROVIDER_LIST"
	TrustedListElementTrustServiceStatusList       TrustedListElement = "TRUST_SERVICE_STATUS_LIST"
	TrustedListElementTSLAddress                   TrustedListElement = "TSL_ADDRESS"
	TrustedListElementTSLInformation               TrustedListElement = "TSL_INFORMATION"
	TrustedListElementTSLInformationExtensions     TrustedListElement = "TSL_INFORMATION_EXTENSIONS"
	TrustedListElementTSLInformationURI            TrustedListElement = "TSL_INFORMATION_URI"
	TrustedListElementTSLLegalNotice               TrustedListElement = "TSL_LEGAL_NOTICE"
	TrustedListElementTSLLocation                  TrustedListElement = "TSL_LOCATION"
	TrustedListElementTSLName                      TrustedListElement = "TSL_NAME"
	TrustedListElementTSLPolicy                    TrustedListElement = "TSL_POLICY"
	TrustedListElementTSLSequenceNumber            TrustedListElement = "TSL_SEQUENCE_NUMBER"
	TrustedListElementTSLType                      TrustedListElement = "TSL_TYPE"
	TrustedListElementTSLVersionIdentifier         TrustedListElement = "TSL_VERSION_IDENTIFIER"
	TrustedListElementTSPService                   TrustedListElement = "TSP_SERVICE"
	TrustedListElementTSPServices                  TrustedListElement = "TSP_SERVICES"
	TrustedListElementTSPServiceDefinitionURI      TrustedListElement = "TSP_SERVICE_DEFINITION_URI"
	TrustedListElementTSPTradeName                 TrustedListElement = "TSP_TRADE_NAME"
	TrustedListElementURI                          TrustedListElement = "URI"
	TrustedListElementX509Certificate              TrustedListElement = "X509_CERTIFICATE"
	TrustedListElementX509SKI                      TrustedListElement = "X509_SKI"
	TrustedListElementX509SubjectName              TrustedListElement = "X509_SUBJECT_NAME"
)

// trustedListElementTagNames maps each constant to its wire tag name (getTagName()).
var trustedListElementTagNames = map[TrustedListElement]string{
	TrustedListElementAdditionalInformation:        "AdditionalInformation",
	TrustedListElementAdditionalServiceInformation: "AdditionalServiceInformation",
	TrustedListElementCountryName:                  "CountryName",
	TrustedListElementDateTime:                     "dateTime",
	TrustedListElementDigitalID:                    "DigitalId",
	TrustedListElementDistributionPoints:           "DistributionPoints",
	TrustedListElementElectronicAddress:            "ElectronicAddress",
	TrustedListElementExpiredCertsRevocationInfo:   "ExpiredCertsRevocationInfo",
	TrustedListElementExtension:                    "Extension",
	TrustedListElementHistoricalInformationPeriod:  "HistoricalInformationPeriod",
	TrustedListElementInformationValue:             "InformationValue",
	TrustedListElementListIssueDateTime:            "ListIssueDateTime",
	TrustedListElementLocality:                     "Locality",
	TrustedListElementName:                         "Name",
	TrustedListElementNextUpdate:                   "NextUpdate",
	TrustedListElementOther:                        "Other",
	TrustedListElementOtherInformation:             "OtherInformation",
	TrustedListElementOtherTSLPointer:              "OtherTSLPointer",
	TrustedListElementPointersToOtherTSL:           "PointersToOtherTSL",
	TrustedListElementPolicyOrLegalNotice:          "PolicyOrLegalNotice",
	TrustedListElementPostalAddress:                "PostalAddress",
	TrustedListElementPostalAddresses:              "PostalAddresses",
	TrustedListElementPostalCode:                   "PostalCode",
	TrustedListElementSchemeExtension:              "SchemeExtensions",
	TrustedListElementSchemeInformation:            "SchemeInformation",
	TrustedListElementSchemeInformationURI:         "SchemeInformationURI",
	TrustedListElementSchemeName:                   "SchemeName",
	TrustedListElementSchemeOperatorAddress:        "SchemeOperatorAddress",
	TrustedListElementSchemeOperatorName:           "SchemeOperatorName",
	TrustedListElementSchemeServiceDefinitionURI:   "SchemeServiceDefinitionURI",
	TrustedListElementSchemeTerritory:              "SchemeTerritory",
	TrustedListElementSchemeTypeCommunityRules:     "SchemeTypeCommunityRules",
	TrustedListElementServiceDigitalIdentity:       "ServiceDigitalIdentity",
	TrustedListElementServiceDigitalIdentities:     "ServiceDigitalIdentities",
	TrustedListElementServiceHistory:               "ServiceHistory",
	TrustedListElementServiceHistoryInstance:       "ServiceHistoryInstance",
	TrustedListElementServiceInformation:           "ServiceInformation",
	TrustedListElementServiceInformationExtensions: "ServiceInformationExtensions",
	TrustedListElementServiceTypeIdentifier:        "ServiceTypeIdentifier",
	TrustedListElementServiceName:                  "ServiceName",
	TrustedListElementServiceStatus:                "ServiceStatus",
	TrustedListElementServiceSupplyPoint:           "ServiceSupplyPoint",
	TrustedListElementServiceSupplyPoints:          "ServiceSupplyPoints",
	TrustedListElementStateOrProvince:              "StateOrProvince",
	TrustedListElementStatusDeterminationApproach:  "StatusDeterminationApproach",
	TrustedListElementStatusStartingTime:           "StatusStartingTime",
	TrustedListElementStreetAddress:                "StreetAddress",
	TrustedListElementTextualInformation:           "TextualInformation",
	TrustedListElementTrustServiceProvider:         "TrustServiceProvider",
	TrustedListElementTrustServiceProviderList:     "TrustServiceProviderList",
	TrustedListElementTrustServiceStatusList:       "TrustServiceStatusList",
	TrustedListElementTSLAddress:                   "TSPAddress",
	TrustedListElementTSLInformation:               "TSPInformation",
	TrustedListElementTSLInformationExtensions:     "TSPInformationExtensions",
	TrustedListElementTSLInformationURI:            "TSPInformationURI",
	TrustedListElementTSLLegalNotice:               "TSLLegalNotice",
	TrustedListElementTSLLocation:                  "TSLLocation",
	TrustedListElementTSLName:                      "TSPName",
	TrustedListElementTSLPolicy:                    "TSLPolicy",
	TrustedListElementTSLSequenceNumber:            "TSLSequenceNumber",
	TrustedListElementTSLType:                      "TSLType",
	TrustedListElementTSLVersionIdentifier:         "TSLVersionIdentifier",
	TrustedListElementTSPService:                   "TSPService",
	TrustedListElementTSPServices:                  "TSPServices",
	TrustedListElementTSPServiceDefinitionURI:      "TSPServiceDefinitionURI",
	TrustedListElementTSPTradeName:                 "TSPTradeName",
	TrustedListElementURI:                          "URI",
	TrustedListElementX509Certificate:              "X509Certificate",
	TrustedListElementX509SKI:                      "X509SKI",
	TrustedListElementX509SubjectName:              "X509SubjectName",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e TrustedListElement) TagName() string {
	return trustedListElementTagNames[e]
}

// Namespace implements common.DSSElement. Ports getNamespace().
func (e TrustedListElement) Namespace() *common.DSSNamespace {
	return TrustedListNamespaceNS
}

// URI implements common.DSSElement. Ports getURI().
func (e TrustedListElement) URI() string {
	return TrustedListNamespaceNS.Uri()
}

// IsSameTagName implements common.DSSElement. Ports isSameTagName(String).
func (e TrustedListElement) IsSameTagName(value string) bool {
	return e.TagName() == value
}

var _ common.DSSElement = TrustedListElement("")
