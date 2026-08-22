// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/checks/QSCDCertificateAtSigningTimeCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// QSCDCertificateAtSigningTimeCheck checks whether the certificate has been
// for QSCD at signing time.
type QSCDCertificateAtSigningTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationSignatureQualification]

	// certificateQualification is the certificate qualification at signing
	// time.
	certificateQualification enumerations.CertificateQualification
}

// NewQSCDCertificateAtSigningTimeCheck is the default constructor. Port of
// QSCDCertificateAtSigningTimeCheck(I18nProvider, XmlValidationSignatureQualification, CertificateQualification, LevelRule).
func NewQSCDCertificateAtSigningTimeCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationSignatureQualification],
	certificateQualification enumerations.CertificateQualification, constraint policy.LevelRule) *QSCDCertificateAtSigningTimeCheck {
	c := &QSCDCertificateAtSigningTimeCheck{
		ChainItemBase:            process.NewChainItemBase(i18nProvider, result, constraint),
		certificateQualification: certificateQualification,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *QSCDCertificateAtSigningTimeCheck) Process() bool {
	return c.certificateQualification.IsQscd()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *QSCDCertificateAtSigningTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagQualQSCDAtST
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *QSCDCertificateAtSigningTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagQualQSCDAtSTANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *QSCDCertificateAtSigningTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *QSCDCertificateAtSigningTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
