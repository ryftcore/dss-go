package token

import (
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

func TestAssertEncryptionAlgorithmValidMismatch(t *testing.T) {
	entry := &fakePrivateKeyEntry{certificate: mustLoadCertificateToken(t, "testdata/generated/eku.crt")}
	// fakePrivateKeyEntry reports RSA; ask for an ECDSA signature algorithm instead.
	err := abstractSignatureTokenConnectionAssertEncryptionAlgorithmValid(enumerations.SignatureAlgorithm_ECDSA_SHA256, entry)
	if err == nil {
		t.Fatal("expected an error for a mismatched EncryptionAlgorithm")
	}
}

func TestAssertEncryptionAlgorithmValidPanicsOnEmptySignatureAlgorithm(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic for an empty SignatureAlgorithm")
		}
	}()
	_ = abstractSignatureTokenConnectionAssertEncryptionAlgorithmValid("", &fakePrivateKeyEntry{})
}

func TestAssertEncryptionAlgorithmValidPanicsOnNilKeyEntry(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic for a nil keyEntry")
		}
	}()
	_ = abstractSignatureTokenConnectionAssertEncryptionAlgorithmValid(enumerations.SignatureAlgorithm_RSA_SHA256, nil)
}

func TestAssertDigestAlgorithmValidMismatch(t *testing.T) {
	digest := model.NewDigest(enumerations.DigestAlgorithm_SHA1, []byte{1, 2, 3})
	err := abstractSignatureTokenConnectionAssertDigestAlgorithmValid(digest, enumerations.SignatureAlgorithm_RSA_SHA256)
	if err == nil {
		t.Fatal("expected an error: SHA1 digest does not match the SHA256 signature algorithm")
	}
}

func TestAssertDigestAlgorithmValidRawAlgorithmAcceptsAnyDigest(t *testing.T) {
	digest := model.NewDigest(enumerations.DigestAlgorithm_SHA1, []byte{1, 2, 3})
	// RSA_RAW carries no digest algorithm, so any digest is accepted.
	if err := abstractSignatureTokenConnectionAssertDigestAlgorithmValid(digest, enumerations.SignatureAlgorithm_RSA_RAW); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSignDigestUnsupportedEncryptionAlgorithmCombination(t *testing.T) {
	entry := &fakePrivateKeyEntryWithEncryption{
		fakePrivateKeyEntry: fakePrivateKeyEntry{certificate: mustLoadCertificateToken(t, "testdata/generated/eku.crt")},
		encryptionAlgorithm: enumerations.EncryptionAlgorithm_EDDSA,
	}
	var base AbstractSignatureTokenConnection
	// EDDSA has no (EDDSA, "") entry in the SignatureAlgorithm table, matching upstream's
	// getRawSignatureAlgorithm(EDDSA) UnsupportedOperationException.
	_, err := base.SignDigest(model.NewDigest(enumerations.DigestAlgorithm_SHA256, []byte{1, 2, 3}), entry)
	if err == nil {
		t.Fatal("expected an error: EdDSA has no raw/digest-signing SignatureAlgorithm")
	}
}

// fakePrivateKeyEntryWithEncryption lets tests override the reported EncryptionAlgorithm.
type fakePrivateKeyEntryWithEncryption struct {
	fakePrivateKeyEntry
	encryptionAlgorithm enumerations.EncryptionAlgorithm
}

func (f *fakePrivateKeyEntryWithEncryption) EncryptionAlgorithm() enumerations.EncryptionAlgorithm {
	return f.encryptionAlgorithm
}
