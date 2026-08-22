// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/vci/checks/SignaturePolicyHashValidCheck.java (DSS 6.5.RC1).
package vci

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SignaturePolicyHashValidCheck checks if signature policy identifier is present and the hash
// matched.
type SignaturePolicyHashValidCheck struct {
	*process.ChainItemBase[*jaxb.XmlVCI]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewSignaturePolicyHashValidCheck is the default constructor. Port of
// SignaturePolicyHashValidCheck(I18nProvider, XmlVCI, SignatureWrapper, LevelRule).
func NewSignaturePolicyHashValidCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlVCI],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *SignaturePolicyHashValidCheck {
	c := &SignaturePolicyHashValidCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SignaturePolicyHashValidCheck) Process() bool {
	return c.signature.IsPolicyPresent() && c.signature.IsPolicyDigestValid()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SignaturePolicyHashValidCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_VCI_ISPM
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SignaturePolicyHashValidCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_VCI_ISPM_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SignaturePolicyHashValidCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SignaturePolicyHashValidCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationPolicyProcessingError
}
