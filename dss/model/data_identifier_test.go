package model

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestDataIdentifierFromBinaries(t *testing.T) {
	identifier := NewDataIdentifier([]byte("abc"))
	want := sha256.Sum256([]byte("abc"))

	if got := hex.EncodeToString(identifier.DigestID().Value()); got != hex.EncodeToString(want[:]) {
		t.Errorf("DigestID().Value() = %s", got)
	}
	if got := identifier.AsXmlID(); got != "D-BA7816BF8F01CFEA414140DE5DAE2223B00361A396177A9CB410FF61F20015AD" {
		t.Errorf("AsXmlID() = %q", got)
	}
}

func TestDataIdentifierWritesTheNameAsUTF16BE(t *testing.T) {
	// DataOutputStream#writeChars writes each Java char as two big-endian bytes, so the name
	// "Aé" contributes 0x0041 0x00E9 before the document digest is appended.
	document := NewInMemoryDocument([]byte("abc"))
	identifier, err := NewDataIdentifierForDocument("Aé", document)
	if err != nil {
		t.Fatal(err)
	}

	documentDigest := sha256.Sum256([]byte("abc"))
	preimage := append([]byte{0x00, 0x41, 0x00, 0xE9}, documentDigest[:]...)
	want := sha256.Sum256(preimage)
	if got := hex.EncodeToString(identifier.DigestID().Value()); got != hex.EncodeToString(want[:]) {
		t.Errorf("DigestID().Value() = %s, want %s", got, hex.EncodeToString(want[:]))
	}
	if got := identifier.AsXmlID(); got != "D-9AE14E59CE7BD7A67172BAE49373ED9E88B605F084B6F5BB7629451D7BAF51EF" {
		t.Errorf("AsXmlID() = %q", got)
	}
}

func TestDataIdentifierBuildParts(t *testing.T) {
	document := NewInMemoryDocument([]byte("abc"))
	documentDigest := sha256.Sum256([]byte("abc"))

	// An empty name contributes nothing, matching Java's null-name branch.
	data, err := dataIdentifierBuild("", document)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(data) != hex.EncodeToString(documentDigest[:]) {
		t.Errorf("build(\"\", doc) = %x", data)
	}

	// A nil document contributes nothing either.
	data, err = dataIdentifierBuild("AB", nil)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(data) != "00410042" {
		t.Errorf("build(\"AB\", nil) = %x", data)
	}

	// Supplementary characters are written as the surrogate pair a Java String holds.
	data, err = dataIdentifierBuild("\U0001D400", nil)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(data) != "d835dc00" {
		t.Errorf("build(supplementary) = %x, want the UTF-16 surrogate pair", data)
	}
}

func TestDataIdentifierUsesAnExistingDigestForDigestDocuments(t *testing.T) {
	// A DigestDocument has no content, so upstream takes the digest it already carries
	// instead of hashing the document.
	value := make([]byte, 32)
	for i := range value {
		value[i] = byte(i + 1)
	}
	document := NewDigestDocumentFromValue(enumerations.DigestAlgorithmSHA256, value)

	identifier, err := NewDataIdentifierForDocument("", document)
	if err != nil {
		t.Fatal(err)
	}
	want := sha256.Sum256(value)
	if got := hex.EncodeToString(identifier.DigestID().Value()); got != hex.EncodeToString(want[:]) {
		t.Errorf("DigestID().Value() = %s, want %s", got, hex.EncodeToString(want[:]))
	}
}

func TestDataIdentifierReportsAMissingExistingDigest(t *testing.T) {
	// getExistingDigest() throws when no digest was added; the Go port returns that error.
	if _, err := NewDataIdentifierForDocument("", NewDigestDocument()); err == nil {
		t.Error("a DigestDocument without a digest must be reported")
	}
}
