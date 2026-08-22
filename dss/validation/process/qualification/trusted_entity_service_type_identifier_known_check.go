// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/usage/checks/TrustedEntityServiceTypeIdentifierKnownCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TrustedEntityServiceTypeIdentifierKnownCheck verifies whether the type
// identifier of the trusted entity service is known.
type TrustedEntityServiceTypeIdentifierKnownCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateApprovalStatus]

	// stiUri is the Service Type Identifier URI.
	stiUri string
}

// NewTrustedEntityServiceTypeIdentifierKnownCheck is the default
// constructor. Port of
// TrustedEntityServiceTypeIdentifierKnownCheck(I18nProvider, XmlValidationCertificateApprovalStatus, String, LevelRule).
func NewTrustedEntityServiceTypeIdentifierKnownCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationCertificateApprovalStatus], stiUri string,
	constraint policy.LevelRule) *TrustedEntityServiceTypeIdentifierKnownCheck {
	c := &TrustedEntityServiceTypeIdentifierKnownCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		stiUri:        stiUri,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TrustedEntityServiceTypeIdentifierKnownCheck) Process() bool {
	status := enumerations.LoTEServiceTypeIdentifierFromURI(c.stiUri)
	return status != nil && status.Label() != "" // Label is present -> defined
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *TrustedEntityServiceTypeIdentifierKnownCheck) BuildAdditionalInfo() *string {
	message := c.I18nProvider.GetMessage(i18n.MessageTag_CERTIFICATE_USAGE_STI, c.stiUri)
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TrustedEntityServiceTypeIdentifierKnownCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_STI_KNOWN
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TrustedEntityServiceTypeIdentifierKnownCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGE_STI_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TrustedEntityServiceTypeIdentifierKnownCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TrustedEntityServiceTypeIdentifierKnownCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
