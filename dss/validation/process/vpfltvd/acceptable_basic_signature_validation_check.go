// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfltvd/checks/AcceptableBasicSignatureValidationCheck.java (DSS 6.5.RC1).
package vpfltvd

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// AcceptableBasicSignatureValidationCheck checks if the signature's basic
// validation result is acceptable.
type AcceptableBasicSignatureValidationCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationProcessLongTermData]

	// content is the signature's basic validation conclusion: the embedded
	// content struct of the caller's concrete result object - the Go stand-in
	// for Java's XmlConstraintsConclusion supertype reference (see chain.go's
	// Result type), since the concrete result type varies by caller
	// (XmlValidationProcessBasicSignature here).
	content *jaxb.XmlConstraintsConclusionContent

	// bbbIndication is the validation Indication.
	bbbIndication enumerations.Indication

	// bbbSubIndication is the validation SubIndication.
	bbbSubIndication enumerations.SubIndication
}

// NewAcceptableBasicSignatureValidationCheck is the default constructor. Port
// of
// AcceptableBasicSignatureValidationCheck(I18nProvider, XmlValidationProcessLongTermData, XmlConstraintsConclusion, LevelRule).
//
// basicSignatureValidation is passed as the embedded
// XmlConstraintsConclusionContent directly (the Go stand-in for Java's
// XmlConstraintsConclusion supertype reference - see chain.go's Result type),
// since the concrete result type varies by caller (XmlValidationProcessBasicSignature
// here).
func NewAcceptableBasicSignatureValidationCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationProcessLongTermData],
	basicSignatureValidation *jaxb.XmlConstraintsConclusionContent,
	constraint policy.LevelRule) *AcceptableBasicSignatureValidationCheck {
	c := &AcceptableBasicSignatureValidationCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		content:       basicSignatureValidation,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AcceptableBasicSignatureValidationCheck) Process() bool {
	if c.content != nil && c.content.Conclusion != nil {
		basicSignatureConclusion := c.content.Conclusion
		c.bbbIndication = basicSignatureConclusion.Indication.Indication()
		if basicSignatureConclusion.SubIndication != nil {
			c.bbbSubIndication = basicSignatureConclusion.SubIndication.SubIndication()
		}
		return process.IsAllowedBasicSignatureValidation(basicSignatureConclusion)
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AcceptableBasicSignatureValidationCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_LTV_ABSV
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *AcceptableBasicSignatureValidationCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_LTV_ABSV_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AcceptableBasicSignatureValidationCheck) FailedIndicationForConclusion() enumerations.Indication {
	return c.bbbIndication
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *AcceptableBasicSignatureValidationCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return c.bbbSubIndication
}

// PreviousErrors returns a list of previous errors occurred in the chain.
// Port of getPreviousErrors().
func (c *AcceptableBasicSignatureValidationCheck) PreviousErrors() []*jaxb.XmlMessage {
	if c.content != nil && c.content.Conclusion != nil {
		return c.content.Conclusion.Errors
	}
	return nil
}
