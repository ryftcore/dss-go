// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/policy/NonASN1SignaturePolicyValidator.java (DSS 6.5.RC1).
package policy

import (
	"fmt"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/signature"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
)

// nonASN1SignaturePolicyValidatorXMLBOM is the UTF-8 BOM ([]byte{0xEF, 0xBB,
// 0xBF}) followed by '<', used to detect an XML document (with a byte-order
// mark) that would otherwise be mistaken for a non-ASN1/non-XML policy.
var nonASN1SignaturePolicyValidatorXMLBOM = []byte{0xEF, 0xBB, 0xBF, '<'}

// NonASN1SignaturePolicyValidator covers the case of non-ASN1 signature
// policies (e.g. a PDF file and its digest).
type NonASN1SignaturePolicyValidator struct {
	AbstractSignaturePolicyValidator
}

// NewNonASN1SignaturePolicyValidator is the default constructor.
func NewNonASN1SignaturePolicyValidator() *NonASN1SignaturePolicyValidator {
	return &NonASN1SignaturePolicyValidator{}
}

// CanValidate reports whether the policy content is neither ASN.1 nor XML.
func (v *NonASN1SignaturePolicyValidator) CanValidate(signaturePolicy *signature.SignaturePolicy) bool {
	policyContent := signaturePolicy.PolicyContent()
	if policyContent == nil {
		return false
	}
	firstByte, err := spi.DSSUtilsReadFirstByte(policyContent)
	if err != nil {
		panic(err)
	}
	startsWithBom, err := nonASN1SignaturePolicyValidatorStartsWithXMLBOM(policyContent)
	if err != nil {
		panic(model.NewDSSErrorMessageCause("Cannot read the InputStream!", err))
	}
	// non ASN1 nor XML policy
	return !spi.DSSASN1UtilsIsASN1SequenceTag(firstByte) && firstByte != '<' && !startsWithBom
}

// nonASN1SignaturePolicyValidatorStartsWithXMLBOM reports whether
// policyContent begins with a UTF-8 BOM followed by '<'.
func nonASN1SignaturePolicyValidatorStartsWithXMLBOM(policyContent model.DSSDocument) (bool, error) {
	stream, err := policyContent.OpenStream()
	if err != nil {
		return false, err
	}
	defer utils.CloseQuietly(stream)
	return utils.StartsWithStream(stream, nonASN1SignaturePolicyValidatorXMLBOM)
}

// Validate recomputes the digest of the policy content and compares it
// against the declared digest.
func (v *NonASN1SignaturePolicyValidator) Validate(signaturePolicy *signature.SignaturePolicy) *signature.SignaturePolicyValidationResult {
	validationResult := signature.NewSignaturePolicyValidationResult()

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
	validationResult.SetDigestAlgorithmsEqual(true)

	recalculatedDigest := v.GetComputedDigest(signaturePolicy.PolicyContent(), digest.Algorithm())
	validationResult.SetDigest(recalculatedDigest)

	if digest.Equals(recalculatedDigest) {
		validationResult.SetDigestValid(true)
		validationResult.SetDigestAlgorithmsEqual(true)
	} else {
		validationResult.AddError(AbstractSignaturePolicyValidatorGeneralErrorKey,
			fmt.Sprintf("The policy digest value (%s) does not match the re-calculated digest value (%s).",
				utils.ToBase64(digest.Value()), utils.ToBase64(recalculatedDigest.Value())))
	}

	return validationResult
}

var _ SignaturePolicyValidator = (*NonASN1SignaturePolicyValidator)(nil)
