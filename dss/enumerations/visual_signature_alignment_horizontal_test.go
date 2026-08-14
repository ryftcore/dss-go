package enumerations

import "testing"

func TestVisualSignatureAlignmentHorizontalValueOf(t *testing.T) {
	for _, v := range VisualSignatureAlignmentHorizontalValues() {
		got, err := VisualSignatureAlignmentHorizontalValueOf(string(v))
		if err != nil {
			t.Errorf("VisualSignatureAlignmentHorizontalValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("VisualSignatureAlignmentHorizontalValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestVisualSignatureAlignmentHorizontalValueOfUnknown(t *testing.T) {
	if _, err := VisualSignatureAlignmentHorizontalValueOf("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
