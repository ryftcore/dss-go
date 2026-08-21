// Tests for ProfileParameters, matching
// dss-document/src/main/java/eu/europa/esig/dss/signature/ProfileParameters.java
// upstream.
package document

import (
	"testing"

	"github.com/utain/esig/dss/model"
)

func TestNewProfileParametersDefaults(t *testing.T) {
	p := NewProfileParameters()
	if p.DeterministicId() != "" {
		t.Fatalf("DeterministicId() = %q, want empty", p.DeterministicId())
	}
	if p.DetachedContents() != nil {
		t.Fatalf("DetachedContents() = %v, want nil", p.DetachedContents())
	}
}

func TestProfileParametersSetters(t *testing.T) {
	p := NewProfileParameters()
	p.SetDeterministicId("id-1")
	if got := p.DeterministicId(); got != "id-1" {
		t.Fatalf("DeterministicId() = %q, want id-1", got)
	}
	docs := []model.DSSDocument{model.NewInMemoryDocument([]byte("a"))}
	p.SetDetachedContents(docs)
	if got := p.DetachedContents(); len(got) != 1 || got[0] != docs[0] {
		t.Fatalf("DetachedContents() = %v, want %v", got, docs)
	}
}

func TestProfileParametersEquals(t *testing.T) {
	a := NewProfileParameters()
	b := NewProfileParameters()
	if !a.Equals(b) {
		t.Fatal("two fresh ProfileParameters should be equal")
	}
	if !a.Equals(a) {
		t.Fatal("a value must equal itself")
	}
	if a.Equals(nil) {
		t.Fatal("a value must not equal nil")
	}

	a.SetDeterministicId("x")
	if a.Equals(b) {
		t.Fatal("differing deterministicId should not be equal")
	}
	b.SetDeterministicId("x")
	if !a.Equals(b) {
		t.Fatal("matching deterministicId should be equal again")
	}

	doc := model.NewInMemoryDocument([]byte("a"))
	a.SetDetachedContents([]model.DSSDocument{doc})
	if a.Equals(b) {
		t.Fatal("differing detachedContents should not be equal")
	}
}

func TestProfileParametersString(t *testing.T) {
	p := NewProfileParameters()
	p.SetDeterministicId("id-1")
	if got := p.String(); got == "" {
		t.Fatal("String() should not be empty")
	}
}
