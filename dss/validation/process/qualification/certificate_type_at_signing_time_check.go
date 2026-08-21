// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/checks/CertificateTypeAtSigningTimeCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateTypeAtSigningTimeCheck checks if the certificate type has been
// successfully identified at best signing time.
type CertificateTypeAtSigningTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationSignatureQualification]

	// certificateQualification is the Certificate Qualification to be
	// checked.
	certificateQualification enumerations.CertificateQualification
}

// NewCertificateTypeAtSigningTimeCheck is the default constructor. Port of
// CertificateTypeAtSigningTimeCheck(I18nProvider, XmlValidationSignatureQualification, CertificateQualification, LevelRule).
func NewCertificateTypeAtSigningTimeCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationSignatureQualification],
	certificateQualification enumerations.CertificateQualification, constraint policy.LevelRule) *CertificateTypeAtSigningTimeCheck {
	c := &CertificateTypeAtSigningTimeCheck{
		ChainItemBase:            process.NewChainItemBase(i18nProvider, result, constraint),
		certificateQualification: certificateQualification,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *CertificateTypeAtSigningTimeCheck) Process() bool {
	return enumerations.CertificateType_UNKNOWN != c.certificateQualification.Type()
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *CertificateTypeAtSigningTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_CERT_TYPE_AT_ST
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *CertificateTypeAtSigningTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_QUAL_CERT_TYPE_AT_ST_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *CertificateTypeAtSigningTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *CertificateTypeAtSigningTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
