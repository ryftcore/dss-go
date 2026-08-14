// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/tsl/TrustedListPath.java (DSS 6.5.RC1).
package definition

import "github.com/utain/esig/dss/xml/common"

// ETSI TS 119 612 Trusted List XML path definitions.
//
// KAT-verified against the Java source (dss-xades 6.5.RC1): each XPathQuery below is
// built from the same element chain as the corresponding upstream field, whose resolved
// query string is reproduced in the comment above it - see trusted_list_path_test.go.
var (
	// ADDITIONAL_SERVICE_INFORMATION_PATH = "./tl:AdditionalServiceInformation"
	TrustedListPath_ADDITIONAL_SERVICE_INFORMATION_PATH = common.FromCurrentPosition(TrustedListElement_ADDITIONAL_SERVICE_INFORMATION)

	// NEXT_UPDATE_PATH = "./tl:SchemeInformation/tl:NextUpdate"
	TrustedListPath_NEXT_UPDATE_PATH = common.FromCurrentPosition(TrustedListElement_SCHEME_INFORMATION, TrustedListElement_NEXT_UPDATE)

	// OTHER_TSL_POINTER_PATH = "./tl:SchemeInformation/tl:PointersToOtherTSL/tl:OtherTSLPointer"
	TrustedListPath_OTHER_TSL_POINTER_PATH = common.FromCurrentPosition(TrustedListElement_SCHEME_INFORMATION,
		TrustedListElement_POINTERS_TO_OTHER_TSL, TrustedListElement_OTHER_TSL_POINTER)

	// SERVICE_DIGITAL_IDENTITY_PATH = "./tl:ServiceDigitalIdentity"
	TrustedListPath_SERVICE_DIGITAL_IDENTITY_PATH = common.FromCurrentPosition(TrustedListElement_SERVICE_DIGITAL_IDENTITY)

	// TSL_VERSION_IDENTIFIER_PATH = "./tl:SchemeInformation/tl:TSLVersionIdentifier"
	TrustedListPath_TSL_VERSION_IDENTIFIER_PATH = common.FromCurrentPosition(TrustedListElement_SCHEME_INFORMATION, TrustedListElement_TSL_VERSION_IDENTIFIER)

	// X509_CERTIFICATE_PATH = "./tl:ServiceDigitalIdentities/tl:ServiceDigitalIdentity/tl:DigitalId/tl:X509Certificate"
	TrustedListPath_X509_CERTIFICATE_PATH = common.FromCurrentPosition(TrustedListElement_SERVICE_DIGITAL_IDENTITIES,
		TrustedListElement_SERVICE_DIGITAL_IDENTITY, TrustedListElement_DIGITAL_ID, TrustedListElement_X509_CERTIFICATE)

	// TSP_SERVICE_INFORMATION_PATH = "./tl:TrustServiceProviderList/tl:TrustServiceProvider/tl:TSPServices/tl:TSPService/tl:ServiceInformation"
	TrustedListPath_TSP_SERVICE_INFORMATION_PATH = common.FromCurrentPosition(TrustedListElement_TRUST_SERVICE_PROVIDER_LIST,
		TrustedListElement_TRUST_SERVICE_PROVIDER, TrustedListElement_TSP_SERVICES, TrustedListElement_TSP_SERVICE, TrustedListElement_SERVICE_INFORMATION)
)
