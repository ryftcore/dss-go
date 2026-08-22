// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/TimestampValidationCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TimestampValidationCheck checks whether the validation of timestamp with a
// Past Signature Validation process succeed. See EN 319 102-1 ch. "5.6.3
// Validation Process for Signatures providing Long Term Availability and
// Integrity of Validation Material" step 5) of the "5.6.3.4 Processing".
type TimestampValidationCheck[T any] struct {
	*process.ChainItemBase[T]

	// timestamp is the timestamp to check.
	timestamp *diagnostic.TimestampWrapper

	// timestampValidationResult is the timestamp validation result.
	timestampValidationResult *jaxb.XmlValidationProcessArchivalDataTimestamp
}

// NewTimestampValidationCheck is the default constructor. Port of
// TimestampValidationCheck(I18nProvider, T, TimestampWrapper, XmlValidationProcessArchivalDataTimestamp, LevelRule).
func NewTimestampValidationCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	timestamp *diagnostic.TimestampWrapper,
	timestampValidationResult *jaxb.XmlValidationProcessArchivalDataTimestamp,
	constraint policy.LevelRule) *TimestampValidationCheck[T] {
	c := &TimestampValidationCheck[T]{
		ChainItemBase:             process.NewChainItemBaseWithId(i18nProvider, result, constraint, timestamp.Id()),
		timestamp:                 timestamp,
		timestampValidationResult: timestampValidationResult,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *TimestampValidationCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockTypeTST
}

// Process performs the check. Port of process().
func (c *TimestampValidationCheck[T]) Process() bool {
	return c.IsValid(&c.timestampValidationResult.XmlConstraintsConclusionContent)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo(): the unsupported-TimestampType error becomes a panic,
// the method being called from the base ChainItem, which cannot propagate one.
func (c *TimestampValidationCheck[T]) BuildAdditionalInfo() *string {
	date := process.GetFormattedDate(c.timestamp.ProductionTime())
	typeTag, err := process.GetTimestampTypeMessageTag(c.timestamp.Type())
	if err != nil {
		panic(err)
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_TIMESTAMP_VALIDATION, typeTag, c.timestamp.Id(), date)
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TimestampValidationCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_ADEST_IBSVPTADC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *TimestampValidationCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_ADEST_IBSVPTADC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TimestampValidationCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return c.timestampValidationResult.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(): the generated SubIndication
// member is a pointer, whose nil is Java's null.
func (c *TimestampValidationCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.timestampValidationResult.Conclusion.SubIndication == nil {
		return ""
	}
	return c.timestampValidationResult.Conclusion.SubIndication.SubIndication()
}
