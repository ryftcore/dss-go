// Ported from dss-enumerations/.../SignerTextPosition.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestSignerTextPositionValueOf(t *testing.T) {
	if len(SignerTextPositionValues()) != 4 {
		t.Fatalf("expected 4 values, got %d", len(SignerTextPositionValues()))
	}
	for _, v := range SignerTextPositionValues() {
		got, err := SignerTextPositionValueOf(string(v))
		if err != nil || got != v {
			t.Errorf("SignerTextPositionValueOf(%q) = %v, %v; want %v, nil", v, got, err, v)
		}
	}
	if _, err := SignerTextPositionValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
