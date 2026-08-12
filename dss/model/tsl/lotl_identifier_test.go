package tsl

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

func TestLOTLIdentifierAsXmlID(t *testing.T) {
	lotlInfo := NewLOTLInfo(nil, nil, nil, "https://example.org/lotl.xml")
	id := NewLOTLIdentifier(&lotlInfo)

	sum := sha256.Sum256([]byte("https://example.org/lotl.xml"))
	expected := "LOTL-" + strings.ToUpper(hex.EncodeToString(sum[:]))
	if id.AsXmlID() != expected {
		t.Fatalf("AsXmlID() = %s, want %s", id.AsXmlID(), expected)
	}
}

func TestLOTLIdentifierString(t *testing.T) {
	lotlInfo := NewLOTLInfo(nil, nil, nil, "https://example.org/lotl.xml")
	id := NewLOTLIdentifier(&lotlInfo)

	if !strings.HasPrefix(id.String(), "LOTLIdentifier:") {
		t.Fatalf("String() = %s, want LOTLIdentifier:... prefix", id.String())
	}
}
