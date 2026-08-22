// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateIssuerNameCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateIssuerNameCheck verifies if the certificate's issuer
// distinguished name matches the subject distinguished name of the issuer.
type CertificateIssuerNameCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificateIssuerNameCheck is the default constructor. Port of
// CertificateIssuerNameCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewCertificateIssuerNameCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *CertificateIssuerNameCheck {
	c := &CertificateIssuerNameCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateIssuerNameCheck) Process() bool {
	return c.certificate.IsMatchingIssuerSubjectName()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateIssuerNameCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_DCIDNMSDNIC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateIssuerNameCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_DCIDNMSDNIC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateIssuerNameCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateIssuerNameCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}
