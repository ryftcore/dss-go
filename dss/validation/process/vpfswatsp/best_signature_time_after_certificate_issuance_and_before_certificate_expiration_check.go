// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/psv/checks/BestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
package vpfswatsp

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// BestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck
// checks if the best-signature-time is in the certificate's validity range.
type BestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck struct {
	*process.ChainItemBase[*jaxb.XmlPSV]

	// controlTime is the best-signature-time.
	controlTime time.Time

	// certificate is the certificate to check.
	certificate *diagnostic.CertificateWrapper

	// currentTimeSubIndication is the current SubIndication.
	currentTimeSubIndication enumerations.SubIndication
}

// NewBestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck
// is the default constructor. Port of
// BestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck(Provider, XmlPSV, Date, CertificateWrapper, SubIndication, LevelRule).
func NewBestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck(
	i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlPSV], controlTime time.Time,
	certificate *diagnostic.CertificateWrapper, currentTimeSubIndication enumerations.SubIndication,
	constraint policy.LevelRule) *BestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck {
	c := &BestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck{
		ChainItemBase:            process.NewChainItemBase(i18nProvider, result, constraint),
		controlTime:              controlTime,
		certificate:              certificate,
		currentTimeSubIndication: currentTimeSubIndication,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *BestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck) Process() bool {
	// inclusive by RFC 5280
	notBefore := c.certificate.NotBefore()
	notAfter := c.certificate.NotAfter()
	return notBefore != nil && c.controlTime.Compare(*notBefore) >= 0 &&
		notAfter != nil && c.controlTime.Compare(*notAfter) <= 0
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *BestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagTSVISCNVABST
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *BestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagTSVISCNVABSTANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *BestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *BestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return c.currentTimeSubIndication
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo(). The " ? " placeholders are upstream's, spaces included.
func (c *BestSignatureTimeAfterCertificateIssuanceAndBeforeCertificateExpirationCheck) BuildAdditionalInfo() *string {
	notBeforeStr := " ? "
	if c.certificate.NotBefore() != nil {
		notBeforeStr = process.GetFormattedDate(c.certificate.NotBefore())
	}
	notAfterStr := " ? "
	if c.certificate.NotAfter() != nil {
		notAfterStr = process.GetFormattedDate(c.certificate.NotAfter())
	}
	validationTime := process.GetFormattedDate(&c.controlTime)
	message := c.I18nProvider.GetMessage(i18n.MessageTagCertificateValidity, validationTime, notBeforeStr, notAfterStr)
	return &message
}
