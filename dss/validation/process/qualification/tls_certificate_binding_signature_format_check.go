// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/TLSCertificateBindingSignatureFormatCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TLSCertificateBindingSignatureFormatCheck verifies whether the format of
// the obtained TLS Certificate Binding signature is allowed.
type TLSCertificateBindingSignatureFormatCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationQWACProcess]

	// signature is the TLS Certificate Binding signature.
	signature *diagnostic.SignatureWrapper
}

// NewTLSCertificateBindingSignatureFormatCheck is the default constructor.
// Port of
// TLSCertificateBindingSignatureFormatCheck(I18nProvider, XmlValidationQWACProcess, SignatureWrapper, LevelRule).
func NewTLSCertificateBindingSignatureFormatCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationQWACProcess],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *TLSCertificateBindingSignatureFormatCheck {
	c := &TLSCertificateBindingSignatureFormatCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TLSCertificateBindingSignatureFormatCheck) Process() bool {
	if c.signature.SignatureFormat() == "" {
		return false
	}
	form, err := c.signature.SignatureFormat().SignatureForm()
	return err == nil && enumerations.SignatureFormJAdES == form
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLSCertificateBindingSignatureFormatCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_TLS_CERT_BINDING_SIG_FORM
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLSCertificateBindingSignatureFormatCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_TLS_CERT_BINDING_SIG_FORM_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLSCertificateBindingSignatureFormatCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TLSCertificateBindingSignatureFormatCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
