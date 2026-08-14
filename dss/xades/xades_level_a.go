// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XAdESLevelA.java (DSS 6.5.RC1).
//
// Java extends XAdESLevelXL and overrides extendSignatures(List); the Go port embeds the -XL
// level, and "super.extendSignatures(signatures)" is the explicit
// a.XAdESLevelXL.ExtendSignatures call. This is the top of the legacy (non-baseline) XAdES
// augmentation ladder.
package xades

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/spi/validation"
)

// XAdESLevelA holds the level A aspects of XAdES.
type XAdESLevelA struct {
	XAdESLevelXL
}

// NewXAdESLevelA is the default constructor for XAdESLevelA.
// Port of XAdESLevelA(CertificateVerifier).
func NewXAdESLevelA(certificateVerifier validation.CertificateVerifier) *XAdESLevelA {
	extension := &XAdESLevelA{}
	extension.InitXAdESLevelA(extension, certificateVerifier)
	return extension
}

// InitXAdESLevelA registers the concrete extension level with this base and forwards to the -XL
// level. Port of the super(certificateVerifier) call of XAdESLevelA(CertificateVerifier).
func (a *XAdESLevelA) InitXAdESLevelA(self XAdESSignatureExtensionOverrides,
	certificateVerifier validation.CertificateVerifier) {
	a.InitXAdESLevelXL(self, certificateVerifier)
}

// ExtendSignatures adds the ArchiveTimeStamp element which is an unsigned property
// qualifying the signature. The hash sent to the TSA (messageImprint) is computed on the
// XAdES-X-L form of the electronic signature and the signed data objects. A XAdES-A form MAY
// contain several ArchiveTimeStamp elements.
// Port of the overridden protected #extendSignatures(List).
func (a *XAdESLevelA) ExtendSignatures(signatures []validation.AdvancedSignature) error {
	if err := a.XAdESLevelXL.ExtendSignatures(signatures); err != nil {
		return err
	}

	signatureRequirementsChecker := a.SignatureRequirementsChecker()
	signatureRequirementsChecker.AssertSignaturesValid(signatures)

	addTimestampValidationData := false

	for _, signature := range signatures {
		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			// Java's (XAdESSignature) cast; a non-XAdES signature would raise a ClassCastException.
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		if _, err := a.InitializeSignatureBuilder(xadesSignature); err != nil {
			return err
		}
		if err := a.assertExtendSignatureToAPossible(); err != nil {
			return err
		}

		if a.XadesSignature.HasLTAProfile() {
			addTimestampValidationData = true
		}
	}

	// Perform signature validation
	var validationDataContainer *validation.ValidationDataContainer
	if addTimestampValidationData {
		container, err := a.DocumentAnalyzer.GetValidationData(signatures)
		if err != nil {
			return err
		}
		validationDataContainer = container
	}

	// Append LTA-level (+ ValidationData)
	for _, signature := range signatures {
		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			// Java's (XAdESSignature) cast; a non-XAdES signature would raise a ClassCastException.
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		if _, err := a.InitializeSignatureBuilder(xadesSignature); err != nil {
			return err
		}
		levelXLUnsignedProperties := a.UnsignedSignaturePropertiesDom.Clone(true)

		if a.XadesSignature.HasLTAProfile() && addTimestampValidationData {
			// must be executed before data removing
			indent, err := a.RemoveLastTimestampAndAnyValidationData()
			if err != nil {
				return err
			}

			validationDataForInclusion := validationDataContainer.AllValidationDataForSignatureForInclusion(signature)
			if err := a.IncorporateTimestampValidationData(validationDataForInclusion, indent); err != nil {
				return err
			}
		}
		if err := a.IncorporateArchiveTimestamp(); err != nil {
			return err
		}

		indented, err := a.IndentIfPrettyPrint(a.UnsignedSignaturePropertiesDom, levelXLUnsignedProperties)
		if err != nil {
			return err
		}
		a.UnsignedSignaturePropertiesDom = indented
	}
	return nil
}

// assertExtendSignatureToAPossible ports the private assertExtendSignatureToAPossible.
func (a *XAdESLevelA) assertExtendSignatureToAPossible() error {
	if enumerations.SignatureLevel_XAdES_A == a.Params.SignatureLevel() {
		return a.AssertDetachedDocumentsContainBinaries()
	}
	return nil
}
