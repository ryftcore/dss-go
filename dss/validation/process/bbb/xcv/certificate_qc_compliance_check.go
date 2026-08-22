// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateQcComplianceCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateQcComplianceCheck checks if the certificate is QC Compliant
// (has the id-etsi-qcs-QcCompliance statement).
type CertificateQcComplianceCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificateQcComplianceCheck is the default constructor. Port of
// CertificateQcComplianceCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewCertificateQcComplianceCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *CertificateQcComplianceCheck {
	c := &CertificateQcComplianceCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateQcComplianceCheck) Process() bool {
	return c.certificate.IsQcCompliance()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateQcComplianceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCICQCC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateQcComplianceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCICQCCANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateQcComplianceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateQcComplianceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}
