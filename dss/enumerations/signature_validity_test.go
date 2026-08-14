// Ported from dss-enumerations/.../SignatureValidity.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestSignatureValidityGet(t *testing.T) {
	trueVal := true
	falseVal := false
	if got := SignatureValidityGet(nil); got != SignatureValidity_NOT_EVALUATED {
		t.Errorf("SignatureValidityGet(nil) = %v, want %v", got, SignatureValidity_NOT_EVALUATED)
	}
	if got := SignatureValidityGet(&trueVal); got != SignatureValidity_VALID {
		t.Errorf("SignatureValidityGet(true) = %v, want %v", got, SignatureValidity_VALID)
	}
	if got := SignatureValidityGet(&falseVal); got != SignatureValidity_INVALID {
		t.Errorf("SignatureValidityGet(false) = %v, want %v", got, SignatureValidity_INVALID)
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
