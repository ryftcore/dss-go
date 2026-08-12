package enumerations

import "testing"

func TestSubContextValueOf(t *testing.T) {
	for _, v := range SubContextValues() {
		got, err := SubContextValueOf(string(v))
		if err != nil {
			t.Errorf("SubContextValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("SubContextValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestSubContextValueOfUnknown(t *testing.T) {
	if _, err := SubContextValueOf("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
