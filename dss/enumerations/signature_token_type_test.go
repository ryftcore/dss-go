package enumerations

import "testing"

func TestSignatureTokenTypeValues(t *testing.T) {
	want := []SignatureTokenType{
		SignatureTokenType_PKCS11, SignatureTokenType_PKCS12, SignatureTokenType_MSCAPI,
		SignatureTokenType_APPLE, SignatureTokenType_JKS, SignatureTokenType_MOCCA,
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
