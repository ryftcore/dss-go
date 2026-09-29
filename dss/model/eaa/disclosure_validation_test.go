package eaa

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/eaa/claim"
)

// newTestDisclosure builds a disclosure carrying its own ComputeDigest closure, as a concrete
// presentation-type embedder would.
func newTestDisclosure(salt []byte, name, value string) *ValidationDisclosure {
	d := NewValidationDisclosure()
	d.Salt = salt
	d.Claim = claim.NewStringWithName(name, value)
	d.ComputeDigest = func(digestAlgorithm enumerations.DigestAlgorithm) model.Digest {
		return model.NewDigest(digestAlgorithm, []byte(value))
	}
	return d
}

// TestDisclosureValidationEqualsIgnoresComputeDigest pins ValidationDisclosure#equals semantics
// (salt and claim only) inside DisclosureValidation#equals: the ComputeDigest closure and the
// digest cache are not part of a disclosure's identity, so two disclosures equal upstream stay
// equal here even though their function fields differ.
func TestDisclosureValidationEqualsIgnoresComputeDigest(t *testing.T) {
	a := NewDisclosureValidationWithDisclosure(newTestDisclosure([]byte{1, 2}, "family_name", "Doe"))
	b := NewDisclosureValidationWithDisclosure(newTestDisclosure([]byte{1, 2}, "family_name", "Doe"))

	// Populate one side's digest cache so the two disclosures differ in unexported state too.
	a.Disclosure().Digest(enumerations.DigestAlgorithmSHA256)

	if !a.Equals(b) {
		t.Fatal("DisclosureValidations over disclosures with the same salt and claim must be equal")
	}
	if !a.Disclosure().Equals(b.Disclosure()) {
		t.Fatal("ValidationDisclosures with the same salt and claim must be equal")
	}
}

func TestDisclosureValidationEqualsDetectsDifferences(t *testing.T) {
	base := NewDisclosureValidationWithDisclosure(newTestDisclosure([]byte{1, 2}, "family_name", "Doe"))

	otherSalt := NewDisclosureValidationWithDisclosure(newTestDisclosure([]byte{9, 9}, "family_name", "Doe"))
	if base.Equals(otherSalt) {
		t.Fatal("a different salt must make the DisclosureValidations unequal")
	}
	otherClaim := NewDisclosureValidationWithDisclosure(newTestDisclosure([]byte{1, 2}, "family_name", "Roe"))
	if base.Equals(otherClaim) {
		t.Fatal("a different claim value must make the DisclosureValidations unequal")
	}

	withNamespace := NewDisclosureValidationWithDisclosure(newTestDisclosure([]byte{1, 2}, "family_name", "Doe"))
	withNamespace.SetNamespace("org.iso.18013.5.1")
	if base.Equals(withNamespace) {
		t.Fatal("a different namespace must make the DisclosureValidations unequal")
	}

	digestID := int64(7)
	withDigestID := NewDisclosureValidationWithDisclosure(newTestDisclosure([]byte{1, 2}, "family_name", "Doe"))
	withDigestID.SetDigestId(&digestID)
	if base.Equals(withDigestID) {
		t.Fatal("a different digest id must make the DisclosureValidations unequal")
	}
	sameDigestID := int64(7)
	withSameDigestID := NewDisclosureValidationWithDisclosure(newTestDisclosure([]byte{1, 2}, "family_name", "Doe"))
	withSameDigestID.SetDigestId(&sameDigestID)
	if !withDigestID.Equals(withSameDigestID) {
		t.Fatal("equal digest ids held by distinct pointers must compare equal")
	}
}

func TestDisclosureValidationEqualsNilDisclosure(t *testing.T) {
	none1 := NewDisclosureValidation()
	none2 := NewDisclosureValidation()
	if !none1.Equals(none2) {
		t.Fatal("two DisclosureValidations without a disclosure must be equal")
	}
	with := NewDisclosureValidationWithDisclosure(newTestDisclosure(nil, "n", "v"))
	if none1.Equals(with) || with.Equals(none1) {
		t.Fatal("a missing disclosure must not equal a present one")
	}
}
