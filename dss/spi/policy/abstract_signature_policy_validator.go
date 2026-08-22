// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/policy/AbstractSignaturePolicyValidator.java (DSS 6.5.RC1).
package policy

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// AbstractSignaturePolicyValidatorGeneralErrorKey is the error key to be
// used for general errors. Port of the protected static final
// GENERAL_ERROR_KEY, exported for use by embedding types in this package.
const AbstractSignaturePolicyValidatorGeneralErrorKey = "general"

// AbstractSignaturePolicyValidator is the base implementation of
// SignaturePolicyValidator, providing the default GetComputedDigest.
// Concrete validators embed this type (Go's analogue of Java inheritance)
// and shadow GetComputedDigest when they need a different digest
// computation.
type AbstractSignaturePolicyValidator struct{}

// GetComputedDigest ports
// AbstractSignaturePolicyValidator#getComputedDigest. Panics wrapping a
// *model.DSSError if the digest cannot be computed (e.g. unable to read the
// policyDocument, or an unsupported digestAlgorithm), mirroring Java's
// unchecked DSSException propagating out of a method with no throws clause.
func (v *AbstractSignaturePolicyValidator) GetComputedDigest(policyDocument model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) model.Digest {
	digest, err := spi.DSSUtilsGetDigest(digestAlgorithm, policyDocument)
	if err != nil {
		panic(err)
	}
	return digest
}
