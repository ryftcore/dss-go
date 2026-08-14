package model

import (
	"encoding/hex"
	"testing"

	"github.com/utain/esig/dss/enumerations"
)

func TestMultipleDigestIdentifierComputesAndCachesDigests(t *testing.T) {
	identifier := NewMultipleDigestIdentifier("TimestampTokenIdentifier", "T-", []byte("abc"))

	if got := hex.EncodeToString(identifier.Binaries()); got != "616263" {
		t.Errorf("Binaries() = %s", got)
	}
	// The SHA-256 digest is pre-populated from the identifier itself.
	if got := identifier.AsXmlID(); got != "T-BA7816BF8F01CFEA414140DE5DAE2223B00361A396177A9CB410FF61F20015AD" {
		t.Errorf("AsXmlID() = %q", got)
	}

	// Known answers for "abc".
	tests := map[enumerations.DigestAlgorithm]string{
		enumerations.DigestAlgorithm_SHA224: "23097d223405d8228642a477bda255b32aadbce4bda0b3f7e36c9da7",
		enumerations.DigestAlgorithm_SHA256: "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
		enumerations.DigestAlgorithm_SHA384: "cb00753f45a35e8bb5a03d699ac65007272c32ab0eded1631a8b605a43ff5bed8086072ba1e7cc2358baeca134c825a7",
		enumerations.DigestAlgorithm_SHA512: "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f",
	}
	for algorithm, want := range tests {
		value, err := identifier.DigestValue(algorithm)
		if err != nil {
			t.Fatalf("DigestValue(%v): %v", algorithm, err)
		}
		if got := hex.EncodeToString(value); got != want {
			t.Errorf("DigestValue(%v) = %s, want %s", algorithm, got, want)
		}
		// The second call is served from the cache and must yield the very same slice.
		cached, err := identifier.DigestValue(algorithm)
		if err != nil {
			t.Fatal(err)
		}
		if &cached[0] != &value[0] {
			t.Errorf("DigestValue(%v) was recomputed instead of cached", algorithm)
		}
	}

	if _, err := identifier.DigestValue(enumerations.DigestAlgorithm_WHIRLPOOL); err == nil {
		t.Error("an algorithm without a Go implementation must be reported")
	}
}

func TestMultipleDigestIdentifierIsMatch(t *testing.T) {
	identifier := NewMultipleDigestIdentifier("TimestampTokenIdentifier", "T-", []byte("abc"))

	sha512, _ := hex.DecodeString("ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f")
	match, err := identifier.IsMatch(NewDigest(enumerations.DigestAlgorithm_SHA512, sha512))
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Error("the SHA-512 digest of the binaries must match")
	}

	match, err = identifier.IsMatch(NewDigest(enumerations.DigestAlgorithm_SHA512, []byte{0x00}))
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Error("a different digest must not match")
	}

	if _, err := identifier.IsMatch(NewDigest(enumerations.DigestAlgorithm_MD2, []byte{0x00})); err == nil {
		t.Error("an algorithm without a Go implementation must be reported")
	}
}
