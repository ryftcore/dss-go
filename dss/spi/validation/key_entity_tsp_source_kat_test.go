package validation

import (
	"crypto"
	"crypto/x509"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/cmscore"
	"github.com/utain/esig/dss/model"
)

// The tokens this file issues are compared with the ones BouncyCastle 1.84 issues from the same
// key, certificate, policy, production time and serial number - testdata/gen/TspFixtures.java
// drives KeyEntityTSPSource's own wiring - so the comparison covers the TSTInfo encoding, the
// signed-attribute set and the RSA PKCS#1 v1.5 signature, all of which are deterministic.

// keyEntityTSPSourceTestPinnedSerial makes the serial number of the issued token deterministic,
// the way a subclass overriding getTimeStampSerialNumber() would.
type keyEntityTSPSourceTestPinnedSerial struct {
	KeyEntityTSPSource
	serial *big.Int
}

// TimeStampSerialNumber returns the pinned serial number.
func (s *keyEntityTSPSourceTestPinnedSerial) TimeStampSerialNumber() (*big.Int, error) {
	return s.serial, nil
}

// keyEntityTSPSourceTestSource builds a source over the testdata TSA key and certificate, with
// the production time and the serial number pinned to the ones the fixtures were issued with.
func keyEntityTSPSourceTestSource(t *testing.T, chain []*x509.Certificate,
	digestAlgorithm enumerations.DigestAlgorithm) *keyEntityTSPSourceTestPinnedSerial {
	t.Helper()
	key, err := x509.ParsePKCS8PrivateKey(timestampTokenKATFile(t, "tsa.key"))
	if err != nil {
		t.Fatalf("unable to read tsa.key: %v", err)
	}
	signer, ok := key.(crypto.Signer)
	if !ok {
		t.Fatal("tsa.key is not a signing key")
	}
	serial, ok := new(big.Int).SetString("123456789012345678901234567890", 10)
	if !ok {
		t.Fatal("unable to build the pinned serial number")
	}
	source := &keyEntityTSPSourceTestPinnedSerial{
		KeyEntityTSPSource: NewKeyEntityTSPSourceBase(),
		serial:             serial,
	}
	source.InitKeyEntityTSPSource(source)
	source.SetPrivateKey(signer)
	source.SetCertificate(chain[0])
	source.SetCertificateChain(chain)
	source.SetTsaPolicy("1.2.3.4.1")
	source.SetDigestAlgorithm(digestAlgorithm)
	// 2021-01-02T03:04:05Z, the production time TspFixtures.java pins.
	source.SetProductionTime(time.Unix(0, 1609556645000*int64(time.Millisecond)).UTC())
	return source
}

// keyEntityTSPSourceTestCertificates loads the DER certificates of testdata.
func keyEntityTSPSourceTestCertificates(t *testing.T, names ...string) []*x509.Certificate {
	t.Helper()
	certificates := make([]*x509.Certificate, 0, len(names))
	for _, name := range names {
		certificate, err := x509.ParseCertificate(timestampTokenKATFile(t, name))
		if err != nil {
			t.Fatalf("unable to parse %s: %v", name, err)
		}
		certificates = append(certificates, certificate)
	}
	return certificates
}

// TestKeyEntityTSPSourceKAT_MatchesBouncyCastle issues the single-certificate fixture again and
// requires the bytes to be identical to the ones BouncyCastle produced.
func TestKeyEntityTSPSourceKAT_MatchesBouncyCastle(t *testing.T) {
	source := keyEntityTSPSourceTestSource(t, keyEntityTSPSourceTestCertificates(t, "tsa.crt"),
		enumerations.DigestAlgorithm_SHA256)

	binary, err := source.TimeStampResponse(enumerations.DigestAlgorithm_SHA512,
		sha512Sum(timestampTokenKATFile(t, "content.bin")))
	if err != nil {
		t.Fatalf("TimeStampResponse() failed: %v", err)
	}
	expected := timestampTokenKATFile(t, "timestamp-token-sha512.tst")
	if hex.EncodeToString(binary.Bytes()) != hex.EncodeToString(expected) {
		t.Errorf("the issued token differs from the BouncyCastle one:\n got %s\nwant %s",
			hex.EncodeToString(binary.Bytes()), hex.EncodeToString(expected))
	}
}

// TestKeyEntityTSPSourceKAT_TwoCertificateChain issues the two-certificate fixture. Its bytes
// cannot be identical - cmscore DER-orders the certificates SET where BouncyCastle keeps the
// insertion order - so the TSTInfo, the signed attributes and the signature are compared instead,
// which is everything the ordering does not touch.
func TestKeyEntityTSPSourceKAT_TwoCertificateChain(t *testing.T) {
	oracle := timestampTokenKATOracle(t)
	source := keyEntityTSPSourceTestSource(t, keyEntityTSPSourceTestCertificates(t, "tsa.crt", "ca.crt"),
		enumerations.DigestAlgorithm_SHA512)

	binary, err := source.TimeStampResponse(enumerations.DigestAlgorithm_SHA256,
		sha256Sum(timestampTokenKATFile(t, "content.bin")))
	if err != nil {
		t.Fatalf("TimeStampResponse() failed: %v", err)
	}

	cms, err := cmscore.ParseCMS(binary.Bytes())
	if err != nil {
		t.Fatalf("the issued token is not a CMS: %v", err)
	}
	if got, want := hex.EncodeToString(sha256Sum(cms.SignedContent())),
		oracle["timestamp-token.tst.tstInfoDER.sha256"]; got != want {
		t.Errorf("TSTInfo digest = %s, want %s", got, want)
	}
	signerInfo := cms.SignerInfos()[0]
	if got, want := hex.EncodeToString(sha256Sum(signerInfo.SignedAttributesDER())),
		oracle["timestamp-token.tst.signedAttributesDER"]; got != want {
		t.Errorf("signed attributes digest = %s, want %s", got, want)
	}
	if got, want := signerInfo.DigestAlgorithm.Algorithm.String(),
		oracle["timestamp-token.tst.signerDigestAlgOID"]; got != want {
		t.Errorf("SignerInfo.digestAlgorithm = %s, want %s", got, want)
	}
	if got, want := signerInfo.SignatureAlgorithm.Algorithm.String(),
		oracle["timestamp-token.tst.signerEncryptionAlgOID"]; got != want {
		t.Errorf("SignerInfo.signatureAlgorithm = %s, want %s", got, want)
	}
	// The whole SignerInfo is byte-identical, signature included: RSA PKCS#1 v1.5 is
	// deterministic and the signed attributes are the same.
	reference, err := cmscore.ParseCMS(timestampTokenKATFile(t, "timestamp-token.tst"))
	if err != nil {
		t.Fatalf("the BouncyCastle fixture is not a CMS: %v", err)
	}
	if hex.EncodeToString(signerInfo.Encoded()) != hex.EncodeToString(reference.SignerInfos()[0].Encoded()) {
		t.Error("the issued SignerInfo differs from the BouncyCastle one")
	}
}

// TestKeyEntityTSPSourceIssuedTokenValidates feeds an issued token back into TimestampToken: it
// has to parse, to match the data it was issued over and to verify against the TSA certificate.
func TestKeyEntityTSPSourceIssuedTokenValidates(t *testing.T) {
	certificates := keyEntityTSPSourceTestCertificates(t, "tsa.crt", "ca.crt")
	source := keyEntityTSPSourceTestSource(t, certificates, enumerations.DigestAlgorithm_SHA256)
	digest := sha256Sum(timestampTokenKATFile(t, "content.bin"))

	binary, err := source.TimeStampResponse(enumerations.DigestAlgorithm_SHA256, digest)
	if err != nil {
		t.Fatalf("TimeStampResponse() failed: %v", err)
	}
	token, err := NewTimestampToken(binary.Bytes(), enumerations.TimestampType_SIGNATURE_TIMESTAMP)
	if err != nil {
		t.Fatalf("the issued token is not a TimestampToken: %v", err)
	}
	if !token.MatchData(digest) {
		t.Error("MatchData() = false for the digest the token was issued over, want true")
	}
	tsa, err := model.NewCertificateToken(certificates[0])
	if err != nil {
		t.Fatalf("unable to build the TSA CertificateToken: %v", err)
	}
	if !token.IsSignedByToken(tsa) {
		t.Fatalf("IsSignedByToken(tsa) = false, want true (reason: %s)", token.InvalidityReason())
	}
	if !token.IsValid() {
		t.Error("IsValid() = false, want true")
	}
	if got := token.GenerationTime().UnixMilli(); got != 1609556645000 {
		t.Errorf("GenerationTime() = %d, want the pinned production time", got)
	}
	if got := token.SignatureAlgorithm(); got != enumerations.SignatureAlgorithm_RSA_SHA256 {
		t.Errorf("SignatureAlgorithm() = %s, want RSA_SHA256", got)
	}
}

// TestKeyEntityTSPSourceRejectsUnacceptedAlgorithm covers the accepted-digest-algorithm guard and
// the default set upstream ships with.
func TestKeyEntityTSPSourceRejectsUnacceptedAlgorithm(t *testing.T) {
	source := keyEntityTSPSourceTestSource(t, keyEntityTSPSourceTestCertificates(t, "tsa.crt"),
		enumerations.DigestAlgorithm_SHA256)

	if _, err := source.TimeStampResponse(enumerations.DigestAlgorithm_SHA1, make([]byte, 20)); err == nil {
		t.Error("TimeStampResponse(SHA1) succeeded, want the unsupported-algorithm error")
	} else if got, want := err.Error(),
		"DigestAlgorithm 'SHA1' is not supported by the KeyEntityTSPSource implementation!"; got != want {
		t.Errorf("error = %q, want %q", got, want)
	}

	source.SetAcceptedDigestAlgorithms([]enumerations.DigestAlgorithm{enumerations.DigestAlgorithm_SHA1})
	if _, err := source.TimeStampResponse(enumerations.DigestAlgorithm_SHA256, make([]byte, 32)); err == nil {
		t.Error("TimeStampResponse(SHA256) succeeded after narrowing the accepted set, want an error")
	}
}

// TestKeyEntityTSPSourceMissingProperties reproduces the Objects.requireNonNull guards, which the
// port turns into panics carrying the Java messages.
func TestKeyEntityTSPSourceMissingProperties(t *testing.T) {
	certificates := keyEntityTSPSourceTestCertificates(t, "tsa.crt")

	cases := map[string]struct {
		configure func(*keyEntityTSPSourceTestPinnedSerial)
		message   string
	}{
		"no policy": {
			configure: func(s *keyEntityTSPSourceTestPinnedSerial) { s.SetTsaPolicy("") },
			message:   "TSAPolicy OID is not defined! Use #setTsaPolicy method.",
		},
		"no key": {
			configure: func(s *keyEntityTSPSourceTestPinnedSerial) { s.SetPrivateKey(nil) },
			message:   "PrivateKey is not defined! Use #setPrivateKey method.",
		},
		"no certificate": {
			configure: func(s *keyEntityTSPSourceTestPinnedSerial) { s.SetCertificate(nil) },
			message:   "Certificate is not defined! Use #setCertificate method.",
		},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			source := keyEntityTSPSourceTestSource(t, certificates, enumerations.DigestAlgorithm_SHA256)
			testCase.configure(source)
			defer func() {
				if recovered := recover(); recovered != testCase.message {
					t.Errorf("recover() = %v, want %q", recovered, testCase.message)
				}
			}()
			_, _ = source.TimeStampResponse(enumerations.DigestAlgorithm_SHA256, make([]byte, 32))
			t.Error("TimeStampResponse() returned, want a panic")
		})
	}
}

// TestKeyEntityTSPSourceSignatureAlgorithm covers the encryption-algorithm compatibility check of
// getSignatureAlgorithm().
func TestKeyEntityTSPSourceSignatureAlgorithm(t *testing.T) {
	source := keyEntityTSPSourceTestSource(t, keyEntityTSPSourceTestCertificates(t, "tsa.crt"),
		enumerations.DigestAlgorithm_SHA384)

	algorithm, err := source.SignatureAlgorithm()
	if err != nil {
		t.Fatalf("SignatureAlgorithm() failed: %v", err)
	}
	if algorithm != enumerations.SignatureAlgorithm_RSA_SHA384 {
		t.Errorf("SignatureAlgorithm() = %s, want RSA_SHA384", algorithm)
	}

	source.SetEncryptionAlgorithm(enumerations.EncryptionAlgorithm_ECDSA)
	if _, err := source.SignatureAlgorithm(); err == nil {
		t.Error("SignatureAlgorithm() accepted ECDSA for an RSA key, want an error")
	}
}

// TestKeyEntityTSPSourceKeyStoreRefusesUnknownType checks the documented refusal of every key
// store format Go cannot read.
func TestKeyEntityTSPSourceKeyStoreRefusesUnknownType(t *testing.T) {
	if _, err := NewKeyEntityTSPSourceFromKeyStore([]byte{0x30}, "JKS", "pwd", "", "pwd"); err == nil {
		t.Error("NewKeyEntityTSPSourceFromKeyStore() accepted a JKS store, want an error")
	}
}
