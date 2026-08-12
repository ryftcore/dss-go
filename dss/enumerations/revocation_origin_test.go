package enumerations

import "testing"

func TestRevocationOriginValues(t *testing.T) {
	want := []RevocationOrigin{
		RevocationOrigin_CMS_SIGNED_DATA,
		RevocationOrigin_REVOCATION_VALUES,
		RevocationOrigin_ATTRIBUTE_REVOCATION_VALUES,
		RevocationOrigin_TIMESTAMP_VALIDATION_DATA,
		RevocationOrigin_ANY_VALIDATION_DATA,
		RevocationOrigin_DSS_DICTIONARY,
		RevocationOrigin_VRI_DICTIONARY,
		RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL,
		RevocationOrigin_EVIDENCE_RECORD,
		RevocationOrigin_INPUT_DOCUMENT,
		RevocationOrigin_EXTERNAL,
		RevocationOrigin_CACHED,
	}
	got := RevocationOriginValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("values[%d] = %v, want %v", i, got[i], w)
		}
	}
}

func TestRevocationOriginIsInternalOrigin(t *testing.T) {
	for _, v := range RevocationOriginValues() {
		want := v != RevocationOrigin_EXTERNAL && v != RevocationOrigin_CACHED
		if got := v.IsInternalOrigin(); got != want {
			t.Errorf("%v.IsInternalOrigin() = %v, want %v", v, got, want)
		}
	}
}
