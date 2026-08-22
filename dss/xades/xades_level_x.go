// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XAdESLevelX.java (DSS 6.5.RC1).
//
// Java extends XAdESLevelC and overrides extendSignatures(List<AdvancedSignature>), calling
// super.extendSignatures first. Go has no method overriding across embedding, so the level chain
// follows the same convention every other level in this package uses: the concrete level embeds
// the previous one, registers itself with the base through InitXAdESLevelX (which forwards to
// InitXAdESLevelC), and ExtendSignatures is reached through the registered overrides. The
// explicit super call becomes a direct call on the embedded LevelC.
//
// Errors: the requirements checks and the message-digest computation return errors here where
// Java throws (PORTING.md: throw -> (T, error)).
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// LevelX represents the implementation of the XAdES level -X extension.
type LevelX struct {
	LevelC
}

// NewXAdESLevelX is the default constructor for LevelX.
// Port of XAdESLevelX(CertificateVerifier).
func NewXAdESLevelX(certificateVerifier validation.CertificateVerifier) *LevelX {
	level := &LevelX{}
	level.InitXAdESLevelX(level, certificateVerifier)
	return level
}

// InitXAdESLevelX registers the concrete level with this base and with LevelC.
func (e *LevelX) InitXAdESLevelX(self SignatureExtensionOverrides,
	certificateVerifier validation.CertificateVerifier) {
	e.InitXAdESLevelC(self, certificateVerifier)
}

// ExtendSignatures adds the xades:SigAndRefsTimeStamp segment to
// xades:UnsignedSignatureProperties. The time-stamp is placed on the digital signature
// (ds:Signature element), the time-stamp(s) present in the XAdES-T form, the certification path
// references and the revocation status references.
//
// A XAdES-X form MAY contain several SigAndRefsTimeStamp elements, obtained from different TSAs.
//
// Port of the overridden protected #extendSignatures(List).
func (e *LevelX) ExtendSignatures(signatures []validation.AdvancedSignature) error {
	if err := e.LevelC.ExtendSignatures(signatures); err != nil {
		return err
	}

	signaturesToExtend := e.extendToXLevelSignatures(signatures)
	if utils.IsCollectionEmpty(signaturesToExtend) {
		return nil
	}

	// document.SignatureRequirementsChecker's assertions have bare returns and panic with error
	// values, so they are called for effect here.
	signatureRequirementsChecker := e.SignatureRequirementsChecker()
	if enumerations.SignatureLevelXAdESX == e.Params.SignatureLevel() {
		signatureRequirementsChecker.AssertExtendToXLevelPossible(signaturesToExtend)
	}
	signatureRequirementsChecker.AssertSignaturesValid(signaturesToExtend)

	for _, signature := range signaturesToExtend {
		xadesSignature, ok := signature.(*Signature)
		if !ok {
			return xadesLevelXUnexpectedSignatureType(signature)
		}
		if _, err := e.InitializeSignatureBuilder(xadesSignature); err != nil {
			return err
		}
		if !e.xLevelExtensionRequired(signature) {
			// Unable to extend due to higher levels covering the current X-level
			continue
		}

		levelCUnsignedProperties := e.UnsignedSignaturePropertiesDom.Clone(true)

		signatureTimestampParameters := e.Params.GetSignatureTimestampParameters()
		digestAlgorithm := signatureTimestampParameters.DigestAlgorithm()
		canonicalizationMethod := signatureTimestampParameters.CanonicalizationMethod()
		// AdvancedSignature.TimestampSource() is invariant in Go, so the covariant Java return
		// (TimestampSource) is recovered by a type assertion - the same shape
		// cades_baseline_requirements_checker.go already uses for TimestampSource.
		timestampSource, ok := e.XadesSignature.TimestampSource().(*TimestampSource)
		if !ok {
			return fmt.Errorf("unexpected timestamp source type %T", e.XadesSignature.TimestampSource())
		}
		messageDigest := timestampSource.GetTimestampX1MessageDigest(digestAlgorithm,
			canonicalizationMethod, e.Params.IsEn319132())
		if err := e.CreateXAdESTimeStampType(enumerations.TimestampTypeValidationDataTimestamp,
			canonicalizationMethod, messageDigest); err != nil {
			return err
		}

		indented, err := e.IndentIfPrettyPrint(e.UnsignedSignaturePropertiesDom, levelCUnsignedProperties)
		if err != nil {
			return err
		}
		e.UnsignedSignaturePropertiesDom = indented
	}
	return nil
}

// extendToXLevelSignatures ports the private getExtendToXLevelSignatures.
func (e *LevelX) extendToXLevelSignatures(
	signatures []validation.AdvancedSignature) []validation.AdvancedSignature {
	signaturesToExtend := make([]validation.AdvancedSignature, 0)
	for _, signature := range signatures {
		if e.xLevelExtensionRequired(signature) {
			signaturesToExtend = append(signaturesToExtend, signature)
		}
	}
	return signaturesToExtend
}

// xLevelExtensionRequired ports the private xLevelExtensionRequired.
func (e *LevelX) xLevelExtensionRequired(signature validation.AdvancedSignature) bool {
	return enumerations.SignatureLevelXAdESX == e.Params.SignatureLevel() || !signature.HasXProfile()
}

// xadesLevelXUnexpectedSignatureType reports a non-XAdES signature reaching this extension, which
// Java's `(XAdESSignature) signature` cast would raise as a ClassCastException.
func xadesLevelXUnexpectedSignatureType(signature validation.AdvancedSignature) error {
	return fmt.Errorf("unexpected signature type %T", signature)
}
