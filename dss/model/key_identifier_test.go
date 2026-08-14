package model

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestKeyIdentifierDigestsTheEncodedKey(t *testing.T) {
	key := NewPublicKeyFromEncoded([]byte("encoded key material"), nil)
	identifier := NewKeyIdentifier(key)

	want := sha256.Sum256(key.Encoded())
	if got := hex.EncodeToString(identifier.DigestID().Value()); got != hex.EncodeToString(want[:]) {
		t.Errorf("DigestID().Value() = %s", got)
	}
	if got := identifier.AsXmlID()[:3]; got != "PK-" {
		t.Errorf("AsXmlID() prefix = %q", got)
	}
	if got := identifier.String()[:14]; got != "KeyIdentifier:" {
		t.Errorf("String() = %q", identifier.String())
	}
}
