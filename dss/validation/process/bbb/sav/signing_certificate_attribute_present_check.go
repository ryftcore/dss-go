// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/SigningCertificateAttributePresentCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// SigningCertificateAttributePresentCheck checks if the signing certificate
// reference is present.
type SigningCertificateAttributePresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// token is the token to verify.
	token diagnostic.TokenProxy
}

// NewSigningCertificateAttributePresentCheck is the default constructor. Port of
// SigningCertificateAttributePresentCheck(I18nProvider, XmlSAV, TokenProxy, LevelRule).
func NewSigningCertificateAttributePresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	token diagnostic.TokenProxy, constraint policy.LevelRule) *SigningCertificateAttributePresentCheck {
	c := &SigningCertificateAttributePresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		token:         token,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *SigningCertificateAttributePresentCheck) Process() bool {
	return c.token.IsSigningCertificateReferencePresent()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *SigningCertificateAttributePresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBICSISASCP
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *SigningCertificateAttributePresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBICSISASCPANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *SigningCertificateAttributePresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *SigningCertificateAttributePresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationSigConstraintsFailure
}
