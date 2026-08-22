// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificatePolicySupportedByQSCDIdsCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificatePolicySupportedByQSCDIdsCheck checks if the certificate has a
// is a supported by QSCD policy identifier.
type CertificatePolicySupportedByQSCDIdsCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificatePolicySupportedByQSCDIdsCheck is the default constructor.
// Port of CertificatePolicySupportedByQSCDIdsCheck(I18nProvider, XmlSubXCV, CertificateWrapper, LevelRule).
func NewCertificatePolicySupportedByQSCDIdsCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.LevelRule) *CertificatePolicySupportedByQSCDIdsCheck {
	c := &CertificatePolicySupportedByQSCDIdsCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
//
// This check only uses the certificate (not the TL): checks in policy id
// extension.
func (c *CertificatePolicySupportedByQSCDIdsCheck) Process() bool {
	return process.IsSupportedByQSCD(c.certificate)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificatePolicySupportedByQSCDIdsCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CMDCIQSCD
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificatePolicySupportedByQSCDIdsCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CMDCIQSCD_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificatePolicySupportedByQSCDIdsCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificatePolicySupportedByQSCDIdsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}
