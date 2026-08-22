// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAANotBeforePresentCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAANotBeforePresentCheck verifies whether the EAA contains the 'not before'
// information.
type EAANotBeforePresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAANotBeforePresentCheck is the default constructor.
func NewEAANotBeforePresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAANotBeforePresentCheck {
	c := &EAANotBeforePresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAANotBeforePresentCheck) Process() bool {
	return c.eaa.EAANotBefore() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAANotBeforePresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_NBF_PRESENT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAANotBeforePresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_NBF_PRESENT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAANotBeforePresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAANotBeforePresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
