// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfbs/checks/BasicValidationProcessCheck.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.validation.process.vpfbs.checks flattens into this single
// Go package vpfbs (collision-checked with the vpfbs root classes), following
// the same checks-subpackage flattening convention used throughout this port.
package vpfbs

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// BasicValidationProcessCheck verifies if the Basic Signature Validation
// Process succeeds.
type BasicValidationProcessCheck[T any] struct {
	*process.ChainItemBase[T]

	// xmlConclusion is the final check conclusion.
	xmlConclusion *jaxb.XmlConclusion
}

// NewBasicValidationProcessCheck is the default constructor. Port of
// BasicValidationProcessCheck(I18nProvider, T, XmlConclusion, TokenProxy, LevelRule).
func NewBasicValidationProcessCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	xmlConclusion *jaxb.XmlConclusion, token diagnostic.TokenProxy, constraint policy.LevelRule) *BasicValidationProcessCheck[T] {
	c := &BasicValidationProcessCheck[T]{
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, token.Id()),
		xmlConclusion: xmlConclusion,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *BasicValidationProcessCheck[T]) Process() bool {
	return c.IsValidConclusion(c.xmlConclusion)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *BasicValidationProcessCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.xmlConclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *BasicValidationProcessCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.xmlConclusion.SubIndication == nil {
		return ""
	}
	return c.xmlConclusion.SubIndication.SubIndication()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *BasicValidationProcessCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_ADEST_ROBVPIIC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *BasicValidationProcessCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_ADEST_ROBVPIIC_ANS
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo(): Java's String.format("%s/%s", indication, subIndication)
// prints the literal "null" for a null subIndication argument, so a nil
// SubIndication is rendered as the string "null" here rather than being
// dropped, to reproduce the same additional-info text.
func (c *BasicValidationProcessCheck[T]) BuildAdditionalInfo() *string {
	if !c.IsValidConclusion(c.xmlConclusion) {
		subIndicationStr := "null"
		if c.xmlConclusion.SubIndication != nil {
			subIndicationStr = string(c.xmlConclusion.SubIndication.SubIndication())
		}
		indication := fmt.Sprintf("%s/%s", c.xmlConclusion.Indication.Indication(), subIndicationStr)
		message := c.I18nProvider.GetMessage(i18n.MessageTag_BASIC_SIGNATURE_VALIDATION_RESULT, indication)
		return &message
	}
	return nil
}
