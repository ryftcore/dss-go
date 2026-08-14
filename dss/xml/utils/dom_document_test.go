package utils

import (
	"bytes"
	"testing"

	"github.com/utain/esig/dss/enumerations"
)

func TestNewDOMDocumentNilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for a nil node")
		}
	}()
	NewDOMDocument(nil)
}

func TestDOMDocumentOpenStreamAndMimeType(t *testing.T) {
	doc, err := DomUtilsBuildDOMFromString(`<r a="1"><c/></r>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d := NewDOMDocumentWithName(doc, "signature.xml")
	if d.Name() != "signature.xml" {
		t.Errorf("Name() = %q", d.Name())
	}
	if d.MimeType() != enumerations.MimeTypeEnum_XML {
		t.Errorf("expected the default MimeType to be XML, got %v", d.MimeType())
	}

	rc, err := d.OpenStream()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rc.Close()
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(rc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// DOMDocument streams DomUtilsWriteDocumentTo's bytes, which are byte-for-byte the
	// identity Transformer's - declaration, standalone and short empty element included.
	if want := `<?xml version="1.0" encoding="UTF-8" standalone="no"?><r a="1"><c/></r>`; buf.String() != want {
		t.Fatalf("serialized content\n got: %q\nwant: %q", buf.String(), want)
	}
}

func TestDOMDocumentDigestIsCached(t *testing.T) {
	doc, err := DomUtilsBuildDOMFromString(`<r/>`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d := NewDOMDocument(doc)

	v1, err := d.DigestValue(enumerations.DigestAlgorithm_SHA256)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v2, err := d.DigestValue(enumerations.DigestAlgorithm_SHA256)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(v1, v2) {
		t.Fatalf("expected cached digest to be stable across calls")
	}
	if len(v1) != 32 {
		t.Fatalf("expected a 32-byte SHA-256 digest, got %d bytes", len(v1))
	}
}

func TestDOMDocumentEquals(t *testing.T) {
	doc1, _ := DomUtilsBuildDOMFromString(`<r a="1"/>`)
	doc2, _ := DomUtilsBuildDOMFromString(`<r a="1"/>`)
	doc3, _ := DomUtilsBuildDOMFromString(`<r a="2"/>`)

	d1 := NewDOMDocument(doc1)
	d2 := NewDOMDocument(doc2)
	d3 := NewDOMDocument(doc3)

	if !d1.Equals(d2) {
		t.Error("expected documents with identical serialized content to be equal")
	}
	if d1.Equals(d3) {
		t.Error("expected documents with different content to not be equal")
	}
	if d1.Equals(nil) {
		t.Error("expected Equals(nil) to be false")
	}
	if !d1.Equals(d1) {
		t.Error("expected a document to equal itself")
	}
}
