// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/TLSCertificateBindingSignatureValidationResultCheck.java (DSS 6.5.RC1).
//
// CROSS-CHUNK ASSUMPTION: SignatureValidationResultCheck (Java package
// qualification.signature.checks) is owned by a sibling porter of this
// shared package and was not present on disk while this file was written;
// its constructor is assumed to follow this port's usual generic-[T any]
// ChainItem shape (NewSignatureValidationResultCheck[T any](I18nProvider,
// *process.Result[T], *jaxb.XmlConclusion, policy.LevelRule)
// *SignatureValidationResultCheck[T]), matching the Java constructor's
// parameter order, but must be reconciled against the sibling porter's
// actual signature once available.
package qualification

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
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
	return i18n.MessageTag_TLS_CERT_BINDING_SIG_VALID
}

// ErrorMessageTag returns the check's error message tag. Port of the
// overridden getErrorMessageTag().
func (c *TLSCertificateBindingSignatureValidationResultCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_TLS_CERT_BINDING_SIG_VALID_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// the overridden getFailedIndicationForConclusion().
func (c *TLSCertificateBindingSignatureValidationResultCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of the overridden getFailedSubIndicationForConclusion(), whose
// default is null.
func (c *TLSCertificateBindingSignatureValidationResultCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
