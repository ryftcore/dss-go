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
		{"ADDITIONAL_INFORMATION", TrustedListElementAdditionalInformation.TagName(), "AdditionalInformation"},
		{"ADDITIONAL_SERVICE_INFORMATION", TrustedListElementAdditionalServiceInformation.TagName(), "AdditionalServiceInformation"},
		{"COUNTRY_NAME", TrustedListElementCountryName.TagName(), "CountryName"},
		{"DATE_TIME", TrustedListElementDateTime.TagName(), "dateTime"},
		{"DIGITAL_ID", TrustedListElementDigitalID.TagName(), "DigitalId"},
		{"DISTRIBUTION_POINTS", TrustedListElementDistributionPoints.TagName(), "DistributionPoints"},
		{"ELECTRONIC_ADDRESS", TrustedListElementElectronicAddress.TagName(), "ElectronicAddress"},
		{"EXPIRED_CERTS_REVOCATION_INFO", TrustedListElementExpiredCertsRevocationInfo.TagName(), "ExpiredCertsRevocationInfo"},
		{"EXTENSION", TrustedListElementExtension.TagName(), "Extension"},
		{"HISTORICAL_INFORMATION_PERIOD", TrustedListElementHistoricalInformationPeriod.TagName(), "HistoricalInformationPeriod"},
		{"INFORMATION_VALUE", TrustedListElementInformationValue.TagName(), "InformationValue"},
		{"LIST_ISSUE_DATE_TIME", TrustedListElementListIssueDateTime.TagName(), "ListIssueDateTime"},
		{"LOCALITY", TrustedListElementLocality.TagName(), "Locality"},
		{"NAME", TrustedListElementName.TagName(), "Name"},
		{"NEXT_UPDATE", TrustedListElementNextUpdate.TagName(), "NextUpdate"},
		{"OTHER", TrustedListElementOther.TagName(), "Other"},
		{"OTHER_INFORMATION", TrustedListElementOtherInformation.TagName(), "OtherInformation"},
		{"OTHER_TSL_POINTER", TrustedListElementOtherTSLPointer.TagName(), "OtherTSLPointer"},
		{"POINTERS_TO_OTHER_TSL", TrustedListElementPointersToOtherTSL.TagName(), "PointersToOtherTSL"},
		{"POLICY_OR_LEGAL_NOTICE", TrustedListElementPolicyOrLegalNotice.TagName(), "PolicyOrLegalNotice"},
		{"POSTAL_ADDRESS", TrustedListElementPostalAddress.TagName(), "PostalAddress"},
		{"POSTAL_ADDRESSES", TrustedListElementPostalAddresses.TagName(), "PostalAddresses"},
		{"POSTAL_CODE", TrustedListElementPostalCode.TagName(), "PostalCode"},
		{"SCHEME_EXTENSION", TrustedListElementSchemeExtension.TagName(), "SchemeExtensions"},
		{"SCHEME_INFORMATION", TrustedListElementSchemeInformation.TagName(), "SchemeInformation"},
		{"SCHEME_INFORMATION_URI", TrustedListElementSchemeInformationURI.TagName(), "SchemeInformationURI"},
		{"SCHEME_NAME", TrustedListElementSchemeName.TagName(), "SchemeName"},
		{"SCHEME_OPERATOR_ADDRESS", TrustedListElementSchemeOperatorAddress.TagName(), "SchemeOperatorAddress"},
		{"SCHEME_OPERATOR_NAME", TrustedListElementSchemeOperatorName.TagName(), "SchemeOperatorName"},
		{"SCHEME_SERVICE_DEFINITION_URI", TrustedListElementSchemeServiceDefinitionURI.TagName(), "SchemeServiceDefinitionURI"},
		{"SCHEME_TERRITORY", TrustedListElementSchemeTerritory.TagName(), "SchemeTerritory"},
		{"SCHEME_TYPE_COMMUNITY_RULES", TrustedListElementSchemeTypeCommunityRules.TagName(), "SchemeTypeCommunityRules"},
		{"SERVICE_DIGITAL_IDENTITY", TrustedListElementServiceDigitalIdentity.TagName(), "ServiceDigitalIdentity"},
		{"SERVICE_DIGITAL_IDENTITIES", TrustedListElementServiceDigitalIdentities.TagName(), "ServiceDigitalIdentities"},
		{"SERVICE_HISTORY", TrustedListElementServiceHistory.TagName(), "ServiceHistory"},
		{"SERVICE_HISTORY_INSTANCE", TrustedListElementServiceHistoryInstance.TagName(), "ServiceHistoryInstance"},
		{"SERVICE_INFORMATION", TrustedListElementServiceInformation.TagName(), "ServiceInformation"},
		{"SERVICE_INFORMATION_EXTENSIONS", TrustedListElementServiceInformationExtensions.TagName(), "ServiceInformationExtensions"},
		{"SERVICE_TYPE_IDENTIFIER", TrustedListElementServiceTypeIdentifier.TagName(), "ServiceTypeIdentifier"},
		{"SERVICE_NAME", TrustedListElementServiceName.TagName(), "ServiceName"},
		{"SERVICE_STATUS", TrustedListElementServiceStatus.TagName(), "ServiceStatus"},
		{"SERVICE_SUPPLY_POINT", TrustedListElementServiceSupplyPoint.TagName(), "ServiceSupplyPoint"},
		{"SERVICE_SUPPLY_POINTS", TrustedListElementServiceSupplyPoints.TagName(), "ServiceSupplyPoints"},
		{"STATE_OR_PROVINCE", TrustedListElementStateOrProvince.TagName(), "StateOrProvince"},
		{"STATUS_DETERMINATION_APPROACH", TrustedListElementStatusDeterminationApproach.TagName(), "StatusDeterminationApproach"},
		{"STATUS_STARTING_TIME", TrustedListElementStatusStartingTime.TagName(), "StatusStartingTime"},
		{"STREET_ADDRESS", TrustedListElementStreetAddress.TagName(), "StreetAddress"},
		{"TEXTUAL_INFORMATION", TrustedListElementTextualInformation.TagName(), "TextualInformation"},
		{"TRUST_SERVICE_PROVIDER", TrustedListElementTrustServiceProvider.TagName(), "TrustServiceProvider"},
		{"TRUST_SERVICE_PROVIDER_LIST", TrustedListElementTrustServiceProviderList.TagName(), "TrustServiceProviderList"},
		{"TRUST_SERVICE_STATUS_LIST", TrustedListElementTrustServiceStatusList.TagName(), "TrustServiceStatusList"},
		{"TSL_ADDRESS", TrustedListElementTSLAddress.TagName(), "TSPAddress"},
		{"TSL_INFORMATION", TrustedListElementTSLInformation.TagName(), "TSPInformation"},
		{"TSL_INFORMATION_EXTENSIONS", TrustedListElementTSLInformationExtensions.TagName(), "TSPInformationExtensions"},
		{"TSL_INFORMATION_URI", TrustedListElementTSLInformationURI.TagName(), "TSPInformationURI"},
		{"TSL_LEGAL_NOTICE", TrustedListElementTSLLegalNotice.TagName(), "TSLLegalNotice"},
		{"TSL_LOCATION", TrustedListElementTSLLocation.TagName(), "TSLLocation"},
		{"TSL_NAME", TrustedListElementTSLName.TagName(), "TSPName"},
		{"TSL_POLICY", TrustedListElementTSLPolicy.TagName(), "TSLPolicy"},
		{"TSL_SEQUENCE_NUMBER", TrustedListElementTSLSequenceNumber.TagName(), "TSLSequenceNumber"},
		{"TSL_TYPE", TrustedListElementTSLType.TagName(), "TSLType"},
		{"TSL_VERSION_IDENTIFIER", TrustedListElementTSLVersionIdentifier.TagName(), "TSLVersionIdentifier"},
		{"TSP_SERVICE", TrustedListElementTSPService.TagName(), "TSPService"},
		{"TSP_SERVICES", TrustedListElementTSPServices.TagName(), "TSPServices"},
		{"TSP_SERVICE_DEFINITION_URI", TrustedListElementTSPServiceDefinitionURI.TagName(), "TSPServiceDefinitionURI"},
		{"TSP_TRADE_NAME", TrustedListElementTSPTradeName.TagName(), "TSPTradeName"},
		{"URI", TrustedListElementURI.TagName(), "URI"},
		{"X509_CERTIFICATE", TrustedListElementX509Certificate.TagName(), "X509Certificate"},
		{"X509_SKI", TrustedListElementX509SKI.TagName(), "X509SKI"},
		{"X509_SUBJECT_NAME", TrustedListElementX509SubjectName.TagName(), "X509SubjectName"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.TagName() = %q, want %q", c.name, c.got, c.want)
		}
	}
	if got := TrustedListNamespaceNS.Uri(); got != "http://uri.etsi.org/02231/v2#" {
		t.Errorf("TrustedListNamespaceNS.Uri() = %q", got)
	}
	if got := TrustedListNamespaceNS.Prefix(); got != "tl" {
		t.Errorf("TrustedListNamespaceNS.Prefix() = %q", got)
	}
}
