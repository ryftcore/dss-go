package enumerations

import "testing"

func TestSignatureProfileValueOf(t *testing.T) {
	for _, v := range SignatureProfileValues() {
		got, err := SignatureProfileValueOf(string(v))
		if err != nil {
			t.Errorf("SignatureProfileValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("SignatureProfileValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestSignatureProfileValueOfUnknown(t *testing.T) {
	if _, err := SignatureProfileValueOf("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}

func TestSignatureProfileValueByName(t *testing.T) {
	got, err := SignatureProfileValueByName("BASELINE-LTA")
	if err != nil {
		t.Fatalf("SignatureProfileValueByName unexpected error: %v", err)
	}
	if got != SignatureProfile_BASELINE_LTA {
		t.Errorf("SignatureProfileValueByName(BASELINE-LTA) = %q, want %q", got, SignatureProfile_BASELINE_LTA)
	}
}

func TestSignatureProfileString(t *testing.T) {
	if got := SignatureProfile_BASELINE_LTA.String(); got != "BASELINE-LTA" {
		t.Errorf("SignatureProfile_BASELINE_LTA.String() = %q, want %q", got, "BASELINE-LTA")
	}
	if got := SignatureProfile_AdES.String(); got != "AdES" {
		t.Errorf("SignatureProfile_AdES.String() = %q, want %q", got, "AdES")
	}
}
