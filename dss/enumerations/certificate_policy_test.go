// Ported from dss-enumerations/.../CertificatePolicy.java (DSS 6.5.RC1).
package enumerations

import "testing"

func TestCertificatePolicy(t *testing.T) {
	cases := []struct {
		v           CertificatePolicy
		description string
		oid         string
	}{
		{CertificatePolicy_QCP_PUBLIC, "qcp-public", "0.4.0.1456.1.2"},
		{CertificatePolicy_QCP_PUBLIC_WITH_SSCD, "qcp-public-with-sscd", "0.4.0.1456.1.1"},
		{CertificatePolicy_NCP, "normalized-certificate-policy", "0.4.0.2042.1.1"},
		{CertificatePolicy_NCPP, "normalized-certificate-policy-sscd", "0.4.0.2042.1.2"},
		{CertificatePolicy_LCP, "lightweight-certificate-policy", "0.4.0.2042.1.3"},
		{CertificatePolicy_EVCP, "extended-validation-certificate-policy", "0.4.0.2042.1.4"},
		{CertificatePolicy_DVCP, "domain-validation-certificate-policy", "0.4.0.2042.1.6"},
		{CertificatePolicy_OVCP, "organizational-validation-certificate-policy", "0.4.0.2042.1.7"},
		{CertificatePolicy_IVCP, "individual-validation-certificate-policy", "0.4.0.2042.1.8"},
		{CertificatePolicy_QCP_NATURAL, "qcp-natural", "0.4.0.194112.1.0"},
		{CertificatePolicy_QCP_LEGAL, "qcp-legal", "0.4.0.194112.1.1"},
		{CertificatePolicy_QCP_NATURAL_QSCD, "qcp-natural-qscd", "0.4.0.194112.1.2"},
		{CertificatePolicy_QCP_LEGAL_QSCD, "qcp-legal-qscd", "0.4.0.194112.1.3"},
		{CertificatePolicy_QCP_WEB, "qcp-web", "0.4.0.194112.1.4"},
		{CertificatePolicy_QNCP_WEB, "qncp-web", "0.4.0.194112.1.5"},
		{CertificatePolicy_QNCP_WEB_GEN, "qncp-web-gen", "0.4.0.194112.1.6"},
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
