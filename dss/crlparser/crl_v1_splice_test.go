package crlparser

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"
	"testing"
	"time"

	"golang.org/x/crypto/cryptobyte"
	cryptobyte_asn1 "golang.org/x/crypto/cryptobyte/asn1"

	"github.com/ryftcore/dss-go/dss/model"
)

// buildVersionlessCRL builds a genuine RFC 5280 v1 CRL - a TBSCertList with NO version field, no
// crlExtensions - listing revokedSerial, signed ecdsa-with-SHA256 by key over the TBSCertList
// bytes it embeds. crypto/x509 can parse and create v2 CRLs only, so the DER is assembled by hand.
func buildVersionlessCRL(t *testing.T, issuer *x509.Certificate, key *ecdsa.PrivateKey, revokedSerial int64) []byte {
	t.Helper()
	ecdsaWithSHA256 := asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 2}
	now := time.Now().UTC().Truncate(time.Second)

	var signatureAlgorithm cryptobyte.Builder
	signatureAlgorithm.AddASN1(cryptobyte_asn1.SEQUENCE, func(b *cryptobyte.Builder) {
		b.AddASN1ObjectIdentifier(ecdsaWithSHA256)
	})
	sigAlg, err := signatureAlgorithm.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	var tbs cryptobyte.Builder
	tbs.AddASN1(cryptobyte_asn1.SEQUENCE, func(b *cryptobyte.Builder) {
		b.AddBytes(sigAlg)
		b.AddBytes(issuer.RawSubject)
		b.AddASN1UTCTime(now.Add(-time.Hour))
		b.AddASN1UTCTime(now.Add(24 * time.Hour))
		b.AddASN1(cryptobyte_asn1.SEQUENCE, func(revoked *cryptobyte.Builder) { // revokedCertificates
			revoked.AddASN1(cryptobyte_asn1.SEQUENCE, func(entry *cryptobyte.Builder) {
				entry.AddASN1Int64(revokedSerial)
				entry.AddASN1UTCTime(now.Add(-30 * time.Minute))
			})
		})
	})
	tbsBytes, err := tbs.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	digest := sha256.Sum256(tbsBytes)
	signature, err := ecdsa.SignASN1(rand.Reader, key, digest[:])
	if err != nil {
		t.Fatal(err)
	}

	var crl cryptobyte.Builder
	crl.AddASN1(cryptobyte_asn1.SEQUENCE, func(b *cryptobyte.Builder) {
		b.AddBytes(tbsBytes)
		b.AddBytes(sigAlg)
		b.AddASN1BitString(signature)
	})
	der, err := crl.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	return der
}

// TestCRLUtilsBuildCRLValidity_VersionlessV1CRL is the dedicated unit test for the v1 splice
// (open question 2 of the review): a genuine versionless CRL is rejected by crypto/x509
// (checked here, so the test cannot silently stop exercising the splice), yet parses, has its
// signature verified against the CRL's ORIGINAL TBSCertList bytes, and answers revocation
// lookups.
func TestCRLUtilsBuildCRLValidity_VersionlessV1CRL(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "V1 CRL Test CA"},
		NotBefore:             time.Now().Add(-48 * time.Hour),
		NotAfter:              time.Now().Add(48 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	issuerCert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	issuerToken, err := model.NewCertificateToken(issuerCert)
	if err != nil {
		t.Fatal(err)
	}

	const revokedSerial = 4242
	crlDER := buildVersionlessCRL(t, issuerCert, key, revokedSerial)

	if _, err := x509.ParseRevocationList(crlDER); err == nil {
		t.Fatal("premise broken: crypto/x509 now parses a versionless CRL, so the splice is no longer exercised by this test")
	}

	crlBinary, err := CRLUtilsBuildCRLBinary(crlDER)
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLBinary: %v", err)
	}
	validity, err := CRLUtilsBuildCRLValidity(crlBinary, issuerToken)
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLValidity: %v", err)
	}
	if !validity.IsSignatureIntact() {
		t.Errorf("signature over the original TBSCertList must verify: %s", validity.SignatureInvalidityReason())
	}
	if !validity.IssuerX509PrincipalMatches() {
		t.Errorf("expected the issuer principal to match")
	}
	if !validity.IsCrlSignKeyUsage() || !validity.IsValid() {
		t.Errorf("expected a valid CRL (cRLSign key usage present)")
	}
	if !bytes.Equal(validity.DerEncoded(), crlDER) {
		t.Errorf("the CRL's original bytes must be kept, not the spliced copy")
	}
	if validity.ThisUpdate() == nil || validity.NextUpdate() == nil {
		t.Errorf("thisUpdate / nextUpdate must be read from the versionless CRL")
	}

	entry := CRLUtilsRevocationInfo(validity, big.NewInt(revokedSerial))
	if entry == nil {
		t.Fatalf("expected serial %d to be revoked", revokedSerial)
	}
	if entry.RevocationDate().IsZero() {
		t.Errorf("expected the revocation date to be read")
	}
	if CRLUtilsRevocationInfo(validity, big.NewInt(revokedSerial+1)) != nil {
		t.Errorf("a serial that is not listed must not be revoked")
	}

	// A tampered TBSCertList must still be caught: the signature is checked over the original
	// bytes, so flipping a byte of the revoked serial (inside the TBS) breaks it.
	tampered := bytes.Clone(crlDER)
	idx := bytes.LastIndex(tampered, []byte{0x10, 0x92}) // INTEGER 4242 = 02 02 10 92
	if idx < 0 {
		t.Fatal("test bug: could not find the revoked serial in the CRL")
	}
	tampered[idx+1] ^= 0x01
	tamperedBinary, err := CRLUtilsBuildCRLBinary(tampered)
	if err != nil {
		t.Fatal(err)
	}
	tamperedValidity, err := CRLUtilsBuildCRLValidity(tamperedBinary, issuerToken)
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLValidity(tampered): %v", err)
	}
	if tamperedValidity.IsSignatureIntact() || tamperedValidity.IsValid() {
		t.Errorf("a tampered versionless CRL must not verify")
	}
}
