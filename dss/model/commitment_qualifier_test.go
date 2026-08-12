package model

import "testing"

func TestCommitmentQualifierRoundTrip(t *testing.T) {
	q := NewCommitmentQualifier()
	q.SetOid("1.2.3.4")
	doc := NewInMemoryDocument([]byte("qualifier content"))
	q.SetContent(doc)

	if q.Oid() != "1.2.3.4" {
		t.Fatalf("Oid() = %q, want 1.2.3.4", q.Oid())
	}
	if q.Content() != doc {
		t.Fatal("Content() did not round-trip")
	}
}
