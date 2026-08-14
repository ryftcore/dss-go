// Ported from dss-enumerations/.../TextWrapping.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestTextWrappingValueOf(t *testing.T) {
	if len(TextWrappingValues()) != 3 {
		t.Fatalf("expected 3 values, got %d", len(TextWrappingValues()))
	}
	for _, v := range TextWrappingValues() {
		got, err := TextWrappingValueOf(string(v))
		if err != nil || got != v {
			t.Errorf("TextWrappingValueOf(%q) = %v, %v; want %v, nil", v, got, err, v)
		}
	}
	if _, err := TextWrappingValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
