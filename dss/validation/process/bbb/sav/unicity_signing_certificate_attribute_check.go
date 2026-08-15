// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/checks/UnicitySigningCertificateAttributeCheck.java (DSS 6.5.RC1).
package sav

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// UnicitySigningCertificateAttributeCheck checks if only one reference to the
// signing certificate reference is present.
type UnicitySigningCertificateAttributeCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// token is the token to verify.
	token diagnostic.TokenProxy
}

// NewUnicitySigningCertificateAttributeCheck is the default constructor. Port of
// UnicitySigningCertificateAttributeCheck(I18nProvider, XmlSAV, TokenProxy, LevelRule).
func NewUnicitySigningCertificateAttributeCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	token diagnostic.TokenProxy, constraint policy.LevelRule) *UnicitySigningCertificateAttributeCheck {
	c := &UnicitySigningCertificateAttributeCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		token:         token,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *UnicitySigningCertificateAttributeCheck) Process() bool {
	return c.token.IsSigningCertificateReferenceUnique()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *UnicitySigningCertificateAttributeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ISASCPU
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *UnicitySigningCertificateAttributeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_ICS_ISASCPU_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *UnicitySigningCertificateAttributeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *UnicitySigningCertificateAttributeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_SIG_CONSTRAINTS_FAILURE
}
