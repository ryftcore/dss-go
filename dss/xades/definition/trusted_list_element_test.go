// KAT test for trusted_list_element.go: every TrustedListElement_* tag name below is
// transcribed verbatim from the upstream Java source (dss-xades 6.5.RC1)
// eu.europa.esig.dss.xades.definition.tsl.TrustedListElement.
package definition

import "testing"

func TestTrustedListElement_KAT(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ADDITIONAL_INFORMATION", TrustedListElement_ADDITIONAL_INFORMATION.TagName(), "AdditionalInformation"},
		{"ADDITIONAL_SERVICE_INFORMATION", TrustedListElement_ADDITIONAL_SERVICE_INFORMATION.TagName(), "AdditionalServiceInformation"},
		{"COUNTRY_NAME", TrustedListElement_COUNTRY_NAME.TagName(), "CountryName"},
		{"DATE_TIME", TrustedListElement_DATE_TIME.TagName(), "dateTime"},
		{"DIGITAL_ID", TrustedListElement_DIGITAL_ID.TagName(), "DigitalId"},
		{"DISTRIBUTION_POINTS", TrustedListElement_DISTRIBUTION_POINTS.TagName(), "DistributionPoints"},
		{"ELECTRONIC_ADDRESS", TrustedListElement_ELECTRONIC_ADDRESS.TagName(), "ElectronicAddress"},
		{"EXPIRED_CERTS_REVOCATION_INFO", TrustedListElement_EXPIRED_CERTS_REVOCATION_INFO.TagName(), "ExpiredCertsRevocationInfo"},
		{"EXTENSION", TrustedListElement_EXTENSION.TagName(), "Extension"},
		{"HISTORICAL_INFORMATION_PERIOD", TrustedListElement_HISTORICAL_INFORMATION_PERIOD.TagName(), "HistoricalInformationPeriod"},
		{"INFORMATION_VALUE", TrustedListElement_INFORMATION_VALUE.TagName(), "InformationValue"},
		{"LIST_ISSUE_DATE_TIME", TrustedListElement_LIST_ISSUE_DATE_TIME.TagName(), "ListIssueDateTime"},
		{"LOCALITY", TrustedListElement_LOCALITY.TagName(), "Locality"},
		{"NAME", TrustedListElement_NAME.TagName(), "Name"},
		{"NEXT_UPDATE", TrustedListElement_NEXT_UPDATE.TagName(), "NextUpdate"},
		{"OTHER", TrustedListElement_OTHER.TagName(), "Other"},
		{"OTHER_INFORMATION", TrustedListElement_OTHER_INFORMATION.TagName(), "OtherInformation"},
		{"OTHER_TSL_POINTER", TrustedListElement_OTHER_TSL_POINTER.TagName(), "OtherTSLPointer"},
		{"POINTERS_TO_OTHER_TSL", TrustedListElement_POINTERS_TO_OTHER_TSL.TagName(), "PointersToOtherTSL"},
		{"POLICY_OR_LEGAL_NOTICE", TrustedListElement_POLICY_OR_LEGAL_NOTICE.TagName(), "PolicyOrLegalNotice"},
		{"POSTAL_ADDRESS", TrustedListElement_POSTAL_ADDRESS.TagName(), "PostalAddress"},
		{"POSTAL_ADDRESSES", TrustedListElement_POSTAL_ADDRESSES.TagName(), "PostalAddresses"},
		{"POSTAL_CODE", TrustedListElement_POSTAL_CODE.TagName(), "PostalCode"},
		{"SCHEME_EXTENSION", TrustedListElement_SCHEME_EXTENSION.TagName(), "SchemeExtensions"},
		{"SCHEME_INFORMATION", TrustedListElement_SCHEME_INFORMATION.TagName(), "SchemeInformation"},
		{"SCHEME_INFORMATION_URI", TrustedListElement_SCHEME_INFORMATION_URI.TagName(), "SchemeInformationURI"},
		{"SCHEME_NAME", TrustedListElement_SCHEME_NAME.TagName(), "SchemeName"},
		{"SCHEME_OPERATOR_ADDRESS", TrustedListElement_SCHEME_OPERATOR_ADDRESS.TagName(), "SchemeOperatorAddress"},
		{"SCHEME_OPERATOR_NAME", TrustedListElement_SCHEME_OPERATOR_NAME.TagName(), "SchemeOperatorName"},
		{"SCHEME_SERVICE_DEFINITION_URI", TrustedListElement_SCHEME_SERVICE_DEFINITION_URI.TagName(), "SchemeServiceDefinitionURI"},
		{"SCHEME_TERRITORY", TrustedListElement_SCHEME_TERRITORY.TagName(), "SchemeTerritory"},
		{"SCHEME_TYPE_COMMUNITY_RULES", TrustedListElement_SCHEME_TYPE_COMMUNITY_RULES.TagName(), "SchemeTypeCommunityRules"},
		{"SERVICE_DIGITAL_IDENTITY", TrustedListElement_SERVICE_DIGITAL_IDENTITY.TagName(), "ServiceDigitalIdentity"},
		{"SERVICE_DIGITAL_IDENTITIES", TrustedListElement_SERVICE_DIGITAL_IDENTITIES.TagName(), "ServiceDigitalIdentities"},
		{"SERVICE_HISTORY", TrustedListElement_SERVICE_HISTORY.TagName(), "ServiceHistory"},
		{"SERVICE_HISTORY_INSTANCE", TrustedListElement_SERVICE_HISTORY_INSTANCE.TagName(), "ServiceHistoryInstance"},
		{"SERVICE_INFORMATION", TrustedListElement_SERVICE_INFORMATION.TagName(), "ServiceInformation"},
		{"SERVICE_INFORMATION_EXTENSIONS", TrustedListElement_SERVICE_INFORMATION_EXTENSIONS.TagName(), "ServiceInformationExtensions"},
		{"SERVICE_TYPE_IDENTIFIER", TrustedListElement_SERVICE_TYPE_IDENTIFIER.TagName(), "ServiceTypeIdentifier"},
		{"SERVICE_NAME", TrustedListElement_SERVICE_NAME.TagName(), "ServiceName"},
		{"SERVICE_STATUS", TrustedListElement_SERVICE_STATUS.TagName(), "ServiceStatus"},
		{"SERVICE_SUPPLY_POINT", TrustedListElement_SERVICE_SUPPLY_POINT.TagName(), "ServiceSupplyPoint"},
		{"SERVICE_SUPPLY_POINTS", TrustedListElement_SERVICE_SUPPLY_POINTS.TagName(), "ServiceSupplyPoints"},
		{"STATE_OR_PROVINCE", TrustedListElement_STATE_OR_PROVINCE.TagName(), "StateOrProvince"},
		{"STATUS_DETERMINATION_APPROACH", TrustedListElement_STATUS_DETERMINATION_APPROACH.TagName(), "StatusDeterminationApproach"},
		{"STATUS_STARTING_TIME", TrustedListElement_STATUS_STARTING_TIME.TagName(), "StatusStartingTime"},
		{"STREET_ADDRESS", TrustedListElement_STREET_ADDRESS.TagName(), "StreetAddress"},
		{"TEXTUAL_INFORMATION", TrustedListElement_TEXTUAL_INFORMATION.TagName(), "TextualInformation"},
		{"TRUST_SERVICE_PROVIDER", TrustedListElement_TRUST_SERVICE_PROVIDER.TagName(), "TrustServiceProvider"},
		{"TRUST_SERVICE_PROVIDER_LIST", TrustedListElement_TRUST_SERVICE_PROVIDER_LIST.TagName(), "TrustServiceProviderList"},
		{"TRUST_SERVICE_STATUS_LIST", TrustedListElement_TRUST_SERVICE_STATUS_LIST.TagName(), "TrustServiceStatusList"},
		{"TSL_ADDRESS", TrustedListElement_TSL_ADDRESS.TagName(), "TSPAddress"},
		{"TSL_INFORMATION", TrustedListElement_TSL_INFORMATION.TagName(), "TSPInformation"},
		{"TSL_INFORMATION_EXTENSIONS", TrustedListElement_TSL_INFORMATION_EXTENSIONS.TagName(), "TSPInformationExtensions"},
		{"TSL_INFORMATION_URI", TrustedListElement_TSL_INFORMATION_URI.TagName(), "TSPInformationURI"},
		{"TSL_LEGAL_NOTICE", TrustedListElement_TSL_LEGAL_NOTICE.TagName(), "TSLLegalNotice"},
		{"TSL_LOCATION", TrustedListElement_TSL_LOCATION.TagName(), "TSLLocation"},
		{"TSL_NAME", TrustedListElement_TSL_NAME.TagName(), "TSPName"},
		{"TSL_POLICY", TrustedListElement_TSL_POLICY.TagName(), "TSLPolicy"},
		{"TSL_SEQUENCE_NUMBER", TrustedListElement_TSL_SEQUENCE_NUMBER.TagName(), "TSLSequenceNumber"},
		{"TSL_TYPE", TrustedListElement_TSL_TYPE.TagName(), "TSLType"},
		{"TSL_VERSION_IDENTIFIER", TrustedListElement_TSL_VERSION_IDENTIFIER.TagName(), "TSLVersionIdentifier"},
		{"TSP_SERVICE", TrustedListElement_TSP_SERVICE.TagName(), "TSPService"},
		{"TSP_SERVICES", TrustedListElement_TSP_SERVICES.TagName(), "TSPServices"},
		{"TSP_SERVICE_DEFINITION_URI", TrustedListElement_TSP_SERVICE_DEFINITION_URI.TagName(), "TSPServiceDefinitionURI"},
		{"TSP_TRADE_NAME", TrustedListElement_TSP_TRADE_NAME.TagName(), "TSPTradeName"},
		{"URI", TrustedListElement_URI.TagName(), "URI"},
		{"X509_CERTIFICATE", TrustedListElement_X509_CERTIFICATE.TagName(), "X509Certificate"},
		{"X509_SKI", TrustedListElement_X509_SKI.TagName(), "X509SKI"},
		{"X509_SUBJECT_NAME", TrustedListElement_X509_SUBJECT_NAME.TagName(), "X509SubjectName"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.TagName() = %q, want %q", c.name, c.got, c.want)
		}
	}
	if got := TrustedListNamespace_NS.Uri(); got != "http://uri.etsi.org/02231/v2#" {
		t.Errorf("TrustedListNamespace_NS.Uri() = %q", got)
	}
	if got := TrustedListNamespace_NS.Prefix(); got != "tl" {
		t.Errorf("TrustedListNamespace_NS.Prefix() = %q", got)
	}
}
