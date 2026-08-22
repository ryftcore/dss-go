// Tests for the four DSSKeyEntryPredicate implementations. testdata/generated/eku.crt was
// produced with:
//
//	openssl req -x509 -newkey rsa:2048 -keyout eku.key -out eku.crt -days 3650 -nodes \
//	    -subj "/CN=Go Port Test EKU/O=DSS Go Port" \
//	    -addext "keyUsage=critical,digitalSignature,keyCertSign" \
//	    -addext "extendedKeyUsage=clientAuth,codeSigning"
//
// (the private key was discarded - these predicates only ever read Certificate()).
package token

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// fakePrivateKeyEntry implements DSSPrivateKeyEntry over a bare certificate, for predicate tests
// that never need to sign anything.
type fakePrivateKeyEntry struct {
	certificate *model.CertificateToken
}

func (f *fakePrivateKeyEntry) Certificate() *model.CertificateToken        { return f.certificate }
func (f *fakePrivateKeyEntry) CertificateChain() []*model.CertificateToken { return nil }
func (f *fakePrivateKeyEntry) EncryptionAlgorithm() enumerations.EncryptionAlgorithm {
	return enumerations.EncryptionAlgorithmRSA
}

func mustLoadCertificateToken(t *testing.T, path string) *model.CertificateToken {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		t.Fatalf("no PEM block in %s", path)
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("ParseCertificate: %v", err)
	}
	token, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("NewCertificateToken: %v", err)
	}
	return token
}

func TestAllKeyEntryPredicate(t *testing.T) {
	predicate := NewAllKeyEntryPredicate()
	entry := &fakePrivateKeyEntry{certificate: mustLoadCertificateToken(t, "testdata/generated/eku.crt")}
	if !predicate(entry) {
		t.Error("AllKeyEntryPredicate rejected an entry; it must accept everything")
	}
}

func TestKeyUsageKeyEntryPredicate(t *testing.T) {
	entry := &fakePrivateKeyEntry{certificate: mustLoadCertificateToken(t, "testdata/generated/eku.crt")}

	if predicate := NewKeyUsageKeyEntryPredicate(enumerations.KeyUsageBitDigitalSignature); !predicate(entry) {
		t.Error("expected DIGITAL_SIGNATURE to match: the certificate carries it")
	}
	if predicate := NewKeyUsageKeyEntryPredicate(enumerations.KeyUsageBitKeyCertSign); !predicate(entry) {
		t.Error("expected KEY_CERT_SIGN to match: the certificate carries it")
	}
	if predicate := NewKeyUsageKeyEntryPredicate(enumerations.KeyUsageBitCRLSign); predicate(entry) {
		t.Error("expected CRL_SIGN not to match: the certificate does not carry it")
	}
	// Matches if ANY of the requested bits is present.
	if predicate := NewKeyUsageKeyEntryPredicate(enumerations.KeyUsageBitCRLSign, enumerations.KeyUsageBitKeyCertSign); !predicate(entry) {
		t.Error("expected a match: one of the two requested bits is present")
	}
}

func TestKeyUsageKeyEntryPredicateNoCertificate(t *testing.T) {
	predicate := NewKeyUsageKeyEntryPredicate(enumerations.KeyUsageBitDigitalSignature)
	if predicate(&fakePrivateKeyEntry{certificate: nil}) {
		t.Error("expected false for a nil certificate")
	}
}

func TestKeyUsageKeyEntryPredicatePanicsOnNil(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic for a nil keyUsages argument")
		}
	}()
	var nilUsages []enumerations.KeyUsageBit
	NewKeyUsageKeyEntryPredicate(nilUsages...)
}

func TestExtendedKeyUsageKeyEntryPredicate(t *testing.T) {
	entry := &fakePrivateKeyEntry{certificate: mustLoadCertificateToken(t, "testdata/generated/eku.crt")}

	if predicate := NewExtendedKeyUsageKeyEntryPredicate(enumerations.ExtendedKeyUsageClientAuth); !predicate(entry) {
		t.Error("expected CLIENT_AUTH to match: the certificate carries it")
	}
	if predicate := NewExtendedKeyUsageKeyEntryPredicate(enumerations.ExtendedKeyUsageCodeSigning); !predicate(entry) {
		t.Error("expected CODE_SIGNING to match: the certificate carries it")
	}
	if predicate := NewExtendedKeyUsageKeyEntryPredicate(enumerations.ExtendedKeyUsageServerAuth); predicate(entry) {
		t.Error("expected SERVER_AUTH not to match: the certificate does not carry it")
	}
}

func TestExtendedKeyUsageKeyEntryPredicateForOIDs(t *testing.T) {
	entry := &fakePrivateKeyEntry{certificate: mustLoadCertificateToken(t, "testdata/generated/eku.crt")}

	// clientAuth OID.
	predicate := NewExtendedKeyUsageKeyEntryPredicateForOIDs("1.3.6.1.5.5.7.3.2")
	if !predicate(entry) {
		t.Error("expected the raw clientAuth OID to match")
	}
	predicate = NewExtendedKeyUsageKeyEntryPredicateForOIDs("1.2.3.4.5.6.7.8.9")
	if predicate(entry) {
		t.Error("expected an unrelated OID not to match")
	}
}

func TestExtendedKeyUsageKeyEntryPredicateNoExtension(t *testing.T) {
	// user_a_rsa.p12's certificate carries no ExtendedKeyUsage extension.
	data := mustReadFixture(t, "testdata/dss-token/src/test/resources/user_a_rsa.p12")
	signatureToken, err := NewPkcs12SignatureTokenFromBytes(data, NewPasswordProtection([]byte("password")))
	if err != nil {
		t.Fatalf("NewPkcs12SignatureTokenFromBytes: %v", err)
	}
	defer signatureToken.Close()
	keys, err := signatureToken.Keys()
	if err != nil || len(keys) != 1 {
		t.Fatalf("Keys: %v (len=%d)", err, len(keys))
	}

	predicate := NewExtendedKeyUsageKeyEntryPredicate(enumerations.ExtendedKeyUsageClientAuth)
	if predicate(keys[0]) {
		t.Error("expected no match: the certificate carries no ExtendedKeyUsage extension")
	}
}

func TestValidAtTimeKeyEntryPredicate(t *testing.T) {
	entry := &fakePrivateKeyEntry{certificate: mustLoadCertificateToken(t, "testdata/generated/eku.crt")}

	notBefore := entry.certificate.NotBefore()
	notAfter := entry.certificate.NotAfter()

	tests := []struct {
		name string
		at   time.Time
		want bool
	}{
		{"well within validity", notBefore.Add(24 * time.Hour), true},
		{"exactly notBefore", notBefore, true},
		{"exactly notAfter", notAfter, true},
		{"before notBefore", notBefore.Add(-time.Second), false},
		{"after notAfter", notAfter.Add(time.Second), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			predicate := NewValidAtTimeKeyEntryPredicateAt(tt.at)
			if got := predicate(entry); got != tt.want {
				t.Errorf("predicate(at=%v) = %v, want %v", tt.at, got, tt.want)
			}
		})
	}
}

func TestValidAtTimeKeyEntryPredicateDefaultsToNow(t *testing.T) {
	entry := &fakePrivateKeyEntry{certificate: mustLoadCertificateToken(t, "testdata/generated/eku.crt")}
	// The fixture is valid for 3650 days from generation, so "now" must fall inside its range.
	if predicate := NewValidAtTimeKeyEntryPredicate(); !predicate(entry) {
		t.Error("expected the current time to fall within the certificate's validity range")
	}
}

func TestValidAtTimeKeyEntryPredicateNoCertificate(t *testing.T) {
	predicate := NewValidAtTimeKeyEntryPredicate()
	if predicate(&fakePrivateKeyEntry{certificate: nil}) {
		t.Error("expected false for a nil certificate")
	}
}
