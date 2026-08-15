// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAARevocationPresentCheck.java (DSS 6.5.RC1).
package eaa

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// EAARevocationPresentCheck verifies whether the EAA revocation claim is
// present.
type EAARevocationPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAARevocationPresentCheck is the default constructor.
func NewEAARevocationPresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAARevocationPresentCheck {
	c := &EAARevocationPresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process(), exposed as public in Java so
// that a caller can consult the concrete result before deciding the next
// chain item (see sav.EAAAcceptanceValidation.statusPresent()).
func (c *EAARevocationPresentCheck) Process() bool {
	return c.eaa.EAAPayload().EAAStatus() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAARevocationPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_PR
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAARevocationPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_PR_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAARevocationPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAARevocationPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_EAA_CONSTRAINTS_FAILURE
}
