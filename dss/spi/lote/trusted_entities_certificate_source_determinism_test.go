package lote

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/lote"
	"github.com/ryftcore/dss-go/dss/model/tsl"
)

// testTokens creates count distinct self-signed certificate tokens.
func testTokens(t *testing.T, count int) []*model.CertificateToken {
	t.Helper()
	tokens := make([]*model.CertificateToken, 0, count)
	for index := 0; index < count; index++ {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		template := &x509.Certificate{
			SerialNumber: big.NewInt(int64(index + 1)),
			Subject:      pkix.Name{CommonName: "lote determinism " + string(rune('A'+index))},
			NotBefore:    time.Now().Add(-time.Hour),
			NotAfter:     time.Now().Add(time.Hour),
		}
		der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
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
		tokens = append(tokens, token)
	}
	return tokens
}

// TestTrustedEntitiesCertificateOrderIsStable requires the certificates of a
// TrustedEntitiesCertificateSource to come out in the same (DSS id) order on every run: the
// setters take a Go map, whose iteration order is random, whereas upstream's HashMap order is
// stable for a given key set (T23-STD-001; spi/tsl.TrustedListsCertificateSource does the same).
func TestTrustedEntitiesCertificateOrderIsStable(t *testing.T) {
	tokens := testTokens(t, 24)
	expected := make([]string, 0, len(tokens))
	for _, token := range tokens {
		expected = append(expected, token.DSSIDAsString())
	}
	sort.Strings(expected)

	certificateIDs := func(source *TrustedEntitiesCertificateSource) []string {
		var ids []string
		for _, certificate := range source.Certificates() {
			ids = append(ids, certificate.DSSIDAsString())
		}
		return ids
	}
	assertOrder := func(what string, got []string) {
		t.Helper()
		if len(got) != len(expected) {
			t.Fatalf("%s: %d certificates, want %d", what, len(got), len(expected))
		}
		for index := range expected {
			if got[index] != expected[index] {
				t.Fatalf("%s: certificate %d is %s, want %s (DSS id order)", what, index, got[index], expected[index])
			}
		}
	}

	for run := 0; run < 5; run++ {
		propertiesByCertificate := make(map[*model.CertificateToken][]*lote.TrustedProperties, len(tokens))
		trustTimeByCertificate := make(map[*model.CertificateToken][]*tsl.CertificateTrustTime, len(tokens))
		for _, token := range tokens {
			propertiesByCertificate[token] = []*lote.TrustedProperties{}
			trustTimeByCertificate[token] = []*tsl.CertificateTrustTime{}
		}

		source := NewTrustedEntitiesCertificateSource()
		source.SetTrustedPropertiesByCertificates(propertiesByCertificate)
		assertOrder("SetTrustedPropertiesByCertificates", certificateIDs(source))

		source = NewTrustedEntitiesCertificateSource()
		source.SetTrustedTimeByCertificates(trustTimeByCertificate)
		assertOrder("SetTrustedTimeByCertificates", certificateIDs(source))
	}
}
