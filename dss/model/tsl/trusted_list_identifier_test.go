package tsl

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

func TestTrustedListIdentifierAsXmlID(t *testing.T) {
	tlInfo := NewTLInfo(nil, nil, nil, "https://example.org/tl.xml")
	id := NewTrustedListIdentifier(tlInfo)

	sum := sha256.Sum256([]byte("https://example.org/tl.xml"))
	expected := "TL-" + strings.ToUpper(hex.EncodeToString(sum[:]))
	if id.AsXmlID() != expected {
		t.Fatalf("AsXmlID() = %s, want %s", id.AsXmlID(), expected)
	}
}

func TestTrustedListIdentifierString(t *testing.T) {
	tlInfo := NewTLInfo(nil, nil, nil, "https://example.org/tl.xml")
	id := NewTrustedListIdentifier(tlInfo)

	if !strings.HasPrefix(id.String(), "TrustedListIdentifier:") {
		t.Fatalf("String() = %s, want TrustedListIdentifier:... prefix", id.String())
	}
}

func TestTrustedListIdentifierEqualsIsUrlBased(t *testing.T) {
	a := NewTrustedListIdentifier(NewTLInfo(nil, nil, nil, "https://example.org/a.xml"))
	b := NewTrustedListIdentifier(NewTLInfo(nil, nil, nil, "https://example.org/a.xml"))
	c := NewTrustedListIdentifier(NewTLInfo(nil, nil, nil, "https://example.org/c.xml"))

	if !a.Equals(b) {
		t.Fatalf("expected identifiers over the same URL to be equal")
	}
	if a.Equals(c) {
		t.Fatalf("expected identifiers over different URLs to be unequal")
	}

	// A TrustedListIdentifier never equals a LOTLIdentifier, even for the same URL: Java's
	// Identifier#equals() compares getClass() first.
	l := NewLOTLIdentifier(&LOTLInfo{TLInfo: *NewTLInfo(nil, nil, nil, "https://example.org/a.xml")})
	if a.Equals(l) {
		t.Fatalf("expected a TrustedListIdentifier to never equal a LOTLIdentifier")
	}
}
