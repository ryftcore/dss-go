// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/checks/TLNotExpiredCheck.java (DSS 6.5.RC1).
//
// Java declares currentTL as XmlTrustSourceList (the supertype shared by XmlTrustedList and
// XmlTrustSourceList itself); the Go port takes the shared embedded content struct directly,
// since every field this check reads (NextUpdate) lives there. TLValidationBlock passes
// &currentTL.XmlTrustSourceListContent, LoTEValidationBlock &currentList.XmlTrustSourceListContent.
package qualification

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	dssjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// TLNotExpiredCheck verifies whether the Trusted List is not expired.
type TLNotExpiredCheck struct {
	*process.ChainItemBase[*jaxb.XmlTLAnalysis]

	// currentTL is the Trusted List to check.
	currentTL *dssjaxb.XmlTrustSourceListContent

	// currentTime is the validation time.
	currentTime time.Time
}

// NewTLNotExpiredCheck is the default constructor. Port of
// TLNotExpiredCheck(I18nProvider, XmlTLAnalysis, XmlTrustSourceList, Date, LevelRule).
func NewTLNotExpiredCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlTLAnalysis],
	currentTL *dssjaxb.XmlTrustSourceListContent, currentTime time.Time,
	constraint policy.LevelRule) *TLNotExpiredCheck {
	c := &TLNotExpiredCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		currentTL:     currentTL,
		currentTime:   currentTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TLNotExpiredCheck) Process() bool {
	nextUpdate := c.currentTL.NextUpdate
	return nextUpdate != nil && nextUpdate.Time().After(c.currentTime)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLNotExpiredCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_TL_EXP
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLNotExpiredCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_TL_EXP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLNotExpiredCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port of
// getFailedSubIndicationForConclusion().
func (c *TLNotExpiredCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
