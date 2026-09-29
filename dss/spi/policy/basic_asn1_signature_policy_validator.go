// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/policy/BasicASN1SignaturePolicyValidator.java (DSS 6.5.RC1).
//
// Deviation: upstream navigates the parsed ASN.1 SEQUENCE via BouncyCastle's
// ASN1Sequence#getObjectAt. This port parses the same TR 102 272 structure
// (SignPolicyHashAlg, SignPolicyInfo, SignPolicyHash, ...) with the BER engine
// of internal/asn1ber (which, like BouncyCastle, accepts the indefinite-length
// form) and indexes its children the same way; an index past the end of the
// SEQUENCE panics, standing in for the ArrayIndexOutOfBoundsException that
// getObjectAt throws. The digest itself is computed by
// spi.DSSASN1UtilsAsn1SignaturePolicyDigest, which owns the byte-exact DER
// re-encoding used for hashing.
package policy

import (
	"bytes"
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
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

		// Java: ASN1Sequence asn1Sequence = DSSASN1Utils.toASN1Primitive(policyBytes);
		// a parse failure (DSSException) or a top-level object that is not a
		// SEQUENCE (ClassCastException) is caught below and marks the digest
		// invalid, leaving asn1Processable false.
		elements := basicASN1SignaturePolicyValidatorTopLevelElements(policyBytes)

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

		// Java: (ASN1Sequence) asn1Sequence.getObjectAt(0) followed by
		// AlgorithmIdentifier.getInstance(...): an empty SEQUENCE, a first
		// element that is not a SEQUENCE or a malformed AlgorithmIdentifier
		// throws, and the catch marks the digest invalid.
		algID, err := asn1ber.AlgorithmIdentifierFromElement(basicASN1SignaturePolicyValidatorObjectAt(elements, 0))
		if err != nil {
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

			// Java: ((ASN1OctetString) asn1Sequence.getObjectAt(2)).getOctets():
			// a SEQUENCE with fewer than three elements (or a third element
			// that is not an OCTET STRING) throws here, after the re-calculated
			// digest has been compared, and the catch marks the digest invalid.
			signPolicyHash := basicASN1SignaturePolicyValidatorObjectAt(elements, 2)
			if !signPolicyHash.IsUniversal(asn1ber.TagOctetString) {
				panic("ASN1Encodable cannot be cast to ASN1OctetString")
			}
			policyDigestValueFromPolicy := signPolicyHash.Octets()
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
// single top-level ASN.1 SEQUENCE, returning its immediate children. Ports
// DSSASN1Utils.toASN1Primitive(policyBytes) followed by the implicit cast to
// ASN1Sequence: like Java it panics (DSSException / ClassCastException) when
// policyBytes do not hold exactly one well-formed element, or when that
// element is not a SEQUENCE.
func basicASN1SignaturePolicyValidatorTopLevelElements(policyBytes []byte) []*asn1ber.Element {
	outer, rest, err := asn1ber.Parse(policyBytes)
	if err != nil {
		panic(model.NewDSSErrorMessageCause("Cannot convert binaries to ASN1Primitive", err))
	}
	if len(rest) != 0 {
		panic(model.NewDSSError("Cannot convert binaries to ASN1Primitive : extra data found after the object"))
	}
	if !outer.IsUniversal(asn1ber.TagSequence) || !outer.IsConstructed() {
		panic("ASN1Primitive cannot be cast to ASN1Sequence")
	}
	return outer.Children()
}

// basicASN1SignaturePolicyValidatorObjectAt ports ASN1Sequence#getObjectAt:
// it panics for an index past the end of the SEQUENCE, standing in for Java's
// ArrayIndexOutOfBoundsException.
func basicASN1SignaturePolicyValidatorObjectAt(elements []*asn1ber.Element, index int) *asn1ber.Element {
	if index >= len(elements) {
		panic(fmt.Sprintf("Index %d out of bounds for length %d", index, len(elements)))
	}
	return elements[index]
}

var _ SignaturePolicyValidator = (*BasicASN1SignaturePolicyValidator)(nil)
