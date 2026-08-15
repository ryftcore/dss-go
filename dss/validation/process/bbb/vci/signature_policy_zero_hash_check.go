// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/vci/checks/SignaturePolicyZeroHashCheck.java (DSS 6.5.RC1).
package vci

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SignaturePolicyZeroHashCheck checks if the signature policy identifier is a zero-hash.
type SignaturePolicyZeroHashCheck struct {
	*process.ChainItemBase[*jaxb.XmlVCI]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewSignaturePolicyZeroHashCheck is the default constructor. Port of
// SignaturePolicyZeroHashCheck(I18nProvider, XmlVCI, SignatureWrapper, LevelRule).
func NewSignaturePolicyZeroHashCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlVCI],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *SignaturePolicyZeroHashCheck {
	c := &SignaturePolicyZeroHashCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SignaturePolicyZeroHashCheck) Process() bool {
	return c.signature.IsPolicyZeroHash()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SignaturePolicyZeroHashCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_VCI_IZHSP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SignaturePolicyZeroHashCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_VCI_IZHSP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SignaturePolicyZeroHashCheck) FailedIndicationForConclusion() enumerations.Indication {
	return "" // Java returns null
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SignaturePolicyZeroHashCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return "" // Java returns null
}
