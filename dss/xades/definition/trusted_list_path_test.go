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
		{"ADDITIONAL_SERVICE_INFORMATION_PATH", TrustedListPath_ADDITIONAL_SERVICE_INFORMATION_PATH.QueryString(), "./tl:AdditionalServiceInformation"},
		{"NEXT_UPDATE_PATH", TrustedListPath_NEXT_UPDATE_PATH.QueryString(), "./tl:SchemeInformation/tl:NextUpdate"},
		{"OTHER_TSL_POINTER_PATH", TrustedListPath_OTHER_TSL_POINTER_PATH.QueryString(), "./tl:SchemeInformation/tl:PointersToOtherTSL/tl:OtherTSLPointer"},
		{"SERVICE_DIGITAL_IDENTITY_PATH", TrustedListPath_SERVICE_DIGITAL_IDENTITY_PATH.QueryString(), "./tl:ServiceDigitalIdentity"},
		{"TSL_VERSION_IDENTIFIER_PATH", TrustedListPath_TSL_VERSION_IDENTIFIER_PATH.QueryString(), "./tl:SchemeInformation/tl:TSLVersionIdentifier"},
		{"X509_CERTIFICATE_PATH", TrustedListPath_X509_CERTIFICATE_PATH.QueryString(), "./tl:ServiceDigitalIdentities/tl:ServiceDigitalIdentity/tl:DigitalId/tl:X509Certificate"},
		{"TSP_SERVICE_INFORMATION_PATH", TrustedListPath_TSP_SERVICE_INFORMATION_PATH.QueryString(), "./tl:TrustServiceProviderList/tl:TrustServiceProvider/tl:TSPServices/tl:TSPService/tl:ServiceInformation"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.QueryString() = %q, want %q", c.name, c.got, c.want)
		}
	}
}
