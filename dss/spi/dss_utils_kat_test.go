// Known-answer / fidelity tests for dss_utils.go. Test vectors are either well-known NIST
// examples (the "abc" digest KATs) or generated locally with OpenSSL 3.0 into
// testdata/dss_utils/ (a self-signed RSA-2048 certificate and PKCS#7 certs-only "p7c" bundles
// in DER and PEM form, built with `openssl req -x509 ...` / `openssl crl2pkcs7 -nocrl ...`).
package spi

import (
	"bytes"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/utils"
)

func TestDSSUtilsDigestKAT(t *testing.T) {
	data := []byte("abc")
	cases := []struct {
		algo enumerations.DigestAlgorithm
		want string
	}{
		{enumerations.DigestAlgorithm_MD5, "900150983cd24fb0d6963f7d28e17f72"},
		{enumerations.DigestAlgorithm_SHA1, "a9993e364706816aba3e25717850c26c9cd0d89d"},
		{enumerations.DigestAlgorithm_SHA256, "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
		{enumerations.DigestAlgorithm_SHA3_256, "3a985da74fe225b2045c172d6bd390bd855f086e3e9d525b46bfe24511431532"},
		// BouncyCastle's SHAKEDigest#getDigestSize() is fixedOutputLength/4, so SHAKE-128
		// squeezes 32 bytes and SHAKE-256 64 - twice the security strength, not once.
		// Captured from upstream DSSUtils#digest on BouncyCastle 1.84.
		{enumerations.DigestAlgorithm_SHAKE128, "5881092dd818bf5cf8a3ddb793fbcba74097d5c526a6d35f97b83351940f2cc8"},
		{enumerations.DigestAlgorithm_SHAKE256, "483366601360a8771c6863080cc4114d8db44530f8f1e1ee4f94ea37e78b5739d5a15bef186a5386c75744c0527e1faa9f8726e462a12a4feb06bd8801e751e4"},
		{enumerations.DigestAlgorithm_SHAKE256_512, "483366601360a8771c6863080cc4114d8db44530f8f1e1ee4f94ea37e78b5739d5a15bef186a5386c75744c0527e1faa9f8726e462a12a4feb06bd8801e751e4"},
	}
	for _, c := range cases {
		got, err := DSSUtilsDigest(c.algo, data)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", c.algo, err)
		}
		if hex.EncodeToString(got) != c.want {
			t.Errorf("%s: got %s want %s", c.algo, hex.EncodeToString(got), c.want)
		}
	}
}

func TestDSSUtilsDigestUnsupported(t *testing.T) {
	for _, algo := range []enumerations.DigestAlgorithm{"MD2", "WHIRLPOOL"} {
		if _, err := DSSUtilsDigest(algo, []byte("x")); err == nil {
			t.Errorf("%s: expected error, got nil", algo)
		}
	}
}

func TestDSSUtilsToHex(t *testing.T) {
	// DSSUtils.toHex delegates to Utils.toHex, which is lowercase (commons-codec
	// Hex.encodeHexString), despite what the upstream javadoc claims - see the DEVIATION note
	// on DSSUtilsToHex.
	if got := DSSUtilsToHex([]byte{0x01, 0xab}); got != "01ab" {
		t.Errorf("got %q", got)
	}
	if got := DSSUtilsToHex(nil); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestDSSUtilsFormatDateToRFC(t *testing.T) {
	date := time.Date(2019, time.November, 19, 17, 28, 15, 0, time.UTC)
	got := DSSUtilsFormatDateToRFC(date)
	want := "2019-11-19T17:28:15Z"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
	if got := DSSUtilsFormatDateToRFC(time.Time{}); got != "N/A" {
		t.Errorf("zero date: got %q want N/A", got)
	}
}

func TestDSSUtilsParseRFCDateRoundTrip(t *testing.T) {
	want := time.Date(2019, time.November, 19, 17, 28, 15, 0, time.UTC)
	str := DSSUtilsFormatDateToRFC(want)
	got := DSSUtilsParseRFCDate(str)
	if !got.Equal(want) {
		t.Errorf("got %v want %v", got, want)
	}
	if !DSSUtilsIsRFCDate(str) {
		t.Errorf("IsRFCDate(%q) = false, want true", str)
	}
	if DSSUtilsIsRFCDate("2019-11-19T17:28:15+02:00") {
		t.Errorf("IsRFCDate should reject a non-Z offset (upstream pattern has a literal 'Z')")
	}
	if !DSSUtilsParseRFCDate("not-a-date").IsZero() {
		t.Errorf("ParseRFCDate(garbage) should return zero Time")
	}
}

func TestDSSUtilsFormatDateToISO8601(t *testing.T) {
	date := time.Date(2019, time.November, 19, 17, 28, 15, 0, time.UTC)
	if got := DSSUtilsFormatDateToISO8601(date); got != "2019-11-19" {
		t.Errorf("got %q", got)
	}
	roundTrip := DSSUtilsParseISO8601Date("2001-01-01")
	want := time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC)
	if !roundTrip.Equal(want) {
		t.Errorf("got %v want %v", roundTrip, want)
	}
	if !DSSUtilsIsISO8601Date("2001-01-01") {
		t.Errorf("expected valid ISO8601 date")
	}
	if DSSUtilsIsISO8601Date("2001-13-01") {
		t.Errorf("expected month 13 to be rejected")
	}
}

func TestDSSUtilsFormatDateWithCustomFormat(t *testing.T) {
	date := time.Date(2019, time.November, 19, 17, 28, 0, 0, time.UTC)
	got := DSSUtilsFormatDateWithCustomFormat(date, "yyyyMMdd-HHmm")
	if got != "20191119-1728" {
		t.Errorf("got %q", got)
	}
}

func TestDSSUtilsUTCDate(t *testing.T) {
	// month is 0-based, matching java.util.Calendar: 0 = January.
	got := DSSUtilsUTCDate(2019, 0, 15)
	want := time.Date(2019, time.January, 15, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestDSSUtilsOidHelpers(t *testing.T) {
	if !DSSUtilsIsUrnOid("urn:oid:1.2.3") {
		t.Errorf("expected urn:oid:1.2.3 to be recognised")
	}
	if DSSUtilsIsUrnOid("1.2.3") {
		t.Errorf("expected 1.2.3 to not be a urn oid")
	}
	if !DSSUtilsIsOidCode("1.3.6.1.4.1.343") {
		t.Errorf("expected valid oid code")
	}
	if DSSUtilsIsOidCode("25.25") {
		t.Errorf("expected 25.25 to be invalid (first arc must be 0-2)")
	}
	if DSSUtilsIsOidCode("http://sample.com") {
		t.Errorf("expected URL to not be a valid oid code")
	}
	if got := DSSUtilsOidCode("urn:oid:1.2.3"); got != "1.2.3" {
		t.Errorf("got %q", got)
	}
	if got := DSSUtilsToUrnOid("1.2.3"); got != "urn:oid:1.2.3" {
		t.Errorf("got %q", got)
	}
	if got := DSSUtilsObjectIdentifierValue("urn:oid:1.2.3"); got != "1.2.3" {
		t.Errorf("got %q", got)
	}
	if got := DSSUtilsObjectIdentifierValue("http://website.com"); got != "http://website.com" {
		t.Errorf("got %q", got)
	}
	if got := DSSUtilsObjectIdentifierValue("1.2.3"); got != "1.2.3" {
		t.Errorf("got %q", got)
	}
}

func TestDSSUtilsTrimAndStrip(t *testing.T) {
	// Java: removes \n and \r characters entirely, then trims: "  a\nb\r  " -> "  ab  " ->
	// "ab".
	got := DSSUtilsTrimWhitespacesAndNewlines("  a\nb\r  ")
	if got != "ab" {
		t.Errorf("got %q want %q", got, "ab")
	}
	stripped, err := DSSUtilsStripFirstLeadingOccurrence("prefix-value", "prefix-")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stripped != "value" {
		t.Errorf("got %q", stripped)
	}
}

func TestDSSUtilsDecodeURI(t *testing.T) {
	if got := DSSUtilsDecodeURI("a%20b"); got != "a b" {
		t.Errorf("got %q", got)
	}
	if got := DSSUtilsDecodeURI("a+b"); got != "a+b" {
		t.Errorf("expected literal '+' to be preserved, got %q", got)
	}
}

func TestDSSUtilsEncodeURIPreservesUnicode(t *testing.T) {
	got := DSSUtilsEncodeURI("café file.txt")
	if got != "café%20file.txt" {
		t.Errorf("got %q", got)
	}
}

func TestDSSUtilsHost(t *testing.T) {
	got := DSSUtilsHost("ldap://ldap.infonotary.com/dc=identity-ca,dc=infonotary,dc=com")
	if got != "ldap.infonotary.com" {
		t.Errorf("got %q", got)
	}
	if got := DSSUtilsHost("https://example.com:8443/path"); got != "example.com" {
		t.Errorf("got %q", got)
	}
}

func TestDSSUtilsCertificatePEMDERRoundTrip(t *testing.T) {
	derCert, err := os.ReadFile("testdata/dss_utils/cert.der")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}
	token, err := DSSUtilsLoadCertificateFromBinary(derCert)
	if err != nil {
		t.Fatalf("LoadCertificateFromBinary(DER): %v", err)
	}

	pemStr := DSSUtilsConvertToPEM(token)
	der2, err := DSSUtilsConvertToDER(pemStr)
	if err != nil {
		t.Fatalf("ConvertToDER: %v", err)
	}
	if !bytes.Equal(der2, derCert) {
		t.Errorf("PEM->DER roundtrip mismatch")
	}

	tokenFromPEM, err := DSSUtilsLoadCertificateFromBinary([]byte(pemStr))
	if err != nil {
		t.Fatalf("LoadCertificateFromBinary(PEM): %v", err)
	}
	if !bytes.Equal(tokenFromPEM.Certificate().Raw, token.Certificate().Raw) {
		t.Errorf("certificate loaded from PEM does not match DER original")
	}
}

func TestDSSUtilsLoadCertificateFromFile(t *testing.T) {
	token, err := DSSUtilsLoadCertificate("testdata/dss_utils/cert.pem")
	if err != nil {
		t.Fatalf("LoadCertificate: %v", err)
	}
	if token.Certificate().Subject.CommonName != "DSS Go Port Test" {
		t.Errorf("got CN=%q", token.Certificate().Subject.CommonName)
	}
}

func TestDSSUtilsLoadCertificateFromBase64(t *testing.T) {
	derCert, err := os.ReadFile("testdata/dss_utils/cert.der")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}
	encoded := utils.ToBase64(derCert)
	token, err := DSSUtilsLoadCertificateFromBase64EncodedString(encoded)
	if err != nil {
		t.Fatalf("LoadCertificateFromBase64EncodedString: %v", err)
	}
	if !bytes.Equal(token.Certificate().Raw, derCert) {
		t.Errorf("mismatch")
	}
}

func TestDSSUtilsLoadCertificateFromP7c(t *testing.T) {
	derCert, err := os.ReadFile("testdata/dss_utils/cert.der")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}

	for _, name := range []string{"testdata/dss_utils/bundle.p7c", "testdata/dss_utils/bundle_pem.p7c"} {
		tokens, err := DSSUtilsLoadCertificateFromP7c(name)
		if err != nil {
			t.Fatalf("%s: LoadCertificateFromP7c: %v", name, err)
		}
		if len(tokens) != 1 {
			t.Fatalf("%s: got %d certificates, want 1", name, len(tokens))
		}
		if !bytes.Equal(tokens[0].Certificate().Raw, derCert) {
			t.Errorf("%s: certificate mismatch", name)
		}
	}
}

func TestDSSUtilsGenerateKid(t *testing.T) {
	token, err := DSSUtilsLoadCertificate("testdata/dss_utils/cert.pem")
	if err != nil {
		t.Fatalf("LoadCertificate: %v", err)
	}
	kid := DSSUtilsGenerateKid(token)
	if len(kid) == 0 {
		t.Errorf("expected a non-empty DER-encoded IssuerSerial")
	}
}
