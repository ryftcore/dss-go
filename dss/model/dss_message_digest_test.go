package model

import (
	"bytes"
	"strings"
	"testing"

	"github.com/utain/esig/dss/enumerations"
)

func TestDSSMessageDigestRoundTripAndString(t *testing.T) {
	value := []byte{0x1, 0x2, 0x3}
	md := NewDSSMessageDigestWithValue(enumerations.DigestAlgorithm_SHA256, value)

	if md.Algorithm() != enumerations.DigestAlgorithm_SHA256 {
		t.Fatalf("Algorithm() = %v", md.Algorithm())
	}
	if !bytes.Equal(md.Value(), value) {
		t.Fatalf("Value() = %x, want %x", md.Value(), value)
	}
	if !strings.HasPrefix(md.String(), "MessageDigest [") {
		t.Fatalf("String() = %q, want prefix 'MessageDigest ['", md.String())
	}
}

func TestCreateEmptyDSSMessageDigest(t *testing.T) {
	md := CreateEmptyDSSMessageDigest()
	if !md.IsEmpty() {
		t.Fatal("expected empty DSSMessageDigest to report IsEmpty()")
	}
}
