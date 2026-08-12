package model

import (
	"crypto/sha256"
	"testing"

	"github.com/utain/esig/dss/enumerations"
)

func TestDigestHexValue(t *testing.T) {
	// Known SHA-256 of the empty string:
	// e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
	sum := sha256.Sum256(nil)
	d := NewDigest(enumerations.DigestAlgorithm_SHA256, sum[:])

	want := "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855"
	if got := d.HexValue(); got != want {
		t.Fatalf("HexValue() = %q, want %q", got, want)
	}
}

func TestDigestHexValueLeadingZeroByte(t *testing.T) {
	// A digest value starting with a 0x00 byte: BigInteger(1, value)
	// normalizes away leading zero bytes, so the hex string must not
	// contain them either (verifying Java-parity, not naive hex.Encode).
	value := []byte{0x00, 0xAB, 0xCD}
	d := NewDigest(enumerations.DigestAlgorithm_SHA1, value)

	want := "ABCD"
	if got := d.HexValue(); got != want {
		t.Fatalf("HexValue() = %q, want %q", got, want)
	}
}

func TestDigestHexValueOddLengthPadded(t *testing.T) {
	// 0x0F alone would render as a single hex digit "f"; Java pads to
	// even length.
	value := []byte{0x0F}
	d := NewDigest(enumerations.DigestAlgorithm_SHA1, value)

	want := "0F"
	if got := d.HexValue(); got != want {
		t.Fatalf("HexValue() = %q, want %q", got, want)
	}
}

func TestDigestBase64Value(t *testing.T) {
	value := []byte("hello")
	d := NewDigest(enumerations.DigestAlgorithm_SHA256, value)

	want := "aGVsbG8="
	if got := d.Base64Value(); got != want {
		t.Fatalf("Base64Value() = %q, want %q", got, want)
	}
}

func TestDigestIsEmpty(t *testing.T) {
	var d Digest
	if !d.IsEmpty() {
		t.Fatal("zero-value Digest should be empty")
	}
	d = NewDigest(enumerations.DigestAlgorithm_SHA256, []byte{1, 2, 3})
	if d.IsEmpty() {
		t.Fatal("populated Digest should not be empty")
	}
}

func TestDigestEqualsAndString(t *testing.T) {
	a := NewDigest(enumerations.DigestAlgorithm_SHA256, []byte{1, 2, 3})
	b := NewDigest(enumerations.DigestAlgorithm_SHA256, []byte{1, 2, 3})
	c := NewDigest(enumerations.DigestAlgorithm_SHA1, []byte{1, 2, 3})

	if !a.Equals(b) {
		t.Fatal("expected equal digests to be Equals()")
	}
	if a.Equals(c) {
		t.Fatal("expected different algorithms to not be Equals()")
	}

	want := "SHA256:#010203"
	if got := a.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	var empty Digest
	if got, want := empty.String(), "?:?"; got != want {
		t.Fatalf("empty String() = %q, want %q", got, want)
	}
}

func TestDigestHexValuePanicsOnNilValue(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil digest value")
		}
	}()
	var d Digest
	d.SetAlgorithm(enumerations.DigestAlgorithm_SHA256)
	d.HexValue()
}
