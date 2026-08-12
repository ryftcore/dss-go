package enumerations

import "testing"

func TestDigestMatcherTypeValues(t *testing.T) {
	want := []DigestMatcherType{
		DigestMatcherType_REFERENCE,
		DigestMatcherType_OBJECT,
		DigestMatcherType_MANIFEST,
		DigestMatcherType_SIGNED_PROPERTIES,
		DigestMatcherType_KEY_INFO,
		DigestMatcherType_SIGNATURE_PROPERTIES,
		DigestMatcherType_XPOINTER,
		DigestMatcherType_MANIFEST_ENTRY,
		DigestMatcherType_COUNTER_SIGNATURE,
		DigestMatcherType_MESSAGE_DIGEST,
		DigestMatcherType_CONTENT_DIGEST,
		DigestMatcherType_JWS_SIGNING_INPUT,
		DigestMatcherType_SIG_D_ENTRY,
		DigestMatcherType_COSE_SIG_STRUCTURE,
		DigestMatcherType_COUNTER_SIGNED_SIGNATURE_VALUE,
		DigestMatcherType_MESSAGE_IMPRINT,
		DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_OBJECT,
		DigestMatcherType_EVIDENCE_RECORD_ORPHAN_REFERENCE,
		DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP,
		DigestMatcherType_EVIDENCE_RECORD_ARCHIVE_TIME_STAMP_SEQUENCE,
		DigestMatcherType_EVIDENCE_RECORD_MASTER_SIGNATURE,
		DigestMatcherType_EAA_DISCLOSURE,
		DigestMatcherType_EAA_NESTED_DISCLOSURE,
		DigestMatcherType_EAA_ORPHAN_SELECTIVELY_DISCLOSABLE_CLAIM,
		DigestMatcherType_EAA_KEY_BINDING,
	}
	got := DigestMatcherTypeValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	seen := map[DigestMatcherType]bool{}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("values[%d] = %v, want %v", i, got[i], w)
		}
		if seen[w] {
			t.Errorf("duplicate value %v", w)
		}
		seen[w] = true
	}
}
