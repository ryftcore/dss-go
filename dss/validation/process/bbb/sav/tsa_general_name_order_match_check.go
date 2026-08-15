// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/TSAGeneralNameOrderMatchCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// TSAGeneralNameOrderMatchCheck checks if the TSTInfo.tsa field value matches
// the timestamp's issuer distinguishing name.
type TSAGeneralNameOrderMatchCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// timestampWrapper is the timestamp to verify.
	timestampWrapper *diagnostic.TimestampWrapper
}

// NewTSAGeneralNameOrderMatchCheck is the default constructor. Port of
// TSAGeneralNameOrderMatchCheck(I18nProvider, XmlSAV, TimestampWrapper, LevelRule).
func NewTSAGeneralNameOrderMatchCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	timestampWrapper *diagnostic.TimestampWrapper, constraint policy.LevelRule) *TSAGeneralNameOrderMatchCheck {
	c := &TSAGeneralNameOrderMatchCheck{
		ChainItemBase:    process.NewChainItemBase(i18nProvider, result, constraint),
		timestampWrapper: timestampWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TSAGeneralNameOrderMatchCheck) Process() bool {
	return c.timestampWrapper.IsTSAGeneralNameOrderMatch()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TSAGeneralNameOrderMatchCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_TAV_DTSAOM
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *TSAGeneralNameOrderMatchCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_TAV_DTSAOM_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TSAGeneralNameOrderMatchCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *TSAGeneralNameOrderMatchCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE
}
