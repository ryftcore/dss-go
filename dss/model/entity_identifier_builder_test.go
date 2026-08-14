package model

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestEntityIdentifierBuilderConcatenatesKeyThenName(t *testing.T) {
	principal, err := NewX500Principal(mustHex(t, "3014311230100603550403130954657374204e616d65"))
	if err != nil {
		t.Fatal(err)
	}
	key := NewPublicKeyFromEncoded([]byte("SubjectPublicKeyInfo"), nil)

	builder := NewEntityIdentifierBuilder(key, principal)
	want := append(append([]byte{}, key.Encoded()...), principal.Encoded()...)
	if got := hex.EncodeToString(builder.BuildBinaries()); got != hex.EncodeToString(want) {
		t.Errorf("BuildBinaries() = %s, want %s", got, hex.EncodeToString(want))
	}
	digest := sha256.Sum256(want)
	identifier := builder.Build()
	if got := hex.EncodeToString(identifier.DigestID().Value()); got != hex.EncodeToString(digest[:]) {
		t.Errorf("Build() digest = %s", got)
	}
	if got := identifier.AsXmlID()[:3]; got != "EK-" {
		t.Errorf("AsXmlID() prefix = %q", got)
	}
}

func TestEntityIdentifierBuilderSkipsMissingParts(t *testing.T) {
	principal, err := NewX500Principal(mustHex(t, "3014311230100603550403130954657374204e616d65"))
	if err != nil {
		t.Fatal(err)
	}
	key := NewPublicKeyFromEncoded([]byte("spki"), nil)

	// Either part may be absent; upstream simply writes nothing for it.
	if got := string(NewEntityIdentifierBuilder(nil, principal).BuildBinaries()); got != string(principal.Encoded()) {
		t.Errorf("a nil key must contribute nothing, got %x", got)
	}
	if got := string(NewEntityIdentifierBuilder(key, nil).BuildBinaries()); got != "spki" {
		t.Errorf("a nil name must contribute nothing, got %q", got)
	}
	if got := NewEntityIdentifierBuilder(nil, nil).BuildBinaries(); len(got) != 0 {
		t.Errorf("both parts absent must build empty binaries, got %x", got)
	}
}
