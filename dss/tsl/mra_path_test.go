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
			MRAPath_CERTIFICATE_CONTENT_DECLARATION_POINTED_PARTY_PATH.QueryString(),
			"./mra:CertificateContentDeclarationPointedParty"},
		{"CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY_PATH",
			MRAPath_CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY_PATH.QueryString(),
			"./mra:CertificateContentDeclarationPointingParty"},
		{"CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_PATH",
			MRAPath_CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_PATH.QueryString(),
			eqInfo + "/mra:CertificateContentReferencesEquivalenceList/mra:CertificateContentReferenceEquivalence"},
		{"CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT_PATH",
			MRAPath_CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT_PATH.QueryString(),
			"./mra:CertificateContentReferenceEquivalenceContext"},
		{"MUTUAL_RECOGNITION_AGREEMENT_INFORMATION_PATH",
			MRAPath_MUTUAL_RECOGNITION_AGREEMENT_INFORMATION_PATH.QueryString(),
			"//mra:MutualRecognitionAgreementInformation"},
		{"QUALIFIER_EQUIVALENCE_LIST_PATH",
			MRAPath_QUALIFIER_EQUIVALENCE_LIST_PATH.QueryString(),
			eqInfo + "/mra:TrustServiceTSLQualificationExtensionEquivalenceList/mra:QualifierEquivalenceList"},
		{"SERVICE_TYPE_IDENTIFIER_PATH",
			MRAPath_SERVICE_TYPE_IDENTIFIER_PATH.QueryString(),
			"./mra:TrustServiceTSLType/tl:ServiceTypeIdentifier"},
		{"TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE_PATH",
			MRAPath_TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE_PATH.QueryString(),
			eqInfo + "/mra:TrustServiceEquivalenceHistory/mra:TrustServiceEquivalenceHistoryInstance"},
		{"TRUST_SERVICE_EQUIVALENCE_INFORMATION_STATUS_PATH",
			MRAPath_TRUST_SERVICE_EQUIVALENCE_INFORMATION_STATUS_PATH.QueryString(),
			eqInfo + "/mra:TrustServiceEquivalenceStatus"},
		{"TRUST_SERVICE_EQUIVALENCE_INFORMATION_STATUS_STARTING_TIME_PATH",
			MRAPath_TRUST_SERVICE_EQUIVALENCE_INFORMATION_STATUS_STARTING_TIME_PATH.QueryString(),
			eqInfo + "/mra:TrustServiceEquivalenceStatusStartingTime"},
		{"TRUST_SERVICE_EQUIVALENCE_STATUS_PATH",
			MRAPath_TRUST_SERVICE_EQUIVALENCE_STATUS_PATH.QueryString(),
			"./mra:TrustServiceEquivalenceStatus"},
		{"TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME_PATH",
			MRAPath_TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME_PATH.QueryString(),
			"./mra:TrustServiceEquivalenceStatusStartingTime"},
		{"TRUST_SERVICE_LEGAL_IDENTIFIER_PATH",
			MRAPath_TRUST_SERVICE_LEGAL_IDENTIFIER_PATH.QueryString(),
			eqInfo + "/mra:TrustServiceLegalIdentifier"},
		{"TRUST_SERVICE_TSL_TYPE_PATH",
			MRAPath_TRUST_SERVICE_TSL_TYPE_PATH.QueryString(),
			"./mra:TrustServiceTSLType"},
		{"TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY_PATH",
			MRAPath_TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY_PATH.QueryString(),
			eqInfo + "/mra:TrustServiceTSLTypeEquivalenceList/mra:TrustServiceTSLTypeListPointedParty"},
		{"TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY_PATH",
			MRAPath_TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY_PATH.QueryString(),
			eqInfo + "/mra:TrustServiceTSLTypeEquivalenceList/mra:TrustServiceTSLTypeListPointingParty"},
		{"TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE_LIST_POINTED_PARTY_SERVICE_STATUS_PATH",
			MRAPath_TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE_LIST_POINTED_PARTY_SERVICE_STATUS_PATH.QueryString(),
			eqInfo + "/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusInvalidEquivalence/mra:TrustServiceTSLStatusListPointedParty/tl:ServiceStatus"},
		{"TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE_LIST_POINTING_PARTY_SERVICE_STATUS_PATH",
			MRAPath_TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE_LIST_POINTING_PARTY_SERVICE_STATUS_PATH.QueryString(),
			eqInfo + "/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusInvalidEquivalence/mra:TrustServiceTSLStatusListPointingParty/tl:ServiceStatus"},
		{"TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE_LIST_POINTED_PARTY_SERVICE_STATUS_PATH",
			MRAPath_TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE_LIST_POINTED_PARTY_SERVICE_STATUS_PATH.QueryString(),
			eqInfo + "/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusValidEquivalence/mra:TrustServiceTSLStatusListPointedParty/tl:ServiceStatus"},
		{"TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE_LIST_POINTING_PARTY_SERVICE_STATUS_PATH",
			MRAPath_TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE_LIST_POINTING_PARTY_SERVICE_STATUS_PATH.QueryString(),
			eqInfo + "/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusValidEquivalence/mra:TrustServiceTSLStatusListPointingParty/tl:ServiceStatus"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.QueryString() = %q, want %q", c.name, c.got, c.want)
		}
	}
}
