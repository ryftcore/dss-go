// Tests for SignatureValueChecker, matching
// dss-document/src/main/java/eu/europa/esig/dss/signature/SignatureValueChecker.java
// upstream.
package document

import (
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

func TestSignatureValueCheckerEnsureSignatureValueEmptyTargetPanics(t *testing.T) {
	c := NewSignatureValueChecker()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for an empty target SignatureAlgorithm")
		}
	}()
	_, _ = c.EnsureSignatureValue(model.NewSignatureValue(), "")
}

func TestSignatureValueCheckerEnsureSignatureValueNilSignatureValue(t *testing.T) {
	c := NewSignatureValueChecker()
	got, err := c.EnsureSignatureValue(nil, enumerations.SignatureAlgorithm_ECDSA_SHA256)
	if err != nil {
		t.Fatalf("EnsureSignatureValue: %s", err)
	}
	if got != nil {
		t.Fatalf("EnsureSignatureValue(nil, ...) = %v, want nil", got)
	}
}

func TestSignatureValueCheckerEnsureSignatureValueMatchingAlgorithmReturnsSame(t *testing.T) {
	c := NewSignatureValueChecker()
	sv := model.NewSignatureValueWithValue(enumerations.SignatureAlgorithm_ECDSA_SHA256, []byte{1, 2, 3})
	got, err := c.EnsureSignatureValue(sv, enumerations.SignatureAlgorithm_ECDSA_SHA256)
	if err != nil {
		t.Fatalf("EnsureSignatureValue: %s", err)
	}
	if got != sv {
		t.Fatalf("EnsureSignatureValue() = %v, want the identical *SignatureValue back", got)
	}
}

func TestSignatureValueCheckerEnsureSignatureValueDigestMismatchErrors(t *testing.T) {
	c := NewSignatureValueChecker()
	// RSA_SHA256 has a SHA256 digest; targeting RSA_SHA512 mismatches the digest, and RSA is not
	// an EC-equivalent encryption algorithm, so no conversion path exists - this must error.
	sv := model.NewSignatureValueWithValue(enumerations.SignatureAlgorithm_RSA_SHA256, []byte{1, 2, 3})
	got, err := c.EnsureSignatureValue(sv, enumerations.SignatureAlgorithm_RSA_SHA512)
	if err == nil {
		t.Fatal("expected error for a digest algorithm mismatch with no conversion path")
	}
	if got != nil {
		t.Fatalf("expected nil result on error, got %v", got)
	}
}

func TestSignatureValueCheckerEnsureSignatureValueUnsupportedConversionErrors(t *testing.T) {
	c := NewSignatureValueChecker()
	// Same digest algorithm (SHA256) but a different, non-EC-equivalent SignatureAlgorithm:
	// Java's DSSException("... Conversion is not supported!") path.
	sv := model.NewSignatureValueWithValue(enumerations.SignatureAlgorithm_RSA_SHA256, []byte{1, 2, 3})
	got, err := c.EnsureSignatureValue(sv, enumerations.SignatureAlgorithm_DSA_SHA256)
	if err == nil {
		t.Fatal("expected error: RSA -> DSA is not a supported conversion")
	}
	if got != nil {
		t.Fatalf("expected nil result on error, got %v", got)
	}
}

func TestSignatureValueCheckerEnsureSignatureValueEmptySignatureAlgorithmMismatchErrors(t *testing.T) {
	c := NewSignatureValueChecker()
	// signatureValue.Algorithm() empty (never set) means signatureDigestAlgorithm is the zero
	// value, which will not match a non-empty target digest algorithm.
	sv := model.NewSignatureValue()
	sv.SetValue([]byte{1, 2, 3})
	got, err := c.EnsureSignatureValue(sv, enumerations.SignatureAlgorithm_RSA_SHA256)
	if err == nil {
		t.Fatal("expected error when signatureValue carries no algorithm at all")
	}
	if got != nil {
		t.Fatalf("expected nil result on error, got %v", got)
	}
}
