// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAAdministrativeIssuanceDatePresentCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// EAAAdministrativeIssuanceDatePresentCheck verifies whether the EAA contains
// an administrative issuance date.
type EAAAdministrativeIssuanceDatePresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAAAdministrativeIssuanceDatePresentCheck is the default constructor.
func NewEAAAdministrativeIssuanceDatePresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAAAdministrativeIssuanceDatePresentCheck {
	c := &EAAAdministrativeIssuanceDatePresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAAdministrativeIssuanceDatePresentCheck) Process() bool {
	return c.eaa.AdministrativeIssuanceDate() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAAdministrativeIssuanceDatePresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_AID_PRESENT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAAdministrativeIssuanceDatePresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_AID_PRESENT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAAdministrativeIssuanceDatePresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAAdministrativeIssuanceDatePresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_EAA_CONSTRAINTS_FAILURE
}
