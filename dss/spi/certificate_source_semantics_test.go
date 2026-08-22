package spi

import (
	"crypto/x509"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/model"
)

// certificateSourceSemanticsTestToken loads one of the DER certificates of testdata/asn1 as a
// fresh CertificateToken; calling it twice for the same file yields two distinct instances of
// the same certificate, which is what the deduplication checks below need.
func certificateSourceSemanticsTestToken(t *testing.T, name string) *model.CertificateToken {
	t.Helper()
	der, err := os.ReadFile(filepath.Join("testdata", "asn1", name))
	if err != nil {
		t.Fatalf("unable to read %s: %v", name, err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("unable to parse %s: %v", name, err)
	}
	token, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("unable to build a CertificateToken for %s: %v", name, err)
	}
	return token
}

// certificateSourceSemanticsTestIDs renders a certificate collection as its sorted DSS Ids.
func certificateSourceSemanticsTestIDs(tokens []*model.CertificateToken) string {
	ids := make([]string, 0, len(tokens))
	for _, token := range tokens {
		ids = append(ids, token.DSSIDAsString())
	}
	sort.Strings(ids)
	return strings.Join(ids, ",")
}

// TestCertificateSourcesAcceptCertificatesFromTheirConstructors checks that every certificate
// source this package builds is usable straight out of its constructor.
//
// Java runs CommonCertificateSource's field initialisers through the implicit super() call of
// every subclass; a Go constructor that only fills its own fields leaves the embedded maps nil
// and the first addCertificate panics with "assignment to entry in nil map". OCSPCertificateSource
// went down that path for every OCSP response carrying a responder certificate.
func TestCertificateSourcesAcceptCertificatesFromTheirConstructors(t *testing.T) {
	certificate := certificateSourceSemanticsTestToken(t, "ocsp_ca.der")

	t.Run("CommonCertificateSource", func(t *testing.T) {
		source := NewCommonCertificateSource()
		source.AddCertificate(certificate)
		if source.NumberOfCertificates() != 1 {
			t.Errorf("NumberOfCertificates() = %d, want 1", source.NumberOfCertificates())
		}
	})

	t.Run("zero value CommonCertificateSource", func(t *testing.T) {
		var source CommonCertificateSource
		source.AddCertificate(certificate)
		if source.NumberOfCertificates() != 1 {
			t.Errorf("NumberOfCertificates() = %d, want 1", source.NumberOfCertificates())
		}
	})

	t.Run("CommonTrustedCertificateSource", func(t *testing.T) {
		source := NewCommonTrustedCertificateSource()
		source.AddCertificate(certificate)
		if !source.IsTrusted(certificate) {
			t.Errorf("expected the added certificate to be trusted")
		}
	})

	t.Run("KidCertificateSource", func(t *testing.T) {
		source := NewKidCertificateSource()
		source.AddCertificate(certificate)
		if source.NumberOfCertificates() != 1 {
			t.Errorf("NumberOfCertificates() = %d, want 1", source.NumberOfCertificates())
		}
	})

	t.Run("KeyStoreCertificateSource", func(t *testing.T) {
		source, err := NewKeyStoreCertificateSource(KeyStoreCertificateSourceTypePEM, nil)
		if err != nil {
			t.Fatalf("NewKeyStoreCertificateSource: %v", err)
		}
		source.AddCertificateToKeyStore(certificate)
		if source.NumberOfCertificates() != 1 {
			t.Errorf("NumberOfCertificates() = %d, want 1", source.NumberOfCertificates())
		}
	})

	t.Run("CommonX509URLCertificateSource", func(t *testing.T) {
		source := NewCommonX509URLCertificateSource()
		source.AddCertificateForURL("http://localhost/cert", certificate)
		if source.NumberOfCertificates() != 1 {
			t.Errorf("NumberOfCertificates() = %d, want 1", source.NumberOfCertificates())
		}
	})

	t.Run("RevocationCertificateSourceBase", func(t *testing.T) {
		source := NewRevocationCertificateSourceBase()
		source.AddCertificate(certificate)
		if source.NumberOfCertificates() != 1 {
			t.Errorf("NumberOfCertificates() = %d, want 1", source.NumberOfCertificates())
		}
	})

	t.Run("zero value TokenCertificateSource", func(t *testing.T) {
		var source TokenCertificateSource
		source.AddCertificate(certificate)
		if source.NumberOfCertificates() != 1 {
			t.Errorf("NumberOfCertificates() = %d, want 1", source.NumberOfCertificates())
		}
	})
}

// TestCommonTrustedCertificateSourceIsTrustedAtTime pins the inherited
// isTrustedAtTime(CertificateToken, Date), whose Java body is `return isTrusted(certificateToken)`
// and therefore dispatches to the trusted source's own isTrusted override.
func TestCommonTrustedCertificateSourceIsTrustedAtTime(t *testing.T) {
	trusted := certificateSourceSemanticsTestToken(t, "ocsp_ca.der")
	untrusted := certificateSourceSemanticsTestToken(t, "ocsp_leaf.der")

	source := NewCommonTrustedCertificateSource()
	source.AddCertificate(trusted)

	controlTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	if !source.IsTrustedAtTime(trusted, controlTime) {
		t.Errorf("IsTrustedAtTime(trusted) = false, want true")
	}
	if source.IsTrustedAtTime(untrusted, controlTime) {
		t.Errorf("IsTrustedAtTime(untrusted) = true, want false")
	}

	// The base source trusts nothing, at any time.
	base := NewCommonCertificateSource()
	base.AddCertificate(trusted)
	if base.IsTrustedAtTime(trusted, controlTime) {
		t.Errorf("CommonCertificateSource.IsTrustedAtTime() = true, want false")
	}
}

// TestCertificateReordererDeduplicatesEqualTokens pins getAllCertificatesOnce's deduplication,
// which relies on List#contains, i.e. on Token#equals (DSS-Id equality) and NOT on object
// identity: the same certificate loaded twice is one entry of the chain.
func TestCertificateReordererDeduplicatesEqualTokens(t *testing.T) {
	ca := certificateSourceSemanticsTestToken(t, "ocsp_ca.der")
	caAgain := certificateSourceSemanticsTestToken(t, "ocsp_ca.der")
	leaf := certificateSourceSemanticsTestToken(t, "ocsp_leaf.der")
	leafAgain := certificateSourceSemanticsTestToken(t, "ocsp_leaf.der")
	if ca == caAgain || leaf == leafAgain {
		t.Fatalf("the fixture loader must return distinct instances")
	}
	expected := certificateSourceSemanticsTestIDs([]*model.CertificateToken{leaf, ca})

	testCases := []struct {
		name      string
		reorderer *CertificateReorderer
	}{
		{
			name:      "duplicate in the chain",
			reorderer: NewCertificateReorderer([]*model.CertificateToken{leaf, ca, caAgain}),
		},
		{
			name:      "signing certificate repeated in the chain",
			reorderer: NewCertificateReordererWithSigningCertificate(leafAgain, []*model.CertificateToken{leaf, ca}),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ordered, err := testCase.reorderer.OrderedCertificates()
			if err != nil {
				t.Fatalf("OrderedCertificates: %v", err)
			}
			if len(ordered) != 2 {
				t.Fatalf("OrderedCertificates() returned %d certificates, want 2", len(ordered))
			}
			if got := certificateSourceSemanticsTestIDs(ordered); got != expected {
				t.Errorf("OrderedCertificates() = %s, want %s", got, expected)
			}
		})
	}

	t.Run("chains keyed by an equal token collapse to one", func(t *testing.T) {
		chains, err := NewCertificateReorderer([]*model.CertificateToken{ca, caAgain}).OrderedCertificateChains()
		if err != nil {
			t.Fatalf("OrderedCertificateChains: %v", err)
		}
		if len(chains) != 1 {
			t.Errorf("OrderedCertificateChains() returned %d chains, want 1", len(chains))
		}
	})
}
