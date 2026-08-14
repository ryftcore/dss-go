package enumerations

import "testing"

type cryptographicSuiteAlgorithmUsageCase struct {
	v   CryptographicSuiteAlgorithmUsage
	uri string
}

func cryptographicSuiteAlgorithmUsageCases() []cryptographicSuiteAlgorithmUsageCase {
	return []cryptographicSuiteAlgorithmUsageCase{
		{CryptographicSuiteAlgorithmUsage_SIGN_DATA, "http://uri.etsi.org/19322/sign_data"},
		{CryptographicSuiteAlgorithmUsage_SIGN_CERTIFICATES, "http://uri.etsi.org/19322/sign_data/sign_certificates"},
		{CryptographicSuiteAlgorithmUsage_SIGN_OCSP, "http://uri.etsi.org/19322/sign_data/sign_ocsp"},
		{CryptographicSuiteAlgorithmUsage_SIGN_TIMESTAMPS, "http://uri.etsi.org/19322/sign_data/sign_timestamps"},
		{CryptographicSuiteAlgorithmUsage_VALIDATE_DATA, "http://uri.etsi.org/19322/sign_data/validate_data"},
		{CryptographicSuiteAlgorithmUsage_VALIDATE_CERTIFICATES, "http://uri.etsi.org/19322/sign_data/validate_data/validate_certificates"},
		{CryptographicSuiteAlgorithmUsage_VALIDATE_OCSP, "http://uri.etsi.org/19322/sign_data/validate_data/validate_ocsp"},
		{CryptographicSuiteAlgorithmUsage_VALIDATE_TIMESTAMPS, "http://uri.etsi.org/19322/sign_data/validate_data/validate_timestamps"},
	}
}

func TestCryptographicSuiteAlgorithmUsageURI(t *testing.T) {
	cases := cryptographicSuiteAlgorithmUsageCases()
	if len(CryptographicSuiteAlgorithmUsageValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(CryptographicSuiteAlgorithmUsageValues()))
	}
	for _, c := range cases {
		if got := c.v.URI(); got != c.uri {
			t.Errorf("%v.URI() = %q, want %q", c.v, got, c.uri)
		}
	}
}

func TestCryptographicSuiteAlgorithmUsageFromURI(t *testing.T) {
	for _, c := range cryptographicSuiteAlgorithmUsageCases() {
		if got := CryptographicSuiteAlgorithmUsageFromURI(c.uri); got != c.v {
			t.Errorf("CryptographicSuiteAlgorithmUsageFromURI(%q) = %v, want %v", c.uri, got, c.v)
		}
	}
	if got := CryptographicSuiteAlgorithmUsageFromURI("nope"); got != "" {
		t.Errorf("CryptographicSuiteAlgorithmUsageFromURI(unknown) = %v, want zero value", got)
	}
	if got := CryptographicSuiteAlgorithmUsageFromURI(""); got != "" {
		t.Errorf("CryptographicSuiteAlgorithmUsageFromURI(\"\") = %v, want zero value", got)
	}
}
