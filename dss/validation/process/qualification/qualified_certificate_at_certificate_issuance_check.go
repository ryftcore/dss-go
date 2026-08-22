// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/checks/QualifiedCertificateAtCertificateIssuanceCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
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
// QualifiedCertificateAtCertificateIssuanceCheck(Provider, XmlValidationSignatureQualification, CertificateQualification, LevelRule).
func NewQualifiedCertificateAtCertificateIssuanceCheck(i18nProvider *i18n.Provider,
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
	return i18n.MessageTagQualQCAtCC
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *QualifiedCertificateAtCertificateIssuanceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQualQCAtCCANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *QualifiedCertificateAtCertificateIssuanceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *QualifiedCertificateAtCertificateIssuanceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
