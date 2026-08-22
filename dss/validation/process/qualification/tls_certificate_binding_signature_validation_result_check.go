// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/TLSCertificateBindingSignatureValidationResultCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TLSCertificateBindingSignatureValidationResultCheck verifies the basic
// validation result of a JAdES signature on the TLS Certificate Binding.
type TLSCertificateBindingSignatureValidationResultCheck struct {
	*SignatureValidationResultCheck[*jaxb.XmlValidationQWACProcess]
}

// NewTLSCertificateBindingSignatureValidationResultCheck is the default
// constructor. Port of
// TLSCertificateBindingSignatureValidationResultCheck(I18nProvider, XmlValidationQWACProcess, XmlConclusion, LevelRule).
func NewTLSCertificateBindingSignatureValidationResultCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationQWACProcess], bindingSignatureBasicValidationConclusion *jaxb.XmlConclusion,
	constraint policy.LevelRule) *TLSCertificateBindingSignatureValidationResultCheck {
	c := &TLSCertificateBindingSignatureValidationResultCheck{
		SignatureValidationResultCheck: NewSignatureValidationResultCheck(i18nProvider, result, bindingSignatureBasicValidationConclusion, constraint),
	}
	c.InitChainItem(c)
	return c
}

// MessageTag returns the check's message tag. Port of the overridden getMessageTag().
func (c *TLSCertificateBindingSignatureValidationResultCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagTLSCertBindingSigValid
}

// ErrorMessageTag returns the check's error message tag. Port of the
// overridden getErrorMessageTag().
func (c *TLSCertificateBindingSignatureValidationResultCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagTLSCertBindingSigValidANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// the overridden getFailedIndicationForConclusion().
func (c *TLSCertificateBindingSignatureValidationResultCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of the overridden getFailedSubIndicationForConclusion(), whose
// default is null.
func (c *TLSCertificateBindingSignatureValidationResultCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
