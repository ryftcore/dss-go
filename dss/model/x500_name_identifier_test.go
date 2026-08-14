package model

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestX500NameIdentifierDigestsTheEncodedName(t *testing.T) {
	der := mustHex(t, "3014311230100603550403130954657374204e616d65")
	principal, err := NewX500Principal(der)
	if err != nil {
		t.Fatal(err)
	}
	identifier := NewX500NameIdentifier(principal)

	want := sha256.Sum256(der)
	if got := hex.EncodeToString(identifier.DigestID().Value()); got != hex.EncodeToString(want[:]) {
		t.Errorf("DigestID().Value() = %s", got)
	}
	if got := identifier.AsXmlID()[:4]; got != "RDN-" {
		t.Errorf("AsXmlID() prefix = %q", got)
	}
}
