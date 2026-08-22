// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/AuthorityKeyIdentifierPresentCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AuthorityKeyIdentifierPresentCheck verifies whether the value of the RFC
// 5280 "4.2.1.1. Authority Key Identifier" certificate extension is present.
type AuthorityKeyIdentifierPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewAuthorityKeyIdentifierPresentCheck is the default constructor. Port of
// AuthorityKeyIdentifierPresentCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewAuthorityKeyIdentifierPresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *AuthorityKeyIdentifierPresentCheck {
	c := &AuthorityKeyIdentifierPresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AuthorityKeyIdentifierPresentCheck) Process() bool {
	return c.certificate.AuthorityKeyIdentifier() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AuthorityKeyIdentifierPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_IAKIP
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *AuthorityKeyIdentifierPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_IAKIP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AuthorityKeyIdentifierPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *AuthorityKeyIdentifierPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}
