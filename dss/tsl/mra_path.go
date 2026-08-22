// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/definition/mra/MRAPath.java (DSS 6.5.RC1).
package tsl

import (
	xadesdefinition "github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/common"
)

// XPath expressions used within the implementation in relation to the MRA (Mutual Recognition
// Agreement) scheme.
//
// Upstream is a class extending AbstractPath purely to reach its static builders; since those
// builders are free package functions in the Go port of dss-xml-common (see
// xml/common/abstract_path.go), the fields become package-level variables here, following the
// precedent of xades/definition's TrustedListPath. Java's public no-arg constructor has no
// counterpart (the class carries no state).
//
// KAT-verified against the Java source: each XPathQuery below is built from the same element
// chain as the corresponding upstream field, whose resolved query string is reproduced in the
// comment above it - see mra_path_test.go.
var (
	// CERTIFICATE_CONTENT_DECLARATION_POINTED_PARTY_PATH = "./mra:CertificateContentDeclarationPointedParty"
	MRAPathCertificateContentDeclarationPointedPartyPath = common.FromCurrentPosition(
		MRAElementCertificateContentDeclarationPointedParty)

	// CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY_PATH = "./mra:CertificateContentDeclarationPointingParty"
	MRAPathCertificateContentDeclarationPointingPartyPath = common.FromCurrentPosition(
		MRAElementCertificateContentDeclarationPointingParty)

	// CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:CertificateContentReferencesEquivalenceList/mra:CertificateContentReferenceEquivalence"
	MRAPathCertificateContentReferencesEquivalencePath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceInformation,
		MRAElementCertificateContentReferencesEquivalenceList,
		MRAElementCertificateContentReferencesEquivalence)

	// CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT_PATH = "./mra:CertificateContentReferenceEquivalenceContext"
	MRAPathCertificateContentReferenceEquivalenceContextPath = common.FromCurrentPosition(
		MRAElementCertificateContentReferenceEquivalenceContext)

	// MUTUAL_RECOGNITION_AGREEMENT_INFORMATION_PATH = "//mra:MutualRecognitionAgreementInformation"
	MRAPathMutualRecognitionAgreementInformationPath = common.All(
		MRAElementMutualRecognitionAgreementInformation)

	// QUALIFIER_EQUIVALENCE_LIST_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLQualificationExtensionEquivalenceList/mra:QualifierEquivalenceList"
	MRAPathQualifierEquivalenceListPath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceInformation,
		MRAElementTrustServiceTSLQualificationExtensionEquivalenceList,
		MRAElementQualifierEquivalenceList)

	// SERVICE_TYPE_IDENTIFIER_PATH = "./mra:TrustServiceTSLType/tl:ServiceTypeIdentifier"
	MRAPathServiceTypeIdentifierPath = common.FromCurrentPosition(
		MRAElementTrustServiceTSLType,
		xadesdefinition.TrustedListElementServiceTypeIdentifier)

	// TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceEquivalenceHistory/mra:TrustServiceEquivalenceHistoryInstance"
	MRAPathTrustServiceEquivalenceHistoryInstancePath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceInformation,
		MRAElementTrustServiceEquivalenceHistory,
		MRAElementTrustServiceEquivalenceHistoryInstance)

	// TRUST_SERVICE_EQUIVALENCE_INFORMATION_STATUS_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceEquivalenceStatus"
	MRAPathTrustServiceEquivalenceInformationStatusPath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceInformation,
		MRAElementTrustServiceEquivalenceStatus)

	// TRUST_SERVICE_EQUIVALENCE_INFORMATION_STATUS_STARTING_TIME_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceEquivalenceStatusStartingTime"
	MRAPathTrustServiceEquivalenceInformationStatusStartingTimePath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceInformation,
		MRAElementTrustServiceEquivalenceStatusStartingTime)

	// TRUST_SERVICE_EQUIVALENCE_STATUS_PATH = "./mra:TrustServiceEquivalenceStatus"
	MRAPathTrustServiceEquivalenceStatusPath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceStatus)

	// TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME_PATH = "./mra:TrustServiceEquivalenceStatusStartingTime"
	MRAPathTrustServiceEquivalenceStatusStartingTimePath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceStatusStartingTime)

	// TRUST_SERVICE_LEGAL_IDENTIFIER_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceLegalIdentifier"
	MRAPathTrustServiceLegalIdentifierPath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceInformation,
		MRAElementTrustServiceLegalIdentifier)

	// TRUST_SERVICE_TSL_TYPE_PATH = "./mra:TrustServiceTSLType"
	MRAPathTrustServiceTSLTypePath = common.FromCurrentPosition(
		MRAElementTrustServiceTSLType)

	// TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLTypeEquivalenceList/mra:TrustServiceTSLTypeListPointedParty"
	MRAPathTrustServiceTSLTypeListPointedPartyPath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceInformation,
		MRAElementTrustServiceTSLTypeEquivalenceList,
		MRAElementTrustServiceTSLTypeListPointedParty)

	// TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLTypeEquivalenceList/mra:TrustServiceTSLTypeListPointingParty"
	MRAPathTrustServiceTSLTypeListPointingPartyPath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceInformation,
		MRAElementTrustServiceTSLTypeEquivalenceList,
		MRAElementTrustServiceTSLTypeListPointingParty)

	// TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE_LIST_POINTED_PARTY_SERVICE_STATUS_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusInvalidEquivalence/mra:TrustServiceTSLStatusListPointedParty/tl:ServiceStatus"
	MRAPathTrustServiceTSLStatusInvalidEquivalenceListPointedPartyServiceStatusPath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceInformation,
		MRAElementTrustServiceTSLStatusEquivalenceList,
		MRAElementTrustServiceTSLStatusInvalidEquivalence,
		MRAElementTrustServiceTSLStatusListPointedParty,
		xadesdefinition.TrustedListElementServiceStatus)

	// TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE_LIST_POINTING_PARTY_SERVICE_STATUS_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusInvalidEquivalence/mra:TrustServiceTSLStatusListPointingParty/tl:ServiceStatus"
	MRAPathTrustServiceTSLStatusInvalidEquivalenceListPointingPartyServiceStatusPath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceInformation,
		MRAElementTrustServiceTSLStatusEquivalenceList,
		MRAElementTrustServiceTSLStatusInvalidEquivalence,
		MRAElementTrustServiceTSLStatusListPointingParty,
		xadesdefinition.TrustedListElementServiceStatus)

	// TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE_LIST_POINTED_PARTY_SERVICE_STATUS_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusValidEquivalence/mra:TrustServiceTSLStatusListPointedParty/tl:ServiceStatus"
	MRAPathTrustServiceTSLStatusValidEquivalenceListPointedPartyServiceStatusPath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceInformation,
		MRAElementTrustServiceTSLStatusEquivalenceList,
		MRAElementTrustServiceTSLStatusValidEquivalence,
		MRAElementTrustServiceTSLStatusListPointedParty,
		xadesdefinition.TrustedListElementServiceStatus)

	// TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE_LIST_POINTING_PARTY_SERVICE_STATUS_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusValidEquivalence/mra:TrustServiceTSLStatusListPointingParty/tl:ServiceStatus"
	MRAPathTrustServiceTSLStatusValidEquivalenceListPointingPartyServiceStatusPath = common.FromCurrentPosition(
		MRAElementTrustServiceEquivalenceInformation,
		MRAElementTrustServiceTSLStatusEquivalenceList,
		MRAElementTrustServiceTSLStatusValidEquivalence,
		MRAElementTrustServiceTSLStatusListPointingParty,
		xadesdefinition.TrustedListElementServiceStatus)
)
