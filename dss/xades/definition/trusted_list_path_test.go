// KAT test for trusted_list_path.go: every TrustedListPath_* query string below is
// transcribed verbatim from the upstream Java source (dss-xades 6.5.RC1)
// eu.europa.esig.dss.xades.definition.tsl.TrustedListPath.
package definition

import "testing"

func TestTrustedListPath_KAT(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ADDITIONAL_SERVICE_INFORMATION_PATH", TrustedListPathAdditionalServiceInformationPath.QueryString(), "./tl:AdditionalServiceInformation"},
		{"NEXT_UPDATE_PATH", TrustedListPathNextUpdatePath.QueryString(), "./tl:SchemeInformation/tl:NextUpdate"},
		{"OTHER_TSL_POINTER_PATH", TrustedListPathOtherTSLPointerPath.QueryString(), "./tl:SchemeInformation/tl:PointersToOtherTSL/tl:OtherTSLPointer"},
		{"SERVICE_DIGITAL_IDENTITY_PATH", TrustedListPathServiceDigitalIdentityPath.QueryString(), "./tl:ServiceDigitalIdentity"},
		{"TSL_VERSION_IDENTIFIER_PATH", TrustedListPathTSLVersionIdentifierPath.QueryString(), "./tl:SchemeInformation/tl:TSLVersionIdentifier"},
		{"X509_CERTIFICATE_PATH", TrustedListPathX509CertificatePath.QueryString(), "./tl:ServiceDigitalIdentities/tl:ServiceDigitalIdentity/tl:DigitalId/tl:X509Certificate"},
		{"TSP_SERVICE_INFORMATION_PATH", TrustedListPathTSPServiceInformationPath.QueryString(), "./tl:TrustServiceProviderList/tl:TrustServiceProvider/tl:TSPServices/tl:TSPService/tl:ServiceInformation"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.QueryString() = %q, want %q", c.name, c.got, c.want)
		}
	}
}
