package enumerations

import "testing"

func TestVisualSignatureRotationValues(t *testing.T) {
	want := []VisualSignatureRotation{
		VisualSignatureRotationNone,
		VisualSignatureRotationAutomatic,
		VisualSignatureRotationRotate90,
		VisualSignatureRotationRotate180,
		VisualSignatureRotationRotate270,
	}
	got := VisualSignatureRotationValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("values[%d] = %v, want %v", i, got[i], w)
		}
	}
}
