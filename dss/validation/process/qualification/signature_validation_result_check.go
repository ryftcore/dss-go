// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/checks/SignatureValidationResultCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
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
// SignatureValidationResultCheck(Provider, T, XmlConclusion, LevelRule).
func NewSignatureValidationResultCheck[T any](i18nProvider *i18n.Provider, result *process.Result[T],
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
	return i18n.MessageTagADESTIBSVPSC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *SignatureValidationResultCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagADESTIBSVPSCANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SignatureValidationResultCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	indication := c.signatureBasicValidationConclusion.Indication.Indication()
	if enumerations.IndicationTotalFailed == indication {
		return enumerations.IndicationFailed
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
