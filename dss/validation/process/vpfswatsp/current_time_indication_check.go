// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/psv/checks/CurrentTimeIndicationCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CurrentTimeIndicationCheck checks if the current state is PASSED.
type CurrentTimeIndicationCheck struct {
	*process.ChainItemBase[*jaxb.XmlPSV]

	// indication is the Indication to check.
	indication enumerations.Indication

	// subIndication is the corresponding SubIndication.
	subIndication enumerations.SubIndication

	// errors are the current errors.
	errors []*jaxb.XmlMessage
}

// NewCurrentTimeIndicationCheck is the default constructor. Port of
// CurrentTimeIndicationCheck(I18nProvider, XmlPSV, Indication, SubIndication, List, LevelRule).
func NewCurrentTimeIndicationCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlPSV],
	indication enumerations.Indication, subIndication enumerations.SubIndication, errors []*jaxb.XmlMessage,
	constraint policy.LevelRule) *CurrentTimeIndicationCheck {
	c := &CurrentTimeIndicationCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		indication:    indication,
		subIndication: subIndication,
		errors:        errors,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CurrentTimeIndicationCheck) Process() bool {
	return enumerations.IndicationPassed == c.indication
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CurrentTimeIndicationCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagPSVIPCVC
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *CurrentTimeIndicationCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagPSVIPCVCANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CurrentTimeIndicationCheck) FailedIndicationForConclusion() enumerations.Indication {
	return c.indication
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CurrentTimeIndicationCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return c.subIndication
}

// PreviousErrors returns a list of previous errors occurred in the chain. Port
// of getPreviousErrors().
func (c *CurrentTimeIndicationCheck) PreviousErrors() []*jaxb.XmlMessage {
	return c.errors
}
