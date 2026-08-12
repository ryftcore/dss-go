package lote

import (
	"crypto/sha256"
	"testing"

	"github.com/utain/esig/dss/model"
)

func TestLoTEIdentifierDigestMatchesURL(t *testing.T) {
	listInfo := NewLoTEInfo(nil, nil, nil, "https://example.org/lote.xml")
	id := NewLoTEIdentifier(listInfo)

	want := sha256.Sum256([]byte("https://example.org/lote.xml"))
	if id.DigestID().HexValue() != model.NewDigest(model.IdentifierDigestAlgorithm, want[:]).HexValue() {
		t.Fatalf("unexpected digest: %s", id.DigestID().HexValue())
	}
	if got := id.AsXmlID(); got[:5] != "LoTE-" {
		t.Fatalf("unexpected prefix in AsXmlID(): %s", got)
	}
	if id.String()[:14] != "LoTEIdentifier" {
		t.Fatalf("unexpected String(): %s", id.String())
	}
}

func TestLoLoTEIdentifierPrefixAndClassName(t *testing.T) {
	listInfo := NewLoTEInfo(nil, nil, nil, "https://example.org/lolote.xml")
	id := NewLoLoTEIdentifier(listInfo)

	if got := id.AsXmlID(); got[:7] != "LoLoTE-" {
		t.Fatalf("unexpected prefix in AsXmlID(): %s", got)
	}
	if id.String()[:16] != "LoLoTEIdentifier" {
		t.Fatalf("unexpected String(): %s", id.String())
	}
}

func TestLoTEAndLoLoTEIdentifiersOfSameURLAreNotEqual(t *testing.T) {
	// Java's Identifier#equals() compares getClass(), so a LoTEIdentifier and a
	// LoLoTEIdentifier built from the same URL must never be equal.
	listInfo := NewLoTEInfo(nil, nil, nil, "https://example.org/same.xml")
	lote := NewLoTEIdentifier(listInfo)
	lolote := NewLoLoTEIdentifier(listInfo)

	if lote.Equals(lolote) {
		t.Fatal("LoTEIdentifier must not equal a LoLoTEIdentifier built from the same URL")
	}
}
