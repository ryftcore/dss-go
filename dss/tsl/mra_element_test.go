// KAT test for mra_element.go / mra_namespace.go: every MRAElement_* tag name below is
// transcribed verbatim from the upstream Java source (dss-tsl-validation 6.5.RC1)
// eu.europa.esig.dss.tsl.definition.mra.MRAElement, and the namespace URI/prefix from
// eu.europa.esig.dss.tsl.definition.mra.MRANamespace.
package tsl

import "testing"

func TestMRAElement_KAT(t *testing.T) {
	cases := []struct {
		name string
		got  MRAElement
		want string
	}{
		{"CERTIFICATE_CONTENT_DECLARATION_POINTED_PARTY", MRAElement_CERTIFICATE_CONTENT_DECLARATION_POINTED_PARTY, "CertificateContentDeclarationPointedParty"},
		{"CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY", MRAElement_CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY, "CertificateContentDeclarationPointingParty"},
		{"CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT", MRAElement_CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT, "CertificateContentReferenceEquivalenceContext"},
		{"CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE", MRAElement_CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE, "CertificateContentReferenceEquivalence"},
		{"CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_LIST", MRAElement_CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_LIST, "CertificateContentReferencesEquivalenceList"},
		{"MUTUAL_RECOGNITION_AGREEMENT_INFORMATION", MRAElement_MUTUAL_RECOGNITION_AGREEMENT_INFORMATION, "MutualRecognitionAgreementInformation"},
		{"QC_CCLEGISLATION", MRAElement_QC_CCLEGISLATION, "QcCClegislation"},
		{"QC_STATEMENT", MRAElement_QC_STATEMENT, "QcStatement"},
		{"QC_STATEMENT_ID", MRAElement_QC_STATEMENT_ID, "QcStatementId"},
		{"QC_STATEMENT_INFO", MRAElement_QC_STATEMENT_INFO, "QcStatementInfo"},
		{"QC_STATEMENT_SET", MRAElement_QC_STATEMENT_SET, "QcStatementSet"},
		{"QC_TYPE", MRAElement_QC_TYPE, "QcType"},
		{"QUALIFIER_EQUIVALENCE", MRAElement_QUALIFIER_EQUIVALENCE, "QualifierEquivalence"},
		{"QUALIFIER_EQUIVALENCE_LIST", MRAElement_QUALIFIER_EQUIVALENCE_LIST, "QualifierEquivalenceList"},
		{"QUALIFIER_POINTED_PARTY", MRAElement_QUALIFIER_POINTED_PARTY, "QualifierPointedParty"},
		{"QUALIFIER_POINTING_PARTY", MRAElement_QUALIFIER_POINTING_PARTY, "QualifierPointingParty"},
		{"TRUST_SERVICE_EQUIVALENCE_HISTORY", MRAElement_TRUST_SERVICE_EQUIVALENCE_HISTORY, "TrustServiceEquivalenceHistory"},
		{"TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE", MRAElement_TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE, "TrustServiceEquivalenceHistoryInstance"},
		{"TRUST_SERVICE_EQUIVALENCE_INFORMATION", MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION, "TrustServiceEquivalenceInformation"},
		{"TRUST_SERVICE_EQUIVALENCE_STATUS", MRAElement_TRUST_SERVICE_EQUIVALENCE_STATUS, "TrustServiceEquivalenceStatus"},
		{"TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME", MRAElement_TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME, "TrustServiceEquivalenceStatusStartingTime"},
		{"TRUST_SERVICE_LEGAL_IDENTIFIER", MRAElement_TRUST_SERVICE_LEGAL_IDENTIFIER, "TrustServiceLegalIdentifier"},
		{"TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_EQUIVALENCE_LIST", MRAElement_TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_EQUIVALENCE_LIST, "TrustServiceTSLQualificationExtensionEquivalenceList"},
		{"TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME", MRAElement_TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME, "TrustServiceTSLQualificationExtensionName"},
		{"TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTED_PARTY", MRAElement_TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTED_PARTY, "TrustServiceTSLQualificationExtensionNamePointedParty"},
		{"TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTING_PARTY", MRAElement_TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTING_PARTY, "TrustServiceTSLQualificationExtensionNamePointingParty"},
		{"TRUST_SERVICE_TSL_STATUS_EQUIVALENCE_LIST", MRAElement_TRUST_SERVICE_TSL_STATUS_EQUIVALENCE_LIST, "TrustServiceTSLStatusEquivalenceList"},
		{"TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE", MRAElement_TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE, "TrustServiceTSLStatusInvalidEquivalence"},
		{"TRUST_SERVICE_TSL_STATUS_LIST_POINTED_PARTY", MRAElement_TRUST_SERVICE_TSL_STATUS_LIST_POINTED_PARTY, "TrustServiceTSLStatusListPointedParty"},
		{"TRUST_SERVICE_TSL_STATUS_LIST_POINTING_PARTY", MRAElement_TRUST_SERVICE_TSL_STATUS_LIST_POINTING_PARTY, "TrustServiceTSLStatusListPointingParty"},
		{"TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE", MRAElement_TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE, "TrustServiceTSLStatusValidEquivalence"},
		{"TRUST_SERVICE_TSL_TYPE", MRAElement_TRUST_SERVICE_TSL_TYPE, "TrustServiceTSLType"},
		{"TRUST_SERVICE_TSL_TYPE_EQUIVALENCE_LIST", MRAElement_TRUST_SERVICE_TSL_TYPE_EQUIVALENCE_LIST, "TrustServiceTSLTypeEquivalenceList"},
		{"TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY", MRAElement_TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY, "TrustServiceTSLTypeListPointedParty"},
		{"TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY", MRAElement_TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY, "TrustServiceTSLTypeListPointingParty"},
	}
	if len(cases) != len(mraElementTagNames) {
		t.Fatalf("KAT covers %d constants, the tag-name table holds %d", len(cases), len(mraElementTagNames))
	}
	for _, c := range cases {
		if got := c.got.TagName(); got != c.want {
			t.Errorf("%s.TagName() = %q, want %q", c.name, got, c.want)
		}
		if string(c.got) != c.name {
			t.Errorf("constant value = %q, want the Java name() %q", string(c.got), c.name)
		}
		if !c.got.IsSameTagName(c.want) {
			t.Errorf("%s.IsSameTagName(%q) = false, want true", c.name, c.want)
		}
		if c.got.IsSameTagName(c.want + "X") {
			t.Errorf("%s.IsSameTagName(%q) = true, want false", c.name, c.want+"X")
		}
		if got := c.got.URI(); got != "http://ec.europa.eu/tools/lotl/mra/schema/v2#" {
			t.Errorf("%s.URI() = %q", c.name, got)
		}
		if c.got.Namespace() != MRANamespace_NS {
			t.Errorf("%s.Namespace() is not MRANamespace_NS", c.name)
		}
	}
}

func TestMRANamespace_KAT(t *testing.T) {
	if got := MRANamespace_NS.Uri(); got != "http://ec.europa.eu/tools/lotl/mra/schema/v2#" {
		t.Errorf("MRANamespace_NS.Uri() = %q", got)
	}
	if got := MRANamespace_NS.Prefix(); got != "mra" {
		t.Errorf("MRANamespace_NS.Prefix() = %q", got)
	}
}
