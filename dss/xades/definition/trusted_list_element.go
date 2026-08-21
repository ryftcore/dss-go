// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/tsl/TrustedListElement.java (DSS 6.5.RC1).
package definition

import "github.com/utain/esig/dss/xml/common"

// TrustedListElement is a list of TS 119 612 XSD Trusted List elements.
type TrustedListElement string

// TrustedListElement constants, one per TS 119 612 Trusted List XSD element name.
const (
	TrustedListElement_ADDITIONAL_INFORMATION         TrustedListElement = "ADDITIONAL_INFORMATION"
	TrustedListElement_ADDITIONAL_SERVICE_INFORMATION TrustedListElement = "ADDITIONAL_SERVICE_INFORMATION"
	TrustedListElement_COUNTRY_NAME                   TrustedListElement = "COUNTRY_NAME"
	TrustedListElement_DATE_TIME                      TrustedListElement = "DATE_TIME"
	TrustedListElement_DIGITAL_ID                     TrustedListElement = "DIGITAL_ID"
	TrustedListElement_DISTRIBUTION_POINTS            TrustedListElement = "DISTRIBUTION_POINTS"
	TrustedListElement_ELECTRONIC_ADDRESS             TrustedListElement = "ELECTRONIC_ADDRESS"
	TrustedListElement_EXPIRED_CERTS_REVOCATION_INFO  TrustedListElement = "EXPIRED_CERTS_REVOCATION_INFO"
	TrustedListElement_EXTENSION                      TrustedListElement = "EXTENSION"
	TrustedListElement_HISTORICAL_INFORMATION_PERIOD  TrustedListElement = "HISTORICAL_INFORMATION_PERIOD"
	TrustedListElement_INFORMATION_VALUE              TrustedListElement = "INFORMATION_VALUE"
	TrustedListElement_LIST_ISSUE_DATE_TIME           TrustedListElement = "LIST_ISSUE_DATE_TIME"
	TrustedListElement_LOCALITY                       TrustedListElement = "LOCALITY"
	TrustedListElement_NAME                           TrustedListElement = "NAME"
	TrustedListElement_NEXT_UPDATE                    TrustedListElement = "NEXT_UPDATE"
	TrustedListElement_OTHER                          TrustedListElement = "OTHER"
	TrustedListElement_OTHER_INFORMATION              TrustedListElement = "OTHER_INFORMATION"
	TrustedListElement_OTHER_TSL_POINTER              TrustedListElement = "OTHER_TSL_POINTER"
	TrustedListElement_POINTERS_TO_OTHER_TSL          TrustedListElement = "POINTERS_TO_OTHER_TSL"
	TrustedListElement_POLICY_OR_LEGAL_NOTICE         TrustedListElement = "POLICY_OR_LEGAL_NOTICE"
	TrustedListElement_POSTAL_ADDRESS                 TrustedListElement = "POSTAL_ADDRESS"
	TrustedListElement_POSTAL_ADDRESSES               TrustedListElement = "POSTAL_ADDRESSES"
	TrustedListElement_POSTAL_CODE                    TrustedListElement = "POSTAL_CODE"
	TrustedListElement_SCHEME_EXTENSION               TrustedListElement = "SCHEME_EXTENSION"
	TrustedListElement_SCHEME_INFORMATION             TrustedListElement = "SCHEME_INFORMATION"
	TrustedListElement_SCHEME_INFORMATION_URI         TrustedListElement = "SCHEME_INFORMATION_URI"
	TrustedListElement_SCHEME_NAME                    TrustedListElement = "SCHEME_NAME"
	TrustedListElement_SCHEME_OPERATOR_ADDRESS        TrustedListElement = "SCHEME_OPERATOR_ADDRESS"
	TrustedListElement_SCHEME_OPERATOR_NAME           TrustedListElement = "SCHEME_OPERATOR_NAME"
	TrustedListElement_SCHEME_SERVICE_DEFINITION_URI  TrustedListElement = "SCHEME_SERVICE_DEFINITION_URI"
	TrustedListElement_SCHEME_TERRITORY               TrustedListElement = "SCHEME_TERRITORY"
	TrustedListElement_SCHEME_TYPE_COMMUNITY_RULES    TrustedListElement = "SCHEME_TYPE_COMMUNITY_RULES"
	TrustedListElement_SERVICE_DIGITAL_IDENTITY       TrustedListElement = "SERVICE_DIGITAL_IDENTITY"
	TrustedListElement_SERVICE_DIGITAL_IDENTITIES     TrustedListElement = "SERVICE_DIGITAL_IDENTITIES"
	TrustedListElement_SERVICE_HISTORY                TrustedListElement = "SERVICE_HISTORY"
	TrustedListElement_SERVICE_HISTORY_INSTANCE       TrustedListElement = "SERVICE_HISTORY_INSTANCE"
	TrustedListElement_SERVICE_INFORMATION            TrustedListElement = "SERVICE_INFORMATION"
	TrustedListElement_SERVICE_INFORMATION_EXTENSIONS TrustedListElement = "SERVICE_INFORMATION_EXTENSIONS"
	TrustedListElement_SERVICE_TYPE_IDENTIFIER        TrustedListElement = "SERVICE_TYPE_IDENTIFIER"
	TrustedListElement_SERVICE_NAME                   TrustedListElement = "SERVICE_NAME"
	TrustedListElement_SERVICE_STATUS                 TrustedListElement = "SERVICE_STATUS"
	TrustedListElement_SERVICE_SUPPLY_POINT           TrustedListElement = "SERVICE_SUPPLY_POINT"
	TrustedListElement_SERVICE_SUPPLY_POINTS          TrustedListElement = "SERVICE_SUPPLY_POINTS"
	TrustedListElement_STATE_OR_PROVINCE              TrustedListElement = "STATE_OR_PROVINCE"
	TrustedListElement_STATUS_DETERMINATION_APPROACH  TrustedListElement = "STATUS_DETERMINATION_APPROACH"
	TrustedListElement_STATUS_STARTING_TIME           TrustedListElement = "STATUS_STARTING_TIME"
	TrustedListElement_STREET_ADDRESS                 TrustedListElement = "STREET_ADDRESS"
	TrustedListElement_TEXTUAL_INFORMATION            TrustedListElement = "TEXTUAL_INFORMATION"
	TrustedListElement_TRUST_SERVICE_PROVIDER         TrustedListElement = "TRUST_SERVICE_PROVIDER"
	TrustedListElement_TRUST_SERVICE_PROVIDER_LIST    TrustedListElement = "TRUST_SERVICE_PROVIDER_LIST"
	TrustedListElement_TRUST_SERVICE_STATUS_LIST      TrustedListElement = "TRUST_SERVICE_STATUS_LIST"
	TrustedListElement_TSL_ADDRESS                    TrustedListElement = "TSL_ADDRESS"
	TrustedListElement_TSL_INFORMATION                TrustedListElement = "TSL_INFORMATION"
	TrustedListElement_TSL_INFORMATION_EXTENSIONS     TrustedListElement = "TSL_INFORMATION_EXTENSIONS"
	TrustedListElement_TSL_INFORMATION_URI            TrustedListElement = "TSL_INFORMATION_URI"
	TrustedListElement_TSL_LEGAL_NOTICE               TrustedListElement = "TSL_LEGAL_NOTICE"
	TrustedListElement_TSL_LOCATION                   TrustedListElement = "TSL_LOCATION"
	TrustedListElement_TSL_NAME                       TrustedListElement = "TSL_NAME"
	TrustedListElement_TSL_POLICY                     TrustedListElement = "TSL_POLICY"
	TrustedListElement_TSL_SEQUENCE_NUMBER            TrustedListElement = "TSL_SEQUENCE_NUMBER"
	TrustedListElement_TSL_TYPE                       TrustedListElement = "TSL_TYPE"
	TrustedListElement_TSL_VERSION_IDENTIFIER         TrustedListElement = "TSL_VERSION_IDENTIFIER"
	TrustedListElement_TSP_SERVICE                    TrustedListElement = "TSP_SERVICE"
	TrustedListElement_TSP_SERVICES                   TrustedListElement = "TSP_SERVICES"
	TrustedListElement_TSP_SERVICE_DEFINITION_URI     TrustedListElement = "TSP_SERVICE_DEFINITION_URI"
	TrustedListElement_TSP_TRADE_NAME                 TrustedListElement = "TSP_TRADE_NAME"
	TrustedListElement_URI                            TrustedListElement = "URI"
	TrustedListElement_X509_CERTIFICATE               TrustedListElement = "X509_CERTIFICATE"
	TrustedListElement_X509_SKI                       TrustedListElement = "X509_SKI"
	TrustedListElement_X509_SUBJECT_NAME              TrustedListElement = "X509_SUBJECT_NAME"
)

// trustedListElementTagNames maps each constant to its wire tag name (getTagName()).
var trustedListElementTagNames = map[TrustedListElement]string{
	TrustedListElement_ADDITIONAL_INFORMATION:         "AdditionalInformation",
	TrustedListElement_ADDITIONAL_SERVICE_INFORMATION: "AdditionalServiceInformation",
	TrustedListElement_COUNTRY_NAME:                   "CountryName",
	TrustedListElement_DATE_TIME:                      "dateTime",
	TrustedListElement_DIGITAL_ID:                     "DigitalId",
	TrustedListElement_DISTRIBUTION_POINTS:            "DistributionPoints",
	TrustedListElement_ELECTRONIC_ADDRESS:             "ElectronicAddress",
	TrustedListElement_EXPIRED_CERTS_REVOCATION_INFO:  "ExpiredCertsRevocationInfo",
	TrustedListElement_EXTENSION:                      "Extension",
	TrustedListElement_HISTORICAL_INFORMATION_PERIOD:  "HistoricalInformationPeriod",
	TrustedListElement_INFORMATION_VALUE:              "InformationValue",
	TrustedListElement_LIST_ISSUE_DATE_TIME:           "ListIssueDateTime",
	TrustedListElement_LOCALITY:                       "Locality",
	TrustedListElement_NAME:                           "Name",
	TrustedListElement_NEXT_UPDATE:                    "NextUpdate",
	TrustedListElement_OTHER:                          "Other",
	TrustedListElement_OTHER_INFORMATION:              "OtherInformation",
	TrustedListElement_OTHER_TSL_POINTER:              "OtherTSLPointer",
	TrustedListElement_POINTERS_TO_OTHER_TSL:          "PointersToOtherTSL",
	TrustedListElement_POLICY_OR_LEGAL_NOTICE:         "PolicyOrLegalNotice",
	TrustedListElement_POSTAL_ADDRESS:                 "PostalAddress",
	TrustedListElement_POSTAL_ADDRESSES:               "PostalAddresses",
	TrustedListElement_POSTAL_CODE:                    "PostalCode",
	TrustedListElement_SCHEME_EXTENSION:               "SchemeExtensions",
	TrustedListElement_SCHEME_INFORMATION:             "SchemeInformation",
	TrustedListElement_SCHEME_INFORMATION_URI:         "SchemeInformationURI",
	TrustedListElement_SCHEME_NAME:                    "SchemeName",
	TrustedListElement_SCHEME_OPERATOR_ADDRESS:        "SchemeOperatorAddress",
	TrustedListElement_SCHEME_OPERATOR_NAME:           "SchemeOperatorName",
	TrustedListElement_SCHEME_SERVICE_DEFINITION_URI:  "SchemeServiceDefinitionURI",
	TrustedListElement_SCHEME_TERRITORY:               "SchemeTerritory",
	TrustedListElement_SCHEME_TYPE_COMMUNITY_RULES:    "SchemeTypeCommunityRules",
	TrustedListElement_SERVICE_DIGITAL_IDENTITY:       "ServiceDigitalIdentity",
	TrustedListElement_SERVICE_DIGITAL_IDENTITIES:     "ServiceDigitalIdentities",
	TrustedListElement_SERVICE_HISTORY:                "ServiceHistory",
	TrustedListElement_SERVICE_HISTORY_INSTANCE:       "ServiceHistoryInstance",
	TrustedListElement_SERVICE_INFORMATION:            "ServiceInformation",
	TrustedListElement_SERVICE_INFORMATION_EXTENSIONS: "ServiceInformationExtensions",
	TrustedListElement_SERVICE_TYPE_IDENTIFIER:        "ServiceTypeIdentifier",
	TrustedListElement_SERVICE_NAME:                   "ServiceName",
	TrustedListElement_SERVICE_STATUS:                 "ServiceStatus",
	TrustedListElement_SERVICE_SUPPLY_POINT:           "ServiceSupplyPoint",
	TrustedListElement_SERVICE_SUPPLY_POINTS:          "ServiceSupplyPoints",
	TrustedListElement_STATE_OR_PROVINCE:              "StateOrProvince",
	TrustedListElement_STATUS_DETERMINATION_APPROACH:  "StatusDeterminationApproach",
	TrustedListElement_STATUS_STARTING_TIME:           "StatusStartingTime",
	TrustedListElement_STREET_ADDRESS:                 "StreetAddress",
	TrustedListElement_TEXTUAL_INFORMATION:            "TextualInformation",
	TrustedListElement_TRUST_SERVICE_PROVIDER:         "TrustServiceProvider",
	TrustedListElement_TRUST_SERVICE_PROVIDER_LIST:    "TrustServiceProviderList",
	TrustedListElement_TRUST_SERVICE_STATUS_LIST:      "TrustServiceStatusList",
	TrustedListElement_TSL_ADDRESS:                    "TSPAddress",
	TrustedListElement_TSL_INFORMATION:                "TSPInformation",
	TrustedListElement_TSL_INFORMATION_EXTENSIONS:     "TSPInformationExtensions",
	TrustedListElement_TSL_INFORMATION_URI:            "TSPInformationURI",
	TrustedListElement_TSL_LEGAL_NOTICE:               "TSLLegalNotice",
	TrustedListElement_TSL_LOCATION:                   "TSLLocation",
	TrustedListElement_TSL_NAME:                       "TSPName",
	TrustedListElement_TSL_POLICY:                     "TSLPolicy",
	TrustedListElement_TSL_SEQUENCE_NUMBER:            "TSLSequenceNumber",
	TrustedListElement_TSL_TYPE:                       "TSLType",
	TrustedListElement_TSL_VERSION_IDENTIFIER:         "TSLVersionIdentifier",
	TrustedListElement_TSP_SERVICE:                    "TSPService",
	TrustedListElement_TSP_SERVICES:                   "TSPServices",
	TrustedListElement_TSP_SERVICE_DEFINITION_URI:     "TSPServiceDefinitionURI",
	TrustedListElement_TSP_TRADE_NAME:                 "TSPTradeName",
	TrustedListElement_URI:                            "URI",
	TrustedListElement_X509_CERTIFICATE:               "X509Certificate",
	TrustedListElement_X509_SKI:                       "X509SKI",
	TrustedListElement_X509_SUBJECT_NAME:              "X509SubjectName",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e TrustedListElement) TagName() string {
	return trustedListElementTagNames[e]
}

// Namespace implements common.DSSElement. Ports getNamespace().
func (e TrustedListElement) Namespace() *common.DSSNamespace {
	return TrustedListNamespace_NS
}

// URI implements common.DSSElement. Ports getURI().
func (e TrustedListElement) URI() string {
	return TrustedListNamespace_NS.Uri()
}

// IsSameTagName implements common.DSSElement. Ports isSameTagName(String).
func (e TrustedListElement) IsSameTagName(value string) bool {
	return e.TagName() == value
}

var _ common.DSSElement = TrustedListElement("")
