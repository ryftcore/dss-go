// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAAdministrativeExpirationDatePresentCheck.java (DSS 6.5.RC1).
package eaa

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// EAAAdministrativeExpirationDatePresentCheck verifies whether the EAA
// contains an administrative expiration date.
type EAAAdministrativeExpirationDatePresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAAAdministrativeExpirationDatePresentCheck is the default constructor.
func NewEAAAdministrativeExpirationDatePresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAAAdministrativeExpirationDatePresentCheck {
	c := &EAAAdministrativeExpirationDatePresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAAdministrativeExpirationDatePresentCheck) Process() bool {
	return c.eaa.AdministrativeExpirationDate() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAAdministrativeExpirationDatePresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_AED_PRESENT
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAAdministrativeExpirationDatePresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_AED_PRESENT_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAAdministrativeExpirationDatePresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAAdministrativeExpirationDatePresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_EAA_CONSTRAINTS_FAILURE
}
