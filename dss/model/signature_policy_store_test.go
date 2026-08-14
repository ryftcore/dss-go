package model

import "testing"

func TestSignaturePolicyStoreRoundTrip(t *testing.T) {
	s := NewSignaturePolicyStore()
	s.SetId("policy-id")
	spDoc := NewSpDocSpecification()
	spDoc.SetId("2.2.25.1")
	s.SetSpDocSpecification(spDoc)
	doc := NewInMemoryDocument([]byte("policy content"))
	s.SetSignaturePolicyContent(doc)
	s.SetSigPolDocLocalURI("http://example.org/local")

	if s.Id() != "policy-id" {
		t.Fatalf("Id() = %q", s.Id())
	}
	if s.SpDocSpecification() != spDoc {
		t.Fatal("SpDocSpecification() did not round-trip")
	}
	if s.SignaturePolicyContent() != doc {
		t.Fatal("SignaturePolicyContent() did not round-trip")
	}
	if s.SigPolDocLocalURI() != "http://example.org/local" {
		t.Fatalf("SigPolDocLocalURI() = %q", s.SigPolDocLocalURI())
	}
}
