// KAT test for xades141_element.go: every XAdES141Element_* tag name below is transcribed
// verbatim from the upstream Java source (dss-xades 6.5.RC1).
package definition

import "testing"

func TestXAdES141Element_KAT(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ANY_VALIDATION_DATA", XAdES141Element_ANY_VALIDATION_DATA.TagName(), "AnyValidationData"},
		{"ARCHIVE_TIMESTAMP", XAdES141Element_ARCHIVE_TIMESTAMP.TagName(), "ArchiveTimeStamp"},
		{"ATTRIBUTE_CERTIFICATE_REFS_V2", XAdES141Element_ATTRIBUTE_CERTIFICATE_REFS_V2.TagName(), "AttributeCertificateRefsV2"},
		{"CERT_REFS", XAdES141Element_CERT_REFS.TagName(), "CertRefs"},
		{"COMPLETE_CERTIFICATE_REFS_V2", XAdES141Element_COMPLETE_CERTIFICATE_REFS_V2.TagName(), "CompleteCertificateRefsV2"},
		{"NEW_SDO_DIGEST_VALUE", XAdES141Element_NEW_SDO_DIGEST_VALUE.TagName(), "NewSDODigestValue"},
		{"ORIGINAL_REF_DIGEST", XAdES141Element_ORIGINAL_REF_DIGEST.TagName(), "OriginalRefDigest"},
		{"RECOMPUTED_DIGEST_VALUE", XAdES141Element_RECOMPUTED_DIGEST_VALUE.TagName(), "RecomputedDigestValue"},
		{"REFS_ONLY_TIMESTAMP_V2", XAdES141Element_REFS_ONLY_TIMESTAMP_V2.TagName(), "RefsOnlyTimeStampV2"},
		{"RENEWED_DIGESTS", XAdES141Element_RENEWED_DIGESTS.TagName(), "RenewedDigests"},
		{"RENEWED_DIGESTS_V2", XAdES141Element_RENEWED_DIGESTS_V2.TagName(), "RenewedDigestsV2"},
		{"SIG_AND_REFS_TIMESTAMP_V2", XAdES141Element_SIG_AND_REFS_TIMESTAMP_V2.TagName(), "SigAndRefsTimeStampV2"},
		{"SIGNATURE_POLICY_DOCUMENT", XAdES141Element_SIGNATURE_POLICY_DOCUMENT.TagName(), "SignaturePolicyDocument"},
		{"SIGNATURE_POLICY_STORE", XAdES141Element_SIGNATURE_POLICY_STORE.TagName(), "SignaturePolicyStore"},
		{"SIG_POL_DOC_LOCAL_URI", XAdES141Element_SIG_POL_DOC_LOCAL_URI.TagName(), "SigPolDocLocalURI"},
		{"SP_DOC_SPECIFICATION", XAdES141Element_SP_DOC_SPECIFICATION.TagName(), "SPDocSpecification"},
		{"TIMESTAMP_VALIDATION_DATA", XAdES141Element_TIMESTAMP_VALIDATION_DATA.TagName(), "TimeStampValidationData"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.TagName() = %q, want %q", c.name, c.got, c.want)
		}
	}
	if got := XAdES141Element_SP_DOC_SPECIFICATION.URI(); got != XAdESNamespace_XADES_141.Uri() {
		t.Errorf("URI() = %q, want %q", got, XAdESNamespace_XADES_141.Uri())
	}
}
