// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/tsl/TrustedListPath.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// ETSI TS 119 612 Trusted List XML path definitions.
//
// KAT-verified against the Java source (dss-xades 6.5.RC1): each XPathQuery below is
// built from the same element chain as the corresponding upstream field, whose resolved
// query string is reproduced in the comment above it - see trusted_list_path_test.go.
var (
	// ADDITIONAL_SERVICE_INFORMATION_PATH = "./tl:AdditionalServiceInformation"
	TrustedListPathAdditionalServiceInformationPath = common.FromCurrentPosition(TrustedListElementAdditionalServiceInformation)

	// NEXT_UPDATE_PATH = "./tl:SchemeInformation/tl:NextUpdate"
	TrustedListPathNextUpdatePath = common.FromCurrentPosition(TrustedListElementSchemeInformation, TrustedListElementNextUpdate)

	// OTHER_TSL_POINTER_PATH = "./tl:SchemeInformation/tl:PointersToOtherTSL/tl:OtherTSLPointer"
	TrustedListPathOtherTSLPointerPath = common.FromCurrentPosition(TrustedListElementSchemeInformation,
		TrustedListElementPointersToOtherTSL, TrustedListElementOtherTSLPointer)

	// SERVICE_DIGITAL_IDENTITY_PATH = "./tl:ServiceDigitalIdentity"
	TrustedListPathServiceDigitalIdentityPath = common.FromCurrentPosition(TrustedListElementServiceDigitalIdentity)

	// TSL_VERSION_IDENTIFIER_PATH = "./tl:SchemeInformation/tl:TSLVersionIdentifier"
	TrustedListPathTSLVersionIdentifierPath = common.FromCurrentPosition(TrustedListElementSchemeInformation, TrustedListElementTSLVersionIdentifier)

	// X509_CERTIFICATE_PATH = "./tl:ServiceDigitalIdentities/tl:ServiceDigitalIdentity/tl:DigitalId/tl:X509Certificate"
	TrustedListPathX509CertificatePath = common.FromCurrentPosition(TrustedListElementServiceDigitalIdentities,
		TrustedListElementServiceDigitalIdentity, TrustedListElementDigitalID, TrustedListElementX509Certificate)

	// TSP_SERVICE_INFORMATION_PATH = "./tl:TrustServiceProviderList/tl:TrustServiceProvider/tl:TSPServices/tl:TSPService/tl:ServiceInformation"
	TrustedListPathTSPServiceInformationPath = common.FromCurrentPosition(TrustedListElementTrustServiceProviderList,
		TrustedListElementTrustServiceProvider, TrustedListElementTSPServices, TrustedListElementTSPService, TrustedListElementServiceInformation)
)
