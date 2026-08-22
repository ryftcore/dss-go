// Ported from dss-enumerations/.../SignatureValidity.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestSignatureValidityGet(t *testing.T) {
	trueVal := true
	falseVal := false
	if got := SignatureValidityGet(nil); got != SignatureValidityNotEvaluated {
		t.Errorf("SignatureValidityGet(nil) = %v, want %v", got, SignatureValidityNotEvaluated)
	}
	if got := SignatureValidityGet(&trueVal); got != SignatureValidityValid {
		t.Errorf("SignatureValidityGet(true) = %v, want %v", got, SignatureValidityValid)
	}
	if got := SignatureValidityGet(&falseVal); got != SignatureValidityInvalid {
		t.Errorf("SignatureValidityGet(false) = %v, want %v", got, SignatureValidityInvalid)
	}
}

func TestSignatureValidityValueOf(t *testing.T) {
	if len(SignatureValidityValues()) != 3 {
		t.Fatalf("expected 3 values, got %d", len(SignatureValidityValues()))
	}
	for _, v := range SignatureValidityValues() {
		got, err := SignatureValidityValueOf(string(v))
		if err != nil || got != v {
			t.Errorf("SignatureValidityValueOf(%q) = %v, %v; want %v, nil", v, got, err, v)
		}
	}
	if _, err := SignatureValidityValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
