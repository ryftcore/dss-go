// Tests for this porter's ROOT-chunk additions: JAdESHeaderParameterNames, JAdESSigningTimeType,
// JAdESTimestampParameters, JAdESSignatureParameters, and the HTTPHeader family. Ported test
// *vectors* are not applicable here (these are constant/value-object classes, not parsers), per
// PORTING.md "port test vectors, not JUnit code" - these are exhaustive table tests over the
// registry-like constant sets, plus behavioural tests mirroring what upstream's getters/setters
// guarantee.
package jades

import (
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

func TestJAdESHeaderParameterNames_Values(t *testing.T) {
	// A representative sample, including both duplicate-spelling aliases this ROOT chunk had to
	// add for already-landed call sites (see the file header note in
	// jades_header_parameter_names.go).
	// JAdESHeaderParameterNamesTstVD/…TstVd and …SigPSt/…SigPst are two spellings of the same
	// constant value (see below) and so cannot both appear as keys of this map literal - each
	// pair's aliasing is asserted separately underneath instead.
	cases := map[string]string{
		JAdESHeaderParameterNamesSigT:          "sigT",
		JAdESHeaderParameterNamesX5tO:          "x5t#o",
		JAdESHeaderParameterNamesSigPid:        "sigPId",
		JAdESHeaderParameterNamesOrgantization: "organization",
		JAdESHeaderParameterNamesEtsiU:         "etsiU",
		JAdESHeaderParameterNamesTstVD:         "tstVD",
		JAdESHeaderParameterNamesSigPSt:        "sigPSt",
		JAdESHeaderParameterNamesCSig:          "cSig",
		JAdESHeaderParameterNamesExp:           "exp",
		JAdESHeaderParameterNamesIat:           "iat",
	}
	for constant, want := range cases {
		if constant != want {
			t.Errorf("constant = %q, want %q", constant, want)
		}
	}
	if JAdESHeaderParameterNamesTstVD != JAdESHeaderParameterNamesTstVd {
		t.Error("JAdESHeaderParameterNamesTstVD and JAdESHeaderParameterNamesTstVd must alias the same value")
	}
	if JAdESHeaderParameterNamesSigPSt != JAdESHeaderParameterNamesSigPst {
		t.Error("JAdESHeaderParameterNamesSigPSt and JAdESHeaderParameterNamesSigPst must alias the same value")
	}
	if JAdESHeaderParameterNamesExp != JWTClaimNamesExp {
		t.Error("JAdESHeaderParameterNamesExp must match JWTClaimNamesExp verbatim (RFC 7519 4.1.4)")
	}
	if JAdESHeaderParameterNamesIat != JWTClaimNamesIat {
		t.Error("JAdESHeaderParameterNamesIat must match JWTClaimNamesIat verbatim (RFC 7519 4.1.6)")
	}
}

func TestJAdESSigningTimeType_Values(t *testing.T) {
	cases := map[JAdESSigningTimeType]string{
		JAdESSigningTimeTypeIAT:  "IAT",
		JAdESSigningTimeTypeSigT: "SIG_T",
		JAdESSigningTimeTypeNone: "NONE",
	}
	for constant, want := range cases {
		if string(constant) != want {
			t.Errorf("constant = %q, want %q", constant, want)
		}
	}
	values := JAdESSigningTimeTypeValues()
	if len(values) != 3 {
		t.Fatalf("JAdESSigningTimeTypeValues() = %v, want 3 entries", values)
	}
}

func TestJAdESTimestampParameters_Defaults(t *testing.T) {
	p := NewJAdESTimestampParameters()
	if p.CanonicalizationMethod() != "" {
		t.Errorf("CanonicalizationMethod() = %q, want empty", p.CanonicalizationMethod())
	}
	// Java's TimestampParameters field initializer is `DigestAlgorithm digestAlgorithm =
	// DigestAlgorithm.SHA512;` (dss-model TimestampParameters.java), not null - the frozen
	// model.NewTimestampParameters() this delegates to already reproduces that default.
	if p.DigestAlgorithm() != enumerations.DigestAlgorithmSHA512 {
		t.Errorf("DigestAlgorithm() = %q, want SHA512", p.DigestAlgorithm())
	}

	p2 := NewJAdESTimestampParametersWithDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	if p2.DigestAlgorithm() != enumerations.DigestAlgorithmSHA256 {
		t.Errorf("DigestAlgorithm() = %v, want SHA256", p2.DigestAlgorithm())
	}
}

func TestJAdESTimestampParameters_SetCanonicalizationMethodPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("SetCanonicalizationMethod should panic - not supported in the current version")
		}
	}()
	NewJAdESTimestampParameters().SetCanonicalizationMethod("http://www.w3.org/2001/10/xml-exc-c14n#")
}

func TestJAdESTimestampParameters_Equals(t *testing.T) {
	a := NewJAdESTimestampParametersWithDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	b := NewJAdESTimestampParametersWithDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	if !a.Equals(b) {
		t.Error("expected equal JAdESTimestampParameters")
	}
	c := NewJAdESTimestampParametersWithDigestAlgorithm(enumerations.DigestAlgorithmSHA512)
	if a.Equals(c) {
		t.Error("expected unequal JAdESTimestampParameters for different digest algorithms")
	}
}

func TestJAdESSignatureParameters_Defaults(t *testing.T) {
	p := NewJAdESSignatureParameters()
	if !p.IsIncludeCertificateChain() {
		t.Error("IsIncludeCertificateChain() default should be true")
	}
	if !p.IsIncludeSignatureType() {
		t.Error("IsIncludeSignatureType() default should be true")
	}
	if !p.IsIncludeKeyIdentifier() {
		t.Error("IsIncludeKeyIdentifier() default should be true")
	}
	if !p.IsBase64UrlEncodedPayload() {
		t.Error("IsBase64UrlEncodedPayload() default should be true")
	}
	if p.SigningCertificateDigestMethod() != enumerations.DigestAlgorithmSHA512 {
		t.Errorf("SigningCertificateDigestMethod() = %v, want SHA512", p.SigningCertificateDigestMethod())
	}
	if p.JwsSerializationType() != enumerations.JWSSerializationTypeCompactSerialization {
		t.Errorf("JwsSerializationType() = %v, want COMPACT_SERIALIZATION", p.JwsSerializationType())
	}
	if p.JadesSigningTimeType() != JAdESSigningTimeTypeIAT {
		t.Errorf("JadesSigningTimeType() = %v, want IAT", p.JadesSigningTimeType())
	}
	if p.IsBase64UrlEncodedEtsiUComponents() != nil {
		t.Error("IsBase64UrlEncodedEtsiUComponents() default should be nil (Java null)")
	}
	if p.ExpirationTime() != nil {
		t.Error("ExpirationTime() default should be nil")
	}
}

func TestJAdESSignatureParameters_SetSignatureLevelRestrictsToJAdES(t *testing.T) {
	p := NewJAdESSignatureParameters()
	p.SetSignatureLevel(enumerations.SignatureLevelJAdESBaselineB)
	if p.SignatureLevel() != enumerations.SignatureLevelJAdESBaselineB {
		t.Errorf("SignatureLevel() = %v, want JAdES_BASELINE_B", p.SignatureLevel())
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("SetSignatureLevel should panic for a non-JAdES level")
		}
	}()
	p.SetSignatureLevel(enumerations.SignatureLevelXAdESBaselineB)
}

func TestJAdESSignatureParameters_LazyTimestampParameters(t *testing.T) {
	p := NewJAdESSignatureParameters()
	if p.GetContentTimestampParameters() == nil {
		t.Fatal("GetContentTimestampParameters() should lazily instantiate")
	}
	if p.GetSignatureTimestampParameters() == nil {
		t.Fatal("GetSignatureTimestampParameters() should lazily instantiate")
	}
	if p.GetArchiveTimestampParameters() == nil {
		t.Fatal("GetArchiveTimestampParameters() should lazily instantiate")
	}
}

func TestJAdESSignatureParameters_SetBase64UrlEncodedEtsiUComponents(t *testing.T) {
	p := NewJAdESSignatureParameters()
	trueVal := true
	p.SetBase64UrlEncodedEtsiUComponents(&trueVal)
	if p.IsBase64UrlEncodedEtsiUComponents() == nil || !*p.IsBase64UrlEncodedEtsiUComponents() {
		t.Error("expected IsBase64UrlEncodedEtsiUComponents() to report true")
	}
}

func TestJAdESSignatureParameters_ExpirationTime(t *testing.T) {
	p := NewJAdESSignatureParameters()
	now := time.Now()
	p.SetExpirationTime(&now)
	if p.ExpirationTime() == nil || !p.ExpirationTime().Equal(now) {
		t.Error("expected ExpirationTime() to round-trip")
	}
}

func TestHTTPHeader(t *testing.T) {
	h := NewHTTPHeader("Digest", "SHA-256=abc")
	if h.Name() != "Digest" {
		t.Errorf("Name() = %q, want Digest", h.Name())
	}
	if h.Value() != "SHA-256=abc" {
		t.Errorf("Value() = %q, want SHA-256=abc", h.Value())
	}
	h.SetValue("SHA-256=def")
	if h.Value() != "SHA-256=def" {
		t.Errorf("Value() after SetValue = %q, want SHA-256=def", h.Value())
	}
	if h.MimeType() != nil {
		t.Error("MimeType() should be nil (not applicable)")
	}

	var doc model.DSSDocument = h
	if doc.Name() != "Digest" {
		t.Error("HTTPHeader should satisfy model.DSSDocument")
	}
}

func TestHTTPHeader_UnsupportedOperationsPanic(t *testing.T) {
	h := NewHTTPHeader("Digest", "SHA-256=abc")

	panics := map[string]func(){
		"SetName":     func() { h.SetName("x") },
		"OpenStream":  func() { h.OpenStream() }, //nolint:errcheck
		"WriteTo":     func() { h.WriteTo(nil) }, //nolint:errcheck
		"SetMimeType": func() { h.SetMimeType(enumerations.MimeTypeEnumBinary) },
		"Save":        func() { h.Save("/tmp/x") },                                  //nolint:errcheck
		"Digest":      func() { h.Digest(enumerations.DigestAlgorithmSHA256) },      //nolint:errcheck
		"DigestValue": func() { h.DigestValue(enumerations.DigestAlgorithmSHA256) }, //nolint:errcheck
	}
	for name, fn := range panics {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("%s should panic", name)
				}
			}()
			fn()
		}()
	}
}

func TestHTTPHeaderDigest(t *testing.T) {
	body := model.NewInMemoryDocument([]byte("hello world"))
	hd := NewHTTPHeaderDigest(body, enumerations.DigestAlgorithmSHA256)

	if hd.Name() != DSSJsonUtilsHTTPHeaderDigest {
		t.Errorf("Name() = %q, want %q", hd.Name(), DSSJsonUtilsHTTPHeaderDigest)
	}
	digestValue, err := body.DigestValue(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatalf("DigestValue: %v", err)
	}
	want := enumerations.DigestAlgorithmSHA256.HttpHeaderAlgo() + "=" + utils.ToBase64(digestValue)
	if hd.Value() != want {
		t.Errorf("Value() = %q, want %q", hd.Value(), want)
	}
	if hd.MessageBodyDocument() != model.DSSDocument(body) {
		t.Error("MessageBodyDocument() should return the original document")
	}

	var doc model.DSSDocument = hd
	if doc.Name() != DSSJsonUtilsHTTPHeaderDigest {
		t.Error("HTTPHeaderDigest should satisfy model.DSSDocument")
	}
}

func TestHTTPHeaderDigest_UnsupportedDigestAlgorithmPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected a panic for a DigestAlgorithm with no RFC 5843 http header mapping")
		}
	}()
	body := model.NewInMemoryDocument([]byte("hello world"))
	NewHTTPHeaderDigest(body, enumerations.DigestAlgorithmSHA384)
}

func TestHTTPHeaderSignatureScope(t *testing.T) {
	named := model.NewInMemoryDocument([]byte("payload"))
	named.SetName("my-doc")
	scope := NewHTTPHeaderSignatureScope(named)
	if scope.DocumentName() != "my-doc" {
		t.Errorf("DocumentName() = %q, want my-doc", scope.DocumentName())
	}
	if scope.Type() != enumerations.SignatureScopeTypeFull {
		t.Errorf("Type() = %v, want FULL", scope.Type())
	}
	if scope.Description(nil) != "Payload value digest" {
		t.Errorf("Description() = %q, want %q", scope.Description(nil), "Payload value digest")
	}

	unnamed := model.NewInMemoryDocument([]byte("payload"))
	scope2 := NewHTTPHeaderSignatureScope(unnamed)
	if scope2.DocumentName() != "HttpHeaders payload" {
		t.Errorf("DocumentName() for unnamed document = %q, want fallback name", scope2.DocumentName())
	}
}

func TestHTTPHeaderMessageBodySignatureScope(t *testing.T) {
	body := model.NewInMemoryDocument([]byte("payload"))
	scope := NewHTTPHeaderMessageBodySignatureScope(body)
	if scope.Description(nil) != "Message body value digest" {
		t.Errorf("Description() = %q, want %q", scope.Description(nil), "Message body value digest")
	}
	if scope.Type() != enumerations.SignatureScopeTypeFull {
		t.Errorf("Type() = %v, want FULL (inherited from HTTPHeaderSignatureScope)", scope.Type())
	}
}
