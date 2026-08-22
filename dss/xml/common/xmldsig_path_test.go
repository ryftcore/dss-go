// KAT test for xmldsig_path.go: every XMLDSigPath_* query string below is dumped verbatim
// from a Java oracle run over dss-xml-common 6.5.RC1 + dss-alert 6.5.RC1
// (java.lang.reflect over XMLDSigPath.class.getDeclaredFields(), each XPathQuery field's
// getQueryString() printed).
package common

import "testing"

func TestXMLDSigPath_KAT(t *testing.T) {
	cases := []struct {
		name string
		got  XPathQuery
		want string
	}{
		{"SIGNATURE_PATH", XMLDSigPathSignaturePath, "./ds:Signature"},
		{"ALL_SIGNATURES_PATH", XMLDSigPathAllSignaturesPath, "//ds:Signature"},
		{"OBJECT_PATH", XMLDSigPathObjectPath, "./ds:Object"},
		{"MANIFEST_PATH", XMLDSigPathManifestPath, "./ds:Object/ds:Manifest"},
		{"SIGNED_INFO_PATH", XMLDSigPathSignedInfoPath, "./ds:SignedInfo"},
		{"SIGNED_INFO_CANONICALIZATION_METHOD", XMLDSigPathSignedInfoCanonicalizationMethod, "./ds:SignedInfo/ds:CanonicalizationMethod"},
		{"SIGNED_INFO_REFERENCE_PATH", XMLDSigPathSignedInfoReferencePath, "./ds:SignedInfo/ds:Reference"},
		{"SIGNATURE_METHOD_PATH", XMLDSigPathSignatureMethodPath, "./ds:SignedInfo/ds:SignatureMethod"},
		{"REFERENCE_PATH", XMLDSigPathReferencePath, "./ds:Reference"},
		{"SIGNATURE_VALUE_PATH", XMLDSigPathSignatureValuePath, "./ds:SignatureValue"},
		{"SIGNATURE_VALUE_ID_PATH", XMLDSigPathSignatureValueIDPath, "./ds:SignatureValue/@Id"},
		{"ALL_SIGNATURE_VALUES_PATH", XMLDSigPathAllSignatureValuesPath, "//ds:SignatureValue"},
		{"KEY_INFO_PATH", XMLDSigPathKeyInfoPath, "./ds:KeyInfo"},
		{"KEY_INFO_X509_DATA", XMLDSigPathKeyInfoX509Data, "./ds:KeyInfo/ds:X509Data"},
		{"KEY_INFO_X509_CERTIFICATE_PATH", XMLDSigPathKeyInfoX509CertificatePath, "./ds:KeyInfo/ds:X509Data/ds:X509Certificate"},
		{"SIGNATURE_PROPERTIES_PATH", XMLDSigPathSignaturePropertiesPath, "./ds:Object/ds:SignatureProperties"},
		{"SIGNATURE_PROPERTY_PATH", XMLDSigPathSignaturePropertyPath, "./ds:Object/ds:SignatureProperties/ds:SignatureProperty"},
		{"DIGEST_METHOD_ALGORITHM_PATH", XMLDSigPathDigestMethodAlgorithmPath, "./ds:DigestMethod/@Algorithm"},
		{"DIGEST_VALUE_PATH", XMLDSigPathDigestValuePath, "./ds:DigestValue"},
		{"CANONICALIZATION_ALGORITHM_PATH", XMLDSigPathCanonicalizationAlgorithmPath, "./ds:CanonicalizationMethod/@Algorithm"},
		{"TRANSFORM_PATH", XMLDSigPathTransformPath, "./ds:Transform"},
		{"TRANSFORMS_PATH", XMLDSigPathTransformsPath, "./ds:Transforms"},
		{"TRANSFORMS_TRANSFORM_PATH", XMLDSigPathTransformsTransformPath, "./ds:Transforms/ds:Transform"},
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
		{"OBJECT_TYPE", XMLDSigPathObjectType, "http://www.w3.org/2000/09/xmldsig#Object"},
		{"MANIFEST_TYPE", XMLDSigPathManifestType, "http://www.w3.org/2000/09/xmldsig#Manifest"},
		{"COUNTER_SIGNATURE_TYPE", XMLDSigPathCounterSignatureType, "http://uri.etsi.org/01903#CountersignedSignature"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
}
