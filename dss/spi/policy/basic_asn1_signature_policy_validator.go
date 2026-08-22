// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/policy/BasicASN1SignaturePolicyValidator.java (DSS 6.5.RC1).
//
// Deviation: upstream navigates the parsed ASN.1 SEQUENCE via BouncyCastle's
// ASN1Sequence#getObjectAt. This port parses the same TR 102 272 structure
// (SignPolicyHashAlg, SignPolicyInfo, SignPolicyHash, ...) with the standard
// library's encoding/asn1 into a slice of asn1.RawValue, indexed the same
// way; the digest itself is computed by spi.DSSASN1UtilsAsn1SignaturePolicyDigest,
// which owns the byte-exact DER re-encoding used for hashing.
package policy

import (
	"bytes"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/signature"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// BasicASN1SignaturePolicyValidator is the default implementation, checking
// only the hash of the policy.
//
// Note: this implementation is not registered as a service to allow
// overriding.
type BasicASN1SignaturePolicyValidator struct {
	AbstractSignaturePolicyValidator
}

// NewBasicASN1SignaturePolicyValidator is the default constructor.
func NewBasicASN1SignaturePolicyValidator() *BasicASN1SignaturePolicyValidator {
	return &BasicASN1SignaturePolicyValidator{}
}

// CanValidate reports whether the policy content starts with an ASN.1
// SEQUENCE tag.
func (v *BasicASN1SignaturePolicyValidator) CanValidate(signaturePolicy *signature.Policy) bool {
	policyContent := signaturePolicy.PolicyContent()
	if policyContent == nil {
		return false
	}
	firstByte, err := spi.DSSUtilsReadFirstByte(policyContent)
	if err != nil {
		panic(err)
	}
	return spi.DSSASN1UtilsIsASN1SequenceTag(firstByte)
}

// Validate parses the ASN.1 signature policy and compares its declared
// digest algorithm and digest value against the ones carried by
// signaturePolicy.
func (v *BasicASN1SignaturePolicyValidator) Validate(signaturePolicy *signature.Policy) *signature.PolicyValidationResult {
	validationResult := signature.NewPolicyValidationResult()

	policyContent := signaturePolicy.PolicyContent()
	if policyContent == nil {
		validationResult.AddError(AbstractSignaturePolicyValidatorGeneralErrorKey, "The signature policy content is not obtained.")
		return validationResult
	}
	validationResult.SetIdentified(true)

	digest := signaturePolicy.Digest()
	if digest.IsEmpty() {
		validationResult.AddError(AbstractSignaturePolicyValidatorGeneralErrorKey, "The policy digest value is not defined.")
		return validationResult
	}
	// valid if no errors occur
	validationResult.SetDigestValid(true)

	func() {
		// Java catches any Exception raised while parsing/re-digesting and
		// downgrades it to a validation error instead of propagating it;
		// mirrored here with a deferred recover, since the ASN.1 helpers
		// used below panic on malformed input.
		defer func() {
			if r := recover(); r != nil {
				validationResult.SetDigestValid(false)
				validationResult.AddError(AbstractSignaturePolicyValidatorGeneralErrorKey, fmt.Sprint(r))
			}
		}()

		policyBytes, err := spi.DSSUtilsToByteArrayOfDocument(policyContent)
		if err != nil {
			panic(err)
		}

		elements, ok := basicASN1SignaturePolicyValidatorTopLevelElements(policyBytes)
		if !ok || len(elements) < 3 {
			// Not a recognisable TR 102 272 SEQUENCE: Java's
			// ASN1Sequence#getObjectAt would throw, which is caught below;
			// silently returning here (no ASN1Processable flag set) mirrors
			// the effect for a non-conforming top-level object.
			return
		}

		validationResult.SetAsn1Processable(true)

		/*
		 * a) If the resulting document is based on TR 102 272 [i.2] (ESI: ASN.1 format for signature policies),
		 * use the digest value present in the SignPolicyDigest element from the resulting document.
		 * Check that the digest algorithm indicated in the SignPolicyDigestAlg from the resulting
		 * document is equal to the digest algorithm indicated in the property.
		 *
		 * b) If the resulting document is based on TR 102 038 [i.3] ((ESI) XML format for signature policies),
		 * use the digest value present in signPolicyHash element from the resulting document.
		 * Check that the digest algorithm indicated in the signPolicyHashAlg from the resulting
		 * document is equal to the digest algorithm indicated in the attribute.
		 */

		var algID pkix.AlgorithmIdentifier
		if _, err := asn1.Unmarshal(elements[0].FullBytes, &algID); err != nil {
			panic(err)
		}
		signPolicyHashAlgFromPolicy, err := enumerations.DigestAlgorithmForOID(algID.Algorithm.String())
		if err != nil {
			panic(err)
		}

		if signPolicyHashAlgFromPolicy == digest.Algorithm() {
			validationResult.SetDigestAlgorithmsEqual(true)

			recalculatedDigest := v.GetComputedDigest(policyContent, signPolicyHashAlgFromPolicy)
			validationResult.SetDigest(recalculatedDigest)

			equal := digest.Equals(recalculatedDigest)
			validationResult.SetDigestValid(equal)
			if !equal {
				validationResult.AddError(AbstractSignaturePolicyValidatorGeneralErrorKey,
					fmt.Sprintf("The policy digest value (%s) does not match the re-calculated digest value (%s).",
						utils.ToBase64(digest.Value()), utils.ToBase64(recalculatedDigest.Value())))
			}

			policyDigestValueFromPolicy := elements[2].Bytes
			equal = bytes.Equal(digest.Value(), policyDigestValueFromPolicy)
			validationResult.SetDigestValid(equal)
			if !equal {
				validationResult.AddError(AbstractSignaturePolicyValidatorGeneralErrorKey,
					fmt.Sprintf("The policy digest value (%s) does not match the digest value from the policy file (%s).",
						utils.ToBase64(digest.Value()), utils.ToBase64(policyDigestValueFromPolicy)))
			}

		} else {
			validationResult.AddError(AbstractSignaturePolicyValidatorGeneralErrorKey,
				fmt.Sprintf("The digest algorithm indicated in the SignPolicyHashAlg from the resulting document (%s) is not equal to the digest algorithm (%s).",
					signPolicyHashAlgFromPolicy, digest.Algorithm()))
			validationResult.SetDigestAlgorithmsEqual(false)
			validationResult.SetDigestValid(false)
		}
	}()

	return validationResult
}

// GetComputedDigest computes the TR 102 272 ASN.1 signature policy digest,
// shadowing AbstractSignaturePolicyValidator.GetComputedDigest. Panics
// wrapping the underlying error on failure (see that method's doc).
func (v *BasicASN1SignaturePolicyValidator) GetComputedDigest(policyDocument model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) model.Digest {
	policyBytes, err := spi.DSSUtilsToByteArrayOfDocument(policyDocument)
	if err != nil {
		panic(err)
	}
	asn1SignaturePolicyDigest, err := spi.DSSASN1UtilsAsn1SignaturePolicyDigest(digestAlgorithm, policyBytes)
	if err != nil {
		panic(err)
	}
	return model.NewDigest(digestAlgorithm, asn1SignaturePolicyDigest)
}

// basicASN1SignaturePolicyValidatorTopLevelElements parses policyBytes as a
// top-level ASN.1 SEQUENCE, returning its immediate children. Ports the
// implicit cast to ASN1Sequence performed by BouncyCastle when navigating
// asn1Sequence.getObjectAt(...); ok is false when policyBytes is not a
// SEQUENCE.
func basicASN1SignaturePolicyValidatorTopLevelElements(policyBytes []byte) (elements []asn1.RawValue, ok bool) {
	var outer asn1.RawValue
	if _, err := asn1.Unmarshal(policyBytes, &outer); err != nil {
		return nil, false
	}
	if outer.Class != asn1.ClassUniversal || outer.Tag != asn1.TagSequence || !outer.IsCompound {
		return nil, false
	}
	if _, err := asn1.Unmarshal(policyBytes, &elements); err != nil {
		return nil, false
	}
	return elements, true
}

var _ SignaturePolicyValidator = (*BasicASN1SignaturePolicyValidator)(nil)
