package token

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/internal/pfx"
)

// pkcs12TestCA is a certificate, its key and the DER that goes into a pfx.Store.
type pkcs12TestCert struct {
	cert *x509.Certificate
	der  []byte
	key  *ecdsa.PrivateKey
}

var pkcs12TestSerial int64

// newPKCS12TestCert issues a certificate for subject, signed by issuer (nil = self-signed).
func newPKCS12TestCert(t *testing.T, subject string, isCA bool, issuer *pkcs12TestCert) *pkcs12TestCert {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pkcs12TestSerial++
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(pkcs12TestSerial),
		Subject:               pkix.Name{CommonName: subject},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  isCA,
	}
	parent, parentKey := template, key
	if issuer != nil {
		parent, parentKey = issuer.cert, issuer.key
	}
	der, err := x509.CreateCertificate(rand.Reader, template, parent, &key.PublicKey, parentKey)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return &pkcs12TestCert{cert: cert, der: der, key: key}
}

func chainSubjects(entry keyStoreEntry) []string {
	subjects := make([]string, len(entry.chain))
	for i, c := range entry.chain {
		subjects[i] = c.Subject.CommonName
	}
	return subjects
}

func assertChain(t *testing.T, entry keyStoreEntry, want ...string) {
	t.Helper()
	got := chainSubjects(entry)
	if len(got) != len(want) {
		t.Fatalf("entry %q: chain %v, want %v", entry.alias, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("entry %q: chain %v, want %v", entry.alias, got, want)
		}
	}
	if entry.certificate != entry.chain[0] {
		t.Errorf("entry %q: certificate must be the head of its chain", entry.alias)
	}
}

// assertOwnKey checks that the entry's leaf certificate carries the public key of its own signer.
func assertOwnKey(t *testing.T, entry keyStoreEntry, owner *pkcs12TestCert) {
	t.Helper()
	pub, ok := entry.privateKey.Public().(*ecdsa.PublicKey)
	if !ok || !pub.Equal(&owner.key.PublicKey) {
		t.Fatalf("entry %q: wrong private key for its entry", entry.alias)
	}
	if !bytes.Equal(entry.certificate.Raw, owner.der) {
		t.Errorf("entry %q: leaf certificate is not the certificate of its key (%s)", entry.alias, entry.certificate.Subject.CommonName)
	}
}

// TestPkcs12BuildKeyStoreSharedChain covers T20-SEC-001. A PFX holding two private keys whose
// certificates are issued under the same intermediate and root must give BOTH entries the complete
// leaf-intermediate-root chain: java.security.KeyStore's PKCS12 provider builds every private key
// entry's chain independently over all certificates of the file, whereas the port consumed the
// shared chain certificates with the first key and left the second key with a chain truncated to
// its leaf - the chain that ends up embedded in the signatures made with that key.
func TestPkcs12BuildKeyStoreSharedChain(t *testing.T) {
	root := newPKCS12TestCert(t, "Root", true, nil)
	intermediate := newPKCS12TestCert(t, "Intermediate", true, root)
	leaf1 := newPKCS12TestCert(t, "Leaf 1", false, intermediate)
	leaf2 := newPKCS12TestCert(t, "Leaf 2", false, intermediate)

	store := &pfx.Store{
		Certificates: []pfx.Certificate{
			{Raw: leaf1.der, LocalKeyID: []byte{1}},
			{Raw: leaf2.der, LocalKeyID: []byte{2}},
			{Raw: intermediate.der},
			{Raw: root.der},
		},
		PrivateKeys: []pfx.PrivateKey{
			{Key: leaf1.key, LocalKeyID: []byte{1}, FriendlyName: "one"},
			{Key: leaf2.key, LocalKeyID: []byte{2}, FriendlyName: "two"},
		},
	}
	ks, err := pkcs12BuildKeyStore(store)
	if err != nil {
		t.Fatal(err)
	}
	if len(ks.entries) != 2 {
		t.Fatalf("got %d entries, want 2", len(ks.entries))
	}
	assertChain(t, ks.entries[0], "Leaf 1", "Intermediate", "Root")
	assertChain(t, ks.entries[1], "Leaf 2", "Intermediate", "Root")
	assertOwnKey(t, ks.entries[0], leaf1)
	assertOwnKey(t, ks.entries[1], leaf2)
}

// TestPkcs12BuildKeyStoreSharedChainWithoutLocalKeyID pins the "first unclaimed certificate"
// fallback for keys without a localKeyId: the CA certificates the first key's chain walked
// through must still not be mistaken for the second key's own certificate.
func TestPkcs12BuildKeyStoreSharedChainWithoutLocalKeyID(t *testing.T) {
	root := newPKCS12TestCert(t, "Root", true, nil)
	intermediate := newPKCS12TestCert(t, "Intermediate", true, root)
	leaf1 := newPKCS12TestCert(t, "Leaf 1", false, intermediate)
	leaf2 := newPKCS12TestCert(t, "Leaf 2", false, intermediate)

	store := &pfx.Store{
		Certificates: []pfx.Certificate{{Raw: leaf1.der}, {Raw: intermediate.der}, {Raw: root.der}, {Raw: leaf2.der}},
		PrivateKeys:  []pfx.PrivateKey{{Key: leaf1.key}, {Key: leaf2.key}},
	}
	ks, err := pkcs12BuildKeyStore(store)
	if err != nil {
		t.Fatal(err)
	}
	assertChain(t, ks.entries[0], "Leaf 1", "Intermediate", "Root")
	assertChain(t, ks.entries[1], "Leaf 2", "Intermediate", "Root")
	assertOwnKey(t, ks.entries[0], leaf1)
	assertOwnKey(t, ks.entries[1], leaf2)
	if ks.entries[0].alias != "1" || ks.entries[1].alias != "2" {
		t.Errorf("aliases = %q, %q, want the positional defaults", ks.entries[0].alias, ks.entries[1].alias)
	}
}

// TestPkcs12BuildKeyStoreCrossCertificateCycle pins the termination guard now that chain
// membership is tracked per chain: two CA certificates that certify each other (cross-
// certification) form an issuer/subject cycle that has no self-signed end.
func TestPkcs12BuildKeyStoreCrossCertificateCycle(t *testing.T) {
	// Build A and B so that A is issued by B and B is issued by A.
	keyA, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	keyB, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	nameA, nameB := pkix.Name{CommonName: "CA A"}, pkix.Name{CommonName: "CA B"}
	templateA := &x509.Certificate{SerialNumber: big.NewInt(9001), Subject: nameA, NotBefore: time.Now().Add(-time.Hour),
		NotAfter: time.Now().Add(time.Hour), BasicConstraintsValid: true, IsCA: true}
	templateB := &x509.Certificate{SerialNumber: big.NewInt(9002), Subject: nameB, NotBefore: time.Now().Add(-time.Hour),
		NotAfter: time.Now().Add(time.Hour), BasicConstraintsValid: true, IsCA: true}
	// A's issuer is B (parent template subject B) and vice versa.
	derA, err := x509.CreateCertificate(rand.Reader, templateA, &x509.Certificate{Subject: nameB}, &keyA.PublicKey, keyB)
	if err != nil {
		t.Fatal(err)
	}
	derB, err := x509.CreateCertificate(rand.Reader, templateB, &x509.Certificate{Subject: nameA}, &keyB.PublicKey, keyA)
	if err != nil {
		t.Fatal(err)
	}
	certA := &pkcs12TestCert{der: derA, key: keyA}
	certA.cert, _ = x509.ParseCertificate(derA)
	leaf := newPKCS12TestCert(t, "Leaf", false, certA)

	store := &pfx.Store{
		Certificates: []pfx.Certificate{{Raw: leaf.der, LocalKeyID: []byte{7}}, {Raw: derA}, {Raw: derB}},
		PrivateKeys:  []pfx.PrivateKey{{Key: leaf.key, LocalKeyID: []byte{7}}},
	}
	ks, err := pkcs12BuildKeyStore(store)
	if err != nil {
		t.Fatal(err)
	}
	assertChain(t, ks.entries[0], "Leaf", "CA A", "CA B")
}
