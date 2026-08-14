package model

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestEntityIdentifier(t *testing.T) {
	identifier := NewEntityIdentifier([]byte("DSS"))
	want := sha256.Sum256([]byte("DSS"))

	if got := identifier.AsXmlID(); got != "EK-B121C379F13BB66F3EBB9021E3F3EBABC115564CBEFF65994C325C7005955DCF" {
		t.Errorf("AsXmlID() = %q", got)
	}
	if got := hex.EncodeToString(identifier.DigestID().Value()); got != hex.EncodeToString(want[:]) {
		t.Errorf("DigestID().Value() = %s", got)
	}
	if got := identifier.String(); got != "EntityIdentifier:SHA256:#B121C379F13BB66F3EBB9021E3F3EBABC115564CBEFF65994C325C7005955DCF" {
		t.Errorf("String() = %q", got)
	}
}
