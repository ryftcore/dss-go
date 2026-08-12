package enumerations

import "testing"

func TestASiCManifestTypeEnumValueOf(t *testing.T) {
	for _, v := range ASiCManifestTypeEnumValues() {
		got, err := ASiCManifestTypeEnumValueOf(string(v))
		if err != nil {
			t.Errorf("ASiCManifestTypeEnumValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("ASiCManifestTypeEnumValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestASiCManifestTypeEnumValueOfUnknown(t *testing.T) {
	if _, err := ASiCManifestTypeEnumValueOf("NOT_A_VALUE"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
