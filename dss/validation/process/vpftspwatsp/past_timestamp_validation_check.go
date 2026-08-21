// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpftspwatsp/checks/PastTimestampValidationCheck.java (DSS 6.5.RC1).
//
// The Java class lives in vpftspwatsp.checks, flattened into this package the
// way the phase 8e package layout does for every ...:checks subpackage.
package vpftspwatsp

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/vpfswatsp"
)

// PastTimestampValidationCheck checks if timestamp's past validation is
// acceptable.
type PastTimestampValidationCheck[T any] struct {
	*vpfswatsp.AbstractPastTokenValidationCheck[T]

	// timestamp is the validated timestamp.
	timestamp *diagnostic.TimestampWrapper
}

// NewPastTimestampValidationCheck is the default constructor. Port of
// PastTimestampValidationCheck(I18nProvider, T, TimestampWrapper, XmlPSV, LevelRule).
func NewPastTimestampValidationCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	timestamp *diagnostic.TimestampWrapper, xmlPSV *jaxb.XmlPSV,
	constraint policy.LevelRule) *PastTimestampValidationCheck[T] {
	c := &PastTimestampValidationCheck[T]{
		AbstractPastTokenValidationCheck: vpfswatsp.NewAbstractPastTokenValidationCheck(
			i18nProvider, result, timestamp, xmlPSV, constraint),
		timestamp: timestamp,
	}
	// Re-register with the outer type so the overridden methods below dispatch
	// correctly.
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *PastTimestampValidationCheck[T]) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockType_TST_PSV
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *PastTimestampValidationCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_PSV_IPTVC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *PastTimestampValidationCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_PSV_IPTVC_ANS
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo(): the unsupported-TimestampType error becomes a panic,
// the method being called from the base ChainItem, which cannot propagate one.
func (c *PastTimestampValidationCheck[T]) BuildAdditionalInfo() *string {
	date := process.GetFormattedDate(c.timestamp.ProductionTime())
	typeTag, err := process.GetTimestampTypeMessageTag(c.timestamp.Type())
	if err != nil {
		panic(err)
	}
	message := c.I18nProvider.GetMessage(i18n.MessageTag_TIMESTAMP_VALIDATION, typeTag, c.timestamp.Id(), date)
	return &message
}
