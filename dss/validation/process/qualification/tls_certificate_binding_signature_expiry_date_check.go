// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/TLSCertificateBindingSignatureExpiryDateCheck.java (DSS 6.5.RC1).
//
// EXTERNAL DEPENDENCY GAP: see tls_certificate_binding_present_in_signature_check.go's
// header for the qwac.GetIdentifiedTLSCertificates assumption this file also
// depends on.
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/qwac"
)

// TLSCertificateBindingSignatureExpiryDateCheck verifies the validity of
// the TLS Certificate Binding signature.
type TLSCertificateBindingSignatureExpiryDateCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationQWACProcess]

	// currentTime is the current validation time.
	currentTime time.Time

	// signature is the TLS Certificate Binding signature.
	signature *diagnostic.SignatureWrapper

	// certificates is the list of certificates derived from the validation
	// process.
	certificates []*diagnostic.CertificateWrapper
}

// NewTLSCertificateBindingSignatureExpiryDateCheck is the default
// constructor. Port of
// TLSCertificateBindingSignatureExpiryDateCheck(I18nProvider, XmlValidationQWACProcess, Date, SignatureWrapper, List, LevelRule).
func NewTLSCertificateBindingSignatureExpiryDateCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationQWACProcess],
	currentTime time.Time, signature *diagnostic.SignatureWrapper, certificates []*diagnostic.CertificateWrapper,
	constraint policy.LevelRule) *TLSCertificateBindingSignatureExpiryDateCheck {
	c := &TLSCertificateBindingSignatureExpiryDateCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		currentTime:   currentTime,
		signature:     signature,
		certificates:  certificates,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
//
// The maximum effective expiry time is whichever is soonest of this field,
// the longest-lived TLS certificate identified in the sigD member payload
// (below), or the notAfter time of the signing certificate.
func (c *TLSCertificateBindingSignatureExpiryDateCheck) Process() bool {
	if c.signature.ExpirationTime() != nil && !c.currentTime.Before(*c.signature.ExpirationTime()) {
		return false
	}
	for _, certificate := range c.getIdentifiedTLSCertificates() {
		if certificate.NotAfter() != nil && !c.currentTime.Before(*certificate.NotAfter()) {
			return false
		}
	}
	if c.signature.SigningCertificate() != nil && c.signature.SigningCertificate().NotAfter() != nil &&
		!c.currentTime.Before(*c.signature.SigningCertificate().NotAfter()) {
		return false
	}
	return true
}

// getIdentifiedTLSCertificates ports the private getIdentifiedTLSCertificates().
func (c *TLSCertificateBindingSignatureExpiryDateCheck) getIdentifiedTLSCertificates() []*diagnostic.CertificateWrapper {
	return qwac.GetIdentifiedTLSCertificates(c.signature, c.certificates)
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *TLSCertificateBindingSignatureExpiryDateCheck) BuildAdditionalInfo() *string {
	if c.signature.ExpirationTime() != nil && !c.currentTime.Before(*c.signature.ExpirationTime()) {
		message := c.I18nProvider.GetMessage(i18n.MessageTagQWACExpiryExp, process.GetFormattedDate(&c.currentTime),
			process.GetFormattedDate(c.signature.ExpirationTime()))
		return &message
	}
	for _, certificate := range c.getIdentifiedTLSCertificates() {
		if certificate.NotAfter() != nil && !c.currentTime.Before(*certificate.NotAfter()) {
			message := c.I18nProvider.GetMessage(i18n.MessageTagQWACExpiryExp, process.GetFormattedDate(&c.currentTime),
				process.GetFormattedDate(certificate.NotAfter()), certificate.Id())
			return &message
		}
	}
	if c.signature.SigningCertificate() != nil && c.signature.SigningCertificate().NotAfter() != nil &&
		!c.currentTime.Before(*c.signature.SigningCertificate().NotAfter()) {
		message := c.I18nProvider.GetMessage(i18n.MessageTagQWACExpiryExp, process.GetFormattedDate(&c.currentTime),
			process.GetFormattedDate(c.signature.SigningCertificate().NotAfter()), c.signature.SigningCertificate().Id())
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLSCertificateBindingSignatureExpiryDateCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagTLSCertBindingSigExpiryDate
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLSCertificateBindingSignatureExpiryDateCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagTLSCertBindingSigExpiryDateANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLSCertificateBindingSignatureExpiryDateCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TLSCertificateBindingSignatureExpiryDateCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
