package model

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/utain/esig/dss/enumerations"
)

func TestFileDocumentMissingFileErrors(t *testing.T) {
	if _, err := NewFileDocument(filepath.Join(t.TempDir(), "does-not-exist.txt")); err == nil {
		t.Fatal("expected error for a non-existent file")
	}
}

func TestFileDocumentOpenStreamAndDigest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")
	content := []byte("file document content")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	doc, err := NewFileDocument(path)
	if err != nil {
		t.Fatalf("NewFileDocument: %v", err)
	}
	if doc.Name() != "sample.txt" {
		t.Fatalf("Name() = %q, want sample.txt", doc.Name())
	}
	if !doc.Exists() {
		t.Fatal("expected Exists() to be true")
	}

	rc, err := doc.OpenStream()
	if err != nil {
		t.Fatalf("OpenStream: %v", err)
	}
	defer rc.Close()
	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Fatalf("got %q, want %q", got, content)
	}

	digestValue, err := doc.DigestValue(enumerations.DigestAlgorithm_SHA256)
	if err != nil {
		t.Fatalf("DigestValue: %v", err)
	}
	if len(digestValue) != 32 {
		t.Fatalf("expected a 32-byte SHA-256 digest, got %d bytes", len(digestValue))
	}
}

func TestFileDocumentSave(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "src.txt")
	content := []byte("save me")
	if err := os.WriteFile(srcPath, content, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	doc, err := NewFileDocument(srcPath)
	if err != nil {
		t.Fatalf("NewFileDocument: %v", err)
	}

	dstPath := filepath.Join(dir, "dst.txt")
	if err := doc.Save(dstPath); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Fatalf("got %q, want %q", got, content)
	}
}
