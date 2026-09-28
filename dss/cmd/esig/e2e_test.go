// End-to-end coverage: every main path of the CLI, driven through run() the
// same way main() drives it, against the small fixtures the dss module
// itself ships (dss/testdata) - no corpus, no network. "tl refresh" needs a
// real network fetch of the live EU List of Trusted Lists, which this suite
// deliberately does not depend on (a test suite that reaches the real
// internet is flaky by construction); it is covered here only for its own
// argument parsing (see main_test.go) and its cache read/write format
// (TestTLCacheRoundTrip and TestTLCacheUsedByValidate below). Actually
// refreshing from ec.europa.eu is exercised by hand.
package main

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/spi"
	"golang.org/x/crypto/pkcs12"
)

// fixture resolves a path under the dss module's own testdata/ directory -
// the small, in-module fixtures also used by the facade's own tests
// (dss/facade_test.go), two directories up from dss/cmd/esig.
func fixture(name string) string {
	return filepath.Join("..", "..", "testdata", name)
}

// testTSAServer starts an httptest.Server that speaks just enough RFC 3161
// to answer this client's tsa.go: it decodes the incoming TimeStampReq with
// this package's own tsaRequest type and issues a genuine token, signed with
// the tsa_ec.p12 fixture's key, over the requested message imprint and
// echoing the request's nonce - which tsa.go insists on, and which
// spi/validation.KeyEntityTSPSource (upstream never sets a nonce) cannot
// produce. It is a test double for a real TSA, not for the library.
func testTSAServer(t *testing.T) *httptest.Server {
	t.Helper()
	tsa := newTestTSA(t)
	return testTSAServerWith(t, func(req tsaRequest) []byte {
		return tsa.token(t, req.MessageImprint.HashAlgorithm.Algorithm, req.MessageImprint.HashedMessage, req.Nonce)
	})
}

// testTSAServerWith is testTSAServer with the token issuance supplied by the
// caller, so a test can answer with a token that does not match the request.
func testTSAServerWith(t *testing.T, issue func(req tsaRequest) []byte) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var req tsaRequest
		if _, err := asn1.Unmarshal(body, &req); err != nil {
			http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
			return
		}
		respBytes, err := asn1.Marshal(tsaResponse{
			Status: tsaPKIStatusInfo{Status: pkiStatusGranted},
			Token:  asn1.RawValue{FullBytes: issue(req)},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(respBytes)
	}))
	t.Cleanup(server.Close)
	return server
}

// testTSA issues RFC 3161 time-stamp tokens from the tsa_ec.p12 fixture.
type testTSA struct {
	key   *ecdsa.PrivateKey
	chain []*x509.Certificate
}

func newTestTSA(t *testing.T) *testTSA {
	t.Helper()
	store, err := os.ReadFile(fixture("tsa_ec.p12"))
	if err != nil {
		t.Fatal(err)
	}
	blocks, err := pkcs12.ToPEM(store, "testpassword")
	if err != nil {
		t.Fatalf("tsa_ec.p12: %v", err)
	}
	tsa := &testTSA{}
	for _, block := range blocks {
		switch block.Type {
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				t.Fatal(err)
			}
			tsa.chain = append(tsa.chain, cert)
		case "PRIVATE KEY":
			if tsa.key, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
				t.Fatal(err)
			}
		}
	}
	if tsa.key == nil || len(tsa.chain) == 0 {
		t.Fatal("tsa_ec.p12: no key entry")
	}
	return tsa
}

// token issues a TimeStampToken over the given message imprint, carrying
// nonce in its TSTInfo when nonce is not nil.
func (tsa *testTSA) token(t *testing.T, hashAlgorithm asn1.ObjectIdentifier, digest []byte, nonce *big.Int) []byte {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Second)
	tstInfo, err := asn1.Marshal(struct {
		Version        int
		Policy         asn1.ObjectIdentifier
		MessageImprint tsaMessageImprint
		SerialNumber   *big.Int
		GenTime        time.Time `asn1:"generalized"`
		Nonce          *big.Int  `asn1:"optional"`
	}{1, asn1.ObjectIdentifier{1, 2, 3, 4, 5, 6, 7, 8, 9},
		tsaMessageImprint{tsaAlgorithmIdentifier{hashAlgorithm}, digest},
		big.NewInt(now.UnixNano()), now, nonce})
	if err != nil {
		t.Fatal(err)
	}
	mustMarshal := func(v any) []byte {
		der, err := asn1.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		return der
	}
	tstDigest := sha256.Sum256(tstInfo)
	certHash := sha256.Sum256(tsa.chain[0].Raw)
	type essCertIDv2 struct{ CertHash []byte }
	type signingCertificateV2 struct{ Certs []essCertIDv2 }
	sha256OID := asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}
	builder := &cmscore.SignerInfoBuilder{
		SID:             cmscore.NewIssuerAndSerialNumberSID(tsa.chain[0].RawIssuer, tsa.chain[0].SerialNumber),
		DigestAlgorithm: asn1ber.NewAlgorithmIdentifier(sha256OID),
		SignedAttributes: cmscore.Attributes{
			cmscore.NewAttribute(cmscore.OIDContentType, mustMarshal(cmscore.OIDCTTSTInfo)),
			cmscore.NewAttribute(cmscore.OIDSigningTime, mustMarshal(now)),
			cmscore.NewAttribute(spi.OIDIdAaSigningCertificateV2,
				mustMarshal(signingCertificateV2{[]essCertIDv2{{certHash[:]}}})),
			cmscore.NewAttribute(cmscore.OIDMessageDigest, mustMarshal(tstDigest[:])),
		},
		SignatureAlgorithm: asn1ber.NewAlgorithmIdentifier(asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 2}),
	}
	toBeSigned := sha256.Sum256(builder.SignedAttributesDER())
	if builder.Signature, err = tsa.key.Sign(rand.Reader, toBeSigned[:], crypto.SHA256); err != nil {
		t.Fatal(err)
	}
	signerInfo, err := builder.Build()
	if err != nil {
		t.Fatal(err)
	}
	certificates := make([]cmscore.CertificateChoice, 0, len(tsa.chain))
	for _, cert := range tsa.chain {
		certificates = append(certificates, cmscore.NewCertificateChoice(cert.Raw))
	}
	cms, err := (&cmscore.SignedDataBuilder{
		EncapContentInfo: cmscore.NewEncapsulatedContentInfo(cmscore.OIDCTTSTInfo, tstInfo),
		Certificates:     certificates,
		SignerInfos:      []*cmscore.SignerInfo{signerInfo},
	}).BuildCMS()
	if err != nil {
		t.Fatal(err)
	}
	return cms.DEREncoded()
}

// TestHTTPTSPSourceRoundTrip exercises the client half of tsa.go directly
// against testTSAServer, independent of "sign"/"extend": a genuine RFC 3161
// request goes out, a genuine token comes back, and it decodes into
// something [dss.Validate] recognises as a valid time-stamp once the TSA's
// own certificate chain is trusted.
func TestHTTPTSPSourceRoundTrip(t *testing.T) {
	server := testTSAServer(t)
	client := newHTTPTSPSource(server.URL)

	digest := make([]byte, 32)
	token, err := client.TimeStampResponse(dss.DigestSHA256, digest)
	if err != nil {
		t.Fatalf("TimeStampResponse: %v", err)
	}
	if len(token.Bytes()) == 0 {
		t.Fatal("TimeStampResponse returned an empty token")
	}
}

func TestHTTPTSPSourceRejectsNon200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "no", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := newHTTPTSPSource(server.URL)
	if _, err := client.TimeStampResponse(dss.DigestSHA256, make([]byte, 32)); err == nil {
		t.Fatal("expected an error for a non-200 response")
	}
}

func TestHTTPTSPSourceRejectsRefusal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respBytes, err := asn1.Marshal(tsaResponse{Status: tsaPKIStatusInfo{Status: 2 /* rejection */}})
		if err != nil {
			t.Fatal(err)
		}
		w.Write(respBytes)
	}))
	defer server.Close()

	client := newHTTPTSPSource(server.URL)
	if _, err := client.TimeStampResponse(dss.DigestSHA256, make([]byte, 32)); err == nil {
		t.Fatal("expected an error for a PKIStatus rejection")
	}
}

// TestHTTPTSPSourceRejectsMismatchedToken pins the request/response matching
// tsa.go does (BouncyCastle's TimeStampResponse.validate for upstream's
// OnlineTSPSource): a granted token over another digest, with another nonce
// or with no nonce at all is refused rather than embedded in a signature.
func TestHTTPTSPSourceRejectsMismatchedToken(t *testing.T) {
	tsa := newTestTSA(t)
	cases := []struct {
		name  string
		issue func(req tsaRequest) []byte
		want  string
	}{
		{"other digest", func(req tsaRequest) []byte {
			return tsa.token(t, req.MessageImprint.HashAlgorithm.Algorithm, make([]byte, 32), req.Nonce)
		}, "message imprint digest"},
		{"other algorithm", func(req tsaRequest) []byte {
			return tsa.token(t, asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 3}, req.MessageImprint.HashedMessage, req.Nonce)
		}, "message imprint algorithm"},
		{"wrong nonce", func(req tsaRequest) []byte {
			return tsa.token(t, req.MessageImprint.HashAlgorithm.Algorithm, req.MessageImprint.HashedMessage,
				new(big.Int).Add(req.Nonce, big.NewInt(1)))
		}, "nonce"},
		{"no nonce", func(req tsaRequest) []byte {
			return tsa.token(t, req.MessageImprint.HashAlgorithm.Algorithm, req.MessageImprint.HashedMessage, nil)
		}, "nonce"},
		{"not a token", func(req tsaRequest) []byte {
			der, _ := asn1.Marshal(42)
			return der
		}, "cannot be parsed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := testTSAServerWith(t, tc.issue)
			digest := bytes.Repeat([]byte{0xA5}, 32)
			_, err := newHTTPTSPSource(server.URL).TimeStampResponse(dss.DigestSHA256, digest)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want one mentioning %q", err, tc.want)
			}
		})
	}
}

// TestHTTPTSPSourceRedactsURLPassword pins that a TSA URL carrying HTTP
// Basic credentials (https://user:secret@host/) never has its password
// echoed into an error message, which the CLI prints to stderr.
func TestHTTPTSPSourceRedactsURLPassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "no", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	const secret = "s3cr3t-tsa-password"
	for _, rawURL := range []string{
		strings.Replace(server.URL, "http://", "http://user:"+secret+"@", 1), // non-200 response
		"http://user:" + secret + "@127.0.0.1:1/",                            // connection refused
		"http://user:" + secret + "@[::1",                                    // unparsable URL
	} {
		_, err := newHTTPTSPSource(rawURL).TimeStampResponse(dss.DigestSHA256, make([]byte, 32))
		if err == nil {
			t.Fatalf("%s: expected an error", rawURL)
		}
		if strings.Contains(err.Error(), secret) {
			t.Errorf("error discloses the TSA password: %v", err)
		}
	}
}

// TestParseOID pins the round trip parseOID needs:
// [enumerations.DigestAlgorithm.OID] output back into an
// asn1.ObjectIdentifier.
func TestParseOID(t *testing.T) {
	oid, err := parseOID(dss.DigestSHA256.OID())
	if err != nil {
		t.Fatalf("parseOID: %v", err)
	}
	want := asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}
	if !oid.Equal(want) {
		t.Errorf("oid = %v, want %v", oid, want)
	}
	if _, err := parseOID("not-an-oid"); err == nil {
		t.Error("expected an error for a malformed OID")
	}
}

// TestSignValidateInspectReport signs the module's own PDF fixture at
// PAdES-B, validates the result twice (untrusted, then anchored), inspects
// it, and renders its SimpleReport back through "esig report" - the
// non-network main paths, chained the way a user actually runs them.
func TestSignValidateInspectReport(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ESIG_TEST_P12_PASSWORD", "testpassword")

	signedPath := filepath.Join(dir, "signed.pdf")
	code, stdout, stderr := runOut(t, "sign", fixture("sample.pdf"),
		"-format", "pades", "-level", "B",
		"-p12", fixture("signer_rsa.p12"), "-p12-pass", "env:ESIG_TEST_P12_PASSWORD",
		"-out", signedPath)
	if code != exitOK {
		t.Fatalf("sign: exit code = %d, want %d; stderr:\n%s", code, exitOK, stderr)
	}
	if !strings.Contains(stdout, signedPath) {
		t.Errorf("sign stdout = %q, want it to name %q", stdout, signedPath)
	}
	if _, err := os.Stat(signedPath); err != nil {
		t.Fatalf("signed document was not written: %v", err)
	}

	// Untrusted: no anchor, so the chain cannot be verified - a validation
	// failure by DSS's own definition (exitValidationFailed), not a runtime
	// error.
	code, stdout, _ = runOut(t, "validate", signedPath)
	if code != exitValidationFailed {
		t.Errorf("validate (untrusted): exit code = %d, want %d (exitValidationFailed)", code, exitValidationFailed)
	}
	if !strings.Contains(stdout, "INDETERMINATE") {
		t.Errorf("validate (untrusted) stdout = %q, want it to report INDETERMINATE", stdout)
	}

	// Anchored: the signing certificate is its own trust anchor (it is
	// self-signed test fixture), so this must pass.
	code, stdout, stderr = runOut(t, "validate", signedPath, "-trust", fixture("signer_rsa.cer"))
	if code != exitOK {
		t.Fatalf("validate (trusted): exit code = %d, want %d; stdout:\n%s\nstderr:\n%s", code, exitOK, stdout, stderr)
	}
	if !strings.Contains(stdout, "TOTAL_PASSED") {
		t.Errorf("validate (trusted) stdout = %q, want it to report TOTAL_PASSED", stdout)
	}

	// Every report format renders without error.
	for _, format := range []string{"simple", "detailed", "diagnostic", "etsi-vr"} {
		code, stdout, stderr := runOut(t, "validate", signedPath, "-trust", fixture("signer_rsa.cer"), "-format", format)
		if code != exitOK {
			t.Errorf("validate -format %s: exit code = %d, want %d; stderr:\n%s", format, code, exitOK, stderr)
		}
		if len(stdout) == 0 {
			t.Errorf("validate -format %s produced no output", format)
		}
	}

	// The SimpleReport, saved to a file, renders through "esig report".
	simplePath := filepath.Join(dir, "simple.xml")
	code, _, stderr = runOut(t, "validate", signedPath, "-trust", fixture("signer_rsa.cer"), "-format", "simple", "-out", simplePath)
	if code != exitOK {
		t.Fatalf("validate -format simple -out: exit code = %d, want %d; stderr:\n%s", code, exitOK, stderr)
	}
	code, stdout, stderr = runOut(t, "report", simplePath, "-render")
	if code != exitOK {
		t.Fatalf("report -render: exit code = %d, want %d; stderr:\n%s", code, exitOK, stderr)
	}
	if !strings.Contains(stdout, "TOTAL_PASSED") {
		t.Errorf("report -render stdout = %q, want it to report TOTAL_PASSED", stdout)
	}

	// inspect reads the diagnostic data regardless of trust.
	code, stdout, stderr = runOut(t, "inspect", signedPath)
	if code != exitOK {
		t.Fatalf("inspect: exit code = %d, want %d; stderr:\n%s", code, exitOK, stderr)
	}
	if !strings.Contains(stdout, "PAdES-BASELINE-B") {
		t.Errorf("inspect stdout = %q, want it to name the signature format", stdout)
	}
}

// TestSignDetachedCAdES covers -detached on "sign" (a binary payload signed
// detached) and DetachedContents on "validate" (the same payload supplied
// back for the detached signature to be checked against).
func TestSignDetachedCAdES(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ESIG_TEST_P12_PASSWORD", "testpassword")

	payload := filepath.Join(dir, "payload.bin")
	if err := os.WriteFile(payload, []byte("payload"), 0o644); err != nil {
		t.Fatal(err)
	}

	signedPath := filepath.Join(dir, "signed.p7s")
	code, _, stderr := runOut(t, "sign", payload,
		"-format", "cades", "-level", "B", "-detached",
		"-p12", fixture("signer_rsa.p12"), "-p12-pass", "env:ESIG_TEST_P12_PASSWORD",
		"-out", signedPath)
	if code != exitOK {
		t.Fatalf("sign -detached: exit code = %d, want %d; stderr:\n%s", code, exitOK, stderr)
	}

	code, stdout, stderr := runOut(t, "validate", signedPath,
		"-detached", payload, "-trust", fixture("signer_rsa.cer"))
	if code != exitOK {
		t.Fatalf("validate -detached: exit code = %d, want %d; stdout:\n%s\nstderr:\n%s", code, exitOK, stdout, stderr)
	}
	if !strings.Contains(stdout, "TOTAL_PASSED") {
		t.Errorf("validate -detached stdout = %q, want TOTAL_PASSED", stdout)
	}
}

// TestSignExtendLevelT covers -level T on both "sign" (T directly) and
// "extend" (B raised to T), both against testTSAServer, and confirms
// "validate" sees the T level and a best-signature-time on the result.
func TestSignExtendLevelT(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ESIG_TEST_P12_PASSWORD", "testpassword")
	tsa := testTSAServer(t)

	xmlPath := filepath.Join(dir, "invoice.xml")
	if err := os.WriteFile(xmlPath, []byte("<invoice><total>42</total></invoice>"), 0o644); err != nil {
		t.Fatal(err)
	}

	bPath := filepath.Join(dir, "invoice-b.xml")
	code, _, stderr := runOut(t, "sign", xmlPath,
		"-format", "xades", "-level", "B",
		"-p12", fixture("signer_rsa.p12"), "-p12-pass", "env:ESIG_TEST_P12_PASSWORD",
		"-out", bPath)
	if code != exitOK {
		t.Fatalf("sign -level B: exit code = %d, want %d; stderr:\n%s", code, exitOK, stderr)
	}

	tPath := filepath.Join(dir, "invoice-t.xml")
	code, _, stderr = runOut(t, "extend", bPath,
		"-format", "xades", "-level", "T", "-tsa", tsa.URL,
		"-out", tPath)
	if code != exitOK {
		t.Fatalf("extend -level T: exit code = %d, want %d; stderr:\n%s", code, exitOK, stderr)
	}

	code, stdout, stderr := runOut(t, "validate", tPath, "-trust", fixture("signer_rsa.cer"))
	if code != exitOK {
		t.Fatalf("validate: exit code = %d, want %d; stdout:\n%s\nstderr:\n%s", code, exitOK, stdout, stderr)
	}
	if !strings.Contains(stdout, "XAdES-BASELINE-T") {
		t.Errorf("validate stdout = %q, want it to report XAdES-BASELINE-T", stdout)
	}
}

// TestASiCContainer covers -format asice/asics and -asic-format on "sign",
// with more than one document going into an ASiC-E container.
func TestASiCContainer(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ESIG_TEST_P12_PASSWORD", "testpassword")

	docA := filepath.Join(dir, "a.bin")
	docB := filepath.Join(dir, "b.json")
	if err := os.WriteFile(docA, []byte("payload"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(docB, []byte(`{"amount":42}`), 0o644); err != nil {
		t.Fatal(err)
	}

	// "sign" only takes one <file>, so this covers the container path with a
	// single document, defaulting to ASiC-S, and the -asic-format override.
	outPath := filepath.Join(dir, "container.sce")
	code, _, stderr := runOut(t, "sign", docA,
		"-format", "asics", "-level", "B", "-asic-format", "cades",
		"-p12", fixture("signer_rsa.p12"), "-p12-pass", "env:ESIG_TEST_P12_PASSWORD",
		"-out", outPath)
	if code != exitOK {
		t.Fatalf("sign asics -asic-format cades: exit code = %d, want %d; stderr:\n%s", code, exitOK, stderr)
	}

	code, stdout, stderr := runOut(t, "validate", outPath, "-trust", fixture("signer_rsa.cer"))
	if code != exitOK {
		t.Fatalf("validate: exit code = %d, want %d; stdout:\n%s\nstderr:\n%s", code, exitOK, stdout, stderr)
	}
	if !strings.Contains(stdout, "CAdES-BASELINE-B") {
		t.Errorf("validate stdout = %q, want it to report CAdES-BASELINE-B", stdout)
	}
}

// TestPasswordMustBeEnvForm pins that a literal password is refused: only
// env:VARNAME is accepted (see resolvePassword's doc comment for why). It
// also pins that the refusal does not echo the literal back - the value is a
// real password precisely when this path fires, and stderr is logged.
func TestPasswordMustBeEnvForm(t *testing.T) {
	// The literal is deliberately not the keystore's own password, so that
	// finding it in stderr can only mean the flag value was echoed.
	const literal = "s3cret-literal-not-env"
	t.Setenv("ESIG_TEST_P12_PASSWORD", "testpassword")
	common := []string{"sign", fixture("sample.pdf"), "-format", "pades", "-level", "B",
		"-p12", fixture("signer_rsa.p12")}

	cases := map[string][]string{
		// -p12-pass is resolved first, so it fails before -pdf-pass is read.
		"-p12-pass": append(append([]string{}, common...), "-p12-pass", literal),
		"-pdf-pass": append(append([]string{}, common...),
			"-p12-pass", "env:ESIG_TEST_P12_PASSWORD", "-pdf-pass", literal),
	}
	for flag, args := range cases {
		t.Run(flag, func(t *testing.T) {
			code, _, stderr := runOut(t, args...)
			if code != exitRuntime {
				t.Errorf("exit code = %d, want %d (exitRuntime); stderr:\n%s", code, exitRuntime, stderr)
			}
			if !strings.Contains(stderr, "env:") {
				t.Errorf("stderr = %q, want it to mention the env: form", stderr)
			}
			if strings.Contains(stderr, literal) {
				t.Errorf("stderr = %q, want it NOT to echo the password back", stderr)
			}
		})
	}
}

// TestSignProtectedPDF covers -pdf-pass end to end on an encrypted
// (password-protected) PDF - upstream's protected/open_protected.pdf test
// resource, whose user password is a single space: "sign" at B, "validate"
// and "inspect" on the result, "extend" to T against testTSAServer, and
// "validate" again. Every step needs the password; the same steps without it
// must fail as a runtime error (the document cannot be opened), not as a
// verdict. The signed and extended outputs stay encrypted, which is what
// "validate" without the password failing on them demonstrates.
func TestSignProtectedPDF(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ESIG_TEST_P12_PASSWORD", "testpassword")
	t.Setenv("ESIG_TEST_PDF_PASSWORD", " ")
	tsa := testTSAServer(t)

	// Without the password the encrypted PDF cannot be opened for signing.
	code, _, stderr := runOut(t, "sign", fixture("upstream/protected/open_protected.pdf"),
		"-format", "pades", "-level", "B",
		"-p12", fixture("signer_rsa.p12"), "-p12-pass", "env:ESIG_TEST_P12_PASSWORD",
		"-out", filepath.Join(dir, "unused.pdf"))
	if code != exitRuntime {
		t.Errorf("sign without -pdf-pass: exit code = %d, want %d (exitRuntime); stderr:\n%s", code, exitRuntime, stderr)
	}
	if !strings.Contains(stderr, "password") {
		t.Errorf("sign without -pdf-pass stderr = %q, want it to mention the password", stderr)
	}

	// -pdf-pass takes the env: form only, like -p12-pass.
	code, _, stderr = runOut(t, "sign", fixture("upstream/protected/open_protected.pdf"),
		"-format", "pades", "-level", "B",
		"-p12", fixture("signer_rsa.p12"), "-p12-pass", "env:ESIG_TEST_P12_PASSWORD",
		"-pdf-pass", " ", "-out", filepath.Join(dir, "unused.pdf"))
	if code != exitRuntime {
		t.Errorf("sign with literal -pdf-pass: exit code = %d, want %d (exitRuntime); stderr:\n%s", code, exitRuntime, stderr)
	}
	if !strings.Contains(stderr, "env:PDF_PASSWORD") {
		t.Errorf("sign with literal -pdf-pass stderr = %q, want it to suggest the env:PDF_PASSWORD form", stderr)
	}

	bPath := filepath.Join(dir, "protected-b.pdf")
	code, stdout, stderr := runOut(t, "sign", fixture("upstream/protected/open_protected.pdf"),
		"-format", "pades", "-level", "B",
		"-p12", fixture("signer_rsa.p12"), "-p12-pass", "env:ESIG_TEST_P12_PASSWORD",
		"-pdf-pass", "env:ESIG_TEST_PDF_PASSWORD", "-out", bPath)
	if code != exitOK {
		t.Fatalf("sign -pdf-pass: exit code = %d, want %d; stderr:\n%s", code, exitOK, stderr)
	}
	if !strings.Contains(stdout, bPath) {
		t.Errorf("sign stdout = %q, want it to name %q", stdout, bPath)
	}

	// The signed PDF is still encrypted: validating it without the password
	// cannot even open it.
	code, _, stderr = runOut(t, "validate", bPath, "-trust", fixture("signer_rsa.cer"))
	if code != exitRuntime {
		t.Errorf("validate without -pdf-pass: exit code = %d, want %d (exitRuntime); stderr:\n%s", code, exitRuntime, stderr)
	}

	code, stdout, stderr = runOut(t, "validate", bPath,
		"-pdf-pass", "env:ESIG_TEST_PDF_PASSWORD", "-trust", fixture("signer_rsa.cer"))
	if code != exitOK {
		t.Fatalf("validate -pdf-pass: exit code = %d, want %d; stdout:\n%s\nstderr:\n%s", code, exitOK, stdout, stderr)
	}
	if !strings.Contains(stdout, "TOTAL_PASSED") || !strings.Contains(stdout, "PAdES-BASELINE-B") {
		t.Errorf("validate -pdf-pass stdout = %q, want TOTAL_PASSED at PAdES-BASELINE-B", stdout)
	}

	code, stdout, stderr = runOut(t, "inspect", bPath, "-pdf-pass", "env:ESIG_TEST_PDF_PASSWORD")
	if code != exitOK {
		t.Fatalf("inspect -pdf-pass: exit code = %d, want %d; stderr:\n%s", code, exitOK, stderr)
	}
	if !strings.Contains(stdout, "PAdES-BASELINE-B") {
		t.Errorf("inspect -pdf-pass stdout = %q, want it to name the signature format", stdout)
	}

	tPath := filepath.Join(dir, "protected-t.pdf")
	code, _, stderr = runOut(t, "extend", bPath,
		"-format", "pades", "-level", "T", "-tsa", tsa.URL,
		"-out", tPath)
	if code != exitRuntime {
		t.Errorf("extend without -pdf-pass: exit code = %d, want %d (exitRuntime); stderr:\n%s", code, exitRuntime, stderr)
	}
	code, _, stderr = runOut(t, "extend", bPath,
		"-format", "pades", "-level", "T", "-tsa", tsa.URL,
		"-pdf-pass", "env:ESIG_TEST_PDF_PASSWORD", "-out", tPath)
	if code != exitOK {
		t.Fatalf("extend -pdf-pass: exit code = %d, want %d; stderr:\n%s", code, exitOK, stderr)
	}

	code, stdout, stderr = runOut(t, "validate", tPath,
		"-pdf-pass", "env:ESIG_TEST_PDF_PASSWORD", "-trust", fixture("signer_rsa.cer"))
	if code != exitOK {
		t.Fatalf("validate extended -pdf-pass: exit code = %d, want %d; stdout:\n%s\nstderr:\n%s", code, exitOK, stdout, stderr)
	}
	if !strings.Contains(stdout, "TOTAL_PASSED") || !strings.Contains(stdout, "PAdES-BASELINE-T") {
		t.Errorf("validate extended stdout = %q, want TOTAL_PASSED at PAdES-BASELINE-T", stdout)
	}

	// -pdf-pass on a document that is not a PDF is a runtime error from the
	// facade, since "validate" detects the format itself and only a PDF can
	// take a password.
	xmlPath := filepath.Join(dir, "not-a-pdf.xml")
	if err := os.WriteFile(xmlPath, []byte("<invoice><total>42</total></invoice>"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, _, stderr = runOut(t, "validate", xmlPath, "-pdf-pass", "env:ESIG_TEST_PDF_PASSWORD")
	if code != exitRuntime {
		t.Errorf("validate -pdf-pass on XML: exit code = %d, want %d (exitRuntime); stderr:\n%s", code, exitRuntime, stderr)
	}
}

// TestSignProtectedPDFNonASCIIPassword signs and validates an AES-128
// (/V 4 /R 4) PDF that pdfbox itself protected with the user password "café"
// (internal/pdf/testdata/password, see its README): the password reaches the
// CLI as UTF-8 text from the environment and has to be hashed the way pdfbox
// hashes it - ISO-8859-1 for that revision - for the document to open at all.
func TestSignProtectedPDFNonASCIIPassword(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ESIG_TEST_P12_PASSWORD", "testpassword")
	t.Setenv("ESIG_TEST_PDF_PASSWORD", "café")
	protected := filepath.Join("..", "..", "internal", "pdf", "testdata", "password", "aes128_r4_latin1.pdf")

	signedPath := filepath.Join(dir, "latin1-b.pdf")
	code, _, stderr := runOut(t, "sign", protected,
		"-format", "pades", "-level", "B",
		"-p12", fixture("signer_rsa.p12"), "-p12-pass", "env:ESIG_TEST_P12_PASSWORD",
		"-pdf-pass", "env:ESIG_TEST_PDF_PASSWORD", "-out", signedPath)
	if code != exitOK {
		t.Fatalf("sign -pdf-pass (non-ASCII): exit code = %d, want %d; stderr:\n%s", code, exitOK, stderr)
	}
	code, stdout, stderr := runOut(t, "validate", signedPath,
		"-pdf-pass", "env:ESIG_TEST_PDF_PASSWORD", "-trust", fixture("signer_rsa.cer"))
	if code != exitOK {
		t.Fatalf("validate -pdf-pass (non-ASCII): exit code = %d, want %d; stdout:\n%s\nstderr:\n%s", code, exitOK, stdout, stderr)
	}
	if !strings.Contains(stdout, "TOTAL_PASSED") {
		t.Errorf("validate stdout = %q, want TOTAL_PASSED", stdout)
	}
}

// TestTLCacheRoundTrip covers the cache format "tl refresh -cache" writes
// and "validate -tl-cache" reads, without a network fetch: it writes the
// module's own test certificate through writeTLCache and reads it back
// through readTLCache directly, then confirms "validate -tl-cache" picks up
// what was written.
func TestTLCacheRoundTrip(t *testing.T) {
	cert, err := dss.LoadCertificate(fixture("signer_rsa.cer"))
	if err != nil {
		t.Fatalf("LoadCertificate: %v", err)
	}
	dir := t.TempDir()
	if err := writeTLCache(dir, []*dss.CertificateToken{cert}); err != nil {
		t.Fatalf("writeTLCache: %v", err)
	}
	got, err := readTLCache(dir)
	if err != nil {
		t.Fatalf("readTLCache: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("readTLCache: got %d certificates, want 1", len(got))
	}
	if !bytes.Equal(got[0].Encoded(), cert.Encoded()) {
		t.Error("readTLCache did not round-trip the certificate's encoding")
	}

	// readTLCache on a directory that was never refreshed is a clear error,
	// not a silent empty result.
	if _, err := readTLCache(t.TempDir()); err == nil {
		t.Error("expected an error reading an empty cache directory")
	}
}

// TestTLCacheUsedByValidate signs a document with the facade directly (no
// network involved), caches its signer's certificate through writeTLCache
// as -tl-cache would hold it after a refresh, and confirms "validate
// -tl-cache" anchors the chain from it.
func TestTLCacheUsedByValidate(t *testing.T) {
	signer, err := dss.OpenPKCS12(fixture("signer_rsa.p12"), "testpassword")
	if err != nil {
		t.Fatalf("OpenPKCS12: %v", err)
	}
	defer signer.Close()

	document := dss.NewDocument("payload.bin", []byte("payload"))
	signed, err := dss.Sign(document, signer, dss.SignOptions{Format: dss.FormatCAdES, Level: dss.LevelB})
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	dir := t.TempDir()
	signedPath := filepath.Join(dir, "signed.p7m")
	if err := signed.Save(signedPath); err != nil {
		t.Fatal(err)
	}

	cacheDir := t.TempDir()
	if err := writeTLCache(cacheDir, signer.CertificateChain()); err != nil {
		t.Fatalf("writeTLCache: %v", err)
	}

	code, stdout, stderr := runOut(t, "validate", signedPath, "-tl-cache", cacheDir)
	if code != exitOK {
		t.Fatalf("validate -tl-cache: exit code = %d, want %d; stdout:\n%s\nstderr:\n%s", code, exitOK, stdout, stderr)
	}
	if !strings.Contains(stdout, "TOTAL_PASSED") {
		t.Errorf("validate -tl-cache stdout = %q, want TOTAL_PASSED", stdout)
	}
}
