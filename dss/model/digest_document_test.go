package model

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/utain/esig/dss/enumerations"
)

func TestDigestDocumentAddAndGetDigestValue(t *testing.T) {
	d := NewDigestDocument()
	value := []byte{0x01, 0x02, 0x03}
	d.AddDigestValue(enumerations.DigestAlgorithm_SHA256, value)

	got, err := d.DigestValue(enumerations.DigestAlgorithm_SHA256)
	if err != nil {
		t.Fatalf("DigestValue: %v", err)
	}
	if !bytes.Equal(got, value) {
		t.Fatalf("got %x, want %x", got, value)
	}
}

func TestDigestDocumentMissingAlgorithmErrors(t *testing.T) {
	d := NewDigestDocument()
	if _, err := d.DigestValue(enumerations.DigestAlgorithm_SHA256); err == nil {
		t.Fatal("expected error for missing digest algorithm")
	}
}

func TestDigestDocumentAddDigestBase64(t *testing.T) {
	value := []byte("digest-bytes")
	encoded := base64.StdEncoding.EncodeToString(value)

	d := NewDigestDocument()
	if err := d.AddDigestBase64(enumerations.DigestAlgorithm_SHA1, encoded); err != nil {
		t.Fatalf("AddDigestBase64: %v", err)
	}
	got, err := d.DigestValue(enumerations.DigestAlgorithm_SHA1)
	if err != nil {
		t.Fatalf("DigestValue: %v", err)
	}
	if !bytes.Equal(got, value) {
		t.Fatalf("got %x, want %x", got, value)
	}
}

func TestDigestDocumentAddDigestBase64Invalid(t *testing.T) {
	d := NewDigestDocument()
	if err := d.AddDigestBase64(enumerations.DigestAlgorithm_SHA1, "not-valid-base64!!"); err == nil {
		t.Fatal("expected error for invalid base64 digest")
	}
}

func TestDigestDocumentOpenStreamUnsupported(t *testing.T) {
	d := NewDigestDocumentFromValue(enumerations.DigestAlgorithm_SHA256, []byte{0xAA})
	if _, err := d.OpenStream(); err == nil {
		t.Fatal("expected OpenStream to error on a digest-only document")
	}
	if err := d.Save("/tmp/should-not-be-written"); err == nil {
		t.Fatal("expected Save to error on a digest-only document")
	}
}

func TestDigestDocumentExistingDigest(t *testing.T) {
	empty := NewDigestDocument()
	if _, err := empty.ExistingDigest(); err != ErrNoDigest {
		t.Fatalf("expected ErrNoDigest, got %v", err)
	}

	value := []byte{0x9, 0x8, 0x7}
	d := NewDigestDocumentFromValue(enumerations.DigestAlgorithm_SHA256, value)
	digest, err := d.ExistingDigest()
	if err != nil {
		t.Fatalf("ExistingDigest: %v", err)
	}
	if digest.Algorithm() != enumerations.DigestAlgorithm_SHA256 {
		t.Fatalf("Algorithm() = %v, want SHA256", digest.Algorithm())
	}
	if !bytes.Equal(digest.Value(), value) {
		t.Fatalf("Value() = %x, want %x", digest.Value(), value)
	}
}

func TestDigestDocumentFromValueWithNameDerivesMimeType(t *testing.T) {
	d := NewDigestDocumentFromValueWithName(enumerations.DigestAlgorithm_SHA256, []byte{0x1}, "file.txt")
	if d.Name() != "file.txt" {
		t.Fatalf("Name() = %q, want file.txt", d.Name())
	}
	if d.MimeType() == nil {
		t.Fatal("expected non-nil MimeType derived from file name")
	}
}
