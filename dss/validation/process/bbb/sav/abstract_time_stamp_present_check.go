// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/AbstractTimeStampPresentCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AbstractTimeStampPresentCheckOverrides declares the abstract getTimestamps()
// that a concrete presence check (TLevelTimeStampCheck, LTALevelTimeStampCheck)
// must supply, dispatched separately from process.ChainItemOverrides because
// AbstractTimeStampPresentCheck itself implements Process() by calling it.
type AbstractTimeStampPresentCheckOverrides interface {
	// Timestamps returns a collection of timestamps to be checked for a
	// presence of a valid one. Port of the abstract getTimestamps().
	Timestamps() []*diagnostic.TimestampWrapper
}

// AbstractTimeStampPresentCheck performs analysis if a valid timestamp from the
// given set is present.
type AbstractTimeStampPresentCheck[T any] struct {
	*process.ChainItemBase[T]

	// bbbs is a map of BasicBuildingBlocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// xmlTimestamps is a collection of XmlTimestamps.
	xmlTimestamps []*jaxb.XmlTimestamp

	// overrides points back at the concrete check; see
	// InitAbstractTimeStampPresentCheck.
	overrides AbstractTimeStampPresentCheckOverrides
}

// NewAbstractTimeStampPresentCheck is the default constructor. Port of
// AbstractTimeStampPresentCheck(I18nProvider, T, Map, Collection, LevelRule).
func NewAbstractTimeStampPresentCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, xmlTimestamps []*jaxb.XmlTimestamp,
	constraint policy.LevelRule) *AbstractTimeStampPresentCheck[T] {
	return &AbstractTimeStampPresentCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		bbbs:          bbbs,
		xmlTimestamps: xmlTimestamps,
	}
}

// InitAbstractTimeStampPresentCheck registers the concrete check with its base
// so that the base can dispatch to Timestamps(). Called by the concrete check's
// constructor before InitChainItem.
func (c *AbstractTimeStampPresentCheck[T]) InitAbstractTimeStampPresentCheck(overrides AbstractTimeStampPresentCheckOverrides) {
	c.overrides = overrides
}

// Process performs the check. Port of process().
func (c *AbstractTimeStampPresentCheck[T]) Process() bool {
	for _, timestamp := range c.overrides.Timestamps() {
		timestampBasicValidation := c.getTimestampBasicValidation(timestamp)
		if timestampBasicValidation != nil && process.IsAllowedBasicTimestampValidation(timestampBasicValidation.Conclusion) {
			if c.IsValidConclusion(timestampBasicValidation.Conclusion) {
				return true
			}
			tstPSV := c.getPastSignatureValidationForTimestamp(timestamp)
			if tstPSV != nil && c.IsValidConclusion(tstPSV.Conclusion) {
				return true
			}
		}
	}
	return false
}

// getTimestampBasicValidation ports the private
// getTimestampBasicValidation(TimestampWrapper).
func (c *AbstractTimeStampPresentCheck[T]) getTimestampBasicValidation(timestamp *diagnostic.TimestampWrapper) *jaxb.XmlValidationProcessBasicTimestamp {
	for _, xmlTimestamp := range c.xmlTimestamps {
		if xmlTimestamp.Id != nil && timestamp.Id() == *xmlTimestamp.Id {
			return xmlTimestamp.ValidationProcessBasicTimestamp
		}
	}
	return nil
}

// getPastSignatureValidationForTimestamp ports the private
// getPastSignatureValidationForTimestamp(TimestampWrapper).
func (c *AbstractTimeStampPresentCheck[T]) getPastSignatureValidationForTimestamp(timestampWrapper *diagnostic.TimestampWrapper) *jaxb.XmlPSV {
	tstBBB := c.bbbs[timestampWrapper.Id()]
	if tstBBB != nil {
		return tstBBB.PSV
	}
	return nil
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AbstractTimeStampPresentCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *AbstractTimeStampPresentCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE
}
