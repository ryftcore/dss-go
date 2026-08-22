// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/vts/checks/ControlTimeCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ControlTimeCheck verifies the control time determined during the Validation
// Time Sliding process.
type ControlTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlVTS]

	// controlTime is the control time; nil is Java's null.
	controlTime *time.Time
}

// NewControlTimeCheck is the default constructor. Port of
// ControlTimeCheck(Provider, XmlVTS, Date, LevelRule).
func NewControlTimeCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlVTS],
	controlTime *time.Time, constraint policy.LevelRule) *ControlTimeCheck {
	c := &ControlTimeCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		controlTime:   controlTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ControlTimeCheck) Process() bool {
	return c.controlTime != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ControlTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagPSVICTD
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ControlTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagPSVICTDANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ControlTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *ControlTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationNoPOE
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo(), whose null result leaves the element absent.
func (c *ControlTimeCheck) BuildAdditionalInfo() *string {
	if c.controlTime != nil {
		message := c.I18nProvider.GetMessage(i18n.MessageTagControlTimeAlone,
			process.GetFormattedDate(c.controlTime))
		return &message
	}
	return nil
}
