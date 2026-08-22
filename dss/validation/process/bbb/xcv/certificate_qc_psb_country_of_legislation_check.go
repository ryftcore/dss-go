// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateQcPSBCountryOfLegislationCheck.java (DSS 6.5.RC1).
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

// CertificateQcPSBCountryOfLegislationCheck checks the
// CountryOfLegislation attribute value of the QcPSB QcStatement.
type CertificateQcPSBCountryOfLegislationCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificateQcPSBCountryOfLegislationCheck is the default constructor.
// Port of CertificateQcPSBCountryOfLegislationCheck(I18nProvider, XmlSubXCV, CertificateWrapper, MultiValuesRule).
func NewCertificateQcPSBCountryOfLegislationCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.MultiValuesRule) *CertificateQcPSBCountryOfLegislationCheck {
	c := &CertificateQcPSBCountryOfLegislationCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateQcPSBCountryOfLegislationCheck) Process() bool {
	if c.certificate.QcPSB() != nil {
		return c.ProcessValueCheck(c.certificate.QcPSB().CountryOfLegislation())
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateQcPSBCountryOfLegislationCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCPSBCLA
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateQcPSBCountryOfLegislationCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCPSBCLAANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateQcPSBCountryOfLegislationCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateQcPSBCountryOfLegislationCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}
