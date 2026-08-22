package enumerations

import "testing"

func TestSignerTextVerticalAlignmentValues(t *testing.T) {
	want := []SignerTextVerticalAlignment{
		SignerTextVerticalAlignmentTop, SignerTextVerticalAlignmentMiddle, SignerTextVerticalAlignmentBottom,
	}
	got := SignerTextVerticalAlignmentValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: got %v, want %v", i, got[i], w)
		}
	}
}
