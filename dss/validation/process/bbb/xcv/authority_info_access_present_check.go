// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/AuthorityInfoAccessPresentCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AuthorityInfoAccessPresentCheck checks if the authority information access
// urls are present.
type AuthorityInfoAccessPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewAuthorityInfoAccessPresentCheck is the default constructor. Port of
// AuthorityInfoAccessPresentCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewAuthorityInfoAccessPresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *AuthorityInfoAccessPresentCheck {
	c := &AuthorityInfoAccessPresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AuthorityInfoAccessPresentCheck) Process() bool {
	return utils.IsCollectionNotEmpty(c.certificate.CAIssuersAccessUrls())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AuthorityInfoAccessPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVAIAPres
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *AuthorityInfoAccessPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVAIAPresANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AuthorityInfoAccessPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *AuthorityInfoAccessPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}
