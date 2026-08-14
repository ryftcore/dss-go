// Tests for AbstractSignatureService, matching
// dss-document/src/main/java/eu/europa/esig/dss/signature/AbstractSignatureService.java at
// /home/user/dss-upstream.
package document

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
)

func TestNewAbstractSignatureServiceNilCertificateVerifierPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil CertificateVerifier")
		}
	}()
	NewAbstractSignatureService[*AbstractSignatureParameters[*model.TimestampParameters], *model.TimestampParameters](nil)
}

func newTestAbstractSignatureService(t *testing.T) *AbstractSignatureService[*AbstractSignatureParameters[*model.TimestampParameters], *model.TimestampParameters] {
	t.Helper()
	s := NewAbstractSignatureService[*AbstractSignatureParameters[*model.TimestampParameters], *model.TimestampParameters](validation.NewCommonCertificateVerifier())
	return &s
}

func TestAbstractSignatureServiceSetTspSource(t *testing.T) {
	s := newTestAbstractSignatureService(t)
	if s.TspSource != nil {
		t.Fatal("TspSource should default to nil")
	}
	// SetTspSource just assigns the field; nil is a legitimate value to round-trip through it.
	s.SetTspSource(nil)
	if s.TspSource != nil {
		t.Fatal("SetTspSource(nil) should leave TspSource nil")
	}
}

func TestAbstractSignatureServiceAssertSigningCertificateValidNoCertificateNoFlagPanics(t *testing.T) {
	s := newTestAbstractSignatureService(t)
	params := NewAbstractSignatureParameters[*model.TimestampParameters]()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic: no signing certificate and GenerateTBSWithoutCertificate is false")
		}
	}()
	s.AssertSigningCertificateValid(&params)
}

func TestAbstractSignatureServiceAssertSigningCertificateValidNoCertificateFlagSetNoPanic(t *testing.T) {
	s := newTestAbstractSignatureService(t)
	params := NewAbstractSignatureParameters[*model.TimestampParameters]()
	params.SetGenerateTBSWithoutCertificate(true)
	// Should return without panicking: the "no certificate but explicitly allowed" escape hatch.
	s.AssertSigningCertificateValid(&params)
}

func TestAbstractSignatureServiceTimestampUnsupportedPanics(t *testing.T) {
	s := newTestAbstractSignatureService(t)
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic: Timestamp is unsupported unless overridden by a concrete format service")
		}
	}()
	s.Timestamp(model.NewInMemoryDocument([]byte("x")), nil)
}

func TestAbstractSignatureServiceIsValidSignatureValueNilArgumentsPanic(t *testing.T) {
	s := newTestAbstractSignatureService(t)
	_, cert := generateDocumentTestCertificate(t)
	tbs := model.NewToBeSigned()
	sv := model.NewSignatureValue()

	cases := []struct {
		name string
		fn   func()
	}{
		{"nil ToBeSigned", func() { s.IsValidSignatureValue(nil, sv, cert) }},
		{"nil SignatureValue", func() { s.IsValidSignatureValue(tbs, nil, cert) }},
		{"nil CertificateToken", func() { s.IsValidSignatureValue(tbs, sv, nil) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatalf("expected panic for %s", tc.name)
				}
			}()
			tc.fn()
		})
	}
}

func TestAbstractSignatureServiceIsValidSignatureValueRoundTrip(t *testing.T) {
	s := newTestAbstractSignatureService(t)
	key, cert := generateDocumentTestCertificate(t)

	toBeSigned := model.NewToBeSignedWithBytes([]byte("hello, DSS"))

	// Ports the ECDSA_SHA256 case AbstractSignatureService actually verifies against: hash and
	// sign with the certificate's own private key, mirroring what a real SignatureAlgorithm
	// implementation does upstream (crypto/ecdsa here since the test cert above is EC).
	digest := sha256.Sum256(toBeSigned.Bytes())
	sigBytes, err := ecdsa.SignASN1(rand.Reader, key, digest[:])
	if err != nil {
		t.Fatalf("sign: %s", err)
	}
	signatureValue := model.NewSignatureValueWithValue(enumerations.SignatureAlgorithm_ECDSA_SHA256, sigBytes)

	if !s.IsValidSignatureValue(toBeSigned, signatureValue, cert) {
		t.Fatal("IsValidSignatureValue() = false, want true for a genuine signature")
	}

	tampered := append([]byte(nil), sigBytes...)
	tampered[len(tampered)-1] ^= 0xFF
	tamperedValue := model.NewSignatureValueWithValue(enumerations.SignatureAlgorithm_ECDSA_SHA256, tampered)
	if s.IsValidSignatureValue(toBeSigned, tamperedValue, cert) {
		t.Fatal("IsValidSignatureValue() = true for a tampered signature, want false")
	}
}

func TestAbstractSignatureServiceGetFinalFileNameVariants(t *testing.T) {
	s := newTestAbstractSignatureService(t)
	doc := model.NewInMemoryDocumentWithName([]byte("x"), "original.txt")

	name1, err := s.GetFinalFileName(doc, enumerations.SigningOperation_SIGN)
	if err != nil {
		t.Fatalf("GetFinalFileName: %s", err)
	}
	if name1 == "" {
		t.Fatal("GetFinalFileName() returned empty name")
	}

	name2, err := s.GetFinalFileNameWithLevel(doc, enumerations.SigningOperation_SIGN, enumerations.SignatureLevel_XAdES_BASELINE_B)
	if err != nil {
		t.Fatalf("GetFinalFileNameWithLevel: %s", err)
	}
	if name2 == "" {
		t.Fatal("GetFinalFileNameWithLevel() returned empty name")
	}

	name3, err := s.GetFinalFileNameWithPackaging(doc, enumerations.SigningOperation_SIGN, enumerations.SignatureLevel_XAdES_BASELINE_B, enumerations.SignaturePackaging_ENVELOPING)
	if err != nil {
		t.Fatalf("GetFinalFileNameWithPackaging: %s", err)
	}
	if name3 == "" {
		t.Fatal("GetFinalFileNameWithPackaging() returned empty name")
	}

	name4, err := s.GetFinalDocumentNameWithMimeType(doc, enumerations.SigningOperation_SIGN, enumerations.SignatureLevel_XAdES_BASELINE_B, enumerations.MimeTypeEnum_BINARY)
	if err != nil {
		t.Fatalf("GetFinalDocumentNameWithMimeType: %s", err)
	}
	if name4 == "" {
		t.Fatal("GetFinalDocumentNameWithMimeType() returned empty name")
	}
}

func TestAbstractSignatureServiceEnsureSignatureValueMatchingAlgorithmPassesThrough(t *testing.T) {
	s := newTestAbstractSignatureService(t)
	sv := model.NewSignatureValueWithValue(enumerations.SignatureAlgorithm_ECDSA_SHA256, []byte{1, 2, 3})
	got, err := s.EnsureSignatureValue(enumerations.SignatureAlgorithm_ECDSA_SHA256, sv)
	if err != nil {
		t.Fatalf("EnsureSignatureValue: %s", err)
	}
	if got != sv {
		t.Fatalf("EnsureSignatureValue() = %v, want the same SignatureValue unchanged (algorithms already match)", got)
	}
}
