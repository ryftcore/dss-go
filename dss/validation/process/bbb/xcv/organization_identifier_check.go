// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/OrganizationIdentifierCheck.java (DSS 6.5.RC1).
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

// OrganizationIdentifierCheck checks if the certificate's
// organizationIdentifier attribute is acceptable.
type OrganizationIdentifierCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewOrganizationIdentifierCheck is the default constructor. Port of
// OrganizationIdentifierCheck(I18nProvider, XmlSubXCV, CertificateWrapper, MultiValuesRule).
func NewOrganizationIdentifierCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.MultiValuesRule) *OrganizationIdentifierCheck {
	c := &OrganizationIdentifierCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *OrganizationIdentifierCheck) Process() bool {
	return c.ProcessValueCheck(c.certificate.OrganizationIdentifier())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *OrganizationIdentifierCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISCGORGAI
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *OrganizationIdentifierCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_ISCGORGAI_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *OrganizationIdentifierCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *OrganizationIdentifierCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CHAIN_CONSTRAINTS_FAILURE
}
