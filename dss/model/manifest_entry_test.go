package model

import "testing"

func TestManifestEntryRoundTrip(t *testing.T) {
	e := NewManifestEntry()
	e.SetUri("doc.xml")
	e.SetFound(true)
	e.SetIntact(true)
	e.SetRootfile(true)
	doc := NewInMemoryDocument([]byte("x"))
	e.SetDocument(doc)

	if e.Uri() != "doc.xml" {
		t.Fatalf("Uri() = %q", e.Uri())
	}
	if !e.IsFound() || !e.IsIntact() || !e.IsRootfile() {
		t.Fatal("expected found/intact/rootfile flags to be true")
	}
	if e.Document() != doc {
		t.Fatal("Document() did not round-trip")
	}
}

func TestManifestEntrySetDocumentNameIsNoOp(t *testing.T) {
	e := NewManifestEntry()
	e.SetDocumentName("ignored.txt")
	if e.Document() != nil {
		t.Fatal("expected SetDocumentName to be a no-op, matching deprecated Java behavior")
	}
}
