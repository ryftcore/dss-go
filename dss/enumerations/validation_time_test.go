package enumerations

import "testing"

func TestValidationTimeValues(t *testing.T) {
	want := []ValidationTime{
		ValidationTime_CERTIFICATE_ISSUANCE_TIME,
		ValidationTime_BEST_SIGNATURE_TIME,
		ValidationTime_VALIDATION_TIME,
		ValidationTime_TIMESTAMP_GENERATION_TIME,
		ValidationTime_TIMESTAMP_POE_TIME,
	}
	got := ValidationTimeValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: got %v, want %v", i, got[i], w)
		}
	}
}
