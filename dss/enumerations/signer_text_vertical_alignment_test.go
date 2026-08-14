package enumerations

import "testing"

func TestSignerTextVerticalAlignmentValues(t *testing.T) {
	want := []SignerTextVerticalAlignment{
		SignerTextVerticalAlignment_TOP, SignerTextVerticalAlignment_MIDDLE, SignerTextVerticalAlignment_BOTTOM,
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
