// Known-answer tests for the signature scope value objects. Every expectation is the literal
// output of the corresponding upstream class (DSS 6.5.RC1), captured by instantiating
// eu.europa.esig.dss.spi.validation.scope.{Full,Container,ContainerContent,Digest}SignatureScope
// on an InMemoryDocument holding scopeKatContent and printing getName(null), getDescription(null),
// getType() and getDigest(SHA256) - so these assertions compare against upstream itself rather
// than against a hand-written expectation.
package scope

import (
	"encoding/hex"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	mscope "github.com/ryftcore/dss-go/dss/model/scope"
)

// scopeKatContent is the document body the Java oracle was run against.
var scopeKatContent = []byte("scope-test-content")

// scopeKatDigestSHA256 is SHA-256 over scopeKatContent, as upstream getDigest(SHA256) reports it.
const scopeKatDigestSHA256 = "e628d1067bb9566463e016ad5482e4760d67c73a4cba2350890fcced4117b91e"

func TestSignatureScopeKnownAnswers(t *testing.T) {
	named := model.NewInMemoryDocumentWithName(scopeKatContent, "doc.txt")
	unnamed := model.NewInMemoryDocument(scopeKatContent)

	for _, testCase := range []struct {
		key         string
		scope       mscope.SignatureScope
		name        string
		description string
		scopeType   enumerations.SignatureScopeType
	}{
		{
			key:         "full.named",
			scope:       NewFullSignatureScope("doc.txt", named),
			name:        "doc.txt",
			description: "Full document",
			scopeType:   enumerations.SignatureScopeTypeFull,
		},
		{
			// The explicit filename wins over the document's own name.
			key:         "full.other",
			scope:       NewFullSignatureScope("other-name", named),
			name:        "other-name",
			description: "Full document",
			scopeType:   enumerations.SignatureScopeTypeFull,
		},
		{
			key:         "full.unnamed",
			scope:       NewFullSignatureScope("unnamed", unnamed),
			name:        "unnamed",
			description: "Full document",
			scopeType:   enumerations.SignatureScopeTypeFull,
		},
		{
			key:         "container",
			scope:       NewContainerSignatureScopeWithName("archive.asice", named),
			name:        "archive.asice",
			description: "ASiCS archive",
			scopeType:   enumerations.SignatureScopeTypeFull,
		},
		{
			// ContainerContentSignatureScope takes its name from the document, not a parameter.
			key:         "containerContent",
			scope:       NewContainerContentSignatureScope(named),
			name:        "doc.txt",
			description: "ASiCS archive content",
			scopeType:   enumerations.SignatureScopeTypeArchived,
		},
		{
			key:         "digest",
			scope:       NewDigestSignatureScope("digest-doc", named),
			name:        "digest-doc",
			description: "Digest of the document content",
			scopeType:   enumerations.SignatureScopeTypeDigest,
		},
	} {
		t.Run(testCase.key, func(t *testing.T) {
			if got := testCase.scope.Name(nil); got != testCase.name {
				t.Errorf("Name() = %q, want %q", got, testCase.name)
			}
			if got := testCase.scope.Description(nil); got != testCase.description {
				t.Errorf("Description() = %q, want %q", got, testCase.description)
			}
			if got := testCase.scope.Type(); got != testCase.scopeType {
				t.Errorf("Type() = %q, want %q", got, testCase.scopeType)
			}
			digest, err := testCase.scope.Digest(enumerations.DigestAlgorithmSHA256)
			if err != nil {
				t.Fatalf("Digest(SHA256): %v", err)
			}
			if got := hex.EncodeToString(digest.Value()); got != scopeKatDigestSHA256 {
				t.Errorf("Digest(SHA256) = %s, want %s", got, scopeKatDigestSHA256)
			}
			if got := digest.Algorithm(); got != enumerations.DigestAlgorithmSHA256 {
				t.Errorf("Digest(SHA256).Algorithm() = %q, want SHA256", got)
			}
		})
	}
}
