package model

import "testing"

func TestManifestFileEntriesAndRootFile(t *testing.T) {
	mf := NewManifestFile()

	root := NewManifestEntry()
	root.SetUri("root.xml")
	root.SetRootfile(true)

	other := NewManifestEntry()
	other.SetUri("other.xml")

	mf.SetEntries([]*ManifestEntry{other, root})

	if !mf.IsDocumentCovered("other.xml") {
		t.Fatal("expected other.xml to be covered")
	}
	if mf.IsDocumentCovered("missing.xml") {
		t.Fatal("expected missing.xml to not be covered")
	}
	got := mf.RootFile()
	if got != root {
		t.Fatal("RootFile() did not return the entry with Rootfile()==true")
	}
}

func TestManifestFileEntriesLazyInit(t *testing.T) {
	mf := NewManifestFile()
	if mf.Entries() == nil {
		t.Fatal("expected Entries() to lazily initialize to a non-nil slice")
	}
	if mf.RootFile() != nil {
		t.Fatal("expected RootFile() to be nil when no entries are set")
	}
}

func TestManifestFileFilenameDelegatesToDocument(t *testing.T) {
	mf := NewManifestFile()
	doc := NewInMemoryDocumentWithName([]byte("x"), "manifest.xml")
	mf.SetDocument(doc)
	if mf.Filename() != "manifest.xml" {
		t.Fatalf("Filename() = %q, want manifest.xml", mf.Filename())
	}
}
