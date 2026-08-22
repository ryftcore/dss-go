// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/SignaturePolicyValidationResult.java (DSS 6.5.RC1).
package signature

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

func TestSignaturePolicyValidationResult_RoundTrip(t *testing.T) {
	r := NewSignaturePolicyValidationResult()
	r.SetIdentified(true)
	r.SetAsn1Processable(true)
	r.SetDigestAlgorithmsEqual(true)
	r.SetDigestValid(true)
	digest := model.NewDigest(enumerations.DigestAlgorithm_SHA256, []byte{1, 2, 3})
	r.SetDigest(digest)

	if !r.IsIdentified() || !r.IsAsn1Processable() || !r.IsDigestAlgorithmsEqual() || !r.IsDigestValid() {
		t.Fatalf("boolean flags did not round-trip")
	}
	if !r.Digest().Equals(digest) {
		t.Fatalf("Digest() = %v, want %v", r.Digest(), digest)
	}
	if r.ProcessingErrors() != "" {
		t.Fatalf("ProcessingErrors() = %q, want empty string with no errors", r.ProcessingErrors())
	}
}

func TestSignaturePolicyValidationResult_ProcessingErrors(t *testing.T) {
	r := NewSignaturePolicyValidationResult()
	r.AddError("step1", "digest mismatch")
	r.AddError("step2", "unknown identifier")

	want := "The errors found on signature policy validation are: at step1: digest mismatch, at step2: unknown identifier"
	if got := r.ProcessingErrors(); got != want {
		t.Fatalf("ProcessingErrors() = %q, want %q", got, want)
	}
}
