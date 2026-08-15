// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAIssuanceDatePresentCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// EAAIssuanceDatePresentCheck verifies whether the EAA Presentation contains
// the issuance date.
type EAAIssuanceDatePresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAAIssuanceDatePresentCheck is the default constructor.
func NewEAAIssuanceDatePresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAAIssuanceDatePresentCheck {
	c := &EAAIssuanceDatePresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAIssuanceDatePresentCheck) Process() bool {
	return c.eaa.EAAIssuedAt() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAIssuanceDatePresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_ISSUANCE_DATE_PRESENT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAIssuanceDatePresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_ISSUANCE_DATE_PRESENT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAIssuanceDatePresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAIssuanceDatePresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_EAA_CONSTRAINTS_FAILURE
}
