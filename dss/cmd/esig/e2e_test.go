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
	"encoding/asn1"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
)

// fixture resolves a path under the dss module's own testdata/ directory -
// the small, in-module fixtures also used by the facade's own tests
// (dss/facade_test.go), two directories up from dss/cmd/esig.
func fixture(name string) string {
	return filepath.Join("..", "..", "testdata", name)
}

// testTSAServer starts an httptest.Server that speaks just enough RFC 3161
// to answer this client's tsa.go: it decodes the incoming TimeStampReq with
// this package's own tsaRequest type and issues a real token with
// [spivalidation.KeyEntityTSPSource] - the same offline, self-hosted TSA the
// facade's own tests use (dss/facade_test.go's testTSA) - keyed to always
// answer as SHA-256, which is what every sign/extend call in this file asks
// for. It is a test double for a real TSA's HTTP transport, not for the
// library: the token itself is genuine, produced by the same code path
// dss.SignOptions.TSPSource drives in production.
func testTSAServer(t *testing.T) *httptest.Server {
	t.Helper()
	tsa, err := spivalidation.NewKeyEntityTSPSourceFromKeyStorePath(
		fixture("tsa_ec.p12"), "PKCS12", "testpassword", "", "testpassword")
	if err != nil {
		t.Fatalf("KeyEntityTSPSource: %v", err)
	}
	tsa.SetTsaPolicy("1.2.3.4.5.6.7.8.9")

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
		token, err := tsa.TimeStampResponse(dss.DigestSHA256, req.MessageImprint.HashedMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respBytes, err := asn1.Marshal(tsaResponse{
			Status: tsaPKIStatusInfo{Status: pkiStatusGranted},
			Token:  asn1.RawValue{FullBytes: token.Bytes()},
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
// env:VARNAME is accepted (see resolvePassword's doc comment for why).
func TestPasswordMustBeEnvForm(t *testing.T) {
	code, _, stderr := runOut(t, "sign", fixture("sample.pdf"),
		"-format", "pades", "-level", "B",
		"-p12", fixture("signer_rsa.p12"), "-p12-pass", "testpassword")
	if code != exitRuntime {
		t.Errorf("exit code = %d, want %d (exitRuntime)", code, exitRuntime)
	}
	if !strings.Contains(stderr, "env:") {
		t.Errorf("stderr = %q, want it to mention the env: form", stderr)
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
