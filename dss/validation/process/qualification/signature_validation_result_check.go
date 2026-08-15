// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/checks/SignatureValidationResultCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SignatureValidationResultCheck performs verification of the Basic
// Signature validation process result.
type SignatureValidationResultCheck[T any] struct {
	*process.ChainItemBase[T]

	// signatureBasicValidationConclusion is the Basic Validation conclusion
	// of the signature.
	signatureBasicValidationConclusion *jaxb.XmlConclusion
}

// NewSignatureValidationResultCheck is the default constructor. Port of
// SignatureValidationResultCheck(I18nProvider, T, XmlConclusion, LevelRule).
func NewSignatureValidationResultCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	basicValidationConclusion *jaxb.XmlConclusion, constraint policy.LevelRule) *SignatureValidationResultCheck[T] {
	c := &SignatureValidationResultCheck[T]{
		ChainItemBase:                      process.NewChainItemBase(i18nProvider, result, constraint),
		signatureBasicValidationConclusion: basicValidationConclusion,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SignatureValidationResultCheck[T]) Process() bool {
	return c.IsValidConclusion(c.signatureBasicValidationConclusion)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SignatureValidationResultCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_ADEST_IBSVPSC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *SignatureValidationResultCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_ADEST_IBSVPSC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SignatureValidationResultCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	indication := c.signatureBasicValidationConclusion.Indication.Indication()
	if enumerations.Indication_TOTAL_FAILED == indication {
		return enumerations.Indication_FAILED
	}
	return indication
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *SignatureValidationResultCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.signatureBasicValidationConclusion.SubIndication == nil {
		return ""
	}
	return c.signatureBasicValidationConclusion.SubIndication.SubIndication()
}
