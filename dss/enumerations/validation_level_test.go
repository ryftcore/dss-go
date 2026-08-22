package enumerations

import "testing"

func TestValidationLevelValueOf(t *testing.T) {
	for _, v := range ValidationLevelValues() {
		got, err := ValidationLevelValueOf(string(v))
		if err != nil {
			t.Fatalf("ValidationLevelValueOf(%q) returned error: %v", v, err)
		}
		if got != v {
			t.Errorf("ValidationLevelValueOf(%q) = %q, want %q", v, got, v)
		}
	}

	if _, err := ValidationLevelValueOf("bogus"); err == nil {
		t.Error("ValidationLevelValueOf(\"bogus\") expected error, got nil")
	}
}

func TestValidationLevelValues(t *testing.T) {
	want := []ValidationLevel{
		ValidationLevelBasicSignatures,
		ValidationLevelTimestamps,
		ValidationLevelLongTermData,
		ValidationLevelArchivalData,
	}
	got := ValidationLevelValues()
	if len(got) != len(want) {
		t.Fatalf("ValidationLevelValues() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("ValidationLevelValues()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
