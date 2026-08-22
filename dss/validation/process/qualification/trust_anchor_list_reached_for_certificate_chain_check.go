// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/checks/TrustAnchorListReachedForCertificateChainCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TrustAnchorListReachedForCertificateChainCheck checks if one of the
// trust anchor lists has been reached for the certificate chain.
type TrustAnchorListReachedForCertificateChainCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationEAAQualification]

	// signingCertificate is the end-entity certificate.
	signingCertificate *diagnostic.CertificateWrapper
}

// NewTrustAnchorListReachedForCertificateChainCheck is the default
// constructor. Port of
// TrustAnchorListReachedForCertificateChainCheck(I18nProvider, XmlValidationEAAQualification, CertificateWrapper, LevelRule).
func NewTrustAnchorListReachedForCertificateChainCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationEAAQualification], signingCertificate *diagnostic.CertificateWrapper,
	constraint policy.LevelRule) *TrustAnchorListReachedForCertificateChainCheck {
	c := &TrustAnchorListReachedForCertificateChainCheck{
		ChainItemBase:      process.NewChainItemBase(i18nProvider, result, constraint),
		signingCertificate: signingCertificate,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TrustAnchorListReachedForCertificateChainCheck) Process() bool {
	return c.signingCertificate != nil &&
		(c.signingCertificate.IsTrustedListReached() || c.signingCertificate.IsListOfTrustedEntitiesReached())
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TrustAnchorListReachedForCertificateChainCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_CERT_TRUST_ANCHOR_LIST_REACHED
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TrustAnchorListReachedForCertificateChainCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_EAA_CERT_TRUST_ANCHOR_LIST_REACHED_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TrustAnchorListReachedForCertificateChainCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TrustAnchorListReachedForCertificateChainCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
