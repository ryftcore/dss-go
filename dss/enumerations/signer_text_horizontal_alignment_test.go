package enumerations

import "testing"

func TestSignerTextHorizontalAlignmentValues(t *testing.T) {
	want := []SignerTextHorizontalAlignment{
		SignerTextHorizontalAlignmentLeft,
		SignerTextHorizontalAlignmentCenter,
		SignerTextHorizontalAlignmentRight,
	}
	got := SignerTextHorizontalAlignmentValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("values[%d] = %v, want %v", i, got[i], w)
		}
	}
}
