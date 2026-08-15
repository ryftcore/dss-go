// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/pid/checks/PIDProviderCertificateAtValidationTimeCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// PIDProviderCertificateAtValidationTimeCheck verifies whether the
// certificate's usage corresponds to a certificate for PID issuance at the
// validation time.
type PIDProviderCertificateAtValidationTimeCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationPIDQualificationProcess]

	// certificateApprovalStatusAtValidationTime is the certificate
	// qualification at signing time.
	certificateApprovalStatusAtValidationTime enumerations.CertificateApprovalStatus
}

// NewPIDProviderCertificateAtValidationTimeCheck is the default
// constructor. Port of
// PIDProviderCertificateAtValidationTimeCheck(I18nProvider, XmlValidationPIDQualificationProcess, CertificateApprovalStatus, LevelRule).
func NewPIDProviderCertificateAtValidationTimeCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationPIDQualificationProcess],
	certificateApprovalStatusAtValidationTime enumerations.CertificateApprovalStatus,
	constraint policy.LevelRule) *PIDProviderCertificateAtValidationTimeCheck {
	c := &PIDProviderCertificateAtValidationTimeCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		certificateApprovalStatusAtValidationTime: certificateApprovalStatusAtValidationTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *PIDProviderCertificateAtValidationTimeCheck) Process() bool {
	return enumerations.CertificateApprovalStatusEnum_PID_PROVIDER == c.certificateApprovalStatusAtValidationTime
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *PIDProviderCertificateAtValidationTimeCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_PID_PROVIDER_AT_VALIDATION_TIME
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *PIDProviderCertificateAtValidationTimeCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_PID_PROVIDER_AT_VALIDATION_TIME_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *PIDProviderCertificateAtValidationTimeCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_FAILED
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *PIDProviderCertificateAtValidationTimeCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
