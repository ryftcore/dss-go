// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/TSAGeneralNameValueMatchCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TSAGeneralNameValueMatchCheck checks if the TSTInfo.tsa field value matches
// the timestamp's issuer distinguishing name.
type TSAGeneralNameValueMatchCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// timestampWrapper is the timestamp to verify.
	timestampWrapper *diagnostic.TimestampWrapper
}

// NewTSAGeneralNameValueMatchCheck is the default constructor. Port of
// TSAGeneralNameValueMatchCheck(I18nProvider, XmlSAV, TimestampWrapper, LevelRule).
func NewTSAGeneralNameValueMatchCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	timestampWrapper *diagnostic.TimestampWrapper, constraint policy.LevelRule) *TSAGeneralNameValueMatchCheck {
	c := &TSAGeneralNameValueMatchCheck{
		ChainItemBase:    process.NewChainItemBase(i18nProvider, result, constraint),
		timestampWrapper: timestampWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TSAGeneralNameValueMatchCheck) Process() bool {
	return c.timestampWrapper.IsTSAGeneralNameMatch()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TSAGeneralNameValueMatchCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBTavDTSAVM
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *TSAGeneralNameValueMatchCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBTavDTSAVMANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TSAGeneralNameValueMatchCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *TSAGeneralNameValueMatchCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}
