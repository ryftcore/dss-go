// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/checks/TLWellSignedCheck.java (DSS 6.5.RC1).
//
// See tl_not_expired_check.go for why currentTL is typed as the shared
// XmlTrustSourceListContent rather than XmlTrustSourceList.
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	dssjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TLWellSignedCheck checks whether the signature of the Trusted List is valid.
type TLWellSignedCheck struct {
	*process.ChainItemBase[*jaxb.XmlTLAnalysis]

	// currentTL is the Trusted List to check.
	currentTL *dssjaxb.XmlTrustSourceListContent
}

// NewTLWellSignedCheck is the default constructor. Port of
// TLWellSignedCheck(I18nProvider, XmlTLAnalysis, XmlTrustSourceList, LevelRule).
func NewTLWellSignedCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlTLAnalysis],
	currentTL *dssjaxb.XmlTrustSourceListContent, constraint policy.LevelRule) *TLWellSignedCheck {
	c := &TLWellSignedCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		currentTL:     currentTL,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TLWellSignedCheck) Process() bool {
	return c.currentTL.WellSigned
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLWellSignedCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualTLWS
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLWellSignedCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQualTLWSANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLWellSignedCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port of
// getFailedSubIndicationForConclusion().
func (c *TLWellSignedCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
