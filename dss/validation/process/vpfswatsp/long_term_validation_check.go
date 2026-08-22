// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/LongTermValidationCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
//
// Java's parameter type XmlConstraintsConclusion is the generated base class;
// its Go stand-in is the embedded XmlConstraintsConclusionContent every
// extension type carries (see process/chain.go).
package vpfswatsp

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// LongTermValidationCheck checks if the long-term validation check is
// acceptable.
type LongTermValidationCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationProcessArchivalData]

	// longTermValidationResult is the long-term validation's conclusion.
	longTermValidationResult *jaxb.XmlConstraintsConclusionContent

	// ltvIndication is the LTV Indication.
	ltvIndication enumerations.Indication

	// ltvSubIndication is the LTV SubIndication.
	ltvSubIndication enumerations.SubIndication
}

// NewLongTermValidationCheck is the default constructor. Port of
// LongTermValidationCheck(I18nProvider, XmlValidationProcessArchivalData, XmlConstraintsConclusion, LevelRule).
func NewLongTermValidationCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationProcessArchivalData],
	longTermValidationResult *jaxb.XmlConstraintsConclusionContent,
	constraint policy.LevelRule) *LongTermValidationCheck {
	c := &LongTermValidationCheck{
		ChainItemBase:            process.NewChainItemBase(i18nProvider, result, constraint),
		longTermValidationResult: longTermValidationResult,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *LongTermValidationCheck) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeLTV
}

// Process performs the check. Port of process().
func (c *LongTermValidationCheck) Process() bool {
	if c.longTermValidationResult != nil && c.longTermValidationResult.Conclusion != nil {
		c.ltvIndication = c.longTermValidationResult.Conclusion.Indication.Indication()
		// XmlConclusion#getSubIndication(): the generated member is a pointer,
		// whose nil is Java's null.
		if c.longTermValidationResult.Conclusion.SubIndication != nil {
			c.ltvSubIndication = c.longTermValidationResult.Conclusion.SubIndication.SubIndication()
		} else {
			c.ltvSubIndication = ""
		}

		return process.IsAllowedValidationWithLongTermData(c.longTermValidationResult.Conclusion)
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *LongTermValidationCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagArchLTVV
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *LongTermValidationCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagArchLTVVANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *LongTermValidationCheck) FailedIndicationForConclusion() enumerations.Indication {
	return c.ltvIndication
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *LongTermValidationCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return c.ltvSubIndication
}

// PreviousErrors returns a list of previous errors occurred in the chain. Port
// of getPreviousErrors(), whose Collections.emptyList() is a nil slice here.
func (c *LongTermValidationCheck) PreviousErrors() []*jaxb.XmlMessage {
	if c.longTermValidationResult != nil && c.longTermValidationResult.Conclusion != nil {
		return c.longTermValidationResult.Conclusion.Errors
	}
	return nil
}
