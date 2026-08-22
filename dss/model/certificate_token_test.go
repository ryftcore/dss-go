package model

import (
	"crypto/x509"
	"encoding/base64"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// rootCertificate is a fixed RSA test certificate; the expectations below are computed from
// its exact DER encoding.
const rootCertificateBase64 = "" +
	"MIIDTzCCAjegAwIBAgICEzcwDQYJKoZIhvcNAQELBQAwRzELMAkGA1UEBhMCTFUxGjAYBgNVBA" +
	"oMEU5vd2luYSBTw7ZsdXRpb25zMRwwGgYDVQQDExNEU1MgR28gUG9ydCBSb290IENBMB4XDTIx" +
	"MDMwNDA1MDYwN1oXDTMxMDMwNDA1MDYwN1owRzELMAkGA1UEBhMCTFUxGjAYBgNVBAoMEU5vd2" +
	"luYSBTw7ZsdXRpb25zMRwwGgYDVQQDExNEU1MgR28gUG9ydCBSb290IENBMIIBIjANBgkqhkiG" +
	"9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4gXVNplGr5jbwl94JS+n3xewJNKJrZc0bZrxFusbjuLcFa" +
	"r9xiz9q+PXm2e8/CSJyfFRA9OaMvJ/nnQUUzrgdJvxNhRKxBi4D/3grj63VBlmN5FUsuYwmTrH" +
	"6dZ4nl7HGxf3tFtE8nVe4DGNAM4USzEacl5g2DY7o1r25p+myin15nk+1EMnhGVuOFPdwYCpd+" +
	"PNZvaD+z/5pEht/485HNGOqehT/EKYUMmBuxlA8DrxBNhIycfYw2PthKDXBIh6rOCJoeDxAcWP" +
	"aXBu+UZbISHhGfyzSrniTPnFZkhNCHh+RkXAXimUjFR17NLGIyWdMpyeEnyqn5YckCgnKrc8Lw" +
	"IDAQABo0UwQzAOBgNVHQ8BAf8EBAMCAYYwEgYDVR0TAQH/BAgwBgEB/wIBAjAdBgNVHQ4EFgQU" +
	"Af72tIMHPtHTlemWdkMVlPxblo0wDQYJKoZIhvcNAQELBQADggEBAKXcsB764gV+ZdBw0FSPAs" +
	"HAnQXJFyZppBCGzBGREE9V9dGssIfpesrLpxLDAXG/Gh8hekaMoQgehEa87roGTkJKreuef9CB" +
	"LVzgL36DtrYK2fS92jq5xylpNQsQcQRXnVk1+XiVDx9SXVijd27pMpHeZ9NkT/FYmkrjZ50SWH" +
	"jGi3hpiA2faugQQFutq2jGkEVUmA2ZrxRBwSDfUdTuQ1LpMqxt/h8jBS1lfczunIghsiggLfRQ" +
	"aVk+IpYFdqH3Sj2oE9PX1p80wzjLaTvpyV7Fbwf4gW8/xvlZzfYD66UVPLFAWI221rW4AOMpNA" +
	"toVaereMOpEVC1a2RL3vQcP+w="

// leafCertificate is a fixed RSA test certificate; the expectations below are computed from
// its exact DER encoding.
const leafCertificateBase64 = "" +
	"MIIDTDCCAjSgAwIBAgIEEjRWeDANBgkqhkiG9w0BAQsFADBHMQswCQYDVQQGEwJMVTEaMBgGA1" +
	"UECgwRTm93aW5hIFPDtmx1dGlvbnMxHDAaBgNVBAMTE0RTUyBHbyBQb3J0IFJvb3QgQ0EwHhcN" +
	"MjEwMzA0MDUwNjA3WhcNMzEwMzA0MDUwNjA3WjBGMQswCQYDVQQGEwJMVTEaMBgGA1UECgwRTm" +
	"93aW5hIFPDtmx1dGlvbnMxGzAZBgNVBAMTEkRTUyBHbyBQb3J0IFNpZ25lcjCCASIwDQYJKoZI" +
	"hvcNAQEBBQADggEPADCCAQoCggEBALb7t8hpydWVsL/EUbgEsSkGO/SO4ZID1mibsyzutJcDks" +
	"9o3JzQCyOuEh2FUIsHvTCFHXNoIOAeYeXE61LcUmgR03co1huOHLZOvKZIKttn1QHyZx2C8OZ8" +
	"7ifzAa/Kb8monM4TxoRbOAESOv4pHn2zS3K6g3/oN3a/Om9IZ8J+ZYOf5xQeJH5Jo7+1E57suX" +
	"fxEkVhQMK2E70KomCrgj0ZAmq+RNec1ohHZUdLB2Uo3KipY4ExSj33mEDMZRMlYzOGY1+0EBDo" +
	"6BVazgHikBwjeqVP36/Stu9eetW9eagEMQhz78caye8sFSHhhl0yp4DAJj1d/29hTELQN8zuq8" +
	"sCAwEAAaNBMD8wDgYDVR0PAQH/BAQDAgbAMAwGA1UdEwEB/wQCMAAwHwYDVR0jBBgwFoAUAf72" +
	"tIMHPtHTlemWdkMVlPxblo0wDQYJKoZIhvcNAQELBQADggEBAA95qJmk1nTd8tfWIyup6bCIdi" +
	"H2tIlfvcKVY3QG7N9TU7QBAkMBK4Yv+MNv6aIgBVoBLa+2OYaegJ5KQLmxL7g1yO0V755zuYYB" +
	"21+X44a3W7GOmmzyMq8sCKXVVVoXhe/15W5ZYQKJALDqXmWqpc/+m97mR1IW0VpwhaBBalI7Di" +
	"XVsV6KHHSLB4GmvcBn1rvqwcroW09J5ig0ryXyLYQ1rY7rUB/74fKJO8k4tK9wDZMni8twyfJ7" +
	"LgXp/kkagCA1EWHih3f7t2zigjryMmhpv6hQAGKbnZwdV2ZulNkY7bM0WbDFGzUVkg23gy+C7n" +
	"QJx8kmdV5/zgBmC5p9bsk="

// certificateTokenFixture parses one of the fixed test certificates.
func certificateTokenFixture(t *testing.T, base64DER string) *CertificateToken {
	t.Helper()
	der, err := base64.StdEncoding.DecodeString(base64DER)
	if err != nil {
		t.Fatalf("bad fixture: %v", err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("ParseCertificate: %v", err)
	}
	token, err := NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("NewCertificateToken: %v", err)
	}
	return token
}

func TestCertificateTokenIdentifiersAreKnownAnswers(t *testing.T) {
	// The expected strings are the SHA-256 digests of the certificate DER (DSS Id) and of the
	// SubjectPublicKeyInfo concatenated with the DER subject name (entity key), rendered the
	// way Digest#getHexValue does.
	tests := []struct {
		name      string
		base64DER string
		dssID     string
		entityKey string
	}{
		{
			name:      "root",
			base64DER: rootCertificateBase64,
			dssID:     "C-E46CC70B1A54986D0E07A28BA9C99468115F71F305894A280CC08DADC31E4AF3",
			entityKey: "EK-36ABC734FED7E18CB0B70D21B8F4286A6FC8E997E84FBD23DBBB321E5EE17939",
		},
		{
			name:      "leaf",
			base64DER: leafCertificateBase64,
			dssID:     "C-2E77AC385070F5D3C9AE8DEB5A3496C3076BFD870342AC8692B02330881EE9B8",
			entityKey: "EK-56A20D7BF3E244228E90E7694BE88196016BCC702E3A8ACFAA22F65BCD4638D5",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token := certificateTokenFixture(t, tc.base64DER)
			if got := token.DSSIDAsString(); got != tc.dssID {
				t.Errorf("DSSIDAsString() = %q, want %q", got, tc.dssID)
			}
			if got := token.Abbreviation(); got != tc.dssID {
				t.Errorf("Abbreviation() = %q, want %q", got, tc.dssID)
			}
			if got := token.EntityKey().AsXmlID(); got != tc.entityKey {
				t.Errorf("EntityKey() = %q, want %q", got, tc.entityKey)
			}
			// The identifier is the digest of the certificate's own DER, unchanged.
			if len(token.Encoded()) == 0 {
				t.Fatal("Encoded() must return the certificate DER")
			}
			if got := NewCertificateTokenIdentifier(token).AsXmlID(); got != tc.dssID {
				t.Errorf("CertificateTokenIdentifier = %q, want %q", got, tc.dssID)
			}
		})
	}
}

func TestCertificateTokenNames(t *testing.T) {
	root := certificateTokenFixture(t, rootCertificateBase64)
	leaf := certificateTokenFixture(t, leafCertificateBase64)

	// Captured from a JDK 21 X500Principal; note the NFKD-decomposed "o" + U+0308 in the
	// canonical form, which a naive lower-casing would not produce.
	const rootSubjectCanonical = "cn=dss go port root ca,o=nowina sölutions,c=lu"
	const rootSubjectRFC2253 = "CN=DSS Go Port Root CA,O=Nowina Sölutions,C=LU"
	const leafSubjectCanonical = "cn=dss go port signer,o=nowina sölutions,c=lu"

	if got := root.Subject().Canonical(); got != rootSubjectCanonical {
		t.Errorf("root Subject().Canonical() = %q, want %q", got, rootSubjectCanonical)
	}
	if got := root.Subject().RFC2253(); got != rootSubjectRFC2253 {
		t.Errorf("root Subject().RFC2253() = %q, want %q", got, rootSubjectRFC2253)
	}
	if got := root.Issuer().Canonical(); got != rootSubjectCanonical {
		t.Errorf("root Issuer().Canonical() = %q, want %q", got, rootSubjectCanonical)
	}
	if got := leaf.Subject().Canonical(); got != leafSubjectCanonical {
		t.Errorf("leaf Subject().Canonical() = %q, want %q", got, leafSubjectCanonical)
	}
	if got := leaf.Issuer().Canonical(); got != rootSubjectCanonical {
		t.Errorf("leaf Issuer().Canonical() = %q, want %q", got, rootSubjectCanonical)
	}
	// getIssuerX500Principal() must return the issuer, keeping its DER byte-for-byte.
	if got := leaf.IssuerX500Principal(); string(got.Encoded()) != string(leaf.Certificate().RawIssuer) {
		t.Error("IssuerX500Principal() must keep the original issuer DER")
	}
	pretty, err := root.Subject().PrettyPrintRFC2253()
	if err != nil {
		t.Fatalf("PrettyPrintRFC2253: %v", err)
	}
	if !strings.HasPrefix(pretty, "commonName=DSS Go Port Root CA,") {
		t.Errorf("PrettyPrintRFC2253() = %q", pretty)
	}
}

func TestCertificateTokenSelfSignedAndSelfIssued(t *testing.T) {
	root := certificateTokenFixture(t, rootCertificateBase64)
	leaf := certificateTokenFixture(t, leafCertificateBase64)

	if !root.IsSelfIssued() {
		t.Error("the root certificate is self-issued")
	}
	if !root.IsSelfSigned() {
		t.Error("the root certificate is self-signed")
	}
	// isSelfSigned() marks the token's signature as VALID.
	if root.SignatureValidity() != enumerations.SignatureValidityValid {
		t.Errorf("root SignatureValidity() = %v", root.SignatureValidity())
	}
	if !root.IsSignatureIntact() || !root.IsValid() {
		t.Error("the root signature should be intact")
	}
	// The cached branch must give the same answer.
	if !root.IsSelfSigned() {
		t.Error("IsSelfSigned() is not stable")
	}

	if leaf.IsSelfIssued() {
		t.Error("the leaf certificate is not self-issued")
	}
	if leaf.IsSelfSigned() {
		t.Error("the leaf certificate is not self-signed")
	}
	if leaf.SignatureValidity() != enumerations.SignatureValidityNotEvaluated {
		t.Errorf("leaf SignatureValidity() = %v before any check", leaf.SignatureValidity())
	}
}

func TestCertificateTokenIsSignedBy(t *testing.T) {
	root := certificateTokenFixture(t, rootCertificateBase64)
	leaf := certificateTokenFixture(t, leafCertificateBase64)

	if leaf.IssuerEntityKey() != nil {
		t.Error("IssuerEntityKey() must be nil before the signer is established")
	}
	if !leaf.IsSignedByToken(root) {
		t.Fatalf("the leaf must be signed by the root, reason: %q", leaf.InvalidityReason())
	}
	if leaf.SignatureValidity() != enumerations.SignatureValidityValid {
		t.Errorf("leaf SignatureValidity() = %v", leaf.SignatureValidity())
	}
	if leaf.InvalidityReason() != "" {
		t.Errorf("InvalidityReason() = %q, want empty", leaf.InvalidityReason())
	}
	if leaf.PublicKeyOfTheSigner() == nil {
		t.Fatal("the signer's public key must be remembered")
	}
	// The issuer entity key is the root's entity key: the same public key and the same name.
	if got, want := leaf.IssuerEntityKey().AsXmlID(), root.EntityKey().AsXmlID(); got != want {
		t.Errorf("IssuerEntityKey() = %q, want %q", got, want)
	}
	// Once the signer is known the short-circuit branch answers from the cached key.
	if !leaf.IsSignedBy(root.PublicKey()) {
		t.Error("the cached signer must still match")
	}
	if leaf.IsSignedBy(leaf.PublicKey()) {
		t.Error("a different key must not match the cached signer")
	}

	// A self-signed certificate does not record itself as the signer.
	if !root.IsSignedByToken(root) {
		t.Error("the root must be signed by itself")
	}
	if root.PublicKeyOfTheSigner() != nil {
		t.Error("a self-signed token must not record the signer's public key")
	}
	// The self-signed override of getIssuerEntityKey() rebuilds the key from the subject.
	if got, want := root.IssuerEntityKey().AsXmlID(), root.EntityKey().AsXmlID(); got != want {
		t.Errorf("root IssuerEntityKey() = %q, want %q", got, want)
	}
}

func TestCertificateTokenCheckIsSignedByRecordsTheFailure(t *testing.T) {
	root := certificateTokenFixture(t, rootCertificateBase64)
	leaf := certificateTokenFixture(t, leafCertificateBase64)

	if got := leaf.CheckIsSignedBy(leaf.PublicKey()); got != enumerations.SignatureValidityInvalid {
		t.Errorf("CheckIsSignedBy(wrong key) = %v", got)
	}
	if leaf.InvalidityReason() == "" {
		t.Error("a failed check must record an invalidity reason")
	}
	if got := leaf.CheckIsSignedBy(root.PublicKey()); got != enumerations.SignatureValidityValid {
		t.Errorf("CheckIsSignedBy(root) = %v", got)
	}
	if leaf.InvalidityReason() != "" {
		t.Errorf("a successful check must clear the invalidity reason, got %q", leaf.InvalidityReason())
	}
}

func TestCertificateTokenCertificateProperties(t *testing.T) {
	root := certificateTokenFixture(t, rootCertificateBase64)
	leaf := certificateTokenFixture(t, leafCertificateBase64)

	if got := root.SerialNumber().String(); got != "4919" {
		t.Errorf("root SerialNumber() = %s", got)
	}
	if got := leaf.SerialNumber().String(); got != "305419896" {
		t.Errorf("leaf SerialNumber() = %s", got)
	}
	if got := root.SignatureAlgorithm(); got != enumerations.SignatureAlgorithmRSASHA256 {
		t.Errorf("SignatureAlgorithm() = %v", got)
	}
	if !root.IsCA() {
		t.Error("the root is a CA")
	}
	if got := root.PathLenConstraint(); got != 2 {
		t.Errorf("root PathLenConstraint() = %d, want 2", got)
	}
	if leaf.IsCA() {
		t.Error("the leaf is not a CA")
	}
	if got := leaf.PathLenConstraint(); got != -1 {
		t.Errorf("leaf PathLenConstraint() = %d, want -1", got)
	}
	if got := len(root.Signature()); got == 0 {
		t.Error("Signature() must return the certificate signature")
	}
	if root.Certificate() == nil {
		t.Error("Certificate() must return the wrapped certificate")
	}

	// getSourceURL/setSourceURL are a plain pair of accessors.
	if root.SourceURL() != "" {
		t.Error("SourceURL() starts empty")
	}
	root.SetSourceURL("http://example.org/ca.crt")
	if root.SourceURL() != "http://example.org/ca.crt" {
		t.Errorf("SourceURL() = %q", root.SourceURL())
	}

	// isEquivalent compares the encoded public keys.
	if !root.IsEquivalent(root) {
		t.Error("a certificate is equivalent to itself")
	}
	if root.IsEquivalent(leaf) {
		t.Error("certificates with different keys are not equivalent")
	}
}

func TestCertificateTokenKeyUsageBits(t *testing.T) {
	root := certificateTokenFixture(t, rootCertificateBase64)
	leaf := certificateTokenFixture(t, leafCertificateBase64)

	want := []enumerations.KeyUsageBit{
		enumerations.KeyUsageBitDigitalSignature,
		enumerations.KeyUsageBitKeyCertSign,
		enumerations.KeyUsageBitCRLSign,
	}
	got := root.KeyUsageBits()
	if len(got) != len(want) {
		t.Fatalf("root KeyUsageBits() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("root KeyUsageBits()[%d] = %v, want %v", i, got[i], want[i])
		}
	}
	if !root.CheckKeyUsage(enumerations.KeyUsageBitKeyCertSign) {
		t.Error("the root must carry keyCertSign")
	}
	if root.CheckKeyUsage(enumerations.KeyUsageBitKeyEncipherment) {
		t.Error("the root must not carry keyEncipherment")
	}

	leafWant := []enumerations.KeyUsageBit{
		enumerations.KeyUsageBitDigitalSignature,
		enumerations.KeyUsageBitNonRepudiation,
	}
	leafGot := leaf.KeyUsageBits()
	if len(leafGot) != len(leafWant) {
		t.Fatalf("leaf KeyUsageBits() = %v, want %v", leafGot, leafWant)
	}
	for i := range leafWant {
		if leafGot[i] != leafWant[i] {
			t.Errorf("leaf KeyUsageBits()[%d] = %v, want %v", i, leafGot[i], leafWant[i])
		}
	}
}

func TestCertificateTokenValidity(t *testing.T) {
	root := certificateTokenFixture(t, rootCertificateBase64)
	notBefore := time.Date(2021, 3, 4, 5, 6, 7, 0, time.UTC)
	notAfter := time.Date(2031, 3, 4, 5, 6, 7, 0, time.UTC)

	if !root.NotBefore().Equal(notBefore) {
		t.Errorf("NotBefore() = %v", root.NotBefore())
	}
	if !root.NotAfter().Equal(notAfter) {
		t.Errorf("NotAfter() = %v", root.NotAfter())
	}
	if !root.CreationDate().Equal(notBefore) {
		t.Errorf("CreationDate() must be notBefore, got %v", root.CreationDate())
	}
	// Both bounds are inclusive, as java.security.cert.X509Certificate#checkValidity is.
	for _, at := range []time.Time{notBefore, notAfter, notBefore.Add(time.Hour)} {
		if !root.IsValidOn(at) {
			t.Errorf("IsValidOn(%v) = false", at)
		}
	}
	for _, at := range []time.Time{notBefore.Add(-time.Nanosecond), notAfter.Add(time.Nanosecond)} {
		if root.IsValidOn(at) {
			t.Errorf("IsValidOn(%v) = true", at)
		}
	}
	// The zero time.Time stands in for Java's null date.
	if root.IsValidOn(time.Time{}) {
		t.Error("IsValidOn(zero time) must be false")
	}
}

func TestCertificateTokenToString(t *testing.T) {
	root := certificateTokenFixture(t, rootCertificateBase64)
	out := root.ToString("")

	for _, want := range []string{
		"CertificateToken[\n",
		"\tDSS Id              : C-E46CC70B1A54986D0E07A28BA9C99468115F71F305894A280CC08DADC31E4AF3\n",
		"\tIdentity Id         : EntityIdentifier:SHA256:#36ABC734FED7E18CB0B70D21B8F4286A6FC8E997E84FBD23DBBB321E5EE17939\n",
		"\tValidity period     : Thu Mar 04 05:06:07 UTC 2021 - Tue Mar 04 05:06:07 UTC 2031\n",
		"\tSubject name        : cn=dss go port root ca,o=nowina sölutions,c=lu\n",
		"\tIssuer subject name : cn=dss go port root ca,o=nowina sölutions,c=lu\n",
		"\tSerial Number       : 4919\n",
		"\tSignature algorithm : RSA_SHA256\n",
		"\t[SELF-SIGNED]\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("ToString() is missing %q\ngot:\n%s", want, out)
		}
	}
	if !strings.HasSuffix(out, "\n]") {
		t.Errorf("ToString() must close the block, got:\n%s", out)
	}
	if root.String() != out {
		t.Error("String() must be ToString(\"\")")
	}
	// Java indents the body one tab deeper and then drops the first character of the
	// indentation before closing the block.
	indented := root.ToString("  ")
	if !strings.HasPrefix(indented, "  CertificateToken[\n") {
		t.Errorf("ToString(indent) = %q...", indented[:40])
	}
	if !strings.HasSuffix(indented, "\n \t]") {
		t.Errorf("ToString(indent) must close with the substring(1) indentation, got %q", indented[len(indented)-8:])
	}
}

func TestCertificateTokenEqualsAndDigest(t *testing.T) {
	root := certificateTokenFixture(t, rootCertificateBase64)
	sameRoot := certificateTokenFixture(t, rootCertificateBase64)
	leaf := certificateTokenFixture(t, leafCertificateBase64)

	if !root.Equals(sameRoot) {
		t.Error("two tokens over the same certificate must be equal")
	}
	if root.Equals(leaf) {
		t.Error("tokens over different certificates must not be equal")
	}
	if root.Equals(nil) {
		t.Error("a token must not equal nil")
	}

	// getDigest delegates to the identifier, so SHA-256 comes from the cache and the other
	// algorithms are computed from the certificate binaries.
	sha256Digest, err := root.Digest(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatal(err)
	}
	if len(sha256Digest) != 32 {
		t.Errorf("SHA-256 digest length = %d", len(sha256Digest))
	}
	sha512Digest, err := root.Digest(enumerations.DigestAlgorithmSHA512)
	if err != nil {
		t.Fatal(err)
	}
	if len(sha512Digest) != 64 {
		t.Errorf("SHA-512 digest length = %d", len(sha512Digest))
	}
	if _, err := root.Digest(enumerations.DigestAlgorithmWHIRLPOOL); err == nil {
		t.Error("an algorithm without a Go implementation must be reported")
	}
}

func TestCertificateTokenPanicsOnMissingCertificate(t *testing.T) {
	defer func() {
		if r := recover(); r != "X509 certificate is missing" {
			t.Errorf("recover() = %v, want the Java requireNonNull message", r)
		}
	}()
	_, _ = NewCertificateToken(nil)
	t.Error("NewCertificateToken(nil) must panic")
}

func TestCertificateTokenSigAlgOIDAndParams(t *testing.T) {
	root := certificateTokenFixture(t, rootCertificateBase64)
	oid, params, err := certificateTokenSigAlgOIDAndParams(root.Encoded())
	if err != nil {
		t.Fatal(err)
	}
	if oid != "1.2.840.113549.1.1.11" {
		t.Errorf("signature algorithm OID = %s", oid)
	}
	// sun.security.x509.AlgorithmId normalizes an explicit ASN.1 NULL to no parameters, and
	// sha256WithRSAEncryption always carries one.
	if params != nil {
		t.Errorf("an explicit NULL must yield nil parameters, got %x", params)
	}
}

func TestCertificateTokenPathLenConstraintForACAWithoutPathLen(t *testing.T) {
	// A CA certificate whose BasicConstraints carries no pathLenConstraint reports
	// Integer.MAX_VALUE, exactly like X509Certificate#getBasicConstraints().
	root := certificateTokenFixture(t, rootCertificateBase64)
	certificate := *root.Certificate()
	certificate.MaxPathLen = -1
	certificate.MaxPathLenZero = false
	token, err := NewCertificateToken(&certificate)
	if err != nil {
		t.Fatal(err)
	}
	if got := token.PathLenConstraint(); got != math.MaxInt32 {
		t.Errorf("PathLenConstraint() = %d, want %d", got, math.MaxInt32)
	}
}
