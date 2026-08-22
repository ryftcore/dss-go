// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificatePS2DQcCompetentAuthorityNameCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
)

// CertificatePS2DQcCompetentAuthorityNameCheck checks the certificate's
// QcPS2D Name.
type CertificatePS2DQcCompetentAuthorityNameCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificatePS2DQcCompetentAuthorityNameCheck is the default constructor.
// Port of CertificatePS2DQcCompetentAuthorityNameCheck(I18nProvider, XmlSubXCV, CertificateWrapper, MultiValuesRule).
func NewCertificatePS2DQcCompetentAuthorityNameCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.MultiValuesRule) *CertificatePS2DQcCompetentAuthorityNameCheck {
	c := &CertificatePS2DQcCompetentAuthorityNameCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificatePS2DQcCompetentAuthorityNameCheck) Process() bool {
	psd2Info := c.certificate.PSD2Info()
	if psd2Info != nil && psd2Info.NcaName() != "" {
		return c.ProcessValueCheck(psd2Info.NcaName())
	}
	// no value present
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificatePS2DQcCompetentAuthorityNameCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CMDCICQCNA
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificatePS2DQcCompetentAuthorityNameCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CMDCICQCNA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificatePS2DQcCompetentAuthorityNameCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificatePS2DQcCompetentAuthorityNameCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}
