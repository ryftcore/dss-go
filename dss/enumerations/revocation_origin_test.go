package enumerations

import "testing"

func TestRevocationOriginValues(t *testing.T) {
	want := []RevocationOrigin{
		RevocationOriginCMSSignedData,
		RevocationOriginRevocationValues,
		RevocationOriginAttributeRevocationValues,
		RevocationOriginTimestampValidationData,
		RevocationOriginAnyValidationData,
		RevocationOriginDSSDictionary,
		RevocationOriginVRIDictionary,
		RevocationOriginAdbeRevocationInfoArchival,
		RevocationOriginEvidenceRecord,
		RevocationOriginInputDocument,
		RevocationOriginExternal,
		RevocationOriginCached,
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
		want := v != RevocationOriginExternal && v != RevocationOriginCached
		if got := v.IsInternalOrigin(); got != want {
			t.Errorf("%v.IsInternalOrigin() = %v, want %v", v, got, want)
		}
	}
}
