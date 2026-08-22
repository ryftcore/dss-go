// KAT test for mra_path.go: every MRAPath_* query string below is the expression the upstream
// Java field builds (dss-tsl-validation 6.5.RC1
// eu.europa.esig.dss.tsl.definition.mra.MRAPath), resolved through the same
// AbstractPath/XPathQueryBuilder rules the xades/definition TrustedListPath KAT relies on.
package tsl

import "testing"

func TestMRAPath_KAT(t *testing.T) {
	const eqInfo = "./mra:TrustServiceEquivalenceInformation"
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"CERTIFICATE_CONTENT_DECLARATION_POINTED_PARTY_PATH",
			MRAPathCertificateContentDeclarationPointedPartyPath.QueryString(),
			"./mra:CertificateContentDeclarationPointedParty"},
		{"CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY_PATH",
			MRAPathCertificateContentDeclarationPointingPartyPath.QueryString(),
			"./mra:CertificateContentDeclarationPointingParty"},
		{"CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_PATH",
			MRAPathCertificateContentReferencesEquivalencePath.QueryString(),
			eqInfo + "/mra:CertificateContentReferencesEquivalenceList/mra:CertificateContentReferenceEquivalence"},
		{"CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT_PATH",
			MRAPathCertificateContentReferenceEquivalenceContextPath.QueryString(),
			"./mra:CertificateContentReferenceEquivalenceContext"},
		{"MUTUAL_RECOGNITION_AGREEMENT_INFORMATION_PATH",
			MRAPathMutualRecognitionAgreementInformationPath.QueryString(),
			"//mra:MutualRecognitionAgreementInformation"},
		{"QUALIFIER_EQUIVALENCE_LIST_PATH",
			MRAPathQualifierEquivalenceListPath.QueryString(),
			eqInfo + "/mra:TrustServiceTSLQualificationExtensionEquivalenceList/mra:QualifierEquivalenceList"},
		{"SERVICE_TYPE_IDENTIFIER_PATH",
			MRAPathServiceTypeIdentifierPath.QueryString(),
			"./mra:TrustServiceTSLType/tl:ServiceTypeIdentifier"},
		{"TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE_PATH",
			MRAPathTrustServiceEquivalenceHistoryInstancePath.QueryString(),
			eqInfo + "/mra:TrustServiceEquivalenceHistory/mra:TrustServiceEquivalenceHistoryInstance"},
		{"TRUST_SERVICE_EQUIVALENCE_INFORMATION_STATUS_PATH",
			MRAPathTrustServiceEquivalenceInformationStatusPath.QueryString(),
			eqInfo + "/mra:TrustServiceEquivalenceStatus"},
		{"TRUST_SERVICE_EQUIVALENCE_INFORMATION_STATUS_STARTING_TIME_PATH",
			MRAPathTrustServiceEquivalenceInformationStatusStartingTimePath.QueryString(),
			eqInfo + "/mra:TrustServiceEquivalenceStatusStartingTime"},
		{"TRUST_SERVICE_EQUIVALENCE_STATUS_PATH",
			MRAPathTrustServiceEquivalenceStatusPath.QueryString(),
			"./mra:TrustServiceEquivalenceStatus"},
		{"TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME_PATH",
			MRAPathTrustServiceEquivalenceStatusStartingTimePath.QueryString(),
			"./mra:TrustServiceEquivalenceStatusStartingTime"},
		{"TRUST_SERVICE_LEGAL_IDENTIFIER_PATH",
			MRAPathTrustServiceLegalIdentifierPath.QueryString(),
			eqInfo + "/mra:TrustServiceLegalIdentifier"},
		{"TRUST_SERVICE_TSL_TYPE_PATH",
			MRAPathTrustServiceTSLTypePath.QueryString(),
			"./mra:TrustServiceTSLType"},
		{"TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY_PATH",
			MRAPathTrustServiceTSLTypeListPointedPartyPath.QueryString(),
			eqInfo + "/mra:TrustServiceTSLTypeEquivalenceList/mra:TrustServiceTSLTypeListPointedParty"},
		{"TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY_PATH",
			MRAPathTrustServiceTSLTypeListPointingPartyPath.QueryString(),
			eqInfo + "/mra:TrustServiceTSLTypeEquivalenceList/mra:TrustServiceTSLTypeListPointingParty"},
		{"TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE_LIST_POINTED_PARTY_SERVICE_STATUS_PATH",
			MRAPathTrustServiceTSLStatusInvalidEquivalenceListPointedPartyServiceStatusPath.QueryString(),
			eqInfo + "/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusInvalidEquivalence/mra:TrustServiceTSLStatusListPointedParty/tl:ServiceStatus"},
		{"TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE_LIST_POINTING_PARTY_SERVICE_STATUS_PATH",
			MRAPathTrustServiceTSLStatusInvalidEquivalenceListPointingPartyServiceStatusPath.QueryString(),
			eqInfo + "/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusInvalidEquivalence/mra:TrustServiceTSLStatusListPointingParty/tl:ServiceStatus"},
		{"TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE_LIST_POINTED_PARTY_SERVICE_STATUS_PATH",
			MRAPathTrustServiceTSLStatusValidEquivalenceListPointedPartyServiceStatusPath.QueryString(),
			eqInfo + "/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusValidEquivalence/mra:TrustServiceTSLStatusListPointedParty/tl:ServiceStatus"},
		{"TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE_LIST_POINTING_PARTY_SERVICE_STATUS_PATH",
			MRAPathTrustServiceTSLStatusValidEquivalenceListPointingPartyServiceStatusPath.QueryString(),
			eqInfo + "/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusValidEquivalence/mra:TrustServiceTSLStatusListPointingParty/tl:ServiceStatus"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.QueryString() = %q, want %q", c.name, c.got, c.want)
		}
	}
}
