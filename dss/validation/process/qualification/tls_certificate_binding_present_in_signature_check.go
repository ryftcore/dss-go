// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/checks/TLSCertificateBindingPresentInSignatureCheck.java (DSS 6.5.RC1).
//
// EXTERNAL DEPENDENCY GAP: eu.europa.esig.dss.validation.qwac.QWACUtils lives
// outside both this porter's manifest and the shared qualification package
// (it is a top-level dss-validation class, Java package
// eu.europa.esig.dss.validation.qwac, not
// eu.europa.esig.dss.validation.process.qualification.*) and had not been
// ported to Go on disk while this file was written. The call below assumes
// a Go package qwac (github.com/ryftcore/dss-go/dss/validation/qwac) exposing
// GetIdentifiedTLSCertificates(*diagnostic.SignatureWrapper, []*diagnostic.CertificateWrapper)
// []*diagnostic.CertificateWrapper, matching this port's static-utility-class
// flattening convention; must be reconciled once that package exists.
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/qwac"
)

// TLSCertificateBindingPresentInSignatureCheck verifies that the TLS
// certificate is present within the TLS Certificate Binding signature.
type TLSCertificateBindingPresentInSignatureCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationQWACProcess]

	// tlsCertificate is the TLS certificate.
	tlsCertificate *diagnostic.CertificateWrapper

	// bindingSignature is the TLS Certificate Binding signature.
	bindingSignature *diagnostic.SignatureWrapper
}

// NewTLSCertificateBindingPresentInSignatureCheck is the default
// constructor. Port of
// TLSCertificateBindingPresentInSignatureCheck(I18nProvider, XmlValidationQWACProcess, CertificateWrapper, SignatureWrapper, LevelRule).
func NewTLSCertificateBindingPresentInSignatureCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlValidationQWACProcess],
	tlsCertificate *diagnostic.CertificateWrapper, bindingSignature *diagnostic.SignatureWrapper,
	constraint policy.LevelRule) *TLSCertificateBindingPresentInSignatureCheck {
	c := &TLSCertificateBindingPresentInSignatureCheck{
		ChainItemBase:    process.NewChainItemBase(i18nProvider, result, constraint),
		tlsCertificate:   tlsCertificate,
		bindingSignature: bindingSignature,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TLSCertificateBindingPresentInSignatureCheck) Process() bool {
	return c.isTLSCertificateIdentified()
}

// isTLSCertificateIdentified ports the private isTLSCertificateIdentified().
func (c *TLSCertificateBindingPresentInSignatureCheck) isTLSCertificateIdentified() bool {
	return utils.IsCollectionNotEmpty(
		qwac.GetIdentifiedTLSCertificates(c.bindingSignature, []*diagnostic.CertificateWrapper{c.tlsCertificate}))
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TLSCertificateBindingPresentInSignatureCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagTLSCertBindingCertIdentified
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TLSCertificateBindingPresentInSignatureCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagTLSCertBindingCertIdentifiedANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TLSCertificateBindingPresentInSignatureCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TLSCertificateBindingPresentInSignatureCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
