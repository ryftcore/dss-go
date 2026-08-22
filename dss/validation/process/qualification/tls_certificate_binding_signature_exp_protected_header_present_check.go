// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/TLSCertificateBindingSignatureExpProtectedHeaderPresentCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TLSCertificateBindingSignatureExpProtectedHeaderPresentCheck verifies
// whether the 'exp' protected header (expiration time) is present.
type TLSCertificateBindingSignatureExpProtectedHeaderPresentCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationQWACProcess]

	// signature is the TLS Certificate Binding signature.
	signature *diagnostic.SignatureWrapper
}

// NewTLSCertificateBindingSignatureExpProtectedHeaderPresentCheck is the
// default constructor. Port of
// TLSCertificateBindingSignatureExpProtectedHeaderPresentCheck(I18nProvider, XmlValidationQWACProcess, SignatureWrapper, LevelRule).
func NewTLSCertificateBindingSignatureExpProtectedHeaderPresentCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationQWACProcess], signature *diagnostic.SignatureWrapper,
	constraint policy.LevelRule) *TLSCertificateBindingSignatureExpProtectedHeaderPresentCheck {
	c := &TLSCertificateBindingSignatureExpProtectedHeaderPresentCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TLSCertificateBindingSignatureExpProtectedHeaderPresentCheck) Process() bool {
	return c.signature.ExpirationTime() != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLSCertificateBindingSignatureExpProtectedHeaderPresentCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_TLS_CERT_BINDING_SIG_EXP
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLSCertificateBindingSignatureExpProtectedHeaderPresentCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_TLS_CERT_BINDING_SIG_EXP_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLSCertificateBindingSignatureExpProtectedHeaderPresentCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TLSCertificateBindingSignatureExpProtectedHeaderPresentCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
