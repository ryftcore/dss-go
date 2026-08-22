// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/vci/checks/SignaturePolicyIdentifiedCheck.java (DSS 6.5.RC1).
package vci

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SignaturePolicyIdentifiedCheck checks if signature policy identifier is present and the policy is
// identified.
type SignaturePolicyIdentifiedCheck struct {
	*process.ChainItemBase[*jaxb.XmlVCI]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewSignaturePolicyIdentifiedCheck is the default constructor. Port of
// SignaturePolicyIdentifiedCheck(I18nProvider, XmlVCI, SignatureWrapper, LevelRule).
func NewSignaturePolicyIdentifiedCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlVCI],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *SignaturePolicyIdentifiedCheck {
	c := &SignaturePolicyIdentifiedCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SignaturePolicyIdentifiedCheck) Process() bool {
	return c.signature.IsPolicyPresent() && c.signature.IsPolicyIdentified()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SignaturePolicyIdentifiedCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_VCI_ISPA
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SignaturePolicyIdentifiedCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_VCI_ISPA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SignaturePolicyIdentifiedCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SignaturePolicyIdentifiedCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSignaturePolicyNotAvailable
}
