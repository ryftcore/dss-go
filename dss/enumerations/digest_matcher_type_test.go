package enumerations

import "testing"

func TestDigestMatcherTypeValues(t *testing.T) {
	want := []DigestMatcherType{
		DigestMatcherTypeReference,
		DigestMatcherTypeObject,
		DigestMatcherTypeManifest,
		DigestMatcherTypeSignedProperties,
		DigestMatcherTypeKeyInfo,
		DigestMatcherTypeSignatureProperties,
		DigestMatcherTypeXPointer,
		DigestMatcherTypeManifestEntry,
		DigestMatcherTypeCounterSignature,
		DigestMatcherTypeMessageDigest,
		DigestMatcherTypeContentDigest,
		DigestMatcherTypeJWSSigningInput,
		DigestMatcherTypeSigDEntry,
		DigestMatcherTypeCoseSigStructure,
		DigestMatcherTypeCounterSignedSignatureValue,
		DigestMatcherTypeMessageImprint,
		DigestMatcherTypeEvidenceRecordArchiveObject,
		DigestMatcherTypeEvidenceRecordOrphanReference,
		DigestMatcherTypeEvidenceRecordArchiveTimeStamp,
		DigestMatcherTypeEvidenceRecordArchiveTimeStampSequence,
		DigestMatcherTypeEvidenceRecordMasterSignature,
		DigestMatcherTypeEAADisclosure,
		DigestMatcherTypeEAANestedDisclosure,
		DigestMatcherTypeEAAOrphanSelectivelyDisclosableClaim,
		DigestMatcherTypeEAAKeyBinding,
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
