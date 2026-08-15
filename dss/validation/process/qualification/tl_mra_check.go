// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/checks/TLMRACheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	dssjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// TLMRACheck checks if the Trusted List is defined with MRA.
type TLMRACheck struct {
	*process.ChainItemBase[*jaxb.XmlTLAnalysis]

	// currentTL is the Trusted List to check.
	currentTL *dssjaxb.XmlTrustedList
}

// NewTLMRACheck is the default constructor. Port of
// TLMRACheck(I18nProvider, XmlTLAnalysis, XmlTrustedList, LevelRule).
func NewTLMRACheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlTLAnalysis],
	currentTL *dssjaxb.XmlTrustedList, constraint policy.LevelRule) *TLMRACheck {
	c := &TLMRACheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		currentTL:     currentTL,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TLMRACheck) Process() bool {
	return c.currentTL.Mra == nil || !*c.currentTL.Mra
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLMRACheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_TL_IMRA
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLMRACheck) ErrorMessageTag() i18n.MessageTag {
	parentTL := c.currentTL.Parent
	if parentTL != nil {
		tslType := parentTL.Type
		if tslType != nil {
			if enumerations.TSLTypeEnum_EUlistofthelists.URI() == *tslType {
				return i18n.MessageTag_QUAL_TL_IMRA_ANS_V1
			} else if enumerations.TSLTypeEnum_AdESlistofthelists.URI() == *tslType {
				return i18n.MessageTag_QUAL_TL_IMRA_ANS_V2
			}
		}
	}
	// default
	return i18n.MessageTag_QUAL_TL_IMRA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLMRACheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port of
// getFailedSubIndicationForConclusion().
func (c *TLMRACheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
