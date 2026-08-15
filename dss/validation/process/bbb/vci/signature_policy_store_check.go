// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/vci/checks/SignaturePolicyStoreCheck.java (DSS 6.5.RC1).
package vci

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// SignaturePolicyStoreCheck checks if a SignaturePolicyStore is present.
type SignaturePolicyStoreCheck struct {
	*process.ChainItemBase[*jaxb.XmlVCI]

	// signature is the signature to check.
	signature *diagnostic.SignatureWrapper
}

// NewSignaturePolicyStoreCheck is the default constructor. Port of
// SignaturePolicyStoreCheck(I18nProvider, XmlVCI, SignatureWrapper, LevelRule).
func NewSignaturePolicyStoreCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlVCI],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *SignaturePolicyStoreCheck {
	c := &SignaturePolicyStoreCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SignaturePolicyStoreCheck) Process() bool {
	return c.signature.IsPolicyStorePresent()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SignaturePolicyStoreCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_VCI_ISPSUPP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SignaturePolicyStoreCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_VCI_ISPSUPP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SignaturePolicyStoreCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SignaturePolicyStoreCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIGNATURE_POLICY_NOT_AVAILABLE
}
