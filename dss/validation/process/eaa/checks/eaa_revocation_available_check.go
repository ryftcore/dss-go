// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/EAARevocationAvailableCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAARevocationAvailableCheck checks if the SVA was able to retrieve a status
// token for the EAA.
type EAARevocationAvailableCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewEAARevocationAvailableCheck is the default constructor.
func NewEAARevocationAvailableCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *EAARevocationAvailableCheck {
	c := &EAARevocationAvailableCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *EAARevocationAvailableCheck) Process() bool {
	return utils.IsCollectionNotEmpty(c.eaa.EAARevocations())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *EAARevocationAvailableCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_AV
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *EAARevocationAvailableCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_REV_AV_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *EAARevocationAvailableCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *EAARevocationAvailableCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
