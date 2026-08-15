// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/X509UrlPresentCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
)

// X509UrlPresentCheck verifies whether a 'x5u' (X.509 URL) header parameter is
// present within the protected header of a signature.
type X509UrlPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// signature is the signature to verify.
	signature *diagnostic.SignatureWrapper
}

// NewX509UrlPresentCheck is the default constructor. Port of
// X509UrlPresentCheck(I18nProvider, XmlSAV, SignatureWrapper, LevelRule).
func NewX509UrlPresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *X509UrlPresentCheck {
	c := &X509UrlPresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *X509UrlPresentCheck) Process() bool {
	return utils.IsCollectionNotEmpty(c.signature.X509UrlReferences())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *X509UrlPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ISAX509UP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *X509UrlPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ISAX509UP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *X509UrlPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *X509UrlPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE
}
