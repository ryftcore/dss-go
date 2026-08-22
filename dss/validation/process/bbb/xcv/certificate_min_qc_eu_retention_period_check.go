// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateMinQcEuRetentionPeriodCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateMinQcEuRetentionPeriodCheck checks the QCEuRetentionPeriod
// constraint.
type CertificateMinQcEuRetentionPeriodCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// constraint is the constraint from the policy file.
	constraint policy.NumericValueRule
}

// NewCertificateMinQcEuRetentionPeriodCheck is the default constructor. Port
// of CertificateMinQcEuRetentionPeriodCheck(I18nProvider, XmlSubXCV, CertificateWrapper, NumericValueRule).
func NewCertificateMinQcEuRetentionPeriodCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.NumericValueRule) *CertificateMinQcEuRetentionPeriodCheck {
	c := &CertificateMinQcEuRetentionPeriodCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
		constraint:    constraint,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateMinQcEuRetentionPeriodCheck) Process() bool {
	qcEuRetentionPeriod := c.certificate.QCEuRetentionPeriod()
	if qcEuRetentionPeriod != nil {
		return float64(*qcEuRetentionPeriod) >= c.constraint.Value()
	}
	// not present
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateMinQcEuRetentionPeriodCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCICQCERPA
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateMinQcEuRetentionPeriodCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCICQCERPAANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateMinQcEuRetentionPeriodCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateMinQcEuRetentionPeriodCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}
