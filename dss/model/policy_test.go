package model

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestPolicyIsEmpty(t *testing.T) {
	p := NewPolicy()
	if !p.IsEmpty() {
		t.Fatal("expected fresh Policy to be empty")
	}
	p.SetId("urn:policy:1")
	if p.IsEmpty() {
		t.Fatal("expected Policy with an Id to not be empty")
	}
}

func TestPolicyIsSPQualifierPresent(t *testing.T) {
	p := NewPolicy()
	if p.IsSPQualifierPresent() {
		t.Fatal("expected fresh Policy to have no SP qualifier")
	}
	p.SetSpuri("http://example.org/policy.pdf")
	if !p.IsSPQualifierPresent() {
		t.Fatal("expected Policy with SpUri to have an SP qualifier")
	}
}

func TestPolicyRoundTripAndEquals(t *testing.T) {
	a := NewPolicy()
	a.SetId("urn:policy:1")
	a.SetQualifier(enumerations.ObjectIdentifierQualifierOIDAsURN)
	a.SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	a.SetDigestValue([]byte{1, 2, 3})
	a.SetDocumentationReferences("ref1", "ref2")

	b := NewPolicy()
	b.SetId("urn:policy:1")
	b.SetQualifier(enumerations.ObjectIdentifierQualifierOIDAsURN)
	b.SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	b.SetDigestValue([]byte{1, 2, 3})
	b.SetDocumentationReferences("ref1", "ref2")

	if !a.Equals(b) {
		t.Fatalf("expected equal policies to be Equals(): a=%s b=%s", a.String(), b.String())
	}

	b.SetId("urn:policy:2")
	if a.Equals(b) {
		t.Fatal("expected different Ids to not be Equals()")
	}
}
