// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/checks/QualifiedCertificateAtSigningTimeCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// QualifiedCertificateAtSigningTimeCheck checks whether the certificate is
// qualified at signing time.
type QualifiedCertificateAtSigningTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationSignatureQualification]

	// qualificationAtSigningTime is the certificate qualification at signing
	// time.
	qualificationAtSigningTime enumerations.CertificateQualification
}

// NewQualifiedCertificateAtSigningTimeCheck is the default constructor. Port
// of
// QualifiedCertificateAtSigningTimeCheck(I18nProvider, XmlValidationSignatureQualification, CertificateQualification, LevelRule).
func NewQualifiedCertificateAtSigningTimeCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationSignatureQualification],
	qualificationAtSigningTime enumerations.CertificateQualification, constraint policy.LevelRule) *QualifiedCertificateAtSigningTimeCheck {
	c := &QualifiedCertificateAtSigningTimeCheck{
		ChainItemBase:              process.NewChainItemBase(i18nProvider, result, constraint),
		qualificationAtSigningTime: qualificationAtSigningTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *QualifiedCertificateAtSigningTimeCheck) Process() bool {
	return c.qualificationAtSigningTime.IsQc()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *QualifiedCertificateAtSigningTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualQCAtST
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *QualifiedCertificateAtSigningTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQualQCAtSTANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *QualifiedCertificateAtSigningTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *QualifiedCertificateAtSigningTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
