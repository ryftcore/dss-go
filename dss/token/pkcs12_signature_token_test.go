// Tests use the same fixtures and passwords as
// dss-token/src/test/java/eu/europa/esig/dss/token/Pkcs12SignatureTokenTest.java (copied into
// testdata/dss-token/src/test/resources per PORTING.md), plus two openssl-generated PKCS12
// files under testdata/generated for the sign/verify round trip. Those two were produced with:
//
//	openssl req -x509 -newkey rsa:2048 -keyout rsa.key -out rsa.crt -days 3650 -nodes \
//	    -subj "/CN=Go Port Test RSA/O=DSS Go Port"
//	openssl pkcs12 -export -legacy -in rsa.crt -inkey rsa.key -out rsa_test.p12 -name "rsa-test" \
//	    -passout pass:testpassword -keypbe PBE-SHA1-3DES -certpbe PBE-SHA1-3DES -macalg SHA1
//	openssl ecparam -name prime256v1 -genkey -noout -out ec.key
//	openssl req -x509 -new -key ec.key -out ec.crt -days 3650 -nodes \
//	    -subj "/CN=Go Port Test EC/O=DSS Go Port"
//	openssl pkcs12 -export -legacy -in ec.crt -inkey ec.key -out ec_test.p12 -name "ec-test" \
//	    -passout pass:testpassword -keypbe PBE-SHA1-3DES -certpbe PBE-SHA1-3DES -macalg SHA1
//
// -legacy/PBE-SHA1-3DES/-macalg SHA1 are deliberate: golang.org/x/crypto/pkcs12 (pinned in
// go.mod) only implements the classic PBE-SHA1-3DES/RC2 encryption schemes and a SHA-1 MAC, not
// the AES/SHA-256 defaults modern OpenSSL produces - see key_store_signature_token_connection.go's
// file header for the full list of PKCS12 gaps this pinned dependency version has.
package token

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"os"
	"strings"
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

func mustReadFixture(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading fixture %s: %v", path, err)
	}
	return data
}

// TestPkcs12SignatureTokenLoadAndSign ports Pkcs12SignatureTokenTest#testPkcs12.
func TestPkcs12SignatureTokenLoadAndSign(t *testing.T) {
	data := mustReadFixture(t, "testdata/dss-token/src/test/resources/user_a_rsa.p12")
	signatureToken, err := NewPkcs12SignatureTokenFromBytes(data, NewPasswordProtection([]byte("password")))
	if err != nil {
		t.Fatalf("NewPkcs12SignatureTokenFromBytes: %v", err)
	}
	defer signatureToken.Close()

	keys, err := signatureToken.Keys()
	if err != nil {
		t.Fatalf("Keys: %v", err)
	}
	if len(keys) == 0 {
		t.Fatal("expected at least one key")
	}

	dssPrivateKeyEntry, ok := keys[0].(*KSPrivateKeyEntry)
	if !ok {
		t.Fatalf("keys[0] is %T, want *KSPrivateKeyEntry", keys[0])
	}
	if dssPrivateKeyEntry.Alias() == "" {
		t.Fatal("expected a non-empty alias")
	}

	entry, err := signatureToken.KeyWithPassword(dssPrivateKeyEntry.Alias(), NewPasswordProtection([]byte("password")))
	if err != nil {
		t.Fatalf("KeyWithPassword: %v", err)
	}
	if entry == nil {
		t.Fatal("expected a non-nil entry")
	}
	if entry.Certificate() == nil {
		t.Fatal("expected a non-nil certificate")
	}
	if len(entry.CertificateChain()) == 0 {
		t.Fatal("expected a non-empty certificate chain")
	}
	if entry.EncryptionAlgorithm() == "" {
		t.Fatal("expected a non-empty EncryptionAlgorithm")
	}

	toBeSigned := model.NewToBeSignedWithBytes([]byte("Hello world"))
	signValue, err := signatureToken.Sign(toBeSigned, enumerations.DigestAlgorithm_SHA256, entry)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if signValue.Algorithm() == "" {
		t.Fatal("expected a non-empty SignatureAlgorithm")
	}
	if len(signValue.Value()) == 0 {
		t.Fatal("expected a non-empty signature value")
	}
}

// TestPkcs12SignatureTokenWrongPassword ports Pkcs12SignatureTokenTest#wrongPassword: the error
// message must match Java's DSSException message exactly (message-based error identity is part
// of the port's contract - see PORTING.md).
func TestPkcs12SignatureTokenWrongPassword(t *testing.T) {
	data := mustReadFixture(t, "testdata/dss-token/src/test/resources/user_a_rsa.p12")
	_, err := NewPkcs12SignatureTokenFromBytes(data, NewPasswordProtection([]byte("wrong password")))
	if err == nil {
		t.Fatal("expected an error")
	}
	dssErr, ok := err.(*model.DSSError)
	if !ok {
		t.Fatalf("error is %T, want *model.DSSError", err)
	}
	if dssErr.Message != "Unable to instantiate KeyStoreSignatureTokenConnection" {
		t.Errorf("Message = %q, want %q", dssErr.Message, "Unable to instantiate KeyStoreSignatureTokenConnection")
	}
}

// TestPkcs12SignatureTokenChainPreserved loads a PKCS12 store whose certificate bag holds a full
// (leaf + intermediate + root) chain and checks it is reconstructed, matching what
// java.security.KeyStore's PKCS12 provider exposes through getCertificateChain().
func TestPkcs12SignatureTokenChainPreserved(t *testing.T) {
	data := mustReadFixture(t, "testdata/dss-token/src/test/resources/good-ecdsa-user.p12")
	signatureToken, err := NewPkcs12SignatureTokenFromBytes(data, NewPasswordProtection([]byte("ks-password")))
	if err != nil {
		t.Fatalf("NewPkcs12SignatureTokenFromBytes: %v", err)
	}
	defer signatureToken.Close()

	keys, err := signatureToken.Keys()
	if err != nil {
		t.Fatalf("Keys: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("len(keys) = %d, want 1", len(keys))
	}
	chain := keys[0].CertificateChain()
	if len(chain) < 2 {
		t.Fatalf("len(chain) = %d, want a multi-certificate chain (leaf + at least one issuer)", len(chain))
	}
	// Each certificate but the last must be issued by the next one in the chain.
	for i := 0; i < len(chain)-1; i++ {
		if chain[i].Issuer().Canonical() != chain[i+1].Subject().Canonical() {
			t.Errorf("chain[%d].Issuer() = %q, chain[%d].Subject() = %q; chain is not linked",
				i, chain[i].Issuer().Canonical(), i+1, chain[i+1].Subject().Canonical())
		}
	}
	if keys[0].EncryptionAlgorithm() != enumerations.EncryptionAlgorithm_ECDSA {
		t.Errorf("EncryptionAlgorithm() = %s, want ECDSA", keys[0].EncryptionAlgorithm())
	}
}

// signVerifyRoundTrip signs a message and a raw digest with every entry keys[0] can use under
// digestAlgorithm, and verifies each signature independently with the Go standard library -
// never through this port's own code - so a self-consistent bug in the signing path cannot also
// make its matching verification pass.
func signVerifyRoundTrip(t *testing.T, tok SignatureTokenConnection, entry DSSPrivateKeyEntry, digestAlgorithm enumerations.DigestAlgorithm) {
	t.Helper()

	message := []byte("The quick brown fox jumps over the lazy dog")
	toBeSigned := model.NewToBeSignedWithBytes(message)

	signValue, err := tok.Sign(toBeSigned, digestAlgorithm, entry)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	verifySignatureValue(t, entry.Certificate().Certificate().PublicKey, message, signValue.Value())

	sum := sha256.Sum256(message)
	digest := model.NewDigest(digestAlgorithm, sum[:])
	digestSignValue, err := tok.SignDigest(digest, entry)
	if err != nil {
		t.Fatalf("SignDigest: %v", err)
	}
	verifySignatureValue(t, entry.Certificate().Certificate().PublicKey, message, digestSignValue.Value())
}

func verifySignatureValue(t *testing.T, publicKey any, message, signature []byte) {
	t.Helper()
	switch pub := publicKey.(type) {
	case *rsa.PublicKey:
		sum := sha256.Sum256(message)
		if err := rsa.VerifyPKCS1v15(pub, crypto.Hash(0), appendDigestInfoPrefixSHA256(sum[:]), signature); err != nil {
			t.Errorf("rsa.VerifyPKCS1v15: %v", err)
		}
	case *ecdsa.PublicKey:
		sum := sha256.Sum256(message)
		if !ecdsa.VerifyASN1(pub, sum[:], signature) {
			t.Error("ecdsa.VerifyASN1: signature does not verify")
		}
	default:
		t.Fatalf("unsupported public key type %T", publicKey)
	}
}

// appendDigestInfoPrefixSHA256 lets rsa.VerifyPKCS1v15 be called with crypto.Hash(0) (raw),
// mirroring how this port signs plain RSA: it always feeds NONEwithRSA an already
// DigestInfo-encoded digest (see AbstractSignatureTokenConnection), so verification must build
// the same DigestInfo before the padding check.
func appendDigestInfoPrefixSHA256(digest []byte) []byte {
	encoded, err := DigestInfoEncoderEncode(enumerations.DigestAlgorithm_SHA256.OID(), digest)
	if err != nil {
		panic(err)
	}
	return encoded
}

func TestPkcs12SignatureTokenRoundTripRSA(t *testing.T) {
	data := mustReadFixture(t, "testdata/generated/rsa_test.p12")
	signatureToken, err := NewPkcs12SignatureTokenFromBytes(data, NewPasswordProtection([]byte("testpassword")))
	if err != nil {
		t.Fatalf("NewPkcs12SignatureTokenFromBytes: %v", err)
	}
	defer signatureToken.Close()

	keys, err := signatureToken.Keys()
	if err != nil {
		t.Fatalf("Keys: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("len(keys) = %d, want 1", len(keys))
	}
	if _, ok := entryPublicKey(keys[0]).(*rsa.PublicKey); !ok {
		t.Fatalf("expected an RSA key, got %T", entryPublicKey(keys[0]))
	}

	signVerifyRoundTrip(t, signatureToken, keys[0], enumerations.DigestAlgorithm_SHA256)
}

func TestPkcs12SignatureTokenRoundTripECDSA(t *testing.T) {
	data := mustReadFixture(t, "testdata/generated/ec_test.p12")
	signatureToken, err := NewPkcs12SignatureTokenFromBytes(data, NewPasswordProtection([]byte("testpassword")))
	if err != nil {
		t.Fatalf("NewPkcs12SignatureTokenFromBytes: %v", err)
	}
	defer signatureToken.Close()

	keys, err := signatureToken.Keys()
	if err != nil {
		t.Fatalf("Keys: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("len(keys) = %d, want 1", len(keys))
	}
	if _, ok := entryPublicKey(keys[0]).(*ecdsa.PublicKey); !ok {
		t.Fatalf("expected an ECDSA key, got %T", entryPublicKey(keys[0]))
	}

	signVerifyRoundTrip(t, signatureToken, keys[0], enumerations.DigestAlgorithm_SHA256)
}

// TestPkcs12SignatureTokenRoundTripRSAPSS exercises the RSASSA-PSS branch of
// AbstractSignatureTokenConnection, verifying with rsa.VerifyPSS from the Go standard library.
func TestPkcs12SignatureTokenRoundTripRSAPSS(t *testing.T) {
	data := mustReadFixture(t, "testdata/generated/rsa_test.p12")
	signatureToken, err := NewPkcs12SignatureTokenFromBytes(data, NewPasswordProtection([]byte("testpassword")))
	if err != nil {
		t.Fatalf("NewPkcs12SignatureTokenFromBytes: %v", err)
	}
	defer signatureToken.Close()

	keys, err := signatureToken.Keys()
	if err != nil || len(keys) != 1 {
		t.Fatalf("Keys: %v (len=%d)", err, len(keys))
	}
	entry := keys[0]
	rsaPub, ok := entryPublicKey(entry).(*rsa.PublicKey)
	if !ok {
		t.Fatalf("expected an RSA key, got %T", entryPublicKey(entry))
	}

	message := []byte("PSS round trip message")
	toBeSigned := model.NewToBeSignedWithBytes(message)
	signValue, err := signatureToken.SignWithSignatureAlgorithm(toBeSigned, enumerations.SignatureAlgorithm_RSA_SSA_PSS_SHA256_MGF1, entry)
	if err != nil {
		t.Fatalf("SignWithSignatureAlgorithm: %v", err)
	}

	sum := sha256.Sum256(message)
	opts := &rsa.PSSOptions{SaltLength: enumerations.DigestAlgorithm_SHA256.SaltLength(), Hash: crypto.SHA256}
	if err := rsa.VerifyPSS(rsaPub, crypto.SHA256, sum[:], signValue.Value(), opts); err != nil {
		t.Errorf("rsa.VerifyPSS: %v", err)
	}
}

func entryPublicKey(entry DSSPrivateKeyEntry) any {
	return entry.Certificate().Certificate().PublicKey
}

// TestPkcs12SignatureTokenChainPreservedX509 double-checks the reconstructed chain against a
// direct crypto/x509 parse of the certificate bag order, independent of model.CertificateToken.
func TestPkcs12SignatureTokenChainPreservedX509(t *testing.T) {
	data := mustReadFixture(t, "testdata/dss-token/src/test/resources/good-ecdsa-user.p12")
	signatureToken, err := NewPkcs12SignatureTokenFromBytes(data, NewPasswordProtection([]byte("ks-password")))
	if err != nil {
		t.Fatalf("NewPkcs12SignatureTokenFromBytes: %v", err)
	}
	defer signatureToken.Close()

	keys, err := signatureToken.Keys()
	if err != nil || len(keys) != 1 {
		t.Fatalf("Keys: %v (len=%d)", err, len(keys))
	}
	for _, cert := range keys[0].CertificateChain() {
		if _, err := x509.ParseCertificate(cert.Certificate().Raw); err != nil {
			t.Errorf("chain certificate does not re-parse: %v", err)
		}
	}
}

// TestKeyStoreUnsupportedKeyTypeFixtures pins the two upstream dss-token fixtures this port
// cannot open, so the capability gap documented in key_store_signature_token_connection.go's
// header stays visible and any dependency change that closes it fails this test loudly.
//
// Both stores open under Java's JCA KeyStore in upstream Pkcs12SignatureToken, yielding one key
// entry with a 3-certificate chain; here both ToPEM (key type it cannot re-encode) and the
// Decode fallback (which demands exactly two safe bags) refuse them.
func TestKeyStoreUnsupportedKeyTypeFixtures(t *testing.T) {
	for _, testCase := range []struct {
		fixture       string
		password      string
		toPEMReason   string
		fallbackCause string
	}{
		{
			fixture:       "Ed25519-good-user.p12",
			password:      "ks-password",
			toPEMReason:   "found unknown private key type in PKCS#8 wrapping",
			fallbackCause: "expected exactly two safe bags in the PFX PDU",
		},
		{
			fixture:       "good-dsa-user.p12",
			password:      "ks-password",
			toPEMReason:   "unknown algorithm: 1.2.840.10040.4.1",
			fallbackCause: "expected exactly two safe bags in the PFX PDU",
		},
	} {
		t.Run(testCase.fixture, func(t *testing.T) {
			data := mustReadFixture(t, "testdata/dss-token/src/test/resources/"+testCase.fixture)
			_, err := NewPkcs12SignatureTokenFromBytes(data, NewPasswordProtection([]byte(testCase.password)))
			if err == nil {
				t.Fatalf("%s now loads: golang.org/x/crypto/pkcs12 gained support for this key type, "+
					"so the LIMITATION note in key_store_signature_token_connection.go and this test "+
					"must be revisited (upstream yields 1 key entry with a 3-certificate chain)", testCase.fixture)
			}
			if !strings.Contains(err.Error(), testCase.toPEMReason) {
				t.Errorf("chain-preserving parse failed for an unexpected reason:\n got %q\n want it to mention %q",
					err, testCase.toPEMReason)
			}
			if !strings.Contains(err.Error(), testCase.fallbackCause) {
				t.Errorf("single-entry fallback failed for an unexpected reason:\n got %q\n want it to mention %q",
					err, testCase.fallbackCause)
			}
		})
	}
}
