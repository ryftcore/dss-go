// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/definition/mra/MRAElement.java (DSS 6.5.RC1).
package tsl

import "github.com/utain/esig/dss/xml/common"

// MRAElement contains a list of MRA (Mutual Recognition Agreement) elements. Java enum ->
// typed string constants whose value is exactly Java's name(), per PORTING.md; the per-constant
// tag name is backed by the package-level lookup table below, following the precedent set by
// xades/definition's TrustedListElement.
type MRAElement string

const (
	MRAElement_CERTIFICATE_CONTENT_DECLARATION_POINTED_PARTY                 MRAElement = "CERTIFICATE_CONTENT_DECLARATION_POINTED_PARTY"
	MRAElement_CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY                MRAElement = "CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY"
	MRAElement_CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT             MRAElement = "CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT"
	MRAElement_CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE                    MRAElement = "CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE"
	MRAElement_CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_LIST               MRAElement = "CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_LIST"
	MRAElement_MUTUAL_RECOGNITION_AGREEMENT_INFORMATION                      MRAElement = "MUTUAL_RECOGNITION_AGREEMENT_INFORMATION"
	MRAElement_QC_CCLEGISLATION                                              MRAElement = "QC_CCLEGISLATION"
	MRAElement_QC_STATEMENT                                                  MRAElement = "QC_STATEMENT"
	MRAElement_QC_STATEMENT_ID                                               MRAElement = "QC_STATEMENT_ID"
	MRAElement_QC_STATEMENT_INFO                                             MRAElement = "QC_STATEMENT_INFO"
	MRAElement_QC_STATEMENT_SET                                              MRAElement = "QC_STATEMENT_SET"
	MRAElement_QC_TYPE                                                       MRAElement = "QC_TYPE"
	MRAElement_QUALIFIER_EQUIVALENCE                                         MRAElement = "QUALIFIER_EQUIVALENCE"
	MRAElement_QUALIFIER_EQUIVALENCE_LIST                                    MRAElement = "QUALIFIER_EQUIVALENCE_LIST"
	MRAElement_QUALIFIER_POINTED_PARTY                                       MRAElement = "QUALIFIER_POINTED_PARTY"
	MRAElement_QUALIFIER_POINTING_PARTY                                      MRAElement = "QUALIFIER_POINTING_PARTY"
	MRAElement_TRUST_SERVICE_EQUIVALENCE_HISTORY                             MRAElement = "TRUST_SERVICE_EQUIVALENCE_HISTORY"
	MRAElement_TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE                    MRAElement = "TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE"
	MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION                         MRAElement = "TRUST_SERVICE_EQUIVALENCE_INFORMATION"
	MRAElement_TRUST_SERVICE_EQUIVALENCE_STATUS                              MRAElement = "TRUST_SERVICE_EQUIVALENCE_STATUS"
	MRAElement_TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME                MRAElement = "TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME"
	MRAElement_TRUST_SERVICE_LEGAL_IDENTIFIER                                MRAElement = "TRUST_SERVICE_LEGAL_IDENTIFIER"
	MRAElement_TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_EQUIVALENCE_LIST    MRAElement = "TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_EQUIVALENCE_LIST"
	MRAElement_TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME                MRAElement = "TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME"
	MRAElement_TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTED_PARTY  MRAElement = "TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTED_PARTY"
	MRAElement_TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTING_PARTY MRAElement = "TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTING_PARTY"
	MRAElement_TRUST_SERVICE_TSL_STATUS_EQUIVALENCE_LIST                     MRAElement = "TRUST_SERVICE_TSL_STATUS_EQUIVALENCE_LIST"
	MRAElement_TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE                  MRAElement = "TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE"
	MRAElement_TRUST_SERVICE_TSL_STATUS_LIST_POINTED_PARTY                   MRAElement = "TRUST_SERVICE_TSL_STATUS_LIST_POINTED_PARTY"
	MRAElement_TRUST_SERVICE_TSL_STATUS_LIST_POINTING_PARTY                  MRAElement = "TRUST_SERVICE_TSL_STATUS_LIST_POINTING_PARTY"
	MRAElement_TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE                    MRAElement = "TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE"
	MRAElement_TRUST_SERVICE_TSL_TYPE                                        MRAElement = "TRUST_SERVICE_TSL_TYPE"
	MRAElement_TRUST_SERVICE_TSL_TYPE_EQUIVALENCE_LIST                       MRAElement = "TRUST_SERVICE_TSL_TYPE_EQUIVALENCE_LIST"
	MRAElement_TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY                     MRAElement = "TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY"
	MRAElement_TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY                    MRAElement = "TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY"
)

// mraElementTagNames maps each constant to its wire tag name (getTagName()).
var mraElementTagNames = map[MRAElement]string{
	MRAElement_CERTIFICATE_CONTENT_DECLARATION_POINTED_PARTY:                 "CertificateContentDeclarationPointedParty",
	MRAElement_CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY:                "CertificateContentDeclarationPointingParty",
	MRAElement_CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT:             "CertificateContentReferenceEquivalenceContext",
	MRAElement_CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE:                    "CertificateContentReferenceEquivalence",
	MRAElement_CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_LIST:               "CertificateContentReferencesEquivalenceList",
	MRAElement_MUTUAL_RECOGNITION_AGREEMENT_INFORMATION:                      "MutualRecognitionAgreementInformation",
	MRAElement_QC_CCLEGISLATION:                                              "QcCClegislation",
	MRAElement_QC_STATEMENT:                                                  "QcStatement",
	MRAElement_QC_STATEMENT_ID:                                               "QcStatementId",
	MRAElement_QC_STATEMENT_INFO:                                             "QcStatementInfo",
	MRAElement_QC_STATEMENT_SET:                                              "QcStatementSet",
	MRAElement_QC_TYPE:                                                       "QcType",
	MRAElement_QUALIFIER_EQUIVALENCE:                                         "QualifierEquivalence",
	MRAElement_QUALIFIER_EQUIVALENCE_LIST:                                    "QualifierEquivalenceList",
	MRAElement_QUALIFIER_POINTED_PARTY:                                       "QualifierPointedParty",
	MRAElement_QUALIFIER_POINTING_PARTY:                                      "QualifierPointingParty",
	MRAElement_TRUST_SERVICE_EQUIVALENCE_HISTORY:                             "TrustServiceEquivalenceHistory",
	MRAElement_TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE:                    "TrustServiceEquivalenceHistoryInstance",
	MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION:                         "TrustServiceEquivalenceInformation",
	MRAElement_TRUST_SERVICE_EQUIVALENCE_STATUS:                              "TrustServiceEquivalenceStatus",
	MRAElement_TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME:                "TrustServiceEquivalenceStatusStartingTime",
	MRAElement_TRUST_SERVICE_LEGAL_IDENTIFIER:                                "TrustServiceLegalIdentifier",
	MRAElement_TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_EQUIVALENCE_LIST:    "TrustServiceTSLQualificationExtensionEquivalenceList",
	MRAElement_TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME:                "TrustServiceTSLQualificationExtensionName",
	MRAElement_TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTED_PARTY:  "TrustServiceTSLQualificationExtensionNamePointedParty",
	MRAElement_TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTING_PARTY: "TrustServiceTSLQualificationExtensionNamePointingParty",
	MRAElement_TRUST_SERVICE_TSL_STATUS_EQUIVALENCE_LIST:                     "TrustServiceTSLStatusEquivalenceList",
	MRAElement_TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE:                  "TrustServiceTSLStatusInvalidEquivalence",
	MRAElement_TRUST_SERVICE_TSL_STATUS_LIST_POINTED_PARTY:                   "TrustServiceTSLStatusListPointedParty",
	MRAElement_TRUST_SERVICE_TSL_STATUS_LIST_POINTING_PARTY:                  "TrustServiceTSLStatusListPointingParty",
	MRAElement_TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE:                    "TrustServiceTSLStatusValidEquivalence",
	MRAElement_TRUST_SERVICE_TSL_TYPE:                                        "TrustServiceTSLType",
	MRAElement_TRUST_SERVICE_TSL_TYPE_EQUIVALENCE_LIST:                       "TrustServiceTSLTypeEquivalenceList",
	MRAElement_TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY:                     "TrustServiceTSLTypeListPointedParty",
	MRAElement_TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY:                    "TrustServiceTSLTypeListPointingParty",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e MRAElement) TagName() string {
	return mraElementTagNames[e]
}

// Namespace implements common.DSSElement. Ports getNamespace(): every MRAElement constant is
// built with MRANamespace.NS.
func (e MRAElement) Namespace() *common.DSSNamespace {
	return MRANamespace_NS
}

// URI implements common.DSSElement. Ports getURI().
func (e MRAElement) URI() string {
	return MRANamespace_NS.Uri()
}

// IsSameTagName implements common.DSSElement. Ports isSameTagName(String).
func (e MRAElement) IsSameTagName(value string) bool {
	return e.TagName() == value
}

var _ common.DSSElement = MRAElement("")
