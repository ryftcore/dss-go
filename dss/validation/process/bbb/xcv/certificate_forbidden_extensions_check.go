// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateForbiddenExtensionsCheck.java (DSS 6.5.RC1).
package xcv

import (
	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb"
)

// CertificateForbiddenExtensionsCheck verifies if the certificate does not
// contain forbidden certificate extensions.
type CertificateForbiddenExtensionsCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificateForbiddenExtensionsCheck is the default constructor. Port of
// CertificateForbiddenExtensionsCheck(I18nProvider, XmlSubXCV, CertificateWrapper, MultiValuesRule).
func NewCertificateForbiddenExtensionsCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.MultiValuesRule) *CertificateForbiddenExtensionsCheck {
	c := &CertificateForbiddenExtensionsCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateForbiddenExtensionsCheck) Process() bool {
	return utils.IsCollectionEmpty(c.usedForbiddenCertificateExtensionsOids())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateForbiddenExtensionsCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_DCCFCE
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateForbiddenExtensionsCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_DCCFCE_ANS
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *CertificateForbiddenExtensionsCheck) BuildErrorMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(c.ErrorMessageTag(), c.usedForbiddenCertificateExtensionsOids())
}

// usedForbiddenCertificateExtensionsOids ports the private
// getUsedForbiddenCertificateExtensionsOids().
func (c *CertificateForbiddenExtensionsCheck) usedForbiddenCertificateExtensionsOids() []string {
	var values []string
	for _, certificateExtension := range c.certificate.CertificateExtensions() {
		var oid string
		if certificateExtension.ExtensionOID() != nil {
			oid = *certificateExtension.ExtensionOID()
		}
		if c.ProcessValueCheck(oid) {
			values = append(values, oid)
		}
	}
	return values
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateForbiddenExtensionsCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateForbiddenExtensionsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE
}
