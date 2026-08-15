// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAOneTimeUseCheck.java (DSS 6.5.RC1).
package eaa

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// EAAOneTimeUseCheck verifies whether the EAA is for one-time use.
type EAAOneTimeUseCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAAOneTimeUseCheck is the default constructor.
func NewEAAOneTimeUseCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAAOneTimeUseCheck {
	c := &EAAOneTimeUseCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAOneTimeUseCheck) Process() bool {
	return !utils.IsTrue(c.eaa.OneTimeUse())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAOneTimeUseCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_OTU
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAOneTimeUseCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_OTU_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAOneTimeUseCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAOneTimeUseCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_EAA_CONSTRAINTS_FAILURE
}
