// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificatePolicyIdsCheck.java (DSS 6.5.RC1).
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

// CertificatePolicyIdsCheck checks if the certificate policy ids are
// acceptable.
type CertificatePolicyIdsCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificatePolicyIdsCheck is the default constructor. Port of
// CertificatePolicyIdsCheck(I18nProvider, XmlSubXCV, CertificateWrapper, MultiValuesRule).
func NewCertificatePolicyIdsCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.MultiValuesRule) *CertificatePolicyIdsCheck {
	c := &CertificatePolicyIdsCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificatePolicyIdsCheck) Process() bool {
	return c.ProcessValuesCheck(c.certificate.PolicyIds())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificatePolicyIdsCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CMDCIPI
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificatePolicyIdsCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_CMDCIPI_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificatePolicyIdsCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificatePolicyIdsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CHAIN_CONSTRAINTS_FAILURE
}
