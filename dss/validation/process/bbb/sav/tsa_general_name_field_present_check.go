// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/TSAGeneralNameFieldPresentCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TSAGeneralNameFieldPresentCheck checks if the TSTInfo.tsa field is present.
type TSAGeneralNameFieldPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// timestampWrapper is the timestamp to verify.
	timestampWrapper *diagnostic.TimestampWrapper
}

// NewTSAGeneralNameFieldPresentCheck is the default constructor. Port of
// TSAGeneralNameFieldPresentCheck(Provider, XmlSAV, TimestampWrapper, LevelRule).
func NewTSAGeneralNameFieldPresentCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSAV],
	timestampWrapper *diagnostic.TimestampWrapper, constraint policy.LevelRule) *TSAGeneralNameFieldPresentCheck {
	c := &TSAGeneralNameFieldPresentCheck{
		ChainItemBase:    process.NewChainItemBase(i18nProvider, result, constraint),
		timestampWrapper: timestampWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TSAGeneralNameFieldPresentCheck) Process() bool {
	return c.timestampWrapper.IsTSAGeneralNamePresent()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TSAGeneralNameFieldPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBTavITSAP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *TSAGeneralNameFieldPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBTavITSAPANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TSAGeneralNameFieldPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *TSAGeneralNameFieldPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}
