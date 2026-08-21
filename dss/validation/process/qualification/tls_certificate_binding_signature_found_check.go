// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/TLSCertificateBindingSignatureFoundCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TLSCertificateBindingSignatureFoundCheck checks whether a TLS Certificate
// Binding signature has been successfully fetched.
type TLSCertificateBindingSignatureFoundCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationQWACProcess]

	// signature is the TLS Certificate Binding signature.
	signature *diagnostic.SignatureWrapper
}

// NewTLSCertificateBindingSignatureFoundCheck is the default constructor.
// Port of
// TLSCertificateBindingSignatureFoundCheck(I18nProvider, XmlValidationQWACProcess, SignatureWrapper, LevelRule).
func NewTLSCertificateBindingSignatureFoundCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationQWACProcess],
	signature *diagnostic.SignatureWrapper, constraint policy.LevelRule) *TLSCertificateBindingSignatureFoundCheck {
	c := &TLSCertificateBindingSignatureFoundCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		signature:     signature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TLSCertificateBindingSignatureFoundCheck) Process() bool {
	return c.signature != nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLSCertificateBindingSignatureFoundCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_TLS_CERT_BINDING_SIG
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLSCertificateBindingSignatureFoundCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_TLS_CERT_BINDING_SIG_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLSCertificateBindingSignatureFoundCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TLSCertificateBindingSignatureFoundCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
