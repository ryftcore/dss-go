// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/checks/CheckSubXCVResult.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CheckSubXCVResult checks if the SubXCV validation result is valid.
type CheckSubXCVResult struct {
	*process.ChainItemBase[*jaxb.XmlXCV]

	// subResult is the SubXCV result.
	subResult *jaxb.XmlSubXCV
}

// NewCheckSubXCVResult is the default constructor. Port of
// CheckSubXCVResult(Provider, XmlXCV, XmlSubXCV, LevelRule).
func NewCheckSubXCVResult(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlXCV],
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
	return jaxb.XmlBlockTypeSubXCV
}

// Process performs the check. Port of process().
func (c *CheckSubXCVResult) Process() bool {
	return c.IsValid(&c.subResult.XmlConstraintsConclusionContent)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CheckSubXCVResult) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVSub
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CheckSubXCVResult) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVSubANS
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *CheckSubXCVResult) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTagCertificateID, c.subResult.Id)
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
