// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XAdESLevelXL.java (DSS 6.5.RC1).
//
// Java extends XAdESLevelX and overrides extendSignatures(List); the Go port embeds the -X level
// (ported in the sibling chunk), and "super.extendSignatures(signatures)" is the explicit
// xl.XAdESLevelX.ExtendSignatures call. XAdESLevelA embeds this type in turn.
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// XAdESLevelXL is the XL profile of a XAdES signature.
type XAdESLevelXL struct {
	XAdESLevelX
}

// NewXAdESLevelXL is the default constructor for XAdESLevelXL.
// Port of XAdESLevelXL(CertificateVerifier).
func NewXAdESLevelXL(certificateVerifier validation.CertificateVerifier) *XAdESLevelXL {
	extension := &XAdESLevelXL{}
	extension.InitXAdESLevelXL(extension, certificateVerifier)
	return extension
}

// InitXAdESLevelXL registers the concrete extension level with this base and forwards to the -X
// level. Port of the super(certificateVerifier) call of XAdESLevelXL(CertificateVerifier).
func (xl *XAdESLevelXL) InitXAdESLevelXL(self XAdESSignatureExtensionOverrides,
	certificateVerifier validation.CertificateVerifier) {
	xl.InitXAdESLevelX(self, certificateVerifier)
}

// ExtendSignatures adds CertificateValues and RevocationValues segments to
// UnsignedSignatureProperties. An XML electronic signature MAY contain at most one
// CertificateValues element and at most one RevocationValues element.
// Port of the overridden protected #extendSignatures(List).
//
// NOTE (upstream, kept verbatim): the guard below is computed over the signatures still
// requiring an XL extension, but every subsequent loop and assertion runs over the full
// signatures list, exactly as Java does.
func (xl *XAdESLevelXL) ExtendSignatures(signatures []validation.AdvancedSignature) error {
	if err := xl.XAdESLevelX.ExtendSignatures(signatures); err != nil {
		return err
	}

	signaturesToExtend := xl.extendToXLLevelSignatures(signatures)
	if utils.IsCollectionEmpty(signaturesToExtend) {
		return nil
	}

	for _, signature := range signatures {
		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			// Java's (XAdESSignature) cast; a non-XAdES signature would raise a ClassCastException.
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		if _, err := xl.InitializeSignatureBuilder(xadesSignature); err != nil {
			return err
		}

		// NOTE: do not force sources reload for certificate and revocation sources
		// in order to ensure the same validation data as on -C level
		xl.XadesSignature.ResetTimestampSource()
	}

	signatureRequirementsChecker := xl.SignatureRequirementsChecker()
	if enumerations.SignatureLevelXAdESXL == xl.Params.SignatureLevel() {
		signatureRequirementsChecker.AssertExtendToXLLevelPossible(signatures)
	}
	signatureRequirementsChecker.AssertSignaturesValid(signaturesToExtend)
	signatureRequirementsChecker.AssertCertificateChainValidForXLLevel(signatures)

	// Perform signature validation
	validationDataContainer, err := xl.DocumentAnalyzer.GetValidationData(signatures)
	if err != nil {
		return err
	}

	for _, signature := range signatures {
		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			// Java's (XAdESSignature) cast; a non-XAdES signature would raise a ClassCastException.
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		if _, err := xl.InitializeSignatureBuilder(xadesSignature); err != nil {
			return err
		}
		if signatureRequirementsChecker.HasALevelOrHigher(signature) {
			// Unable to extend due to higher levels covering the current XL-level
			continue
		}

		indent, err := xl.RemoveOldCertificateValues()
		if err != nil {
			return err
		}
		if err := xl.RemoveOldRevocationValues(); err != nil {
			return err
		}

		levelXUnsignedProperties := xl.UnsignedSignaturePropertiesDom.Clone(true)

		validationDataForInclusion := validationDataContainer.AllValidationDataForSignatureForInclusion(signature)

		certificateValuesToAdd := validationDataForInclusion.CertificateTokens()
		crlsToAdd := validationDataForInclusion.CrlTokens()
		ocspsToAdd := validationDataForInclusion.OcspTokens()

		if err := xl.IncorporateCertificateValuesWithIndent(xl.UnsignedSignaturePropertiesDom,
			certificateValuesToAdd, indent); err != nil {
			return err
		}
		if err := xl.IncorporateRevocationValuesWithIndent(xl.UnsignedSignaturePropertiesDom,
			crlsToAdd, ocspsToAdd, indent); err != nil {
			return err
		}

		indented, err := xl.IndentIfPrettyPrint(xl.UnsignedSignaturePropertiesDom, levelXUnsignedProperties)
		if err != nil {
			return err
		}
		xl.UnsignedSignaturePropertiesDom = indented
	}
	return nil
}

// extendToXLLevelSignatures ports the private getExtendToXLLevelSignatures.
func (xl *XAdESLevelXL) extendToXLLevelSignatures(
	signatures []validation.AdvancedSignature) []validation.AdvancedSignature {
	signaturesToExtend := make([]validation.AdvancedSignature, 0)
	for _, signature := range signatures {
		if xl.xlLevelExtensionRequired(signature) {
			signaturesToExtend = append(signaturesToExtend, signature)
		}
	}
	return signaturesToExtend
}

// xlLevelExtensionRequired ports the private xlLevelExtensionRequired.
func (xl *XAdESLevelXL) xlLevelExtensionRequired(signature validation.AdvancedSignature) bool {
	return enumerations.SignatureLevelXAdESXL == xl.Params.SignatureLevel() || !signature.HasAProfile()
}
