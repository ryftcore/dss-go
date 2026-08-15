// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/checks/RevocationAfterCertificateIssuanceCheck.java (DSS 6.5.RC1).
package xcv

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// RevocationAfterCertificateIssuanceCheck checks whether the concerned
// certificate has existed at the time of revocation data generation.
type RevocationAfterCertificateIssuanceCheck struct {
	*process.ChainItemBase[*jaxb.XmlRAC]

	// certificate is the certificate in question.
	certificate *diagnostic.CertificateWrapper

	// revocationData is the revocation data to check.
	revocationData *diagnostic.RevocationWrapper
}

// NewRevocationAfterCertificateIssuanceCheck is the default constructor. Port of
// RevocationAfterCertificateIssuanceCheck(I18nProvider, XmlRAC, CertificateWrapper,
// RevocationWrapper, LevelRule).
func NewRevocationAfterCertificateIssuanceCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlRAC], certificate *diagnostic.CertificateWrapper,
	revocationData *diagnostic.RevocationWrapper, constraint policy.LevelRule) *RevocationAfterCertificateIssuanceCheck {
	c := &RevocationAfterCertificateIssuanceCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:    certificate,
		revocationData: revocationData,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationAfterCertificateIssuanceCheck) Process() bool {
	certNotBefore := c.certificate.NotBefore()
	thisUpdate := c.revocationData.ThisUpdate()
	return certNotBefore != nil && thisUpdate != nil && !certNotBefore.After(*thisUpdate)
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *RevocationAfterCertificateIssuanceCheck) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTag_REVOCATION_INFO,
		process.GetFormattedDate(c.revocationData.ThisUpdate()),
		process.GetFormattedDate(c.certificate.NotBefore()),
		process.GetFormattedDate(c.certificate.NotAfter()))
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationAfterCertificateIssuanceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_AFTER_CERT_NOT_BEFORE
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *RevocationAfterCertificateIssuanceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_AFTER_CERT_NOT_BEFORE_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationAfterCertificateIssuanceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *RevocationAfterCertificateIssuanceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE
}
