// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/LongTermAvailabilityAndIntegrityValidationMaterialCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note, and long_term_validation_check.go
// for the XmlConstraintsConclusion stand-in.
package vpfswatsp

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// LongTermAvailabilityAndIntegrityValidationMaterialCheck verifies whether the
// term availability and integrity of validation material is present within the
// signature.
type LongTermAvailabilityAndIntegrityValidationMaterialCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationProcessArchivalData]

	// signature is the signature to be verified.
	signature *diagnostic.SignatureWrapper

	// longTermValidationResult is the long-term validation's conclusion.
	longTermValidationResult *jaxb.XmlConstraintsConclusionContent

	// ltvIndication is the LTV Indication.
	ltvIndication enumerations.Indication

	// ltvSubIndication is the LTV SubIndication.
	ltvSubIndication enumerations.SubIndication
}

// NewLongTermAvailabilityAndIntegrityValidationMaterialCheck is the default
// constructor. Port of
// LongTermAvailabilityAndIntegrityValidationMaterialCheck(I18nProvider, XmlValidationProcessArchivalData, SignatureWrapper, XmlConstraintsConclusion, LevelRule).
func NewLongTermAvailabilityAndIntegrityValidationMaterialCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationProcessArchivalData], signature *diagnostic.SignatureWrapper,
	longTermValidationResult *jaxb.XmlConstraintsConclusionContent,
	constraint policy.LevelRule) *LongTermAvailabilityAndIntegrityValidationMaterialCheck {
	c := &LongTermAvailabilityAndIntegrityValidationMaterialCheck{
		ChainItemBase:            process.NewChainItemBase(i18nProvider, result, constraint),
		signature:                signature,
		longTermValidationResult: longTermValidationResult,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *LongTermAvailabilityAndIntegrityValidationMaterialCheck) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeLTA
}

// Process performs the check. Port of process().
func (c *LongTermAvailabilityAndIntegrityValidationMaterialCheck) Process() bool {
	if c.longTermValidationResult != nil && c.longTermValidationResult.Conclusion != nil {
		c.ltvIndication = c.longTermValidationResult.Conclusion.Indication.Indication()
		// XmlConclusion#getSubIndication(): the generated member is a pointer,
		// whose nil is Java's null.
		if c.longTermValidationResult.Conclusion.SubIndication != nil {
			c.ltvSubIndication = c.longTermValidationResult.Conclusion.SubIndication.SubIndication()
		} else {
			c.ltvSubIndication = ""
		}
	}
	return process.IsLongTermAvailabilityAndIntegrityMaterialPresent(c.signature)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *LongTermAvailabilityAndIntegrityValidationMaterialCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_ARCH_LTAIVMP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *LongTermAvailabilityAndIntegrityValidationMaterialCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_ARCH_LTAIVMP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *LongTermAvailabilityAndIntegrityValidationMaterialCheck) FailedIndicationForConclusion() enumerations.Indication {
	return c.ltvIndication
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *LongTermAvailabilityAndIntegrityValidationMaterialCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return c.ltvSubIndication
}

// PreviousErrors returns a list of previous errors occurred in the chain. Port
// of getPreviousErrors(), whose Collections.emptyList() is a nil slice here.
//
// Java copies the conclusion's errors into a fresh ArrayList before appending
// its own error message, leaving the conclusion untouched; the Go copy is
// explicit for the same reason.
func (c *LongTermAvailabilityAndIntegrityValidationMaterialCheck) PreviousErrors() []*jaxb.XmlMessage {
	if c.longTermValidationResult != nil && c.longTermValidationResult.Conclusion != nil {
		conclusionErrors := c.longTermValidationResult.Conclusion.Errors
		errors := make([]*jaxb.XmlMessage, len(conclusionErrors))
		copy(errors, conclusionErrors)
		if utils.IsCollectionNotEmpty(errors) {
			errors = append(errors, c.BuildErrorMessage())
		}
		return errors
	}
	return nil
}
