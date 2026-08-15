// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificatePS2DQcCompetentAuthorityIdCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb"
)

// CertificatePS2DQcCompetentAuthorityIdCheck checks the certificate's QcPS2D
// Id.
type CertificatePS2DQcCompetentAuthorityIdCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificatePS2DQcCompetentAuthorityIdCheck is the default constructor.
// Port of CertificatePS2DQcCompetentAuthorityIdCheck(I18nProvider, XmlSubXCV, CertificateWrapper, MultiValuesRule).
func NewCertificatePS2DQcCompetentAuthorityIdCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.MultiValuesRule) *CertificatePS2DQcCompetentAuthorityIdCheck {
	c := &CertificatePS2DQcCompetentAuthorityIdCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificatePS2DQcCompetentAuthorityIdCheck) Process() bool {
	psd2Info := c.certificate.PSD2Info()
	if psd2Info != nil && psd2Info.NcaId() != "" {
		return c.ProcessValueCheck(psd2Info.NcaId())
	}
	// no value present
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificatePS2DQcCompetentAuthorityIdCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CMDCICQCIA
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificatePS2DQcCompetentAuthorityIdCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CMDCICQCIA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificatePS2DQcCompetentAuthorityIdCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificatePS2DQcCompetentAuthorityIdCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CHAIN_CONSTRAINTS_FAILURE
}
