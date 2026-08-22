// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/isc/checks/DigestValueMatchCheck.java (DSS 6.5.RC1).
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

// DigestValueMatchCheck checks if the digest value matches for a signing
// certificate reference.
type DigestValueMatchCheck struct {
	*process.ChainItemBase[*jaxb.XmlISC]

	// token is the token to verify.
	token diagnostic.TokenProxy
}

// NewDigestValueMatchCheck is the default constructor. Port of
// DigestValueMatchCheck(I18nProvider, XmlISC, TokenProxy, LevelRule).
func NewDigestValueMatchCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlISC],
	token diagnostic.TokenProxy, constraint policy.LevelRule) *DigestValueMatchCheck {
	c := &DigestValueMatchCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		token:         token,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *DigestValueMatchCheck) Process() bool {
	signingCertificateReferences := c.token.SigningCertificateReferences()
	if utils.IsCollectionNotEmpty(signingCertificateReferences) {
		for _, reference := range signingCertificateReferences {
			if reference.IsDigestValuePresent() && reference.IsDigestValueMatch() {
				return true
			}
		}
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *DigestValueMatchCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ICDVV
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *DigestValueMatchCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ICDVV_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *DigestValueMatchCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *DigestValueMatchCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationNoSigningCertificateFound
}
