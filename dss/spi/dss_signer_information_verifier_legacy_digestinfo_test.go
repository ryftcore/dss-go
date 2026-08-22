// Regression tests for SignerInformationVerifier#Verify's RSA PKCS#1 v1.5 fallback: a DigestInfo
// whose digest AlgorithmIdentifier omits the NULL parameter must verify (BouncyCastle's
// RSADigestSigner accepts it, so upstream DSS reports such signatures intact), while every other
// deviation - wrong content, tampered signature, trailing bytes inside the DigestInfo, broken
// padding - must still be rejected. The point of the negative cases is that the fallback widened
// exactly one thing, the DigestInfo spelling, and nothing else about PKCS#1 v1.5.
package spi

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// legacyDigestInfoContent is the payload every case below signs or claims to have signed.
var legacyDigestInfoContent = []byte("CMS SignerInfo signed attributes stand-in")

// legacyDigestInfoKey is one RSA key shared by the whole file; generating a 2048-bit key per
// case would dominate the runtime for no added coverage.
var legacyDigestInfoKey = func() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return key
}()

// legacyDigestInfoWithoutNullParameter builds SEQUENCE { SEQUENCE { OID }, OCTET STRING digest },
// i.e. the non-canonical DigestInfo real-world signers produce, for the given content.
func legacyDigestInfoWithoutNullParameter(t *testing.T, content []byte) []byte {
	t.Helper()
	digest, err := DSSUtilsDigest(enumerations.DigestAlgorithmSHA256, content)
	if err != nil {
		t.Fatalf("digesting: %s", err)
	}
	objectIdentifier, err := asn1ber.OIDFromString(enumerations.DigestAlgorithmSHA256.OID())
	if err != nil {
		t.Fatalf("parsing the SHA-256 OID: %s", err)
	}
	return asn1ber.WriteSequence(append(
		asn1ber.WriteSequence(asn1ber.EncodeOID(objectIdentifier)),
		asn1ber.WriteTLV(asn1ber.TagOctetString, digest)...))
}

// legacyDigestInfoVerifier is the verifier under test, bound to the shared key's public half.
func legacyDigestInfoVerifier() *SignerInformationVerifier {
	return &SignerInformationVerifier{publicKey: &legacyDigestInfoKey.PublicKey}
}

// TestVerifyAcceptsDigestInfoWithoutNullParameter is the fixed defect: before the fallback, this
// signature - mathematically valid, and accepted by upstream DSS through BouncyCastle - was
// reported as broken.
func TestVerifyAcceptsDigestInfoWithoutNullParameter(t *testing.T) {
	digestInfo := legacyDigestInfoWithoutNullParameter(t, legacyDigestInfoContent)
	// crypto.Hash(0) signs the supplied DigestInfo verbatim, producing exactly the encoding the
	// real-world fixtures in cades/testdata/upstream carry.
	signature, err := rsa.SignPKCS1v15(rand.Reader, legacyDigestInfoKey, crypto.Hash(0), digestInfo)
	if err != nil {
		t.Fatalf("signing: %s", err)
	}
	if err := legacyDigestInfoVerifier().Verify(
		enumerations.SignatureAlgorithmRSASHA256, legacyDigestInfoContent, signature); err != nil {
		t.Errorf("a DigestInfo without the NULL parameter must verify: %s", err)
	}
}

// TestVerifyStillAcceptsCanonicalDigestInfo checks the fallback did not disturb the normal path.
func TestVerifyStillAcceptsCanonicalDigestInfo(t *testing.T) {
	digest, err := DSSUtilsDigest(enumerations.DigestAlgorithmSHA256, legacyDigestInfoContent)
	if err != nil {
		t.Fatalf("digesting: %s", err)
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, legacyDigestInfoKey, crypto.SHA256, digest)
	if err != nil {
		t.Fatalf("signing: %s", err)
	}
	if err := legacyDigestInfoVerifier().Verify(
		enumerations.SignatureAlgorithmRSASHA256, legacyDigestInfoContent, signature); err != nil {
		t.Errorf("a canonical signature must still verify: %s", err)
	}
}

// TestVerifyRejectsEverythingElse pins that the fallback widened the DigestInfo spelling and
// nothing else: content substitution, signature tampering, and - the case that would matter for
// a Bleichenbacher-style forgery - a DigestInfo padded out with trailing bytes are all still
// rejected, because crypto/rsa rebuilds and compares the whole EM block.
func TestVerifyRejectsEverythingElse(t *testing.T) {
	digestInfo := legacyDigestInfoWithoutNullParameter(t, legacyDigestInfoContent)
	signature, err := rsa.SignPKCS1v15(rand.Reader, legacyDigestInfoKey, crypto.Hash(0), digestInfo)
	if err != nil {
		t.Fatalf("signing: %s", err)
	}

	t.Run("different content", func(t *testing.T) {
		if err := legacyDigestInfoVerifier().Verify(
			enumerations.SignatureAlgorithmRSASHA256, []byte("something else entirely"), signature); err == nil {
			t.Error("a signature over different content must be rejected")
		}
	})

	t.Run("tampered signature", func(t *testing.T) {
		tampered := make([]byte, len(signature))
		copy(tampered, signature)
		tampered[len(tampered)-1] ^= 0x01
		if err := legacyDigestInfoVerifier().Verify(
			enumerations.SignatureAlgorithmRSASHA256, legacyDigestInfoContent, tampered); err == nil {
			t.Error("a tampered signature must be rejected")
		}
	})

	t.Run("trailing bytes after the DigestInfo", func(t *testing.T) {
		// The classic forgery shape: a well-formed DigestInfo followed by attacker-chosen bytes
		// inside the same padded block. Signed here with the real key only to prove the verifier
		// rejects the *shape*, not merely an unsigned blob.
		padded := append(append([]byte{}, digestInfo...), 0xde, 0xad, 0xbe, 0xef)
		forged, err := rsa.SignPKCS1v15(rand.Reader, legacyDigestInfoKey, crypto.Hash(0), padded)
		if err != nil {
			t.Fatalf("signing: %s", err)
		}
		if err := legacyDigestInfoVerifier().Verify(
			enumerations.SignatureAlgorithmRSASHA256, legacyDigestInfoContent, forged); err == nil {
			t.Error("a DigestInfo with trailing bytes must be rejected")
		}
	})

	t.Run("wrong digest algorithm", func(t *testing.T) {
		if err := legacyDigestInfoVerifier().Verify(
			enumerations.SignatureAlgorithmRSASHA512, legacyDigestInfoContent, signature); err == nil {
			t.Error("verifying under a different digest algorithm must be rejected")
		}
	})
}

// TestVerifyFallbackScope checks the fallback helper itself refuses to run outside RSA PKCS#1
// v1.5 - notably for RSASSA-PSS, whose encoding carries no DigestInfo at all.
func TestVerifyFallbackScope(t *testing.T) {
	for _, signatureAlgorithm := range []enumerations.SignatureAlgorithm{
		enumerations.SignatureAlgorithmRSASSAPSSSHA256MGF1,
		enumerations.SignatureAlgorithmECDSASHA256,
		enumerations.SignatureAlgorithmED25519,
	} {
		if err := dssSignerInformationVerifierVerifyRSAWithoutDigestInfoNullParameter(
			&legacyDigestInfoKey.PublicKey, signatureAlgorithm, legacyDigestInfoContent, []byte{0x00}); err == nil {
			t.Errorf("%s: the fallback must not apply", signatureAlgorithm)
		}
	}
}
