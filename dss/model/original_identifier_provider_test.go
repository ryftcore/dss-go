package model

import "testing"

func TestOriginalIdentifierProviderReturnsTheXmlID(t *testing.T) {
	provider := NewOriginalIdentifierProvider()
	token := certificateTokenFixture(t, rootCertificateBase64)

	if got, want := provider.IDAsString(token), token.DSSIDAsString(); got != want {
		t.Errorf("IDAsString() = %q, want %q", got, want)
	}
	if got := provider.IDAsString(token); got != "C-E46CC70B1A54986D0E07A28BA9C99468115F71F305894A280CC08DADC31E4AF3" {
		t.Errorf("IDAsString() = %q", got)
	}
}
