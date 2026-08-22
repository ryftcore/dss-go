// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/RevocationInfoAccessPresentCheck.java (DSS 6.5.RC1).
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

// RevocationInfoAccessPresentCheck checks if the revocation access points are
// present in the certificate.
type RevocationInfoAccessPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewRevocationInfoAccessPresentCheck is the default constructor. Port of
// RevocationInfoAccessPresentCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewRevocationInfoAccessPresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *RevocationInfoAccessPresentCheck {
	c := &RevocationInfoAccessPresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationInfoAccessPresentCheck) Process() bool {
	return utils.IsCollectionNotEmpty(c.certificate.CRLDistributionPoints()) || utils.IsCollectionNotEmpty(c.certificate.OCSPAccessUrls())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationInfoAccessPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_PRES
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *RevocationInfoAccessPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_PRES_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationInfoAccessPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *RevocationInfoAccessPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}
