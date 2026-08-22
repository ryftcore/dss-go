package token

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// TestDigestInfoEncoderEncodeKAT is a known-answer test: the standard PKCS#1 v1.5 DigestInfo
// prefix for SHA-256 (RFC 8017 Note 1) followed by SHA-256("abc"), independently verified via
// `openssl dgst -sha256` and a hand-built ASN.1 prefix.
func TestDigestInfoEncoderEncodeKAT(t *testing.T) {
	digest := sha256.Sum256([]byte("abc"))
	want, err := hex.DecodeString("3031300d060960864801650304020105000420" +
		"ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad")
	if err != nil {
		t.Fatalf("bad hex fixture: %v", err)
	}

	got, err := DigestInfoEncoderEncode(enumerations.DigestAlgorithm_SHA256.OID(), digest[:])
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("Encode mismatch:\n got  %x\n want %x", got, want)
	}

	if !DigestInfoEncoderIsEncoded(got) {
		t.Fatalf("IsEncoded(encoded DigestInfo) = false, want true")
	}
}

func TestDigestInfoEncoderIsEncoded(t *testing.T) {
	digest := sha256.Sum256([]byte("abc"))
	encoded, err := DigestInfoEncoderEncode(enumerations.DigestAlgorithm_SHA1.OID(), digest[:20])
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}

	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{"encoded DigestInfo", encoded, true},
		{"nil", nil, false},
		{"too short", []byte{0x30, 0x03, 0x30, 0x01}, false},
		{"raw digest, not encoded", digest[:], false},
		{"not a sequence", append([]byte{0x31}, encoded[1:]...), false},
		{"truncated", encoded[:len(encoded)-1], false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DigestInfoEncoderIsEncoded(tt.data); got != tt.want {
				t.Errorf("IsEncoded(%x) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestDigestInfoEncoderEncodePanicsOnNilInputs(t *testing.T) {
	t.Run("empty OID", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic")
			}
		}()
		_, _ = DigestInfoEncoderEncode("", []byte{1})
	})
	t.Run("nil digest", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic")
			}
		}()
		_, _ = DigestInfoEncoderEncode(enumerations.DigestAlgorithm_SHA256.OID(), nil)
	})
}

func TestDigestInfoEncoderEncodeInvalidOID(t *testing.T) {
	if _, err := DigestInfoEncoderEncode("not-an-oid", []byte{1, 2, 3}); err == nil {
		t.Fatal("expected an error for a malformed OID")
	}
}
