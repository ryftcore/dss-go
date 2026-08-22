// Ported from dss-model/.../eaa/ValidationDisclosure.java (DSS 6.5.RC1).
package eaa

import (
	"bytes"
	"reflect"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/eaa/claim"
)

// ValidationDisclosure is the generic base of an EAA Disclosure on
// validation, designed for embedding by concrete presentation-type
// disclosures (mdoc, SD-JWT, ...). Java's protected `salt`/`claim` fields
// are exported (Salt/Claim) so embedders can assign them directly, per
// PORTING.md's "protected fields → exported fields" rule.
//
// Java's abstract `computeDigest(DigestAlgorithm)` method becomes
// the ComputeDigest function field, which embedders MUST set (typically
// from their own constructor) before calling Digest. Likewise,
// Namespace/DigestId are overridable in Java (default nil, mdoc-only);
// NamespaceFunc/DigestIdFunc let an embedder opt into non-default
// behavior; leaving them nil reproduces Java's base-class default.
type ValidationDisclosure struct {
	// Salt is the salt value (protected field in Java).
	Salt []byte

	// Claim is the value of the disclosure claim (protected field in
	// Java).
	Claim claim.Claim

	// digestMap caches computed digest values, keyed by DigestAlgorithm.
	digestMap map[enumerations.DigestAlgorithm]model.Digest

	// ComputeDigest implements the abstract computeDigest(DigestAlgorithm)
	// method: computes digest according to the rules for the given EAA
	// presentation type. Required non-nil to call Digest with a non-empty
	// algorithm.
	ComputeDigest func(digestAlgorithm enumerations.DigestAlgorithm) model.Digest

	// NamespaceFunc optionally overrides the default Namespace() (mdoc
	// only). Leave nil to reproduce Java's default (returns "").
	NamespaceFunc func() string

	// DigestIdFunc optionally overrides the default DigestId() (mdoc
	// only). Leave nil to reproduce Java's default (returns nil).
	DigestIdFunc func() *int64
}

// NewValidationDisclosure ports the default (protected, no-arg)
// constructor.
func NewValidationDisclosure() *ValidationDisclosure {
	return &ValidationDisclosure{}
}

// Name gets the name of the disclosure claim. Ports
// ValidationDisclosure#getName.
func (v *ValidationDisclosure) Name() string {
	if v.Claim != nil {
		return v.Claim.Name()
	}
	return ""
}

// Namespace gets the applicable namespace (mdoc only). Ports
// ValidationDisclosure#getNamespace.
func (v *ValidationDisclosure) Namespace() string {
	if v.NamespaceFunc != nil {
		return v.NamespaceFunc()
	}
	return ""
}

// DigestId gets the digest ID for the issuer data authentication (mdoc
// only). Ports ValidationDisclosure#getDigestId.
func (v *ValidationDisclosure) DigestId() *int64 {
	if v.DigestIdFunc != nil {
		return v.DigestIdFunc()
	}
	return nil
}

// Digest gets the digest value for digestAlgorithm, computing and caching
// it via ComputeDigest on first request. Ports
// ValidationDisclosure#getDigest. Passing "" (Java null) returns an empty
// Digest without invoking ComputeDigest or touching the cache.
func (v *ValidationDisclosure) Digest(digestAlgorithm enumerations.DigestAlgorithm) model.Digest {
	if digestAlgorithm == "" {
		return model.Digest{}
	}
	if d, ok := v.digestMap[digestAlgorithm]; ok {
		return d
	}
	d := v.ComputeDigest(digestAlgorithm)
	if v.digestMap == nil {
		v.digestMap = make(map[enumerations.DigestAlgorithm]model.Digest)
	}
	v.digestMap[digestAlgorithm] = d
	return d
}

// Equals ports ValidationDisclosure#equals. NOTE: uses bytes.Equal for
// Salt, which (unlike Java's Arrays.equals) treats a nil and an empty
// slice as equal; this is an accepted minor deviation for a field that
// upstream never observably sets to a non-nil empty array.
func (v *ValidationDisclosure) Equals(other *ValidationDisclosure) bool {
	if v == other {
		return true
	}
	if other == nil {
		return false
	}
	return bytes.Equal(v.Salt, other.Salt) && reflect.DeepEqual(v.Claim, other.Claim)
}
