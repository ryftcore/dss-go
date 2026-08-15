package qwac

import "testing"

func TestGetTLSCertificateBindingUrl_FindsMatchingRel(t *testing.T) {
	headers := map[string][]string{
		"Link": {`<https://example.com/binding>; rel="tls-certificate-binding"`},
	}
	if got := GetTLSCertificateBindingUrl(headers); got != "https://example.com/binding" {
		t.Fatalf("GetTLSCertificateBindingUrl() = %q, want %q", got, "https://example.com/binding")
	}
}

func TestGetTLSCertificateBindingUrl_IgnoresOtherRel(t *testing.T) {
	headers := map[string][]string{
		"Link": {`<https://example.com/other>; rel="alternate"`},
	}
	if got := GetTLSCertificateBindingUrl(headers); got != "" {
		t.Fatalf("GetTLSCertificateBindingUrl() = %q, want empty", got)
	}
}

func TestGetTLSCertificateBindingUrl_NilOrMissingHeader(t *testing.T) {
	if got := GetTLSCertificateBindingUrl(nil); got != "" {
		t.Fatalf("GetTLSCertificateBindingUrl(nil) = %q, want empty", got)
	}
	if got := GetTLSCertificateBindingUrl(map[string][]string{"Other": {"x"}}); got != "" {
		t.Fatalf("GetTLSCertificateBindingUrl(no Link) = %q, want empty", got)
	}
}

func TestGetTLSCertificateBindingUrl_SkipsUnparsableValueAndFindsNext(t *testing.T) {
	headers := map[string][]string{
		"Link": {
			`not-a-valid-link-header`,
			`<https://example.com/binding>; rel="tls-certificate-binding"`,
		},
	}
	if got := GetTLSCertificateBindingUrl(headers); got != "https://example.com/binding" {
		t.Fatalf("GetTLSCertificateBindingUrl() = %q, want %q", got, "https://example.com/binding")
	}
}

func TestGetTLSCertificateBindingUrl_MultipleLinksInOneValuePicksMatching(t *testing.T) {
	headers := map[string][]string{
		"Link": {`<https://example.com/alt>; rel="alternate", <https://example.com/binding>; rel="tls-certificate-binding"`},
	}
	if got := GetTLSCertificateBindingUrl(headers); got != "https://example.com/binding" {
		t.Fatalf("GetTLSCertificateBindingUrl() = %q, want %q", got, "https://example.com/binding")
	}
}
