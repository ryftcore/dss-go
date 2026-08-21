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
	MRAPath_CERTIFICATE_CONTENT_DECLARATION_POINTED_PARTY_PATH = common.FromCurrentPosition(
		MRAElement_CERTIFICATE_CONTENT_DECLARATION_POINTED_PARTY)

	// CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY_PATH = "./mra:CertificateContentDeclarationPointingParty"
	MRAPath_CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY_PATH = common.FromCurrentPosition(
		MRAElement_CERTIFICATE_CONTENT_DECLARATION_POINTING_PARTY)

	// CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:CertificateContentReferencesEquivalenceList/mra:CertificateContentReferenceEquivalence"
	MRAPath_CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION,
		MRAElement_CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE_LIST,
		MRAElement_CERTIFICATE_CONTENT_REFERENCES_EQUIVALENCE)

	// CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT_PATH = "./mra:CertificateContentReferenceEquivalenceContext"
	MRAPath_CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT_PATH = common.FromCurrentPosition(
		MRAElement_CERTIFICATE_CONTENT_REFERENCE_EQUIVALENCE_CONTEXT)

	// MUTUAL_RECOGNITION_AGREEMENT_INFORMATION_PATH = "//mra:MutualRecognitionAgreementInformation"
	MRAPath_MUTUAL_RECOGNITION_AGREEMENT_INFORMATION_PATH = common.All(
		MRAElement_MUTUAL_RECOGNITION_AGREEMENT_INFORMATION)

	// QUALIFIER_EQUIVALENCE_LIST_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLQualificationExtensionEquivalenceList/mra:QualifierEquivalenceList"
	MRAPath_QUALIFIER_EQUIVALENCE_LIST_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION,
		MRAElement_TRUST_SERVICE_TSL_QUALIFICATION_EXTENSION_EQUIVALENCE_LIST,
		MRAElement_QUALIFIER_EQUIVALENCE_LIST)

	// SERVICE_TYPE_IDENTIFIER_PATH = "./mra:TrustServiceTSLType/tl:ServiceTypeIdentifier"
	MRAPath_SERVICE_TYPE_IDENTIFIER_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_TSL_TYPE,
		xadesdefinition.TrustedListElement_SERVICE_TYPE_IDENTIFIER)

	// TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceEquivalenceHistory/mra:TrustServiceEquivalenceHistoryInstance"
	MRAPath_TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION,
		MRAElement_TRUST_SERVICE_EQUIVALENCE_HISTORY,
		MRAElement_TRUST_SERVICE_EQUIVALENCE_HISTORY_INSTANCE)

	// TRUST_SERVICE_EQUIVALENCE_INFORMATION_STATUS_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceEquivalenceStatus"
	MRAPath_TRUST_SERVICE_EQUIVALENCE_INFORMATION_STATUS_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION,
		MRAElement_TRUST_SERVICE_EQUIVALENCE_STATUS)

	// TRUST_SERVICE_EQUIVALENCE_INFORMATION_STATUS_STARTING_TIME_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceEquivalenceStatusStartingTime"
	MRAPath_TRUST_SERVICE_EQUIVALENCE_INFORMATION_STATUS_STARTING_TIME_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION,
		MRAElement_TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME)

	// TRUST_SERVICE_EQUIVALENCE_STATUS_PATH = "./mra:TrustServiceEquivalenceStatus"
	MRAPath_TRUST_SERVICE_EQUIVALENCE_STATUS_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_STATUS)

	// TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME_PATH = "./mra:TrustServiceEquivalenceStatusStartingTime"
	MRAPath_TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_STATUS_STARTING_TIME)

	// TRUST_SERVICE_LEGAL_IDENTIFIER_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceLegalIdentifier"
	MRAPath_TRUST_SERVICE_LEGAL_IDENTIFIER_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION,
		MRAElement_TRUST_SERVICE_LEGAL_IDENTIFIER)

	// TRUST_SERVICE_TSL_TYPE_PATH = "./mra:TrustServiceTSLType"
	MRAPath_TRUST_SERVICE_TSL_TYPE_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_TSL_TYPE)

	// TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLTypeEquivalenceList/mra:TrustServiceTSLTypeListPointedParty"
	MRAPath_TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION,
		MRAElement_TRUST_SERVICE_TSL_TYPE_EQUIVALENCE_LIST,
		MRAElement_TRUST_SERVICE_TSL_TYPE_LIST_POINTED_PARTY)

	// TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLTypeEquivalenceList/mra:TrustServiceTSLTypeListPointingParty"
	MRAPath_TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION,
		MRAElement_TRUST_SERVICE_TSL_TYPE_EQUIVALENCE_LIST,
		MRAElement_TRUST_SERVICE_TSL_TYPE_LIST_POINTING_PARTY)

	// TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE_LIST_POINTED_PARTY_SERVICE_STATUS_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusInvalidEquivalence/mra:TrustServiceTSLStatusListPointedParty/tl:ServiceStatus"
	MRAPath_TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE_LIST_POINTED_PARTY_SERVICE_STATUS_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION,
		MRAElement_TRUST_SERVICE_TSL_STATUS_EQUIVALENCE_LIST,
		MRAElement_TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE,
		MRAElement_TRUST_SERVICE_TSL_STATUS_LIST_POINTED_PARTY,
		xadesdefinition.TrustedListElement_SERVICE_STATUS)

	// TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE_LIST_POINTING_PARTY_SERVICE_STATUS_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusInvalidEquivalence/mra:TrustServiceTSLStatusListPointingParty/tl:ServiceStatus"
	MRAPath_TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE_LIST_POINTING_PARTY_SERVICE_STATUS_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION,
		MRAElement_TRUST_SERVICE_TSL_STATUS_EQUIVALENCE_LIST,
		MRAElement_TRUST_SERVICE_TSL_STATUS_INVALID_EQUIVALENCE,
		MRAElement_TRUST_SERVICE_TSL_STATUS_LIST_POINTING_PARTY,
		xadesdefinition.TrustedListElement_SERVICE_STATUS)

	// TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE_LIST_POINTED_PARTY_SERVICE_STATUS_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusValidEquivalence/mra:TrustServiceTSLStatusListPointedParty/tl:ServiceStatus"
	MRAPath_TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE_LIST_POINTED_PARTY_SERVICE_STATUS_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION,
		MRAElement_TRUST_SERVICE_TSL_STATUS_EQUIVALENCE_LIST,
		MRAElement_TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE,
		MRAElement_TRUST_SERVICE_TSL_STATUS_LIST_POINTED_PARTY,
		xadesdefinition.TrustedListElement_SERVICE_STATUS)

	// TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE_LIST_POINTING_PARTY_SERVICE_STATUS_PATH =
	// "./mra:TrustServiceEquivalenceInformation/mra:TrustServiceTSLStatusEquivalenceList/mra:TrustServiceTSLStatusValidEquivalence/mra:TrustServiceTSLStatusListPointingParty/tl:ServiceStatus"
	MRAPath_TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE_LIST_POINTING_PARTY_SERVICE_STATUS_PATH = common.FromCurrentPosition(
		MRAElement_TRUST_SERVICE_EQUIVALENCE_INFORMATION,
		MRAElement_TRUST_SERVICE_TSL_STATUS_EQUIVALENCE_LIST,
		MRAElement_TRUST_SERVICE_TSL_STATUS_VALID_EQUIVALENCE,
		MRAElement_TRUST_SERVICE_TSL_STATUS_LIST_POINTING_PARTY,
		xadesdefinition.TrustedListElement_SERVICE_STATUS)
)
