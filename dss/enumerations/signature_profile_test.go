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
	if got != SignatureProfileBaselineLTA {
		t.Errorf("SignatureProfileValueByName(BASELINE-LTA) = %q, want %q", got, SignatureProfileBaselineLTA)
	}
}

func TestSignatureProfileString(t *testing.T) {
	if got := SignatureProfileBaselineLTA.String(); got != "BASELINE-LTA" {
		t.Errorf("SignatureProfileBaselineLTA.String() = %q, want %q", got, "BASELINE-LTA")
	}
	if got := SignatureProfileAdES.String(); got != "AdES" {
		t.Errorf("SignatureProfileAdES.String() = %q, want %q", got, "AdES")
	}
}
