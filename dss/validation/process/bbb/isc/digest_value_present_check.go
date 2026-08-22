// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/isc/checks/DigestValuePresentCheck.java (DSS 6.5.RC1).
package isc

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// DigestValuePresentCheck checks if the digest value is present for a signing
// certificate reference.
type DigestValuePresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlISC]

	// token is the token to verify.
	token diagnostic.TokenProxy
}

// NewDigestValuePresentCheck is the default constructor. Port of
// DigestValuePresentCheck(I18nProvider, XmlISC, TokenProxy, LevelRule).
func NewDigestValuePresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlISC],
	token diagnostic.TokenProxy, constraint policy.LevelRule) *DigestValuePresentCheck {
	c := &DigestValuePresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		token:         token,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *DigestValuePresentCheck) Process() bool {
	signingCertificateReferences := c.token.SigningCertificateReferences()
	if utils.IsCollectionNotEmpty(signingCertificateReferences) {
		for _, reference := range signingCertificateReferences {
			if reference.IsDigestValuePresent() {
				return true
			}
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *DigestValuePresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ISACDP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *DigestValuePresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ISACDP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *DigestValuePresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *DigestValuePresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationNoSigningCertificateFound
}
