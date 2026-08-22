// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/checks/TLFreshnessCheck.java (DSS 6.5.RC1).
//
// See tl_not_expired_check.go for why currentTL is typed as the shared
// XmlTrustSourceListContent rather than XmlTrustSourceList.
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	dssjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TLFreshnessCheck verifies whether the Trusted List is fresh.
type TLFreshnessCheck struct {
	*process.ChainItemBase[*jaxb.XmlTLAnalysis]

	// currentTL is the Trusted List to check.
	currentTL *dssjaxb.XmlTrustSourceListContent

	// currentTime is the validation time.
	currentTime time.Time

	// durationRule is the constraint defining the maximum freshness time.
	durationRule policy.DurationRule
}

// NewTLFreshnessCheck is the default constructor. Port of
// TLFreshnessCheck(I18nProvider, XmlTLAnalysis, XmlTrustSourceList, Date, DurationRule).
func NewTLFreshnessCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlTLAnalysis],
	currentTL *dssjaxb.XmlTrustSourceListContent, currentTime time.Time,
	durationRule policy.DurationRule) *TLFreshnessCheck {
	c := &TLFreshnessCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, durationRule),
		currentTL:     currentTL,
		currentTime:   currentTime,
		durationRule:  durationRule,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TLFreshnessCheck) Process() bool {
	maxFreshness := c.getMaxFreshness()
	validationDateTime := c.currentTime.UnixMilli()
	limit := validationDateTime - maxFreshness

	lastLoading := c.currentTL.LastLoading
	return lastLoading != nil && lastLoading.Time().UnixMilli() > limit
}

// getMaxFreshness ports the private getMaxFreshness().
func (c *TLFreshnessCheck) getMaxFreshness() int64 {
	return c.durationRule.Duration()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLFreshnessCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_TL_FRESH
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLFreshnessCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_TL_FRESH_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLFreshnessCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port of
// getFailedSubIndicationForConclusion().
func (c *TLFreshnessCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
