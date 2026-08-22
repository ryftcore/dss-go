// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateQcEuLimitValueCurrencyCheck.java (DSS 6.5.RC1).
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

// CertificateQcEuLimitValueCurrencyCheck checks the minimal allowed
// QCLimitValue statement is defined with an acceptable currency.
type CertificateQcEuLimitValueCurrencyCheck struct {
	*bbb.AbstractValueCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificateQcEuLimitValueCurrencyCheck is the default constructor. Port
// of CertificateQcEuLimitValueCurrencyCheck(Provider, XmlSubXCV, CertificateWrapper, ValueRule).
func NewCertificateQcEuLimitValueCurrencyCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.ValueRule) *CertificateQcEuLimitValueCurrencyCheck {
	c := &CertificateQcEuLimitValueCurrencyCheck{
		AbstractValueCheckItem: bbb.NewAbstractValueCheckItem(i18nProvider, result, constraint),
		certificate:            certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateQcEuLimitValueCurrencyCheck) Process() bool {
	qcLimitValue := c.certificate.QCLimitValue()
	if qcLimitValue != nil {
		return c.ProcessValueCheck(qcLimitValue.Currency())
	}
	// not present
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateQcEuLimitValueCurrencyCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCICQCLVHAC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateQcEuLimitValueCurrencyCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCICQCLVHACANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateQcEuLimitValueCurrencyCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateQcEuLimitValueCurrencyCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}
