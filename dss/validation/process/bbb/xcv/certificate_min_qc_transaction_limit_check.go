// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateMinQcTransactionLimitCheck.java (DSS 6.5.RC1).
package xcv

import (
	"math"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateMinQcTransactionLimitCheck checks the minimal allowed QC
// transaction limit for the certificate.
type CertificateMinQcTransactionLimitCheck struct {
	*process.ChainItemBase[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// constraint is the constraint from the policy file.
	constraint policy.NumericValueRule
}

// NewCertificateMinQcTransactionLimitCheck is the default constructor. Port
// of CertificateMinQcTransactionLimitCheck(I18nProvider, XmlSubXCV, CertificateWrapper, NumericValueRule).
func NewCertificateMinQcTransactionLimitCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.NumericValueRule) *CertificateMinQcTransactionLimitCheck {
	c := &CertificateMinQcTransactionLimitCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:   certificate,
		constraint:    constraint,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateMinQcTransactionLimitCheck) Process() bool {
	qcLimitValue := c.certificate.QCLimitValue()
	if qcLimitValue != nil {
		/*
		 * EN 319 412-5 (ch. 4.3.2 QCStatement regarding limits on the value of transactions) :
		 *
		 * -- value = amount * 10^exponent
		 */
		value := float64(qcLimitValue.Amount()) * math.Pow(10, float64(qcLimitValue.Exponent()))
		return value >= c.constraint.Value()
	}
	// not present
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateMinQcTransactionLimitCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CMDCICQCLVA
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateMinQcTransactionLimitCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CMDCICQCLVA_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateMinQcTransactionLimitCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateMinQcTransactionLimitCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CHAIN_CONSTRAINTS_FAILURE
}
