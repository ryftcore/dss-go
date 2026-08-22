package enumerations

import "testing"

func TestValidationTimeValues(t *testing.T) {
	want := []ValidationTime{
		ValidationTimeCertificateIssuanceTime,
		ValidationTimeBESTSignatureTime,
		ValidationTimeValidationTime,
		ValidationTimeTimestampGenerationTime,
		ValidationTimeTimestampPOETime,
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
