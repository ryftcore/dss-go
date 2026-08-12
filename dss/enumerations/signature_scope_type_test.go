package enumerations

import "testing"

func TestSignatureScopeTypeValueOf(t *testing.T) {
	for _, v := range SignatureScopeTypeValues() {
		got, err := SignatureScopeTypeValueOf(string(v))
		if err != nil {
			t.Errorf("SignatureScopeTypeValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("SignatureScopeTypeValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestSignatureScopeTypeValueOfUnknown(t *testing.T) {
	if _, err := SignatureScopeTypeValueOf("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
