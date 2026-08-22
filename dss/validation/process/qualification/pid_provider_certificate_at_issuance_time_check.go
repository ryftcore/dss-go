// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/pid/checks/PIDProviderCertificateAtIssuanceTimeCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// PIDProviderCertificateAtIssuanceTimeCheck verifies whether the
// certificate's usage corresponds to a certificate for PID issuance at the
// certificate issuance time.
type PIDProviderCertificateAtIssuanceTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationPIDQualificationProcess]

	// certificateApprovalStatusAtIssuanceTime is the certificate
	// qualification at signing time.
	certificateApprovalStatusAtIssuanceTime enumerations.CertificateApprovalStatus
}

// NewPIDProviderCertificateAtIssuanceTimeCheck is the default constructor.
// Port of
// PIDProviderCertificateAtIssuanceTimeCheck(I18nProvider, XmlValidationPIDQualificationProcess, CertificateApprovalStatus, LevelRule).
func NewPIDProviderCertificateAtIssuanceTimeCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationPIDQualificationProcess],
	certificateApprovalStatusAtIssuanceTime enumerations.CertificateApprovalStatus,
	constraint policy.LevelRule) *PIDProviderCertificateAtIssuanceTimeCheck {
	c := &PIDProviderCertificateAtIssuanceTimeCheck{
		ChainItemBase:                           process.NewChainItemBase(i18nProvider, result, constraint),
		certificateApprovalStatusAtIssuanceTime: certificateApprovalStatusAtIssuanceTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *PIDProviderCertificateAtIssuanceTimeCheck) Process() bool {
	return enumerations.CertificateApprovalStatusEnumPIDProvider == c.certificateApprovalStatusAtIssuanceTime
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *PIDProviderCertificateAtIssuanceTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_PID_PROVIDER_AT_ISSUANCE_TIME
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *PIDProviderCertificateAtIssuanceTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_PID_PROVIDER_AT_ISSUANCE_TIME_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *PIDProviderCertificateAtIssuanceTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *PIDProviderCertificateAtIssuanceTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
