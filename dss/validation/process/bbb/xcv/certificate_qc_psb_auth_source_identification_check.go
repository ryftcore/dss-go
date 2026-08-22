// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateQcPSBAuthSourceIdentificationCheck.java (DSS 6.5.RC1).
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

// CertificateQcPSBAuthSourceIdentificationCheck checks the
// AuthSourceIdentification attribute value of the QcPSB QcStatement.
type CertificateQcPSBAuthSourceIdentificationCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificateQcPSBAuthSourceIdentificationCheck is the default
// constructor. Port of CertificateQcPSBAuthSourceIdentificationCheck(I18nProvider, XmlSubXCV, CertificateWrapper, MultiValuesRule).
func NewCertificateQcPSBAuthSourceIdentificationCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.MultiValuesRule) *CertificateQcPSBAuthSourceIdentificationCheck {
	c := &CertificateQcPSBAuthSourceIdentificationCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateQcPSBAuthSourceIdentificationCheck) Process() bool {
	if c.certificate.QcPSB() != nil {
		return c.ProcessValueCheck(c.certificate.QcPSB().AuthSourceIdentification())
	}
	return false
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateQcPSBAuthSourceIdentificationCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCPSBASIA
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateQcPSBAuthSourceIdentificationCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVCMDCPSBASIAANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateQcPSBAuthSourceIdentificationCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateQcPSBAuthSourceIdentificationCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationChainConstraintsFailure
}
