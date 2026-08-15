// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/TLSCertificateBindingQWACCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// TLSCertificateBindingQWACCheck checks whether the certificate presented
// in the binding is QWAC.
type TLSCertificateBindingQWACCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationQWACProcess]

	// qwacProfile is the QWAC Profile of the binding certificate.
	qwacProfile enumerations.QWACProfile
}

// NewTLSCertificateBindingQWACCheck is the default constructor. Port of
// TLSCertificateBindingQWACCheck(I18nProvider, XmlValidationQWACProcess, QWACProfile, LevelRule).
func NewTLSCertificateBindingQWACCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationQWACProcess],
	qwacProfile enumerations.QWACProfile, constraint policy.LevelRule) *TLSCertificateBindingQWACCheck {
	c := &TLSCertificateBindingQWACCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		qwacProfile:   qwacProfile,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TLSCertificateBindingQWACCheck) Process() bool {
	return enumerations.QWACProfile_QWAC_2 == c.qwacProfile
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLSCertificateBindingQWACCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_TLS_CERT_BINDING_QWAC2
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLSCertificateBindingQWACCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_TLS_CERT_BINDING_QWAC2_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLSCertificateBindingQWACCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TLSCertificateBindingQWACCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
