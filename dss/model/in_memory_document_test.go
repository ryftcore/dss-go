package model

import (
	"bytes"
	"crypto/sha256"
	"io"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

func TestInMemoryDocumentOpenStreamRoundTrip(t *testing.T) {
	data := []byte("hello dss")
	doc := NewInMemoryDocument(data)

	rc, err := doc.OpenStream()
	if err != nil {
		t.Fatalf("OpenStream: %v", err)
	}
	defer rc.Close()

	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Fatalf("got %q, want %q", got, data)
	}
}

func TestInMemoryDocumentNilBytesPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on nil bytes")
		}
	}()
	NewInMemoryDocument(nil)
}

func TestInMemoryDocumentWithNameDerivesMimeType(t *testing.T) {
	doc := NewInMemoryDocumentWithName([]byte("<xml/>"), "sample.xml")
	if doc.Name() != "sample.xml" {
		t.Fatalf("Name() = %q, want sample.xml", doc.Name())
	}
	if doc.MimeType() == nil {
		t.Fatal("expected non-nil MimeType derived from file name")
	}
}

func TestInMemoryDocumentDigestValueMatchesSHA256AndCaches(t *testing.T) {
	data := []byte("The quick brown fox jumps over the lazy dog")
	doc := NewInMemoryDocument(data)

	want := sha256.Sum256(data)
	got, err := doc.DigestValue(enumerations.DigestAlgorithm_SHA256)
	if err != nil {
		t.Fatalf("DigestValue: %v", err)
	}
	if !bytes.Equal(got, want[:]) {
		t.Fatalf("digest mismatch: got %x, want %x", got, want)
	}

	// Second call must hit the cache and return the identical bytes.
	got2, err := doc.DigestValue(enumerations.DigestAlgorithm_SHA256)
	if err != nil {
		t.Fatalf("DigestValue (cached): %v", err)
	}
	if !bytes.Equal(got2, want[:]) {
		t.Fatalf("cached digest mismatch: got %x, want %x", got2, want)
	}
}

func TestInMemoryDocumentDigestUnsupportedAlgorithm(t *testing.T) {
	doc := NewInMemoryDocument([]byte("data"))
	if _, err := doc.DigestValue(enumerations.DigestAlgorithm_MD2); err == nil {
		t.Fatal("expected error for unsupported digest algorithm MD2")
	}
}

func TestInMemoryDocumentWriteToAndSave(t *testing.T) {
	data := []byte("write-to content")
	doc := NewInMemoryDocument(data)

	var buf bytes.Buffer
	n, err := doc.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo: %v", err)
	}
	if n != int64(len(data)) {
		t.Fatalf("WriteTo wrote %d bytes, want %d", n, len(data))
	}
	if !bytes.Equal(buf.Bytes(), data) {
		t.Fatalf("WriteTo content = %q, want %q", buf.Bytes(), data)
	}
}

func TestInMemoryDocumentEquals(t *testing.T) {
	a := NewInMemoryDocumentWithName([]byte("x"), "a.txt")
	b := NewInMemoryDocumentWithName([]byte("x"), "a.txt")
	c := NewInMemoryDocumentWithName([]byte("y"), "a.txt")

	if !a.Equals(b) {
		t.Fatal("expected a.Equals(b) to be true")
	}
	if a.Equals(c) {
		t.Fatal("expected a.Equals(c) to be false (different bytes)")
	}
}

func TestCreateEmptyDocument(t *testing.T) {
	doc := CreateEmptyDocument()
	if len(doc.Bytes()) != 0 {
		t.Fatalf("expected empty bytes, got %v", doc.Bytes())
	}
}
