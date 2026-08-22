package model

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// TestDigestDocumentExistingDigestDeterministic guards against defect #2 of the Phase 2b audit
// ("Nondeterministic output ordering"): ExistingDigest() picks an arbitrary entry out of
// digestMap; ranging a bare Go map for a "first" pick is randomized on every run, unlike
// Java's HashMap.entrySet().iterator().next() (arbitrary but stable within a JVM run), so
// digestMap is now insertion-ordered and ExistingDigest() picks the first-added algorithm.
func TestDigestDocumentExistingDigestDeterministic(t *testing.T) {
	build := func() enumerations.DigestAlgorithm {
		d := NewDigestDocument()
		for _, alg := range []enumerations.DigestAlgorithm{
			enumerations.DigestAlgorithmSHA512,
			enumerations.DigestAlgorithmSHA1,
			enumerations.DigestAlgorithmSHA384,
			enumerations.DigestAlgorithmSHA3256,
			enumerations.DigestAlgorithmSHA256,
		} {
			d.AddDigestValue(alg, []byte("digest-for-"+string(alg)))
		}
		digest, err := d.ExistingDigest()
		if err != nil {
			t.Fatalf("ExistingDigest: %v", err)
		}
		return digest.Algorithm()
	}
	want := build()
	if want != enumerations.DigestAlgorithmSHA512 {
		t.Fatalf("ExistingDigest() algorithm = %s, want the first-added algorithm SHA512", want)
	}
	for i := 0; i < 25; i++ {
		got := build()
		if got != want {
			t.Fatalf("run %d: ExistingDigest() algorithm = %s, want %s", i, got, want)
		}
	}
}
