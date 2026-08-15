// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/checks/CheckSubXCVResult.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// CheckSubXCVResult checks if the SubXCV validation result is valid.
type CheckSubXCVResult struct {
	*process.ChainItemBase[*jaxb.XmlXCV]

	// subResult is the SubXCV result.
	subResult *jaxb.XmlSubXCV
}

// NewCheckSubXCVResult is the default constructor. Port of
// CheckSubXCVResult(I18nProvider, XmlXCV, XmlSubXCV, LevelRule).
func NewCheckSubXCVResult(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlXCV],
	subResult *jaxb.XmlSubXCV, constraint policy.LevelRule) *CheckSubXCVResult {
	c := &CheckSubXCVResult{
		ChainItemBase: process.NewChainItemBaseWithId(i18nProvider, result, constraint, subResult.Id),
		subResult:     subResult,
	}
	c.InitChainItem(c)
	return c
}

// BlockType returns the validating block type. Port of getBlockType().
func (c *CheckSubXCVResult) BlockType() jaxb.XmlBlockType {
	return jaxb.XmlBlockType_SUB_XCV
}

// Process performs the check. Port of process().
func (c *CheckSubXCVResult) Process() bool {
	return c.IsValid(&c.subResult.XmlConstraintsConclusionContent)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CheckSubXCVResult) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_SUB
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CheckSubXCVResult) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_SUB_ANS
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *CheckSubXCVResult) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTag_CERTIFICATE_ID, c.subResult.Id)
	return &message
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CheckSubXCVResult) FailedIndicationForConclusion() enumerations.Indication {
	return c.subResult.Conclusion.Indication.Indication()
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CheckSubXCVResult) FailedSubIndicationForConclusion() enumerations.SubIndication {
	if c.subResult.Conclusion.SubIndication == nil {
		return ""
	}
	return c.subResult.Conclusion.SubIndication.SubIndication()
}
