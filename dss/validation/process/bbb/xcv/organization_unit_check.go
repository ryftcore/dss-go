// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/OrganizationUnitCheck.java (DSS 6.5.RC1).
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

// OrganizationUnitCheck checks if the certificate's organization unit is
// acceptable.
type OrganizationUnitCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewOrganizationUnitCheck is the default constructor. Port of
// OrganizationUnitCheck(I18nProvider, XmlSubXCV, CertificateWrapper, MultiValuesRule).
func NewOrganizationUnitCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.MultiValuesRule) *OrganizationUnitCheck {
	c := &OrganizationUnitCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *OrganizationUnitCheck) Process() bool {
	return c.ProcessValueCheck(c.certificate.OrganizationalUnit())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *OrganizationUnitCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVISCGORGAU
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *OrganizationUnitCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVISCGORGAUANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *OrganizationUnitCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *OrganizationUnitCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}
