// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/checks/QualifiedCertificateAtCertificateIssuanceCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// QualifiedCertificateAtCertificateIssuanceCheck checks whether the
// certificate is qualified at certificate issuance time.
type QualifiedCertificateAtCertificateIssuanceCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationSignatureQualification]

	// qualificationAtIssuance is the certificate qualification at issuance
	// time.
	qualificationAtIssuance enumerations.CertificateQualification
}

// NewQualifiedCertificateAtCertificateIssuanceCheck is the default
// constructor. Port of
// QualifiedCertificateAtCertificateIssuanceCheck(I18nProvider, XmlValidationSignatureQualification, CertificateQualification, LevelRule).
func NewQualifiedCertificateAtCertificateIssuanceCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationSignatureQualification],
	qualificationAtIssuance enumerations.CertificateQualification, constraint policy.LevelRule) *QualifiedCertificateAtCertificateIssuanceCheck {
	c := &QualifiedCertificateAtCertificateIssuanceCheck{
		ChainItemBase:           process.NewChainItemBase(i18nProvider, result, constraint),
		qualificationAtIssuance: qualificationAtIssuance,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *QualifiedCertificateAtCertificateIssuanceCheck) Process() bool {
	return c.qualificationAtIssuance.IsQc()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *QualifiedCertificateAtCertificateIssuanceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_QC_AT_CC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *QualifiedCertificateAtCertificateIssuanceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_QC_AT_CC_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *QualifiedCertificateAtCertificateIssuanceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *QualifiedCertificateAtCertificateIssuanceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
