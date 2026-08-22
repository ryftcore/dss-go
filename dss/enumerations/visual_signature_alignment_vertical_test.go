package enumerations

import "testing"

func TestVisualSignatureAlignmentVerticalValues(t *testing.T) {
	want := []VisualSignatureAlignmentVertical{
		VisualSignatureAlignmentVerticalNone,
		VisualSignatureAlignmentVerticalTop,
		VisualSignatureAlignmentVerticalMiddle,
		VisualSignatureAlignmentVerticalBottom,
	}
	got := VisualSignatureAlignmentVerticalValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("values[%d] = %v, want %v", i, got[i], w)
		}
	}
}
