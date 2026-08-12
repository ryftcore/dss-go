package enumerations

import "testing"

func TestValidationModelValues(t *testing.T) {
	want := []ValidationModel{
		ValidationModel_SHELL,
		ValidationModel_CHAIN,
		ValidationModel_HYBRID,
	}
	got := ValidationModelValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("values[%d] = %v, want %v", i, got[i], w)
		}
	}
}
