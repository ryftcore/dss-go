// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/SignatureDigestReference.java (DSS 6.5.RC1).
package signature

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

func TestSignatureDigestReference_RoundTrip(t *testing.T) {
	digest := model.NewDigest(enumerations.DigestAlgorithm_SHA256, []byte{1, 2, 3})
	ref := NewSignatureDigestReferenceWithCanonicalization("http://www.w3.org/2006/12/xml-c14n11", digest)

	if got, want := ref.CanonicalizationMethod(), "http://www.w3.org/2006/12/xml-c14n11"; got != want {
		t.Fatalf("CanonicalizationMethod() = %q, want %q", got, want)
	}
	if got := ref.DigestAlgorithm(); got != enumerations.DigestAlgorithm_SHA256 {
		t.Fatalf("DigestAlgorithm() = %v", got)
	}
	if string(ref.DigestValue()) != string([]byte{1, 2, 3}) {
		t.Fatalf("DigestValue() = %v", ref.DigestValue())
	}
}

func TestSignatureDigestReference_Equals(t *testing.T) {
	digest := model.NewDigest(enumerations.DigestAlgorithm_SHA256, []byte{1, 2, 3})
	a := NewSignatureDigestReference(digest)
	b := NewSignatureDigestReference(digest)
	c := NewSignatureDigestReferenceWithCanonicalization("c14n", digest)

	if !a.Equals(b) {
		t.Fatalf("Equals() = false for equal digests and no canonicalization method")
	}
	if a.Equals(c) {
		t.Fatalf("Equals() = true when canonicalization methods differ")
	}
	if a.Equals(nil) {
		t.Fatalf("Equals(nil) = true")
	}
}
