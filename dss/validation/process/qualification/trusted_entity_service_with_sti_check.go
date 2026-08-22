// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/usage/checks/TrustedEntityServiceWithStiCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TrustedEntityServiceWithStiCheck verifies whether the trusted entity
// services with the given STI exist.
type TrustedEntityServiceWithStiCheck struct {
	*process.ChainItemBase[*jaxb.XmlValidationCertificateApprovalStatus]

	// trustedServicesWithSti is the list of TrustedEntityServiceWrappers at
	// control time.
	trustedServicesWithSti []*diagnostic.TrustedEntityServiceWrapper

	// stiUri is the Service Type Identifier URI.
	stiUri string
}

// NewTrustedEntityServiceWithStiCheck is the default constructor. Port of
// TrustedEntityServiceWithStiCheck(I18nProvider, XmlValidationCertificateApprovalStatus, List, String, LevelRule).
func NewTrustedEntityServiceWithStiCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlValidationCertificateApprovalStatus], trustedServicesWithSti []*diagnostic.TrustedEntityServiceWrapper,
	stiUri string, constraint policy.LevelRule) *TrustedEntityServiceWithStiCheck {
	c := &TrustedEntityServiceWithStiCheck{
		ChainItemBase:          process.NewChainItemBase(i18nProvider, result, constraint),
		trustedServicesWithSti: trustedServicesWithSti,
		stiUri:                 stiUri,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *TrustedEntityServiceWithStiCheck) Process() bool {
	return utils.IsCollectionNotEmpty(c.trustedServicesWithSti)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *TrustedEntityServiceWithStiCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagCertUsageSti
}

// BuildConstraintMessage builds a constraint message. Port of
// buildConstraintMessage().
func (c *TrustedEntityServiceWithStiCheck) BuildConstraintMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(c.MessageTag(), c.getStiUserFriendlyLabel())
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *TrustedEntityServiceWithStiCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagCertUsageStiANS
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *TrustedEntityServiceWithStiCheck) BuildErrorMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(c.ErrorMessageTag(), c.getStiUserFriendlyLabel())
}

// getStiUserFriendlyLabel ports the private getStiUserFriendlyLabel().
func (c *TrustedEntityServiceWithStiCheck) getStiUserFriendlyLabel() string {
	sti := enumerations.LoTEServiceTypeIdentifierFromURI(c.stiUri)
	if sti != nil && sti.Label() != "" {
		return sti.Label()
	}
	return c.stiUri
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *TrustedEntityServiceWithStiCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *TrustedEntityServiceWithStiCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
