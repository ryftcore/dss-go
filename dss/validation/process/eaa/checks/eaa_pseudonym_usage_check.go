// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAAPseudonymUsageCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
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
		message := c.I18nProvider.GetMessage(i18n.MessageTagPseudo, c.eaa.HolderPseudonym())
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAAPseudonymUsageCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAPseudoUsed
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAAPseudonymUsageCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAPseudoUsedANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAAPseudonymUsageCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAAPseudonymUsageCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
