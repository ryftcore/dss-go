// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/checks/EAAQualifiedSignatureOrSealCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAAQualifiedSignatureOrSealCheck verifies whether the EAA has been
// created with a qualified electronic signature or seal.
type EAAQualifiedSignatureOrSealCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationEAAQualificationProcess]

	// signature is the signature to be checked.
	signature *diagnostic.SignatureWrapper

	// signatureQualification is the Signature Qualification to be checked.
	signatureQualification enumerations.SignatureQualification
}

// NewEAAQualifiedSignatureOrSealCheck is the default constructor. Port of
// EAAQualifiedSignatureOrSealCheck(I18nProvider, XmlValidationEAAQualificationProcess, SignatureWrapper, SignatureQualification, LevelRule).
func NewEAAQualifiedSignatureOrSealCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationEAAQualificationProcess],
	signature *diagnostic.SignatureWrapper, signatureQualification enumerations.SignatureQualification,
	constraint policy.LevelRule) *EAAQualifiedSignatureOrSealCheck {
	c := &EAAQualifiedSignatureOrSealCheck{
		ChainItemBase:          process.NewChainItemBase(i18nProvider, result, constraint),
		signature:              signature,
		signatureQualification: signatureQualification,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAQualifiedSignatureOrSealCheck) Process() bool {
	// Indeterminate statuses are handled separately
	return enumerations.SignatureQualificationQESig == c.signatureQualification ||
		enumerations.SignatureQualificationQESeal == c.signatureQualification ||
		enumerations.SignatureQualificationIndeterminateQESig == c.signatureQualification ||
		enumerations.SignatureQualificationIndeterminateQESeal == c.signatureQualification
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAQualifiedSignatureOrSealCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAASigQual
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *EAAQualifiedSignatureOrSealCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAASigQualANS
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *EAAQualifiedSignatureOrSealCheck) BuildErrorMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(c.ErrorMessageTag(), c.signatureQualification.Readable())
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *EAAQualifiedSignatureOrSealCheck) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTagSignatureID, c.signature.Id())
	return &message
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAQualifiedSignatureOrSealCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *EAAQualifiedSignatureOrSealCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
