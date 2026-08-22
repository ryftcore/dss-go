// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/usage/checks/TrustedEntityServiceStatusKnownCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TrustedEntityServiceStatusKnownCheck verifies whether the status of the
// trusted entity service is known.
type TrustedEntityServiceStatusKnownCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateApprovalStatus]

	// serviceStatusUri is the Service Status URI.
	serviceStatusUri string
}

// NewTrustedEntityServiceStatusKnownCheck is the default constructor. Port of
// TrustedEntityServiceStatusKnownCheck(I18nProvider, XmlValidationCertificateApprovalStatus, String, LevelRule).
func NewTrustedEntityServiceStatusKnownCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationCertificateApprovalStatus], serviceStatusUri string,
	constraint policy.LevelRule) *TrustedEntityServiceStatusKnownCheck {
	c := &TrustedEntityServiceStatusKnownCheck{
		ChainItemBase:    process.NewChainItemBase(i18nProvider, result, constraint),
		serviceStatusUri: serviceStatusUri,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TrustedEntityServiceStatusKnownCheck) Process() bool {
	status := enumerations.LoTEServiceStatusFromURI(c.serviceStatusUri)
	return status != nil && status.Label() != "" // Label is present -> defined
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *TrustedEntityServiceStatusKnownCheck) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTag_CERTIFICATE_USAGE_STATUS, c.serviceStatusUri)
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TrustedEntityServiceStatusKnownCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_STATUS_KNOWN
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TrustedEntityServiceStatusKnownCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_STATUS_KNOWN_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TrustedEntityServiceStatusKnownCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TrustedEntityServiceStatusKnownCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
