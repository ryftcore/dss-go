// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/scope/SignatureScope.java (DSS 6.5.RC1).
package scope

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// fakeSignatureScope is a minimal concrete SignatureScope, mirroring how a real subclass
// (e.g. FullSignatureScope, outside this manifest) would embed SignatureScopeBase.
type fakeSignatureScope struct {
	SignatureScopeBase
	description string
	scopeType   enumerations.SignatureScopeType
}

func newFakeSignatureScope(document model.DSSDocument, description string, scopeType enumerations.SignatureScopeType) *fakeSignatureScope {
	return &fakeSignatureScope{
		SignatureScopeBase: NewSignatureScopeBase(document),
		description:        description,
		scopeType:          scopeType,
	}
}

func (f *fakeSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return f.description
}

func (f *fakeSignatureScope) Type() enumerations.SignatureScopeType {
	return f.scopeType
}

var _ SignatureScope = (*fakeSignatureScope)(nil)

func TestSignatureScope_NameDefaultsToDocumentName(t *testing.T) {
	doc := model.NewInMemoryDocumentWithName([]byte("hello world"), "doc.txt")
	s := newFakeSignatureScope(doc, "a full scope", enumerations.SignatureScopeTypeFull)

	if got, want := s.DocumentName(), "doc.txt"; got != want {
		t.Fatalf("DocumentName() = %q, want %q", got, want)
	}
	if got, want := s.Name(nil), "doc.txt"; got != want {
		t.Fatalf("Name(nil) = %q, want %q", got, want)
	}
	if got, want := s.Description(nil), "a full scope"; got != want {
		t.Fatalf("Description(nil) = %q, want %q", got, want)
	}
	if got, want := s.Type(), enumerations.SignatureScopeTypeFull; got != want {
		t.Fatalf("Type() = %v, want %v", got, want)
	}
	if got := s.Transformations(); got != nil {
		t.Fatalf("Transformations() = %v, want nil", got)
	}
}

func TestSignatureScope_DigestOfPlainDocument(t *testing.T) {
	doc := model.NewInMemoryDocumentWithName([]byte("hello world"), "doc.txt")
	s := newFakeSignatureScope(doc, "", enumerations.SignatureScopeTypeFull)

	digest, err := s.Digest(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatalf("Digest() error = %v", err)
	}
	want, err := doc.DigestValue(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatalf("DigestValue() error = %v", err)
	}
	if string(digest.Value()) != string(want) {
		t.Fatalf("Digest().Value() = %x, want %x", digest.Value(), want)
	}
}

func TestSignatureScope_DigestOfDigestDocumentReturnsExistingDigest(t *testing.T) {
	digestDoc := model.NewDigestDocumentFromValueWithName(enumerations.DigestAlgorithmSHA256, []byte("precomputed"), "doc.txt")
	s := newFakeSignatureScope(digestDoc, "", enumerations.SignatureScopeTypeDigest)

	digest, err := s.Digest(enumerations.DigestAlgorithmSHA1)
	if err != nil {
		t.Fatalf("Digest() error = %v", err)
	}
	if digest.Algorithm() != enumerations.DigestAlgorithmSHA256 {
		t.Fatalf("Digest().Algorithm() = %v, want the existing SHA256 digest's algorithm, not the requested SHA1", digest.Algorithm())
	}
	if string(digest.Value()) != "precomputed" {
		t.Fatalf("Digest().Value() = %q, want %q", digest.Value(), "precomputed")
	}
}

func TestSignatureScope_DigestOfNilDocumentIsEmpty(t *testing.T) {
	// Uses the (name, document) constructor directly: unlike newFakeSignatureScope, the
	// single-argument SignatureScope(DSSDocument) constructor dereferences document.Name(),
	// which cannot be done on a nil document (matching Java's NullPointerException there).
	s := &fakeSignatureScope{
		SignatureScopeBase: NewSignatureScopeBaseWithName("", nil),
		scopeType:          enumerations.SignatureScopeTypeFull,
	}
	digest, err := s.Digest(enumerations.DigestAlgorithmSHA256)
	if err != nil {
		t.Fatalf("Digest() error = %v", err)
	}
	if !digest.IsEmpty() {
		t.Fatalf("Digest() = %v, want the zero Digest for a scope without a document", digest)
	}
}

func TestSignatureScope_ChildrenAndEquals(t *testing.T) {
	doc := model.NewInMemoryDocumentWithName([]byte("hello world"), "doc.txt")
	parent := newFakeSignatureScope(doc, "", enumerations.SignatureScopeTypeFull)
	child := newFakeSignatureScope(doc, "", enumerations.SignatureScopeTypePartial)

	if got := len(parent.Children()); got != 0 {
		t.Fatalf("Children() = %d entries, want 0 before adding any", got)
	}
	parent.AddChildSignatureScope(child)
	children := parent.Children()
	if len(children) != 1 || children[0] != SignatureScope(child) {
		t.Fatalf("Children() = %v, want [%v]", children, child)
	}

	sameDoc := newFakeSignatureScope(doc, "different description", enumerations.SignatureScopeTypeFull)
	if !parent.Equals(sameDoc) {
		t.Fatalf("Equals() = false for two scopes over the same document and name")
	}
	otherDoc := model.NewInMemoryDocumentWithName([]byte("other content"), "other.txt")
	different := newFakeSignatureScope(otherDoc, "", enumerations.SignatureScopeTypeFull)
	if parent.Equals(different) {
		t.Fatalf("Equals() = true for two scopes over different documents")
	}
	if parent.Equals(nil) {
		t.Fatalf("Equals(nil) = true")
	}
}

func TestSignatureScope_DSSIDAsString(t *testing.T) {
	doc := model.NewInMemoryDocumentWithName([]byte("hello world"), "doc.txt")
	s := newFakeSignatureScope(doc, "", enumerations.SignatureScopeTypeFull)

	want, err := model.NewDataIdentifierForDocument("doc.txt", doc)
	if err != nil {
		t.Fatalf("NewDataIdentifierForDocument() error = %v", err)
	}
	if got := s.DSSIDAsString(); got != want.AsXmlID() {
		t.Fatalf("DSSIDAsString() = %q, want %q", got, want.AsXmlID())
	}
	// DSSID() must be idempotent (cached).
	if s.DSSID() != s.DSSID() {
		t.Fatalf("DSSID() is not cached across calls")
	}
}
