// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/SignaturePolicy.java (DSS 6.5.RC1).
package signature

import (
	"reflect"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

func TestSignaturePolicy_DefaultConstructorIsImplicit(t *testing.T) {
	p := NewPolicy()
	if got, want := p.Identifier(), "IMPLICIT_POLICY"; got != want {
		t.Fatalf("Identifier() = %q, want %q", got, want)
	}
}

func TestSignaturePolicy_RoundTrip(t *testing.T) {
	p := NewPolicyWithIdentifier("1.2.3.4")
	doc := model.NewInMemoryDocument([]byte("policy"))
	digest := model.NewDigest(enumerations.DigestAlgorithmSHA256, []byte{1, 2, 3})
	userNotice := model.NewUserNotice()
	userNotice.SetExplicitText("read me")
	docSpec := model.NewSpDocSpecification()
	docSpec.SetId("1.2.3.4.5")
	validationResult := NewPolicyValidationResult()

	p.SetDescription("a policy")
	p.SetPolicyContent(doc)
	p.SetDigest(digest)
	p.SetDocumentationReferences([]string{"ref1"})
	p.SetZeroHash(true)
	p.SetHashAsInTechnicalSpecification(true)
	p.SetURI("http://example.org/policy.pdf")
	p.SetUserNotice(userNotice)
	p.SetDocSpecification(docSpec)
	p.SetValidationResult(validationResult)

	if got, want := p.Identifier(), "1.2.3.4"; got != want {
		t.Fatalf("Identifier() = %q, want %q", got, want)
	}
	if got, want := p.Description(), "a policy"; got != want {
		t.Fatalf("Description() = %q, want %q", got, want)
	}
	if p.PolicyContent() != doc {
		t.Fatalf("PolicyContent() did not round-trip")
	}
	if !p.Digest().Equals(digest) {
		t.Fatalf("Digest() = %v, want %v", p.Digest(), digest)
	}
	if got, want := p.DocumentationReferences(), []string{"ref1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("DocumentationReferences() = %v, want %v", got, want)
	}
	if !p.IsZeroHash() || !p.IsHashAsInTechnicalSpecification() {
		t.Fatalf("boolean flags did not round-trip")
	}
	if got, want := p.URI(), "http://example.org/policy.pdf"; got != want {
		t.Fatalf("URI() = %q, want %q", got, want)
	}
	if p.UserNotice() != userNotice {
		t.Fatalf("UserNotice() did not round-trip")
	}
	if p.DocSpecification() != docSpec {
		t.Fatalf("DocSpecification() did not round-trip")
	}
	if p.ValidationResult() != validationResult {
		t.Fatalf("ValidationResult() did not round-trip")
	}
	if got, want := p.TransformsDescription(), []string{}; !reflect.DeepEqual(got, want) {
		t.Fatalf("TransformsDescription() = %v, want empty slice", got)
	}
}
