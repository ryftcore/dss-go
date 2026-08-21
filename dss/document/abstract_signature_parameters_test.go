// Tests for AbstractSignatureParameters, matching
// dss-document/src/main/java/eu/europa/esig/dss/signature/AbstractSignatureParameters.java
// upstream.
package document

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

func generateDocumentTestCertificate(t *testing.T) (*ecdsa.PrivateKey, *model.CertificateToken) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %s", err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "document test signer"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %s", err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parse certificate: %s", err)
	}
	token, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("NewCertificateToken: %s", err)
	}
	return key, token
}

func TestAbstractSignatureParametersDetachedContentsPriority(t *testing.T) {
	p := NewAbstractSignatureParameters[*model.TimestampParameters]()

	// Nothing set: returns an empty (non-nil) slice.
	if got := p.DetachedContents(); got == nil || len(got) != 0 {
		t.Fatalf("DetachedContents() = %v, want empty slice", got)
	}

	// Setting detachedContents directly is used once populated.
	own := []model.DSSDocument{model.NewInMemoryDocument([]byte("own"))}
	p.SetDetachedContents(own)
	if got := p.DetachedContents(); len(got) != 1 || got[0] != own[0] {
		t.Fatalf("DetachedContents() = %v, want %v", got, own)
	}

	// The internal context (once it carries detached contents) takes priority over the
	// directly-set field.
	ctxDocs := []model.DSSDocument{model.NewInMemoryDocument([]byte("ctx"))}
	p.GetContext().SetDetachedContents(ctxDocs)
	if got := p.DetachedContents(); len(got) != 1 || got[0] != ctxDocs[0] {
		t.Fatalf("DetachedContents() = %v, want context contents %v (context takes priority)", got, ctxDocs)
	}
}

func TestAbstractSignatureParametersSetSigningCertificateDerivesEncryptionAlgorithm(t *testing.T) {
	p := NewAbstractSignatureParameters[*model.TimestampParameters]()
	_, cert := generateDocumentTestCertificate(t)

	p.SetSigningCertificate(cert)

	if got := p.SigningCertificate(); got != cert {
		t.Fatalf("SigningCertificate() = %v, want %v", got, cert)
	}
	if got := p.EncryptionAlgorithm(); got != enumerations.EncryptionAlgorithm_ECDSA {
		t.Fatalf("EncryptionAlgorithm() = %v, want ECDSA (derived from the EC public key)", got)
	}
}

func TestAbstractSignatureParametersSetCertificateChainFromTokensDedupsAndSkipsNil(t *testing.T) {
	p := NewAbstractSignatureParameters[*model.TimestampParameters]()
	_, cert1 := generateDocumentTestCertificate(t)
	_, cert2 := generateDocumentTestCertificate(t)

	p.SetCertificateChainFromTokens(cert1, nil, cert2, cert1)

	chain := p.CertificateChain()
	if len(chain) != 2 {
		t.Fatalf("len(CertificateChain()) = %d, want 2 (dedup cert1, skip nil)", len(chain))
	}
	if chain[0] != cert1 || chain[1] != cert2 {
		t.Fatalf("CertificateChain() = %v, want [cert1, cert2] in insertion order", chain)
	}
}

func TestAbstractSignatureParametersSetCertificateChainReplaces(t *testing.T) {
	p := NewAbstractSignatureParameters[*model.TimestampParameters]()
	_, cert1 := generateDocumentTestCertificate(t)
	p.SetCertificateChainFromTokens(cert1)

	_, cert2 := generateDocumentTestCertificate(t)
	p.SetCertificateChain([]*model.CertificateToken{cert2})

	chain := p.CertificateChain()
	if len(chain) != 1 || chain[0] != cert2 {
		t.Fatalf("CertificateChain() = %v, want [cert2] (SetCertificateChain replaces)", chain)
	}
}

func TestAbstractSignatureParametersGetDeterministicIdIsCachedAndDeterministic(t *testing.T) {
	p := NewAbstractSignatureParameters[*model.TimestampParameters]()
	_, cert := generateDocumentTestCertificate(t)
	p.SetSigningCertificate(cert)
	signingDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	p.BLevel().SetSigningDate(&signingDate)

	id1 := p.GetDeterministicId()
	if id1 == "" {
		t.Fatal("GetDeterministicId() returned empty id")
	}
	id2 := p.GetDeterministicId()
	if id1 != id2 {
		t.Fatalf("GetDeterministicId() not stable across calls: %q != %q", id1, id2)
	}

	// Reinit clears the cached context, so a fresh id is computed (and, for identical inputs,
	// the same deterministic value comes back out).
	p.Reinit()
	id3 := p.GetDeterministicId()
	if id3 != id1 {
		t.Fatalf("GetDeterministicId() after Reinit() = %q, want same deterministic value %q", id3, id1)
	}
}

func TestAbstractSignatureParametersReinitClearsContext(t *testing.T) {
	p := NewAbstractSignatureParameters[*model.TimestampParameters]()
	p.GetContext().SetDeterministicId("stale")
	p.Reinit()
	if got := p.GetContext().DeterministicId(); got != "" {
		t.Fatalf("GetContext().DeterministicId() after Reinit() = %q, want empty (fresh context)", got)
	}
}

func TestAbstractSignatureParametersEquals(t *testing.T) {
	p1 := NewAbstractSignatureParameters[*model.TimestampParameters]()
	p2 := NewAbstractSignatureParameters[*model.TimestampParameters]()
	// BLevelParameters defaults its signing date to time.Now(), which differs by nanoseconds
	// between the two otherwise-identical instances constructed above; pin it so the two
	// parameter sets start out genuinely equal.
	signingDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	p1.BLevel().SetSigningDate(&signingDate)
	p2.BLevel().SetSigningDate(&signingDate)
	if !p1.Equals(&p2) {
		t.Fatal("two AbstractSignatureParameters with identical fields should be equal")
	}
	if !p1.Equals(&p1) {
		t.Fatal("a value must equal itself")
	}
	if p1.Equals(nil) {
		t.Fatal("a value must not equal nil")
	}

	_, cert := generateDocumentTestCertificate(t)
	p1.SetSigningCertificate(cert)
	if p1.Equals(&p2) {
		t.Fatal("differing signingCertificate should not be equal")
	}
	p2.SetSigningCertificate(cert)
	if !p1.Equals(&p2) {
		t.Fatal("matching signingCertificate should be equal again")
	}

	p1.SetSignedData([]byte{1, 2, 3})
	if p1.Equals(&p2) {
		t.Fatal("differing signedData should not be equal")
	}
}

func TestAbstractSignatureParametersString(t *testing.T) {
	p := NewAbstractSignatureParameters[*model.TimestampParameters]()
	if got := p.String(); got == "" {
		t.Fatal("String() should not be empty")
	}
}
