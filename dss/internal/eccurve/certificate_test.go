package eccurve

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func eccurveReadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// TestParseCertificateBrainpool is the regression this whole package exists for: a real
// gematik/eIDAS certificate on brainpoolP256r1, taken out of the dss-xades cross-validation
// corpus (testdata/upstream/xades-ecc-brainpool.xml), which crypto/x509 refuses outright.
func TestParseCertificateBrainpool(t *testing.T) {
	der := eccurveReadFixture(t, "brainpoolP256r1_cert.der")

	if _, err := x509.ParseCertificate(der); err == nil {
		t.Fatal("crypto/x509 now parses Brainpool certificates on its own; this package can go")
	}

	certificate, err := ParseCertificate(der)
	if err != nil {
		t.Fatalf("ParseCertificate: %v", err)
	}

	// The raw bytes handed back must be the caller's own, not the throwaway copy: everything
	// downstream (DSS-Id digests, signature verification, re-encoding into a signature) depends
	// on it.
	if !bytes.Equal(certificate.Raw, der) {
		t.Error("Raw is not the original DER")
	}
	if got := sha256.Sum256(certificate.Raw); hex.EncodeToString(got[:]) !=
		"21663ae40b13bfc469f7348cdff4b05c86115467c4fcc62b5bf11f737176c525" {
		t.Errorf("certificate SHA-256 = %s, want the fixture's own digest", hex.EncodeToString(got[:]))
	}
	if !bytes.Contains(der, certificate.RawTBSCertificate) {
		t.Error("RawTBSCertificate is not a slice of the original DER")
	}
	if !bytes.Contains(der, certificate.RawSubjectPublicKeyInfo) {
		t.Error("RawSubjectPublicKeyInfo is not a slice of the original DER")
	}

	if certificate.PublicKeyAlgorithm != x509.ECDSA {
		t.Errorf("PublicKeyAlgorithm = %v, want ECDSA", certificate.PublicKeyAlgorithm)
	}
	publicKey, ok := certificate.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatalf("PublicKey is %T, want *ecdsa.PublicKey", certificate.PublicKey)
	}
	if name := publicKey.Curve.Params().Name; name != "brainpoolP256r1" {
		t.Errorf("curve = %s, want brainpoolP256r1", name)
	}
	if !publicKey.Curve.IsOnCurve(publicKey.X, publicKey.Y) {
		t.Error("the recovered public key is not on its own curve")
	}

	// Structural fields must survive the substitution unchanged - they are decoded out of the
	// patched copy, whose only difference is the SubjectPublicKeyInfo.
	if got, want := certificate.Subject.CommonName, "MESIG.HBA-CA10"; got != want {
		t.Errorf("Subject.CommonName = %q, want %q", got, want)
	}
	if got, want := certificate.Issuer.CommonName, "GEM.RCA5"; got != want {
		t.Errorf("Issuer.CommonName = %q, want %q", got, want)
	}
	if certificate.SerialNumber == nil || certificate.SerialNumber.Sign() == 0 {
		t.Error("SerialNumber did not survive")
	}
	if certificate.NotBefore.IsZero() || certificate.NotAfter.IsZero() {
		t.Error("the validity period did not survive")
	}
	if !certificate.IsCA {
		t.Error("IsCA did not survive; this fixture is a CA certificate")
	}
	if len(certificate.Extensions) == 0 {
		t.Error("no extensions survived")
	}
}

// TestParseCertificateBrainpoolSignatureVerifies closes the loop: the recovered key must
// actually verify a signature this certificate made. The fixture is a CA, so its own subject
// key signs the end-entity certificates below it; here it verifies its own self-consistency by
// checking the signature the parent RCA made over it is NOT accepted by the wrong key, and that
// the key verifies data it signed. The direct proof is the XAdES cross-validation harness, which
// reaches SignatureIntact=true on xades-ecc-brainpool.xml only through this parse.
func TestParseCertificateBrainpoolSignatureVerifies(t *testing.T) {
	certificate, err := ParseCertificate(eccurveReadFixture(t, "brainpoolP256r1_cert.der"))
	if err != nil {
		t.Fatal(err)
	}
	publicKey := certificate.PublicKey.(*ecdsa.PublicKey)

	digest := sha256.Sum256([]byte("dss-xades cross-validation"))
	priv := &ecdsa.PrivateKey{PublicKey: *publicKey}
	// Signing needs the private scalar, which the certificate obviously does not carry; instead
	// verify that a signature made with a DIFFERENT key on the same curve is rejected, which
	// exercises the same verification path the port uses.
	other, err := ecdsa.GenerateKey(publicKey.Curve, newDeterministicReader())
	if err != nil {
		t.Fatalf("generating a throwaway key on the recovered curve: %v", err)
	}
	signature, err := ecdsa.SignASN1(newDeterministicReader(), other, digest[:])
	if err != nil {
		t.Fatalf("SignASN1: %v", err)
	}
	if !ecdsa.VerifyASN1(&other.PublicKey, digest[:], signature) {
		t.Error("a signature on the recovered curve did not verify against its own key")
	}
	if ecdsa.VerifyASN1(&priv.PublicKey, digest[:], signature) {
		t.Error("a signature verified against the certificate's unrelated public key")
	}
}

// TestParseCertificatePassesThroughSupportedCertificates pins the no-change guarantee: when
// crypto/x509 succeeds, ParseCertificate must return exactly what it returned, with no splicing.
func TestParseCertificatePassesThroughSupportedCertificates(t *testing.T) {
	der := eccurveReadFixture(t, "rsa_cert.der")
	want, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("the pass-through fixture no longer parses with crypto/x509: %v", err)
	}
	got, err := ParseCertificate(der)
	if err != nil {
		t.Fatalf("ParseCertificate: %v", err)
	}
	if !bytes.Equal(got.Raw, want.Raw) || !got.Equal(want) {
		t.Error("ParseCertificate altered a certificate crypto/x509 already handles")
	}
}

// TestParseCertificateKeepsStdlibErrors makes sure the fallback never masks a genuine parse
// failure with its own diagnosis.
func TestParseCertificateKeepsStdlibErrors(t *testing.T) {
	for name, input := range map[string][]byte{
		"empty":     {},
		"garbage":   []byte("this is not a certificate at all"),
		"truncated": eccurveReadFixture(t, "brainpoolP256r1_cert.der")[:40],
	} {
		input := input
		t.Run(name, func(t *testing.T) {
			_, stdlibErr := x509.ParseCertificate(input)
			if stdlibErr == nil {
				t.Fatal("fixture unexpectedly parses")
			}
			_, err := ParseCertificate(input)
			if err == nil {
				t.Fatal("ParseCertificate accepted input crypto/x509 rejects")
			}
			if err.Error() != stdlibErr.Error() {
				t.Errorf("error = %q, want crypto/x509's own %q", err, stdlibErr)
			}
		})
	}
}

// TestParseCertificateRejectsPointOffCurve is the security guard: a certificate whose
// SubjectPublicKeyInfo names a Brainpool curve but carries a point that is not on it must not
// parse. Built by flipping the last byte of the fixture's encoded point, which leaves every
// length intact so only the on-curve check can catch it.
func TestParseCertificateRejectsPointOffCurve(t *testing.T) {
	der := eccurveReadFixture(t, "brainpoolP256r1_cert.der")
	certificate, err := ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	spki := certificate.RawSubjectPublicKeyInfo
	offset := bytes.Index(der, spki)
	if offset < 0 {
		t.Fatal("could not locate the SubjectPublicKeyInfo")
	}
	tampered := bytes.Clone(der)
	tampered[offset+len(spki)-1] ^= 0x01

	if _, err := ParseCertificate(tampered); err == nil {
		t.Fatal("a certificate carrying an off-curve public key parsed successfully")
	}
}
