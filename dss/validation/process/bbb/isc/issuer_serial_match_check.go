// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/isc/checks/IssuerSerialMatchCheck.java (DSS 6.5.RC1).
package isc

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// IssuerSerialMatchCheck checks if the issuer serial matches for a signing
// certificate reference.
type IssuerSerialMatchCheck struct {
	*process.ChainItemBase[*jaxb.XmlISC]

	// token is the token to verify.
	token diagnostic.TokenProxy
}

// NewIssuerSerialMatchCheck is the default constructor. Port of
// IssuerSerialMatchCheck(I18nProvider, XmlISC, TokenProxy, LevelRule).
func NewIssuerSerialMatchCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlISC],
	token diagnostic.TokenProxy, constraint policy.LevelRule) *IssuerSerialMatchCheck {
	c := &IssuerSerialMatchCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		token:         token,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *IssuerSerialMatchCheck) Process() bool {
	signingCertificateReference := c.token.SigningCertificateReference()
	if signingCertificateReference != nil {
		return signingCertificateReference.IsIssuerSerialMatch()
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *IssuerSerialMatchCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBICSAIDNASNE
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *IssuerSerialMatchCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBICSAIDNASNEANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *IssuerSerialMatchCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *IssuerSerialMatchCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationNoSigningCertificateFound
}
