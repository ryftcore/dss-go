// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/DisclosurePresentCheck.java (DSS 6.5.RC1).
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

// DisclosurePresentCheck verifies whether the signature contains at least one
// provided disclosure.
type DisclosurePresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlFC]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper
}

// NewDisclosurePresentCheck is the default constructor.
func NewDisclosurePresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlFC],
	eaaWrapper *diagnostic.EAAWrapper, constraint policy.LevelRule) *DisclosurePresentCheck {
	c := &DisclosurePresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:           eaaWrapper,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *DisclosurePresentCheck) Process() bool {
	digestMatchers := c.eaa.DigestMatchers()
	if utils.IsCollectionEmpty(digestMatchers) {
		return false
	}
	for _, d := range digestMatchers {
		if d.Type != nil && enumerations.DigestMatcherType_EAA_DISCLOSURE == d.Type.DigestMatcherType() {
			return true
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *DisclosurePresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_DPEAAP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *DisclosurePresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_DPEAAP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *DisclosurePresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *DisclosurePresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_FORMAT_FAILURE
}
