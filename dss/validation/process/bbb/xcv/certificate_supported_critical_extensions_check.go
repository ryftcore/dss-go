// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/checks/CertificateSupportedCriticalExtensionsCheck.java (DSS 6.5.RC1).
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

// CertificateSupportedCriticalExtensionsCheck verifies if the certificate
// does not contain any of the certificate extensions listed within a list of
// unsupported certificate extensions.
type CertificateSupportedCriticalExtensionsCheck struct {
	*bbb.AbstractMultiValuesCheckItem[*jaxb.XmlSubXCV]

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper
}

// NewCertificateSupportedCriticalExtensionsCheck is the default constructor.
// Port of CertificateSupportedCriticalExtensionsCheck(I18nProvider, XmlSubXCV, CertificateWrapper, MultiValuesRule).
func NewCertificateSupportedCriticalExtensionsCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSubXCV],
	certificate *diagnostic.CertificateWrapper, constraint policy.MultiValuesRule) *CertificateSupportedCriticalExtensionsCheck {
	c := &CertificateSupportedCriticalExtensionsCheck{
		AbstractMultiValuesCheckItem: bbb.NewAbstractMultiValuesCheckItem(i18nProvider, result, constraint),
		certificate:                  certificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateSupportedCriticalExtensionsCheck) Process() bool {
	return utils.IsCollectionEmpty(c.unsupportedCertificateExtensionsOids())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateSupportedCriticalExtensionsCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVDCCUCE
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateSupportedCriticalExtensionsCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBXCVDCCUCEANS
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage(): the
// arg is pre-rendered as Java's List#toString() would ("[a, b, c]") - see the
// identical note in certificate_forbidden_extensions_check.go.
func (c *CertificateSupportedCriticalExtensionsCheck) BuildErrorMessage() *jaxb.XmlMessage {
	oids := c.unsupportedCertificateExtensionsOids()
	return c.BuildXmlMessage(c.ErrorMessageTag(), "["+strings.Join(oids, ", ")+"]")
}

// unsupportedCertificateExtensionsOids ports the private
// getUnsupportedCertificateExtensionsOids().
func (c *CertificateSupportedCriticalExtensionsCheck) unsupportedCertificateExtensionsOids() []string {
	var values []string
	for _, certificateExtension := range c.certificate.CertificateExtensions() {
		critical := certificateExtension.ExtensionCritical()
		var oid string
		if certificateExtension.ExtensionOID() != nil {
			oid = *certificateExtension.ExtensionOID()
		}
		if critical != nil && *critical && !c.ProcessValueCheck(oid) {
			values = append(values, oid)
		}
	}
	return values
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateSupportedCriticalExtensionsCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *CertificateSupportedCriticalExtensionsCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationCertificateChainGeneralFailure
}
