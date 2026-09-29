package jades

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	modelscope "github.com/ryftcore/dss-go/dss/model/scope"
)

// TestJavaSplit pins java.lang.String#split(String) semantics: trailing empty strings are removed.
func TestJavaSplit(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"SHA-256=abc", []string{"SHA-256", "abc"}},
		{"SHA-256=abc=", []string{"SHA-256", "abc"}},
		{"SHA-256=abc==", []string{"SHA-256", "abc"}},
		{"SHA-256=a=b", []string{"SHA-256", "a", "b"}},
		{"SHA-256=", []string{"SHA-256"}},
		{"=abc", []string{"", "abc"}},
		{"=", []string{}},
		{"", []string{""}},
		{"noseparator", []string{"noseparator"}},
	}
	for _, c := range cases {
		got := javaSplit(c.in, "=")
		if len(got) != len(c.want) {
			t.Errorf("javaSplit(%q) = %q, want %q", c.in, got, c.want)
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("javaSplit(%q) = %q, want %q", c.in, got, c.want)
				break
			}
		}
	}
}

// TestSignatureScopeFinderHTTPHeaderDigestScope checks that the 'Digest' HTTP header value that
// HTTPHeaderDigest itself produces (RFC 3230 'algo=' + padded base64, e.g. "SHA-256=<43 chars>=")
// is parsed back, so the message body scope is reported (J17C-STD-001), both for a
// HTTPHeaderDigest and for a plain HTTPHeader named 'Digest' (Java: `document instanceof
// HTTPHeader`).
func TestSignatureScopeFinderHTTPHeaderDigestScope(t *testing.T) {
	body := model.NewInMemoryDocument([]byte("hello world"))
	finder := NewSignatureScopeFinder()

	hd := NewHTTPHeaderDigest(body, enumerations.DigestAlgorithmSHA256)
	if v := hd.Value(); v[len(v)-1] != '=' {
		t.Fatalf("test premise: SHA-256 instance digest %q is expected to be padded", v)
	}
	other := NewHTTPHeader("Content-Type", "application/json")

	wantDigest, err := body.Digest(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatal(err)
	}
	checkMessageBodyDigest := func(t *testing.T, digestHeader model.DSSDocument) {
		t.Helper()
		scopes := finder.getHttpHeaderSignatureScope([]model.DSSDocument{other, digestHeader})
		if len(scopes) != 2 {
			t.Fatalf("expected the payload scope and the message body scope, got %d scopes", len(scopes))
		}
		if _, ok := scopes[1].(*HTTPHeaderMessageBodySignatureScope); !ok {
			t.Fatalf("second scope is %T, want *HTTPHeaderMessageBodySignatureScope", scopes[1])
		}
		gotDigest, err := scopes[1].Digest(enumerations.DigestAlgorithmSHA256)
		if err != nil {
			t.Fatal(err)
		}
		if gotDigest.HexValue() != wantDigest.HexValue() {
			t.Errorf("message body scope digest = %s, want %s", gotDigest.HexValue(), wantDigest.HexValue())
		}
	}

	t.Run("HTTPHeaderDigest", func(t *testing.T) {
		checkMessageBodyDigest(t, hd)
	})

	t.Run("plain HTTPHeader named Digest", func(t *testing.T) {
		checkMessageBodyDigest(t, NewHTTPHeader(DSSJsonUtilsHTTPHeaderDigest, hd.Value()))
	})

	t.Run("not conformant Digest value", func(t *testing.T) {
		for _, value := range []string{"SHA-256", "SHA-256=", "SHA-256=a=b", "UNKNOWN-ALGO=AAAA", ""} {
			var scopes []modelscope.SignatureScope = finder.getHttpHeaderSignatureScope(
				[]model.DSSDocument{NewHTTPHeader(DSSJsonUtilsHTTPHeaderDigest, value)})
			if len(scopes) != 1 {
				t.Errorf("Digest value %q: expected only the payload scope, got %d scopes", value, len(scopes))
			}
		}
	})
}
