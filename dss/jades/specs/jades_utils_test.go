// Ported from specs-jades/src/test/java/eu/europa/esig/jades/JAdESUtilsTest.java (DSS 6.5.RC1),
// as vectors (testdata/, copied from specs-jades/src/test/resources) plus the assertions
// rewritten in Go. Message parity with the everit/jsonsKema Java library is NOT a contract for
// this port; only presence/absence of an error - identified by a recognizable
// substring naming the offending property, exactly as assertErrorFound does upstream - is
// verified.
package specs

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/internal/corpustest"
)

func readTestdata(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(corpustest.Path(t, name))
	if err != nil {
		t.Fatalf("reading testdata/%s: %v", name, err)
	}
	return string(data)
}

func assertErrorFound(t *testing.T, errs []string, substring string) {
	t.Helper()
	for _, e := range errs {
		if strings.Contains(e, substring) {
			return
		}
	}
	t.Errorf("expected an error mentioning %q, got: %v", substring, errs)
}

func assertNoErrors(t *testing.T, errs []string) {
	t.Helper()
	if len(errs) != 0 {
		t.Errorf("expected no validation errors, got: %v", errs)
	}
}

func decodeBase64URLHeader(t *testing.T, protectedBase64 string) string {
	t.Helper()
	// The protected header base64url encoding is unpadded (RFC 7515 section 2), same as the
	// rest of this port's JOSE base64url handling.
	decoded, err := base64.RawURLEncoding.DecodeString(protectedBase64)
	if err != nil {
		// tolerate padded input too, matching java.util.Base64's more permissive decoder.
		decoded, err = base64.URLEncoding.DecodeString(protectedBase64)
		if err != nil {
			t.Fatalf("decoding protected header: %v", err)
		}
	}
	return string(decoded)
}

// jwsView is a minimal typed view over a parsed JWS JSON document, enough to reach into
// "protected"/"header"/"signatures" the way the Java test's JsonObjectWrapper does.
type jwsView map[string]any

func parseJWS(t *testing.T, jsonText string) jwsView {
	t.Helper()
	var v jwsView
	if err := json.Unmarshal([]byte(jsonText), &v); err != nil {
		t.Fatalf("parsing JWS JSON: %v", err)
	}
	return v
}

func (v jwsView) protectedString(t *testing.T) string {
	t.Helper()
	s, ok := v["protected"].(string)
	if !ok {
		t.Fatalf("missing string 'protected' field in %v", v)
	}
	return s
}

func (v jwsView) headerJSON(t *testing.T) string {
	t.Helper()
	header, ok := v["header"]
	if !ok {
		t.Fatalf("missing 'header' field in %v", v)
	}
	data, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("re-marshalling header: %v", err)
	}
	return string(data)
}

func (v jwsView) signatures(t *testing.T) []jwsView {
	t.Helper()
	raw, ok := v["signatures"].([]any)
	if !ok {
		t.Fatalf("missing 'signatures' array in %v", v)
	}
	out := make([]jwsView, 0, len(raw))
	for _, s := range raw {
		m, ok := s.(map[string]any)
		if !ok {
			t.Fatalf("signatures entry is not an object: %v", s)
		}
		out = append(out, jwsView(m))
	}
	return out
}

// validateSignature ports the private #validateSignature(JsonObjectWrapper): the protected
// header must validate cleanly against the JAdES protected header schema, and the unprotected
// header must validate cleanly against the JAdES unprotected header schema.
func validateSignature(t *testing.T, signature jwsView) {
	t.Helper()

	protectedString := decodeBase64URLHeader(t, signature.protectedString(t))
	errs := JAdESProtectedHeaderUtilsInstance().ValidateAgainstSchema(protectedString)
	assertNoErrors(t, errs)

	headerJSON := signature.headerJSON(t)
	errs = JAdESUnprotectedHeaderUtilsInstance().ValidateAgainstSchema(headerJSON)
	assertNoErrors(t, errs)
}

func TestJAdESUtils_JSONFlattened(t *testing.T) {
	jsonText := readTestdata(t, "jades-lta.json")
	errs := JAdESUtilsInstance().ValidateAgainstSchema(jsonText)
	assertNoErrors(t, errs)

	validateSignature(t, parseJWS(t, jsonText))
}

func TestJAdESUtils_JSONFlattenedInvalid(t *testing.T) {
	jsonText := readTestdata(t, "jades-lta-invalid.json")
	jws := parseJWS(t, jsonText)

	errs := JAdESUtilsInstance().ValidateAgainstSchema(jsonText)
	assertErrorFound(t, errs, "evilPayload")

	protectedString := decodeBase64URLHeader(t, jws.protectedString(t))
	errs = JAdESProtectedHeaderUtilsInstance().ValidateAgainstSchema(protectedString)
	assertErrorFound(t, errs, "x5t")

	headerJSON := jws.headerJSON(t)
	errs = JAdESUnprotectedHeaderUtilsInstance().ValidateAgainstSchema(headerJSON)
	assertErrorFound(t, errs, "x509Cert")
}

func TestJAdESUtils_JSONSerialization(t *testing.T) {
	jsonText := readTestdata(t, "jades-with-sigPSt.json")
	errs := JAdESUtilsInstance().ValidateAgainstSchema(jsonText)
	assertNoErrors(t, errs)

	jws := parseJWS(t, jsonText)
	signatures := jws.signatures(t)
	if len(signatures) != 1 {
		t.Fatalf("expected 1 signature, got %d", len(signatures))
	}
	validateSignature(t, signatures[0])
}

func TestJAdESUtils_JSONSerializationInvalid(t *testing.T) {
	jsonText := readTestdata(t, "jades-with-sigPSt-invalid.json")
	jws := parseJWS(t, jsonText)

	errs := JAdESUtilsInstance().ValidateAgainstSchema(jsonText)
	assertErrorFound(t, errs, "signature")

	signatures := jws.signatures(t)
	if len(signatures) != 1 {
		t.Fatalf("expected 1 signature, got %d", len(signatures))
	}
	signature := signatures[0]

	protectedString := decodeBase64URLHeader(t, signature.protectedString(t))
	errs = JAdESProtectedHeaderUtilsInstance().ValidateAgainstSchema(protectedString)
	assertErrorFound(t, errs, "hashAV")

	headerJSON := signature.headerJSON(t)
	errs = JAdESUnprotectedHeaderUtilsInstance().ValidateAgainstSchema(headerJSON)
	assertErrorFound(t, errs, "tstokens")
	assertErrorFound(t, errs, "sigPSt")
}
