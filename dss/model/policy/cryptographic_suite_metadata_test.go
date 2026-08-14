// Ported from dss-model/.../model/policy/crypto/CryptographicSuiteMetadata.java (DSS 6.5.RC1).
package policy

import (
	"testing"
	"time"
)

func TestCryptographicSuiteMetadata_RoundTrip(t *testing.T) {
	m := NewCryptographicSuiteMetadata()
	m.SetPolicyName("ETSI TS 119 312")
	m.SetPolicyOID("0.4.0.2231.1")
	m.SetPolicyURI("https://example.org/policy")
	m.SetPublisherName("ETSI")
	m.SetPublisherAddress("650 Route des Lucioles")
	m.SetPublisherURI("https://etsi.org")
	issued := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	m.SetPolicyIssueDate(&issued)
	next := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	m.SetNextUpdate(&next)
	m.SetUsage("validation")
	m.SetVersion("1.0")
	m.SetLang("en")
	m.SetId("id-1")

	if m.PolicyName() != "ETSI TS 119 312" {
		t.Fatalf("PolicyName() = %q", m.PolicyName())
	}
	if m.PolicyIssueDate() == nil || !m.PolicyIssueDate().Equal(issued) {
		t.Fatalf("PolicyIssueDate() = %v, want %v", m.PolicyIssueDate(), issued)
	}
	if m.Id() != "id-1" {
		t.Fatalf("Id() = %q", m.Id())
	}

	other := NewCryptographicSuiteMetadata()
	*other = *m
	if !m.Equals(other) {
		t.Fatalf("Equals() = false for equal instances")
	}
	other.SetVersion("2.0")
	if m.Equals(other) {
		t.Fatalf("Equals() = true after divergence")
	}
}
