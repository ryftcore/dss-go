// Ported from dss-enumerations/.../CertificatePolicy.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestCertificatePolicy(t *testing.T) {
	cases := []struct {
		v           CertificatePolicy
		description string
		oid         string
	}{
		{CertificatePolicyQCPPublic, "qcp-public", "0.4.0.1456.1.2"},
		{CertificatePolicyQCPPublicWithSSCD, "qcp-public-with-sscd", "0.4.0.1456.1.1"},
		{CertificatePolicyNCP, "normalized-certificate-policy", "0.4.0.2042.1.1"},
		{CertificatePolicyNCPP, "normalized-certificate-policy-sscd", "0.4.0.2042.1.2"},
		{CertificatePolicyLCP, "lightweight-certificate-policy", "0.4.0.2042.1.3"},
		{CertificatePolicyEVCP, "extended-validation-certificate-policy", "0.4.0.2042.1.4"},
		{CertificatePolicyDVCP, "domain-validation-certificate-policy", "0.4.0.2042.1.6"},
		{CertificatePolicyOVCP, "organizational-validation-certificate-policy", "0.4.0.2042.1.7"},
		{CertificatePolicyIVCP, "individual-validation-certificate-policy", "0.4.0.2042.1.8"},
		{CertificatePolicyQCPNatural, "qcp-natural", "0.4.0.194112.1.0"},
		{CertificatePolicyQCPLegal, "qcp-legal", "0.4.0.194112.1.1"},
		{CertificatePolicyQCPNaturalQSCD, "qcp-natural-qscd", "0.4.0.194112.1.2"},
		{CertificatePolicyQCPLegalQSCD, "qcp-legal-qscd", "0.4.0.194112.1.3"},
		{CertificatePolicyQCPWeb, "qcp-web", "0.4.0.194112.1.4"},
		{CertificatePolicyQNCPWeb, "qncp-web", "0.4.0.194112.1.5"},
		{CertificatePolicyQNCPWebGen, "qncp-web-gen", "0.4.0.194112.1.6"},
	}
	if len(CertificatePolicyValues()) != len(cases) {
		t.Fatalf("expected %d values, got %d", len(cases), len(CertificatePolicyValues()))
	}
	for _, c := range cases {
		if got := c.v.OID(); got != c.oid {
			t.Errorf("%v.OID() = %q, want %q", c.v, got, c.oid)
		}
		if got := c.v.Description(); got != c.description {
			t.Errorf("%v.Description() = %q, want %q", c.v, got, c.description)
		}
		got, err := CertificatePolicyValueOf(string(c.v))
		if err != nil || got != c.v {
			t.Errorf("CertificatePolicyValueOf(%q) = %v, %v; want %v, nil", c.v, got, err, c.v)
		}
	}
	if _, err := CertificatePolicyValueOf("NOPE"); err == nil {
		t.Error("expected error for unknown name")
	}
}
