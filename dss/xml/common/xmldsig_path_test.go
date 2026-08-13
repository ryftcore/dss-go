// KAT test for xmldsig_path.go: every XMLDSigPath_* query string below is dumped verbatim
// from a Java oracle run over dss-xml-common 6.5.RC1 + dss-alert 6.5.RC1
// (java.lang.reflect over XMLDSigPath.class.getDeclaredFields(), each XPathQuery field's
// getQueryString() printed), per PORTING.md's "exhaustive table test" rule for
// registry-like tables and this phase's explicit KAT requirement for AbstractPath's output.
package common

import "testing"

func TestXMLDSigPath_KAT(t *testing.T) {
	cases := []struct {
		name string
		got  XPathQuery
		want string
	}{
		{"SIGNATURE_PATH", XMLDSigPath_SIGNATURE_PATH, "./ds:Signature"},
		{"ALL_SIGNATURES_PATH", XMLDSigPath_ALL_SIGNATURES_PATH, "//ds:Signature"},
		{"OBJECT_PATH", XMLDSigPath_OBJECT_PATH, "./ds:Object"},
		{"MANIFEST_PATH", XMLDSigPath_MANIFEST_PATH, "./ds:Object/ds:Manifest"},
		{"SIGNED_INFO_PATH", XMLDSigPath_SIGNED_INFO_PATH, "./ds:SignedInfo"},
		{"SIGNED_INFO_CANONICALIZATION_METHOD", XMLDSigPath_SIGNED_INFO_CANONICALIZATION_METHOD, "./ds:SignedInfo/ds:CanonicalizationMethod"},
		{"SIGNED_INFO_REFERENCE_PATH", XMLDSigPath_SIGNED_INFO_REFERENCE_PATH, "./ds:SignedInfo/ds:Reference"},
		{"SIGNATURE_METHOD_PATH", XMLDSigPath_SIGNATURE_METHOD_PATH, "./ds:SignedInfo/ds:SignatureMethod"},
		{"REFERENCE_PATH", XMLDSigPath_REFERENCE_PATH, "./ds:Reference"},
		{"SIGNATURE_VALUE_PATH", XMLDSigPath_SIGNATURE_VALUE_PATH, "./ds:SignatureValue"},
		{"SIGNATURE_VALUE_ID_PATH", XMLDSigPath_SIGNATURE_VALUE_ID_PATH, "./ds:SignatureValue/@Id"},
		{"ALL_SIGNATURE_VALUES_PATH", XMLDSigPath_ALL_SIGNATURE_VALUES_PATH, "//ds:SignatureValue"},
		{"KEY_INFO_PATH", XMLDSigPath_KEY_INFO_PATH, "./ds:KeyInfo"},
		{"KEY_INFO_X509_DATA", XMLDSigPath_KEY_INFO_X509_DATA, "./ds:KeyInfo/ds:X509Data"},
		{"KEY_INFO_X509_CERTIFICATE_PATH", XMLDSigPath_KEY_INFO_X509_CERTIFICATE_PATH, "./ds:KeyInfo/ds:X509Data/ds:X509Certificate"},
		{"SIGNATURE_PROPERTIES_PATH", XMLDSigPath_SIGNATURE_PROPERTIES_PATH, "./ds:Object/ds:SignatureProperties"},
		{"SIGNATURE_PROPERTY_PATH", XMLDSigPath_SIGNATURE_PROPERTY_PATH, "./ds:Object/ds:SignatureProperties/ds:SignatureProperty"},
		{"DIGEST_METHOD_ALGORITHM_PATH", XMLDSigPath_DIGEST_METHOD_ALGORITHM_PATH, "./ds:DigestMethod/@Algorithm"},
		{"DIGEST_VALUE_PATH", XMLDSigPath_DIGEST_VALUE_PATH, "./ds:DigestValue"},
		{"CANONICALIZATION_ALGORITHM_PATH", XMLDSigPath_CANONICALIZATION_ALGORITHM_PATH, "./ds:CanonicalizationMethod/@Algorithm"},
		{"TRANSFORM_PATH", XMLDSigPath_TRANSFORM_PATH, "./ds:Transform"},
		{"TRANSFORMS_PATH", XMLDSigPath_TRANSFORMS_PATH, "./ds:Transforms"},
		{"TRANSFORMS_TRANSFORM_PATH", XMLDSigPath_TRANSFORMS_TRANSFORM_PATH, "./ds:Transforms/ds:Transform"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.got.QueryString(); got != c.want {
				t.Errorf("%s.QueryString() = %q, want %q", c.name, got, c.want)
			}
		})
	}
}

func TestXMLDSigPath_TypeConstants(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"OBJECT_TYPE", XMLDSigPath_OBJECT_TYPE, "http://www.w3.org/2000/09/xmldsig#Object"},
		{"MANIFEST_TYPE", XMLDSigPath_MANIFEST_TYPE, "http://www.w3.org/2000/09/xmldsig#Manifest"},
		{"COUNTER_SIGNATURE_TYPE", XMLDSigPath_COUNTER_SIGNATURE_TYPE, "http://uri.etsi.org/01903#CountersignedSignature"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
}
