// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAPseudonymUsageCheck.java (DSS 6.5.RC1).
package eaa

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// EAAPseudonymUsageCheck verifies whether the EAA uses a pseudonym.
type EAAPseudonymUsageCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAAPseudonymUsageCheck is the default constructor.
func NewEAAPseudonymUsageCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAAPseudonymUsageCheck {
	c := &EAAPseudonymUsageCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAAPseudonymUsageCheck) Process() bool {
	return c.eaa.HolderPseudonym() == ""
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *EAAPseudonymUsageCheck) BuildAdditionalInfo() *string {
	if c.eaa.HolderPseudonym() != "" {
		message := c.I18nProvider.GetMessage(i18n.MessageTag_PSEUDO, c.eaa.HolderPseudonym())
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAPseudonymUsageCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_PSEUDO_USED
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAPseudonymUsageCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_PSEUDO_USED_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAPseudonymUsageCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAPseudonymUsageCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_EAA_CONSTRAINTS_FAILURE
}
