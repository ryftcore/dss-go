package tsl

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestCertificateContentEquivalenceRoundTrip(t *testing.T) {
	c := NewCertificateContentEquivalence()
	c.SetContext(enumerations.MRAEquivalenceContext_QC_COMPLIANCE)
	oids := NewQCStatementOids()
	oids.SetQcTypeIds([]string{"0.4.0.1862.1.6.1"})
	c.SetContentReplacement(oids)

	if c.Context() != enumerations.MRAEquivalenceContext_QC_COMPLIANCE {
		t.Fatalf("unexpected Context: %v", c.Context())
	}
	if c.ContentReplacement() != oids {
		t.Fatalf("unexpected ContentReplacement: %v", c.ContentReplacement())
	}
	if c.Condition() != nil {
		t.Fatalf("expected nil Condition, got %v", c.Condition())
	}
}

func TestCertificateContentEquivalenceEquals(t *testing.T) {
	a := NewCertificateContentEquivalence()
	a.SetContext(enumerations.MRAEquivalenceContext_QC_TYPE)
	b := NewCertificateContentEquivalence()
	b.SetContext(enumerations.MRAEquivalenceContext_QC_TYPE)

	if !a.Equals(b) {
		t.Fatalf("expected equal CertificateContentEquivalence values")
	}

	b.SetContext(enumerations.MRAEquivalenceContext_QC_COMPLIANCE)
	if a.Equals(b) {
		t.Fatalf("expected unequal CertificateContentEquivalence values after context change")
	}
	if a.Equals(nil) {
		t.Fatalf("expected Equals(nil) to be false")
	}
}
