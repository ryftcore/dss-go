package enumerations

import "testing"

func TestImageScalingValueOf(t *testing.T) {
	for _, v := range ImageScalingValues() {
		got, err := ImageScalingValueOf(string(v))
		if err != nil {
			t.Errorf("ImageScalingValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("ImageScalingValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestImageScalingValueOfUnknown(t *testing.T) {
	if _, err := ImageScalingValueOf("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
