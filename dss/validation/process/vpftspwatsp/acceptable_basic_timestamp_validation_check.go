// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpftspwatsp/checks/AcceptableBasicTimestampValidationCheck.java (DSS 6.5.RC1).
//
// Java's parameter type XmlConstraintsConclusion is the generated base class;
// its Go stand-in is the embedded XmlConstraintsConclusionContent every
// extension type carries (see process/chain.go).
package vpftspwatsp

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AcceptableBasicTimestampValidationCheck checks if a result of a Basic
// Signature Validation process for a timestamp token is acceptable.
type AcceptableBasicTimestampValidationCheck[T any] struct {
	*process.ChainItemBase[T]

	// basicTimestampValidation is the signature's basic validation conclusion.
	basicTimestampValidation *jaxb.XmlConstraintsConclusionContent

	// bbbIndication is the validation Indication.
	bbbIndication enumerations.Indication

	// bbbSubIndication is the validation SubIndication.
	bbbSubIndication enumerations.SubIndication
}

// NewAcceptableBasicTimestampValidationCheck is the default constructor. Port
// of AcceptableBasicTimestampValidationCheck(I18nProvider, T, XmlConstraintsConclusion, LevelRule).
func NewAcceptableBasicTimestampValidationCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	basicTimestampValidation *jaxb.XmlConstraintsConclusionContent,
	constraint policy.LevelRule) *AcceptableBasicTimestampValidationCheck[T] {
	c := &AcceptableBasicTimestampValidationCheck[T]{
		ChainItemBase:            process.NewChainItemBase(i18nProvider, result, constraint),
		basicTimestampValidation: basicTimestampValidation,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *AcceptableBasicTimestampValidationCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockType_TST_BBB
}

// Process performs the check. Port of process().
func (c *AcceptableBasicTimestampValidationCheck[T]) Process() bool {
	if c.basicTimestampValidation != nil && c.basicTimestampValidation.Conclusion != nil {
		conclusion := c.basicTimestampValidation.Conclusion
		c.bbbIndication = conclusion.Indication.Indication()
		// XmlConclusion#getSubIndication(): the generated member is a pointer,
		// whose nil is Java's null.
		if conclusion.SubIndication != nil {
			c.bbbSubIndication = conclusion.SubIndication.SubIndication()
		} else {
			c.bbbSubIndication = ""
		}

		return process.IsAllowedBasicTimestampValidation(conclusion)
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AcceptableBasicTimestampValidationCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_ARCH_IRTVBBA
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *AcceptableBasicTimestampValidationCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_ARCH_IRTVBBA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AcceptableBasicTimestampValidationCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.bbbIndication
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *AcceptableBasicTimestampValidationCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return c.bbbSubIndication
}

// PreviousErrors returns a list of previous errors occurred in the chain. Port
// of getPreviousErrors(), whose Collections.emptyList() is a nil slice here.
func (c *AcceptableBasicTimestampValidationCheck[T]) PreviousErrors() []*jaxb.XmlMessage {
	if c.basicTimestampValidation != nil && c.basicTimestampValidation.Conclusion != nil {
		return c.basicTimestampValidation.Conclusion.Errors
	}
	return nil
}
