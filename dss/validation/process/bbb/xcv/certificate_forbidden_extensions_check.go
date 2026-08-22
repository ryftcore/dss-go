// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateForbiddenExtensionsCheck.java (DSS 6.5.RC1).
package xcv

import (
	"strings"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb"
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

// BuildErrorMessage builds an error message. Port of buildErrorMessage(): the
// arg is pre-rendered as Java's List#toString() would ("[a, b, c]") since the
// shared i18n.messageFormatArgString falls back to fmt.Sprint for a non-string
// arg, which renders a []string as "[a b c]" (no commas) - a frozen-package
// (i18n) formatting gap flagged for a later phase, worked around here at the
// call site rather than by editing that frozen package.
func (c *CertificateForbiddenExtensionsCheck) BuildErrorMessage() *jaxb.XmlMessage {
	oids := c.usedForbiddenCertificateExtensionsOids()
	return c.BuildXmlMessage(c.ErrorMessageTag(), "["+strings.Join(oids, ", ")+"]")
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
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateForbiddenExtensionsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}
