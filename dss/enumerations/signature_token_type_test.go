package enumerations

import "testing"

func TestSignatureTokenTypeValues(t *testing.T) {
	want := []SignatureTokenType{
		SignatureTokenTypePKCS11, SignatureTokenTypePKCS12, SignatureTokenTypeMSCAPI,
		SignatureTokenTypeApple, SignatureTokenTypeJKS, SignatureTokenTypeMOCCA,
	}
	got := SignatureTokenTypeValues()
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: got %v, want %v", i, got[i], w)
		}
	}
}
