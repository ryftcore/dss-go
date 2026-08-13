package pfx

import (
	"crypto/dsa" //nolint:staticcheck
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"os"
	"testing"
)

func mustReadFixture(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading fixture %s: %v", path, err)
	}
	return data
}

func TestLoadEd25519Fixture(t *testing.T) {
	data := mustReadFixture(t, "testdata/Ed25519-good-user.p12")
	store, err := Load(data, "ks-password")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(store.PrivateKeys) != 1 {
		t.Fatalf("len(PrivateKeys) = %d, want 1", len(store.PrivateKeys))
	}
	key, ok := store.PrivateKeys[0].Key.(ed25519.PrivateKey)
	if !ok {
		t.Fatalf("PrivateKeys[0].Key is %T, want ed25519.PrivateKey", store.PrivateKeys[0].Key)
	}
	if len(key) != ed25519.PrivateKeySize {
		t.Errorf("len(key) = %d, want %d", len(key), ed25519.PrivateKeySize)
	}
	if len(store.Certificates) != 3 {
		t.Fatalf("len(Certificates) = %d, want 3", len(store.Certificates))
	}
	for i, cert := range store.Certificates {
		if _, err := x509.ParseCertificate(cert.Raw); err != nil {
			t.Errorf("Certificates[%d] does not parse: %v", i, err)
		}
	}
	leafPub := key.Public().(ed25519.PublicKey)
	leafCert, err := x509.ParseCertificate(store.Certificates[0].Raw)
	if err != nil {
		t.Fatalf("parsing leaf certificate: %v", err)
	}
	certPub, ok := leafCert.PublicKey.(ed25519.PublicKey)
	if !ok || string(certPub) != string(leafPub) {
		t.Errorf("leaf certificate public key does not match the private key's public half")
	}
}

func TestLoadDSAFixture(t *testing.T) {
	data := mustReadFixture(t, "testdata/good-dsa-user.p12")
	store, err := Load(data, "ks-password")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(store.PrivateKeys) != 1 {
		t.Fatalf("len(PrivateKeys) = %d, want 1", len(store.PrivateKeys))
	}
	key, ok := store.PrivateKeys[0].Key.(*dsa.PrivateKey)
	if !ok {
		t.Fatalf("PrivateKeys[0].Key is %T, want *dsa.PrivateKey", store.PrivateKeys[0].Key)
	}
	if len(store.Certificates) != 3 {
		t.Fatalf("len(Certificates) = %d, want 3", len(store.Certificates))
	}

	// Find the certificate matching this key's LocalKeyID, and check the recomputed Y matches
	// the certificate's own public key exactly - that is the strongest evidence the hand-rolled
	// DSA PKCS#8 parse (P, Q, G, X) is correct, since Y = G^X mod P is never itself encoded in a
	// PKCS#8 DSA PrivateKeyInfo.
	var leaf *x509.Certificate
	for _, cert := range store.Certificates {
		if string(cert.LocalKeyID) == string(store.PrivateKeys[0].LocalKeyID) {
			parsed, err := x509.ParseCertificate(cert.Raw)
			if err != nil {
				t.Fatalf("parsing matched certificate: %v", err)
			}
			leaf = parsed
			break
		}
	}
	if leaf == nil {
		t.Fatal("no certificate matched the private key's LocalKeyID")
	}
	certPub, ok := leaf.PublicKey.(*dsa.PublicKey)
	if !ok {
		t.Fatalf("leaf certificate public key is %T, want *dsa.PublicKey", leaf.PublicKey)
	}
	if key.Y.Cmp(certPub.Y) != 0 {
		t.Errorf("recomputed Y does not match the certificate's public key")
	}
	if key.P.Cmp(certPub.P) != 0 || key.Q.Cmp(certPub.Q) != 0 || key.G.Cmp(certPub.G) != 0 {
		t.Errorf("DSA domain parameters do not match the certificate's public key")
	}

	// crypto/dsa.Sign/Verify round trip, independent of the token package.
	hash := []byte("0123456789012345678901234567890123456789") // 40 bytes, more than Q needs
	r, s, err := dsa.Sign(rand.Reader, key, hash)
	if err != nil {
		t.Fatalf("dsa.Sign: %v", err)
	}
	if !dsa.Verify(&key.PublicKey, hash, r, s) {
		t.Error("dsa.Verify: signature does not verify")
	}
}

func TestLoadECDSAFixtureChain(t *testing.T) {
	data := mustReadFixture(t, "testdata/good-ecdsa-user.p12")
	store, err := Load(data, "ks-password")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(store.PrivateKeys) != 1 {
		t.Fatalf("len(PrivateKeys) = %d, want 1", len(store.PrivateKeys))
	}
	if _, ok := store.PrivateKeys[0].Key.(*ecdsa.PrivateKey); !ok {
		t.Fatalf("PrivateKeys[0].Key is %T, want *ecdsa.PrivateKey", store.PrivateKeys[0].Key)
	}
	if len(store.Certificates) != 3 {
		t.Fatalf("len(Certificates) = %d, want 3", len(store.Certificates))
	}
}

func TestLoadUserARSAFixture(t *testing.T) {
	data := mustReadFixture(t, "testdata/user_a_rsa.p12")
	store, err := Load(data, "password")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(store.PrivateKeys) != 1 {
		t.Fatalf("len(PrivateKeys) = %d, want 1", len(store.PrivateKeys))
	}
	if _, ok := store.PrivateKeys[0].Key.(*rsa.PrivateKey); !ok {
		t.Fatalf("PrivateKeys[0].Key is %T, want *rsa.PrivateKey", store.PrivateKeys[0].Key)
	}
}

func TestLoadGeneratedLegacyFixtures(t *testing.T) {
	for _, testCase := range []struct {
		fixture string
		keyType string
	}{
		{"testdata/rsa_test.p12", "rsa"},
		{"testdata/ec_test.p12", "ec"},
	} {
		t.Run(testCase.fixture, func(t *testing.T) {
			data := mustReadFixture(t, testCase.fixture)
			store, err := Load(data, "testpassword")
			if err != nil {
				t.Fatalf("Load: %v", err)
			}
			if len(store.PrivateKeys) != 1 {
				t.Fatalf("len(PrivateKeys) = %d, want 1", len(store.PrivateKeys))
			}
			if len(store.Certificates) != 1 {
				t.Fatalf("len(Certificates) = %d, want 1", len(store.Certificates))
			}
			switch testCase.keyType {
			case "rsa":
				if _, ok := store.PrivateKeys[0].Key.(*rsa.PrivateKey); !ok {
					t.Fatalf("Key is %T, want *rsa.PrivateKey", store.PrivateKeys[0].Key)
				}
			case "ec":
				if _, ok := store.PrivateKeys[0].Key.(*ecdsa.PrivateKey); !ok {
					t.Fatalf("Key is %T, want *ecdsa.PrivateKey", store.PrivateKeys[0].Key)
				}
			}
		})
	}
}

func TestLoadWrongPassword(t *testing.T) {
	data := mustReadFixture(t, "testdata/user_a_rsa.p12")
	if _, err := Load(data, "wrong password"); err == nil {
		t.Fatal("expected an error")
	}
}

// TestDeriveKeyMaterialAgainstReference checks the RFC 7292 Appendix B.2 key derivation against
// a value computed by an independent implementation of the same pseudocode (not this package's
// own), for the same salt/password/iteration-count/purpose/size inputs - the strongest
// available check on the multi-block ("I" mutation) code path, which every fixture above only
// exercises indirectly.
func TestDeriveKeyMaterialAgainstReference(t *testing.T) {
	salt := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	password := []byte{0, 65, 0, 66, 0, 0}
	got := deriveKeyMaterial(pfxHashSHA1, 1, salt, password, 1000, 24)
	want := "7c16f906bb3a79c8b9dd07a82ab37db68358313a2d6f56e3"
	if hex.EncodeToString(got) != want {
		t.Errorf("got  %x\nwant %s", got, want)
	}
}

// TestLoadPBES2Fixture exercises the modern openssl/JDK default: PBES2 with PBKDF2 and
// AES-256-CBC privacy, and a SHA-256 MacData - neither of which any other fixture here uses.
func TestLoadPBES2Fixture(t *testing.T) {
	data := mustReadFixture(t, "testdata/modrsa_test.p12")
	store, err := Load(data, "modernpass")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(store.PrivateKeys) != 1 {
		t.Fatalf("len(PrivateKeys) = %d, want 1", len(store.PrivateKeys))
	}
	if _, ok := store.PrivateKeys[0].Key.(*rsa.PrivateKey); !ok {
		t.Fatalf("Key is %T, want *rsa.PrivateKey", store.PrivateKeys[0].Key)
	}
	if len(store.Certificates) != 1 {
		t.Fatalf("len(Certificates) = %d, want 1", len(store.Certificates))
	}
}

// TestLoadRC4Fixture exercises the RC4 stream-cipher legacy PBE scheme, which - unlike 3DES and
// RC2 - derives no IV and needs no block padding.
func TestLoadRC4Fixture(t *testing.T) {
	data := mustReadFixture(t, "testdata/rc4_test.p12")
	store, err := Load(data, "rc4pass")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(store.PrivateKeys) != 1 {
		t.Fatalf("len(PrivateKeys) = %d, want 1", len(store.PrivateKeys))
	}
	if _, ok := store.PrivateKeys[0].Key.(*rsa.PrivateKey); !ok {
		t.Fatalf("Key is %T, want *rsa.PrivateKey", store.PrivateKeys[0].Key)
	}
}
