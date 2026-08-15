// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/AbstractTimeStampTypeCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// AbstractTimeStampTypeCheckOverrides declares the abstract getTimestampType()
// that a concrete presence-of-timestamp-type check (ArchiveTimeStampCheck,
// DocumentTimeStampCheck, SignatureTimeStampCheck, ValidationDataTimeStampCheck,
// ValidationDataRefsOnlyTimeStampCheck) must supply, dispatched separately from
// process.ChainItemOverrides because AbstractTimeStampTypeCheck itself
// implements Process() by calling it.
type AbstractTimeStampTypeCheckOverrides interface {
	// TimestampType returns the associated TimestampType to be verified
	// against. Port of the abstract getTimestampType().
	TimestampType() enumerations.TimestampType
}

// AbstractTimeStampTypeCheck verifies a presence of a time-stamp token in a
// signature of the given time-stamp type.
type AbstractTimeStampTypeCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper

	// overrides points back at the concrete check; see InitAbstractTimeStampTypeCheck.
	overrides AbstractTimeStampTypeCheckOverrides
}

// NewAbstractTimeStampTypeCheck is the default constructor. Port of
// AbstractTimeStampTypeCheck(I18nProvider, XmlSAV, SignatureWrapper, LevelRule).
func NewAbstractTimeStampTypeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *AbstractTimeStampTypeCheck {
	return &AbstractTimeStampTypeCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
}

// InitAbstractTimeStampTypeCheck registers the concrete check with its base so
// that the base can dispatch to TimestampType(). Called by the concrete check's
// constructor before InitChainItem.
func (c *AbstractTimeStampTypeCheck) InitAbstractTimeStampTypeCheck(overrides AbstractTimeStampTypeCheckOverrides) {
	c.overrides = overrides
}

// Process performs the check. Port of process().
func (c *AbstractTimeStampTypeCheck) Process() bool {
	for _, timestampWrapper := range c.signature.TimestampList() {
		if c.overrides.TimestampType() == timestampWrapper.Type() {
			return true
		}
	}
	return false
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AbstractTimeStampTypeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *AbstractTimeStampTypeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE
}
