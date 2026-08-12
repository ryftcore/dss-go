package model

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestTokenIdentifierFromBinaries(t *testing.T) {
	identifier := NewTokenIdentifier("TimestampTokenIdentifier", "T-", []byte("token binaries"))
	want := sha256.Sum256([]byte("token binaries"))

	if got := hex.EncodeToString(identifier.DigestID().Value()); got != hex.EncodeToString(want[:]) {
		t.Errorf("DigestID().Value() = %s", got)
	}
	if got := identifier.AsXmlID()[:2]; got != "T-" {
		t.Errorf("AsXmlID() prefix = %q", got)
	}
	if got := string(identifier.Binaries()); got != "token binaries" {
		t.Errorf("Binaries() = %q", got)
	}
}

func TestTokenIdentifierFromToken(t *testing.T) {
	// The identifier is computed over the token's encoded form, i.e. the certificate DER.
	token := certificateTokenFixture(t, rootCertificateBase64)
	identifier := NewTokenIdentifierFromToken("CertificateTokenIdentifier", "C-", token)

	want := sha256.Sum256(token.Encoded())
	if got := hex.EncodeToString(identifier.DigestID().Value()); got != hex.EncodeToString(want[:]) {
		t.Errorf("DigestID().Value() = %s", got)
	}
	if got := identifier.AsXmlID(); got != "C-E46CC70B1A54986D0E07A28BA9C99468115F71F305894A280CC08DADC31E4AF3" {
		t.Errorf("AsXmlID() = %q", got)
	}
}
