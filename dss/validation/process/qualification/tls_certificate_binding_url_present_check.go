// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/TLSCertificateBindingUrlPresentCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TLSCertificateBindingUrlPresentCheck checks whether an HTTP 'Link'
// response header (as defined in IETF RFC 8288 [6]), with a rel value of
// tls-certificate-binding, has been found.
type TLSCertificateBindingUrlPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationQWACProcess]

	// tlsCertificateBindingUrl is the TLS Certificate Binding URL.
	tlsCertificateBindingUrl string
}

// NewTLSCertificateBindingUrlPresentCheck is the default constructor. Port of
// TLSCertificateBindingUrlPresentCheck(I18nProvider, XmlValidationQWACProcess, String, LevelRule).
func NewTLSCertificateBindingUrlPresentCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationQWACProcess],
	tlsCertificateBindingUrl string, constraint policy.LevelRule) *TLSCertificateBindingUrlPresentCheck {
	c := &TLSCertificateBindingUrlPresentCheck{
		ChainItemBase:            process.NewChainItemBase(i18nProvider, result, constraint),
		tlsCertificateBindingUrl: tlsCertificateBindingUrl,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TLSCertificateBindingUrlPresentCheck) Process() bool {
	return c.tlsCertificateBindingUrl != ""
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLSCertificateBindingUrlPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagTLSCertBindingURL
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLSCertificateBindingUrlPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagTLSCertBindingURLANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLSCertificateBindingUrlPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TLSCertificateBindingUrlPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
