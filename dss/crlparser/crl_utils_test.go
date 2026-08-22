package crlparser

import (
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// crlparserReadFile reads a testdata file, failing the test on error.
func crlparserReadFile(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("reading testdata/%s: %v", name, err)
	}
	return data
}

// crlparserLoadCertificateToken parses a PEM certificate fixture into a CertificateToken.
func crlparserLoadCertificateToken(t *testing.T, name string) *model.CertificateToken {
	t.Helper()
	data := crlparserReadFile(t, name)
	block, _ := pem.Decode(data)
	if block == nil {
		t.Fatalf("testdata/%s: not a PEM file", name)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("testdata/%s: parsing certificate: %v", name, err)
	}
	token, err := model.NewCertificateToken(cert)
	if err != nil {
		t.Fatalf("testdata/%s: building CertificateToken: %v", name, err)
	}
	return token
}

func TestCRLUtilsBuildCRLBinary(t *testing.T) {
	der := crlparserReadFile(t, "ca.crl.der")
	pemBytes := crlparserReadFile(t, "ca.crl.pem")

	t.Run("DER", func(t *testing.T) {
		binary, err := CRLUtilsBuildCRLBinary(der)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(binary.Binaries()) != string(der) {
			t.Fatalf("expected the DER bytes to be preserved verbatim")
		}
	})

	t.Run("PEM", func(t *testing.T) {
		binary, err := CRLUtilsBuildCRLBinary(pemBytes)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(binary.Binaries()) != string(der) {
			t.Fatalf("expected the PEM to convert to the same DER bytes as the DER fixture")
		}
	})

	t.Run("empty", func(t *testing.T) {
		if _, err := CRLUtilsBuildCRLBinary(nil); err == nil {
			t.Fatalf("expected an error for empty input")
		}
	})

	t.Run("garbage", func(t *testing.T) {
		garbage := crlparserReadFile(t, "notacrl.bin")
		// notacrl.bin is 32 random bytes; it is astronomically unlikely to start with 0x30 or
		// '-', which is what this test means to exercise. Skip the rare collision rather than
		// flake.
		if garbage[0] == crlUtilsDERSequenceTag || garbage[0] == '-' {
			t.Skip("random fixture happens to start with a DER/PEM marker byte")
		}
		if _, err := CRLUtilsBuildCRLBinary(garbage); err == nil {
			t.Fatalf("expected an error for content that is neither DER nor PEM")
		}
	})
}

func TestPemToDerConverterConvert(t *testing.T) {
	der := crlparserReadFile(t, "ca.crl.der")
	pemBytes := crlparserReadFile(t, "ca.crl.pem")

	got, err := PemToDerConverterConvert(pemBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != string(der) {
		t.Fatalf("PEM conversion did not reproduce the DER fixture")
	}

	if _, err := PemToDerConverterConvert([]byte("not a pem file")); err == nil {
		t.Fatalf("expected an error for non-PEM content")
	}
}

func TestCRLUtilsBuildCRLValidity_Valid(t *testing.T) {
	issuer := crlparserLoadCertificateToken(t, "ca.crt")
	crlBinary, err := CRLUtilsBuildCRLBinary(crlparserReadFile(t, "ca.crl.der"))
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLBinary: %v", err)
	}

	validity, err := CRLUtilsBuildCRLValidity(crlBinary, issuer)
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLValidity: %v", err)
	}

	if !validity.IssuerX509PrincipalMatches() {
		t.Errorf("expected the issuer principal to match")
	}
	if !validity.IsSignatureIntact() {
		t.Errorf("expected the signature to be intact, invalidity reason: %q", validity.SignatureInvalidityReason())
	}
	if !validity.IsCrlSignKeyUsage() {
		t.Errorf("expected the issuer to carry the cRLSign key usage")
	}
	if !validity.IsValid() {
		t.Errorf("expected the CRLValidity to be valid")
	}
	if validity.SignatureAlgorithm() != enumerations.SignatureAlgorithm_RSA_SHA256 {
		t.Errorf("SignatureAlgorithm() = %v, want RSA_SHA256", validity.SignatureAlgorithm())
	}
	if validity.ThisUpdate() == nil {
		t.Errorf("expected ThisUpdate to be set")
	}
	if validity.NextUpdate() == nil {
		t.Errorf("expected NextUpdate to be set")
	}
	if validity.CRLNumber() == nil || validity.CRLNumber().Cmp(big.NewInt(4096)) != 0 {
		t.Errorf("CRLNumber() = %v, want 4096", validity.CRLNumber())
	}
	if validity.URL() != "http://crl.dss-go-port.test/ca.crl" {
		t.Errorf("URL() = %q, want the IDP full name URI", validity.URL())
	}
	if !validity.AreCriticalExtensionsOidNotEmpty() {
		t.Errorf("expected the critical issuingDistributionPoint OID to be recorded")
	}
	if validity.IsUnknownCriticalExtension() {
		t.Errorf("did not expect an unknown critical extension (URL is set, no reasonFlags, not every onlyXxx flag)")
	}
	if validity.X509CRL() == nil {
		t.Errorf("expected the parsed RevocationList to be cached")
	}
	if validity.IssuerToken() != issuer {
		t.Errorf("expected the issuer token to be recorded once the signature validated")
	}
	if validity.CrlBinary() != crlBinary {
		t.Errorf("expected CrlBinary() to return back the input CRLBinary")
	}
	if len(validity.DerEncoded()) == 0 {
		t.Errorf("expected DerEncoded() to be non-empty")
	}
}

func TestCRLUtilsBuildCRLValidity_WrongIssuer(t *testing.T) {
	crlBinary, err := CRLUtilsBuildCRLBinary(crlparserReadFile(t, "ca.crl.der"))
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLBinary: %v", err)
	}
	unrelated := crlparserLoadCertificateToken(t, "other.crt")

	validity, err := CRLUtilsBuildCRLValidity(crlBinary, unrelated)
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLValidity: %v", err)
	}
	if validity.IssuerX509PrincipalMatches() {
		t.Errorf("did not expect the issuer principal to match an unrelated CA")
	}
	if validity.IsValid() {
		t.Errorf("did not expect the CRLValidity to be valid")
	}
}

func TestCRLUtilsBuildCRLValidity_WrongKey(t *testing.T) {
	crlBinary, err := CRLUtilsBuildCRLBinary(crlparserReadFile(t, "ca.crl.der"))
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLBinary: %v", err)
	}
	// Same subject DN as the real CA, but a different key pair: the issuer principal matches
	// while the signature does not validate.
	impostor := crlparserLoadCertificateToken(t, "ca_impostor.crt")

	validity, err := CRLUtilsBuildCRLValidity(crlBinary, impostor)
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLValidity: %v", err)
	}
	if !validity.IssuerX509PrincipalMatches() {
		t.Errorf("expected the issuer principal to match (same subject DN)")
	}
	if validity.IsSignatureIntact() {
		t.Errorf("did not expect the signature to validate against the wrong key")
	}
	if validity.SignatureInvalidityReason() == "" {
		t.Errorf("expected a signature invalidity reason to be recorded")
	}
	if validity.IsCrlSignKeyUsage() {
		t.Errorf("crlSignKeyUsage should stay false: it is only set once the signature is intact")
	}
	if validity.IsValid() {
		t.Errorf("did not expect the CRLValidity to be valid")
	}
}

func TestCRLUtilsBuildCRLValidity_NoCRLSignKeyUsage(t *testing.T) {
	crlBinary, err := CRLUtilsBuildCRLBinary(crlparserReadFile(t, "ca.crl.der"))
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLBinary: %v", err)
	}
	// Same subject and key as the real CA, but its own certificate does not carry the cRLSign
	// key usage.
	noSign := crlparserLoadCertificateToken(t, "ca_no_crlsign.crt")

	validity, err := CRLUtilsBuildCRLValidity(crlBinary, noSign)
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLValidity: %v", err)
	}
	if !validity.IssuerX509PrincipalMatches() {
		t.Errorf("expected the issuer principal to match")
	}
	if !validity.IsSignatureIntact() {
		t.Errorf("expected the signature to validate (same key pair)")
	}
	if validity.IsCrlSignKeyUsage() {
		t.Errorf("did not expect the cRLSign key usage to be present")
	}
	if validity.SignatureInvalidityReason() == "" {
		t.Errorf("expected the missing key usage reason to be recorded")
	}
	if validity.IsValid() {
		t.Errorf("did not expect the CRLValidity to be valid")
	}
}

func TestCRLUtilsBuildCRLValidity_Minimal(t *testing.T) {
	// ca_minimal2.crl.der is a V2 CRL with no revoked certificates and a single non-critical
	// extension (authorityKeyIdentifier): no crlNumber, no distribution point, no critical
	// extensions. crypto/x509.ParseRevocationList only supports V2 CRLs, so this fixture
	// stands in for upstream's "nothing set" case rather than an actual V1 (versionless) CRL.
	issuer := crlparserLoadCertificateToken(t, "ca.crt")
	crlBinary, err := CRLUtilsBuildCRLBinary(crlparserReadFile(t, "ca_minimal2.crl.der"))
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLBinary: %v", err)
	}

	validity, err := CRLUtilsBuildCRLValidity(crlBinary, issuer)
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLValidity: %v", err)
	}
	if validity.AreCriticalExtensionsOidNotEmpty() {
		t.Errorf("did not expect any critical extensions")
	}
	if validity.IsUnknownCriticalExtension() {
		t.Errorf("an empty critical extension set can never be 'unknown'")
	}
	if validity.CRLNumber() != nil {
		t.Errorf("did not expect a CRL Number extension")
	}
	if validity.URL() != "" {
		t.Errorf("did not expect a distribution point URL")
	}
	if !validity.IsValid() {
		t.Errorf("expected the CRLValidity to be valid")
	}
	if CRLUtilsRevocationInfo(validity, big.NewInt(1)) != nil {
		t.Errorf("expected no revoked entries on the minimal CRL")
	}
}

func TestCRLUtilsBuildCRLValidity_MalformedCRL(t *testing.T) {
	issuer := crlparserLoadCertificateToken(t, "ca.crt")
	garbage := crlparserReadFile(t, "notacrl.bin")
	// Force the DER-looking branch by prefixing the SEQUENCE tag; the content underneath is
	// still garbage and must fail to parse as a CertificateList.
	fake := append([]byte{crlUtilsDERSequenceTag, byte(len(garbage))}, garbage...)

	crlBinary, err := CRLUtilsBuildCRLBinary(fake)
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLBinary: %v", err)
	}
	if _, err := CRLUtilsBuildCRLValidity(crlBinary, issuer); err == nil {
		t.Fatalf("expected an error parsing a malformed CRL")
	}
}

func TestCRLUtilsRevocationInfo(t *testing.T) {
	issuer := crlparserLoadCertificateToken(t, "ca.crt")
	crlBinary, err := CRLUtilsBuildCRLBinary(crlparserReadFile(t, "ca.crl.der"))
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLBinary: %v", err)
	}
	validity, err := CRLUtilsBuildCRLValidity(crlBinary, issuer)
	if err != nil {
		t.Fatalf("CRLUtilsBuildCRLValidity: %v", err)
	}

	t.Run("revoked with reason", func(t *testing.T) {
		entry := CRLUtilsRevocationInfo(validity, big.NewInt(0x1001))
		if entry == nil {
			t.Fatalf("expected serial 0x1001 to be revoked")
		}
		if entry.SerialNumber().Cmp(big.NewInt(0x1001)) != 0 {
			t.Errorf("SerialNumber() = %v, want 0x1001", entry.SerialNumber())
		}
		if entry.RevocationDate().IsZero() {
			t.Errorf("expected a non-zero revocation date")
		}
		reason := entry.RevocationReason()
		if reason == nil {
			t.Fatalf("expected a revocation reason to be present")
		}
		if enumerations.RevocationReasonFromInt(*reason) != enumerations.RevocationReason_KEY_COMPROMISE {
			t.Errorf("RevocationReason() = %d, want keyCompromise (1)", *reason)
		}
	})

	t.Run("revoked without reason", func(t *testing.T) {
		entry := CRLUtilsRevocationInfo(validity, big.NewInt(0x1002))
		if entry == nil {
			t.Fatalf("expected serial 0x1002 to be revoked")
		}
		if entry.RevocationReason() != nil {
			t.Errorf("expected no revocation reason to be present, got %v", *entry.RevocationReason())
		}
	})

	t.Run("not revoked", func(t *testing.T) {
		if entry := CRLUtilsRevocationInfo(validity, big.NewInt(0x9999)); entry != nil {
			t.Errorf("expected no entry for an unrevoked serial, got %v", entry)
		}
	})

	t.Run("reparses when the cache is empty", func(t *testing.T) {
		uncached := NewCRLValidity(crlBinary)
		entry := CRLUtilsRevocationInfo(uncached, big.NewInt(0x1001))
		if entry == nil {
			t.Fatalf("expected the fallback re-parse path to still find the entry")
		}
	})
}
