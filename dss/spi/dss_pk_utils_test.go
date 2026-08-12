package spi

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/utain/esig/dss/model"
)

// dssPKUtilsTestPublicKey marshals a Go key into a SubjectPublicKeyInfo and wraps it the way
// a parsed certificate would.
func dssPKUtilsTestPublicKey(t *testing.T, key any) *model.PublicKey {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		t.Fatalf("unable to marshal the key: %v", err)
	}
	publicKey, err := model.NewPublicKey(der)
	if err != nil {
		t.Fatalf("unable to parse the SubjectPublicKeyInfo: %v", err)
	}
	return publicKey
}

// TestDSSPKUtilsPublicKeySize checks the key sizes for every key type the port recognises.
func TestDSSPKUtilsPublicKeySize(t *testing.T) {
	// The fixture CA carries a 2048-bit RSA key.
	der, err := os.ReadFile(filepath.Join("testdata", "asn1", "ocsp_ca.der"))
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	rsaKey := model.NewPublicKeyFromEncoded(certificate.RawSubjectPublicKeyInfo, certificate.PublicKey)
	if got := DSSPKUtilsPublicKeySize(rsaKey); got != 2048 {
		t.Errorf("RSA: got %d, want 2048", got)
	}
	if got := DSSPKUtilsStringPublicKeySize(rsaKey); got != "2048" {
		t.Errorf("RSA: got %q, want \"2048\"", got)
	}

	// EC keys are sized by the field size of their curve.
	for _, entry := range []struct {
		curve elliptic.Curve
		size  int
	}{
		{elliptic.P256(), 256},
		{elliptic.P384(), 384},
		{elliptic.P521(), 521},
	} {
		key, err := ecdsa.GenerateKey(entry.curve, rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		publicKey := dssPKUtilsTestPublicKey(t, &key.PublicKey)
		if got := DSSPKUtilsPublicKeySize(publicKey); got != entry.size {
			t.Errorf("%s: got %d, want %d", entry.curve.Params().Name, got, entry.size)
		}
	}

	// UPSTREAM QUIRK: the EdDSA and XDH branches report a BYTE count, not a bit length.
	edPublic, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	edKey := dssPKUtilsTestPublicKey(t, edPublic)
	if len(edKey.Encoded()) != 44 {
		t.Fatalf("an Ed25519 SubjectPublicKeyInfo is 44 bytes, got %d", len(edKey.Encoded()))
	}
	if got := DSSPKUtilsPublicKeySize(edKey); got != 32 {
		t.Errorf("Ed25519: got %d, want 32 (the upstream byte count)", got)
	}

	xKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	xPublicKey := dssPKUtilsTestPublicKey(t, xKey.PublicKey())
	if got := DSSPKUtilsPublicKeySize(xPublicKey); got != 32 {
		t.Errorf("X25519: got %d, want 32 (the upstream byte count)", got)
	}
}

// TestDSSPKUtilsPublicKeySizeOfUnparsableKeys checks the OID fallback for the Ed448/X448 keys
// crypto/x509 cannot decode.
func TestDSSPKUtilsPublicKeySizeOfUnparsableKeys(t *testing.T) {
	for _, entry := range []struct {
		name     string
		oid      asn1.ObjectIdentifier
		keyBytes int
		want     int
	}{
		{"Ed448", asn1.ObjectIdentifier{1, 3, 101, 113}, 57, 57},
		{"X448", asn1.ObjectIdentifier{1, 3, 101, 111}, 56, 56},
	} {
		algorithm, err := asn1.Marshal(entry.oid)
		if err != nil {
			t.Fatal(err)
		}
		algorithm = append([]byte{0x30, byte(len(algorithm))}, algorithm...)
		bitString := append([]byte{0x03, byte(entry.keyBytes + 1), 0x00}, make([]byte, entry.keyBytes)...)
		body := append(algorithm, bitString...)
		spki := append([]byte{0x30, byte(len(body))}, body...)

		if _, err := x509.ParsePKIXPublicKey(spki); err == nil {
			t.Fatalf("%s: crypto/x509 unexpectedly parses this key; the fallback test is stale", entry.name)
		}
		publicKey := model.NewPublicKeyFromEncoded(spki, nil)
		if got := DSSPKUtilsPublicKeySize(publicKey); got != entry.want {
			t.Errorf("%s: got %d, want %d", entry.name, got, entry.want)
		}
	}
}

// TestDSSPKUtilsPublicKeySizeOfDSA checks the DSA branch against a hand-built
// SubjectPublicKeyInfo, since crypto/x509 can parse a DSA key but not marshal one.
func TestDSSPKUtilsPublicKeySizeOfDSA(t *testing.T) {
	// A 1024-bit p, a 160-bit q and matching g/y; the values need not form a valid group
	// for crypto/x509 to decode them, only to be positive.
	p := new(big.Int).Lsh(big.NewInt(1), 1023)
	p.Add(p, big.NewInt(2891))
	q := new(big.Int).Lsh(big.NewInt(1), 159)
	q.Add(q, big.NewInt(7))
	g := big.NewInt(2)
	y := big.NewInt(3)

	publicKeyValue, err := asn1.Marshal(y)
	if err != nil {
		t.Fatal(err)
	}
	spki, err := asn1.Marshal(struct {
		Algorithm struct {
			Algorithm  asn1.ObjectIdentifier
			Parameters struct{ P, Q, G *big.Int }
		}
		PublicKey asn1.BitString
	}{
		Algorithm: struct {
			Algorithm  asn1.ObjectIdentifier
			Parameters struct{ P, Q, G *big.Int }
		}{
			Algorithm:  asn1.ObjectIdentifier{1, 2, 840, 10040, 4, 1},
			Parameters: struct{ P, Q, G *big.Int }{P: p, Q: q, G: g},
		},
		PublicKey: asn1.BitString{Bytes: publicKeyValue, BitLength: 8 * len(publicKeyValue)},
	})
	if err != nil {
		t.Fatal(err)
	}

	publicKey, err := model.NewPublicKey(spki)
	if err != nil {
		t.Fatalf("unable to parse the DSA SubjectPublicKeyInfo: %v", err)
	}
	if got := DSSPKUtilsPublicKeySize(publicKey); got != 1024 {
		t.Errorf("DSA: got %d, want 1024", got)
	}
}

// TestDSSPKUtilsUnknownKey checks the "?" reported for a key the port cannot size.
func TestDSSPKUtilsUnknownKey(t *testing.T) {
	if got := DSSPKUtilsStringPublicKeySize(nil); got != "?" {
		t.Errorf("a missing key must be reported as \"?\", got %q", got)
	}
	// An RSA key whose SubjectPublicKeyInfo is unreadable and whose parsed key is absent.
	unknown := model.NewPublicKeyFromEncoded([]byte{0x30, 0x00}, nil)
	if got := DSSPKUtilsPublicKeySize(unknown); got != -1 {
		t.Errorf("an unknown key infrastructure must be reported as -1, got %d", got)
	}
	if got := DSSPKUtilsStringPublicKeySize(unknown); got != "?" {
		t.Errorf("an unknown key infrastructure must be reported as \"?\", got %q", got)
	}
}

// TestDSSPKUtilsStringPublicKeySizeOfToken checks the Token overload, which falls back to the
// certificate's own key when the token is self-signed.
func TestDSSPKUtilsStringPublicKeySizeOfToken(t *testing.T) {
	der, err := os.ReadFile(filepath.Join("testdata", "asn1", "ocsp_ca.der"))
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	token, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatal(err)
	}
	// The fixture CA is self-signed, so its own key sizes it even before any signature check.
	if !token.IsSelfSigned() {
		t.Fatal("the fixture CA is expected to be self-signed")
	}
	if got := DSSPKUtilsStringPublicKeySizeOfToken(token); got != "2048" {
		t.Errorf("got %q, want \"2048\"", got)
	}

	// A token that is neither self-signed nor has an established signer reports "?".
	leafDER, err := os.ReadFile(filepath.Join("testdata", "asn1", "ocsp_leaf.der"))
	if err != nil {
		t.Fatal(err)
	}
	leafCertificate, err := x509.ParseCertificate(leafDER)
	if err != nil {
		t.Fatal(err)
	}
	leaf, err := model.NewCertificateToken(leafCertificate)
	if err != nil {
		t.Fatal(err)
	}
	if got := DSSPKUtilsStringPublicKeySizeOfToken(leaf); got != "?" {
		t.Errorf("got %q, want \"?\"", got)
	}
	// Once the signer is established, the issuer's key size is reported.
	if !leaf.IsSignedByToken(token) {
		t.Fatal("the fixture leaf is expected to be signed by the fixture CA")
	}
	if got := DSSPKUtilsStringPublicKeySizeOfToken(leaf); got != "2048" {
		t.Errorf("got %q, want \"2048\"", got)
	}
}
