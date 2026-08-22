// Ported from
// dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/policy/XMLSignaturePolicyValidator.java
// (DSS 6.5.RC1).
//
// Package layout: eu.europa.esig.dss.xades.validation.policy is grouped, along with
// root/dataobject/dom/evidencerecord/reference/signature/validation/validation.scope/
// validation.timestamp, into one Go package `xades`, hence no `policy` subpackage here.
//
// Santuario replacement: upstream's getDigestAfterTransforms constructs
// org.apache.xml.security.transforms.Transforms directly from the ds:Transforms element and
// calls performTransforms(XMLSignatureInput) on a whole-document input. Per internal/xmldsig's
// doc.go table (PerformTransforms / transforms.Transforms#performTransforms), this is
// xmldsig.PerformTransforms(xmldsig.NewNodeData(document), transformsElement, "", false,
// xmldsig.DefaultRegistry()).
//
// DEVIATION (documented precedent: xades_signature.go's "Deviations" section and
// xades_reference_validation.go's header note the same Go-has-no-virtual-dispatch-through-
// embedding limitation): spi/policy.SignaturePolicyValidator (frozen) declares
// CanValidate/Validate over the concrete base type *signature.SignaturePolicy, not an interface,
// so a caller can never hand this validator a value whose Go runtime type is
// *XAdESSignaturePolicy - only ever the embedded *signature.SignaturePolicy field, which carries
// no trace of the wrapping type. Upstream's `signaturePolicy instanceof XAdESSignaturePolicy`
// check is therefore reproduced as xmlSignaturePolicyValidatorTransforms below, which always
// answers nil through this frozen entry point; the getDigestAfterTransforms path is real and
// exercised by its own callers (were any to hold a concrete *XAdESSignaturePolicy directly), but
// unreachable from Validate() as currently wired. Flagged for the integrator per PORTING.md.
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/signature"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// XMLSignaturePolicyValidatorXMLErrorKey is the error key to be used for XML processing related
// issues. Port of the protected static final XML_ERROR_KEY.
const XMLSignaturePolicyValidatorXMLErrorKey = "xmlProcessing"

// XMLSignaturePolicyValidator validates an XML Signature Policy. Port of the class
// XMLSignaturePolicyValidator, extending spi/policy.AbstractSignaturePolicyValidator.
type XMLSignaturePolicyValidator struct {
	policy.AbstractSignaturePolicyValidator
}

// NewXMLSignaturePolicyValidator is the port of the default constructor.
func NewXMLSignaturePolicyValidator() *XMLSignaturePolicyValidator {
	return &XMLSignaturePolicyValidator{}
}

// CanValidate is the port of canValidate(SignaturePolicy).
func (v *XMLSignaturePolicyValidator) CanValidate(signaturePolicy *signature.SignaturePolicy) bool {
	policyContent := signaturePolicy.PolicyContent()
	if policyContent != nil {
		startsWithXmlPreamble, err := xmlutils.DomUtilsStartsWithXmlPreambleDocument(policyContent)
		if err != nil {
			return false
		}
		return startsWithXmlPreamble
	}
	return false
}

// Validate is the port of validate(SignaturePolicy).
func (v *XMLSignaturePolicyValidator) Validate(signaturePolicy *signature.SignaturePolicy) *signature.SignaturePolicyValidationResult {
	validationResult := signature.NewSignaturePolicyValidationResult()

	policyContent := signaturePolicy.PolicyContent()
	if policyContent == nil {
		validationResult.AddError(policy.AbstractSignaturePolicyValidatorGeneralErrorKey, "The signature policy content is not obtained.")
		return validationResult
	}
	validationResult.SetIdentified(true)

	digest := signaturePolicy.Digest()
	if digest.IsEmpty() {
		validationResult.AddError(policy.AbstractSignaturePolicyValidatorGeneralErrorKey, "The policy digest value is not defined.")
		return validationResult
	}
	validationResult.SetDigestAlgorithmsEqual(true)

	var recalculatedDigest model.Digest
	transforms := xmlSignaturePolicyValidatorTransforms(signaturePolicy)

	if transforms != nil {
		var err error
		recalculatedDigest, err = v.GetDigestAfterTransforms(signaturePolicy.PolicyContent(), digest.Algorithm(), transforms)
		if err != nil {
			errorMessage := fmt.Sprintf("Unable to perform transforms on an XML Policy. Reason : %s", err.Error())
			validationResult.AddError(XMLSignaturePolicyValidatorXMLErrorKey, errorMessage)
			recalculatedDigest = model.Digest{}
		}
	} else {
		recalculatedDigest = v.GetComputedDigest(signaturePolicy.PolicyContent(), digest.Algorithm())
	}
	validationResult.SetDigest(recalculatedDigest)

	if !recalculatedDigest.IsEmpty() {
		if digest.Equals(recalculatedDigest) {
			validationResult.SetDigestValid(true)
		} else {
			validationResult.AddError(policy.AbstractSignaturePolicyValidatorGeneralErrorKey,
				fmt.Sprintf("The policy digest value (%s) does not match the re-calculated digest value (%s).",
					utils.ToBase64(digest.Value()), utils.ToBase64(recalculatedDigest.Value())))
		}
	}

	return validationResult
}

// xmlSignaturePolicyValidatorTransforms ports the `signaturePolicy instanceof
// XAdESSignaturePolicy` check and subsequent xadesSignaturePolicy.getTransforms() call; see this
// file's header for why it always returns nil through the frozen SignaturePolicyValidator entry
// point.
func xmlSignaturePolicyValidatorTransforms(signaturePolicy *signature.SignaturePolicy) *xmldom.Node {
	return nil
}

// GetDigestAfterTransforms computes the Digest after processing of the given transformsElement.
// Port of the public getDigestAfterTransforms(DSSDocument, DigestAlgorithm, Element).
func (v *XMLSignaturePolicyValidator) GetDigestAfterTransforms(policyDocument model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm, transformsElement *xmldom.Node) (model.Digest, error) {
	if transformsElement != nil {
		document, err := xmlutils.DomUtilsBuildDOMFromDocument(policyDocument)
		if err != nil {
			return model.Digest{}, err
		}
		xmlSignatureInput := xmldsig.NewNodeData(document)

		xmlSignatureInputOut, err := xmldsig.PerformTransforms(xmlSignatureInput, transformsElement, "", false, xmldsig.DefaultRegistry())
		if err != nil {
			return model.Digest{}, err
		}
		bytesToBeDigested, err := xmlSignatureInputOut.Bytes()
		if err != nil {
			return model.Digest{}, err
		}
		digestValue, err := spi.DSSUtilsDigest(digestAlgorithm, bytesToBeDigested)
		if err != nil {
			return model.Digest{}, err
		}
		return model.NewDigest(digestAlgorithm, digestValue), nil
	}
	return v.GetComputedDigest(policyDocument, digestAlgorithm), nil
}

// compile-time assertion: an XMLSignaturePolicyValidator satisfies spi/policy.SignaturePolicyValidator.
var _ policy.SignaturePolicyValidator = (*XMLSignaturePolicyValidator)(nil)
