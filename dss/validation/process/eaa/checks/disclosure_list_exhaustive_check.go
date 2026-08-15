// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/DisclosureListExhaustiveCheck.java (DSS 6.5.RC1).
package checks

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// DisclosureListExhaustiveCheck verifies whether disclosures have been
// provided for all selectively disclosable claim hashes present within the
// EAA's payload.
type DisclosureListExhaustiveCheck struct {
	*process.ChainItemBase[*jaxb.XmlFC]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewDisclosureListExhaustiveCheck is the default constructor.
func NewDisclosureListExhaustiveCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlFC],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *DisclosureListExhaustiveCheck {
	c := &DisclosureListExhaustiveCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *DisclosureListExhaustiveCheck) Process() bool {
	digestMatchers := c.eaa.DigestMatchers()
	if utils.IsCollectionEmpty(digestMatchers) {
		return true
	}
	for _, d := range digestMatchers {
		if !d.DataFound {
			return false
		}
	}
	return true
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *DisclosureListExhaustiveCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_DLEEAAP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *DisclosureListExhaustiveCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_DLEEAAP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *DisclosureListExhaustiveCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *DisclosureListExhaustiveCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
