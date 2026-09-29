package scope

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

// incomparableDocument is a DSSDocument implementation, other than the ones of package model,
// that is a value type holding a slice, so two of its values cannot be compared with ==. The
// embedded nil interface supplies the DSSDocument methods, none of which is called.
type incomparableDocument struct {
	model.DSSDocument
	parts [][]byte
}

// pointerDocument is a comparable custom DSSDocument.
type pointerDocument struct {
	model.DSSDocument
}

// TestDocumentsEqualFallbackNeverPanics pins the fallback for DSSDocument implementations other
// than InMemoryDocument, DigestDocument and FileDocument: Java's default Object#equals is
// reference equality and cannot fail, so neither may Go's (T22D3-STD-001).
func TestDocumentsEqualFallbackNeverPanics(t *testing.T) {
	first := incomparableDocument{parts: [][]byte{{1}}}
	second := incomparableDocument{parts: [][]byte{{1}}}
	if documentsEqual(first, second) {
		t.Error("two distinct values of an incomparable document type were reported equal")
	}

	same := &pointerDocument{}
	other := &pointerDocument{}
	if !documentsEqual(same, same) {
		t.Error("a document is not equal to itself")
	}
	if documentsEqual(same, other) {
		t.Error("two distinct custom documents were reported equal")
	}

	if documentsEqual(model.NewInMemoryDocument([]byte("a")), same) {
		t.Error("an in-memory document was reported equal to a custom document")
	}
	if !documentsEqual(model.NewInMemoryDocument([]byte("a")), model.NewInMemoryDocument([]byte("a"))) {
		t.Error("in-memory documents with the same content were reported different")
	}
}
