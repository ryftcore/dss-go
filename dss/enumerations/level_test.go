package enumerations

import "testing"

func TestLevelValueOf(t *testing.T) {
	for _, v := range LevelValues() {
		got, err := LevelValueOf(string(v))
		if err != nil {
			t.Errorf("LevelValueOf(%q) unexpected error: %v", v, err)
		}
		if got != v {
			t.Errorf("LevelValueOf(%q) = %q, want %q", v, got, v)
		}
	}
}

func TestLevelValueOfUnknown(t *testing.T) {
	if _, err := LevelValueOf("bogus"); err == nil {
		t.Error("expected error for unknown value, got nil")
	}
}
