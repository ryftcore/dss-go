// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/isc/checks/SigningCertificateRecognitionCheck.java (DSS 6.5.RC1).
package isc

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SigningCertificateRecognitionCheck checks if a signing certificate is
// identified.
type SigningCertificateRecognitionCheck struct {
	*process.ChainItemBase[*jaxb.XmlISC]

	// token is the token to verify.
	token diagnostic.TokenProxy
}

// NewSigningCertificateRecognitionCheck is the default constructor. Port of
// SigningCertificateRecognitionCheck(I18nProvider, XmlISC, TokenProxy, LevelRule).
func NewSigningCertificateRecognitionCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlISC],
	token diagnostic.TokenProxy, constraint policy.LevelRule) *SigningCertificateRecognitionCheck {
	c := &SigningCertificateRecognitionCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		token:         token,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SigningCertificateRecognitionCheck) Process() bool {
	signingCertificate := c.token.SigningCertificate()
	return signingCertificate != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SigningCertificateRecognitionCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ISCI
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SigningCertificateRecognitionCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ISCI_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SigningCertificateRecognitionCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SigningCertificateRecognitionCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_NO_SIGNING_CERTIFICATE_FOUND
}
