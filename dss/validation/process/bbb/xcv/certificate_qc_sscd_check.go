// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateQcSSCDCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateQcSSCDCheck checks if the certificate is supported by QCSD
// (has the id-etsi-qcs-QcSSCD statement).
type CertificateQcSSCDCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificateQcSSCDCheck is the default constructor. Port of
// CertificateQcSSCDCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewCertificateQcSSCDCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *CertificateQcSSCDCheck {
	c := &CertificateQcSSCDCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateQcSSCDCheck) Process() bool {
	return c.certificate.IsSupportedByQSCD()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateQcSSCDCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CMDCICSQCSSCD
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateQcSSCDCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CMDCICSQCSSCD_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateQcSSCDCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateQcSSCDCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CHAIN_CONSTRAINTS_FAILURE
}
