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
		{"ANY_VALIDATION_DATA", XAdES141ElementAnyValidationData.TagName(), "AnyValidationData"},
		{"ARCHIVE_TIMESTAMP", XAdES141ElementArchiveTimestamp.TagName(), "ArchiveTimeStamp"},
		{"ATTRIBUTE_CERTIFICATE_REFS_V2", XAdES141ElementAttributeCertificateRefsV2.TagName(), "AttributeCertificateRefsV2"},
		{"CERT_REFS", XAdES141ElementCertRefs.TagName(), "CertRefs"},
		{"COMPLETE_CERTIFICATE_REFS_V2", XAdES141ElementCompleteCertificateRefsV2.TagName(), "CompleteCertificateRefsV2"},
		{"NEW_SDO_DIGEST_VALUE", XAdES141ElementNewSDODigestValue.TagName(), "NewSDODigestValue"},
		{"ORIGINAL_REF_DIGEST", XAdES141ElementOriginalRefDigest.TagName(), "OriginalRefDigest"},
		{"RECOMPUTED_DIGEST_VALUE", XAdES141ElementRecomputedDigestValue.TagName(), "RecomputedDigestValue"},
		{"REFS_ONLY_TIMESTAMP_V2", XAdES141ElementRefsOnlyTimestampV2.TagName(), "RefsOnlyTimeStampV2"},
		{"RENEWED_DIGESTS", XAdES141ElementRenewedDigests.TagName(), "RenewedDigests"},
		{"RENEWED_DIGESTS_V2", XAdES141ElementRenewedDigestsV2.TagName(), "RenewedDigestsV2"},
		{"SIG_AND_REFS_TIMESTAMP_V2", XAdES141ElementSigAndRefsTimestampV2.TagName(), "SigAndRefsTimeStampV2"},
		{"SIGNATURE_POLICY_DOCUMENT", XAdES141ElementSignaturePolicyDocument.TagName(), "SignaturePolicyDocument"},
		{"SIGNATURE_POLICY_STORE", XAdES141ElementSignaturePolicyStore.TagName(), "SignaturePolicyStore"},
		{"SIG_POL_DOC_LOCAL_URI", XAdES141ElementSigPolDocLocalURI.TagName(), "SigPolDocLocalURI"},
		{"SP_DOC_SPECIFICATION", XAdES141ElementSPDocSpecification.TagName(), "SPDocSpecification"},
		{"TIMESTAMP_VALIDATION_DATA", XAdES141ElementTimestampValidationData.TagName(), "TimeStampValidationData"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.TagName() = %q, want %q", c.name, c.got, c.want)
		}
	}
	if got := XAdES141ElementSPDocSpecification.URI(); got != XAdESNamespaceXAdES141.Uri() {
		t.Errorf("URI() = %q, want %q", got, XAdESNamespaceXAdES141.Uri())
	}
}
