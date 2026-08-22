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
		{"CERTIFICATE_CONTENT_DECLARATION_POINTED_PARTY", MRAElementCertificateContentDeclarationPointedParty, "CertificateContentDeclarationPointedParty"},
		{"CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY", MRAElementCertificateContentDeclarationPointingParty, "CertificateContentDeclarationPointingParty"},
		{"CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT", MRAElementCertificateContentReferenceEquivalenceContext, "CertificateContentReferenceEquivalenceContext"},
		{"CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE", MRAElementCertificateContentReferencesEquivalence, "CertificateContentReferenceEquivalence"},
		{"CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_LIST", MRAElementCertificateContentReferencesEquivalenceList, "CertificateContentReferencesEquivalenceList"},
		{"MUTUAL_RECOGNITION_AGREEMENT_INFORMATION", MRAElementMutualRecognitionAgreementInformation, "MutualRecognitionAgreementInformation"},
		{"QC_CCLEGISLATION", MRAElementQCCCLegislation, "QcCClegislation"},
		{"QC_STATEMENT", MRAElementQCStatement, "QcStatement"},
		{"QC_STATEMENT_ID", MRAElementQCStatementID, "QcStatementId"},
		{"QC_STATEMENT_INFO", MRAElementQCStatementInfo, "QcStatementInfo"},
		{"QC_STATEMENT_SET", MRAElementQCStatementSet, "QcStatementSet"},
		{"QC_TYPE", MRAElementQCType, "QcType"},
		{"QUALIFIER_EQUIVALENCE", MRAElementQualifierEquivalence, "QualifierEquivalence"},
		{"QUALIFIER_EQUIVALENCE_LIST", MRAElementQualifierEquivalenceList, "QualifierEquivalenceList"},
		{"QUALIFIER_POINTED_PARTY", MRAElementQualifierPointedParty, "QualifierPointedParty"},
		{"QUALIFIER_POINTING_PARTY", MRAElementQualifierPointingParty, "QualifierPointingParty"},
		{"TRUST_SERVICE_EQUIVALENCE_HISTORY", MRAElementTrustServiceEquivalenceHistory, "TrustServiceEquivalenceHistory"},
		{"TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE", MRAElementTrustServiceEquivalenceHistoryInstance, "TrustServiceEquivalenceHistoryInstance"},
		{"TRUST_SERVICE_EQUIVALENCE_INFORMATION", MRAElementTrustServiceEquivalenceInformation, "TrustServiceEquivalenceInformation"},
		{"TRUST_SERVICE_EQUIVALENCE_STATUS", MRAElementTrustServiceEquivalenceStatus, "TrustServiceEquivalenceStatus"},
		{"TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME", MRAElementTrustServiceEquivalenceStatusStartingTime, "TrustServiceEquivalenceStatusStartingTime"},
		{"TRUST_SERVICE_LEGAL_IDENTIFIER", MRAElementTrustServiceLegalIdentifier, "TrustServiceLegalIdentifier"},
		{"TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_EQUIVALENCE_LIST", MRAElementTrustServiceTSLQualificationExtensionEquivalenceList, "TrustServiceTSLQualificationExtensionEquivalenceList"},
		{"TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME", MRAElementTrustServiceTSLQualificationExtensionName, "TrustServiceTSLQualificationExtensionName"},
		{"TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTED_PARTY", MRAElementTrustServiceTSLQualificationExtensionNamePointedParty, "TrustServiceTSLQualificationExtensionNamePointedParty"},
		{"TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTING_PARTY", MRAElementTrustServiceTSLQualificationExtensionNamePointingParty, "TrustServiceTSLQualificationExtensionNamePointingParty"},
		{"TRUST_SERVICE_TSL_STATUS_EQUIVALENCE_LIST", MRAElementTrustServiceTSLStatusEquivalenceList, "TrustServiceTSLStatusEquivalenceList"},
		{"TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE", MRAElementTrustServiceTSLStatusInvalidEquivalence, "TrustServiceTSLStatusInvalidEquivalence"},
		{"TRUST_SERVICE_TSL_STATUS_LIST_POINTED_PARTY", MRAElementTrustServiceTSLStatusListPointedParty, "TrustServiceTSLStatusListPointedParty"},
		{"TRUST_SERVICE_TSL_STATUS_LIST_POINTING_PARTY", MRAElementTrustServiceTSLStatusListPointingParty, "TrustServiceTSLStatusListPointingParty"},
		{"TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE", MRAElementTrustServiceTSLStatusValidEquivalence, "TrustServiceTSLStatusValidEquivalence"},
		{"TRUST_SERVICE_TSL_TYPE", MRAElementTrustServiceTSLType, "TrustServiceTSLType"},
		{"TRUST_SERVICE_TSL_TYPE_EQUIVALENCE_LIST", MRAElementTrustServiceTSLTypeEquivalenceList, "TrustServiceTSLTypeEquivalenceList"},
		{"TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY", MRAElementTrustServiceTSLTypeListPointedParty, "TrustServiceTSLTypeListPointedParty"},
		{"TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY", MRAElementTrustServiceTSLTypeListPointingParty, "TrustServiceTSLTypeListPointingParty"},
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
		if c.got.Namespace() != MRANamespaceNS {
			t.Errorf("%s.Namespace() is not MRANamespaceNS", c.name)
		}
	}
}

func TestMRANamespace_KAT(t *testing.T) {
	if got := MRANamespaceNS.Uri(); got != "http://ec.europa.eu/tools/lotl/mra/schema/v2#" {
		t.Errorf("MRANamespaceNS.Uri() = %q", got)
	}
	if got := MRANamespaceNS.Prefix(); got != "mra" {
		t.Errorf("MRANamespaceNS.Prefix() = %q", got)
	}
}
