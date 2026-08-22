// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/definition/mra/MRAElement.java (DSS 6.5.RC1).
package tsl

import "github.com/ryftcore/dss-go/dss/xml/common"

// MRAElement contains a list of MRA (Mutual Recognition Agreement) elements. Java enum ->
// typed string constants whose value is exactly Java's name(), per PORTING.md; the per-constant
// tag name is backed by the package-level lookup table below, following the precedent set by
// xades/definition's TrustedListElement.
type MRAElement string

const (
	MRAElementCertificateContentDeclarationPointedParty              MRAElement = "CERTIFICATE_CONTENT_DECLARATION_POINTED_PARTY"
	MRAElementCertificateContentDeclarationPointingParty             MRAElement = "CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY"
	MRAElementCertificateContentReferenceEquivalenceContext          MRAElement = "CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT"
	MRAElementCertificateContentReferencesEquivalence                MRAElement = "CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE"
	MRAElementCertificateContentReferencesEquivalenceList            MRAElement = "CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_LIST"
	MRAElementMutualRecognitionAgreementInformation                  MRAElement = "MUTUAL_RECOGNITION_AGREEMENT_INFORMATION"
	MRAElementQCCCLegislation                                        MRAElement = "QC_CCLEGISLATION"
	MRAElementQCStatement                                            MRAElement = "QC_STATEMENT"
	MRAElementQCStatementID                                          MRAElement = "QC_STATEMENT_ID"
	MRAElementQCStatementInfo                                        MRAElement = "QC_STATEMENT_INFO"
	MRAElementQCStatementSet                                         MRAElement = "QC_STATEMENT_SET"
	MRAElementQCType                                                 MRAElement = "QC_TYPE"
	MRAElementQualifierEquivalence                                   MRAElement = "QUALIFIER_EQUIVALENCE"
	MRAElementQualifierEquivalenceList                               MRAElement = "QUALIFIER_EQUIVALENCE_LIST"
	MRAElementQualifierPointedParty                                  MRAElement = "QUALIFIER_POINTED_PARTY"
	MRAElementQualifierPointingParty                                 MRAElement = "QUALIFIER_POINTING_PARTY"
	MRAElementTrustServiceEquivalenceHistory                         MRAElement = "TRUST_SERVICE_EQUIVALENCE_HISTORY"
	MRAElementTrustServiceEquivalenceHistoryInstance                 MRAElement = "TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE"
	MRAElementTrustServiceEquivalenceInformation                     MRAElement = "TRUST_SERVICE_EQUIVALENCE_INFORMATION"
	MRAElementTrustServiceEquivalenceStatus                          MRAElement = "TRUST_SERVICE_EQUIVALENCE_STATUS"
	MRAElementTrustServiceEquivalenceStatusStartingTime              MRAElement = "TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME"
	MRAElementTrustServiceLegalIdentifier                            MRAElement = "TRUST_SERVICE_LEGAL_IDENTIFIER"
	MRAElementTrustServiceTSLQualificationExtensionEquivalenceList   MRAElement = "TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_EQUIVALENCE_LIST"
	MRAElementTrustServiceTSLQualificationExtensionName              MRAElement = "TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME"
	MRAElementTrustServiceTSLQualificationExtensionNamePointedParty  MRAElement = "TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTED_PARTY"
	MRAElementTrustServiceTSLQualificationExtensionNamePointingParty MRAElement = "TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_NAME_POINTING_PARTY"
	MRAElementTrustServiceTSLStatusEquivalenceList                   MRAElement = "TRUST_SERVICE_TSL_STATUS_EQUIVALENCE_LIST"
	MRAElementTrustServiceTSLStatusInvalidEquivalence                MRAElement = "TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE"
	MRAElementTrustServiceTSLStatusListPointedParty                  MRAElement = "TRUST_SERVICE_TSL_STATUS_LIST_POINTED_PARTY"
	MRAElementTrustServiceTSLStatusListPointingParty                 MRAElement = "TRUST_SERVICE_TSL_STATUS_LIST_POINTING_PARTY"
	MRAElementTrustServiceTSLStatusValidEquivalence                  MRAElement = "TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE"
	MRAElementTrustServiceTSLType                                    MRAElement = "TRUST_SERVICE_TSL_TYPE"
	MRAElementTrustServiceTSLTypeEquivalenceList                     MRAElement = "TRUST_SERVICE_TSL_TYPE_EQUIVALENCE_LIST"
	MRAElementTrustServiceTSLTypeListPointedParty                    MRAElement = "TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY"
	MRAElementTrustServiceTSLTypeListPointingParty                   MRAElement = "TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY"
)

// mraElementTagNames maps each constant to its wire tag name (getTagName()).
var mraElementTagNames = map[MRAElement]string{
	MRAElementCertificateContentDeclarationPointedParty:              "CertificateContentDeclarationPointedParty",
	MRAElementCertificateContentDeclarationPointingParty:             "CertificateContentDeclarationPointingParty",
	MRAElementCertificateContentReferenceEquivalenceContext:          "CertificateContentReferenceEquivalenceContext",
	MRAElementCertificateContentReferencesEquivalence:                "CertificateContentReferenceEquivalence",
	MRAElementCertificateContentReferencesEquivalenceList:            "CertificateContentReferencesEquivalenceList",
	MRAElementMutualRecognitionAgreementInformation:                  "MutualRecognitionAgreementInformation",
	MRAElementQCCCLegislation:                                        "QcCClegislation",
	MRAElementQCStatement:                                            "QcStatement",
	MRAElementQCStatementID:                                          "QcStatementId",
	MRAElementQCStatementInfo:                                        "QcStatementInfo",
	MRAElementQCStatementSet:                                         "QcStatementSet",
	MRAElementQCType:                                                 "QcType",
	MRAElementQualifierEquivalence:                                   "QualifierEquivalence",
	MRAElementQualifierEquivalenceList:                               "QualifierEquivalenceList",
	MRAElementQualifierPointedParty:                                  "QualifierPointedParty",
	MRAElementQualifierPointingParty:                                 "QualifierPointingParty",
	MRAElementTrustServiceEquivalenceHistory:                         "TrustServiceEquivalenceHistory",
	MRAElementTrustServiceEquivalenceHistoryInstance:                 "TrustServiceEquivalenceHistoryInstance",
	MRAElementTrustServiceEquivalenceInformation:                     "TrustServiceEquivalenceInformation",
	MRAElementTrustServiceEquivalenceStatus:                          "TrustServiceEquivalenceStatus",
	MRAElementTrustServiceEquivalenceStatusStartingTime:              "TrustServiceEquivalenceStatusStartingTime",
	MRAElementTrustServiceLegalIdentifier:                            "TrustServiceLegalIdentifier",
	MRAElementTrustServiceTSLQualificationExtensionEquivalenceList:   "TrustServiceTSLQualificationExtensionEquivalenceList",
	MRAElementTrustServiceTSLQualificationExtensionName:              "TrustServiceTSLQualificationExtensionName",
	MRAElementTrustServiceTSLQualificationExtensionNamePointedParty:  "TrustServiceTSLQualificationExtensionNamePointedParty",
	MRAElementTrustServiceTSLQualificationExtensionNamePointingParty: "TrustServiceTSLQualificationExtensionNamePointingParty",
	MRAElementTrustServiceTSLStatusEquivalenceList:                   "TrustServiceTSLStatusEquivalenceList",
	MRAElementTrustServiceTSLStatusInvalidEquivalence:                "TrustServiceTSLStatusInvalidEquivalence",
	MRAElementTrustServiceTSLStatusListPointedParty:                  "TrustServiceTSLStatusListPointedParty",
	MRAElementTrustServiceTSLStatusListPointingParty:                 "TrustServiceTSLStatusListPointingParty",
	MRAElementTrustServiceTSLStatusValidEquivalence:                  "TrustServiceTSLStatusValidEquivalence",
	MRAElementTrustServiceTSLType:                                    "TrustServiceTSLType",
	MRAElementTrustServiceTSLTypeEquivalenceList:                     "TrustServiceTSLTypeEquivalenceList",
	MRAElementTrustServiceTSLTypeListPointedParty:                    "TrustServiceTSLTypeListPointedParty",
	MRAElementTrustServiceTSLTypeListPointingParty:                   "TrustServiceTSLTypeListPointingParty",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e MRAElement) TagName() string {
	return mraElementTagNames[e]
}

// Namespace implements common.DSSElement. Ports getNamespace(): every MRAElement constant is
// built with MRANamespace.NS.
func (e MRAElement) Namespace() *common.DSSNamespace {
	return MRANamespaceNS
}

// URI implements common.DSSElement. Ports getURI().
func (e MRAElement) URI() string {
	return MRANamespaceNS.Uri()
}

// IsSameTagName implements common.DSSElement. Ports isSameTagName(String).
func (e MRAElement) IsSameTagName(value string) bool {
	return e.TagName() == value
}

var _ common.DSSElement = MRAElement("")
